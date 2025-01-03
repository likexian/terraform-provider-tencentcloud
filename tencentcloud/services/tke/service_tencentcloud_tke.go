package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tke2 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20220501"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type ClusterBasicSetting struct {
	ClusterId               string
	ClusterOs               string
	ClusterOsType           string
	ClusterVersion          string
	ClusterName             string
	ClusterDescription      string
	ClusterLevel            *string
	AutoUpgradeClusterLevel *bool
	VpcId                   string
	ProjectId               int64
	ClusterNodeNum          int64
	ClusterStatus           string
	SubnetId                string
	Tags                    map[string]string
}

type ClusterAdvancedSettings struct {
	Ipvs                    bool
	AsEnabled               bool
	EnableCustomizedPodCIDR bool
	BasePodNumber           int64
	ContainerRuntime        string
	RuntimeVersion          string
	NodeNameType            string
	ExtraArgs               ClusterExtraArgs
	NetworkType             string
	IsNonStaticIpMode       bool
	DeletionProtection      bool
	KubeProxyMode           string
	Property                string
	OsCustomizeType         string
	VpcCniType              string
}

type ClusterExtraArgs struct {
	KubeAPIServer         []string
	KubeControllerManager []string
	KubeScheduler         []string
}

type RunInstancesForNode struct {
	Master []string
	Work   []string
}

type InstanceAdvancedSettings struct {
	MountTarget        string
	DockerGraphPath    string
	PreStartUserScript string
	UserScript         string
	Unschedulable      int64
	DesiredPodNum      int64
	Labels             []*tke.Label
	DataDisks          []*tke.DataDisk
	ExtraArgs          tke.InstanceExtraArgs
}

type ClusterCidrSettings struct {
	ClusterCidr               string
	IgnoreClusterCidrConflict bool
	MaxNodePodNum             int64
	MaxClusterServiceNum      int64
	ServiceCIDR               string
	EniSubnetIds              []string
	ClaimExpiredSeconds       int64
}

type ClusterInfo struct {
	ClusterBasicSetting
	ClusterCidrSettings
	ClusterAdvancedSettings

	DeployType  string
	CreatedTime string
}

type InstanceInfo struct {
	InstanceId                   string
	InstanceRole                 string
	InstanceState                string
	FailedReason                 string
	NodePoolId                   string
	CreatedTime                  string
	InstanceAdvancedSettings     *tke.InstanceAdvancedSettings
	InstanceDataDiskMountSetting *tke.InstanceDataDiskMountSetting
	LanIp                        string
}

type PrometheusConfigIds struct {
	InstanceId  string
	ClusterType string
	ClusterId   string
}

type Switch struct {
	ClusterId *string     `json:"ClusterId,omitempty" name:"ClusterId"`
	Audit     *SwitchInfo `json:"Audit,omitempty" name:"Audit"`
	Event     *SwitchInfo `json:"Event,omitempty" name:"Event"`
	Log       *SwitchInfo `json:"Log,omitempty" name:"Log"`
}

type SwitchInfo struct {
	Enable      *bool   `json:"Enable,omitempty" name:"Enable"`
	LogsetId    *string `json:"LogsetId,omitempty" name:"LogsetId"`
	TopicId     *string `json:"TopicId,omitempty" name:"TopicId"`
	Version     *string `json:"Version,omitempty" name:"Version"`
	UpgradeAble *bool   `json:"UpgradeAble,omitempty" name:"UpgradeAble"`
}

type DescribeLogSwitchesResponseParams struct {
	SwitchSet []*Switch `json:"SwitchSet,omitempty" name:"SwitchSet"`
	RequestId *string   `json:"RequestId,omitempty" name:"RequestId"`
}

type DescribeLogSwitchesResponse struct {
	tchttp.BaseResponse
	Response *DescribeLogSwitchesResponseParams `json:"Response"`
}

func NewTkeService(client *connectivity.TencentCloudClient) TkeService {
	return TkeService{client: client}
}

type TkeService struct {
	client *connectivity.TencentCloudClient
}

func (me *TkeService) DescribeClusterInstances(ctx context.Context, id string) (masters []InstanceInfo, workers []InstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClusterInstancesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &id
	masters = make([]InstanceInfo, 0, 100)
	workers = make([]InstanceInfo, 0, 100)
	var offset int64 = 0
	var limit int64 = 20
	var has = map[string]bool{}
	var total int64 = -1

getMoreData:
	if total >= 0 && offset >= total {
		return
	}
	request.Limit = &limit
	request.Offset = &offset
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterInstances(request)
	if err != nil {
		errRet = err
		return
	}
	if total < 0 {
		total = int64(*response.Response.TotalCount)
	}

	if len(response.Response.InstanceSet) > 0 {
		offset += limit
	} else {
		// get empty set, we're done
		return
	}

	for _, item := range response.Response.InstanceSet {
		if has[*item.InstanceId] {
			errRet = fmt.Errorf("get repeated instance_id[%s] when doing DescribeClusterInstances", *item.InstanceId)
			return
		}
		has[*item.InstanceId] = true
		instanceInfo := InstanceInfo{
			InstanceId:               *item.InstanceId,
			InstanceRole:             *item.InstanceRole,
			InstanceState:            *item.InstanceState,
			FailedReason:             *item.FailedReason,
			InstanceAdvancedSettings: item.InstanceAdvancedSettings,
		}
		if item.CreatedTime != nil {
			instanceInfo.CreatedTime = *item.CreatedTime
		}
		if item.NodePoolId != nil {
			instanceInfo.NodePoolId = *item.NodePoolId
		}
		if item.LanIP != nil {
			instanceInfo.LanIp = *item.LanIP
		}
		if instanceInfo.InstanceRole == TKE_ROLE_WORKER {
			workers = append(workers, instanceInfo)
		} else {
			masters = append(masters, instanceInfo)
		}
	}
	goto getMoreData

}

func (me *TkeService) DescribeClusterInstancesByRole(ctx context.Context, id, role string) (masters []InstanceInfo, workers []InstanceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClusterInstancesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &id
	request.InstanceRole = &role
	masters = make([]InstanceInfo, 0, 100)
	workers = make([]InstanceInfo, 0, 100)
	var offset int64 = 0
	var limit int64 = 20
	var has = map[string]bool{}
	var total int64 = -1

	for {
		if total >= 0 && offset >= total {
			break
		}
		request.Limit = &limit
		request.Offset = &offset
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTkeClient().DescribeClusterInstances(request)
		if err != nil {
			errRet = err
			return
		}
		if total < 0 {
			total = int64(*response.Response.TotalCount)
		}

		if len(response.Response.InstanceSet) == 0 {
			// get empty set, we're done
			break
		}

		offset += limit

		for _, item := range response.Response.InstanceSet {
			if has[*item.InstanceId] {
				errRet = fmt.Errorf("get repeated instance_id[%s] when doing DescribeClusterInstances", *item.InstanceId)
				return
			}
			has[*item.InstanceId] = true
			instanceInfo := InstanceInfo{
				InstanceId:               *item.InstanceId,
				InstanceRole:             *item.InstanceRole,
				InstanceState:            *item.InstanceState,
				FailedReason:             *item.FailedReason,
				InstanceAdvancedSettings: item.InstanceAdvancedSettings,
			}
			if item.CreatedTime != nil {
				instanceInfo.CreatedTime = *item.CreatedTime
			}
			if item.NodePoolId != nil {
				instanceInfo.NodePoolId = *item.NodePoolId
			}
			if item.LanIP != nil {
				instanceInfo.LanIp = *item.LanIP
			}
			if instanceInfo.InstanceRole == TKE_ROLE_WORKER {
				workers = append(workers, instanceInfo)
			} else {
				masters = append(masters, instanceInfo)
			}
		}
	}

	return
}

