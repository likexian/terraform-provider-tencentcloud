/*
Provides a resource to create a tke tmpRecordRule

Example Usage

```hcl

resource "tencentcloud_tke_tmp_record_rule_yaml" "foo" {
  instance_id       = ""
  content           = ""   # yaml format
}

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTkeTmpRecordRuleYaml() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTkeTmpRecordRuleYamlRead,
		Create: resourceTencentCloudTkeTmpRecordRuleYamlCreate,
		Update: resourceTencentCloudTkeTmpRecordRuleYamlUpdate,
		Delete: resourceTencentCloudTkeTmpRecordRuleYamlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance Id.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Contents of record rules in yaml format.",
			},

			"total": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "count of returned lists.",
			},

			"record_rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of record rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the instance.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time of record rule.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Used for the argument, if the configuration comes to the template, the template id.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Contents of record rules in yaml format.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An ID identify the cluster, like cls-xxxxxx.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTkeTmpRecordRuleYamlCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_record_rule_yaml.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tke.NewCreatePrometheusRecordRuleYamlRequest()

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	tmpRecordRuleName := ""
	if v, ok := d.GetOk("content"); ok {
		if m, err := YamlParser(v.(string)); err != nil {
			log.Printf("[CRITAL]%s check yaml syntax failed, error:%+v", logId, err)
			return err
		} else {
			metadata := m["metadata"]
			if metadata != nil {
				if metadata.(map[interface{}]interface{})["name"] != nil {
					tmpRecordRuleName = metadata.(map[interface{}]interface{})["name"].(string)
				}
			}
		}

		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().CreatePrometheusRecordRuleYaml(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tke tmpRecordRule failed, reason:%+v", logId, err)
		return err
	}

	instanceId := *request.InstanceId
	d.SetId(strings.Join([]string{instanceId, tmpRecordRuleName}, FILED_SP))
	return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
}

func resourceTencentCloudTkeTmpRecordRuleYamlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_record_rule_yaml.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	name := ids[1]

	recordRuleService := RecordRuleService{client: meta.(*TencentCloudClient).apiV3Conn}
	request, err := recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, name)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			request, err = recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, name)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	recordRules := request.Response.Records
	if len(recordRules) == 0 {
		d.SetId("")
		return nil
	}

	records := make([]map[string]interface{}, 0, len(recordRules))
	for _, recordRule := range recordRules {
		var infoMap = map[string]interface{}{}
		infoMap["name"] = *recordRule.Name
		infoMap["update_time"] = *recordRule.UpdateTime
		infoMap["template_id"] = *recordRule.TemplateId
		infoMap["content"] = *recordRule.Content
		infoMap["cluster_id"] = *recordRule.ClusterId
		records = append(records, infoMap)
	}

	_ = d.Set("prometheus_record_rule_yaml_items", records)
	_ = d.Set("total", request.Response.Total)

	return nil
}

func resourceTencentCloudTkeTmpRecordRuleYamlUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_record_rule_yaml.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tke.NewModifyPrometheusRecordRuleYamlRequest()

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	request.InstanceId = &ids[0]
	request.Name = &ids[1]

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("name") {
		return fmt.Errorf("`name` do not support change now.")
	}

	//if d.HasChange("content") {
	//	if v, ok := d.GetOk("content"); ok {
	//		request.Content = helper.String(v.(string))
	//
	//		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
	//			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().ModifyPrometheusRecordRuleYaml(request)
	//			if e != nil {
	//				return retryError(e)
	//			} else {
	//				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
	//					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
	//			}
	//			return nil
	//		})
	//
	//		if err != nil {
	//			return err
	//		}
	//
	//		return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
	//	}
	//}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().ModifyPrometheusRecordRuleYaml(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			return err
		}

		return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
	}

	return nil
}

func resourceTencentCloudTkeTmpRecordRuleYamlDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_record_rule_yaml.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	service := RecordRuleService{client: meta.(*TencentCloudClient).apiV3Conn}
	if err := service.DeletePrometheusRecordRuleYaml(ctx, ids[0], ids[1]); err != nil {
		return err
	}

	return nil
}
