---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_launch_template"
sidebar_current: "docs-tencentcloud-resource-cvm_launch_template"
description: |-
  Provides a resource to create a cvm launch template
---

# tencentcloud_cvm_launch_template

Provides a resource to create a cvm launch template

## Example Usage

```hcl
data "tencentcloud_images" "my_favorite_image" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_cvm_launch_template" "demo" {
  launch_template_name = "test"
  placement {
    zone       = "ap-guangzhou-6"
    project_id = 0
  }
  image_id = data.tencentcloud_images.my_favorite_image.images.0.image_id
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required, String, ForceNew) Image ID.
* `launch_template_name` - (Required, String, ForceNew) The name of launch template.
* `placement` - (Required, List, ForceNew) The location of instance.
* `action_timer` - (Optional, List, ForceNew) Timed task.
* `cam_role_name` - (Optional, String, ForceNew) The role name of CAM.
* `client_token` - (Optional, String, ForceNew) A string to used guarantee request idempotency.
* `data_disks` - (Optional, List, ForceNew) Data disk configuration information of the instance.
* `disable_api_termination` - (Optional, Bool, ForceNew) Instance destruction protection flag.
* `disaster_recover_group_ids` - (Optional, Set: [`String`], ForceNew) The ID of disaster recover group.
* `dry_run` - (Optional, Bool, ForceNew) Whether to preflight only this request, true or false.
* `enhanced_service` - (Optional, List, ForceNew) Enhanced service. If this parameter is not specified, cloud monitoring and cloud security services will be enabled by default in public images.
* `host_name` - (Optional, String, ForceNew) The host name of CVM.
* `hpc_cluster_id` - (Optional, String, ForceNew) The ID of HPC cluster.
* `instance_charge_prepaid` - (Optional, List, ForceNew) The configuration of charge prepaid.
* `instance_charge_type` - (Optional, String, ForceNew) The charge type of instance. Default value: POSTPAID_BY_HOUR.
* `instance_count` - (Optional, Int, ForceNew) The number of instances purchased.
* `instance_market_options` - (Optional, List, ForceNew) The marketplace options of instance.
* `instance_name` - (Optional, String, ForceNew) The name of instance. If you do not specify an instance display name, 'Unnamed' is displayed by default.
* `instance_type` - (Optional, String, ForceNew) The type of the instance. If this parameter is not specified, the system will dynamically specify the default model according to the resource sales in the current region.
* `internet_accessible` - (Optional, List, ForceNew) The information settings of public network bandwidth. If you do not specify this parameter, the default Internet bandwidth is 0 Mbps.
* `launch_template_version_description` - (Optional, String, ForceNew) Instance launch template version description.
* `login_settings` - (Optional, List, ForceNew) The login settings of instance. By default, passwords are randomly generated and notified to users via internal messages.
* `security_group_ids` - (Optional, Set: [`String`], ForceNew) The security group ID of instance. If this parameter is not specified, the default security group is bound.
* `system_disk` - (Optional, List, ForceNew) System disk configuration information of the instance. If this parameter is not specified, it is assigned according to the system default.
* `tag_specification` - (Optional, List, ForceNew) Tag description list.
* `user_data` - (Optional, String, ForceNew) The data of users.
* `virtual_private_cloud` - (Optional, List, ForceNew) The configuration information of VPC. If this parameter is not specified, the basic network is used by default.

The `action_timer` object supports the following:

* `action_time` - (Optional, String) Execution time.
* `externals` - (Optional, List) Extended data.
* `timer_action` - (Optional, String) Timer name.

The `automation_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool) Whether to enable TencentCloud Automation Tools(TAT), TRUE or FALSE.

The `data_disks` object supports the following:

