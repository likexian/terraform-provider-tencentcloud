package apigateway_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcapigateway "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/apigateway"

	"context"
	"fmt"
	"strings"
	"testing"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudAPIGateWayStrategyAttachment_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testApiStrategyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testApiStrategyAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testApiStrategyAttachmentExists("tencentcloud_api_gateway_strategy_attachment.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "service_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "strategy_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "environment_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_strategy_attachment.test", "bind_api_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_api_gateway_strategy_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testApiStrategyAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_api_gateway_strategy_attachment" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("IP strategy attachment id is broken, id is %s", rs.Primary.ID)
		}
		serviceId := idSplit[0]
		strategyId := idSplit[1]
		bindApiId := idSplit[2]

		has, err := service.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound.InvalidIPStrategy" {
					return nil
				}
			}
			return err
		}

		if has {
			return fmt.Errorf("[CHECK][IP strategy][Destroy] check: IP strategy still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testApiStrategyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := svcapigateway.NewAPIGatewayService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][IP strategy][Exists] check:  %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][IP strategy][Exists] check: id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 4 {
			return fmt.Errorf("IP strategy attachment id is broken, id is %s", rs.Primary.ID)
		}
		serviceId := idSplit[0]
		strategyId := idSplit[1]
		bindApiId := idSplit[2]
		has, err := service.DescribeStrategyAttachment(ctx, serviceId, strategyId, bindApiId)
		if err != nil {
			return err
		}

		if !has {
			return fmt.Errorf("[CHECK][IP strategy][Exists] check: not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAPIGatewayServiceAttachmentBase = `
resource "tencentcloud_api_gateway_service" "service" {
  	service_name = "attach_service"
  	protocol     = "http&https"
  	net_type     = ["INNER", "OUTER"]
  	ip_version   = "IPv4"
}

resource "tencentcloud_api_gateway_ip_strategy" "test"{
    service_id    = tencentcloud_api_gateway_service.service.id
    strategy_name = "attach_strategy"
    strategy_type = "BLACK"
    strategy_data = "9.9.9.9"
}

resource "tencentcloud_api_gateway_api" "api" {
    service_id            = tencentcloud_api_gateway_service.service.id
    api_name              = "attach_api"
    api_desc              = "my hello api update"
    auth_type             = "SECRET"
    protocol              = "HTTP"
    enable_cors           = true
    request_config_path   = "/user/info"
    request_config_method = "POST"
    request_parameters {
    	name          = "email"
        position      = "QUERY"
        type          = "string"
        desc          = "your email please?"
        default_value = "tom@qq.com"
        required      = true
    }
    service_config_type      = "HTTP"
    service_config_timeout   = 10
    service_config_url       = "http://www.tencent.com"
    service_config_path      = "/user"
    service_config_method    = "POST"
    response_type            = "XML"
    response_success_example = "<note>success</note>"
    response_fail_example    = "<note>fail</note>"
    response_error_codes {
    	code           = 20
        msg            = "system error"
       	desc           = "system error code"
       	converted_code = 10
        need_convert   = true
	}
}

resource "tencentcloud_api_gateway_service_release" "service" {
  service_id       = tencentcloud_api_gateway_api.api.service_id
  environment_name = "release"
  release_desc     = "test service release"
}
`

const testApiStrategyAttachment_basic = testAPIGatewayServiceAttachmentBase + `
resource "tencentcloud_api_gateway_strategy_attachment" "test"{
   service_id       = tencentcloud_api_gateway_service_release.service.service_id
   strategy_id      = tencentcloud_api_gateway_ip_strategy.test.strategy_id 
   environment_name = "release"
   bind_api_id      = tencentcloud_api_gateway_api.api.id
}
`
