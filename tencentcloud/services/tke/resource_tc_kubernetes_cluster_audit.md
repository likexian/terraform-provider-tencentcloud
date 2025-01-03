Provides a resource to create a kubernetes cluster audit

Example Usage

Automatic creation of log sets and topics

```hcl
resource "tencentcloud_kubernetes_cluster_audit" "example" {
  cluster_id              = "cls-fdy7hm1q"
  delete_logset_and_topic = true
}
```

Manually fill in log sets and topics

```hcl
resource "tencentcloud_kubernetes_cluster_audit" "example" {
  cluster_id              = "cls-fdy7hm1q"
  logset_id               = "30d32c56-e650-4175-9c70-5280cddee48c"
  topic_id                = "cfc056ca-517f-46fd-be68-9c5cad518b2f"
  topic_region            = "ap-guangzhou"
  delete_logset_and_topic = false
}
```

Import

kubernetes cluster audit can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_cluster_audit.example cls-fdy7hm1q
```
