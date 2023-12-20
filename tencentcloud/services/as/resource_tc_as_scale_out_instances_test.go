package as_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixAsScaleOutInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsScaleOutInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_scale_out_instances.scale_out_instances", "id")),
			},
			{
				ResourceName:      "tencentcloud_as_scale_out_instances.scale_out_instances",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAsScaleOutInstances = `

resource "tencentcloud_as_scale_out_instances" "scale_out_instances" {
  auto_scaling_group_id = "asg-519acdug"
  scale_out_number = 1
}

`
