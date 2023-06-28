package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverDatabaseTDEResource_basic -v
func TestAccTencentCloudSqlserverDatabaseTDEResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatabaseTDE,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_database_tde.database_tde", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_database_tde.database_tde",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverDatabaseTDE = `
resource "tencentcloud_sqlserver_database_tde" "database_tde" {
  instance_id = "mssql-qelbzgwf"
  db_tde_encrypt {
    db_name    = "keep_tde_db"
    encryption = "enable"
  }
  db_tde_encrypt {
    db_name    = "keep_tde_db2"
    encryption = "disable"
  }
}
`
