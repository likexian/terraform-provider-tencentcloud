---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_bandwidth_package"
sidebar_current: "docs-tencentcloud-resource-vpc_bandwidth_package"
description: |-
  Provides a resource to create a vpc bandwidth_package
---

# tencentcloud_vpc_bandwidth_package

Provides a resource to create a vpc bandwidth_package

## Example Usage

```hcl
resource "tencentcloud_vpc_bandwidth_package" "example" {
  network_type           = "BGP"
  charge_type            = "TOP5_POSTPAID_BY_MONTH"
  bandwidth_package_name = "tf-example"
  tags = {
    "createdBy" = "terraform"
  }
}
```

### PrePaid Bandwidth Package

```hcl
resource "tencentcloud_vpc_bandwidth_package" "bandwidth_package" {
  network_type           = "BGP"
  charge_type            = "FIXED_PREPAID_BY_MONTH"
  bandwidth_package_name = "test-001"
  time_span              = 3
  internet_max_bandwidth = 100
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bandwidth_package_name` - (Optional, String) Bandwidth package name.
* `charge_type` - (Optional, String) Bandwidth package billing type, default: TOP5_POSTPAID_BY_MONTH. Optional value: `TOP5_POSTPAID_BY_MONTH`: TOP5 billed by monthly postpaid; `PERCENT95_POSTPAID_BY_MONTH`: 95 billed monthly postpaid; `FIXED_PREPAID_BY_MONTH`: Monthly prepaid billing (Type FIXED_PREPAID_BY_MONTH product API capability is under construction); `BANDWIDTH_POSTPAID_BY_DAY`: bandwidth billed by daily postpaid; `ENHANCED95_POSTPAID_BY_MONTH`: enhanced 95 billed monthly postpaid.
* `internet_max_bandwidth` - (Optional, Int) Bandwidth packet speed limit size. Unit: Mbps, -1 means no speed limit.
* `network_type` - (Optional, String) Bandwidth packet type, default:BGP, optional:- `BGP`: common BGP shared bandwidth package- `HIGH_QUALITY_BGP`: High Quality BGP Shared Bandwidth Package.
* `tags` - (Optional, Map) Tag description list.
* `time_span` - (Optional, Int) The purchase duration of the prepaid monthly bandwidth package, unit: month, value range: 1~60.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc bandwidth_package can be imported using the id, e.g.
```
$ terraform import tencentcloud_vpc_bandwidth_package.bandwidth_package bandwidthPackage_id
```

