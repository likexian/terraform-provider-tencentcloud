package cvm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmLaunchTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmLaunchTemplateCreate,
		Read:   resourceTencentCloudCvmLaunchTemplateRead,
		Delete: resourceTencentCloudCvmLaunchTemplateDelete,
		Schema: map[string]*schema.Schema{
			"launch_template_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of launch template.",
			},

			"placement": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The location of instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The available zone ID of the instance.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The project ID of the instance.",
						},
						"host_ids": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "The CDH ID list of the instance(input).",
						},
						"host_ips": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Deprecated:  "It has been deprecated from version 1.81.108.",
							Description: "Specify the host machine ip.",
						},
					},
				},
			},

			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID.",
			},

			"launch_template_version_description": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance launch template version description.",
			},

			"instance_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The type of the instance. If this parameter is not specified, the system will dynamically specify the default model according to the resource sales in the current region.",
			},

			"system_disk": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "System disk configuration information of the instance. If this parameter is not specified, it is assigned according to the system default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of system disk.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "System disk ID.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The size of system disk.",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud Dedicated Cluster(CDC) ID.",
						},
					},
				},
			},

			"data_disks": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Data disk configuration information of the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The size of the data disk.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of data disk.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk ID.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the data disk is destroyed along with the instance, true or false.",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk snapshot ID.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the data disk is encrypted, TRUE or FALSE.",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The id of custom CMK.",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Cloud disk performance, MB/s.",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud Dedicated Cluster(CDC) ID.",
						},
					},
				},
			},

			"virtual_private_cloud": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The configuration information of VPC. If this parameter is not specified, the basic network is used by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of subnet.",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is it used as a Public network gateway, TRUE or FALSE.",
						},
						"private_ip_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The address of private ip.",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of ipv6 addresses for Elastic Network Interface.",
						},
					},
				},
			},

			"internet_accessible": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The information settings of public network bandwidth. If you do not specify this parameter, the default Internet bandwidth is 0 Mbps.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"internet_charge_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of internet charge.",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Internet outbound bandwidth upper limit, Mbps.",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to allocate public network IP, TRUE or FALSE.",
						},
						"bandwidth_package_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of bandwidth package.",
						},
					},
				},
			},

			"instance_count": {
				Optional:    true,
				ForceNew:    true,
				Default:     1,
				Type:        schema.TypeInt,
				Description: "The number of instances purchased.",
			},

			"instance_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of instance. If you do not specify an instance display name, 'Unnamed' is displayed by default.",
			},

			"login_settings": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The login settings of instance. By default, passwords are randomly generated and notified to users via internal messages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The login password of instance.",
						},
						"key_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "List of key ID.",
						},
						"keep_image_login": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Keep the original settings of the mirror.",
						},
					},
				},
			},

			"security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The security group ID of instance. If this parameter is not specified, the default security group is bound.",
			},

			"enhanced_service": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Enhanced service. If this parameter is not specified, cloud monitoring and cloud security services will be enabled by default in public images.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Enable cloud security service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to enable cloud security service, TRUE or FALSE.",
									},
								},
							},
						},
						"monitor_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Enable cloud monitor service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to enable cloud monitor service, TRUE or FALSE.",
									},
								},
							},
						},
						"automation_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Enable TencentCloud Automation Tools(TAT).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to enable TencentCloud Automation Tools(TAT), TRUE or FALSE.",
									},
								},
							},
						},
					},
				},
			},

			"client_token": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A string to used guarantee request idempotency.",
			},

			"host_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The host name of CVM.",
			},

			"action_timer": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Timed task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timer_action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Timer name.",
						},
						"action_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Execution time.",
						},
						"externals": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Extended data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"release_address": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Release address.",
									},
									"unsupport_networks": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Unsupported network type.",
									},
									"storage_block_attr": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "HDD local storage attributes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The type of HDD local storage.",
												},
												"min_size": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The minimum capacity of HDD local storage.",
												},
												"max_size": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The maximum capacity of HDD local storage.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"disaster_recover_group_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ID of disaster recover group.",
			},

			"tag_specification": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Tag description list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of resource.",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Tag list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key of tag.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The value of tag.",
									},
								},
							},
						},
					},
				},
			},

			"instance_market_options": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The marketplace options of instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spot_options": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Bidding related options.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_price": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Bidding.",
									},
									"spot_instance_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Bidding request type, currently only supported type: one-time.",
									},
								},
							},
						},
						"market_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Market option type, currently only supports value: spot.",
						},
					},
				},
			},

			"user_data": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The data of users.",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to preflight only this request, true or false.",
			},

			"cam_role_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The role name of CAM.",
			},

			"hpc_cluster_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of HPC cluster.",
			},

			"instance_charge_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The charge type of instance. Default value: POSTPAID_BY_HOUR.",
			},

			"instance_charge_prepaid": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The configuration of charge prepaid.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The period of purchasing instances.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Automatic renew flag.",
						},
					},
				},
			},

			"disable_api_termination": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Instance destruction protection flag.",
			},

			"tags": {
				Type:        schema.TypeMap,
				ForceNew:    true,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudCvmLaunchTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_launch_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request          = cvm.NewCreateLaunchTemplateRequest()
		response         = cvm.NewCreateLaunchTemplateResponse()
		launchTemplateId string
	)
	if v, ok := d.GetOk("launch_template_name"); ok {
		request.LaunchTemplateName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "placement"); ok {
		placement := cvm.Placement{}
		if v, ok := dMap["zone"]; ok {
			placement.Zone = helper.String(v.(string))
		}
		if v, ok := dMap["project_id"]; ok {
			placement.ProjectId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["host_ids"]; ok {
			hostIdsSet := v.(*schema.Set).List()
			for i := range hostIdsSet {
				hostIds := hostIdsSet[i].(string)
				placement.HostIds = append(placement.HostIds, &hostIds)
			}
		}
		// It has been deprecated from version 1.81.108
		//if v, ok := dMap["host_ips"]; ok {
		//	hostIpsSet := v.(*schema.Set).List()
		//	for i := range hostIpsSet {
		//		hostIps := hostIpsSet[i].(string)
		//		placement.HostIps = append(placement.HostIps, &hostIps)
		//	}
		//}
		if v, ok := dMap["host_id"]; ok {
			placement.HostId = helper.String(v.(string))
		}
		request.Placement = &placement
	}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("launch_template_version_description"); ok {
		request.LaunchTemplateVersionDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "system_disk"); ok {
		systemDisk := cvm.SystemDisk{}
		if v, ok := dMap["disk_type"]; ok {
			systemDisk.DiskType = helper.String(v.(string))
		}
		if v, ok := dMap["disk_id"]; ok {
			systemDisk.DiskId = helper.String(v.(string))
		}
		if v, ok := dMap["disk_size"]; ok {
			systemDisk.DiskSize = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["cdc_id"]; ok {
			systemDisk.CdcId = helper.String(v.(string))
		}
		request.SystemDisk = &systemDisk
	}

	if v, ok := d.GetOk("data_disks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataDisk := cvm.DataDisk{}
			if v, ok := dMap["disk_size"]; ok {
				dataDisk.DiskSize = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["disk_type"]; ok {
				dataDisk.DiskType = helper.String(v.(string))
			}
			if v, ok := dMap["disk_id"]; ok {
				dataDisk.DiskId = helper.String(v.(string))
			}
			if v, ok := dMap["delete_with_instance"]; ok {
				dataDisk.DeleteWithInstance = helper.Bool(v.(bool))
			}
			if v, ok := dMap["snapshot_id"]; ok {
				dataDisk.SnapshotId = helper.String(v.(string))
			}
			if v, ok := dMap["encrypt"]; ok {
				dataDisk.Encrypt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["kms_key_id"]; ok {
				dataDisk.KmsKeyId = helper.String(v.(string))
			}
			if v, ok := dMap["throughput_performance"]; ok {
				dataDisk.ThroughputPerformance = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["cdc_id"]; ok {
				dataDisk.CdcId = helper.String(v.(string))
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			virtualPrivateCloud.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			virtualPrivateCloud.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["as_vpc_gateway"]; ok {
			virtualPrivateCloud.AsVpcGateway = helper.Bool(v.(bool))
		}
		if v, ok := dMap["private_ip_addresses"]; ok {
			privateIpAddressesSet := v.(*schema.Set).List()
			for i := range privateIpAddressesSet {
				privateIpAddresses := privateIpAddressesSet[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		request.VirtualPrivateCloud = &virtualPrivateCloud
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "internet_accessible"); ok {
		internetAccessible := cvm.InternetAccessible{}
		if v, ok := dMap["internet_charge_type"].(string); ok && v != "" {
			internetAccessible.InternetChargeType = helper.String(v)
		}
		if v, ok := dMap["internet_max_bandwidth_out"].(int); ok && v != 0 {
			internetAccessible.InternetMaxBandwidthOut = helper.IntInt64(v)
		}
		if v, ok := dMap["public_ip_assigned"].(bool); ok && v {
			internetAccessible.PublicIpAssigned = helper.Bool(v)
		}
		if v, ok := dMap["bandwidth_package_id"].(string); ok && v != "" {
			internetAccessible.BandwidthPackageId = helper.String(v)
		}
		request.InternetAccessible = &internetAccessible
	}

	if v, _ := d.GetOk("instance_count"); v != nil {
		request.InstanceCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "login_settings"); ok {
		loginSettings := cvm.LoginSettings{}
		if v, ok := dMap["password"]; ok {
			loginSettings.Password = helper.String(v.(string))
		}
		if v, ok := dMap["key_ids"]; ok {
			keyIdsSet := v.(*schema.Set).List()
			for i := range keyIdsSet {
				keyIds := keyIdsSet[i].(string)
				loginSettings.KeyIds = append(loginSettings.KeyIds, &keyIds)
			}
		}
		if v, ok := dMap["keep_image_login"]; ok {
			loginSettings.KeepImageLogin = helper.String(v.(string))
		}
		request.LoginSettings = &loginSettings
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "enhanced_service"); ok {
		enhancedService := cvm.EnhancedService{}
		if securityServiceMap, ok := helper.InterfaceToMap(dMap, "security_service"); ok {
			runSecurityServiceEnabled := cvm.RunSecurityServiceEnabled{}
			if v, ok := securityServiceMap["enabled"]; ok {
				runSecurityServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.SecurityService = &runSecurityServiceEnabled
		}
		if monitorServiceMap, ok := helper.InterfaceToMap(dMap, "monitor_service"); ok {
			runMonitorServiceEnabled := cvm.RunMonitorServiceEnabled{}
			if v, ok := monitorServiceMap["enabled"]; ok {
				runMonitorServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.MonitorService = &runMonitorServiceEnabled
		}
		if automationServiceMap, ok := helper.InterfaceToMap(dMap, "automation_service"); ok {
			runAutomationServiceEnabled := cvm.RunAutomationServiceEnabled{}
			if v, ok := automationServiceMap["enabled"]; ok {
				runAutomationServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.AutomationService = &runAutomationServiceEnabled
		}
		request.EnhancedService = &enhancedService
	}

	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host_name"); ok {
		request.HostName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "action_timer"); ok {
		actionTimer := cvm.ActionTimer{}
		if v, ok := dMap["timer_action"]; ok {
			actionTimer.TimerAction = helper.String(v.(string))
		}
		if v, ok := dMap["action_time"]; ok {
			actionTimer.ActionTime = helper.String(v.(string))
		}
		if externalsMap, ok := helper.InterfaceToMap(dMap, "externals"); ok {
			externals := cvm.Externals{}
			if v, ok := externalsMap["release_address"]; ok {
				externals.ReleaseAddress = helper.Bool(v.(bool))
			}
			if v, ok := externalsMap["unsupport_networks"]; ok {
				unsupportNetworksSet := v.(*schema.Set).List()
				for i := range unsupportNetworksSet {
					unsupportNetworks := unsupportNetworksSet[i].(string)
					externals.UnsupportNetworks = append(externals.UnsupportNetworks, &unsupportNetworks)
				}
			}
			if storageBlockAttrMap, ok := helper.InterfaceToMap(externalsMap, "storage_block_attr"); ok {
				storageBlock := cvm.StorageBlock{}
				if v, ok := storageBlockAttrMap["type"]; ok {
					storageBlock.Type = helper.String(v.(string))
				}
				if v, ok := storageBlockAttrMap["min_size"]; ok {
					storageBlock.MinSize = helper.IntInt64(v.(int))
				}
				if v, ok := storageBlockAttrMap["max_size"]; ok {
					storageBlock.MaxSize = helper.IntInt64(v.(int))
				}
				externals.StorageBlockAttr = &storageBlock
			}
			actionTimer.Externals = &externals
		}
		request.ActionTimer = &actionTimer
	}

	if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
		disasterRecoverGroupIdsSet := v.(*schema.Set).List()
		for i := range disasterRecoverGroupIdsSet {
			disasterRecoverGroupIds := disasterRecoverGroupIdsSet[i].(string)
			request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &disasterRecoverGroupIds)
		}
	}

	if v, ok := d.GetOk("tag_specification"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagSpecification := cvm.TagSpecification{}
			if v, ok := dMap["resource_type"]; ok {
				tagSpecification.ResourceType = helper.String(v.(string))
			}
			if v, ok := dMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tag := cvm.Tag{}
					if v, ok := tagsMap["key"]; ok {
						tag.Key = helper.String(v.(string))
					}
					if v, ok := tagsMap["value"]; ok {
						tag.Value = helper.String(v.(string))
					}
					tagSpecification.Tags = append(tagSpecification.Tags, &tag)
				}
			}
			request.TagSpecification = append(request.TagSpecification, &tagSpecification)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_market_options"); ok {
		instanceMarketOptionsRequest := cvm.InstanceMarketOptionsRequest{}
		if spotOptionsMap, ok := helper.InterfaceToMap(dMap, "spot_options"); ok {
			spotMarketOptions := cvm.SpotMarketOptions{}
			if v, ok := spotOptionsMap["max_price"]; ok {
				spotMarketOptions.MaxPrice = helper.String(v.(string))
			}
			if v, ok := spotOptionsMap["spot_instance_type"]; ok {
				spotMarketOptions.SpotInstanceType = helper.String(v.(string))
			}
			instanceMarketOptionsRequest.SpotOptions = &spotMarketOptions
		}
		if v, ok := dMap["market_type"]; ok {
			instanceMarketOptionsRequest.MarketType = helper.String(v.(string))
		}
		request.InstanceMarketOptions = &instanceMarketOptionsRequest
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("hpc_cluster_id"); ok {
		request.HpcClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_charge_prepaid"); ok {
		instanceChargePrepaid := cvm.InstanceChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			instanceChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			instanceChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		request.InstanceChargePrepaid = &instanceChargePrepaid
	}

	if v, _ := d.GetOk("disable_api_termination"); v != nil {
		request.DisableApiTermination = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().CreateLaunchTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm launchTemplate failed, reason:%+v", logId, err)
		return err
	}

	launchTemplateId = *response.Response.LaunchTemplateId
	d.SetId(launchTemplateId)

	return resourceTencentCloudCvmLaunchTemplateRead(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_launch_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	launchTemplateId := d.Id()

	launchTemplate, err := service.DescribeCvmLaunchTemplateById(ctx, launchTemplateId)
	if err != nil {
		return err
	}

	if launchTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmLaunchTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if launchTemplate.LaunchTemplateName != nil {
		_ = d.Set("launch_template_name", launchTemplate.LaunchTemplateName)
	}

	launchTemplateVersion, err := service.DescribeLaunchTemplateVersionsById(ctx, launchTemplateId)
	if err != nil {
		return err
	}

	if launchTemplateVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmLaunchTemplateVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if launchTemplateVersion.LaunchTemplateVersionData != nil && launchTemplateVersion.LaunchTemplateVersionData.Placement != nil {
		placementMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.Placement.Zone != nil {
			placementMap["zone"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.Zone
		}

		if launchTemplateVersion.LaunchTemplateVersionData.Placement.ProjectId != nil {
			placementMap["project_id"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.ProjectId
		}

		if launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIds != nil {
			placementMap["host_ids"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIds
		}

		// It has been deprecated from version 1.81.108
		//if launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIps != nil {
		//	placementMap["host_ips"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIps
		//}

		if launchTemplateVersion.LaunchTemplateVersionData.Placement.HostId != nil {
			placementMap["host_id"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.HostId
		}

		_ = d.Set("placement", []interface{}{placementMap})
	}

	if launchTemplateVersion.LaunchTemplateVersionData != nil && launchTemplateVersion.LaunchTemplateVersionData.ImageId != nil {
		_ = d.Set("image_id", launchTemplateVersion.LaunchTemplateVersionData.ImageId)
	}

	if launchTemplateVersion.LaunchTemplateVersionDescription != nil {
		_ = d.Set("launch_template_version_description", launchTemplateVersion.LaunchTemplateVersionDescription)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.InstanceType != nil {
		_ = d.Set("instance_type", launchTemplateVersion.LaunchTemplateVersionData.InstanceType)
	}

	if launchTemplateVersion.LaunchTemplateVersionData != nil && launchTemplateVersion.LaunchTemplateVersionData.SystemDisk != nil {
		systemDiskMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskType != nil {
			systemDiskMap["disk_type"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskType
		}

		if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskId != nil {
			systemDiskMap["disk_id"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskId
		}

		if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskSize != nil {
			systemDiskMap["disk_size"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskSize
		}

		if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.CdcId != nil {
			systemDiskMap["cdc_id"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.CdcId
		}

		_ = d.Set("system_disk", []interface{}{systemDiskMap})
	}

	if launchTemplateVersion.LaunchTemplateVersionData.DataDisks != nil {
		dataDisksList := []interface{}{}
		for _, dataDisk := range launchTemplateVersion.LaunchTemplateVersionData.DataDisks {
			dataDisksMap := map[string]interface{}{}

			if dataDisk.DiskSize != nil {
				dataDisksMap["disk_size"] = dataDisk.DiskSize
			}

			if dataDisk.DiskType != nil {
				dataDisksMap["disk_type"] = dataDisk.DiskType
			}

			if dataDisk.DiskId != nil {
				dataDisksMap["disk_id"] = dataDisk.DiskId
			}

			if dataDisk.DeleteWithInstance != nil {
				dataDisksMap["delete_with_instance"] = dataDisk.DeleteWithInstance
			}

			if dataDisk.SnapshotId != nil {
				dataDisksMap["snapshot_id"] = dataDisk.SnapshotId
			}

			if dataDisk.Encrypt != nil {
				dataDisksMap["encrypt"] = dataDisk.Encrypt
			}

			if dataDisk.KmsKeyId != nil {
				dataDisksMap["kms_key_id"] = dataDisk.KmsKeyId
			}

			if dataDisk.ThroughputPerformance != nil {
				dataDisksMap["throughput_performance"] = dataDisk.ThroughputPerformance
			}

			if dataDisk.CdcId != nil {
				dataDisksMap["cdc_id"] = dataDisk.CdcId
			}

			dataDisksList = append(dataDisksList, dataDisksMap)
		}

		_ = d.Set("data_disks", dataDisksList)

	}

	if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud != nil {
		virtualPrivateCloudMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.VpcId != nil {
			virtualPrivateCloudMap["vpc_id"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.VpcId
		}

		if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.SubnetId != nil {
			virtualPrivateCloudMap["subnet_id"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.SubnetId
		}

		if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.AsVpcGateway != nil {
			virtualPrivateCloudMap["as_vpc_gateway"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.AsVpcGateway
		}

		if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.PrivateIpAddresses != nil {
			virtualPrivateCloudMap["private_ip_addresses"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.PrivateIpAddresses
		}

		if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.Ipv6AddressCount != nil {
			virtualPrivateCloudMap["ipv6_address_count"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.Ipv6AddressCount
		}

		_ = d.Set("virtual_private_cloud", []interface{}{virtualPrivateCloudMap})
	}

	if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible != nil {
		internetAccessibleMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetChargeType != nil {
			internetAccessibleMap["internet_charge_type"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetChargeType
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetMaxBandwidthOut != nil {
			internetAccessibleMap["internet_max_bandwidth_out"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetMaxBandwidthOut
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.PublicIpAssigned != nil {
			internetAccessibleMap["public_ip_assigned"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.PublicIpAssigned
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.BandwidthPackageId != nil {
			internetAccessibleMap["bandwidth_package_id"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.BandwidthPackageId
		}

		_ = d.Set("internet_accessible", []interface{}{internetAccessibleMap})
	}

	if launchTemplateVersion.LaunchTemplateVersionData.InstanceCount != nil {
		_ = d.Set("instance_count", launchTemplateVersion.LaunchTemplateVersionData.InstanceCount)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.InstanceName != nil {
		_ = d.Set("instance_name", launchTemplateVersion.LaunchTemplateVersionData.InstanceName)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings != nil {
		loginSettingsMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.Password != nil {
			loginSettingsMap["password"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.Password
		}

		if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeyIds != nil {
			loginSettingsMap["key_ids"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeyIds
		}

		if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeepImageLogin != nil {
			loginSettingsMap["keep_image_login"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeepImageLogin
		}

		_ = d.Set("login_settings", []interface{}{loginSettingsMap})
	}

	if launchTemplateVersion.LaunchTemplateVersionData.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", launchTemplateVersion.LaunchTemplateVersionData.SecurityGroupIds)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService != nil {
		enhancedServiceMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService != nil {
			securityServiceMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService.Enabled != nil {
				securityServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService.Enabled
			}

			enhancedServiceMap["security_service"] = []interface{}{securityServiceMap}
		}

		if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService != nil {
			monitorServiceMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService.Enabled != nil {
				monitorServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService.Enabled
			}

			enhancedServiceMap["monitor_service"] = []interface{}{monitorServiceMap}
		}

		if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService != nil {
			automationServiceMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService.Enabled != nil {
				automationServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService.Enabled
			}

			enhancedServiceMap["automation_service"] = []interface{}{automationServiceMap}
		}

		_ = d.Set("enhanced_service", []interface{}{enhancedServiceMap})
	}

	if launchTemplateVersion.LaunchTemplateVersionData.ClientToken != nil {
		_ = d.Set("client_token", launchTemplateVersion.LaunchTemplateVersionData.ClientToken)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.HostName != nil {
		_ = d.Set("host_name", launchTemplateVersion.LaunchTemplateVersionData.HostName)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer != nil {
		actionTimerMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.TimerAction != nil {
			actionTimerMap["timer_action"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.TimerAction
		}

		if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.ActionTime != nil {
			actionTimerMap["action_time"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.ActionTime
		}

		if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals != nil {
			externalsMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.ReleaseAddress != nil {
				externalsMap["release_address"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.ReleaseAddress
			}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.UnsupportNetworks != nil {
				externalsMap["unsupport_networks"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.UnsupportNetworks
			}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr != nil {
				storageBlockAttrMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.Type != nil {
					storageBlockAttrMap["type"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.Type
				}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MinSize != nil {
					storageBlockAttrMap["min_size"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MinSize
				}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MaxSize != nil {
					storageBlockAttrMap["max_size"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MaxSize
				}

				externalsMap["storage_block_attr"] = []interface{}{storageBlockAttrMap}
			}

			actionTimerMap["externals"] = []interface{}{externalsMap}
		}

		_ = d.Set("action_timer", []interface{}{actionTimerMap})
	}

	if launchTemplateVersion.LaunchTemplateVersionData.DisasterRecoverGroupIds != nil {
		_ = d.Set("disaster_recover_group_ids", launchTemplateVersion.LaunchTemplateVersionData.DisasterRecoverGroupIds)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.TagSpecification != nil {
		tagSpecificationList := []interface{}{}
		for _, tagSpecification := range launchTemplateVersion.LaunchTemplateVersionData.TagSpecification {
			tagSpecificationMap := map[string]interface{}{}

			if tagSpecification.ResourceType != nil {
				tagSpecificationMap["resource_type"] = tagSpecification.ResourceType
			}

			if tagSpecification.Tags != nil {
				tagsList := []interface{}{}
				for _, tag := range tagSpecification.Tags {
					tagsMap := map[string]interface{}{}

					if tag.Key != nil {
						tagsMap["key"] = tag.Key
					}

					if tag.Value != nil {
						tagsMap["value"] = tag.Value
					}

					tagsList = append(tagsList, tagsMap)
				}

				tagSpecificationMap["tags"] = []interface{}{tagsList}
			}

			tagSpecificationList = append(tagSpecificationList, tagSpecificationMap)
		}

		_ = d.Set("tag_specification", tagSpecificationList)

	}

	if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions != nil {
		instanceMarketOptionsMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions != nil {
			spotOptionsMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.MaxPrice != nil {
				spotOptionsMap["max_price"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.MaxPrice
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.SpotInstanceType != nil {
				spotOptionsMap["spot_instance_type"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.SpotInstanceType
			}

			instanceMarketOptionsMap["spot_options"] = []interface{}{spotOptionsMap}
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.MarketType != nil {
			instanceMarketOptionsMap["market_type"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.MarketType
		}

		_ = d.Set("instance_market_options", []interface{}{instanceMarketOptionsMap})
	}

	if launchTemplateVersion.LaunchTemplateVersionData.UserData != nil {
		_ = d.Set("user_data", launchTemplateVersion.LaunchTemplateVersionData.UserData)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.CamRoleName != nil {
		_ = d.Set("cam_role_name", launchTemplateVersion.LaunchTemplateVersionData.CamRoleName)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.HpcClusterId != nil {
		_ = d.Set("hpc_cluster_id", launchTemplateVersion.LaunchTemplateVersionData.HpcClusterId)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargeType != nil {
		_ = d.Set("instance_charge_type", launchTemplateVersion.LaunchTemplateVersionData.InstanceChargeType)
	}

	if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid != nil {
		instanceChargePrepaidMap := map[string]interface{}{}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.Period != nil {
			instanceChargePrepaidMap["period"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.Period
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.RenewFlag != nil {
			instanceChargePrepaidMap["renew_flag"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.RenewFlag
		}

		_ = d.Set("instance_charge_prepaid", []interface{}{instanceChargePrepaidMap})
	}

	// if launchTemplateVersion.LaunchTemplateVersionData.DisableApiTermination != nil {
	// 	_ = d.Set("disable_api_termination", launchTemplateVersion.LaunchTemplateVersionData.DisableApiTermination)
	// }

	return nil
}

func resourceTencentCloudCvmLaunchTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_launch_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	launchTemplateId := d.Id()

	if err := service.DeleteCvmLaunchTemplateById(ctx, launchTemplateId); err != nil {
		return err
	}

	return nil
}
