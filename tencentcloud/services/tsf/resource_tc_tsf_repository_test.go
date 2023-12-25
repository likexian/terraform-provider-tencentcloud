package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTsfRepositoryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfRepository,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_repository.repository", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_repository.repository",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfRepository = `

resource "tencentcloud_tsf_repository" "repository" {
    repository_name = ""
  repository_type = ""
  bucket_name = ""
  bucket_region = ""
  directory = ""
  repository_desc = ""
    }

`
