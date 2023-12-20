package chdfs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudChdfsAccessGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccChdfsAccessGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_chdfs_access_group.access_group", "id")),
			},
			{
				Config: testAccChdfsAccessGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_chdfs_access_group.access_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_access_group.access_group", "access_group_name", "testAccessGroupTotal"),
					resource.TestCheckResourceAttr("tencentcloud_chdfs_access_group.access_group", "description", "test access group total"),
				),
			},
			{
				ResourceName:      "tencentcloud_chdfs_access_group.access_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccChdfsAccessGroup = `

resource "tencentcloud_chdfs_access_group" "access_group" {
  access_group_name = "testAccessGroup"
  vpc_type          = 1
  vpc_id            = "vpc-4owdpnwr"
  description       = "test access group"
}

`

const testAccChdfsAccessGroupUpdate = `

resource "tencentcloud_chdfs_access_group" "access_group" {
  access_group_name = "testAccessGroupTotal"
  vpc_type          = 1
  vpc_id            = "vpc-4owdpnwr"
  description       = "test access group total"
}

`
