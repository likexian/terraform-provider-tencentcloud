package emr

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewEMRService(client *connectivity.TencentCloudClient) EMRService {
	return EMRService{client: client}
}

type EMRService struct {
	client *connectivity.TencentCloudClient
}

func (me *EMRService) UpdateInstance(ctx context.Context, request *emr.ScaleOutInstanceRequest) (id string, err error) {
	logId := tccommon.GetLogId(ctx)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseEmrClient().ScaleOutInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	id = *response.Response.InstanceId
	return
}

func (me *EMRService) DeleteInstance(ctx context.Context, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewTerminateInstanceRequest()
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = common.StringPtr(v.(string))
	}
	ratelimit.Check(request.GetAction())
	//API: https://cloud.tencent.com/document/api/589/34261
	_, err := me.client.UseEmrClient().TerminateInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	return nil
}

func (me *EMRService) CreateInstance(ctx context.Context, d *schema.ResourceData) (id string, err error) {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewCreateInstanceRequest()

	if v, ok := d.GetOk("auto_renew"); ok {
		request.AutoRenew = common.Uint64Ptr((uint64)(v.(int)))
	}

	if v, ok := d.GetOk("product_id"); ok {
		request.ProductId = common.Uint64Ptr((uint64)(v.(int)))
	}

	if v, ok := d.GetOk("vpc_settings"); ok {
		value := v.(map[string]interface{})
		var vpcId string
		var subnetId string

		if subV, ok := value["vpc_id"]; ok {
			vpcId = subV.(string)
		}
		if subV, ok := value["subnet_id"]; ok {
			subnetId = subV.(string)
		}
		vpcSettings := &emr.VPCSettings{VpcId: &vpcId, SubnetId: &subnetId}
		request.VPCSettings = vpcSettings
	}

	if v, ok := d.GetOk("softwares"); ok {
		softwares := v.(*schema.Set).List()
		request.Software = make([]*string, 0)
		for _, software := range softwares {
			request.Software = append(request.Software, common.StringPtr(software.(string)))
		}
	}

	if v, ok := d.GetOk("resource_spec"); ok {
		tmpResourceSpec := v.([]interface{})
		resourceSpec := tmpResourceSpec[0].(map[string]interface{})
		request.ResourceSpec = &emr.NewResourceSpec{}
		for k, v := range resourceSpec {
			if k == "master_resource_spec" {
				if len(v.([]interface{})) > 0 {
					request.ResourceSpec.MasterResourceSpec = ParseResource(v.([]interface{})[0].(map[string]interface{}))
				}
			} else if k == "core_resource_spec" {
				if len(v.([]interface{})) > 0 {
					request.ResourceSpec.CoreResourceSpec = ParseResource(v.([]interface{})[0].(map[string]interface{}))
				}
			} else if k == "task_resource_spec" {
				if len(v.([]interface{})) > 0 {
					request.ResourceSpec.TaskResourceSpec = ParseResource(v.([]interface{})[0].(map[string]interface{}))
				}
			} else if k == "master_count" {
				request.ResourceSpec.MasterCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "core_count" {
				request.ResourceSpec.CoreCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "task_count" {
				request.ResourceSpec.TaskCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "common_resource_spec" {
				if len(v.([]interface{})) > 0 {
					request.ResourceSpec.CommonResourceSpec = ParseResource(v.([]interface{})[0].(map[string]interface{}))
				}
			} else if k == "common_count" {
				request.ResourceSpec.CommonCount = common.Int64Ptr((int64)(v.(int)))
			}
		}
	}

	if v, ok := d.GetOk("support_ha"); ok {
		request.SupportHA = common.Uint64Ptr((uint64)(v.(int)))
	} else {
		request.SupportHA = common.Uint64Ptr(0)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = common.StringPtr(v.(string))
	}

	needMasterWan := d.Get("need_master_wan").(string)
	request.NeedMasterWan = common.StringPtr(needMasterWan)
	payMode := d.Get("pay_mode")
	request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
	if v, ok := d.GetOk("placement"); ok {
		request.Placement = &emr.Placement{}
		placement := v.(map[string]interface{})

		if projectId, ok := placement["project_id"]; ok {
			projectIdInt64, _ := strconv.ParseInt(projectId.(string), 10, 64)
			request.Placement.ProjectId = common.Int64Ptr(projectIdInt64)
		} else {
			request.Placement.ProjectId = common.Int64Ptr(0)
		}
		if zone, ok := placement["zone"]; ok {
			request.Placement.Zone = common.StringPtr(zone.(string))
		}
	}

	if v, ok := d.GetOk("placement_info"); ok {
		request.Placement = &emr.Placement{}
		placementList := v.([]interface{})
		placement := placementList[0].(map[string]interface{})

		if v, ok := placement["project_id"]; ok {
			projectId := v.(int)
			request.Placement.ProjectId = helper.IntInt64(projectId)
		} else {
			request.Placement.ProjectId = helper.IntInt64(0)
		}
		if zone, ok := placement["zone"]; ok {
			request.Placement.Zone = common.StringPtr(zone.(string))
		}
	}

	if v, ok := d.GetOk("time_span"); ok {
		request.TimeSpan = common.Uint64Ptr((uint64)(v.(int)))
	}
	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = common.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("login_settings"); ok {
		request.LoginSettings = &emr.LoginSettings{}
		loginSettings := v.(map[string]interface{})
		if password, ok := loginSettings["password"]; ok {
			request.LoginSettings.Password = common.StringPtr(password.(string))
		}
		if publicKeyId, ok := loginSettings["public_key_id"]; ok {
			request.LoginSettings.PublicKeyId = common.StringPtr(publicKeyId.(string))
		}
	}
	if v, ok := d.GetOk("sg_id"); ok {
		request.SgId = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOk("extend_fs_field"); ok {
		request.ExtendFsField = common.StringPtr(v.(string))
	}
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		emrTags := make([]*emr.Tag, 0)
		for k, v := range tags {
			tagKey := k
			tagValue := v
			emrTags = append(emrTags, &emr.Tag{
				TagKey:   helper.String(tagKey),
				TagValue: helper.String(tagValue),
			})
		}
		request.Tags = emrTags
	}

	ratelimit.Check(request.GetAction())
	//API: https://cloud.tencent.com/document/api/589/34261
	response, err := me.client.UseEmrClient().CreateInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	id = *response.Response.InstanceId
	return
}

func (me *EMRService) DescribeInstances(ctx context.Context, filters map[string]interface{}) (clusters []*emr.ClusterInstancesInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewDescribeInstancesRequest()

	ratelimit.Check(request.GetAction())
	// API: https://cloud.tencent.com/document/api/589/41707
	if v, ok := filters["instance_ids"]; ok {
		instances := v.([]interface{})
		request.InstanceIds = make([]*string, 0)
		for _, instance := range instances {
			request.InstanceIds = append(request.InstanceIds, common.StringPtr(instance.(string)))
		}
	}
	if v, ok := filters["display_strategy"]; ok {
		request.DisplayStrategy = common.StringPtr(v.(string))
	}
	if v, ok := filters["project_id"]; ok {
		request.ProjectId = common.Int64Ptr(v.(int64))
	}
	response, err := me.client.UseEmrClient().DescribeInstances(request)

	if err != nil {
		if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkError.Code == "ResourceNotFound.ClusterNotFound" {
				return
			}
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	clusters = response.Response.ClusterList
	return
}

func (me *EMRService) DescribeInstancesById(ctx context.Context, instanceId string, displayStrategy string) (clusters []*emr.ClusterInstancesInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewDescribeInstancesRequest()

	ratelimit.Check(request.GetAction())
	// API: https://cloud.tencent.com/document/api/589/41707
	request.ProjectId = helper.IntInt64(-1)
	request.InstanceIds = make([]*string, 0)
	request.InstanceIds = append(request.InstanceIds, common.StringPtr(instanceId))
	request.DisplayStrategy = common.StringPtr(displayStrategy)

	response, err := me.client.UseEmrClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	clusters = response.Response.ClusterList
	return
}

func (me *EMRService) DescribeClusterNodes(ctx context.Context, instanceId, nodeFlag, hardwareResourceType string, offset, limit int) (nodes []*emr.NodeHardwareInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := emr.NewDescribeClusterNodesRequest()

	ratelimit.Check(request.GetAction())
	// API: https://cloud.tencent.com/document/api/589/41707
	request.InstanceId = &instanceId
	request.NodeFlag = &nodeFlag
	request.HardwareResourceType = &hardwareResourceType
	request.Limit = helper.IntInt64(limit)
	request.Offset = helper.IntInt64(offset)
	response, err := me.client.UseEmrClient().DescribeClusterNodes(request)

	if err != nil {
		if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkError.Code == "ResourceNotFound.ClusterNotFound" {
				return
			}
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	nodes = response.Response.NodeList
	return
}

func (me *EMRService) ModifyResourcesTags(ctx context.Context, region string, instanceId string, oldTags, newTags map[string]interface{}) error {
	resourceName := tccommon.BuildTagResourceName("emr", "emr-instance", region, instanceId)
	rTags, dTags := svctag.DiffTags(oldTags, newTags)
	tagService := svctag.NewTagService(me.client)
	if err := tagService.ModifyTags(ctx, resourceName, rTags, dTags); err != nil {
		return err
	}

	addTags := make([]*emr.Tag, 0)
	modifyTags := make([]*emr.Tag, 0)
	deleteTags := make([]*emr.Tag, 0)
	for k, v := range newTags {
		tagKey := k
		tageValue := v.(string)
		_, ok := oldTags[tagKey]
		if !ok {
			addTags = append(addTags, &emr.Tag{
				TagKey:   &tagKey,
				TagValue: &tageValue,
			})
		} else if oldTags[tagKey].(string) != tageValue {
			modifyTags = append(modifyTags, &emr.Tag{
				TagKey:   &tagKey,
				TagValue: &tageValue,
			})
		}
	}
	for k, v := range oldTags {
		tagKey := k
		tageValue := v.(string)
		_, ok := newTags[tagKey]
		if !ok {
			deleteTags = append(deleteTags, &emr.Tag{
				TagKey:   &tagKey,
				TagValue: &tageValue,
			})
		}
	}
	modifyResourceTags := &emr.ModifyResourceTags{
		Resource:       helper.String(resourceName),
		ResourceId:     helper.String(instanceId),
		ResourceRegion: helper.String(region),
	}
	if len(addTags) > 0 {
		modifyResourceTags.AddTags = addTags
	}
	if len(modifyTags) > 0 {
		modifyResourceTags.ModifyTags = modifyTags
	}
	if len(deleteTags) > 0 {
		modifyResourceTags.DeleteTags = deleteTags
	}

	request := emr.NewModifyResourcesTagsRequest()
	ratelimit.Check(request.GetAction())
	request.ModifyType = helper.String("Cluster")
	request.ModifyResourceTagsInfoList = []*emr.ModifyResourceTags{modifyResourceTags}

	response, err := me.client.UseEmrClient().ModifyResourcesTags(request)
	if err != nil {
		return err
	}
	if response != nil && response.Response != nil && len(response.Response.FailList) > 0 {
		return fmt.Errorf("file resource list: %v", response.Response.FailList)
	}
	return nil
}

func (me *EMRService) DescribeEmrUserManagerById(ctx context.Context, instanceId string, userName string) (userManager *emr.DescribeUsersForUserManagerResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDescribeUsersForUserManagerRequest()
	request.InstanceId = &instanceId
	request.UserManagerFilter = &emr.UserManagerFilter{
		UserName: &userName,
	}
	request.PageNo = helper.IntInt64(0)
	request.PageSize = helper.IntInt64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeUsersForUserManager(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	userManager = response.Response
	return
}

func (me *EMRService) DeleteEmrUserManagerById(ctx context.Context, instanceId string, userName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := emr.NewDeleteUserManagerUserListRequest()
	request.InstanceId = &instanceId
	request.UserNameList = []*string{&userName}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DeleteUserManagerUserList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *EMRService) DescribeEmrCvmQuotaByFilter(ctx context.Context, param map[string]interface{}) (cvmQuota *emr.DescribeCvmQuotaResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = emr.NewDescribeCvmQuotaRequest()
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
		if k == "ZoneId" {
			request.ZoneId = v.(*int64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseEmrClient().DescribeCvmQuota(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	cvmQuota = response.Response
	return
}

func (me *EMRService) DescribeEmrAutoScaleRecordsByFilter(ctx context.Context, param map[string]interface{}) (autoScaleRecords []*emr.AutoScaleRecord, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = emr.NewDescribeAutoScaleRecordsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*emr.KeyValue)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseEmrClient().DescribeAutoScaleRecords(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RecordList) < 1 {
			break
		}
		autoScaleRecords = append(autoScaleRecords, response.Response.RecordList...)
		if len(response.Response.RecordList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
