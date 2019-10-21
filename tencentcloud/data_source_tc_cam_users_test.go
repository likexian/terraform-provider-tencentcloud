package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudCamUsersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCamUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCamUsersDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCamUserExists("tencentcloud_cam_user.user"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_users.users", "user_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_users.users", "user_list.0.remark", "test"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_users.users", "user_list.0.name", "cam-user-tests"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_users.users", "user_list.0.console_login", "true"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_users.users", "user_list.0.phone_num", "13631555963"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_users.users", "user_list.0.country_code", "86"),
					resource.TestCheckResourceAttr("data.tencentcloud_cam_users.users", "user_list.0.email", "1234@qq.com"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_users.users", "user_list.0.uid"),
				),
			},
		},
	})
}

const testAccCamUsersDataSource_basic = `
resource "tencentcloud_cam_user" "user" {
  name                = "cam-user-tests"
  remark              = "test"
  console_login       = true
  use_api             = true
  need_reset_password = true
  password            = "Gail@1234"
  phone_num           = "13631555963"
  country_code        = "86"
  email               = "1234@qq.com"
}
  
data "tencentcloud_cam_users" "users" {
  name = "${tencentcloud_cam_user.user.id}"
}
`
