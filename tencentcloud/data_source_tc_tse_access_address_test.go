package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTseAccessAddressDataSource_basic -v
func TestAccTencentCloudTseAccessAddressDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseAccessAddressDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tse_access_address.access_address"),
					resource.TestCheckResourceAttr("data.tencentcloud_tse_access_address.access_address", "engine_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tse_access_address.access_address", "env_address_infos.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_tse_access_address.access_address", "env_address_infos.0.enable_config_internet", "false"),
					resource.TestCheckResourceAttr("data.tencentcloud_tse_access_address.access_address", "env_address_infos.0.enable_config_intranet", "false"),
				),
			},
		},
	})
}

const testAccTseAccessAddressDataSource = testAccTseInstance + `

data "tencentcloud_tse_access_address" "access_address" {
  instance_id = tencentcloud_tse_instance.instance.id
  engine_region = "ap-guangzhou"
}

`
