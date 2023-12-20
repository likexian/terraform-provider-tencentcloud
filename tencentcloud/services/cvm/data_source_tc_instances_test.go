package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDataSourceInstancesBase(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceInstancesBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.foo", "instance_list.0.instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.foo", "instance_list.0.instance_name", tcacctest.DefaultInsName),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.foo", "instance_list.0.instance_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.foo", "instance_list.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.foo", "instance_list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.foo", "instance_list.0.availability_zone"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.foo", "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.foo", "instance_list.0.system_disk_type"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceInstancesBase = tcacctest.InstanceCommonTestCase + `
data "tencentcloud_instances" "foo" {
  instance_id = tencentcloud_instance.default.id
  instance_name = tencentcloud_instance.default.instance_name
}
`
