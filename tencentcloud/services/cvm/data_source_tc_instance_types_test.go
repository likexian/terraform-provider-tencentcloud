package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudInstanceTypesDataSource_basic -v
func TestAccTencentCloudInstanceTypesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceTypesDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.availability_zone", "ap-guangzhou-3"),
				),
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudInstanceTypesDataSource_sell -v
func TestAccTencentCloudInstanceTypesDataSource_sell(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceTypesDataSourceConfigSell,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.cpu_core_count", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.memory_size", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.example", "instance_types.0.family", "SA2"),
				),
			},
		},
	})
}

const testAccTencentCloudInstanceTypesDataSourceConfigBasic = `
data "tencentcloud_instance_types" "example" {
  availability_zone = "ap-guangzhou-3"
  cpu_core_count = 4
  memory_size    = 8
}
`

const testAccTencentCloudInstanceTypesDataSourceConfigSell = `
data "tencentcloud_instance_types" "example" {
  cpu_core_count = 2
  memory_size    = 2
  exclude_sold_out = true

  filter{
	name = "instance-family"
    values = ["SA2"]
  }

  filter{
	name = "zone"
    values = ["ap-guangzhou-3"]
  }
}
`
