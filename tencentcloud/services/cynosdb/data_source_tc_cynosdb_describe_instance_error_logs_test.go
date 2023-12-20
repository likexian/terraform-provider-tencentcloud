package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbDescribeInstanceErrorLogsDataSource_basic -v
func TestAccTencentCloudCynosdbDescribeInstanceErrorLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbDescribeInstanceErrorLogsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_describe_instance_error_logs.describe_instance_error_logs"),
				),
			},
		},
	})
}

const testAccCynosdbDescribeInstanceErrorLogsDataSource = `
data "tencentcloud_cynosdb_describe_instance_error_logs" "describe_instance_error_logs" {
  instance_id   = "cynosdbmysql-ins-afqx1hy0"
  start_time    = "2023-06-01 15:04:05"
  end_time      = "2023-06-19 15:04:05"
  order_by      = "Timestamp"
  order_by_type = "DESC"
  log_levels    = ["note", "warning"]
  key_words     = ["Aborted"]
}
`
