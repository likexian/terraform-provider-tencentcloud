package ci

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCiMediaSuperResolutionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaSuperResolutionTemplateCreate,
		Read:   resourceTencentCloudCiMediaSuperResolutionTemplateRead,
		Update: resourceTencentCloudCiMediaSuperResolutionTemplateUpdate,
		Delete: resourceTencentCloudCiMediaSuperResolutionTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "bucket name.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"resolution": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resolution Options sdtohd: Standard Definition to Ultra Definition, hdto4k: HD to 4K.",
			},

			"enable_scale_up": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Auto scaling switch, off by default.",
			},

			"version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "version, default value Base, Base: basic version, Enhance: enhanced version.",
			},
		},
	}
}

func resourceTencentCloudCiMediaSuperResolutionTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_super_resolution_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaSuperResolutionTemplateOptions{
			Tag: "SuperResolution",
		}
		bucket     string
		templateId string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("resolution"); ok {
		request.Resolution = v.(string)
	}

	if v, ok := d.GetOk("enable_scale_up"); ok {
		request.EnableScaleUp = v.(string)
	}

	if v, ok := d.GetOk("version"); ok {
		request.Version = v.(string)
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaSuperResolutionTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaSnapshotTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSuperResolutionTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaSuperResolutionTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSuperResolutionTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_super_resolution_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	template, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if template.Name != "" {
		_ = d.Set("name", template.Name)
	}

	if template.SuperResolution != nil {
		mediaSuperResolutionTemplate := template.SuperResolution
		if mediaSuperResolutionTemplate.Resolution != "" {
			_ = d.Set("resolution", mediaSuperResolutionTemplate.Resolution)
		}

		if mediaSuperResolutionTemplate.EnableScaleUp != "" {
			_ = d.Set("enable_scale_up", mediaSuperResolutionTemplate.EnableScaleUp)
		}

		if mediaSuperResolutionTemplate.Version != "" {
			_ = d.Set("version", mediaSuperResolutionTemplate.Version)
		}
	}

	return nil
}

func resourceTencentCloudCiMediaSuperResolutionTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_super_resolution_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaSuperResolutionTemplateOptions{
		Tag: "SuperResolution",
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("resolution"); ok {
		request.Resolution = v.(string)
	}

	if d.HasChange("enable_scale_up") {
		if v, ok := d.GetOk("enable_scale_up"); ok {
			request.EnableScaleUp = v.(string)
		}
	}

	if d.HasChange("version") {
		if v, ok := d.GetOk("version"); ok {
			request.Version = v.(string)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaSuperResolutionTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaSuperResolutionTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSuperResolutionTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaSuperResolutionTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSuperResolutionTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_super_resolution_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteCiMediaTemplateById(ctx, bucket, templateId); err != nil {
		return err
	}

	return nil
}
