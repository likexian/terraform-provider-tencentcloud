package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudCamGroupPolicyAttachment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamGroupPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamGroupPolicyAttachment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamGroupPolicyAttachmentExists("tencentcloud_cam_group_policy_attachment.group_policy_attachment_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_group_policy_attachment.group_policy_attachment_basic", "group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_group_policy_attachment.group_policy_attachment_basic", "policy_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_group_policy_attachment.group_policy_attachment_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamGroupPolicyAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	camService := CamService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_group_policy_attachment" {
			continue
		}

		_, err := camService.DescribeGroupPolicyAttachmentById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("CAM group policy attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamGroupPolicyAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("CAM group policy attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("CAM group policy attachment id is not set")
		}
		camService := CamService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := camService.DescribeGroupPolicyAttachmentById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccCamGroupPolicyAttachment_basic = `
resource "tencentcloud_cam_group" "group" {
  name   = "cam-group-test2"
  remark = "test"
}
  
resource "tencentcloud_cam_policy" "policy" {
  name        = "cam-policy-test9"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
	description = "test"
}
  
resource "tencentcloud_cam_group_policy_attachment" "group_policy_attachment_basic" {
  group_id  = "${tencentcloud_cam_group.group.id}"
  policy_id = "${tencentcloud_cam_policy.policy.id}"
}
`
