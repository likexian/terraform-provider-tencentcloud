Provides a resource to creating direct connect gateway instance.

Example Usage

If network_type is VPC

```hcl
// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create dc gateway
resource "tencentcloud_dc_gateway" "example" {
  name                = "tf-example"
  network_instance_id = tencentcloud_vpc.vpc.id
  network_type        = "VPC"
  gateway_type        = "NORMAL"
}
```

If network_type is CCN

```hcl
// create ccn
resource "tencentcloud_ccn" "ccn" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

// create dc gateway
resource "tencentcloud_dc_gateway" "example" {
  name                = "tf-example"
  network_instance_id = tencentcloud_ccn.ccn.id
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}
```

Import

Direct connect gateway instance can be imported, e.g.

```
$ terraform import tencentcloud_dc_gateway.example dcg-dr1y0hu7
```