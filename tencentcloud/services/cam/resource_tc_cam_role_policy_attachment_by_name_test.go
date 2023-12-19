package cam_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cam"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudCamRolePolicyAttachmentByNameResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCamRolePolicyAttachmentByNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRolePolicyAttachmentByName_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamRolePolicyAttachmentByNameExists("tencentcloud_cam_role_policy_attachment_by_name.role_policy_attachment_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_policy_attachment_by_name.role_policy_attachment_basic", "role_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role_policy_attachment_by_name.role_policy_attachment_basic", "policy_name"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_role_policy_attachment_by_name.role_policy_attachment_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamRolePolicyAttachmentByNameDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	camService := cam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_role_policy_attachment_by_name" {
			continue
		}
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) < 2 {
			return fmt.Errorf("RolePolicyAttachmentId is invalid!")
		}
		roleName, policyName := items[0], items[1]
		params := make(map[string]interface{})
		params["policy_name"] = policyName
		instance, err := camService.DescribeRolePolicyAttachmentByName(ctx, roleName, params)

		if err == nil && instance != nil {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Desctroy] check: CAM role policy attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamRolePolicyAttachmentByNameExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment id is not set")
		}
		camService := cam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		items := strings.Split(rs.Primary.ID, "#")
		if len(items) < 2 {
			return fmt.Errorf("RolePolicyAttachmentId is invalid!")
		}
		roleName, policyName := items[0], items[1]
		params := make(map[string]interface{})
		params["policy_name"] = policyName
		instance, err := camService.DescribeRolePolicyAttachmentByName(ctx, roleName, params)

		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CAM role policy attachment][Exist] check: CAM role policy attachment %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

// need to add policy resource definition
func testAccCamRolePolicyAttachmentByName_basic() string {
	return tcacctest.DefaultCamVariables + `
resource "tencentcloud_cam_role_policy_attachment_by_name" "role_policy_attachment_basic" {
  role_name   = var.cam_role_basic
  policy_name = var.cam_policy_basic
}
`
}
