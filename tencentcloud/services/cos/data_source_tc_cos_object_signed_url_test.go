package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosObjectSignedUrlDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosObjectSignedUrlDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cos_object_signed_url.cos_object_signed_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cos_object_signed_url.cos_object_signed_url", "signed_url"),
				),
			},
		},
	})
}

const testAccCosObjectSignedUrlDataSource = `
data "tencentcloud_cos_object_signed_url" "cos_object_signed_url" {
  bucket = "keep-test-1308919341"
  path   = "path/to/file"
  headers = {
    Content-Type = "text/plain"
  }
  queries = {
    prefix = "xxx"
  }
}
`
