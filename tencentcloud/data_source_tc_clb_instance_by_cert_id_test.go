package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClbInstanceByCertIdDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceByCertIdDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_clb_instance_by_cert_id.instance_by_cert_id")),
			},
		},
	})
}

const testAccClbInstanceByCertIdDataSource = `

data "tencentcloud_clb_instance_by_cert_id" "instance_by_cert_id" {
  cert_ids = 
  }

`
