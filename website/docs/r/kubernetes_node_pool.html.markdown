---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_node_pool"
sidebar_current: "docs-tencentcloud-resource-kubernetes_node_pool"
description: |-
  Provide a resource to create an auto scaling group for kubernetes cluster.
---

# tencentcloud_kubernetes_node_pool

Provide a resource to create an auto scaling group for kubernetes cluster.

~> **NOTE:**  We recommend the usage of one cluster with essential worker config + node pool to manage cluster and nodes. Its a more flexible way than manage worker config with tencentcloud_kubernetes_cluster, tencentcloud_kubernetes_scale_worker or exist node management of `tencentcloud_kubernetes_attachment`. Cause some unchangeable parameters of `worker_config` may cause the whole cluster resource `force new`.

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.31.0.0/16"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

//this is the cluster with empty worker config
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf-tke-unit-test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32
  cluster_version         = "1.18.4"
  cluster_deploy_type     = "MANAGED_CLUSTER"
}

//this is one example of managing node using node pool
resource "tencentcloud_kubernetes_node_pool" "mynodepool" {
  name                     = "mynodepool"
  cluster_id               = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size                 = 6
  min_size                 = 1
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids               = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy             = "INCREMENTAL_INTERVALS"
  desired_capacity         = 4
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    security_group_ids = ["sg-24vswocp"]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false

  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  taints {
    key    = "test_taint"
    value  = "taint_value"
    effect = "PreferNoSchedule"
  }

  taints {
    key    = "test_taint2"
    value  = "taint_value2"
    effect = "PreferNoSchedule"
  }

  node_config {
    extra_args = [
      "root-dir=/var/lib/kubelet"
    ]
  }
}
```

Using Spot CVM Instance

```hcl
resource "tencentcloud_kubernetes_node_pool" "mynodepool" {
  name                     = "mynodepool"
  cluster_id               = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size                 = 6
  min_size                 = 1
  vpc_id                   = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids               = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy             = "INCREMENTAL_INTERVALS"
  desired_capacity         = 4
  enable_auto_scale        = true
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type        = var.default_instance_type
    system_disk_type     = "CLOUD_PREMIUM"
    system_disk_size     = "50"
    security_group_ids   = ["sg-24vswocp"]
    instance_charge_type = "SPOTPAID"
    spot_instance_type   = "one-time"
    spot_max_price       = "1000"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_config` - (Required) Auto scaling config parameters.
* `cluster_id` - (Required, ForceNew) ID of the cluster.
* `max_size` - (Required) Maximum number of node.
* `min_size` - (Required) Minimum number of node.
* `name` - (Required, ForceNew) Name of the node pool. The name does not exceed 25 characters, and only supports Chinese, English, numbers, underscores, separators (`-`) and decimal points.
* `vpc_id` - (Required, ForceNew) ID of VPC network.
* `default_cooldown` - (Optional) Seconds of scaling group cool down. Default value is `300`.
* `delete_keep_instance` - (Optional) Indicate to keep the CVM instance when delete the node pool. Default is `true`.
* `desired_capacity` - (Optional) Desired capacity ot the node. If `enable_auto_scale` is set `true`, this will be a computed parameter.
* `enable_auto_scale` - (Optional) Indicate whether to enable auto scaling or not.
* `labels` - (Optional) Labels of kubernetes node pool created nodes. The label key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').
* `multi_zone_subnet_policy` - (Optional, ForceNew) Multi-availability zone/subnet policy. Valid values: PRIORITY and EQUALITY. Default value: PRIORITY.
* `node_config` - (Optional) Node config.
* `node_os_type` - (Optional) The image version of the node. Valida values are `DOCKER_CUSTOMIZE` and `GENERAL`. Default is `GENERAL`. This parameter will only affect new nodes, not including the existing nodes.
* `node_os` - (Optional) Operating system of the cluster, the available values include: `tlinux2.4x86_64`, `ubuntu18.04.1x86_64`, `ubuntu16.04.1 LTSx86_64`, `centos7.6.0_x64` and `centos7.2x86_64`. Default is 'tlinux2.4x86_64'. This parameter will only affect new nodes, not including the existing nodes.
* `retry_policy` - (Optional, ForceNew) Available values for retry policies include `IMMEDIATE_RETRY` and `INCREMENTAL_INTERVALS`.
* `scaling_group_name` - (Optional) Name of relative scaling group.
* `scaling_group_project_id` - (Optional) Project ID the scaling group belongs to.
* `scaling_mode` - (Optional, ForceNew) Auto scaling mode. Valid values are `CLASSIC_SCALING`(scaling by create/destroy instances), `WAKE_UP_STOPPED_SCALING`(Boot priority for expansion. When expanding the capacity, the shutdown operation is given priority to the shutdown of the instance. If the number of instances is still lower than the expected number of instances after the startup, the instance will be created, and the method of destroying the instance will still be used for shrinking).
* `subnet_ids` - (Optional, ForceNew) ID list of subnet, and for VPC it is required.
* `taints` - (Optional) Taints of kubernetes node pool created nodes.
* `termination_policies` - (Optional) Policy of scaling group termination. Available values: `["OLDEST_INSTANCE"]`, `["NEWEST_INSTANCE"]`.
* `unschedulable` - (Optional, ForceNew) Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.
* `zones` - (Optional) List of auto scaling group available zones, for Basic network it is required.

The `auto_scaling_config` object supports the following:

* `instance_type` - (Required, ForceNew) Specified types of CVM instance.
* `backup_instance_types` - (Optional) Backup CVM instance types if specified instance type sold out or mismatch.
* `bandwidth_package_id` - (Optional) bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.
* `data_disk` - (Optional) Configurations of data disk.
* `enhanced_monitor_service` - (Optional, ForceNew) To specify whether to enable cloud monitor service. Default is TRUE.
* `enhanced_security_service` - (Optional, ForceNew) To specify whether to enable cloud security service. Default is TRUE.
* `instance_charge_type_prepaid_period` - (Optional) The tenancy (in month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.
* `instance_charge_type_prepaid_renew_flag` - (Optional) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional) Charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`. The default is `POSTPAID_BY_HOUR`. NOTE: `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.
* `internet_charge_type` - (Optional) Charge types for network traffic. Valid value: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
* `internet_max_bandwidth_out` - (Optional) Max bandwidth of Internet access in Mbps. Default is `0`.
* `key_ids` - (Optional, ForceNew) ID list of keys.
* `password` - (Optional, ForceNew) Password to access.
* `public_ip_assigned` - (Optional) Specify whether to assign an Internet IP address.
* `security_group_ids` - (Optional) Security groups to which a CVM instance belongs.
* `spot_instance_type` - (Optional) Type of spot instance, only support `one-time` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `spot_max_price` - (Optional) Max price of a spot instance, is the format of decimal string, for example "0.50". Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `system_disk_size` - (Optional) Volume of system disk in GB. Default is `50`.
* `system_disk_type` - (Optional) Type of a CVM disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`. Default is `CLOUD_PREMIUM`.

The `data_disk` object supports the following:

* `auto_format_and_mount` - (Optional, ForceNew) Indicate whether to auto format and mount or not. Default is `false`.
* `disk_partition` - (Optional, ForceNew) The name of the device or partition to mount.
* `disk_size` - (Optional, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, ForceNew) Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD`.
* `file_system` - (Optional, ForceNew) File system, e.g. `ext3/ext4/xfs`.
* `mount_target` - (Optional, ForceNew) Mount target.

The `data_disk` object supports the following:

* `disk_size` - (Optional, ForceNew) Volume of disk in GB. Default is `0`.
* `disk_type` - (Optional, ForceNew) Types of disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`.
* `snapshot_id` - (Optional, ForceNew) Data disk snapshot ID.

The `node_config` object supports the following:

* `data_disk` - (Optional, ForceNew) Configurations of data disk.
* `desired_pod_num` - (Optional, ForceNew) Indicate to set desired pod number in node. valid when the cluster is podCIDR.
* `docker_graph_path` - (Optional, ForceNew) Docker graph path. Default is `/var/lib/docker`.
* `extra_args` - (Optional, ForceNew) Custom parameter information related to the node. This is a white-list parameter.
* `is_schedule` - (Optional, ForceNew) Indicate to schedule the adding node or not. Default is true.
* `mount_target` - (Optional, ForceNew) Mount target. Default is not mounting.
* `user_data` - (Optional, ForceNew) Base64-encoded User Data text, the length limit is 16KB.

The `taints` object supports the following:

* `effect` - (Required) Effect of the taint. Valid values are: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.
* `key` - (Required) Key of the taint. The taint key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').
* `value` - (Required) Value of the taint.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `auto_scaling_group_id` - The auto scaling group ID.
* `launch_config_id` - The launch config ID.
* `node_count` - The total node count.
* `status` - Status of the node pool.


