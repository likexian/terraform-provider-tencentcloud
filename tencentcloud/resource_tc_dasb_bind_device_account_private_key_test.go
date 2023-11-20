package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbBindDeviceAccountPrivateKeyResource_basic -v
func TestAccTencentCloudNeedFixDasbBindDeviceAccountPrivateKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbBindDeviceAccountPrivateKey,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_bind_device_account_private_key.bind_device_account_private_key", "id"),
				),
			},
		},
	})
}

const testAccDasbBindDeviceAccountPrivateKey = `
resource "tencentcloud_dasb_bind_device_account_private_key" "example" {
  id                   = 16
  private_key          = "MIICXAIBAAKBgQCqGKukO1De7zhZj6+H0qtjTkVxwTCpvKe4eCZ0FPqri0cb2JZfXJ/DgYSF6vUpwmJG8wVQZKjeGcjDOL5UlsuusFncCzWBQ7RKNUSesmQRMSGkVb1/3j+skZ6UtW+5u09lHNsj6tQ51s1SPrCBkedbNf0Tp0GbMJDyR4e9T04ZZwIDAQABAoGAFijko56+qGyN8M0RVyaRAXz++xTqHBLh"
  private_key_password = "TerraformPassword"
}
`
