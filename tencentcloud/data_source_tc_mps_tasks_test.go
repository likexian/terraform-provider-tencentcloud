package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsTasksDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsTasksDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_tasks.tasks"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_tasks.tasks", "task_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_tasks.tasks", "task_set.0.task_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_tasks.tasks", "task_set.0.task_type"),
				),
			},
		},
	})
}

const testAccMpsTasksDataSource = `

data "tencentcloud_mps_tasks" "tasks" {
  status = "FINISH"
  limit  = 20
}

`
