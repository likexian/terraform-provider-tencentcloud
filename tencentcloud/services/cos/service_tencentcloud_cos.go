package cos

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/beevik/etree"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewCosService(client *connectivity.TencentCloudClient) CosService {
	return CosService{client: client}
}

type CosService struct {
	client *connectivity.TencentCloudClient
}

type CosBucketDomainCertItem struct {
	bucket     string
	domainName string
}

const (
	CERT_ENABLED  = "Enabled"
	CERT_DISABLED = "Disabled"
)

const PUBLIC_GRANTEE = "http://cam.qcloud.com/groups/global/AllUsers"

func (me *CosService) HeadObject(ctx context.Context, bucket, key string) (info *s3.HeadObjectOutput, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	ratelimit.Check("HeadObject")
	response, err := me.client.UseCosClient().HeadObject(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "head object", request.String(), err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "head object", request.String(), response.String())

	return response, nil
}

func (me *CosService) DeleteObject(ctx context.Context, bucket, key string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.String(), errRet.Error())
		}
	}()
	ratelimit.Check("DeleteObject")
	response, err := me.client.UseCosClient().DeleteObject(&request)
	if err != nil {
		errRet = fmt.Errorf("cos delete object error: %s, bucket: %s, object: %s", err.Error(), bucket, key)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "delete object", request.String(), response.String())

	return nil
}

func (me *CosService) PutObjectAcl(ctx context.Context, bucket, key, acl string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.PutObjectAclInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    aws.String(acl),
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put object acl", request.String(), errRet.Error())
		}
	}()
	ratelimit.Check("PutObjectAcl")
	response, err := me.client.UseCosClient().PutObjectAcl(&request)
	if err != nil {
		errRet = fmt.Errorf("cos put object acl error: %s, bucket: %s, object: %s", err.Error(), bucket, key)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put object acl", request.String(), response.String())

	return nil
}

// PutBucket - base on aws s3
func (me *CosService) PutBucket(ctx context.Context, bucket, acl string, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket", request.String(), errRet.Error())
		}
	}()
	ratelimit.Check("CreateBucket")
	response, err := me.client.UseCosClientNew(cdcId).CreateBucket(&request)

	if err != nil {
		errRet = fmt.Errorf("cos put bucket error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s], endpoint %s\n",
		logId, "put bucket", request.String(), response.String(), me.client.UseCosClientNew(cdcId).Endpoint)

	return nil
}

// TencentCosPutBucket - To support MAZ config, We use tencentcloud cos sdk instead of aws s3
func (me *CosService) TencentCosPutBucket(ctx context.Context, bucket string, opt *cos.BucketPutOptions, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	req, _ := json.Marshal(opt)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request [%s], reason[%s]\n",
				logId, "put bucket", req, errRet.Error())
		}
	}()

	ratelimit.Check("TencentcloudCosPutBucket")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.Put(ctx, opt)

	if err != nil {
		errRet = fmt.Errorf("cos put bucket error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	resp, _ := json.Marshal(response.Response.Body)

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s], baseUrl %s\n",
		logId, "put bucket", req, resp, me.client.UseTencentCosClientNew(bucket, cdcId).BaseURL.BucketURL)

	return nil
}

func (me *CosService) TencentCosBucketGetLocation(ctx context.Context, bucket string, opt *cos.BucketPutOptions) (location string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	req, _ := json.Marshal(opt)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request [%s], reason[%s]\n",
				logId, "put bucket", req, errRet.Error())
		}
	}()

	ratelimit.Check("TencentcloudCosPutBucket")
	result, response, err := me.client.UseTencentCosClient(bucket).Bucket.GetLocation(ctx)
	if err != nil {
		errRet = fmt.Errorf("cos get location error: %s, bucket: %s", err.Error(), bucket)
		return
	}
	location = result.Location
	resp, _ := json.Marshal(response.Response.Body)

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s], baseUrl %s\n",
		logId, "put bucket", req, resp, me.client.UseTencentCosClient(bucket).BaseURL.BucketURL)

	return
}

func (me *CosService) TencentCosPutBucketACLBody(
	ctx context.Context,
	bucket string,
	reqBody string,
	header string,
	cdcId string,
) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	acl := &cos.ACLXml{}
	opt := &cos.BucketPutACLOptions{}
	if cdcId == "" && me.client.CosDomain == "" {
		if reqBody != "" {
			err := xml.Unmarshal([]byte(reqBody), acl)
			if err != nil {
				errRet = fmt.Errorf("cos [PutBucketACLBody] XML Unmarshal error: %s, bucket: %s", err.Error(), bucket)
				return
			}

			opt.Body = acl
		} else if header != "" {
			opt.Header = &cos.ACLHeaderOptions{
				XCosACL: header,
			}
		}
	} else {
		err := xml.Unmarshal([]byte(reqBody), acl)
		if err != nil {
			errRet = fmt.Errorf("cos [PutBucketACLBody] XML Unmarshal error: %s, bucket: %s", err.Error(), bucket)
			return
		}

		var (
			uin         string
			fullControl string
			read        string
			write       string
			readAcp     string
			writeAcp    string
		)

		for _, v := range acl.AccessControlList {
			tmpList := regexp.MustCompile(`\d+`).FindAllString(v.Grantee.ID, 1)
			if len(tmpList) > 0 {
				uin = tmpList[0]
			}

			if v.Permission == "FULL_CONTROL" {
				fullControl = fmt.Sprintf("id=\"%s\"", uin)
			}

			if v.Permission == "READ" {
				read = fmt.Sprintf("id=\"%s\"", uin)
			}

			if v.Permission == "WRITE" {
				write = fmt.Sprintf("id=\"%s\"", uin)
			}

			if v.Permission == "READ_ACP" {
				readAcp = fmt.Sprintf("id=\"%s\"", uin)
			}

			if v.Permission == "WRITE_ACP" {
				writeAcp = fmt.Sprintf("id=\"%s\"", uin)
			}
		}

		opt.Header = &cos.ACLHeaderOptions{
			XCosACL:              header,
			XCosGrantFullControl: fullControl,
			XCosGrantRead:        read,
			XCosGrantWrite:       write,
			XCosGrantReadACP:     readAcp,
			XCosGrantWriteACP:    writeAcp,
		}
	}

	req, _ := json.Marshal(opt)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request [%s], reason[%s]\n",
				logId, "PutBucketACLBody", req, errRet.Error())
		}
	}()

	ratelimit.Check("TencentcloudCosPutBucketACLBody")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.PutACL(ctx, opt)

	if err != nil {
		errRet = fmt.Errorf("cos [PutBucketACLBody] error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	resp, _ := json.Marshal(response.Response.Body)

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "PutBucketACLBody", req, resp)

	return nil
}

