package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCcnCrossBorderComplianceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCrossBorderComplianceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ccn_cross_border_compliance.cross_border_compliance")),
			},
		},
	})
}

const testAccVpcCrossBorderComplianceDataSource = `

data "tencentcloud_ccn_cross_border_compliance" "cross_border_compliance" {
  service_provider = "UNICOM"
  compliance_id = 10002
  email = "test@tencent.com"
  service_start_date = "2020-07-29"
  service_end_date = "2021-07-29"
  state = "APPROVED"
}

`
