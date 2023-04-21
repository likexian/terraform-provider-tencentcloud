package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type LightHouseService struct {
	client *connectivity.TencentCloudClient
}

// lighthouse instance

func (me *LightHouseService) DescribeLighthouseInstanceById(ctx context.Context, instanceId string) (instance *lighthouse.Instance, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceIds = append(request.InstanceIds, helper.String(instanceId))
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*lighthouse.Instance, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseLighthouseClient().DescribeInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceSet) < 1 {
			break
		}
		instances = append(instances, response.Response.InstanceSet...)
		if len(response.Response.InstanceSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	instance = instances[0]

	return
}

func (me *LightHouseService) DeleteLighthouseInstanceById(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewTerminateInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, &id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseLighthouseClient().TerminateInstances(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) IsolateLighthouseInstanceById(ctx context.Context, id string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewIsolateInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, &id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseLighthouseClient().IsolateInstances(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) LighthouseBlueprintStateRefreshFunc(blueprintId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseBlueprintById(ctx, blueprintId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.BlueprintState), nil
	}
}

func (me *LightHouseService) DeleteLighthouseBlueprintById(ctx context.Context, blueprintId string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDeleteBlueprintsRequest()
	request.BlueprintIds = []*string{&blueprintId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DeleteBlueprints(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) DescribeLighthouseBlueprintById(ctx context.Context, blueprintId string) (blueprint *lighthouse.Blueprint, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeBlueprintsRequest()
	request.BlueprintIds = []*string{&blueprintId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeBlueprints(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.BlueprintSet) < 1 {
		return
	}

	blueprint = response.Response.BlueprintSet[0]
	return
}

func (me *LightHouseService) DescribeLighthouseFirewallRuleById(ctx context.Context, instance_id string) (firewallRules []*lighthouse.FirewallRuleInfo, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeFirewallRulesRequest()

	request.InstanceId = helper.String(instance_id)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	limit := 50
	offset := 0
	firewallRules = make([]*lighthouse.FirewallRuleInfo, 0)
	for {
		ratelimit.Check(request.GetAction())
		request.Limit = helper.IntInt64(limit)
		request.Offset = helper.IntInt64(offset)
		response, err := me.client.UseLighthouseClient().DescribeFirewallRules(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		firewallRules = append(firewallRules, response.Response.FirewallRuleSet...)

		if len(response.Response.FirewallRuleSet) < limit {
			break
		}
		offset += limit
	}

	return
}

func (me *LightHouseService) DescribeLighthouseFirewallRulesTemplateByFilter(ctx context.Context) (firewallRulesTemplate []*lighthouse.FirewallRuleInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = lighthouse.NewDescribeFirewallRulesTemplateRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeFirewallRulesTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil && len(response.Response.FirewallRuleSet) != 0 {
		firewallRulesTemplate = append(firewallRulesTemplate, response.Response.FirewallRuleSet...)
	}
	return
}

func (me *LightHouseService) DescribeLighthouseDiskBackupById(ctx context.Context, diskBackupId string) (diskBackup *lighthouse.DiskBackup, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeDiskBackupsRequest()
	request.DiskBackupIds = []*string{&diskBackupId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeDiskBackups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DiskBackupSet) < 1 {
		return
	}

	diskBackup = response.Response.DiskBackupSet[0]
	return
}

func (me *LightHouseService) DeleteLighthouseDiskBackupById(ctx context.Context, diskBackupId string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDeleteDiskBackupsRequest()
	request.DiskBackupIds = []*string{&diskBackupId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DeleteDiskBackups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *LightHouseService) LighthouseDiskBackupStateRefreshFunc(diskBackupId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseDiskBackupById(ctx, diskBackupId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.DiskBackupState), nil
	}
}

func (me *LightHouseService) LighthouseApplyDiskBackupStateRefreshFunc(diskBackupId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseDiskBackupById(ctx, diskBackupId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.LatestOperationState), nil
	}
}

func (me *LightHouseService) DescribeLighthouseDiskAttachmentById(ctx context.Context, diskId string) (diskAttachment *lighthouse.Disk, errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDescribeDisksRequest()
	request.DiskIds = []*string{&diskId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DescribeDisks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DiskSet) < 1 {
		return
	}

	diskAttachment = response.Response.DiskSet[0]
	return
}

func (me *LightHouseService) LighthouseDiskAttachmentStateRefreshFunc(diskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeLighthouseDiskAttachmentById(ctx, diskId)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.DiskState), nil
	}
}

func (me *LightHouseService) DeleteLighthouseDiskAttachmentById(ctx context.Context, diskId string) (errRet error) {
	logId := getLogId(ctx)

	request := lighthouse.NewDetachDisksRequest()
	request.DiskIds = []*string{&diskId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseLighthouseClient().DetachDisks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
