package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbReloadProxyNodeResource_basic -v
func TestAccTencentCloudCynosdbReloadProxyNodeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckCynosdbProxyDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbReloadProxyNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_reload_proxy_node.reload_proxy_node", "id")),
			},
		},
	})
}

const testAccCynosdbReloadProxyNode = testAccCynosdbProxy + `
resource "tencentcloud_cynosdb_reload_proxy_node" "reload_proxy_node" {
  cluster_id     = tencentcloud_cynosdb_proxy.proxy.id
  proxy_group_id = tencentcloud_cynosdb_proxy.proxy.proxy_group_id
}
`
