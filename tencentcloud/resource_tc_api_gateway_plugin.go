package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAPIGatewayPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayPluginCreate,
		Read:   resourceTencentCloudAPIGatewayPluginRead,
		Update: resourceTencentCloudAPIGatewayPluginUpdate,
		Delete: resourceTencentCloudAPIGatewayPluginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"plugin_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Name of the user define plugin. It must start with a letter and end with letter or number, the rest can contain letters, numbers and dashes(-). The length range is from 2 to 50.",
			},
			"plugin_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Type of plugin. Now support IPControl, TrafficControl, Cors, CustomReq, CustomAuth, Routing, TrafficControlByParameter, CircuitBreaker, ProxyCache.",
			},
			"plugin_data": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Statement to define plugin.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description of plugin.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayPluginCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = apiGateway.NewCreatePluginRequest()
		response = apiGateway.NewCreatePluginResponse()
	)

	if v, ok := d.GetOk("plugin_name"); ok {
		request.PluginName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plugin_type"); ok {
		request.PluginType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plugin_data"); ok {
		request.PluginData = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().CreatePlugin(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create apiGateway plugin failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(*response.Response.Result.PluginId)

	return resourceTencentCloudAPIGatewayPluginRead(d, meta)
}

func resourceTencentCloudAPIGatewayPluginRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		pluginId = d.Id()
	)

	plugin, err := service.DescribeApiGatewayPluginById(ctx, pluginId)
	if err != nil {
		return err
	}

	if plugin == nil {
		d.SetId("")
		return fmt.Errorf("resource `ApiGatewayPlugin` %s does not exist", d.Id())
	}

	if plugin.PluginName != nil {
		_ = d.Set("plugin_name", plugin.PluginName)
	}

	if plugin.PluginType != nil {
		_ = d.Set("plugin_type", plugin.PluginType)
	}

	if plugin.PluginData != nil {
		_ = d.Set("plugin_data", plugin.PluginData)
	}

	if plugin.Description != nil {
		_ = d.Set("description", plugin.Description)
	}

	return nil
}

func resourceTencentCloudAPIGatewayPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		request  = apiGateway.NewModifyPluginRequest()
		pluginId = d.Id()
	)

	unsupportedUpdateFields := []string{"plugin_type"}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_api_gateway_plugin update on %s is not support yet", field)
		}
	}

	request.PluginId = &pluginId

	if d.HasChange("plugin_name") {
		if v, ok := d.GetOk("plugin_name"); ok {
			request.PluginName = helper.String(v.(string))
		}
	}

	if d.HasChange("plugin_data") {
		if v, ok := d.GetOk("plugin_data"); ok {
			request.PluginData = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAPIGatewayClient().ModifyPlugin(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update apiGateway plugin failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudAPIGatewayPluginRead(d, meta)
}

func resourceTencentCloudAPIGatewayPluginDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_plugin.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		ctx      = context.WithValue(context.TODO(), logIdKey, logId)
		service  = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		pluginId = d.Id()
	)

	if err := service.DeleteApiGatewayPluginById(ctx, pluginId); err != nil {
		return err
	}

	return nil
}
