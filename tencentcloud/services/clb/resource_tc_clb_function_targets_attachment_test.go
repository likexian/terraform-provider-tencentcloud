package clb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbFunctionTargetsAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbFunctionTargetsAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clb_function_targets_attachment.function_targets", "id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_function_targets_attachment.function_targets", "function_targets.0.weight", "10"),
				),
			},
			{
				Config: testAccClbFunctionTargetsAttachmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_clb_function_targets_attachment.function_targets", "function_targets.0.weight", "20"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_function_targets_attachment.function_targets",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbFunctionTargetsAttachment = `

resource "tencentcloud_clb_function_targets_attachment" "function_targets" {
  domain           = "xxx.com"
  listener_id      = "lbl-nonkgvc2"
  load_balancer_id = "lb-5dnrkgry"
  url              = "/"

  function_targets {
    weight = 10

    function {
      function_name           = "keep-1676351130"
      function_namespace      = "default"
      function_qualifier      = "$LATEST"
      function_qualifier_type = "VERSION"
    }
  }
}

`

const testAccClbFunctionTargetsAttachmentUpdate = `

resource "tencentcloud_clb_function_targets_attachment" "function_targets" {
  domain           = "xxx.com"
  listener_id      = "lbl-nonkgvc2"
  load_balancer_id = "lb-5dnrkgry"
  url              = "/"

  function_targets {
    weight = 20

    function {
      function_name           = "keep-1676351130"
      function_namespace      = "default"
      function_qualifier      = "$LATEST"
      function_qualifier_type = "VERSION"
    }
  }
}

`
