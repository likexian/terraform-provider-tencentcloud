package clb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"

	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudClbRedirection_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbRedirectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbRedirection_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbRedirectionExists("tencentcloud_clb_redirection.redirection_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_redirection.redirection_basic", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_redirection.redirection_basic", "source_listener_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_redirection.redirection_basic", "target_listener_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_redirection.redirection_basic", "source_rule_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_redirection.redirection_basic", "target_rule_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbRedirection_auto(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbRedirectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbRedirection_auto, tcacctest.DefaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbRedirectionExists("tencentcloud_clb_redirection.redirection_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_redirection.redirection_basic", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_redirection.redirection_basic", "source_listener_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_redirection.redirection_basic", "source_rule_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_redirection.redirection_basic", "is_auto_rewrite", "true"),
				),
			},
		},
	})
}

func testAccCheckClbRedirectionDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_redirection" {
			continue
		}
		time.Sleep(5 * time.Second)
		instance, err := clbService.DescribeRedirectionById(ctx, rs.Primary.ID)
		if instance != nil && len(*instance) > 0 && err == nil {
			return fmt.Errorf("[CHECK][CLB redirection][Destroy] check: CLB redirection still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbRedirectionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB redirection][Exists] check: CLB redirection %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB redirection][Create] check: CLB redirection id is not set")
		}
		clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := clbService.DescribeRedirectionById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance == nil || len(*instance) == 0 {
			return fmt.Errorf("[CHECK][CLB redirection][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbRedirection_basic = `
resource "tencentcloud_clb_instance" "clb_basic_redirection" {
  network_type = "OPEN"
  clb_name     = "tf-clb-redirection-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic_redirection.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic_redirection.id
  listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_listener" "listener_target" {
  clb_id        = tencentcloud_clb_instance.clb_basic_redirection.id
  port          = 44
  protocol      = "HTTP"
  listener_name = "listener_basic1"
}

resource "tencentcloud_clb_listener_rule" "rule_target" {
  clb_id              = tencentcloud_clb_instance.clb_basic_redirection.id
  listener_id         = tencentcloud_clb_listener.listener_target.listener_id
  domain              = "abcd.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}

resource "tencentcloud_clb_redirection" "redirection_basic" {
  clb_id             = tencentcloud_clb_instance.clb_basic_redirection.id
  source_listener_id = tencentcloud_clb_listener.listener_basic.listener_id
  target_listener_id = tencentcloud_clb_listener.listener_target.listener_id
  source_rule_id     = tencentcloud_clb_listener_rule.rule_basic.rule_id
  target_rule_id     = tencentcloud_clb_listener_rule.rule_target.rule_id
    is_auto_rewrite	 = false
}
`

const testAccClbRedirection_auto = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-redirection-auto"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 443
  protocol      = "HTTPS"
  listener_name = "listener_basic"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}


resource "tencentcloud_clb_redirection" "redirection_basic" {
  clb_id             = tencentcloud_clb_instance.clb_basic.id
  target_listener_id = tencentcloud_clb_listener.listener_basic.listener_id
  target_rule_id     = tencentcloud_clb_listener_rule.rule_basic.rule_id
  is_auto_rewrite	 = true
}
`
