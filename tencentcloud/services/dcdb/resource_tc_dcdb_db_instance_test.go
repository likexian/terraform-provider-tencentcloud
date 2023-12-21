package dcdb_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dcdb"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dcdb_db_instance", &resource.Sweeper{
		Name: "tencentcloud_dcdb_db_instance",
		F:    testSweepDcdbDbInstance,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dcdb_db_instance
func testSweepDcdbDbInstance(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	dcdbService := svcdcdb.NewDcdbService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

	instances, err := dcdbService.DescribeDcdbInstancesByFilter(ctx, nil)
	if err != nil {
		return err
	}
	if instances == nil {
		return fmt.Errorf("dcdb db instance not exists.")
	}

	for _, v := range instances {
		delId := *v.InstanceId
		delName := *v.InstanceName

		if strings.HasPrefix(delName, "test_dcdb_") {
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				err := dcdbService.DeleteDcdbDbInstanceById(ctx, delId)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] delete dcdb db instance %s failed! reason:[%s]", delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudNeedFixDcdbDbInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDCDBDbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbDbInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCDBDbInstanceExists("tencentcloud_dcdb_db_instance.db_instance"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "instance_name", "test_dcdb_db_instance"),
					// resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "zones.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "shard_memory", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "shard_storage", "10"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "shard_node_count", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "shard_count", "2"),

					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "period", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_instance.db_instance", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_instance.db_instance", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "db_version_id", "8.0"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "resource_tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "resource_tags.0.tag_key", "aaa"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "resource_tags.0.tag_value", "bbb"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.0.param", "character_set_server"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.0.value", "utf8mb4"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.1.param", "lower_case_table_names"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.1.value", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.2.param", "sync_mode"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.2.value", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.3.param", "innodb_page_size"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "init_params.3.value", "16384"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "project_id", "0"),
					// resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "extranet_access", "true"),
				),
			},
			{
				Config: testAccDcdbDbInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCDBDbInstanceExists("tencentcloud_dcdb_db_instance.db_instance"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "period", "2"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "instance_name", "test_dcdb_db_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_instance.db_instance", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_instance.db_instance", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "project_id", tcacctest.DefaultProjectId),
					// resource.TestCheckResourceAttr("tencentcloud_dcdb_db_instance.db_instance", "extranet_access", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_db_instance.db_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDCDBDbInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dcdbService := svcdcdb.NewDcdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dcdb_db_instance" {
			continue
		}

		ret, err := dcdbService.DescribeDcdbDbInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if *ret.TotalCount > 0 || len(ret.Instances) > 0 {
			return fmt.Errorf("dcdb db instance still exist, instanceId: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDCDBDbInstanceExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("dcdb db instance  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dcdb db instance id is not set")
		}

		dcdbService := svcdcdb.NewDcdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		ret, err := dcdbService.DescribeDcdbDbInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if *ret.TotalCount == 0 || len(ret.Instances) == 0 {
			return fmt.Errorf("dcdb db instance not found, instanceId: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccDcdbDbInstance_vpc_config = tcacctest.DefaultAzVariable + `
data "tencentcloud_security_groups" "internal" {
	name = "default"
  }
  
  data "tencentcloud_vpc_instances" "vpc" {
	name = "Default-VPC"
  }
  
  data "tencentcloud_vpc_subnets" "subnet" {
	vpc_id = data.tencentcloud_vpc_instances.vpc.instance_list.0.vpc_id
  }
  
  resource "tencentcloud_vpc" "vpc" {
	cidr_block = "172.18.111.0/24"
	name       = "test-pg-network-vpc"
  }
  
  resource "tencentcloud_subnet" "subnet" {
	availability_zone = var.default_az
	cidr_block        = "172.18.111.0/24"
	name              = "test-pg-network-sub1"
	vpc_id            = tencentcloud_vpc.vpc.id
  }
  
  locals {
	vpc_id        = data.tencentcloud_vpc_subnets.subnet.instance_list.0.vpc_id
	subnet_id     = data.tencentcloud_vpc_subnets.subnet.instance_list.0.subnet_id
	sg_id         = data.tencentcloud_security_groups.internal.security_groups.0.security_group_id
	new_vpc_id    = tencentcloud_subnet.subnet.vpc_id
	new_subnet_id = tencentcloud_subnet.subnet.id
  }
  
`

const testAccDcdbDbInstance_basic = testAccDcdbDbInstance_vpc_config + `

resource "tencentcloud_dcdb_db_instance" "db_instance" {
  instance_name = "test_dcdb_db_instance"
  zones = [var.default_az]
  period = 1
  shard_memory = "2"
  shard_storage = "10"
  shard_node_count = "2"
  shard_count = "2"
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
  db_version_id = "8.0"
  resource_tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
  init_params {
	 param = "character_set_server"
	 value = "utf8mb4"
  }
  init_params {
	param = "lower_case_table_names"
	value = "1"
  }
  init_params {
	param = "sync_mode"
	value = "2"
  }
  init_params {
	param = "innodb_page_size"
	value = "16384"
  }
  security_group_ids = [local.sg_id]
  project_id = 0
//   extranet_access = true
}

`

const testAccDcdbDbInstance_update = testAccDcdbDbInstance_vpc_config + tcacctest.DefaultProjectVariable + `

resource "tencentcloud_dcdb_db_instance" "db_instance" {
  instance_name = "test_dcdb_db_instance_CHANGED"
  zones = [var.default_az]
  period = 2
  shard_memory = "2"
  shard_storage = "10"
  shard_node_count = "2"
  shard_count = "2"
  vpc_id = local.new_vpc_id
  subnet_id = local.new_subnet_id
  vip = "172.18.111.10"
  db_version_id = "8.0"
  resource_tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
  init_params {
	 param = "character_set_server"
	 value = "utf8mb4"
  }
  init_params {
	param = "lower_case_table_names"
	value = "1"
  }
  init_params {
	param = "sync_mode"
	value = "2"
  }
  init_params {
	param = "innodb_page_size"
	value = "16384"
  }
  security_group_ids = [local.sg_id]
  project_id = var.default_project
//   extranet_access = false
}

`
