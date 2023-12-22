package kms_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudKmsKeyDataSource(t *testing.T) {
	t.Parallel()
	dataSourceName := "data.tencentcloud_kms_keys.test"
	rName := fmt.Sprintf("tf-testacc-kms-key-%s", acctest.RandString(13))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKmsKeyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.key_state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.key_usage"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.creator_uin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.key_rotation_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.owner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.next_rotate_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.origin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "key_list.0.valid_to"),
				),
			},
		},
	})
}

func testAccDataSourceKmsKeyConfig(rName string) string {
	return fmt.Sprintf(`
resource "tencentcloud_kms_key" "test" {
  	alias = %[1]q
	description = %[1]q
  	is_enabled = false
	key_rotation_enabled = true
}
data "tencentcloud_kms_keys" "test" {
  search_key_alias = tencentcloud_kms_key.test.alias
}
`, rName)
}
