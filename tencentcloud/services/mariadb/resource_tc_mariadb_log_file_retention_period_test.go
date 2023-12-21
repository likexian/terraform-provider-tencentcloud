package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbLogFileRetentionPeriod_basic -v
func TestAccTencentCloudMariadbLogFileRetentionPeriod_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbLogFileRetentionPeriod,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_log_file_retention_period.logFileRetentionPeriod", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mariadb_log_file_retention_period.logFileRetentionPeriod", "days", "8"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_log_file_retention_period.logFileRetentionPeriod",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbLogFileRetentionPeriod = testAccMariadbHourDbInstance + `

resource "tencentcloud_mariadb_log_file_retention_period" "logFileRetentionPeriod" {
  instance_id = tencentcloud_mariadb_hour_db_instance.basic.id
  days = "8"
}

`
