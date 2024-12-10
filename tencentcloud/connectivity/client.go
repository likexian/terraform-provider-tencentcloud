package connectivity

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	intlProfile "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
	mdl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/mdl/v20200326"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	api "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/api/v20201106"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
	apm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apm/v20210622"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	cdwdoris "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228"
	cdwpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	ciam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ciam/v20220331"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	controlcenter "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/controlcenter/v20230110"
	csip "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/csip/v20221121"
	cvmv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	domain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/domain/v20180808"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	postgre "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	region "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	sslCertificate "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	tcaplusdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"
	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	tcss "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcss/v20201101"
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	thpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/thpc/v20230321"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tkev20220501 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20220501"
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss/v20180426"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

//internal version: replace import begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace import end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

const (
	PROVIDER_CVM_REQUEST_TIMEOUT = "TENCENTCLOUD_CVM_REQUEST_TIMEOUT"
	PROVIDER_CBS_REQUEST_TIMEOUT = "TENCENTCLOUD_CBS_REQUEST_TIMEOUT"
)

// TencentCloudClient is client for all TencentCloud service
type TencentCloudClient struct {
	Credential *common.Credential
	Region     string
	Protocol   string
	Domain     string
	CosDomain  string

	cosConn            *s3.S3
	tencentCosConn     *cos.Client
	mysqlConn          *cdb.Client
	redisConn          *redis.Client
	asConn             *as.Client
	vpcConn            *vpc.Client
	cbsConn            *cbs.Client
	cvmv20170312Conn   *cvmv20170312.Client
	clbConn            *clb.Client
	dayuConn           *dayu.Client
	dcConn             *dc.Client
	tagConn            *tag.Client
	mongodbConn        *mongodb.Client
	tkev20180525Conn   *tkev20180525.Client
	tdmqConn           *tdmq.Client
	tcrConn            *tcr.Client
	camConn            *cam.Client
	stsConn            *sts.Client
	gaapConn           *gaap.Client
	sslConn            *ssl.Client
	cfsConn            *cfs.Client
	scfConn            *scf.Client
	tcaplusConn        *tcaplusdb.Client
	cdnConn            *cdn.Client
	monitorConn        *monitor.Client
	esConn             *es.Client
	sqlserverConn      *sqlserver.Client
	postgreConn        *postgre.Client
	ckafkaConn         *ckafka.Client
	auditConn          *audit.Client
	cynosConn          *cynosdb.Client
	vodConn            *vod.Client
	apiGatewayConn     *apigateway.Client
	sslCertificateConn *sslCertificate.Client
	kmsConn            *kms.Client
	ssmConn            *ssm.Client
	apiConn            *api.Client
	emrConn            *emr.Client
	clsConn            *cls.Client
	dnsPodConn         *dnspod.Client
	privateDnsConn     *privatedns.Client
	antiddosConn       *antiddos.Client
	domainConn         *domain.Client
	lighthouseConn     *lighthouse.Client
	temConn            *tem.Client
	teoConn            *teo.Client
	tcmConn            *tcm.Client
	cssConn            *css.Client
	sesConn            *ses.Client
	dcdbConn           *dcdb.Client
	smsConn            *sms.Client
	catConn            *cat.Client
	mariadbConn        *mariadb.Client
	rumConn            *rum.Client
	ptsConn            *pts.Client
	tatConn            *tat.Client
	organizationConn   *organization.Client
	tdcpgConn          *tdcpg.Client
	dbbrainConn        *dbbrain.Client
	dtsConn            *dts.Client
	ciConn             *cos.Client
	cosBatchConn       *cos.Client
	tsfConn            *tsf.Client
	mpsConn            *mps.Client
	cwpConn            *cwp.Client
	chdfsConn          *chdfs.Client
	mdlConn            *mdl.Client
	apmConn            *apm.Client
	ciamConn           *ciam.Client
	tseConn            *tse.Client
	cdwchConn          *cdwch.Client
	ebConn             *eb.Client
	dlcConn            *dlc.Client
	wedataConn         *wedata.Client
	wafConn            *waf.Client
	cfwConn            *cfw.Client
	oceanusConn        *oceanus.Client
	dasbConn           *dasb.Client
	trocketConn        *trocket.Client
	biConn             *bi.Client
	cdwpgConn          *cdwpg.Client
	csipConn           *csip.Client
	regionConn         *region.Client
	//internal version: replace client begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace client end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	tkev20220501Conn  *tkev20220501.Client
	cdcConn           *cdc.Client
	cdwdorisConn      *cdwdoris.Client
	controlcenterConn *controlcenter.Client
	thpcConn          *thpc.Client
	//omit nil client
	omitNilConn             *common.Client
	emrv20190103Conn        *emr.Client
	teov20220901Conn        *teo.Client
	sslv20191205Conn        *sslCertificate.Client
	postgresv20170312Conn   *postgre.Client
	cfwv20190904Conn        *cfw.Client
	ccnv20170312Conn        *vpc.Client
	tcssv20201101Conn       *tcss.Client
	cloudauditv20190319Conn *audit.Client
	privatednsv20201028Conn *privatedns.Client
	wafv20180125Conn        *waf.Client
	camv20190116Conn        *cam.Client
	clsv20201016Conn        *cls.Client
	postgresqlv20170312Conn *postgre.Client
}

// NewClientProfile returns a new ClientProfile
func (me *TencentCloudClient) NewClientProfile(timeout int) *profile.ClientProfile {
	cpf := profile.NewClientProfile()

	// all request use method POST
	cpf.HttpProfile.ReqMethod = "POST"
	// request timeout
	cpf.HttpProfile.ReqTimeout = timeout
	// request protocol
	cpf.HttpProfile.Scheme = me.Protocol
	// request domain
	cpf.HttpProfile.RootDomain = me.Domain
	// default language
	cpf.Language = "en-US"

	return cpf
}

// NewClientIntlProfile returns a new ClientProfile
func (me *TencentCloudClient) NewClientIntlProfile(timeout int) *intlProfile.ClientProfile {
	cpf := intlProfile.NewClientProfile()

	// all request use method POST
	cpf.HttpProfile.ReqMethod = "POST"
	// request timeout
	cpf.HttpProfile.ReqTimeout = timeout
	// request protocol
	cpf.HttpProfile.Scheme = me.Protocol
	// request domain
	cpf.HttpProfile.RootDomain = me.Domain
	// default language
	cpf.Language = "en-US"

	return cpf
}

func (me *TencentCloudClient) UseCosClientNew(cdcId ...string) *s3.S3 {
	if cdcId[0] == "" {
		return me.UseCosClient()
	} else {
		return me.UseCosCdcClient(cdcId[0])
	}
}

// UseCosClient returns cos client for service
func (me *TencentCloudClient) UseCosClient() *s3.S3 {
	if me.cosConn != nil {
		return me.cosConn
	}

	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			cosUrl := fmt.Sprintf("https://cos.%s.myqcloud.com", region)
			if me.CosDomain != "" {
				cosUrl = me.CosDomain
			}
			return endpoints.ResolvedEndpoint{
				URL:           cosUrl,
				SigningRegion: region,
			}, nil
		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	creds := credentials.NewStaticCredentials(me.Credential.SecretId, me.Credential.SecretKey, me.Credential.Token)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      creds,
		Region:           aws.String(me.Region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	}))

	return s3.New(sess)
}

