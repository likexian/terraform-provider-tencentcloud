package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -timeout=0 -test.run TestAccTencentCloudNeedFixSqlserverGeneralCloudRoInstanceResource_basic -v
func TestAccTencentCloudNeedFixSqlserverGeneralCloudRoInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverGeneralCloudRoInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCloudRoInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverRoInstanceExists("tencentcloud_sqlserver_general_cloud_ro_instance.general_cloud_ro_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_ro_instance.general_cloud_ro_instance", "id"),
				),
			},
			{
				Config: testAccSqlserverGeneralCloudRoInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverRoInstanceExists("tencentcloud_sqlserver_general_cloud_ro_instance.general_cloud_ro_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_ro_instance.general_cloud_ro_instance", "id"),
				),
			},
		},
	})
}

func testAccCheckSqlserverGeneralCloudRoInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_general_cloud_ro_instance" {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		roInstanceId := idSplit[1]

		result, err := service.DescribeSqlserverGeneralCloudRoInstanceById(ctx, roInstanceId)
		if err != nil {
			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound.InstanceNotFound" {
					return nil
				}
			}

			return err
		}

		if result != nil {
			return fmt.Errorf("sqlserver general_cloud_ro_instance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverRoInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}

		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		roInstanceId := idSplit[1]

		result, err := service.DescribeSqlserverGeneralCloudRoInstanceById(ctx, roInstanceId)
		if err != nil {
			return err
		}

		if result != nil {
			return nil
		} else {
			return fmt.Errorf("sqlserver general_cloud_ro_instance %s is not found", rs.Primary.ID)
		}
	}
}

const testAccSqlserverGeneralCloudRoInstance = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example" {
  name                 = "tf_example"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}

resource "tencentcloud_sqlserver_general_cloud_ro_instance" "example" {
  instance_id                      = tencentcloud_sqlserver_general_cloud_instance.example.id
  zone                             = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  read_only_group_type             = 2
  read_only_group_name             = "test-ro-group"
  read_only_group_is_offline_delay = 1
  read_only_group_max_delay_time   = 10
  read_only_group_min_in_group     = 1
  memory                           = 4
  storage                          = 100
  cpu                              = 2
  machine_type                     = "CLOUD_BSSD"
  instance_charge_type             = "POSTPAID"
  subnet_id                        = tencentcloud_subnet.subnet.id
  vpc_id                           = tencentcloud_vpc.vpc.id
  security_group_list              = [tencentcloud_security_group.security_group.id]
  collation                        = "Chinese_PRC_CI_AS"
  time_zone                        = "China Standard Time"
  resource_tags                    = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
`

const testAccSqlserverGeneralCloudRoInstanceUpdate = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example" {
  instance_id                      = tencentcloud_sqlserver_general_cloud_instance.example.id
  zone                             = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  read_only_group_type             = 2
  read_only_group_name             = "test-ro-group"
  read_only_group_is_offline_delay = 1
  read_only_group_max_delay_time   = 10
  read_only_group_min_in_group     = 1
  memory                           = 8
  storage                          = 200
  cpu                              = 4
  machine_type                     = "CLOUD_BSSD"
  instance_charge_type             = "POSTPAID"
  subnet_id                        = tencentcloud_subnet.subnet.id
  vpc_id                           = tencentcloud_vpc.vpc.id
  security_group_list              = [tencentcloud_security_group.security_group.id]
  collation                        = "Chinese_PRC_CI_AS"
  time_zone                        = "China Standard Time"
  resource_tags                    = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
`
