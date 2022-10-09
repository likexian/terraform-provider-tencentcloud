package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_teo_zone
	resource.AddTestSweepers("tencentcloud_teo_application_proxy", &resource.Sweeper{
		Name: "tencentcloud_teo_application_proxy",
		F:    testSweepApplicationProxy,
	})
}

func testSweepApplicationProxy(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(region)
	client := cli.(*TencentCloudClient).apiV3Conn
	service := TeoService{client}

	zoneId := defaultZoneId

	for {
		proxy, err := service.DescribeTeoApplicationProxy(ctx, zoneId, "")
		if err != nil {
			return err
		}

		if proxy == nil {
			return nil
		}

		err = service.DeleteTeoApplicationProxyById(ctx, zoneId, *proxy.ProxyId)
		if err != nil {
			return err
		}
	}
}

// go test -i; go test -test.run TestAccTencentCloudTeoApplicationProxy_basic -v
func TestAccTencentCloudTeoApplicationProxy_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplicationProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoApplicationProxy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationProxyExists("tencentcloud_teo_application_proxy.basic"),
					//resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "zone_name", "tf-teo.com"),
					//resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "plan_type", "ent_with_bot"),
					//resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "type", "full"),
					//resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "paused", "false"),
					//resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "cname_speed_up", "enabled"),
					//resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "vanity_name_servers.#", "1"),
					//resource.TestCheckResourceAttr("tencentcloud_teo_application_proxy.basic", "vanity_name_servers.0.switch", "on"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_application_proxy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckApplicationProxyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_application_proxy" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		agents, err := service.DescribeTeoApplicationProxy(ctx, zoneId, proxyId)
		if agents != nil {
			return fmt.Errorf("zone ApplicationProxy %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckApplicationProxyExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		zoneId := idSplit[0]
		proxyId := idSplit[1]

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		agents, err := service.DescribeTeoApplicationProxy(ctx, zoneId, proxyId)
		if agents == nil {
			return fmt.Errorf("zone ApplicationProxy %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoApplicationProxyVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}`

const testAccTeoApplicationProxy = testAccTeoApplicationProxyVar + `

resource "tencentcloud_teo_application_proxy" "basic" {
  zone_id = var.zone_id

  accelerate_type      = 1
  security_type        = 1
  plat_type            = "domain"
  proxy_name           = "test-instance"
  proxy_type           = "instance"
  session_persist_time = 2400
}

`