// UseCosClient returns cos client for service with CDC
func (me *TencentCloudClient) UseCosCdcClient(cdcId string) *s3.S3 {
	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			endpointUrl := fmt.Sprintf("https://%s.cos-cdc.%s.myqcloud.com", cdcId, region)
			return endpoints.ResolvedEndpoint{
				URL:           endpointUrl,
				SigningRegion: region,
			}, nil
		}
		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	creds := credentials.NewStaticCredentials(me.Credential.SecretId, me.Credential.SecretKey, me.Credential.Token)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      creds,
		Region:           aws.String(me.Region),
		EndpointResolver: endpoints.ResolverFunc(resolver),
	}))

	return s3.New(sess)
}

func (me *TencentCloudClient) UseTencentCosClientNew(bucket string, cdcId ...string) *cos.Client {
	if cdcId[0] == "" {
		return me.UseTencentCosClient(bucket)
	} else {
		return me.UseTencentCosCdcClient(bucket, cdcId[0])
	}
}

// UseTencentCosClient tencent cloud own client for service instead of aws
func (me *TencentCloudClient) UseTencentCosClient(bucket string) *cos.Client {
	cosUrl := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, me.Region)
	if me.CosDomain != "" {
		parsedURL, _ := url.Parse(me.CosDomain)
		parsedURL.Host = bucket + "." + parsedURL.Host
		cosUrl = parsedURL.String()
	}

	u, _ := url.Parse(cosUrl)

	if me.tencentCosConn != nil && me.tencentCosConn.BaseURL.BucketURL == u {
		return me.tencentCosConn
	}

	baseUrl := &cos.BaseURL{
		BucketURL: u,
	}

	me.tencentCosConn = cos.NewClient(baseUrl, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:     me.Credential.SecretId,
			SecretKey:    me.Credential.SecretKey,
			SessionToken: me.Credential.Token,
		},
	})

	return me.tencentCosConn
}

// UseTencentCosClient tencent cloud own client for service instead of aws with CDC
func (me *TencentCloudClient) UseTencentCosCdcClient(bucket string, cdcId string) *cos.Client {
	var u *url.URL
	u, _ = url.Parse(fmt.Sprintf("https://%s.%s.cos-cdc.%s.myqcloud.com", bucket, cdcId, me.Region))

	if me.tencentCosConn != nil && me.tencentCosConn.BaseURL.BucketURL == u {
		return me.tencentCosConn
	}

	baseUrl := &cos.BaseURL{
		BucketURL: u,
	}

	me.tencentCosConn = cos.NewClient(baseUrl, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:     me.Credential.SecretId,
			SecretKey:    me.Credential.SecretKey,
			SessionToken: me.Credential.Token,
		},
	})

	return me.tencentCosConn
}

// UseMysqlClient returns mysql(cdb) client for service
func (me *TencentCloudClient) UseMysqlClient(iacExtInfo ...IacExtInfo) *cdb.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	// if me.mysqlConn != nil {
	// 	me.mysqlConn.WithHttpTransport(&logRoundTripper)
	// 	return me.mysqlConn
	// }

	cpf := me.NewClientProfile(300)
	me.mysqlConn, _ = cdb.NewClient(me.Credential, me.Region, cpf)
	me.mysqlConn.WithHttpTransport(&logRoundTripper)

	return me.mysqlConn
}

func (me *TencentCloudClient) UseMysqlClientRegion(region string, iacExtInfo ...IacExtInfo) *cdb.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	cpf := me.NewClientProfile(300)
	if region != "" {
		me.mysqlConn, _ = cdb.NewClient(me.Credential, region, cpf)
	} else {
		me.mysqlConn, _ = cdb.NewClient(me.Credential, me.Region, cpf)
	}
	me.mysqlConn.WithHttpTransport(&logRoundTripper)

	return me.mysqlConn
}

// UseRedisClient returns redis client for service
func (me *TencentCloudClient) UseRedisClient() *redis.Client {
	if me.redisConn != nil {
		return me.redisConn
	}

	cpf := me.NewClientProfile(300)
	me.redisConn, _ = redis.NewClient(me.Credential, me.Region, cpf)
	me.redisConn.WithHttpTransport(&LogRoundTripper{})

	return me.redisConn
}

// UseAsClient returns as client for service
func (me *TencentCloudClient) UseAsClient() *as.Client {
	if me.asConn != nil {
		return me.asConn
	}

	cpf := me.NewClientProfile(300)
	me.asConn, _ = as.NewClient(me.Credential, me.Region, cpf)
	me.asConn.WithHttpTransport(&LogRoundTripper{})

	return me.asConn
}

// UseVpcClient returns vpc client for service
func (me *TencentCloudClient) UseVpcClient(iacExtInfo ...IacExtInfo) *vpc.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.vpcConn != nil {
		me.vpcConn.WithHttpTransport(&logRoundTripper)
		return me.vpcConn
	}

	cpf := me.NewClientProfile(300)
	me.vpcConn, _ = vpc.NewClient(me.Credential, me.Region, cpf)
	me.vpcConn.WithHttpTransport(&logRoundTripper)

	return me.vpcConn
}

func (me *TencentCloudClient) UseOmitNilClient(module string) *common.Client {
	secretId := me.Credential.SecretId
	secretKey := me.Credential.SecretKey
	token := me.Credential.Token
	region := me.Region
	var credential common.CredentialIface
	if token != "" {
		credential = common.NewTokenCredential(secretId, secretKey, token)
	} else {
		credential = common.NewCredential(secretId, secretKey)
	}

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = fmt.Sprintf("%s.tencentcloudapi.com", module)
	cpf.HttpProfile.ReqMethod = "POST"
	me.omitNilConn = common.NewCommonClient(credential, region, cpf).WithLogger(log.Default())

	return me.omitNilConn
}

// UseCbsClient returns cbs client for service
func (me *TencentCloudClient) UseCbsClient(iacExtInfo ...IacExtInfo) *cbs.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.cbsConn != nil {
		me.cbsConn.WithHttpTransport(&logRoundTripper)
		return me.cbsConn
	}

	var reqTimeout = getEnvDefault(PROVIDER_CBS_REQUEST_TIMEOUT, 300)
	cpf := me.NewClientProfile(reqTimeout)
	me.cbsConn, _ = cbs.NewClient(me.Credential, me.Region, cpf)
	me.cbsConn.WithHttpTransport(&logRoundTripper)

	return me.cbsConn
}

// UseDcClient returns dc client for service
func (me *TencentCloudClient) UseDcClient() *dc.Client {
	if me.dcConn != nil {
		return me.dcConn
	}

	cpf := me.NewClientProfile(300)
	me.dcConn, _ = dc.NewClient(me.Credential, me.Region, cpf)
	me.dcConn.WithHttpTransport(&LogRoundTripper{})

	return me.dcConn
}

// UseMongodbClient returns mongodb client for service
func (me *TencentCloudClient) UseMongodbClient(iacExtInfo ...IacExtInfo) *mongodb.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.mongodbConn != nil {
		me.mongodbConn.WithHttpTransport(&logRoundTripper)
		return me.mongodbConn
	}

	cpf := me.NewClientProfile(300)
	me.mongodbConn, _ = mongodb.NewClient(me.Credential, me.Region, cpf)
	me.mongodbConn.WithHttpTransport(&logRoundTripper)

	return me.mongodbConn
}

