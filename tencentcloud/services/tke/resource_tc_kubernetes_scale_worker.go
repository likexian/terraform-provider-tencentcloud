package tke

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTkeScaleWorker() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTkeScaleWorkerCreate,
		Read:   resourceTencentCloudTkeScaleWorkerRead,
		Delete: resourceTencentCloudTkeScaleWorkerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				importFlag = true
				err := resourceTencentCloudTkeScaleWorkerRead(d, m)
				if err != nil {
					return nil, fmt.Errorf("failed to import resource")
				}

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the cluster.",
			},
			"worker_config": {
				Type:     schema.TypeList,
				ForceNew: true,
				MaxItems: 1,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: TkeCvmCreateInfo(),
				},
				Description: "Deploy the machine configuration information of the 'WORK' service, and create <=20 units for common users.",
			},
			//advanced instance settings
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Labels of kubernetes scale worker created nodes.",
			},
			"extra_args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom parameter information related to the node.",
			},
			"gpu_args": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: TKEGpuArgsSetting(),
				},
				Description: "GPU driver parameters.",
			},
			"unschedulable": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "Set whether the added node participates in scheduling. The default value is 0, which means participating in scheduling; non-0 means not participating in scheduling. After the node initialization is completed, you can execute kubectl uncordon nodename to join the node in scheduling.",
			},
			"desired_pod_num": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Description: "Indicate to set desired pod number in current node. Valid when the cluster enable customized pod cidr.",
			},
			"docker_graph_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Docker graph path. Default is `/var/lib/docker`.",
			},
			"mount_target": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Mount target. Default is not mounting.",
			},
			"data_disk": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				MaxItems:    11,
				Description: "Configurations of data disk.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							ForceNew:     true,
							Optional:     true,
							Default:      svcas.SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
							Description:  "Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD` and `CLOUD_HSSD` and `CLOUD_TSSD`.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							ForceNew:    true,
							Optional:    true,
							Default:     0,
							Description: "Volume of disk in GB. Default is `0`.",
						},
						"file_system": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Optional:    true,
							Default:     "",
							Description: "File system, e.g. `ext3/ext4/xfs`.",
						},
						"auto_format_and_mount": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     false,
							Description: "Indicate whether to auto format and mount or not. Default is `false`.",
						},
						"mount_target": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Default:     "",
							Description: "Mount target.",
						},
					},
				},
			},
			"pre_start_user_script": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Base64-encoded user script, executed before initializing the node, currently only effective for adding existing nodes.",
			},
			"user_script": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Base64 encoded user script, this script will be executed after the k8s component is run. The user needs to ensure that the script is reentrant and retry logic. The script and its generated log files can be viewed in the /data/ccs_userscript/ path of the node, if required. The node needs to be initialized before it can be added to the schedule. It can be used with the unschedulable parameter. After the final initialization of userScript is completed, add the kubectl uncordon nodename --kubeconfig=/root/.kube/config command to add the node to the schedule.",
			},
			// Computed values
			"worker_instances_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: tkeCvmState(),
				},
				Description: "An information list of kubernetes cluster 'WORKER'. Each element contains the following attributes:",
			},
		},
	}
}

func resourceTencentCloudTkeScaleWorkerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_scale_worker.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var cvms RunInstancesForNode
	var iAdvanced tke.InstanceAdvancedSettings
	cvms.Work = []string{}

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Get("cluster_id").(string)
	if clusterId == "" {
		return fmt.Errorf("`cluster_id` is empty.")
	}

	info, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, err = service.DescribeCluster(ctx, clusterId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if !has {
		return fmt.Errorf("cluster [%s] is not exist.", clusterId)
	}

	dMap := make(map[string]interface{}, 5)
	//mount_target, docker_graph_path, data_disk, extra_args, desired_pod_num
	iAdvancedParas := []string{"mount_target", "docker_graph_path", "extra_args", "data_disk", "desired_pod_num", "gpu_args"}
	for _, k := range iAdvancedParas {
		if v, ok := d.GetOk(k); ok {
			dMap[k] = v
		}
	}
	iAdvanced = tkeGetInstanceAdvancedPara(dMap, meta)

	iAdvanced.Labels = GetTkeLabels(d, "labels")
	if temp, ok := d.GetOk("unschedulable"); ok {
		iAdvanced.Unschedulable = helper.Int64(int64(temp.(int)))
	}

	if v, ok := d.GetOk("pre_start_user_script"); ok {
		iAdvanced.PreStartUserScript = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_script"); ok {
		iAdvanced.UserScript = helper.String(v.(string))
	}

	if workers, ok := d.GetOk("worker_config"); ok {
		workerList := workers.([]interface{})
		for index := range workerList {
			worker := workerList[index].(map[string]interface{})
			paraJson, _, err := tkeGetCvmRunInstancesPara(worker, meta, info.VpcId, info.ProjectId)
			if err != nil {
				return err
			}
			cvms.Work = append(cvms.Work, paraJson)
		}
	}
	if len(cvms.Work) != 1 {
		return fmt.Errorf("only one additional configuration of virtual machines is now supported now, " +
			"so len(cvms.Work) should be 1")
	}

	instanceIds, err := service.CreateClusterInstances(ctx, clusterId, cvms.Work[0], iAdvanced)
	if err != nil {
		return err
	}

	workerInstancesList := make([]map[string]interface{}, 0, len(instanceIds))
	for _, v := range instanceIds {
		if v == "" {
			return fmt.Errorf("CreateClusterInstances return one instanceId is empty")
		}
		infoMap := make(map[string]interface{})
		infoMap["instance_id"] = v
		infoMap["instance_role"] = TKE_ROLE_WORKER
		workerInstancesList = append(workerInstancesList, infoMap)
	}

	if err = d.Set("worker_instances_list", workerInstancesList); err != nil {
		return err
	}

	//修改id设置,不符合id规则
	id := clusterId + tccommon.FILED_SP + strings.Join(instanceIds, tccommon.FILED_SP)
	d.SetId(id)

	//wait for LANIP
	time.Sleep(tccommon.ReadRetryTimeout)
	return resourceTencentCloudTkeScaleWorkerRead(d, meta)
}

