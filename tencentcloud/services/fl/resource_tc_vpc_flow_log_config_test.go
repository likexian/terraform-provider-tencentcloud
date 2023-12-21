package fl_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcFlowLogConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcFlowLogConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_flow_log_config.flow_log_config", "id")),
			},
			{
				Config: testAccVpcFlowLogConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_flow_log_config.flow_log_config", "enable", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_vpc_flow_log_config.flow_log_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcFlowLogConfig = `

resource "tencentcloud_vpc_flow_log_config" "flow_log_config" {
  flow_log_id = "fl-geg2keoj"
  enable = false
}

`

const testAccVpcFlowLogConfigUpdate = `

resource "tencentcloud_vpc_flow_log_config" "flow_log_config" {
  flow_log_id = "fl-geg2keoj"
  enable = true
}

`
