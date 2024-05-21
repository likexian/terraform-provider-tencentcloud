package cdb_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

var testAccTencentCloudMysqlPrivilegeType = "tencentcloud_mysql_privilege"
var testAccTencentCloudMysqlPrivilegeName = testAccTencentCloudMysqlPrivilegeType + ".privilege"

func TestAccTencentCloudMysqlPrivilegeResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccMysqlPrivilegeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlPrivilege,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccMysqlPrivilegeExists,
					resource.TestCheckResourceAttrSet(testAccTencentCloudMysqlPrivilegeName, "mysql_id"),
					resource.TestCheckResourceAttrSet(testAccTencentCloudMysqlPrivilegeName, "account_name"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "global.#", "1"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "table.#", "1"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "column.#", "2"),
					resource.TestCheckTypeSetElemAttr(testAccTencentCloudMysqlPrivilegeName, "global.*", "TRIGGER"),
				),
			},
			{
				Config: testAccMysqlPrivilegeUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccMysqlPrivilegeExists,
					resource.TestCheckResourceAttrSet(testAccTencentCloudMysqlPrivilegeName, "mysql_id"),
					resource.TestCheckResourceAttrSet(testAccTencentCloudMysqlPrivilegeName, "account_name"),
					resource.TestCheckTypeSetElemAttr(testAccTencentCloudMysqlPrivilegeName, "global.*", "TRIGGER"),

					//diff
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "global.#", "2"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "table.#", "2"),
					resource.TestCheckResourceAttr(testAccTencentCloudMysqlPrivilegeName, "column.#", "0"),
					resource.TestCheckTypeSetElemAttr(testAccTencentCloudMysqlPrivilegeName, "global.*", "SELECT"),
				),
			},
		},
	})
}

func testAccMysqlPrivilegeExists(s *terraform.State) error {

	rs, ok := s.RootModule().Resources[testAccTencentCloudMysqlPrivilegeName]
	if !ok {
		return fmt.Errorf("resource %s is not found", testAccTencentCloudMysqlPrivilegeName)
	}

	var privilegeId localcdb.ResourceTencentCloudMysqlPrivilegeId

	if err := json.Unmarshal([]byte(rs.Primary.ID), &privilegeId); err != nil {
		return fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
	}

	request := cdb.NewDescribeAccountPrivilegesRequest()
	request.InstanceId = &privilegeId.MysqlId
	request.User = &privilegeId.AccountName
	request.Host = &privilegeId.AccountHost

	var response *cdb.DescribeAccountPrivilegesResponse
	var inErr, outErr error

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, inErr = tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().DescribeAccountPrivileges(request)
		if inErr != nil {
			if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
					return resource.NonRetryableError(fmt.Errorf("mysql account not exists in mysql"))
				}
				if sdkErr.Code == "InvalidParameter" && strings.Contains(sdkErr.GetMessage(), "instance not found") {
					return resource.NonRetryableError(fmt.Errorf("mysql account not exists in mysql"))
				}
				if sdkErr.Code == "InternalError.TaskError" && strings.Contains(sdkErr.Message, "User does not exist") {
					return resource.NonRetryableError(fmt.Errorf("mysql account not exists in mysql"))
				}
			}
			return tccommon.RetryError(inErr, tccommon.InternalError)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if response == nil || response.Response == nil {
		return errors.New("sdk DescribeAccountPrivileges return error,miss Response")
	}

	if len(response.Response.GlobalPrivileges) > 0 ||
		len(response.Response.ColumnPrivileges) > 0 ||
		len(response.Response.TablePrivileges) > 0 ||
		len(response.Response.DatabasePrivileges) > 0 {
		return nil
	}
	return fmt.Errorf("set privilege return nil")
}

