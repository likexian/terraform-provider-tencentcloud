package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRumProjectDataSource -v
func TestAccTencentCloudRumProjectDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRumProject,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_project.project"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_project.project", "project_set.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.create_time", "2022-11-16 18:16:01"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.creator", "100027012454"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.desc", "Automated testing, do not delete"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.enable_url_group", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.instance_id", "rum-pasZKEI3RLgakj"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.instance_key", "pasZKEI3RLgakj"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.instance_name", "keep-rum"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.is_star", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.key", "ZEYrYfvaYQ30jRdmPx"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.name", "keep-project"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.pid", "131363"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.project_status", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.rate", "100"),
					resource.TestCheckResourceAttr("data.tencentcloud_rum_project.project", "project_set.0.type", "web"),
				),
			},
		},
	})
}

const testAccDataSourceRumProject = `

data "tencentcloud_rum_project" "project" {
	instance_id = "rum-pasZKEI3RLgakj"
}

`
