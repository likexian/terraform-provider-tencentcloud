package tse

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewTseService(client *connectivity.TencentCloudClient) TseService {
	return TseService{client: client}
}

type TseService struct {
	client *connectivity.TencentCloudClient
}

func (me *TseService) DescribeTseInstanceById(ctx context.Context, instanceId string) (instance *tse.SREInstance, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeSREInstancesRequest()
	filter := &tse.Filter{
		Name:   helper.String("InstanceId"),
		Values: []*string{&instanceId},
	}
	request.Filters = append(request.Filters, filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeSREInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Content) < 1 {
		return
	}

	instance = response.Response.Content[0]
	return
}

func (me *TseService) CheckTseInstanceStatusById(ctx context.Context, instanceId, operate string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := me.DescribeTseInstanceById(ctx, instanceId)
		if e != nil {
			return resource.NonRetryableError(e)
		}

		if operate == "create" {
			if instance == nil {
				return resource.NonRetryableError(fmt.Errorf("instance %s instance not exists", instanceId))
			}

			if *instance.Status == "creating" || *instance.Status == "restarting" {
				return resource.RetryableError(fmt.Errorf("create instance status is %v,start retrying ...", *instance.Status))
			}
			if *instance.Status == "running" {
				return nil
			}
		}

		if operate == "update" {
			if instance == nil {
				return resource.NonRetryableError(fmt.Errorf("instance %s instance not exists", instanceId))
			}

			if *instance.Status == "updating" || *instance.Status == "restarting" {
				return resource.RetryableError(fmt.Errorf("update instance status is %v,start retrying ...", *instance.Status))
			}
			if *instance.Status == "running" {
				return nil
			}
		}

		if operate == "delete" {
			if instance == nil {
				return nil
			}

			if *instance.Status == "destroying" {
				return resource.RetryableError(fmt.Errorf("delete instance status is %v,start retrying ...", *instance.Status))
			}
		}

		return resource.NonRetryableError(fmt.Errorf("instance status is %v,we won't wait for it finish", *instance.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s %s instance fail, reason:%s\n ", logId, operate, err.Error())
		errRet = err
		return
	}

	return
}

func (me *TseService) DeleteTseInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteEngineRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteEngine(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseAccessAddressByFilter(ctx context.Context, param map[string]interface{}) (accessAddress *tse.DescribeSREInstanceAccessAddressResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeSREInstanceAccessAddressRequest()
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
		if k == "VpcId" {
			request.VpcId = v.(*string)
		}
		if k == "SubnetId" {
			request.SubnetId = v.(*string)
		}
		if k == "Workload" {
			request.Workload = v.(*string)
		}
		if k == "EngineRegion" {
			request.EngineRegion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTseClient().DescribeSREInstanceAccessAddress(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	accessAddress = response.Response

	return
}

func (me *TseService) DescribeTseNacosReplicasByFilter(ctx context.Context, param map[string]interface{}) (nacosReplicas []*tse.NacosReplica, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeNacosReplicasRequest()
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
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeNacosReplicas(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Replicas) < 1 {
			break
		}
		nacosReplicas = append(nacosReplicas, response.Response.Replicas...)
		if len(response.Response.Replicas) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseZookeeperReplicasByFilter(ctx context.Context, param map[string]interface{}) (zookeeperReplicas []*tse.ZookeeperReplica, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeZookeeperReplicasRequest()
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
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeZookeeperReplicas(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Replicas) < 1 {
			break
		}
		zookeeperReplicas = append(zookeeperReplicas, response.Response.Replicas...)
		if len(response.Response.Replicas) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseZookeeperServerInterfacesByFilter(ctx context.Context, param map[string]interface{}) (zookeeperServerInterfaces []*tse.ZookeeperServerInterface, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeZookeeperServerInterfacesRequest()
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
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeZookeeperServerInterfaces(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Content) < 1 {
			break
		}
		zookeeperServerInterfaces = append(zookeeperServerInterfaces, response.Response.Content...)
		if len(response.Response.Content) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseNacosServerInterfacesByFilter(ctx context.Context, instanceId string) (nacosServerInterfaces []*tse.NacosServerInterface, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeNacosServerInterfacesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeNacosServerInterfaces(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Content) < 1 {
			break
		}
		nacosServerInterfaces = append(nacosServerInterfaces, response.Response.Content...)
		if len(response.Response.Content) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseGatewayNodesByFilter(ctx context.Context, param map[string]interface{}) (gatewayNodes []*tse.CloudNativeAPIGatewayNode, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeCloudNativeAPIGatewayNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayId" {
			request.GatewayId = v.(*string)
		}
		if k == "GroupId" {
			request.GroupId = v.(*string)
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
		response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayNodes(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Result.NodeList) < 1 {
			break
		}
		gatewayNodes = append(gatewayNodes, response.Response.Result.NodeList...)
		if len(response.Response.Result.NodeList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TseService) DescribeTseGatewayCanaryRulesByFilter(ctx context.Context, param map[string]interface{}) (gatewayCanaryRules *tse.CloudAPIGatewayCanaryRuleList, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeCloudNativeAPIGatewayCanaryRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayId" {
			request.GatewayId = v.(*string)
		}
		if k == "ServiceId" {
			request.ServiceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
	)
	canaryRules := make([]*tse.CloudNativeAPIGatewayCanaryRule, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayCanaryRules(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response.Result == nil || len(response.Response.Result.CanaryRuleList) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		canaryRules = append(canaryRules, response.Response.Result.CanaryRuleList...)
		if len(response.Response.Result.CanaryRuleList) < int(limit) {
			break
		}

		offset += limit
	}

	gatewayCanaryRules = &tse.CloudAPIGatewayCanaryRuleList{
		TotalCount:     &total,
		CanaryRuleList: canaryRules,
	}

	return
}

func (me *TseService) DescribeTseGatewayRoutesByFilter(ctx context.Context, param map[string]interface{}) (gatewayRoutes *tse.KongServiceRouteList, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeCloudNativeAPIGatewayRoutesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayId" {
			request.GatewayId = v.(*string)
		}
		if k == "ServiceName" {
			request.ServiceName = v.(*string)
		}
		if k == "RouteName" {
			request.RouteName = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tse.ListFilter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
	)
	route := make([]*tse.KongRoutePreview, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayRoutes(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response.Result == nil || len(response.Response.Result.RouteList) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		route = append(route, response.Response.Result.RouteList...)
		if len(response.Response.Result.RouteList) < int(limit) {
			break
		}

		offset += limit
	}

	gatewayRoutes = &tse.KongServiceRouteList{
		TotalCount: &total,
		RouteList:  route,
	}

	return
}

func (me *TseService) DescribeTseGatewayServicesByFilter(ctx context.Context, param map[string]interface{}) (gatewayServices *tse.KongServices, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeCloudNativeAPIGatewayServicesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayId" {
			request.GatewayId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tse.ListFilter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
	)
	services := make([]*tse.KongServicePreview, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayServices(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response.Result == nil || len(response.Response.Result.ServiceList) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		services = append(services, response.Response.Result.ServiceList...)
		if len(response.Response.Result.ServiceList) < int(limit) {
			break
		}

		offset += limit
	}

	gatewayServices = &tse.KongServices{
		TotalCount:  &total,
		ServiceList: services,
	}

	return
}

func (me *TseService) DescribeTseCngwServiceById(ctx context.Context, gatewayId, name string) (cngwService *tse.KongServiceDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeOneCloudNativeAPIGatewayServiceRequest()
	request.GatewayId = &gatewayId
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeOneCloudNativeAPIGatewayService(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	cngwService = response.Response.Result
	return
}

func (me *TseService) DeleteTseCngwServiceById(ctx context.Context, gatewayId, name string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteCloudNativeAPIGatewayServiceRequest()
	request.GatewayId = &gatewayId
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteCloudNativeAPIGatewayService(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseCngwServiceRateLimitById(ctx context.Context, gatewayId string, name string) (cngwServiceRateLimit *tse.CloudNativeAPIGatewayRateLimitDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeCloudNativeAPIGatewayServiceRateLimitRequest()
	request.GatewayId = &gatewayId
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayServiceRateLimit(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	cngwServiceRateLimit = response.Response.Result
	return
}

func (me *TseService) DeleteTseCngwServiceRateLimitById(ctx context.Context, gatewayId string, name string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteCloudNativeAPIGatewayServiceRateLimitRequest()
	request.GatewayId = &gatewayId
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteCloudNativeAPIGatewayServiceRateLimit(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseCngwRouteById(ctx context.Context, gatewayId string, serviceID string, routeName string) (cngwRoute *tse.KongRoutePreview, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeCloudNativeAPIGatewayRoutesRequest()
	request.GatewayId = &gatewayId
	request.RouteName = &routeName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayRoutes(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response.Result == nil || len(response.Response.Result.RouteList) < 1 {
		return
	}

	for _, v := range response.Response.Result.RouteList {
		if *v.ServiceID == serviceID {
			cngwRoute = v
			return
		}
	}

	return
}

func (me *TseService) DeleteTseCngwRouteById(ctx context.Context, gatewayId string, routeName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteCloudNativeAPIGatewayRouteRequest()
	request.GatewayId = &gatewayId
	request.Name = &routeName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteCloudNativeAPIGatewayRoute(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseCngwRouteRateLimitById(ctx context.Context, gatewayId, routeId string) (cngwRouteRateLimit *tse.CloudNativeAPIGatewayRateLimitDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeCloudNativeAPIGatewayRouteRateLimitRequest()
	request.GatewayId = &gatewayId
	request.Id = &routeId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayRouteRateLimit(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	cngwRouteRateLimit = response.Response.Result
	return
}

func (me *TseService) DeleteTseCngwRouteRateLimitById(ctx context.Context, gatewayId, routeId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteCloudNativeAPIGatewayRouteRateLimitRequest()
	request.GatewayId = &gatewayId
	request.Id = &routeId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteCloudNativeAPIGatewayRouteRateLimit(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseCngwCanaryRuleById(ctx context.Context, gatewayId string, serviceId string, priority string) (cngwCanaryRule *tse.CloudNativeAPIGatewayCanaryRule, errRet error) {
	logId := tccommon.GetLogId(ctx)

	priorityInt64, err := strconv.ParseInt(priority, 10, 64)
	if err != nil {
		return nil, err
	}

	request := tse.NewDescribeCloudNativeAPIGatewayCanaryRulesRequest()
	request.GatewayId = &gatewayId
	request.ServiceId = &serviceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayCanaryRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil || len(response.Response.Result.CanaryRuleList) < 1 {
		return
	}

	for _, v := range response.Response.Result.CanaryRuleList {
		if *v.Priority == priorityInt64 {
			cngwCanaryRule = v
			return
		}
	}

	return
}

func (me *TseService) DeleteTseCngwCanaryRuleById(ctx context.Context, gatewayId string, serviceId string, priority string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	priorityInt64, err := strconv.ParseInt(priority, 10, 64)
	if err != nil {
		return err
	}

	request := tse.NewDeleteCloudNativeAPIGatewayCanaryRuleRequest()
	request.GatewayId = &gatewayId
	request.ServiceId = &serviceId
	request.Priority = &priorityInt64

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteCloudNativeAPIGatewayCanaryRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseCngwGatewayById(ctx context.Context, gatewayId string) (cngwGateway *tse.DescribeCloudNativeAPIGatewayResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeCloudNativeAPIGatewayRequest()
	request.GatewayId = &gatewayId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeCloudNativeAPIGateway(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Result == nil {
		return
	}

	cngwGateway = response.Response.Result
	return
}

func (me *TseService) DeleteTseCngwGatewayById(ctx context.Context, gatewayId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteCloudNativeAPIGatewayRequest()
	request.GatewayId = &gatewayId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteCloudNativeAPIGateway(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) CheckTseNativeAPIGatewayStatusById(ctx context.Context, gatewayId, operate string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		gateway, e := me.DescribeTseCngwGatewayById(ctx, gatewayId)
		if e != nil && operate != "delete" {
			return resource.NonRetryableError(e)
		}

		if operate == "create" {
			if gateway == nil {
				return resource.NonRetryableError(fmt.Errorf("gateway %s not exists", gatewayId))
			}

			if *gateway.Status == "Creating" || *gateway.Status == "InstallingPlugin" || *gateway.Status == "CreatedWithNoPlugin" {
				return resource.RetryableError(fmt.Errorf("create gateway status is %v,start retrying ...", *gateway.Status))
			}
			if *gateway.Status == "Running" {
				return nil
			}
		}

		if operate == "update" {
			if gateway == nil {
				return resource.NonRetryableError(fmt.Errorf("gateway %s not exists", gatewayId))
			}

			if *gateway.Status == "Modifying" {
				return resource.RetryableError(fmt.Errorf("update gateway status is %v,start retrying ...", *gateway.Status))
			}
			if *gateway.Status == "Running" {
				return nil
			}
		}

		if operate == "delete" {
			if e != nil {
				if sdkErr, ok := e.(*sdkError.TencentCloudSDKError); ok {
					if sdkErr.Code == "ResourceNotFound.InstanceNotFound" {
						return nil
					}
				}
			}
			if gateway == nil {
				return nil
			}

			if *gateway.Status == "Deleting" {
				return resource.RetryableError(fmt.Errorf("delete gateway status is %v,start retrying ...", *gateway.Status))
			}
		}
		return resource.NonRetryableError(fmt.Errorf("gateway status is %v,we won't wait for it finish", *gateway.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s %s gateway fail, reason:%s\n ", logId, operate, err.Error())
		errRet = err
		return
	}

	return
}

func (me *TseService) DescribeTseCngwGroupById(ctx context.Context, gatewayId string, groupId string) (cngwGroup *tse.NativeGatewayServerGroup, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeNativeGatewayServerGroupsRequest()
	request.GatewayId = &gatewayId
	request.Filters = append(
		request.Filters,
		&tse.Filter{
			Name:   helper.String("GroupId"),
			Values: []*string{&groupId},
		},
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeNativeGatewayServerGroups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.Result == nil || len(response.Response.Result.GatewayGroupList) < 1 {
		return
	}

	cngwGroup = response.Response.Result.GatewayGroupList[0]
	return
}

func (me *TseService) DeleteTseCngwGroupById(ctx context.Context, gatewayId string, groupId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteNativeGatewayServerGroupRequest()
	request.GatewayId = &gatewayId
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteNativeGatewayServerGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseGatewaysByFilter(ctx context.Context, param map[string]interface{}) (gateways *tse.ListCloudNativeAPIGatewayResult, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeCloudNativeAPIGatewaysRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*tse.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
	)
	gateway := make([]*tse.DescribeCloudNativeAPIGatewayResult, 0)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeCloudNativeAPIGateways(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.GatewayList) < 1 {
			break
		}
		total = *response.Response.Result.TotalCount
		gateway = append(gateway, response.Response.Result.GatewayList...)
		if len(response.Response.Result.GatewayList) < int(limit) {
			break
		}

		offset += limit
	}

	gateways = &tse.ListCloudNativeAPIGatewayResult{
		TotalCount:  &total,
		GatewayList: gateway,
	}

	return
}
func (me *TseService) DescribeTseGroupsByFilter(ctx context.Context, param map[string]interface{}) (groups *tse.NativeGatewayServerGroups, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeNativeGatewayServerGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayId" {
			request.GatewayId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tse.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
		total  uint64
	)
	group := make([]*tse.NativeGatewayServerGroup, 0)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeNativeGatewayServerGroups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.GatewayGroupList) < 1 {
			break
		}

		total = *response.Response.Result.TotalCount
		group = append(group, response.Response.Result.GatewayGroupList...)
		if len(response.Response.Result.GatewayGroupList) < int(limit) {
			break
		}

		offset += limit
	}

	groups = &tse.NativeGatewayServerGroups{
		TotalCount:       &total,
		GatewayGroupList: group,
	}

	return
}

func (me *TseService) CheckTseNativeAPIGatewayGroupStatusById(ctx context.Context, gatewayId, groupId, operate string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	time.Sleep(5 * time.Second)
	err := resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		gateway, e := me.DescribeTseCngwGroupById(ctx, gatewayId, groupId)
		if e != nil && operate != "delete" {
			return resource.NonRetryableError(e)
		}

		if operate == "create" {
			if gateway == nil {
				return resource.NonRetryableError(fmt.Errorf("group %s not exists", groupId))
			}

			if *gateway.Status == "Creating" {
				return resource.RetryableError(fmt.Errorf("create group status is %v,start retrying ...", *gateway.Status))
			}
			if *gateway.Status == "Running" {
				return nil
			}
		}

		if operate == "update" {
			if gateway == nil {
				return resource.NonRetryableError(fmt.Errorf("group %s not exists", groupId))
			}

			if *gateway.Status == "Modifying" || *gateway.Status == "UpdatingSpec" {
				return resource.RetryableError(fmt.Errorf("update group status is %v,start retrying ...", *gateway.Status))
			}
			if *gateway.Status == "Running" {
				return nil
			}
		}

		if operate == "delete" {
			if e != nil {
				if sdkErr, ok := e.(*sdkError.TencentCloudSDKError); ok {
					if sdkErr.Code == "ResourceNotFound.InstanceNotFound" {
						return nil
					}
				}
			}
			if gateway == nil {
				return nil
			}

			if *gateway.Status == "Running" || *gateway.Status == "Deleting" {
				return resource.RetryableError(fmt.Errorf("delete group status is %v,start retrying ...", *gateway.Status))
			}
		}
		return resource.NonRetryableError(fmt.Errorf("group status is %v,we won't wait for it finish", *gateway.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s %s group fail, reason:%s\n ", logId, operate, err.Error())
		errRet = err
		return
	}

	return
}

func (me *TseService) DescribeTseGatewayCertificatesByFilter(ctx context.Context, param map[string]interface{}) (gatewayCertificates *tse.KongCertificatesList, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tse.NewDescribeCloudNativeAPIGatewayCertificatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "GatewayId" {
			request.GatewayId = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*tse.ListFilter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
		total  int64
	)
	certificates := make([]*tse.KongCertificatesPreview, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayCertificates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.CertificatesList) < 1 {
			break
		}

		total = *response.Response.Result.Total
		certificates = append(certificates, response.Response.Result.CertificatesList...)
		if len(response.Response.Result.CertificatesList) < int(limit) {
			break
		}

		offset += limit
	}

	gatewayCertificates = &tse.KongCertificatesList{
		Total:            &total,
		CertificatesList: certificates,
	}

	return
}

func (me *TseService) DescribeTseCngwCertificateById(ctx context.Context, gatewayId string, id string) (cngwCertificate *tse.KongCertificatesPreview, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeCloudNativeAPIGatewayCertificateDetailsRequest()
	request.GatewayId = &gatewayId
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeCloudNativeAPIGatewayCertificateDetails(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response.Result == nil || response.Response.Result.Cert == nil {
		return
	}

	cngwCertificate = response.Response.Result.Cert
	return
}

func (me *TseService) DeleteTseCngwCertificateById(ctx context.Context, gatewayId string, id string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteCloudNativeAPIGatewayCertificateRequest()
	request.GatewayId = &gatewayId
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteCloudNativeAPIGatewayCertificate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseWafProtectionById(ctx context.Context, gatewayId string) (wafProtection *tse.DescribeWafProtectionResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeWafProtectionRequest()
	request.GatewayId = &gatewayId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeWafProtection(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	wafProtection = response.Response.Result
	return
}

func (me *TseService) DescribeTseWafDomainsById(ctx context.Context, gatewayId string) (wafDomains *tse.DescribeWafDomainsResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeWafDomainsRequest()
	request.GatewayId = &gatewayId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeWafDomains(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	wafDomains = response.Response.Result
	return
}

func (me *TseService) DeleteTseWafDomainsById(ctx context.Context, gatewayId string, domain string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteWafDomainsRequest()
	request.GatewayId = &gatewayId
	request.Domains = []*string{&domain}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteWafDomains(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) DescribeTseCngwStrategyBindGroupById(ctx context.Context, gatewayId, strategyId, groupId string) (cngwStrategyBindGroup *tse.CloudNativeAPIGatewayStrategyBindingGroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeAutoScalerResourceStrategyBindingGroupsRequest()
	request.GatewayId = &gatewayId
	request.StrategyId = &strategyId

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
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTseClient().DescribeAutoScalerResourceStrategyBindingGroups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if response == nil || response.Response == nil || response.Response.Result == nil || len(response.Response.Result.GroupInfos) < 1 {
			break
		}

		for _, v := range response.Response.Result.GroupInfos {
			if *v.GroupId == groupId {
				cngwStrategyBindGroup = v
				return
			}
		}

		offset += limit
	}

	return nil, nil
}

func (me *TseService) DescribeTseCngwNetworkById(ctx context.Context, gatewayId, groupId, networkId string) (cngwNetwork *tse.DescribePublicNetworkResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribePublicNetworkRequest()
	request.GatewayId = &gatewayId
	request.GroupId = &groupId
	request.NetworkId = &networkId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribePublicNetwork(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	cngwNetwork = response.Response.Result
	return
}

func (me *TseService) DeleteTseCngwNetworkById(ctx context.Context, gatewayId, groupId, networkId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteCloudNativeAPIGatewayPublicNetworkRequest()
	request.GatewayId = &gatewayId
	request.GroupId = &groupId

	object, err := me.DescribeTseCngwNetworkById(ctx, gatewayId, groupId, networkId)
	if err != nil {
		return err
	}

	if object == nil || object.PublicNetwork == nil {
		return nil
	}

	request.Vip = object.PublicNetwork.Vip

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteCloudNativeAPIGatewayPublicNetwork(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TseService) TseCngwNetworkStateRefreshFunc(gatewayId, groupId, networkId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		object, err := me.DescribeTseCngwNetworkById(ctx, gatewayId, groupId, networkId)
		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object.PublicNetwork.Status), nil
	}
}

func (me *TseService) DescribeTseCngwStrategyById(ctx context.Context, gatewayId string, strategyId string) (cngwStrategy *tse.CloudNativeAPIGatewayStrategy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDescribeAutoScalerResourceStrategiesRequest()
	request.GatewayId = &gatewayId
	request.StrategyId = &strategyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DescribeAutoScalerResourceStrategies(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.Result == nil || response.Response.Result.StrategyList == nil {
		return
	}

	cngwStrategy = response.Response.Result.StrategyList[0]
	return
}
func (me *TseService) DeleteTseCngwStrategyById(ctx context.Context, gatewayId string, strategyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tse.NewDeleteAutoScalerResourceStrategyRequest()
	request.GatewayId = &gatewayId
	request.StrategyId = &strategyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTseClient().DeleteAutoScalerResourceStrategy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
