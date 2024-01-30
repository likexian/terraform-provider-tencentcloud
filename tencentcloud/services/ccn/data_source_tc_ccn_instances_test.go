package ccn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudCcnV3InstancesBasic(t *testing.T) {
	t.Parallel()
	keyName := "data.tencentcloud_ccn_instances.name_instances"
	keyId := "data.tencentcloud_ccn_instances.id_instances"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudCcnInstances,

				Check: resource.ComposeTestCheckFunc(
					//id filter
					tcacctest.AccCheckTencentCloudDataSourceID(keyId),
					resource.TestCheckResourceAttr(keyId, "instance_list.#", "1"),
					resource.TestCheckResourceAttr(keyId, "instance_list.0.name", "ci-temp-test-ccn"),
					resource.TestCheckResourceAttr(keyId, "instance_list.0.description", "ci-temp-test-ccn-des"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.ccn_id"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.qos"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.state"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.charge_type"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.bandwidth_limit_type"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.attachment_list.#"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.create_time"),

					//name filter ,Every VPC with a "guagua_vpc_instance_test" name will be found
					tcacctest.AccCheckTencentCloudDataSourceID(keyName),
					resource.TestCheckResourceAttrSet(keyName, "instance_list.#"),
					resource.TestCheckResourceAttrSet(keyName, "instance_list.0.name"),
					resource.TestCheckResourceAttrSet(keyName, "instance_list.0.description"),
					resource.TestCheckResourceAttrSet(keyName, "instance_list.0.ccn_id"),
					resource.TestCheckResourceAttrSet(keyName, "instance_list.0.qos"),
					resource.TestCheckResourceAttrSet(keyName, "instance_list.0.state"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.charge_type"),
					resource.TestCheckResourceAttrSet(keyId, "instance_list.0.bandwidth_limit_type"),
					resource.TestCheckResourceAttrSet(keyName, "instance_list.0.attachment_list.#"),
					resource.TestCheckResourceAttrSet(keyName, "instance_list.0.create_time"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudCcnInstances = `
resource tencentcloud_ccn main {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
  charge_type = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}

data tencentcloud_ccn_instances id_instances {
  ccn_id = tencentcloud_ccn.main.id
}

data tencentcloud_ccn_instances name_instances {
  name = tencentcloud_ccn.main.name
}

`
