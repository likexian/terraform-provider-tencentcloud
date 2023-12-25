package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudHaVipsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudHaVipsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ha_vips.havips"),
					resource.TestCheckResourceAttr("data.tencentcloud_ha_vips.havips", "ha_vip_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_ha_vips.havips", "ha_vip_list.0.name", "terraform_test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ha_vips.havips", "ha_vip_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ha_vips.havips", "ha_vip_list.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ha_vips.havips", "ha_vip_list.0.vip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ha_vips.havips", "ha_vip_list.0.state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ha_vips.havips", "ha_vip_list.0.create_time"),
				),
			},
		},
	})
}

const testAccTencentCloudHaVipsDataSourceConfig_basic = tcacctest.DefaultVpcVariable + `
resource "tencentcloud_ha_vip" "havip" {
  name      = "terraform_test"
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}

data "tencentcloud_ha_vips" "havips" {
  id = tencentcloud_ha_vip.havip.id
}
`
