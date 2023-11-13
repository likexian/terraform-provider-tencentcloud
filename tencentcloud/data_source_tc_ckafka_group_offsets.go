/*
Use this data source to query detailed information of ckafka group_offsets

Example Usage

```hcl
data "tencentcloud_ckafka_group_offsets" "group_offsets" {
  instance_id = "InstanceId"
  group = "groupName"
  topics =
  search_word = "topicName"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCkafkaGroupOffsets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaGroupOffsetsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"group": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Kafka consumer group name.",
			},

			"topics": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "An array of topic names subscribed by the group, if there is no such array, it means all topic information under the specified group.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy match topicName.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of matching results.",
						},
						"topic_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The topic array, where each element is a json object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TopicName.",
									},
									"partitions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "He topic partition array, where each element is a json object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"partition": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Topic partitionId.",
												},
												"offset": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The offset of the position.",
												},
												"metadata": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When consumers submit messages, they can pass in metadata for other purposes. Currently, it is usually an empty string.",
												},
												"error_code": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ErrorCode.",
												},
												"log_end_offset": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The latest offset of the current partition.",
												},
												"lag": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The number of unconsumed messages.",
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCkafkaGroupOffsetsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_group_offsets.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group"); ok {
		paramMap["Group"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topics"); ok {
		topicsSet := v.(*schema.Set).List()
		paramMap["Topics"] = helper.InterfacesStringsPoint(topicsSet)
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.GroupOffsetResponse

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaGroupOffsetsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		groupOffsetResponseMap := map[string]interface{}{}

		if result.TotalCount != nil {
			groupOffsetResponseMap["total_count"] = result.TotalCount
		}

		if result.TopicList != nil {
			topicListList := []interface{}{}
			for _, topicList := range result.TopicList {
				topicListMap := map[string]interface{}{}

				if topicList.Topic != nil {
					topicListMap["topic"] = topicList.Topic
				}

				if topicList.Partitions != nil {
					partitionsList := []interface{}{}
					for _, partitions := range topicList.Partitions {
						partitionsMap := map[string]interface{}{}

						if partitions.Partition != nil {
							partitionsMap["partition"] = partitions.Partition
						}

						if partitions.Offset != nil {
							partitionsMap["offset"] = partitions.Offset
						}

						if partitions.Metadata != nil {
							partitionsMap["metadata"] = partitions.Metadata
						}

						if partitions.ErrorCode != nil {
							partitionsMap["error_code"] = partitions.ErrorCode
						}

						if partitions.LogEndOffset != nil {
							partitionsMap["log_end_offset"] = partitions.LogEndOffset
						}

						if partitions.Lag != nil {
							partitionsMap["lag"] = partitions.Lag
						}

						partitionsList = append(partitionsList, partitionsMap)
					}

					topicListMap["partitions"] = []interface{}{partitionsList}
				}

				topicListList = append(topicListList, topicListMap)
			}

			groupOffsetResponseMap["topic_list"] = []interface{}{topicListList}
		}

		ids = append(ids, *result.InstanceId)
		_ = d.Set("result", groupOffsetResponseMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), groupOffsetResponseMap); e != nil {
			return e
		}
	}
	return nil
}
