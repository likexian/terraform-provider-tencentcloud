package clb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"

	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudClbListener_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_basic", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "listener_name", "listener_basic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "port", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "scheduler", "WRR"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListenerResource_PortRange(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_portRange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_port_range"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_port_range", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_port_range", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_port_range", "listener_name", "listener_port_range"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_port_range", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_port_range", "port", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_port_range", "end_port", "6"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_port_range", "scheduler", "WRR"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_port_range",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListener_tcp_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "100"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "target_type", "NODE"),
				),
			},
			{
				Config: testAccClbListener_tcp_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "60"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "target_type", "TARGETGROUP"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_tcp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListenerTCPWithTCP(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_tcp_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "100"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_type", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_port", "200"),
				),
			},
			{
				Config: testAccClbListener_tcp_update_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "60"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_type", "TCP"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_tcp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListenerTCPWithHTTP(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_tcp_http,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "100"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_type", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_code", "16"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_path", "/"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_domain", "www.tencent.com"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_method", "HEAD"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_version", "HTTP/1.1"),
				),
			},
			{
				Config: testAccClbListener_tcp_update_http,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "60"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_type", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_port", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_code", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_path", ""),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_domain", ""),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_method", "GET"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_http_version", "HTTP/1.0"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_tcp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListenerTCPWithCustomer(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_tcp_customer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "100"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_type", "CUSTOM"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_context_type", "HEX"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_send_context", "0123456789ABCDEF"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_recv_context", "ABCD"),
				),
			},
			{
				Config: testAccClbListener_tcp_customer_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "60"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_type", "CUSTOM"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_context_type", "TEXT"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_send_context", "/get/test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_recv_context", "http_1xx"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_tcp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListener_https(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbListener_https, tcacctest.DefaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_https"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_https", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "listener_name", "listener_https"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "port", "77"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_id", tcacctest.DefaultSshCertificate),
				),
			},
			{
				Config: fmt.Sprintf(testAccClbListener_https_update, tcacctest.DefaultSshCertificateB),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_https"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_https", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "listener_name", "listener_https_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "port", "33"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_id", tcacctest.DefaultSshCertificateB),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_https",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClbListener_tcpssl(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccClbListener_tcpssl, tcacctest.DefaultSshCertificate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcpssl"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcpssl", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "protocol", "TCP_SSL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "listener_name", "listener_tcpssl"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "certificate_id", tcacctest.DefaultSshCertificate),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_time_out", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_interval_time", "100"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_health_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_unhealth_num", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccClbListener_tcpssl_update, tcacctest.DefaultSshCertificateB),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcpssl"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcpssl", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "protocol", "TCP_SSL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "listener_name", "listener_tcpssl_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "certificate_id", tcacctest.DefaultSshCertificateB),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_switch", "true"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_time_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcpssl", "health_check_unhealth_num", "3"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_tcpssl",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccTencentCloudClbListener_udp(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_udp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_basic_udp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_basic_udp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "listener_name", "listener_basic_udp"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "port", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "scheduler", "WRR"),
				),
			},
			{
				Config: testAccClbListener_udpUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_basic_udp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_basic_udp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "protocol", "UDP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "listener_name", "listener_basic_udp_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "port", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic_udp", "health_check_type", "CUSTOM"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_basic_udp",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccCheckClbListenerDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_listener" {
			continue
		}
		time.Sleep(5 * time.Second)
		resourceId := rs.Primary.ID
		items := strings.Split(resourceId, tccommon.FILED_SP)
		itemLength := len(items)
		listenerId := items[itemLength-1]
		clbId := rs.Primary.Attributes["clb_id"]
		if itemLength == 2 && clbId != "" {
			clbId = items[0]
		}
		instance, err := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if instance != nil && err == nil {
			return fmt.Errorf("[CHECK][CLB listener][Destroy] check: CLB listener still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbListenerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB listener][Exists] check: CLB listener %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB listener][Exists] check: CLB listener id is not set")
		}
		clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		resourceId := rs.Primary.ID
		items := strings.Split(resourceId, tccommon.FILED_SP)
		itemLength := len(items)
		listenerId := items[itemLength-1]
		clbId := rs.Primary.Attributes["clb_id"]
		if itemLength == 2 && clbId != "" {
			clbId = items[0]
		}
		instance, err := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CLB listener][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbListener_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  port                = 1
  protocol            = "TCP"
  listener_name       = "listener_basic"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
`
const testAccClbListener_udp = `
resource "tencentcloud_clb_instance" "clb_basic_udp" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-basic"
}

resource "tencentcloud_clb_listener" "listener_basic_udp" {
  clb_id              = tencentcloud_clb_instance.clb_basic_udp.id
  port                = 1
  protocol            = "UDP"
  listener_name       = "listener_basic_udp"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
`
const testAccClbListener_udpUpdate = `
resource "tencentcloud_clb_instance" "clb_basic_udp" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-basic"
}

