package tsf_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctsf "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApiRateLimitRuleResource_basic -v
func TestAccTencentCloudTsfApiRateLimitRuleResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfApiRateLimitRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiRateLimitRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApiRateLimitRuleExists("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "api_id", tcacctest.DefaultTsfApiId),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "max_qps", "10"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule", "usable_status", "enabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfApiRateLimitRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_api_rate_limit_rule" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		apiId := idSplit[0]
		ruleId := idSplit[1]

		res, err := service.DescribeTsfApiRateLimitRuleById(ctx, apiId, ruleId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf ApiRateLimitRule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApiRateLimitRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}
		apiId := idSplit[0]
		ruleId := idSplit[1]

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfApiRateLimitRuleById(ctx, apiId, ruleId)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf ApiRateLimitRule %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApiRateLimitRuleVar = `
variable "api_id" {
	default = "` + tcacctest.DefaultTsfApiId + `"
}
`

const testAccTsfApiRateLimitRule = testAccTsfApiRateLimitRuleVar + `

resource "tencentcloud_tsf_api_rate_limit_rule" "api_rate_limit_rule" {
	api_id = var.api_id
	max_qps = 10
	usable_status = "enabled"
}

`
