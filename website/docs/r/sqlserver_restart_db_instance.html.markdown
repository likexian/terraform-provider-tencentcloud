---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_restart_db_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_restart_db_instance"
description: |-
  Provides a resource to create a sqlserver restart_db_instance
---

# tencentcloud_sqlserver_restart_db_instance

Provides a resource to create a sqlserver restart_db_instance

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "tf-example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  project_id        = 0
  memory            = 16
  storage           = 40
}

resource "tencentcloud_sqlserver_restart_db_instance" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver restart_db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_restart_db_instance.restart_db_instance restart_db_instance_id
```

