package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixRedisUpgradeCacheVersionOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisUpgradeCacheVersionOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_redis_upgrade_cache_version_operation.upgrade_cache_version_operation", "id")),
			},
		},
	})
}

const testAccRedisUpgradeCacheVersionOperation = `

resource "tencentcloud_redis_upgrade_cache_version_operation" "upgrade_cache_version_operation" {
  instance_id = "crs-c1nl9rpv"
  current_redis_version = "5.0.0"
  upgrade_redis_version = "5.0.0"
  instance_type_upgrade_now = 1
}

`
