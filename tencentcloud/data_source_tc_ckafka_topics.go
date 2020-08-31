/*
Use this data source to query detailed information of ckafka topic instances.

Example Usage

```hcl
resource "tencentcloud_ckafka_topic" "foo" {
	instance_id						= "ckafka-f9ife4zz"
	topic_name						= "example"
	note							= "topic note"
	replica_num						= 2
	partition_num					= 1
	enable_white_list				= 1
	ip_white_list    				= ["ip1","ip2"]
	clean_up_policy					= "delete"
	sync_replica_min_num			= 1
	unclean_leader_election_enable  = false
	segment							= 3600000
	retention						= 60000
	max_message_bytes				= 0
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCkafkaTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaTopicRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Ckafka instance ID.",
			},
			"topic_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 64),
				Description:  "Name of the CKafka topic. It must start with a letter, the rest can contain letters, numbers and dashes(-). The length range is from 1 to 64.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// computed
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of instances. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the CKafka topic.",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the CKafka topic. It must start with a letter, the rest can contain letters, numbers and dashes(-). The length range is from 1 to 64.",
						},
						"partition_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of partition.",
						},
						"replica_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of replica, the maximum is 3.",
						},
						"note": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subject note is a string of no more than 64 characters. It must start with a letter, and the remaining part can contain letters, numbers and dashes (-).",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the topic instance.",
						},
						"enable_white_list": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "IP Whitelist switch, 1: open; 0: close.",
						},
						"ip_white_list_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP Whitelist count.",
						},
						"forward_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Periodic frequency of data backup to cos.",
						},
						"forward_cos_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data backup cos bucket: the bucket address that is dumped to cos.",
						},
						"forward_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data backup cos status: 1 do not open data backup, 0 open data backup.",
						},
						"retention": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Message can be selected. Retention time, unit ms, the current minimum value is 60000ms.",
						},
						"sync_replica_min_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Min number of sync replicas, Default is 1.",
						},
						"clean_up_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Clear log policy, log clear mode, the default is delete. delete: logs are deleted according to the storage time, compact: logs are compressed according to the key, compact, delete: logs are compressed according to the key and will be deleted according to the storage time.",
						},
						"unclean_leader_election_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to allow unsynchronized replicas to be selected as leader, false: not allowed, true: allowed, not allowed by default.",
						},
						"max_message_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max message bytes.",
						},
						"segment": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Segment scrolling time, in ms, the current minimum is 3600000ms.",
						},
						"segment_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of bytes rolled by shard.",
						},
						//"config": {
						//	Type:        schema.TypeList,
						//	Computed:    true,
						//	Description: "A list of instances. Each element contains the following attributes.",
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"retention": {
						//				Type:        schema.TypeInt,
						//				Computed:    true,
						//				Description: "Message can be selected. Retention time, unit ms, the current minimum value is 60000ms.",
						//			},
						//			"min_in_sync_replicas": {
						//				Type:        schema.TypeInt,
						//				Computed:    true,
						//				Description: "Min number of sync replicas, Default is 1.",
						//			},
						//			"clean_up_policy": {
						//				Type:        schema.TypeString,
						//				Computed:    true,
						//				Description: "Clear log policy, log clear mode, the default is delete. delete: logs are deleted according to the storage time, compact: logs are compressed according to the key, compact, delete: logs are compressed according to the key and will be deleted according to the storage time.",
						//			},
						//			"unclean_leader_election_enable": {
						//				Type:        schema.TypeInt,
						//				Computed:    true,
						//				Description: "Whether to allow unsynchronized replicas to be selected as leader, false: not allowed, true: allowed, not allowed by default.",
						//			},
						//			"max_message_bytes": {
						//				Type:        schema.TypeInt,
						//				Computed:    true,
						//				Description: "Max message bytes.",
						//			},
						//			"segment_ms": {
						//				Type:        schema.TypeInt,
						//				Computed:    true,
						//				Description: "Segment scrolling time, in ms, the current minimum is 3600000ms.",
						//			},
						//			"segment_bytes": {
						//				Type:        schema.TypeInt,
						//				Computed:    true,
						//				Description: "Number of bytes rolled by shard.",
						//			},
						//		},
						//	},
						//},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCkafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ckafka_topic.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var instanceId, topicName string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("topic_name"); ok {
		topicName = v.(string)
	}
	ckafkcService := CkafkaService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	topicDetails, err := ckafkcService.DescribeCkafkaTopics(ctx, instanceId, topicName)
	if err != nil {
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(topicDetails))
	ids := make([]string, 0, len(topicDetails))

	for _, topic := range topicDetails {
		//configs := []*ckafka.Config{topic.Config}
		instance := map[string]interface{}{
			"topic_name":                     topic.TopicName,
			"topic_id":                       topic.TopicId,
			"partition_num":                  topic.PartitionNum,
			"replica_num":                    topic.ReplicaNum,
			"note":                           topic.Note,
			"create_time":                    helper.FormatUnixTime(uint64(*topic.CreateTime)),
			"enable_white_list":              topic.EnableWhiteList,
			"ip_white_list_count":            topic.IpWhiteListCount,
			"forward_interval":               topic.ForwardInterval,
			"forward_cos_bucket":             topic.ForwardCosBucket,
			"forward_status":                 topic.ForwardStatus,
			"retention":                      topic.Config.Retention,
			"sync_replica_min_num":           topic.Config.MinInsyncReplicas,
			"clean_up_policy":                topic.Config.CleanUpPolicy,
			"unclean_leader_election_enable": topic.Config.UncleanLeaderElectionEnable,
			"max_message_bytes":              topic.Config.MaxMessageBytes,
			"segment":                        topic.Config.SegmentMs,
			"segment_bytes":                  topic.Config.SegmentBytes,
			//"config":              configs,
		}
		resourceId := instanceId + FILED_SP + *topic.TopicName
		instanceList = append(instanceList, instance)
		ids = append(ids, resourceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set ckafka topic instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceList); err != nil {
			return err
		}
	}

	return nil
}
