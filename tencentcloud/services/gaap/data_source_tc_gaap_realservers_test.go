package gaap_test

import (
	"regexp"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudGaapRealservers_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapRealserversBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_realservers.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_realservers.foo", "realservers.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_realservers.foo", "realservers.0.ip", "1.1.11.19"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.project_id"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapRealservers_domain(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapRealserversDomain,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_realservers.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_realservers.foo", "realservers.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_realservers.foo", "realservers.0.domain", "www.qq343.com"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.project_id"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudGaapRealservers_name(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapRealserversName,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_realservers.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_gaap_realservers.foo", "realservers.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_realservers.foo", "realservers.0.name", "ci-test-gaap-realserver"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_realservers.foo", "realservers.0.project_id"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudGaapRealserversBasic = `
resource tencentcloud_gaap_realserver "foo" {
  ip   = "1.1.11.19"
  name = "ci-test-gaap-realserver"
}

data tencentcloud_gaap_realservers "foo" {
  ip = tencentcloud_gaap_realserver.foo.ip
}
`

const TestAccDataSourceTencentCloudGaapRealserversDomain = `
resource tencentcloud_gaap_realserver "foo" {
  domain = "www.qq343.com"
  name   = "ci-test-gaap-realserver"
}

data tencentcloud_gaap_realservers "foo" {
  domain = tencentcloud_gaap_realserver.foo.domain
}
`

const TestAccDataSourceTencentCloudGaapRealserversName = `
resource tencentcloud_gaap_realserver "foo" {
  domain = "www.tencent343.com"
  name   = "ci-test-gaap-realserver"
}

data tencentcloud_gaap_realservers "foo" {
  name = tencentcloud_gaap_realserver.foo.name
}
`
