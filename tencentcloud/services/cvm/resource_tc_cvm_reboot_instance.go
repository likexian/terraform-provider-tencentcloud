// Code generated by iacg; DO NOT EDIT.
package cvm

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmRebootInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmRebootInstanceCreate,
		Read:   resourceTencentCloudCvmRebootInstanceRead,
		Delete: resourceTencentCloudCvmRebootInstanceDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"force_reboot": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"stop_type"},
				Description:   "This parameter has been disused. We recommend using StopType instead. Note that ForceReboot and StopType parameters cannot be specified at the same time. Whether to forcibly restart an instance after a normal restart fails. Valid values are `TRUE` and `FALSE`. Default value: FALSE.",
				Deprecated:    "It has been deprecated from version 1.81.21. Please use `stop_type` instead.",
			},

			"stop_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"force_reboot"},
				Description:   "Shutdown type. Valid values: `SOFT`: soft shutdown; `HARD`: hard shutdown; `SOFT_FIRST`: perform a soft shutdown first, and perform a hard shutdown if the soft shutdown fails. Default value: SOFT.",
			},
		},
	}
}

func resourceTencentCloudCvmRebootInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_reboot_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
	)
	var (
		request  = cvm.NewRebootInstancesRequest()
		response = cvm.NewRebootInstancesResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceIds = []*string{helper.String(v.(string))}
	}

	if v, ok := d.GetOkExists("force_reboot"); ok {
		request.ForceReboot = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("stop_type"); ok {
		request.StopType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RebootInstancesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm reboot instance failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(instanceId)

	return resourceTencentCloudCvmRebootInstanceRead(d, meta)
}

func resourceTencentCloudCvmRebootInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_reboot_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmRebootInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_reboot_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
