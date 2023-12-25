package tke_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctke "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var (
	testTkeServerlessNodePoolResourceKey = "tencentcloud_kubernetes_serverless_node_pool.pool_example"
)

const (
	clusterIdForTkeTestEnvKey = "TKE_CLUSTER_ID_FOR_SEVER_LESS_NODE_POOL_TEST"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_serverless_node_pool
	resource.AddTestSweepers("tencentcloud_serverless_node_pool", &resource.Sweeper{
		Name: "tencentcloud_serverless_node_pool",
		F:    testServerlessNodePoolSweep,
	})
}

func TestAccTencentCloudKubernetesServerlessNodePoolResource_Basic(t *testing.T) {
	t.Parallel()

	tkeClusterId := tcacctest.DefaultTkeClusterId
	envClusterId := os.Getenv(clusterIdForTkeTestEnvKey)
	if strings.HasPrefix(envClusterId, "cls-") {
		tkeClusterId = envClusterId
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckServerlessNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: getTestAccTkeServerlessNodePoolConfig(tkeClusterId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(testTkeServerlessNodePoolResourceKey),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "name", "hello-world"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "labels.key1", "value1"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "taints.#", "1"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "taints.0.key", "no-eip-instance"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "taints.0.value", "yes"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "serverless_nodes.0.display_name", "serverless_node1"),
				),
			},
			{
				Config: getTestAccTkeServerlessNodePoolUpdateConfig(tkeClusterId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "name", "hello-world-2"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "labels.key2", "value2"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "taints.#", "2"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "taints.1.key", "no-cbs-instance"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "taints.1.value", "no"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "taints.1.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(testTkeServerlessNodePoolResourceKey, "serverless_nodes.0.display_name", "serverless_node2"),
				),
			},
		},
	})
}

func testServerlessNodePoolSweep(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tkeClusterId := tcacctest.DefaultTkeClusterId
	envClusterId := os.Getenv(clusterIdForTkeTestEnvKey)
	if strings.HasPrefix(envClusterId, "cls-") {
		tkeClusterId = envClusterId
	}
	log.Printf("testServerlessNodePoolSweep region %s, clusterId %s", region, tkeClusterId)

	cli, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return err
	}
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svctke.NewTkeService(client)

	request := tke.NewDescribeClusterVirtualNodePoolsRequest()
	request.ClusterId = helper.String(tkeClusterId)
	response, err := client.UseTkeClient().DescribeClusterVirtualNodePools(request)
	if err != nil {
		log.Printf("Query %s serverless node pool fail: %s", tkeClusterId, err.Error())
		return err
	}
	nodePools := response.Response.NodePoolSet
	if len(nodePools) == 0 {
		return nil
	}
	for i := range nodePools {
		poolId := *nodePools[i].NodePoolId
		poolName := nodePools[i].Name
		if poolName == nil {
			continue
		}

		if !nodePoolNameReg.MatchString(*poolName) {
			continue
		}
		delReq := tke.NewDeleteClusterVirtualNodePoolRequest()
		delReq.ClusterId = common.StringPtr(tkeClusterId)
		delReq.NodePoolIds = common.StringPtrs([]string{poolId})
		err := service.DeleteClusterVirtualNodePool(ctx, delReq)
		if err != nil {
			continue
		}
	}
	return nil
}

func testAccCheckServerlessNodePoolDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tkeService := svctke.NewTkeService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	tkeClusterId := tcacctest.DefaultTkeClusterId
	envClusterId := os.Getenv(clusterIdForTkeTestEnvKey)
	if strings.HasPrefix(envClusterId, "cls-") {
		tkeClusterId = envClusterId
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_kubernetes_serverless_node_pool" {
			continue
		}
		respNodePool, has, err := tkeService.DescribeServerlessNodePoolByClusterIdAndNodePoolId(ctx, tkeClusterId, rs.Primary.ID)

		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				respNodePool, has, err = tkeService.DescribeServerlessNodePoolByClusterIdAndNodePoolId(ctx, tkeClusterId, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return err
		}

		if has && *respNodePool.LifeState != "deleting" {
			return fmt.Errorf("tke serverless node pool instance still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func getTestAccTkeServerlessNodePoolConfig(clusterId string) string {
	return fmt.Sprintf(testAccTkeServerlessNodePoolTemplate, clusterId)
}

func getTestAccTkeServerlessNodePoolUpdateConfig(clusterId string) string {
	return fmt.Sprintf(testAccTkeServerlessNodePoolUpdateTemplate, clusterId)
}

const (
	testAccTkeServerlessNodePoolTemplate = `
data "tencentcloud_kubernetes_clusters" "existed_cluster" {
  cluster_id = "%s"
}

data "tencentcloud_security_groups" "sg" {
  name = "default"
}

data "tencentcloud_vpc_instances" "vpc_cluster" {
  vpc_id = data.tencentcloud_kubernetes_clusters.existed_cluster.list.0.vpc_id
}

resource "tencentcloud_kubernetes_serverless_node_pool" "pool_example" {
  cluster_id = data.tencentcloud_kubernetes_clusters.existed_cluster.list.0.cluster_id
  name = "hello-world"
  serverless_nodes {
    display_name = "serverless_node1"
    subnet_id = data.tencentcloud_vpc_instances.vpc_cluster.instance_list.0.subnet_ids.0
  }
  labels = {
    "key1" = "value1"
  }
  taints {
    key = "no-eip-instance"
    value = "yes"
    effect = "NoSchedule"
  }
  security_group_ids = [data.tencentcloud_security_groups.sg.id]
}
`

	testAccTkeServerlessNodePoolUpdateTemplate = `
data "tencentcloud_kubernetes_clusters" "existed_cluster" {
  cluster_id = "%s"
}

data "tencentcloud_security_groups" "sg" {
  name = "default"
}

data "tencentcloud_vpc_instances" "vpc_cluster" {
  vpc_id = data.tencentcloud_kubernetes_clusters.existed_cluster.list.0.vpc_id
}

resource "tencentcloud_kubernetes_serverless_node_pool" "pool_example" {
  cluster_id = data.tencentcloud_kubernetes_clusters.existed_cluster.list.0.cluster_id
  name = "hello-world-2"
  serverless_nodes {
    display_name = "serverless_node2"
    subnet_id = data.tencentcloud_vpc_instances.vpc_cluster.instance_list.0.subnet_ids.0
  }
  labels = {
    "key2" = "value2"
  }
  taints {
    key = "no-eip-instance"
    value = "yes"
    effect = "NoSchedule"
  }
  taints {
    key = "no-cbs-instance"
    value = "no"
    effect = "NoSchedule"
  }	
  security_group_ids = [data.tencentcloud_security_groups.sg.id]
}
`
)