func (me *CosService) HeadBucket(ctx context.Context, bucket string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("HeadBucket")
	response, err := me.client.UseCosClient().HeadBucket(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "head bucket", request.String(), err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "head bucket", request.String(), response.String())

	return nil
}

func (me *CosService) TencentcloudHeadBucket(ctx context.Context, bucket string, cdcId string) (code int, header http.Header, errRet error) {
	logId := tccommon.GetLogId(ctx)

	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.Head(ctx)

	if response != nil {
		code = response.StatusCode
		header = response.Header
	}

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
			logId, "HeadBucket", err.Error())
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success\n",
		logId, "HeadBucket")

	return
}

func (me *CosService) DeleteBucket(ctx context.Context, bucket string, forced bool, versioned bool, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	if forced {
		log.Printf("[DEBUG]%s api[%s] triggered, bucket [%s], versioned [%v]\n",
			logId, "ForceCleanObject", bucket, versioned)
		err := me.ForceCleanObject(ctx, bucket, versioned, cdcId)
		if err != nil {
			return err
		}
	}

	request := s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("DeleteBucket")
	response, err := me.client.UseCosClientNew(cdcId).DeleteBucket(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "delete bucket", request.String(), err.Error())
		return fmt.Errorf("cos delete bucket error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "delete bucket", request.String(), response.String())

	return nil
}

func (me *CosService) ForceCleanObject(ctx context.Context, bucket string, versioned bool, cdcId string) error {
	logId := tccommon.GetLogId(ctx)

	// Get the object list of bucket with all versions
	verOpt := cos.BucketGetObjectVersionsOptions{}
	objList, resp, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.GetObjectVersions(ctx, &verOpt)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, resp body [%s], reason[%s]\n",
			logId, "GetObjectVersions", resp.Body, err.Error())
		return fmt.Errorf("cos force clean object error: %s, bucket: %s", err.Error(), bucket)
	}
	if objList.IsTruncated {
		return fmt.Errorf("cos force clean object error: the list of objects is truncated and the bucket[%s] needs to be deleted manually!!!", bucket)
	}

	verCnt := len(objList.Version)
	markerCnt := len(objList.DeleteMarker)
	log.Printf("[DEBUG][ForceCleanObject]%s api[%s] success, get [%v] versions of object, get [%v] deleteMarker, versioned[%v].\n", logId, "GetObjectVersions", verCnt, markerCnt, versioned)

	delCnt := verCnt + markerCnt
	if delCnt == 0 {
		return nil
	}

	delObjs := make([]cos.Object, 0, delCnt)
	if versioned {
		//add the versions
		for _, v := range objList.Version {
			delObjs = append(delObjs, cos.Object{
				Key:       v.Key,
				VersionId: v.VersionId,
			})
		}
		// add the delete-marker
		for _, m := range objList.DeleteMarker {
			delObjs = append(delObjs, cos.Object{
				Key:       m.Key,
				VersionId: m.VersionId,
			})
		}
	} else {
		for _, v := range objList.Version {
			delObjs = append(delObjs, cos.Object{
				Key: v.Key,
			})
		}
	}

	opt := cos.ObjectDeleteMultiOptions{
		Quiet:   true,
		Objects: delObjs,
	}

	// Multi-delete by specified object.
	result, resp, err := me.client.UseTencentCosClientNew(bucket, cdcId).Object.DeleteMulti(ctx, &opt)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, resp body [%s], reason[%s], opt[%v]\n",
			logId, "DeleteMulti", resp.Body, err.Error(), opt)
		return fmt.Errorf("cos force clean object error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG][ForceCleanObject]%s api[%s] completed, removed [%v] versions of object. [%v] failed to remove.\n",
		logId, "DeleteMulti", len(result.DeletedObjects), len(result.Errors))

	// Clean the failed removal version.
	if len(result.Errors) > 0 {
		log.Printf("[CRITAL]%s api[%s] it still [%v] objects have not been removed, need try DeleteMulti again.\n",
			logId, "DeleteMulti", len(result.Errors))

		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			unDelObjs := make([]cos.Object, 0, len(result.Errors))
			for _, v := range result.Errors {
				unDelObjs = append(unDelObjs, cos.Object{
					Key:       v.Key,
					VersionId: v.VersionId,
				})
			}
			unDelOpt := cos.ObjectDeleteMultiOptions{
				Quiet:   true,
				Objects: unDelObjs,
			}

			result, resp, err := me.client.UseTencentCosClientNew(bucket, cdcId).Object.DeleteMulti(ctx, &unDelOpt)
			if err != nil {
				log.Printf("[CRITAL][retry]%s api[%s] fail, resp body [%s], reason[%s]\n",
					logId, "DeleteMulti ", resp.Body, err.Error())
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			if len(result.Errors) > 0 {
				return resource.RetryableError(fmt.Errorf("[CRITAL][retry]%s api[%s] it still %v objects have not been removed, need try DeleteMulti again.\n",
					logId, "DeleteMulti", len(result.Errors)))
			}
			return nil
		}); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG][ForceCleanObject]%s api[%s] success, [%v] objects have been cleaned.\n",
		logId, "ForceCleanObject", len(result.DeletedObjects))
	return nil
}

