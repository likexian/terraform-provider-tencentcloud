---
subcategory: "Content Delivery Network(CDN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdn_domain"
sidebar_current: "docs-tencentcloud-resource-cdn_domain"
description: |-
  Provides a resource to create a CDN domain.
---

# tencentcloud_cdn_domain

Provides a resource to create a CDN domain.

~> **NOTE:** To disable most of configuration with switch, just modify switch argument to off instead of remove the whole block

## Example Usage

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain         = "xxxx.com"
  service_type   = "web"
  area           = "mainland"
  full_url_cache = false

  origin {
    origin_type          = "ip"
    origin_list          = ["127.0.0.1"]
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"

    force_redirect {
      switch               = "on"
      redirect_type        = "http"
      redirect_status_code = 302
    }
  }

  tags = {
    hello = "world"
  }
}
```

### Example Usage of cdn uses cache and request headers

```hcl
resource "tencentcloud_cdn_domain" "foo" {
  domain       = "xxxx.com"
  service_type = "web"
  area         = "mainland"
  # full_url_cache = true # Deprecated, use cache_key below.
  cache_key {
    full_url_cache = "on"
  }
  range_origin_switch = "off"

  rule_cache {
    cache_time      = 10000
    no_cache_switch = "on"
    re_validate     = "on"
  }

  request_header {
    switch = "on"

    header_rules {
      header_mode  = "add"
      header_name  = "tf-header-name"
      header_value = "tf-header-value"
      rule_type    = "all"
      rule_paths   = ["*"]
    }
  }

  origin {
    origin_type          = "ip"
    origin_list          = ["127.0.0.1"]
    origin_pull_protocol = "follow"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"

    force_redirect {
      switch               = "on"
      redirect_type        = "http"
      redirect_status_code = 302
    }
  }

  tags = {
    hello = "world"
  }
}
```

### Example Usage of COS bucket url as origin

```hcl
resource "tencentcloud_cos_bucket" "bucket" {
  # Bucket format should be [custom name]-[appid].
  bucket = "demo-bucket-1251234567"
  acl    = "private"
}

