package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlDeployGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDeployGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_deploy_group.deploy_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_deploy_group.deploy_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlDeployGroup = `

resource "tencentcloud_mysql_deploy_group" "deploy_group" {
  deploy_group_name = "terrform-deploy"
  description       = "deploy test"
  limit_num         = 1
  dev_class         = ["TS85"]
}

`
