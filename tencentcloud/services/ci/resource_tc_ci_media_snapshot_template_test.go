package ci_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localci "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ci"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCiMediaSnapshotTemplateResource_basic -v
func TestAccTencentCloudCiMediaSnapshotTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCiMediaSnapshotTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaSnapshotTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaSnapshotTemplateTemplateExists("tencentcloud_ci_media_snapshot_template.media_snapshot_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "bucket", tcacctest.DefaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "name", "snapshot_template_test"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "snapshot.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "snapshot.0.count", "10"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "snapshot.0.snapshot_out_mode", "SnapshotAndSprite"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "snapshot.0.sprite_snapshot_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "snapshot.0.sprite_snapshot_config.0.color", "White"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "snapshot.0.sprite_snapshot_config.0.columns", "10"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_snapshot_template.media_snapshot_template", "snapshot.0.sprite_snapshot_config.0.lines", "10"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_snapshot_template.media_snapshot_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaSnapshotTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_snapshot_template" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		bucket := idSplit[0]
		templateId := idSplit[1]

		res, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
		if err != nil {
			return err
		}

		if res != nil {
			return fmt.Errorf("ci media snapshot template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaSnapshotTemplateTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media snapshot template %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		bucket := idSplit[0]
		templateId := idSplit[1]

		result, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
		if err != nil {
			return err
		}

		if result == nil {
			return fmt.Errorf("ci media snapshot template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaSnapshotTemplateVar = `
variable "bucket" {
	default = "` + tcacctest.DefaultCiBucket + `"
  }

`

const testAccCiMediaSnapshotTemplate = testAccCiMediaSnapshotTemplateVar + `

resource "tencentcloud_ci_media_snapshot_template" "media_snapshot_template" {
    bucket = var.bucket
  	name = "snapshot_template_test"
  	snapshot {
      count = "10"
      snapshot_out_mode = "SnapshotAndSprite"
      sprite_snapshot_config {
        color = "White"
        columns = "10"
        lines = "10"
        margin = "10"
        padding = "10"
      }
  	}
}

`
