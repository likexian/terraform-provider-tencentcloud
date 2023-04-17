package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisRenewInstanceOperationResource_basic -v
func TestAccTencentCloudRedisRenewInstanceOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisRenewInstanceOperation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_renew_instance_operation.renew_instance_operation", "id"),
				),
			},
		},
	})
}

func testAccRedisRenewInstanceOperation() string {
	return testAccRedisInstancePrepaidBasic() + `
resource "tencentcloud_redis_renew_instance_operation" "renew_instance_operation" {
	instance_id = tencentcloud_redis_instance.redis_prepaid_instance_test.id
	period = 1
	modify_pay_mode = "prepaid"
}`
}
