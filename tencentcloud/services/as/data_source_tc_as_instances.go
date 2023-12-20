package as

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID of the cloud server (CVM) to be queried. The limit is 100 per request.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions. If there are multiple Filters, the relationship between Filters is a logical AND (AND) relationship. If there are multiple Values in the same Filter, the relationship between Values under the same Filter is a logical OR (OR) relationship.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields to be filtered. Valid names: `instance-id`: Filters by instance ID, `auto-scaling-group-id`: Filter by scaling group ID.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Value of the field.",
						},
					},
				},
			},

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of instance details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"auto_scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling group ID.",
						},
						"auto_scaling_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling group name.",
						},
						"launch_configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Launch configuration ID.",
						},
						"launch_configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Launch configuration name.",
						},
						"life_cycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Life cycle state. Please refer to the link for field value details: https://cloud.tencent.com/document/api/377/20453#Instance.",
						},
						"health_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health status, the valid values are HEALTHY and UNHEALTHY.",
						},
						"protected_from_scale_in": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable scale in protection.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available zone.",
						},
						"creation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Valid values: `AUTO_CREATION`, `MANUAL_ATTACHING`.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the instance joined the group.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type.",
						},
						"version_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Version ID.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudAsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_as_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*as.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := as.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := AsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceList []*as.Instance

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAsInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	tmpList := make([]map[string]interface{}, 0, len(instanceList))

	if instanceList != nil {
		for _, instance := range instanceList {
			instanceMap := map[string]interface{}{}

			if instance.InstanceId != nil {
				instanceMap["instance_id"] = instance.InstanceId
			}

			if instance.AutoScalingGroupId != nil {
				instanceMap["auto_scaling_group_id"] = instance.AutoScalingGroupId
			}

			if instance.AutoScalingGroupName != nil {
				instanceMap["auto_scaling_group_name"] = instance.AutoScalingGroupName
			}

			if instance.LaunchConfigurationId != nil {
				instanceMap["launch_configuration_id"] = instance.LaunchConfigurationId
			}

			if instance.LaunchConfigurationName != nil {
				instanceMap["launch_configuration_name"] = instance.LaunchConfigurationName
			}

			if instance.LifeCycleState != nil {
				instanceMap["life_cycle_state"] = instance.LifeCycleState
			}

			if instance.HealthStatus != nil {
				instanceMap["health_status"] = instance.HealthStatus
			}

			if instance.ProtectedFromScaleIn != nil {
				instanceMap["protected_from_scale_in"] = instance.ProtectedFromScaleIn
			}

			if instance.Zone != nil {
				instanceMap["zone"] = instance.Zone
			}

			if instance.CreationType != nil {
				instanceMap["creation_type"] = instance.CreationType
			}

			if instance.AddTime != nil {
				instanceMap["add_time"] = instance.AddTime
			}

			if instance.InstanceType != nil {
				instanceMap["instance_type"] = instance.InstanceType
			}

			if instance.VersionNumber != nil {
				instanceMap["version_number"] = instance.VersionNumber
			}

			ids = append(ids, *instance.InstanceId)
			tmpList = append(tmpList, instanceMap)
		}

		_ = d.Set("instance_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
