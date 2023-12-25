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

func TestAccTencentCloudMonitorRecordingRule_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckRecordingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRecordingRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordingRuleExists("tencentcloud_monitor_tmp_recording_rule.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_recording_rule.basic", "name", "recording_rule-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_recording_rule.basic", "rule_state", "2"),
				),
			},
			{
				Config: testRecordingRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordingRuleExists("tencentcloud_monitor_tmp_recording_rule.update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_recording_rule.update", "name", "recording_rule-update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_recording_rule.update", "rule_state", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_recording_rule.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRecordingRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_recording_rule" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		instance, err := service.DescribeMonitorRecordingRuleById(ctx, ids[0], ids[1])
		if err != nil {
			return err
		}

		if instance != nil && *instance.RuleState != 1 {
			return fmt.Errorf("instance %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckRecordingRuleExists(r string) resource.TestCheckFunc {
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

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeMonitorRecordingRuleById(ctx, ids[0], ids[1])
		if err != nil {
			return err
		}

		if instance == nil || *instance.RuleState == 1 {
			return fmt.Errorf("instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testRecordingRuleVar = `
variable "prometheus_id" {
  default = "` + tcacctest.DefaultPrometheusId + `"
}
`
const testRecordingRule_basic = testRecordingRuleVar + `
resource "tencentcloud_monitor_tmp_recording_rule" "basic" {
  name			= "recording_rule-test"
  instance_id	= var.prometheus_id
  rule_state	= 2
  group			= <<EOF
---
name: example-test
rules:
  - record: job:http_inprogress_requests:sum
    expr: sum by (job) (http_inprogress_requests)
EOF
}`

const testRecordingRule_update = testRecordingRuleVar + `
resource "tencentcloud_monitor_tmp_recording_rule" "update" {
  name			= "recording_rule-update"
  instance_id	= var.prometheus_id
  rule_state	= 3
  group			= <<EOF
---
name: example-test-update
rules:
  - record: job:http_inprogress_requests:sum
    expr: sum by (job) (http_inprogress_requests)
EOF
}`
