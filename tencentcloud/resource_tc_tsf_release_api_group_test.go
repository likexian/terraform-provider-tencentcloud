package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfReleaseApiGroupResource_basic -v
func TestAccTencentCloudTsfReleaseApiGroupResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTsfUnitNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfReleaseApiGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tsf_release_api_group.release_api_group", "id"),
				),
			},
		},
	})
}

const testAccTsfReleaseApiGroup = `

resource "tencentcloud_tsf_start_container_group" "start_container_group" {
	group_id = "group-ynd95rea"
	operate = "stop"
}

`
