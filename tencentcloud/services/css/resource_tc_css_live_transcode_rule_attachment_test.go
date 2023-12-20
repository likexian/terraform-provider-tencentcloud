package css_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccss "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/css"

	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_css_live_transcode_rule_attachment", &resource.Sweeper{
		Name: "tencentcloud_css_live_transcode_rule_attachment",
		F:    testSweepCssLiveTranscodeRuleAttachment,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_css_live_transcode_rule_attachment
func testSweepCssLiveTranscodeRuleAttachment(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	cssService := svccss.NewCssService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

	temps, err := cssService.DescribeCssLiveTranscodeTemplates(ctx)
	if err != nil {
		return err
	}
	if temps == nil {
		return fmt.Errorf("live transcode template attachment not exists.")
	}

	for _, v := range temps {
		delName := v.TemplateName
		delId := v.TemplateId

		if strings.HasPrefix(*delName, tcacctest.DefaultCSSPrefix) {
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				err := cssService.DeleteCssLiveTranscodeTemplateById(ctx, delId)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] sweeper live transcode template attachment %s:%v failed! reason:[%s]", *delName, *delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudCssLiveTranscodeRuleAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	baseTime := time.Now().UTC().Add(10 * time.Hour)
	startTime := baseTime.Format(time.RFC3339)
	endTime := baseTime.Add(1 * time.Hour).Format(time.RFC3339)
	// startTimeNew := baseTime.Add(30 * time.Minute).Format(time.RFC3339)
	// endTimeNew := baseTime.Add(2 * time.Hour).Format(time.RFC3339)
	liveUrl := "rtmp://5000.liveplay.myqcloud.com/live/stream1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCssLiveTranscodeRuleAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssLiveTranscodeRuleAttachment, tcacctest.DefaultCSSLiveType, liveUrl, tcacctest.DefaultCSSDomainName, tcacctest.DefaultCSSAppName, tcacctest.DefaultCSSStreamName, startTime, endTime, tcacctest.DefaultCSSOperator),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssLiveTranscodeRuleAttachmentExists("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment", "domain_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment", "app_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment", "stream_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment", "template_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment", "update_time"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_live_transcode_rule_attachment.live_transcode_rule_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCssLiveTranscodeRuleAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_css_live_transcode_rule_attachment" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		domainName := idSplit[0]
		templateId := idSplit[3]

		temp, err := cssService.DescribeCssLiveTranscodeRuleAttachment(ctx, helper.String(domainName), helper.String(templateId))
		if err != nil {
			return nil
		}

		if temp != nil {
			return fmt.Errorf("css live transcode template attachment still exist, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCssLiveTranscodeRuleAttachmentExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("css live transcode template attachment %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("css live transcode template attachment id is not set")
		}

		cssService := svccss.NewCssService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		domainName := idSplit[0]
		templateId := idSplit[3]

		temp, err := cssService.DescribeCssLiveTranscodeRuleAttachment(ctx, helper.String(domainName), helper.String(templateId))
		if err != nil {
			return err
		}

		if temp == nil {
			return fmt.Errorf("css live transcode template attachment not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCssPullstreamtask = `
  resource "tencentcloud_css_pull_stream_task" "task" {
	source_type = "%s"
	source_urls = ["%s"]
	domain_name = "%s"
	app_name = "%s"
	stream_name = "%s"
	start_time = "%s"
	end_time = "%s"
	operator = "%s"
	comment = "This is a e2e test case."
  }
`
const testAccCssLiveTranscodeTemp = `
resource "tencentcloud_css_live_transcode_template" "temp" {
  template_name = "tf1080p"
  acodec = "aac"
  video_bitrate = 100
  vcodec = "origin"
  description = "This_is_a_tf_test_temp."
  need_video = 1
  need_audio = 1
}
`

const testAccCssLiveTranscodeRuleAttachment = testAccCssPullstreamtask + testAccCssLiveTranscodeTemp + `
resource "tencentcloud_css_live_transcode_rule_attachment" "live_transcode_rule_attachment" {
  domain_name = tencentcloud_css_pull_stream_task.task.domain_name
  app_name = tencentcloud_css_pull_stream_task.task.app_name
  stream_name = tencentcloud_css_pull_stream_task.task.stream_name
  template_id = tencentcloud_css_live_transcode_template.temp.id
}

`
