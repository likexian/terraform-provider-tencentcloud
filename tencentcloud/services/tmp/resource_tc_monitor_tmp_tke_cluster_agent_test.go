package tmp_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_monitor_tmp_tke_cluster_agent
	resource.AddTestSweepers("tencentcloud_monitor_tmp_tke_cluster_agent", &resource.Sweeper{
		Name: "tencentcloud_monitor_tmp_tke_cluster_agent",
		F:    testSweepClusterAgent,
	})
}
func testSweepClusterAgent(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(region)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcmonitor.NewMonitorService(client)

	instanceId := tcacctest.ClusterPrometheusId
	clusterId := tcacctest.TkeClusterIdAgent
	clusterType := tcacctest.TkeClusterTypeAgent

	agents, err := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
	if err != nil {
		return err
	}

	if agents != nil {
		return nil
	}

	err = service.DeletePrometheusClusterAgent(ctx, instanceId, clusterId, clusterType)
	if err != nil {
		return err
	}

	return nil
}

// go test -i; go test -test.run TestAccTencentCloudMonitorClusterAgent_basic -v
func TestAccTencentCloudMonitorClusterAgent_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClusterAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testClusterAgentYaml_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterAgentExists("tencentcloud_monitor_tmp_tke_cluster_agent.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_cluster_agent.basic", "agents.0.cluster_id", "cls-9ae9qo9k"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_cluster_agent.basic", "agents.0.cluster_type", "eks"),
				),
			},
		},
	})
}

func testAccCheckClusterAgentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_cluster_agent" {
			continue
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		clusterId := items[1]
		clusterType := items[2]
		agents, err := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
		if agents != nil {
			return fmt.Errorf("cluster agent %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckClusterAgentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("instance id is not set")
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		clusterId := items[1]
		clusterType := items[2]
		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		agents, err := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
		if agents == nil {
			return fmt.Errorf("cluster agent %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testClusterAgentYamlVar = `
variable "prometheus_id" {
  default = "` + tcacctest.ClusterPrometheusId + `"
}
variable "default_region" {
  default = "` + tcacctest.DefaultRegion + `"
}
variable "agent_cluster_id" {
  default = "` + tcacctest.TkeClusterIdAgent + `"
}
variable "agent_cluster_type" {
  default = "` + tcacctest.TkeClusterTypeAgent + `"
}`

const testClusterAgentYaml_basic = testClusterAgentYamlVar + `
resource "tencentcloud_monitor_tmp_tke_cluster_agent" "basic" {
  instance_id = var.prometheus_id
  agents {
    region          = var.default_region
    cluster_type    = var.agent_cluster_type
    cluster_id      = var.agent_cluster_id
    enable_external = false
  }
}`
