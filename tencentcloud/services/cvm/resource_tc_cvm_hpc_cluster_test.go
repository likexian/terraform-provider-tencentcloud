package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmHpcClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config:    testAccCvmHpcCluster,
				PreConfig: func() { tcacctest.AccStepSetRegion(t, "ap-beijing") },
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_hpc_cluster.hpc_cluster", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_hpc_cluster.hpc_cluster", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_hpc_cluster.hpc_cluster", "remark"),
				),
			},
			{
				ResourceName:      "tencentcloud_cvm_hpc_cluster.hpc_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmHpcCluster = `

resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
  zone = "ap-beijing-6"
  name = "terraform-test"
  remark = "create for test"
}

`
