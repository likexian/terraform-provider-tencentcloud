package tag_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTagAttachmentResource_basic -v
func TestAccTencentCloudTagAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTagAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTagResourceTag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTagAttachmentExists("tencentcloud_tag_attachment.tag_attachment"),
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.tag_attachment", "tag_key", "test_terraform_tagAttachment_key"),
					resource.TestCheckResourceAttr("tencentcloud_tag_attachment.tag_attachment", "tag_value", "Terraform_tagAttachment_value"),
					resource.TestCheckResourceAttrSet("tencentcloud_tag_attachment.tag_attachment", "resource")),
			},
			{
				ResourceName:      "tencentcloud_tag_attachment.tag_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccCheckTagAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tag_attachment" {
			continue
		}
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svctag.NewTagService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		tags, err := service.DescribeTagTagAttachmentById(ctx, rs.Primary.Attributes["tag_key"],
			rs.Primary.Attributes["tag_value"], rs.Primary.Attributes["resource"])
		if err != nil {
			return err
		}
		if tags == nil {
			return nil
		}
		return fmt.Errorf("delete tagAttachment key %s fail, still on server", rs.Primary.Attributes["tag_key"])
	}
	return nil
}

func testAccCheckTagAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctag.NewTagService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTagTagAttachmentById(ctx, rs.Primary.Attributes["tag_key"],
			rs.Primary.Attributes["tag_value"], rs.Primary.Attributes["resource"])
		if err != nil {
			return err
		}
		if res != nil && res.Resource != nil && res.Tags != nil {
			return nil
		}

		return fmt.Errorf("tagAttachment %s not found on server", rs.Primary.Attributes["tag_key"])
	}
}

const testAccTagResourceTag = tcacctest.DefaultCvmModificationVariable + `
data "tencentcloud_user_info" "info" {}

locals {
  uin = data.tencentcloud_user_info.info.uin
}

resource "tencentcloud_tag_attachment" "tag_attachment" {
  tag_key = "test_terraform_tagAttachment_key"
  tag_value = "Terraform_tagAttachment_value"
  resource = "qcs::cvm:ap-guangzhou:uin/${local.uin}:instance/${var.cvm_id}"
}

`
