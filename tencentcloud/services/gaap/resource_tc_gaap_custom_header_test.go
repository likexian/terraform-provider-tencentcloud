package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapCustomHeaderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCustomHeader,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_custom_header.custom_header", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_custom_header.custom_header", "rule_id", "rule-hddrxgpd"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_custom_header.custom_header", "headers.#", "2"),
				),
			},
			{
				Config: testAccGaapCustomHeaderUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_custom_header.custom_header", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_custom_header.custom_header", "rule_id", "rule-hddrxgpd"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_custom_header.custom_header", "headers.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_custom_header.custom_header",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGaapCustomHeader = `
resource "tencentcloud_gaap_custom_header" "custom_header" {
  rule_id = "rule-hddrxgpd"
  headers {
    header_name  = "HeaderName1"
    header_value = "HeaderValue1"
  }
  headers {
    header_name  = "HeaderName2"
    header_value = "HeaderValue2"
  }
}
`

const testAccGaapCustomHeaderUpdate = `
resource "tencentcloud_gaap_custom_header" "custom_header" {
  rule_id = "rule-hddrxgpd"
  headers {
    header_name  = "HeaderName1"
    header_value = "HeaderValue1"
  }
}
`
