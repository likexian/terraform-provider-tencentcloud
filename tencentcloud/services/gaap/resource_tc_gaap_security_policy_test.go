package gaap_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcgaap "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/gaap"

	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudGaapSecurityPolicy_basic(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityPolicyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityPolicyExists("tencentcloud_gaap_security_policy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_security_policy.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityPolicy_disable(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityPolicyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityPolicyExists("tencentcloud_gaap_security_policy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "ACCEPT"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
				),
			},
			{
				Config: testAccGaapSecurityPolicyDisable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityPolicyExists("tencentcloud_gaap_security_policy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "false"),
				),
			},
		},
	})
}

func TestAccTencentCloudGaapSecurityPolicy_drop(t *testing.T) {
	id := new(string)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGaapSecurityPolicyDestroy(id),
		Steps: []resource.TestStep{
			{
				Config: testAccGaapSecurityPolicyDrop,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaapSecurityPolicyExists("tencentcloud_gaap_security_policy.foo", id),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "action", "DROP"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_security_policy.foo", "enable", "true"),
				),
			},
		},
	})
}

func testAccCheckGaapSecurityPolicyExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no listener ID is set")
		}

		service := svcgaap.NewGaapService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		_, _, _, exist, err := service.DescribeSecurityPolicy(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if !exist {
			return fmt.Errorf("security policy not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

func testAccCheckGaapSecurityPolicyDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcgaap.NewGaapService(client)

		_, _, _, exist, err := service.DescribeSecurityPolicy(context.TODO(), *id)
		if err != nil {
			return err
		}

		if exist {
			return errors.New("security policy still exists")
		}

		return nil
	}
}

var testAccGaapSecurityPolicyBasic = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}
`, "link-asnxs27x") //tcacctest.DefaultGaapProxyId)

var testAccGaapSecurityPolicyDisable = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
  enable   = false
}
`, "link-asnxs27x") //tcacctest.DefaultGaapProxyId)

var testAccGaapSecurityPolicyDrop = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "DROP"
}
`, "link-asnxs27x") //tcacctest.DefaultGaapProxyId)
