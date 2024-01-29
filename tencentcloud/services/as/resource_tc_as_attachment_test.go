package as_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"

	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func init() {
	resource.AddTestSweepers("tencentcloud_as_attachment", &resource.Sweeper{
		Name: "tencentcloud_as_attachment",
		F:    testSweepAsAttachment,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_as_attachment
func testSweepAsAttachment(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	asService := svcas.NewAsService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())

	scalingGroups, err := asService.DescribeAutoScalingGroupByFilter(ctx, "", "", "", nil)
	if err != nil {
		return fmt.Errorf("list scaling group error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range scalingGroups {
		if !tccommon.CheckResourcePersist(*v.AutoScalingGroupName, *v.CreatedTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.AutoScalingGroupId,
				Name: *v.AutoScalingGroupName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:        *v.AutoScalingGroupId,
			Name:      *v.AutoScalingGroupName,
			CreatTime: *v.CreatedTime,
		})
	}
	tccommon.ProcessScanCloudResources(resources, nonKeepResources, "as", "attachment")

	for _, v := range scalingGroups {
		scalingGroupId := *v.AutoScalingGroupId
		scalingGroupName := *v.AutoScalingGroupName
		if !strings.HasPrefix(scalingGroupName, "tf-as-") {
			continue
		}

		var instanceIds []string
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, errRet := asService.DescribeAutoScalingAttachment(ctx, scalingGroupId, true)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			instanceIds = result
			return nil
		})
		if err != nil {
			return err
		}
		if len(instanceIds) == 0 {
			continue
		}

		if err = asService.DetachInstances(ctx, scalingGroupId, instanceIds); err != nil {
			log.Printf("[ERROR] delete scaling group %s error: %s", scalingGroupId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudAsAttachment(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAsAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAsAttachment(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsAttachmentExists("tencentcloud_as_attachment.attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_attachment.attachment", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_attachment.attachment", "instance_ids.#", "1"),
				),
			},
			// test update case
			{
				Config: testAccAsAttachment_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAsAttachmentExists("tencentcloud_as_attachment.attachment"),
					resource.TestCheckResourceAttrSet("tencentcloud_as_attachment.attachment", "scaling_group_id"),
					resource.TestCheckResourceAttr("tencentcloud_as_attachment.attachment", "instance_ids.#", "2"),
				),
			},
		},
	})
}

func testAccCheckAsAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("auto scaling attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("auto scaling attachment id is not set")
		}
		asService := svcas.NewAsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instances, err := asService.DescribeAutoScalingAttachment(ctx, rs.Primary.ID, false)
		if err != nil {
			return err
		}
		if len(instances) < 1 {
			return fmt.Errorf("auto scaling attachement not exists: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckAsAttachmentDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	asService := svcas.NewAsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_as_attachment" {
			continue
		}

		instances, err := asService.DescribeAutoScalingAttachment(ctx, rs.Primary.ID, false)
		if err != nil {
			if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == svcas.AsScalingGroupNotFound {
					return nil
				}
			}
			return err
		}
		if len(instances) > 0 {
			return fmt.Errorf("auto scaling attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAsAttachment() string {
	return tcacctest.DefaultAsVariable + `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = var.availability_zone
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-attachment-config"
  image_id           = "img-2lr9q49h"
  instance_types     = [data.tencentcloud_instance_types.default.instance_types.0.instance_type]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-attachment-group"
  configuration_id   = tencentcloud_as_scaling_config.launch_configuration.id
  max_size           = 5
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_instance" "cvm_instance" {
  instance_name     = "tf_as_instance"
  availability_zone = var.availability_zone
  image_id          = "img-2lr9q49h"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_SSD"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_as_attachment" "attachment" {
  scaling_group_id = tencentcloud_as_scaling_group.scaling_group.id
  instance_ids     = [tencentcloud_instance.cvm_instance.id]
}
`
}

func testAccAsAttachment_update() string {
	return tcacctest.DefaultAsVariable + `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = var.availability_zone
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-attachment-config"
  image_id           = "img-2lr9q49h"
  instance_types     = [data.tencentcloud_instance_types.default.instance_types.0.instance_type]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-attachment-group"
  configuration_id   = tencentcloud_as_scaling_config.launch_configuration.id
  max_size           = 5
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_instance" "cvm_instance" {
  instance_name     = "tf_as_instance"
  availability_zone = var.availability_zone
  image_id          = "img-2lr9q49h"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_SSD"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_instance" "cvm_instance_1" {
  instance_name     = "tf_as_instance_1"
  availability_zone = var.availability_zone
  image_id          = "img-2lr9q49h"
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_SSD"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_as_attachment" "attachment" {
  scaling_group_id = tencentcloud_as_scaling_group.scaling_group.id
  instance_ids     = [tencentcloud_instance.cvm_instance.id, tencentcloud_instance.cvm_instance_1.id]
}
`
}
