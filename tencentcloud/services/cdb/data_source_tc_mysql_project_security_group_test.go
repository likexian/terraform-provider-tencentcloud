package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlProjectSecurityGroupDataSource_basic -v
func TestAccTencentCloudMysqlProjectSecurityGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlProjectSecurityGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_project_security_group.project_security_group"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_project_security_group.project_security_group", "groups.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_project_security_group.project_security_group", "groups.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_project_security_group.project_security_group", "groups.0.security_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_project_security_group.project_security_group", "groups.0.security_group_name"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_project_security_group.project_security_group", "groups.0.security_group_remark"),
				),
			},
		},
	})
}

const testAccMysqlProjectSecurityGroupDataSource = `

data "tencentcloud_mysql_project_security_group" "project_security_group" {

}

`