# Create cdn domain
resource "tencentcloud_cdn_domain" "cdn" {
  domain       = "abc.com"
  service_type = "web"
  area         = "mainland"
  # full_url_cache = false # Deprecated
  cache_key {
    full_url_cache = "off"
  }

  origin {
    origin_type          = "cos"
    origin_list          = [tencentcloud_cos_bucket.bucket.cos_bucket_url]
    server_name          = tencentcloud_cos_bucket.bucket.cos_bucket_url
    origin_pull_protocol = "follow"
    cos_private_access   = "on"
  }

  https_config {
    https_switch         = "off"
    http2_switch         = "off"
    ocsp_stapling_switch = "off"
    spdy_switch          = "off"
    verify_client        = "off"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Name of the acceleration domain.
* `origin` - (Required, List) Origin server configuration. It's a list and consist of at most one item.
* `service_type` - (Required, String, ForceNew) Acceleration domain name service type. `web`: static acceleration, `download`: download acceleration, `media`: streaming media VOD acceleration, `hybrid`: hybrid acceleration, `dynamic`: dynamic acceleration.
* `area` - (Optional, String) Domain name acceleration region. `mainland`: acceleration inside mainland China, `overseas`: acceleration outside mainland China, `global`: global acceleration. Overseas acceleration service must be enabled to use overseas acceleration and global acceleration.
* `authentication` - (Optional, List) Specify timestamp hotlink protection configuration, NOTE: only one type can choose for the sub elements.
* `aws_private_access` - (Optional, List) Access authentication for S3 origin.
* `band_width_alert` - (Optional, List) Bandwidth cap configuration.
* `cache_key` - (Optional, List) Cache key configuration (Ignore Query String configuration). NOTE: All of `full_url_cache` default value is `on`.
* `compression` - (Optional, List) Smart compression configurations.
* `downstream_capping` - (Optional, List) Downstream capping configuration.
* `error_page` - (Optional, List) Error page configurations.
* `explicit_using_dry_run` - (Optional, Bool) Used for validate only by store arguments to request json string as expected, WARNING: if set to `true`, NO Cloud Api will be invoked but store as local data, do not use this argument unless you really know what you are doing.
* `follow_redirect_switch` - (Optional, String) 301/302 redirect following switch, available values: `on`, `off` (default).
* `full_url_cache` - (Optional, Bool, **Deprecated**) Use `cache_key` -> `full_url_cache` instead. Whether to enable full-path cache. Default value is `true`.
* `https_config` - (Optional, List) HTTPS acceleration configuration. It's a list and consist of at most one item.
* `hw_private_access` - (Optional, List) Access authentication for OBS origin.
* `ip_filter` - (Optional, List) Specify Ip filter configurations.
* `ip_freq_limit` - (Optional, List) Specify Ip frequency limit configurations.
* `ipv6_access_switch` - (Optional, String) ipv6 access configuration switch. Only available when area set to `mainland`. Valid values are `on` and `off`. Default value is `off`.
* `max_age` - (Optional, List) Browser cache configuration. (This feature is in beta and not generally available yet).
* `offline_cache_switch` - (Optional, String) Offline cache switch, available values: `on`, `off` (default).
* `origin_pull_optimization` - (Optional, List) Cross-border linkage optimization configuration. (This feature is in beta and not generally available yet).
* `origin_pull_timeout` - (Optional, List) Cross-border linkage optimization configuration.
* `oss_private_access` - (Optional, List) Access authentication for OSS origin.
* `others_private_access` - (Optional, List) Object storage back-to-source authentication of other vendors.
* `post_max_size` - (Optional, List) Maximum post size configuration.
* `project_id` - (Optional, Int) The project CDN belongs to, default to 0.
* `qn_private_access` - (Optional, List) Access authentication for OBS origin.
* `quic_switch` - (Optional, String) QUIC switch, available values: `on`, `off` (default).
* `range_origin_switch` - (Optional, String) Sharding back to source configuration switch. Valid values are `on` and `off`. Default value is `on`.
* `referer` - (Optional, List) Referer configuration.
* `request_header` - (Optional, List) Request header configuration. It's a list and consist of at most one item.
* `response_header_cache_switch` - (Optional, String) Response header cache switch, available values: `on`, `off` (default).
* `response_header` - (Optional, List) Response header configurations.
* `rule_cache` - (Optional, List) Advanced path cache configuration.
* `seo_switch` - (Optional, String) SEO switch, available values: `on`, `off` (default).
* `specific_config_mainland` - (Optional, String) Specific configuration for mainland, NOTE: Both specifying full schema or using it is superfluous, please use cloud api parameters json passthroughs, check the [Data Types](https://www.tencentcloud.com/document/api/228/31739#MainlandConfig) for more details.
* `specific_config_overseas` - (Optional, String) Specific configuration for oversea, NOTE: Both specifying full schema or using it is superfluous, please use cloud api parameters json passthroughs, check the [Data Types](https://www.tencentcloud.com/document/api/228/31739#OverseaConfig) for more details.
* `status_code_cache` - (Optional, List) Status code cache configurations.
* `tags` - (Optional, Map) Tags of cdn domain.
* `video_seek_switch` - (Optional, String) Video seek switch, available values: `on`, `off` (default).

The `authentication` object supports the following:

* `switch` - (Optional, String) Authentication switching, available values: `on`, `off`.
* `type_a` - (Optional, List) Timestamp hotlink protection mode A configuration.
* `type_b` - (Optional, List) Timestamp hotlink protection mode B configuration. NOTE: according to upgrading of TencentCloud Platform, TypeB is unavailable for now.
* `type_c` - (Optional, List) Timestamp hotlink protection mode C configuration.
* `type_d` - (Optional, List) Timestamp hotlink protection mode D configuration.

The `aws_private_access` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `access_key` - (Optional, String) Access ID.
* `bucket` - (Optional, String) Bucket.
* `region` - (Optional, String) Region.
* `secret_key` - (Optional, String) Key.

The `band_width_alert` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `alert_percentage` - (Optional, Int) Alert percentage.
* `alert_switch` - (Optional, String) Switch alert.
* `bps_threshold` - (Optional, Int) threshold of bps.
* `counter_measure` - (Optional, String) Counter measure.
* `metric` - (Optional, String) Metric.
* `statistic_item` - (Optional, List) Specify statistic item configuration.

The `cache_key` object supports the following:

* `full_url_cache` - (Optional, String) Whether to enable full-path cache, values `on` (DEFAULT ON), `off`.
* `ignore_case` - (Optional, String) Specifies whether the cache key is case sensitive.
* `key_rules` - (Optional, List) Path-specific cache key configuration.
* `query_string` - (Optional, List) Request parameter contained in CacheKey.

The `cache_rules` object of `status_code_cache` supports the following:

* `cache_time` - (Required, Int) Status code cache expiration time (in seconds).
* `status_code` - (Required, String) Code of status cache. available values: `403`, `404`.

The `capping_rules` object of `downstream_capping` supports the following:

* `kbps_threshold` - (Required, Int) Capping rule kbps threshold.
* `rule_paths` - (Required, List) List of capping rule path.
* `rule_type` - (Required, String) Capping rule type.

The `client_certificate_config` object of `https_config` supports the following:

* `certificate_content` - (Required, String) Client Certificate PEM format, requires Base64 encoding.

The `compression_rules` object of `compression` supports the following:

* `algorithms` - (Required, List) List of algorithms, available: `gzip` and `brotli`.
* `compress` - (Required, Bool) Must be set as true, enables compression.
* `max_length` - (Required, Int) The maximum file size to trigger compression (in bytes).
* `min_length` - (Required, Int) The minimum file size to trigger compression (in bytes).
* `file_extensions` - (Optional, List) List of file extensions like `jpg`, `txt`.
* `rule_paths` - (Optional, List) List of rule paths for each `rule_type`: `*` for `all`, file ext like `jpg` for `file`, `/dir/like/` for `directory` and `/path/index.html` for `path`.
* `rule_type` - (Optional, String) Rule type, available: `all`, `file`, `directory`, `path`, `contentType`.

The `compression` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `compression_rules` - (Optional, List) List of compression rules.

The `downstream_capping` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `capping_rules` - (Optional, List) List of capping rule.

The `error_page` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `page_rules` - (Optional, List) List of error page rule.

The `filter_rules` object of `ip_filter` supports the following:

* `filter_type` - (Required, String) Ip filter `blacklist`/`whitelist` type of filter rules.
* `filters` - (Required, List) Ip filter rule list, supports IPs in X.X.X.X format, or /8, /16, /24 format IP ranges. Up to 50 allowlists or blocklists can be entered.
* `rule_paths` - (Required, List) Content list for each `rule_type`: `*` for `all`, file ext like `jpg` for `file`, `/dir/like/` for `directory` and `/path/index.html` for `path`.
* `rule_type` - (Required, String) Ip filter rule type of filter rules, available: `all`, `file`, `directory`, `path`.

The `force_redirect` object of `https_config` supports the following:

* `carry_headers` - (Optional, String) Whether to return the newly added header during force redirection. Values: `on`, `off`.
* `redirect_status_code` - (Optional, Int) Forced redirect status code. Valid values are `301` and `302`. When `switch` setting `off`, this property does not need to be set or set to `302`. Default value is `302`.
* `redirect_type` - (Optional, String) Forced redirect type. Valid values are `http` and `https`. `http` means a forced redirect from HTTPS to HTTP, `https` means a forced redirect from HTTP to HTTPS. When `switch` setting `off`, this property does not need to be set or set to `http`. Default value is `http`.
* `switch` - (Optional, String) Forced redirect configuration switch. Valid values are `on` and `off`. Default value is `off`.

The `header_rules` object of `request_header` supports the following:

* `header_mode` - (Required, String) Http header setting method. The following types are supported: `add`: add a head, if a head already exists, there will be a duplicate head, `del`: delete the head.
* `header_name` - (Required, String) Http header name.
* `header_value` - (Required, String) Http header value, optional when Mode is `del`, Required when Mode is `add`/`set`.
* `rule_paths` - (Required, List) Matching content under the corresponding type of CacheType: `all`: fill *, `file`: fill in the suffix name, such as jpg, txt, `directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html.
* `rule_type` - (Required, String) Rule type. The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, `directory`: the specified path takes effect, `path`: specify the absolute path to take effect.

The `header_rules` object of `response_header` supports the following:

* `header_mode` - (Required, String) Response header mode.
* `header_name` - (Required, String) response header name of rule.
* `header_value` - (Required, String) response header value of rule.
* `rule_paths` - (Required, List) response rule paths of rule.
* `rule_type` - (Required, String) response rule type of rule.

The `https_config` object supports the following:

* `https_switch` - (Required, String) HTTPS configuration switch. Valid values are `on` and `off`.
* `client_certificate_config` - (Optional, List) Client certificate configuration information.
* `force_redirect` - (Optional, List) Configuration of forced HTTP or HTTPS redirects.
* `http2_switch` - (Optional, String) HTTP2 configuration switch. Valid values are `on` and `off`. and default value is `off`.
* `ocsp_stapling_switch` - (Optional, String) OCSP configuration switch. Valid values are `on` and `off`. and default value is `off`.
* `server_certificate_config` - (Optional, List) Server certificate configuration information.
* `spdy_switch` - (Optional, String) Spdy configuration switch. Valid values are `on` and `off`. and default value is `off`. This parameter is for white-list customer.
* `tls_versions` - (Optional, List) Tls version settings, only support some Advanced domain names, support settings TLSv1, TLSv1.1, TLSv1.2, TLSv1.3, when modifying must open consecutive versions.
* `verify_client` - (Optional, String) Client certificate authentication feature. Valid values are `on` and `off`. and default value is `off`.

The `hw_private_access` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `access_key` - (Optional, String) Access ID.
* `bucket` - (Optional, String) Bucket.
* `secret_key` - (Optional, String) Key.

The `ip_filter` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `filter_rules` - (Optional, List) Ip filter rules, This feature is only available to selected beta customers.
* `filter_type` - (Optional, String) IP `blacklist`/`whitelist` type.
* `filters` - (Optional, List) Ip filter list, Supports IPs in X.X.X.X format, or /8, /16, /24 format IP ranges. Up to 50 allowlists or blocklists can be entered.
* `return_code` - (Optional, Int) Return code, available values: 400-499.

The `ip_freq_limit` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `qps` - (Optional, Int) Sets the limited number of requests per second, 514 will be returned for requests that exceed the limit.

The `key_rules` object of `cache_key` supports the following:

* `query_string` - (Required, List) Request parameter contained in CacheKey.
* `rule_paths` - (Required, List) List of rule paths for each `key_rules`: `/` for `index`, file ext like `jpg` for `file`, `/dir/like/` for `directory` and `/path/index.html` for `path`.
* `rule_type` - (Required, String) Rule type, available: `file`, `directory`, `path`, `index`.
* `full_url_cache` - (Optional, String) Whether to enable full-path cache, values `on` (DEFAULT ON), `off`.
* `ignore_case` - (Optional, String) Whether caches are case insensitive.
* `rule_tag` - (Optional, String) Specify rule tag, default value is `user`.

The `max_age_rules` object of `max_age` supports the following:

* `max_age_contents` - (Required, List) List of rule paths for each `max_age_type`: `*` for `all`, file ext like `jpg` for `file`, `/dir/like/` for `directory` and `/path/index.html` for `path`.
* `max_age_time` - (Required, Int) Max Age time in seconds, this can set to `0` that stands for no cache.
* `max_age_type` - (Required, String) The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, `directory`: the specified path takes effect, `path`: specify the absolute path to take effect, `index`: home page.
* `follow_origin` - (Optional, String) Whether to follow origin, values: `on`/`off`, if set to `on`, the `max_age_time` will be ignored.

The `max_age` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `max_age_rules` - (Optional, List) List of Max Age rule configuration.

The `origin_pull_optimization` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `optimization_type` - (Optional, String) Optimization type, values: `OVToCN` - Overseas to CN, `CNToOV` CN to Overseas.

The `origin_pull_timeout` object supports the following:

* `connect_timeout` - (Required, Int) The origin-pull connection timeout (in seconds). Valid range: 5-60.
* `receive_timeout` - (Required, Int) The origin-pull receipt timeout (in seconds). Valid range: 10-60.

The `origin` object supports the following:

* `origin_list` - (Required, List) Master origin server list. Valid values can be ip or domain name. When modifying the origin server, you need to enter the corresponding `origin_type`.
* `origin_type` - (Required, String) Master origin server type. The following types are supported: `domain`: Domain name, `domainv6`: IPv6 domain name, `cos`: COS bucket address, `third_party`: Third-party object storage origin, `igtm`: IGTM origin, `ip`: IP address, `ipv6`: One IPv6 address, `ip_ipv6`: Multiple IPv4 addresses and one IPv6 address, `ip_domain`: IP addresses and domain names (only available to beta users), `ip_domainv6`: Multiple IPv4 addresses and one IPv6 domain name, `ipv6_domain`: Multiple IPv6 addresses and one domain name, `ipv6_domainv6`: Multiple IPv6 addresses and one IPv6 domain name, `domain_domainv6`: Multiple IPv4 domain names and one IPv6 domain name, `ip_ipv6_domain`: Multiple IPv4 and IPv6 addresses and one domain name, `ip_ipv6_domainv6`: Multiple IPv4 and IPv6 addresses and one IPv6 domain name, `ip_domain_domainv6`: Multiple IPv4 addresses and IPv4 domain names and one IPv6 domain name, `ipv6_domain_domainv6`: Multiple IPv4 domain names and IPv6 addresses and one IPv6 domain name, `ip_ipv6_domain_domainv6`: Multiple IPv4 and IPv6 addresses and IPv4 domain names and one IPv6 domain name.
* `backup_origin_list` - (Optional, List) Backup origin server list. Valid values can be ip or domain name. When modifying the backup origin server, you need to enter the corresponding `backup_origin_type`.
* `backup_origin_type` - (Optional, String) Backup origin server type, which supports the following types: `domain`: domain name type, `ip`: IP list used as origin server, `ipv6_domain`: Multiple IPv6 addresses and one domain name, `ip_ipv6`: Multiple IPv4 addresses and one IPv6 address, `ip_ipv6_domain`: Multiple IPv4 and IPv6 addresses and one domain name.
* `backup_server_name` - (Optional, String) Host header used when accessing the backup origin server. If left empty, the ServerName of master origin server will be used by default.
* `cos_private_access` - (Optional, String) When OriginType is COS, you can specify if access to private buckets is allowed. Valid values are `on` and `off`. and default value is `off`.
* `origin_company` - (Optional, String) Object storage back to the source vendor. Required when the source station type is a third-party storage source station (third_party). Optional values include the following: `aws_s3`: AWS S3; `ali_oss`: Alibaba Cloud OSS; `hw_obs`: Huawei OBS; `qiniu_kodo`: Qiniu Cloud kodo; `others`: other vendors' object storage, only supports object storage compatible with AWS signature algorithm, such as Tencent Cloud Financial Zone COS. Example value: `hw_obs`.
* `origin_pull_protocol` - (Optional, String) Origin-pull protocol configuration. `http`: forced HTTP origin-pull, `follow`: protocol follow origin-pull, `https`: forced HTTPS origin-pull. This only supports origin server port 443 for origin-pull.
* `server_name` - (Optional, String) Host header used when accessing the master origin server. If left empty, the acceleration domain name will be used by default.

The `oss_private_access` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `access_key` - (Optional, String) Access ID.
* `bucket` - (Optional, String) Bucket.
* `region` - (Optional, String) Region.
* `secret_key` - (Optional, String) Key.

The `others_private_access` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `access_key` - (Optional, String) Access ID.
* `bucket` - (Optional, String) Bucket.
* `region` - (Optional, String) Region.
* `secret_key` - (Optional, String) Key.

The `page_rules` object of `error_page` supports the following:

* `redirect_code` - (Required, Int) Redirect code of error page rules.
* `redirect_url` - (Required, String) Redirect url of error page rules.
* `status_code` - (Required, Int) Status code of error page rules.

The `post_max_size` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `max_size` - (Optional, Int) Maximum size in MB, value range is `[1, 200]`.

The `qn_private_access` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `access_key` - (Optional, String) Access ID.
* `secret_key` - (Optional, String) Key.

The `query_string` object of `cache_key` supports the following:

* `action` - (Optional, String) Include/exclude query parameters. Values: `includeAll` (Default), `excludeAll`, `includeCustom`, `excludeCustom`.
* `reorder` - (Optional, String) Whether to sort again, values `on`, `off` (Default).
* `switch` - (Optional, String) Whether to use QueryString as part of CacheKey, values `on`, `off` (Default).
* `value` - (Optional, String) Array of included/excluded query strings (separated by `;`).

The `query_string` object of `key_rules` supports the following:

* `action` - (Optional, String) Specify key rule QS action, values: `includeCustom`, `excludeCustom`.
* `switch` - (Optional, String) Whether to use QueryString as part of CacheKey, values `on`, `off` (Default).
* `value` - (Optional, String) Array of included/excluded query strings (separated by `;`).

The `referer_rules` object of `referer` supports the following:

* `allow_empty` - (Required, Bool) Whether to allow emptpy.
* `referer_type` - (Required, String) Referer type.
* `referers` - (Required, List) Referer list.
* `rule_paths` - (Required, List) Referer rule path list.
* `rule_type` - (Required, String) Referer rule type.

The `referer` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `referer_rules` - (Optional, List) List of referer rules.

The `request_header` object supports the following:

* `header_rules` - (Optional, List) Custom request header configuration rules.
* `switch` - (Optional, String) Custom request header configuration switch. Valid values are `on` and `off`. and default value is `off`.

The `response_header` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `header_rules` - (Optional, List) List of response header rule.

The `rule_cache` object supports the following:

* `cache_time` - (Required, Int) Cache expiration time setting, the unit is second, the maximum can be set to 365 days.
* `compare_max_age` - (Optional, String) Advanced cache expiration configuration. When it is turned on, it will compare the max-age value returned by the origin site with the cache expiration time set in CacheRules, and take the minimum value to cache at the node. Valid values are `on` and `off`. Default value is `off`.
* `follow_origin_switch` - (Optional, String) Follow the source station configuration switch. Valid values are `on` and `off`.
* `heuristic_cache_switch` - (Optional, String) Specify whether to enable heuristic cache, only available while `follow_origin_switch` enabled, values: `on`, `off` (Default).
* `heuristic_cache_time` - (Optional, Int) Specify heuristic cache time in second, only available while `follow_origin_switch` and `heuristic_cache_switch` enabled.
* `ignore_cache_control` - (Optional, String) Force caching. After opening, the no-store and no-cache resources returned by the origin site will also be cached in accordance with the CacheRules rules. Valid values are `on` and `off`. Default value is `off`.
* `ignore_set_cookie` - (Optional, String) Ignore the Set-Cookie header of the origin site. Valid values are `on` and `off`. Default value is `off`. This parameter is for white-list customer.
* `no_cache_switch` - (Optional, String) Cache configuration switch. Valid values are `on` and `off`.
* `re_validate` - (Optional, String) Always check back to origin. Valid values are `on` and `off`. Default value is `off`.
* `rule_paths` - (Optional, List) Matching content under the corresponding type of CacheType: `all`: fill *, `file`: fill in the suffix name, such as jpg, txt, `directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html, `index`: fill /.
* `rule_type` - (Optional, String) Rule type. The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, `directory`: the specified path takes effect, `path`: specify the absolute path to take effect, `index`: home page.
* `switch` - (Optional, String) Cache configuration switch. Valid values are `on` and `off`.

The `server_certificate_config` object of `https_config` supports the following:

* `certificate_content` - (Optional, String) Server certificate information. This is required when uploading an external certificate, which should contain the complete certificate chain.
* `certificate_id` - (Optional, String) Server certificate ID.
* `message` - (Optional, String) Certificate remarks.
* `private_key` - (Optional, String) Server key information. This is required when uploading an external certificate.

The `statistic_item` object of `band_width_alert` supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `alert_percentage` - (Optional, Int) Alert percentage.
* `alert_switch` - (Optional, String) Switch alert.
* `bps_threshold` - (Optional, Int) threshold of bps.
* `counter_measure` - (Optional, String) Counter measure, values: `RETURN_404`, `RESOLVE_DNS_TO_ORIGIN`.
* `cycle` - (Optional, Int) Cycle of checking in minutes, values `60`, `1440`.
* `metric` - (Optional, String) Metric.
* `type` - (Optional, String) Type of statistic item.
* `unblock_time` - (Optional, Int) Time of auto unblock.

The `status_code_cache` object supports the following:

* `switch` - (Required, String) Configuration switch, available values: `on`, `off` (default).
* `cache_rules` - (Optional, List) List of cache rule.

The `type_a` object of `authentication` supports the following:

* `expire_time` - (Required, Int) Signature expiration time in second. The maximum value is 630720000.
* `file_extensions` - (Required, List) File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.
* `filter_type` - (Required, String) Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.
* `secret_key` - (Required, String) The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.
* `sign_param` - (Required, String) Signature parameter name. Only upper and lower-case letters, digits, and underscores (_) are allowed. It cannot start with a digit. Length limit: 1-100 characters.
* `backup_secret_key` - (Optional, String) Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.

The `type_b` object of `authentication` supports the following:

* `expire_time` - (Required, Int) Signature expiration time in second. The maximum value is 630720000.
* `file_extensions` - (Required, List) File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.
* `filter_type` - (Required, String) Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.
* `secret_key` - (Required, String) The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.
* `backup_secret_key` - (Optional, String) Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.

The `type_c` object of `authentication` supports the following:

* `expire_time` - (Required, Int) Signature expiration time in second. The maximum value is 630720000.
* `file_extensions` - (Required, List) File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.
* `filter_type` - (Required, String) Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.
* `secret_key` - (Required, String) The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.
* `backup_secret_key` - (Optional, String) Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.
* `time_format` - (Optional, String) Timestamp formation, available values: `dec`, `hex`.

The `type_d` object of `authentication` supports the following:

* `expire_time` - (Required, Int) Signature expiration time in second. The maximum value is 630720000.
* `file_extensions` - (Required, List) File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.
* `filter_type` - (Required, String) Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.
* `secret_key` - (Required, String) The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.
* `backup_secret_key` - (Optional, String) Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.
* `time_format` - (Optional, String) Timestamp formation, available values: `dec`, `hex`.
* `time_param` - (Optional, String) Timestamp parameter name. Only upper and lower-case letters, digits, and underscores (_) are allowed. It cannot start with a digit. Length limit: 1-100 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME address of domain name.
* `create_time` - Creation time of domain name.
* `dry_run_create_result` - Used for store `dry_run` request json.
* `dry_run_update_result` - Used for store `dry_run` update request json.
* `status` - Acceleration service status.


## Import

CDN domain can be imported using the id, e.g.

```
$ terraform import tencentcloud_cdn_domain.foo xxxx.com
```