func resourceTencentCloudTkeScaleWorkerRead(d *schema.ResourceData, meta interface{}) error {

	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_scale_worker.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, tccommon.GetLogId(tccommon.ContextNil))
	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var (
		items                  = strings.Split(d.Id(), tccommon.FILED_SP)
		oldWorkerInstancesList = d.Get("worker_instances_list").([]interface{})
		instanceMap            = make(map[string]bool)
		clusterId              = ""
	)

	if importFlag {
		clusterId = items[0]
		if len(items[1:]) >= 2 {
			return fmt.Errorf("only one additional configuration of virtual machines is now supported now, " +
				"so should be 1")
		}
		infoMap := map[string]interface{}{
			"instance_id": items[1],
		}
		oldWorkerInstancesList = append(oldWorkerInstancesList, infoMap)
	} else {
		clusterId = d.Get("cluster_id").(string)
	}

	if clusterId == "" {
		return fmt.Errorf("tke.`cluster_id` is empty.")
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, err := service.DescribeCluster(ctx, clusterId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if !has {
			d.SetId("")
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, v := range oldWorkerInstancesList {
		infoMap, ok := v.(map[string]interface{})
		if !ok || infoMap["instance_id"] == nil {
			return fmt.Errorf("worker_instances_list is broken.")
		}
		instanceId, ok := infoMap["instance_id"].(string)
		if !ok || instanceId == "" {
			return fmt.Errorf("worker_instances_list.instance_id is broken.")
		}
		if instanceMap[instanceId] {
			continue
		}
		instanceMap[instanceId] = true
	}

	_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, workers, err = service.DescribeClusterInstances(ctx, clusterId)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	newWorkerInstancesList := make([]map[string]interface{}, 0, len(workers))
	labelsMap := make(map[string]string)
	instanceIds := make([]*string, 0)
	for sub, cvm := range workers {
		if _, ok := instanceMap[cvm.InstanceId]; !ok {
			continue
		}
		instanceIds = append(instanceIds, &workers[sub].InstanceId)
		tempMap := make(map[string]interface{})
		tempMap["instance_id"] = cvm.InstanceId
		tempMap["instance_role"] = cvm.InstanceRole
		tempMap["instance_state"] = cvm.InstanceState
		tempMap["failed_reason"] = cvm.FailedReason
		tempMap["lan_ip"] = cvm.LanIp

		newWorkerInstancesList = append(newWorkerInstancesList, tempMap)
		if cvm.InstanceAdvancedSettings != nil {
			if cvm.InstanceAdvancedSettings.Labels != nil {
				for _, v := range cvm.InstanceAdvancedSettings.Labels {
					labelsMap[helper.PString(v.Name)] = helper.PString(v.Value)
				}
			}

			_ = d.Set("unschedulable", helper.PInt64(cvm.InstanceAdvancedSettings.Unschedulable))
			_ = d.Set("pre_start_user_script", helper.PString(cvm.InstanceAdvancedSettings.PreStartUserScript))
			_ = d.Set("user_script", helper.PString(cvm.InstanceAdvancedSettings.UserScript))

			if importFlag {
				_ = d.Set("docker_graph_path", helper.PString(cvm.InstanceAdvancedSettings.DockerGraphPath))
				_ = d.Set("desired_pod_num", helper.PInt64(cvm.InstanceAdvancedSettings.DesiredPodNumber))
				_ = d.Set("mount_target", helper.PString(cvm.InstanceAdvancedSettings.MountTarget))
			}

			if cvm.InstanceAdvancedSettings.DataDisks != nil && len(cvm.InstanceAdvancedSettings.DataDisks) > 0 {
				dataDisks := make([]interface{}, 0, len(cvm.InstanceAdvancedSettings.DataDisks))
				for i := range cvm.InstanceAdvancedSettings.DataDisks {
					item := cvm.InstanceAdvancedSettings.DataDisks[i]
					disk := make(map[string]interface{})
					disk["disk_type"] = helper.PString(item.DiskType)
					disk["disk_size"] = helper.PInt64(item.DiskSize)
					disk["file_system"] = helper.PString(item.FileSystem)
					disk["auto_format_and_mount"] = helper.PBool(item.AutoFormatAndMount)
					disk["mount_target"] = helper.PString(item.MountTarget)
					disk["disk_partition"] = helper.PString(item.MountTarget)
					dataDisks = append(dataDisks, disk)
				}
				if importFlag {
					_ = d.Set("data_disk", dataDisks)
				}
			}

			if cvm.InstanceAdvancedSettings.GPUArgs != nil {
				setting := cvm.InstanceAdvancedSettings.GPUArgs

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

				if importFlag {
					if driverEmptyFlag || cudaEmptyFlag || cudnnEmptyFlag || customDriverEmptyFlag {
						_ = d.Set("gpu_args", []interface{}{gpuArgs})
					}

				}
			}
		}
	}

	//worker_config
	var instances []*cvm.Instance
	var errRet error
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instances, errRet = cvmService.DescribeInstanceByFilter(ctx, instanceIds, nil)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	instanceList := make([]interface{}, 0, len(instances))
	for _, instance := range instances {
		mapping := map[string]interface{}{
			"count":                               1,
			"instance_charge_type_prepaid_period": 1,
			"instance_type":                       helper.PString(instance.InstanceType),
			"subnet_id":                           helper.PString(instance.VirtualPrivateCloud.SubnetId),
			"availability_zone":                   helper.PString(instance.Placement.Zone),
			"instance_name":                       helper.PString(instance.InstanceName),
			"instance_charge_type":                helper.PString(instance.InstanceChargeType),
			"system_disk_type":                    helper.PString(instance.SystemDisk.DiskType),
			"system_disk_size":                    helper.PInt64(instance.SystemDisk.DiskSize),
			"internet_charge_type":                helper.PString(instance.InternetAccessible.InternetChargeType),
			"bandwidth_package_id":                helper.PString(instance.InternetAccessible.BandwidthPackageId),
			"internet_max_bandwidth_out":          helper.PInt64(instance.InternetAccessible.InternetMaxBandwidthOut),
			"security_group_ids":                  helper.StringsInterfaces(instance.SecurityGroupIds),
			"img_id":                              helper.PString(instance.ImageId),
		}

		if instance.RenewFlag != nil && helper.PString(instance.InstanceChargeType) == "PREPAID" {
			mapping["instance_charge_type_prepaid_renew_flag"] = helper.PString(instance.RenewFlag)
		} else {
			mapping["instance_charge_type_prepaid_renew_flag"] = ""
		}
		if helper.PInt64(instance.InternetAccessible.InternetMaxBandwidthOut) > 0 {
			mapping["public_ip_assigned"] = true
		}

		if instance.CamRoleName != nil {
			mapping["cam_role_name"] = instance.CamRoleName
		}
		if instance.LoginSettings != nil {
			if instance.LoginSettings.KeyIds != nil && len(instance.LoginSettings.KeyIds) > 0 {
				mapping["key_ids"] = helper.StringsInterfaces(instance.LoginSettings.KeyIds)
			}
			if instance.LoginSettings.Password != nil {
				mapping["password"] = helper.PString(instance.LoginSettings.Password)
			}
		}
		if instance.DisasterRecoverGroupId != nil && helper.PString(instance.DisasterRecoverGroupId) != "" {
			mapping["disaster_recover_group_ids"] = []string{helper.PString(instance.DisasterRecoverGroupId)}
		}
		if instance.HpcClusterId != nil {
			mapping["hpc_cluster_id"] = helper.PString(instance.HpcClusterId)
		}

		dataDisks := make([]interface{}, 0, len(instance.DataDisks))
		for _, v := range instance.DataDisks {
			dataDisk := map[string]interface{}{
				"disk_type":   helper.PString(v.DiskType),
				"disk_size":   helper.PInt64(v.DiskSize),
				"snapshot_id": helper.PString(v.DiskId),
				"encrypt":     helper.PBool(v.Encrypt),
				"kms_key_id":  helper.PString(v.KmsKeyId),
			}
			dataDisks = append(dataDisks, dataDisk)
		}

		mapping["data_disk"] = dataDisks
		instanceList = append(instanceList, mapping)
	}
	if importFlag {
		_ = d.Set("worker_config", instanceList)
	}

	// The machines I generated was deleted by others.
	if len(newWorkerInstancesList) == 0 {
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("labels", labelsMap)
	_ = d.Set("worker_instances_list", newWorkerInstancesList)

	return nil
}
func resourceTencentCloudTkeScaleWorkerDelete(d *schema.ResourceData, meta interface{}) error {

	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_scale_worker.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Get("cluster_id").(string)

	if clusterId == "" {
		return fmt.Errorf("`cluster_id` is empty.")
	}

	_, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = service.DescribeCluster(ctx, clusterId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}
	// The cluster has been deleted
	if !has {
		return nil
	}
	workerInstancesList := d.Get("worker_instances_list").([]interface{})

	instanceMap := make(map[string]bool)

	for _, v := range workerInstancesList {

		infoMap, ok := v.(map[string]interface{})

		if !ok || infoMap["instance_id"] == nil {
			return fmt.Errorf("worker_instances_list is broken.")
		}
		instanceId, ok := infoMap["instance_id"].(string)
		if !ok || instanceId == "" {
			return fmt.Errorf("worker_instances_list.instance_id is broken.")
		}

		if instanceMap[instanceId] {
			log.Printf("[WARN]The same instance id exists in the list")
		}

		instanceMap[instanceId] = true

	}

	_, workers, err := service.DescribeClusterInstances(ctx, clusterId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, workers, err = service.DescribeClusterInstances(ctx, clusterId)

			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}

			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	needDeletes := []string{}
	for _, cvm := range workers {
		if _, ok := instanceMap[cvm.InstanceId]; ok {
			needDeletes = append(needDeletes, cvm.InstanceId)
		}
	}
	// The machines I generated was deleted by others.
	if len(needDeletes) == 0 {
		return nil
	}

	err = service.DeleteClusterInstances(ctx, clusterId, needDeletes)
	if err != nil {
		err = resource.Retry(3*tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err = service.DeleteClusterInstances(ctx, clusterId, needDeletes)

			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}

				if e.GetCode() == "InternalError.Param" &&
					strings.Contains(e.GetMessage(), `PARAM_ERROR[some instances []is not in right state`) {
					return nil
				}
			}

			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})
	}
	return err
}
