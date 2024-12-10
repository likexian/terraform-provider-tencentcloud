---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instance_network_access"
sidebar_current: "docs-tencentcloud-resource-postgresql_instance_network_access"
description: |-
  Provides a resource to create a postgres postgresql_instance_network_access
---

# tencentcloud_postgresql_instance_network_access

Provides a resource to create a postgres postgresql_instance_network_access

## Example Usage

```hcl
resource "tencentcloud_postgresql_instance_network_access" "postgresql_instance_network_access" {
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID in the format of postgres-6bwgamo3.
* `is_assign_vip` - (Required, Bool, ForceNew) Whether to manually assign the VIP. Valid values: `true` (manually assign), `false` (automatically assign).
* `subnet_id` - (Required, String, ForceNew) Subnet ID.
* `vpc_id` - (Required, String, ForceNew) Unified VPC ID.
* `vip` - (Optional, String, ForceNew) Target VIP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgres postgresql_instance_network_access can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_instance_network_access.postgresql_instance_network_access postgresql_instance_network_access_id
```

