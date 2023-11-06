package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodDownloadSnapshotOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDownloadSnapshotOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_download_snapshot_operation.download_snapshot", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_download_snapshot_operation.download_snapshot", "snapshot_id", "87910DFF"),
					resource.TestCheckResourceAttrSet("tencentcloud_dnspod_download_snapshot_operation.download_snapshot", "cos_url"),
				),
			},
		},
	})
}

const testAccDnspodDownloadSnapshotOperation = `
resource "tencentcloud_dnspod_download_snapshot_operation" "download_snapshot" {
  domain = "iac-tf.cloud"
  snapshot_id = "87910DFF"
}
`