// UseClbClient returns clb client for service
func (me *TencentCloudClient) UseClbClient(iacExtInfo ...IacExtInfo) *clb.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.clbConn != nil {
		me.clbConn.WithHttpTransport(&logRoundTripper)
		return me.clbConn
	}

	cpf := me.NewClientProfile(300)
	me.clbConn, _ = clb.NewClient(me.Credential, me.Region, cpf)
	me.clbConn.WithHttpTransport(&logRoundTripper)

	return me.clbConn
}

// UseCvmClient returns cvm client for service
func (me *TencentCloudClient) UseCvmClient(iacExtInfo ...IacExtInfo) *cvmv20170312.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.cvmv20170312Conn != nil {
		me.cvmv20170312Conn.WithHttpTransport(&logRoundTripper)
		return me.cvmv20170312Conn
	}

	var reqTimeout = getEnvDefault(PROVIDER_CVM_REQUEST_TIMEOUT, 300)
	cpf := me.NewClientProfile(reqTimeout)
	me.cvmv20170312Conn, _ = cvmv20170312.NewClient(me.Credential, me.Region, cpf)
	me.cvmv20170312Conn.WithHttpTransport(&logRoundTripper)

	return me.cvmv20170312Conn
}

// UseCvmV20170312Client returns cvm client for service
func (me *TencentCloudClient) UseCvmV20170312Client(iacExtInfo ...IacExtInfo) *cvmv20170312.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.cvmv20170312Conn != nil {
		me.cvmv20170312Conn.WithHttpTransport(&logRoundTripper)
		return me.cvmv20170312Conn
	}

	var reqTimeout = getEnvDefault(PROVIDER_CVM_REQUEST_TIMEOUT, 300)
	cpf := me.NewClientProfile(reqTimeout)
	me.cvmv20170312Conn, _ = cvmv20170312.NewClient(me.Credential, me.Region, cpf)
	me.cvmv20170312Conn.WithHttpTransport(&logRoundTripper)

	return me.cvmv20170312Conn
}

// UseTagClient returns tag client for service
func (me *TencentCloudClient) UseTagClient() *tag.Client {
	if me.tagConn != nil {
		return me.tagConn
	}

	cpf := me.NewClientProfile(300)
	me.tagConn, _ = tag.NewClient(me.Credential, me.Region, cpf)
	me.tagConn.WithHttpTransport(&LogRoundTripper{})

	return me.tagConn
}

// UseTkeClient returns tke client for service
func (me *TencentCloudClient) UseTkeClient(iacExtInfo ...IacExtInfo) *tkev20180525.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.tkev20180525Conn != nil {
		me.tkev20180525Conn.WithHttpTransport(&logRoundTripper)
		return me.tkev20180525Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.tkev20180525Conn, _ = tkev20180525.NewClient(me.Credential, me.Region, cpf)
	me.tkev20180525Conn.WithHttpTransport(&logRoundTripper)

	return me.tkev20180525Conn
}

// UseTkeV20180525Client returns tke client for service
func (me *TencentCloudClient) UseTkeV20180525Client(iacExtInfo ...IacExtInfo) *tkev20180525.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.tkev20180525Conn != nil {
		me.tkev20180525Conn.WithHttpTransport(&logRoundTripper)
		return me.tkev20180525Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.tkev20180525Conn, _ = tkev20180525.NewClient(me.Credential, me.Region, cpf)
	me.tkev20180525Conn.WithHttpTransport(&logRoundTripper)

	return me.tkev20180525Conn
}

// UseTdmqClient returns Tdmq client for service
func (me *TencentCloudClient) UseTdmqClient(iacExtInfo ...IacExtInfo) *tdmq.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.tdmqConn != nil {
		me.tdmqConn.WithHttpTransport(&logRoundTripper)
		return me.tdmqConn
	}

	cpf := me.NewClientProfile(300)
	me.tdmqConn, _ = tdmq.NewClient(me.Credential, me.Region, cpf)
	me.tdmqConn.WithHttpTransport(&logRoundTripper)

	return me.tdmqConn
}

// UseGaapClient returns gaap client for service
func (me *TencentCloudClient) UseGaapClient(iacExtInfo ...IacExtInfo) *gaap.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.gaapConn != nil {
		me.gaapConn.WithHttpTransport(&logRoundTripper)
		return me.gaapConn
	}

	cpf := me.NewClientProfile(300)
	me.gaapConn, _ = gaap.NewClient(me.Credential, me.Region, cpf)
	me.gaapConn.WithHttpTransport(&logRoundTripper)

	return me.gaapConn
}

// UseSslClient returns ssl client for service
func (me *TencentCloudClient) UseSslClient() *ssl.Client {
	if me.sslConn != nil {
		return me.sslConn
	}

	cpf := me.NewClientProfile(300)
	me.sslConn, _ = ssl.NewClient(me.Credential, me.Region, cpf)
	me.sslConn.WithHttpTransport(&LogRoundTripper{})

	return me.sslConn
}

// UseCamClient returns cam client for service
func (me *TencentCloudClient) UseCamClient() *cam.Client {
	if me.camConn != nil {
		return me.camConn
	}

	cpf := me.NewClientProfile(300)
	me.camConn, _ = cam.NewClient(me.Credential, me.Region, cpf)
	me.camConn.WithHttpTransport(&LogRoundTripper{})

	return me.camConn
}

// UseStsClient returns sts client for service
func (me *TencentCloudClient) UseStsClient(stsExtInfo ...StsExtInfo) *sts.Client {
	/*
		me.Credential will changed, don't cache it
		if me.stsConn != nil {
			return me.stsConn
		}
	*/

	var logRoundTripper LogRoundTripper
	if len(stsExtInfo) != 0 {
		logRoundTripper.Authorization = stsExtInfo[0].Authorization
	}

	cpf := me.NewClientProfile(300)
	me.stsConn, _ = sts.NewClient(me.Credential, me.Region, cpf)
	me.stsConn.WithHttpTransport(&logRoundTripper)

	return me.stsConn
}

// UseCfsClient returns cfs client for service
func (me *TencentCloudClient) UseCfsClient() *cfs.Client {
	if me.cfsConn != nil {
		return me.cfsConn
	}

	cpf := me.NewClientProfile(300)
	me.cfsConn, _ = cfs.NewClient(me.Credential, me.Region, cpf)
	me.cfsConn.WithHttpTransport(&LogRoundTripper{})

	return me.cfsConn
}

// UseScfClient returns scf client for service
func (me *TencentCloudClient) UseScfClient(iacExtInfo ...IacExtInfo) *scf.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.scfConn != nil {
		me.scfConn.WithHttpTransport(&logRoundTripper)
		return me.scfConn
	}

	cpf := me.NewClientProfile(300)
	me.scfConn, _ = scf.NewClient(me.Credential, me.Region, cpf)
	me.scfConn.WithHttpTransport(&logRoundTripper)

	return me.scfConn
}

// UseTcaplusClient returns tcaplush client for service
func (me *TencentCloudClient) UseTcaplusClient() *tcaplusdb.Client {
	if me.tcaplusConn != nil {
		return me.tcaplusConn
	}

	cpf := me.NewClientProfile(300)
	me.tcaplusConn, _ = tcaplusdb.NewClient(me.Credential, me.Region, cpf)
	me.tcaplusConn.WithHttpTransport(&LogRoundTripper{})

	return me.tcaplusConn
}

