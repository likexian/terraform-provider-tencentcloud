package tcm_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTcmClusterAttachment_basic -v
func TestAccTencentCloudTcmClusterAttachment_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClusterAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmClusterAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterAttachmentExists("tencentcloud_tcm_cluster_attachment.basic"),
					// resource.TestCheckResourceAttrSet("tencentcloud_tcm_cluster_attachment.basic", "mesh_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "cluster_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "cluster_list.0.cluster_id", tcacctest.DefaultMeshClusterId),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "cluster_list.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "cluster_list.0.role", "REMOTE"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "cluster_list.0.vpc_id", tcacctest.DefaultMeshVpcId),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "cluster_list.0.subnet_id", tcacctest.DefaultMeshSubnetId),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "cluster_list.0.type", "EKS"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_cluster_attachment.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClusterAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctcm.NewTcmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcm_cluster_attachment" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}
		meshId := ids[0]
		clusterId := ids[1]

		mesh, err := service.DescribeTcmMesh(ctx, meshId)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound" {
					return nil
				}
			}
			return err
		}

		if mesh != nil {
			if len(mesh.Mesh.ClusterList) > 0 {
				for _, v := range mesh.Mesh.ClusterList {
					if *v.ClusterId == clusterId {
						return fmt.Errorf("clusterList %s still exists", rs.Primary.ID)
					}
				}
			}
		}
	}

	return nil
}

func testAccCheckClusterAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}
		meshId := ids[0]
		clusterId := ids[1]

		service := svctcm.NewTcmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		mesh, err := service.DescribeTcmMesh(ctx, meshId)
		if err != nil {
			return err
		}

		if mesh == nil {
			return fmt.Errorf("mesh %s is not found", rs.Primary.ID)
		}
		if len(mesh.Mesh.ClusterList) > 0 {
			for _, v := range mesh.Mesh.ClusterList {
				if *v.ClusterId == clusterId {
					return nil
				}
			}
			return fmt.Errorf("mesh clusterList %s is not found", rs.Primary.ID)
		} else {
			return fmt.Errorf("clusterList %s is not found", rs.Primary.ID)
		}
	}
}

const testAccTcmClusterAttachmentVar = `
variable "cluster_id" {
  default = "` + tcacctest.DefaultMeshClusterId + `"
}
variable "vpc_id" {
  default = "` + tcacctest.DefaultMeshVpcId + `"
}
variable "subnet_id" {
  default = "` + tcacctest.DefaultMeshSubnetId + `"
}
`

const testAccTcmClusterAttachment = testAccTcmClusterAttachmentVar + `

resource "tencentcloud_tcm_mesh" "basic" {
	display_name = "test_mesh"
	mesh_version = "1.12.5"
	type = "HOSTED"
	config {
	  istio {
		outbound_traffic_policy = "ALLOW_ANY"
		disable_policy_checks = true
		enable_pilot_http = true
		disable_http_retry = true
		smart_dns {
		  istio_meta_dns_capture = true
		  istio_meta_dns_auto_allocate = true
		}
	  }
	  tracing {
		  enable = true
		  sampling = 1
		  apm {
			  enable = false
		  }
		  zipkin {
			  address = "10.0.0.1:1000"
		  }
	  }
	}
	tag_list {
	  key = "key"
	  value = "value"
	  passthrough = false
	}
  }

resource "tencentcloud_tcm_cluster_attachment" "basic" {
  mesh_id = tencentcloud_tcm_mesh.basic.id
  cluster_list {
    cluster_id = var.cluster_id
    region = "ap-guangzhou"
    role = "REMOTE"
    vpc_id = var.vpc_id
    subnet_id = var.subnet_id
    type = "EKS"
  }
}

`
