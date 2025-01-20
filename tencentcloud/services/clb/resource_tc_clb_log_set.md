Provides a resource to create an exclusive CLB Logset.

Example Usage

```hcl
resource "tencentcloud_clb_log_set" "example" {
  logset_name = "tf-example"
  logset_type = "ACCESS"
  period      = 7
}
```

Import

CLB log set can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_logset.example.example 4eb9e3a8-9c42-4b32-9ddf-e215e9c92764
```