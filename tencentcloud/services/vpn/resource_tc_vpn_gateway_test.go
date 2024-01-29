package vpn_test

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
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func init() {
	resource.AddTestSweepers("tencentcloud_vpn_gateway", &resource.Sweeper{
		Name: "tencentcloud_vpn_gateway",
		F:    testSweepVpnGateway,
	})
}

func testSweepVpnGateway(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta)

	vpcService := svcvpc.NewVpcService(client.GetAPIV3Conn())

	instances, err := vpcService.DescribeVpnGwByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range instances {
		if !tccommon.CheckResourcePersist(*v.VpnGatewayName, *v.CreatedTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.VpnGatewayId,
				Name: *v.VpnGatewayName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:        *v.VpnGatewayId,
			Name:      *v.VpnGatewayName,
			CreatTime: *v.CreatedTime,
		})
	}
	tccommon.ProcessScanCloudResources(resources, nonKeepResources, "vpn", "gateway")

	for _, v := range instances {

		vpnGwId := *v.VpnGatewayId
		vpnName := *v.VpnGatewayName
		now := time.Now()
		createTime := tccommon.StringToTime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()

		if strings.HasPrefix(vpnName, tcacctest.KeepResource) || strings.HasPrefix(vpnName, tcacctest.DefaultResource) {
			continue
		}

		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = vpcService.DeleteVpnGateway(ctx, vpnGwId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", vpnGwId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudVpnGatewayResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("tencentcloud_vpn_gateway.my_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "tags.test", "tf"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway.my_cgw", "state"),
				),
			},
			{
				Config: testAccVpnGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("tencentcloud_vpn_gateway.my_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "name", "terraform_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "bandwidth", "5"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "tags.test", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway.my_cgw", "state"),
				),
			},
			{
				Config: testAccCcnVpnGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("tencentcloud_vpn_gateway.my_ccn_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_ccn_cgw", "name", "terraform_ccn_vpngw_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_ccn_cgw", "bandwidth", "5"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_ccn_cgw", "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_ccn_cgw", "tags.test", "tf-ccn-vpngw"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway.my_ccn_cgw", "state"),
				),
			},
		},
	})
}

func testAccCheckVpnGatewayDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_gateway" {
			continue
		}
		request := vpc.NewDescribeVpnGatewaysRequest()
		request.VpnGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeVpnGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeVpnGateways(request)
			if e != nil {
				ee, ok := e.(*errors.TencentCloudSDKError)
				if !ok {
					return tccommon.RetryError(e)
				}
				if ee.Code == "ResourceNotFound" {
					log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return resource.NonRetryableError(e)
				} else {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return tccommon.RetryError(e)
				}
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n", logId, err.Error())
			ee, ok := err.(*errors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == svcvpc.VPCNotFound {
				return nil
			} else {
				return err
			}
		} else {
			if len(response.Response.VpnGatewaySet) != 0 {
				return fmt.Errorf("VPN gateway id is still exists")
			}
		}

	}
	return nil
}

func testAccCheckVpnGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN gateway instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN gateway id is not set")
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		request := vpc.NewDescribeVpnGatewaysRequest()
		request.VpnGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeVpnGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeVpnGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(response.Response.VpnGatewaySet) != 1 {
			return fmt.Errorf("VPN gateway id is not found")
		}
		return nil
	}
}

const testAccVpnGatewayConfig = `
# Create VPC
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "terraform_test"
  vpc_id    = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  bandwidth = 10

  tags = {
    test = "tf"
  }
}
`
const testAccVpnGatewayConfigUpdate = `
# Create VPC and Subnet
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "terraform_update"
  vpc_id    = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  bandwidth = 5

  tags = {
    test = "test"
  }
}
`

const testAccCcnVpnGatewayConfig = `
# Create VPNGW of CCN type
resource "tencentcloud_vpn_gateway" "my_ccn_cgw" {
  name      = "terraform_ccn_vpngw_test"
  bandwidth = 5
  type      = "CCN"

  tags = {
    test = "tf-ccn-vpngw"
  }
}
`
