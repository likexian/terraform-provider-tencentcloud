package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMongodbInstanceBackupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceBackup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance_backup.instance_backup", "id")),
			},
			{
				ResourceName:      "tencentcloud_mongodb_instance_backup.instance_backup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMongodbInstanceBackup = `

resource "tencentcloud_mongodb_instance_backup" "instance_backup" {
  instance_id = "cmgo-9d0p6umb"
  backup_method = 0
  backup_remark = "my backup"
}

`
