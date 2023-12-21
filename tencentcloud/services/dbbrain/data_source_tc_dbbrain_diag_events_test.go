package dbbrain_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainDiagEventsDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02T15:04:05+08:00")
	endTime := time.Now().In(loc).Format("2006-01-02T15:04:05+08:00")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainDiagEventsDataSource, tcacctest.DefaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_diag_events.diag_events"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.diag_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.event_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.severity"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.outline"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.diag_item"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_events.diag_events", "list.0.region"),
				),
			},
		},
	})
}

const testAccDbbrainDiagEventsDataSource = `

data "tencentcloud_dbbrain_diag_events" "diag_events" {
  instance_ids = ["%s"]
  start_time = "%s"
  end_time = "%s"
  severities = [1,4,5]
}

`
