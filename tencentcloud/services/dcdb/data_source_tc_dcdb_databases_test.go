package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDCDBDatabasesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		// CheckDestroy: testAccCheckDcdbAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbDatabases, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_databases.databases"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_databases.databases", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_databases.databases", "list.0.db_name", "mysql"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_databases.databases", "list.1.db_name", "performance_schema"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_databases.databases", "list.2.db_name", "query_rewrite"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_databases.databases", "list.3.db_name", "sys"),
				),
			},
		},
	})
}

const testAccDataSourceDcdbDatabases = `
data "tencentcloud_dcdb_databases" "databases" {
	instance_id = "%s" # use the hard code before the dcdb_instance resource is ready.
}

`
