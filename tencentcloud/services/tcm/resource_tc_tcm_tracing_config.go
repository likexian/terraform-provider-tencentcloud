package tcm

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
)

func ResourceTencentCloudTcmTracingConfig() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTcmTracingConfigRead,
		Create: resourceTencentCloudTcmTracingConfigCreate,
		Update: resourceTencentCloudTcmTracingConfigUpdate,
		Delete: resourceTencentCloudTcmTracingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"mesh_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mesh ID.",
			},

			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether enable tracing.",
			},

			"apm": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "APM config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether enable APM.",
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance id of the APM.",
						},
					},
				},
			},

			"sampling": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Tracing sampling, 0.0-1.0.",
			},

			"zipkin": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Third party zipkin config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zipkin address.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcmTracingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcm_tracing_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var meshId string
	if v, ok := d.GetOk("mesh_id"); ok {
		meshId = v.(string)
	}

	d.SetId(meshId)
	return resourceTencentCloudTcmTracingConfigUpdate(d, meta)
}

func resourceTencentCloudTcmTracingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcm_tracing_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	meshId := d.Id()

	response, err := service.DescribeTcmMesh(ctx, meshId)

	if err != nil {
		return err
	}

	if response == nil {
		d.SetId("")
		return fmt.Errorf("resource `tracingConfig` %s does not exist", meshId)
	}

	mesh := response.Mesh
	if mesh.MeshId != nil {
		_ = d.Set("mesh_id", mesh.MeshId)
	}

	tracing := mesh.Config.Tracing
	if tracing != nil {
		if tracing.Enable != nil {
			_ = d.Set("enable", tracing.Enable)
		}
		apmMap := map[string]interface{}{}
		if tracing.APM.Enable != nil {
			apmMap["enable"] = tracing.APM.Enable
		}
		if tracing.APM.Region != nil {
			apmMap["region"] = tracing.APM.Region
		}
		if tracing.APM.InstanceId != nil {
			apmMap["instance_id"] = tracing.APM.InstanceId
		}

		_ = d.Set("apm", []interface{}{apmMap})
	}

	if tracing.Sampling != nil {
		_ = d.Set("sampling", tracing.Sampling)
	}

	if tracing.Zipkin != nil {
		zipkinMap := map[string]interface{}{}
		if tracing.Zipkin.Address != nil {
			zipkinMap["address"] = tracing.Zipkin.Address
		}

		_ = d.Set("zipkin", []interface{}{zipkinMap})
	}

	return nil
}

func resourceTencentCloudTcmTracingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcm_tracing_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tcm.NewModifyTracingConfigRequest()

	meshId := d.Id()

	request.MeshId = &meshId

	if v, ok := d.GetOk("enable"); ok {
		request.Enable = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "apm"); ok {
		aPM := tcm.APM{}
		if v, ok := dMap["enable"]; ok {
			aPM.Enable = helper.Bool(v.(bool))
		}
		if v, ok := dMap["region"]; ok {
			aPM.Region = helper.String(v.(string))
		}
		if v, ok := dMap["instance_id"]; ok {
			aPM.InstanceId = helper.String(v.(string))
		}

		request.APM = &aPM
	}

	if v, ok := d.GetOk("sampling"); ok {
		request.Sampling = helper.Float64(v.(float64))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "zipkin"); ok {
		tracingZipkin := tcm.TracingZipkin{}
		if v, ok := dMap["address"]; ok {
			tracingZipkin.Address = helper.String(v.(string))
		}

		request.Zipkin = &tracingZipkin
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcmClient().ModifyTracingConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tcm tracingConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmTracingConfigRead(d, meta)
}

func resourceTencentCloudTcmTracingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcm_tracing_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
