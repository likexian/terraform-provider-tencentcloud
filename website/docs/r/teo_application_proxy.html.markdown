---
subcategory: "Teo"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_application_proxy"
sidebar_current: "docs-tencentcloud-resource-teo_application_proxy"
description: |-
  Provides a resource to create a teo applicationProxy
---

# tencentcloud_teo_application_proxy

Provides a resource to create a teo applicationProxy

## Example Usage

```hcl
resource "tencentcloud_teo_application_proxy" "app0" {
  zone_id   = tencentcloud_teo_zone.sfurnace_work.id
  zone_name = "sfurnace.work"

  accelerate_type      = 1
  security_type        = 1
  plat_type            = "domain"
  proxy_name           = "www.sfurnace.work"
  proxy_type           = "hostname"
  session_persist_time = 2400
}
```

## Argument Reference

The following arguments are supported:

* `accelerate_type` - (Required, Int) - 0: Disable acceleration.- 1: Enable acceleration.
* `plat_type` - (Required, String) Scheduling mode.- ip: Anycast IP.- domain: CNAME.
* `proxy_name` - (Required, String) Layer-4 proxy name.
* `security_type` - (Required, Int) - 0: Disable security protection.- 1: Enable security protection.
* `zone_id` - (Required, String) Site ID.
* `zone_name` - (Required, String) Site name.
* `proxy_type` - (Optional, String) Specifies how a layer-4 proxy is created.- hostname: Subdomain name.- instance: Instance.
* `session_persist_time` - (Optional, Int) Session persistence duration. Value range: 30-3600 (in seconds).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `host_id` - ID of the layer-7 domain name.
* `proxy_id` - Proxy ID.
* `schedule_value` - Scheduling information.
* `update_time` - Last modification date.


## Import

teo applicationProxy can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_application_proxy.applicationProxy applicationProxy_id
```

