package ccn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCcnRoutesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCcnRoutes,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ccn_routes.ccn_routes", "id")),
			},
			{
				Config: testAccVpcCcnRoutesUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ccn_routes.ccn_routes", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ccn_routes.ccn_routes", "switch", "on"),
				),
			},
			{
				ResourceName:      "tencentcloud_ccn_routes.ccn_routes",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcCcnRoutes = `

resource "tencentcloud_ccn_routes" "ccn_routes" {
  ccn_id = "ccn-0bbkedsb"
  route_id = "ccnr-9sqye2qg"
  switch = "off"
}

`

const testAccVpcCcnRoutesUpdate = `

resource "tencentcloud_ccn_routes" "ccn_routes" {
  ccn_id = "ccn-0bbkedsb"
  route_id = "ccnr-9sqye2qg"
  switch = "on"
}

`
