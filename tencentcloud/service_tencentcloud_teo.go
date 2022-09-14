package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TeoService struct {
	client *connectivity.TencentCloudClient
}

func (me *TeoService) DescribeTeoZone(ctx context.Context, zoneId string) (zone *teo.Zone, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeZonesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&teo.AdvancedFilter{
			Name:   helper.String("ZoneId"),
			Values: []*string{&zoneId},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*teo.Zone, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTeoClient().DescribeZones(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Zones) < 1 {
			break
		}
		instances = append(instances, response.Response.Zones...)
		if len(response.Response.Zones) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	zone = instances[0]

	return
}

func (me *TeoService) DeleteTeoZoneById(ctx context.Context, zoneId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteZoneRequest()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteZone(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoDnsRecord(ctx context.Context, zoneId, name string) (dnsRecord *teo.DnsRecord, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeDnsRecordsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneId = &zoneId
	request.Filters = append(
		request.Filters,
		&teo.AdvancedFilter{
			Name:   helper.String("name"),
			Values: []*string{&name},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*teo.DnsRecord, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTeoClient().DescribeDnsRecords(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DnsRecords) < 1 {
			break
		}
		instances = append(instances, response.Response.DnsRecords...)
		if len(response.Response.DnsRecords) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	dnsRecord = instances[0]

	return

}

func (me *TeoService) DeleteTeoDnsRecordById(ctx context.Context, zoneId, dnsRecordId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteDnsRecordsRequest()

	request.ZoneId = &zoneId
	request.DnsRecordIds = []*string{&dnsRecordId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteDnsRecords(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoLoadBalancing(ctx context.Context, zoneId, loadBalancingId string) (loadBalancing *teo.LoadBalancing, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeLoadBalancingRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&teo.AdvancedFilter{
			Name:   helper.String("ZoneId"),
			Values: []*string{&zoneId},
		},
	)
	request.Filters = append(
		request.Filters,
		&teo.AdvancedFilter{
			Name:   helper.String("LoadBalancingId"),
			Values: []*string{&loadBalancingId},
		},
	)

	var offset uint64 = 0
	var pageSize uint64 = 100
	loadBalancings := make([]*teo.LoadBalancing, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTeoClient().DescribeLoadBalancing(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}
		loadBalancings = append(loadBalancings, response.Response.Data...)
		if len(response.Response.Data) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(loadBalancings) < 1 {
		return
	}
	loadBalancing = loadBalancings[0]

	return
}

func (me *TeoService) DeleteTeoLoadBalancingById(ctx context.Context, zoneId string, loadBalancingId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteLoadBalancingRequest()
	request.ZoneId = &zoneId
	request.LoadBalancingId = &loadBalancingId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteLoadBalancing(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoOriginGroup(ctx context.Context, zoneId, originGroupId string) (originGroup *teo.OriginGroup, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeOriginGroupRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&teo.AdvancedFilter{
			Name:   helper.String("ZoneId"),
			Values: []*string{&zoneId},
		},
	)
	request.Filters = append(
		request.Filters,
		&teo.AdvancedFilter{
			Name:   helper.String("OriginGroupId"),
			Values: []*string{&originGroupId},
		},
	)

	var offset uint64 = 0
	var pageSize uint64 = 100
	originGroups := make([]*teo.OriginGroup, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTeoClient().DescribeOriginGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.OriginGroups) < 1 {
			break
		}
		originGroups = append(originGroups, response.Response.OriginGroups...)
		if len(response.Response.OriginGroups) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(originGroups) < 1 {
		return
	}
	originGroup = originGroups[0]

	return
}

func (me *TeoService) DeleteTeoOriginGroupById(ctx context.Context, zoneId, originGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteOriginGroupRequest()
	request.ZoneId = &zoneId
	request.OriginGroupId = &originGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteOriginGroup(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoRuleEngine(ctx context.Context, zoneId, ruleId string) (ruleEngine *teo.RuleItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.ZoneId = &zoneId
	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("RuleId"),
			Values: []*string{&ruleId},
		},
	)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DescribeRules(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil && response.Response.RuleItems != nil {
		for _, v := range response.Response.RuleItems {
			if *v.RuleId == ruleId {
				ruleEngine = v
				return
			}
		}
	}

	return

}

func (me *TeoService) DeleteTeoRuleEngineById(ctx context.Context, zoneId, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteRulesRequest()

	request.ZoneId = &zoneId
	request.RuleIds = []*string{&ruleId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteRules(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoApplicationProxy(ctx context.Context, zoneId, proxyId string) (applicationProxy *teo.ApplicationProxy, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeApplicationProxiesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("ZoneId"),
			Values: []*string{&zoneId},
		},
	)
	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("ProxyId"),
			Values: []*string{&proxyId},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 100
	instances := make([]*teo.ApplicationProxy, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTeoClient().DescribeApplicationProxies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ApplicationProxies) < 1 {
			break
		}
		instances = append(instances, response.Response.ApplicationProxies...)
		if len(response.Response.ApplicationProxies) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}
	applicationProxy = instances[0]

	return
}

func (me *TeoService) DeleteTeoApplicationProxyById(ctx context.Context, zoneId, proxyId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteApplicationProxyRequest()

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteApplicationProxy(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoApplicationProxyRule(ctx context.Context, zoneId, proxyId, ruleId string) (applicationProxyRule *teo.ApplicationProxy, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeApplicationProxiesRequest()
	)

	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("ZoneId"),
			Values: []*string{&zoneId},
		},
	)
	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("ProxyId"),
			Values: []*string{&proxyId},
		},
	)
	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("RuleId"),
			Values: []*string{&ruleId},
		},
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DescribeApplicationProxies(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ApplicationProxies) < 1 {
		return
	}
	applicationProxyRule = response.Response.ApplicationProxies[0]
	return
}

func (me *TeoService) DeleteTeoApplicationProxyRuleById(ctx context.Context, zoneId, proxyId, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := teo.NewDeleteApplicationProxyRuleRequest()

	request.ZoneId = &zoneId
	request.ProxyId = &proxyId
	request.RuleId = &ruleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTeoClient().DeleteApplicationProxyRule(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TeoService) DescribeTeoZoneSetting(ctx context.Context, zoneId string) (zoneSetting *teo.DescribeZoneSettingResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeZoneSettingRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId

	response, err := me.client.UseTeoClient().DescribeZoneSetting(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	zoneSetting = response.Response
	return
}

func (me *TeoService) DescribeTeoSecurityPolicy(ctx context.Context, zoneId, entity string) (securityPolicy *teo.DescribeSecurityPolicyResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeSecurityPolicyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId
	request.Entity = &entity

	response, err := me.client.UseTeoClient().DescribeSecurityPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	securityPolicy = response.Response
	return
}

// TODO
//func (me *TeoService) DescribeTeoHostCertificate(ctx context.Context, zoneId, host string) (hostCertificate *teo.DescribeHostsCertificateResponseParams, errRet error) {
//	var (
//		logId   = getLogId(ctx)
//		request = teo.NewDescribeHostsCertificateRequest()
//	)
//
//	defer func() {
//		if errRet != nil {
//			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
//				logId, "query object", request.ToJsonString(), errRet.Error())
//		}
//	}()
//
//	request.Filters = append(
//		request.Filters,
//		&teo.Filter{
//			Name:    helper.String("ZoneId"),
//			Values: []*string{&zoneId},
//		},
//	)
//	request.Filters = append(
//		request.Filters,
//		&teo.Filter{
//			Key:    helper.String("Host"),
//			Values: []*string{&host},
//		},
//	)
//	ratelimit.Check(request.GetAction())
//
//	var offset int64 = 0
//	var pageSize int64 = 100
//	instances := make([]*teo.hostCertificateInfo, 0)
//
//	for {
//		request.Offset = &offset
//		request.Limit = &pageSize
//		ratelimit.Check(request.GetAction())
//		response, err := me.client.UseTeoClient().DescribeHostsCertificate(request)
//		if err != nil {
//			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
//				logId, request.GetAction(), request.ToJsonString(), err.Error())
//			errRet = err
//			return
//		}
//		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
//			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
//
//		if response == nil || len(response.Response.hostCertificates) < 1 {
//			break
//		}
//		instances = append(instances, response.Response....)
//		if len(response.Response.) < int(pageSize) {
//			break
//		}
//		offset += pageSize
//	}
//
//	if len(instances) < 1 {
//		return
//	}
//	hostCertificate = instances[0]
//
//	return
//
//}

func (me *TeoService) DescribeTeoDnsSec(ctx context.Context, zoneId string) (dnsSec *teo.DescribeDnssecResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeDnssecRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId

	response, err := me.client.UseTeoClient().DescribeDnssec(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	dnsSec = response.Response
	return
}

func (me *TeoService) DescribeTeoDefaultCertificate(ctx context.Context, zoneId, certId string) (defaultCertificate *teo.DefaultServerCertInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeDefaultCertificatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("ZoneId"),
			Values: []*string{&zoneId},
		},
	)
	request.Filters = append(
		request.Filters,
		&teo.Filter{
			Name:   helper.String("CertId"),
			Values: []*string{&certId},
		},
	)

	var offset int64 = 0
	var pageSize int64 = 100
	certificates := make([]*teo.DefaultServerCertInfo, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTeoClient().DescribeDefaultCertificates(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DefaultServerCertInfo) < 1 {
			break
		}
		certificates = append(certificates, response.Response.DefaultServerCertInfo...)
		if len(response.Response.DefaultServerCertInfo) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(certificates) < 1 {
		return
	}
	defaultCertificate = certificates[0]

	return
}

func (me *TeoService) DescribeTeoDdosPolicy(ctx context.Context, zoneId string, policyId int64) (ddosPolicy *teo.DescribeDDoSPolicyResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeDDoSPolicyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId
	request.PolicyId = &policyId

	response, err := me.client.UseTeoClient().DescribeDDoSPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	ddosPolicy = response.Response
	return
}

func (me *TeoService) DescribeZoneDDoSPolicy(ctx context.Context, zoneId string) (ddosPolicy *teo.DescribeZoneDDoSPolicyResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeZoneDDoSPolicyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	if zoneId != "" {
		request.ZoneId = &zoneId
	}

	response, err := me.client.UseTeoClient().DescribeZoneDDoSPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	ddosPolicy = response.Response
	return
}

func (me *TeoService) DescribeAvailablePlans(ctx context.Context) (availablePlans *teo.DescribeAvailablePlansResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeAvailablePlansRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseTeoClient().DescribeAvailablePlans(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	availablePlans = response.Response
	return
}

func (me *TeoService) DescribeTeoRuleEnginePriority(ctx context.Context, zoneId string) (ruleEnginePriority []*teo.RuleItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = teo.NewDescribeRulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.ZoneId = &zoneId

	response, err := me.client.UseTeoClient().DescribeRules(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	ruleEnginePriority = response.Response.RuleItems
	return
}
