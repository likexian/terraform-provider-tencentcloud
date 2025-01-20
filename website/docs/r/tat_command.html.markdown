---
subcategory: "TencentCloud Automation Tools(TAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tat_command"
sidebar_current: "docs-tencentcloud-resource-tat_command"
description: |-
  Provides a resource to create a TAT command
---

# tencentcloud_tat_command

Provides a resource to create a TAT command

## Example Usage

```hcl
resource "tencentcloud_tat_command" "example" {
  username          = "root"
  command_name      = "tf-example"
  content           = <<EOF
#!/bin/bash
if [ "$(id -u)" != "0" ]; then
    echo "Please run this script as the root user." >&2
    exit 1
fi
ps aux
EOF
  description       = "Terraform demo."
  command_type      = "SHELL"
  working_directory = "/root"
  timeout           = 50
  tags {
    key   = "createBy"
    value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `command_name` - (Required, String) Command name. The name can be up to 60 bytes, and contain [a-z], [A-Z], [0-9] and [_-.].
* `content` - (Required, String) Base64-encoded command. The maximum length is 64 KB.
* `command_type` - (Optional, String) Command type. `SHELL` and `POWERSHELL` are supported. The default value is `SHELL`.
* `default_parameters` - (Optional, String) The default value of the custom parameter value when it is enabled. The field type is JSON encoded string. For example, {"varA": "222"}.`key` is the name of the custom parameter and value is the default value. Both `key` and `value` are strings.If no parameter value is provided in the `InvokeCommand` API, the default value is used.Up to 20 custom parameters are supported.The name of the custom parameter cannot exceed 64 characters and can contain [a-z], [A-Z], [0-9] and [-_].
* `description` - (Optional, String) Command description. The maximum length is 120 characters.
* `enable_parameter` - (Optional, Bool) Whether to enable the custom parameter feature.This cannot be modified once created.Default value: `false`.
* `output_cos_bucket_url` - (Optional, String) The COS bucket URL for uploading logs. The URL must start with `https`, such as `https://BucketName-123454321.cos.ap-beijing.myqcloud.com`.
* `output_cos_key_prefix` - (Optional, String) The COS bucket directory where the logs are saved. Check below for the rules of the directory name.1. It must be a combination of number, letters, and visible characters. Up to 60 characters are allowed.2. Use a slash (/) to create a subdirectory.3. Consecutive dots (.) and slashes (/) are not allowed. It can not start with a slash (/).
* `tags` - (Optional, List) Tags bound to the command. At most 10 tags are allowed.
* `timeout` - (Optional, Int) Command timeout period. Default value: 60 seconds. Value range: [1, 86400].
* `username` - (Optional, String) The username used to execute the command on the CVM or Lighthouse instance.The principle of least privilege is the best practice for permission management. We recommend you execute TAT commands as a general user. By default, the root user is used to execute commands on Linux and the System user is used on Windows.
* `working_directory` - (Optional, String) Command execution path. The default value is /root for `SHELL` commands and C:/Program Files/qcloudtat_agent/workdir for `POWERSHELL` commands.

The `tags` object supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_by` - Command creator. `TAT` indicates a public command and `USER` indicates a personal command.
* `created_time` - Command creation time.
* `formatted_description` - Formatted description of the command. This parameter is an empty string for user commands and contains values for public commands.
* `updated_time` - Command update time.


## Import

tat command can be imported using the id, e.g.
```
$ terraform import tencentcloud_tat_command.example cmd-6fydo27j
```

