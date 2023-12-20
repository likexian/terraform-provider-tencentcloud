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
	"github.com/tencentyun/cos-go-sdk-v5"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_ci_media_video_montage_template
	resource.AddTestSweepers("tencentcloud_ci_media_video_montage_template", &resource.Sweeper{
		Name: "tencentcloud_ci_media_video_montage_template",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := localci.NewCiService(client)

			response, _, err := client.UseCiClient(tcacctest.DefaultCiBucket).CI.DescribeMediaTemplate(ctx, &cos.DescribeMediaTemplateOptions{
				Name: "video_montage_template",
			})
			if err != nil {
				return err
			}

			for _, v := range response.TemplateList {
				err := service.DeleteCiMediaTemplateById(ctx, tcacctest.DefaultCiBucket, v.TemplateId)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudCiMediaVideoMontageTemplateResource_basic -v
func TestAccTencentCloudCiMediaVideoMontageTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCiMediaVideoMontageTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaVideoMontageTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaVideoMontageTemplateExists("tencentcloud_ci_media_video_montage_template.media_video_montage_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "bucket", tcacctest.DefaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "name", "video_montage_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "duration", "10.5"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio.0.codec", "aac"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio.0.samplerate", "44100"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio.0.bitrate", "128"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio.0.channels", "4"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio.0.remove", "false"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "video.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "video.0.codec", "H.264"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "video.0.width", "1280"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "video.0.bitrate", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "video.0.fps", "25"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "container.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "container.0.format", "mp4"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio_mix.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio_mix.0.audio_source", "https://"+tcacctest.DefaultCiBucket+".cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio_mix.0.mix_mode", "Once"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_video_montage_template.media_video_montage_template", "audio_mix.0.replace", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_video_montage_template.media_video_montage_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaVideoMontageTemplateDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_video_montage_template" {
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
			return fmt.Errorf("ci media video montage template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaVideoMontageTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := localci.NewCiService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media video montage template %s is not found", re)
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
			return fmt.Errorf("ci media video montage template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaVideoMontageTemplateVar = `
variable "bucket" {
	default = "` + tcacctest.DefaultCiBucket + `"
  }
`

const testAccCiMediaVideoMontageTemplate = testAccCiMediaVideoMontageTemplateVar + `

resource "tencentcloud_ci_media_video_montage_template" "media_video_montage_template" {
	bucket = var.bucket
	name = "video_montage_template"
	duration = "10.5"
	audio {
		codec = "aac"
		samplerate = "44100"
		bitrate = "128"
		channels = "4"
		remove = "false"
	}
	video {
		codec = "H.264"
		width = "1280"
		height = ""
		bitrate = "1000"
		fps = "25"
		crf = ""
		remove = ""
	}
	container {
		format = "mp4"
	}
	audio_mix {
		audio_source = "https://` + tcacctest.DefaultCiBucket + `.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
		mix_mode = "Once"
		replace = "true"
	}
}

`
