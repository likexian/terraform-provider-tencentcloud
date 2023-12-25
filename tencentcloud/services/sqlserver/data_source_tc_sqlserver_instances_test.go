package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataSqlserverInstancesName = "data.tencentcloud_sqlserver_instances.example"

// go test -i; go test -test.run TestAccDataSourceTencentCloudSqlserverInstances -v
func TestAccDataSourceTencentCloudSqlserverInstances(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSqlserverInstancesBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.#"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.id"),
					resource.TestCheckResourceAttr(testDataSqlserverInstancesName, "instance_list.0.charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.engine_version"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.memory"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.storage"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.vip"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.vport"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.status"),
					resource.TestCheckResourceAttrSet(testDataSqlserverInstancesName, "instance_list.0.used_storage"),
				),
			},
		},
	})
}

var testAccTencentCloudDataSqlserverInstancesBasic = `
data "tencentcloud_sqlserver_instances" "example"{
  name = "keep"
}
`
