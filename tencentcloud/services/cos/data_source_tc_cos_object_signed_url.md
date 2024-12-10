Use this data source to query the signed url of the COS object.

Example Usage

```hcl
data "tencentcloud_cos_object_signed_url" "cos_object_signed_url" {
  bucket = "xxxxxx"
  path   = "path/to/file"
}
```
