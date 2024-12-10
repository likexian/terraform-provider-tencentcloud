package postgresql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudPostgresqlInstanceNetworkAccessResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccPostgresqlInstanceNetworkAccess,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_network_access.postgresql_instance_network_access", "id")),
		}, {
			ResourceName:      "tencentcloud_postgresql_instance_network_access.postgresql_instance_network_access",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccPostgresqlInstanceNetworkAccess = `

resource "tencentcloud_postgresql_instance_network_access" "postgresql_instance_network_access" {
}
`
