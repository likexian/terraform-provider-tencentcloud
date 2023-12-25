package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudVpcACL_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudVpcACLInstances,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_acls.default"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc_acls.default", "name", "test_acl"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudVpcACLInstances = `
data "tencentcloud_vpc_instances" "test" {
	is_default = true
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id  = data.tencentcloud_vpc_instances.test.instance_list.0.vpc_id
    name  	= "test_acl"
	ingress = [
		"ACCEPT#192.168.1.0/24#80#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#80#TCP",
	]
}  

data "tencentcloud_vpc_acls" "default" {
	name = "test_acl"
}
`
