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

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCiMediaVoiceSeparateTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaVoiceSeparateTemplateCreate,
		Read:   resourceTencentCloudCiMediaVoiceSeparateTemplateRead,
		Update: resourceTencentCloudCiMediaVoiceSeparateTemplateUpdate,
		Delete: resourceTencentCloudCiMediaVoiceSeparateTemplateDelete,
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

			"audio_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Output audio IsAudio: output human voice, IsBackground: output background sound, AudioAndBackground: output vocal and background sound.",
			},

			"audio_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "audio configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Codec format, value aac, mp3, flac, amr.",
						},
						"samplerate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sampling Rate- 1: Unit: Hz- 2: Optional 8000, 11025, 22050, 32000, 44100, 48000, 96000- 3: When Codec is set to aac/flac, 8000 is not supported- 4: When Codec is set to mp3, 8000 and 96000 are not supported- 5: When Codec is set to amr, only 8000 is supported.",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Original audio bit rate, unit: Kbps, Value range: [8, 1000].",
						},
						"channels": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "number of channels- When Codec is set to aac/flac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3, support 1, 2- When Codec is set to amr, only 1 is supported.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaVoiceSeparateTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_voice_separate_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaVoiceSeparateTemplateOptions{
			Tag: "VoiceSeparate",
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

	if v, ok := d.GetOk("audio_mode"); ok {
		request.AudioMode = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio_config"); ok {
		audioConfig := cos.AudioConfig{}
		if v, ok := dMap["codec"]; ok {
			audioConfig.Codec = v.(string)
		}
		if v, ok := dMap["samplerate"]; ok {
			audioConfig.Samplerate = v.(string)
		}
		if v, ok := dMap["bitrate"]; ok {
			audioConfig.Bitrate = v.(string)
		}
		if v, ok := dMap["channels"]; ok {
			audioConfig.Channels = v.(string)
		}
		request.AudioConfig = &audioConfig
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaVoiceSeparateTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaVoiceSeparateTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaVoiceSeparateTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaVoiceSeparateTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVoiceSeparateTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_voice_separate_template.read")()
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

	if template.VoiceSeparate != nil {
		mediaVoiceSeparateTemplate := template.VoiceSeparate
		if mediaVoiceSeparateTemplate.AudioMode != "" {
			_ = d.Set("audio_mode", mediaVoiceSeparateTemplate.AudioMode)
		}

		if mediaVoiceSeparateTemplate.AudioConfig != nil {
			audioConfigMap := map[string]interface{}{}

			if mediaVoiceSeparateTemplate.AudioConfig.Codec != "" {
				audioConfigMap["codec"] = mediaVoiceSeparateTemplate.AudioConfig.Codec
			}

			if mediaVoiceSeparateTemplate.AudioConfig.Samplerate != "" {
				audioConfigMap["samplerate"] = mediaVoiceSeparateTemplate.AudioConfig.Samplerate
			}

			if mediaVoiceSeparateTemplate.AudioConfig.Bitrate != "" {
				audioConfigMap["bitrate"] = mediaVoiceSeparateTemplate.AudioConfig.Bitrate
			}

			if mediaVoiceSeparateTemplate.AudioConfig.Channels != "" {
				audioConfigMap["channels"] = mediaVoiceSeparateTemplate.AudioConfig.Channels
			}

			_ = d.Set("audio_config", []interface{}{audioConfigMap})
		}
	}

	return nil
}

func resourceTencentCloudCiMediaVoiceSeparateTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_voice_separate_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaVoiceSeparateTemplateOptions{
		Tag: "VoiceSeparate",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("audio_mode"); ok {
		request.AudioMode = v.(string)
	}

	if d.HasChange("audio_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "audio_config"); ok {
			audioConfig := cos.AudioConfig{}
			if v, ok := dMap["codec"]; ok {
				audioConfig.Codec = v.(string)
			}
			if v, ok := dMap["samplerate"]; ok {
				audioConfig.Samplerate = v.(string)
			}
			if v, ok := dMap["bitrate"]; ok {
				audioConfig.Bitrate = v.(string)
			}
			if v, ok := dMap["channels"]; ok {
				audioConfig.Channels = v.(string)
			}
			request.AudioConfig = &audioConfig
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaVoiceSeparateTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaVoiceSeparateTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaVoiceSeparateTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaVoiceSeparateTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVoiceSeparateTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_voice_separate_template.delete")()
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
