package cvm_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_instance", &resource.Sweeper{
		Name: "tencentcloud_instance",
		F:    testSweepCvmInstance,
	})
}

func testSweepCvmInstance(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	cvmService := svccvm.NewCvmService(client)

	instances, err := cvmService.DescribeInstanceByFilter(ctx, nil, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range instances {
		if !tccommon.CheckResourcePersist(*v.InstanceName, *v.CreatedTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.InstanceId,
				Name: *v.InstanceName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:         *v.InstanceId,
			Name:       *v.InstanceName,
			CreateTime: *v.CreatedTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "RunInstances")

	for _, v := range instances {
		instanceId := *v.InstanceId
		//instanceName := *v.InstanceName
		now := time.Now()
		createTime := tccommon.StringToTime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()

		//if strings.HasPrefix(instanceName, tcacctest.KeepResource) || strings.HasPrefix(instanceName, tcacctest.DefaultResource) {
		//	continue
		//}

		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = cvmService.DeleteInstance(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudInstanceResource_Basic(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.cvm_basic"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckResourceAttrSet(id, "vpc_id"),
					resource.TestCheckResourceAttrSet(id, "subnet_id"),
					resource.TestCheckResourceAttrSet(id, "project_id"),
					resource.TestCheckResourceAttr(id, "tags.hostname", "tci"),
				),
			},
			{
				Config: testAccTencentCloudInstanceModifyInstanceType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "instance_type"),
				),
			},
			{
				ResourceName:            id,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"disable_monitor_service", "disable_security_service", "hostname", "password", "force_delete"},
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_PrepaidBasic(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.cvm_prepaid_basic"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstancePrepaidBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "private_ip"),
					resource.TestCheckResourceAttrSet(id, "vpc_id"),
					resource.TestCheckResourceAttrSet(id, "subnet_id"),
					resource.TestCheckResourceAttrSet(id, "project_id"),
					resource.TestCheckResourceAttr(id, "tags.hostname", "tci"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithDataDisk(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithDataDisk,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "system_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_snapshot_id", ""),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "100"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithDataDiskUpdate,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "system_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "system_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "150"),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_snapshot_id", ""),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "150"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithNetwork(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithNetworkFalse("false"),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckNoResourceAttr(id, "public_ip"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithNetwork("true", 5),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "internet_max_bandwidth_out", "5"),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "public_ip"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithPrivateIP(t *testing.T) {
	t.Parallel()
	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithPrivateIP,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithKeyPairs(t *testing.T) {
	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithKeyPair_withoutKeyPair,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithKeyPair(
					"[tencentcloud_key_pair.key_pair_0.id, tencentcloud_key_pair.key_pair_1.id]",
				),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "key_ids.#", "2"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 5)
				},
				Config: testAccTencentCloudInstanceWithKeyPair("[tencentcloud_key_pair.key_pair_2.id]"),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "key_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithPassword(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithPassword("TF_test_123"),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "password"),
				),
			},
			{
				PreConfig: func() {
					time.Sleep(time.Second * 5)
				},
				Config: testAccTencentCloudInstanceWithPassword("TF_test_123456"),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "password"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithImageLogin(t *testing.T) {

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithImageLogin,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "keep_image_login", "true"),
					resource.TestCheckResourceAttr(id, "disable_api_termination", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithName(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithName(tcacctest.DefaultInsName),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "instance_name", tcacctest.DefaultInsName),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithName(tcacctest.DefaultInsNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "instance_name", tcacctest.DefaultInsNameUpdate),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithHostname(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithHostname,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "hostname", tcacctest.DefaultInsName),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithSecurityGroup(t *testing.T) {
	t.Parallel()

	instanceId := "tencentcloud_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: instanceId,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithSecurityGroup(`["sg-cm7fbbf3"]`),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(instanceId),
					testAccCheckTencentCloudInstanceExists(instanceId),
					resource.TestCheckResourceAttr(instanceId, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(instanceId, "security_groups.#", "1"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithSecurityGroup(`[
					"sg-cm7fbbf3",
					"sg-kensue7b"
				]`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceId, "security_groups.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithOrderlySecurityGroup(t *testing.T) {
	t.Parallel()

	instanceId := "tencentcloud_instance.cvm_with_orderly_sg"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: instanceId,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceOrderlySecurityGroups,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(instanceId),

					resource.TestCheckResourceAttr(instanceId, "orderly_security_groups.0", "sg-cm7fbbf3"),
					resource.TestCheckResourceAttr(instanceId, "orderly_security_groups.1", "sg-kensue7b"),
					resource.TestCheckResourceAttr(instanceId, "orderly_security_groups.2", "sg-05f7wnhn"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithTags(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithTags(`{
					"hello" = "world"
					"happy" = "hour"
				}`),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "tags.hello", "world"),
					resource.TestCheckResourceAttr(id, "tags.happy", "hour"),
				),
			},
			{
				Config: testAccTencentCloudInstanceWithTags(`{
					"hello" = "hello"
				}`),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttr(id, "tags.hello", "hello"),
					resource.TestCheckNoResourceAttr(id, "tags.happy"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithPlacementGroup(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithPlacementGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
					resource.TestCheckResourceAttrSet(id, "placement_group_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_WithSpotpaid(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithSpotpaid,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_DataDiskOrder(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceWithDataDiskOrder,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "data_disks.0.data_disk_size", "100"),
					resource.TestCheckResourceAttr(id, "data_disks.1.data_disk_size", "50"),
					resource.TestCheckResourceAttr(id, "data_disks.2.data_disk_size", "70"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_DataDiskByCbs(t *testing.T) {
	t.Parallel()

	id := "tencentcloud_instance.cvm_add_data_disk_by_cbs"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceAddDataDiskByCbs,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
				),
			},
		},
	})
}

func TestAccTencentCloudNeedFixInstancePostpaidToPrepaid(t *testing.T) {

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstancePostPaid,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
			{
				Config: testAccTencentCloudInstanceBasicToPrepaid,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_period", "1"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceResource_PrepaidFallbackToPostpaid(t *testing.T) {

	id := "tencentcloud_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:      func() { tcacctest.AccPreCheck(t) },
		IDRefreshName: id,
		Providers:     tcacctest.AccProviders,
		CheckDestroy:  testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceBasicToPrepaid,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_period", "1"),
					resource.TestCheckResourceAttr(id, "instance_charge_type_prepaid_renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
				),
			},
			{
				Config: testAccTencentCloudInstancePostPaid,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID(id),
					testAccCheckTencentCloudInstanceExists(id),
					resource.TestCheckResourceAttr(id, "instance_status", "RUNNING"),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no security group ID is set")
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		sg, err := service.DescribeSecurityGroup(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if sg == nil {
			return fmt.Errorf("security group not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckTencentCloudInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cvm instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cvm instance id is not set")
		}

		cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, err = cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("cvm instance id is not found")
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := svccvm.NewCvmService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_instance" {
			continue
		}

		instance, err := cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, err = cvmService.DescribeInstanceById(ctx, rs.Primary.ID)
				if err != nil {
					return tccommon.RetryError(err)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		if instance != nil && *instance.InstanceState != svccvm.CVM_STATUS_SHUTDOWN && *instance.InstanceState != svccvm.CVM_STATUS_TERMINATING {
			return fmt.Errorf("cvm instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccTencentCloudInstanceBasic = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-basic-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-basic-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0

  tags = {
    hostname = "tci"
  }
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstancePrepaidBasic = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-prepaid-basic-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-prepaid-basic-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "cvm_prepaid_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
  instance_charge_type                    = "PREPAID"
  instance_charge_type_prepaid_period     = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  force_delete = true
  tags = {
    hostname = "tci"
  }
}
`

const testAccTencentCloudInstanceWithDataDiskOrder = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-with-cbs-order-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-with-cbs-order-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0

  data_disks {
    data_disk_size         = 100
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
  data_disks {
    data_disk_size         = 50
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
  data_disks {
    data_disk_size         = 70
    data_disk_type         = "CLOUD_PREMIUM"
    delete_with_instance   = true
  }
}
`

const testAccTencentCloudInstanceAddDataDiskByCbs = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-attach-cbs-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-attach-cbs-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "cvm_add_data_disk_by_cbs" {
  instance_name     = "cvm-add-data-disk-by-cbs"
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}

resource "tencentcloud_cbs_storage" "cbs_disk1" {
	storage_name = "cbs_disk1"
	storage_type = "CLOUD_SSD"
	storage_size = 200
	availability_zone = var.availability_cvm_zone
	project_id = 0
	encrypt = false
}
resource "tencentcloud_cbs_storage" "cbs_disk2" {
	storage_name = "cbs_disk2"
	storage_type = "CLOUD_SSD"
	storage_size = 100
	availability_zone = var.availability_cvm_zone
	project_id = 0
	encrypt = false
}
resource "tencentcloud_cbs_storage_attachment" "attachment_cbs_disk1" {
	storage_id = tencentcloud_cbs_storage.cbs_disk1.id
	instance_id = tencentcloud_instance.cvm_add_data_disk_by_cbs.id
}
resource "tencentcloud_cbs_storage_attachment" "attachment_cbs_disk2" {
	storage_id = tencentcloud_cbs_storage.cbs_disk2.id
	instance_id = tencentcloud_instance.cvm_add_data_disk_by_cbs.id
}
`

const testAccTencentCloudInstancePostPaid = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 2
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "` + tcacctest.DefaultInsName + `"
  availability_zone = "` + tcacctest.DefaultAZone + `"
  image_id          = "` + tcacctest.DefaultTkeOSImageId + `"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  force_delete = true
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstanceBasicToPrepaid = `
data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 2
  memory_size    = 2
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "` + tcacctest.DefaultInsName + `"
  availability_zone = "` + tcacctest.DefaultAZone + `"
  image_id          = "` + tcacctest.DefaultTkeOSImageId + `"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  instance_charge_type       = "PREPAID"
  instance_charge_type_prepaid_period = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  force_delete = true
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstanceModifyInstanceType = tcacctest.DefaultInstanceVariable + `
data "tencentcloud_instance_types" "new_type" {
	availability_zone = var.availability_cvm_zone
  
	cpu_core_count = 2
	memory_size    = 2
  }

resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-basic-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-basic-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "cvm_basic" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.new_type.instance_types.0.instance_type
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0

  tags = {
    hostname = "tci"
  }
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstanceWithDataDisk = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    delete_with_instance  = true
	// encrypt = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 100
    # data_disk_snapshot_id = "snap-nvzu3dmh"
    delete_with_instance  = true
  }

  disable_security_service = true
  disable_monitor_service  = true
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

const testAccTencentCloudInstanceWithDataDiskUpdate = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type

  system_disk_type = "CLOUD_PREMIUM"
  system_disk_size = 100

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  } 
   
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }

  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }



  disable_security_service = true
  disable_monitor_service  = true
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`

func testAccTencentCloudInstanceWithNetworkFalse(hasPublicIp string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		hasPublicIp,
	)
}

func testAccTencentCloudInstanceWithNetwork(hasPublicIp string, maxBandWidthOut int64) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  internet_max_bandwidth_out = %d
  allocate_public_ip         = %s
  system_disk_type           = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		maxBandWidthOut, hasPublicIp,
	)
}

const testAccTencentCloudInstanceWithPrivateIP = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_vpc" "vpc" {
	name       = "cvm-with-privateip-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "cvm-with-privateip-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = var.availability_cvm_zone
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  private_ip        = "10.0.0.123"
}
`

const testAccTencentCloudInstanceWithKeyPair_withoutKeyPair = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
	instance_name     = var.instance_name
	availability_zone = var.availability_cvm_zone
	image_id          = data.tencentcloud_images.default.images.0.image_id
	instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
	system_disk_type  = "CLOUD_PREMIUM"
	lifecycle {
		ignore_changes = [instance_type]
	}
}
`

func testAccTencentCloudInstanceWithKeyPair(keyIds string) string {

	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_key_pair" "key_pair_0" {
  key_name = "key_pair_0"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

resource "tencentcloud_key_pair" "key_pair_1" {
  key_name = "key_pair_1"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzwYE6KI8uULEvSNA2k1tlsLtMDe+x1Saw6yL3V1mk9NFws0K2BshYqsnP/BlYiGZv/Nld5xmGoA9LupOcUpyyGGSHZdBrMx1Dz9ajewe7kGowRWwwMAHTlzh9+iqeg/v6P5vW6EwK4hpGWgv06vGs3a8CzfbHu1YRbZAO/ysp3ymdL+vGvw/vzC0T+YwPMisn9wFD5FTlJ+Em6s9PzxqR/41t4YssmCwUV78ZoYL8CyB0emuB8wALvcXbdUVxMxpBEHd5U6ZP5+HPxU2WFbWqiFCuErLIZRuxFw8L/Ot+JOyNnadN1XU4crYDX5cML1i/ExXKVIDoBaLtgAJOpyeP"
}

resource "tencentcloud_key_pair" "key_pair_2" {
  key_name = "key_pair_2"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDJ1zyoM55pKxJptZBKceZSEypPN7BOunqBR1Qj3Tz5uImJ+dwfKzggu8PGcbHtuN8D2n1BH/GDkiGFaz/sIYUJWWZudcdut+ra32MqUvk953Sztf12rsFC1+lZ1CYEgon8Lt6ehxn+61tsS31yfUmpL1mq2vuca7J0NLdPMpxIYkGlifyAMISMmxi/m7gPYpbdZTmhQQS2aOhuLm+B4MwtTvT58jqNzIaFU0h5sqAvGQfzI5pcxwYvFTeQeXjJZfaYapDHN0MAg0b/vIWWNrDLv7dlv//OKBIaL0LIzIGQS8XXhF3HlyqfDuf3bjLBIKzYGSV/DRqlEsGBgzinJZXvJZug5oq1n2njDFsdXEvL6fYsP4WLvBLiQlceQ7oXi7m5nfrwFTaX+mpo7dUOR9AcyQ1AAgCcM67orB4E33ycaArGHtpjnCnWUjqQ+yCj4EXsD4yOL77wGsmhkbboVNnYAD9MJWsFP03hZE7p/RHY0C5NfLPT3mL45oZxBpC5mis="
}

resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  key_ids           = %s
  system_disk_type  = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		keyIds,
	)
}

func testAccTencentCloudInstanceWithPassword(password string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  password                   = "%s"
  system_disk_type           = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		password,
	)
}

const testAccTencentCloudInstanceWithImageLogin = tcacctest.DefaultInstanceVariable + `
data "tencentcloud_images" "zoo" {
  image_type = ["PRIVATE_IMAGE"]
}
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.zoo.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  keep_image_login 			 = true
  system_disk_type           = "CLOUD_PREMIUM"
  disable_api_termination    = false
}
`

func testAccTencentCloudInstanceWithName(instanceName string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = "%s"
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		instanceName,
	)
}

const testAccTencentCloudInstanceWithHostname = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  hostname          = var.instance_name
  system_disk_type  = "CLOUD_PREMIUM"
}
`

func testAccTencentCloudInstanceWithSecurityGroup(ids string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name              = var.instance_name
  availability_zone          = var.availability_cvm_zone
  image_id                   = data.tencentcloud_images.default.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  security_groups            = %s
  lifecycle {
	ignore_changes = [instance_type]
  }
}
`,
		ids,
	)
}