// UseDayuClient returns dayu client for service
func (me *TencentCloudClient) UseDayuClient() *dayu.Client {
	if me.dayuConn != nil {
		return me.dayuConn
	}

	cpf := me.NewClientProfile(300)
	me.dayuConn, _ = dayu.NewClient(me.Credential, me.Region, cpf)
	me.dayuConn.WithHttpTransport(&LogRoundTripper{})

	return me.dayuConn
}

// UseCdnClient returns cdn client for service
func (me *TencentCloudClient) UseCdnClient(iacExtInfo ...IacExtInfo) *cdn.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.cdnConn != nil {
		me.cdnConn.WithHttpTransport(&logRoundTripper)
		return me.cdnConn
	}

	cpf := me.NewClientProfile(300)
	me.cdnConn, _ = cdn.NewClient(me.Credential, me.Region, cpf)
	me.cdnConn.WithHttpTransport(&logRoundTripper)

	return me.cdnConn
}

// UseMonitorClient returns monitor client for service
func (me *TencentCloudClient) UseMonitorClient() *monitor.Client {
	if me.monitorConn != nil {
		return me.monitorConn
	}

	cpf := me.NewClientProfile(300)
	me.monitorConn, _ = monitor.NewClient(me.Credential, me.Region, cpf)
	me.monitorConn.WithHttpTransport(&LogRoundTripper{})

	return me.monitorConn
}

// UseEsClient returns es client for service
func (me *TencentCloudClient) UseEsClient(iacExtInfo ...IacExtInfo) *es.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.esConn != nil {
		me.esConn.WithHttpTransport(&logRoundTripper)
		return me.esConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.esConn, _ = es.NewClient(me.Credential, me.Region, cpf)
	me.esConn.WithHttpTransport(&logRoundTripper)

	return me.esConn
}

// UsePostgresqlClient returns postgresql client for service
func (me *TencentCloudClient) UsePostgresqlClient(iacExtInfo ...IacExtInfo) *postgre.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.postgreConn != nil {
		me.postgreConn.WithHttpTransport(&logRoundTripper)
		return me.postgreConn
	}

	cpf := me.NewClientProfile(300)
	me.postgreConn, _ = postgre.NewClient(me.Credential, me.Region, cpf)
	me.postgreConn.WithHttpTransport(&logRoundTripper)

	return me.postgreConn
}

// UseSqlserverClient returns sqlserver client for service
func (me *TencentCloudClient) UseSqlserverClient(iacExtInfo ...IacExtInfo) *sqlserver.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.sqlserverConn != nil {
		me.sqlserverConn.WithHttpTransport(&logRoundTripper)
		return me.sqlserverConn
	}

	cpf := me.NewClientProfile(300)
	me.sqlserverConn, _ = sqlserver.NewClient(me.Credential, me.Region, cpf)
	me.sqlserverConn.WithHttpTransport(&logRoundTripper)

	return me.sqlserverConn
}

// UseCkafkaClient returns ckafka client for service
func (me *TencentCloudClient) UseCkafkaClient(iacExtInfo ...IacExtInfo) *ckafka.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.ckafkaConn != nil {
		me.ckafkaConn.WithHttpTransport(&logRoundTripper)
		return me.ckafkaConn
	}

	cpf := me.NewClientProfile(300)
	me.ckafkaConn, _ = ckafka.NewClient(me.Credential, me.Region, cpf)
	me.ckafkaConn.WithHttpTransport(&logRoundTripper)

	return me.ckafkaConn
}

// UseAuditClient returns audit client for service
func (me *TencentCloudClient) UseAuditClient() *audit.Client {
	if me.auditConn != nil {
		return me.auditConn
	}

	cpf := me.NewClientProfile(300)
	me.auditConn, _ = audit.NewClient(me.Credential, me.Region, cpf)
	me.auditConn.WithHttpTransport(&LogRoundTripper{})

	return me.auditConn
}

// UseCynosdbClient returns cynosdb client for service
func (me *TencentCloudClient) UseCynosdbClient() *cynosdb.Client {
	if me.cynosConn != nil {
		return me.cynosConn
	}

	cpf := me.NewClientProfile(300)
	me.cynosConn, _ = cynosdb.NewClient(me.Credential, me.Region, cpf)
	me.cynosConn.WithHttpTransport(&LogRoundTripper{})

	return me.cynosConn
}

// UseVodClient returns vod client for service
func (me *TencentCloudClient) UseVodClient() *vod.Client {
	if me.vodConn != nil {
		return me.vodConn
	}

	cpf := me.NewClientProfile(300)
	me.vodConn, _ = vod.NewClient(me.Credential, me.Region, cpf)
	me.vodConn.WithHttpTransport(&LogRoundTripper{})

	return me.vodConn
}

// UseAPIGatewayClient returns apigateway client for service
func (me *TencentCloudClient) UseAPIGatewayClient() *apigateway.Client {
	if me.apiGatewayConn != nil {
		return me.apiGatewayConn
	}

	cpf := me.NewClientProfile(300)
	me.apiGatewayConn, _ = apigateway.NewClient(me.Credential, me.Region, cpf)
	me.apiGatewayConn.WithHttpTransport(&LogRoundTripper{})

	return me.apiGatewayConn
}

// UseTCRClient returns apigateway client for service
func (me *TencentCloudClient) UseTCRClient(iacExtInfo ...IacExtInfo) *tcr.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.tcrConn != nil {
		me.tcrConn.WithHttpTransport(&logRoundTripper)
		return me.tcrConn
	}

	cpf := me.NewClientProfile(300)
	me.tcrConn, _ = tcr.NewClient(me.Credential, me.Region, cpf)
	me.tcrConn.WithHttpTransport(&logRoundTripper)

	return me.tcrConn
}

// UseSSLCertificateClient returns SSL Certificate client for service
func (me *TencentCloudClient) UseSSLCertificateClient() *sslCertificate.Client {
	if me.sslCertificateConn != nil {
		return me.sslCertificateConn
	}

	cpf := me.NewClientProfile(300)
	me.sslCertificateConn, _ = sslCertificate.NewClient(me.Credential, me.Region, cpf)
	me.sslCertificateConn.WithHttpTransport(&LogRoundTripper{})

	return me.sslCertificateConn
}

// UseKmsClient returns KMS client for service
func (me *TencentCloudClient) UseKmsClient() *kms.Client {
	if me.kmsConn != nil {
		return me.kmsConn
	}

	cpf := me.NewClientProfile(300)
	me.kmsConn, _ = kms.NewClient(me.Credential, me.Region, cpf)
	me.kmsConn.WithHttpTransport(&LogRoundTripper{})

	return me.kmsConn
}

// UseSsmClient returns SSM client for service
func (me *TencentCloudClient) UseSsmClient() *ssm.Client {
	if me.ssmConn != nil {
		return me.ssmConn
	}

	cpf := me.NewClientProfile(300)
	me.ssmConn, _ = ssm.NewClient(me.Credential, me.Region, cpf)
	me.ssmConn.WithHttpTransport(&LogRoundTripper{})

	return me.ssmConn
}

// UseApiClient return API client for service
func (me *TencentCloudClient) UseApiClient() *api.Client {
	if me.apiConn != nil {
		return me.apiConn
	}
	cpf := me.NewClientProfile(300)
	me.apiConn, _ = api.NewClient(me.Credential, me.Region, cpf)
	me.apiConn.WithHttpTransport(&LogRoundTripper{})

	return me.apiConn
}

