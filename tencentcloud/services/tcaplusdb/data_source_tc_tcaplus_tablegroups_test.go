package tcaplusdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTcaplusGroupsName = "data.tencentcloud_tcaplus_tablegroups.id_test"

func TestAccTencentCloudTcaplusGroupsData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTcaplusGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTcaplusGroupsBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "cluster_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "tablegroup_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.#"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.0.tablegroup_name"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.0.table_count"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.0.tablegroup_id"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.0.total_size"),
					resource.TestCheckResourceAttrSet(testDataTcaplusGroupsName, "list.0.create_time"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTcaplusGroupsBasic = tcacctest.DefaultTcaPlusData + `

data "tencentcloud_tcaplus_tablegroups" "id_test" {
   cluster_id         = local.tcaplus_id
   tablegroup_id      = local.tcaplus_table_group_id
}
`
