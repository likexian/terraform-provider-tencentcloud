package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDCDBDatabaseObjectsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDCDBDatabaseObjectsDataSource, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_database_objects.database_objects"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_database_objects.database_objects", "tables.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_database_objects.database_objects", "views.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_database_objects.database_objects", "procs.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_database_objects.database_objects", "funcs.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_database_objects.database_objects", "db_name", "tf_test_db"),
				),
			},
		},
	})
}

const testAccDCDBDatabaseObjectsDataSource = `

data "tencentcloud_dcdb_database_objects" "database_objects" {
	instance_id = "%s"
	db_name = "tf_test_db"
}

`
