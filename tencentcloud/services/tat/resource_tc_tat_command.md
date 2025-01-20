Provides a resource to create a TAT command

Example Usage

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

Import

tat command can be imported using the id, e.g.
```
$ terraform import tencentcloud_tat_command.example cmd-6fydo27j
```