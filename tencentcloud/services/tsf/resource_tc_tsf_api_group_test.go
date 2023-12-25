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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApiGroupResource_basic -v
func TestAccTencentCloudTsfApiGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTsfApiGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTsfApiGroupExists("tencentcloud_tsf_api_group.api_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_api_group.api_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "group_name", "terraform_test_group"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "group_context", "/terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "auth_type", "none"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "description", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "group_type", "ms"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "namespace_name_key_position", "path"),
					resource.TestCheckResourceAttr("tencentcloud_tsf_api_group.api_group", "service_name_key_position", "path"),
				),
			},
			{
				ResourceName:      "tencentcloud_tsf_api_group.api_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTsfApiGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tsf_api_group" {
			continue
		}

		res, err := service.DescribeTsfApiGroupById(ctx, rs.Primary.ID)
		if err != nil {
			code := err.(*sdkErrors.TencentCloudSDKError).Code
			if code == "InvalidParameterValue.GatewayParameterInvalid" {
				return nil
			}
			return err
		}

		if res != nil {
			return fmt.Errorf("tsf api group %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckTsfApiGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svctsf.NewTsfService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTsfApiGroupById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res == nil {
			return fmt.Errorf("tsf api group %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccTsfApiGroupVar = `
variable "gateway_instance_id" {
	default = "` + tcacctest.DefaultTsfGateway + `"
}
`

const testAccTsfApiGroup = testAccTsfApiGroupVar + `

resource "tencentcloud_tsf_api_group" "api_group" {
	group_name = "terraform_test_group"
	group_context = "/terraform-test"
	auth_type = "none"
	description = "terraform-test"
	group_type = "ms"
	gateway_instance_id = var.gateway_instance_id
	# namespace_name_key = "path"
	# service_name_key = "path"
	namespace_name_key_position = "path"
	service_name_key_position = "path"
  }

`