func (me *CosService) GetBucketCors(ctx context.Context, bucket string, cdcId string) (corsRules []map[string]interface{}, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.GetBucketCorsInput{
		Bucket: aws.String(bucket),
	}

	ratelimit.Check("GetBucketCors")
	response, err := me.client.UseCosClientNew(cdcId).GetBucketCors(&request)

	if err != nil {
		awsError, ok := err.(awserr.Error)
		if !ok || awsError.Code() != "NoSuchCORSConfiguration" {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "get bucket cors", request.String(), err.Error())
			errRet = fmt.Errorf("cos get bucket cors error: %s, bucket: %s", err.Error(), bucket)
			return
		}
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket cors", request.String(), response.String())

	corsRules = make([]map[string]interface{}, 0, len(response.CORSRules))
	if len(response.CORSRules) > 0 {
		for _, value := range response.CORSRules {
			rule := make(map[string]interface{})
			rule["allowed_origins"] = helper.StringsInterfaces(value.AllowedOrigins)
			rule["allowed_methods"] = helper.StringsInterfaces(value.AllowedMethods)
			rule["allowed_headers"] = helper.StringsInterfaces(value.AllowedHeaders)

			if value.ExposeHeaders != nil {
				rule["expose_headers"] = helper.StringsInterfaces(value.ExposeHeaders)
			}
			if value.MaxAgeSeconds != nil {
				rule["max_age_seconds"] = int(*value.MaxAgeSeconds)
			}

			corsRules = append(corsRules, rule)
		}
	}
	return
}

