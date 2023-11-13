/*
Provides a resource to create a cls alarm

Example Usage

```hcl
resource "tencentcloud_cls_alarm" "alarm" {
  name = "alarm"
  alarm_targets {
		topic_id = "5cd3a17e-fb0b-418c-afd7-77b365397426"
		query = "* | select count(*) as count"
		number = 1
		start_time_offset = 0
		end_time_offset = 0
		logset_id = "5cd3a17e-1111-418c-afd7-77b365397426"

  }
  monitor_time {
		type = "Period"
		time = 1

  }
  condition = "$1&gt;100"
  trigger_count = 5
  alarm_period = 5
  alarm_notice_ids =
  status = true
  message_template = "test"
  call_back {
		body = "test"
		headers =

  }
  analysis {
		name = "analysis"
		type = "query"
		content = "content"
		config_info {
			key = "key"
			value = "value"
		}

  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

cls alarm can be imported using the id, e.g.

```
terraform import tencentcloud_cls_alarm.alarm alarm_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudClsAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsAlarmCreate,
		Read:   resourceTencentCloudClsAlarmRead,
		Update: resourceTencentCloudClsAlarmUpdate,
		Delete: resourceTencentCloudClsAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Log alarm name.",
			},

			"alarm_targets": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "List of alarm target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Topic id.",
						},
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Query rules.",
						},
						"number": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The number of alarm object.",
						},
						"start_time_offset": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Search start time of offset.",
						},
						"end_time_offset": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Search end time of offset.",
						},
						"logset_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Logset id.",
						},
					},
				},
			},

			"monitor_time": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Monitor task execution time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Period for periodic execution, Fixed for regular execution.",
						},
						"time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Time period or point in time.",
						},
					},
				},
			},

			"condition": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Triggering conditions.",
			},

			"trigger_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Continuous cycle.",
			},

			"alarm_period": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Alarm repeat cycle.",
			},

			"alarm_notice_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of alarm notice id.",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the alarm policy.",
			},

			"message_template": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User define alarm notice.",
			},

			"call_back": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "User define callback.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"body": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Callback body.",
						},
						"headers": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Callback headers.",
						},
					},
				},
			},

			"analysis": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Multidimensional analysis.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Analysis name.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Analysis type.",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Analysis content.",
						},
						"config_info": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudClsAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_alarm.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateAlarmRequest()
		response = cls.NewCreateAlarmResponse()
		alarmId  string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alarm_targets"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			alarmTarget := cls.AlarmTarget{}
			if v, ok := dMap["topic_id"]; ok {
				alarmTarget.TopicId = helper.String(v.(string))
			}
			if v, ok := dMap["query"]; ok {
				alarmTarget.Query = helper.String(v.(string))
			}
			if v, ok := dMap["number"]; ok {
				alarmTarget.Number = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["start_time_offset"]; ok {
				alarmTarget.StartTimeOffset = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["end_time_offset"]; ok {
				alarmTarget.EndTimeOffset = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["logset_id"]; ok {
				alarmTarget.LogsetId = helper.String(v.(string))
			}
			request.AlarmTargets = append(request.AlarmTargets, &alarmTarget)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "monitor_time"); ok {
		monitorTime := cls.MonitorTime{}
		if v, ok := dMap["type"]; ok {
			monitorTime.Type = helper.String(v.(string))
		}
		if v, ok := dMap["time"]; ok {
			monitorTime.Time = helper.IntInt64(v.(int))
		}
		request.MonitorTime = &monitorTime
	}

	if v, ok := d.GetOk("condition"); ok {
		request.Condition = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trigger_count"); ok {
		request.TriggerCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("alarm_period"); ok {
		request.AlarmPeriod = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("alarm_notice_ids"); ok {
		alarmNoticeIdsSet := v.(*schema.Set).List()
		for i := range alarmNoticeIdsSet {
			alarmNoticeIds := alarmNoticeIdsSet[i].(string)
			request.AlarmNoticeIds = append(request.AlarmNoticeIds, &alarmNoticeIds)
		}
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("message_template"); ok {
		request.MessageTemplate = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "call_back"); ok {
		callBackInfo := cls.CallBackInfo{}
		if v, ok := dMap["body"]; ok {
			callBackInfo.Body = helper.String(v.(string))
		}
		if v, ok := dMap["headers"]; ok {
			headersSet := v.(*schema.Set).List()
			for i := range headersSet {
				headers := headersSet[i].(string)
				callBackInfo.Headers = append(callBackInfo.Headers, &headers)
			}
		}
		request.CallBack = &callBackInfo
	}

	if v, ok := d.GetOk("analysis"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			analysisDimensional := cls.AnalysisDimensional{}
			if v, ok := dMap["name"]; ok {
				analysisDimensional.Name = helper.String(v.(string))
			}
			if v, ok := dMap["type"]; ok {
				analysisDimensional.Type = helper.String(v.(string))
			}
			if v, ok := dMap["content"]; ok {
				analysisDimensional.Content = helper.String(v.(string))
			}
			if v, ok := dMap["config_info"]; ok {
				for _, item := range v.([]interface{}) {
					configInfoMap := item.(map[string]interface{})
					alarmAnalysisConfig := cls.AlarmAnalysisConfig{}
					if v, ok := configInfoMap["key"]; ok {
						alarmAnalysisConfig.Key = helper.String(v.(string))
					}
					if v, ok := configInfoMap["value"]; ok {
						alarmAnalysisConfig.Value = helper.String(v.(string))
					}
					analysisDimensional.ConfigInfo = append(analysisDimensional.ConfigInfo, &alarmAnalysisConfig)
				}
			}
			request.Analysis = append(request.Analysis, &analysisDimensional)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateAlarm(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cls alarm failed, reason:%+v", logId, err)
		return err
	}

	alarmId = *response.Response.AlarmId
	d.SetId(alarmId)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::cls:%s:uin/:alarm/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsAlarmRead(d, meta)
}

func resourceTencentCloudClsAlarmRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_alarm.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	alarmId := d.Id()

	alarm, err := service.DescribeClsAlarmById(ctx, alarmId)
	if err != nil {
		return err
	}

	if alarm == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsAlarm` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if alarm.Name != nil {
		_ = d.Set("name", alarm.Name)
	}

	if alarm.AlarmTargets != nil {
		alarmTargetsList := []interface{}{}
		for _, alarmTargets := range alarm.AlarmTargets {
			alarmTargetsMap := map[string]interface{}{}

			if alarm.AlarmTargets.TopicId != nil {
				alarmTargetsMap["topic_id"] = alarm.AlarmTargets.TopicId
			}

			if alarm.AlarmTargets.Query != nil {
				alarmTargetsMap["query"] = alarm.AlarmTargets.Query
			}

			if alarm.AlarmTargets.Number != nil {
				alarmTargetsMap["number"] = alarm.AlarmTargets.Number
			}

			if alarm.AlarmTargets.StartTimeOffset != nil {
				alarmTargetsMap["start_time_offset"] = alarm.AlarmTargets.StartTimeOffset
			}

			if alarm.AlarmTargets.EndTimeOffset != nil {
				alarmTargetsMap["end_time_offset"] = alarm.AlarmTargets.EndTimeOffset
			}

			if alarm.AlarmTargets.LogsetId != nil {
				alarmTargetsMap["logset_id"] = alarm.AlarmTargets.LogsetId
			}

			alarmTargetsList = append(alarmTargetsList, alarmTargetsMap)
		}

		_ = d.Set("alarm_targets", alarmTargetsList)

	}

	if alarm.MonitorTime != nil {
		monitorTimeMap := map[string]interface{}{}

		if alarm.MonitorTime.Type != nil {
			monitorTimeMap["type"] = alarm.MonitorTime.Type
		}

		if alarm.MonitorTime.Time != nil {
			monitorTimeMap["time"] = alarm.MonitorTime.Time
		}

		_ = d.Set("monitor_time", []interface{}{monitorTimeMap})
	}

	if alarm.Condition != nil {
		_ = d.Set("condition", alarm.Condition)
	}

	if alarm.TriggerCount != nil {
		_ = d.Set("trigger_count", alarm.TriggerCount)
	}

	if alarm.AlarmPeriod != nil {
		_ = d.Set("alarm_period", alarm.AlarmPeriod)
	}

	if alarm.AlarmNoticeIds != nil {
		_ = d.Set("alarm_notice_ids", alarm.AlarmNoticeIds)
	}

	if alarm.Status != nil {
		_ = d.Set("status", alarm.Status)
	}

	if alarm.MessageTemplate != nil {
		_ = d.Set("message_template", alarm.MessageTemplate)
	}

	if alarm.CallBack != nil {
		callBackMap := map[string]interface{}{}

		if alarm.CallBack.Body != nil {
			callBackMap["body"] = alarm.CallBack.Body
		}

		if alarm.CallBack.Headers != nil {
			callBackMap["headers"] = alarm.CallBack.Headers
		}

		_ = d.Set("call_back", []interface{}{callBackMap})
	}

	if alarm.Analysis != nil {
		analysisList := []interface{}{}
		for _, analysis := range alarm.Analysis {
			analysisMap := map[string]interface{}{}

			if alarm.Analysis.Name != nil {
				analysisMap["name"] = alarm.Analysis.Name
			}

			if alarm.Analysis.Type != nil {
				analysisMap["type"] = alarm.Analysis.Type
			}

			if alarm.Analysis.Content != nil {
				analysisMap["content"] = alarm.Analysis.Content
			}

			if alarm.Analysis.ConfigInfo != nil {
				configInfoList := []interface{}{}
				for _, configInfo := range alarm.Analysis.ConfigInfo {
					configInfoMap := map[string]interface{}{}

					if configInfo.Key != nil {
						configInfoMap["key"] = configInfo.Key
					}

					if configInfo.Value != nil {
						configInfoMap["value"] = configInfo.Value
					}

					configInfoList = append(configInfoList, configInfoMap)
				}

				analysisMap["config_info"] = []interface{}{configInfoList}
			}

			analysisList = append(analysisList, analysisMap)
		}

		_ = d.Set("analysis", analysisList)

	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cls", "alarm", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudClsAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_alarm.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cls.NewModifyAlarmRequest()

	alarmId := d.Id()

	request.AlarmId = &alarmId

	immutableArgs := []string{"name", "alarm_targets", "monitor_time", "condition", "trigger_count", "alarm_period", "alarm_notice_ids", "status", "message_template", "call_back", "analysis"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("alarm_targets") {
		if v, ok := d.GetOk("alarm_targets"); ok {
			for _, item := range v.([]interface{}) {
				alarmTarget := cls.AlarmTarget{}
				if v, ok := dMap["topic_id"]; ok {
					alarmTarget.TopicId = helper.String(v.(string))
				}
				if v, ok := dMap["query"]; ok {
					alarmTarget.Query = helper.String(v.(string))
				}
				if v, ok := dMap["number"]; ok {
					alarmTarget.Number = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["start_time_offset"]; ok {
					alarmTarget.StartTimeOffset = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["end_time_offset"]; ok {
					alarmTarget.EndTimeOffset = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["logset_id"]; ok {
					alarmTarget.LogsetId = helper.String(v.(string))
				}
				request.AlarmTargets = append(request.AlarmTargets, &alarmTarget)
			}
		}
	}

	if d.HasChange("monitor_time") {
		if dMap, ok := helper.InterfacesHeadMap(d, "monitor_time"); ok {
			monitorTime := cls.MonitorTime{}
			if v, ok := dMap["type"]; ok {
				monitorTime.Type = helper.String(v.(string))
			}
			if v, ok := dMap["time"]; ok {
				monitorTime.Time = helper.IntInt64(v.(int))
			}
			request.MonitorTime = &monitorTime
		}
	}

	if d.HasChange("condition") {
		if v, ok := d.GetOk("condition"); ok {
			request.Condition = helper.String(v.(string))
		}
	}

	if d.HasChange("trigger_count") {
		if v, ok := d.GetOkExists("trigger_count"); ok {
			request.TriggerCount = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("alarm_period") {
		if v, ok := d.GetOkExists("alarm_period"); ok {
			request.AlarmPeriod = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("alarm_notice_ids") {
		if v, ok := d.GetOk("alarm_notice_ids"); ok {
			alarmNoticeIdsSet := v.(*schema.Set).List()
			for i := range alarmNoticeIdsSet {
				alarmNoticeIds := alarmNoticeIdsSet[i].(string)
				request.AlarmNoticeIds = append(request.AlarmNoticeIds, &alarmNoticeIds)
			}
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("message_template") {
		if v, ok := d.GetOk("message_template"); ok {
			request.MessageTemplate = helper.String(v.(string))
		}
	}

	if d.HasChange("call_back") {
		if dMap, ok := helper.InterfacesHeadMap(d, "call_back"); ok {
			callBackInfo := cls.CallBackInfo{}
			if v, ok := dMap["body"]; ok {
				callBackInfo.Body = helper.String(v.(string))
			}
			if v, ok := dMap["headers"]; ok {
				headersSet := v.(*schema.Set).List()
				for i := range headersSet {
					headers := headersSet[i].(string)
					callBackInfo.Headers = append(callBackInfo.Headers, &headers)
				}
			}
			request.CallBack = &callBackInfo
		}
	}

	if d.HasChange("analysis") {
		if v, ok := d.GetOk("analysis"); ok {
			for _, item := range v.([]interface{}) {
				analysisDimensional := cls.AnalysisDimensional{}
				if v, ok := dMap["name"]; ok {
					analysisDimensional.Name = helper.String(v.(string))
				}
				if v, ok := dMap["type"]; ok {
					analysisDimensional.Type = helper.String(v.(string))
				}
				if v, ok := dMap["content"]; ok {
					analysisDimensional.Content = helper.String(v.(string))
				}
				if v, ok := dMap["config_info"]; ok {
					for _, item := range v.([]interface{}) {
						configInfoMap := item.(map[string]interface{})
						alarmAnalysisConfig := cls.AlarmAnalysisConfig{}
						if v, ok := configInfoMap["key"]; ok {
							alarmAnalysisConfig.Key = helper.String(v.(string))
						}
						if v, ok := configInfoMap["value"]; ok {
							alarmAnalysisConfig.Value = helper.String(v.(string))
						}
						analysisDimensional.ConfigInfo = append(analysisDimensional.ConfigInfo, &alarmAnalysisConfig)
					}
				}
				request.Analysis = append(request.Analysis, &analysisDimensional)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyAlarm(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cls alarm failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("cls", "alarm", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsAlarmRead(d, meta)
}

func resourceTencentCloudClsAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_alarm.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	alarmId := d.Id()

	if err := service.DeleteClsAlarmById(ctx, alarmId); err != nil {
		return err
	}

	return nil
}
