package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTdmqMessageDataSource_basic -v
func TestAccTencentCloudNeedFixTdmqMessageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqMessageDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_message.message"),
				),
			},
		},
	})
}

const testAccTdmqMessageDataSource = `
data "tencentcloud_tdmq_message" "message" {
  cluster_id     = "rocketmq-rkrbm52djmro"
  environment_id = "keep_ns"
  topic_name     = "keep-topic"
  msg_id         = "A9FE8D0567FE15DB97425FC08EEF0000"
  query_dlq_msg  = false
}
`
