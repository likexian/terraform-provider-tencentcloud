---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_storage_attachment"
sidebar_current: "docs-tencentcloud-resource-cbs_storage_attachment"
description: |-
  Provides a CBS storage attachment resource.
---

# tencentcloud_cbs_storage_attachment

Provides a CBS storage attachment resource.

## Example Usage

```hcl
resource "tencentcloud_cbs_storage_attachment" "attachment" {
  storage_id  = "disk-kdt0sq6m"
  instance_id = "ins-jqlegd42"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the CVM instance.
* `storage_id` - (Required, ForceNew) ID of the mounted CBS.