func testAccMysqlPrivilegeDestroy(s *terraform.State) error {
	rs, ok := s.RootModule().Resources[testAccTencentCloudMysqlPrivilegeName]
	if !ok {
		return fmt.Errorf("resource %s is not found", testAccTencentCloudMysqlPrivilegeName)
	}

	var privilegeId localcdb.ResourceTencentCloudMysqlPrivilegeId

	if err := json.Unmarshal([]byte(rs.Primary.ID), &privilegeId); err != nil {
		return fmt.Errorf("Local data[terraform.tfstate] corruption,can not got old account privilege id")
	}

	mysqlService := localcdb.NewMysqlService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	instance, err := mysqlService.DescribeDBInstanceById(tccommon.ContextNil, privilegeId.MysqlId)

	if err != nil {
		return err
	}

	if instance == nil {
		return nil
	}

	request := cdb.NewDescribeAccountPrivilegesRequest()
	request.InstanceId = &privilegeId.MysqlId
	request.User = &privilegeId.AccountName
	request.Host = &privilegeId.AccountHost

	var response *cdb.DescribeAccountPrivilegesResponse
	var inErr, outErr error

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, inErr = tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().DescribeAccountPrivileges(request)
		if inErr != nil {
			if sdkErr, ok := inErr.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == localcdb.MysqlInstanceIdNotFound {
					return nil
				}
				if sdkErr.Code == "InvalidParameter" && strings.Contains(sdkErr.GetMessage(), "instance not found") {
					return nil
				}
				if sdkErr.Code == "InternalError.TaskError" && strings.Contains(sdkErr.Message, "User does not exist") {
					return nil
				}
				if sdkErr.Code == "InvalidParameterValue.UserNotExistError" {
					return nil
				}
			}
			return tccommon.RetryError(inErr, tccommon.InternalError)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if response == nil || response.Response == nil {
		return nil
	}

	if len(response.Response.GlobalPrivileges) > 0 ||
		len(response.Response.ColumnPrivileges) > 0 ||
		len(response.Response.TablePrivileges) > 0 ||
		len(response.Response.DatabasePrivileges) > 0 {
		return fmt.Errorf("privilege is still exist")
	}
	return nil
}

const testAccMysqlPrivilege = testAccMysql + `
resource "tencentcloud_mysql_account" "mysql_account" {
  mysql_id    = tencentcloud_mysql_instance.mysql.id
  name        = "test11priv"
  host        = "119.168.110.%%"
  password    = "test1234"
  description = "test from terraform"
}

resource "tencentcloud_mysql_privilege" "privilege" {
  mysql_id     = tencentcloud_mysql_instance.mysql.id
  account_name = tencentcloud_mysql_account.mysql_account.name
  account_host = tencentcloud_mysql_account.mysql_account.host
  global       = ["TRIGGER"]
  database {
    privileges    = ["SELECT"]
    database_name = "performance_schema"
  }
  table {
    privileges    = ["SELECT", "INSERT", "UPDATE"]
    database_name = "mysql"
    table_name    = "user"
  }
  column {
    privileges    = ["SELECT"]
    database_name = "mysql"
    table_name    = "user"
    column_name   = "host"
  }

  column {
    privileges    = ["SELECT"]
    database_name = "mysql"
    table_name    = "user"
    column_name   = "user"
  }
}`

const testAccMysqlPrivilegeUpdate = testAccMysql + `
resource "tencentcloud_mysql_account" "mysql_account" {
  mysql_id    = tencentcloud_mysql_instance.mysql.id
  name        = "test11priv"
  host        = "119.168.110.%%"
  password    = "test1234"
  description = "test from terraform"
}

resource "tencentcloud_mysql_privilege" "privilege" {
  mysql_id     = tencentcloud_mysql_instance.mysql.id
  account_name = tencentcloud_mysql_account.mysql_account.name
  account_host = tencentcloud_mysql_account.mysql_account.host
  global       = ["TRIGGER","SELECT"]
  table {
    privileges    = ["SELECT"]
    database_name = "mysql"
    table_name    = "user"
  }
  table {
    privileges    = ["SELECT"]
    database_name = "mysql"
    table_name    = "db"
  }
}`
