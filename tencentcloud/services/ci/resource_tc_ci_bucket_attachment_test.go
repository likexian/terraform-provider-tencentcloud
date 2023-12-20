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

// go test -i; go test -test.run TestAccTencentCloudCiBucketAttachmentResource_basic -v
func TestAccTencentCloudCiBucketAttachmentResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCiBucketAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiBucketAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiBucketAttachmentExists("tencentcloud_ci_bucket_attachment.bucket_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_bucket_attachment.bucket_attachment", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_bucket_attachment.bucket_attachment", "ci_status", "on"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_bucket_attachment.bucket_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiBucketAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_bucket_attachment" {
			continue
		}

		res, err := service.DescribeCiBucketById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if res != nil {
			status := res.CIStatus
			if res.CIStatus == "on" {
				return fmt.Errorf("ci bucket still exist, Id: %v, status:%s", rs.Primary.ID, status)
			}
		}
	}
	return nil
}

func testAccCheckCiBucketAttachmentExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci bucket %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}

		result, err := service.DescribeCiBucketById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("ci bucket not found, Id: %v", rs.Primary.ID)
		}

		if result != nil {
			status := result.CIStatus
			if result.CIStatus == "off" {
				return fmt.Errorf("ci bucket unbound, Id: %v, status:%s", rs.Primary.ID, status)
			}
		}
		return nil
	}
}

const testAccCiBucketAttachment = `

resource "tencentcloud_ci_bucket_attachment" "bucket_attachment" {
  bucket = "terraform-ci-test-1308919341"
}

`
