package chdfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixChdfsLifeCycleRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsLifeCycleRule,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_life_cycle_rule.life_cycle_rule", "id")),
			},
			{
				Config: testAccChdfsLifeCycleRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_chdfs_life_cycle_rule.life_cycle_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_life_cycle_rule.life_cycle_rule", "life_cycle_rule.0.life_cycle_rule_name", "terraform-for-test"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_life_cycle_rule.life_cycle_rule", "life_cycle_rule.0.path", "/terraform"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_life_cycle_rule.life_cycle_rule", "life_cycle_rule.0.status", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_chdfs_life_cycle_rule.life_cycle_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsLifeCycleRule = `

resource "tencentcloud_chdfs_life_cycle_rule" "life_cycle_rule" {
  file_system_id = "f14mpfy5lh4e"

  life_cycle_rule {
    life_cycle_rule_name = "terraform-test"
    path                 = "/test"
    status               = 1

    transitions {
      days = 30
      type = 1
    }
  }
}

`

const testAccChdfsLifeCycleRuleUpdate = `

resource "tencentcloud_chdfs_life_cycle_rule" "life_cycle_rule" {
  file_system_id = "f14mpfy5lh4e"

  life_cycle_rule {
    life_cycle_rule_name = "terraform-for-test"
    path                 = "/terraform"
    status               = 2

    transitions {
      days = 30
      type = 1
    }
  }
}

`
