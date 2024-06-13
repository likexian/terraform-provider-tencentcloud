package teo_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_teo_zone
	resource.AddTestSweepers("tencentcloud_teo_zone", &resource.Sweeper{
		Name: "tencentcloud_teo_zone",
		F:    testSweepZone,
	})
}

func testSweepZone(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(region)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcteo.NewTeoService(client)

	zone, err := service.DescribeTeoZone(ctx, "")
	if err != nil {
		return err
	}

	if zone == nil {
		return nil
	}

	err = service.DeleteTeoZoneById(ctx, *zone.ZoneId)
	if err != nil {
		return err
	}

	return nil
}

// go test -test.run TestAccTencentCloudTeoZone_basic -v
func TestAccTencentCloudTeoZone_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PRIVATE) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoZone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("tencentcloud_teo_zone.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "zone_name", "tf-teo.xyz"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "alias_zone_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "paused", "false"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "plan_id", "edgeone-2kfv1h391n6w"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "type", "partial"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "tags.勿动", "TF测试"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "tags.占用人", "arunma"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "ownership_verification.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "ownership_verification.0.dns_verification.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "ownership_verification.0.dns_verification.0.record_type", "TXT"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone.basic", "ownership_verification.0.dns_verification.0.record_value"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone.basic", "ownership_verification.0.dns_verification.0.subdomain"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone.basic", "status"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_zone.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoZoneUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists("tencentcloud_teo_zone.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "zone_name", "tf-teo.xyz"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "area", "overseas"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "alias_zone_name", "tf-test-up"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "paused", "true"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "plan_id", "edgeone-2kfv1h391n6w"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "type", "partial"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "tags.勿动", "TF测试"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "tags.占用人", "arunma"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "ownership_verification.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "ownership_verification.0.dns_verification.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_teo_zone.basic", "ownership_verification.0.dns_verification.0.record_type", "TXT"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone.basic", "ownership_verification.0.dns_verification.0.record_value"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_zone.basic", "ownership_verification.0.dns_verification.0.subdomain"),
				),
			},
		},
	})
}

func testAccCheckZoneDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_teo_zone" {
			continue
		}

		agents, err := service.DescribeTeoZone(ctx, rs.Primary.ID)
		if agents != nil {
			return fmt.Errorf("zone %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckZoneExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := svcteo.NewTeoService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		agents, err := service.DescribeTeoZone(ctx, rs.Primary.ID)
		if agents == nil {
			return fmt.Errorf("zone %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoZoneVar = `
variable "zone_name" {
  default = "tf-teo.xyz"
}

variable "plan_id" {
  default = "edgeone-2kfv1h391n6w"
}`

const testAccTeoZone = testAccTeoZoneVar + `

resource "tencentcloud_teo_zone" "basic" {
	area            = "overseas"
	alias_zone_name = "tf-test"
	paused          = false
	plan_id         = var.plan_id
	tags = {
	  "勿动"  = "TF测试"
	  "占用人" = "arunma"
	}
	type      = "partial"
	zone_name = var.zone_name
  }

`

const testAccTeoZoneUp = testAccTeoZoneVar + `

resource "tencentcloud_teo_zone" "basic" {
	area            = "overseas"
	alias_zone_name = "tf-test-up"
	paused          = true
	plan_id         = var.plan_id
	tags = {
	  "勿动"  = "TF测试"
	  "占用人" = "arunma"
	}
	type      = "partial"
	zone_name = var.zone_name
  }

`
