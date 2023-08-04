package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testSqlserverAccountResourceName = "tencentcloud_sqlserver_account"
var testSqlserverAccountResourceKey = testSqlserverAccountResourceName + ".test"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_account
	resource.AddTestSweepers("tencentcloud_sqlserver_account", &resource.Sweeper{
		Name: "tencentcloud_sqlserver_account",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn

			service := SqlserverService{client}

			db, err := service.DescribeSqlserverInstances(ctx, "", defaultSQLServerName, -1, "", "", -1)

			if err != nil {
				return err
			}

			if len(db) == 0 {
				return fmt.Errorf("%s not exists", defaultSQLServerName)
			}

			instanceId := *db[0].InstanceId

			accounts, _ := service.DescribeSqlserverAccounts(ctx, instanceId)

			for i := range accounts {
				account := accounts[i]
				name := *account.Name
				created, err := time.Parse("2006-01-02 15:04:05", *account.CreateTime)
				if err != nil {
					created = time.Time{}
				}
				if isResourcePersist(name, &created) {
					continue
				}
				err = service.DeleteSqlserverAccount(ctx, instanceId, name)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudSqlserverAccountResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountExists(testSqlserverAccountResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "name", "tf_sqlserver_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "password", "testt123"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "update_time"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "is_admin", "false"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "status"),
				),
			},
			{
				ResourceName:            testSqlserverAccountResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "is_admin"},
			},

			{
				Config: testAccSqlserverAccountUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverAccountExists(testSqlserverAccountResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "name", "tf_sqlserver_account"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "password", "test1233"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "remark", "testt"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "update_time"),
					resource.TestCheckResourceAttr(testSqlserverAccountResourceKey, "is_admin", "false"),
					resource.TestCheckResourceAttrSet(testSqlserverAccountResourceKey, "status"),
				),
			},
		},
	})
}

func testAccCheckSqlserverAccountDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testSqlserverAccountResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, FILED_SP)
		if len(idStrs) != 2 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		name := idStrs[1]

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeSqlserverAccountById(ctx, instanceId, name)

		if err != nil {
			return err
		}

		if !has {
			return nil
		} else {
			return fmt.Errorf("delete SQL Server account %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		id := rs.Primary.ID
		idStrs := strings.Split(id, FILED_SP)
		if len(idStrs) != 2 {
			return fmt.Errorf("invalid SQL server account id %s", id)
		}
		instanceId := idStrs[0]
		name := idStrs[1]

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeSqlserverAccountById(ctx, instanceId, name)
		if err != nil {
			_, has, err = service.DescribeSqlserverAccountById(ctx, instanceId, name)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("SQL Server account %s is not found", rs.Primary.ID)
		}
	}
}

const testAccSqlserverAccount string = CommonPresetSQLServer + `
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

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_account"
  password    = "Qwer@234"
  remark      = "test-remark"
}
`

const testAccSqlserverAccountUpdate string = CommonPresetSQLServer + `
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

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_account" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
  name        = "tf_example_account"
  password    = "Qwer@234Update"
  remark      = "test-remark-update"
}
`