* `disk_size` - (Required, Int) The size of the data disk.
* `cdc_id` - (Optional, String) Cloud Dedicated Cluster(CDC) ID.
* `delete_with_instance` - (Optional, Bool) Whether the data disk is destroyed along with the instance, true or false.
* `disk_id` - (Optional, String) Data disk ID.
* `disk_type` - (Optional, String) The type of data disk.
* `encrypt` - (Optional, Bool) Whether the data disk is encrypted, TRUE or FALSE.
* `kms_key_id` - (Optional, String) The id of custom CMK.
* `snapshot_id` - (Optional, String) Data disk snapshot ID.
* `throughput_performance` - (Optional, Int) Cloud disk performance, MB/s.

The `enhanced_service` object supports the following:

* `automation_service` - (Optional, List) Enable TencentCloud Automation Tools(TAT).
* `monitor_service` - (Optional, List) Enable cloud monitor service.
* `security_service` - (Optional, List) Enable cloud security service.

The `externals` object of `action_timer` supports the following:

* `release_address` - (Optional, Bool) Release address.
* `storage_block_attr` - (Optional, List) HDD local storage attributes.
* `unsupport_networks` - (Optional, Set) Unsupported network type.

The `instance_charge_prepaid` object supports the following:

* `period` - (Required, Int) The period of purchasing instances.
* `renew_flag` - (Optional, String) Automatic renew flag.

The `instance_market_options` object supports the following:

* `spot_options` - (Required, List) Bidding related options.
* `market_type` - (Optional, String) Market option type, currently only supports value: spot.

The `internet_accessible` object supports the following:

* `bandwidth_package_id` - (Optional, String) The ID of bandwidth package.
* `internet_charge_type` - (Optional, String) The type of internet charge.
* `internet_max_bandwidth_out` - (Optional, Int) Internet outbound bandwidth upper limit, Mbps.
* `public_ip_assigned` - (Optional, Bool) Whether to allocate public network IP, TRUE or FALSE.

The `login_settings` object supports the following:

* `keep_image_login` - (Optional, String) Keep the original settings of the mirror.
* `key_ids` - (Optional, Set) List of key ID.
* `password` - (Optional, String) The login password of instance.

The `monitor_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool) Whether to enable cloud monitor service, TRUE or FALSE.

The `placement` object supports the following:

* `zone` - (Required, String) The available zone ID of the instance.
* `host_ids` - (Optional, Set) The CDH ID list of the instance(input).
* `host_ips` - (Optional, Set, **Deprecated**) It has been deprecated from version 1.81.108. Specify the host machine ip.
* `project_id` - (Optional, Int) The project ID of the instance.

The `security_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool) Whether to enable cloud security service, TRUE or FALSE.

The `spot_options` object of `instance_market_options` supports the following:

* `max_price` - (Required, String) Bidding.
* `spot_instance_type` - (Optional, String) Bidding request type, currently only supported type: one-time.

The `storage_block_attr` object of `externals` supports the following:

* `max_size` - (Required, Int) The maximum capacity of HDD local storage.
* `min_size` - (Required, Int) The minimum capacity of HDD local storage.
* `type` - (Required, String) The type of HDD local storage.

The `system_disk` object supports the following:

* `cdc_id` - (Optional, String) Cloud Dedicated Cluster(CDC) ID.
* `disk_id` - (Optional, String) System disk ID.
* `disk_size` - (Optional, Int) The size of system disk.
* `disk_type` - (Optional, String) The type of system disk.

The `tag_specification` object supports the following:

* `resource_type` - (Required, String) The type of resource.
* `tags` - (Required, List) Tag list.

The `tags` object of `tag_specification` supports the following:

* `key` - (Required, String) The key of tag.
* `value` - (Required, String) The value of tag.

The `virtual_private_cloud` object supports the following:

* `subnet_id` - (Required, String) The id of subnet.
* `vpc_id` - (Required, String) The id of VPC.
* `as_vpc_gateway` - (Optional, Bool) Is it used as a Public network gateway, TRUE or FALSE.
* `ipv6_address_count` - (Optional, Int) The number of ipv6 addresses for Elastic Network Interface.
* `private_ip_addresses` - (Optional, Set) The address of private ip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