func (me *TkeService) DescribeClusters(ctx context.Context, id string, name string) (clusterInfos []ClusterInfo, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClustersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if id != "" && name != "" {
		errRet = fmt.Errorf("cluster_id, cluster_name only one can be set one")
		return
	}

	if id != "" {
		request.ClusterIds = []*string{&id}
	}

	if name != "" {
		filter := &tke.Filter{
			Name:   helper.String("ClusterName"),
			Values: []*string{&name},
		}
		request.Filters = []*tke.Filter{filter}
	}

	response, err := me.client.UseTkeClient().DescribeClusters(request)

	if err != nil {
		errRet = err
		return
	}

	lenClusters := len(response.Response.Clusters)

	if lenClusters == 0 {
		return
	}
	clusterInfos = make([]ClusterInfo, 0, lenClusters)

	for index := range response.Response.Clusters {
		cluster := response.Response.Clusters[index]
		var clusterInfo ClusterInfo

		clusterInfo.ClusterId = *cluster.ClusterId
		clusterInfo.ClusterOs = *cluster.ClusterOs
		clusterInfo.ClusterVersion = *cluster.ClusterVersion
		clusterInfo.ClusterDescription = *cluster.ClusterDescription
		clusterInfo.ClusterName = *cluster.ClusterName
		clusterInfo.ClusterStatus = *cluster.ClusterStatus
		clusterInfo.ClusterLevel = cluster.ClusterLevel
		clusterInfo.AutoUpgradeClusterLevel = cluster.AutoUpgradeClusterLevel

		clusterInfo.ProjectId = int64(*cluster.ProjectId)
		clusterInfo.VpcId = *cluster.ClusterNetworkSettings.VpcId
		clusterInfo.ClusterNodeNum = int64(*cluster.ClusterNodeNum)

		clusterInfo.IgnoreClusterCidrConflict = *cluster.ClusterNetworkSettings.IgnoreClusterCIDRConflict
		clusterInfo.ClusterCidr = *cluster.ClusterNetworkSettings.ClusterCIDR
		clusterInfo.MaxClusterServiceNum = int64(*cluster.ClusterNetworkSettings.MaxClusterServiceNum)
		clusterInfo.EniSubnetIds = common.StringValues(cluster.ClusterNetworkSettings.Subnets)

		clusterInfo.MaxNodePodNum = int64(*cluster.ClusterNetworkSettings.MaxNodePodNum)
		clusterInfo.DeployType = strings.ToUpper(*cluster.ClusterType)
		clusterInfo.Ipvs = *cluster.ClusterNetworkSettings.Ipvs
		clusterInfo.CreatedTime = *cluster.CreatedTime

		projectMap, err := helper.JsonToMap(*cluster.Property)
		if err != nil {
			errRet = err
			return
		}
		if projectMap != nil {
			if projectMap["VpcCniType"] != nil {
				vpcCniType := projectMap["VpcCniType"].(string)
				clusterInfo.VpcCniType = vpcCniType
			}
		}

		if len(cluster.TagSpecification) > 0 {
			clusterInfo.Tags = make(map[string]string)
			for _, tag := range cluster.TagSpecification[0].Tags {
				clusterInfo.Tags[*tag.Key] = *tag.Value
			}
		}

		clusterInfos = append(clusterInfos, clusterInfo)
	}
	return
}

func (me *TkeService) DescribeCluster(ctx context.Context, id string) (
	clusterInfo ClusterInfo,
	has bool,
	errRet error,
) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClustersRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterIds = []*string{&id}

	var iacExtInfo connectivity.IacExtInfo
	iacExtInfo.InstanceId = id
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient(iacExtInfo).DescribeClusters(request)

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.Clusters) == 0 {
		return
	}

	has = true
	cluster := response.Response.Clusters[0]
	clusterInfo.ClusterId = *cluster.ClusterId
	clusterInfo.ClusterOs = *cluster.ClusterOs
	clusterInfo.ClusterVersion = *cluster.ClusterVersion
	clusterInfo.ClusterDescription = *cluster.ClusterDescription
	clusterInfo.ClusterName = *cluster.ClusterName
	clusterInfo.ClusterStatus = *cluster.ClusterStatus
	clusterInfo.ClusterLevel = cluster.ClusterLevel
	clusterInfo.AutoUpgradeClusterLevel = cluster.AutoUpgradeClusterLevel

	clusterInfo.ProjectId = int64(*cluster.ProjectId)
	clusterInfo.VpcId = *cluster.ClusterNetworkSettings.VpcId
	clusterInfo.ClusterNodeNum = int64(*cluster.ClusterNodeNum)

	clusterInfo.DeployType = strings.ToUpper(*cluster.ClusterType)
	clusterInfo.Ipvs = *cluster.ClusterNetworkSettings.Ipvs
	clusterInfo.Property = helper.PString(cluster.Property)
	clusterInfo.OsCustomizeType = helper.PString(cluster.OsCustomizeType)
	clusterInfo.ContainerRuntime = helper.PString(cluster.ContainerRuntime)
	clusterInfo.DeletionProtection = helper.PBool(cluster.DeletionProtection)
	clusterInfo.RuntimeVersion = helper.PString(cluster.RuntimeVersion)
	if cluster.ClusterNetworkSettings != nil {
		clusterInfo.KubeProxyMode = helper.PString(cluster.ClusterNetworkSettings.KubeProxyMode)
		clusterInfo.IgnoreClusterCidrConflict = helper.PBool(cluster.ClusterNetworkSettings.IgnoreClusterCIDRConflict)
		clusterInfo.ClusterCidr = helper.PString(cluster.ClusterNetworkSettings.ClusterCIDR)
		clusterInfo.MaxClusterServiceNum = int64(helper.PUint64(cluster.ClusterNetworkSettings.MaxClusterServiceNum))
		clusterInfo.MaxNodePodNum = int64(helper.PUint64(cluster.ClusterNetworkSettings.MaxNodePodNum))
		clusterInfo.ServiceCIDR = helper.PString(cluster.ClusterNetworkSettings.ServiceCIDR)
	}
	clusterInfo.EniSubnetIds = common.StringValues(cluster.ClusterNetworkSettings.Subnets)

	projectMap, err := helper.JsonToMap(*cluster.Property)
	if err != nil {
		errRet = err
		return
	}
	if projectMap != nil {
		if projectMap["VpcCniType"] != nil {
			vpcCniType := projectMap["VpcCniType"].(string)
			clusterInfo.VpcCniType = vpcCniType
		}
		if projectMap["NetworkType"] != nil {
			networkType := projectMap["NetworkType"].(string)
			clusterInfo.NetworkType = networkType
		}
	}

	if len(cluster.TagSpecification) > 0 {
		clusterInfo.Tags = make(map[string]string)
		for _, tag := range cluster.TagSpecification[0].Tags {
			clusterInfo.Tags[*tag.Key] = *tag.Value
		}
	}

	return
}

