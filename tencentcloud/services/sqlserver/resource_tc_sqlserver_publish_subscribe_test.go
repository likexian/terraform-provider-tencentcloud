package sqlserver_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcsqlserver "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sqlserver"

	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_sqlserver_publish_subscribe
	resource.AddTestSweepers("tencentcloud_sqlserver_publish_subscribe", &resource.Sweeper{
		Name: "tencentcloud_sqlserver_publish_subscribe",
		F:    testAccTencentCloudSQLServerPubSubSweeper,
	})
}

func testAccTencentCloudSQLServerPubSubSweeper(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcsqlserver.NewSqlserverService(client)
	instance, err := service.DescribeSqlserverInstances(ctx, "", tcacctest.DefaultPubSQLServerName, -1, "", "", 1)
	if err != nil {
		return err
	}
	subInstances, err := service.DescribeSqlserverInstances(ctx, "", tcacctest.DefaultSubSQLServerName, -1, "", "", 1)

	if err != nil {
		return err
	}

	pubInstanceId := *instance[0].InstanceId
	subInstanceId := *subInstances[0].InstanceId

	testAccUnsubscribePubDB(ctx, &service, pubInstanceId)

	database, err := service.DescribeDBsOfInstance(ctx, subInstanceId)
	if err != nil {
		return err
	}

	if len(database) == 0 {
		log.Printf("no DBs in %s", subInstanceId)
		return nil
	}

	for i := range database {
		item := database[i]
		created := time.Time{}
		name := *item.Name
		if item.CreateTime != nil {
			created = tccommon.StringToTime(*item.CreateTime)
		}
		if name != tcacctest.DefaultSQLServerPubSubDB || tcacctest.IsResourcePersist("", &created) {
			continue
		}
		if err = service.DeleteSqlserverDB(ctx, subInstanceId, []*string{item.Name}); err != nil {
			log.Printf("err: %s", err.Error())
		}
	}
	return err
}

func testAccCheckPubSubsExists() error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(tcacctest.DefaultRegion)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcsqlserver.NewSqlserverService(client)
	instance, err := service.DescribeSqlserverInstances(ctx, "", tcacctest.DefaultPubSQLServerName, -1, "", "", 1)
	if err != nil {
		return err
	}

	pubInstanceId := *instance[0].InstanceId

	pubsubs, _ := service.DescribeSqlserverPublishSubscribes(ctx, map[string]interface{}{
		"instanceId": pubInstanceId,
	})

	if len(pubsubs) > 0 {
		return fmt.Errorf("pubsub of %s still exists", tcacctest.DefaultPubSQLServerName)
	}
	return nil
}

func testAccUnsubscribePubDB(ctx context.Context, service *svcsqlserver.SqlserverService, instanceId string) {

	pubsubs, _ := service.DescribeSqlserverPublishSubscribes(ctx, map[string]interface{}{
		"instanceId": instanceId,
	})

	if len(pubsubs) == 0 {
		log.Printf("NO pubsub result")
		return
	}

	pubSubId := *pubsubs[0].Id

	pubSub := &sqlserver.PublishSubscribe{
		Id: &pubSubId,
	}
	tuples := []*sqlserver.DatabaseTuple{
		{
			PublishDatabase:   helper.String(tcacctest.DefaultSQLServerPubSubDB),
			SubscribeDatabase: helper.String(tcacctest.DefaultSQLServerPubSubDB),
		},
	}
	err := service.DeletePublishSubscribe(ctx, pubSub, tuples)
	if err != nil {
		fmt.Printf("[ERROR] %s", err.Error())
	}
}

// go test -i; go test -test.run TestAccTencentCloudSqlserverPublishSubscribeResource -v
func TestAccTencentCloudSqlserverPublishSubscribeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
			if err := testAccCheckPubSubsExists(); err != nil {
				t.Errorf("Precheck failed: %s", err.Error())
			}
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSqlserverPublishSubscribeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverPublishSubscribe_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverPublishSubscribeExists("tencentcloud_sqlserver_publish_subscribe.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_publish_subscribe.example", "publish_instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_publish_subscribe.example", "subscribe_instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "publish_subscribe_name", "example"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "database_tuples.#", "1"),
				),
			},
			{
				Config: testAccSqlserverPublishSubscribe_basic_update_name,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSqlserverPublishSubscribeExists("tencentcloud_sqlserver_publish_subscribe.example"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_publish_subscribe.example", "publish_subscribe_name", "example1"),
				),
			},
			{
				ResourceName:            "tencentcloud_sqlserver_publish_subscribe.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_subscribe_db"},
			},
		},
	})
}

func testAccCheckSqlserverPublishSubscribeDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	sqlserverService := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_sqlserver_publish_subscribe" {
			continue
		}
		split := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(split) < 2 {
			continue
		}
		_, has, err := sqlserverService.DescribeSqlserverPublishSubscribeById(ctx, split[0], split[1])
		if err != nil {
			return err
		}
		if has {
			return fmt.Errorf("SQL Server Publish Subscribe %s  still exists", split[0]+tccommon.FILED_SP+split[1])
		}
	}
	return nil
}

func testAccCheckSqlserverPublishSubscribeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("SQL Server Publish Subscribe %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("SQL Server Publish Subscribe id is not set")
		}

		sqlserverService := svcsqlserver.NewSqlserverService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		split := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(split) < 2 {
			return fmt.Errorf("SQL Server Publish Subscribe is not set: %s", rs.Primary.ID)
		}
		_, has, err := sqlserverService.DescribeSqlserverPublishSubscribeById(ctx, split[0], split[1])
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("SQL Server Publish Subscribe %s is not found", rs.Primary.ID)
		}
		return nil
	}
}

const testAccSqlserverPublishSubscribe_basic = tcacctest.CommonPubSubSQLServer + `
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
  publish_instance_id    = "mssql-qelbzgwf"
  subscribe_instance_id  = "mssql-jdk2pwld"
  publish_subscribe_name = "example"
  delete_subscribe_db    = false
  database_tuples {
    publish_database   = local.sqlserver_pub_db
    subscribe_database = local.sqlserver_sub_db
  }
}
`

const testAccSqlserverPublishSubscribe_basic_update_name = tcacctest.CommonPubSubSQLServer + `
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
  publish_instance_id    = "mssql-qelbzgwf"
  subscribe_instance_id  = "mssql-jdk2pwld"
  publish_subscribe_name = "example_update"
  delete_subscribe_db    = false
  database_tuples {
    publish_database   = local.sqlserver_pub_db
    subscribe_database = local.sqlserver_sub_db
  }
}
`
