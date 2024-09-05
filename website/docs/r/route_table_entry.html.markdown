---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_route_table_entry"
sidebar_current: "docs-tencentcloud-resource-route_table_entry"
description: |-
  Provides a resource to create an entry of a routing table.
---

# tencentcloud_route_table_entry

Provides a resource to create an entry of a routing table.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet"
  cidr_block        = "10.0.12.0/24"
  availability_zone = var.availability_zone
  route_table_id    = tencentcloud_route_table.example.id
}

# create route table
resource "tencentcloud_route_table" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
  name   = "tf-example"
}

# create route table entry
resource "tencentcloud_route_table_entry" "example" {
  route_table_id         = tencentcloud_route_table.example.id
  destination_cidr_block = "10.4.4.0/24"
  next_type              = "EIP"
  next_hub               = "0"
  description            = "describe"
}

# output
output "item_id" {
  value = tencentcloud_route_table_entry.example.route_item_id
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required, String, ForceNew) Destination address block.
* `next_hub` - (Required, String, ForceNew) ID of next-hop gateway. Note: when `next_type` is EIP, `next_hub` should be `0`.
* `next_type` - (Required, String, ForceNew) Type of next-hop. Valid values: `CVM`, `VPN`, `DIRECTCONNECT`, `PEERCONNECTION`, `HAVIP`, `NAT`, `NORMAL_CVM`, `EIP` and `LOCAL_GATEWAY`.
* `route_table_id` - (Required, String, ForceNew) ID of routing table to which this entry belongs.
* `description` - (Optional, String, ForceNew) Description of the routing table entry.
* `disabled` - (Optional, Bool) Whether the entry is disabled, default is `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `route_item_id` - ID of route table entry.


## Import

Route table entry can be imported using the id, e.g.

```
$ terraform import tencentcloud_route_table_entry.example 3065857.rtb-b050fg94
```