func (me *CosService) GetBucketLifecycle(ctx context.Context, bucket string, cdcId string) (lifecycleRules []map[string]interface{}, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("GetBucketLifecycleConfiguration")
	response, err := me.client.UseCosClientNew(cdcId).GetBucketLifecycleConfiguration(&request)

	if err != nil {
		awsError, ok := err.(awserr.Error)
		if !ok || awsError.Code() != "NoSuchLifecycleConfiguration" {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "get bucket lifecycle", request.String(), err.Error())
			errRet = fmt.Errorf("cos get bucket cors error: %s, bucket: %s", err.Error(), bucket)
			return
		}
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket lifecycle", request.String(), response.String())

	lifecycleRules = make([]map[string]interface{}, 0, len(response.Rules))
	if len(response.Rules) > 0 {
		for _, value := range response.Rules {
			rule := make(map[string]interface{})

			if value.ID != nil {
				rule["id"] = *value.ID
			}

			// filter_prefix
			if value.Filter != nil {
				if value.Filter.And != nil && value.Filter.And.Prefix != nil &&
					*value.Filter.And.Prefix != "" {
					rule["filter_prefix"] = *value.Filter.And.Prefix
				} else if value.Filter.Prefix != nil && *value.Filter.Prefix != "" {
					rule["filter_prefix"] = *value.Filter.Prefix
				}
			}
			// transition
			if len(value.Transitions) > 0 {
				transitions := make([]interface{}, 0, len(value.Transitions))
				for _, v := range value.Transitions {
					t := make(map[string]interface{})
					if v.Date != nil {
						t["date"] = (*v.Date).Format("2006-01-02")
					}
					if v.Days != nil {
						t["days"] = int(*v.Days)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["transition"] = schema.NewSet(transitionHash, transitions)
			}
			// expiration
			if value.Expiration != nil {
				e := make(map[string]interface{})
				if value.Expiration.Date != nil {
					e["date"] = (*value.Expiration.Date).Format("2006-01-02")
				}
				if value.Expiration.Days != nil {
					e["days"] = int(*value.Expiration.Days)
				}
				if value.Expiration.ExpiredObjectDeleteMarker != nil {
					e["delete_marker"] = *value.Expiration.ExpiredObjectDeleteMarker
				}
				rule["expiration"] = schema.NewSet(expirationHash, []interface{}{e})
			}

			// transition
			if len(value.NoncurrentVersionTransitions) > 0 {
				transitions := make([]interface{}, 0, len(value.NoncurrentVersionTransitions))
				for _, v := range value.NoncurrentVersionTransitions {
					t := make(map[string]interface{})
					if v.NoncurrentDays != nil {
						t["non_current_days"] = int(*v.NoncurrentDays)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["non_current_transition"] = schema.NewSet(transitionHash, transitions)
			}
			// non current expiration
			if value.NoncurrentVersionExpiration != nil {
				e := make(map[string]interface{})
				if value.NoncurrentVersionExpiration.NoncurrentDays != nil {
					e["non_current_days"] = int(*value.NoncurrentVersionExpiration.NoncurrentDays)
				}
				rule["non_current_expiration"] = schema.NewSet(nonCurrentExpirationHash, []interface{}{e})
			}
			// abortIncompleteMultipartUpload
			if value.AbortIncompleteMultipartUpload != nil {
				e := make(map[string]interface{})
				if value.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil {
					e["days_after_initiation"] = int(*value.AbortIncompleteMultipartUpload.DaysAfterInitiation)
				}
				rule["abort_incomplete_multipart_upload"] = schema.NewSet(abortIncompleteMultipartUploadHash, []interface{}{e})
			}

			lifecycleRules = append(lifecycleRules, rule)
		}
	}
	return
}

func (me *CosService) GetDataSourceBucketLifecycle(ctx context.Context, bucket string) (lifecycleRules []map[string]interface{}, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
	}

	ratelimit.Check("GetBucketLifecycleConfiguration")
	response, err := me.client.UseCosClient().GetBucketLifecycleConfiguration(&request)
	if err != nil {
		awsError, ok := err.(awserr.Error)
		if !ok || awsError.Code() != "NoSuchLifecycleConfiguration" {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "get bucket lifecycle", request.String(), err.Error())
			errRet = fmt.Errorf("cos get bucket cors error: %s, bucket: %s", err.Error(), bucket)
			return
		}
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket lifecycle", request.String(), response.String())

	lifecycleRules = make([]map[string]interface{}, 0, len(response.Rules))
	if len(response.Rules) > 0 {
		for _, value := range response.Rules {
			rule := make(map[string]interface{})

			// filter_prefix
			if value.Filter != nil {
				if value.Filter.And != nil && value.Filter.And.Prefix != nil &&
					*value.Filter.And.Prefix != "" {
					rule["filter_prefix"] = *value.Filter.And.Prefix
				} else if value.Filter.Prefix != nil && *value.Filter.Prefix != "" {
					rule["filter_prefix"] = *value.Filter.Prefix
				}
			}
			// transition
			if len(value.Transitions) > 0 {
				transitions := make([]interface{}, 0, len(value.Transitions))
				for _, v := range value.Transitions {
					t := make(map[string]interface{})
					if v.Date != nil {
						t["date"] = (*v.Date).Format("2006-01-02")
					}
					if v.Days != nil {
						t["days"] = int(*v.Days)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["transition"] = transitions
			}
			// expiration
			if value.Expiration != nil {
				e := make(map[string]interface{})
				if value.Expiration.Date != nil {
					e["date"] = (*value.Expiration.Date).Format("2006-01-02")
				}
				if value.Expiration.Days != nil {
					e["days"] = int(*value.Expiration.Days)
				}
				rule["expiration"] = []interface{}{e}
			}
			// non current transition
			if len(value.NoncurrentVersionTransitions) > 0 {
				transitions := make([]interface{}, 0, len(value.NoncurrentVersionTransitions))
				for _, v := range value.NoncurrentVersionTransitions {
					t := make(map[string]interface{})
					if v.NoncurrentDays != nil {
						t["non_current_days"] = int(*v.NoncurrentDays)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["non_current_transition"] = transitions
			}
			// non current expiration
			if value.NoncurrentVersionExpiration != nil {
				e := make(map[string]interface{})
				if value.NoncurrentVersionExpiration.NoncurrentDays != nil {
					e["non_current_days"] = int(*value.NoncurrentVersionExpiration.NoncurrentDays)
				}
				rule["non_current_expiration"] = []interface{}{e}
			}
			// abortIncompleteMultipartUpload
			if value.AbortIncompleteMultipartUpload != nil {
				e := make(map[string]interface{})
				if value.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil {
					e["days_after_initiation"] = int(*value.AbortIncompleteMultipartUpload.DaysAfterInitiation)
				}
				rule["abort_incomplete_multipart_upload"] = []interface{}{e}
			}

			lifecycleRules = append(lifecycleRules, rule)
		}
	}
	return
}

func (me *CosService) GetBucketWebsite(ctx context.Context, bucket string, cdcId string) (websites []map[string]interface{}, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.GetBucketWebsiteInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("GetBucketWebsite")
	response, err := me.client.UseCosClientNew(cdcId).GetBucketWebsite(&request)

	if err != nil {
		awsError, ok := err.(awserr.Error)
		if ok && awsError.Code() == "NoSuchWebsiteConfiguration" {
			return
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get bucket website", request.String(), err.Error())
		errRet = fmt.Errorf("cos get bucket website error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket website", request.String(), response.String())

	websites = make([]map[string]interface{}, 0, 1)
	website := make(map[string]interface{})
	if response.IndexDocument != nil {
		website["index_document"] = *response.IndexDocument.Suffix
	}
	if response.ErrorDocument != nil {
		website["error_document"] = *response.ErrorDocument.Key
	}
	if response.RedirectAllRequestsTo != nil {
		website["redirect_all_requests_to"] = *response.RedirectAllRequestsTo.Protocol
	}
	if len(website) > 0 {
		websites = append(websites, website)
	}

	return
}

func (me *CosService) GetBucketEncryption(ctx context.Context, bucket string, cdcId string) (encryption string, kmsId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("GetBucketEncryption")
	response, err := me.client.UseCosClientNew(cdcId).GetBucketEncryption(&request)

	if err != nil {
		awsError, ok := err.(awserr.Error)
		if ok && awsError.Code() == "NoSuchEncryptionConfiguration" {
			return
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get bucket encryption", request.String(), err.Error())
		errRet = fmt.Errorf("cos get bucket encryption error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket encryption", request.String(), response.String())

	if len(response.ServerSideEncryptionConfiguration.Rules) > 0 {
		encryption = *response.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm
		kMSMasterKeyID := response.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.KMSMasterKeyID
		if kMSMasterKeyID != nil {
			kmsId = *kMSMasterKeyID
		}
	}
	return
}

func (me *CosService) GetBucketVersioning(ctx context.Context, bucket string, cdcId string) (versioningEnable bool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("GetBucketVersioning")
	response, err := me.client.UseCosClientNew(cdcId).GetBucketVersioning(&request)

	if err != nil {
		awsError, ok := err.(awserr.Error)
		if ok && awsError.Code() == "NoSuchVersioningConfiguration" {
			return
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get bucket versioning", request.String(), err.Error())
		errRet = fmt.Errorf("cos get bucket versioning error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket versioning", request.String(), response.String())

	if response.Status == nil || *response.Status == "Suspended" {
		versioningEnable = false
	} else if *response.Status == "Enabled" {
		versioningEnable = true
	}

	return
}

func (me *CosService) GetBucketAccleration(ctx context.Context, bucket string, cdcId string) (accelerationEnable bool, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.GetBucketAccelerateConfigurationInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("GetBucketAccelerateConfiguration")
	response, err := me.client.UseCosClientNew(cdcId).GetBucketAccelerateConfiguration(&request)

	if err != nil {
		awsError, ok := err.(awserr.Error)
		if ok && awsError.Code() == "NoSuchAccelerateConfiguration" {
			return
		}
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get bucket acceleration", request.String(), err.Error())
		errRet = fmt.Errorf("cos get bucket acceleration error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket acceleration", request.String(), response.String())

	if response.Status == nil || *response.Status == "Suspended" {
		accelerationEnable = false
	} else if *response.Status == "Enabled" {
		accelerationEnable = true
	}

	return
}

func (me *CosService) GetBucketLogStatus(ctx context.Context, bucket string, cdcId string) (logEnable bool, logTargetBucket string, logPrefix string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.GetBucketLoggingInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("GetBucketVersioning")
	response, err := me.client.UseCosClientNew(cdcId).GetBucketLogging(&request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get bucket log status", request.String(), err.Error())
		errRet = fmt.Errorf("cos get bucket log status error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket log status", request.String(), response.String())

	if response.LoggingEnabled == nil || response.LoggingEnabled.TargetBucket == nil || *response.LoggingEnabled.TargetBucket == "" || response.LoggingEnabled.TargetPrefix == nil || *response.LoggingEnabled.TargetPrefix == "" {
		logEnable = false
	} else {
		logEnable = true
		logTargetBucket = *response.LoggingEnabled.TargetBucket
		logPrefix = *response.LoggingEnabled.TargetPrefix
	}

	return
}

func (me *CosService) ListBuckets(ctx context.Context) (buckets []*s3.Bucket, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.ListBucketsInput{}
	ratelimit.Check("ListBuckets")
	response, err := me.client.UseCosClient().ListBuckets(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get bucket list", request.String(), err.Error())
		errRet = fmt.Errorf("cos get bucket list error: %s", err.Error())
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket list", request.String(), response.String())

	buckets = response.Buckets
	return
}

func (me *CosService) ListObjects(ctx context.Context, bucket string) (objects []*s3.Object, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("ListObjects")
	response, err := me.client.UseCosClient().ListObjects(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get object list", request.String(), err.Error())
		errRet = fmt.Errorf("cos get object list error: %s", err.Error())
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get object list", request.String(), response.String())

	objects = response.Contents
	return
}

// SetBucketTags if len(tags) == 0, only delete tags
func (me *CosService) SetBucketTags(ctx context.Context, bucket string, tags map[string]string, cdcId string) error {
	logId := tccommon.GetLogId(ctx)

	deleteReq := &s3.DeleteBucketTaggingInput{Bucket: aws.String(bucket)}

	ratelimit.Check("DeleteBucketTagging")
	deleteResp, err := me.client.UseCosClientNew(cdcId).DeleteBucketTagging(deleteReq)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, "delete olg tags", deleteReq.String(), err)
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, "delete olg tags", deleteReq.String(), deleteResp.String())

	if len(tags) == 0 {
		return nil
	}

	putReq := &s3.PutBucketTaggingInput{
		Bucket:  aws.String(bucket),
		Tagging: new(s3.Tagging),
	}

	for k, v := range tags {
		putReq.Tagging.TagSet = append(putReq.Tagging.TagSet, &s3.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	ratelimit.Check("PutBucketTagging")

	resp, err := me.client.UseCosClientNew(cdcId).PutBucketTagging(putReq)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, "put new tags", deleteReq.String(), err)
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, "put new tags", putReq.String(), resp.String())

	return nil
}

func (me *CosService) GetBucketTags(ctx context.Context, bucket string, cdcId string) (map[string]string, error) {
	logId := tccommon.GetLogId(ctx)

	req := &s3.GetBucketTaggingInput{Bucket: aws.String(bucket)}

	ratelimit.Check("GetBucketTagging")
	resp, err := me.client.UseCosClientNew(cdcId).GetBucketTagging(req)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); !ok || awsErr.Code() != "404" {
			return nil, nil
		}

		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, "get tags", req.String(), err)
		return nil, err
	}

	tags := make(map[string]string, len(resp.TagSet))
	for _, t := range resp.TagSet {
		tags[*t.Key] = *t.Value
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, "get tags", req.String(), resp.String())

	return tags, nil
}

func (me *CosService) GetObjectTags(ctx context.Context, bucket string, key string) (map[string]string, error) {
	logId := tccommon.GetLogId(ctx)

	req := &s3.GetObjectTaggingInput{
		Bucket: &bucket,
		Key:    &key,
	}

	ratelimit.Check("GetObjectTagging")
	resp, err := me.client.UseCosClient().GetObjectTagging(req)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); !ok || awsErr.Code() != "404" {
			return nil, nil
		}

		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, "get object tags", req.String(), err)
		return nil, err
	}

	tags := make(map[string]string, len(resp.TagSet))

	for _, tag := range resp.TagSet {
		tags[*tag.Key] = *tag.Value
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, "get tags", req.String(), resp.String())

	return tags, nil
}

// SetObjectTags same as delete Bucket Tags
func (me *CosService) SetObjectTags(ctx context.Context, bucket string, key string, tags map[string]string) error {
	logId := tccommon.GetLogId(ctx)

	deleteReq := &s3.DeleteObjectTaggingInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	ratelimit.Check("DeleteObjectTagging")

	deleteResp, err := me.client.UseCosClient().DeleteObjectTagging(deleteReq)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]",
			logId, "delete olg object tags", deleteReq.String(), err)
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, "delete olg object tags", deleteReq.String(), deleteResp.String())

	if len(tags) == 0 {
		return nil
	}

	putReq := &s3.PutObjectTaggingInput{
		Key:     aws.String(key),
		Bucket:  aws.String(bucket),
		Tagging: new(s3.Tagging),
	}

	for k, v := range tags {
		putReq.Tagging.TagSet = append(putReq.Tagging.TagSet, &s3.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	ratelimit.Check("PutObjectTagging")

	resp, err := me.client.UseCosClient().PutObjectTagging(putReq)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%v]",
			logId, "put new object tags", deleteReq.String(), err)
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, "put new object tags", putReq.String(), resp.String())

	return nil
}

func (me *CosService) PutBucketPolicy(ctx context.Context, bucket, policy string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.PutBucketPolicyInput{
		Bucket: &bucket,
		Policy: &policy,
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket", request.String(), errRet.Error())
		}
	}()
	ratelimit.Check("PutBucketPolicy")
	response, err := me.client.UseCosClient().PutBucketPolicy(&request)
	if err != nil {
		errRet = fmt.Errorf("cos put bucket policy error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket policy", request.String(), response.String())

	return nil
}

func (me *CosService) DescribePolicyByBucket(ctx context.Context, bucket string) (bucketPolicy string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := s3.GetBucketPolicyInput{Bucket: aws.String(bucket)}

	ratelimit.Check("GetBucketPolicy")
	response, err := me.client.UseCosClient().GetBucketPolicy(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "get bucket policy", request.String(), err.Error())
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "get bucket policy", request.String(), response.String())
	bucketPolicy = *response.Policy
	return
}

func (me *CosService) DeleteBucketPolicy(ctx context.Context, bucket string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := s3.DeleteBucketPolicyInput{
		Bucket: aws.String(bucket),
	}
	ratelimit.Check("DeleteBucketPolicy")
	response, err := me.client.UseCosClient().DeleteBucketPolicy(&request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "delete bucket policy", request.String(), err.Error())
		return fmt.Errorf("cos delete bucket policy error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "delete bucket policy", request.String(), response.String())

	return nil
}

func (me *CosService) GetBucketACL(ctx context.Context, bucket string, cdcId string) (result *cos.BucketGetACLResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId, "GetBucketACL", errRet.Error())
		}
	}()

	ratelimit.Check("TencentcloudCosPutBucketACL")
	acl, _, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.GetACL(ctx)

	if err != nil {
		errRet = fmt.Errorf("cos [GetBucketACL] error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	aclXML, err := xml.Marshal(acl)

	if err != nil {
		errRet = fmt.Errorf("cos [GetBucketACL] xml marshal error: %s, bucket: %s", err.Error(), bucket)
		return nil, errRet
	}

	log.Printf("[DEBUG]%s api[%s] success, response body:\n%s\n",
		logId, "GetBucketACL", aclXML)

	result = acl

	return
}

func GetBucketPublicACL(acl *cos.BucketGetACLResult) string {
	var publicRead, publicWrite bool

	for i := range acl.AccessControlList {
		item := acl.AccessControlList[i]

		if item.Grantee.URI == PUBLIC_GRANTEE && item.Permission == "READ" {
			publicRead = true
		}

		if item.Grantee.URI == PUBLIC_GRANTEE && item.Permission == "WRITE" {
			publicWrite = true
		}
	}

	if publicRead && !publicWrite {
		return s3.ObjectCannedACLPublicRead
	}

	if publicRead && publicWrite {
		return s3.ObjectCannedACLPublicReadWrite
	}

	return s3.ObjectCannedACLPrivate
}

func (me *CosService) GetBucketPullOrigin(ctx context.Context, bucket string) (result []map[string]interface{}, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId, "GetBucketPullOrigin", errRet.Error())
		}
	}()

	ratelimit.Check("TencentcloudCosGetBucketPullOrigin")
	originConfig, response, err := me.client.UseTencentCosClient(bucket).Bucket.GetOrigin(ctx)

	if response.StatusCode == 404 {
		return make([]map[string]interface{}, 0), nil
	}

	if err != nil {
		errRet = fmt.Errorf("cos [GetBucketPullOrigin] error: %s, bucket: %s", err.Error(), bucket)
		return nil, errRet
	}

	resp, _ := json.Marshal(originConfig)

	log.Printf("[DEBUG]%s api[%s] success, request body response body [%s]\n",
		logId, "GetBucketPullOrigin", resp)

	rules := make([]map[string]interface{}, 0)

	for _, rule := range originConfig.Rule {
		item := make(map[string]interface{})
		item["priority"] = helper.Int(rule.RulePriority)
		item["host"] = helper.String(rule.OriginInfo.HostInfo.HostName)

		if rule.OriginCondition != nil {
			item["prefix"] = helper.String(rule.OriginCondition.Prefix)
		}

		if rule.OriginType == "Mirror" {
			item["sync_back_to_source"] = helper.Bool(true)
		} else if rule.OriginType == "Proxy" {
			item["sync_back_to_source"] = helper.Bool(false)
		}

		if rule.OriginParameter != nil {
			if rule.OriginParameter.HttpHeader != nil {
				if len(rule.OriginParameter.HttpHeader.NewHttpHeaders) != 0 {
					headers := make(map[string]interface{})
					for _, header := range rule.OriginParameter.HttpHeader.NewHttpHeaders {
						headers[header.Key] = helper.String(header.Value)
					}
					item["custom_http_headers"] = headers
				}

				if len(rule.OriginParameter.HttpHeader.FollowHttpHeaders) != 0 {
					headers := schema.NewSet(func(i interface{}) int {
						return helper.HashString(i.(string))
					}, nil)
					for _, header := range rule.OriginParameter.HttpHeader.FollowHttpHeaders {
						headers.Add(header.Key)
					}
					item["follow_http_headers"] = headers
				}

			}
			item["protocol"] = helper.String(rule.OriginParameter.Protocol)
			item["follow_redirection"] = helper.Bool(rule.OriginParameter.FollowRedirection)
			item["follow_query_string"] = helper.Bool(rule.OriginParameter.FollowQueryString)
		}

		if rule.OriginInfo.FileInfo != nil {
			item["host"] = helper.String(rule.OriginInfo.HostInfo.HostName)
			//item["redirect_prefix"] = helper.String(rule.OriginInfo.FileInfo.Prefix)
			//item["redirect_suffix"] = helper.String(rule.OriginInfo.FileInfo.Suffix)
		}

		rules = append(rules, item)
	}

	return rules, nil
}

func (me *CosService) PutBucketPullOrigin(ctx context.Context, bucket string, rules []cos.BucketOriginRule, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	opt := &cos.BucketPutOriginOptions{
		Rule: rules,
	}
	ratelimit.Check("PutBucketPullOrigin")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.PutOrigin(ctx, opt)

	req, _ := json.Marshal(opt)
	resp, _ := json.Marshal(response.Response.Body)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request:[%s], reason[%s]\n",
				logId, "PutBucketPullOrigin", req, errRet.Error())
		}
	}()

	if err != nil {
		errRet = fmt.Errorf("[PutBucketPullOrigin] error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[PutBucketPullOrigin] success, request body [%s], response body [%s]\n",
		logId, req, resp)

	return nil
}

