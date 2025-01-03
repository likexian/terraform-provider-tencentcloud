---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cluster_audit"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cluster_audit"
description: |-
  Provides a resource to create a kubernetes cluster audit
---

# tencentcloud_kubernetes_cluster_audit

Provides a resource to create a kubernetes cluster audit

## Example Usage

### Automatic creation of log sets and topics

```hcl
resource "tencentcloud_kubernetes_cluster_audit" "example" {
  cluster_id              = "cls-fdy7hm1q"
  delete_logset_and_topic = true
}
```

### Manually fill in log sets and topics

```hcl
resource "tencentcloud_kubernetes_cluster_audit" "example" {
  cluster_id              = "cls-fdy7hm1q"
  logset_id               = "30d32c56-e650-4175-9c70-5280cddee48c"
  topic_id                = "cfc056ca-517f-46fd-be68-9c5cad518b2f"
  topic_region            = "ap-guangzhou"
  delete_logset_and_topic = false
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `delete_logset_and_topic` - (Optional, Bool, ForceNew) `true` means to delete the log set and topic created by default when cluster audit is turned off; `false` means not to delete. Default is `false`.
* `logset_id` - (Optional, String, ForceNew) CLS logset ID.
* `topic_id` - (Optional, String, ForceNew) CLS topic ID.
* `topic_region` - (Optional, String, ForceNew) The region where the topic is located defaults to the current region of the cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

kubernetes cluster audit can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_cluster_audit.example cls-fdy7hm1q
```

