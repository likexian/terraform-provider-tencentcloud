package tmp_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudMonitorInstance_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMonInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("tencentcloud_monitor_tmp_instance.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.basic", "instance_name", "demo-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.basic", "data_retention_time", "30"),
				),
			},
			{
				Config: testInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("tencentcloud_monitor_tmp_instance.update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.update", "instance_name", "demo-test-update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.update", "data_retention_time", "30"),
				),
			},
			//{
			//	ResourceName:      "tencentcloud_monitor_tmp_instance.basic",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}

func testAccCheckMonInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_instance" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		instance, err := service.DescribeMonitorTmpInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance != nil {
			status := strconv.FormatInt(*instance.InstanceStatus, 10)
			if strings.Contains("5,6,8,9", status) {
				return nil
			}
			return fmt.Errorf("instance %s still exists: %v", rs.Primary.ID, *instance.InstanceStatus)
		}
	}

	return nil
}

func testAccCheckInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeMonitorTmpInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil || *instance.InstanceStatus != 2 {
			return fmt.Errorf("instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testInstanceVar = tcacctest.DefaultAzVariable + `
variable "vpc_id" {
  default = "` + tcacctest.DefaultTmpVpcId + `"
}
variable "subnet_id" {
  default = "` + tcacctest.DefaultTmpSubnetId + `"
}
`
const testInstance_basic = testInstanceVar + `
resource "tencentcloud_monitor_tmp_instance" "basic" {
 instance_name 		= "demo-test"
 vpc_id 				= var.vpc_id
 subnet_id				= var.subnet_id
 data_retention_time	= 30
 zone 					= var.default_az
 tags = {
   "createdBy" = "terraform"
 }
}`

const testInstance_update = testInstanceVar + `
resource "tencentcloud_monitor_tmp_instance" "update" {
 instance_name 		= "demo-test-update"
 vpc_id 				= var.vpc_id
 subnet_id				= var.subnet_id
 data_retention_time	= 30
 zone 					= var.default_az
 tags = {
   "createdBy" = "terraform"
 }
}`