func (me *CosService) DeleteBucketPullOrigin(ctx context.Context, bucket string, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId, "DeleteBucketPullOrigin", errRet.Error())
		}
	}()

	ratelimit.Check("DeleteBucketPullOrigin")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.DeleteOrigin(ctx)

	resp, _ := json.Marshal(response.Response.Body)

	if err != nil {
		errRet = fmt.Errorf("[DeleteBucketPullOrigin] error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[DeleteBucketPullOrigin] success, response body [%s]\n",
		logId, resp)

	return nil
}

func (me *CosService) GetBucketOriginDomain(ctx context.Context, bucket string) (result []map[string]interface{}, errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId, "GetBucketOriginDomain", errRet.Error())
		}
	}()

	ratelimit.Check("TencentcloudCosGetBucketOriginDomain")
	domain, response, err := me.client.UseTencentCosClient(bucket).Bucket.GetDomain(ctx)

	if response != nil && response.StatusCode == 404 {
		log.Printf("[WARN] [GetBucketOriginDomain] returns %d, %s", 404, err)
		return make([]map[string]interface{}, 0), nil
	}

	if err != nil {
		errRet = fmt.Errorf("cos [GetBucketOriginDomain] error: %s, bucket: %s", err.Error(), bucket)
		return nil, errRet
	}

	rules := make([]map[string]interface{}, 0)

	for _, rule := range domain.Rules {
		item := make(map[string]interface{})
		item["domain"] = helper.String(rule.Name)
		item["status"] = helper.String(rule.Status)
		item["type"] = helper.String(rule.Type)
		rules = append(rules, item)
	}

	resp, _ := json.Marshal(response.Response.Body)

	log.Printf("[DEBUG]%s api[%s] success, request body response body [%s]\n",
		logId, "GetBucketOriginDomain", resp)

	return rules, nil
}

