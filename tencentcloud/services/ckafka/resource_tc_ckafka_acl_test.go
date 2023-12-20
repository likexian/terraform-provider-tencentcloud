package ckafka_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localckafka "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ckafka"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCkafkaAclResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCkafkaAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaAcl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCkafkaAclExists("tencentcloud_ckafka_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "resource_type", "TOPIC"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "operation_type", "WRITE"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "permission_type", "ALLOW"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "host", "10.10.10.0"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "resource_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "principal"),
				),
			},
			{
				ResourceName:      "tencentcloud_ckafka_acl.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCkafkaAclExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkaService := localckafka.NewCkafkaService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ckafka acl %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ckafka acl id is not set")
		}

		_, has, err := ckafkaService.DescribeAclByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("ckafka acl doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckCkafkaAclDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ckafkaService := localckafka.NewCkafkaService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ckafka_acl" {
			continue
		}

		_, has, err := ckafkaService.DescribeAclByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("ckafka acl still exists: %s", rs.Primary.ID)
	}
	return nil
}

const testAccCkafkaAcl = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_user" "foo" {
	instance_id  = var.instance_id
	account_name = "tf-test-acl-resource"
	password     = "test1234"
  }

resource "tencentcloud_ckafka_topic" "kafka_topic_acl" {
	instance_id                     = var.instance_id
	topic_name                      = "ckafka-topic-acl-test"
	replica_num                     = 2
	partition_num                   = 1
	note                            = "test topic"
	enable_white_list               = true
	ip_white_list                   = ["192.168.1.1"]
	clean_up_policy                 = "delete"
	sync_replica_min_num            = 1
	unclean_leader_election_enable  = false
	segment                         = 86400000
	retention                       = 60000
	max_message_bytes               = 8388608
}

resource "tencentcloud_ckafka_acl" foo {
  instance_id     = var.instance_id
  resource_type   = "TOPIC"
  resource_name   = tencentcloud_ckafka_topic.kafka_topic_acl.topic_name
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "10.10.10.0"
  principal       = tencentcloud_ckafka_user.foo.account_name
}
`