func testAccTencentCloudInstanceWithTags(tags string) string {
	return fmt.Sprintf(
		tcacctest.DefaultInstanceVariable+`
resource "tencentcloud_instance" "foo" {
  instance_name     = var.instance_name
  availability_zone = var.availability_cvm_zone
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  data_disks {
    data_disk_type        = "CLOUD_PREMIUM"
    data_disk_size        = 150
    delete_with_instance  = true
  }
  lifecycle {
	ignore_changes = [instance_type]
  }
  tags = %s
}
`,
		tags,
	)
}

const testAccTencentCloudInstanceWithPlacementGroup = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name      = var.instance_name
  availability_zone  = var.availability_cvm_zone
  image_id           = data.tencentcloud_images.default.images.0.image_id
  instance_type      = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type   = "CLOUD_PREMIUM"
  placement_group_id = "ps-1y147q03"
}
`

const testAccTencentCloudInstanceWithSpotpaid = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "foo" {
  instance_name        = var.instance_name
  availability_zone    = var.availability_cvm_zone
  image_id             = data.tencentcloud_images.default.images.0.image_id
  instance_type        = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type     = "CLOUD_PREMIUM"
  instance_charge_type = "SPOTPAID"
  spot_instance_type   = "ONE-TIME"
  spot_max_price       = "0.5"
}
`

const testAccTencentCloudInstanceOrderlySecurityGroups = tcacctest.DefaultInstanceVariable + `
resource "tencentcloud_instance" "cvm_with_orderly_sg" {
	instance_name              = "test-orderly-sg-cvm"
	availability_zone          = var.availability_cvm_zone
	image_id                   = data.tencentcloud_images.default.images.0.image_id
	instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
	system_disk_type           = "CLOUD_PREMIUM"
	orderly_security_groups    = ["sg-cm7fbbf3", "sg-kensue7b", "sg-05f7wnhn"]
}
`