// UseEmrClient return EMR client for service
func (me *TencentCloudClient) UseEmrClient() *emr.Client {
	if me.emrConn != nil {
		return me.emrConn
	}
	cpf := me.NewClientProfile(300)
	me.emrConn, _ = emr.NewClient(me.Credential, me.Region, cpf)
	me.emrConn.WithHttpTransport(&LogRoundTripper{})

	return me.emrConn
}

// UseClsClient return CLS client for service
func (me *TencentCloudClient) UseClsClient(iacExtInfo ...IacExtInfo) *cls.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.clsConn != nil {
		me.clsConn.WithHttpTransport(&logRoundTripper)
		return me.clsConn
	}

	cpf := me.NewClientProfile(300)
	me.clsConn, _ = cls.NewClient(me.Credential, me.Region, cpf)
	me.clsConn.WithHttpTransport(&logRoundTripper)

	return me.clsConn
}

// UseLighthouseClient return Lighthouse client for service
func (me *TencentCloudClient) UseLighthouseClient(iacExtInfo ...IacExtInfo) *lighthouse.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.lighthouseConn != nil {
		me.lighthouseConn.WithHttpTransport(&logRoundTripper)
		return me.lighthouseConn
	}

	cpf := me.NewClientProfile(300)
	me.lighthouseConn, _ = lighthouse.NewClient(me.Credential, me.Region, cpf)
	me.lighthouseConn.WithHttpTransport(&logRoundTripper)

	return me.lighthouseConn
}

// UseDnsPodClient return DnsPod client for service
func (me *TencentCloudClient) UseDnsPodClient() *dnspod.Client {
	if me.dnsPodConn != nil {
		return me.dnsPodConn
	}
	cpf := me.NewClientProfile(300)
	me.dnsPodConn, _ = dnspod.NewClient(me.Credential, me.Region, cpf)
	me.dnsPodConn.WithHttpTransport(&LogRoundTripper{})

	return me.dnsPodConn
}

// UsePrivateDnsClient return PrivateDns client for service
func (me *TencentCloudClient) UsePrivateDnsClient(iacExtInfo ...IacExtInfo) *privatedns.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.privateDnsConn != nil {
		me.privateDnsConn.WithHttpTransport(&logRoundTripper)
		return me.privateDnsConn
	}

	cpf := me.NewClientProfile(300)
	me.privateDnsConn, _ = privatedns.NewClient(me.Credential, me.Region, cpf)
	me.privateDnsConn.WithHttpTransport(&logRoundTripper)

	return me.privateDnsConn
}

// UseDomainClient return Domain client for service
func (me *TencentCloudClient) UseDomainClient() *domain.Client {
	if me.domainConn != nil {
		return me.domainConn
	}
	cpf := me.NewClientProfile(300)
	me.domainConn, _ = domain.NewClient(me.Credential, me.Region, cpf)
	me.domainConn.WithHttpTransport(&LogRoundTripper{})

	return me.domainConn
}

// UseAntiddosClient returns antiddos client for service
func (me *TencentCloudClient) UseAntiddosClient() *antiddos.Client {
	if me.antiddosConn != nil {
		return me.antiddosConn
	}

	cpf := me.NewClientProfile(300)
	me.antiddosConn, _ = antiddos.NewClient(me.Credential, me.Region, cpf)
	me.antiddosConn.WithHttpTransport(&LogRoundTripper{})

	return me.antiddosConn
}

// UseTemClient returns tem client for service
func (me *TencentCloudClient) UseTemClient() *tem.Client {
	if me.temConn != nil {
		return me.temConn
	}

	cpf := me.NewClientProfile(300)
	me.temConn, _ = tem.NewClient(me.Credential, me.Region, cpf)
	me.temConn.WithHttpTransport(&LogRoundTripper{})

	return me.temConn
}

// UseTeoClient returns teo client for service
func (me *TencentCloudClient) UseTeoClient(iacExtInfo ...IacExtInfo) *teo.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.teoConn != nil {
		me.teoConn.WithHttpTransport(&logRoundTripper)
		return me.teoConn
	}

	cpf := me.NewClientProfile(300)
	me.teoConn, _ = teo.NewClient(me.Credential, me.Region, cpf)
	me.teoConn.WithHttpTransport(&logRoundTripper)

	return me.teoConn
}

// UseTcmClient returns Tcm client for service
func (me *TencentCloudClient) UseTcmClient() *tcm.Client {
	if me.tcmConn != nil {
		return me.tcmConn
	}

	cpf := me.NewClientProfile(300)
	me.tcmConn, _ = tcm.NewClient(me.Credential, me.Region, cpf)
	me.tcmConn.WithHttpTransport(&LogRoundTripper{})

	return me.tcmConn
}

// UseCssClient returns css client for service
func (me *TencentCloudClient) UseCssClient() *css.Client {
	if me.cssConn != nil {
		return me.cssConn
	}

	cpf := me.NewClientProfile(300)
	me.cssConn, _ = css.NewClient(me.Credential, me.Region, cpf)
	me.cssConn.WithHttpTransport(&LogRoundTripper{})

	return me.cssConn
}

// UseSesClient returns Ses client for service
func (me *TencentCloudClient) UseSesClient() *ses.Client {
	if me.sesConn != nil {
		return me.sesConn
	}

	cpf := me.NewClientProfile(300)
	me.sesConn, _ = ses.NewClient(me.Credential, me.Region, cpf)
	me.sesConn.WithHttpTransport(&LogRoundTripper{})

	return me.sesConn
}

// UseDcdbClient returns dcdb client for service
func (me *TencentCloudClient) UseDcdbClient() *dcdb.Client {
	if me.dcdbConn != nil {
		return me.dcdbConn
	}

	cpf := me.NewClientProfile(300)
	me.dcdbConn, _ = dcdb.NewClient(me.Credential, me.Region, cpf)
	me.dcdbConn.WithHttpTransport(&LogRoundTripper{})

	return me.dcdbConn
}

// UseSmsClient returns Sms client for service
func (me *TencentCloudClient) UseSmsClient() *sms.Client {
	if me.smsConn != nil {
		return me.smsConn
	}

	cpf := me.NewClientProfile(300)
	me.smsConn, _ = sms.NewClient(me.Credential, me.Region, cpf)
	me.smsConn.WithHttpTransport(&LogRoundTripper{})

	return me.smsConn
}

// UseCatClient returns Cat client for service
func (me *TencentCloudClient) UseCatClient() *cat.Client {
	if me.catConn != nil {
		return me.catConn
	}

	cpf := me.NewClientProfile(300)
	me.catConn, _ = cat.NewClient(me.Credential, me.Region, cpf)
	me.catConn.WithHttpTransport(&LogRoundTripper{})

	return me.catConn
}

// UseMariadbClient returns mariadb client for service
func (me *TencentCloudClient) UseMariadbClient(iacExtInfo ...IacExtInfo) *mariadb.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.mariadbConn != nil {
		me.mariadbConn.WithHttpTransport(&logRoundTripper)
		return me.mariadbConn
	}

	cpf := me.NewClientProfile(300)
	me.mariadbConn, _ = mariadb.NewClient(me.Credential, me.Region, cpf)
	me.mariadbConn.WithHttpTransport(&logRoundTripper)

	return me.mariadbConn
}

