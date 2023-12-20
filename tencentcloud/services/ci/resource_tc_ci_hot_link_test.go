package ci_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localci "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ci"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCiHotLinkResource_basic -v
func TestAccTencentCloudCiHotLinkResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCiHotLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiHotLink,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiHotLinkExists("tencentcloud_ci_hot_link.hot_link"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_hot_link.hot_link", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_hot_link.hot_link", "type", "white"),
					resource.TestCheckResourceAttr("tencentcloud_ci_hot_link.hot_link", "url.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_hot_link.hot_link",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiHotLinkDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_hot_link" {
			continue
		}

		res, err := service.DescribeCiHotLinkById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("ci hot link still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiHotLinkExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci hot link %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}

		result, err := service.DescribeCiHotLinkById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("ci hot link not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiHotLinkVar = `
variable "bucket" {
	default = "` + tcacctest.DefaultCiBucket + `"
  }

`
const testAccCiHotLink = testAccCiHotLinkVar + `

resource "tencentcloud_ci_hot_link" "hot_link" {
	bucket = var.bucket
	url = ["10.0.0.1", "10.0.0.2"]
	type = "white"
}

`
