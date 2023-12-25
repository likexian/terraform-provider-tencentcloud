package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCvmRenewHostResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRenewHost,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_renew_host.renew_host", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_renew_host.renew_host",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmRenewHost = `

resource "tencentcloud_cvm_renew_host" "renew_host" {
  host_ids = 
  host_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_MANUAL_RENEW"

  }
}

`
