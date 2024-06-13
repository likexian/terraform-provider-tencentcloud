package emr_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	svccdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
	svcemr "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/emr"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_emr
	resource.AddTestSweepers("tencentcloud_emr", &resource.Sweeper{
		Name: "tencentcloud_emr",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			sharedClient, err := tcacctest.SharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

			emrService := svcemr.NewEMRService(client)
			filters := make(map[string]interface{})
			filters["display_strategy"] = svcemr.DisplayStrategyIsclusterList
			clusters, err := emrService.DescribeInstances(ctx, filters)
			if err != nil {
				return nil
			}

			// add scanning resources
			var resources, nonKeepResources []*tccommon.ResourceInstance
			for _, v := range clusters {
				if !tccommon.CheckResourcePersist(*v.ClusterId, *v.AddTime) {
					nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
						Id:   *v.ClusterId,
						Name: *v.ClusterName,
					})
				}
				resources = append(resources, &tccommon.ResourceInstance{
					Id:         *v.ClusterId,
					Name:       *v.ClusterName,
					CreateTime: *v.AddTime,
				})
			}
			tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateInstance")

			for _, cluster := range clusters {
				clusterName := *cluster.ClusterName
				if strings.HasPrefix(clusterName, tcacctest.KeepResource) || strings.HasPrefix(clusterName, tcacctest.DefaultResource) {
					continue
				}
				now := time.Now()
				createTime := tccommon.StringToTime(*cluster.AddTime)
				interval := now.Sub(createTime).Minutes()
				// less than 30 minute, not delete
				if tccommon.NeedProtect == 1 && int64(interval) < 30 {
					continue
				}
				metaDB := cluster.MetaDb
				instanceId := *cluster.ClusterId
				request := emr.NewTerminateInstanceRequest()
				request.InstanceId = &instanceId
				if _, err = client.UseEmrClient().TerminateInstance(request); err != nil {
					return nil
				}
				err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
					clusters, err := emrService.DescribeInstancesById(ctx, instanceId, svcemr.DisplayStrategyIsclusterList)

					if e, ok := err.(*errors.TencentCloudSDKError); ok {
						if e.GetCode() == "InternalError.ClusterNotFound" {
							return nil
						}
						if e.GetCode() == "UnauthorizedOperation" {
							return nil
						}
					}

					if len(clusters) > 0 {
						status := *(clusters[0].Status)
						if status != svcemr.EmrInternetStatusDeleted {
							return resource.RetryableError(
								fmt.Errorf("%v create cluster endpoint status still is %v", instanceId, status))
						}
					}

					if err != nil {
						return resource.RetryableError(err)
					}
					return nil
				})
				if err != nil {
					return nil
				}

				if metaDB != nil && *metaDB != "" {
					// remove metadb
					mysqlService := svccdb.NewMysqlService(client)

					err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						err := mysqlService.OfflineIsolatedInstances(ctx, *metaDB)
						if err != nil {
							return tccommon.RetryError(err, tccommon.InternalError)
						}
						return nil
					})

					if err != nil {
						return nil
					}
				}
			}
			return nil
		},
	})
}

var testEmrClusterResourceKey = "tencentcloud_emr_cluster.emrrrr"

func TestAccTencentCloudEmrClusterResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testEmrBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmrExists(testEmrClusterResourceKey),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "product_id", "38"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "vpc_settings.vpc_id", tcacctest.DefaultEMRVpcId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "vpc_settings.subnet_id", tcacctest.DefaultEMRSubnetId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "softwares.#", "5"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "support_ha", "0"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "instance_name", "emr-test-demo"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "login_settings.password", "Tencent@cloud123"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "time_span", "3600"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "time_unit", "s"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "pay_mode", "0"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "placement_info.0.zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "placement_info.0.project_id", "0"),
					resource.TestCheckResourceAttrSet(testEmrClusterResourceKey, "instance_id"),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "sg_id", tcacctest.DefaultEMRSgId),
					resource.TestCheckResourceAttr(testEmrClusterResourceKey, "tags.emr-key", "emr-value"),
				),
			},
			{
				ResourceName:            testEmrClusterResourceKey,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"display_strategy", "placement", "time_span", "time_unit", "login_settings"},
			},
		},
	})
}

func testAccCheckEmrExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("emr cluster %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("emr cluster id is not set")
		}

		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		service := svcemr.NewEMRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		instanceId := rs.Primary.ID
		clusters, err := service.DescribeInstancesById(ctx, instanceId, svcemr.DisplayStrategyIsclusterList)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				clusters, err = service.DescribeInstancesById(ctx, instanceId, svcemr.DisplayStrategyIsclusterList)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}

		if err != nil {
			return nil
		}
		if len(clusters) <= 0 {
			return fmt.Errorf("emr cluster create fail")
		} else {
			log.Printf("[DEBUG]emr cluster  %s create  ok", rs.Primary.ID)
			return nil
		}

	}
}

const testEmrBasic = tcacctest.DefaultEMRVariable + `
data "tencentcloud_instance_types" "cvm4c8m" {
	exclude_sold_out=true
	cpu_core_count=4
	memory_size=8
    filter {
      name   = "instance-charge-type"
      values = ["POSTPAID_BY_HOUR"]
    }
    filter {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
}

resource "tencentcloud_emr_cluster" "emrrrr" {
	product_id=38
	vpc_settings={
	  vpc_id=var.vpc_id
	  subnet_id=var.subnet_id
	}
	softwares = [
	  "hdfs-2.8.5",
	  "knox-1.6.1",
	  "openldap-2.4.44",
	  "yarn-2.8.5",
	  "zookeeper-3.6.3",
	]
	support_ha=0
	instance_name="emr-test-demo"
	resource_spec {
	  master_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
	  }
	  core_resource_spec {
		mem_size=8192
		cpu=4
		disk_size=100
		disk_type="CLOUD_PREMIUM"
		spec="CVM.${data.tencentcloud_instance_types.cvm4c8m.instance_types.0.family}"
		storage_type=5
		root_size=50
	  }
	  master_count=1
	  core_count=2
	}
	login_settings={
	  password="Tencent@cloud123"
	}
	time_span=3600
	time_unit="s"
	pay_mode=0
	placement_info {
	  zone="ap-guangzhou-3"
	  project_id=0
	}
	sg_id=var.sg_id
	tags = {
        emr-key = "emr-value"
    }
  }
`
