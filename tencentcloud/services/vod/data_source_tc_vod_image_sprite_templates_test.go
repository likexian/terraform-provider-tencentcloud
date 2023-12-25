package vod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudVodImageSpriteTemplates(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodImageSpriteTemplates,

				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vod_image_sprite_templates.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.sample_type", "Percent"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.sample_interval", "10"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.row_count", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.column_count", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.name", "tf-sprite"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.comment", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.fill_type", "stretch"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.width", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.height", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.resolution_adaptive", "false"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_image_sprite_templates.foo", "template_list.0.update_time"),
				),
			},
		},
	})
}

const testAccVodImageSpriteTemplates = testAccVodImageSpriteTemplate + `
data "tencentcloud_vod_image_sprite_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_image_sprite_template.foo.id
}
`
