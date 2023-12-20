package apigateway

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAPIGatewayAPIApps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayAPIAppsRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"api_app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Api app ID.",
			},

			"api_app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Api app name.",
			},

			"api_app_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of ApiApp.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ApiApp ID.",
						},
						"api_app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ApiApp Name.",
						},
						"api_app_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ApiApp key.",
						},
						"api_app_secret": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ApiApp secret.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ApiApp create time.",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ApiApp modified time.",
						},
						"api_app_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ApiApp description.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayAPIAppsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_api_apps.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                = tccommon.GetLogId(tccommon.ContextNil)
		ctx                  = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService    = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		apiAppId, apiAppName string
		apiApps              []*apigateway.ApiAppInfo
	)

	if v, ok := d.GetOk("api_app_id"); ok {
		apiAppId = v.(string)
	}

	if v, ok := d.GetOk("api_app_name"); ok {
		apiAppName = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := apiGatewayService.DescribeApiAppList(ctx, apiAppId, apiAppName)
		if e != nil {
			return tccommon.RetryError(e)
		}
		apiApps = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read api_gateway apiApps failed, reason:%+v", logId, err)
		return err
	}

	apiAppList := []interface{}{}
	ids := make([]string, 0, len(apiApps))
	if apiApps != nil {
		for _, item := range apiApps {
			docMap := map[string]interface{}{}
			if item.ApiAppId != nil {
				docMap["api_app_id"] = item.ApiAppId
			}
			if item.ApiAppName != nil {
				docMap["api_app_name"] = item.ApiAppName
			}
			if item.ApiAppKey != nil {
				docMap["api_app_key"] = item.ApiAppKey
			}
			if item.ApiAppSecret != nil {
				docMap["api_app_secret"] = item.ApiAppSecret
			}
			if item.CreatedTime != nil {
				docMap["created_time"] = item.CreatedTime
			}
			if item.ModifiedTime != nil {
				docMap["modified_time"] = item.ModifiedTime
			}
			if item.ApiAppDesc != nil {
				docMap["api_app_desc"] = item.ApiAppDesc
			}
			apiAppList = append(apiAppList, docMap)
			ids = append(ids, *item.ApiAppId)
		}
		_ = d.Set("api_app_list", apiAppList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), apiAppList); e != nil {
			return e
		}
	}

	return nil
}
