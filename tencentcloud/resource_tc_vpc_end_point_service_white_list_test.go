package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestTencentCloudVpcEndPointServiceWhiteListResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPointServiceWhiteList,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point_service_white_list.end_point_service_white_list", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_end_point_service_white_list.end_point_service_white_list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcEndPointServiceWhiteList = `

resource "tencentcloud_vpc_end_point_service_white_list" "end_point_service_white_list" {
  user_uin = "100020512675"
  end_point_service_id = "vpcsvc-69y13tdb"
  description = "terraform for test"
}

`
