package sms_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSmsTemplate_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_SMS) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSmsTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sms_template.template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sms_template.template", "template_name", "Template By Terraform"),
					resource.TestCheckResourceAttr("tencentcloud_sms_template.template", "template_content", "Template Content"),
					resource.TestCheckResourceAttr("tencentcloud_sms_template.template", "international", "0"),
					resource.TestCheckResourceAttr("tencentcloud_sms_template.template", "sms_type", "0"),
					resource.TestCheckResourceAttr("tencentcloud_sms_template.template", "remark", "terraform test"),
				),
			},
		},
	})
}

const testAccSmsTemplate = `

resource "tencentcloud_sms_template" "template" {
  template_name = "Template By Terraform"
  template_content = "Template Content"
  international = 0
  sms_type = 0
  remark = "terraform test"
}

`
