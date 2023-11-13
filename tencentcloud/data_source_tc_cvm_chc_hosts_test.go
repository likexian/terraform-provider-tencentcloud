package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmChcHostsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcHostsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_chc_hosts.chc_hosts")),
			},
		},
	})
}

const testAccCvmChcHostsDataSource = `

data "tencentcloud_cvm_chc_hosts" "chc_hosts" {
  chc_ids = 
  filters {
		name = ""
		values = 

  }
  }

`
