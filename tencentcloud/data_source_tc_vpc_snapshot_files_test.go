package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcSnapshotFilesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSnapshotFilesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_snapshot_files.snapshot_files")),
			},
		},
	})
}

const testAccVpcSnapshotFilesDataSource = `

data "tencentcloud_vpc_snapshot_files" "snapshot_files" {
  business_type = "securitygroup"
  instance_id   = "sg-902tl7t7"
  start_date    = "2022-10-10 00:00:00"
  end_date      = "2023-10-30 19:00:00"
}

`
