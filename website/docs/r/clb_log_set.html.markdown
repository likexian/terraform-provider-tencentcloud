---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_log_set"
sidebar_current: "docs-tencentcloud-resource-clb_log_set"
description: |-
  Provides a resource to create an exclusive CLB Logset.
---

# tencentcloud_clb_log_set

Provides a resource to create an exclusive CLB Logset.

## Example Usage

```hcl
resource "tencentcloud_clb_log_set" "example" {
  logset_name = "tf-example"
  logset_type = "ACCESS"
  period      = 7
}
```

## Argument Reference

The following arguments are supported:

* `logset_name` - (Optional, String, ForceNew) Logset name, which must be unique among all CLS logsets; default value: clb_logset.
* `logset_type` - (Optional, String, ForceNew) Logset type. Valid values: ACCESS (access logs; default value) and HEALTH (health check logs).
* `period` - (Optional, Int, ForceNew) Logset retention period in days. Maximun value is `90`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Logset creation time.
* `name` - (**Deprecated**) It has been deprecated from version 1.81.162+. Please use `logset_name` instead. Logset name, which unique and fixed `clb_logset` among all CLS logsets.
* `topic_count` - Number of log topics in logset.


## Import

CLB log set can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_logset.example.example 4eb9e3a8-9c42-4b32-9ddf-e215e9c92764
```

