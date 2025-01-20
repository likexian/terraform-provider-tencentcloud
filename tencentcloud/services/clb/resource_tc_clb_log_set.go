package clb

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudClbLogSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbLogSetCreate,
		Read:   resourceTencentCloudClbLogSetRead,
		Delete: resourceTencentCloudClbLogSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"logset_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Logset name, which must be unique among all CLS logsets; default value: clb_logset.",
			},
			"logset_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Logset type. Valid values: ACCESS (access logs; default value) and HEALTH (health check logs).",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Logset retention period in days. Maximun value is `90`.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "It has been deprecated from version 1.81.162+. Please use `logset_name` instead.",
				Description: "Logset name, which unique and fixed `clb_logset` among all CLS logsets.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Logset creation time.",
			},
			"topic_count": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Number of log topics in logset.",
			},
		},
	}
}

func resourceTencentCloudClbLogSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_log_set.create")()
	defer clbActionMu.Unlock()
	clbActionMu.Lock()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	_, _, err := service.DescribeClbLogSet(ctx)
	if err != nil {
		return err
	}

	request := clb.NewCreateClsLogSetRequest()
	response := clb.NewCreateClsLogSetResponse()
	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_type"); ok {
		request.LogsetType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntUint64(v.(int))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().CreateClsLogSet(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("Create cls logset failed. Response is nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return tccommon.RetryError(errors.WithStack(retryErr))
			}
		}

		response = result
		return nil
	})

	if err != nil {
		return err
	}

	if response.Response.LogsetId == nil {
		return fmt.Errorf("LogsetId is nil.")
	}

	d.SetId(*response.Response.LogsetId)

	return resourceTencentCloudClbLogSetRead(d, meta)
}

func resourceTencentCloudClbLogSetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_log_set.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svccls.NewClsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id      = d.Id()
	)

	info, err := service.DescribeClsLogset(ctx, id)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return fmt.Errorf("resource `Logset` %s does not exist", id)
	}

	if info.LogsetName != nil {
		_ = d.Set("logset_name", info.LogsetName)
		_ = d.Set("name", info.LogsetName)
	}

	if info.CreateTime != nil {
		_ = d.Set("create_time", info.CreateTime)
	}

	if info.TopicCount != nil {
		_ = d.Set("topic_count", helper.Int64ToStrPoint(*info.TopicCount))
	}

	return nil
}

func resourceTencentCloudClbLogSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_log_set.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svccls.NewClsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		id      = d.Id()
	)

	if err := service.DeleteClsLogsetById(ctx, id); err != nil {
		return err
	}

	return nil
}
