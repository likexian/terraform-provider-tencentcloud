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

func ResourceTencentCloudCvmProgramFpgaImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmProgramFpgaImageCreate,
		Read:   resourceTencentCloudCvmProgramFpgaImageRead,
		Delete: resourceTencentCloudCvmProgramFpgaImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID information of the instance.",
			},

			"fpga_url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS URL address of the FPGA image file.",
			},

			"dbd_fs": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "The DBDF number of the FPGA card on the instance, if left blank, the FPGA image will be burned to all FPGA cards owned by the instance by default.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dry_run": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Trial run, will not perform the actual burning action, the default is False.",
			},
		},
	}
}

func resourceTencentCloudCvmProgramFpgaImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_program_fpga_image.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
	)
	var (
		request  = cvm.NewProgramFpgaImageRequest()
		response = cvm.NewProgramFpgaImageResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fpga_url"); ok {
		request.FPGAUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dbd_fs"); ok {
		dBDFsSet := v.(*schema.Set).List()
		for i := range dBDFsSet {
			dBDFs := dBDFsSet[i].(string)
			request.DBDFs = append(request.DBDFs, helper.String(dBDFs))
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ProgramFpgaImageWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm program fpga image failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	d.SetId(instanceId)

	return resourceTencentCloudCvmProgramFpgaImageRead(d, meta)
}

func resourceTencentCloudCvmProgramFpgaImageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_program_fpga_image.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmProgramFpgaImageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_program_fpga_image.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
