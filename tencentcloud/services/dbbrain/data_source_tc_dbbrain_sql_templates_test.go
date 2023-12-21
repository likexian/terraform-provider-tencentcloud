package dbbrain_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainSqlTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainSqlTemplatesDataSource, tcacctest.DefaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_sql_templates.sql_templates"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_sql_templates.sql_templates", "instance_id", tcacctest.DefaultDbBrainInstanceId),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_sql_templates.sql_templates", "schema", "tf_ci_test"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_sql_templates.sql_templates", "sql_text", "select sleep(5);"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_sql_templates.sql_templates", "product", "mysql"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_sql_templates.sql_templates", "sql_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_sql_templates.sql_templates", "sql_template"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_sql_templates.sql_templates", "sql_id"),
				),
			},
		},
	})
}

const testAccDbbrainSqlTemplatesDataSource = `

data "tencentcloud_dbbrain_sql_templates" "sql_templates" {
  instance_id = "%s"
  schema = "tf_ci_test"
  sql_text = "select sleep(5);"
  product = "mysql"
}

`
