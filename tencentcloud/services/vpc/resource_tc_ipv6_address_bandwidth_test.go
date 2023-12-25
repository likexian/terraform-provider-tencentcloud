package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudIpv6AddressBandwidthResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpv6AddressBandwidth,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ipv6_address_bandwidth.ipv6_address", "id")),
			},
			{
				Config: testAccIpv6AddressBandwidthUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ipv6_address_bandwidth.ipv6_address", "id")),
					resource.TestCheckResourceAttr("tencentcloud_ipv6_address_bandwidth.ipv6_address", "internet_max_bandwidth_out", "8"),
				),
			},
		},
	})
}

const testAccIpv6AddressBandwidth = `

resource "tencentcloud_ipv6_address_bandwidth" "ipv6_address" {
  ipv6_address               = "2402:4e00:1019:9400:0:9905:a90b:2ef0"
  internet_max_bandwidth_out = 6
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
}

`

const testAccIpv6AddressBandwidthUpdate = `

resource "tencentcloud_ipv6_address_bandwidth" "ipv6_address" {
  ipv6_address               = "2402:4e00:1019:9400:0:9905:a90b:2ef0"
  internet_max_bandwidth_out = 8
  internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
}

`
