package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudTdmqSubscribeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckTdmqSubscribeDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqSubscribe,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTdmqSubscribeExists("tencentcloud_tcmq_subscribe.subscribe"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcmq_subscribe.subscribe", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tcmq_subscribe.subscribe", "topic_name", "test_subscribe_topic"),
					resource.TestCheckResourceAttr("tencentcloud_tcmq_subscribe.subscribe", "subscription_name", "test_subscribe"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcmq_subscribe.subscribe",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTdmqSubscribeDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcmq_subscribe" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		topicName := idSplit[0]
		subscriptionName := idSplit[1]
		subscribe, err := service.DescribeTcmqSubscribeById(ctx, topicName, subscriptionName)
		if subscribe != nil {
			return fmt.Errorf("TdmqSubscribe instance still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTdmqSubscribeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TdmqSubscribe %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TdmqSubscribe id is not set")
		}

		service := TcmqService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		topicName := idSplit[0]
		subscriptionName := idSplit[1]
		subscribe, err := service.DescribeTcmqSubscribeById(ctx, topicName, subscriptionName)
		if subscribe == nil {
			return fmt.Errorf("TdmqSubscribe %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTdmqSubscribe = `
resource "tencentcloud_tcmq_topic" "topic" {
	topic_name = "test_subscribe_topic"
}

resource "tencentcloud_tcmq_subscribe" "subscribe" {
	topic_name = tencentcloud_tcmq_topic.topic.topic_name
	subscription_name = "test_subscribe"
	protocol = "http"
	endpoint = "http://mikatong.com"
}
`
