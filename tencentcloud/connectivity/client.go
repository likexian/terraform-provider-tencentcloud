package connectivity

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	dasb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dasb/v20191018"

	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"

	ciam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ciam/v20220331"

	"github.com/tencentyun/cos-go-sdk-v5"

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
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"
	cdwpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwpg/v20201230"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	chdfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/chdfs/v20201112"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
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
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	postgre "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
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
	tdcpg "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdcpg/v20211118"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wss/v20180426"
)

//internal version: replace import begin
//internal version: replace import end

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

	cosConn            *s3.S3
	tencentCosConn     *cos.Client
	mysqlConn          *cdb.Client
	redisConn          *redis.Client
	asConn             *as.Client
	vpcConn            *vpc.Client
	cbsConn            *cbs.Client
	cvmConn            *cvm.Client
	clbConn            *clb.Client
	dayuConn           *dayu.Client
	dcConn             *dc.Client
	tagConn            *tag.Client
	mongodbConn        *mongodb.Client
	tkeConn            *tke.Client
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
	//internal version: replace client begin
	//internal version: replace client end
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

// UseCosClient returns cos client for service
func (me *TencentCloudClient) UseCosClient() *s3.S3 {
	if me.cosConn != nil {
		return me.cosConn
	}

	resolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL:           fmt.Sprintf("https://cos.%s.myqcloud.com", region),
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

// UseTencentCosClient tencent cloud own client for service instead of aws
func (me *TencentCloudClient) UseTencentCosClient(bucket string) *cos.Client {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, me.Region))

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
func (me *TencentCloudClient) UseMysqlClient() *cdb.Client {
	if me.mysqlConn != nil {
		return me.mysqlConn
	}

	cpf := me.NewClientProfile(300)
	me.mysqlConn, _ = cdb.NewClient(me.Credential, me.Region, cpf)
	me.mysqlConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseVpcClient() *vpc.Client {
	if me.vpcConn != nil {
		return me.vpcConn
	}

	cpf := me.NewClientProfile(300)
	me.vpcConn, _ = vpc.NewClient(me.Credential, me.Region, cpf)
	me.vpcConn.WithHttpTransport(&LogRoundTripper{})

	return me.vpcConn
}

// UseCbsClient returns cbs client for service
func (me *TencentCloudClient) UseCbsClient() *cbs.Client {
	if me.cbsConn != nil {
		return me.cbsConn
	}

	var reqTimeout = getEnvDefault(PROVIDER_CBS_REQUEST_TIMEOUT, 300)
	cpf := me.NewClientProfile(reqTimeout)
	me.cbsConn, _ = cbs.NewClient(me.Credential, me.Region, cpf)
	me.cbsConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseMongodbClient() *mongodb.Client {
	if me.mongodbConn != nil {
		return me.mongodbConn
	}

	cpf := me.NewClientProfile(300)
	me.mongodbConn, _ = mongodb.NewClient(me.Credential, me.Region, cpf)
	me.mongodbConn.WithHttpTransport(&LogRoundTripper{})

	return me.mongodbConn
}

// UseClbClient returns clb client for service
func (me *TencentCloudClient) UseClbClient() *clb.Client {
	if me.clbConn != nil {
		return me.clbConn
	}

	cpf := me.NewClientProfile(300)
	me.clbConn, _ = clb.NewClient(me.Credential, me.Region, cpf)
	me.clbConn.WithHttpTransport(&LogRoundTripper{})

	return me.clbConn
}

// UseCvmClient returns cvm client for service
func (me *TencentCloudClient) UseCvmClient() *cvm.Client {
	if me.cvmConn != nil {
		return me.cvmConn
	}

	var reqTimeout = getEnvDefault(PROVIDER_CVM_REQUEST_TIMEOUT, 300)
	cpf := me.NewClientProfile(reqTimeout)
	me.cvmConn, _ = cvm.NewClient(me.Credential, me.Region, cpf)
	me.cvmConn.WithHttpTransport(&LogRoundTripper{})

	return me.cvmConn
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
func (me *TencentCloudClient) UseTkeClient() *tke.Client {
	if me.tkeConn != nil {
		return me.tkeConn
	}

	cpf := me.NewClientProfile(300)
	me.tkeConn, _ = tke.NewClient(me.Credential, me.Region, cpf)
	me.tkeConn.WithHttpTransport(&LogRoundTripper{})

	return me.tkeConn
}

// UseTdmqClient returns Tdmq client for service
func (me *TencentCloudClient) UseTdmqClient() *tdmq.Client {
	if me.tdmqConn != nil {
		return me.tdmqConn
	}

	cpf := me.NewClientProfile(300)
	me.tdmqConn, _ = tdmq.NewClient(me.Credential, me.Region, cpf)
	me.tdmqConn.WithHttpTransport(&LogRoundTripper{})

	return me.tdmqConn
}

// UseGaapClient returns gaap client for service
func (me *TencentCloudClient) UseGaapClient() *gaap.Client {
	if me.gaapConn != nil {
		return me.gaapConn
	}

	cpf := me.NewClientProfile(300)
	me.gaapConn, _ = gaap.NewClient(me.Credential, me.Region, cpf)
	me.gaapConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseStsClient() *sts.Client {
	/*
		me.Credential will changed, don't cache it
		if me.stsConn != nil {
			return me.stsConn
		}
	*/

	cpf := me.NewClientProfile(300)
	me.stsConn, _ = sts.NewClient(me.Credential, me.Region, cpf)
	me.stsConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseScfClient() *scf.Client {
	if me.scfConn != nil {
		return me.scfConn
	}

	cpf := me.NewClientProfile(300)
	me.scfConn, _ = scf.NewClient(me.Credential, me.Region, cpf)
	me.scfConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseCdnClient() *cdn.Client {
	if me.cdnConn != nil {
		return me.cdnConn
	}

	cpf := me.NewClientProfile(300)
	me.cdnConn, _ = cdn.NewClient(me.Credential, me.Region, cpf)
	me.cdnConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseEsClient() *es.Client {
	if me.esConn != nil {
		return me.esConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.esConn, _ = es.NewClient(me.Credential, me.Region, cpf)
	me.esConn.WithHttpTransport(&LogRoundTripper{})

	return me.esConn
}

// UsePostgresqlClient returns postgresql client for service
func (me *TencentCloudClient) UsePostgresqlClient() *postgre.Client {
	if me.postgreConn != nil {
		return me.postgreConn
	}

	cpf := me.NewClientProfile(300)
	me.postgreConn, _ = postgre.NewClient(me.Credential, me.Region, cpf)
	me.postgreConn.WithHttpTransport(&LogRoundTripper{})

	return me.postgreConn
}

// UseSqlserverClient returns sqlserver client for service
func (me *TencentCloudClient) UseSqlserverClient() *sqlserver.Client {
	if me.sqlserverConn != nil {
		return me.sqlserverConn
	}

	cpf := me.NewClientProfile(300)
	me.sqlserverConn, _ = sqlserver.NewClient(me.Credential, me.Region, cpf)
	me.sqlserverConn.WithHttpTransport(&LogRoundTripper{})

	return me.sqlserverConn
}

// UseCkafkaClient returns ckafka client for service
func (me *TencentCloudClient) UseCkafkaClient() *ckafka.Client {
	if me.ckafkaConn != nil {
		return me.ckafkaConn
	}

	cpf := me.NewClientProfile(300)
	me.ckafkaConn, _ = ckafka.NewClient(me.Credential, me.Region, cpf)
	me.ckafkaConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseTCRClient() *tcr.Client {
	if me.tcrConn != nil {
		return me.tcrConn
	}

	cpf := me.NewClientProfile(300)
	me.tcrConn, _ = tcr.NewClient(me.Credential, me.Region, cpf)
	me.tcrConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseClsClient() *cls.Client {
	if me.clsConn != nil {
		return me.clsConn
	}
	cpf := me.NewClientProfile(300)
	me.clsConn, _ = cls.NewClient(me.Credential, me.Region, cpf)
	me.clsConn.WithHttpTransport(&LogRoundTripper{})

	return me.clsConn
}

// UseLighthouseClient return Lighthouse client for service
func (me *TencentCloudClient) UseLighthouseClient() *lighthouse.Client {
	if me.lighthouseConn != nil {
		return me.lighthouseConn
	}
	cpf := me.NewClientProfile(300)
	me.lighthouseConn, _ = lighthouse.NewClient(me.Credential, me.Region, cpf)
	me.lighthouseConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UsePrivateDnsClient() *privatedns.Client {
	if me.privateDnsConn != nil {
		return me.privateDnsConn
	}
	cpf := me.NewClientProfile(300)
	me.privateDnsConn, _ = privatedns.NewClient(me.Credential, me.Region, cpf)
	me.privateDnsConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseTeoClient() *teo.Client {
	if me.teoConn != nil {
		return me.teoConn
	}

	cpf := me.NewClientProfile(300)
	me.teoConn, _ = teo.NewClient(me.Credential, me.Region, cpf)
	me.teoConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseMariadbClient() *mariadb.Client {
	if me.mariadbConn != nil {
		return me.mariadbConn
	}

	cpf := me.NewClientProfile(300)
	me.mariadbConn, _ = mariadb.NewClient(me.Credential, me.Region, cpf)
	me.mariadbConn.WithHttpTransport(&LogRoundTripper{})

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
func (me *TencentCloudClient) UseTdcpgClient() *tdcpg.Client {
	if me.tdcpgConn != nil {
		return me.tdcpgConn
	}

	cpf := me.NewClientProfile(300)
	me.tdcpgConn, _ = tdcpg.NewClient(me.Credential, me.Region, cpf)
	me.tdcpgConn.WithHttpTransport(&LogRoundTripper{})

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
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos-control.%s.myqcloud.com", uin, me.Region))

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
func (me *TencentCloudClient) UseTseClient() *tse.Client {
	if me.tseConn != nil {
		return me.tseConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.tseConn, _ = tse.NewClient(me.Credential, me.Region, cpf)
	me.tseConn.WithHttpTransport(&LogRoundTripper{})

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

func (me *TencentCloudClient) UseWafClient() *waf.Client {
	if me.wafConn != nil {
		return me.wafConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.wafConn, _ = waf.NewClient(me.Credential, me.Region, cpf)
	me.wafConn.WithHttpTransport(&LogRoundTripper{})

	return me.wafConn
}

func (me *TencentCloudClient) UseCfwClient() *cfw.Client {
	if me.cfwConn != nil {
		return me.cfwConn
	}

	cpf := me.NewClientProfile(300)
	cpf.Language = "zh-CN"
	me.cfwConn, _ = cfw.NewClient(me.Credential, me.Region, cpf)
	me.cfwConn.WithHttpTransport(&LogRoundTripper{})

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

//internal version: replace useClient begin
//internal version: replace useClient end

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
