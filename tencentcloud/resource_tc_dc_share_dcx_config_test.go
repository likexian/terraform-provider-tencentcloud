package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcShareDcxConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcShareDcxConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dc_share_dcx_config.share_dcx_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_dc_share_dcx_config.share_dcx_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcShareDcxConfig = `

resource "tencentcloud_dc_share_dcx_config" "share_dcx_config" {
  direct_connect_tunnel_id = "dcx-test1234"
  enable = true
}

`
