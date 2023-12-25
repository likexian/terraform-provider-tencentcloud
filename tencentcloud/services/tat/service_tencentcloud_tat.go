package tat

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewTatService(client *connectivity.TencentCloudClient) TatService {
	return TatService{client: client}
}

type TatService struct {
	client *connectivity.TencentCloudClient
}

func (me *TatService) DescribeTatCommand(ctx context.Context, commandId string) (command *tat.Command, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tat.NewDescribeCommandsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&tat.Filter{
			Name:   helper.String("command-id"),
			Values: []*string{&commandId},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 100
	commandInfo := make([]*tat.Command, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTatClient().DescribeCommands(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CommandSet) < 1 {
			break
		}
		commandInfo = append(commandInfo, response.Response.CommandSet...)
		if len(response.Response.CommandSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(commandInfo) < 1 {
		return
	}
	command = commandInfo[0]

	return

}

func (me *TatService) DeleteTatCommandById(ctx context.Context, commandId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tat.NewDeleteCommandRequest()

	request.CommandId = &commandId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTatClient().DeleteCommand(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TatService) DescribeTatInvoker(ctx context.Context, invokerId string) (invoker *tat.Invoker, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tat.NewDescribeInvokersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&tat.Filter{
			Name:   helper.String("invoker-id"),
			Values: []*string{&invokerId},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 100
	invokerSet := make([]*tat.Invoker, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTatClient().DescribeInvokers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InvokerSet) < 1 {
			break
		}
		invokerSet = append(invokerSet, response.Response.InvokerSet...)
		if len(response.Response.InvokerSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(invokerSet) < 1 {
		return
	}
	invoker = invokerSet[0]

	return

}

func (me *TatService) DeleteTatInvokerById(ctx context.Context, invokerId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tat.NewDeleteInvokerRequest()

	request.InvokerId = &invokerId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTatClient().DeleteInvoker(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TatService) DescribeTatCommandByFilter(ctx context.Context, param map[string]interface{}) (command []*tat.Command, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tat.NewDescribeCommandsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "command_id" {
			request.Filters = append(
				request.Filters,
				&tat.Filter{
					Name:   helper.String("command-id"),
					Values: []*string{v.(*string)},
				},
			)
		}

		if k == "command_name" {
			request.Filters = append(
				request.Filters,
				&tat.Filter{
					Name:   helper.String("command-name"),
					Values: []*string{v.(*string)},
				},
			)
		}

		if k == "command_type" {
			request.Filters = append(
				request.Filters,
				&tat.Filter{
					Name:   helper.String("command-type"),
					Values: []*string{v.(*string)},
				},
			)
		}

		if k == "created_by" {
			request.Filters = append(
				request.Filters,
				&tat.Filter{
					Name:   helper.String("created-by"),
					Values: []*string{v.(*string)},
				},
			)
		}
	}

	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTatClient().DescribeCommands(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CommandSet) < 1 {
			break
		}
		command = append(command, response.Response.CommandSet...)
		if len(response.Response.CommandSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *TatService) DescribeTatInvokerByFilter(ctx context.Context, param map[string]interface{}) (invoker []*tat.Invoker, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tat.NewDescribeInvokersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "invoker_id" {
			request.Filters = append(
				request.Filters,
				&tat.Filter{
					Name:   helper.String("invoker-id"),
					Values: []*string{v.(*string)},
				},
			)
		}

		if k == "command_id" {
			request.Filters = append(
				request.Filters,
				&tat.Filter{
					Name:   helper.String("command-id"),
					Values: []*string{v.(*string)},
				},
			)
		}

		if k == "type" {
			request.Filters = append(
				request.Filters,
				&tat.Filter{
					Name:   helper.String("type"),
					Values: []*string{v.(*string)},
				},
			)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTatClient().DescribeInvokers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InvokerSet) < 1 {
			break
		}
		invoker = append(invoker, response.Response.InvokerSet...)
		if len(response.Response.InvokerSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *TatService) DescribeTatAgentByFilter(ctx context.Context, param map[string]interface{}) (agent []*tat.AutomationAgentInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tat.NewDescribeAutomationAgentStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceIds" {
			request.InstanceIds = v.([]*string)
		}
		if k == "filters" {
			request.Filters = v.([]*tat.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTatClient().DescribeAutomationAgentStatus(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AutomationAgentSet) < 1 {
			break
		}
		agent = append(agent, response.Response.AutomationAgentSet...)
		if len(response.Response.AutomationAgentSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TatService) DescribeTatInvokerRecordsByFilter(ctx context.Context, param map[string]interface{}) (invokerRecords []*tat.InvokerRecord, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tat.NewDescribeInvokerRecordsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InvokerIds" {
			request.InvokerIds = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTatClient().DescribeInvokerRecords(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InvokerRecordSet) < 1 {
			break
		}
		invokerRecords = append(invokerRecords, response.Response.InvokerRecordSet...)
		if len(response.Response.InvokerRecordSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TatService) DescribeTatInvokerConfigById(ctx context.Context, invokerId string) (invokerAttachment *tat.Invoker, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tat.NewDescribeInvokersRequest()
	request.InvokerIds = []*string{&invokerId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*tat.Invoker, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTatClient().DescribeInvokers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InvokerSet) < 1 {
			break
		}
		instances = append(instances, response.Response.InvokerSet...)
		if len(response.Response.InvokerSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	invokerAttachment = instances[0]
	return
}

func (me *TatService) DescribeTatInvocationById(ctx context.Context, invocationId string) (invocation *tat.Invocation, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tat.NewDescribeInvocationsRequest()
	request.InvocationIds = []*string{&invocationId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*tat.Invocation, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTatClient().DescribeInvocations(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InvocationSet) < 1 {
			break
		}
		instances = append(instances, response.Response.InvocationSet...)
		if len(response.Response.InvocationSet) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	invocation = instances[0]
	return
}

func (me *TatService) DeleteTatInvocationById(ctx context.Context, invocationId, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tat.NewCancelInvocationRequest()
	request.InvocationId = &invocationId
	request.InstanceIds = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTatClient().CancelInvocation(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TatService) DescribeTatInvocationTaskByFilter(ctx context.Context, param map[string]interface{}) (invocationTaskDatasource []*tat.InvocationTask, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tat.NewDescribeInvocationTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InvocationTaskIds" {
			request.InvocationTaskIds = v.([]*string)
		}
		if k == "filters" {
			request.Filters = v.([]*tat.Filter)
		}
		if k == "HideOutput" {
			request.HideOutput = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTatClient().DescribeInvocationTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InvocationTaskSet) < 1 {
			break
		}
		invocationTaskDatasource = append(invocationTaskDatasource, response.Response.InvocationTaskSet...)
		if len(response.Response.InvocationTaskSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