resource "tencentcloud_clb_listener" "listener_basic_udp" {
  clb_id              = tencentcloud_clb_instance.clb_basic_udp.id
  port                = 1
  protocol            = "UDP"
  listener_name       = "listener_basic_udp_update"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
  health_check_type   =  "CUSTOM"
}
`
const testAccClbListener_portRange = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-port-range"
}

resource "tencentcloud_clb_listener" "listener_port_range" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  port                = 1
  end_port            = 6
  protocol            = "TCP"
  listener_name       = "listener_port_range"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "NODE"
}
`

const testAccClbListener_tcp = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-tcp"
}

resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  target_type         = "NODE"
}
`

const testAccClbListener_tcp_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-tcp"
}

resource "tencentcloud_clb_listener" "listener_tcp"{
  clb_id = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp_update"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 20
  health_check_interval_time = 200
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  session_expire_time        = 60
  scheduler                  = "WRR"
  target_type         = "TARGETGROUP"
}
`

const testAccClbListener_tcpssl = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-tcpssl"
}

resource "tencentcloud_clb_listener" "listener_tcpssl" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcpssl"
  port                       = 44
  protocol                   = "TCP_SSL"
  certificate_ssl_mode       = "UNIDIRECTIONAL"
  certificate_id             = "%s"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  scheduler                  = "WRR"
  target_type         = "TARGETGROUP"
}
`
const testAccClbListener_tcpssl_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-tcpssl"
}

resource "tencentcloud_clb_listener" "listener_tcpssl"{
  clb_id = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcpssl_update"
  port                       = 44
  protocol                   = "TCP_SSL"
  certificate_ssl_mode       = "UNIDIRECTIONAL"
  certificate_id             = "%s"
  health_check_switch        = true
  health_check_time_out      = 20
  health_check_interval_time = 200
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  scheduler                  = "WRR"
  target_type         = "TARGETGROUP"
}
`
const testAccClbListener_https = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-https"
}

resource "tencentcloud_clb_listener" "listener_https" {
  clb_id               = tencentcloud_clb_instance.clb_basic.id
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
  sni_switch           = false
}
`

const testAccClbListener_https_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-https"
}

resource "tencentcloud_clb_listener" "listener_https" {
  clb_id               = tencentcloud_clb_instance.clb_basic.id
  listener_name        = "listener_https_update"
  port                 = 33
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "%s"
  sni_switch           = false
}
`

const clb_tcp_instance = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-tcp-tcp"
}
`

const testAccClbListener_tcp_tcp = clb_tcp_instance + `
resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "TCP"
  health_check_port          = 200
}
`

const testAccClbListener_tcp_update_tcp = clb_tcp_instance + `
resource "tencentcloud_clb_listener" "listener_tcp"{
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp_update"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 20
  health_check_interval_time = 200
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  session_expire_time        = 60
  scheduler                  = "WRR"
  health_check_type          = "TCP"
  health_check_port          = 333
}
`

const clb_http_instance = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-tcp-http"
}
`

const testAccClbListener_tcp_http = clb_http_instance + `
resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "HTTP"
  health_check_http_domain   = "www.tencent.com"
  health_check_http_code     = 16
  health_check_http_version  = "HTTP/1.1"
  health_check_http_method   = "HEAD"
  health_check_http_path     = "/"
}
`

const testAccClbListener_tcp_update_http = clb_http_instance + `
resource "tencentcloud_clb_listener" "listener_tcp"{
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp_update"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 20
  health_check_interval_time = 200
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  session_expire_time        = 60
  scheduler                  = "WRR"
  health_check_port          = 200
  health_check_type          = "HTTP"
  health_check_http_code     = 2
  health_check_http_version  = "HTTP/1.0"
  health_check_http_method   = "GET"
}
`

const clb_customer_instance = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-tcp-customer"
}
`

const testAccClbListener_tcp_customer = clb_customer_instance + `
resource "tencentcloud_clb_listener" "listener_tcp"{
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
  health_check_type          = "CUSTOM"
  health_check_context_type  = "HEX"
  health_check_send_context  = "0123456789ABCDEF"
  health_check_recv_context  = "ABCD"
}
`

const testAccClbListener_tcp_customer_update = clb_customer_instance + `
resource "tencentcloud_clb_listener" "listener_tcp"{
  clb_id                     = tencentcloud_clb_instance.clb_basic.id
  listener_name              = "listener_tcp_update"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = true
  health_check_time_out      = 20
  health_check_interval_time = 200
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  session_expire_time        = 60
  scheduler                  = "WRR"
  health_check_type          = "CUSTOM"
  health_check_context_type  = "TEXT"
  health_check_send_context  = "/get/test"
  health_check_recv_context  = "http_1xx"
}
`