func (me *CosService) PutBucketOriginDomain(ctx context.Context, bucket string, rules []cos.BucketDomainRule, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	opt := &cos.BucketPutDomainOptions{
		Rules: rules,
	}
	ratelimit.Check("PutBucketOriginDomain")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.PutDomain(ctx, opt)

	req, _ := json.Marshal(opt)
	resp, _ := json.Marshal(response.Response.Body)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request:[%s], reason[%s]\n",
				logId, "PutBucketOriginDomain", req, errRet.Error())
		}
	}()

	if err != nil {
		errRet = fmt.Errorf("[PutBucketOriginDomain] error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[PutBucketOriginDomain] success, request body [%s], response body [%s]\n",
		logId, req, resp)

	return nil
}

func (me *CosService) DeleteBucketOriginDomain(ctx context.Context, bucket string, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId, "DeleteBucketOriginDomain", errRet.Error())
		}
	}()

	ratelimit.Check("DeleteBucketOriginDomain")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.DeleteDomain(ctx)

	resp, _ := json.Marshal(response.Response.Body)

	if err != nil {
		errRet = fmt.Errorf("[DeleteBucketOriginDomain] error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	log.Printf("[DEBUG]%s api[DeleteBucketOriginDomain] success, response body [%s]\n",
		logId, resp)

	return nil
}

