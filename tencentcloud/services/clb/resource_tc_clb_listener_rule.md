Provides a resource to create a CLB listener rule.

-> **NOTE:** This resource only be applied to the HTTP or HTTPS listeners.

Example Usage

Create a single domain listener rule

```hcl
resource "tencentcloud_clb_listener_rule" "example" {
  listener_id                = "lbl-hh141sn9"
  clb_id                     = "lb-k2zjp9lv"
  domain                     = "example.com"
  url                        = "/"
  health_check_switch        = true
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  health_check_http_code     = 2
  health_check_http_path     = "/"
  health_check_http_domain   = "check.com"
  health_check_http_method   = "GET"
  certificate_ssl_mode       = "MUTUAL"
  certificate_id             = "VjANRdz8"
  certificate_ca_id          = "VfqO4zkB"
  session_expire_time        = 30
  scheduler                  = "WRR"
}
```

Create a listener rule for domain lists

```hcl
resource "tencentcloud_clb_listener_rule" "example" {
  listener_id                = "lbl-2qzcv7oq"
  clb_id                     = "lb-l6cp6jt4"
  domains                    = ["example1.com", "example2.com"]
  url                        = "/"
  health_check_switch        = true
  health_check_interval_time = 5
  health_check_health_num    = 3
  health_check_unhealth_num  = 3
  health_check_port          = 8080
  health_check_time_out      = 2
  health_check_http_code     = 15
  health_check_type          = "HTTP"
  health_check_http_path     = "/"
  health_check_http_domain   = "check.com"
  health_check_http_method   = "GET"
  scheduler                  = "WRR"
  multi_cert_info {
    ssl_mode = "UNIDIRECTIONAL"
    cert_id_list = [
      "LCYouprI",
      "JVO1alRN",
    ]
  }
}
```

Import

CLB listener rule can be imported using the id (version >= 1.47.0), e.g.

```
$ terraform import tencentcloud_clb_listener_rule.example lb-k2zjp9lv#lbl-hh141sn9#loc-agg236ys
```
