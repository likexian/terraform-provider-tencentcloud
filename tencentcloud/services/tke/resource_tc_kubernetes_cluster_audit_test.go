package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterAuditResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesClusterAudit,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_audit.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_audit.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_audit.example", "delete_logset_and_topic"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_cluster_audit.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccKubernetesClusterAudit = `
resource "tencentcloud_kubernetes_cluster_audit" "example" {
  cluster_id              = "cls-fdy7hm1q"
  delete_logset_and_topic = true
}
`