func (me *CosService) GetBucketReplication(ctx context.Context, bucket string, cdcId string) (result *cos.GetBucketReplicationResult, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId, "GetBucketReplication", errRet.Error())
		}
	}()

	ratelimit.Check("GetBucketReplication")
	result, response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.GetBucketReplication(ctx)

	if response.StatusCode == 404 {
		log.Printf("[WARN]%s, api[%s] returns %d", logId, "GetBucketReplication", response.StatusCode)
		return
	}

	resp, _ := json.Marshal(response.Response.Body)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] response body [%s]\n",
		logId, "GetBucketReplication", resp)

	return
}

func (me *CosService) PutBucketReplication(ctx context.Context, bucket string, role string, rules []cos.BucketReplicationRule, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	option := &cos.PutBucketReplicationOptions{
		Role: role,
		Rule: rules,
	}

	request, _ := xml.Marshal(option)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request: %s reason[%s]\n",
				logId, "PutBucketReplication", request, errRet.Error())
		}
	}()

	ratelimit.Check("PutBucketReplication")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.PutBucketReplication(ctx, option)

	resp, _ := json.Marshal(response.Response.Body)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] response body [%s]\n",
		logId, "PutBucketReplication", resp)

	return
}

func (me *CosService) DeleteBucketReplication(ctx context.Context, bucket string, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]\n",
				logId, "DeleteBucketReplication", errRet.Error())
		}
	}()

	ratelimit.Check("DeleteBucketReplication")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.DeleteBucketReplication(ctx)

	resp, _ := json.Marshal(response.Response.Body)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] response body [%s]\n",
		logId, "DeleteBucketReplication", resp)

	return
}

func (me *CosService) DescribeCosBucketDomainCertificate(ctx context.Context, certId string) (result *cos.BucketGetDomainCertificateResult, bucket string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	ids, err := me.parseCertId(certId)
	if err != nil {
		errRet = err
		return
	}

	bucket = ids.bucket
	domainName := ids.domainName
	option := &cos.BucketGetDomainCertificateOptions{
		DomainName: domainName,
	}
	request, _ := xml.Marshal(option)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request[%s], reason[%s]\n",
				logId, "GetDomainCertificate", request, errRet.Error())
		}
	}()

	result, response, err := me.client.UseTencentCosClient(bucket).Bucket.GetDomainCertificate(ctx, option)
	resp, _ := json.Marshal(response.Response.Body)
	if response.StatusCode == 404 {
		log.Printf("[WARN]%s, api[%s] returns %d", logId, "GetDomainCertificate", response.StatusCode)
		return
	}

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request [%s], response body [%s], result [%s]\n",
		logId, "GetDomainCertificate", request, resp, result)

	return
}

func (me *CosService) DeleteCosBucketDomainCertificate(ctx context.Context, certId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	ids, err := me.parseCertId(certId)
	if err != nil {
		errRet = err
		return
	}

	bucket := ids.bucket
	domainName := ids.domainName
	option := &cos.BucketDeleteDomainCertificateOptions{
		DomainName: domainName,
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, option [%s], reason[%s]\n",
				logId, "DeleteDomainCertificate", option, errRet.Error())
		}
	}()

	ratelimit.Check("DeleteDomainCertificate")
	response, err := me.client.UseTencentCosClient(bucket).Bucket.DeleteDomainCertificate(ctx, option)

	if err != nil {
		errRet = err
		return err
	}

	resp, _ := json.Marshal(response.Response.Body)

	log.Printf("[DEBUG]%s api[%s] success, option [%s], response body [%s]\n",
		logId, "DeleteDomainCertificate", option, resp)

	return
}

