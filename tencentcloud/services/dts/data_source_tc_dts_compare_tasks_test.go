package dts_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDtsCompareTasksDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDtsCompareTasks, tcacctest.DefaultDTSJobId, tcacctest.DefaultDTSJobId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dts_compare_tasks.compare_tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dts_compare_tasks.compare_tasks", "list.#"),
				),
			},
		},
	})
}

const testAccDataSourceDtsCompareTasks = `
resource "tencentcloud_dts_compare_task" "task" {
	job_id = "%s"
	task_name = "tf_test_compare_task"
	objects {
	  object_mode = "all"
	}
  }

data "tencentcloud_dts_compare_tasks" "compare_tasks" {
  job_id = "%s"
  }

`
