package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlVerifyRootAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlVerifyRootAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_verify_root_account.verify_root_account", "id")),
			},
		},
	})
}

const testAccMysqlVerifyRootAccount = `

resource "tencentcloud_mysql_verify_root_account" "verify_root_account" {
  instance_id = ""
  password = ""
}

`
