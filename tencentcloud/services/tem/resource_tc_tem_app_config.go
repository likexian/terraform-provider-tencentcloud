package tem

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTemAppConfig() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemAppConfigRead,
		Create: resourceTencentCloudTemAppConfigCreate,
		Update: resourceTencentCloudTemAppConfigUpdate,
		Delete: resourceTencentCloudTemAppConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "environment ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "appConfig name.",
			},

			"config_data": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "payload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemAppConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_app_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = tem.NewCreateConfigDataRequest()
		environmentId string
		name          string
	)

	if v, ok := d.GetOk("environment_id"); ok {
		environmentId = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_data"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			pair := tem.Pair{}
			if v, ok := dMap["key"]; ok {
				pair.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				pair.Value = helper.String(v.(string))
			}
			request.Data = append(request.Data, &pair)
		}

	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTemClient().CreateConfigData(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tem appConfig failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(environmentId + tccommon.FILED_SP + name)
	return resourceTencentCloudTemAppConfigRead(d, meta)
}

func resourceTencentCloudTemAppConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_appConfig.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	name := idSplit[1]

	appConfig, err := service.DescribeTemAppConfig(ctx, environmentId, name)

	if err != nil {
		return err
	}

	if appConfig == nil {
		d.SetId("")
		return fmt.Errorf("resource `appConfig` %s does not exist", name)
	}

	_ = d.Set("environment_id", environmentId)

	if appConfig.Name != nil {
		_ = d.Set("name", appConfig.Name)
	}

	if appConfig.Data != nil {
		dataList := []interface{}{}
		for _, data := range appConfig.Data {
			dataMap := map[string]interface{}{}
			if data.Key != nil {
				dataMap["key"] = data.Key
			}
			if data.Value != nil {
				dataMap["value"] = data.Value
			}

			dataList = append(dataList, dataMap)
		}
		_ = d.Set("config_data", dataList)
	}

	return nil
}

func resourceTencentCloudTemAppConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_app_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tem.NewModifyConfigDataRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	name := idSplit[1]

	request.EnvironmentId = &environmentId
	request.Name = &name

	if d.HasChange("environment_id") {
		return fmt.Errorf("`environment_id` do not support change now.")
	}

	if d.HasChange("name") {
		return fmt.Errorf("`name` do not support change now.")
	}

	if d.HasChange("config_data") {
		if v, ok := d.GetOk("config_data"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				pair := tem.Pair{}
				if v, ok := dMap["key"]; ok {
					pair.Key = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					pair.Value = helper.String(v.(string))
				}
				request.Data = append(request.Data, &pair)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTemClient().ModifyConfigData(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTemAppConfigRead(d, meta)
}

func resourceTencentCloudTemAppConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_app_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	name := idSplit[1]

	if err := service.DeleteTemAppConfigById(ctx, environmentId, name); err != nil {
		return err
	}

	return nil
}