// UsePtsClient returns pts client for service
func (me *TencentCloudClient) UsePtsClient() *pts.Client {
	if me.ptsConn != nil {
		return me.ptsConn
	}

	cpf := me.NewClientProfile(300)
	me.ptsConn, _ = pts.NewClient(me.Credential, me.Region, cpf)
	me.ptsConn.WithHttpTransport(&LogRoundTripper{})

	return me.ptsConn
}

// UseTatClient returns tat client for service
func (me *TencentCloudClient) UseTatClient() *tat.Client {
	if me.tatConn != nil {
		return me.tatConn
	}

	cpf := me.NewClientProfile(300)
	me.tatConn, _ = tat.NewClient(me.Credential, me.Region, cpf)
	me.tatConn.WithHttpTransport(&LogRoundTripper{})

	return me.tatConn
}

// UseOrganizationClient returns organization client for service
func (me *TencentCloudClient) UseOrganizationClient() *organization.Client {
	if me.organizationConn != nil {
		return me.organizationConn
	}

	cpf := me.NewClientProfile(300)
	me.organizationConn, _ = organization.NewClient(me.Credential, me.Region, cpf)
	me.organizationConn.WithHttpTransport(&LogRoundTripper{})

	return me.organizationConn
}

// UseTdcpgClient returns tdcpg client for service
func (me *TencentCloudClient) UseTdcpgClient(iacExtInfo ...IacExtInfo) *tdcpg.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.tdcpgConn != nil {
		me.tdcpgConn.WithHttpTransport(&logRoundTripper)
		return me.tdcpgConn
	}

	cpf := me.NewClientProfile(300)
	me.tdcpgConn, _ = tdcpg.NewClient(me.Credential, me.Region, cpf)
	me.tdcpgConn.WithHttpTransport(&logRoundTripper)

	return me.tdcpgConn
}

// UseDbbrainClient returns dbbrain client for service
func (me *TencentCloudClient) UseDbbrainClient() *dbbrain.Client {
	if me.dbbrainConn != nil {
		return me.dbbrainConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.dbbrainConn, _ = dbbrain.NewClient(me.Credential, me.Region, cpf)
	me.dbbrainConn.WithHttpTransport(&LogRoundTripper{})

	return me.dbbrainConn
}

// UseRumClient returns rum client for service
func (me *TencentCloudClient) UseRumClient() *rum.Client {
	if me.rumConn != nil {
		return me.rumConn
	}

	cpf := me.NewClientProfile(300)
	me.rumConn, _ = rum.NewClient(me.Credential, me.Region, cpf)
	me.rumConn.WithHttpTransport(&LogRoundTripper{})

	return me.rumConn
}

// UseDtsClient returns dts client for service
func (me *TencentCloudClient) UseDtsClient() *dts.Client {
	if me.dtsConn != nil {
		return me.dtsConn
	}

	cpf := me.NewClientProfile(300)
	me.dtsConn, _ = dts.NewClient(me.Credential, me.Region, cpf)
	me.dtsConn.WithHttpTransport(&LogRoundTripper{})

	return me.dtsConn
}

// UseCosBatchClient returns ci client for service
func (me *TencentCloudClient) UseCosBatchClient(uin string) *cos.Client {
	cosUrl := fmt.Sprintf("https://%s.cos-control.%s.myqcloud.com", uin, me.Region)
	if me.CosDomain != "" {
		cosUrl = me.CosDomain
	}

	u, _ := url.Parse(cosUrl)

	if me.cosBatchConn != nil && me.cosBatchConn.BaseURL.BatchURL == u {
		return me.cosBatchConn
	}

	baseUrl := &cos.BaseURL{
		BatchURL: u,
	}

	me.cosBatchConn = cos.NewClient(baseUrl, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:     me.Credential.SecretId,
			SecretKey:    me.Credential.SecretKey,
			SessionToken: me.Credential.Token,
		},
	})

	return me.cosBatchConn
}

// UseCiClient returns ci client for service
func (me *TencentCloudClient) UseCiClient(bucket string) *cos.Client {
	u, _ := url.Parse(fmt.Sprintf("https://%s.ci.%s.myqcloud.com", bucket, me.Region))

	if me.ciConn != nil && me.ciConn.BaseURL.CIURL == u {
		return me.ciConn
	}

	baseUrl := &cos.BaseURL{
		CIURL: u,
	}

	me.ciConn = cos.NewClient(baseUrl, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:     me.Credential.SecretId,
			SecretKey:    me.Credential.SecretKey,
			SessionToken: me.Credential.Token,
		},
	})

	return me.ciConn
}

// UsePicClient returns pic client for service
func (me *TencentCloudClient) UsePicClient(bucket string) *cos.Client {
	u, _ := url.Parse(fmt.Sprintf("https://%s.pic.%s.myqcloud.com", bucket, me.Region))

	if me.ciConn != nil && me.ciConn.BaseURL.CIURL == u {
		return me.ciConn
	}

	baseUrl := &cos.BaseURL{
		CIURL: u,
	}

	me.ciConn = cos.NewClient(baseUrl, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:     me.Credential.SecretId,
			SecretKey:    me.Credential.SecretKey,
			SessionToken: me.Credential.Token,
		},
	})

	return me.ciConn
}

