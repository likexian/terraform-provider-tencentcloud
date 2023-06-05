package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbActivateHourDbInstanceResource_basic -v
func TestAccTencentCloudMariadbActivateHourDbInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbActivateHourDbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_activate_hour_db_instance.activate_hour_db_instance", "id"),
				),
			},
		},
	})
}

const testAccMariadbActivateHourDbInstance = `
resource "tencentcloud_mariadb_activate_hour_db_instance" "activate_hour_db_instance" {
  instance_id = "tdsql-9vqvls95"
}
`
