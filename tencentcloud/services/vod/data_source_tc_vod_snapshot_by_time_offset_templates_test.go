package vod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudVodSnapshotByTimeOffsetTemplates(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSnapshotByTimeOffsetTemplates,

				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.name", "tf-snapshot"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.width", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.height", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.resolution_adaptive", "false"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.format", "png"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.comment", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.fill_type", "white"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_snapshot_by_time_offset_templates.foo", "template_list.0.update_time"),
				),
			},
		},
	})
}

const testAccVodSnapshotByTimeOffsetTemplates = testAccVodSnapshotByTimeOffsetTemplate + `
data "tencentcloud_vod_snapshot_by_time_offset_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
}
`
