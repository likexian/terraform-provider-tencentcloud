package tse

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTseCngwRouteRateLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwRouteRateLimitCreate,
		Read:   resourceTencentCloudTseCngwRouteRateLimitRead,
		Update: resourceTencentCloudTseCngwRouteRateLimitUpdate,
		Delete: resourceTencentCloudTseCngwRouteRateLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"route_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Route id, or route name.",
			},

			"limit_detail": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "rate limit configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "status of service rate limit.",
						},
						"qps_thresholds": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "qps threshold.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "qps threshold unit.Reference value:`second`,`minute`,`hour`,`day`,`month`,`year`.",
									},
									"max": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "the max threshold.",
									},
								},
							},
						},
						"limit_by": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "basis for service rate limit.Reference value:`ip`,`service`,`consumer`,`credential`,`path`,`header`.",
						},
						"response_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "response strategy.Reference value:`url`: forward request according to url,`text`: response configuration,`default`: return directly.",
						},
						"hide_client_headers": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "whether to hide the headers of client.",
						},
						"is_delay": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "whether to enable request queuing.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "request paths that require rate limit.",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "request headers that require rate limit.",
						},
						"external_redis": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "external redis information, maybe null.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redis_host": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "redis ip, maybe null.",
									},
									"redis_password": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "redis password, maybe null.",
									},
									"redis_port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "redis port, maybe null.",
									},
									"redis_timeout": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "redis timeout, unit: `ms`, maybe null.",
									},
								},
							},
						},
						"policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "counter policy.Reference value:`local`,`redis`,`external_redis`.",
						},
						"rate_limit_response": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "response configuration, the response strategy is text, maybe null.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"body": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "custom response body, maybe bull.",
									},
									"headers": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "headrs.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "key of header.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "value of header.",
												},
											},
										},
									},
									"http_status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "http status code.",
									},
								},
							},
						},
						"rate_limit_response_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "request forwarding address, maybe null.",
						},
						"line_up_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "queue time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTseCngwRouteRateLimitCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_route_rate_limit.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayRouteRateLimitRequest()
		gatewayId string
		routeId   string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_id"); ok {
		routeId = v.(string)
		request.Id = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "limit_detail"); ok {
		cloudNativeAPIGatewayRateLimitDetail := tse.CloudNativeAPIGatewayRateLimitDetail{}
		if v, ok := dMap["enabled"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.Enabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["qps_thresholds"]; ok {
			for _, item := range v.([]interface{}) {
				qpsThresholdsMap := item.(map[string]interface{})
				qpsThreshold := tse.QpsThreshold{}
				if v, ok := qpsThresholdsMap["unit"]; ok {
					qpsThreshold.Unit = helper.String(v.(string))
				}
				if v, ok := qpsThresholdsMap["max"]; ok {
					qpsThreshold.Max = helper.IntInt64(v.(int))
				}
				cloudNativeAPIGatewayRateLimitDetail.QpsThresholds = append(cloudNativeAPIGatewayRateLimitDetail.QpsThresholds, &qpsThreshold)
			}
		}
		if v, ok := dMap["limit_by"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.LimitBy = helper.String(v.(string))
		}
		if v, ok := dMap["response_type"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.ResponseType = helper.String(v.(string))
		}
		if v, ok := dMap["hide_client_headers"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.HideClientHeaders = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_delay"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.IsDelay = helper.Bool(v.(bool))
		}
		if v, ok := dMap["path"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.Path = helper.String(v.(string))
		}
		if v, ok := dMap["header"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.Header = helper.String(v.(string))
		}
		if externalRedisMap, ok := helper.InterfaceToMap(dMap, "external_redis"); ok {
			externalRedis := tse.ExternalRedis{}
			if v, ok := externalRedisMap["redis_host"]; ok {
				externalRedis.RedisHost = helper.String(v.(string))
			}
			if v, ok := externalRedisMap["redis_password"]; ok {
				externalRedis.RedisPassword = helper.String(v.(string))
			}
			if v, ok := externalRedisMap["redis_port"]; ok {
				externalRedis.RedisPort = helper.IntInt64(v.(int))
			}
			if v, ok := externalRedisMap["redis_timeout"]; ok {
				externalRedis.RedisTimeout = helper.IntInt64(v.(int))
			}
			cloudNativeAPIGatewayRateLimitDetail.ExternalRedis = &externalRedis
		}
		if v, ok := dMap["policy"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.Policy = helper.String(v.(string))
		}
		if rateLimitResponseMap, ok := helper.InterfaceToMap(dMap, "rate_limit_response"); ok {
			rateLimitResponse := tse.RateLimitResponse{}
			if v, ok := rateLimitResponseMap["body"]; ok {
				rateLimitResponse.Body = helper.String(v.(string))
			}
			if v, ok := rateLimitResponseMap["headers"]; ok {
				for _, item := range v.([]interface{}) {
					headersMap := item.(map[string]interface{})
					kVMapping := tse.KVMapping{}
					if v, ok := headersMap["key"]; ok {
						kVMapping.Key = helper.String(v.(string))
					}
					if v, ok := headersMap["value"]; ok {
						kVMapping.Value = helper.String(v.(string))
					}
					rateLimitResponse.Headers = append(rateLimitResponse.Headers, &kVMapping)
				}
			}
			if v, ok := rateLimitResponseMap["http_status"]; ok {
				rateLimitResponse.HttpStatus = helper.IntInt64(v.(int))
			}
			cloudNativeAPIGatewayRateLimitDetail.RateLimitResponse = &rateLimitResponse
		}
		if v, ok := dMap["rate_limit_response_url"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.RateLimitResponseUrl = helper.String(v.(string))
		}
		if v, ok := dMap["line_up_time"]; ok {
			cloudNativeAPIGatewayRateLimitDetail.LineUpTime = helper.IntInt64(v.(int))
		}
		request.LimitDetail = &cloudNativeAPIGatewayRateLimitDetail
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().CreateCloudNativeAPIGatewayRouteRateLimit(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwRouteRateLimit failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(gatewayId + tccommon.FILED_SP + routeId)

	return resourceTencentCloudTseCngwRouteRateLimitRead(d, meta)
}

func resourceTencentCloudTseCngwRouteRateLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_route_rate_limit.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeId := idSplit[1]

	cngwRouteRateLimit, err := service.DescribeTseCngwRouteRateLimitById(ctx, gatewayId, routeId)
	if err != nil {
		return err
	}

	if cngwRouteRateLimit == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwRouteRateLimit` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("route_id", routeId)

	if cngwRouteRateLimit != nil {
		limitDetailMap := map[string]interface{}{}

		if cngwRouteRateLimit.Enabled != nil {
			limitDetailMap["enabled"] = cngwRouteRateLimit.Enabled
		}

		if cngwRouteRateLimit.QpsThresholds != nil {
			qpsThresholdsList := []interface{}{}
			for _, qpsThresholds := range cngwRouteRateLimit.QpsThresholds {
				qpsThresholdsMap := map[string]interface{}{}

				if qpsThresholds.Unit != nil {
					qpsThresholdsMap["unit"] = qpsThresholds.Unit
				}

				if qpsThresholds.Max != nil {
					qpsThresholdsMap["max"] = qpsThresholds.Max
				}

				qpsThresholdsList = append(qpsThresholdsList, qpsThresholdsMap)
			}

			limitDetailMap["qps_thresholds"] = qpsThresholdsList
		}

		if cngwRouteRateLimit.LimitBy != nil {
			limitDetailMap["limit_by"] = cngwRouteRateLimit.LimitBy
		}

		if cngwRouteRateLimit.ResponseType != nil {
			limitDetailMap["response_type"] = cngwRouteRateLimit.ResponseType
		}

		if cngwRouteRateLimit.HideClientHeaders != nil {
			limitDetailMap["hide_client_headers"] = cngwRouteRateLimit.HideClientHeaders
		}

		if cngwRouteRateLimit.IsDelay != nil {
			limitDetailMap["is_delay"] = cngwRouteRateLimit.IsDelay
		}

		if cngwRouteRateLimit.Path != nil {
			limitDetailMap["path"] = cngwRouteRateLimit.Path
		}

		if cngwRouteRateLimit.Header != nil {
			limitDetailMap["header"] = cngwRouteRateLimit.Header
		}

		if cngwRouteRateLimit.ExternalRedis != nil {
			externalRedisMap := map[string]interface{}{}

			if cngwRouteRateLimit.ExternalRedis.RedisHost != nil {
				externalRedisMap["redis_host"] = cngwRouteRateLimit.ExternalRedis.RedisHost
			}

			if cngwRouteRateLimit.ExternalRedis.RedisPassword != nil {
				externalRedisMap["redis_password"] = cngwRouteRateLimit.ExternalRedis.RedisPassword
			}

			if cngwRouteRateLimit.ExternalRedis.RedisPort != nil {
				externalRedisMap["redis_port"] = cngwRouteRateLimit.ExternalRedis.RedisPort
			}

			if cngwRouteRateLimit.ExternalRedis.RedisTimeout != nil {
				externalRedisMap["redis_timeout"] = cngwRouteRateLimit.ExternalRedis.RedisTimeout
			}

			limitDetailMap["external_redis"] = []interface{}{externalRedisMap}
		}

		if cngwRouteRateLimit.Policy != nil {
			limitDetailMap["policy"] = cngwRouteRateLimit.Policy
		}

		if cngwRouteRateLimit.RateLimitResponse != nil {
			rateLimitResponseMap := map[string]interface{}{}

			if cngwRouteRateLimit.RateLimitResponse.Body != nil {
				rateLimitResponseMap["body"] = cngwRouteRateLimit.RateLimitResponse.Body
			}

			if cngwRouteRateLimit.RateLimitResponse.Headers != nil {
				headersList := []interface{}{}
				for _, headers := range cngwRouteRateLimit.RateLimitResponse.Headers {
					headersMap := map[string]interface{}{}

					if headers.Key != nil {
						headersMap["key"] = headers.Key
					}

					if headers.Value != nil {
						headersMap["value"] = headers.Value
					}

					headersList = append(headersList, headersMap)
				}

				rateLimitResponseMap["headers"] = headersList
			}

			if cngwRouteRateLimit.RateLimitResponse.HttpStatus != nil {
				rateLimitResponseMap["http_status"] = cngwRouteRateLimit.RateLimitResponse.HttpStatus
			}

			limitDetailMap["rate_limit_response"] = []interface{}{rateLimitResponseMap}
		}

		if cngwRouteRateLimit.RateLimitResponseUrl != nil {
			limitDetailMap["rate_limit_response_url"] = cngwRouteRateLimit.RateLimitResponseUrl
		}

		if cngwRouteRateLimit.LineUpTime != nil {
			limitDetailMap["line_up_time"] = cngwRouteRateLimit.LineUpTime
		}

		_ = d.Set("limit_detail", []interface{}{limitDetailMap})
	}

	return nil
}

func resourceTencentCloudTseCngwRouteRateLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_route_rate_limit.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tse.NewModifyCloudNativeAPIGatewayRouteRateLimitRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeId := idSplit[1]

	request.GatewayId = &gatewayId
	request.Id = &routeId

	immutableArgs := []string{"gateway_id", "route_id", "limit_detail"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("limit_detail") {
		if dMap, ok := helper.InterfacesHeadMap(d, "limit_detail"); ok {
			cloudNativeAPIGatewayRateLimitDetail := tse.CloudNativeAPIGatewayRateLimitDetail{}
			if v, ok := dMap["enabled"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := dMap["qps_thresholds"]; ok {
				for _, item := range v.([]interface{}) {
					qpsThresholdsMap := item.(map[string]interface{})
					qpsThreshold := tse.QpsThreshold{}
					if v, ok := qpsThresholdsMap["unit"]; ok {
						qpsThreshold.Unit = helper.String(v.(string))
					}
					if v, ok := qpsThresholdsMap["max"]; ok {
						qpsThreshold.Max = helper.IntInt64(v.(int))
					}
					cloudNativeAPIGatewayRateLimitDetail.QpsThresholds = append(cloudNativeAPIGatewayRateLimitDetail.QpsThresholds, &qpsThreshold)
				}
			}
			if v, ok := dMap["limit_by"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.LimitBy = helper.String(v.(string))
			}
			if v, ok := dMap["response_type"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.ResponseType = helper.String(v.(string))
			}
			if v, ok := dMap["hide_client_headers"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.HideClientHeaders = helper.Bool(v.(bool))
			}
			if v, ok := dMap["is_delay"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.IsDelay = helper.Bool(v.(bool))
			}
			if v, ok := dMap["path"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.Path = helper.String(v.(string))
			}
			if v, ok := dMap["header"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.Header = helper.String(v.(string))
			}
			if externalRedisMap, ok := helper.InterfaceToMap(dMap, "external_redis"); ok {
				externalRedis := tse.ExternalRedis{}
				if v, ok := externalRedisMap["redis_host"]; ok {
					externalRedis.RedisHost = helper.String(v.(string))
				}
				if v, ok := externalRedisMap["redis_password"]; ok {
					externalRedis.RedisPassword = helper.String(v.(string))
				}
				if v, ok := externalRedisMap["redis_port"]; ok {
					externalRedis.RedisPort = helper.IntInt64(v.(int))
				}
				if v, ok := externalRedisMap["redis_timeout"]; ok {
					externalRedis.RedisTimeout = helper.IntInt64(v.(int))
				}
				cloudNativeAPIGatewayRateLimitDetail.ExternalRedis = &externalRedis
			}
			if v, ok := dMap["policy"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.Policy = helper.String(v.(string))
			}
			if rateLimitResponseMap, ok := helper.InterfaceToMap(dMap, "rate_limit_response"); ok {
				rateLimitResponse := tse.RateLimitResponse{}
				if v, ok := rateLimitResponseMap["body"]; ok {
					rateLimitResponse.Body = helper.String(v.(string))
				}
				if v, ok := rateLimitResponseMap["headers"]; ok {
					for _, item := range v.([]interface{}) {
						headersMap := item.(map[string]interface{})
						kVMapping := tse.KVMapping{}
						if v, ok := headersMap["key"]; ok {
							kVMapping.Key = helper.String(v.(string))
						}
						if v, ok := headersMap["value"]; ok {
							kVMapping.Value = helper.String(v.(string))
						}
						rateLimitResponse.Headers = append(rateLimitResponse.Headers, &kVMapping)
					}
				}
				if v, ok := rateLimitResponseMap["http_status"]; ok {
					rateLimitResponse.HttpStatus = helper.IntInt64(v.(int))
				}
				cloudNativeAPIGatewayRateLimitDetail.RateLimitResponse = &rateLimitResponse
			}
			if v, ok := dMap["rate_limit_response_url"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.RateLimitResponseUrl = helper.String(v.(string))
			}
			if v, ok := dMap["line_up_time"]; ok {
				cloudNativeAPIGatewayRateLimitDetail.LineUpTime = helper.IntInt64(v.(int))
			}
			request.LimitDetail = &cloudNativeAPIGatewayRateLimitDetail
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().ModifyCloudNativeAPIGatewayRouteRateLimit(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwRouteRateLimit failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwRouteRateLimitRead(d, meta)
}

func resourceTencentCloudTseCngwRouteRateLimitDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_route_rate_limit.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	routeId := idSplit[1]

	if err := service.DeleteTseCngwRouteRateLimitById(ctx, gatewayId, routeId); err != nil {
		return err
	}

	return nil
}
