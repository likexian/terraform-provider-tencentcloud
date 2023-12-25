package tem_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctem "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tem"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
)

// go test -i; go test -test.run TestAccTencentCloudTemApplicationServiceResource_basic -v
func TestAccTencentCloudTemApplicationServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTemApplicationServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemApplicationService,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemApplicationServiceExists("tencentcloud_tem_application_service.application_service"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "environment_id", tcacctest.DefaultEnvironmentId),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "application_id", tcacctest.DefaultApplicationId),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.type", "CLUSTER"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.service_name", "terraform-test-0"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_application_service.application_service", "service.0.ip"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.port_mapping_item_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.port_mapping_item_list.0.port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.port_mapping_item_list.0.target_port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.port_mapping_item_list.0.protocol", "TCP"),
				),
			},
			{
				Config: testAccTemApplicationServiceUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTemApplicationServiceExists("tencentcloud_tem_application_service.application_service"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "environment_id", tcacctest.DefaultEnvironmentId),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "application_id", tcacctest.DefaultApplicationId),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.type", "EXTERNAL"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.service_name", "terraform-test-0"),
					resource.TestCheckResourceAttrSet("tencentcloud_tem_application_service.application_service", "service.0.ip"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.port_mapping_item_list.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.port_mapping_item_list.0.port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.port_mapping_item_list.0.target_port", "80"),
					resource.TestCheckResourceAttr("tencentcloud_tem_application_service.application_service", "service.0.port_mapping_item_list.0.protocol", "TCP"),
				),
			},
			{
				ResourceName:      "tencentcloud_tem_application_service.application_service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTemApplicationServiceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tem_application_service" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		applicationId := idSplit[1]
		serviceName := idSplit[2]

		res, err := service.DescribeTemApplicationServiceById(ctx, environmentId, applicationId)
		if res == nil || res.Result == nil {
			for _, v := range res.Result.ServicePortMappingList {
				if *v.ServiceName == serviceName {
					return fmt.Errorf("tem applicationService %s still exists", rs.Primary.ID)
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTemApplicationServiceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		environmentId := idSplit[0]
		applicationId := idSplit[1]
		serviceName := idSplit[2]

		service := svctem.NewTemService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		res, err := service.DescribeTemApplicationServiceById(ctx, environmentId, applicationId)

		if res == nil || res.Result == nil {
			var applicationService *tem.ServicePortMapping
			for _, v := range res.Result.ServicePortMappingList {
				if *v.ServiceName == serviceName {
					applicationService = v
				}
			}
			if applicationService == nil {
				return fmt.Errorf("tem applicationService %s is not found", rs.Primary.ID)
			}
			return nil
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTemApplicationServiceVar = `
variable "environment_id" {
  default = "` + tcacctest.DefaultEnvironmentId + `"
}
variable "application_id" {
	default = "` + tcacctest.DefaultApplicationId + `"
  }
`

const testAccTemApplicationService = testAccTemApplicationServiceVar + `

resource "tencentcloud_tem_application_service" "application_service" {
	environment_id = var.environment_id
	application_id = var.application_id
	service {
		type = "CLUSTER"
		service_name = "terraform-test-0"
		port_mapping_item_list {
			port = 80
			target_port = 80
			protocol = "TCP"
		}
	}
}

`

const testAccTemApplicationServiceUp = testAccTemApplicationServiceVar + `

resource "tencentcloud_tem_application_service" "application_service" {
	environment_id = var.environment_id
	application_id = var.application_id
	service {
		type = "EXTERNAL"
		service_name = "terraform-test-0"
		port_mapping_item_list {
			port = 80
			target_port = 80
			protocol = "TCP"
		}
	}
}

`
