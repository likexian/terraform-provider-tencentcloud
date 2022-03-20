package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testSqlserverInstanceResourceName = "tencentcloud_sqlserver_instance"
var testSqlserverInstanceResourceKey = testSqlserverInstanceResourceName + ".test"

func init() {
	resource.AddTestSweepers("tencentcloud_sqlserver_instance", &resource.Sweeper{
		Name: "tencentcloud_sqlserver_instance",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := SqlserverService{client: client}
			instances, err := service.DescribeSqlserverInstances(ctx, "", -1, defaultVpcId, defaultSubnetId, 1)

			if err != nil {
				return err
			}

			err = batchDeleteSQLServerInstances(ctx, service, instances)

			if err != nil {
				return err
			}

			return nil
		},
	})
}

func batchDeleteSQLServerInstances(ctx context.Context, service SqlserverService, instances []*sqlserver.DBInstance) error {
	wg := sync.WaitGroup{}

	wg.Add(len(instances))
	for i := range instances {
		go func(i int) {
			defer wg.Done()
			id := *instances[i].InstanceId
			name := *instances[i].Name
			if strings.HasPrefix(name, "preset_sqlserver") {
				return
			}

			var outErr, inErr error
			var has bool

			outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, inErr = service.DescribeSqlserverInstanceById(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})

			if outErr != nil {
				return
			}

			if !has {
				return
			}

			//terminate sql instance
			outErr = service.TerminateSqlserverInstance(ctx, id)

			if outErr != nil {
				outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					inErr = service.TerminateSqlserverInstance(ctx, id)
					if inErr != nil {
						return retryError(inErr)
					}
					return nil
				})
			}

			if outErr != nil {
				return
			}

			outErr = service.DeleteSqlserverInstance(ctx, id)

			if outErr != nil {
				outErr = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					inErr = service.DeleteSqlserverInstance(ctx, id)
					if inErr != nil {
						return retryError(inErr)
					}
					return nil
				})
			}

			if outErr != nil {
				return
			}

			outErr = service.RecycleDBInstance(ctx, id)
			if outErr != nil {
				return
			}

			outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				_, has, inErr := service.DescribeSqlserverInstanceById(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				if has {
					inErr = fmt.Errorf("delete SQL Server instance %s fail, instance still exists from SDK DescribeSqlserverInstanceById", id)
					return resource.RetryableError(inErr)
				}
				return nil
			})

			if outErr != nil {
				return
			}
		}(i)
	}

	wg.Wait()
	return nil
}

func TestAccTencentCloudSqlserverInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "name", "tf_sqlserver_instance"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "memory", "2"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "maintenance_time_span", "3"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "storage", "10"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vip"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vport"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "status"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "tags.test", "test"),
				),
			},
			{
				ResourceName:            testSqlserverInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"multi_zones"},
			},

			{
				Config: testAccSqlserverInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "name", "tf_sqlserver_instance_update"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "maintenance_time_span", "4"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "storage", "20"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "project_id", defaultProjectId),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vip"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vport"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "status"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "security_groups.#", "0"),
					//resource.TestCheckNoResourceAttr(testSqlserverInstanceResourceKey, "tags.test"),
					//resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "tags.abc", "abc"),
				),
			},
		},
	})
}

func TestAccTencentCloudSqlserverInstanceMultiClusterResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverInstanceMultiCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "name", "tf_sqlserver_instance_multi"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "memory", "2"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "maintenance_time_span", "3"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "storage", "10"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "ha_type", "CLUSTER"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vip"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "vport"),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "status"),
					resource.TestCheckResourceAttr(testSqlserverInstanceResourceKey, "tags.test", "test"),
				),
			},
			{
				ResourceName:            testSqlserverInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"multi_zones"},
			},
		},
	})
}

func testAccCheckSqlserverInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testSqlserverInstanceResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeSqlserverInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		} else {
			return fmt.Errorf("delete SQL Server instance %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSqlserverInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := SqlserverService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribeSqlserverInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribeSqlserverInstanceById(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("SQL Server instance %s is not found", rs.Primary.ID)
		}
	}
}

const testAccSqlserverInstanceBasic = `
data "tencentcloud_availability_zones_by_product" "zone" {
  product = "sqlserver"
}

locals {
  az = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
  az1 = data.tencentcloud_availability_zones_by_product.zone.zones[1].name
}
`

const testAccSqlserverInstance string = testAccSqlserverInstanceBasic + `
resource "tencentcloud_sqlserver_instance" "test" {
  name                          = "tf_sqlserver_instance"
  availability_zone             = local.az1
  charge_type                   = "POSTPAID_BY_HOUR"
  vpc_id                        = "` + defaultVpcId + `"
  subnet_id                     = "` + defaultSubnetId + `"
  project_id                    = 0
  memory                        = 2
  storage                       = 10
  maintenance_week_set          = [1,2,3]
  maintenance_start_time        = "09:00"
  maintenance_time_span         = 3
  security_groups               = ["` + defaultSecurityGroup + `"]

  tags = {
    "test"                      = "test"
  }
}
`

const testAccSqlserverInstanceUpdate string = testAccSqlserverInstanceBasic + `
resource "tencentcloud_sqlserver_instance" "test" {
  name                      = "tf_sqlserver_instance_update"
  availability_zone         = local.az1
  charge_type               = "POSTPAID_BY_HOUR"
  vpc_id                    = "` + defaultVpcId + `"
  subnet_id                 = "` + defaultSubnetId + `"
  project_id                = ` + defaultProjectId + `
  memory                    = 4
  storage                   = 20
  maintenance_week_set      = [2,3,4]
  maintenance_start_time    = "08:00"
  maintenance_time_span     = 4

  tags = {
    abc                     = "abc"
  }
}
`

const testAccSqlserverInstanceMultiCluster string = testAccSqlserverInstanceBasic + `
resource "tencentcloud_sqlserver_instance" "test" {
  name                          = "tf_sqlserver_instance_multi"
  availability_zone             = local.az1
  charge_type                   = "POSTPAID_BY_HOUR"
  engine_version                = "2017"
  vpc_id                        = "` + defaultVpcId + `"
  subnet_id                     = "` + defaultSubnetId + `"
  project_id                    = 0
  memory                        = 2
  storage                       = 10
  multi_zones                   = true
  ha_type                       = "CLUSTER"
  maintenance_week_set          = [1,2,3]
  maintenance_start_time        = "09:00"
  maintenance_time_span         = 3
  tags = {
    "test"                      = "test"
  }
}
`
