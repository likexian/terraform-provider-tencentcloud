package vod_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcvod "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vod"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_vod_snapshot_template
	resource.AddTestSweepers("tencentcloud_vod_snapshot_template", &resource.Sweeper{
		Name: "tencentcloud_vod_snapshot_template",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta)
			vodService := svcvod.NewVodService(client.GetAPIV3Conn())
			filter := make(map[string]interface{})
			templates, e := vodService.DescribeSnapshotByTimeOffsetTemplatesByFilter(ctx, filter)
			if e != nil {
				return nil
			}
			for _, template := range templates {
				ee := vodService.DeleteSnapshotByTimeOffsetTemplate(ctx, strconv.FormatUint(*template.Definition, 10), uint64(0))
				if ee != nil {
					continue
				}
			}

			spriteTemplates, spriteErr := vodService.DescribeImageSpriteTemplatesByFilter(ctx, filter)
			if spriteErr != nil {
				return nil
			}
			for _, spriteTemplate := range spriteTemplates {
				ee := vodService.DeleteImageSpriteTemplate(ctx, strconv.FormatUint(*spriteTemplate.Definition, 10), uint64(0))
				if ee != nil {
					continue
				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudVodSnapshotByTimeOffsetTemplateResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVodSnapshotByTimeOffsetTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSnapshotByTimeOffsetTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVodSnapshotByTimeOffsetTemplateExists("tencentcloud_vod_snapshot_by_time_offset_template.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "name", "tf-snapshot"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "width", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "height", "128"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "resolution_adaptive", "false"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "format", "png"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "comment", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "fill_type", "white"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_snapshot_by_time_offset_template.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_snapshot_by_time_offset_template.foo", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vod_snapshot_by_time_offset_template.foo", "type"),
				),
			},
			{
				Config: testAccVodSnapshotByTimeOffsetTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "name", "tf-snapshot-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "width", "129"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "height", "129"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "resolution_adaptive", "true"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "format", "jpg"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "comment", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_vod_snapshot_by_time_offset_template.foo", "fill_type", "gauss"),
				),
			},
			{
				ResourceName:      "tencentcloud_vod_snapshot_by_time_offset_template.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckVodSnapshotByTimeOffsetTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	vodService := svcvod.NewVodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vod_snapshot_by_time_offset_template" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		subAppId := helper.StrToInt(idSplit[0])
		definition := idSplit[1]
		filter := map[string]interface{}{
			"definitions": []string{definition},
			"sub_appid":   subAppId,
		}

		templates, err := vodService.DescribeSnapshotByTimeOffsetTemplatesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(templates) == 0 || len(templates) != 1 {
			return nil
		}
		return fmt.Errorf("vod snapshot by time offset template still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckVodSnapshotByTimeOffsetTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vod snapshot by time offset template %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vod snapshot by time offset template id is not set")
		}
		vodService := svcvod.NewVodService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		subAppId := helper.StrToInt(idSplit[0])
		definition := idSplit[1]
		filter := map[string]interface{}{
			"definitions": []string{definition},
			"sub_appid":   subAppId,
		}
		templates, err := vodService.DescribeSnapshotByTimeOffsetTemplatesByFilter(ctx, filter)
		if err != nil {
			return err
		}
		if len(templates) == 0 || len(templates) != 1 {
			return fmt.Errorf("vod snapshot by time offset template doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccVodSnapshotByTimeOffsetTemplate = `
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "sbtot-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  width               = 128
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}
`

const testAccVodSnapshotByTimeOffsetTemplateUpdate = `
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "sbtot-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot-update"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  width               = 129
  height              = 129
  resolution_adaptive = true
  format              = "jpg"
  comment             = "test-update"
  fill_type           = "gauss"
}
`
