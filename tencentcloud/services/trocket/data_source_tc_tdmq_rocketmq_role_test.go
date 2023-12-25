package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTdmqRocketmqRoleDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTdmqRocketmqRole,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tdmq_rocketmq_role.role"),
					resource.TestCheckResourceAttr("data.tencentcloud_tdmq_rocketmq_role.role", "role_sets.#", "1"),
				),
			},
		},
	})
}

const testAccDataSourceTdmqRocketmqRole = `
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
	cluster_name = "test_rocketmq_datasource_role"
	remark = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_role" "role" {
  role_name = "test_rocketmq_role"
  remark = "test rocketmq role"
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
}

data "tencentcloud_tdmq_rocketmq_role" "role" {
  role_name = tencentcloud_tdmq_rocketmq_role.role.role_name
  cluster_id = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
}
`