func (me *TkeService) DescribeClusterCommonNames(ctx context.Context, request *tke.DescribeClusterCommonNamesRequest) (commonNames []*tke.CommonName, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterCommonNames(request)

	if err != nil {
		errRet = err
		return
	}

	commonNames = response.Response.CommonNames

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DescribeClusterLevelAttribute(ctx context.Context, id string) (clusterLevels []*tke.ClusterLevelAttribute, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClusterLevelAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	if id != "" {
		request.ClusterID = &id
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterLevelAttribute(request)

	if err != nil {
		errRet = err
		return
	}

	clusterLevels = response.Response.Items

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DescribeClusterConfig(ctx context.Context, id string, isPublic bool) (config string, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClusterKubeconfigRequest()
	if isPublic {
		request.IsExtranet = &isPublic
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &id

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterKubeconfig(request)

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		return
	}

	config = *response.Response.Kubeconfig
	return
}

func (me *TkeService) GetUpgradeInstanceResult(ctx context.Context, id string) (
	done bool,
	errRet error,
) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewGetUpgradeInstanceProgressRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &id

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().GetUpgradeInstanceProgress(request)

	if err != nil {
		errRet = err
		return
	}

	lifeState := *response.Response.LifeState

	// all instances success, lifeState=done
	if lifeState == "done" {
		return true, nil
	} else if lifeState != "process" {
		return false, fmt.Errorf("upgrade instances failed, tke response lifeState is:%s", lifeState)
	}

	// parent lifeState=process, check whether all instances in processing.
	for _, inst := range response.Response.Instances {
		if *inst.LifeState == "done" || *inst.LifeState == "pending" {
			continue
		}
		if *inst.LifeState != "process" {
			return false, fmt.Errorf("upgrade instances failed, "+
				"instanceId:%s, lifeState is:%s", *inst.InstanceID, *inst.LifeState)
		}
		// instance lifeState=process, check whether failed or not.
		for _, detail := range inst.Detail {
			if *detail.LifeState == "failed" {
				return false, fmt.Errorf("upgrade instances failed, "+
					"instanceId:%s, detail.lifeState is:%s", *inst.InstanceID, *detail.LifeState)
			}
		}
	}

	return
}

func (me *TkeService) CreateCluster(ctx context.Context,
	basic ClusterBasicSetting,
	advanced ClusterAdvancedSettings,
	cvms RunInstancesForNode,
	iAdvanced InstanceAdvancedSettings,
	cidrSetting ClusterCidrSettings,
	tags map[string]string,
	existedInstance []*tke.ExistedInstancesForNode,
	overrideSettings *OverrideSettings,
	iDiskMountSettings []*tke.InstanceDataDiskMountSetting,
	extensionAddons []*tke.ExtensionAddon,
) (id string, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewCreateClusterRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterBasicSettings = &tke.ClusterBasicSettings{}
	request.ClusterBasicSettings.ClusterOs = &basic.ClusterOs
	request.ClusterBasicSettings.ClusterVersion = &basic.ClusterVersion
	request.ClusterBasicSettings.ProjectId = &basic.ProjectId
	request.ClusterBasicSettings.VpcId = &basic.VpcId
	request.ClusterBasicSettings.ClusterDescription = &basic.ClusterDescription
	request.ClusterBasicSettings.ClusterName = &basic.ClusterName
	request.ClusterBasicSettings.OsCustomizeType = &basic.ClusterOsType
	request.ClusterBasicSettings.ClusterLevel = basic.ClusterLevel
	if basic.AutoUpgradeClusterLevel != nil {
		request.ClusterBasicSettings.AutoUpgradeClusterLevel = &tke.AutoUpgradeClusterLevel{
			IsAutoUpgrade: basic.AutoUpgradeClusterLevel,
		}
	}
	if basic.SubnetId != "" {
		request.ClusterBasicSettings.SubnetId = &basic.SubnetId
	}
	for k, v := range tags {
		if len(request.ClusterBasicSettings.TagSpecification) == 0 {
			request.ClusterBasicSettings.TagSpecification = []*tke.TagSpecification{{
				ResourceType: helper.String("cluster"),
			}}
		}

		request.ClusterBasicSettings.TagSpecification[0].Tags = append(request.ClusterBasicSettings.TagSpecification[0].Tags, &tke.Tag{
			Key:   helper.String(k),
			Value: helper.String(v),
		})
	}

	request.ClusterAdvancedSettings = &tke.ClusterAdvancedSettings{}
	request.ClusterAdvancedSettings.IPVS = &advanced.Ipvs
	request.ClusterAdvancedSettings.AsEnabled = &advanced.AsEnabled
	request.ClusterAdvancedSettings.ContainerRuntime = &advanced.ContainerRuntime
	request.ClusterAdvancedSettings.RuntimeVersion = &advanced.RuntimeVersion
	request.ClusterAdvancedSettings.NodeNameType = &advanced.NodeNameType
	request.ClusterAdvancedSettings.EnableCustomizedPodCIDR = &advanced.EnableCustomizedPodCIDR
	request.ClusterAdvancedSettings.BasePodNumber = &advanced.BasePodNumber
	request.ClusterAdvancedSettings.ExtraArgs = &tke.ClusterExtraArgs{
		KubeAPIServer:         common.StringPtrs(advanced.ExtraArgs.KubeAPIServer),
		KubeControllerManager: common.StringPtrs(advanced.ExtraArgs.KubeControllerManager),
		KubeScheduler:         common.StringPtrs(advanced.ExtraArgs.KubeScheduler),
	}
	request.ClusterAdvancedSettings.NetworkType = &advanced.NetworkType
	request.ClusterAdvancedSettings.IsNonStaticIpMode = &advanced.IsNonStaticIpMode
	request.ClusterAdvancedSettings.DeletionProtection = &advanced.DeletionProtection
	request.ClusterAdvancedSettings.KubeProxyMode = &advanced.KubeProxyMode
	request.ClusterAdvancedSettings.VpcCniType = &advanced.VpcCniType

	request.InstanceAdvancedSettings = &tke.InstanceAdvancedSettings{}
	request.InstanceAdvancedSettings.MountTarget = &iAdvanced.MountTarget
	request.InstanceAdvancedSettings.DockerGraphPath = &iAdvanced.DockerGraphPath
	request.InstanceAdvancedSettings.UserScript = &iAdvanced.UserScript
	request.InstanceAdvancedSettings.Unschedulable = &iAdvanced.Unschedulable
	request.InstanceAdvancedSettings.DesiredPodNumber = &iAdvanced.DesiredPodNum
	if len(iAdvanced.ExtraArgs.Kubelet) > 0 {
		request.InstanceAdvancedSettings.ExtraArgs = &iAdvanced.ExtraArgs
	}

	if len(iAdvanced.Labels) > 0 {
		request.InstanceAdvancedSettings.Labels = iAdvanced.Labels
	}

	if len(iAdvanced.DataDisks) > 0 {
		request.InstanceAdvancedSettings.DataDisks = iAdvanced.DataDisks
	}

	if len(extensionAddons) > 0 {
		request.ExtensionAddons = extensionAddons
	}

	if overrideSettings != nil {
		if len(overrideSettings.Master)+len(overrideSettings.Work) > 0 &&
			len(overrideSettings.Master)+len(overrideSettings.Work) != (len(cvms.Master)+len(cvms.Work)) {
			return "", fmt.Errorf("len(overrideSettings) != (len(cvms.Master)+len(cvms.Work))")
		}
	}

	request.RunInstancesForNode = []*tke.RunInstancesForNode{}

	if len(cvms.Master) != 0 {

		var node tke.RunInstancesForNode
		node.NodeRole = helper.String(TKE_ROLE_MASTER_ETCD)
		node.RunInstancesPara = []*string{}
		request.ClusterType = helper.String(TKE_DEPLOY_TYPE_INDEPENDENT)
		for v := range cvms.Master {
			node.RunInstancesPara = append(node.RunInstancesPara, &cvms.Master[v])
			if overrideSettings != nil && len(overrideSettings.Master) != 0 {
				node.InstanceAdvancedSettingsOverrides = append(node.InstanceAdvancedSettingsOverrides, &overrideSettings.Master[v])
			}
		}
		request.RunInstancesForNode = append(request.RunInstancesForNode, &node)

	} else {
		request.ClusterType = helper.String(TKE_DEPLOY_TYPE_MANAGED)
	}

	if len(cvms.Work) != 0 {
		var node tke.RunInstancesForNode
		node.NodeRole = helper.String(TKE_ROLE_WORKER)
		node.RunInstancesPara = []*string{}
		for v := range cvms.Work {
			node.RunInstancesPara = append(node.RunInstancesPara, &cvms.Work[v])
			if overrideSettings != nil && len(overrideSettings.Work) != 0 {
				node.InstanceAdvancedSettingsOverrides = append(node.InstanceAdvancedSettingsOverrides, &overrideSettings.Work[v])
			}
		}
		request.RunInstancesForNode = append(request.RunInstancesForNode, &node)
	}

	if len(iDiskMountSettings) != 0 {
		request.InstanceDataDiskMountSettings = iDiskMountSettings
	}

	request.ClusterCIDRSettings = &tke.ClusterCIDRSettings{}

	maxNodePodNum := uint64(cidrSetting.MaxNodePodNum)
	request.ClusterCIDRSettings.MaxNodePodNum = &maxNodePodNum

	maxClusterServiceNum := uint64(cidrSetting.MaxClusterServiceNum)
	request.ClusterCIDRSettings.MaxClusterServiceNum = &maxClusterServiceNum
	request.ClusterCIDRSettings.ClusterCIDR = &cidrSetting.ClusterCidr
	request.ClusterCIDRSettings.IgnoreClusterCIDRConflict = &cidrSetting.IgnoreClusterCidrConflict
	request.ClusterCIDRSettings.ServiceCIDR = &cidrSetting.ServiceCIDR
	request.ClusterCIDRSettings.EniSubnetIds = common.StringPtrs(cidrSetting.EniSubnetIds)
	request.ClusterCIDRSettings.ClaimExpiredSeconds = &cidrSetting.ClaimExpiredSeconds

	if len(existedInstance) > 0 {
		request.ExistedInstancesForNode = existedInstance
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().CreateCluster(request)

	if err != nil {
		errRet = err
		return
	}

	id = *response.Response.ClusterId
	return
}

func (me *TkeService) CreateClusterInstances(ctx context.Context,
	id string, runInstancePara string,
	iAdvanced tke.InstanceAdvancedSettings) (instanceIds []string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewCreateClusterInstancesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.RunInstancePara = &runInstancePara

	request.InstanceAdvancedSettings = &iAdvanced

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().CreateClusterInstances(request)

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateClusterInstances return nil response")
		return
	}

	instanceIds = make([]string, 0, len(response.Response.InstanceIdSet))

	for _, v := range response.Response.InstanceIdSet {

		instanceIds = append(instanceIds, *v)
	}
	return
}

func (me *TkeService) CheckOneOfClusterNodeReady(ctx context.Context, clusterId string, mustHaveWorkers bool) error {
	logId := tccommon.GetLogId(ctx)
	return resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		_, workers, err := me.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		// check serverless node
		virtualNodes, err := me.DescribeClusterVirtualNode(ctx, clusterId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if len(workers) == 0 && len(virtualNodes) == 0 {
			if mustHaveWorkers {
				log.Printf("[WARN]%s waiting for workers in the cluster[%s] to be created.\n",
					logId, clusterId)
				return resource.RetryableError(fmt.Errorf("waiting for workers created"))
			}
			return nil
		}

		for i := range workers {
			worker := workers[i]
			if worker.InstanceState == "running" {
				return nil
			}
		}
		for i := range virtualNodes {
			virtualNode := virtualNodes[i]
			if virtualNode.Phase != nil && *virtualNode.Phase == "Running" {
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("cluster %s waiting for one of the workers ready", clusterId))
	})
}

/*
if cluster is creating, return error:TencentCloudSDKError] Code=tccommon.InternalError.ClusterState
*/
func (me *TkeService) DeleteClusterInstances(ctx context.Context, id string, instanceIds []string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDeleteClusterInstancesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.InstanceIds = make([]*string, 0, len(instanceIds))

	for index := range instanceIds {
		request.InstanceIds = append(request.InstanceIds, &instanceIds[index])
	}

	request.InstanceDeleteMode = helper.String("terminate")
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().DeleteClusterInstances(request)
	return err
}

func (me *TkeService) DeleteCluster(ctx context.Context, id string) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewDeleteClusterRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.InstanceDeleteMode = helper.String("terminate")

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().DeleteCluster(request)

	return err
}

func (me *TkeService) DescribeClusterSecurity(ctx context.Context, id string) (ret *tke.DescribeClusterSecurityResponse, errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClusterSecurityRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.ClusterId = &id

	return me.client.UseTkeClient().DescribeClusterSecurity(request)
}

func (me *TkeService) CreateClusterAsGroup(ctx context.Context, id, groupPara, configPara string, labels []*tke.Label, iAdvanced InstanceAdvancedSettings) (asGroupId string, errRet error) {
	return "", fmt.Errorf("Cluster AS Group has OFFLINE")
}

func (me *TkeService) DescribeClusterAsGroupsByGroupId(ctx context.Context, id string, groupId string) (clusterAsGroupSet *tke.ClusterAsGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClusterAsGroupsRequest()

	request.ClusterId = &id
	request.AutoScalingGroupIds = []*string{&groupId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterAsGroups(request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), err.Error())
		errRet = err
		return
	}

	if len(response.Response.ClusterAsGroupSet) > 0 {
		clusterAsGroupSet = response.Response.ClusterAsGroupSet[0]
	}
	return
}

