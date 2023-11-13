/*
Use this data source to query detailed information of ckafka topic_flow_ranking

Example Usage

```hcl
data "tencentcloud_ckafka_topic_flow_ranking" "topic_flow_ranking" {
  instance_id = "InstanceId"
  ranking_type = "PRO"
  begin_date = "2021-05-13T07:23:11+08:00"
  end_date = "2021-05-14T07:23:11+08:00"
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

func dataSourceTencentCloudCkafkaTopicFlowRanking() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaTopicFlowRankingRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"ranking_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "RankingType.",
			},

			"begin_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "BeginDate.",
			},

			"end_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "EndDate.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_flow": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TopicFlow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TopicId.",
									},
									"topic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TopicName.",
									},
									"partition_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "PartitionNum.",
									},
									"replica_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ReplicaNum.",
									},
									"topic_traffic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TopicTraffic.",
									},
									"message_heap": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Topic MessageHeap.",
									},
								},
							},
						},
						"consume_speed": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "ConsumeSpeed.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"consumer_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ConsumerGroupName.",
									},
									"speed": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Speed.",
									},
								},
							},
						},
						"topic_message_heap": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TopicMessageHeapRanking.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TopicId.",
									},
									"topic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TopicName.",
									},
									"partition_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "PartitionNum.",
									},
									"replica_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ReplicaNum.",
									},
									"topic_traffic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TopicTraffic.",
									},
									"message_heap": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Topic MessageHeap.",
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

func dataSourceTencentCloudCkafkaTopicFlowRankingRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_topic_flow_ranking.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ranking_type"); ok {
		paramMap["RankingType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("begin_date"); ok {
		paramMap["BeginDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*ckafka.TopicFlowRankingResult

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaTopicFlowRankingByFilter(ctx, paramMap)
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
		topicFlowRankingResultMap := map[string]interface{}{}

		if result.TopicFlow != nil {
			topicFlowList := []interface{}{}
			for _, topicFlow := range result.TopicFlow {
				topicFlowMap := map[string]interface{}{}

				if topicFlow.TopicId != nil {
					topicFlowMap["topic_id"] = topicFlow.TopicId
				}

				if topicFlow.TopicName != nil {
					topicFlowMap["topic_name"] = topicFlow.TopicName
				}

				if topicFlow.PartitionNum != nil {
					topicFlowMap["partition_num"] = topicFlow.PartitionNum
				}

				if topicFlow.ReplicaNum != nil {
					topicFlowMap["replica_num"] = topicFlow.ReplicaNum
				}

				if topicFlow.TopicTraffic != nil {
					topicFlowMap["topic_traffic"] = topicFlow.TopicTraffic
				}

				if topicFlow.MessageHeap != nil {
					topicFlowMap["message_heap"] = topicFlow.MessageHeap
				}

				topicFlowList = append(topicFlowList, topicFlowMap)
			}

			topicFlowRankingResultMap["topic_flow"] = []interface{}{topicFlowList}
		}

		if result.ConsumeSpeed != nil {
			consumeSpeedList := []interface{}{}
			for _, consumeSpeed := range result.ConsumeSpeed {
				consumeSpeedMap := map[string]interface{}{}

				if consumeSpeed.ConsumerGroupName != nil {
					consumeSpeedMap["consumer_group_name"] = consumeSpeed.ConsumerGroupName
				}

				if consumeSpeed.Speed != nil {
					consumeSpeedMap["speed"] = consumeSpeed.Speed
				}

				consumeSpeedList = append(consumeSpeedList, consumeSpeedMap)
			}

			topicFlowRankingResultMap["consume_speed"] = []interface{}{consumeSpeedList}
		}

		if result.TopicMessageHeap != nil {
			topicMessageHeapList := []interface{}{}
			for _, topicMessageHeap := range result.TopicMessageHeap {
				topicMessageHeapMap := map[string]interface{}{}

				if topicMessageHeap.TopicId != nil {
					topicMessageHeapMap["topic_id"] = topicMessageHeap.TopicId
				}

				if topicMessageHeap.TopicName != nil {
					topicMessageHeapMap["topic_name"] = topicMessageHeap.TopicName
				}

				if topicMessageHeap.PartitionNum != nil {
					topicMessageHeapMap["partition_num"] = topicMessageHeap.PartitionNum
				}

				if topicMessageHeap.ReplicaNum != nil {
					topicMessageHeapMap["replica_num"] = topicMessageHeap.ReplicaNum
				}

				if topicMessageHeap.TopicTraffic != nil {
					topicMessageHeapMap["topic_traffic"] = topicMessageHeap.TopicTraffic
				}

				if topicMessageHeap.MessageHeap != nil {
					topicMessageHeapMap["message_heap"] = topicMessageHeap.MessageHeap
				}

				topicMessageHeapList = append(topicMessageHeapList, topicMessageHeapMap)
			}

			topicFlowRankingResultMap["topic_message_heap"] = []interface{}{topicMessageHeapList}
		}

		ids = append(ids, *result.InstanceId)
		_ = d.Set("result", topicFlowRankingResultMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), topicFlowRankingResultMap); e != nil {
			return e
		}
	}
	return nil
}
