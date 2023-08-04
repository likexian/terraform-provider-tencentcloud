package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrDescribeWebhookTriggerLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrDescribeWebhookTriggerLogsDataSource,
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_webhook_trigger_logs.my_logs"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "registry_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "trigger_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.trigger_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.event_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.notify_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.detail"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.update_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_webhook_trigger_logs.my_logs", "logs.0.status"),
				),
			},
		},
	})
}

const testAccTcrDescribeWebhookTriggerLogsDataSource = TCRDataSource + `

data "tencentcloud_tcr_webhook_trigger_logs" "my_logs" {
  registry_id = local.tcr_id
  namespace = local.tcr_ns_name
  trigger_id = 1
    tags = {
    "createdBy" = "terraform"
  }
}

`
