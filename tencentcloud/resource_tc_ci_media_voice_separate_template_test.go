package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCiMediaVoiceSeparateTemplateResource_basic -v
func TestAccTencentCloudCiMediaVoiceSeparateTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCiMediaVoiceSeparateTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaVoiceSeparateTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCiMediaVoiceSeparateTemplateExists("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template"),
					resource.TestCheckResourceAttrSet("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "bucket", defaultCiBucket),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "name", "voice_separate_template"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "audio_mode", "IsAudio"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "audio_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "audio_config.0.codec", "aac"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "audio_config.0.samplerate", "44100"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "audio_config.0.bitrate", "128"),
					resource.TestCheckResourceAttr("tencentcloud_ci_media_voice_separate_template.media_voice_separate_template", "audio_config.0.channels", "4"),
				),
			},
			{
				ResourceName:      "tencentcloud_ci_media_voice_separate_template.media_voice_separate_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCiMediaVoiceSeparateTemplateDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ci_media_voice_separate_template" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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
			return fmt.Errorf("ci media video separate template still exist, Id: %v\n", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCiMediaVoiceSeparateTemplateExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := CiService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("ci media video separate template %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf(" id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
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
			return fmt.Errorf("ci media video separate template not found, Id: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCiMediaVoiceSeparateTemplateVar = `
variable "bucket" {
	default = "` + defaultCiBucket + `"
  }
`

const testAccCiMediaVoiceSeparateTemplate = testAccCiMediaVoiceSeparateTemplateVar + `

resource "tencentcloud_ci_media_voice_separate_template" "media_voice_separate_template" {
	bucket = var.bucket
	name = "voice_separate_template"
	audio_mode = "IsAudio"
	audio_config {
		codec = "aac"
		samplerate = "44100"
		bitrate = "128"
		channels = "4"
	}
  }

`
