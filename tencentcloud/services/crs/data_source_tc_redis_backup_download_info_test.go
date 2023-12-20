package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisBackupDownloadInfoDataSource_basic -v
func TestAccTencentCloudRedisBackupDownloadInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisBackupDownloadInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_redis_backup_download_info.backup_download_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup_download_info.backup_download_info", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_redis_backup_download_info.backup_download_info", "backup_infos.#", "1"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup_download_info.backup_download_info", "backup_infos.0.download_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup_download_info.backup_download_info", "backup_infos.0.file_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup_download_info.backup_download_info", "backup_infos.0.file_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_backup_download_info.backup_download_info", "backup_infos.0.inner_download_url"),
				),
			},
		},
	})
}

const testAccRedisBackupDownloadInfoDataSource = testAccRedisBackupDataSource + `

data "tencentcloud_redis_backup_download_info" "backup_download_info" {
	instance_id = var.instance_id
	backup_id = data.tencentcloud_redis_backup.backup.backup_set.0.backup_id
	# limit_type = "NoLimit"
	# vpc_comparison_symbol = "In"
	# ip_comparison_symbol = "In"
	# limit_vpc {
	  # 	region = "ap-guangzhou"
	  # 	vpc_list = [""]
	# }
	# limit_ip = [""] 
}

`
