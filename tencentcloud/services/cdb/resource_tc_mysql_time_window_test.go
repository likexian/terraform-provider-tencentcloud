package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlTimeWindowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlTimeWindow,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_time_window.time_window", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_time_window.time_window",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlTimeWindow = `

resource "tencentcloud_mysql_time_window" "time_window" {
  instance_id    = "cdb-fitq5t9h"
  max_delay_time = 10
  time_ranges    = [
    "01:00-02:01"
  ]
  weekdays       = [
    "friday",
    "monday",
    "saturday",
    "thursday",
    "tuesday",
    "wednesday",
  ]
}

`
