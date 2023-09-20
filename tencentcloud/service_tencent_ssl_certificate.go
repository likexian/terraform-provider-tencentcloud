package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"math"

	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SSLService struct {
	client *connectivity.TencentCloudClient
}

func (me *SSLService) ApplyCertificate(ctx context.Context, request *ssl.ApplyCertificateRequest) (id string, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSSLCertificateClient().ApplyCertificate(request)

	if err != nil {
		errRet = err
		return
	}

	if response.Response.CertificateId != nil {
		id = *response.Response.CertificateId
	} else {
		errRet = fmt.Errorf("[%s] error, no certificate id response: %s", request.GetAction(), response.ToJsonString())
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SSLService) CreateCertificate(ctx context.Context, request *ssl.CreateCertificateRequest) (certificateId, dealId string, errRet error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	response, err := client.CreateCertificate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		if len(response.Response.CertificateIds) > 0 {
			certificateId = *response.Response.CertificateIds[0]
		}
		if len(response.Response.DealIds) > 0 {
			dealId = *response.Response.DealIds[0]
		}
		return
	}
	errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
	return
}

func (me *SSLService) CommitCertificateInformation(ctx context.Context, request *ssl.CommitCertificateInformationRequest) (errRet error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	response, err := client.CommitCertificateInformation(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *SSLService) DescribeCertificateDetail(ctx context.Context, request *ssl.DescribeCertificateDetailRequest) (response *ssl.DescribeCertificateDetailResponse, err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	response, err = client.DescribeCertificateDetail(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	return
}

func (me *SSLService) ModifyCertificateAlias(ctx context.Context, request *ssl.ModifyCertificateAliasRequest) (err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.ModifyCertificateAliasResponse

	response, err = client.ModifyCertificateAlias(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	return
}

func (me *SSLService) ModifyCertificateProject(ctx context.Context, request *ssl.ModifyCertificateProjectRequest) (err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.ModifyCertificateProjectResponse

	response, err = client.ModifyCertificateProject(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	for _, v := range response.Response.FailCertificates {
		if *v == *request.CertificateIdList[0] {
			err = fmt.Errorf("failed to modify the project. certificateId=%s", *request.CertificateIdList[0])
			return
		}
	}
	return
}

func (me *SSLService) DeleteCertificate(ctx context.Context, request *ssl.DeleteCertificateRequest) (deleteResult bool, err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.DeleteCertificateResponse

	response, err = client.DeleteCertificate(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	deleteResult = *response.Response.DeleteResult
	return
}

func (me *SSLService) CancelCertificateOrder(ctx context.Context, request *ssl.CancelCertificateOrderRequest) (err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.CancelCertificateOrderResponse

	response, err = client.CancelCertificateOrder(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	return
}

func (me *SSLService) SubmitCertificateInformation(ctx context.Context, request *ssl.SubmitCertificateInformationRequest) (err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.SubmitCertificateInformationResponse

	response, err = client.SubmitCertificateInformation(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	return
}

func (me *SSLService) UploadConfirmLetter(ctx context.Context, request *ssl.UploadConfirmLetterRequest) (err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.UploadConfirmLetterResponse

	response, err = client.UploadConfirmLetter(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	return
}

func (me *SSLService) UploadCertificate(ctx context.Context, request *ssl.UploadCertificateRequest) (id string, err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	var response *ssl.UploadCertificateResponse
	response, err = client.UploadCertificate(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, request.GetAction(), request.ToJsonString(), err)
		return
	}

	if response == nil || response.Response == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}

	if response.Response.CertificateId == nil {
		err = fmt.Errorf("api[%s] return id is nil", request.GetAction())
		log.Printf("[CRITAL]%s %v", logId, err)
		return
	}

	id = *response.Response.CertificateId
	return
}

func (me *SSLService) DescribeCertificates(ctx context.Context, request *ssl.DescribeCertificatesRequest) (certificateList []*ssl.Certificates, err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()

	offset := 0
	pageSize := 100
	certificateList = make([]*ssl.Certificates, 0)
	var response *ssl.DescribeCertificatesResponse
	for {
		request.Offset = helper.IntUint64(offset)
		request.Limit = helper.IntUint64(pageSize)
		ratelimit.Check(request.GetAction())
		response, err = client.DescribeCertificates(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || len(response.Response.Certificates) == 0 {
			break
		}
		certificateList = append(certificateList, response.Response.Certificates...)
		if len(response.Response.Certificates) < pageSize {
			break
		}
		offset += pageSize
	}

	return
}

func (me *SSLService) checkCertificateType(ctx context.Context, certId string, checkType string) (bool, error) {

	//get certificate by id

	request := ssl.NewDescribeCertificateDetailRequest()
	request.CertificateId = helper.String(certId)
	certificate, err := me.DescribeCertificateDetail(ctx, request)
	if err != nil {
		return false, err
	}

	if certificate != nil && certificate.Response != nil && *certificate.Response.CertificateType == checkType {
		return true, nil
	} else {
		if certificate == nil || certificate.Response == nil || certificate.Response.CertificateId == nil {
			return false, fmt.Errorf("certificate id %s is not found", certId)
		}
		return false, nil
	}

}
func (me *SSLService) ModifyCertificateResubmit(ctx context.Context, request *ssl.ModifyCertificateResubmitRequest) (err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()
	ratelimit.Check(request.GetAction())

	response, err := client.ModifyCertificateResubmit(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil || response.Response.CertificateId == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if *response.Response.CertificateId != *request.CertificateId {
		err = fmt.Errorf("TencentCloud SDK %s eertificates are inconsistent, request[%s], response[%s]",
			request.GetAction(), *request.CertificateId, *response.Response.CertificateId)
		return
	}
	return
}
func (me *SSLService) CancelAuditCertificate(ctx context.Context, request *ssl.CancelAuditCertificateRequest) (err error) {
	logId := getLogId(ctx)
	client := me.client.UseSSLCertificateClient()

	response, err := client.CancelAuditCertificate(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return
	}
	if response == nil || response.Response == nil || response.Response.Result == nil {
		err = fmt.Errorf("TencentCloud SDK %s return empty response", request.GetAction())
		return
	}
	if !*response.Response.Result {
		err = fmt.Errorf("TencentCloud SDK %s CancelAudit failed", request.GetAction())
		return err
	}

	return
}
func (me *SSLService) getCertificateStatus(ctx context.Context, certificateId string) (uint64, error) {
	describeRequest := ssl.NewDescribeCertificateDetailRequest()
	describeRequest.CertificateId = &certificateId

	describeResponse, err := me.DescribeCertificateDetail(ctx, describeRequest)
	if err != nil {
		return math.MaxUint64, err
	}
	if describeResponse == nil || describeResponse.Response == nil {
		err := fmt.Errorf("TencentCloud SDK %s return empty response", describeRequest.GetAction())
		return math.MaxUint64, err
	}
	if describeResponse.Response.Status == nil {
		err := fmt.Errorf("api[%s] certificate status is nil", describeRequest.GetAction())
		return math.MaxUint64, err
	}

	return *describeResponse.Response.Status, nil
}
