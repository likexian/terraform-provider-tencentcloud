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

// go test -i; go test -test.run TestAccTencentCloudMonitorTmpTkeBasicConfigResource_basic -v
func TestAccTencentCloudMonitorTmpTkeBasicConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpTkeBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmpTkeBasicConfigExists("tencentcloud_monitor_tmp_tke_basic_config.tmp_tke_basic_config"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_basic_config.tmp_tke_basic_config", "name", "cadvisor"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_basic_config.tmp_tke_basic_config", "metrics_name.#", "2"),
				),
			},
		},
	})
}

func testAccCheckTmpTkeBasicConfigExists(r string) resource.TestCheckFunc {
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
		if len(items) != 4 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		clusterType := items[1]
		clusterId := items[2]
		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		agents, err := service.DescribeTkeTmpBasicConfigById(ctx, clusterId, clusterType, instanceId)
		if agents == nil {
			return fmt.Errorf("config %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccMonitorTmpTkeBasicConfigVar = `
variable "prometheus_id" {
	default = "` + tcacctest.DefaultPrometheusId + `"
  }
variable "cluster_id" {
  default = "cls-2trvpflc"
}
variable "cluster_type" {
  default = "` + tcacctest.TkeClusterTypeAgent + `"
}`

const testAccMonitorTmpTkeBasicConfig = testAccMonitorTmpTkeBasicConfigVar + `

resource "tencentcloud_monitor_tmp_tke_basic_config" "tmp_tke_basic_config" {
  instance_id  = var.prometheus_id
  cluster_type = var.cluster_type
  cluster_id   = var.cluster_id
  name = "cadvisor"
  metrics_name = ["container_cpu_usage_seconds_total", "container_fs_limit_bytes"]

}

`
