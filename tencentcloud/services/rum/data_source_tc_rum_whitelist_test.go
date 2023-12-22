package rum_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRumWhitelistDataSource -v
func TestAccTencentCloudRumWhitelistDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRumWhitelist,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_whitelist.whitelist"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_whitelist.whitelist", "whitelist_set.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_whitelist.whitelist", "whitelist_set.0.remark", "keep-whitelist"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_whitelist.whitelist", "whitelist_set.0.ttl", "100027012454"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_whitelist.whitelist", "whitelist_set.0.whitelist_uin", "keep-whitelist"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_whitelist.whitelist", "whitelist_set.0.wid", "11696"),
				),
			},
		},
	})
}

const testAccDataSourceRumWhitelist = `

data "tencentcloud_rum_whitelist" "whitelist" {
	instance_id = "rum-pasZKEI3RLgakj"
}

`
