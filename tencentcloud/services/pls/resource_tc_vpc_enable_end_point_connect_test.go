package pls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcEnableEndPointConnectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEnableEndPointConnect,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_enable_end_point_connect.enable_end_point_connect", "id")),
			},
		},
	})
}

const testAccVpcEnableEndPointConnect = `

resource "tencentcloud_vpc_enable_end_point_connect" "enable_end_point_connect" {
  end_point_service_id = "vpcsvc-98jddhcz"
  end_point_id         = ["vpce-6q0ftmke"]
  accept_flag          = true
}

`
