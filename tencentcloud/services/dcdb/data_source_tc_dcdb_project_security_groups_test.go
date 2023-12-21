package dcdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbProjectSecurityGroupsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbProjectSecurityGroupsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_project_security_groups.project_security_groups"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_project_security_groups.project_security_groups", "product", "dcdb"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_project_security_groups.project_security_groups", "project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_project_security_groups.project_security_groups", "groups.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_project_security_groups.project_security_groups", "groups.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_project_security_groups.project_security_groups", "groups.0.security_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_project_security_groups.project_security_groups", "groups.0.security_group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_project_security_groups.project_security_groups", "groups.0.inbound.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_project_security_groups.project_security_groups", "groups.0.outbound.#"),
				),
			},
		},
	})
}

const testAccDcdbProjectSecurityGroupsDataSource = `

data "tencentcloud_dcdb_project_security_groups" "project_security_groups" {
  product    = "dcdb"
  project_id = 0
}

`
