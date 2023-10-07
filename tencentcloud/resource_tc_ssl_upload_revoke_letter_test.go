package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudSslUploadRevokeLetterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSslUploadRevokeLetter,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ssl_upload_revoke_letter.upload_revoke_letter", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ssl_upload_revoke_letter.upload_revoke_letter", "certificate_id", "8xRYdDlc"),
					resource.TestCheckResourceAttrSet("tencentcloud_ssl_upload_revoke_letter.upload_revoke_letter", "revoke_letter"),
				),
			},
			{
				ResourceName:      "tencentcloud_ssl_upload_revoke_letter.upload_revoke_letter",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSslUploadRevokeLetter = `

resource "tencentcloud_ssl_upload_revoke_letter" "upload_revoke_letter" {
  certificate_id = "8xRYdDlc"
  revoke_letter = filebase64("./c.pdf")
}

`
