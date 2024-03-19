package bh

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDasbResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDasbResourceCreate,
		Read:   resourceTencentCloudDasbResourceRead,
		Update: resourceTencentCloudDasbResourceUpdate,
		Delete: resourceTencentCloudDasbResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"deploy_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy region.",
			},
			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy resource vpcId.",
			},
			"subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy resource subnetId.",
			},
			"resource_edition": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(RESOURCE_EDITION),
				Description:  "Resource type.Value:standard/pro.",
			},
			"resource_node": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of resource nodes.",
			},
			"time_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Billing cycle, only support m: month. This field is mandatory, fill in m.",
			},
			"time_span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Billing time. This field is mandatory, with a minimum value of 1.",
			},
			"auto_renew_flag": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Automatic renewal. 1 is auto renew flag, 0 is not.",
			},
			"deploy_zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Deploy zone.",
			},
			"cidr_block": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subnet segments that require service activation.",
			},
			"vpc_cidr_block": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The network segment corresponding to the VPC that requires service activation.",
			},
			"package_bandwidth": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of bandwidth expansion packets (4M), The set value is an integer multiple of 4.",
			},
			//"package_node": {
			//	Optional:    true,
			//	Computed:    true,
			//	Type:        schema.TypeInt,
			//	Description: "Number of authorized point extension packages (50 points). Cannot exceed 100.",
			//},
		},
	}
}

func resourceTencentCloudDasbResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_resource.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		request         = dasb.NewCreateResourceRequest()
		response        = dasb.NewCreateResourceResponse()
		deployRequest   = dasb.NewDeployResourceRequest()
		describeRequest = dasb.NewDescribeResourcesRequest()
		modifyRequest   = dasb.NewModifyResourceRequest()
		resourceId      string
		vpcId           string
		subnetId        string
		deployRegion    string
		deployZone      string
		cidrBlock       string
		vpcCidrBlock    string
	)

	if v, ok := d.GetOk("deploy_region"); ok {
		request.DeployRegion = helper.String(v.(string))
		deployRegion = v.(string)
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
		vpcId = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
		subnetId = v.(string)
	}

	if v, ok := d.GetOk("resource_edition"); ok {
		request.ResourceEdition = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("resource_node"); ok {
		request.ResourceNode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	request.PayMode = helper.IntInt64(1)

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("deploy_zone"); ok {
		request.DeployZone = helper.String(v.(string))
		deployZone = v.(string)
	}

	if v, ok := d.GetOk("cidr_block"); ok {
		cidrBlock = v.(string)
	}

	if v, ok := d.GetOk("vpc_cidr_block"); ok {
		vpcCidrBlock = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().CreateResource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || *result.Response.ResourceId == "" {
			e = fmt.Errorf("dasb Resource not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb Resource failed, reason:%+v", logId, err)
		return err
	}

	resourceId = *response.Response.ResourceId
	d.SetId(resourceId)

	// deploy resource
	deployRequest.ResourceId = helper.String(resourceId)
	deployRequest.ApCode = helper.String(deployRegion)
	deployRequest.Zone = helper.String(deployZone)
	deployRequest.VpcId = helper.String(vpcId)
	deployRequest.SubnetId = helper.String(subnetId)
	deployRequest.CidrBlock = helper.String(cidrBlock)
	deployRequest.VpcCidrBlock = helper.String(vpcCidrBlock)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().DeployResource(deployRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deployRequest.GetAction(), deployRequest.ToJsonString(), deployRequest.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("dasb Resource deploy error")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s deploy dasb Resource failed, reason:%+v", logId, err)
		return err
	}

	// wait
	describeRequest.ResourceIds = helper.Strings([]string{resourceId})
	err = resource.Retry(tccommon.WriteRetryTimeout*6, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().DescribeResources(describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil || len(result.Response.ResourceSet) != 1 {
			e = fmt.Errorf("dasb Resource not exists")
			return resource.NonRetryableError(e)
		}

		if *result.Response.ResourceSet[0].Status == 4 {
			e = fmt.Errorf("dasb Resource deploy error")
			return resource.NonRetryableError(e)
		}

		if *result.Response.ResourceSet[0].Status == 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("dasb Resource is still in running, state %d", *result.Response.ResourceSet[0].Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dasb Resource failed, reason:%+v", logId, err)
		return err
	}

	// modify
	if v, ok := d.GetOkExists("package_bandwidth"); ok {
		modifyRequest.PackageBandwidth = helper.IntInt64(v.(int))
	}

	//if v, ok := d.GetOkExists("package_node"); ok {
	//	modifyRequest.PackageNode = helper.IntInt64(v.(int))
	//}

	if modifyRequest.PackageBandwidth != nil {
		modifyRequest.ResourceId = &resourceId
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().ModifyResource(modifyRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyRequest.GetAction(), modifyRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update dasb Resource failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDasbResourceRead(d, meta)
}

func resourceTencentCloudDasbResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = DasbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resourceId = d.Id()
	)

	Resource, err := service.DescribeDasbResourceById(ctx, resourceId)
	if err != nil {
		return err
	}

	if Resource == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DasbResource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if Resource.ApCode != nil {
		_ = d.Set("deploy_region", Resource.ApCode)
	}

	if Resource.VpcId != nil {
		_ = d.Set("vpc_id", Resource.VpcId)
	}

	if Resource.SubnetId != nil {
		_ = d.Set("subnet_id", Resource.SubnetId)
	}

	if Resource.SvArgs != nil {
		svArgs := strings.Split(*Resource.SvArgs, "_")
		var tmpStr string
		for _, item := range svArgs {
			if item == RESOURCE_EDITION_PRO || item == RESOURCE_EDITION_STANDARD {
				tmpStr = item
				break
			}
		}

		_ = d.Set("resource_edition", tmpStr)
	}

	if Resource.Nodes != nil {
		_ = d.Set("resource_node", Resource.Nodes)
	}

	if Resource.RenewFlag != nil {
		_ = d.Set("auto_renew_flag", Resource.RenewFlag)
	}

	if Resource.Zone != nil {
		_ = d.Set("deploy_zone", Resource.Zone)
	}

	if Resource.CidrBlock != nil {
		_ = d.Set("cidr_block", Resource.CidrBlock)
	}

	if Resource.VpcCidrBlock != nil {
		_ = d.Set("vpc_cidr_block", Resource.VpcCidrBlock)
	}

	if Resource.PackageBandwidth != nil {
		_ = d.Set("package_bandwidth", Resource.PackageBandwidth)
	}

	//if Resource.PackageNode != nil {
	//	_ = d.Set("package_node", Resource.PackageNode)
	//}

	return nil
}

func resourceTencentCloudDasbResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_resource.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = dasb.NewModifyResourceRequest()
		resourceId = d.Id()
	)

	immutableArgs := []string{"deploy_region", "vpc_id", "subnet_id", "time_unit", "time_span", "pay_mode", "deploy_zone", "cidr_block", "vpc_cidr_block"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.ResourceId = &resourceId
	if d.HasChange("resource_edition") {
		if v, ok := d.GetOk("resource_edition"); ok {
			request.ResourceEdition = helper.String(v.(string))
		}
	}

	if d.HasChange("resource_node") {
		if v, ok := d.GetOkExists("resource_node"); ok {
			request.ResourceNode = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("auto_renew_flag") {
		if v, ok := d.GetOkExists("auto_renew_flag"); ok {
			request.AutoRenewFlag = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("package_bandwidth") {
		if v, ok := d.GetOkExists("package_bandwidth"); ok {
			request.PackageBandwidth = helper.IntInt64(v.(int))
		}
	}

	//if d.HasChange("package_node") {
	//	if v, ok := d.GetOkExists("package_node"); ok {
	//		request.PackageNode = helper.IntInt64(v.(int))
	//	}
	//}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDasbClient().ModifyResource(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update dasb Resource failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDasbResourceRead(d, meta)
}

func resourceTencentCloudDasbResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dasb_resource.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return fmt.Errorf("tencentcloud dasb resource not supported delete, please contact the work order for processing")
}
