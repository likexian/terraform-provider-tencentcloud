package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosObjectRestoreOperationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosObjectRestoreOperation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_object_restore_operation.object_restore", "id"),
				),
			},
		},
	})
}

const testAccCosObjectRestoreOperation = `
resource "tencentcloud_cos_object_restore_operation" "object_restore" {
    bucket = "keep-test-1308919341"
    key = "test-restore.txt"
    tier = "Expedited"
    days = 2
}
`
