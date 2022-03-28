package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testPostgresqlInstanceResourceName = "tencentcloud_postgresql_instance"
var testPostgresqlInstanceResourceKey = testPostgresqlInstanceResourceName + ".test"

func init() {
	resource.AddTestSweepers(testPostgresqlInstanceResourceName, &resource.Sweeper{
		Name: testPostgresqlInstanceResourceName,
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			postgresqlService := PostgresqlService{client: client}
			vpcService := VpcService{client: client}

			instances, err := postgresqlService.DescribePostgresqlInstances(ctx, nil)
			if err != nil {
				return err
			}

			var vpcs []string

			for _, v := range instances {
				id := *v.DBInstanceId
				name := *v.DBInstanceName
				vpcId := *v.VpcId
				if strings.HasPrefix(name, keepResource) || strings.HasPrefix(name, defaultResource) {
					continue
				}
				err := postgresqlService.IsolatePostgresqlInstance(ctx, id)
				if err != nil {
					continue
				}
				err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
					instance, has, err := postgresqlService.DescribePostgresqlInstanceById(ctx, id)
					if err != nil {
						return retryError(err)
					}
					if !has {
						return resource.NonRetryableError(fmt.Errorf("instance %s removed", id))
					}
					if *instance.DBInstanceStatus != "isolated" {
						return resource.RetryableError(fmt.Errorf("waiting for instance isolated, now is %s", *instance.DBInstanceStatus))
					}
					return nil
				})
				if err != nil {
					continue
				}
				err = postgresqlService.DeletePostgresqlInstance(ctx, id)
				if err != nil {
					continue
				}
				vpcs = append(vpcs, vpcId)
			}

			for _, v := range vpcs {
				_ = vpcService.DeleteVpc(ctx, v)
			}

			return nil
		},
	})
}

func TestAccTencentCloudPostgresqlInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPostgresqlInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "name", "tf_postsql_instance"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "storage", "100"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "public_access_switch", "false"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "root_password", "t1qaA2k1wgvfa3?ZZZ"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_port"),
					//resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "tags.tf", "test"),
				),
			},
			{
				ResourceName:            testPostgresqlInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "spec_code", "public_access_switch", "charset"},
			},

			{
				Config: testAccPostgresqlInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "name", "tf_postsql_instance_update"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "storage", "250"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "create_time"),
					// FIXME After PGSQL fixed can reopen case
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "public_access_switch", "true"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "root_password", "t1qaA2k1wgvfa3?ZZZ"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "public_access_host"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "public_access_port"),
					//resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "tags.tf", "teest"),
				),
			},
		},
	})
}

func TestAccTencentCloudPostgresqlMAZInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPostgresqlInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlMAZInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					// SDK 1.0 cannot provide set test expected "db_node_set.*.role" , "Primary"
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "db_node_set.#", "2"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "availability_zone", "ap-guangzhou-6"),
				),
			},
			{
				ResourceName:            testPostgresqlInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "spec_code", "public_access_switch", "charset"},
			},

			{
				Config: testAccPostgresqlMAZInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "db_node_set.#", "2"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "availability_zone", "ap-guangzhou-7"),
				),
			},
		},
	})
}

func testAccCheckPostgresqlInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testPostgresqlInstanceResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)

		if !has {
			return nil
		} else {
			if err != nil {
				return err
			}
			return fmt.Errorf("delete postgresql instance %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckPostgresqlInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("postgresql instance %s is not found", rs.Primary.ID)
		}
	}
}

const testAccPostgresqlInstanceBasic = `
data "tencentcloud_availability_zones_by_product" "zone" {
  product = "postgres"
}
`

const testAccPostgresqlInstance string = testAccPostgresqlInstanceBasic + `
resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_instance"
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id                   = "` + defaultVpcId + `"
  subnet_id = "` + defaultSubnetId + `"
  engine_version		= "10.4"
  root_password                 = "t1qaA2k1wgvfa3?ZZZ"
  charset = "LATIN1"
  project_id = 0
  memory = 4
  storage = 100

	tags = {
		tf = "test"
	}
}
`

const testAccPostgresqlInstanceUpdate string = testAccPostgresqlInstanceBasic + `
resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_instance_update"
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id                   = "` + defaultVpcId + `"
  subnet_id = "` + defaultSubnetId + `"
  engine_version		= "10.4"
  root_password                 = "t1qaA2k1wgvfa3?ZZZ"
  charset = "LATIN1"
  project_id = 0
  public_access_switch = true
  memory = 4
  storage = 250

	tags = {
		tf = "teest"
	}
}
`

const testAccPostgresqlMAZInstance string = `
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/24"
  name       = "test-pg-vpc"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.0.0/24"
  name              = "pg-sub1"
  vpc_id            = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_maz_instance"
  availability_zone = "ap-guangzhou-6"
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version		= "10.4"
  root_password                 = "t1qaA2k1wgvfa3?ZZZ"
  charset = "LATIN1"
  memory = 4
  storage = 100
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }
  db_node_set {
    zone = "ap-guangzhou-7"
  }
}
`

const testAccPostgresqlMAZInstanceUpdate string = `
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/24"
  name       = "test-pg-vpc"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.0.0/24"
  name              = "pg-sub1"
  vpc_id            = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_maz_instance"
  availability_zone = "ap-guangzhou-6"
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version		= "10.4"
  root_password                 = "t1qaA2k1wgvfa3?ZZZ"
  charset = "LATIN1"
  memory = 4
  storage = 250
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }
  db_node_set {
    zone = "ap-guangzhou-6"
  }
}
`
