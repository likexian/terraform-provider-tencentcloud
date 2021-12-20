package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudKafkaTopic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccTencentCloudKafkaTopicDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaTopicInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaTopicInstanceExists("tencentcloud_ckafka_topic.kafka_topic"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "instance_id", "ckafka-f9ife4zz"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "topic_name", "ckafka-topic-tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "note", "this is test ckafka topic"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "replica_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "partition_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "enable_white_list", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "ip_white_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "ip_white_list.0", "192.168.1.1"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "clean_up_policy", "delete"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "sync_replica_min_num", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_topic.kafka_topic", "unclean_leader_election_enable"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "segment", "3600000"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "retention", "60000"),
				),
			},
			{
				Config: testAccKafkaTopicInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaTopicInstanceExists("tencentcloud_ckafka_topic.kafka_topic"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "instance_id", "ckafka-f9ife4zz"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "note", "this is test topic_update"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "partition_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "enable_white_list", "true"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "clean_up_policy", "compact"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "sync_replica_min_num", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_topic.kafka_topic", "unclean_leader_election_enable"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "segment", "4000000"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "retention", "70000"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_topic.kafka_topic", "max_message_bytes", "8388608"),
				),
			},
			{
				ResourceName:      "tencentcloud_ckafka_topic.kafka_topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTencentCloudKafkaTopicDestory(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkcService := CkafkaService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, r := range s.RootModule().Resources {
		if r.Type != "tencentcloud_ckafka_topic" {
			continue
		}
		split := strings.Split(r.Primary.ID, FILED_SP)
		if len(split) < 2 {
			continue
		}
		_, has, error := ckafkcService.DescribeCkafkaTopicByName(ctx, split[0], split[1])
		if error != nil {
			return error
		}
		if !has {
			return nil
		}
		return fmt.Errorf("ckafka topic still exists: %s", r.Primary.ID)
	}
	return nil
}

func testAccCheckKafkaTopicInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ckafka topic %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ckafka topic id is not set")
		}
		ckafkcService := CkafkaService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		split := strings.Split(rs.Primary.ID, FILED_SP)
		if len(split) < 2 {
			return fmt.Errorf("ckafka topic is not set: %s", rs.Primary.ID)
		}
		var exist bool
		outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, inErr := ckafkcService.DescribeCkafkaTopicByName(ctx, split[0], split[1])
			if inErr != nil {
				return retryError(inErr)
			}
			exist = has
			return nil
		})
		if outErr != nil {
			return outErr
		}
		if !exist {
			return fmt.Errorf("ckafka topic doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccKafkaTopicInstance = `
resource "tencentcloud_ckafka_topic" "kafka_topic" {
	instance_id                         = "ckafka-f9ife4zz"
	topic_name                          = "ckafka-topic-tf-test"
	note                                = "this is test ckafka topic"
	replica_num                         = 2
	partition_num                       = 2
	enable_white_list                   = true
	ip_white_list                       = ["192.168.1.1"]
	clean_up_policy                     = "delete"
	sync_replica_min_num                = 1
	unclean_leader_election_enable      = false
	segment                             = 3600000
	retention                           = 60000
	max_message_bytes                   = 0
}
`

const testAccKafkaTopicInstanceUpdate = `
resource "tencentcloud_ckafka_topic" "kafka_topic" {
	instance_id                         = "ckafka-f9ife4zz"
	topic_name                          = "ckafka-topic-tf-test"
	note                                = "this is test topic_update"
	replica_num                         = 2
	partition_num                       = 2
	enable_white_list                   = true
	ip_white_list                       = ["192.168.1.2"]
	clean_up_policy                     = "compact"
	sync_replica_min_num                = 2
	unclean_leader_election_enable      = true
	segment                             = 4000000
	retention                           = 70000
	max_message_bytes                   = 8388608
}
`
