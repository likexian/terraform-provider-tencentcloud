package cynosdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudCynosdbReadonlyInstance() *schema.Resource {
	instanceInfo := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Cluster ID which the readonly instance belongs to.",
		},
		"instance_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of instance.",
		},
		"force_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicate whether to delete readonly instance directly or not. Default is false. If set true, instance will be deleted instead of staying recycle bin. Note: works for both `PREPAID` and `POSTPAID_BY_HOUR` cluster.",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "ID of the VPC.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "ID of the subnet within this VPC.",
		},
	}
	basic := TencentCynosdbInstanceBaseInfo()
	delete(basic, "instance_id")
	delete(basic, "instance_name")
	for k, v := range basic {
		instanceInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudCynosdbReadonlyInstanceCreate,
		Read:   resourceTencentCloudCynosdbReadonlyInstanceRead,
		Update: resourceTencentCloudCynosdbReadonlyInstanceUpdate,
		Delete: resourceTencentCloudCynosdbReadonlyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: instanceInfo,
	}
}

func resourceTencentCloudCynosdbReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_readonly_instance.create")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		client         = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		cynosdbService = CynosdbService{client: client}

		request = cynosdb.NewAddInstancesRequest()
	)

	// instance info
	request.ClusterId = helper.String(d.Get("cluster_id").(string))
	request.InstanceName = helper.String(d.Get("instance_name").(string))
	request.Cpu = helper.IntInt64(d.Get("instance_cpu_core").(int))
	request.Memory = helper.IntInt64(d.Get("instance_memory_size").(int))
	request.ReadOnlyCount = helper.Int64(1)

	// vpc
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	var response *cynosdb.AddInstancesResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().AddInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if response != nil && response.Response != nil && len(response.Response.DealNames) < 1 {
		return fmt.Errorf("cynosdb cluster id count isn't 1")
	}

	dealName := response.Response.DealNames[0]
	dealReq := cynosdb.NewDescribeResourcesByDealNameRequest()
	dealRes := cynosdb.NewDescribeResourcesByDealNameResponse()
	dealReq.DealName = dealName
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		dealRes, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeResourcesByDealName(dealReq)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "InvalidParameterValue.DealNameNotFound" {
					return resource.RetryableError(fmt.Errorf("DealName[%s] Not Found, retry... reason: %s", *dealName, err.Error()))
				}
			}
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if dealRes != nil && dealRes.Response != nil && len(dealRes.Response.BillingResourceInfos) != 1 && len(dealRes.Response.BillingResourceInfos[0].InstanceIds) != 1 {
		return fmt.Errorf("cynosdb readonly instance id count isn't 1")
	}

	id := *dealRes.Response.BillingResourceInfos[0].InstanceIds[0]
	d.SetId(id)

	// set maintenance info
	var weekdays []interface{}
	if v, ok := d.GetOk("instance_maintain_weekdays"); ok {
		weekdays = v.(*schema.Set).List()
	} else {
		weekdays = []interface{}{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	}
	reqWeekdays := make([]*string, 0, len(weekdays))
	for _, v := range weekdays {
		reqWeekdays = append(reqWeekdays, helper.String(v.(string)))
	}
	startTime := int64(d.Get("instance_maintain_start_time").(int))
	duration := int64(d.Get("instance_maintain_duration").(int))
	err = cynosdbService.ModifyMaintainPeriodConfig(ctx, id, startTime, duration, reqWeekdays)
	if err != nil {
		return err
	}

	return resourceTencentCloudCynosdbReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_readonly_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cynosdbService := CynosdbService{client: client}
	clusterId, instance, has, err := cynosdbService.DescribeInstanceById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("instance_cpu_core", instance.Cpu)
	_ = d.Set("instance_memory_size", instance.Memory)
	_ = d.Set("instance_name", instance.InstanceName)
	_ = d.Set("instance_status", instance.Status)
	_ = d.Set("instance_storage_size", instance.Storage)
	if instance.VpcId != nil {
		_ = d.Set("vpc_id", instance.VpcId)
	}
	if instance.SubnetId != nil {
		_ = d.Set("subnet_id", instance.SubnetId)
	}

	maintain, err := cynosdbService.DescribeMaintainPeriod(ctx, id)
	if err != nil {
		return err
	}
	_ = d.Set("instance_maintain_weekdays", maintain.Response.MaintainWeekDays)
	_ = d.Set("instance_maintain_start_time", maintain.Response.MaintainStartTime)
	_ = d.Set("instance_maintain_duration", maintain.Response.MaintainDuration)

	return nil
}

func resourceTencentCloudCynosdbReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_readonly_instance.update")()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		instanceId     = d.Id()
		client         = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		cynosdbService = CynosdbService{client: client}
	)

	d.Partial(true)

	if d.HasChange("instance_cpu_core") || d.HasChange("instance_memory_size") {
		cpu := int64(d.Get("instance_cpu_core").(int))
		memory := int64(d.Get("instance_memory_size").(int))
		err := cynosdbService.UpgradeInstance(ctx, instanceId, cpu, memory)
		if err != nil {
			return err
		}

		errUpdate := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, infos, has, e := cynosdbService.DescribeInstanceById(ctx, instanceId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if !has {
				return resource.NonRetryableError(fmt.Errorf("[CRITAL]%s updating cynosdb cluster instance failed, instance doesn't exist", logId))
			}

			cpuReal := *infos.Cpu
			memReal := *infos.Memory
			if cpu != cpuReal || memory != memReal {
				return resource.RetryableError(fmt.Errorf("[CRITAL] updating cynosdb instance, current cpu and memory values: %d, %d, waiting for them becoming new value: %d, %d", cpuReal, memReal, cpu, memory))
			}
			return nil
		})
		if errUpdate != nil {
			return errUpdate
		}

	}

	if d.HasChange("instance_maintain_weekdays") || d.HasChange("instance_maintain_start_time") || d.HasChange("instance_maintain_duration") {
		weekdays := d.Get("instance_maintain_weekdays").(*schema.Set).List()
		reqWeekdays := make([]*string, 0, len(weekdays))
		for _, v := range weekdays {
			reqWeekdays = append(reqWeekdays, helper.String(v.(string)))
		}
		startTime := int64(d.Get("instance_maintain_start_time").(int))
		duration := int64(d.Get("instance_maintain_duration").(int))
		err := cynosdbService.ModifyMaintainPeriodConfig(ctx, instanceId, startTime, duration, reqWeekdays)
		if err != nil {
			return err
		}

	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		return fmt.Errorf("`vpc_id`, `subnet_id` do not support change now.")
	}

	d.Partial(false)

	return resourceTencentCloudCynosdbReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudCynosdbReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_readonly_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	clusterId := d.Get("cluster_id").(string)
	cynosdbService := CynosdbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	forceDelete := d.Get("force_delete").(bool)

	var err error
	if err = cynosdbService.IsolateInstance(ctx, clusterId, instanceId); err != nil {
		return err
	}

	if forceDelete {
		errUpdate := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, _, has, e := cynosdbService.DescribeInstanceById(ctx, instanceId)
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if has {
				return resource.RetryableError(fmt.Errorf("[CRITAL]%s actual example during removal, heavy new essay", logId))
			}

			return nil
		})
		if errUpdate != nil {
			return errUpdate
		}
		if err = cynosdbService.OfflineInstance(ctx, clusterId, instanceId); err != nil {
			return err
		}
	}

	return nil
}
