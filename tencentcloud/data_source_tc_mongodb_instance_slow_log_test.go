package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceSlowLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceSlowLogDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_slow_log.instance_slow_log")),
			},
		},
	})
}

const testAccMongodbInstanceSlowLogDataSource = `

data "tencentcloud_mongodb_instance_slow_log" "instance_slow_log" {
  instance_id = "cmgo-gwqk8669"
  start_time = "2013-05-10 10:00:00"
  end_time = "2013-05-10 12:47:00"
  slow_ms = 100
  format = "json"
}

`
