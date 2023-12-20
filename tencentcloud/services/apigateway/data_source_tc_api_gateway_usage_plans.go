package apigateway

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func DataSourceTencentCloudAPIGatewayUsagePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayUsagePlansRead,

		Schema: map[string]*schema.Schema{
			"usage_plan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the usage plan.",
			},
			"usage_plan_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the usage plan.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			// Computed values.
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of usage plans.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"usage_plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the usage plan.",
						},
						"usage_plan_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the usage plan.",
						},
						"usage_plan_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom usage plan description.",
						},
						"max_request_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of requests allowed. Valid value formats: `-1`, `[1,99999999]`. The default value is -1, which indicates no limit.",
						},
						"max_request_num_pre_sec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Limit of requests per second. Valid values formats: `-1`, `[1,2000]`. The default value is -1, which indicates no limit.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time in the format of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayUsagePlansRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_usage_plans.read")

	var (
		logId                      = tccommon.GetLogId(tccommon.ContextNil)
		ctx                        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService          = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		infos                      []*apigateway.UsagePlanStatusInfo
		list                       []map[string]interface{}
		usagePlanId, usagePlanName string
		err                        error
	)

	if v, ok := d.GetOk("usage_plan_id"); ok {
		usagePlanId = v.(string)
	}
	if v, ok := d.GetOk("usage_plan_name"); ok {
		usagePlanName = v.(string)
	}

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeUsagePlansStatus(ctx, usagePlanId, usagePlanName)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		var infoMap = make(map[string]interface{}, 7)
		infoMap["usage_plan_id"] = info.UsagePlanId
		infoMap["usage_plan_name"] = info.UsagePlanName
		infoMap["usage_plan_desc"] = info.UsagePlanDesc
		infoMap["max_request_num"] = info.MaxRequestNum
		infoMap["max_request_num_pre_sec"] = info.MaxRequestNumPreSec
		infoMap["modify_time"] = info.ModifiedTime
		infoMap["create_time"] = info.CreatedTime

		list = append(list, infoMap)
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{usagePlanId, usagePlanName}, tccommon.FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
