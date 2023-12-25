package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcNetDetectStateCheckDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcNetDetectStateCheckDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_net_detect_state_check.net_detect_state_check")),
			},
		},
	})
}

const testAccVpcNetDetectStateCheckDataSource = `

data "tencentcloud_vpc_net_detect_state_check" "net_detect_state_check" {
  net_detect_id         = "netd-12345678"
  detect_destination_ip = [
    "10.0.0.3",
    "10.0.0.2"
  ]
  next_hop_type        = "NORMAL_CVM"
  next_hop_destination = "10.0.0.4"
}

`
