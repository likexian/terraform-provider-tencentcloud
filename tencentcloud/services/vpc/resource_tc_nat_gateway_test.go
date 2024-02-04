package vpc_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func init() {
	resource.AddTestSweepers("tencentcloud_nat", &resource.Sweeper{
		Name: "tencentcloud_nat",
		F:    testSweepNatInstance,
	})
}

func testSweepNatInstance(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()

	vpcService := svcvpc.NewVpcService(client)

	instances, err := vpcService.DescribeNatGatewayByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range instances {
		if !tccommon.CheckResourcePersist(*v.NatGatewayName, *v.CreatedTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.NatGatewayId,
				Name: *v.NatGatewayName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:        *v.NatGatewayId,
			Name:      *v.NatGatewayName,
			CreatTime: *v.CreatedTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateNatGateway")

	for _, v := range instances {
		instanceId := *v.NatGatewayId
		instanceName := v.NatGatewayName

		now := time.Now()

		createTime := tccommon.StringToTime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()
		if instanceName != nil {
			if strings.HasPrefix(*instanceName, tcacctest.KeepResource) || strings.HasPrefix(*instanceName, tcacctest.DefaultResource) {
				continue
			}
		}

		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteNatGateway(ctx, instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", instanceId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudNatGateway_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("tencentcloud_nat_gateway.my_nat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "max_concurrent", "3000000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "bandwidth", "500"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "assigned_eip_set.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "tags.tf", "test"),
				),
			},
			{
				Config: testAccNatGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists("tencentcloud_nat_gateway.my_nat"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "name", "new_name"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "max_concurrent", "10000000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "bandwidth", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "assigned_eip_set.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_nat_gateway.my_nat", "tags.tf", "teest"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_nat_gateway" {
			continue
		}
		request := vpc.NewDescribeNatGatewaysRequest()
		request.NatGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeNatGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewaySet) != 0 {
			return fmt.Errorf("nat gateway id is still exists")
		}

	}
	return nil
}

func testAccCheckNatGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("nat gateway instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("nat gateway id is not set")
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		request := vpc.NewDescribeNatGatewaysRequest()
		request.NatGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeNatGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeNatGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(response.Response.NatGatewaySet) != 1 {
			return fmt.Errorf("nat gateway id is not found")
		}
		return nil
	}
}

const testAccNatGatewayConfig = `
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  name           = "terraform_test"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
  tags = {
    tf = "test"
  }
}
`
const testAccNatGatewayConfigUpdate = `
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
# Create EIP 
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "new_eip" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  name           = "new_name"
  max_concurrent = 10000000
  bandwidth      = 1000

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.new_eip.public_ip,
  ]
  tags = {
    tf = "teest"
  }
}
`
