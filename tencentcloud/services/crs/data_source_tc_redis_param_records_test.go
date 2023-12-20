package crs_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudRedisParamRecordsDataSource_basic -v
func TestAccTencentCloudRedisParamRecordsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisParamRecordsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_redis_param_records.param_records"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_param_records.param_records", "instance_param_history.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_param_records.param_records", "instance_param_history.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_param_records.param_records", "instance_param_history.0.new_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_param_records.param_records", "instance_param_history.0.param_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_param_records.param_records", "instance_param_history.0.pre_value"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_redis_param_records.param_records", "instance_param_history.0.status"),
				),
			},
		},
	})
}

const testAccRedisParamRecordsDataSource = `

data "tencentcloud_redis_param_records" "param_records" {
	instance_id = "crs-jf4ico4v"
}

`
