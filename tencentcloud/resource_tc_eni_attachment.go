/*
Provides a resource to detailed information of attached backend server to an ENI.

Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
  name       = "ci-test-eni-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  availability_zone = "ap-guangzhou-3"
  name              = "ci-test-eni-subnet"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = "${tencentcloud_vpc.foo.id}"
  subnet_id   = "${tencentcloud_subnet.foo.id}"
  description = "eni desc"
  ipv4_count  = 1
}

data "tencentcloud_image" "my_favorite_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 1
}

resource "tencentcloud_instance" "foo" {
  instance_name            = "ci-test-eni-attach"
  availability_zone        = "ap-guangzhou-3"
  image_id                 = "${data.tencentcloud_image.my_favorite_image.image_id}"
  instance_type            = "${data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type}"
  system_disk_type         = "CLOUD_PREMIUM"
  disable_security_service = true
  disable_monitor_service  = true
  vpc_id                   = "${tencentcloud_vpc.foo.id}"
  subnet_id                = "${tencentcloud_subnet.foo.id}"
}

resource "tencentcloud_eni_attachment" "foo" {
  eni_id      = "${tencentcloud_eni.foo.id}"
  instance_id = "${tencentcloud_instance.foo.id}"
}
```

Import

ENI attachment can be imported using the id, e.g.

```
  $ terraform import tencentcloud_eni_attachment.foo eni-gtlvkjvz+ins-0h3a5new
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func resourceTencentCloudEniAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEniAttachmentCreate,
		Read:   resourceTencentCloudEniAttachmentRead,
		Delete: resourceTencentCloudEniAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"eni_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the ENI.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the instance which bind the ENI.",
			},
		},
	}
}

func resourceTencentCloudEniAttachmentCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_attachment.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	eniId := d.Get("eni_id").(string)
	cvmId := d.Get("instance_id").(string)

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	if err := service.AttachEniToCvm(ctx, eniId, cvmId); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s+%s", eniId, cvmId))

	return resourceTencentCloudEniAttachmentRead(d, m)
}

func resourceTencentCloudEniAttachmentRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_attachment.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	split := strings.Split(id, "+")
	if len(split) != 2 {
		log.Printf("[CRITAL]%s id %s is invalid", logId, id)
		d.SetId("")
		return nil
	}

	eniId := split[0]

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	enis, err := service.DescribeEniById(ctx, []string{eniId})
	if err != nil {
		return err
	}

	var eni *vpc.NetworkInterface
	for _, e := range enis {
		if e.NetworkInterfaceId == nil {
			return errors.New("eni id is nil")
		}

		if *e.NetworkInterfaceId == eniId {
			eni = e
			break
		}
	}

	if eni == nil {
		d.SetId("")
		return nil
	}

	if eni.NetworkInterfaceId == nil {
		return errors.New("eni id is nil")
	}

	if eni.Attachment == nil {
		d.SetId("")
		return nil
	}

	if eni.Attachment.InstanceId == nil {
		return errors.New("eni attach instance id is nil")
	}

	d.Set("eni_id", eni.NetworkInterfaceId)
	d.Set("instance_id", eni.Attachment.InstanceId)

	return nil
}

func resourceTencentCloudEniAttachmentDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_eni_attachment.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	split := strings.Split(id, "+")
	if len(split) != 2 {
		log.Printf("[CRITAL]%s id %s is invalid", logId, id)
		d.SetId("")
		return nil
	}

	eniId, cvmId := split[0], split[1]

	service := VpcService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DetachEniFromCvm(ctx, eniId, cvmId)
}
