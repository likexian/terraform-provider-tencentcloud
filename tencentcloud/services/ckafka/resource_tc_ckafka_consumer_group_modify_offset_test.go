package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaConsumerGroupModifyOffsetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaConsumerGroupModifyOffset,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_consumer_group_modify_offset.consumer_group_modify_offset", "id")),
			},
		},
	})
}

const testAccCkafkaConsumerGroupModifyOffset = `
resource "tencentcloud_ckafka_consumer_group_modify_offset" "consumer_group_modify_offset" {
	instance_id = "ckafka-vv7wpvae"
	group = "keep-group"
	offset = 0
	strategy = 2
	topics = ["keep-topic"]
}
`
