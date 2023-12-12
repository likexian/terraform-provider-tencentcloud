Provides a resource to create a ses template.

Example Usage

Create a ses text template

```hcl
resource "tencentcloud_ses_template" "example" {
  template_name = "tf_example_ses_temp""
  template_content {
    text = "example for the ses template"
  }
}

```

Create a ses html template

```hcl
resource "tencentcloud_ses_template" "example" {
  template_name = "tf_example_ses_temp"
  template_content {
    html = <<-EOT
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>mail title</title>
</head>
<body>
<div class="container">
  <h1>Welcome to our service! </h1>
  <p>Dear user,</p>
  <p>Thank you for using Tencent Cloud:</p>
  <p><a href="https://cloud.tencent.com/document/product/1653">https://cloud.tencent.com/document/product/1653</a></p>
  <p>If you did not request this email, please ignore it. </p>
  <p><strong>from the iac team</strong></p>
</div>
</body>
</html>
    EOT
  }
}

```

Import

ses template can be imported using the id, e.g.
```
$ terraform import tencentcloud_ses_template.example template_id
```