func (me *CosService) parseCertId(configId string) (ret *CosBucketDomainCertItem, err error) {
	idSplit := strings.Split(configId, tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return nil, fmt.Errorf("id is broken,%s", configId)
	}

	bucket := idSplit[0]
	domain := idSplit[1]
	if bucket == "" || domain == "" {
		return nil, fmt.Errorf("id is broken,%s", configId)
	}

	ret = &CosBucketDomainCertItem{bucket, domain}
	return
}

func (me *CosService) DescribeCosBucketRefererById(ctx context.Context, bucket string) (*cos.BucketGetRefererResult, error) {
	var errRet error
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetCIGuetzli", bucket, errRet.Error())
		}
	}()

	resRaw, err := tccommon.RetryWithContext(ctx, tccommon.ReadRetryTimeout, func(ctx context.Context) (interface{}, error) {
		res, _, err := me.client.UseTencentCosClient(bucket).Bucket.GetReferer(ctx)
		return res, err
	})

	if err != nil {
		errRet = err
		return nil, errRet
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "GetCIGuetzli", bucket)

	return resRaw.(*cos.BucketGetRefererResult), nil
}

func (me *CosService) DescribeCosBucketVersionById(ctx context.Context, bucket string) (*cos.BucketGetVersionResult, error) {
	var errRet error
	logId := tccommon.GetLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetCIGuetzli", bucket, errRet.Error())
		}
	}()

	resRaw, err := tccommon.RetryWithContext(ctx, tccommon.ReadRetryTimeout, func(ctx context.Context) (interface{}, error) {
		res, _, err := me.client.UseTencentCosClient(bucket).Bucket.GetVersioning(ctx)
		return res, err
	})

	if err != nil {
		errRet = err
		return nil, errRet
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "GetCIGuetzli", bucket)

	return resRaw.(*cos.BucketGetVersionResult), nil
}

func (me *CosService) BucketPutIntelligentTiering(ctx context.Context, bucket string, opt *cos.BucketPutIntelligentTieringOptions, cdcId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	req, _ := json.Marshal(opt)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request [%s], reason[%s]\n",
				logId, "put bucket", req, errRet.Error())
		}
	}()

	ratelimit.Check("BucketPutIntelligentTiering")
	response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.PutIntelligentTiering(ctx, opt)

	if err != nil {
		errRet = fmt.Errorf("cos bucket put intelligent tiering error: %s, bucket: %s", err.Error(), bucket)
		return
	}

	resp, _ := json.Marshal(response.Response.Body)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] response body [%s]\n",
		logId, "BucketPutIntelligentTiering", resp)

	return nil
}

func (me *CosService) BucketGetIntelligentTiering(ctx context.Context, bucket string, cdcId string) (result *cos.BucketGetIntelligentTieringResult, errRet error) {
	logId := tccommon.GetLogId(ctx)

	ratelimit.Check("BucketGetIntelligentTiering")
	intelligentTieringResult, response, err := me.client.UseTencentCosClientNew(bucket, cdcId).Bucket.GetIntelligentTiering(ctx)

	resp, _ := json.Marshal(response.Response.Body)
	if response.StatusCode == 404 {
		log.Printf("[WARN]%s, api[%s] returns %d", logId, "GetDomainCertificate", response.StatusCode)
		return
	}

	if err != nil {
		return
	}
	result = intelligentTieringResult
	log.Printf("[DEBUG]%s api[%s] success, request [%s], response body [%s]\n",
		logId, "GetIntelligentTiering", "", resp)
	return
}

/*
The ideal sequence COS wants.
Priority 1: permission priority: Read first, then handle WRITE, FullControl, WRITE_ACP, last is the READ_ACP
Priority 2: type priority: CanonicalUser first, then Group
*/
func (me *CosService) transACLBodyOrderly(ctx context.Context, rawAclBody string) (orderlyAclBody string, errRet error) {
	// logId := tccommon.GetLogId(ctx)

	rawXmlDoc := etree.NewDocument()
	orderXmlDoc := etree.NewDocument()

	if err := rawXmlDoc.ReadFromString(rawAclBody); err != nil {
		return "", fmt.Errorf("[CRITAL]read raw xml from string error: %v", err)
	}

	rawRoot := rawXmlDoc.SelectElement("AccessControlPolicy")
	orderedRoot := orderXmlDoc.CreateElement("AccessControlPolicy")

	orderedOwner := orderedRoot.CreateElement("Owner")
	for _, ownerChild := range rawRoot.FindElements("//Owner/*") {
		orderedOwner.AddChild(ownerChild)
	}

	orderedACL := orderedRoot.CreateElement("AccessControlList")

	// by combination of permissionSeq and granteeTypeSeq
	for _, perSeq := range COSACLPermissionSeq {
		for _, typeSeq := range COSACLGranteeTypeSeq {
			for _, grantEle := range rawRoot.FindElements(fmt.Sprintf("//Grant[Permission='%s']", perSeq)) {
				granteeEle := grantEle.SelectElement("Grantee")
				if granteeEle != nil {
					if granteeEle.SelectAttrValue("type", "unknown") == typeSeq {
						orderedACL.AddChild(grantEle)
						break
					}
				}
			}
		}
	}

	buf := &bytes.Buffer{}
	orderXmlDoc.Indent(2)
	_, err := orderXmlDoc.WriteTo(buf)
	if err != nil {
		return "", fmt.Errorf("transACLBodyOrderly write xml to buffer failed, error: %v", err)
	}
	orderlyAclBody = buf.String()

	// keep for debug the algo
	// for _, grant := range orderedACL.FindElements("//Grant") {
	// 	grantee := grant.SelectElement("Grantee")
	// 	if grantee != nil {
	// 		// fmt.Printf("===:[%s]====\n", grantee.Tag)
	// 		id := grantee.SelectElement("ID")
	// 		if id != nil {
	// 			fmt.Printf("type:[%s]", id.Text())
	// 		}

	// 		uri := grantee.SelectElement("URI")
	// 		if uri != nil {
	// 			fmt.Printf(" type:[%s]", uri.Text())
	// 		}

	// 		permission := grant.SelectElement("Permission")
	// 		if permission != nil {
	// 			fmt.Printf(" permission:[%s]\n", permission.Text())
	// 		}
	// 	}
	// }
	return
}