// UseTsfClient returns tsf client for service
func (me *TencentCloudClient) UseTsfClient() *tsf.Client {
	if me.tsfConn != nil {
		return me.tsfConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.tsfConn, _ = tsf.NewClient(me.Credential, me.Region, cpf)
	me.tsfConn.WithHttpTransport(&LogRoundTripper{})

	return me.tsfConn
}

// UseMpsClient returns mps client for service
func (me *TencentCloudClient) UseMpsClient() *mps.Client {
	if me.mpsConn != nil {
		return me.mpsConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.mpsConn, _ = mps.NewClient(me.Credential, me.Region, cpf)
	me.mpsConn.WithHttpTransport(&LogRoundTripper{})

	return me.mpsConn
}

// UseCwpClient returns tke client for service
func (me *TencentCloudClient) UseCwpClient() *cwp.Client {
	if me.cwpConn != nil {
		return me.cwpConn
	}

	cpf := me.NewClientProfile(300)
	me.cwpConn, _ = cwp.NewClient(me.Credential, me.Region, cpf)
	me.cwpConn.WithHttpTransport(&LogRoundTripper{})

	return me.cwpConn
}

// UseChdfsClient returns chdfs client for service
func (me *TencentCloudClient) UseChdfsClient() *chdfs.Client {
	if me.chdfsConn != nil {
		return me.chdfsConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.chdfsConn, _ = chdfs.NewClient(me.Credential, me.Region, cpf)
	me.chdfsConn.WithHttpTransport(&LogRoundTripper{})

	return me.chdfsConn
}

// UseMdlClient returns mdl client for service
func (me *TencentCloudClient) UseMdlClient() *mdl.Client {
	if me.mdlConn != nil {
		return me.mdlConn
	}

	cpf := me.NewClientIntlProfile(300)
	cpf.Language = "zh-CN"
	me.mdlConn, _ = mdl.NewClient(me.Credential, me.Region, cpf)
	me.mdlConn.WithHttpTransport(&LogRoundTripper{})

	return me.mdlConn
}

// UseApmClient returns apm client for service
func (me *TencentCloudClient) UseApmClient() *apm.Client {
	if me.apmConn != nil {
		return me.apmConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.apmConn, _ = apm.NewClient(me.Credential, me.Region, cpf)
	me.apmConn.WithHttpTransport(&LogRoundTripper{})

	return me.apmConn
}

// UseCiamClient returns ciam client for service
func (me *TencentCloudClient) UseCiamClient() *ciam.Client {
	if me.ciamConn != nil {
		return me.ciamConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.ciamConn, _ = ciam.NewClient(me.Credential, me.Region, cpf)
	me.ciamConn.WithHttpTransport(&LogRoundTripper{})

	return me.ciamConn
}

// UseTseClient returns tse client for service
func (me *TencentCloudClient) UseTseClient(iacExtInfo ...IacExtInfo) *tse.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.tseConn != nil {
		me.tseConn.WithHttpTransport(&logRoundTripper)
		return me.tseConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.tseConn, _ = tse.NewClient(me.Credential, me.Region, cpf)
	me.tseConn.WithHttpTransport(&logRoundTripper)

	return me.tseConn
}

// UseCdwchClient returns cdwch client for service
func (me *TencentCloudClient) UseCdwchClient() *cdwch.Client {
	if me.cdwchConn != nil {
		return me.cdwchConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.cdwchConn, _ = cdwch.NewClient(me.Credential, me.Region, cpf)
	me.cdwchConn.WithHttpTransport(&LogRoundTripper{})

	return me.cdwchConn
}

// UseEbClient returns eb client for service
func (me *TencentCloudClient) UseEbClient() *eb.Client {
	if me.ebConn != nil {
		return me.ebConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.ebConn, _ = eb.NewClient(me.Credential, me.Region, cpf)
	me.ebConn.WithHttpTransport(&LogRoundTripper{})

	return me.ebConn
}

// UseDlcClient returns eb client for service
func (me *TencentCloudClient) UseDlcClient() *dlc.Client {
	if me.dlcConn != nil {
		return me.dlcConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.dlcConn, _ = dlc.NewClient(me.Credential, me.Region, cpf)
	me.dlcConn.WithHttpTransport(&LogRoundTripper{})

	return me.dlcConn
}

// UseWedataClient returns eb client for service
func (me *TencentCloudClient) UseWedataClient() *wedata.Client {
	if me.wedataConn != nil {
		return me.wedataConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.wedataConn, _ = wedata.NewClient(me.Credential, me.Region, cpf)
	me.wedataConn.WithHttpTransport(&LogRoundTripper{})

	return me.wedataConn
}

func (me *TencentCloudClient) UseWafClient(iacExtInfo ...IacExtInfo) *waf.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.wafConn != nil {
		me.wafConn.WithHttpTransport(&logRoundTripper)
		return me.wafConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.wafConn, _ = waf.NewClient(me.Credential, me.Region, cpf)
	me.wafConn.WithHttpTransport(&logRoundTripper)

	return me.wafConn
}

func (me *TencentCloudClient) UseCfwClient(iacExtInfo ...IacExtInfo) *cfw.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.cfwConn != nil {
		me.cfwConn.WithHttpTransport(&logRoundTripper)
		return me.cfwConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.cfwConn, _ = cfw.NewClient(me.Credential, me.Region, cpf)
	me.cfwConn.WithHttpTransport(&logRoundTripper)

	return me.cfwConn
}

func (me *TencentCloudClient) UseOceanusClient() *oceanus.Client {
	if me.oceanusConn != nil {
		return me.oceanusConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.oceanusConn, _ = oceanus.NewClient(me.Credential, me.Region, cpf)
	me.oceanusConn.WithHttpTransport(&LogRoundTripper{})

	return me.oceanusConn
}

func (me *TencentCloudClient) UseDasbClient() *dasb.Client {
	if me.dasbConn != nil {
		return me.dasbConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.dasbConn, _ = dasb.NewClient(me.Credential, me.Region, cpf)
	me.dasbConn.WithHttpTransport(&LogRoundTripper{})

	return me.dasbConn
}

// UseTrocketClient returns trocket client for service
func (me *TencentCloudClient) UseTrocketClient() *trocket.Client {
	if me.trocketConn != nil {
		return me.trocketConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.trocketConn, _ = trocket.NewClient(me.Credential, me.Region, cpf)
	me.trocketConn.WithHttpTransport(&LogRoundTripper{})

	return me.trocketConn
}

// UseBiClient returns bi client for service
func (me *TencentCloudClient) UseBiClient() *bi.Client {
	if me.biConn != nil {
		return me.biConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.biConn, _ = bi.NewClient(me.Credential, me.Region, cpf)
	me.biConn.WithHttpTransport(&LogRoundTripper{})

	return me.biConn
}

// UseCdwpgClient returns cdwpg client for service
func (me *TencentCloudClient) UseCdwpgClient() *cdwpg.Client {
	if me.cdwpgConn != nil {
		return me.cdwpgConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.cdwpgConn, _ = cdwpg.NewClient(me.Credential, me.Region, cpf)
	me.cdwpgConn.WithHttpTransport(&LogRoundTripper{})

	return me.cdwpgConn
}

// UseCsipClient returns csip client for service
func (me *TencentCloudClient) UseCsipClient() *csip.Client {
	if me.csipConn != nil {
		return me.csipConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.csipConn, _ = csip.NewClient(me.Credential, me.Region, cpf)
	me.csipConn.WithHttpTransport(&LogRoundTripper{})

	return me.csipConn
}

// UseRegionClient returns region client for service
func (me *TencentCloudClient) UseRegionClient() *region.Client {
	if me.regionConn != nil {
		return me.regionConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.regionConn, _ = region.NewClient(me.Credential, me.Region, cpf)
	me.regionConn.WithHttpTransport(&LogRoundTripper{})

	return me.regionConn
}

//internal version: replace useClient begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace useClient end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

func getEnvDefault(key string, defVal int) int {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	timeOut, err := strconv.Atoi(val)
	if err != nil {
		panic("TENCENTCLOUD_XXX_REQUEST_TIMEOUT must be int.")
	}
	return timeOut
}

// UseTke2Client returns tke client for service
func (me *TencentCloudClient) UseTke2Client(iacExtInfo ...IacExtInfo) *tkev20220501.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.tkev20220501Conn != nil {
		me.tkev20220501Conn.WithHttpTransport(&logRoundTripper)
		return me.tkev20220501Conn
	}

	cpf := me.NewClientProfile(300)
	me.tkev20220501Conn, _ = tkev20220501.NewClient(me.Credential, me.Region, cpf)
	me.tkev20220501Conn.WithHttpTransport(&logRoundTripper)

	return me.tkev20220501Conn
}

// UseTkeV20220501Client returns tke client for service
func (me *TencentCloudClient) UseTkeV20220501Client(iacExtInfo ...IacExtInfo) *tkev20220501.Client {
	var logRoundTripper LogRoundTripper
	if len(iacExtInfo) != 0 {
		logRoundTripper.InstanceId = iacExtInfo[0].InstanceId
	}

	if me.tkev20220501Conn != nil {
		me.tkev20220501Conn.WithHttpTransport(&logRoundTripper)
		return me.tkev20220501Conn
	}

	cpf := me.NewClientProfile(300)
	me.tkev20220501Conn, _ = tkev20220501.NewClient(me.Credential, me.Region, cpf)
	me.tkev20220501Conn.WithHttpTransport(&logRoundTripper)

	return me.tkev20220501Conn
}

// UseCdcClient returns tem client for service
func (me *TencentCloudClient) UseCdcClient() *cdc.Client {
	if me.cdcConn != nil {
		return me.cdcConn
	}

	cpf := me.NewClientProfile(300)
	me.cdcConn, _ = cdc.NewClient(me.Credential, me.Region, cpf)
	me.cdcConn.WithHttpTransport(&LogRoundTripper{})

	return me.cdcConn
}

// UseCdwdoris return CDWDORIS client for service
func (me *TencentCloudClient) UseCdwdorisV20211228Client() *cdwdoris.Client {
	if me.cdwdorisConn != nil {
		return me.cdwdorisConn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.cdwdorisConn, _ = cdwdoris.NewClient(me.Credential, me.Region, cpf)
	me.cdwdorisConn.WithHttpTransport(&LogRoundTripper{})

	return me.cdwdorisConn
}

// UseControlcenter return CONTROLCENTER client for service
func (me *TencentCloudClient) UseControlcenterV20230110Client() *controlcenter.Client {
	if me.controlcenterConn != nil {
		return me.controlcenterConn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.controlcenterConn, _ = controlcenter.NewClient(me.Credential, me.Region, cpf)
	me.controlcenterConn.WithHttpTransport(&LogRoundTripper{})

	return me.controlcenterConn
}

// UseThpcClient return THPC client for service
func (me *TencentCloudClient) UseThpcV20230321Client() *thpc.Client {
	if me.thpcConn != nil {
		return me.thpcConn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.thpcConn, _ = thpc.NewClient(me.Credential, me.Region, cpf)
	me.thpcConn.WithHttpTransport(&LogRoundTripper{})

	return me.thpcConn
}

// UseEmrV20190103Client return EMR client for service
func (me *TencentCloudClient) UseEmrV20190103Client() *emr.Client {
	if me.emrv20190103Conn != nil {
		return me.emrv20190103Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.emrv20190103Conn, _ = emr.NewClient(me.Credential, me.Region, cpf)
	me.emrv20190103Conn.WithHttpTransport(&LogRoundTripper{})

	return me.emrv20190103Conn
}

// UseTeoV20220901Client return TEO client for service
func (me *TencentCloudClient) UseTeoV20220901Client() *teo.Client {
	if me.teov20220901Conn != nil {
		return me.teov20220901Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.teov20220901Conn, _ = teo.NewClient(me.Credential, me.Region, cpf)
	me.teov20220901Conn.WithHttpTransport(&LogRoundTripper{})

	return me.teov20220901Conn
}

// UseSslV20191205Client return SSL client for service
func (me *TencentCloudClient) UseSslV20191205Client() *sslCertificate.Client {
	if me.sslv20191205Conn != nil {
		return me.sslv20191205Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.sslv20191205Conn, _ = sslCertificate.NewClient(me.Credential, me.Region, cpf)
	me.sslv20191205Conn.WithHttpTransport(&LogRoundTripper{})

	return me.sslv20191205Conn
}

// UsePostgresV20170312Client return POSTGRES client for service
func (me *TencentCloudClient) UsePostgresV20170312Client() *postgre.Client {
	if me.postgresv20170312Conn != nil {
		return me.postgresv20170312Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.postgresv20170312Conn, _ = postgre.NewClient(me.Credential, me.Region, cpf)
	me.postgresv20170312Conn.WithHttpTransport(&LogRoundTripper{})

	return me.postgresv20170312Conn
}

// UseCfwV20190904Client return CFW client for service
func (me *TencentCloudClient) UseCfwV20190904Client() *cfw.Client {
	if me.cfwv20190904Conn != nil {
		return me.cfwv20190904Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.cfwv20190904Conn, _ = cfw.NewClient(me.Credential, me.Region, cpf)
	me.cfwv20190904Conn.WithHttpTransport(&LogRoundTripper{})

	return me.cfwv20190904Conn
}

// UseCcnV20170312Client return CCN client for service
func (me *TencentCloudClient) UseCcnV20170312Client() *vpc.Client {
	if me.ccnv20170312Conn != nil {
		return me.ccnv20170312Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.ccnv20170312Conn, _ = vpc.NewClient(me.Credential, me.Region, cpf)
	me.ccnv20170312Conn.WithHttpTransport(&LogRoundTripper{})

	return me.ccnv20170312Conn
}

// UseTcssV20201101Client return TCSS client for service
func (me *TencentCloudClient) UseTcssV20201101Client() *tcss.Client {
	if me.tcssv20201101Conn != nil {
		return me.tcssv20201101Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.tcssv20201101Conn, _ = tcss.NewClient(me.Credential, me.Region, cpf)
	me.tcssv20201101Conn.WithHttpTransport(&LogRoundTripper{})

	return me.tcssv20201101Conn
}

// UseCloudauditV20190319Client return CLOUDAUDIT client for service
func (me *TencentCloudClient) UseCloudauditV20190319Client() *audit.Client {
	if me.cloudauditv20190319Conn != nil {
		return me.cloudauditv20190319Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.cloudauditv20190319Conn, _ = audit.NewClient(me.Credential, me.Region, cpf)
	me.cloudauditv20190319Conn.WithHttpTransport(&LogRoundTripper{})

	return me.cloudauditv20190319Conn
}

// UsePrivatednsV20201028Client return PRIVATEDNS client for service
func (me *TencentCloudClient) UsePrivatednsV20201028Client() *privatedns.Client {
	if me.privatednsv20201028Conn != nil {
		return me.privatednsv20201028Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.privatednsv20201028Conn, _ = privatedns.NewClient(me.Credential, me.Region, cpf)
	me.privatednsv20201028Conn.WithHttpTransport(&LogRoundTripper{})

	return me.privatednsv20201028Conn
}

// UseWafV20180125Client return WAF client for service
func (me *TencentCloudClient) UseWafV20180125Client() *waf.Client {
	if me.wafv20180125Conn != nil {
		return me.wafv20180125Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.wafv20180125Conn, _ = waf.NewClient(me.Credential, me.Region, cpf)
	me.wafv20180125Conn.WithHttpTransport(&LogRoundTripper{})

	return me.wafv20180125Conn
}

// UseCamV20190116Client return CAM client for service
func (me *TencentCloudClient) UseCamV20190116Client() *cam.Client {
	if me.camv20190116Conn != nil {
		return me.camv20190116Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.camv20190116Conn, _ = cam.NewClient(me.Credential, me.Region, cpf)
	me.camv20190116Conn.WithHttpTransport(&LogRoundTripper{})

	return me.camv20190116Conn
}

// UseClsV20201016Client return CLS client for service
func (me *TencentCloudClient) UseClsV20201016Client() *cls.Client {
	if me.clsv20201016Conn != nil {
		return me.clsv20201016Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.clsv20201016Conn, _ = cls.NewClient(me.Credential, me.Region, cpf)
	me.clsv20201016Conn.WithHttpTransport(&LogRoundTripper{})

	return me.clsv20201016Conn
}

// UsePostgresqlV20170312Client return POSTGRESQL client for service
func (me *TencentCloudClient) UsePostgresqlV20170312Client() *postgre.Client {
	if me.postgresqlv20170312Conn != nil {
		return me.postgresqlv20170312Conn
	}
	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.postgresqlv20170312Conn, _ = postgre.NewClient(me.Credential, me.Region, cpf)
	me.postgresqlv20170312Conn.WithHttpTransport(&LogRoundTripper{})

	return me.postgresqlv20170312Conn
}
