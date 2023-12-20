package chdfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudChdfsAccessRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsAccessRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_access_rule.access_rule", "id")),
			},
			{
				Config: testAccChdfsAccessRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_chdfs_access_rule.access_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_access_rule.access_rule", "access_rule.0.address", "10.0.0.1"),
				),
			},
			{
				ResourceName:      "tencentcloud_chdfs_access_rule.access_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsAccessRule = `

resource "tencentcloud_chdfs_access_rule" "access_rule" {
  access_group_id = "ag-bvmzrbsm"

  access_rule {
    access_mode    = 2
    address        = "10.0.1.1"
    priority       = 12
  }
}

`

const testAccChdfsAccessRuleUpdate = `

resource "tencentcloud_chdfs_access_rule" "access_rule" {
  access_group_id = "ag-bvmzrbsm"

  access_rule {
    access_mode    = 1
    address        = "10.0.0.1"
    priority       = 10
  }
}

`
