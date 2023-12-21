package dbbrain

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewDbbrainService(client *connectivity.TencentCloudClient) DbbrainService {
	return DbbrainService{client: client}
}

type DbbrainService struct {
	client *connectivity.TencentCloudClient
}

func (me *DbbrainService) DescribeDbbrainSqlFilter(ctx context.Context, instanceId, filterId *string) (sqlFilter *dbbrain.SQLFilter, errRet error) {
	param := make(map[string]interface{})
	if instanceId != nil {
		param["instance_id"] = instanceId
	}
	if filterId != nil {
		param["filter_ids"] = []*int64{helper.StrToInt64Point(*filterId)}
	}

	ret, errRet := me.DescribeDbbrainSqlFiltersByFilter(ctx, param)
	if errRet != nil {
		return
	}
	if ret != nil {
		return ret[0], nil
	}
	return
}

func (me *DbbrainService) DescribeDbbrainSqlFiltersByFilter(ctx context.Context, param map[string]interface{}) (sqlFilters []*dbbrain.SQLFilter, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeSqlFiltersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query objects", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}

		if k == "filter_ids" {
			request.FilterIds = v.([]*int64)
		}

		if k == "statuses" {
			request.Statuses = v.([]*string)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDbbrainClient().DescribeSqlFilters(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		sqlFilters = append(sqlFilters, response.Response.Items...)
		if len(response.Response.Items) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *DbbrainService) getSessionToken(ctx context.Context, instanceId, user, pw, product *string) (sessionToken *string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := dbbrain.NewVerifyUserAccountRequest()

	request.InstanceId = instanceId
	request.User = user
	request.Password = pw
	if product != nil {
		request.Product = product
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "VerifyUserAccount", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().VerifyUserAccount(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	sessionToken = response.Response.SessionToken
	return
}

func (me *DbbrainService) DeleteDbbrainSqlFilterById(ctx context.Context, instanceId, filterId, sessionToken *string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dbbrain.NewDeleteSqlFiltersRequest()

	request.InstanceId = instanceId
	request.FilterIds = []*int64{helper.StrToInt64Point(*filterId)}
	request.SessionToken = sessionToken

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().DeleteSqlFilters(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DbbrainService) DescribeDbbrainSecurityAuditLogExportTask(ctx context.Context, secAuditGroupId, asyncRequestId, product *string) (task *dbbrain.SecLogExportTaskInfo, errRet error) {
	param := make(map[string]interface{})
	if secAuditGroupId != nil {
		param["sec_audit_group_id"] = secAuditGroupId
	}
	if asyncRequestId != nil {
		param["async_request_ids"] = []*uint64{helper.StrToUint64Point(*asyncRequestId)}
	}
	if product != nil {
		param["product"] = product
	} else {
		param["product"] = helper.String("mysql")
	}

	ret, errRet := me.DescribeDbbrainSecurityAuditLogExportTasksByFilter(ctx, param)
	if errRet != nil {
		return
	}
	if ret != nil {
		return ret[0], nil
	}
	return
}

func (me *DbbrainService) DescribeDbbrainSecurityAuditLogExportTasksByFilter(ctx context.Context, param map[string]interface{}) (securityAuditLogExportTasks []*dbbrain.SecLogExportTaskInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeSecurityAuditLogExportTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query objects", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "sec_audit_group_id" {
			request.SecAuditGroupId = v.(*string)
		}

		if k == "product" {
			request.Product = v.(*string)
		}

		if k == "async_request_ids" {
			request.AsyncRequestIds = v.([]*uint64)
		}
	}
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseDbbrainClient().DescribeSecurityAuditLogExportTasks(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Tasks) < 1 {
			break
		}
		securityAuditLogExportTasks = append(securityAuditLogExportTasks, response.Response.Tasks...)
		if len(response.Response.Tasks) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *DbbrainService) DeleteDbbrainSecurityAuditLogExportTaskById(ctx context.Context, secAuditGroupId, asyncRequestId, product *string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dbbrain.NewDeleteSecurityAuditLogExportTasksRequest()

	request.SecAuditGroupId = secAuditGroupId
	request.AsyncRequestIds = []*uint64{helper.StrToUint64Point(*asyncRequestId)}
	if product != nil {
		request.Product = product
	} else {
		request.Product = helper.String("mysql")
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().DeleteSecurityAuditLogExportTasks(request)
	if err != nil {
		errRet = err
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DbbrainService) DescribeDbbrainDiagEventsByFilter(ctx context.Context, param map[string]interface{}) (diagEvents []*dbbrain.DiagHistoryEventItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeDBDiagEventsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_ids" {
			request.InstanceIds = v.([]*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "severities" {
			request.Severities = v.([]*int64)
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
		response, err := me.client.UseDbbrainClient().DescribeDBDiagEvents(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		diagEvents = append(diagEvents, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DbbrainService) DescribeDbbrainDiagEventByFilter(ctx context.Context, param map[string]interface{}) (diagEvent *dbbrain.DescribeDBDiagEventResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeDBDiagEventRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "event_id" {
			request.EventId = v.(*int64)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeDBDiagEvent(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		diagEvent = response.Response
	}

	return
}

func (me *DbbrainService) DescribeDbbrainDiagHistoryByFilter(ctx context.Context, param map[string]interface{}) (diagHistory []*dbbrain.DiagHistoryEventItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeDBDiagHistoryRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeDBDiagHistory(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		diagHistory = response.Response.Events
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSecurityAuditLogDownloadUrlsByFilter(ctx context.Context, param map[string]interface{}) (securityAuditLogDownloadUrls []*string, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeSecurityAuditLogDownloadUrlsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "sec_audit_group_id" {
			request.SecAuditGroupId = v.(*string)
		}
		if k == "async_request_id" {
			request.AsyncRequestId = v.(*uint64)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeSecurityAuditLogDownloadUrls(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		securityAuditLogDownloadUrls = response.Response.Urls
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogTimeSeriesStatsByFilter(ctx context.Context, param map[string]interface{}) (slowLogTimeSeriesStats *dbbrain.DescribeSlowLogTimeSeriesStatsResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeSlowLogTimeSeriesStatsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeSlowLogTimeSeriesStats(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil {
		slowLogTimeSeriesStats = response.Response
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogTopSqlsByFilter(ctx context.Context, param map[string]interface{}) (slowLogTopSqls []*dbbrain.SlowLogTopSqlItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeSlowLogTopSqlsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "sort_by" {
			request.SortBy = v.(*string)
		}
		if k == "order_by" {
			request.OrderBy = v.(*string)
		}
		if k == "limit" {
			request.Limit = v.(*int64)
		}
		if k == "offset" {
			request.Offset = v.(*int64)
		}
		if k == "schema_list" {
			request.SchemaList = v.([]*dbbrain.SchemaItem)
		}
		if k == "product" {
			request.Product = v.(*string)
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
		response, err := me.client.UseDbbrainClient().DescribeSlowLogTopSqls(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Rows) < 1 {
			break
		}
		slowLogTopSqls = append(slowLogTopSqls, response.Response.Rows...)
		if len(response.Response.Rows) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogUserHostStatsByFilter(ctx context.Context, param map[string]interface{}) (slowLogUserHostStats []*dbbrain.SlowLogHost, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeSlowLogUserHostStatsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "start_time" {
			request.StartTime = v.(*string)
		}
		if k == "end_time" {
			request.EndTime = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
		if k == "md5" {
			request.Md5 = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeSlowLogUserHostStats(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		slowLogUserHostStats = response.Response.Items
	}

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogUserSqlAdviceByFilter(ctx context.Context, param map[string]interface{}) (slowLogUserSqlAdvice *dbbrain.DescribeUserSqlAdviceResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeUserSqlAdviceRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "sql_text" {
			request.SqlText = v.(*string)
		}
		if k == "schema" {
			request.Schema = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeUserSqlAdvice(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		slowLogUserSqlAdvice = response.Response
	}

	return
}

func (me *DbbrainService) DescribeDbbrainDbDiagReportTaskById(ctx context.Context, asyncRequestId *int64, instanceId string, product string) (dbDiagReportTask *dbbrain.HealthReportTask, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dbbrain.NewDescribeDBDiagReportTasksRequest()
	request.InstanceIds = []*string{helper.String(instanceId)}
	request.Product = &product

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeDBDiagReportTasks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if asyncRequestId != nil {
		for _, task := range response.Response.Tasks {
			if *task.AsyncRequestId == *asyncRequestId {
				dbDiagReportTask = task
				return
			}
		}
		return nil, fmt.Errorf("[ERROR]%sThe asyncRequestId[%v] not found in the qurey results. \n", logId, *asyncRequestId)
	}

	dbDiagReportTask = response.Response.Tasks[0]
	return
}

func (me *DbbrainService) DeleteDbbrainDbDiagReportTaskById(ctx context.Context, asyncRequestId int64, instanceId string, product string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dbbrain.NewDeleteDBDiagReportTasksRequest()
	request.AsyncRequestIds = []*int64{helper.Int64(asyncRequestId)}
	request.InstanceId = &instanceId
	request.Product = &product

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DeleteDBDiagReportTasks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DbbrainService) DbbrainDbDiagReportTaskStateRefreshFunc(asyncRequestId *int64, instanceId string, product string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeDbbrainDbDiagReportTaskById(ctx, asyncRequestId, instanceId, product)
		if err != nil {
			return nil, "", err
		}

		return object, helper.Int64ToStr(*object.Progress), nil
	}
}

func (me *DbbrainService) DescribeDbbrainTdsqlAuditLogById(ctx context.Context, asyncRequestId *string, instanceId string, product string) (tdsqlAuditLog []*dbbrain.AuditLogFile, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dbbrain.NewDescribeAuditLogFilesRequest()
	request.InstanceId = &instanceId
	request.Product = &product
	request.NodeRequestType = &product

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeAuditLogFiles(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Items) < 1 {
		return
	}

	if asyncRequestId != nil {
		for _, item := range response.Response.Items {
			if *item.AsyncRequestId == helper.StrToInt64(*asyncRequestId) {
				tdsqlAuditLog = []*dbbrain.AuditLogFile{item}
				return
			}
		}
	}

	tdsqlAuditLog = response.Response.Items
	return
}

func (me *DbbrainService) DeleteDbbrainTdsqlAuditLogById(ctx context.Context, asyncRequestId string, instanceId string, product string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := dbbrain.NewDeleteAuditLogFileRequest()
	request.AsyncRequestId = helper.StrToInt64Point(asyncRequestId)
	request.InstanceId = &instanceId
	request.Product = &product
	request.NodeRequestType = &product

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DeleteAuditLogFile(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DbbrainService) DescribeDbbrainHealthScoresByFilter(ctx context.Context, param map[string]interface{}) (healthScores *dbbrain.HealthScoreInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeHealthScoreRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "instance_id" {
			request.InstanceId = v.(*string)
		}
		if k == "time" {
			request.Time = v.(*string)
		}
		if k == "product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeHealthScore(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response.Data == nil {
		return
	}
	healthScores = response.Response.Data

	return
}

func (me *DbbrainService) DescribeDbbrainSlowLogsByFilter(ctx context.Context, param map[string]interface{}) (SlowLogs []*dbbrain.SlowLogInfoItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeSlowLogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Product" {
			request.Product = v.(*string)
		}
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "Md5" {
			request.Md5 = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
		if k == "Db" {
			request.DB = v.([]*string)
		}
		if k == "Key" {
			request.Key = v.([]*string)
		}
		if k == "User" {
			request.User = v.([]*string)
		}
		if k == "Ip" {
			request.Ip = v.([]*string)
		}
		if k == "Time" {
			request.Time = v.([]*int64)
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
		response, err := me.client.UseDbbrainClient().DescribeSlowLogs(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Rows) < 1 {
			break
		}
		SlowLogs = append(SlowLogs, response.Response.Rows...)
		if len(response.Response.Rows) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DbbrainService) DescribeDbbrainDbSpaceStatusByFilter(ctx context.Context, param map[string]interface{}) (DbSpaceStatus *dbbrain.DescribeDBSpaceStatusResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeDBSpaceStatusRequest()
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
		if k == "RangeDays" {
			request.RangeDays = v.(*int64)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeDBSpaceStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	return response.Response, nil
}

func (me *DbbrainService) DescribeDbbrainSqlTemplatesByFilter(ctx context.Context, param map[string]interface{}) (SqlTemplate *dbbrain.DescribeSqlTemplateResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeSqlTemplateRequest()
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
		if k == "Schema" {
			request.Schema = v.(*string)
		}
		if k == "SqlText" {
			request.SqlText = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeSqlTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	return response.Response, nil
}

func (me *DbbrainService) DescribeDbbrainTopSpaceSchemasByFilter(ctx context.Context, param map[string]interface{}) (TopSpaceSchemas []*dbbrain.SchemaSpaceData, Timestamp *int64, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeTopSpaceSchemasRequest()
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
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().DescribeTopSpaceSchemas(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.TopSpaceSchemas) < 1 {
		return
	}
	TopSpaceSchemas = response.Response.TopSpaceSchemas
	Timestamp = response.Response.Timestamp

	return
}

func (me *DbbrainService) DescribeDbbrainTopSpaceSchemaTimeSeriesByFilter(ctx context.Context, param map[string]interface{}) (TopSpaceSchemaTimeSeries []*dbbrain.SchemaSpaceTimeSeries, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeTopSpaceSchemaTimeSeriesRequest()
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
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "StartDate" {
			request.StartDate = v.(*string)
		}
		if k == "EndDate" {
			request.EndDate = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().DescribeTopSpaceSchemaTimeSeries(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.TopSpaceSchemaTimeSeries) < 1 {
		return
	}
	TopSpaceSchemaTimeSeries = response.Response.TopSpaceSchemaTimeSeries

	return
}

func (me *DbbrainService) DescribeDbbrainTopSpaceTableTimeSeriesByFilter(ctx context.Context, param map[string]interface{}) (TopSpaceTableTimeSeries []*dbbrain.TableSpaceTimeSeries, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeTopSpaceTableTimeSeriesRequest()
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
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "StartDate" {
			request.StartDate = v.(*string)
		}
		if k == "EndDate" {
			request.EndDate = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDbbrainClient().DescribeTopSpaceTableTimeSeries(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.TopSpaceTableTimeSeries) < 1 {
		return
	}
	TopSpaceTableTimeSeries = response.Response.TopSpaceTableTimeSeries

	return
}

func (me *DbbrainService) DescribeDbbrainTopSpaceTablesByFilter(ctx context.Context, param map[string]interface{}) (TopSpaceTables []*dbbrain.TableSpaceData, Timestamp *int64, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeTopSpaceTablesRequest()
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
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDbbrainClient().DescribeTopSpaceTables(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.TopSpaceTables) < 1 {
		return
	}
	TopSpaceTables = response.Response.TopSpaceTables
	Timestamp = response.Response.Timestamp

	return
}

func (me *DbbrainService) DescribeDbbrainDiagDbInstancesByFilter(ctx context.Context, param map[string]interface{}) (items []*dbbrain.InstanceInfo, DbScanStatus *int64, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeDiagDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "IsSupported" {
			request.IsSupported = v.(*bool)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
		if k == "InstanceNames" {
			request.InstanceNames = v.([]*string)
		}
		if k == "InstanceIds" {
			request.InstanceIds = v.([]*string)
		}
		if k == "Regions" {
			request.Regions = v.([]*string)
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
		response, err := me.client.UseDbbrainClient().DescribeDiagDBInstances(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		items = append(items, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DbbrainService) DescribeDbbrainMysqlProcessListByFilter(ctx context.Context, param map[string]interface{}) (mysqlProcessList []*dbbrain.MySqlProcess, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeMySqlProcessListRequest()
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
		if k == "ID" {
			request.ID = v.(*uint64)
		}
		if k == "User" {
			request.User = v.(*string)
		}
		if k == "Host" {
			request.Host = v.(*string)
		}
		if k == "DB" {
			request.DB = v.(*string)
		}
		if k == "State" {
			request.State = v.(*string)
		}
		if k == "Command" {
			request.Command = v.(*string)
		}
		if k == "Time" {
			request.Time = v.(*uint64)
		}
		if k == "Info" {
			request.Info = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		limit uint64 = 20
	)
	request.Limit = &limit
	response, err := me.client.UseDbbrainClient().DescribeMySqlProcessList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.ProcessList) < 1 {
		return
	}
	mysqlProcessList = response.Response.ProcessList

	return
}

func (me *DbbrainService) DescribeDbbrainNoPrimaryKeyTablesByFilter(ctx context.Context, param map[string]interface{}) (tables []*dbbrain.Table, resp *dbbrain.DescribeNoPrimaryKeyTablesResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeNoPrimaryKeyTablesRequest()
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
		if k == "Date" {
			request.Date = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
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
		response, err := me.client.UseDbbrainClient().DescribeNoPrimaryKeyTables(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.NoPrimaryKeyTables) < 1 {
			break
		}
		tables = append(tables, response.Response.NoPrimaryKeyTables...)
		resp = response.Response
		if len(response.Response.NoPrimaryKeyTables) < int(limit) || *response.Response.NoPrimaryKeyTableRecordCount < limit {
			break
		}

		offset += limit
	}

	return
}

func (me *DbbrainService) DescribeDbbrainRedisTopBigKeysByFilter(ctx context.Context, param map[string]interface{}) (redisTopBigKeys []*dbbrain.RedisKeySpaceData, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeRedisTopBigKeysRequest()
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
		if k == "Date" {
			request.Date = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
		if k == "SortBy" {
			request.SortBy = v.(*string)
		}
		if k == "KeyType" {
			request.KeyType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		limit int64 = 20
	)
	request.Limit = &limit
	response, err := me.client.UseDbbrainClient().DescribeRedisTopBigKeys(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.TopKeys) < 1 {
		return
	}
	redisTopBigKeys = response.Response.TopKeys

	return
}

func (me *DbbrainService) DescribeDbbrainRedisTopKeyPrefixListByFilter(ctx context.Context, param map[string]interface{}) (redisTopKeyPrefixList []*dbbrain.RedisPreKeySpaceData, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = dbbrain.NewDescribeRedisTopKeyPrefixListRequest()
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
		if k == "Date" {
			request.Date = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		limit int64 = 20
	)
	request.Limit = &limit
	response, err := me.client.UseDbbrainClient().DescribeRedisTopKeyPrefixList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Items) < 1 {
		return
	}
	redisTopKeyPrefixList = response.Response.Items

	return
}
