package ciam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCiamUserGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiamUserGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ciam_user_group.user_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_ciam_user_group.user_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiamUserGroup = `

resource "tencentcloud_ciam_user_store" "user_store" {
  user_pool_name = "tf_user_store_test"
  user_pool_desc = "for terraform test"
  user_pool_logo = "https://ciam-prd-1302490086.cos.ap-guangzhou.myqcloud.com/temporary/92630252a2c5422d9663db5feafd619b.png"
}

resource "tencentcloud_ciam_user_group" "user_group" {
  display_name  = "tf_user_group"
  user_store_id = tencentcloud_ciam_user_store.user_store.id
  description   = "for terrafrom test"
}

`
