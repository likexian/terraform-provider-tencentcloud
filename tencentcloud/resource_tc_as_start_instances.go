/*
Provides a resource to create a as start_instances

Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "zone"
    values = [data.tencentcloud_availability_zones_by_product.zones.zones.0.name]
  }

  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  image_id          = data.tencentcloud_images.image.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "user"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
}

# Attachment Instance
resource "tencentcloud_as_attachment" "attachment" {
  scaling_group_id = tencentcloud_as_scaling_group.example.id
  instance_ids     = [tencentcloud_instance.example.id]
}

resource "tencentcloud_as_start_instances" "start_instances" {
  auto_scaling_group_id = tencentcloud_as_scaling_group.example.id
  instance_ids          = tencentcloud_as_attachment.attachment.instance_ids
}

```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAsStartInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsStartInstancesCreate,
		Read:   resourceTencentCloudAsStartInstancesRead,
		Delete: resourceTencentCloudAsStartInstancesDelete,
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Launch configuration ID.",
			},

			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of cvm instances to start.",
			},
		},
	}
}

func resourceTencentCloudAsStartInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_as_start_instances.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = as.NewStartAutoScalingInstancesRequest()
		response   = as.NewStartAutoScalingInstancesResponse()
		activityId string
	)
	if v, ok := d.GetOk("auto_scaling_group_id"); ok {
		request.AutoScalingGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().StartAutoScalingInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate as startInstances failed, reason:%+v", logId, err)
		return nil
	}

	activityId = *response.Response.ActivityId

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err = resource.Retry(4*readRetryTimeout, func() *resource.RetryError {
		status, err := service.DescribeActivityById(ctx, activityId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if status == SCALING_GROUP_ACTIVITY_STATUS_INIT || status == SCALING_GROUP_ACTIVITY_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("remove status is running(%s)", status))
		}
		if status == SCALING_GROUP_ACTIVITY_STATUS_SUCCESSFUL {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("remove status is failed(%s)", status))
	})
	if err != nil {
		return err
	}

	d.SetId(activityId)

	return resourceTencentCloudAsStartInstancesRead(d, meta)
}

func resourceTencentCloudAsStartInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_start_instances.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsStartInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_start_instances.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
