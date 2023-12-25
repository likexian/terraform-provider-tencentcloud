package tsf_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctsf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfClusterResource_basic -v
func TestAccTencentCloudTsfClusterResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfClusterExists("tencentcloud_tsf_cluster.cluster"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_cluster.cluster", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "cluster_name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "vpc_id", "vpc-kphn8u93"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "cluster_cidr", "9.165.120.0/24"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "cluster_desc", "test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "tsf_region_id", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "cluster_version", "1.18.4"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "max_node_pod_num", "32"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "max_cluster_service_num", "128"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_cluster.cluster", "tags.createdBy", "terraform"),
				),
			},
			// {
			// 	ResourceName:      "tencentcloud_tsf_cluster.cluster",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckTsfClusterDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_cluster" {
			continue
		}

		res, err := service.DescribeTsfClusterById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf cluster %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfClusterExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfClusterById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf cluster %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfCluster = `

resource "tencentcloud_tsf_cluster" "cluster" {
	cluster_name = "terraform-test"
	cluster_type = "C"
	vpc_id = "vpc-kphn8u93"
	cluster_cidr = "9.165.120.0/24"
	cluster_desc = "test"
	tsf_region_id = "ap-guangzhou"
	cluster_version = "1.18.4"
	max_node_pod_num = 32
	max_cluster_service_num = 128
	tags = {
	  "createdBy" = "terraform"
	}
}
`
