
package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMariadbDbInstancesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMariadbDbInstances,
				Check: resource.ComposeTestCheckFunc(
				  testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_db_instances.db_instances"),
				),
			},
		},
	})
}

const testAccDataSourceMariadbDbInstances = `

data "tencentcloud_mariadb_db_instances" "db_instances" {
  instance_ids = ""
  project_ids = ""
  search_name = ""
  vpc_id = ""
  subnet_id = ""
  }

`
