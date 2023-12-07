package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoRuleEngine() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoRuleEngineRead,
		Create: resourceTencentCloudTeoRuleEngineCreate,
		Update: resourceTencentCloudTeoRuleEngineUpdate,
		Delete: resourceTencentCloudTeoRuleEngineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the site.",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID.",
			},

			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule name (1 to 255 characters).",
			},

			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule status. Values: `enable`: Enabled; `disable`: Disabled.",
			},

			"tags": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "rule tag list.",
			},

			"rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule items list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"or": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "OR Conditions list of the rule. Rule would be triggered if any of the condition is true.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"and": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "AND Conditions list of the rule. Rule would be triggered if all conditions are true.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Condition operator. Valid values are `equal`, `notequal`.",
												},
												"target": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Condition target. Valid values:- `host`: Host of the URL.- `filename`: filename of the URL.- `extension`: file extension of the URL.- `full_url`: full url.- `url`: path of the URL.",
												},
												"values": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Required:    true,
													Description: "Condition Value.",
												},
												"ignore_case": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to ignore the case of the parameter value, the default value is false.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The parameter name corresponding to the matching type is valid when the Target value is the following, and the valid value cannot be empty: `query_string` (query string): The parameter name of the query string in the URL request under the current site, such as lang and version in lang=cn&version=1; `request_header` (HTTP request header): HTTP request header field name, such as Accept-Language in Accept-Language:zh-CN,zh;q=0.9.",
												},
											},
										},
									},
								},
							},
						},
						"actions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Actions list of the rule. See details in data source `rule_engine_setting`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"normal_action": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Define a normal action.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Action name.",
												},
												"parameters": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Action parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Parameter Name.",
															},
															"values": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Required:    true,
																Description: "Parameter Values.",
															},
														},
													},
												},
											},
										},
									},
									"rewrite_action": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Define a rewrite action.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Action name.",
												},
												"parameters": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Action parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"action": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Action to take on the HEADER. Valid values: `add`, `del`, `set`.",
															},
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Target HEADER name.",
															},
															"values": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Required:    true,
																Description: "Parameter Value.",
															},
														},
													},
												},
											},
										},
									},
									"code_action": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Define a code action.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Action name.",
												},
												"parameters": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Action parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Parameter Name.",
															},
															"values": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Required:    true,
																Description: "Parameter Values.",
															},
															"status_code": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "HTTP status code to use.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"sub_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Actions list of the rule. See details in data source `rule_engine_setting`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tags": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "rule tag list.",
									},
									"rules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Rule items list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"or": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "OR Conditions list of the rule. Rule would be triggered if any of the condition is true.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"and": {
																Type:        schema.TypeList,
																Required:    true,
																Description: "AND Conditions list of the rule. Rule would be triggered if all conditions are true.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"operator": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Condition operator. Valid values are `equal`, `notequal`.",
																		},
																		"target": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Condition target. Valid values:- `host`: Host of the URL.- `filename`: filename of the URL.- `extension`: file extension of the URL.- `full_url`: full url.- `url`: path of the URL.",
																		},
																		"values": {
																			Type: schema.TypeSet,
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																			Required:    true,
																			Description: "Condition Value.",
																		},
																		"ignore_case": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Whether to ignore the case of the parameter value, the default value is false.",
																		},
																		"name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The parameter name corresponding to the matching type is valid when the Target value is the following, and the valid value cannot be empty:- `query_string` (query string): The parameter name of the query string in the URL request under the current site, such as lang and version in lang=cn&version=1; `request_header` (HTTP request header): HTTP request header field name, such as Accept-Language in Accept-Language:zh-CN,zh;q=0.9.",
																		},
																	},
																},
															},
														},
													},
												},
												"actions": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Actions list of the rule. See details in data source `rule_engine_setting`.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"normal_action": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Define a normal action.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"action": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Action name.",
																		},
																		"parameters": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Action parameters.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Parameter Name.",
																					},
																					"values": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																						Required:    true,
																						Description: "Parameter Values.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"rewrite_action": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Define a rewrite action.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"action": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Action name.",
																		},
																		"parameters": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Action parameters.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"action": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Action to take on the HEADER. Valid values: `add`, `del`, `set`.",
																					},
																					"name": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Target HEADER name.",
																					},
																					"values": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																						Required:    true,
																						Description: "Parameter Value.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"code_action": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Define a code action.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"action": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Action name.",
																		},
																		"parameters": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Action parameters.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Parameter Name.",
																					},
																					"values": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																						Required:    true,
																						Description: "Parameter Values.",
																					},
																					"status_code": {
																						Type:        schema.TypeInt,
																						Required:    true,
																						Description: "HTTP status code to use.",
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoRuleEngineCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateRuleRequest()
		response *teo.CreateRuleResponse
		zoneId   string
		ruleId   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.(*schema.Set).List()
		for i := range tagsSet {
			tags := tagsSet[i].(string)
			request.Tags = append(request.Tags, &tags)
		}
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			ruleItem := teo.Rule{}
			if v, ok := dMap["or"]; ok {
				for _, item := range v.([]interface{}) {
					ConditionsMap := item.(map[string]interface{})
					ruleAndConditions := teo.RuleAndConditions{}
					if v, ok := ConditionsMap["and"]; ok {
						for _, item := range v.([]interface{}) {
							ConditionsMap := item.(map[string]interface{})
							ruleCondition := teo.RuleCondition{}
							if v, ok := ConditionsMap["operator"]; ok {
								ruleCondition.Operator = helper.String(v.(string))
							}
							if v, ok := ConditionsMap["target"]; ok {
								ruleCondition.Target = helper.String(v.(string))
							}
							if v, ok := ConditionsMap["ignore_case"]; ok {
								ruleCondition.IgnoreCase = helper.Bool(v.(bool))
							}
							if v, ok := ConditionsMap["name"]; ok {
								ruleCondition.Name = helper.String(v.(string))
							}
							if v, ok := ConditionsMap["values"]; ok {
								valuesSet := v.(*schema.Set).List()
								for i := range valuesSet {
									values := valuesSet[i].(string)
									ruleCondition.Values = append(ruleCondition.Values, &values)
								}
							}
							ruleAndConditions.Conditions = append(ruleAndConditions.Conditions, &ruleCondition)
						}
					}
					ruleItem.Conditions = append(ruleItem.Conditions, &ruleAndConditions)
				}
			}
			if v, ok := dMap["actions"]; ok {
				for _, item := range v.([]interface{}) {
					ActionsMap := item.(map[string]interface{})
					ruleAction := teo.Action{}
					if NormalActionMap, ok := helper.InterfaceToMap(ActionsMap, "normal_action"); ok {
						ruleNormalAction := teo.NormalAction{}
						if v, ok := NormalActionMap["action"]; ok {
							ruleNormalAction.Action = helper.String(v.(string))
						}
						if v, ok := NormalActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								ParametersMap := item.(map[string]interface{})
								ruleNormalActionParams := teo.RuleNormalActionParams{}
								if v, ok := ParametersMap["name"]; ok {
									ruleNormalActionParams.Name = helper.String(v.(string))
								}
								if v, ok := ParametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleNormalActionParams.Values = append(ruleNormalActionParams.Values, &values)
									}
								}
								ruleNormalAction.Parameters = append(ruleNormalAction.Parameters, &ruleNormalActionParams)
							}
						}
						ruleAction.NormalAction = &ruleNormalAction
					}
					if RewriteActionMap, ok := helper.InterfaceToMap(ActionsMap, "rewrite_action"); ok {
						ruleRewriteAction := teo.RewriteAction{}
						if v, ok := RewriteActionMap["action"]; ok {
							ruleRewriteAction.Action = helper.String(v.(string))
						}
						if v, ok := RewriteActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								ParametersMap := item.(map[string]interface{})
								ruleRewriteActionParams := teo.RuleRewriteActionParams{}
								if v, ok := ParametersMap["action"]; ok {
									ruleRewriteActionParams.Action = helper.String(v.(string))
								}
								if v, ok := ParametersMap["name"]; ok {
									ruleRewriteActionParams.Name = helper.String(v.(string))
								}
								if v, ok := ParametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleRewriteActionParams.Values = append(ruleRewriteActionParams.Values, &values)
									}
								}
								ruleRewriteAction.Parameters = append(ruleRewriteAction.Parameters, &ruleRewriteActionParams)
							}
						}
						ruleAction.RewriteAction = &ruleRewriteAction
					}
					if CodeActionMap, ok := helper.InterfaceToMap(ActionsMap, "code_action"); ok {
						ruleCodeAction := teo.CodeAction{}
						if v, ok := CodeActionMap["action"]; ok {
							ruleCodeAction.Action = helper.String(v.(string))
						}
						if v, ok := CodeActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								ParametersMap := item.(map[string]interface{})
								ruleCodeActionParams := teo.RuleCodeActionParams{}
								if v, ok := ParametersMap["name"]; ok {
									ruleCodeActionParams.Name = helper.String(v.(string))
								}
								if v, ok := ParametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleCodeActionParams.Values = append(ruleCodeActionParams.Values, &values)
									}
								}
								if v, ok := ParametersMap["status_code"]; ok {
									ruleCodeActionParams.StatusCode = helper.IntInt64(v.(int))
								}
								ruleCodeAction.Parameters = append(ruleCodeAction.Parameters, &ruleCodeActionParams)
							}
						}
						ruleAction.CodeAction = &ruleCodeAction
					}
					ruleItem.Actions = append(ruleItem.Actions, &ruleAction)
				}
			}
			if v, ok := dMap["sub_rules"]; ok {
				for _, item := range v.([]interface{}) {
					subRulesMap := item.(map[string]interface{})
					subRuleItem := teo.SubRuleItem{}
					if v, ok := subRulesMap["tags"]; ok {
						tagsSet := v.(*schema.Set).List()
						for i := range tagsSet {
							tags := tagsSet[i].(string)
							subRuleItem.Tags = append(subRuleItem.Tags, &tags)
						}
					}
					if v, ok := subRulesMap["rules"]; ok {
						for _, item := range v.([]interface{}) {
							rulesMap := item.(map[string]interface{})
							ruleItem := teo.SubRule{}
							if v, ok := rulesMap["or"]; ok {
								for _, item := range v.([]interface{}) {
									orMap := item.(map[string]interface{})
									ruleAndConditions := teo.RuleAndConditions{}
									if v, ok := orMap["and"]; ok {
										for _, item := range v.([]interface{}) {
											andMap := item.(map[string]interface{})
											ruleCondition := teo.RuleCondition{}
											if v, ok := andMap["operator"]; ok {
												ruleCondition.Operator = helper.String(v.(string))
											}
											if v, ok := andMap["target"]; ok {
												ruleCondition.Target = helper.String(v.(string))
											}
											if v, ok := andMap["values"]; ok {
												valuesSet := v.(*schema.Set).List()
												for i := range valuesSet {
													values := valuesSet[i].(string)
													ruleCondition.Values = append(ruleCondition.Values, &values)
												}
											}
											if v, ok := andMap["ignore_case"]; ok {
												ruleCondition.IgnoreCase = helper.Bool(v.(bool))
											}
											if v, ok := andMap["name"]; ok {
												ruleCondition.Name = helper.String(v.(string))
											}
											ruleAndConditions.Conditions = append(ruleAndConditions.Conditions, &ruleCondition)
										}
									}
									ruleItem.Conditions = append(ruleItem.Conditions, &ruleAndConditions)
								}
							}
							if v, ok := rulesMap["actions"]; ok {
								for _, item := range v.([]interface{}) {
									actionsMap := item.(map[string]interface{})
									ruleAction := teo.Action{}
									if normalActionMap, ok := helper.InterfaceToMap(actionsMap, "normal_action"); ok {
										ruleNormalAction := teo.NormalAction{}
										if v, ok := normalActionMap["action"]; ok {
											ruleNormalAction.Action = helper.String(v.(string))
										}
										if v, ok := normalActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleNormalActionParams := teo.RuleNormalActionParams{}
												if v, ok := parametersMap["name"]; ok {
													ruleNormalActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleNormalActionParams.Values = append(ruleNormalActionParams.Values, &values)
													}
												}
												ruleNormalAction.Parameters = append(ruleNormalAction.Parameters, &ruleNormalActionParams)
											}
										}
										ruleAction.NormalAction = &ruleNormalAction
									}
									if rewriteActionMap, ok := helper.InterfaceToMap(actionsMap, "rewrite_action"); ok {
										ruleRewriteAction := teo.RewriteAction{}
										if v, ok := rewriteActionMap["action"]; ok {
											ruleRewriteAction.Action = helper.String(v.(string))
										}
										if v, ok := rewriteActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleRewriteActionParams := teo.RuleRewriteActionParams{}
												if v, ok := parametersMap["action"]; ok {
													ruleRewriteActionParams.Action = helper.String(v.(string))
												}
												if v, ok := parametersMap["name"]; ok {
													ruleRewriteActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleRewriteActionParams.Values = append(ruleRewriteActionParams.Values, &values)
													}
												}
												ruleRewriteAction.Parameters = append(ruleRewriteAction.Parameters, &ruleRewriteActionParams)
											}
										}
										ruleAction.RewriteAction = &ruleRewriteAction
									}
									if codeActionMap, ok := helper.InterfaceToMap(actionsMap, "code_action"); ok {
										ruleCodeAction := teo.CodeAction{}
										if v, ok := codeActionMap["action"]; ok {
											ruleCodeAction.Action = helper.String(v.(string))
										}
										if v, ok := codeActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleCodeActionParams := teo.RuleCodeActionParams{}
												if v, ok := parametersMap["name"]; ok {
													ruleCodeActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleCodeActionParams.Values = append(ruleCodeActionParams.Values, &values)
													}
												}
												if v, ok := parametersMap["status_code"]; ok {
													ruleCodeActionParams.StatusCode = helper.IntInt64(v.(int))
												}
												ruleCodeAction.Parameters = append(ruleCodeAction.Parameters, &ruleCodeActionParams)
											}
										}
										ruleAction.CodeAction = &ruleCodeAction
									}
									ruleItem.Actions = append(ruleItem.Actions, &ruleAction)
								}
							}
							subRuleItem.Rules = append(subRuleItem.Rules, &ruleItem)
						}
					}
					ruleItem.SubRules = append(ruleItem.SubRules, &subRuleItem)
				}
			}

			request.Rules = append(request.Rules, &ruleItem)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo ruleEngine failed, reason:%+v", logId, err)
		return err
	}

	ruleId = *response.Response.RuleId

	d.SetId(zoneId + FILED_SP + ruleId)
	return resourceTencentCloudTeoRuleEngineRead(d, meta)
}

func resourceTencentCloudTeoRuleEngineRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	ruleEngine, err := service.DescribeTeoRuleEngine(ctx, zoneId, ruleId)

	if err != nil {
		return err
	}

	if ruleEngine == nil {
		d.SetId("")
		return fmt.Errorf("resource `ruleEngine` %s does not exist", ruleId)
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("rule_id", ruleId)

	if ruleEngine.RuleName != nil {
		_ = d.Set("rule_name", ruleEngine.RuleName)
	}

	if ruleEngine.Status != nil {
		_ = d.Set("status", ruleEngine.Status)
	}

	if ruleEngine.Tags != nil {
		_ = d.Set("tags", ruleEngine.Tags)
	}

	if ruleEngine.Rules != nil {
		rulesList := []interface{}{}
		for _, rules := range ruleEngine.Rules {
			rulesMap := map[string]interface{}{}
			if rules.Conditions != nil {
				conditionsList := []interface{}{}
				for _, conditions := range rules.Conditions {
					conditionsMap := map[string]interface{}{}
					if conditions.Conditions != nil {
						conditionsList := []interface{}{}
						for _, conditions := range conditions.Conditions {
							conditionsMap := map[string]interface{}{}
							if conditions.Operator != nil {
								conditionsMap["operator"] = conditions.Operator
							}
							if conditions.Target != nil {
								conditionsMap["target"] = conditions.Target
							}
							if conditions.Values != nil {
								conditionsMap["values"] = conditions.Values
							}
							if conditions.IgnoreCase != nil {
								conditionsMap["ignore_case"] = conditions.IgnoreCase
							}
							if conditions.Name != nil {
								conditionsMap["name"] = conditions.Name
							}

							conditionsList = append(conditionsList, conditionsMap)
						}
						conditionsMap["and"] = conditionsList
					}

					conditionsList = append(conditionsList, conditionsMap)
				}
				rulesMap["or"] = conditionsList
			}
			if rules.Actions != nil {
				actionsList := []interface{}{}
				for _, actions := range rules.Actions {
					actionsMap := map[string]interface{}{}
					if actions.NormalAction != nil {
						normalActionMap := map[string]interface{}{}
						if actions.NormalAction.Action != nil {
							normalActionMap["action"] = actions.NormalAction.Action
						}
						if actions.NormalAction.Parameters != nil {
							parametersList := []interface{}{}
							for _, parameters := range actions.NormalAction.Parameters {
								parametersMap := map[string]interface{}{}
								if parameters.Name != nil {
									parametersMap["name"] = parameters.Name
								}
								if parameters.Values != nil {
									parametersMap["values"] = parameters.Values
								}

								parametersList = append(parametersList, parametersMap)
							}
							normalActionMap["parameters"] = parametersList
						}

						actionsMap["normal_action"] = []interface{}{normalActionMap}
					}
					if actions.RewriteAction != nil {
						rewriteActionMap := map[string]interface{}{}
						if actions.RewriteAction.Action != nil {
							rewriteActionMap["action"] = actions.RewriteAction.Action
						}
						if actions.RewriteAction.Parameters != nil {
							parametersList := []interface{}{}
							for _, parameters := range actions.RewriteAction.Parameters {
								parametersMap := map[string]interface{}{}
								if parameters.Action != nil {
									parametersMap["action"] = parameters.Action
								}
								if parameters.Name != nil {
									parametersMap["name"] = parameters.Name
								}
								if parameters.Values != nil {
									parametersMap["values"] = parameters.Values
								}

								parametersList = append(parametersList, parametersMap)
							}
							rewriteActionMap["parameters"] = parametersList
						}

						actionsMap["rewrite_action"] = []interface{}{rewriteActionMap}
					}
					if actions.CodeAction != nil {
						codeActionMap := map[string]interface{}{}
						if actions.CodeAction.Action != nil {
							codeActionMap["action"] = actions.CodeAction.Action
						}
						if actions.CodeAction.Parameters != nil {
							parametersList := []interface{}{}
							for _, parameters := range actions.CodeAction.Parameters {
								parametersMap := map[string]interface{}{}
								if parameters.Name != nil {
									parametersMap["name"] = parameters.Name
								}
								if parameters.Values != nil {
									parametersMap["values"] = parameters.Values
								}
								if parameters.StatusCode != nil {
									parametersMap["status_code"] = parameters.StatusCode
								}

								parametersList = append(parametersList, parametersMap)
							}
							codeActionMap["parameters"] = parametersList
						}

						actionsMap["code_action"] = []interface{}{codeActionMap}
					}

					actionsList = append(actionsList, actionsMap)
				}
				rulesMap["actions"] = actionsList
			}

			if rules.SubRules != nil {
				subRulesList := []interface{}{}
				for _, subRules := range rules.SubRules {
					subRulesMap := map[string]interface{}{}

					if subRules.Tags != nil {
						subRulesMap["tags"] = subRules.Tags
					}
					if subRules.Rules != nil {
						rulesList := []interface{}{}
						for _, subRule := range subRules.Rules {
							rulesMap := map[string]interface{}{}
							if subRule.Conditions != nil {
								orList := []interface{}{}
								for _, or := range subRule.Conditions {
									orMap := map[string]interface{}{}
									if or.Conditions != nil {
										andList := []interface{}{}
										for _, and := range or.Conditions {
											andMap := map[string]interface{}{}
											if and.Operator != nil {
												andMap["operator"] = and.Operator
											}
											if and.Target != nil {
												andMap["target"] = and.Target
											}
											if and.Values != nil {
												andMap["values"] = and.Values
											}
											if and.IgnoreCase != nil {
												andMap["ignore_case"] = and.IgnoreCase
											}
											if and.Name != nil {
												andMap["name"] = and.Name
											}

											andList = append(andList, andMap)
										}
										orMap["and"] = andList
									}
									orList = append(orList, orMap)
								}
								rulesMap["or"] = orList
							}
							if subRule.Actions != nil {
								actionsList := []interface{}{}
								for _, actions := range subRule.Actions {
									actionsMap := map[string]interface{}{}
									if actions.NormalAction != nil {
										normalActionMap := map[string]interface{}{}
										if actions.NormalAction.Action != nil {
											normalActionMap["action"] = actions.NormalAction.Action
										}
										if actions.NormalAction.Parameters != nil {
											parametersList := []interface{}{}
											for _, parameters := range actions.NormalAction.Parameters {
												parametersMap := map[string]interface{}{}
												if parameters.Name != nil {
													parametersMap["name"] = parameters.Name
												}
												if parameters.Values != nil {
													parametersMap["values"] = parameters.Values
												}
												parametersList = append(parametersList, parametersMap)
											}
											normalActionMap["parameters"] = parametersList
										}
										actionsMap["normal_action"] = []interface{}{normalActionMap}
									}
									if actions.RewriteAction != nil {
										rewriteActionMap := map[string]interface{}{}
										if actions.RewriteAction.Action != nil {
											rewriteActionMap["action"] = actions.RewriteAction.Action
										}
										if actions.RewriteAction.Parameters != nil {
											parametersList := []interface{}{}
											for _, parameters := range actions.RewriteAction.Parameters {
												parametersMap := map[string]interface{}{}
												if parameters.Action != nil {
													parametersMap["action"] = parameters.Action
												}
												if parameters.Name != nil {
													parametersMap["name"] = parameters.Name
												}
												if parameters.Values != nil {
													parametersMap["values"] = parameters.Values
												}
												parametersList = append(parametersList, parametersMap)
											}
											rewriteActionMap["parameters"] = parametersList
										}
										actionsMap["rewrite_action"] = []interface{}{rewriteActionMap}
									}
									if actions.CodeAction != nil {
										codeActionMap := map[string]interface{}{}
										if actions.CodeAction.Action != nil {
											codeActionMap["action"] = actions.CodeAction.Action
										}
										if actions.CodeAction.Parameters != nil {
											parametersList := []interface{}{}
											for _, parameters := range actions.CodeAction.Parameters {
												parametersMap := map[string]interface{}{}
												if parameters.Name != nil {
													parametersMap["name"] = parameters.Name
												}
												if parameters.Values != nil {
													parametersMap["values"] = parameters.Values
												}
												if parameters.StatusCode != nil {
													parametersMap["status_code"] = parameters.StatusCode
												}
												parametersList = append(parametersList, parametersMap)
											}
											codeActionMap["parameters"] = parametersList
										}
										actionsMap["code_action"] = []interface{}{codeActionMap}
									}
									actionsList = append(actionsList, actionsMap)
								}
								rulesMap["actions"] = actionsList
							}
							rulesList = append(rulesList, rulesMap)
						}
						subRulesMap["rules"] = rulesList
					}
					subRulesList = append(subRulesList, subRulesMap)
				}

				rulesMap["sub_rules"] = subRulesList
			}

			rulesList = append(rulesList, rulesMap)
		}
		_ = d.Set("rules", rulesList)
	}

	return nil
}

func resourceTencentCloudTeoRuleEngineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := teo.NewModifyRuleRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	request.ZoneId = &zoneId
	request.RuleId = &ruleId

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.(*schema.Set).List()
		for i := range tagsSet {
			tags := tagsSet[i].(string)
			request.Tags = append(request.Tags, &tags)
		}
	}

	if v, ok := d.GetOk("rules"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			ruleItem := teo.Rule{}
			if v, ok := dMap["or"]; ok {
				for _, item := range v.([]interface{}) {
					ConditionsMap := item.(map[string]interface{})
					ruleAndConditions := teo.RuleAndConditions{}
					if v, ok := ConditionsMap["and"]; ok {
						for _, item := range v.([]interface{}) {
							ConditionsMap := item.(map[string]interface{})
							ruleCondition := teo.RuleCondition{}
							if v, ok := ConditionsMap["operator"]; ok {
								ruleCondition.Operator = helper.String(v.(string))
							}
							if v, ok := ConditionsMap["target"]; ok {
								ruleCondition.Target = helper.String(v.(string))
							}
							if v, ok := ConditionsMap["values"]; ok {
								valuesSet := v.(*schema.Set).List()
								for i := range valuesSet {
									values := valuesSet[i].(string)
									ruleCondition.Values = append(ruleCondition.Values, &values)
								}
							}
							if v, ok := ConditionsMap["ignore_case"]; ok {
								ruleCondition.IgnoreCase = helper.Bool(v.(bool))
							}
							if v, ok := ConditionsMap["name"]; ok {
								ruleCondition.Name = helper.String(v.(string))
							}
							ruleAndConditions.Conditions = append(ruleAndConditions.Conditions, &ruleCondition)
						}
					}
					ruleItem.Conditions = append(ruleItem.Conditions, &ruleAndConditions)
				}
			}
			if v, ok := dMap["actions"]; ok {
				for _, item := range v.([]interface{}) {
					ActionsMap := item.(map[string]interface{})
					ruleAction := teo.Action{}
					if NormalActionMap, ok := helper.InterfaceToMap(ActionsMap, "normal_action"); ok {
						ruleNormalAction := teo.NormalAction{}
						if v, ok := NormalActionMap["action"]; ok {
							ruleNormalAction.Action = helper.String(v.(string))
						}
						if v, ok := NormalActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								ParametersMap := item.(map[string]interface{})
								ruleNormalActionParams := teo.RuleNormalActionParams{}
								if v, ok := ParametersMap["name"]; ok {
									ruleNormalActionParams.Name = helper.String(v.(string))
								}
								if v, ok := ParametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleNormalActionParams.Values = append(ruleNormalActionParams.Values, &values)
									}
								}
								ruleNormalAction.Parameters = append(ruleNormalAction.Parameters, &ruleNormalActionParams)
							}
						}
						ruleAction.NormalAction = &ruleNormalAction
					}
					if RewriteActionMap, ok := helper.InterfaceToMap(ActionsMap, "rewrite_action"); ok {
						ruleRewriteAction := teo.RewriteAction{}
						if v, ok := RewriteActionMap["action"]; ok {
							ruleRewriteAction.Action = helper.String(v.(string))
						}
						if v, ok := RewriteActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								ParametersMap := item.(map[string]interface{})
								ruleRewriteActionParams := teo.RuleRewriteActionParams{}
								if v, ok := ParametersMap["action"]; ok {
									ruleRewriteActionParams.Action = helper.String(v.(string))
								}
								if v, ok := ParametersMap["name"]; ok {
									ruleRewriteActionParams.Name = helper.String(v.(string))
								}
								if v, ok := ParametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleRewriteActionParams.Values = append(ruleRewriteActionParams.Values, &values)
									}
								}
								ruleRewriteAction.Parameters = append(ruleRewriteAction.Parameters, &ruleRewriteActionParams)
							}
						}
						ruleAction.RewriteAction = &ruleRewriteAction
					}
					if CodeActionMap, ok := helper.InterfaceToMap(ActionsMap, "code_action"); ok {
						ruleCodeAction := teo.CodeAction{}
						if v, ok := CodeActionMap["action"]; ok {
							ruleCodeAction.Action = helper.String(v.(string))
						}
						if v, ok := CodeActionMap["parameters"]; ok {
							for _, item := range v.([]interface{}) {
								ParametersMap := item.(map[string]interface{})
								ruleCodeActionParams := teo.RuleCodeActionParams{}
								if v, ok := ParametersMap["name"]; ok {
									ruleCodeActionParams.Name = helper.String(v.(string))
								}
								if v, ok := ParametersMap["values"]; ok {
									valuesSet := v.(*schema.Set).List()
									for i := range valuesSet {
										values := valuesSet[i].(string)
										ruleCodeActionParams.Values = append(ruleCodeActionParams.Values, &values)
									}
								}
								if v, ok := ParametersMap["status_code"]; ok {
									ruleCodeActionParams.StatusCode = helper.IntInt64(v.(int))
								}
								ruleCodeAction.Parameters = append(ruleCodeAction.Parameters, &ruleCodeActionParams)
							}
						}
						ruleAction.CodeAction = &ruleCodeAction
					}
					ruleItem.Actions = append(ruleItem.Actions, &ruleAction)
				}
			}
			if v, ok := dMap["sub_rules"]; ok {
				for _, item := range v.([]interface{}) {
					subRulesMap := item.(map[string]interface{})
					subRuleItem := teo.SubRuleItem{}
					if v, ok := subRulesMap["tags"]; ok {
						tagsSet := v.(*schema.Set).List()
						for i := range tagsSet {
							tags := tagsSet[i].(string)
							subRuleItem.Tags = append(subRuleItem.Tags, &tags)
						}
					}
					if v, ok := subRulesMap["rules"]; ok {
						for _, item := range v.([]interface{}) {
							rulesMap := item.(map[string]interface{})
							ruleItem := teo.SubRule{}
							if v, ok := rulesMap["or"]; ok {
								for _, item := range v.([]interface{}) {
									orMap := item.(map[string]interface{})
									ruleAndConditions := teo.RuleAndConditions{}
									if v, ok := orMap["and"]; ok {
										for _, item := range v.([]interface{}) {
											andMap := item.(map[string]interface{})
											ruleCondition := teo.RuleCondition{}
											if v, ok := andMap["operator"]; ok {
												ruleCondition.Operator = helper.String(v.(string))
											}
											if v, ok := andMap["target"]; ok {
												ruleCondition.Target = helper.String(v.(string))
											}
											if v, ok := andMap["values"]; ok {
												valuesSet := v.(*schema.Set).List()
												for i := range valuesSet {
													values := valuesSet[i].(string)
													ruleCondition.Values = append(ruleCondition.Values, &values)
												}
											}
											if v, ok := andMap["ignore_case"]; ok {
												ruleCondition.IgnoreCase = helper.Bool(v.(bool))
											}
											if v, ok := andMap["name"]; ok {
												ruleCondition.Name = helper.String(v.(string))
											}
											ruleAndConditions.Conditions = append(ruleAndConditions.Conditions, &ruleCondition)
										}
									}
									ruleItem.Conditions = append(ruleItem.Conditions, &ruleAndConditions)
								}
							}
							if v, ok := rulesMap["actions"]; ok {
								for _, item := range v.([]interface{}) {
									actionsMap := item.(map[string]interface{})
									ruleAction := teo.Action{}
									if normalActionMap, ok := helper.InterfaceToMap(actionsMap, "normal_action"); ok {
										ruleNormalAction := teo.NormalAction{}
										if v, ok := normalActionMap["action"]; ok {
											ruleNormalAction.Action = helper.String(v.(string))
										}
										if v, ok := normalActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleNormalActionParams := teo.RuleNormalActionParams{}
												if v, ok := parametersMap["name"]; ok {
													ruleNormalActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleNormalActionParams.Values = append(ruleNormalActionParams.Values, &values)
													}
												}
												ruleNormalAction.Parameters = append(ruleNormalAction.Parameters, &ruleNormalActionParams)
											}
										}
										ruleAction.NormalAction = &ruleNormalAction
									}
									if rewriteActionMap, ok := helper.InterfaceToMap(actionsMap, "rewrite_action"); ok {
										ruleRewriteAction := teo.RewriteAction{}
										if v, ok := rewriteActionMap["action"]; ok {
											ruleRewriteAction.Action = helper.String(v.(string))
										}
										if v, ok := rewriteActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleRewriteActionParams := teo.RuleRewriteActionParams{}
												if v, ok := parametersMap["action"]; ok {
													ruleRewriteActionParams.Action = helper.String(v.(string))
												}
												if v, ok := parametersMap["name"]; ok {
													ruleRewriteActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleRewriteActionParams.Values = append(ruleRewriteActionParams.Values, &values)
													}
												}
												ruleRewriteAction.Parameters = append(ruleRewriteAction.Parameters, &ruleRewriteActionParams)
											}
										}
										ruleAction.RewriteAction = &ruleRewriteAction
									}
									if codeActionMap, ok := helper.InterfaceToMap(actionsMap, "code_action"); ok {
										ruleCodeAction := teo.CodeAction{}
										if v, ok := codeActionMap["action"]; ok {
											ruleCodeAction.Action = helper.String(v.(string))
										}
										if v, ok := codeActionMap["parameters"]; ok {
											for _, item := range v.([]interface{}) {
												parametersMap := item.(map[string]interface{})
												ruleCodeActionParams := teo.RuleCodeActionParams{}
												if v, ok := parametersMap["name"]; ok {
													ruleCodeActionParams.Name = helper.String(v.(string))
												}
												if v, ok := parametersMap["values"]; ok {
													valuesSet := v.(*schema.Set).List()
													for i := range valuesSet {
														values := valuesSet[i].(string)
														ruleCodeActionParams.Values = append(ruleCodeActionParams.Values, &values)
													}
												}
												if v, ok := parametersMap["status_code"]; ok {
													ruleCodeActionParams.StatusCode = helper.IntInt64(v.(int))
												}
												ruleCodeAction.Parameters = append(ruleCodeAction.Parameters, &ruleCodeActionParams)
											}
										}
										ruleAction.CodeAction = &ruleCodeAction
									}
									ruleItem.Actions = append(ruleItem.Actions, &ruleAction)
								}
							}
							subRuleItem.Rules = append(subRuleItem.Rules, &ruleItem)
						}
					}
					ruleItem.SubRules = append(ruleItem.SubRules, &subRuleItem)
				}
			}
			request.Rules = append(request.Rules, &ruleItem)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyRule(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo ruleEngine failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoRuleEngineRead(d, meta)
}

func resourceTencentCloudTeoRuleEngineDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_rule_engine.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	err := resource.Retry(5*time.Second, func() *resource.RetryError {
		if e := service.DeleteTeoRuleEngineById(ctx, zoneId, ruleId); e != nil {
			return retryError(e, "InternalError")
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo ruleEngine failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
