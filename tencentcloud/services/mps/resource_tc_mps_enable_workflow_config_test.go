package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsEnableWorkflowConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsEnableWorkflowConfig_enable,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_workflow_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_workflow_config.config", "workflow_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_enable_workflow_config.config", "enabled", "true"),
				),
			},
			{
				Config: testAccMpsEnableWorkflowConfig_disable,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_workflow_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_workflow_config.config", "workflow_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_enable_workflow_config.config", "enabled", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_mps_enable_workflow_config.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsWorkflow_basic = `

resource "tencentcloud_mps_workflow" "example" {
  output_dir    = "/"
  task_priority = 0
  workflow_name = "tf-workflow-enable-config"

  media_process_task {
    adaptive_dynamic_streaming_task_set {
      definition             = 12
      output_object_path     = "/out"
      segment_object_name    = "/out"
      sub_stream_object_name = "/out/out/"

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }

    snapshot_by_time_offset_task_set {
      definition          = 10
      ext_time_offset_set = [
        "1s",
      ]
      output_object_path  = "/snapshot/"
      time_offset_set     = []

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }

    animated_graphic_task_set {
      definition         = 20000
      end_time_offset    = 0
      output_object_path = "/test/"
      start_time_offset  = 0

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }
  }

  ai_analysis_task {
    definition = 20
  }

  ai_content_review_task {
    definition = 20
  }

  ai_recognition_task {
    definition = 20
  }

  output_storage {
    type = "COS"

    cos_output_storage {
      bucket = "cos-lock-1308919341"
      region = "ap-guangzhou"
    }
  }

  trigger {
    type = "CosFileUpload"

    cos_file_upload_trigger {
      bucket = "cos-lock-1308919341"
      dir    = "/"
      region = "ap-guangzhou"
    }
  }
}

`

const testAccMpsEnableWorkflowConfig_enable = testAccMpsWorkflow_basic + `

resource "tencentcloud_mps_enable_workflow_config" "config" {
  workflow_id = tencentcloud_mps_workflow.example.id
  enabled = true
}

`

const testAccMpsEnableWorkflowConfig_disable = testAccMpsWorkflow_basic + `

resource "tencentcloud_mps_enable_workflow_config" "config" {
  workflow_id = tencentcloud_mps_workflow.example.id
  enabled = false
}

`
