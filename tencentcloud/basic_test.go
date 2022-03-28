package tencentcloud

import "os"

/*
---------------------------------------------------
If you want to run through the test cases,
the following must be changed to your resource id.
---------------------------------------------------
*/

var appid string = os.Getenv("TENCENTCLOUD_APPID")
var ownerUin string = os.Getenv("TENCENTCLOUD_OWNER_UIN")

const (
	keepResource    = "keep"
	defaultResource = "Default"
)

const (
	defaultRegion      = "ap-guangzhou"
	defaultVpcId       = "vpc-86v957zb"
	defaultVpcCidr     = "172.16.0.0/16"
	defaultVpcCidrLess = "172.16.0.0/18"

	defaultAZone          = "ap-guangzhou-3"
	defaultSubnetId       = "subnet-enm92y0m"
	defaultSubnetCidr     = "172.16.0.0/20"
	defaultSubnetCidrLess = "172.16.0.0/22"

	defaultInsName       = "tf-ci-test"
	defaultInsNameUpdate = "tf-ci-test-update"

	defaultSshCertificate  = "f8kGFR2T"
	defaultSshCertificateB = "fbW9Spiy"

	defaultDayuBgp    = "bgp-000006mq"
	defaultDayuBgpMul = "bgp-0000008o"
	defaultDayuBgpIp  = "bgpip-00000294"
	defaultDayuNet    = "net-0000007e"

	defaultGaapProxyId = "link-4yb9g6tb"

	defaultSecurityGroup  = "sg-ijato2x1"
	defaultSecurityGroup2 = "sg-51rgzop1"

	defaultProjectId   = "1250480"
	defaultDayuBgpIdV2 = "bgpip-000004x0"
	defaultDayuBgpIpV2 = "119.28.217.253"

	defaultTkeOSImageId   = "img-2lr9q49h"
	defaultTkeOSImageName = "tlinux2.2(tkernel3)x86_64"
)

// Tke Exclusive Network Environment
const (
	tkeExclusiveVpcId     = "vpc-391sv4w3"
	tkeExclusiveVpcName   = "keep_tke_exclusive_vpc"
	tkeExclusiveSubnetId  = "subnet-ljyn7h30"
	defaultTkeClusterId   = "cls-ely08ic4"
	defaultTkeClusterName = "preset_tke_cluster"
)

/*
---------------------------------------------------
The following are common test case used as templates.
---------------------------------------------------
*/

const defaultVpcVariable = `
variable "instance_name" {
  default = "` + defaultInsName + `"
}

variable "instance_name_update" {
  default = "` + defaultInsNameUpdate + `"
}

variable "availability_region" {
  default = "` + defaultRegion + `"
}

variable "availability_zone" {
  default = "` + defaultAZone + `"
}

variable "vpc_id" {
  default = "` + defaultVpcId + `"
}

variable "vpc_cidr" {
  default = "` + defaultVpcCidr + `"
}

variable "vpc_cidr_less" {
  default = "` + defaultVpcCidrLess + `"
}

variable "subnet_id" {
  default = "` + defaultSubnetId + `"
}

variable "subnet_cidr" {
  default = "` + defaultSubnetCidr + `"
}

variable "subnet_cidr_less" {
  default = "` + defaultSubnetCidrLess + `"
}
`

const defaultInstanceVariable = defaultVpcVariable + `
data "tencentcloud_availability_zones" "default" {
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "instance-family"
    values = ["S1"]
  }

  cpu_core_count = 1
  memory_size    = 1
}
`

const (
	defaultMySQLName = "preset_mysql"
)

// ref with `local.mysql_id`
const CommonPresetMysql = `

variable "availability_zone" {
  default = "` + defaultAZone + `"
}
variable "region" {
  default = "` + defaultRegion + `"
}

data "tencentcloud_mysql_instance" "mysql" {
  instance_name = "` + defaultMySQLName + `"
}

locals {
  mysql_id = data.tencentcloud_mysql_instance.mysql.instance_list.0.mysql_id
}
`

const instanceCommonTestCase = defaultInstanceVariable + `
resource "tencentcloud_instance" "default" {
  instance_name              = var.instance_name
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.default.images.1.image_id
  instance_type              = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10
  vpc_id                     = var.vpc_id
  subnet_id                  = var.subnet_id
}
`

const mysqlInstanceCommonTestCase = defaultVpcVariable + `
resource "tencentcloud_mysql_instance" "default" {
  mem_size = 1000
  volume_size = 25
  instance_name = var.instance_name
  engine_version = "5.7"
  root_password = "0153Y474"
  availability_zone = var.availability_zone
  force_delete = true
}
`
const mysqlInstanceHighPerformanceTestCase = defaultVpcVariable + `
resource "tencentcloud_mysql_instance" "default" {
  mem_size = 1000
  volume_size = 50
  instance_name = var.instance_name
  engine_version = "5.7"
  root_password = "0153Y474"
  availability_zone = var.availability_zone
  force_delete = true
}
`

const mysqlInstanceHighPerformancePrepaidTestCase = defaultVpcVariable + `
resource "tencentcloud_mysql_instance" "default" {
  mem_size = 1000
  volume_size = 50
  pay_type = 0
  instance_name = var.instance_name
  engine_version = "5.7"
  root_password = "0153Y474"
  availability_zone = var.availability_zone
  force_delete = true
}
`
const defaultSCFCosBucket = `
data "tencentcloud_user_info" "info" {}

data "tencentcloud_cos_buckets" "buckets" {
  bucket_prefix = "preset-scf-bucket-${data.tencentcloud_user_info.info.app_id}"
}

locals {
  bucket_name = data.tencentcloud_cos_buckets.buckets.bucket_list.0.bucket
  bucket_url = data.tencentcloud_cos_buckets.buckets.bucket_list.0.cos_bucket_url
}
`

const defaultScfNamespace = "preset-scf-namespace"

const defaultFileSystemName = "preset_cfs"

const defaultFileSystem = `
data "tencentcloud_cfs_file_systems" "fs" {
  name = "` + defaultFileSystemName + `"
}

# doesn't support datasource for now
variable "mount_id" {
  default = "cfs-iobiaxtj"
}

locals {
  cfs = data.tencentcloud_cfs_file_systems.fs.file_system_list.0
  cfs_id = local.cfs.file_system_id
  access_group_id = local.cfs.access_group_id
}`
