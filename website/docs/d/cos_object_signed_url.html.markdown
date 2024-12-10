---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_object_signed_url"
sidebar_current: "docs-tencentcloud-datasource-cos_object_signed_url"
description: |-
  Use this data source to query the signed url of the COS object.
---

# tencentcloud_cos_object_signed_url

Use this data source to query the signed url of the COS object.

## Example Usage

```hcl
data "tencentcloud_cos_object_signed_url" "cos_object_signed_url" {
  bucket = "xxxxxx"
  path   = "path/to/file"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) Name of the bucket.
* `path` - (Required, String) The full path to the object inside the bucket.
* `method` - (Optional, String) Method, GET or PUT. Default value is GET.
* `duration` - (Optional, String) Duration of signed url. Default value is 1m.
* `headers` - (Optional, Map) Request headers.
* `queries` - (Optional, Map) Request query parameters.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `signed_url` - Signed URL.


