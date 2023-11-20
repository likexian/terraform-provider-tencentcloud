/*
Provide a resource to create an auto scaling group for kubernetes cluster.

~> **NOTE:**  We recommend the usage of one cluster with essential worker config + node pool to manage cluster and nodes. Its a more flexible way than manage worker config with tencentcloud_kubernetes_cluster, tencentcloud_kubernetes_scale_worker or exist node management of `tencentcloud_kubernetes_attachment`. Cause some unchangeable parameters of `worker_config` may cause the whole cluster resource `force new`.

~> **NOTE:**  In order to ensure the integrity of customer data, if you destroy nodepool instance, it will keep the cvm instance associate with nodepool by default. If you want to destroy together, please set `delete_keep_instance` to `false`.

~> **NOTE:**  In order to ensure the integrity of customer data, if the cvm instance was destroyed due to shrinking, it will keep the cbs associate with cvm by default. If you want to destroy together, please set `delete_with_instance` to `true`.

Example Usage

```hcl

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.31.0.0/16"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

//this is the cluster with empty worker config
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf-tke-unit-test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32
  cluster_version         = "1.18.4"
  cluster_deploy_type = "MANAGED_CLUSTER"
}

//this is one example of managing node using node pool
resource "tencentcloud_kubernetes_node_pool" "mynodepool" {
  name = "mynodepool"
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size = 6
  min_size = 1
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 4
  enable_auto_scale    = true
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    orderly_security_group_ids = ["sg-24vswocp"]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false
	host_name                  = "12.123.0.0"
	host_name_style            = "ORIGINAL"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  taints {
	key = "test_taint"
    value = "taint_value"
    effect = "PreferNoSchedule"
  }

  taints {
	key = "test_taint2"
    value = "taint_value2"
    effect = "PreferNoSchedule"
  }

  node_config {
      extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
  }
}
```

Using Spot CVM Instance
```hcl
resource "tencentcloud_kubernetes_node_pool" "mynodepool" {
  name = "mynodepool"
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size = 6
  min_size = 1
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 4
  enable_auto_scale    = true
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    orderly_security_group_ids = ["sg-24vswocp", "sg-3qntci2v", "sg-7y1t2wax"]
	instance_charge_type = "SPOTPAID"
    spot_instance_type = "one-time"
    spot_max_price = "1000"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

}

```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var importFlag = false

// merge `instance_type` to `backup_instance_types` as param `instance_types`
func getNodePoolInstanceTypes(d *schema.ResourceData) []*string {
	configParas := d.Get("auto_scaling_config").([]interface{})
	dMap := configParas[0].(map[string]interface{})
	instanceType := dMap["instance_type"]
	currInsType := instanceType.(string)
	v, ok := dMap["backup_instance_types"]
	backupInstanceTypes := v.([]interface{})
	instanceTypes := make([]*string, 0)
	if !ok || len(backupInstanceTypes) == 0 {
		instanceTypes = append(instanceTypes, &currInsType)
		return instanceTypes
	}
	headType := backupInstanceTypes[0].(string)
	if headType != currInsType {
		instanceTypes = append(instanceTypes, &currInsType)
	}
	for i := range backupInstanceTypes {
		insType := backupInstanceTypes[i].(string)
		instanceTypes = append(instanceTypes, &insType)
	}

	return instanceTypes
}

func composedKubernetesAsScalingConfigPara() map[string]*schema.Schema {
	needSchema := map[string]*schema.Schema{
		"instance_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specified types of CVM instance.",
		},
		"backup_instance_types": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Backup CVM instance types if specified instance type sold out or mismatch.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"system_disk_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
			ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
			Description:  "Type of a CVM disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`. Default is `CLOUD_PREMIUM`.",
		},
		"system_disk_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      50,
			ValidateFunc: validateIntegerInRange(20, 1024),
			Description:  "Volume of system disk in GB. Default is `50`.",
		},
		"data_disk": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Configurations of data disk.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:     schema.TypeString,
						Optional: true,
						//ForceNew:     true,
						Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
						ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
						Description:  "Types of disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`.",
					},
					"disk_size": {
						Type:     schema.TypeInt,
						Optional: true,
						//ForceNew:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is `0`.",
					},
					"snapshot_id": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Data disk snapshot ID.",
					},
					"delete_with_instance": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Indicates whether the disk remove after instance terminated. Default is `false`.",
					},
					"encrypt": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Specify whether to encrypt data disk, default: false. NOTE: Make sure the instance type is offering and the cam role `QcloudKMSAccessForCVMRole` was provided.",
					},
					"throughput_performance": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD` and `data_size` > 460GB.",
					},
				},
			},
		},
		// payment
		"instance_charge_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`. The default is `POSTPAID_BY_HOUR`. NOTE: `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.",
		},
		"instance_charge_type_prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue(CVM_PREPAID_PERIOD),
			Description:  "The tenancy (in month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
		},
		"instance_charge_type_prepaid_renew_flag": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validateAllowedStringValue(CVM_PREPAID_RENEW_FLAG),
			Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
		},
		"spot_instance_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateAllowedStringValue([]string{"one-time"}),
			Description:  "Type of spot instance, only support `one-time` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.",
		},
		"spot_max_price": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringNumber,
			Description:  "Max price of a spot instance, is the format of decimal string, for example \"0.50\". Note: it only works when instance_charge_type is set to `SPOTPAID`.",
		},
		"internet_charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
			ValidateFunc: validateAllowedStringValue(INTERNET_CHARGE_ALLOW_TYPE),
			Description:  "Charge types for network traffic. Valid value: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
		},
		"internet_max_bandwidth_out": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Max bandwidth of Internet access in Mbps. Default is `0`.",
		},
		"bandwidth_package_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.",
		},
		"public_ip_assigned": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Specify whether to assign an Internet IP address.",
		},
		"password": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ForceNew:      true,
			ValidateFunc:  validateAsConfigPassword,
			ConflictsWith: []string{"auto_scaling_config.0.key_ids"},
			Description:   "Password to access.",
		},
		"key_ids": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			Elem:          &schema.Schema{Type: schema.TypeString},
			ConflictsWith: []string{"auto_scaling_config.0.password"},
			Description:   "ID list of keys.",
		},
		"security_group_ids": {
			Type:          schema.TypeSet,
			Optional:      true,
			Computed:      true,
			Elem:          &schema.Schema{Type: schema.TypeString},
			ConflictsWith: []string{"auto_scaling_config.0.orderly_security_group_ids"},
			Deprecated:    "The order of elements in this field cannot be guaranteed. Use `orderly_security_group_ids` instead.",
			Description:   "(**Deprecated**) The order of elements in this field cannot be guaranteed. Use `orderly_security_group_ids` instead. Security groups to which a CVM instance belongs.",
		},
		"orderly_security_group_ids": {
			Type:          schema.TypeList,
			Optional:      true,
			Computed:      true,
			Elem:          &schema.Schema{Type: schema.TypeString},
			ConflictsWith: []string{"auto_scaling_config.0.security_group_ids"},
			Description:   "Ordered security groups to which a CVM instance belongs.",
		},
		"enhanced_security_service": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
			//ForceNew:    true,
			Description: "To specify whether to enable cloud security service. Default is TRUE.",
		},
		"enhanced_monitor_service": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
		},
		"cam_role_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Name of cam role.",
		},
		"instance_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Instance name, no more than 60 characters. For usage, refer to `InstanceNameSettings` in https://www.tencentcloud.com/document/product/377/31001.",
		},
		"host_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The hostname of the cloud server, dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows instances are not supported. Examples of other types (Linux, etc.): The character length is [2, 40], multiple periods are allowed, and there is a paragraph between the dots, and each paragraph is allowed to consist of letters (unlimited case), numbers and dashes (-). Pure numbers are not allowed. For usage, refer to `HostNameSettings` in https://www.tencentcloud.com/document/product/377/31001.",
		},
		"host_name_style": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The style of the host name of the cloud server, the value range includes ORIGINAL and UNIQUE, and the default is ORIGINAL. For usage, refer to `HostNameSettings` in https://www.tencentcloud.com/document/product/377/31001.",
		},
	}

	return needSchema
}

func resourceTencentCloudKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesNodePoolCreate,
		Read:   resourceKubernetesNodePoolRead,
		Delete: resourceKubernetesNodePoolDelete,
		Update: resourceKubernetesNodePoolUpdate,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the node pool. The name does not exceed 25 characters, and only supports Chinese, English, numbers, underscores, separators (`-`) and decimal points.",
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Maximum number of node.",
			},
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Minimum number of node.",
			},
			"desired_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Desired capacity of the node. If `enable_auto_scale` is set `true`, this will be a computed parameter.",
			},
			"enable_auto_scale": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate whether to enable auto scaling or not.",
			},
			"retry_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Available values for retry policies include `IMMEDIATE_RETRY` and `INCREMENTAL_INTERVALS`.",
				Default:     SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
				ValidateFunc: validateAllowedStringValue([]string{SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
					SCALING_GROUP_RETRY_POLICY_INCREMENTAL_INTERVALS, SCALING_GROUP_RETRY_POLICY_NO_RETRY}),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of VPC network.",
			},
			"subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID list of subnet, and for VPC it is required.",
			},
			"scaling_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Auto scaling mode. Valid values are `CLASSIC_SCALING`(scaling by create/destroy instances), " +
					"`WAKE_UP_STOPPED_SCALING`(Boot priority for expansion. When expanding the capacity, the shutdown operation is given priority to the shutdown of the instance." +
					" If the number of instances is still lower than the expected number of instances after the startup, the instance will be created, and the method of destroying the instance will still be used for shrinking)" +
					".",
			},
			"multi_zone_subnet_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{MultiZoneSubnetPolicyPriority,
					MultiZoneSubnetPolicyEquality}),
				Description: "Multi-availability zone/subnet policy. Valid values: PRIORITY and EQUALITY. Default value: PRIORITY.",
			},
			"node_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: TkeInstanceAdvancedSetting(),
				},
				Description: "Node config.",
			},
			"auto_scaling_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: composedKubernetesAsScalingConfigPara(),
				},
				Description: "Auto scaling config parameters.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels of kubernetes node pool created nodes. The label key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').",
			},
			"unschedulable": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.",
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of the taint. The taint key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the taint.",
						},
						"effect": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Effect of the taint. Valid values are: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.",
						},
					},
				},
				Description: "Taints of kubernetes node pool created nodes.",
			},
			"delete_keep_instance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate to keep the CVM instance when delete the node pool. Default is `true`.",
			},
			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the node pool deletion protection is enabled.",
			},
			"node_os": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tlinux2.4x86_64",
				Description: "Operating system of the cluster. Please refer to [TencentCloud Documentation](https://www.tencentcloud.com/document/product/457/46750?lang=en&pg=#list-of-public-images-supported-by-tke) for available values. Default is 'tlinux2.4x86_64'. This parameter will only affect new nodes, not including the existing nodes.",
			},
			"node_os_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "GENERAL",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("node_os"); ok {
						if strings.Contains(v.(string), "img-") {
							return true
						}
					}
					return false
				},
				Description: "The image version of the node. Valida values are `DOCKER_CUSTOMIZE` and `GENERAL`. Default is `GENERAL`. This parameter will only affect new nodes, not including the existing nodes.",
			},
			// asg pass through arguments
			"scaling_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of relative scaling group.",
			},
			"zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of auto scaling group available zones, for Basic network it is required.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scaling_group_project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Project ID the scaling group belongs to.",
			},
			"default_cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Seconds of scaling group cool down. Default value is `300`.",
			},
			"termination_policies": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Policy of scaling group termination. Available values: `[\"OLDEST_INSTANCE\"]`, `[\"NEWEST_INSTANCE\"]`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Node pool tag specifications, will passthroughs to the scaling instances.",
			},
			//computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the node pool.",
			},
			"node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total node count.",
			},
			"autoscaling_added_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total of autoscaling added node.",
			},
			"manually_added_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total of manually added node.",
			},
			"launch_config_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The launch config ID.",
			},
			"auto_scaling_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The auto scaling group ID.",
			},
		},
		Importer: &schema.ResourceImporter{
			//State: schema.ImportStatePassthrough,
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				importFlag = true
				err := resourceKubernetesNodePoolRead(d, m)
				if err != nil {
					return nil, fmt.Errorf("failed to import resource")
				}

				return []*schema.ResourceData{d}, nil
			},
		},
		//compare to console, miss cam_role and running_version and lock_initial_node and security_proof
	}
}

//this function composes every single parameter to an as scale parameter with json string format
func composeParameterToAsScalingGroupParaSerial(d *schema.ResourceData) (string, error) {
	var (
		result string
		errRet error
	)

	request := as.NewCreateAutoScalingGroupRequest()

	//this is an empty string
	request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	request.MinSize = helper.IntUint64(d.Get("min_size").(int))

	if *request.MinSize > *request.MaxSize {
		return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
	}

	request.VpcId = helper.String(d.Get("vpc_id").(string))

	if v, ok := d.GetOk("desired_capacity"); ok {
		request.DesiredCapacity = helper.IntUint64(v.(int))
		if *request.DesiredCapacity > *request.MaxSize ||
			*request.DesiredCapacity < *request.MinSize {
			return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
		}

	}

	if v, ok := d.GetOk("retry_policy"); ok {
		request.RetryPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIds := v.([]interface{})
		request.SubnetIds = helper.InterfacesStringsPoint(subnetIds)
	}

	if v, ok := d.GetOk("scaling_mode"); ok {
		request.ServiceSettings = &as.ServiceSettings{ScalingMode: helper.String(v.(string))}
	}

	if v, ok := d.GetOk("multi_zone_subnet_policy"); ok {
		request.MultiZoneSubnetPolicy = helper.String(v.(string))
	}

	result = request.ToJsonString()

	return result, errRet
}

//This function is used to specify tke as group launch config, similar to kubernetesAsScalingConfigParaSerial, but less parameter
func composedKubernetesAsScalingConfigParaSerial(dMap map[string]interface{}, meta interface{}) (string, error) {
	var (
		result string
		errRet error
	)

	request := as.NewCreateLaunchConfigurationRequest()

	instanceType := dMap["instance_type"].(string)
	request.InstanceType = &instanceType

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := dMap["system_disk_type"]; ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := dMap["data_disk"]; ok {
		dataDisks := v.([]interface{})
		//request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			deleteWithInstance, dOk := value["delete_with_instance"].(bool)
			encrypt, eOk := value["encrypt"].(bool)
			throughputPerformance := value["throughput_performance"].(int)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
			}
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			if dOk {
				dataDisk.DeleteWithInstance = &deleteWithInstance
			}
			if eOk {
				dataDisk.Encrypt = &encrypt
			}
			if throughputPerformance > 0 {
				dataDisk.ThroughputPerformance = helper.IntUint64(throughputPerformance)
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := dMap["internet_charge_type"]; ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := dMap["bandwidth_package_id"]; ok {
		if v.(string) != "" {
			request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
	}
	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	request.LoginSettings = &as.LoginSettings{}

	if v, ok := dMap["password"]; ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	if v, ok := dMap["key_ids"]; ok {
		keyIds := v.([]interface{})
		//request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
	}

	if request.LoginSettings.Password != nil && *request.LoginSettings.Password == "" {
		request.LoginSettings.Password = nil
	}

	if request.LoginSettings.Password == nil && len(request.LoginSettings.KeyIds) == 0 {
		errRet = fmt.Errorf("Parameters `key_ids` and `password` should be set one")
		return result, errRet
	}

	if request.LoginSettings.Password != nil && len(request.LoginSettings.KeyIds) != 0 {
		errRet = fmt.Errorf("Parameters `key_ids` and `password` can only be supported one")
		return result, errRet
	}

	if v, ok := dMap["security_group_ids"]; ok {
		if list := v.(*schema.Set).List(); len(list) > 0 {
			errRet = fmt.Errorf("The parameter `security_group_ids` has an issue that the actual order of the security group may be inconsistent with the order of your tf code, which will cause your service to be inaccessible. Please use `orderly_security_group_ids` instead.")
			return result, errRet
		}
	}

	if v, ok := dMap["orderly_security_group_ids"]; ok {
		if list := v.([]interface{}); len(list) > 0 {
			request.SecurityGroupIds = helper.InterfacesStringsPoint(list)
		}
	}

	request.EnhancedService = &as.EnhancedService{}

	if v, ok := dMap["enhanced_security_service"]; ok {
		securityService := v.(bool)
		request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := dMap["enhanced_monitor_service"]; ok {
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &as.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}

	chargeType, ok := dMap["instance_charge_type"].(string)
	if !ok || chargeType == "" {
		chargeType = INSTANCE_CHARGE_TYPE_POSTPAID
	}

	if chargeType == INSTANCE_CHARGE_TYPE_SPOTPAID {
		spotMaxPrice := dMap["spot_max_price"].(string)
		spotInstanceType := dMap["spot_instance_type"].(string)
		request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
			MarketType: helper.String("spot"),
			SpotOptions: &as.SpotMarketOptions{
				MaxPrice:         &spotMaxPrice,
				SpotInstanceType: &spotInstanceType,
			},
		}
	}

	if chargeType == INSTANCE_CHARGE_TYPE_PREPAID {
		period := dMap["instance_charge_type_prepaid_period"].(int)
		renewFlag := dMap["instance_charge_type_prepaid_renew_flag"].(string)
		request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: &renewFlag,
		}
	}

	request.InstanceChargeType = &chargeType

	if v, ok := dMap["cam_role_name"]; ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := dMap["instance_name"]; ok && v != "" {
		request.InstanceNameSettings = &as.InstanceNameSettings{
			InstanceName: helper.String(v.(string)),
		}
	}

	if v, ok := dMap["host_name"]; ok && v != "" {
		if request.HostNameSettings == nil {
			request.HostNameSettings = &as.HostNameSettings{
				HostName: helper.String(v.(string)),
			}
		} else {
			request.HostNameSettings.HostName = helper.String(v.(string))
		}
	}

	if v, ok := dMap["host_name_style"]; ok && v != "" {
		if request.HostNameSettings != nil {
			request.HostNameSettings.HostNameStyle = helper.String(v.(string))
		} else {
			request.HostNameSettings = &as.HostNameSettings{
				HostNameStyle: helper.String(v.(string)),
			}
		}
	}
	result = request.ToJsonString()
	return result, errRet
}

func composeAsLaunchConfigModifyRequest(d *schema.ResourceData, launchConfigId string) (*as.ModifyLaunchConfigurationAttributesRequest, error) {
	launchConfigRaw := d.Get("auto_scaling_config").([]interface{})
	dMap := launchConfigRaw[0].(map[string]interface{})
	request := as.NewModifyLaunchConfigurationAttributesRequest()
	request.LaunchConfigurationId = &launchConfigId

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := dMap["system_disk_type"]; ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := dMap["data_disk"]; ok {
		dataDisks := v.([]interface{})
		//request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			deleteWithInstance, dOk := value["delete_with_instance"].(bool)
			encrypt, eOk := value["encrypt"].(bool)
			throughputPerformance := value["throughput_performance"].(int)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
			}
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			if dOk {
				dataDisk.DeleteWithInstance = &deleteWithInstance
			}
			if eOk {
				dataDisk.Encrypt = &encrypt
			}
			if throughputPerformance > 0 {
				dataDisk.ThroughputPerformance = helper.IntUint64(throughputPerformance)
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := dMap["internet_charge_type"]; ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := dMap["bandwidth_package_id"]; ok {
		if v.(string) != "" {
			request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
	}
	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	if d.HasChange("auto_scaling_config.0.security_group_ids") {
		if v, ok := dMap["security_group_ids"]; ok {
			if list := v.(*schema.Set).List(); len(list) > 0 {
				errRet := fmt.Errorf("The parameter `security_group_ids` has an issue that the actual order of the security group may be inconsistent with the order of your tf code, which will cause your service to be inaccessible. You can check whether the order of your current security groups meets your expectations through the TencentCloud Console, then use `orderly_security_group_ids` field to update them.")
				return nil, errRet
			}
		}
	}

	if d.HasChange("auto_scaling_config.0.orderly_security_group_ids") {
		if v, ok := dMap["orderly_security_group_ids"]; ok {
			if list := v.([]interface{}); len(list) > 0 {
				request.SecurityGroupIds = helper.InterfacesStringsPoint(list)
			}
		}
	}

	chargeType, ok := dMap["instance_charge_type"].(string)

	if !ok || chargeType == "" {
		chargeType = INSTANCE_CHARGE_TYPE_POSTPAID
	}

	if chargeType == INSTANCE_CHARGE_TYPE_SPOTPAID {
		spotMaxPrice := dMap["spot_max_price"].(string)
		spotInstanceType := dMap["spot_instance_type"].(string)
		request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
			MarketType: helper.String("spot"),
			SpotOptions: &as.SpotMarketOptions{
				MaxPrice:         &spotMaxPrice,
				SpotInstanceType: &spotInstanceType,
			},
		}
	}

	if chargeType == INSTANCE_CHARGE_TYPE_PREPAID {
		period := dMap["instance_charge_type_prepaid_period"].(int)
		renewFlag := dMap["instance_charge_type_prepaid_renew_flag"].(string)
		request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: &renewFlag,
		}
	}

	if v, ok := dMap["instance_name"]; ok && v != "" {
		request.InstanceNameSettings = &as.InstanceNameSettings{
			InstanceName: helper.String(v.(string)),
		}
	}

	if v, ok := dMap["host_name"]; ok && v != "" {
		if request.HostNameSettings == nil {
			request.HostNameSettings = &as.HostNameSettings{
				HostName: helper.String(v.(string)),
			}
		} else {
			request.HostNameSettings.HostName = helper.String(v.(string))
		}
	}

	if v, ok := dMap["host_name_style"]; ok && v != "" {
		if request.HostNameSettings != nil {
			request.HostNameSettings.HostNameStyle = helper.String(v.(string))
		} else {
			request.HostNameSettings = &as.HostNameSettings{
				HostNameStyle: helper.String(v.(string)),
			}
		}
	}

	// set enhanced_security_service if necessary
	if v, ok := dMap["enhanced_security_service"]; ok {
		securityService := v.(bool)
		if request.EnhancedService != nil {
			request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
				Enabled: helper.Bool(securityService),
			}
		} else {
			request.EnhancedService = &as.EnhancedService{
				SecurityService: &as.RunSecurityServiceEnabled{
					Enabled: helper.Bool(securityService),
				},
			}
		}

	}

	request.InstanceChargeType = &chargeType

	return request, nil
}

func desiredCapacityOutRange(d *schema.ResourceData) bool {
	capacity := d.Get("desired_capacity").(int)
	minSize := d.Get("min_size").(int)
	maxSize := d.Get("max_size").(int)
	return capacity > maxSize || capacity < minSize
}

func resourceKubernetesNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_node_pool.read")()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		asService = AsService{client: meta.(*TencentCloudClient).apiV3Conn}
		items     = strings.Split(d.Id(), FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	_, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = service.DescribeCluster(ctx, clusterId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	//Describe Node Pool
	var (
		nodePool *tke.NodePool
	)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		nodePool, has, err = service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		status := *nodePool.AutoscalingGroupStatus

		if status == "enabling" || status == "disabling" {
			return resource.RetryableError(fmt.Errorf("node pool status is %s, retrying", status))
		}

		return nil
	})

	if err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", nodePool.Name)
	_ = d.Set("status", nodePool.LifeState)
	AutoscalingAddedTotal := *nodePool.NodeCountSummary.AutoscalingAdded.Total
	ManuallyAddedTotal := *nodePool.NodeCountSummary.ManuallyAdded.Total
	_ = d.Set("autoscaling_added_total", AutoscalingAddedTotal)
	_ = d.Set("manually_added_total", ManuallyAddedTotal)
	_ = d.Set("node_count", AutoscalingAddedTotal+ManuallyAddedTotal)
	_ = d.Set("auto_scaling_group_id", nodePool.AutoscalingGroupId)
	_ = d.Set("launch_config_id", nodePool.LaunchConfigurationId)
	if _, ok := d.GetOkExists("unschedulable"); !ok && importFlag {
		_ = d.Set("unschedulable", nodePool.Unschedulable)
	}
	//set not force new parameters
	if nodePool.MaxNodesNum != nil {
		_ = d.Set("max_size", nodePool.MaxNodesNum)
	}
	if nodePool.MinNodesNum != nil {
		_ = d.Set("min_size", nodePool.MinNodesNum)
	}
	if nodePool.DesiredNodesNum != nil {
		_ = d.Set("desired_capacity", nodePool.DesiredNodesNum)
	}
	if nodePool.AutoscalingGroupStatus != nil {
		_ = d.Set("enable_auto_scale", *nodePool.AutoscalingGroupStatus == "enabled")
	}
	//修复自定义镜像返回信息的不一致
	if nodePool.ImageId != nil && *nodePool.ImageId != "" {
		_ = d.Set("node_os", nodePool.ImageId)
	} else {
		if nodePool.NodePoolOs != nil {
			_ = d.Set("node_os", nodePool.NodePoolOs)
		}
		if nodePool.OsCustomizeType != nil {
			_ = d.Set("node_os_type", nodePool.OsCustomizeType)
		}
	}

	if tags := nodePool.Tags; tags != nil {
		tagMap := make(map[string]string)
		for i := range tags {
			tag := tags[i]
			tagMap[*tag.Key] = *tag.Value
		}
		_ = d.Set("tags", tagMap)
	}

	if nodePool.DeletionProtection != nil {
		_ = d.Set("deletion_protection", nodePool.DeletionProtection)
	}

	//set composed struct
	lables := make(map[string]interface{}, len(nodePool.Labels))
	for _, v := range nodePool.Labels {
		lables[*v.Name] = *v.Value
	}
	_ = d.Set("labels", lables)

	// set launch config
	launchCfg, hasLC, err := asService.DescribeLaunchConfigurationById(ctx, *nodePool.LaunchConfigurationId)

	if hasLC > 0 {
		launchConfig := make(map[string]interface{})
		if launchCfg.InstanceTypes != nil {
			insTypes := launchCfg.InstanceTypes
			launchConfig["instance_type"] = insTypes[0]
			backupInsTypes := insTypes[1:]
			if len(backupInsTypes) > 0 {
				launchConfig["backup_instance_types"] = helper.StringsInterfaces(backupInsTypes)
			}
		} else {
			launchConfig["instance_type"] = launchCfg.InstanceType
		}
		if launchCfg.SystemDisk.DiskType != nil {
			launchConfig["system_disk_type"] = launchCfg.SystemDisk.DiskType
		}
		if launchCfg.SystemDisk.DiskSize != nil {
			launchConfig["system_disk_size"] = launchCfg.SystemDisk.DiskSize
		}
		if launchCfg.InternetAccessible.InternetChargeType != nil {
			launchConfig["internet_charge_type"] = launchCfg.InternetAccessible.InternetChargeType
		}
		if launchCfg.InternetAccessible.InternetMaxBandwidthOut != nil {
			launchConfig["internet_max_bandwidth_out"] = launchCfg.InternetAccessible.InternetMaxBandwidthOut
		}
		if launchCfg.InternetAccessible.BandwidthPackageId != nil {
			launchConfig["bandwidth_package_id"] = launchCfg.InternetAccessible.BandwidthPackageId
		}
		if launchCfg.InternetAccessible.PublicIpAssigned != nil {
			launchConfig["public_ip_assigned"] = launchCfg.InternetAccessible.PublicIpAssigned
		}
		if launchCfg.InstanceChargeType != nil {
			launchConfig["instance_charge_type"] = launchCfg.InstanceChargeType
			if *launchCfg.InstanceChargeType == INSTANCE_CHARGE_TYPE_SPOTPAID && launchCfg.InstanceMarketOptions != nil {
				launchConfig["spot_instance_type"] = launchCfg.InstanceMarketOptions.SpotOptions.SpotInstanceType
				launchConfig["spot_max_price"] = launchCfg.InstanceMarketOptions.SpotOptions.MaxPrice
			}
			if *launchCfg.InstanceChargeType == INSTANCE_CHARGE_TYPE_PREPAID && launchCfg.InstanceChargePrepaid != nil {
				launchConfig["instance_charge_type_prepaid_period"] = launchCfg.InstanceChargePrepaid.Period
				launchConfig["instance_charge_type_prepaid_renew_flag"] = launchCfg.InstanceChargePrepaid.RenewFlag
			}
		}
		if len(launchCfg.DataDisks) > 0 {
			dataDisks := make([]map[string]interface{}, 0, len(launchCfg.DataDisks))
			for i := range launchCfg.DataDisks {
				item := launchCfg.DataDisks[i]
				disk := make(map[string]interface{})
				disk["disk_type"] = *item.DiskType
				disk["disk_size"] = *item.DiskSize
				if item.SnapshotId != nil {
					disk["snapshot_id"] = *item.SnapshotId
				}
				if item.DeleteWithInstance != nil {
					disk["delete_with_instance"] = *item.DeleteWithInstance
				}
				if item.Encrypt != nil {
					disk["encrypt"] = *item.Encrypt
				}
				if item.ThroughputPerformance != nil {
					disk["throughput_performance"] = *item.ThroughputPerformance
				}
				dataDisks = append(dataDisks, disk)
			}
			launchConfig["data_disk"] = dataDisks
		}
		if launchCfg.LoginSettings != nil {
			launchConfig["key_ids"] = helper.StringsInterfaces(launchCfg.LoginSettings.KeyIds)
		}
		// keep existing password in new launchConfig object
		if v, ok := d.GetOk("auto_scaling_config.0.password"); ok {
			launchConfig["password"] = v.(string)
		}

		if launchCfg.SecurityGroupIds != nil {
			launchConfig["security_group_ids"] = helper.StringsInterfaces(launchCfg.SecurityGroupIds)
			launchConfig["orderly_security_group_ids"] = helper.StringsInterfaces(launchCfg.SecurityGroupIds)
		}

		enableSecurity := launchCfg.EnhancedService.SecurityService.Enabled
		enableMonitor := launchCfg.EnhancedService.MonitorService.Enabled
		// Only declared or diff from exist will set.
		if _, ok := d.GetOk("enhanced_security_service"); ok || enableSecurity != nil {
			launchConfig["enhanced_security_service"] = *enableSecurity
		}
		if _, ok := d.GetOk("enhanced_monitor_service"); ok || enableMonitor != nil {
			launchConfig["enhanced_monitor_service"] = *enableMonitor
		}
		if _, ok := d.GetOk("cam_role_name"); ok || launchCfg.CamRoleName != nil {
			launchConfig["cam_role_name"] = launchCfg.CamRoleName
		}
		if launchCfg.InstanceNameSettings != nil && launchCfg.InstanceNameSettings.InstanceName != nil {
			launchConfig["instance_name"] = launchCfg.InstanceNameSettings.InstanceName
		}
		if launchCfg.HostNameSettings != nil && launchCfg.HostNameSettings.HostName != nil {
			launchConfig["host_name"] = launchCfg.HostNameSettings.HostName
		}
		if launchCfg.HostNameSettings != nil && launchCfg.HostNameSettings.HostNameStyle != nil {
			launchConfig["host_name_style"] = launchCfg.HostNameSettings.HostNameStyle
		}

		asgConfig := make([]interface{}, 0, 1)
		asgConfig = append(asgConfig, launchConfig)
		if err := d.Set("auto_scaling_config", asgConfig); err != nil {
			return err
		}
	}

	nodeConfig := make(map[string]interface{})
	nodeConfigs := make([]interface{}, 0, 1)
	if nodePool.DataDisks != nil && len(nodePool.DataDisks) > 0 {
		dataDisks := make([]interface{}, 0, len(nodePool.DataDisks))
		for i := range nodePool.DataDisks {
			item := nodePool.DataDisks[i]
			disk := make(map[string]interface{})
			disk["disk_type"] = helper.PString(item.DiskType)
			disk["disk_size"] = helper.PInt64(item.DiskSize)
			disk["file_system"] = helper.PString(item.FileSystem)
			disk["auto_format_and_mount"] = helper.PBool(item.AutoFormatAndMount)
			disk["mount_target"] = helper.PString(item.MountTarget)
			disk["disk_partition"] = helper.PString(item.MountTarget)
			dataDisks = append(dataDisks, disk)
		}
		nodeConfig["data_disk"] = dataDisks
	}

	if helper.PInt64(nodePool.DesiredPodNum) != 0 {
		nodeConfig["desired_pod_num"] = helper.PInt64(nodePool.DesiredPodNum)
	}

	if helper.PInt64(nodePool.Unschedulable) != 0 {
		nodeConfig["is_schedule"] = false
	} else {
		nodeConfig["is_schedule"] = true
	}

	if helper.PString(nodePool.DockerGraphPath) != "" {
		nodeConfig["docker_graph_path"] = helper.PString(nodePool.DockerGraphPath)
	} else {
		nodeConfig["docker_graph_path"] = "/var/lib/docker"
	}

	if importFlag {
		if nodePool.ExtraArgs != nil && len(nodePool.ExtraArgs.Kubelet) > 0 {
			extraArgs := make([]string, 0)
			for i := range nodePool.ExtraArgs.Kubelet {
				extraArgs = append(extraArgs, helper.PString(nodePool.ExtraArgs.Kubelet[i]))
			}
			nodeConfig["extra_args"] = extraArgs
		}

		if helper.PString(nodePool.UserScript) != "" {
			nodeConfig["user_data"] = helper.PString(nodePool.UserScript)
		}

		if nodePool.GPUArgs != nil {
			setting := nodePool.GPUArgs
			var driverEmptyFlag, cudaEmptyFlag, cudnnEmptyFlag, customDriverEmptyFlag bool
			gpuArgs := map[string]interface{}{
				"mig_enable": helper.PBool(setting.MIGEnable),
			}

			if !isDriverEmpty(setting.Driver) {
				driverEmptyFlag = true
				driver := map[string]interface{}{
					"version": helper.PString(setting.Driver.Version),
					"name":    helper.PString(setting.Driver.Name),
				}
				gpuArgs["driver"] = driver
			}

			if !isCUDAEmpty(setting.CUDA) {
				cudaEmptyFlag = true
				cuda := map[string]interface{}{
					"version": helper.PString(setting.CUDA.Version),
					"name":    helper.PString(setting.CUDA.Name),
				}
				gpuArgs["cuda"] = cuda
			}

			if !isCUDNNEmpty(setting.CUDNN) {
				cudnnEmptyFlag = true
				cudnn := map[string]interface{}{
					"version":  helper.PString(setting.CUDNN.Version),
					"name":     helper.PString(setting.CUDNN.Name),
					"doc_name": helper.PString(setting.CUDNN.DocName),
					"dev_name": helper.PString(setting.CUDNN.DevName),
				}
				gpuArgs["cudnn"] = cudnn
			}

			if !isCustomDriverEmpty(setting.CustomDriver) {
				customDriverEmptyFlag = true
				customDriver := map[string]interface{}{
					"address": helper.PString(setting.CustomDriver.Address),
				}
				gpuArgs["custom_driver"] = customDriver
			}
			if driverEmptyFlag || cudaEmptyFlag || cudnnEmptyFlag || customDriverEmptyFlag {
				nodeConfig["gpu_args"] = []map[string]interface{}{gpuArgs}
			}
		}
		nodeConfigs = append(nodeConfigs, nodeConfig)
		_ = d.Set("node_config", nodeConfigs)
		importFlag = false
	}

	// Relative scaling group status
	asg, hasAsg, err := asService.DescribeAutoScalingGroupById(ctx, *nodePool.AutoscalingGroupId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			asg, hasAsg, err = asService.DescribeAutoScalingGroupById(ctx, *nodePool.AutoscalingGroupId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if hasAsg > 0 {
		_ = d.Set("scaling_group_name", asg.AutoScalingGroupName)
		_ = d.Set("zones", asg.ZoneSet)
		_ = d.Set("scaling_group_project_id", asg.ProjectId)
		_ = d.Set("default_cooldown", asg.DefaultCooldown)
		_ = d.Set("termination_policies", helper.StringsInterfaces(asg.TerminationPolicySet))
		_ = d.Set("vpc_id", asg.VpcId)
		_ = d.Set("retry_policy", asg.RetryPolicy)
		_ = d.Set("subnet_ids", helper.StringsInterfaces(asg.SubnetIdSet))
		if v, ok := d.GetOk("scaling_mode"); ok {
			if asg.ServiceSettings != nil && asg.ServiceSettings.ScalingMode != nil {
				_ = d.Set("scaling_mode", helper.PString(asg.ServiceSettings.ScalingMode))
			} else {
				_ = d.Set("scaling_mode", v.(string))
			}
		}
		// If not check, the diff between computed and default empty value leads to force replacement
		if _, ok := d.GetOk("multi_zone_subnet_policy"); ok {
			_ = d.Set("multi_zone_subnet_policy", asg.MultiZoneSubnetPolicy)
		}
	}
	if v, ok := d.GetOkExists("delete_keep_instance"); ok {
		_ = d.Set("delete_keep_instance", v.(bool))
	} else {
		_ = d.Set("delete_keep_instance", true)
	}

	taints := make([]map[string]interface{}, len(nodePool.Taints))
	for i, v := range nodePool.Taints {
		taint := map[string]interface{}{
			"key":    *v.Key,
			"value":  *v.Value,
			"effect": *v.Effect,
		}
		taints[i] = taint
	}
	_ = d.Set("taints", taints)

	return nil
}

func resourceKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_node_pool.create")()
	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		clusterId       = d.Get("cluster_id").(string)
		enableAutoScale = d.Get("enable_auto_scale").(bool)
		configParas     = d.Get("auto_scaling_config").([]interface{})
		name            = d.Get("name").(string)
		iAdvanced       tke.InstanceAdvancedSettings
	)
	if len(configParas) != 1 {
		return fmt.Errorf("need only one auto_scaling_config")
	}

	groupParaStr, err := composeParameterToAsScalingGroupParaSerial(d)
	if err != nil {
		return err
	}

	configParaStr, err := composedKubernetesAsScalingConfigParaSerial(configParas[0].(map[string]interface{}), meta)
	if err != nil {
		return err
	}

	labels := GetTkeLabels(d, "labels")
	taints := GetTkeTaints(d, "taints")

	//compose InstanceAdvancedSettings
	if workConfig, ok := helper.InterfacesHeadMap(d, "node_config"); ok {
		iAdvanced = tkeGetInstanceAdvancedPara(workConfig, meta)
	}

	if temp, ok := d.GetOk("extra_args"); ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		for _, extraArg := range extraArgs {
			iAdvanced.ExtraArgs.Kubelet = append(iAdvanced.ExtraArgs.Kubelet, &extraArg)
		}
	}
	if temp, ok := d.GetOk("unschedulable"); ok {
		iAdvanced.Unschedulable = helper.Int64(int64(temp.(int)))
	}

	nodeOs := d.Get("node_os").(string)
	nodeOsType := d.Get("node_os_type").(string)
	//自定镜像不能指定节点操作系统类型
	if strings.Contains(nodeOs, "img-") {
		nodeOsType = ""
	}

	deletionProtection := d.Get("deletion_protection").(bool)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	nodePoolId, err := service.CreateClusterNodePool(ctx, clusterId, name, groupParaStr, configParaStr, enableAutoScale, nodeOs, nodeOsType, labels, taints, iAdvanced, deletionProtection)
	if err != nil {
		return err
	}

	d.SetId(clusterId + FILED_SP + nodePoolId)

	// wait for status ok
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		nodePool, _, errRet := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if nodePool != nil && *nodePool.LifeState == "normal" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("node pool status is %s, retry...", *nodePool.LifeState))
	})
	if err != nil {
		return err
	}

	instanceTypes := getNodePoolInstanceTypes(d)

	if len(instanceTypes) != 0 {
		err := service.ModifyClusterNodePoolInstanceTypes(ctx, clusterId, nodePoolId, instanceTypes)
		if err != nil {
			return err
		}
	}

	//modify os, instanceTypes and image
	err = resourceKubernetesNodePoolUpdate(d, meta)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_node_pool.update")()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		client     = meta.(*TencentCloudClient).apiV3Conn
		service    = TkeService{client: client}
		asService  = AsService{client: client}
		cvmService = CvmService{client: client}
		items      = strings.Split(d.Id(), FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	d.Partial(true)

	// LaunchConfig
	if d.HasChange("auto_scaling_config") {
		nodePool, _, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if err != nil {
			return err
		}
		launchConfigId := *nodePool.LaunchConfigurationId
		//  change as config here
		request, composeError := composeAsLaunchConfigModifyRequest(d, launchConfigId)
		if composeError != nil {
			return composeError
		}
		_, err = client.UseAsClient().ModifyLaunchConfigurationAttributes(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return err
		}

		// change existed cvm security service if necessary
		if err := ModifySecurityServiceOfCvmInNodePool(ctx, d, &service, &cvmService, client, clusterId, *nodePool.NodePoolId); err != nil {
			return err
		}

	}

	var capacityHasChanged = false
	// assuming
	// min 1 max 6 desired 2
	// to
	// min 3 max 6 desired 5
	// modify min/max first will cause error, this case must upgrade desired first
	if d.HasChange("desired_capacity") || !desiredCapacityOutRange(d) {
		desiredCapacity := int64(d.Get("desired_capacity").(int))
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePoolDesiredCapacity(ctx, clusterId, nodePoolId, desiredCapacity)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		capacityHasChanged = true
	}

	// ModifyClusterNodePool
	if d.HasChanges(
		"min_size",
		"max_size",
		"name",
		"labels",
		"taints",
		"deletion_protection",
		"enable_auto_scale",
		"node_os_type",
		"node_os",
		"tags",
	) {
		maxSize := int64(d.Get("max_size").(int))
		minSize := int64(d.Get("min_size").(int))
		enableAutoScale := d.Get("enable_auto_scale").(bool)
		deletionProtection := d.Get("deletion_protection").(bool)
		name := d.Get("name").(string)
		labels := GetTkeLabels(d, "labels")
		taints := GetTkeTaints(d, "taints")
		tags := helper.GetTags(d, "tags")
		nodeOs := d.Get("node_os").(string)
		nodeOsType := d.Get("node_os_type").(string)
		//自定镜像不能指定节点操作系统类型
		if strings.Contains(nodeOs, "img-") {
			nodeOsType = ""
		}
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePool(ctx, clusterId, nodePoolId, name, enableAutoScale, minSize, maxSize, nodeOs, nodeOsType, labels, taints, tags, deletionProtection)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	// ModifyScalingGroup
	if d.HasChange("scaling_group_name") ||
		d.HasChange("zones") ||
		d.HasChange("scaling_group_project_id") ||
		d.HasChange("multi_zone_subnet_policy") ||
		d.HasChange("default_cooldown") ||
		d.HasChange("termination_policies") {

		nodePool, _, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if err != nil {
			return err
		}

		var (
			request               = as.NewModifyAutoScalingGroupRequest()
			scalingGroupId        = *nodePool.AutoscalingGroupId
			name                  = d.Get("scaling_group_name").(string)
			projectId             = d.Get("scaling_group_project_id").(int)
			defaultCooldown       = d.Get("default_cooldown").(int)
			multiZoneSubnetPolicy = d.Get("multi_zone_subnet_policy").(string)
		)

		request.AutoScalingGroupId = &scalingGroupId

		if name != "" {
			request.AutoScalingGroupName = &name
		}

		if multiZoneSubnetPolicy != "" {
			request.MultiZoneSubnetPolicy = &multiZoneSubnetPolicy
		}

		// It is safe to use Get() with default value 0.
		request.ProjectId = helper.IntUint64(projectId)

		if defaultCooldown != 0 {
			request.DefaultCooldown = helper.IntUint64(defaultCooldown)
		}

		if v, ok := d.GetOk("zones"); ok {
			request.Zones = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		if v, ok := d.GetOk("termination_policies"); ok {
			request.TerminationPolicies = helper.InterfacesStringsPoint(v.([]interface{}))
		}

		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := asService.ModifyAutoScalingGroup(ctx, request)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})

		if err != nil {
			return err
		}

	}

	if d.HasChange("desired_capacity") && !capacityHasChanged {
		desiredCapacity := int64(d.Get("desired_capacity").(int))
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePoolDesiredCapacity(ctx, clusterId, nodePoolId, desiredCapacity)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("auto_scaling_config.0.backup_instance_types") {
		instanceTypes := getNodePoolInstanceTypes(d)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePoolInstanceTypes(ctx, clusterId, nodePoolId, instanceTypes)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		_ = d.Set("auto_scaling_config.0.backup_instance_types", instanceTypes)
	}
	d.Partial(false)

	return resourceKubernetesNodePoolRead(d, meta)
}

func resourceKubernetesNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_node_pool.delete")()

	var (
		logId              = getLogId(contextNil)
		ctx                = context.WithValue(context.TODO(), logIdKey, logId)
		service            = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		items              = strings.Split(d.Id(), FILED_SP)
		deleteKeepInstance = d.Get("delete_keep_instance").(bool)
		deletionProtection = d.Get("deletion_protection").(bool)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	if deletionProtection {
		return fmt.Errorf("deletion protection was enabled, please set `deletion_protection` to `false` and apply first")
	}

	//delete as group
	hasDelete := false
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err := service.DeleteClusterNodePool(ctx, clusterId, nodePoolId, deleteKeepInstance)

		if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkErr.Code == "InternalError.Param" && strings.Contains(sdkErr.Message, "Not Found") {
				hasDelete = true
				return nil
			}
		}
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if hasDelete {
		return nil
	}

	// wait for delete ok
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		nodePool, has, errRet := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if errRet != nil {
			errCode := errRet.(*sdkErrors.TencentCloudSDKError).Code
			if errCode == "InternalError.UnexpectedInternal" {
				return nil
			}
			return retryError(errRet, InternalError)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("node pool %s still alive, status %s", nodePoolId, *nodePool.LifeState))
		}
		return nil
	})

	return err
}

func isCUDNNEmpty(cudnn *tke.CUDNN) bool {
	return cudnn == nil || (helper.PString(cudnn.Version) == "" && helper.PString(cudnn.Name) == "" && helper.PString(cudnn.DocName) == "" && helper.PString(cudnn.DevName) == "")
}

func isCUDAEmpty(cuda *tke.DriverVersion) bool {
	return cuda == nil || (helper.PString(cuda.Version) == "" && helper.PString(cuda.Name) == "")
}

func isDriverEmpty(driver *tke.DriverVersion) bool {
	return driver == nil || (helper.PString(driver.Version) == "" && helper.PString(driver.Name) == "")
}

func isCustomDriverEmpty(customDriver *tke.CustomDriver) bool {
	return customDriver == nil || helper.PString(customDriver.Address) == ""
}
