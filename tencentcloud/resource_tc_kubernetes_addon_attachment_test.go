package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTkeAddonAttachmentResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeAddonAttachment(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("resource_tc_kubernetes_addon_attachment.cbs", "response_body"),
					resource.TestCheckResourceAttr("resource_tc_kubernetes_addon_attachment.cbs", "name", "cbs"),
					resource.TestCheckResourceAttr("resource_tc_kubernetes_addon_attachment.cbs", "version", "1.0.0"),
				),
			},
		},
	})
}

func testAccTkeAddonAttachment() string {
	return `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}


data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32
  cluster_os			  = "tlinux2.2(tkernel3)x86_64"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_addon_attachment" "cbs" {
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  name = "cbs"
  version = "1.0.0"
}
`
}
