package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixTdmqProInstancesDataSource_basic -v
func TestAccTencentCloudNeedFixTdmqProInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqProInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_pro_instances.pro_instances"),
				),
			},
			{
				Config: testAccTdmqProInstancesDataSourcelFilter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_pro_instances.pro_instances"),
				),
			},
		},
	})
}

const testAccTdmqProInstancesDataSource = `
data "tencentcloud_tdmq_pro_instances" "pro_instances" {
}
`

const testAccTdmqProInstancesDataSourcelFilter = `
data "tencentcloud_tdmq_pro_instances" "pro_instances" {
  filters {
    name   = ""
    values = ""
  }
}
`