func (me *TkeService) DeleteClusterAsGroups(ctx context.Context, id, asGroupId string) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewDeleteClusterAsGroupsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.AutoScalingGroupIds = []*string{&asGroupId}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().DeleteClusterAsGroups(request)
	if err != nil {
		errRet = err
	}
	return
}

/*
open internet access
*/
func (me *TkeService) CreateClusterEndpoint(ctx context.Context, id string, subnetId, securityGroupId string, internet bool, domain string, extensiveParameters string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewCreateClusterEndpointRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id

	request.IsExtranet = &internet

	if subnetId != "" {
		request.SubnetId = &subnetId
	}

	if securityGroupId != "" && internet {
		request.SecurityGroup = &securityGroupId
	}

	if domain != "" {
		request.Domain = helper.String(domain)
	}

	if extensiveParameters != "" {
		request.ExtensiveParameters = helper.String(extensiveParameters)
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().CreateClusterEndpoint(request)
	if err != nil {
		errRet = err
	}
	return
}

func (me *TkeService) DescribeClusterEndpointStatus(ctx context.Context, id string, isExtranet bool) (status string, message string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterEndpointStatusRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.IsExtranet = &isExtranet

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().DescribeClusterEndpointStatus(request)
	if err != nil {
		errRet = err
		return
	}
	if response.Response == nil || response.Response.Status == nil {
		errRet = fmt.Errorf("sdk DescribeClusterEndpointStatus return empty status")
		return
	}
	status = *response.Response.Status
	message = status
	return
}
func (me *TkeService) DescribeClusterEndpoints(ctx context.Context, id string) (response tke.DescribeClusterEndpointsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterEndpointsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	ratelimit.Check(request.GetAction())

	resp, err := me.client.UseTkeClient().DescribeClusterEndpoints(request)
	if err != nil {
		errRet = err
		return
	}
	if resp.Response == nil {
		errRet = fmt.Errorf("sdk DescribeClusterEndpoints return empty")
		return
	}

	return *resp.Response, errRet
}

func (me *TkeService) DeleteClusterEndpoint(ctx context.Context, id string, isInternet bool) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDeleteClusterEndpointRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.IsExtranet = &isInternet

	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().DeleteClusterEndpoint(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) ModifyClusterEndpointSP(ctx context.Context, id string, securityPolicies []string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewModifyClusterEndpointSPRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.SecurityPolicies = helper.Strings(securityPolicies)

	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().ModifyClusterEndpointSP(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) ModifyClusterEndpointSG(ctx context.Context, id string, securityGroup string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewModifyClusterEndpointSPRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.SecurityGroup = &securityGroup

	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().ModifyClusterEndpointSP(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) ModifyClusterAttribute(ctx context.Context, id string, projectId int64, clusterName, clusterDesc, clusterLevel string, autoUpgradeClusterLevel bool) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewModifyClusterAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.ProjectId = &projectId
	request.ClusterName = &clusterName
	request.ClusterDesc = &clusterDesc

	if clusterLevel != "" {
		request.ClusterLevel = &clusterLevel
	}

	request.AutoUpgradeClusterLevel = &tke.AutoUpgradeClusterLevel{
		IsAutoUpgrade: &autoUpgradeClusterLevel,
	}

	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().ModifyClusterAttribute(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) EnableVpcCniNetworkType(ctx context.Context, id string, vpcCniType string, enableStaticIp bool, subnets []string, expiredSeconds uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewEnableVpcCniNetworkTypeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.ClusterId = &id
	request.VpcCniType = &vpcCniType
	request.EnableStaticIp = &enableStaticIp
	request.Subnets = common.StringPtrs(subnets)
	request.ExpiredSeconds = &expiredSeconds

	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().EnableVpcCniNetworkType(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) DescribeEnableVpcCniProgress(ctx context.Context, id string) (status, message string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeEnableVpcCniProgressRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().DescribeEnableVpcCniProgress(request)
	if err != nil {
		errRet = err
		return
	}

	status = *response.Response.Status
	message = *response.Response.ErrorMessage
	return
}

func (me *TkeService) DisableVpcCniNetworkType(ctx context.Context, id string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDisableVpcCniNetworkTypeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id

	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().DisableVpcCniNetworkType(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) AddVpcCniSubnets(ctx context.Context, id string, subnets []string, vpcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewAddVpcCniSubnetsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.SubnetIds = common.StringPtrs(subnets)
	request.VpcId = &vpcId

	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().AddVpcCniSubnets(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) ModifyClusterVersion(ctx context.Context, id string, clusterVersion string, extraArgs interface{}) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewUpdateClusterVersionRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()

	if extraArgs != nil && len(extraArgs.([]interface{})) > 0 {
		// the first elem is in use
		extraInterface := extraArgs.([]interface{})
		extraMap := extraInterface[0].(map[string]interface{})

		kas := make([]*string, 0)
		if kaArgs, exist := extraMap["kube_apiserver"]; exist {
			args := kaArgs.([]interface{})
			for index := range args {
				str := args[index].(string)
				kas = append(kas, &str)
			}
		}
		kcms := make([]*string, 0)
		if kcmArgs, exist := extraMap["kube_controller_manager"]; exist {
			args := kcmArgs.([]interface{})
			for index := range args {
				str := args[index].(string)
				kcms = append(kcms, &str)
			}
		}
		kss := make([]*string, 0)
		if ksArgs, exist := extraMap["kube_scheduler"]; exist {
			args := ksArgs.([]interface{})
			for index := range args {
				str := args[index].(string)
				kss = append(kss, &str)
			}
		}

		request.ExtraArgs = &tke.ClusterExtraArgs{
			KubeAPIServer:         kas,
			KubeControllerManager: kcms,
			KubeScheduler:         kss,
		}
	}

	request.ClusterId = &id
	request.DstVersion = &clusterVersion

	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().UpdateClusterVersion(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) DescribeKubernetesAvailableClusterVersionsByFilter(ctx context.Context, param map[string]interface{}) (ret *tke.DescribeAvailableClusterVersionResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke.NewDescribeAvailableClusterVersionRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "ClusterIds" {
			request.ClusterIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeAvailableClusterVersion(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if err := dataSourceTencentCloudKubernetesAvailableClusterVersionsReadPostRequest0(ctx, request, response); err != nil {
		return nil, err
	}

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *TkeService) CheckClusterVersion(ctx context.Context, id string, clusterVersion string) (isOk bool, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeAvailableClusterVersionRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	ratelimit.Check(request.GetAction())

	resp, err := me.client.UseTkeClient().DescribeAvailableClusterVersion(request)
	if err != nil {
		errRet = err
		return
	}

	if resp == nil || resp.Response == nil || resp.Response.Versions == nil {
		return
	}
	versions := resp.Response.Versions
	for _, v := range versions {
		if *v == clusterVersion {
			isOk = true
			return
		}
	}

	return
}

func (me *TkeService) CheckInstancesUpgradeAble(ctx context.Context, id string, upgradeType string) (instanceIds []string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewCheckInstancesUpgradeAbleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.UpgradeType = &upgradeType
	ratelimit.Check(request.GetAction())

	resp, err := me.client.UseTkeClient().CheckInstancesUpgradeAble(request)
	if err != nil {
		errRet = err
		return
	}

	if resp == nil || resp.Response == nil || resp.Response.UpgradeAbleInstances == nil {
		return
	}
	for _, inst := range resp.Response.UpgradeAbleInstances {
		instanceIds = append(instanceIds, *inst.InstanceId)
	}

	return
}

func (me *TkeService) UpgradeClusterInstances(ctx context.Context, id string, upgradeType string, instanceIds []string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewUpgradeClusterInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	op := "create"
	request.Operation = &op
	request.ClusterId = &id
	request.UpgradeType = &upgradeType
	request.InstanceIds = helper.Strings(instanceIds)
	ratelimit.Check(request.GetAction())

	_, err := me.client.UseTkeClient().UpgradeClusterInstances(request)
	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *TkeService) DescribeImages(ctx context.Context) (imageIds []string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeImagesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeImages(request)
	if err != nil {
		errRet = err
		return
	}

	for _, image := range response.Response.ImageInstanceSet {
		imageIds = append(imageIds, *image.ImageId)
	}
	return
}

func GetTkeLabels(d *schema.ResourceData, k string) []*tke.Label {
	labels := make([]*tke.Label, 0)
	if raw, ok := d.GetOk(k); ok {
		for k, v := range raw.(map[string]interface{}) {
			labels = append(labels, &tke.Label{Name: helper.String(k), Value: helper.String(v.(string))})
		}
	}
	return labels
}

func GetTkeTaints(d *schema.ResourceData, k string) []*tke.Taint {
	taints := make([]*tke.Taint, 0)
	if raw, ok := d.GetOk(k); ok {
		for _, v := range raw.([]interface{}) {
			vv := v.(map[string]interface{})
			taints = append(taints, &tke.Taint{Key: helper.String(vv["key"].(string)), Value: helper.String(vv["value"].(string)), Effect: helper.String(vv["effect"].(string))})
		}
	}
	return taints
}

func GetTkeTags(d *schema.ResourceData, k string) []*tke.Tag {
	tags := make([]*tke.Tag, 0)
	if raw, ok := d.GetOk(k); ok {
		for k, v := range raw.(map[string]interface{}) {
			tags = append(tags, &tke.Tag{Key: helper.String(k), Value: helper.String(v.(string))})
		}
	}
	return tags
}

func (me *TkeService) ModifyClusterAsGroupAttribute(ctx context.Context, id, asGroupId string, maxSize, minSize int64) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewModifyClusterAsGroupAttributeRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.ClusterAsGroupAttribute = &tke.ClusterAsGroupAttribute{
		AutoScalingGroupId: &asGroupId,
		AutoScalingGroupRange: &tke.AutoScalingGroupRange{
			MaxSize: &maxSize,
			MinSize: &minSize,
		},
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().ModifyClusterAsGroupAttribute(request)
	if err != nil {
		errRet = err
	}
	return
}

func (me *TkeService) CreateClusterNodePool(ctx context.Context, clusterId, name, groupPara, configPara string, enableAutoScale bool, nodeOs string, nodeOsType string, labels []*tke.Label, taints []*tke.Taint, iAdvanced tke.InstanceAdvancedSettings, deletionProtection bool, tags []*tke.Tag) (asGroupId string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewCreateClusterNodePoolRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.Name = &name
	request.AutoScalingGroupPara = &groupPara
	request.LaunchConfigurePara = &configPara
	request.InstanceAdvancedSettings = &iAdvanced
	request.EnableAutoscale = &enableAutoScale
	request.DeletionProtection = &deletionProtection
	request.NodePoolOs = &nodeOs
	if nodeOsType != "" {
		request.OsCustomizeType = &nodeOsType
	}

	if len(labels) > 0 {
		request.Labels = labels
	}

	if len(taints) > 0 {
		request.Taints = taints
	}

	if len(tags) > 0 {
		request.Tags = tags
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().CreateClusterNodePool(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.NodePoolId == nil {
		errRet = fmt.Errorf("CreateClusterNodePool return nil response")
		return
	}

	asGroupId = *response.Response.NodePoolId
	return
}

func (me *TkeService) ModifyClusterNodePool(ctx context.Context, clusterId, nodePoolId string, name string, enableAutoScale bool, minSize int64, maxSize int64, nodeOs string, nodeOsType string, labels []*tke.Label, taints []*tke.Taint, tags map[string]string, deletionProtection bool) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewModifyClusterNodePoolRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.NodePoolId = &nodePoolId
	request.Taints = taints
	request.Labels = labels
	request.EnableAutoscale = &enableAutoScale
	request.DeletionProtection = &deletionProtection
	request.MaxNodesNum = &maxSize
	request.MinNodesNum = &minSize
	request.Name = &name
	request.OsName = &nodeOs
	request.OsCustomizeType = &nodeOsType

	if len(labels) > 0 {
		request.Labels = labels
	}

	if len(tags) > 0 {
		for k, v := range tags {
			key := k
			val := v
			request.Tags = append(request.Tags, &tke.Tag{
				Key:   &key,
				Value: &val,
			})
		}
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().ModifyClusterNodePool(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) ModifyClusterNodePoolDesiredCapacity(ctx context.Context, clusterId, nodePoolId string, desiredCapacity int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewModifyNodePoolDesiredCapacityAboutAsgRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.NodePoolId = &nodePoolId
	request.DesiredCapacity = &desiredCapacity

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().ModifyNodePoolDesiredCapacityAboutAsg(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) ModifyClusterNodePoolInstanceTypes(ctx context.Context, clusterId, nodePoolId string, instanceTypes []*string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewModifyNodePoolInstanceTypesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &clusterId
	request.NodePoolId = &nodePoolId
	request.InstanceTypes = instanceTypes

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().ModifyNodePoolInstanceTypes(request)
	if err != nil {
		errRet = err
		return
	}
	return
}

func (me *TkeService) DeleteClusterNodePool(ctx context.Context, id, nodePoolId string, deleteKeepInstance bool) (errRet error) {

	logId := tccommon.GetLogId(ctx)
	request := tke.NewDeleteClusterNodePoolRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &id
	request.NodePoolIds = []*string{&nodePoolId}
	request.KeepInstance = &deleteKeepInstance

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().DeleteClusterNodePool(request)
	if err != nil {
		errRet = err
	}
	return
}

func (me *TkeService) DescribeNodePool(ctx context.Context, clusterId string, nodePoolId string) (
	nodePool *tke.NodePool,
	has bool,
	errRet error,
) {

	logId := tccommon.GetLogId(ctx)
	//the error code of cluster not exist is tccommon.InternalError
	//check cluster exist first
	_, clusterHas, err := me.DescribeCluster(ctx, clusterId)
	if err != nil {
		errRet = err
		return
	}
	if !clusterHas {
		return
	}

	request := tke.NewDescribeClusterNodePoolDetailRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = helper.String(clusterId)
	request.NodePoolId = helper.String(nodePoolId)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterNodePoolDetail(request)

	if err != nil {
		errRet = err
		return
	}

	if response.Response.NodePool == nil {
		return
	}

	has = true
	nodePool = response.Response.NodePool

	return
}

// node pool global config
func (me *TkeService) ModifyClusterNodePoolGlobalConfig(ctx context.Context, request *tke.ModifyClusterAsGroupOptionAttributeRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().ModifyClusterAsGroupOptionAttribute(request)
	if err != nil {
		errRet = err
	}
	return
}

func (me *TkeService) DescribeClusterNodePoolGlobalConfig(ctx context.Context, clusterId string) (
	npGlobalConfig *tke.ClusterAsGroupOption,
	errRet error,
) {

	logId := tccommon.GetLogId(ctx)
	//the error code of cluster not exist is tccommon.InternalError
	//check cluster exist first
	_, clusterHas, err := me.DescribeCluster(ctx, clusterId)
	if err != nil {
		errRet = err
		return
	}
	if !clusterHas {
		return
	}

	request := tke.NewDescribeClusterAsGroupOptionRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = helper.String(clusterId)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterAsGroupOption(request)

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil || response.Response.ClusterAsGroupOption == nil {
		return
	}

	npGlobalConfig = response.Response.ClusterAsGroupOption

	return
}

func (me *TkeService) WaitForAuthenticationOptionsUpdateSuccess(ctx context.Context, id string) (info *tke.ServiceAccountAuthenticationOptions, oidc *tke.OIDCConfigAuthenticationOptions, errRet error) {
	err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		options, state, config, err := me.DescribeClusterAuthenticationOptions(ctx, id)
		info = options
		oidc = config

		if err != nil {
			return resource.NonRetryableError(err)
		}

		if state == "Success" {
			return nil
		}

		if state == "Updating" {
			return resource.RetryableError(fmt.Errorf("state is %s, retry", state))
		}

		return resource.NonRetryableError(fmt.Errorf("update failed: %s", state))
	})

	if err != nil {
		errRet = err
		return
	}
	return
}

// DescribeClusterAuthenticationOptions
// Field `ServiceAccounts.AutoCreateDiscoveryAnonymousAuth` will always return null by design
func (me *TkeService) DescribeClusterAuthenticationOptions(ctx context.Context, id string) (options *tke.ServiceAccountAuthenticationOptions, state string, oidcConfig *tke.OIDCConfigAuthenticationOptions, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClusterAuthenticationOptionsRequest()
	request.ClusterId = helper.String(id)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	res, err := me.client.UseTkeClient().DescribeClusterAuthenticationOptions(request)
	if err != nil {
		errRet = err
	}

	if res.Response != nil {
		state = *res.Response.LatestOperationState
		options = res.Response.ServiceAccounts
		oidcConfig = res.Response.OIDCConfig
	}

	return
}

func (me *TkeService) ModifyClusterAuthenticationOptions(ctx context.Context, request *tke.ModifyClusterAuthenticationOptionsRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().ModifyClusterAuthenticationOptions(request)
	if err != nil {
		errRet = err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) ModifyDeletionProtection(ctx context.Context, id string, enable bool) (errRet error) {
	var (
		logId  = tccommon.GetLogId(ctx)
		action string
	)

	if enable {
		request := tke.NewEnableClusterDeletionProtectionRequest()
		request.ClusterId = &id
		action = request.GetAction()
		ratelimit.Check(action)
		response, err := me.client.UseTkeClient().EnableClusterDeletionProtection(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	} else {
		request := tke.NewDisableClusterDeletionProtectionRequest()
		request.ClusterId = &id
		action = request.GetAction()
		ratelimit.Check(action)
		response, err := me.client.UseTkeClient().DisableClusterDeletionProtection(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId, action, errRet.Error())
		}
	}()

	return
}

func (me *TkeService) AcquireClusterAdminRole(ctx context.Context, clusterId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewAcquireClusterAdminRoleRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &clusterId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().AcquireClusterAdminRole(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DescribeExternalNodeSupportConfig(ctx context.Context, clusterId string) (resp *tke.DescribeExternalNodeSupportConfigResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeExternalNodeSupportConfigRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &clusterId

	ratelimit.Check(request.GetAction())
	res, err := me.client.UseTkeClient().DescribeExternalNodeSupportConfig(request)

	if err != nil {
		errRet = err
		return
	}
	if res == nil || res.Response == nil && res.Response.RequestId == nil {
		return nil, errors.New("invalid response")
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), res.ToJsonString())

	return res.Response, nil
}

func (me *TkeService) DescribeClusterExtraArgs(ctx context.Context, clusterId string) (extraArgs *tke.ClusterExtraArgs, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeClusterExtraArgsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &clusterId

	ratelimit.Check(request.GetAction())
	res, err := me.client.UseTkeClient().DescribeClusterExtraArgs(request)

	if err != nil {
		errRet = err
		return
	}
	if res == nil || res.Response == nil && res.Response.RequestId == nil {
		return nil, errors.New("invalid response")
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), res.ToJsonString())

	return res.Response.ClusterExtraArgs, nil
}

func (me *TkeService) DescribeIPAMD(ctx context.Context, clusterId string) (resp *tke.DescribeIPAMDResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeIPAMDRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterId = &clusterId

	ratelimit.Check(request.GetAction())
	res, err := me.client.UseTkeClient().DescribeIPAMD(request)

	if err != nil {
		errRet = err
		return
	}
	if res == nil || res.Response == nil && res.Response.RequestId == nil {
		return nil, errors.New("invalid response")
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), res.ToJsonString())

	return res.Response, nil

}

func (me *TkeService) DescribeLogSwitches(ctx context.Context, clusterId string) (resp []*tke.Switch, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDescribeLogSwitchesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.ClusterIds = helper.StringsStringsPoint([]string{clusterId})
	request.ClusterType = helper.String("tke")
	ratelimit.Check(request.GetAction())
	res, err := me.client.UseTkeClient().DescribeLogSwitches(request)

	if err != nil {
		errRet = err
		return
	}
	if res == nil || res.Response == nil && res.Response.RequestId == nil {
		return nil, errors.New("invalid response")
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), res.ToJsonString())

	return res.Response.SwitchSet, nil
}

func (me *TkeService) SwitchLogAgent(ctx context.Context, clusterId, rootDir string, enable bool) error {
	if enable {
		request := tke.NewInstallLogAgentRequest()
		request.ClusterId = &clusterId
		if rootDir != "" {
			request.KubeletRootDir = &rootDir
		}
		return me.InstallLogAgent(ctx, request)
	}
	request := tke.NewUninstallLogAgentRequest()
	request.ClusterId = &clusterId
	return me.UninstallLogAgent(ctx, request)
}

func (me *TkeService) SwitchEventPersistence(ctx context.Context, clusterId, logSetId, topicId string,
	enable, deleteEventLog bool) error {
	if enable {
		request := tke.NewEnableEventPersistenceRequest()
		request.ClusterId = &clusterId
		if logSetId != "" {
			request.LogsetId = &logSetId
		}
		if topicId != "" {
			request.TopicId = &topicId
		}
		return me.EnableEventPersistence(ctx, request)
	}

	request := tke.NewDisableEventPersistenceRequest()
	request.ClusterId = &clusterId
	request.DeleteLogSetAndTopic = &deleteEventLog
	return me.DisableEventPersistence(ctx, request)
}

func (me *TkeService) SwitchClusterAudit(ctx context.Context, clusterId, logSetId, topicId string,
	enable, deleteAuditLog bool) error {
	if enable {
		request := tke.NewEnableClusterAuditRequest()
		request.ClusterId = &clusterId
		if logSetId != "" {
			request.LogsetId = &logSetId
		}
		if topicId != "" {
			request.TopicId = &topicId
		}
		return me.EnableClusterAudit(ctx, request)
	}
	request := tke.NewDisableClusterAuditRequest()
	request.ClusterId = &clusterId
	request.DeleteLogSetAndTopic = &deleteAuditLog
	return me.DisableClusterAudit(ctx, request)
}

func (me *TkeService) InstallLogAgent(ctx context.Context, request *tke.InstallLogAgentRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().InstallLogAgent(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) UninstallLogAgent(ctx context.Context, request *tke.UninstallLogAgentRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().UninstallLogAgent(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) EnableEventPersistence(ctx context.Context, request *tke.EnableEventPersistenceRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().EnableEventPersistence(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DisableEventPersistence(ctx context.Context, request *tke.DisableEventPersistenceRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DisableEventPersistence(request)

	if err != nil {
		code := err.(*sdkErrors.TencentCloudSDKError).Code
		if code == "InternalError.KubernetesDeleteOperationError" {
			return
		}
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) EnableClusterAudit(ctx context.Context, request *tke.EnableClusterAuditRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().EnableClusterAudit(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DisableClusterAudit(ctx context.Context, request *tke.DisableClusterAuditRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DisableClusterAudit(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DescribeServerlessNodePoolByClusterIdAndNodePoolId(ctx context.Context, clusterId, nodePoolId string) (instance *tke.VirtualNodePool, has bool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterVirtualNodePoolsRequest()
	request.ClusterId = common.StringPtr(clusterId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterVirtualNodePools(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response != nil && len(response.Response.NodePoolSet) > 0 {
		for _, nodePool := range response.Response.NodePoolSet {
			if nodePool != nil && nodePool.NodePoolId != nil && *nodePool.NodePoolId == nodePoolId {
				has = true
				instance = nodePool
			}
		}
	}
	return
}

func (me *TkeService) CreateClusterVirtualNodePool(ctx context.Context, request *tke.CreateClusterVirtualNodePoolRequest) (id string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().CreateClusterVirtualNodePool(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil && response.Response.NodePoolId != nil {
		id = *response.Response.NodePoolId
	}

	return
}

func (me *TkeService) DeleteClusterVirtualNodePool(ctx context.Context, request *tke.DeleteClusterVirtualNodePoolRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DeleteClusterVirtualNodePool(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) ModifyClusterVirtualNodePool(ctx context.Context, request *tke.ModifyClusterVirtualNodePoolRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().ModifyClusterVirtualNodePool(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DescribeClusterVirtualNode(ctx context.Context, clusterId string) (virtualNodes []tke.VirtualNode, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterVirtualNodeRequest()
	request.ClusterId = common.StringPtr(clusterId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeClusterVirtualNode(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response != nil && len(response.Response.Nodes) > 0 {
		for _, node := range response.Response.Nodes {
			if node != nil {
				virtualNodes = append(virtualNodes, *node)
			}
		}
	}
	return
}

func ModifySecurityServiceOfCvmInNodePool(ctx context.Context, d *schema.ResourceData, tkeSvc *TkeService, cvmSvc *svccvm.CvmService, client *connectivity.TencentCloudClient, clusterId string, nodePoolId string) error {
	logId := tccommon.GetLogId(ctx)

	if d.HasChange("auto_scaling_config.0.enhanced_security_service") {
		workersInsIdOfNodePool := make([]string, 0)
		_, workers, err := tkeSvc.DescribeClusterInstances(ctx, clusterId)
		if err != nil {
			return err
		}
		for _, worker := range workers {
			if worker.NodePoolId != "" && worker.NodePoolId == nodePoolId {
				workersInsIdOfNodePool = append(workersInsIdOfNodePool, worker.InstanceId)
			}
		}

		const BatchProcessedInsLimit = 100 // limit 100 items to change each request
		var (
			launchConfigRaw []interface{}
			dMap            map[string]interface{}
		)
		if raw, ok := d.GetOk("auto_scaling_config"); ok {
			launchConfigRaw = raw.([]interface{})
			dMap = launchConfigRaw[0].(map[string]interface{})
		}

		if v, ok := dMap["enhanced_security_service"]; ok && !v.(bool) {
			// uninstall, cwp/DeleteMachine, need uuid
			// https://cloud.tencent.com/document/product/296/19844
			for i := 0; i < len(workersInsIdOfNodePool); i += BatchProcessedInsLimit {
				var reqInstanceIds []string
				if i+BatchProcessedInsLimit <= len(workersInsIdOfNodePool) {
					reqInstanceIds = workersInsIdOfNodePool[i : i+BatchProcessedInsLimit]
				} else {
					reqInstanceIds = workersInsIdOfNodePool[i:]
				}
				// get uuid
				instanceSet, err := cvmSvc.DescribeInstanceSetByIds(ctx, helper.StrListValToStr(reqInstanceIds))
				if err != nil {
					return err
				}
				// call cwp/DeleteMachine
				for _, ins := range instanceSet {
					requestDeleteMachine := cwp.NewDeleteMachineRequest()
					requestDeleteMachine.Uuid = ins.Uuid
					if _, err := client.UseCwpClient().DeleteMachine(requestDeleteMachine); err != nil {
						log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
							logId, requestDeleteMachine.GetAction(), requestDeleteMachine.ToJsonString(), err.Error())
						return err
					}
				}
			}
		} else {
			// default is true, install security agent
			// tat/InvokeCommand, CommandId=cmd-d8jj2skv, instanceId is enough
			// https://cloud.tencent.com/document/product/1340/52678
			for i := 0; i < len(workersInsIdOfNodePool); i += BatchProcessedInsLimit {
				var reqInstanceIds []string
				if i+BatchProcessedInsLimit <= len(workersInsIdOfNodePool) {
					reqInstanceIds = workersInsIdOfNodePool[i : i+BatchProcessedInsLimit]
				} else {
					reqInstanceIds = workersInsIdOfNodePool[i:]
				}
				requestInvokeCommand := tat.NewInvokeCommandRequest()
				requestInvokeCommand.InstanceIds = helper.StringsStringsPoint(reqInstanceIds)
				requestInvokeCommand.CommandId = helper.String(InstallSecurityAgentCommandId)
				requestInvokeCommand.Parameters = helper.String("{}")
				requestInvokeCommand.Timeout = helper.Uint64(60)
				_, err := client.UseTatClient().InvokeCommand(requestInvokeCommand)
				if err != nil {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, requestInvokeCommand.GetAction(), requestInvokeCommand.ToJsonString(), err.Error())
					return err
				}
			}
		}
	}
	return nil
}

func ModifyClusterInternetOrIntranetAccess(ctx context.Context, d *schema.ResourceData, tkeSvc *TkeService,
	isInternet bool, enable bool, sg string, subnetId string, domain string) error {

	id := d.Id()
	var accessType string
	if isInternet {
		accessType = "cluster internet"
	} else {
		accessType = "cluster intranet"
	}
	// open access
	if enable {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := tkeSvc.CreateClusterEndpoint(ctx, id, subnetId, sg, isInternet, domain, "")
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			status, message, inErr := tkeSvc.DescribeClusterEndpointStatus(ctx, id, isInternet)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			if status == TkeInternetStatusCreating {
				return resource.RetryableError(
					fmt.Errorf("%s create %s endpoint status still is %s", id, accessType, status))
			}
			if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
				return nil
			}
			return resource.NonRetryableError(
				fmt.Errorf("%s create %s endpoint error ,status is %s,message is %s", id, accessType, status, message))
		})
		if err != nil {
			return err
		}
	} else { // close access
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr := tkeSvc.DeleteClusterEndpoint(ctx, id, isInternet)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			status, message, inErr := tkeSvc.DescribeClusterEndpointStatus(ctx, id, isInternet)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			if status == TkeInternetStatusDeleting {
				return resource.RetryableError(
					fmt.Errorf("%s close %s endpoint status still is %s", id, accessType, status))
			}
			if status == TkeInternetStatusNotfound || status == TkeInternetStatusDeleted || status == TkeInternetStatusCreated {
				return nil
			}
			return resource.NonRetryableError(
				fmt.Errorf("%s close %s endpoint error ,status is %s,message is %s", id, accessType, status, message))
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (me *TkeService) createBackupStorageLocation(ctx context.Context, request *tke.CreateBackupStorageLocationRequest) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().CreateBackupStorageLocation(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DescribeBackupStorageLocations(ctx context.Context, names []string) (locations []*tke.BackupStorageLocation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeBackupStorageLocationsRequest()
	request.Names = common.StringPtrs(names)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTkeClient().DescribeBackupStorageLocations(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response != nil && len(response.Response.BackupStorageLocationSet) > 0 {
		locations = append(locations, response.Response.BackupStorageLocationSet...)
	}
	return

}

func (me *TkeService) DeleteBackupStorageLocation(ctx context.Context, name string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := tke.NewDeleteBackupStorageLocationRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.Name = common.StringPtr(name)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTkeClient().DeleteBackupStorageLocation(request)
	return err
}

func (me *TkeService) DescribeTkeEncryptionProtectionById(ctx context.Context, clusterId string) (encryptionProtection *tke.DescribeEncryptionStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeEncryptionStatusRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().DescribeEncryptionStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	encryptionProtection = response.Response
	return
}

func (me *TkeService) DeleteTkeEncryptionProtectionById(ctx context.Context, clusterId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDisableEncryptionProtectionRequest()
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().DisableEncryptionProtection(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) TkeEncryptionProtectionStateRefreshFunc(clusterId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		logId := tccommon.GetLogId(ctx)

		request := tke.NewDescribeEncryptionStatusRequest()
		request.ClusterId = helper.String(clusterId)

		var errRet error
		defer func() {
			if errRet != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), errRet.Error())
			}
		}()

		object, err := me.client.UseTkeClient().DescribeEncryptionStatus(request)

		if err != nil {
			return nil, "", err
		}

		if err != nil {
			errRet = err
			return object, "", err
		}

		return object, helper.PString(object.Response.Status), nil
	}
}

func (me *TkeService) DescribeKubernetesClusterInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret []*tke.Instance, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke.NewDescribeClusterInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "InstanceIds" {
			request.InstanceIds = v.([]*string)
		}
		if k == "InstanceRole" {
			request.InstanceRole = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tke.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTkeV20180525Client().DescribeClusterInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		ret = append(ret, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TkeService) DescribeKubernetesClusterNodePoolsByFilter(ctx context.Context, param map[string]interface{}) (ret []*tke.NodePool, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke.NewDescribeClusterNodePoolsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tke.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterNodePools(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.NodePoolSet) < 1 {
		return
	}

	ret = response.Response.NodePoolSet
	return
}

func (me *TkeService) DescribeKubernetesAddonById(ctx context.Context, clusterId string, addonName string) (ret *tke.Addon, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeAddonRequest()
	request.ClusterId = &clusterId
	request.AddonName = &addonName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().DescribeAddon(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Addons) < 1 {
		return
	}

	ret = response.Response.Addons[0]
	return
}

func (me *TkeService) DescribeKubernetesClustersByFilter(ctx context.Context, param map[string]interface{}) (ret []*tke.Cluster, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke.NewDescribeClustersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	if err := dataSourceTencentCloudKubernetesClustersReadPreRequest0(ctx, request); err != nil {
		return nil, err
	}

	response, err := me.client.UseTkeV20180525Client().DescribeClusters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Clusters) < 1 {
		return
	}

	ret = response.Response.Clusters
	return
}

func (me *TkeService) DescribeKubernetesClusterLevelsByFilter(ctx context.Context, param map[string]interface{}) (ret []*tke.ClusterLevelAttribute, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke.NewDescribeClusterLevelAttributeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterID" {
			request.ClusterID = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterLevelAttribute(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Items) < 1 {
		return
	}

	ret = response.Response.Items
	return
}

func (me *TkeService) DescribeKubernetesClusterCommonNamesByFilter(ctx context.Context, param map[string]interface{}) (ret []*tke.CommonName, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke.NewDescribeClusterCommonNamesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SubaccountUins" {
			request.SubaccountUins = v.([]*string)
		}
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "RoleIds" {
			request.RoleIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterCommonNames(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.CommonNames) < 1 {
		return
	}

	ret = response.Response.CommonNames
	return
}

func (me *TkeService) DescribeKubernetesClusterAuthenticationOptionsByFilter(ctx context.Context, param map[string]interface{}) (ret *tke.DescribeClusterAuthenticationOptionsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke.NewDescribeClusterAuthenticationOptionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterAuthenticationOptions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesChartsByFilter(ctx context.Context, param map[string]interface{}) (ret []*tke.AppChart, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke.NewGetTkeAppChartListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Kind" {
			request.Kind = v.(*string)
		}
		if k == "Arch" {
			request.Arch = v.(*string)
		}
		if k == "ClusterType" {
			request.ClusterType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().GetTkeAppChartList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AppCharts) < 1 {
		return
	}

	ret = response.Response.AppCharts
	return
}

func (me *TkeService) DescribeKubernetesEncryptionProtectionById(ctx context.Context, clusterId string) (ret *tke.DescribeEncryptionStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeEncryptionStatusRequest()
	request.ClusterId = helper.String(clusterId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeEncryptionStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesClusterAttachmentById(ctx context.Context, clusterId string) (ret *tke.Cluster, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClustersRequest()
	request.ClusterIds = []*string{helper.String(clusterId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().DescribeClusters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Clusters) < 1 {
		return
	}

	ret = response.Response.Clusters[0]
	return
}

func (me *TkeService) DescribeKubernetesClusterAttachmentById1(ctx context.Context, instanceId string) (ret *cvm.Instance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{helper.String(instanceId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCvmClient().DescribeInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceSet) < 1 {
		return
	}

	ret = response.Response.InstanceSet[0]
	return
}

func (me *TkeService) DescribeKubernetesClusterAttachmentById2(ctx context.Context, instanceId string, clusterId string) (ret *tke.Instance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterInstancesRequest()
	request.ClusterId = helper.String(clusterId)
	request.InstanceIds = []*string{helper.String(instanceId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().DescribeClusterInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceSet) < 1 {
		return
	}

	for _, info := range response.Response.InstanceSet {
		if info.InstanceId != nil && *info.InstanceId == instanceId {
			ret = info
			break
		}
	}
	return
}

func (me *TkeService) DescribeKubernetesBackupStorageLocationById(ctx context.Context, name string) (ret *tke.BackupStorageLocation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeBackupStorageLocationsRequest()
	request.Names = []*string{helper.String(name)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeBackupStorageLocations(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.BackupStorageLocationSet) < 1 {
		return
	}

	for _, info := range response.Response.BackupStorageLocationSet {
		if info.Name != nil && *info.Name == name {
			ret = info
			break
		}
	}
	return
}

func (me *TkeService) DescribeKubernetesClusterById(ctx context.Context, clusterId string) (ret *tke.Cluster, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClustersRequest()
	request.ClusterIds = []*string{helper.String(clusterId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Clusters) < 1 {
		return
	}

	ret = response.Response.Clusters[0]
	return
}

func (me *TkeService) DescribeKubernetesClusterById1(ctx context.Context, clusterId string) (ret *tke.DescribeClusterInstancesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterInstancesRequest()
	request.ClusterId = helper.String(clusterId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesClusterById2(ctx context.Context, clusterId string) (ret *tke.DescribeClusterSecurityResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterSecurityRequest()
	request.ClusterId = helper.String(clusterId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterSecurity(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesNodePoolById(ctx context.Context, clusterId string) (ret *tke.Cluster, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClustersRequest()
	request.ClusterIds = []*string{helper.String(clusterId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Clusters) < 1 {
		return
	}

	ret = response.Response.Clusters[0]
	return
}

func (me *TkeService) DescribeKubernetesNodePoolById1(ctx context.Context, clusterId string, nodePoolId string) (ret *tke.NodePool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterNodePoolDetailRequest()
	request.ClusterId = helper.String(clusterId)
	request.NodePoolId = helper.String(nodePoolId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterNodePoolDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.NodePool
	return
}

func (me *TkeService) DescribeKubernetesServerlessNodePoolById(ctx context.Context, clusterId string, nodePoolId string) (ret *tke.VirtualNodePool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterVirtualNodePoolsRequest()
	request.ClusterId = helper.String(clusterId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterVirtualNodePools(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.NodePoolSet) < 1 {
		return
	}

	for _, info := range response.Response.NodePoolSet {
		if info.NodePoolId != nil && *info.NodePoolId == nodePoolId {
			ret = info
			break
		}
	}
	return
}

func (me *TkeService) DescribeKubernetesAuthAttachmentById(ctx context.Context, clusterId string) (ret *tke.DescribeClusterAuthenticationOptionsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterAuthenticationOptionsRequest()
	request.ClusterId = helper.String(clusterId)

	if err := resourceTencentCloudKubernetesAuthAttachmentReadPostFillRequest0(ctx, request); err != nil {
		return nil, err
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterAuthenticationOptions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TkeService) DeleteKubernetesAuthAttachmentById(ctx context.Context, clusterId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewModifyClusterAuthenticationOptionsRequest()
	request.ClusterId = &clusterId
	request.ServiceAccounts = &tke.ServiceAccountAuthenticationOptions{
		Issuer:  helper.String("https://kubernetes.default.svc.cluster.local"),
		JWKSURI: helper.String(""),
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeClient().ModifyClusterAuthenticationOptions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TkeService) DescribeKubernetesScaleWorkerById(ctx context.Context, clusterId string) (ret *tke.DescribeClustersResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClustersRequest()
	request.ClusterIds = []*string{helper.String(clusterId)}

	if err := resourceTencentCloudKubernetesScaleWorkerReadPostFillRequest0(ctx, request); err != nil {
		return nil, err
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if err := resourceTencentCloudKubernetesScaleWorkerReadPostRequest0(ctx, request, response); err != nil {
		return nil, err
	}

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesScaleWorkerById1(ctx context.Context) (ret *cvm.DescribeInstancesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cvm.NewDescribeInstancesRequest()

	if err := resourceTencentCloudKubernetesScaleWorkerReadPostFillRequest1(ctx, request); err != nil {
		return nil, err
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := resourceTencentCloudKubernetesScaleWorkerReadPreRequest1(ctx, request)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if err := resourceTencentCloudKubernetesScaleWorkerReadPostRequest1(ctx, request, response); err != nil {
		return nil, err
	}

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesNativeNodePoolById(ctx context.Context, clusterId string, nodePoolId string) (ret *tke2.NodePool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke2.NewDescribeNodePoolsRequest()
	request.ClusterId = &clusterId
	filter := &tke2.Filter{
		Name:   helper.String("NodePoolsId"),
		Values: []*string{&nodePoolId},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	var instances []*tke2.NodePool
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTke2Client().DescribeNodePools(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.NodePools) < 1 {
			break
		}
		instances = append(instances, response.Response.NodePools...)
		if len(response.Response.NodePools) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}

	ret = instances[0]
	return
}

func (me *TkeService) DescribeKubernetesClusterNativeNodePoolsByFilter(ctx context.Context, param map[string]interface{}) (ret []*tke2.NodePool, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tke2.NewDescribeNodePoolsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ClusterId" {
			request.ClusterId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tke2.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTke2Client().DescribeNodePools(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.NodePools) < 1 {
		return
	}

	ret = response.Response.NodePools
	return
}

func (me *TkeService) DescribeKubernetesAddonAttachmentById(ctx context.Context) (ret *tke.ForwardApplicationRequestV3ResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewForwardApplicationRequestV3Request()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	if err := resourceTencentCloudKubernetesAddonAttachmentReadPreRequest0(ctx, request); err != nil {
		return nil, err
	}

	response, err := me.client.UseTkeV20180525Client().ForwardApplicationRequestV3(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesHealthCheckPolicyById(ctx context.Context, clusterId string, name string) (ret *tke2.HealthCheckPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke2.NewDescribeHealthCheckPoliciesRequest()
	request.ClusterId = helper.String(clusterId)
	filter := &tke2.Filter{
		Name:   helper.String("HealthCheckPolicyName"),
		Values: []*string{helper.String(name)},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	var instances []*tke2.HealthCheckPolicy
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTkeV20220501Client().DescribeHealthCheckPolicies(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.HealthCheckPolicies) < 1 {
			break
		}
		instances = append(instances, response.Response.HealthCheckPolicies...)
		if len(response.Response.HealthCheckPolicies) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}

	ret = instances[0]
	return
}

func (me *TkeService) DescribeKubernetesLogConfigById(ctx context.Context, clusterId string, logConfigName string, clusterType string) (ret *tke.DescribeLogConfigsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeLogConfigsRequest()
	request.ClusterId = helper.String(clusterId)
	request.ClusterType = helper.String(clusterType)
	request.LogConfigNames = helper.String(logConfigName)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeLogConfigs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesClusterMasterAttachmentById(ctx context.Context, clusterId string) (ret *tke.Cluster, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClustersRequest()
	request.ClusterIds = []*string{helper.String(clusterId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusters(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Clusters) < 1 {
		return
	}

	ret = response.Response.Clusters[0]
	return
}

func (me *TkeService) DescribeKubernetesClusterMasterAttachmentById1(ctx context.Context, instanceId string) (ret *cvm.Instance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{helper.String(instanceId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCvmV20170312Client().DescribeInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceSet) < 1 {
		return
	}

	ret = response.Response.InstanceSet[0]
	return
}

func (me *TkeService) DescribeKubernetesClusterMasterAttachmentById2(ctx context.Context, clusterId string, instanceId string, nodeRole string) (ret *tke.DescribeClusterInstancesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeClusterInstancesRequest()
	request.ClusterId = helper.String(clusterId)
	request.InstanceIds = []*string{helper.String(instanceId)}
	request.InstanceRole = helper.String(nodeRole)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeClusterInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *TkeService) DescribeKubernetesClusterAuditById(ctx context.Context, clusterId string) (ret *tke.DescribeLogSwitchesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tke.NewDescribeLogSwitchesRequest()
	request.ClusterIds = []*string{helper.String(clusterId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTkeV20180525Client().DescribeLogSwitches(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}
