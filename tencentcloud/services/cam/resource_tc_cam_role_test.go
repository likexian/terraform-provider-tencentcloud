package cam_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	tccam "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cam"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_cam_role
	resource.AddTestSweepers("tencentcloud_cam_role", &resource.Sweeper{
		Name: "tencentcloud_cam_role",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()

			service := tccam.NewCamService(client)

			groups, err := service.DescribeRolesByFilter(ctx, nil)
			if err != nil {
				return err
			}

			// add scanning resources
			var resources, nonKeepResources []*tccommon.ResourceInstance
			for _, v := range groups {
				if !tccommon.CheckResourcePersist(*v.RoleName, *v.AddTime) {
					nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
						Id:   *v.RoleId,
						Name: *v.RoleName,
					})
				}
				resources = append(resources, &tccommon.ResourceInstance{
					Id:         *v.RoleId,
					Name:       *v.RoleName,
					CreateTime: *v.AddTime,
				})
			}
			tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateRole")

			for _, v := range groups {
				name := *v.RoleName

				if !strings.HasPrefix(name, "cam-role-test") {
					continue
				}

				request := cam.NewDeleteRoleRequest()
				request.RoleName = v.RoleName
				request.RoleId = v.RoleId
				if _, err := client.UseCamClient().DeleteRole(request); err != nil {
					log.Printf("[%s] error, request: %s \nreason: %s ", request.GetAction(), request.ToJsonString(), err.Error())
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudCamRole_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCamRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamRole_basic(tcacctest.OwnerUin),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamRoleExists("tencentcloud_cam_role.role_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role_basic", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role_basic", "document"),
				),
			}, {
				Config: testAccCamRole_update(tcacctest.OwnerUin),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCamRoleExists("tencentcloud_cam_role.role_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role_basic", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_role.role_basic", "document"),
				),
			},
			{
				ResourceName:      "tencentcloud_cam_role.role_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCamRoleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	camService := tccam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cam_role" {
			continue
		}

		instance, err := camService.DescribeRoleById(ctx, rs.Primary.ID)
		if err == nil && instance != nil {
			return fmt.Errorf("[CHECK][CAM role][Destroy] check: CAM role still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCamRoleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CAM role][Exists] check: CAM role %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CAM role][Exists] check: CAM role id is not set")
		}
		camService := tccam.NewCamService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := camService.DescribeRoleById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CAM role][Exists] check: CAM role %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCamRole_basic(uin string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cam_role" "role_basic" {
	name          = "cam-role-test1"
	document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/%s:uin/%s\"]}}]}"
	description   = "test"
	console_login = true
}`, uin, uin)
}

func testAccCamRole_update(uin string) string {
	return fmt.Sprintf(`
resource "tencentcloud_cam_role" "role_basic" {
  name          = "cam-role-test1"
  document      = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/%s:uin/%s\"]}},{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"qcs\":[\"qcs::cam::uin/%s:uin/%s\"]}}]}"
  console_login = false
}`, uin, uin, uin, uin)
}
