package tencentcloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
	sdkcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	commonJson "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/json"
	sdkprofile "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sdksts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/apigateway"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/apm"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/audit"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bh"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bi"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cam"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cat"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ccn"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdc"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdh"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdn"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdwch"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdwdoris"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdwpg"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cfs"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cfw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/chdfs"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ci"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ciam"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ckafka"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/controlcenter"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cos"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/csip"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/css"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cwp"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayu"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayuv2"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbbrain"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dc"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dcdb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dcg"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dlc"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dnspod"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/domain"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dts"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/eb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/emr"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/es"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fl"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/gaap"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/kms"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/lighthouse"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mariadb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mdl"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mps"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/oceanus"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/pls"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/privatedns"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/project"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/pts"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/rum"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/scf"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ses"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sms"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sqlserver"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssl"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/sts"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tat"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcaplusdb"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcm"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcmg"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcmq"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tco"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcr"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcss"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdcpg"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tem"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/thpc"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tke"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tmp"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tpulsar"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/trabbit"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/trocket"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tse"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tsf"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vod"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpn"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/waf"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/wedata"
)

const (
	PROVIDER_SECRET_ID      = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY     = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_SECURITY_TOKEN = "TENCENTCLOUD_SECURITY_TOKEN"
	PROVIDER_REGION         = "TENCENTCLOUD_REGION"
	PROVIDER_PROTOCOL       = "TENCENTCLOUD_PROTOCOL"
	PROVIDER_DOMAIN         = "TENCENTCLOUD_DOMAIN"
	PROVIDER_COS_DOMAIN     = "TENCENTCLOUD_COS_DOMAIN"
	//internal version: replace envYunti begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace envYunti end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	PROVIDER_ASSUME_ROLE_ARN                = "TENCENTCLOUD_ASSUME_ROLE_ARN"
	PROVIDER_ASSUME_ROLE_SESSION_NAME       = "TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME"
	PROVIDER_ASSUME_ROLE_SESSION_DURATION   = "TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION"
	PROVIDER_ASSUME_ROLE_EXTERNAL_ID        = "TENCENTCLOUD_ASSUME_ROLE_EXTERNAL_ID"
	PROVIDER_ASSUME_ROLE_SAML_ASSERTION     = "TENCENTCLOUD_ASSUME_ROLE_SAML_ASSERTION"
	PROVIDER_ASSUME_ROLE_PRINCIPAL_ARN      = "TENCENTCLOUD_ASSUME_ROLE_PRINCIPAL_ARN"
	PROVIDER_ASSUME_ROLE_WEB_IDENTITY_TOKEN = "TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN"
	PROVIDER_SHARED_CREDENTIALS_DIR         = "TENCENTCLOUD_SHARED_CREDENTIALS_DIR"
	PROVIDER_PROFILE                        = "TENCENTCLOUD_PROFILE"
	PROVIDER_CAM_ROLE_NAME                  = "TENCENTCLOUD_CAM_ROLE_NAME"
	POD_OIDC_TKE_REGION                     = "TKE_REGION"
	POD_OIDC_TKE_WEB_IDENTITY_TOKEN_FILE    = "TKE_WEB_IDENTITY_TOKEN_FILE"
	POD_OIDC_TKE_PROVIDER_ID                = "TKE_PROVIDER_ID"
	POD_OIDC_TKE_ROLE_ARN                   = "TKE_ROLE_ARN"
)

const (
	DEFAULT_REGION  = "ap-guangzhou"
	DEFAULT_PROFILE = "default"
)

type TencentCloudClient struct {
	apiV3Conn *connectivity.TencentCloudClient
}

var _ tccommon.ProviderMeta = &TencentCloudClient{}

func init() {
	commonJson.OmitBehaviour = commonJson.OmitEmpty
}

// GetAPIV3Conn 返回访问云 API 的客户端连接对象
func (meta *TencentCloudClient) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return meta.apiV3Conn
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"secret_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_ID, nil),
				Description: "This is the TencentCloud access key. It can also be sourced from the `TENCENTCLOUD_SECRET_ID` environment variable.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_KEY, nil),
				Description: "This is the TencentCloud secret key. It can also be sourced from the `TENCENTCLOUD_SECRET_KEY` environment variable.",
				Sensitive:   true,
			},
			"security_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECURITY_TOKEN, nil),
				Description: "TencentCloud Security Token of temporary access credentials. It can be sourced from the `TENCENTCLOUD_SECURITY_TOKEN` environment variable. Notice: for supported products, please refer to: [temporary key supported products](https://intl.cloud.tencent.com/document/product/598/10588).",
				Sensitive:   true,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_REGION, nil),
				Description: "This is the TencentCloud region. It can also be sourced from the `TENCENTCLOUD_REGION` environment variables. The default input value is ap-guangzhou.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(PROVIDER_PROTOCOL, "HTTPS"),
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"HTTP", "HTTPS"}),
				Description:  "The protocol of the API request. Valid values: `HTTP` and `HTTPS`. Default is `HTTPS`.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_DOMAIN, nil),
				Description: "The root domain of the API request, Default is `tencentcloudapi.com`.",
			},
			"cos_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_COS_DOMAIN, nil),
				Description: "The cos domain of the API request, Default is `https://cos.{region}.myqcloud.com`, Other Examples: `https://cluster-123456.cos-cdc.ap-guangzhou.myqcloud.com`.",
			},
			//internal version: replace enableBpass begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			//internal version: replace enableBpass end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			"assume_role": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "The `assume_role` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_ARN, nil),
							Description: "The ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN`.",
						},
						"session_name": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_SESSION_NAME, nil),
							Description: "The session name to use when making the AssumeRole call. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`.",
						},
						"session_duration": {
							Type:     schema.TypeInt,
							Required: true,
							DefaultFunc: func() (interface{}, error) {
								if v := os.Getenv(PROVIDER_ASSUME_ROLE_SESSION_DURATION); v != "" {
									return strconv.Atoi(v)
								}
								return 7200, nil
							},
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 43200),
							Description:  "The duration of the session when making the AssumeRole call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`.",
						},
						"policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A more restrictive policy when making the AssumeRole call. Its content must not contains `principal` elements. Notice: more syntax references, please refer to: [policies syntax logic](https://intl.cloud.tencent.com/document/product/598/10603).",
						},
						"external_id": {
							Type:        schema.TypeString,
							Optional:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_EXTERNAL_ID, nil),
							Description: "External role ID, which can be obtained by clicking the role name in the CAM console. It can contain 2-128 letters, digits, and symbols (=,.@:/-). Regex: [\\w+=,.@:/-]*. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_EXTERNAL_ID`.",
						},
					},
				},
			},
			"assume_role_with_saml": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"assume_role_with_web_identity"},
				Description:   "The `assume_role_with_saml` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"saml_assertion": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_SAML_ASSERTION, nil),
							Description: "SAML assertion information encoded in base64. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SAML_ASSERTION`.",
						},
						"principal_arn": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_PRINCIPAL_ARN, nil),
							Description: "Player Access Description Name. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_PRINCIPAL_ARN`.",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_ARN, nil),
							Description: "The ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN`.",
						},
						"session_name": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_SESSION_NAME, nil),
							Description: "The session name to use when making the AssumeRole call. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`.",
						},
						"session_duration": {
							Type:     schema.TypeInt,
							Required: true,
							DefaultFunc: func() (interface{}, error) {
								if v := os.Getenv(PROVIDER_ASSUME_ROLE_SESSION_DURATION); v != "" {
									return strconv.Atoi(v)
								}
								return 7200, nil
							},
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 43200),
							Description:  "The duration of the session when making the AssumeRoleWithSAML call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`.",
						},
					},
				},
			},
			"enable_pod_oidc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable pod oidc.",
			},
			"assume_role_with_web_identity": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"assume_role_with_saml"},
				Description:   "The `assume_role_with_web_identity` block. If provided, terraform will attempt to assume this role using the supplied credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"web_identity_token": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_WEB_IDENTITY_TOKEN, nil),
							Description: "OIDC token issued by IdP. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN`.",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_ARN, nil),
							Description: "The ARN of the role to assume. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_ARN`.",
						},
						"session_name": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ASSUME_ROLE_SESSION_NAME, nil),
							Description: "The session name to use when making the AssumeRole call. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_NAME`.",
						},
						"session_duration": {
							Type:     schema.TypeInt,
							Required: true,
							DefaultFunc: func() (interface{}, error) {
								if v := os.Getenv(PROVIDER_ASSUME_ROLE_SESSION_DURATION); v != "" {
									return strconv.Atoi(v)
								}
								return 7200, nil
							},
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 43200),
							Description:  "The duration of the session when making the AssumeRoleWithWebIdentity call. Its value ranges from 0 to 43200(seconds), and default is 7200 seconds. It can be sourced from the `TENCENTCLOUD_ASSUME_ROLE_SESSION_DURATION`.",
						},
					},
				},
			},
			"shared_credentials_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SHARED_CREDENTIALS_DIR, nil),
				Description: "The directory of the shared credentials. It can also be sourced from the `TENCENTCLOUD_SHARED_CREDENTIALS_DIR` environment variable. If not set this defaults to ~/.tccli.",
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_PROFILE, nil),
				Description: "The profile name as set in the shared credentials. It can also be sourced from the `TENCENTCLOUD_PROFILE` environment variable. If not set, the default profile created with `tccli configure` will be used.",
			},
			"cam_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_CAM_ROLE_NAME, nil),
				Description: "The name of the CVM instance CAM role. It can be sourced from the `TENCENTCLOUD_CAM_ROLE_NAME` environment variable.",
			},
			"allowed_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"forbidden_account_ids", "assume_role_with_saml", "assume_role_with_web_identity"},
				Description:   "List of allowed TencentCloud account IDs to prevent you from mistakenly using the wrong one (and potentially end up destroying a live environment). Conflicts with `forbidden_account_ids`, If use `assume_role_with_saml` or `assume_role_with_web_identity`, it is not supported.",
			},
			"forbidden_account_ids": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				ConflictsWith: []string{"allowed_account_ids", "assume_role_with_saml", "assume_role_with_web_identity"},
				Description:   "List of forbidden TencentCloud account IDs to prevent you from mistakenly using the wrong one (and potentially end up destroying a live environment). Conflicts with `allowed_account_ids`, If use `assume_role_with_saml` or `assume_role_with_web_identity`, it is not supported.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"tencentcloud_availability_regions":                         common.DataSourceTencentCloudAvailabilityRegions(),
			"tencentcloud_emr":                                          emr.DataSourceTencentCloudEmr(),
			"tencentcloud_emr_nodes":                                    emr.DataSourceTencentCloudEmrNodes(),
			"tencentcloud_emr_cvm_quota":                                emr.DataSourceTencentCloudEmrCvmQuota(),
			"tencentcloud_emr_auto_scale_records":                       emr.DataSourceTencentCloudEmrAutoScaleRecords(),
			"tencentcloud_availability_zones":                           common.DataSourceTencentCloudAvailabilityZones(),
			"tencentcloud_availability_zones_by_product":                common.DataSourceTencentCloudAvailabilityZonesByProduct(),
			"tencentcloud_projects":                                     project.DataSourceTencentCloudProjects(),
			"tencentcloud_instances":                                    cvm.DataSourceTencentCloudInstances(),
			"tencentcloud_instances_set":                                cvm.DataSourceTencentCloudInstancesSet(),
			"tencentcloud_reserved_instances":                           cvm.DataSourceTencentCloudReservedInstances(),
			"tencentcloud_placement_groups":                             cvm.DataSourceTencentCloudPlacementGroups(),
			"tencentcloud_key_pairs":                                    cvm.DataSourceTencentCloudKeyPairs(),
			"tencentcloud_image":                                        cvm.DataSourceTencentCloudImage(),
			"tencentcloud_images":                                       cvm.DataSourceTencentCloudImages(),
			"tencentcloud_image_from_family":                            cvm.DataSourceTencentCloudImageFromFamily(),
			"tencentcloud_instance_types":                               cvm.DataSourceTencentCloudInstanceTypes(),
			"tencentcloud_reserved_instance_configs":                    cvm.DataSourceTencentCloudReservedInstanceConfigs(),
			"tencentcloud_vpc_instances":                                vpc.DataSourceTencentCloudVpcInstances(),
			"tencentcloud_vpc_subnets":                                  vpc.DataSourceTencentCloudVpcSubnets(),
			"tencentcloud_vpc_route_tables":                             vpc.DataSourceTencentCloudVpcRouteTables(),
			"tencentcloud_vpc":                                          vpc.DataSourceTencentCloudVpc(),
			"tencentcloud_vpc_acls":                                     vpc.DataSourceTencentCloudVpcAcls(),
			"tencentcloud_vpc_bandwidth_package_quota":                  vpc.DataSourceTencentCloudVpcBandwidthPackageQuota(),
			"tencentcloud_vpc_bandwidth_package_bill_usage":             vpc.DataSourceTencentCloudVpcBandwidthPackageBillUsage(),
			"tencentcloud_vpc_account_attributes":                       vpc.DataSourceTencentCloudVpcAccountAttributes(),
			"tencentcloud_vpc_classic_link_instances":                   vpc.DataSourceTencentCloudVpcClassicLinkInstances(),
			"tencentcloud_vpc_gateway_flow_monitor_detail":              vpc.DataSourceTencentCloudVpcGatewayFlowMonitorDetail(),
			"tencentcloud_vpc_gateway_flow_qos":                         vpc.DataSourceTencentCloudVpcGatewayFlowQos(),
			"tencentcloud_vpc_cvm_instances":                            vpc.DataSourceTencentCloudVpcCvmInstances(),
			"tencentcloud_vpc_net_detect_states":                        vpc.DataSourceTencentCloudVpcNetDetectStates(),
			"tencentcloud_vpc_network_interface_limit":                  vpc.DataSourceTencentCloudVpcNetworkInterfaceLimit(),
			"tencentcloud_vpc_private_ip_addresses":                     vpc.DataSourceTencentCloudVpcPrivateIpAddresses(),
			"tencentcloud_vpc_product_quota":                            vpc.DataSourceTencentCloudVpcProductQuota(),
			"tencentcloud_vpc_resource_dashboard":                       vpc.DataSourceTencentCloudVpcResourceDashboard(),
			"tencentcloud_vpc_route_conflicts":                          vpc.DataSourceTencentCloudVpcRouteConflicts(),
			"tencentcloud_vpc_security_group_limits":                    vpc.DataSourceTencentCloudVpcSecurityGroupLimits(),
			"tencentcloud_vpc_security_group_references":                vpc.DataSourceTencentCloudVpcSecurityGroupReferences(),
			"tencentcloud_vpc_sg_snapshot_file_content":                 vpc.DataSourceTencentCloudVpcSgSnapshotFileContent(),
			"tencentcloud_vpc_snapshot_files":                           vpc.DataSourceTencentCloudVpcSnapshotFiles(),
			"tencentcloud_vpc_subnet_resource_dashboard":                vpc.DataSourceTencentCloudVpcSubnetResourceDashboard(),
			"tencentcloud_vpc_template_limits":                          vpc.DataSourceTencentCloudVpcTemplateLimits(),
			"tencentcloud_vpc_limits":                                   vpc.DataSourceTencentCloudVpcLimits(),
			"tencentcloud_vpc_used_ip_address":                          vpc.DataSourceTencentCloudVpcUsedIpAddress(),
			"tencentcloud_vpc_net_detect_state_check":                   vpc.DataSourceTencentCloudVpcNetDetectStateCheck(),
			"tencentcloud_subnet":                                       vpc.DataSourceTencentCloudSubnet(),
			"tencentcloud_route_table":                                  vpc.DataSourceTencentCloudRouteTable(),
			"tencentcloud_domains":                                      domain.DataSourceTencentCloudDomains(),
			"tencentcloud_eip":                                          cvm.DataSourceTencentCloudEip(),
			"tencentcloud_eips":                                         cvm.DataSourceTencentCloudEips(),
			"tencentcloud_eip_address_quota":                            cvm.DataSourceTencentCloudEipAddressQuota(),
			"tencentcloud_eip_network_account_type":                     cvm.DataSourceTencentCloudEipNetworkAccountType(),
			"tencentcloud_enis":                                         vpc.DataSourceTencentCloudEnis(),
			"tencentcloud_nats":                                         vpc.DataSourceTencentCloudNats(),
			"tencentcloud_dnats":                                        vpc.DataSourceTencentCloudDnats(),
			"tencentcloud_nat_gateways":                                 vpc.DataSourceTencentCloudNatGateways(),
			"tencentcloud_nat_gateway_snats":                            vpc.DataSourceTencentCloudNatGatewaySnats(),
			"tencentcloud_nat_dc_route":                                 vpc.DataSourceTencentCloudNatDcRoute(),
			"tencentcloud_oceanus_resource_related_job":                 oceanus.DataSourceTencentCloudOceanusResourceRelatedJob(),
			"tencentcloud_oceanus_savepoint_list":                       oceanus.DataSourceTencentCloudOceanusSavepointList(),
			"tencentcloud_oceanus_system_resource":                      oceanus.DataSourceTencentCloudOceanusSystemResource(),
			"tencentcloud_oceanus_work_spaces":                          oceanus.DataSourceTencentCloudOceanusWorkSpaces(),
			"tencentcloud_oceanus_clusters":                             oceanus.DataSourceTencentCloudOceanusClusters(),
			"tencentcloud_oceanus_tree_jobs":                            oceanus.DataSourceTencentCloudOceanusTreeJobs(),
			"tencentcloud_oceanus_tree_resources":                       oceanus.DataSourceTencentCloudOceanusTreeResources(),
			"tencentcloud_oceanus_job_submission_log":                   oceanus.DataSourceTencentCloudOceanusJobSubmissionLog(),
			"tencentcloud_oceanus_check_savepoint":                      oceanus.DataSourceTencentCloudOceanusCheckSavepoint(),
			"tencentcloud_oceanus_job_events":                           oceanus.DataSourceTencentCloudOceanusJobEvents(),
			"tencentcloud_oceanus_meta_table":                           oceanus.DataSourceTencentCloudOceanusMetaTable(),
			"tencentcloud_vpn_customer_gateways":                        vpn.DataSourceTencentCloudVpnCustomerGateways(),
			"tencentcloud_vpn_gateways":                                 vpn.DataSourceTencentCloudVpnGateways(),
			"tencentcloud_vpn_gateway_routes":                           vpn.DataSourceTencentCloudVpnGatewayRoutes(),
			"tencentcloud_vpn_connections":                              vpn.DataSourceTencentCloudVpnConnections(),
			"tencentcloud_vpn_customer_gateway_vendors":                 vpn.DataSourceTencentCloudVpnCustomerGatewayVendors(),
			"tencentcloud_vpn_default_health_check_ip":                  vpn.DataSourceTencentCloudVpnDefaultHealthCheckIp(),
			"tencentcloud_ha_vips":                                      vpc.DataSourceTencentCloudHaVips(),
			"tencentcloud_ha_vip_eip_attachments":                       vpc.DataSourceTencentCloudHaVipEipAttachments(),
			"tencentcloud_ccn_instances":                                ccn.DataSourceTencentCloudCcnInstances(),
			"tencentcloud_ccn_routes":                                   ccn.DataSourceTencentCloudCcnRoutes(),
			"tencentcloud_ccn_bandwidth_limits":                         ccn.DataSourceTencentCloudCcnBandwidthLimits(),
			"tencentcloud_ccn_cross_border_compliance":                  ccn.DataSourceTencentCloudCcnCrossBorderCompliance(),
			"tencentcloud_ccn_tenant_instances":                         ccn.DataSourceTencentCloudCcnTenantInstance(),
			"tencentcloud_ccn_cross_border_flow_monitor":                ccn.DataSourceTencentCloudCcnCrossBorderFlowMonitor(),
			"tencentcloud_ccn_cross_border_region_bandwidth_limits":     ccn.DataSourceTencentCloudCcnCrossBorderRegionBandwidthLimits(),
			"tencentcloud_ccn_route_table_input_policies":               ccn.DataSourceTencentCloudCcnRouteTableInputPolicies(),
			"tencentcloud_dc_instances":                                 dc.DataSourceTencentCloudDcInstances(),
			"tencentcloud_dc_access_points":                             dc.DataSourceTencentCloudDcAccessPoints(),
			"tencentcloud_dc_internet_address_quota":                    dc.DataSourceTencentCloudDcInternetAddressQuota(),
			"tencentcloud_dc_internet_address_statistics":               dc.DataSourceTencentCloudDcInternetAddressStatistics(),
			"tencentcloud_dc_public_direct_connect_tunnel_routes":       dc.DataSourceTencentCloudDcPublicDirectConnectTunnelRoutes(),
			"tencentcloud_dcx_instances":                                dc.DataSourceTencentCloudDcxInstances(),
			"tencentcloud_dc_gateway_instances":                         dcg.DataSourceTencentCloudDcGatewayInstances(),
			"tencentcloud_dc_gateway_ccn_routes":                        dcg.DataSourceTencentCloudDcGatewayCCNRoutes(),
			"tencentcloud_security_group":                               vpc.DataSourceTencentCloudSecurityGroup(),
			"tencentcloud_security_groups":                              vpc.DataSourceTencentCloudSecurityGroups(),
			"tencentcloud_kubernetes_clusters":                          tke.DataSourceTencentCloudKubernetesClusters(),
			"tencentcloud_kubernetes_charts":                            tke.DataSourceTencentCloudKubernetesCharts(),
			"tencentcloud_kubernetes_cluster_levels":                    tke.DataSourceTencentCloudKubernetesClusterLevels(),
			"tencentcloud_kubernetes_cluster_common_names":              tke.DataSourceTencentCloudKubernetesClusterCommonNames(),
			"tencentcloud_kubernetes_cluster_authentication_options":    tke.DataSourceTencentCloudKubernetesClusterAuthenticationOptions(),
			"tencentcloud_kubernetes_available_cluster_versions":        tke.DataSourceTencentCloudKubernetesAvailableClusterVersions(),
			"tencentcloud_eks_clusters":                                 tke.DataSourceTencentCloudEKSClusters(),
			"tencentcloud_eks_cluster_credential":                       tke.DataSourceTencentCloudEksClusterCredential(),
			"tencentcloud_container_clusters":                           tke.DataSourceTencentCloudContainerClusters(),
			"tencentcloud_container_cluster_instances":                  tke.DataSourceTencentCloudContainerClusterInstances(),
			"tencentcloud_mysql_backup_list":                            cdb.DataSourceTencentCloudMysqlBackupList(),
			"tencentcloud_mysql_zone_config":                            cdb.DataSourceTencentCloudMysqlZoneConfig(),
			"tencentcloud_mysql_parameter_list":                         cdb.DataSourceTencentCloudMysqlParameterList(),
			"tencentcloud_mysql_default_params":                         cdb.DataSourceTencentCloudMysqlDefaultParams(),
			"tencentcloud_mysql_instance":                               cdb.DataSourceTencentCloudMysqlInstance(),
			"tencentcloud_mysql_backup_overview":                        cdb.DataSourceTencentCloudMysqlBackupOverview(),
			"tencentcloud_mysql_backup_summaries":                       cdb.DataSourceTencentCloudMysqlBackupSummaries(),
			"tencentcloud_mysql_bin_log":                                cdb.DataSourceTencentCloudMysqlBinLog(),
			"tencentcloud_mysql_binlog_backup_overview":                 cdb.DataSourceTencentCloudMysqlBinlogBackupOverview(),
			"tencentcloud_mysql_clone_list":                             cdb.DataSourceTencentCloudMysqlCloneList(),
			"tencentcloud_mysql_data_backup_overview":                   cdb.DataSourceTencentCloudMysqlDataBackupOverview(),
			"tencentcloud_mysql_db_features":                            cdb.DataSourceTencentCloudMysqlDbFeatures(),
			"tencentcloud_mysql_inst_tables":                            cdb.DataSourceTencentCloudMysqlInstTables(),
			"tencentcloud_mysql_instance_charset":                       cdb.DataSourceTencentCloudMysqlInstanceCharset(),
			"tencentcloud_mysql_instance_info":                          cdb.DataSourceTencentCloudMysqlInstanceInfo(),
			"tencentcloud_mysql_instance_param_record":                  cdb.DataSourceTencentCloudMysqlInstanceParamRecord(),
			"tencentcloud_mysql_instance_reboot_time":                   cdb.DataSourceTencentCloudMysqlInstanceRebootTime(),
			"tencentcloud_mysql_proxy_custom":                           cdb.DataSourceTencentCloudMysqlProxyCustom(),
			"tencentcloud_mysql_rollback_range_time":                    cdb.DataSourceTencentCloudMysqlRollbackRangeTime(),
			"tencentcloud_mysql_slow_log":                               cdb.DataSourceTencentCloudMysqlSlowLog(),
			"tencentcloud_mysql_slow_log_data":                          cdb.DataSourceTencentCloudMysqlSlowLogData(),
			"tencentcloud_mysql_supported_privileges":                   cdb.DataSourceTencentCloudMysqlSupportedPrivileges(),
			"tencentcloud_mysql_switch_record":                          cdb.DataSourceTencentCloudMysqlSwitchRecord(),
			"tencentcloud_mysql_user_task":                              cdb.DataSourceTencentCloudMysqlUserTask(),
			"tencentcloud_mysql_databases":                              cdb.DataSourceTencentCloudMysqlDatabases(),
			"tencentcloud_mysql_error_log":                              cdb.DataSourceTencentCloudMysqlErrorLog(),
			"tencentcloud_mysql_project_security_group":                 cdb.DataSourceTencentCloudMysqlProjectSecurityGroup(),
			"tencentcloud_mysql_ro_min_scale":                           cdb.DataSourceTencentCloudMysqlRoMinScale(),
			"tencentcloud_cos_bucket_object":                            cos.DataSourceTencentCloudCosBucketObject(),
			"tencentcloud_cos_buckets":                                  cos.DataSourceTencentCloudCosBuckets(),
			"tencentcloud_cos_batchs":                                   cos.DataSourceTencentCloudCosBatchs(),
			"tencentcloud_cos_bucket_inventorys":                        cos.DataSourceTencentCloudCosBucketInventorys(),
			"tencentcloud_cos_bucket_multipart_uploads":                 cos.DataSourceTencentCloudCosBucketMultipartUploads(),
			"tencentcloud_cfs_file_systems":                             cfs.DataSourceTencentCloudCfsFileSystems(),
			"tencentcloud_cfs_access_groups":                            cfs.DataSourceTencentCloudCfsAccessGroups(),
			"tencentcloud_cfs_access_rules":                             cfs.DataSourceTencentCloudCfsAccessRules(),
			"tencentcloud_cfs_mount_targets":                            cfs.DataSourceTencentCloudCfsMountTargets(),
			"tencentcloud_cfs_file_system_clients":                      cfs.DataSourceTencentCloudCfsFileSystemClients(),
			"tencentcloud_cfs_available_zone":                           cfs.DataSourceTencentCloudCfsAvailableZone(),
			"tencentcloud_redis_zone_config":                            crs.DataSourceTencentCloudRedisZoneConfig(),
			"tencentcloud_redis_instances":                              crs.DataSourceTencentCloudRedisInstances(),
			"tencentcloud_redis_backup":                                 crs.DataSourceTencentCloudRedisBackup(),
			"tencentcloud_redis_backup_download_info":                   crs.DataSourceTencentCloudRedisBackupDownloadInfo(),
			"tencentcloud_redis_param_records":                          crs.DataSourceTencentCloudRedisRecordsParam(),
			"tencentcloud_redis_instance_shards":                        crs.DataSourceTencentCloudRedisInstanceShards(),
			"tencentcloud_redis_instance_zone_info":                     crs.DataSourceTencentCloudRedisInstanceZoneInfo(),
			"tencentcloud_redis_instance_task_list":                     crs.DataSourceTencentCloudRedisInstanceTaskList(),
			"tencentcloud_redis_instance_node_info":                     crs.DataSourceTencentCloudRedisInstanceNodeInfo(),
			"tencentcloud_redis_clusters":                               crs.DataSourceTencentCloudRedisClusters(),
			"tencentcloud_as_scaling_configs":                           as.DataSourceTencentCloudAsScalingConfigs(),
			"tencentcloud_as_scaling_groups":                            as.DataSourceTencentCloudAsScalingGroups(),
			"tencentcloud_as_scaling_policies":                          as.DataSourceTencentCloudAsScalingPolicies(),
			"tencentcloud_cbs_storages":                                 cbs.DataSourceTencentCloudCbsStorages(),
			"tencentcloud_cbs_storages_set":                             cbs.DataSourceTencentCloudCbsStoragesSet(),
			"tencentcloud_cbs_snapshots":                                cbs.DataSourceTencentCloudCbsSnapshots(),
			"tencentcloud_cbs_snapshot_policies":                        cbs.DataSourceTencentCloudCbsSnapshotPolicies(),
			"tencentcloud_clb_instances":                                clb.DataSourceTencentCloudClbInstances(),
			"tencentcloud_clb_listeners":                                clb.DataSourceTencentCloudClbListeners(),
			"tencentcloud_clb_listener_rules":                           clb.DataSourceTencentCloudClbListenerRules(),
			"tencentcloud_clb_attachments":                              clb.DataSourceTencentCloudClbServerAttachments(),
			"tencentcloud_clb_redirections":                             clb.DataSourceTencentCloudClbRedirections(),
			"tencentcloud_clb_target_groups":                            clb.DataSourceTencentCloudClbTargetGroups(),
			"tencentcloud_clb_cluster_resources":                        clb.DataSourceTencentCloudClbClusterResources(),
			"tencentcloud_clb_cross_targets":                            clb.DataSourceTencentCloudClbCrossTargets(),
			"tencentcloud_clb_exclusive_clusters":                       clb.DataSourceTencentCloudClbExclusiveClusters(),
			"tencentcloud_clb_idle_instances":                           clb.DataSourceTencentCloudClbIdleInstances(),
			"tencentcloud_clb_listeners_by_targets":                     clb.DataSourceTencentCloudClbListenersByTargets(),
			"tencentcloud_clb_instance_by_cert_id":                      clb.DataSourceTencentCloudClbInstanceByCertId(),
			"tencentcloud_clb_instance_traffic":                         clb.DataSourceTencentCloudClbInstanceTraffic(),
			"tencentcloud_clb_instance_detail":                          clb.DataSourceTencentCloudClbInstanceDetail(),
			"tencentcloud_clb_resources":                                clb.DataSourceTencentCloudClbResources(),
			"tencentcloud_clb_target_group_list":                        clb.DataSourceTencentCloudClbTargetGroupList(),
			"tencentcloud_clb_target_health":                            clb.DataSourceTencentCloudClbTargetHealth(),
			"tencentcloud_elasticsearch_instances":                      es.DataSourceTencentCloudElasticsearchInstances(),
			"tencentcloud_elasticsearch_instance_logs":                  es.DataSourceTencentCloudElasticsearchInstanceLogs(),
			"tencentcloud_elasticsearch_instance_operations":            es.DataSourceTencentCloudElasticsearchInstanceOperations(),
			"tencentcloud_elasticsearch_logstash_instance_logs":         es.DataSourceTencentCloudElasticsearchLogstashInstanceLogs(),
			"tencentcloud_elasticsearch_logstash_instance_operations":   es.DataSourceTencentCloudElasticsearchLogstashInstanceOperations(),
			"tencentcloud_elasticsearch_views":                          es.DataSourceTencentCloudElasticsearchViews(),
			"tencentcloud_elasticsearch_diagnose":                       es.DataSourceTencentCloudElasticsearchDiagnose(),
			"tencentcloud_elasticsearch_instance_plugin_list":           es.DataSourceTencentCloudElasticsearchInstancePluginList(),
			"tencentcloud_elasticsearch_describe_index_list":            es.DataSourceTencentCloudElasticsearchDescribeIndexList(),
			"tencentcloud_mongodb_zone_config":                          mongodb.DataSourceTencentCloudMongodbZoneConfig(),
			"tencentcloud_mongodb_instances":                            mongodb.DataSourceTencentCloudMongodbInstances(),
			"tencentcloud_mongodb_instance_backups":                     mongodb.DataSourceTencentCloudMongodbInstanceBackups(),
			"tencentcloud_mongodb_instance_connections":                 mongodb.DataSourceTencentCloudMongodbInstanceConnections(),
			"tencentcloud_mongodb_instance_current_op":                  mongodb.DataSourceTencentCloudMongodbInstanceCurrentOp(),
			"tencentcloud_mongodb_instance_params":                      mongodb.DataSourceTencentCloudMongodbInstanceParams(),
			"tencentcloud_mongodb_instance_slow_log":                    mongodb.DataSourceTencentCloudMongodbInstanceSlowLog(),
			"tencentcloud_dayu_cc_https_policies":                       dayu.DataSourceTencentCloudDayuCCHttpsPolicies(),
			"tencentcloud_dayu_cc_http_policies":                        dayu.DataSourceTencentCloudDayuCCHttpPolicies(),
			"tencentcloud_dayu_ddos_policies":                           dayu.DataSourceTencentCloudDayuDdosPolicies(),
			"tencentcloud_dayu_ddos_policy_cases":                       dayu.DataSourceTencentCloudDayuDdosPolicyCases(),
			"tencentcloud_dayu_ddos_policy_attachments":                 dayu.DataSourceTencentCloudDayuDdosPolicyAttachments(),
			"tencentcloud_dayu_l4_rules":                                dayu.DataSourceTencentCloudDayuL4Rules(),
			"tencentcloud_dayu_l4_rules_v2":                             dayuv2.DataSourceTencentCloudDayuL4RulesV2(),
			"tencentcloud_dayu_l7_rules":                                dayu.DataSourceTencentCloudDayuL7Rules(),
			"tencentcloud_dayu_l7_rules_v2":                             dayuv2.DataSourceTencentCloudDayuL7RulesV2(),
			"tencentcloud_antiddos_pending_risk_info":                   dayuv2.DataSourceTencentCloudAntiddosPendingRiskInfo(),
			"tencentcloud_antiddos_overview_index":                      dayuv2.DataSourceTencentCloudAntiddosOverviewIndex(),
			"tencentcloud_antiddos_overview_ddos_trend":                 dayuv2.DataSourceTencentCloudAntiddosOverviewDdosTrend(),
			"tencentcloud_antiddos_overview_ddos_event_list":            dayuv2.DataSourceTencentCloudAntiddosOverviewDdosEventList(),
			"tencentcloud_antiddos_overview_cc_trend":                   dayuv2.DataSourceTencentCloudAntiddosOverviewCcTrend(),
			"tencentcloud_gaap_proxies":                                 gaap.DataSourceTencentCloudGaapProxies(),
			"tencentcloud_gaap_realservers":                             gaap.DataSourceTencentCloudGaapRealservers(),
			"tencentcloud_gaap_layer4_listeners":                        gaap.DataSourceTencentCloudGaapLayer4Listeners(),
			"tencentcloud_gaap_layer7_listeners":                        gaap.DataSourceTencentCloudGaapLayer7Listeners(),
			"tencentcloud_gaap_http_domains":                            gaap.DataSourceTencentCloudGaapHttpDomains(),
			"tencentcloud_gaap_http_rules":                              gaap.DataSourceTencentCloudGaapHttpRules(),
			"tencentcloud_gaap_security_policies":                       gaap.DataSourceTencentCloudGaapSecurityPolices(),
			"tencentcloud_gaap_security_rules":                          gaap.DataSourceTencentCloudGaapSecurityRules(),
			"tencentcloud_gaap_certificates":                            gaap.DataSourceTencentCloudGaapCertificates(),
			"tencentcloud_gaap_domain_error_pages":                      gaap.DataSourceTencentCloudGaapDomainErrorPageInfoList(),
			"tencentcloud_gaap_access_regions":                          gaap.DataSourceTencentCloudGaapAccessRegions(),
			"tencentcloud_gaap_access_regions_by_dest_region":           gaap.DataSourceTencentCloudGaapAccessRegionsByDestRegion(),
			"tencentcloud_gaap_black_header":                            gaap.DataSourceTencentCloudGaapBlackHeader(),
			"tencentcloud_gaap_country_area_mapping":                    gaap.DataSourceTencentCloudGaapCountryAreaMapping(),
			"tencentcloud_gaap_custom_header":                           gaap.DataSourceTencentCloudGaapCustomHeader(),
			"tencentcloud_gaap_dest_regions":                            gaap.DataSourceTencentCloudGaapDestRegions(),
			"tencentcloud_gaap_proxy_detail":                            gaap.DataSourceTencentCloudGaapProxyDetail(),
			"tencentcloud_gaap_proxy_groups":                            gaap.DataSourceTencentCloudGaapProxyGroups(),
			"tencentcloud_gaap_proxy_group_statistics":                  gaap.DataSourceTencentCloudGaapProxyGroupStatistics(),
			"tencentcloud_gaap_proxy_statistics":                        gaap.DataSourceTencentCloudGaapProxyStatistics(),
			"tencentcloud_gaap_real_servers_status":                     gaap.DataSourceTencentCloudGaapRealServersStatus(),
			"tencentcloud_gaap_rule_real_servers":                       gaap.DataSourceTencentCloudGaapRuleRealServers(),
			"tencentcloud_gaap_resources_by_tag":                        gaap.DataSourceTencentCloudGaapResourcesByTag(),
			"tencentcloud_gaap_region_and_price":                        gaap.DataSourceTencentCloudGaapRegionAndPrice(),
			"tencentcloud_gaap_proxy_and_statistics_listeners":          gaap.DataSourceTencentCloudGaapProxyAndStatisticsListeners(),
			"tencentcloud_gaap_proxies_status":                          gaap.DataSourceTencentCloudGaapProxiesStatus(),
			"tencentcloud_gaap_listener_statistics":                     gaap.DataSourceTencentCloudGaapListenerStatistics(),
			"tencentcloud_gaap_listener_real_servers":                   gaap.DataSourceTencentCloudGaapListenerRealServers(),
			"tencentcloud_gaap_group_and_statistics_proxy":              gaap.DataSourceTencentCloudGaapGroupAndStatisticsProxy(),
			"tencentcloud_gaap_domain_error_page_infos":                 gaap.DataSourceTencentCloudGaapDomainErrorPageInfos(),
			"tencentcloud_gaap_check_proxy_create":                      gaap.DataSourceTencentCloudGaapCheckProxyCreate(),
			"tencentcloud_ssl_certificates":                             ssl.DataSourceTencentCloudSslCertificates(),
			"tencentcloud_ssl_describe_certificate":                     ssl.DataSourceTencentCloudSslDescribeCertificate(),
			"tencentcloud_ssl_describe_companies":                       ssl.DataSourceTencentCloudSslDescribeCompanies(),
			"tencentcloud_ssl_describe_host_api_gateway_instance_list":  ssl.DataSourceTencentCloudSslDescribeHostApiGatewayInstanceList(),
			"tencentcloud_ssl_describe_host_cdn_instance_list":          ssl.DataSourceTencentCloudSslDescribeHostCdnInstanceList(),
			"tencentcloud_ssl_describe_host_clb_instance_list":          ssl.DataSourceTencentCloudSslDescribeHostClbInstanceList(),
			"tencentcloud_ssl_describe_host_cos_instance_list":          ssl.DataSourceTencentCloudSslDescribeHostCosInstanceList(),
			"tencentcloud_ssl_describe_host_ddos_instance_list":         ssl.DataSourceTencentCloudSslDescribeHostDdosInstanceList(),
			"tencentcloud_ssl_describe_host_lighthouse_instance_list":   ssl.DataSourceTencentCloudSslDescribeHostLighthouseInstanceList(),
			"tencentcloud_ssl_describe_host_live_instance_list":         ssl.DataSourceTencentCloudSslDescribeHostLiveInstanceList(),
			"tencentcloud_ssl_describe_host_teo_instance_list":          ssl.DataSourceTencentCloudSslDescribeHostTeoInstanceList(),
			"tencentcloud_ssl_describe_host_tke_instance_list":          ssl.DataSourceTencentCloudSslDescribeHostTkeInstanceList(),
			"tencentcloud_ssl_describe_host_update_record":              ssl.DataSourceTencentCloudSslDescribeHostUpdateRecord(),
			"tencentcloud_ssl_describe_host_update_record_detail":       ssl.DataSourceTencentCloudSslDescribeHostUpdateRecordDetail(),
			"tencentcloud_ssl_describe_host_vod_instance_list":          ssl.DataSourceTencentCloudSslDescribeHostVodInstanceList(),
			"tencentcloud_ssl_describe_host_waf_instance_list":          ssl.DataSourceTencentCloudSslDescribeHostWafInstanceList(),
			"tencentcloud_ssl_describe_manager_detail":                  ssl.DataSourceTencentCloudSslDescribeManagerDetail(),
			"tencentcloud_ssl_describe_managers":                        ssl.DataSourceTencentCloudSslDescribeManagers(),
			"tencentcloud_ssl_describe_host_deploy_record":              ssl.DataSourceTencentCloudSslDescribeHostDeployRecord(),
			"tencentcloud_ssl_describe_host_deploy_record_detail":       ssl.DataSourceTencentCloudSslDescribeHostDeployRecordDetail(),
			"tencentcloud_cam_roles":                                    cam.DataSourceTencentCloudCamRoles(),
			"tencentcloud_cam_users":                                    cam.DataSourceTencentCloudCamUsers(),
			"tencentcloud_cam_groups":                                   cam.DataSourceTencentCloudCamGroups(),
			"tencentcloud_cam_group_memberships":                        cam.DataSourceTencentCloudCamGroupMemberships(),
			"tencentcloud_cam_policies":                                 cam.DataSourceTencentCloudCamPolicies(),
			"tencentcloud_cam_role_policy_attachments":                  cam.DataSourceTencentCloudCamRolePolicyAttachments(),
			"tencentcloud_cam_user_policy_attachments":                  cam.DataSourceTencentCloudCamUserPolicyAttachments(),
			"tencentcloud_cam_group_policy_attachments":                 cam.DataSourceTencentCloudCamGroupPolicyAttachments(),
			"tencentcloud_cam_saml_providers":                           cam.DataSourceTencentCloudCamSAMLProviders(),
			"tencentcloud_cam_list_entities_for_policy":                 cam.DataSourceTencentCloudCamListEntitiesForPolicy(),
			"tencentcloud_cam_account_summary":                          cam.DataSourceTencentCloudCamAccountSummary(),
			"tencentcloud_cam_oidc_config":                              cam.DataSourceTencentCloudCamOidcConfig(),
			"tencentcloud_user_info":                                    cam.DataSourceTencentCloudUserInfo(),
			"tencentcloud_cam_sub_accounts":                             cam.DataSourceTencentCloudCamSubAccounts(),
			"tencentcloud_cam_role_detail":                              cam.DataSourceTencentCloudCamRoleDetail(),
			"tencentcloud_cdn_domains":                                  cdn.DataSourceTencentCloudCdnDomains(),
			"tencentcloud_cdn_domain_verifier":                          cdn.DataSourceTencentCloudCdnDomainVerifyRecord(),
			"tencentcloud_scf_functions":                                scf.DataSourceTencentCloudScfFunctions(),
			"tencentcloud_scf_namespaces":                               scf.DataSourceTencentCloudScfNamespaces(),
			"tencentcloud_scf_account_info":                             scf.DataSourceTencentCloudScfAccountInfo(),
			"tencentcloud_scf_async_event_management":                   scf.DataSourceTencentCloudScfAsyncEventManagement(),
			"tencentcloud_scf_triggers":                                 scf.DataSourceTencentCloudScfTriggers(),
			"tencentcloud_scf_async_event_status":                       scf.DataSourceTencentCloudScfAsyncEventStatus(),
			"tencentcloud_scf_function_address":                         scf.DataSourceTencentCloudScfFunctionAddress(),
			"tencentcloud_scf_request_status":                           scf.DataSourceTencentCloudScfRequestStatus(),
			"tencentcloud_scf_function_aliases":                         scf.DataSourceTencentCloudScfFunctionAliases(),
			"tencentcloud_scf_layer_versions":                           scf.DataSourceTencentCloudScfLayerVersions(),
			"tencentcloud_scf_layers":                                   scf.DataSourceTencentCloudScfLayers(),
			"tencentcloud_scf_function_versions":                        scf.DataSourceTencentCloudScfFunctionVersions(),
			"tencentcloud_scf_logs":                                     scf.DataSourceTencentCloudScfLogs(),
			"tencentcloud_tcaplus_clusters":                             tcaplusdb.DataSourceTencentCloudTcaplusClusters(),
			"tencentcloud_tcaplus_tablegroups":                          tcaplusdb.DataSourceTencentCloudTcaplusTableGroups(),
			"tencentcloud_tcaplus_tables":                               tcaplusdb.DataSourceTencentCloudTcaplusTables(),
			"tencentcloud_tcaplus_idls":                                 tcaplusdb.DataSourceTencentCloudTcaplusIdls(),
			"tencentcloud_monitor_policy_conditions":                    monitor.DataSourceTencentCloudMonitorPolicyConditions(),
			"tencentcloud_monitor_data":                                 monitor.DataSourceTencentCloudMonitorData(),
			"tencentcloud_monitor_product_event":                        monitor.DataSourceTencentCloudMonitorProductEvent(),
			"tencentcloud_monitor_binding_objects":                      monitor.DataSourceTencentCloudMonitorBindingObjects(),
			"tencentcloud_monitor_policy_groups":                        monitor.DataSourceTencentCloudMonitorPolicyGroups(),
			"tencentcloud_monitor_product_namespace":                    monitor.DataSourceTencentCloudMonitorProductNamespace(),
			"tencentcloud_monitor_alarm_notices":                        monitor.DataSourceTencentCloudMonitorAlarmNotices(),
			"tencentcloud_monitor_alarm_metric":                         monitor.DataSourceTencentCloudMonitorAlarmMetric(),
			"tencentcloud_monitor_alarm_policy":                         monitor.DataSourceTencentCloudMonitorAlarmPolicy(),
			"tencentcloud_monitor_alarm_history":                        monitor.DataSourceTencentCloudMonitorAlarmHistory(),
			"tencentcloud_monitor_alarm_basic_alarms":                   monitor.DataSourceTencentCloudMonitorAlarmBasicAlarms(),
			"tencentcloud_monitor_alarm_basic_metric":                   monitor.DataSourceTencentCloudMonitorAlarmBasicMetric(),
			"tencentcloud_monitor_alarm_conditions_template":            monitor.DataSourceTencentCloudMonitorAlarmConditionsTemplate(),
			"tencentcloud_monitor_grafana_plugin_overviews":             tcmg.DataSourceTencentCloudMonitorGrafanaPluginOverviews(),
			"tencentcloud_monitor_alarm_notice_callbacks":               monitor.DataSourceTencentCloudMonitorAlarmNoticeCallbacks(),
			"tencentcloud_monitor_alarm_all_namespaces":                 monitor.DataSourceTencentCloudMonitorAlarmAllNamespaces(),
			"tencentcloud_monitor_alarm_monitor_type":                   monitor.DataSourceTencentCloudMonitorAlarmMonitorType(),
			"tencentcloud_monitor_statistic_data":                       monitor.DataSourceTencentCloudMonitorStatisticData(),
			"tencentcloud_monitor_tmp_regions":                          tmp.DataSourceTencentCloudMonitorTmpRegions(),
			"tencentcloud_monitor_tmp_instances":                        tmp.DataSourceTencentCloudMonitorTmpInstances(),
			"tencentcloud_postgresql_instances":                         postgresql.DataSourceTencentCloudPostgresqlInstances(),
			"tencentcloud_postgresql_specinfos":                         postgresql.DataSourceTencentCloudPostgresqlSpecinfos(),
			"tencentcloud_postgresql_xlogs":                             postgresql.DataSourceTencentCloudPostgresqlXlogs(),
			"tencentcloud_postgresql_parameter_templates":               postgresql.DataSourceTencentCloudPostgresqlParameterTemplates(),
			"tencentcloud_postgresql_readonly_groups":                   postgresql.DataSourceTencentCloudPostgresqlReadonlyGroups(),
			"tencentcloud_postgresql_base_backups":                      postgresql.DataSourceTencentCloudPostgresqlBaseBackups(),
			"tencentcloud_postgresql_log_backups":                       postgresql.DataSourceTencentCloudPostgresqlLogBackups(),
			"tencentcloud_postgresql_backup_download_urls":              postgresql.DataSourceTencentCloudPostgresqlBackupDownloadUrls(),
			"tencentcloud_postgresql_db_instance_classes":               postgresql.DataSourceTencentCloudPostgresqlDbInstanceClasses(),
			"tencentcloud_postgresql_default_parameters":                postgresql.DataSourceTencentCloudPostgresqlDefaultParameters(),
			"tencentcloud_postgresql_recovery_time":                     postgresql.DataSourceTencentCloudPostgresqlRecoveryTime(),
			"tencentcloud_postgresql_regions":                           postgresql.DataSourceTencentCloudPostgresqlRegions(),
			"tencentcloud_postgresql_db_instance_versions":              postgresql.DataSourceTencentCloudPostgresqlDbInstanceVersions(),
			"tencentcloud_postgresql_zones":                             postgresql.DataSourceTencentCloudPostgresqlZones(),
			"tencentcloud_postgresql_account_privileges":                postgresql.DataSourceTencentCloudPostgresqlAccountPrivileges(),
			"tencentcloud_postgresql_dedicated_clusters":                postgresql.DataSourceTencentCloudPostgresqlDedicatedClusters(),
			"tencentcloud_sqlserver_zone_config":                        sqlserver.DataSourceTencentCloudSqlserverZoneConfig(),
			"tencentcloud_sqlserver_instances":                          sqlserver.DataSourceTencentCloudSqlserverInstances(),
			"tencentcloud_sqlserver_backups":                            sqlserver.DataSourceTencentCloudSqlserverBackups(),
			"tencentcloud_sqlserver_dbs":                                sqlserver.DataSourceTencentCloudSqlserverDBs(),
			"tencentcloud_sqlserver_accounts":                           sqlserver.DataSourceTencentCloudSqlserverAccounts(),
			"tencentcloud_sqlserver_account_db_attachments":             sqlserver.DataSourceTencentCloudSqlserverAccountDBAttachments(),
			"tencentcloud_sqlserver_readonly_groups":                    sqlserver.DataSourceTencentCloudSqlserverReadonlyGroups(),
			"tencentcloud_sqlserver_backup_commands":                    sqlserver.DataSourceTencentCloudSqlserverBackupCommands(),
			"tencentcloud_sqlserver_backup_by_flow_id":                  sqlserver.DataSourceTencentCloudSqlserverBackupByFlowId(),
			"tencentcloud_sqlserver_backup_upload_size":                 sqlserver.DataSourceTencentCloudSqlserverBackupUploadSize(),
			"tencentcloud_sqlserver_cross_region_zone":                  sqlserver.DataSourceTencentCloudSqlserverCrossRegionZone(),
			"tencentcloud_sqlserver_db_charsets":                        sqlserver.DataSourceTencentCloudSqlserverDBCharsets(),
			"tencentcloud_ckafka_users":                                 ckafka.DataSourceTencentCloudCkafkaUsers(),
			"tencentcloud_ckafka_acls":                                  ckafka.DataSourceTencentCloudCkafkaAcls(),
			"tencentcloud_ckafka_topics":                                ckafka.DataSourceTencentCloudCkafkaTopics(),
			"tencentcloud_ckafka_instances":                             ckafka.DataSourceTencentCloudCkafkaInstances(),
			"tencentcloud_ckafka_connect_resource":                      ckafka.DataSourceTencentCloudCkafkaConnectResource(),
			"tencentcloud_ckafka_region":                                ckafka.DataSourceTencentCloudCkafkaRegion(),
			"tencentcloud_ckafka_datahub_topic":                         ckafka.DataSourceTencentCloudCkafkaDatahubTopic(),
			"tencentcloud_ckafka_datahub_group_offsets":                 ckafka.DataSourceTencentCloudCkafkaDatahubGroupOffsets(),
			"tencentcloud_ckafka_datahub_task":                          ckafka.DataSourceTencentCloudCkafkaDatahubTask(),
			"tencentcloud_ckafka_group":                                 ckafka.DataSourceTencentCloudCkafkaGroup(),
			"tencentcloud_ckafka_group_offsets":                         ckafka.DataSourceTencentCloudCkafkaGroupOffsets(),
			"tencentcloud_ckafka_group_info":                            ckafka.DataSourceTencentCloudCkafkaGroupInfo(),
			"tencentcloud_ckafka_task_status":                           ckafka.DataSourceTencentCloudCkafkaTaskStatus(),
			"tencentcloud_ckafka_topic_flow_ranking":                    ckafka.DataSourceTencentCloudCkafkaTopicFlowRanking(),
			"tencentcloud_ckafka_topic_produce_connection":              ckafka.DataSourceTencentCloudCkafkaTopicProduceConnection(),
			"tencentcloud_ckafka_topic_subscribe_group":                 ckafka.DataSourceTencentCloudCkafkaTopicSubscribeGroup(),
			"tencentcloud_ckafka_topic_sync_replica":                    ckafka.DataSourceTencentCloudCkafkaTopicSyncReplica(),
			"tencentcloud_ckafka_zone":                                  ckafka.DataSourceTencentCloudCkafkaZone(),
			"tencentcloud_audit_cos_regions":                            audit.DataSourceTencentCloudAuditCosRegions(),
			"tencentcloud_audit_key_alias":                              audit.DataSourceTencentCloudAuditKeyAlias(),
			"tencentcloud_audits":                                       audit.DataSourceTencentCloudAudits(),
			"tencentcloud_audit_events":                                 audit.DataSourceTencentCloudAuditEvents(),
			"tencentcloud_cynosdb_clusters":                             cynosdb.DataSourceTencentCloudCynosdbClusters(),
			"tencentcloud_cynosdb_instances":                            cynosdb.DataSourceTencentCloudCynosdbInstances(),
			"tencentcloud_cynosdb_zone_config":                          cynosdb.DataSourceTencentCloudCynosdbZoneConfig(),
			"tencentcloud_cynosdb_instance_slow_queries":                cynosdb.DataSourceTencentCloudCynosdbInstanceSlowQueries(),
			"tencentcloud_vod_adaptive_dynamic_streaming_templates":     vod.DataSourceTencentCloudVodAdaptiveDynamicStreamingTemplates(),
			"tencentcloud_vod_image_sprite_templates":                   vod.DataSourceTencentCloudVodImageSpriteTemplates(),
			"tencentcloud_vod_procedure_templates":                      vod.DataSourceTencentCloudVodProcedureTemplates(),
			"tencentcloud_vod_snapshot_by_time_offset_templates":        vod.DataSourceTencentCloudVodSnapshotByTimeOffsetTemplates(),
			"tencentcloud_vod_super_player_configs":                     vod.DataSourceTencentCloudVodSuperPlayerConfigs(),
			"tencentcloud_sqlserver_publish_subscribes":                 sqlserver.DataSourceTencentCloudSqlserverPublishSubscribes(),
			"tencentcloud_sqlserver_instance_param_records":             sqlserver.DataSourceTencentCloudSqlserverInstanceParamRecords(),
			"tencentcloud_sqlserver_project_security_groups":            sqlserver.DataSourceTencentCloudSqlserverProjectSecurityGroups(),
			"tencentcloud_sqlserver_regions":                            sqlserver.DataSourceTencentCloudSqlserverRegions(),
			"tencentcloud_sqlserver_rollback_time":                      sqlserver.DataSourceTencentCloudSqlserverRollbackTime(),
			"tencentcloud_sqlserver_slowlogs":                           sqlserver.DataSourceTencentCloudSqlserverSlowlogs(),
			"tencentcloud_sqlserver_upload_backup_info":                 sqlserver.DataSourceTencentCloudSqlserverUploadBackupInfo(),
			"tencentcloud_sqlserver_upload_incremental_info":            sqlserver.DataSourceTencentCloudSqlserverUploadIncrementalInfo(),
			"tencentcloud_api_gateway_usage_plans":                      apigateway.DataSourceTencentCloudAPIGatewayUsagePlans(),
			"tencentcloud_api_gateway_ip_strategies":                    apigateway.DataSourceTencentCloudAPIGatewayIpStrategy(),
			"tencentcloud_api_gateway_customer_domains":                 apigateway.DataSourceTencentCloudAPIGatewayCustomerDomains(),
			"tencentcloud_api_gateway_usage_plan_environments":          apigateway.DataSourceTencentCloudAPIGatewayUsagePlanEnvironments(),
			"tencentcloud_api_gateway_throttling_services":              apigateway.DataSourceTencentCloudAPIGatewayThrottlingServices(),
			"tencentcloud_api_gateway_throttling_apis":                  apigateway.DataSourceTencentCloudAPIGatewayThrottlingApis(),
			"tencentcloud_api_gateway_apis":                             apigateway.DataSourceTencentCloudAPIGatewayAPIs(),
			"tencentcloud_api_gateway_services":                         apigateway.DataSourceTencentCloudAPIGatewayServices(),
			"tencentcloud_api_gateway_api_keys":                         apigateway.DataSourceTencentCloudAPIGatewayAPIKeys(),
			"tencentcloud_api_gateway_plugins":                          apigateway.DataSourceTencentCloudAPIGatewayPlugins(),
			"tencentcloud_api_gateway_upstreams":                        apigateway.DataSourceTencentCloudAPIGatewayUpstreams(),
			"tencentcloud_api_gateway_api_usage_plans":                  apigateway.DataSourceTencentCloudAPIGatewayApiUsagePlans(),
			"tencentcloud_api_gateway_api_app_service":                  apigateway.DataSourceTencentCloudAPIGatewayApiAppService(),
			"tencentcloud_api_gateway_bind_api_apps_status":             apigateway.DataSourceTencentCloudApiGatewayBindApiAppsStatus(),
			"tencentcloud_api_gateway_api_app_api":                      apigateway.DataSourceTencentCloudApiGatewayApiAppApi(),
			"tencentcloud_api_gateway_api_plugins":                      apigateway.DataSourceTencentCloudApiGatewayApiPlugins(),
			"tencentcloud_api_gateway_service_release_versions":         apigateway.DataSourceTencentCloudApiGatewayServiceReleaseVersions(),
			"tencentcloud_api_gateway_service_environment_list":         apigateway.DataSourceTencentCloudApiGatewayServiceEnvironmentList(),
			"tencentcloud_sqlserver_basic_instances":                    sqlserver.DataSourceTencentCloudSqlserverBasicInstances(),
			"tencentcloud_sqlserver_query_xevent":                       sqlserver.DataSourceTencentCloudSqlserverQueryXevent(),
			"tencentcloud_sqlserver_ins_attribute":                      sqlserver.DataSourceTencentCloudSqlserverInsAttribute(),
			"tencentcloud_sqlserver_desc_ha_log":                        sqlserver.DataSourceTencentCloudSqlserverDescHaLog(),
			"tencentcloud_tcr_instances":                                tcr.DataSourceTencentCloudTCRInstances(),
			"tencentcloud_tcr_namespaces":                               tcr.DataSourceTencentCloudTCRNamespaces(),
			"tencentcloud_tcr_tokens":                                   tcr.DataSourceTencentCloudTCRTokens(),
			"tencentcloud_tcr_vpc_attachments":                          tcr.DataSourceTencentCloudTCRVPCAttachments(),
			"tencentcloud_tcr_repositories":                             tcr.DataSourceTencentCloudTCRRepositories(),
			"tencentcloud_tcr_webhook_trigger_logs":                     tcr.DataSourceTencentCloudTcrWebhookTriggerLogs(),
			"tencentcloud_tcr_images":                                   tcr.DataSourceTencentCloudTcrImages(),
			"tencentcloud_tcr_image_manifests":                          tcr.DataSourceTencentCloudTcrImageManifests(),
			"tencentcloud_tcr_tag_retention_execution_tasks":            tcr.DataSourceTencentCloudTcrTagRetentionExecutionTasks(),
			"tencentcloud_tcr_tag_retention_executions":                 tcr.DataSourceTencentCloudTcrTagRetentionExecutions(),
			"tencentcloud_tcr_replication_instance_create_tasks":        tcr.DataSourceTencentCloudTcrReplicationInstanceCreateTasks(),
			"tencentcloud_tcr_replication_instance_sync_status":         tcr.DataSourceTencentCloudTcrReplicationInstanceSyncStatus(),
			"tencentcloud_address_templates":                            vpc.DataSourceTencentCloudAddressTemplates(),
			"tencentcloud_address_template_groups":                      vpc.DataSourceTencentCloudAddressTemplateGroups(),
			"tencentcloud_protocol_templates":                           vpc.DataSourceTencentCloudProtocolTemplates(),
			"tencentcloud_protocol_template_groups":                     vpc.DataSourceTencentCloudProtocolTemplateGroups(),
			"tencentcloud_classic_elastic_public_ipv6s":                 vpc.DataSourceTencentCloudClassicElasticPublicIpv6s(),
			"tencentcloud_elastic_public_ipv6s":                         vpc.DataSourceTencentCloudElasticPublicIpv6s(),
			"tencentcloud_kms_keys":                                     kms.DataSourceTencentCloudKmsKeys(),
			"tencentcloud_kms_public_key":                               kms.DataSourceTencentCloudKmsPublicKey(),
			"tencentcloud_kms_get_parameters_for_import":                kms.DataSourceTencentCloudKmsGetParametersForImport(),
			"tencentcloud_kms_describe_keys":                            kms.DataSourceTencentCloudKmsDescribeKeys(),
			"tencentcloud_kms_white_box_key_details":                    kms.DataSourceTencentCloudKmsWhiteBoxKeyDetails(),
			"tencentcloud_kms_list_keys":                                kms.DataSourceTencentCloudKmsListKeys(),
			"tencentcloud_kms_white_box_decrypt_key":                    kms.DataSourceTencentCloudKmsWhiteBoxDecryptKey(),
			"tencentcloud_kms_white_box_device_fingerprints":            kms.DataSourceTencentCloudKmsWhiteBoxDeviceFingerprints(),
			"tencentcloud_kms_list_algorithms":                          kms.DataSourceTencentCloudKmsListAlgorithms(),
			"tencentcloud_ssm_products":                                 ssm.DataSourceTencentCloudSsmProducts(),
			"tencentcloud_ssm_secrets":                                  ssm.DataSourceTencentCloudSsmSecrets(),
			"tencentcloud_ssm_secret_versions":                          ssm.DataSourceTencentCloudSsmSecretVersions(),
			"tencentcloud_ssm_rotation_detail":                          ssm.DataSourceTencentCloudSsmRotationDetail(),
			"tencentcloud_ssm_rotation_history":                         ssm.DataSourceTencentCloudSsmRotationHistory(),
			"tencentcloud_ssm_service_status":                           ssm.DataSourceTencentCloudSsmServiceStatus(),
			"tencentcloud_ssm_ssh_key_pair_value":                       ssm.DataSourceTencentCloudSsmSshKeyPairValue(),
			"tencentcloud_cdh_instances":                                cdh.DataSourceTencentCloudCdhInstances(),
			"tencentcloud_dayu_eip":                                     dayuv2.DataSourceTencentCloudDayuEip(),
			"tencentcloud_teo_zone_available_plans":                     teo.DataSourceTencentCloudTeoZoneAvailablePlans(),
			"tencentcloud_teo_rule_engine_settings":                     teo.DataSourceTencentCloudTeoRuleEngineSettings(),
			"tencentcloud_sts_caller_identity":                          sts.DataSourceTencentCloudStsCallerIdentity(),
			"tencentcloud_dcdb_instances":                               dcdb.DataSourceTencentCloudDcdbInstances(),
			"tencentcloud_dcdb_accounts":                                dcdb.DataSourceTencentCloudDcdbAccounts(),
			"tencentcloud_dcdb_databases":                               dcdb.DataSourceTencentCloudDcdbDatabases(),
			"tencentcloud_dcdb_parameters":                              dcdb.DataSourceTencentCloudDcdbParameters(),
			"tencentcloud_dcdb_shards":                                  dcdb.DataSourceTencentCloudDcdbShards(),
			"tencentcloud_dcdb_security_groups":                         dcdb.DataSourceTencentCloudDcdbSecurityGroups(),
			"tencentcloud_dcdb_database_objects":                        dcdb.DataSourceTencentCloudDcdbDatabaseObjects(),
			"tencentcloud_dcdb_database_tables":                         dcdb.DataSourceTencentCloudDcdbDatabaseTables(),
			"tencentcloud_dcdb_file_download_url":                       dcdb.DataSourceTencentCloudDcdbFileDownloadUrl(),
			"tencentcloud_dcdb_log_files":                               dcdb.DataSourceTencentCloudDcdbLogFiles(),
			"tencentcloud_dcdb_instance_node_info":                      dcdb.DataSourceTencentCloudDcdbInstanceNodeInfo(),
			"tencentcloud_dcdb_orders":                                  dcdb.DataSourceTencentCloudDcdbOrders(),
			"tencentcloud_dcdb_price":                                   dcdb.DataSourceTencentCloudDcdbPrice(),
			"tencentcloud_dcdb_project_security_groups":                 dcdb.DataSourceTencentCloudDcdbProjectSecurityGroups(),
			"tencentcloud_dcdb_projects":                                dcdb.DataSourceTencentCloudDcdbProjects(),
			"tencentcloud_dcdb_renewal_price":                           dcdb.DataSourceTencentCloudDcdbRenewalPrice(),
			"tencentcloud_dcdb_sale_info":                               dcdb.DataSourceTencentCloudDcdbSaleInfo(),
			"tencentcloud_dcdb_shard_spec":                              dcdb.DataSourceTencentCloudDcdbShardSpec(),
			"tencentcloud_dcdb_slow_logs":                               dcdb.DataSourceTencentCloudDcdbSlowLogs(),
			"tencentcloud_dcdb_upgrade_price":                           dcdb.DataSourceTencentCloudDcdbUpgradePrice(),
			"tencentcloud_mariadb_db_instances":                         mariadb.DataSourceTencentCloudMariadbDbInstances(),
			"tencentcloud_mariadb_accounts":                             mariadb.DataSourceTencentCloudMariadbAccounts(),
			"tencentcloud_mariadb_security_groups":                      mariadb.DataSourceTencentCloudMariadbSecurityGroups(),
			"tencentcloud_mariadb_database_objects":                     mariadb.DataSourceTencentCloudMariadbDatabaseObjects(),
			"tencentcloud_mariadb_databases":                            mariadb.DataSourceTencentCloudMariadbDatabases(),
			"tencentcloud_mariadb_database_table":                       mariadb.DataSourceTencentCloudMariadbDatabaseTable(),
			"tencentcloud_mariadb_dcn_detail":                           mariadb.DataSourceTencentCloudMariadbDcnDetail(),
			"tencentcloud_mariadb_file_download_url":                    mariadb.DataSourceTencentCloudMariadbFileDownloadUrl(),
			"tencentcloud_mariadb_flow":                                 mariadb.DataSourceTencentCloudMariadbFlow(),
			"tencentcloud_mariadb_instance_node_info":                   mariadb.DataSourceTencentCloudMariadbInstanceNodeInfo(),
			"tencentcloud_mariadb_instance_specs":                       mariadb.DataSourceTencentCloudMariadbInstanceSpecs(),
			"tencentcloud_mariadb_log_files":                            mariadb.DataSourceTencentCloudMariadbLogFiles(),
			"tencentcloud_mariadb_orders":                               mariadb.DataSourceTencentCloudMariadbOrders(),
			"tencentcloud_mariadb_price":                                mariadb.DataSourceTencentCloudMariadbPrice(),
			"tencentcloud_mariadb_project_security_groups":              mariadb.DataSourceTencentCloudMariadbProjectSecurityGroups(),
			"tencentcloud_mariadb_renewal_price":                        mariadb.DataSourceTencentCloudMariadbRenewalPrice(),
			"tencentcloud_mariadb_sale_info":                            mariadb.DataSourceTencentCloudMariadbSaleInfo(),
			"tencentcloud_mariadb_slow_logs":                            mariadb.DataSourceTencentCloudMariadbSlowLogs(),
			"tencentcloud_mariadb_upgrade_price":                        mariadb.DataSourceTencentCloudMariadbUpgradePrice(),
			"tencentcloud_mps_schedules":                                mps.DataSourceTencentCloudMpsSchedules(),
			"tencentcloud_mps_tasks":                                    mps.DataSourceTencentCloudMpsTasks(),
			"tencentcloud_mps_parse_live_stream_process_notification":   mps.DataSourceTencentCloudMpsParseLiveStreamProcessNotification(),
			"tencentcloud_mps_parse_notification":                       mps.DataSourceTencentCloudMpsParseNotification(),
			"tencentcloud_mps_media_meta_data":                          mps.DataSourceTencentCloudMpsMediaMetaData(),
			"tencentcloud_tdcpg_clusters":                               tdcpg.DataSourceTencentCloudTdcpgClusters(),
			"tencentcloud_tdcpg_instances":                              tdcpg.DataSourceTencentCloudTdcpgInstances(),
			"tencentcloud_cat_probe_data":                               cat.DataSourceTencentCloudCatProbeData(),
			"tencentcloud_cat_node":                                     cat.DataSourceTencentCloudCatNode(),
			"tencentcloud_cat_metric_data":                              cat.DataSourceTencentCloudCatMetricData(),
			"tencentcloud_rum_project":                                  rum.DataSourceTencentCloudRumProject(),
			"tencentcloud_rum_offline_log_config":                       rum.DataSourceTencentCloudRumOfflineLogConfig(),
			"tencentcloud_rum_whitelist":                                rum.DataSourceTencentCloudRumWhitelist(),
			"tencentcloud_rum_taw_instance":                             rum.DataSourceTencentCloudRumTawInstance(),
			"tencentcloud_rum_custom_url":                               rum.DataSourceTencentCloudRumCustomUrl(),
			"tencentcloud_rum_event_url":                                rum.DataSourceTencentCloudRumEventUrl(),
			"tencentcloud_rum_fetch_url_info":                           rum.DataSourceTencentCloudRumFetchUrlInfo(),
			"tencentcloud_rum_fetch_url":                                rum.DataSourceTencentCloudRumFetchUrl(),
			"tencentcloud_rum_group_log":                                rum.DataSourceTencentCloudRumGroupLog(),
			"tencentcloud_rum_log_list":                                 rum.DataSourceTencentCloudRumLogList(),
			"tencentcloud_rum_log_url_statistics":                       rum.DataSourceTencentCloudRumLogUrlStatistics(),
			"tencentcloud_rum_performance_page":                         rum.DataSourceTencentCloudRumPerformancePage(),
			"tencentcloud_rum_pv_url_info":                              rum.DataSourceTencentCloudRumPvUrlInfo(),
			"tencentcloud_rum_pv_url_statistics":                        rum.DataSourceTencentCloudRumPvUrlStatistics(),
			"tencentcloud_rum_report_count":                             rum.DataSourceTencentCloudRumReportCount(),
			"tencentcloud_rum_log_stats_log_list":                       rum.DataSourceTencentCloudRumLogStatsLogList(),
			"tencentcloud_rum_scores":                                   rum.DataSourceTencentCloudRumScores(),
			"tencentcloud_rum_set_url_statistics":                       rum.DataSourceTencentCloudRumSetUrlStatistics(),
			"tencentcloud_rum_sign":                                     rum.DataSourceTencentCloudRumSign(),
			"tencentcloud_rum_static_project":                           rum.DataSourceTencentCloudRumStaticProject(),
			"tencentcloud_rum_static_resource":                          rum.DataSourceTencentCloudRumStaticResource(),
			"tencentcloud_rum_static_url":                               rum.DataSourceTencentCloudRumStaticUrl(),
			"tencentcloud_rum_taw_area":                                 rum.DataSourceTencentCloudRumTawArea(),
			"tencentcloud_rum_web_vitals_page":                          rum.DataSourceTencentCloudRumWebVitalsPage(),
			"tencentcloud_rum_log_export":                               rum.DataSourceTencentCloudRumLogExport(),
			"tencentcloud_rum_log_export_list":                          rum.DataSourceTencentCloudRumLogExportList(),
			"tencentcloud_dnspod_records":                               dnspod.DataSourceTencentCloudDnspodRecords(),
			"tencentcloud_dnspod_domain_list":                           dnspod.DataSourceTencentCloudDnspodDomainList(),
			"tencentcloud_dnspod_domain_analytics":                      dnspod.DataSourceTencentCloudDnspodDomainAnalytics(),
			"tencentcloud_dnspod_domain_log_list":                       dnspod.DataSourceTencentCloudDnspodDomainLogList(),
			"tencentcloud_dnspod_record_analytics":                      dnspod.DataSourceTencentCloudDnspodRecordAnalytics(),
			"tencentcloud_dnspod_record_line_list":                      dnspod.DataSourceTencentCloudDnspodRecordLineList(),
			"tencentcloud_dnspod_record_list":                           dnspod.DataSourceTencentCloudDnspodRecordList(),
			"tencentcloud_dnspod_record_type":                           dnspod.DataSourceTencentCloudDnspodRecordType(),
			"tencentcloud_subdomain_validate_status":                    dnspod.DataSourceTencentCloudSubdomainValidateStatus(),
			"tencentcloud_tat_command":                                  tat.DataSourceTencentCloudTatCommand(),
			"tencentcloud_tat_invoker":                                  tat.DataSourceTencentCloudTatInvoker(),
			"tencentcloud_tat_invoker_records":                          tat.DataSourceTencentCloudTatInvokerRecords(),
			"tencentcloud_tat_agent":                                    tat.DataSourceTencentCloudTatAgent(),
			"tencentcloud_tat_invocation_task":                          tat.DataSourceTencentCloudTatInvocationTask(),
			"tencentcloud_dbbrain_sql_filters":                          dbbrain.DataSourceTencentCloudDbbrainSqlFilters(),
			"tencentcloud_dbbrain_security_audit_log_export_tasks":      dbbrain.DataSourceTencentCloudDbbrainSecurityAuditLogExportTasks(),
			"tencentcloud_dbbrain_diag_event":                           dbbrain.DataSourceTencentCloudDbbrainDiagEvent(),
			"tencentcloud_dbbrain_diag_events":                          dbbrain.DataSourceTencentCloudDbbrainDiagEvents(),
			"tencentcloud_dbbrain_diag_history":                         dbbrain.DataSourceTencentCloudDbbrainDiagHistory(),
			"tencentcloud_dbbrain_security_audit_log_download_urls":     dbbrain.DataSourceTencentCloudDbbrainSecurityAuditLogDownloadUrls(),
			"tencentcloud_dbbrain_slow_log_time_series_stats":           dbbrain.DataSourceTencentCloudDbbrainSlowLogTimeSeriesStats(),
			"tencentcloud_dbbrain_slow_log_top_sqls":                    dbbrain.DataSourceTencentCloudDbbrainSlowLogTopSqls(),
			"tencentcloud_dbbrain_slow_log_user_host_stats":             dbbrain.DataSourceTencentCloudDbbrainSlowLogUserHostStats(),
			"tencentcloud_dbbrain_slow_log_user_sql_advice":             dbbrain.DataSourceTencentCloudDbbrainSlowLogUserSqlAdvice(),
			"tencentcloud_dbbrain_slow_logs":                            dbbrain.DataSourceTencentCloudDbbrainSlowLogs(),
			"tencentcloud_dbbrain_health_scores":                        dbbrain.DataSourceTencentCloudDbbrainHealthScores(),
			"tencentcloud_dbbrain_sql_templates":                        dbbrain.DataSourceTencentCloudDbbrainSqlTemplates(),
			"tencentcloud_dbbrain_db_space_status":                      dbbrain.DataSourceTencentCloudDbbrainDbSpaceStatus(),
			"tencentcloud_dbbrain_top_space_schemas":                    dbbrain.DataSourceTencentCloudDbbrainTopSpaceSchemas(),
			"tencentcloud_dbbrain_top_space_tables":                     dbbrain.DataSourceTencentCloudDbbrainTopSpaceTables(),
			"tencentcloud_dbbrain_top_space_schema_time_series":         dbbrain.DataSourceTencentCloudDbbrainTopSpaceSchemaTimeSeries(),
			"tencentcloud_dbbrain_top_space_table_time_series":          dbbrain.DataSourceTencentCloudDbbrainTopSpaceTableTimeSeries(),
			"tencentcloud_dbbrain_diag_db_instances":                    dbbrain.DataSourceTencentCloudDbbrainDiagDbInstances(),
			"tencentcloud_dbbrain_mysql_process_list":                   dbbrain.DataSourceTencentCloudDbbrainMysqlProcessList(),
			"tencentcloud_dbbrain_no_primary_key_tables":                dbbrain.DataSourceTencentCloudDbbrainNoPrimaryKeyTables(),
			"tencentcloud_dbbrain_redis_top_big_keys":                   dbbrain.DataSourceTencentCloudDbbrainRedisTopBigKeys(),
			"tencentcloud_dbbrain_redis_top_key_prefix_list":            dbbrain.DataSourceTencentCloudDbbrainRedisTopKeyPrefixList(),
			"tencentcloud_dts_sync_jobs":                                dts.DataSourceTencentCloudDtsSyncJobs(),
			"tencentcloud_dts_compare_tasks":                            dts.DataSourceTencentCloudDtsCompareTasks(),
			"tencentcloud_dts_migrate_jobs":                             dts.DataSourceTencentCloudDtsMigrateJobs(),
			"tencentcloud_dts_migrate_db_instances":                     dts.DataSourceTencentCloudDtsMigrateDbInstances(),
			"tencentcloud_tdmq_rocketmq_cluster":                        trocket.DataSourceTencentCloudTdmqRocketmqCluster(),
			"tencentcloud_tdmq_rocketmq_namespace":                      trocket.DataSourceTencentCloudTdmqRocketmqNamespace(),
			"tencentcloud_tdmq_rocketmq_topic":                          trocket.DataSourceTencentCloudTdmqRocketmqTopic(),
			"tencentcloud_tdmq_rocketmq_role":                           trocket.DataSourceTencentCloudTdmqRocketmqRole(),
			"tencentcloud_tdmq_rocketmq_group":                          trocket.DataSourceTencentCloudTdmqRocketmqGroup(),
			"tencentcloud_tdmq_environment_attributes":                  tpulsar.DataSourceTencentCloudTdmqEnvironmentAttributes(),
			"tencentcloud_tdmq_publisher_summary":                       tpulsar.DataSourceTencentCloudTdmqPublisherSummary(),
			"tencentcloud_tdmq_publishers":                              tpulsar.DataSourceTencentCloudTdmqPublishers(),
			"tencentcloud_tdmq_rabbitmq_node_list":                      trabbit.DataSourceTencentCloudTdmqRabbitmqNodeList(),
			"tencentcloud_tdmq_rabbitmq_vip_instance":                   trabbit.DataSourceTencentCloudTdmqRabbitmqVipInstance(),
			"tencentcloud_tdmq_vip_instance":                            trocket.DataSourceTencentCloudTdmqVipInstance(),
			"tencentcloud_tdmq_rocketmq_messages":                       trocket.DataSourceTencentCloudTdmqRocketmqMessages(),
			"tencentcloud_tdmq_pro_instances":                           tpulsar.DataSourceTencentCloudTdmqProInstances(),
			"tencentcloud_tdmq_pro_instance_detail":                     tpulsar.DataSourceTencentCloudTdmqProInstanceDetail(),
			"tencentcloud_tcmq_queue":                                   tcmq.DataSourceTencentCloudTcmqQueue(),
			"tencentcloud_tcmq_topic":                                   tcmq.DataSourceTencentCloudTcmqTopic(),
			"tencentcloud_tcmq_subscribe":                               tcmq.DataSourceTencentCloudTcmqSubscribe(),
			"tencentcloud_as_instances":                                 as.DataSourceTencentCloudAsInstances(),
			"tencentcloud_as_advices":                                   as.DataSourceTencentCloudAsAdvices(),
			"tencentcloud_as_limits":                                    as.DataSourceTencentCloudAsLimits(),
			"tencentcloud_as_last_activity":                             as.DataSourceTencentCloudAsLastActivity(),
			"tencentcloud_cynosdb_accounts":                             cynosdb.DataSourceTencentCloudCynosdbAccounts(),
			"tencentcloud_cynosdb_cluster_instance_groups":              cynosdb.DataSourceTencentCloudCynosdbClusterInstanceGroups(),
			"tencentcloud_cynosdb_cluster_params":                       cynosdb.DataSourceTencentCloudCynosdbClusterParams(),
			"tencentcloud_cynosdb_param_templates":                      cynosdb.DataSourceTencentCloudCynosdbParamTemplates(),
			"tencentcloud_cynosdb_zone":                                 cynosdb.DataSourceTencentCloudCynosdbZone(),
			"tencentcloud_cvm_instances_modification":                   cvm.DataSourceTencentCloudCvmInstancesModification(),
			"tencentcloud_cynosdb_audit_logs":                           cynosdb.DataSourceTencentCloudCynosdbAuditLogs(),
			"tencentcloud_cynosdb_backup_download_url":                  cynosdb.DataSourceTencentCloudCynosdbBackupDownloadUrl(),
			"tencentcloud_cynosdb_binlog_download_url":                  cynosdb.DataSourceTencentCloudCynosdbBinlogDownloadUrl(),
			"tencentcloud_cynosdb_cluster_detail_databases":             cynosdb.DataSourceTencentCloudCynosdbClusterDetailDatabases(),
			"tencentcloud_cynosdb_cluster_param_logs":                   cynosdb.DataSourceTencentCloudCynosdbClusterParamLogs(),
			"tencentcloud_cynosdb_cluster":                              cynosdb.DataSourceTencentCloudCynosdbCluster(),
			"tencentcloud_cynosdb_describe_instance_slow_queries":       cynosdb.DataSourceTencentCloudCynosdbDescribeInstanceSlowQueries(),
			"tencentcloud_cynosdb_describe_instance_error_logs":         cynosdb.DataSourceTencentCloudCynosdbDescribeInstanceErrorLogs(),
			"tencentcloud_cynosdb_account_all_grant_privileges":         cynosdb.DataSourceTencentCloudCynosdbAccountAllGrantPrivileges(),
			"tencentcloud_cynosdb_resource_package_list":                cynosdb.DataSourceTencentCloudCynosdbResourcePackageList(),
			"tencentcloud_cynosdb_project_security_groups":              cynosdb.DataSourceTencentCloudCynosdbProjectSecurityGroups(),
			"tencentcloud_cynosdb_resource_package_sale_specs":          cynosdb.DataSourceTencentCloudCynosdbResourcePackageSaleSpecs(),
			"tencentcloud_cynosdb_rollback_time_range":                  cynosdb.DataSourceTencentCloudCynosdbRollbackTimeRange(),
			"tencentcloud_cynosdb_proxy_node":                           cynosdb.DataSourceTencentCloudCynosdbProxyNode(),
			"tencentcloud_cynosdb_proxy_version":                        cynosdb.DataSourceTencentCloudCynosdbProxyVersion(),
			"tencentcloud_css_domains":                                  css.DataSourceTencentCloudCssDomains(),
			"tencentcloud_css_backup_stream":                            css.DataSourceTencentCloudCssBackupStream(),
			"tencentcloud_css_deliver_log_down_list":                    css.DataSourceTencentCloudCssDeliverLogDownList(),
			"tencentcloud_css_monitor_report":                           css.DataSourceTencentCloudCssMonitorReport(),
			"tencentcloud_css_pad_templates":                            css.DataSourceTencentCloudCssPadTemplates(),
			"tencentcloud_css_pull_stream_task_status":                  css.DataSourceTencentCloudCssPullStreamTaskStatus(),
			"tencentcloud_css_stream_monitor_list":                      css.DataSourceTencentCloudCssStreamMonitorList(),
			"tencentcloud_css_time_shift_record_detail":                 css.DataSourceTencentCloudCssTimeShiftRecordDetail(),
			"tencentcloud_css_time_shift_stream_list":                   css.DataSourceTencentCloudCssTimeShiftStreamList(),
			"tencentcloud_css_watermarks":                               css.DataSourceTencentCloudCssWatermarks(),
			"tencentcloud_css_xp2p_detail_info_list":                    css.DataSourceTencentCloudCssXp2pDetailInfoList(),
			"tencentcloud_chdfs_access_groups":                          chdfs.DataSourceTencentCloudChdfsAccessGroups(),
			"tencentcloud_chdfs_mount_points":                           chdfs.DataSourceTencentCloudChdfsMountPoints(),
			"tencentcloud_chdfs_file_systems":                           chdfs.DataSourceTencentCloudChdfsFileSystems(),
			"tencentcloud_tcm_mesh":                                     tcm.DataSourceTencentCloudTcmMesh(),
			"tencentcloud_lighthouse_firewall_rules_template":           lighthouse.DataSourceTencentCloudLighthouseFirewallRulesTemplate(),
			"tencentcloud_cvm_instance_vnc_url":                         cvm.DataSourceTencentCloudCvmInstanceVncUrl(),
			"tencentcloud_cvm_disaster_recover_group_quota":             cvm.DataSourceTencentCloudCvmDisasterRecoverGroupQuota(),
			"tencentcloud_cvm_chc_hosts":                                cvm.DataSourceTencentCloudCvmChcHosts(),
			"tencentcloud_cvm_chc_denied_actions":                       cvm.DataSourceTencentCloudCvmChcDeniedActions(),
			"tencentcloud_cvm_image_quota":                              cvm.DataSourceTencentCloudCvmImageQuota(),
			"tencentcloud_cvm_import_image_os":                          cvm.DataSourceTencentCloudCvmImportImageOs(),
			"tencentcloud_tsf_application":                              tsf.DataSourceTencentCloudTsfApplication(),
			"tencentcloud_tsf_application_config":                       tsf.DataSourceTencentCloudTsfApplicationConfig(),
			"tencentcloud_tsf_application_file_config":                  tsf.DataSourceTencentCloudTsfApplicationFileConfig(),
			"tencentcloud_tsf_application_public_config":                tsf.DataSourceTencentCloudTsfApplicationPublicConfig(),
			"tencentcloud_cvm_image_share_permission":                   cvm.DataSourceTencentCloudCvmImageSharePermission(),
			"tencentcloud_tsf_cluster":                                  tsf.DataSourceTencentCloudTsfCluster(),
			"tencentcloud_tsf_microservice":                             tsf.DataSourceTencentCloudTsfMicroservice(),
			"tencentcloud_tsf_unit_rules":                               tsf.DataSourceTencentCloudTsfUnitRules(),
			"tencentcloud_tsf_config_summary":                           tsf.DataSourceTencentCloudTsfConfigSummary(),
			"tencentcloud_tsf_delivery_config_by_group_id":              tsf.DataSourceTencentCloudTsfDeliveryConfigByGroupId(),
			"tencentcloud_tsf_delivery_configs":                         tsf.DataSourceTencentCloudTsfDeliveryConfigs(),
			"tencentcloud_tsf_public_config_summary":                    tsf.DataSourceTencentCloudTsfPublicConfigSummary(),
			"tencentcloud_tsf_api_group":                                tsf.DataSourceTencentCloudTsfApiGroup(),
			"tencentcloud_tsf_application_attribute":                    tsf.DataSourceTencentCloudTsfApplicationAttribute(),
			"tencentcloud_tsf_business_log_configs":                     tsf.DataSourceTencentCloudTsfBusinessLogConfigs(),
			"tencentcloud_tsf_api_detail":                               tsf.DataSourceTencentCloudTsfApiDetail(),
			"tencentcloud_tsf_microservice_api_version":                 tsf.DataSourceTencentCloudTsfMicroserviceApiVersion(),
			"tencentcloud_tsf_repository":                               tsf.DataSourceTencentCloudTsfRepository(),
			"tencentcloud_tsf_pod_instances":                            tsf.DataSourceTencentCloudTsfPodInstances(),
			"tencentcloud_tsf_gateway_all_group_apis":                   tsf.DataSourceTencentCloudTsfGatewayAllGroupApis(),
			"tencentcloud_tsf_group_gateways":                           tsf.DataSourceTencentCloudTsfGroupGateways(),
			"tencentcloud_tsf_usable_unit_namespaces":                   tsf.DataSourceTencentCloudTsfUsableUnitNamespaces(),
			"tencentcloud_tsf_group_instances":                          tsf.DataSourceTencentCloudTsfGroupInstances(),
			"tencentcloud_tsf_group_config_release":                     tsf.DataSourceTencentCloudTsfGroupConfigRelease(),
			"tencentcloud_tsf_container_group":                          tsf.DataSourceTencentCloudTsfContainerGroup(),
			"tencentcloud_tsf_groups":                                   tsf.DataSourceTencentCloudTsfGroups(),
			"tencentcloud_tsf_ms_api_list":                              tsf.DataSourceTencentCloudTsfMsApiList(),
			"tencentcloud_lighthouse_bundle":                            lighthouse.DataSourceTencentCloudLighthouseBundle(),
			"tencentcloud_api_gateway_api_docs":                         apigateway.DataSourceTencentCloudAPIGatewayAPIDocs(),
			"tencentcloud_api_gateway_api_apps":                         apigateway.DataSourceTencentCloudAPIGatewayAPIApps(),
			"tencentcloud_tse_access_address":                           tse.DataSourceTencentCloudTseAccessAddress(),
			"tencentcloud_tse_nacos_replicas":                           tse.DataSourceTencentCloudTseNacosReplicas(),
			"tencentcloud_tse_nacos_server_interfaces":                  tse.DataSourceTencentCloudTseNacosServerInterfaces(),
			"tencentcloud_tse_zookeeper_replicas":                       tse.DataSourceTencentCloudTseZookeeperReplicas(),
			"tencentcloud_tse_zookeeper_server_interfaces":              tse.DataSourceTencentCloudTseZookeeperServerInterfaces(),
			"tencentcloud_tse_groups":                                   tse.DataSourceTencentCloudTseGroups(),
			"tencentcloud_tse_gateways":                                 tse.DataSourceTencentCloudTseGateways(),
			"tencentcloud_tse_gateway_nodes":                            tse.DataSourceTencentCloudTseGatewayNodes(),
			"tencentcloud_tse_gateway_routes":                           tse.DataSourceTencentCloudTseGatewayRoutes(),
			"tencentcloud_tse_gateway_canary_rules":                     tse.DataSourceTencentCloudTseGatewayCanaryRules(),
			"tencentcloud_tse_gateway_services":                         tse.DataSourceTencentCloudTseGatewayServices(),
			"tencentcloud_tse_gateway_certificates":                     tse.DataSourceTencentCloudTseGatewayCertificates(),
			"tencentcloud_lighthouse_modify_instance_bundle":            lighthouse.DataSourceTencentCloudLighthouseModifyInstanceBundle(),
			"tencentcloud_lighthouse_zone":                              lighthouse.DataSourceTencentCloudLighthouseZone(),
			"tencentcloud_lighthouse_scene":                             lighthouse.DataSourceTencentCloudLighthouseScene(),
			"tencentcloud_lighthouse_all_scene":                         lighthouse.DataSourceTencentCloudLighthouseAllScene(),
			"tencentcloud_lighthouse_reset_instance_blueprint":          lighthouse.DataSourceTencentCloudLighthouseResetInstanceBlueprint(),
			"tencentcloud_lighthouse_region":                            lighthouse.DataSourceTencentCloudLighthouseRegion(),
			"tencentcloud_lighthouse_instance_vnc_url":                  lighthouse.DataSourceTencentCloudLighthouseInstanceVncUrl(),
			"tencentcloud_lighthouse_instance_traffic_package":          lighthouse.DataSourceTencentCloudLighthouseInstanceTrafficPackage(),
			"tencentcloud_lighthouse_instance_disk_num":                 lighthouse.DataSourceTencentCloudLighthouseInstanceDiskNum(),
			"tencentcloud_lighthouse_instance_blueprint":                lighthouse.DataSourceTencentCloudLighthouseInstanceBlueprint(),
			"tencentcloud_lighthouse_disk_config":                       lighthouse.DataSourceTencentCloudLighthouseDiskConfig(),
			"tencentcloud_lighthouse_disks":                             lighthouse.DataSourceTencentCloudLighthouseInstanceDisks(),
			"tencentcloud_clickhouse_backup_jobs":                       cdwch.DataSourceTencentCloudClickhouseBackupJobs(),
			"tencentcloud_clickhouse_backup_job_detail":                 cdwch.DataSourceTencentCloudClickhouseBackupJobDetail(),
			"tencentcloud_clickhouse_backup_tables":                     cdwch.DataSourceTencentCloudClickhouseBackupTables(),
			"tencentcloud_cls_shipper_tasks":                            cls.DataSourceTencentCloudClsShipperTasks(),
			"tencentcloud_cls_machines":                                 cls.DataSourceTencentCloudClsMachines(),
			"tencentcloud_cls_machine_group_configs":                    cls.DataSourceTencentCloudClsMachineGroupConfigs(),
			"tencentcloud_cls_logsets":                                  cls.DataSourceTencentCloudClsLogsets(),
			"tencentcloud_eb_search":                                    eb.DataSourceTencentCloudEbSearch(),
			"tencentcloud_eb_bus":                                       eb.DataSourceTencentCloudEbBus(),
			"tencentcloud_eb_event_rules":                               eb.DataSourceTencentCloudEbEventRules(),
			"tencentcloud_eb_platform_event_names":                      eb.DataSourceTencentCloudEbPlatformEventNames(),
			"tencentcloud_eb_platform_event_patterns":                   eb.DataSourceTencentCloudEbPlatformEventPatterns(),
			"tencentcloud_eb_platform_products":                         eb.DataSourceTencentCloudEbPlatformProducts(),
			"tencentcloud_eb_plateform_event_template":                  eb.DataSourceTencentCloudEbPlateformEventTemplate(),
			"tencentcloud_wedata_rule_templates":                        wedata.DataSourceTencentCloudWedataRuleTemplates(),
			"tencentcloud_wedata_data_source_list":                      wedata.DataSourceTencentCloudWedataDataSourceList(),
			"tencentcloud_wedata_data_source_without_info":              wedata.DataSourceTencentCloudWedataDataSourceWithoutInfo(),
			"tencentcloud_private_dns_records":                          privatedns.DataSourceTencentCloudPrivateDnsRecords(),
			"tencentcloud_private_dns_private_zone_list":                privatedns.DataSourceTencentCloudPrivateDnsPrivateZoneList(),
			"tencentcloud_private_dns_forward_rules":                    privatedns.DataSourceTencentCloudPrivateDnsForwardRules(),
			"tencentcloud_private_dns_end_points":                       privatedns.DataSourceTencentCloudPrivateDnsEndPoints(),
			"tencentcloud_waf_ciphers":                                  waf.DataSourceTencentCloudWafCiphers(),
			"tencentcloud_waf_tls_versions":                             waf.DataSourceTencentCloudWafTlsVersions(),
			"tencentcloud_waf_domains":                                  waf.DataSourceTencentCloudWafDomains(),
			"tencentcloud_waf_find_domains":                             waf.DataSourceTencentCloudWafFindDomains(),
			"tencentcloud_waf_ports":                                    waf.DataSourceTencentCloudWafPorts(),
			"tencentcloud_waf_user_domains":                             waf.DataSourceTencentCloudWafUserDomains(),
			"tencentcloud_waf_attack_log_histogram":                     waf.DataSourceTencentCloudWafAttackLogHistogram(),
			"tencentcloud_waf_attack_log_list":                          waf.DataSourceTencentCloudWafAttackLogList(),
			"tencentcloud_waf_attack_overview":                          waf.DataSourceTencentCloudWafAttackOverview(),
			"tencentcloud_waf_attack_total_count":                       waf.DataSourceTencentCloudWafAttackTotalCount(),
			"tencentcloud_waf_peak_points":                              waf.DataSourceTencentCloudWafPeakPoints(),
			"tencentcloud_waf_instance_qps_limit":                       waf.DataSourceTencentCloudWafInstanceQpsLimit(),
			"tencentcloud_waf_user_clb_regions":                         waf.DataSourceTencentCloudWafUserClbRegions(),
			"tencentcloud_cfw_nat_fw_switches":                          cfw.DataSourceTencentCloudCfwNatFwSwitches(),
			"tencentcloud_cfw_vpc_fw_switches":                          cfw.DataSourceTencentCloudCfwVpcFwSwitches(),
			"tencentcloud_cfw_edge_fw_switches":                         cfw.DataSourceTencentCloudCfwEdgeFwSwitches(),
			"tencentcloud_cwp_machines_simple":                          cwp.DataSourceTencentCloudCwpMachinesSimple(),
			"tencentcloud_ses_receivers":                                ses.DataSourceTencentCloudSesReceivers(),
			"tencentcloud_ses_send_tasks":                               ses.DataSourceTencentCloudSesSendTasks(),
			"tencentcloud_ses_email_identities":                         ses.DataSourceTencentCloudSesEmailIdentities(),
			"tencentcloud_ses_black_email_address":                      ses.DataSourceTencentCloudSesBlackEmailAddress(),
			"tencentcloud_ses_statistics_report":                        ses.DataSourceTencentCloudSesStatisticsReport(),
			"tencentcloud_ses_send_email_status":                        ses.DataSourceTencentCloudSesSendEmailStatus(),
			"tencentcloud_organization_org_financial_by_member":         tco.DataSourceTencentCloudOrganizationOrgFinancialByMember(),
			"tencentcloud_organization_org_financial_by_month":          tco.DataSourceTencentCloudOrganizationOrgFinancialByMonth(),
			"tencentcloud_organization_org_financial_by_product":        tco.DataSourceTencentCloudOrganizationOrgFinancialByProduct(),
			"tencentcloud_organization_org_auth_node":                   tco.DataSourceTencentCloudOrganizationOrgAuthNode(),
			"tencentcloud_organization_members":                         tco.DataSourceTencentCloudOrganizationMembers(),
			"tencentcloud_organization_services":                        tco.DataSourceTencentCloudOrganizationServices(),
			"tencentcloud_identity_center_groups":                       tco.DataSourceTencentCloudIdentityCenterGroups(),
			"tencentcloud_identity_center_role_configurations":          tco.DataSourceTencentCloudIdentityCenterRoleConfigurations(),
			"tencentcloud_identity_center_users":                        tco.DataSourceTencentCloudIdentityCenterUsers(),
			"tencentcloud_organization_nodes":                           tco.DataSourceTencentCloudOrganizationNodes(),
			"tencentcloud_organization_org_share_unit_resources":        tco.DataSourceTencentCloudOrganizationOrgShareUnitResources(),
			"tencentcloud_organization_org_share_units":                 tco.DataSourceTencentCloudOrganizationOrgShareUnits(),
			"tencentcloud_organization_org_share_unit_members":          tco.DataSourceTencentCloudOrganizationOrgShareUnitMembers(),
			"tencentcloud_role_configuration_provisionings":             tco.DataSourceTencentCloudRoleConfigurationProvisionings(),
			"tencentcloud_pts_scenario_with_jobs":                       pts.DataSourceTencentCloudPtsScenarioWithJobs(),
			"tencentcloud_cam_list_attached_user_policy":                cam.DataSourceTencentCloudCamListAttachedUserPolicy(),
			"tencentcloud_cam_secret_last_used_time":                    cam.DataSourceTencentCloudCamSecretLastUsedTime(),
			"tencentcloud_cam_policy_granting_service_access":           cam.DataSourceTencentCloudCamPolicyGrantingServiceAccess(),
			"tencentcloud_cam_group_user_account":                       cam.DataSourceTencentCloudCamGroupUserAccount(),
			"tencentcloud_dlc_check_data_engine_image_can_be_rollback":  dlc.DataSourceTencentCloudDlcCheckDataEngineImageCanBeRollback(),
			"tencentcloud_dlc_check_data_engine_image_can_be_upgrade":   dlc.DataSourceTencentCloudDlcCheckDataEngineImageCanBeUpgrade(),
			"tencentcloud_dlc_describe_user_type":                       dlc.DataSourceTencentCloudDlcDescribeUserType(),
			"tencentcloud_dlc_describe_user_info":                       dlc.DataSourceTencentCloudDlcDescribeUserInfo(),
			"tencentcloud_dlc_describe_user_roles":                      dlc.DataSourceTencentCloudDlcDescribeUserRoles(),
			"tencentcloud_dlc_describe_data_engine":                     dlc.DataSourceTencentCloudDlcDescribeDataEngine(),
			"tencentcloud_dlc_describe_data_engine_image_versions":      dlc.DataSourceTencentCloudDlcDescribeDataEngineImageVersions(),
			"tencentcloud_dlc_describe_data_engine_python_spark_images": dlc.DataSourceTencentCloudDlcDescribeDataEnginePythonSparkImages(),
			"tencentcloud_dlc_describe_engine_usage_info":               dlc.DataSourceTencentCloudDlcDescribeEngineUsageInfo(),
			"tencentcloud_dlc_describe_work_group_info":                 dlc.DataSourceTencentCloudDlcDescribeWorkGroupInfo(),
			"tencentcloud_dlc_check_data_engine_config_pairs_validity":  dlc.DataSourceTencentCloudDlcCheckDataEngineConfigPairsValidity(),
			"tencentcloud_dlc_describe_updatable_data_engines":          dlc.DataSourceTencentCloudDlcDescribeUpdatableDataEngines(),
			"tencentcloud_dlc_describe_data_engine_events":              dlc.DataSourceTencentCloudDlcDescribeDataEngineEvents(),
			"tencentcloud_bi_project":                                   bi.DataSourceTencentCloudBiProject(),
			"tencentcloud_bi_user_project":                              bi.DataSourceTencentCloudBiUserProject(),
			"tencentcloud_antiddos_basic_device_status":                 antiddos.DataSourceTencentCloudAntiddosBasicDeviceStatus(),
			"tencentcloud_antiddos_bgp_biz_trend":                       antiddos.DataSourceTencentCloudAntiddosBgpBizTrend(),
			"tencentcloud_antiddos_list_listener":                       antiddos.DataSourceTencentCloudAntiddosListListener(),
			"tencentcloud_antiddos_overview_attack_trend":               antiddos.DataSourceTencentCloudAntiddosOverviewAttackTrend(),
			"tencentcloud_kubernetes_cluster_instances":                 tke.DataSourceTencentCloudKubernetesClusterInstances(),
			"tencentcloud_kubernetes_cluster_node_pools":                tke.DataSourceTencentCloudKubernetesClusterNodePools(),
			"tencentcloud_clickhouse_spec":                              cdwch.DataSourceTencentCloudClickhouseSpec(),
			"tencentcloud_clickhouse_instance_shards":                   cdwch.DataSourceTencentCloudClickhouseInstanceShards(),
			"tencentcloud_clickhouse_instance_nodes":                    cdwch.DataSourceTencentCloudClickhouseInstanceNodes(),
			"tencentcloud_organization_org_share_area":                  tco.DataSourceTencentCloudOrganizationOrgShareArea(),
			"tencentcloud_kubernetes_cluster_native_node_pools":         tke.DataSourceTencentCloudKubernetesClusterNativeNodePools(),
			"tencentcloud_cdc_dedicated_cluster_hosts":                  cdc.DataSourceTencentCloudCdcDedicatedClusterHosts(),
			"tencentcloud_cdc_dedicated_cluster_instance_types":         cdc.DataSourceTencentCloudCdcDedicatedClusterInstanceTypes(),
			"tencentcloud_cdc_dedicated_cluster_orders":                 cdc.DataSourceTencentCloudCdcDedicatedClusterOrders(),
			"tencentcloud_cdwdoris_instances":                           cdwdoris.DataSourceTencentCloudCdwdorisInstances(),
			"tencentcloud_lite_hbase_instances":                         emr.DataSourceTencentCloudLiteHbaseInstances()},

		ResourcesMap: map[string]*schema.Resource{
			"tencentcloud_project":                                                                  project.ResourceTencentCloudProject(),
			"tencentcloud_emr_cluster":                                                              emr.ResourceTencentCloudEmrCluster(),
			"tencentcloud_emr_user_manager":                                                         emr.ResourceTencentCloudEmrUserManager(),
			"tencentcloud_instance":                                                                 cvm.ResourceTencentCloudInstance(),
			"tencentcloud_instance_set":                                                             cvm.ResourceTencentCloudInstanceSet(),
			"tencentcloud_reserved_instance":                                                        cvm.ResourceTencentCloudReservedInstance(),
			"tencentcloud_key_pair":                                                                 cvm.ResourceTencentCloudKeyPair(),
			"tencentcloud_placement_group":                                                          cvm.ResourceTencentCloudPlacementGroup(),
			"tencentcloud_cbs_snapshot":                                                             cbs.ResourceTencentCloudCbsSnapshot(),
			"tencentcloud_cbs_snapshot_policy":                                                      cbs.ResourceTencentCloudCbsSnapshotPolicy(),
			"tencentcloud_cbs_storage":                                                              cbs.ResourceTencentCloudCbsStorage(),
			"tencentcloud_cbs_storage_set":                                                          cbs.ResourceTencentCloudCbsStorageSet(),
			"tencentcloud_cbs_storage_attachment":                                                   cbs.ResourceTencentCloudCbsStorageAttachment(),
			"tencentcloud_cbs_storage_set_attachment":                                               cbs.ResourceTencentCloudCbsStorageSetAttachment(),
			"tencentcloud_cbs_snapshot_policy_attachment":                                           cbs.ResourceTencentCloudCbsSnapshotPolicyAttachment(),
			"tencentcloud_vpc":                                                                      vpc.ResourceTencentCloudVpcInstance(),
			"tencentcloud_vpc_acl":                                                                  vpc.ResourceTencentCloudVpcACL(),
			"tencentcloud_vpc_acl_attachment":                                                       vpc.ResourceTencentCloudVpcAclAttachment(),
			"tencentcloud_vpc_network_acl_quintuple":                                                vpc.ResourceTencentCloudVpcNetworkAclQuintuple(),
			"tencentcloud_vpc_notify_routes":                                                        vpc.ResourceTencentCloudVpcNotifyRoutes(),
			"tencentcloud_vpc_bandwidth_package":                                                    vpc.ResourceTencentCloudVpcBandwidthPackage(),
			"tencentcloud_vpc_bandwidth_package_attachment":                                         vpc.ResourceTencentCloudVpcBandwidthPackageAttachment(),
			"tencentcloud_vpc_traffic_package":                                                      vpc.ResourceTencentCloudVpcTrafficPackage(),
			"tencentcloud_vpc_snapshot_policy":                                                      vpc.ResourceTencentCloudVpcSnapshotPolicy(),
			"tencentcloud_vpc_snapshot_policy_attachment":                                           vpc.ResourceTencentCloudVpcSnapshotPolicyAttachment(),
			"tencentcloud_vpc_snapshot_policy_config":                                               vpc.ResourceTencentCloudVpcSnapshotPolicyConfig(),
			"tencentcloud_vpc_net_detect":                                                           vpc.ResourceTencentCloudVpcNetDetect(),
			"tencentcloud_vpc_flow_log_config":                                                      fl.ResourceTencentCloudVpcFlowLogConfig(),
			"tencentcloud_vpc_classic_link_attachment":                                              vpc.ResourceTencentCloudVpcClassicLinkAttachment(),
			"tencentcloud_vpc_dhcp_ip":                                                              vpc.ResourceTencentCloudVpcDhcpIp(),
			"tencentcloud_vpc_ipv6_cidr_block":                                                      vpc.ResourceTencentCloudVpcIpv6CidrBlock(),
			"tencentcloud_vpc_ipv6_subnet_cidr_block":                                               vpc.ResourceTencentCloudVpcIpv6SubnetCidrBlock(),
			"tencentcloud_vpc_ipv6_eni_address":                                                     vpc.ResourceTencentCloudVpcIpv6EniAddress(),
			"tencentcloud_vpc_dhcp_associate_address":                                               vpc.ResourceTencentCloudVpcDhcpAssociateAddress(),
			"tencentcloud_vpc_local_gateway":                                                        vpc.ResourceTencentCloudVpcLocalGateway(),
			"tencentcloud_vpc_resume_snapshot_instance":                                             vpc.ResourceTencentCloudVpcResumeSnapshotInstance(),
			"tencentcloud_ipv6_address_bandwidth":                                                   vpc.ResourceTencentCloudIpv6AddressBandwidth(),
			"tencentcloud_subnet":                                                                   vpc.ResourceTencentCloudVpcSubnet(),
			"tencentcloud_route_entry":                                                              vpc.ResourceTencentCloudRouteEntry(),
			"tencentcloud_route_table_entry":                                                        vpc.ResourceTencentCloudVpcRouteEntry(),
			"tencentcloud_route_table":                                                              vpc.ResourceTencentCloudVpcRouteTable(),
			"tencentcloud_route_table_association":                                                  vpc.ResourceTencentCloudRouteTableAssociation(),
			"tencentcloud_dnat":                                                                     vpc.ResourceTencentCloudDnat(),
			"tencentcloud_nat_gateway":                                                              vpc.ResourceTencentCloudNatGateway(),
			"tencentcloud_nat_gateway_snat":                                                         vpc.ResourceTencentCloudNatGatewaySnat(),
			"tencentcloud_nat_refresh_nat_dc_route":                                                 vpc.ResourceTencentCloudNatRefreshNatDcRoute(),
			"tencentcloud_vpc_private_nat_gateway":                                                  vpc.ResourceTencentCloudVpcPrivateNatGateway(),
			"tencentcloud_oceanus_job":                                                              oceanus.ResourceTencentCloudOceanusJob(),
			"tencentcloud_oceanus_job_config":                                                       oceanus.ResourceTencentCloudOceanusJobConfig(),
			"tencentcloud_oceanus_job_copy":                                                         oceanus.ResourceTencentCloudOceanusJobCopy(),
			"tencentcloud_oceanus_run_job":                                                          oceanus.ResourceTencentCloudOceanusRunJob(),
			"tencentcloud_oceanus_stop_job":                                                         oceanus.ResourceTencentCloudOceanusStopJob(),
			"tencentcloud_oceanus_trigger_job_savepoint":                                            oceanus.ResourceTencentCloudOceanusTriggerJobSavepoint(),
			"tencentcloud_oceanus_resource":                                                         oceanus.ResourceTencentCloudOceanusResource(),
			"tencentcloud_oceanus_resource_config":                                                  oceanus.ResourceTencentCloudOceanusResourceConfig(),
			"tencentcloud_oceanus_work_space":                                                       oceanus.ResourceTencentCloudOceanusWorkSpace(),
			"tencentcloud_oceanus_folder":                                                           oceanus.ResourceTencentCloudOceanusFolder(),
			"tencentcloud_tag":                                                                      tag.ResourceTencentCloudTag(),
			"tencentcloud_tag_attachment":                                                           tag.ResourceTencentCloudTagAttachment(),
			"tencentcloud_eip":                                                                      cvm.ResourceTencentCloudEip(),
			"tencentcloud_eip_association":                                                          cvm.ResourceTencentCloudEipAssociation(),
			"tencentcloud_eip_address_transform":                                                    cvm.ResourceTencentCloudEipAddressTransform(),
			"tencentcloud_eip_public_address_adjust":                                                cvm.ResourceTencentCloudEipPublicAddressAdjust(),
			"tencentcloud_eip_normal_address_return":                                                cvm.ResourceTencentCloudEipNormalAddressReturn(),
			"tencentcloud_eni":                                                                      vpc.ResourceTencentCloudEni(),
			"tencentcloud_eni_attachment":                                                           vpc.ResourceTencentCloudEniAttachment(),
			"tencentcloud_eni_sg_attachment":                                                        vpc.ResourceTencentCloudEniSgAttachment(),
			"tencentcloud_eni_ipv6_address":                                                         vpc.ResourceTencentCloudEniIpv6Address(),
			"tencentcloud_eni_ipv4_address":                                                         vpc.ResourceTencentCloudEniIpv4Address(),
			"tencentcloud_ccn":                                                                      ccn.ResourceTencentCloudCcn(),
			"tencentcloud_ccn_attachment":                                                           ccn.ResourceTencentCloudCcnAttachment(),
			"tencentcloud_ccn_bandwidth_limit":                                                      ccn.ResourceTencentCloudCcnBandwidthLimit(),
			"tencentcloud_ccn_route_table":                                                          ccn.ResourceTencentCloudCcnRouteTable(),
			"tencentcloud_ccn_route_table_input_policies":                                           ccn.ResourceTencentCloudCcnRouteTableInputPolicies(),
			"tencentcloud_ccn_route_table_broadcast_policies":                                       ccn.ResourceTencentCloudCcnRouteTableBroadcastPolicies(),
			"tencentcloud_ccn_route_table_selection_policies":                                       ccn.ResourceTencentCloudCcnRouteTableSelectionPolicies(),
			"tencentcloud_ccn_routes":                                                               ccn.ResourceTencentCloudCcnRoutes(),
			"tencentcloud_ccn_route_table_associate_instance_config":                                ccn.ResourceTencentCloudCcnRouteTableAssociateInstanceConfig(),
			"tencentcloud_ccn_instances_accept_attach":                                              ccn.ResourceTencentCloudCcnInstancesAcceptAttach(),
			"tencentcloud_ccn_instances_reject_attach":                                              ccn.ResourceTencentCloudCcnInstancesRejectAttach(),
			"tencentcloud_ccn_instances_reset_attach":                                               ccn.ResourceTencentCloudCcnInstancesResetAttach(),
			"tencentcloud_dc_instance":                                                              dc.ResourceTencentCloudDcInstance(),
			"tencentcloud_dcx":                                                                      dc.ResourceTencentCloudDcxInstance(),
			"tencentcloud_dcx_extra_config":                                                         dc.ResourceTencentCloudDcxExtraConfig(),
			"tencentcloud_dc_share_dcx_config":                                                      dc.ResourceTencentCloudDcShareDcxConfig(),
			"tencentcloud_dc_internet_address":                                                      dc.ResourceTencentCloudDcInternetAddress(),
			"tencentcloud_dc_internet_address_config":                                               dc.ResourceTencentCloudDcInternetAddressConfig(),
			"tencentcloud_dc_gateway":                                                               dcg.ResourceTencentCloudDcGatewayInstance(),
			"tencentcloud_dc_gateway_ccn_route":                                                     dcg.ResourceTencentCloudDcGatewayCcnRouteInstance(),
			"tencentcloud_dc_gateway_attachment":                                                    dcg.ResourceTencentCloudDcGatewayAttachment(),
			"tencentcloud_vpn_customer_gateway":                                                     vpn.ResourceTencentCloudVpnCustomerGateway(),
			"tencentcloud_vpn_gateway":                                                              vpn.ResourceTencentCloudVpnGateway(),
			"tencentcloud_vpn_gateway_route":                                                        vpn.ResourceTencentCloudVpnGatewayRoute(),
			"tencentcloud_vpn_connection":                                                           vpn.ResourceTencentCloudVpnConnection(),
			"tencentcloud_vpn_ssl_server":                                                           vpn.ResourceTencentCloudVpnSslServer(),
			"tencentcloud_vpn_ssl_client":                                                           vpn.ResourceTencentCloudVpnSslClient(),
			"tencentcloud_vpn_connection_reset":                                                     vpn.ResourceTencentCloudVpnConnectionReset(),
			"tencentcloud_vpn_customer_gateway_configuration_download":                              vpn.ResourceTencentCloudVpnCustomerGatewayConfigurationDownload(),
			"tencentcloud_vpn_gateway_ssl_client_cert":                                              vpn.ResourceTencentCloudVpnGatewaySslClientCert(),
			"tencentcloud_vpn_gateway_ccn_routes":                                                   vpn.ResourceTencentCloudVpnGatewayCcnRoutes(),
			"tencentcloud_ha_vip":                                                                   vpc.ResourceTencentCloudHaVip(),
			"tencentcloud_ha_vip_eip_attachment":                                                    vpc.ResourceTencentCloudHaVipEipAttachment(),
			"tencentcloud_security_group":                                                           vpc.ResourceTencentCloudSecurityGroup(),
			"tencentcloud_security_group_rule":                                                      vpc.ResourceTencentCloudSecurityGroupRule(),
			"tencentcloud_security_group_rule_set":                                                  vpc.ResourceTencentCloudSecurityGroupRuleSet(),
			"tencentcloud_security_group_lite_rule":                                                 vpc.ResourceTencentCloudSecurityGroupLiteRule(),
			"tencentcloud_lb":                                                                       clb.ResourceTencentCloudLB(),
			"tencentcloud_alb_server_attachment":                                                    clb.ResourceTencentCloudAlbServerAttachment(),
			"tencentcloud_clb_instance":                                                             clb.ResourceTencentCloudClbInstance(),
			"tencentcloud_clb_listener":                                                             clb.ResourceTencentCloudClbListener(),
			"tencentcloud_clb_listener_rule":                                                        clb.ResourceTencentCloudClbListenerRule(),
			"tencentcloud_clb_listener_default_domain":                                              clb.ResourceTencentCloudClbListenerDefaultDomain(),
			"tencentcloud_clb_attachment":                                                           clb.ResourceTencentCloudClbServerAttachment(),
			"tencentcloud_clb_redirection":                                                          clb.ResourceTencentCloudClbRedirection(),
			"tencentcloud_clb_target_group":                                                         clb.ResourceTencentCloudClbTargetGroup(),
			"tencentcloud_clb_target_group_instance_attachment":                                     clb.ResourceTencentCloudClbTGAttachmentInstance(),
			"tencentcloud_clb_target_group_attachment":                                              clb.ResourceTencentCloudClbTargetGroupAttachment(),
			"tencentcloud_clb_log_set":                                                              clb.ResourceTencentCloudClbLogSet(),
			"tencentcloud_clb_log_topic":                                                            clb.ResourceTencentCloudClbLogTopic(),
			"tencentcloud_clb_customized_config":                                                    clb.ResourceTencentCloudClbCustomizedConfig(),
			"tencentcloud_clb_snat_ip":                                                              clb.ResourceTencentCloudClbSnatIp(),
			"tencentcloud_clb_function_targets_attachment":                                          clb.ResourceTencentCloudClbFunctionTargetsAttachment(),
			"tencentcloud_clb_instance_mix_ip_target_config":                                        clb.ResourceTencentCloudClbInstanceMixIpTargetConfig(),
			"tencentcloud_clb_instance_sla_config":                                                  clb.ResourceTencentCloudClbInstanceSlaConfig(),
			"tencentcloud_clb_replace_cert_for_lbs":                                                 clb.ResourceTencentCloudClbReplaceCertForLbs(),
			"tencentcloud_clb_security_group_attachment":                                            clb.ResourceTencentCloudClbSecurityGroupAttachment(),
			"tencentcloud_container_cluster":                                                        tke.ResourceTencentCloudContainerCluster(),
			"tencentcloud_container_cluster_instance":                                               tke.ResourceTencentCloudContainerClusterInstance(),
			"tencentcloud_kubernetes_cluster":                                                       tke.ResourceTencentCloudKubernetesCluster(),
			"tencentcloud_kubernetes_cluster_endpoint":                                              tke.ResourceTencentCloudTkeClusterEndpoint(),
			"tencentcloud_eks_cluster":                                                              tke.ResourceTencentCloudEksCluster(),
			"tencentcloud_eks_container_instance":                                                   tke.ResourceTencentCloudEksContainerInstance(),
			"tencentcloud_kubernetes_addon_attachment":                                              tke.ResourceTencentCloudKubernetesAddonAttachment(),
			"tencentcloud_kubernetes_auth_attachment":                                               tke.ResourceTencentCloudKubernetesAuthAttachment(),
			"tencentcloud_kubernetes_as_scaling_group":                                              tke.ResourceTencentCloudKubernetesAsScalingGroup(),
			"tencentcloud_kubernetes_scale_worker":                                                  tke.ResourceTencentCloudKubernetesScaleWorker(),
			"tencentcloud_kubernetes_cluster_attachment":                                            tke.ResourceTencentCloudKubernetesClusterAttachment(),
			"tencentcloud_kubernetes_node_pool":                                                     tke.ResourceTencentCloudKubernetesNodePool(),
			"tencentcloud_kubernetes_backup_storage_location":                                       tke.ResourceTencentCloudKubernetesBackupStorageLocation(),
			"tencentcloud_kubernetes_serverless_node_pool":                                          tke.ResourceTencentCloudKubernetesServerlessNodePool(),
			"tencentcloud_kubernetes_encryption_protection":                                         tke.ResourceTencentCloudKubernetesEncryptionProtection(),
			"tencentcloud_kubernetes_cluster_master_attachment":                                     tke.ResourceTencentCloudKubernetesClusterMasterAttachment(),
			"tencentcloud_kubernetes_cluster_audit":                                                 tke.ResourceTencentCloudKubernetesClusterAudit(),
			"tencentcloud_mysql_backup_policy":                                                      cdb.ResourceTencentCloudMysqlBackupPolicy(),
			"tencentcloud_mysql_account":                                                            cdb.ResourceTencentCloudMysqlAccount(),
			"tencentcloud_mysql_account_privilege":                                                  cdb.ResourceTencentCloudMysqlAccountPrivilege(),
			"tencentcloud_mysql_privilege":                                                          cdb.ResourceTencentCloudMysqlPrivilege(),
			"tencentcloud_mysql_instance":                                                           cdb.ResourceTencentCloudMysqlInstance(),
			"tencentcloud_mysql_database":                                                           cdb.ResourceTencentCloudMysqlDatabase(),
			"tencentcloud_mysql_readonly_instance":                                                  cdb.ResourceTencentCloudMysqlReadonlyInstance(),
			"tencentcloud_mysql_dr_instance":                                                        cdb.ResourceTencentCloudMysqlDrInstance(),
			"tencentcloud_mysql_time_window":                                                        cdb.ResourceTencentCloudMysqlTimeWindow(),
			"tencentcloud_mysql_param_template":                                                     cdb.ResourceTencentCloudMysqlParamTemplate(),
			"tencentcloud_mysql_security_groups_attachment":                                         cdb.ResourceTencentCloudMysqlSecurityGroupsAttachment(),
			"tencentcloud_mysql_deploy_group":                                                       cdb.ResourceTencentCloudMysqlDeployGroup(),
			"tencentcloud_mysql_local_binlog_config":                                                cdb.ResourceTencentCloudMysqlLocalBinlogConfig(),
			"tencentcloud_mysql_audit_log_file":                                                     cdb.ResourceTencentCloudMysqlAuditLogFile(),
			"tencentcloud_mysql_backup_download_restriction":                                        cdb.ResourceTencentCloudMysqlBackupDownloadRestriction(),
			"tencentcloud_mysql_renew_db_instance_operation":                                        cdb.ResourceTencentCloudMysqlRenewDbInstanceOperation(),
			"tencentcloud_mysql_backup_encryption_status":                                           cdb.ResourceTencentCloudMysqlBackupEncryptionStatus(),
			"tencentcloud_mysql_db_import_job_operation":                                            cdb.ResourceTencentCloudMysqlDbImportJobOperation(),
			"tencentcloud_mysql_dr_instance_to_mater":                                               cdb.ResourceTencentCloudMysqlDrInstanceToMater(),
			"tencentcloud_mysql_instance_encryption_operation":                                      cdb.ResourceTencentCloudMysqlInstanceEncryptionOperation(),
			"tencentcloud_mysql_isolate_instance":                                                   cdb.ResourceTencentCloudMysqlIsolateInstance(),
			"tencentcloud_mysql_password_complexity":                                                cdb.ResourceTencentCloudMysqlPasswordComplexity(),
			"tencentcloud_mysql_remote_backup_config":                                               cdb.ResourceTencentCloudMysqlRemoteBackupConfig(),
			"tencentcloud_mysql_restart_db_instances_operation":                                     cdb.ResourceTencentCloudMysqlRestartDbInstancesOperation(),
			"tencentcloud_mysql_switch_for_upgrade":                                                 cdb.ResourceTencentCloudMysqlSwitchForUpgrade(),
			"tencentcloud_mysql_rollback":                                                           cdb.ResourceTencentCloudMysqlRollback(),
			"tencentcloud_mysql_rollback_stop":                                                      cdb.ResourceTencentCloudMysqlRollbackStop(),
			"tencentcloud_mysql_ro_group":                                                           cdb.ResourceTencentCloudMysqlRoGroup(),
			"tencentcloud_mysql_ro_instance_ip":                                                     cdb.ResourceTencentCloudMysqlRoInstanceIp(),
			"tencentcloud_mysql_ro_group_load_operation":                                            cdb.ResourceTencentCloudMysqlRoGroupLoadOperation(),
			"tencentcloud_mysql_switch_master_slave_operation":                                      cdb.ResourceTencentCloudMysqlSwitchMasterSlaveOperation(),
			"tencentcloud_mysql_proxy":                                                              cdb.ResourceTencentCloudMysqlProxy(),
			"tencentcloud_mysql_reset_root_account":                                                 cdb.ResourceTencentCloudMysqlResetRootAccount(),
			"tencentcloud_mysql_verify_root_account":                                                cdb.ResourceTencentCloudMysqlVerifyRootAccount(),
			"tencentcloud_mysql_reload_balance_proxy_node":                                          cdb.ResourceTencentCloudMysqlReloadBalanceProxyNode(),
			"tencentcloud_mysql_ro_start_replication":                                               cdb.ResourceTencentCloudMysqlRoStartReplication(),
			"tencentcloud_mysql_ro_stop_replication":                                                cdb.ResourceTencentCloudMysqlRoStopReplication(),
			"tencentcloud_mysql_switch_proxy":                                                       cdb.ResourceTencentCloudMysqlSwitchProxy(),
			"tencentcloud_mysql_ssl":                                                                cdb.ResourceTencentCloudMysqlSsl(),
			"tencentcloud_mysql_cls_log_attachment":                                                 cdb.ResourceTencentCloudMysqlClsLogAttachment(),
			"tencentcloud_cos_bucket":                                                               cos.ResourceTencentCloudCosBucket(),
			"tencentcloud_cos_bucket_object":                                                        cos.ResourceTencentCloudCosBucketObject(),
			"tencentcloud_cos_bucket_referer":                                                       cos.ResourceTencentCloudCosBucketReferer(),
			"tencentcloud_cos_bucket_version":                                                       cos.ResourceTencentCloudCosBucketVersion(),
			"tencentcloud_cfs_file_system":                                                          cfs.ResourceTencentCloudCfsFileSystem(),
			"tencentcloud_cfs_access_group":                                                         cfs.ResourceTencentCloudCfsAccessGroup(),
			"tencentcloud_cfs_access_rule":                                                          cfs.ResourceTencentCloudCfsAccessRule(),
			"tencentcloud_cfs_auto_snapshot_policy":                                                 cfs.ResourceTencentCloudCfsAutoSnapshotPolicy(),
			"tencentcloud_cfs_auto_snapshot_policy_attachment":                                      cfs.ResourceTencentCloudCfsAutoSnapshotPolicyAttachment(),
			"tencentcloud_cfs_snapshot":                                                             cfs.ResourceTencentCloudCfsSnapshot(),
			"tencentcloud_cfs_user_quota":                                                           cfs.ResourceTencentCloudCfsUserQuota(),
			"tencentcloud_cfs_sign_up_cfs_service":                                                  cfs.ResourceTencentCloudCfsSignUpCfsService(),
			"tencentcloud_redis_instance":                                                           crs.ResourceTencentCloudRedisInstance(),
			"tencentcloud_redis_backup_config":                                                      crs.ResourceTencentCloudRedisBackupConfig(),
			"tencentcloud_redis_account":                                                            crs.ResourceTencentCloudRedisAccount(),
			"tencentcloud_redis_param_template":                                                     crs.ResourceTencentCloudRedisParamTemplate(),
			"tencentcloud_redis_connection_config":                                                  crs.ResourceTencentCloudRedisConnectionConfig(),
			"tencentcloud_redis_param":                                                              crs.ResourceTencentCloudRedisParam(),
			"tencentcloud_redis_read_only":                                                          crs.ResourceTencentCloudRedisReadOnly(),
			"tencentcloud_redis_ssl":                                                                crs.ResourceTencentCloudRedisSsl(),
			"tencentcloud_redis_backup_download_restriction":                                        crs.ResourceTencentCloudRedisBackupDownloadRestriction(),
			"tencentcloud_redis_clear_instance_operation":                                           crs.ResourceTencentCloudRedisClearInstanceOperation(),
			"tencentcloud_redis_renew_instance_operation":                                           crs.ResourceTencentCloudRedisRenewInstanceOperation(),
			"tencentcloud_redis_startup_instance_operation":                                         crs.ResourceTencentCloudRedisStartupInstanceOperation(),
			"tencentcloud_redis_upgrade_cache_version_operation":                                    crs.ResourceTencentCloudRedisUpgradeCacheVersionOperation(),
			"tencentcloud_redis_upgrade_multi_zone_operation":                                       crs.ResourceTencentCloudRedisUpgradeMultiZoneOperation(),
			"tencentcloud_redis_upgrade_proxy_version_operation":                                    crs.ResourceTencentCloudRedisUpgradeProxyVersionOperation(),
			"tencentcloud_redis_maintenance_window":                                                 crs.ResourceTencentCloudRedisMaintenanceWindow(),
			"tencentcloud_redis_replica_readonly":                                                   crs.ResourceTencentCloudRedisReplicaReadonly(),
			"tencentcloud_redis_switch_master":                                                      crs.ResourceTencentCloudRedisSwitchMaster(),
			"tencentcloud_redis_replicate_attachment":                                               crs.ResourceTencentCloudRedisReplicateAttachment(),
			"tencentcloud_redis_backup_operation":                                                   crs.ResourceTencentCloudRedisBackupOperation(),
			"tencentcloud_redis_security_group_attachment":                                          crs.ResourceTencentCloudRedisSecurityGroupAttachment(),
			"tencentcloud_redis_log_delivery":                                                       crs.ResourceTencentCloudRedisLogDelivery(),
			"tencentcloud_as_load_balancer":                                                         as.ResourceTencentCloudAsLoadBalancer(),
			"tencentcloud_as_scaling_config":                                                        as.ResourceTencentCloudAsScalingConfig(),
			"tencentcloud_as_scaling_group":                                                         as.ResourceTencentCloudAsScalingGroup(),
			"tencentcloud_as_scaling_group_status":                                                  as.ResourceTencentCloudAsScalingGroupStatus(),
			"tencentcloud_as_attachment":                                                            as.ResourceTencentCloudAsAttachment(),
			"tencentcloud_as_scaling_policy":                                                        as.ResourceTencentCloudAsScalingPolicy(),
			"tencentcloud_as_schedule":                                                              as.ResourceTencentCloudAsSchedule(),
			"tencentcloud_as_lifecycle_hook":                                                        as.ResourceTencentCloudAsLifecycleHook(),
			"tencentcloud_as_notification":                                                          as.ResourceTencentCloudAsNotification(),
			"tencentcloud_as_remove_instances":                                                      as.ResourceTencentCloudAsRemoveInstances(),
			"tencentcloud_as_protect_instances":                                                     as.ResourceTencentCloudAsProtectInstances(),
			"tencentcloud_as_start_instances":                                                       as.ResourceTencentCloudAsStartInstances(),
			"tencentcloud_as_stop_instances":                                                        as.ResourceTencentCloudAsStopInstances(),
			"tencentcloud_as_scale_in_instances":                                                    as.ResourceTencentCloudAsScaleInInstances(),
			"tencentcloud_as_scale_out_instances":                                                   as.ResourceTencentCloudAsScaleOutInstances(),
			"tencentcloud_as_execute_scaling_policy":                                                as.ResourceTencentCloudAsExecuteScalingPolicy(),
			"tencentcloud_as_complete_lifecycle":                                                    as.ResourceTencentCloudAsCompleteLifecycle(),
			"tencentcloud_as_start_instance_refresh":                                                as.ResourceTencentCloudAsStartInstanceRefresh(),
			"tencentcloud_mongodb_instance":                                                         mongodb.ResourceTencentCloudMongodbInstance(),
			"tencentcloud_mongodb_sharding_instance":                                                mongodb.ResourceTencentCloudMongodbShardingInstance(),
			"tencentcloud_mongodb_instance_account":                                                 mongodb.ResourceTencentCloudMongodbInstanceAccount(),
			"tencentcloud_mongodb_instance_backup":                                                  mongodb.ResourceTencentCloudMongodbInstanceBackup(),
			"tencentcloud_mongodb_instance_backup_download_task":                                    mongodb.ResourceTencentCloudMongodbInstanceBackupDownloadTask(),
			"tencentcloud_mongodb_instance_transparent_data_encryption":                             mongodb.ResourceTencentCloudMongodbInstanceTransparentDataEncryption(),
			"tencentcloud_mongodb_instance_backup_rule":                                             mongodb.ResourceTencentCloudMongodbInstanceBackupRule(),
			"tencentcloud_dayu_cc_http_policy":                                                      dayu.ResourceTencentCloudDayuCCHttpPolicy(),
			"tencentcloud_dayu_cc_https_policy":                                                     dayu.ResourceTencentCloudDayuCCHttpsPolicy(),
			"tencentcloud_dayu_ddos_policy":                                                         dayu.ResourceTencentCloudDayuDdosPolicy(),
			"tencentcloud_dayu_cc_policy_v2":                                                        dayuv2.ResourceTencentCloudDayuCCPolicyV2(),
			"tencentcloud_dayu_ddos_policy_v2":                                                      dayuv2.ResourceTencentCloudDayuDdosPolicyV2(),
			"tencentcloud_dayu_ddos_policy_case":                                                    dayu.ResourceTencentCloudDayuDdosPolicyCase(),
			"tencentcloud_dayu_ddos_policy_attachment":                                              dayu.ResourceTencentCloudDayuDdosPolicyAttachment(),
			"tencentcloud_dayu_l4_rule":                                                             dayu.ResourceTencentCloudDayuL4Rule(),
			"tencentcloud_dayu_l4_rule_v2":                                                          dayuv2.ResourceTencentCloudDayuL4RuleV2(),
			"tencentcloud_dayu_l7_rule":                                                             dayu.ResourceTencentCloudDayuL7Rule(),
			"tencentcloud_dayu_l7_rule_v2":                                                          dayuv2.ResourceTencentCloudDayuL7RuleV2(),
			"tencentcloud_dayu_eip":                                                                 dayuv2.ResourceTencentCloudDayuEip(),
			"tencentcloud_gaap_proxy":                                                               gaap.ResourceTencentCloudGaapProxy(),
			"tencentcloud_gaap_realserver":                                                          gaap.ResourceTencentCloudGaapRealserver(),
			"tencentcloud_gaap_layer4_listener":                                                     gaap.ResourceTencentCloudGaapLayer4Listener(),
			"tencentcloud_gaap_layer7_listener":                                                     gaap.ResourceTencentCloudGaapLayer7Listener(),
			"tencentcloud_gaap_http_domain":                                                         gaap.ResourceTencentCloudGaapHttpDomain(),
			"tencentcloud_gaap_http_rule":                                                           gaap.ResourceTencentCloudGaapHttpRule(),
			"tencentcloud_gaap_certificate":                                                         gaap.ResourceTencentCloudGaapCertificate(),
			"tencentcloud_gaap_security_policy":                                                     gaap.ResourceTencentCloudGaapSecurityPolicy(),
			"tencentcloud_gaap_security_rule":                                                       gaap.ResourceTencentCloudGaapSecurityRule(),
			"tencentcloud_gaap_domain_error_page":                                                   gaap.ResourceTencentCloudGaapDomainErrorPageInfo(),
			"tencentcloud_gaap_global_domain_dns":                                                   gaap.ResourceTencentCloudGaapGlobalDomainDns(),
			"tencentcloud_gaap_global_domain":                                                       gaap.ResourceTencentCloudGaapGlobalDomain(),
			"tencentcloud_gaap_custom_header":                                                       gaap.ResourceTencentCloudGaapCustomHeader(),
			"tencentcloud_gaap_proxy_group":                                                         gaap.ResourceTencentCloudGaapProxyGroup(),
			"tencentcloud_ssl_certificate":                                                          ssl.ResourceTencentCloudSslCertificate(),
			"tencentcloud_ssl_pay_certificate":                                                      ssl.ResourceTencentCloudSSLInstance(),
			"tencentcloud_ssl_free_certificate":                                                     ssl.ResourceTencentCloudSSLFreeCertificate(),
			"tencentcloud_cam_role":                                                                 cam.ResourceTencentCloudCamRole(),
			"tencentcloud_cam_role_by_name":                                                         cam.ResourceTencentCloudCamRoleByName(),
			"tencentcloud_cam_user":                                                                 cam.ResourceTencentCloudCamUser(),
			"tencentcloud_cam_policy":                                                               cam.ResourceTencentCloudCamPolicy(),
			"tencentcloud_cam_policy_by_name":                                                       cam.ResourceTencentCloudCamPolicyByName(),
			"tencentcloud_cam_role_policy_attachment":                                               cam.ResourceTencentCloudCamRolePolicyAttachment(),
			"tencentcloud_cam_role_policy_attachment_by_name":                                       cam.ResourceTencentCloudCamRolePolicyAttachmentByName(),
			"tencentcloud_cam_user_policy_attachment":                                               cam.ResourceTencentCloudCamUserPolicyAttachment(),
			"tencentcloud_cam_group_policy_attachment":                                              cam.ResourceTencentCloudCamGroupPolicyAttachment(),
			"tencentcloud_cam_group":                                                                cam.ResourceTencentCloudCamGroup(),
			"tencentcloud_cam_oidc_sso":                                                             cam.ResourceTencentCloudCamOIDCSSO(),
			"tencentcloud_cam_role_sso":                                                             cam.ResourceTencentCloudCamRoleSSO(),
			"tencentcloud_cam_group_membership":                                                     cam.ResourceTencentCloudCamGroupMembership(),
			"tencentcloud_cam_saml_provider":                                                        cam.ResourceTencentCloudCamSAMLProvider(),
			"tencentcloud_cam_service_linked_role":                                                  cam.ResourceTencentCloudCamServiceLinkedRole(),
			"tencentcloud_cam_mfa_flag":                                                             cam.ResourceTencentCloudCamMfaFlag(),
			"tencentcloud_cam_access_key":                                                           cam.ResourceTencentCloudCamAccessKey(),
			"tencentcloud_cam_user_saml_config":                                                     cam.ResourceTencentCloudCamUserSamlConfig(),
			"tencentcloud_cam_tag_role_attachment":                                                  cam.ResourceTencentCloudCamTagRoleAttachment(),
			"tencentcloud_cam_policy_version":                                                       cam.ResourceTencentCloudCamPolicyVersion(),
			"tencentcloud_cam_set_policy_version_config":                                            cam.ResourceTencentCloudCamSetPolicyVersionConfig(),
			"tencentcloud_cam_user_permission_boundary_attachment":                                  cam.ResourceTencentCloudCamUserPermissionBoundaryAttachment(),
			"tencentcloud_cam_role_permission_boundary_attachment":                                  cam.ResourceTencentCloudCamRolePermissionBoundaryAttachment(),
			"tencentcloud_organization_quit_organization_operation":                                 tco.ResourceTencentCloudOrganizationQuitOrganizationOperation(),
			"tencentcloud_ciam_user_group":                                                          ciam.ResourceTencentCloudCiamUserGroup(),
			"tencentcloud_ciam_user_store":                                                          ciam.ResourceTencentCloudCiamUserStore(),
			"tencentcloud_scf_function":                                                             scf.ResourceTencentCloudScfFunction(),
			"tencentcloud_scf_function_version":                                                     scf.ResourceTencentCloudScfFunctionVersion(),
			"tencentcloud_scf_function_event_invoke_config":                                         scf.ResourceTencentCloudScfFunctionEventInvokeConfig(),
			"tencentcloud_scf_reserved_concurrency_config":                                          scf.ResourceTencentCloudScfReservedConcurrencyConfig(),
			"tencentcloud_scf_provisioned_concurrency_config":                                       scf.ResourceTencentCloudScfProvisionedConcurrencyConfig(),
			"tencentcloud_scf_invoke_function":                                                      scf.ResourceTencentCloudScfInvokeFunction(),
			"tencentcloud_scf_sync_invoke_function":                                                 scf.ResourceTencentCloudScfSyncInvokeFunction(),
			"tencentcloud_scf_terminate_async_event":                                                scf.ResourceTencentCloudScfTerminateAsyncEvent(),
			"tencentcloud_scf_namespace":                                                            scf.ResourceTencentCloudScfNamespace(),
			"tencentcloud_scf_layer":                                                                scf.ResourceTencentCloudScfLayer(),
			"tencentcloud_scf_function_alias":                                                       scf.ResourceTencentCloudScfFunctionAlias(),
			"tencentcloud_scf_trigger_config":                                                       scf.ResourceTencentCloudScfTriggerConfig(),
			"tencentcloud_scf_custom_domain":                                                        scf.ResourceTencentCloudScfCustomDomain(),
			"tencentcloud_tcaplus_cluster":                                                          tcaplusdb.ResourceTencentCloudTcaplusCluster(),
			"tencentcloud_tcaplus_tablegroup":                                                       tcaplusdb.ResourceTencentCloudTcaplusTableGroup(),
			"tencentcloud_tcaplus_idl":                                                              tcaplusdb.ResourceTencentCloudTcaplusIdl(),
			"tencentcloud_tcaplus_table":                                                            tcaplusdb.ResourceTencentCloudTcaplusTable(),
			"tencentcloud_cdn_domain":                                                               cdn.ResourceTencentCloudCdnDomain(),
			"tencentcloud_cdn_url_push":                                                             cdn.ResourceTencentCloudUrlPush(),
			"tencentcloud_cdn_url_purge":                                                            cdn.ResourceTencentCloudUrlPurge(),
			"tencentcloud_monitor_policy_group":                                                     monitor.ResourceTencentCloudMonitorPolicyGroup(),
			"tencentcloud_monitor_binding_object":                                                   monitor.ResourceTencentCloudMonitorBindingObject(),
			"tencentcloud_monitor_policy_binding_object":                                            monitor.ResourceTencentCloudMonitorPolicyBindingObject(),
			"tencentcloud_monitor_binding_receiver":                                                 monitor.ResourceTencentCloudMonitorBindingAlarmReceiver(),
			"tencentcloud_monitor_alarm_policy":                                                     monitor.ResourceTencentCloudMonitorAlarmPolicy(),
			"tencentcloud_monitor_alarm_notice":                                                     monitor.ResourceTencentCloudMonitorAlarmNotice(),
			"tencentcloud_monitor_alarm_policy_set_default":                                         monitor.ResourceTencentCloudMonitorAlarmPolicySetDefault(),
			"tencentcloud_monitor_tmp_instance":                                                     tmp.ResourceTencentCloudMonitorTmpInstance(),
			"tencentcloud_monitor_tmp_cvm_agent":                                                    tmp.ResourceTencentCloudMonitorTmpCvmAgent(),
			"tencentcloud_monitor_tmp_scrape_job":                                                   tmp.ResourceTencentCloudMonitorTmpScrapeJob(),
			"tencentcloud_monitor_tmp_exporter_integration":                                         tmp.ResourceTencentCloudMonitorTmpExporterIntegration(),
			"tencentcloud_monitor_tmp_alert_rule":                                                   tmp.ResourceTencentCloudMonitorTmpAlertRule(),
			"tencentcloud_monitor_tmp_recording_rule":                                               tmp.ResourceTencentCloudMonitorTmpRecordingRule(),
			"tencentcloud_monitor_tmp_multiple_writes":                                              tmp.ResourceTencentCloudMonitorTmpMultipleWrites(),
			"tencentcloud_monitor_tmp_alert_group":                                                  tmp.ResourceTencentCloudMonitorTmpAlertGroup(),
			"tencentcloud_monitor_tmp_tke_template":                                                 tmp.ResourceTencentCloudMonitorTmpTkeTemplate(),
			"tencentcloud_monitor_tmp_tke_template_attachment":                                      tmp.ResourceTencentCloudMonitorTmpTkeTemplateAttachment(),
			"tencentcloud_monitor_tmp_tke_alert_policy":                                             tmp.ResourceTencentCloudMonitorTmpTkeAlertPolicy(),
			"tencentcloud_monitor_tmp_tke_basic_config":                                             tmp.ResourceTencentCloudMonitorTmpTkeBasicConfig(),
			"tencentcloud_monitor_tmp_tke_cluster_agent":                                            tmp.ResourceTencentCloudMonitorTmpTkeClusterAgent(),
			"tencentcloud_monitor_tmp_tke_config":                                                   tmp.ResourceTencentCloudMonitorTmpTkeConfig(),
			"tencentcloud_monitor_tmp_tke_record_rule_yaml":                                         tmp.ResourceTencentCloudMonitorTmpTkeRecordRuleYaml(),
			"tencentcloud_monitor_tmp_tke_global_notification":                                      tmp.ResourceTencentCloudMonitorTmpTkeGlobalNotification(),
			"tencentcloud_monitor_tmp_manage_grafana_attachment":                                    tmp.ResourceTencentCloudMonitorTmpManageGrafanaAttachment(),
			"tencentcloud_monitor_grafana_instance":                                                 tcmg.ResourceTencentCloudMonitorGrafanaInstance(),
			"tencentcloud_monitor_grafana_integration":                                              tcmg.ResourceTencentCloudMonitorGrafanaIntegration(),
			"tencentcloud_monitor_grafana_notification_channel":                                     tcmg.ResourceTencentCloudMonitorGrafanaNotificationChannel(),
			"tencentcloud_monitor_grafana_plugin":                                                   tcmg.ResourceTencentCloudMonitorGrafanaPlugin(),
			"tencentcloud_monitor_grafana_sso_account":                                              tcmg.ResourceTencentCloudMonitorGrafanaSsoAccount(),
			"tencentcloud_monitor_tmp_grafana_config":                                               tcmg.ResourceTencentCloudMonitorTmpGrafanaConfig(),
			"tencentcloud_monitor_grafana_dns_config":                                               tcmg.ResourceTencentCloudMonitorGrafanaDnsConfig(),
			"tencentcloud_monitor_grafana_env_config":                                               tcmg.ResourceTencentCloudMonitorGrafanaEnvConfig(),
			"tencentcloud_monitor_grafana_whitelist_config":                                         tcmg.ResourceTencentCloudMonitorGrafanaWhitelistConfig(),
			"tencentcloud_monitor_grafana_sso_cam_config":                                           tcmg.ResourceTencentCloudMonitorGrafanaSsoCamConfig(),
			"tencentcloud_monitor_grafana_sso_config":                                               tcmg.ResourceTencentCloudMonitorGrafanaSsoConfig(),
			"tencentcloud_monitor_grafana_version_upgrade":                                          tcmg.ResourceTencentCloudMonitorGrafanaVersionUpgrade(),
			"tencentcloud_mongodb_standby_instance":                                                 mongodb.ResourceTencentCloudMongodbStandbyInstance(),
			"tencentcloud_elasticsearch_instance":                                                   es.ResourceTencentCloudElasticsearchInstance(),
			"tencentcloud_elasticsearch_security_group":                                             es.ResourceTencentCloudElasticsearchSecurityGroup(),
			"tencentcloud_elasticsearch_logstash":                                                   es.ResourceTencentCloudElasticsearchLogstash(),
			"tencentcloud_elasticsearch_logstash_pipeline":                                          es.ResourceTencentCloudElasticsearchLogstashPipeline(),
			"tencentcloud_elasticsearch_restart_logstash_instance_operation":                        es.ResourceTencentCloudElasticsearchRestartLogstashInstanceOperation(),
			"tencentcloud_elasticsearch_start_logstash_pipeline_operation":                          es.ResourceTencentCloudElasticsearchStartLogstashPipelineOperation(),
			"tencentcloud_elasticsearch_stop_logstash_pipeline_operation":                           es.ResourceTencentCloudElasticsearchStopLogstashPipelineOperation(),
			"tencentcloud_elasticsearch_index":                                                      es.ResourceTencentCloudElasticsearchIndex(),
			"tencentcloud_elasticsearch_restart_instance_operation":                                 es.ResourceTencentCloudElasticsearchRestartInstanceOperation(),
			"tencentcloud_elasticsearch_restart_kibana_operation":                                   es.ResourceTencentCloudElasticsearchRestartKibanaOperation(),
			"tencentcloud_elasticsearch_restart_nodes_operation":                                    es.ResourceTencentCloudElasticsearchRestartNodesOperation(),
			"tencentcloud_elasticsearch_diagnose":                                                   es.ResourceTencentCloudElasticsearchDiagnose(),
			"tencentcloud_elasticsearch_diagnose_instance":                                          es.ResourceTencentCloudElasticsearchDiagnoseInstance(),
			"tencentcloud_elasticsearch_update_plugins_operation":                                   es.ResourceTencentCloudElasticsearchUpdatePluginsOperation(),
			"tencentcloud_postgresql_instance":                                                      postgresql.ResourceTencentCloudPostgresqlInstance(),
			"tencentcloud_postgresql_readonly_instance":                                             postgresql.ResourceTencentCloudPostgresqlReadonlyInstance(),
			"tencentcloud_postgresql_readonly_group":                                                postgresql.ResourceTencentCloudPostgresqlReadonlyGroup(),
			"tencentcloud_postgresql_readonly_attachment":                                           postgresql.ResourceTencentCloudPostgresqlReadonlyAttachment(),
			"tencentcloud_postgresql_parameter_template":                                            postgresql.ResourceTencentCloudPostgresqlParameterTemplate(),
			"tencentcloud_postgresql_base_backup":                                                   postgresql.ResourceTencentCloudPostgresqlBaseBackup(),
			"tencentcloud_postgresql_backup_plan_config":                                            postgresql.ResourceTencentCloudPostgresqlBackupPlanConfig(),
			"tencentcloud_postgresql_security_group_config":                                         postgresql.ResourceTencentCloudPostgresqlSecurityGroupConfig(),
			"tencentcloud_postgresql_backup_download_restriction_config":                            postgresql.ResourceTencentCloudPostgresqlBackupDownloadRestrictionConfig(),
			"tencentcloud_postgresql_restart_db_instance_operation":                                 postgresql.ResourceTencentCloudPostgresqlRestartDbInstanceOperation(),
			"tencentcloud_postgresql_renew_db_instance_operation":                                   postgresql.ResourceTencentCloudPostgresqlRenewDbInstanceOperation(),
			"tencentcloud_postgresql_isolate_db_instance_operation":                                 postgresql.ResourceTencentCloudPostgresqlIsolateDbInstanceOperation(),
			"tencentcloud_postgresql_disisolate_db_instance_operation":                              postgresql.ResourceTencentCloudPostgresqlDisisolateDbInstanceOperation(),
			"tencentcloud_postgresql_rebalance_readonly_group_operation":                            postgresql.ResourceTencentCloudPostgresqlRebalanceReadonlyGroupOperation(),
			"tencentcloud_postgresql_delete_log_backup_operation":                                   postgresql.ResourceTencentCloudPostgresqlDeleteLogBackupOperation(),
			"tencentcloud_postgresql_modify_account_remark_operation":                               postgresql.ResourceTencentCloudPostgresqlModifyAccountRemarkOperation(),
			"tencentcloud_postgresql_modify_switch_time_period_operation":                           postgresql.ResourceTencentCloudPostgresqlModifySwitchTimePeriodOperation(),
			"tencentcloud_postgresql_instance_ha_config":                                            postgresql.ResourceTencentCloudPostgresqlInstanceHAConfig(),
			"tencentcloud_postgresql_account":                                                       postgresql.ResourceTencentCloudPostgresqlAccount(),
			"tencentcloud_postgresql_account_privileges_operation":                                  postgresql.ResourceTencentCloudPostgresqlAccountPrivilegesOperation(),
			"tencentcloud_postgresql_apply_parameter_template_operation":                            postgresql.ResourceTencentCloudPostgresqlApplyParameterTemplateOperation(),
			"tencentcloud_postgresql_clone_db_instance":                                             postgresql.ResourceTencentCloudPostgresqlCloneDbInstance(),
			"tencentcloud_sqlserver_instance":                                                       sqlserver.ResourceTencentCloudSqlserverInstance(),
			"tencentcloud_sqlserver_db":                                                             sqlserver.ResourceTencentCloudSqlserverDB(),
			"tencentcloud_sqlserver_account":                                                        sqlserver.ResourceTencentCloudSqlserverAccount(),
			"tencentcloud_sqlserver_account_db_attachment":                                          sqlserver.ResourceTencentCloudSqlserverAccountDBAttachment(),
			"tencentcloud_sqlserver_readonly_instance":                                              sqlserver.ResourceTencentCloudSqlserverReadonlyInstance(),
			"tencentcloud_sqlserver_migration":                                                      sqlserver.ResourceTencentCloudSqlserverMigration(),
			"tencentcloud_sqlserver_config_backup_strategy":                                         sqlserver.ResourceTencentCloudSqlserverConfigBackupStrategy(),
			"tencentcloud_sqlserver_general_backup":                                                 sqlserver.ResourceTencentCloudSqlserverGeneralBackup(),
			"tencentcloud_sqlserver_general_clone":                                                  sqlserver.ResourceTencentCloudSqlserverGeneralClone(),
			"tencentcloud_sqlserver_full_backup_migration":                                          sqlserver.ResourceTencentCloudSqlserverFullBackupMigration(),
			"tencentcloud_sqlserver_incre_backup_migration":                                         sqlserver.ResourceTencentCloudSqlserverIncreBackupMigration(),
			"tencentcloud_sqlserver_business_intelligence_file":                                     sqlserver.ResourceTencentCloudSqlserverBusinessIntelligenceFile(),
			"tencentcloud_sqlserver_business_intelligence_instance":                                 sqlserver.ResourceTencentCloudSqlserverBusinessIntelligenceInstance(),
			"tencentcloud_sqlserver_general_communication":                                          sqlserver.ResourceTencentCloudSqlserverGeneralCommunication(),
			"tencentcloud_sqlserver_general_cloud_instance":                                         sqlserver.ResourceTencentCloudSqlserverGeneralCloudInstance(),
			"tencentcloud_sqlserver_complete_expansion":                                             sqlserver.ResourceTencentCloudSqlserverCompleteExpansion(),
			"tencentcloud_sqlserver_config_database_cdc":                                            sqlserver.ResourceTencentCloudSqlserverConfigDatabaseCDC(),
			"tencentcloud_sqlserver_config_database_ct":                                             sqlserver.ResourceTencentCloudSqlserverConfigDatabaseCT(),
			"tencentcloud_sqlserver_config_database_mdf":                                            sqlserver.ResourceTencentCloudSqlserverConfigDatabaseMdf(),
			"tencentcloud_sqlserver_config_instance_param":                                          sqlserver.ResourceTencentCloudSqlserverConfigInstanceParam(),
			"tencentcloud_sqlserver_config_instance_ro_group":                                       sqlserver.ResourceTencentCloudSqlserverConfigInstanceRoGroup(),
			"tencentcloud_sqlserver_config_instance_security_groups":                                sqlserver.ResourceTencentCloudSqlserverConfigInstanceSecurityGroups(),
			"tencentcloud_sqlserver_renew_db_instance":                                              sqlserver.ResourceTencentCloudSqlserverRenewDBInstance(),
			"tencentcloud_sqlserver_renew_postpaid_db_instance":                                     sqlserver.ResourceTencentCloudSqlserverRenewPostpaidDBInstance(),
			"tencentcloud_sqlserver_restart_db_instance":                                            sqlserver.ResourceTencentCloudSqlserverRestartDBInstance(),
			"tencentcloud_sqlserver_config_terminate_db_instance":                                   sqlserver.ResourceTencentCloudSqlserverConfigTerminateDBInstance(),
			"tencentcloud_sqlserver_restore_instance":                                               sqlserver.ResourceTencentCloudSqlserverRestoreInstance(),
			"tencentcloud_sqlserver_rollback_instance":                                              sqlserver.ResourceTencentCloudSqlserverRollbackInstance(),
			"tencentcloud_sqlserver_start_backup_full_migration":                                    sqlserver.ResourceTencentCloudSqlserverStartBackupFullMigration(),
			"tencentcloud_sqlserver_start_backup_incremental_migration":                             sqlserver.ResourceTencentCloudSqlserverStartBackupIncrementalMigration(),
			"tencentcloud_sqlserver_start_xevent":                                                   sqlserver.ResourceTencentCloudSqlserverStartXevent(),
			"tencentcloud_ckafka_instance":                                                          ckafka.ResourceTencentCloudCkafkaInstance(),
			"tencentcloud_ckafka_user":                                                              ckafka.ResourceTencentCloudCkafkaUser(),
			"tencentcloud_ckafka_acl":                                                               ckafka.ResourceTencentCloudCkafkaAcl(),
			"tencentcloud_ckafka_topic":                                                             ckafka.ResourceTencentCloudCkafkaTopic(),
			"tencentcloud_ckafka_datahub_topic":                                                     ckafka.ResourceTencentCloudCkafkaDatahubTopic(),
			"tencentcloud_ckafka_connect_resource":                                                  ckafka.ResourceTencentCloudCkafkaConnectResource(),
			"tencentcloud_ckafka_renew_instance":                                                    ckafka.ResourceTencentCloudCkafkaRenewInstance(),
			"tencentcloud_ckafka_acl_rule":                                                          ckafka.ResourceTencentCloudCkafkaAclRule(),
			"tencentcloud_ckafka_consumer_group":                                                    ckafka.ResourceTencentCloudCkafkaConsumerGroup(),
			"tencentcloud_ckafka_consumer_group_modify_offset":                                      ckafka.ResourceTencentCloudCkafkaConsumerGroupModifyOffset(),
			"tencentcloud_ckafka_datahub_task":                                                      ckafka.ResourceTencentCloudCkafkaDatahubTask(),
			"tencentcloud_ckafka_route":                                                             ckafka.ResourceTencentCloudCkafkaRoute(),
			"tencentcloud_audit_track":                                                              audit.ResourceTencentCloudAuditTrack(),
			"tencentcloud_events_audit_track":                                                       audit.ResourceTencentCloudEventsAuditTrack(),
			"tencentcloud_image":                                                                    cvm.ResourceTencentCloudImage(),
			"tencentcloud_cynosdb_proxy":                                                            cynosdb.ResourceTencentCloudCynosdbProxy(),
			"tencentcloud_cynosdb_reload_proxy_node":                                                cynosdb.ResourceTencentCloudCynosdbReloadProxyNode(),
			"tencentcloud_cynosdb_cluster_resource_packages_attachment":                             cynosdb.ResourceTencentCloudCynosdbClusterResourcePackagesAttachment(),
			"tencentcloud_cynosdb_cluster":                                                          cynosdb.ResourceTencentCloudCynosdbCluster(),
			"tencentcloud_cynosdb_readonly_instance":                                                cynosdb.ResourceTencentCloudCynosdbReadonlyInstance(),
			"tencentcloud_cynosdb_cluster_password_complexity":                                      cynosdb.ResourceTencentCloudCynosdbClusterPasswordComplexity(),
			"tencentcloud_cynosdb_export_instance_error_logs":                                       cynosdb.ResourceTencentCloudCynosdbExportInstanceErrorLogs(),
			"tencentcloud_cynosdb_export_instance_slow_queries":                                     cynosdb.ResourceTencentCloudCynosdbExportInstanceSlowQueries(),
			"tencentcloud_cynosdb_account_privileges":                                               cynosdb.ResourceTencentCloudCynosdbAccountPrivileges(),
			"tencentcloud_cynosdb_account":                                                          cynosdb.ResourceTencentCloudCynosdbAccount(),
			"tencentcloud_cynosdb_binlog_save_days":                                                 cynosdb.ResourceTencentCloudCynosdbBinlogSaveDays(),
			"tencentcloud_cynosdb_cluster_databases":                                                cynosdb.ResourceTencentCloudCynosdbClusterDatabases(),
			"tencentcloud_cynosdb_instance_param":                                                   cynosdb.ResourceTencentCloudCynosdbInstanceParam(),
			"tencentcloud_cynosdb_isolate_instance":                                                 cynosdb.ResourceTencentCloudCynosdbIsolateInstance(),
			"tencentcloud_cynosdb_param_template":                                                   cynosdb.ResourceTencentCloudCynosdbParamTemplate(),
			"tencentcloud_cynosdb_resource_package":                                                 cynosdb.ResourceTencentCloudCynosdbResourcePackage(),
			"tencentcloud_cynosdb_restart_instance":                                                 cynosdb.ResourceTencentCloudCynosdbRestartInstance(),
			"tencentcloud_cynosdb_roll_back_cluster":                                                cynosdb.ResourceTencentCloudCynosdbRollBackCluster(),
			"tencentcloud_cynosdb_wan":                                                              cynosdb.ResourceTencentCloudCynosdbWan(),
			"tencentcloud_cynosdb_cluster_slave_zone":                                               cynosdb.ResourceTencentCloudCynosdbClusterSlaveZone(),
			"tencentcloud_cynosdb_read_only_instance_exclusive_access":                              cynosdb.ResourceTencentCloudCynosdbReadOnlyInstanceExclusiveAccess(),
			"tencentcloud_cynosdb_proxy_end_point":                                                  cynosdb.ResourceTencentCloudCynosdbProxyEndPoint(),
			"tencentcloud_cynosdb_upgrade_proxy_version":                                            cynosdb.ResourceTencentCloudCynosdbUpgradeProxyVersion(),
			"tencentcloud_vod_adaptive_dynamic_streaming_template":                                  vod.ResourceTencentCloudVodAdaptiveDynamicStreamingTemplate(),
			"tencentcloud_vod_image_sprite_template":                                                vod.ResourceTencentCloudVodImageSpriteTemplate(),
			"tencentcloud_vod_procedure_template":                                                   vod.ResourceTencentCloudVodProcedureTemplate(),
			"tencentcloud_vod_snapshot_by_time_offset_template":                                     vod.ResourceTencentCloudVodSnapshotByTimeOffsetTemplate(),
			"tencentcloud_vod_super_player_config":                                                  vod.ResourceTencentCloudVodSuperPlayerConfig(),
			"tencentcloud_vod_sub_application":                                                      vod.ResourceTencentCloudVodSubApplication(),
			"tencentcloud_vod_sample_snapshot_template":                                             vod.ResourceTencentCloudVodSampleSnapshotTemplate(),
			"tencentcloud_vod_transcode_template":                                                   vod.ResourceTencentCloudVodTranscodeTemplate(),
			"tencentcloud_vod_watermark_template":                                                   vod.ResourceTencentCloudVodWatermarkTemplate(),
			"tencentcloud_vod_event_config":                                                         vod.ResourceTencentCloudVodEventConfig(),
			"tencentcloud_sqlserver_publish_subscribe":                                              sqlserver.ResourceTencentCloudSqlserverPublishSubscribe(),
			"tencentcloud_api_gateway_usage_plan":                                                   apigateway.ResourceTencentCloudAPIGatewayUsagePlan(),
			"tencentcloud_api_gateway_usage_plan_attachment":                                        apigateway.ResourceTencentCloudAPIGatewayUsagePlanAttachment(),
			"tencentcloud_api_gateway_api":                                                          apigateway.ResourceTencentCloudAPIGatewayAPI(),
			"tencentcloud_api_gateway_service":                                                      apigateway.ResourceTencentCloudAPIGatewayService(),
			"tencentcloud_api_gateway_custom_domain":                                                apigateway.ResourceTencentCloudAPIGatewayCustomDomain(),
			"tencentcloud_api_gateway_ip_strategy":                                                  apigateway.ResourceTencentCloudAPIGatewayIPStrategy(),
			"tencentcloud_api_gateway_strategy_attachment":                                          apigateway.ResourceTencentCloudAPIGatewayStrategyAttachment(),
			"tencentcloud_api_gateway_api_key":                                                      apigateway.ResourceTencentCloudAPIGatewayAPIKey(),
			"tencentcloud_api_gateway_api_key_attachment":                                           apigateway.ResourceTencentCloudAPIGatewayAPIKeyAttachment(),
			"tencentcloud_api_gateway_service_release":                                              apigateway.ResourceTencentCloudAPIGatewayServiceRelease(),
			"tencentcloud_api_gateway_plugin":                                                       apigateway.ResourceTencentCloudAPIGatewayPlugin(),
			"tencentcloud_api_gateway_plugin_attachment":                                            apigateway.ResourceTencentCloudAPIGatewayPluginAttachment(),
			"tencentcloud_api_gateway_upstream":                                                     apigateway.ResourceTencentCloudAPIGatewayUpstream(),
			"tencentcloud_api_gateway_api_app_attachment":                                           apigateway.ResourceTencentCloudAPIGatewayApiAppAttachment(),
			"tencentcloud_api_gateway_update_service":                                               apigateway.ResourceTencentCloudAPIGatewayUpdateService(),
			"tencentcloud_sqlserver_basic_instance":                                                 sqlserver.ResourceTencentCloudSqlserverBasicInstance(),
			"tencentcloud_sqlserver_instance_tde":                                                   sqlserver.ResourceTencentCloudSqlserverInstanceTDE(),
			"tencentcloud_sqlserver_database_tde":                                                   sqlserver.ResourceTencentCloudSqlserverDatabaseTDE(),
			"tencentcloud_sqlserver_general_cloud_ro_instance":                                      sqlserver.ResourceTencentCloudSqlserverGeneralCloudRoInstance(),
			"tencentcloud_sqlserver_instance_ssl":                                                   sqlserver.ResourceTencentCloudSqlserverInstanceSsl(),
			"tencentcloud_tcr_instance":                                                             tcr.ResourceTencentCloudTcrInstance(),
			"tencentcloud_tcr_namespace":                                                            tcr.ResourceTencentCloudTcrNamespace(),
			"tencentcloud_tcr_repository":                                                           tcr.ResourceTencentCloudTcrRepository(),
			"tencentcloud_tcr_token":                                                                tcr.ResourceTencentCloudTcrToken(),
			"tencentcloud_tcr_vpc_attachment":                                                       tcr.ResourceTencentCloudTcrVpcAttachment(),
			"tencentcloud_tcr_tag_retention_rule":                                                   tcr.ResourceTencentCloudTcrTagRetentionRule(),
			"tencentcloud_tcr_webhook_trigger":                                                      tcr.ResourceTencentCloudTcrWebhookTrigger(),
			"tencentcloud_tcr_manage_replication_operation":                                         tcr.ResourceTencentCloudTcrManageReplicationOperation(),
			"tencentcloud_tcr_customized_domain":                                                    tcr.ResourceTencentCloudTcrCustomizedDomain(),
			"tencentcloud_tcr_immutable_tag_rule":                                                   tcr.ResourceTencentCloudTcrImmutableTagRule(),
			"tencentcloud_tcr_delete_image_operation":                                               tcr.ResourceTencentCloudTcrDeleteImageOperation(),
			"tencentcloud_tcr_create_image_signature_operation":                                     tcr.ResourceTencentCloudTcrCreateImageSignatureOperation(),
			"tencentcloud_tcr_tag_retention_execution_config":                                       tcr.ResourceTencentCloudTcrTagRetentionExecutionConfig(),
			"tencentcloud_tcr_service_account":                                                      tcr.ResourceTencentCloudTcrServiceAccount(),
			"tencentcloud_tdmq_instance":                                                            tpulsar.ResourceTencentCloudTdmqInstance(),
			"tencentcloud_tdmq_namespace":                                                           tpulsar.ResourceTencentCloudTdmqNamespace(),
			"tencentcloud_tdmq_topic":                                                               tpulsar.ResourceTencentCloudTdmqTopic(),
			"tencentcloud_tdmq_topic_with_full_id":                                                  tpulsar.ResourceTencentCloudTdmqTopicWithFullId(),
			"tencentcloud_tdmq_role":                                                                tpulsar.ResourceTencentCloudTdmqRole(),
			"tencentcloud_tdmq_namespace_role_attachment":                                           tpulsar.ResourceTencentCloudTdmqNamespaceRoleAttachment(),
			"tencentcloud_tdmq_rabbitmq_user":                                                       trabbit.ResourceTencentCloudTdmqRabbitmqUser(),
			"tencentcloud_tdmq_rabbitmq_virtual_host":                                               trabbit.ResourceTencentCloudTdmqRabbitmqVirtualHost(),
			"tencentcloud_tdmq_rabbitmq_vip_instance":                                               trabbit.ResourceTencentCloudTdmqRabbitmqVipInstance(),
			"tencentcloud_tdmq_send_rocketmq_message":                                               trocket.ResourceTencentCloudTdmqSendRocketmqMessage(),
			"tencentcloud_tdmq_professional_cluster":                                                tpulsar.ResourceTencentCloudTdmqProfessionalCluster(),
			"tencentcloud_tdmq_subscription":                                                        tpulsar.ResourceTencentCloudTdmqSubscription(),
			"tencentcloud_cos_bucket_policy":                                                        cos.ResourceTencentCloudCosBucketPolicy(),
			"tencentcloud_cos_bucket_domain_certificate_attachment":                                 cos.ResourceTencentCloudCosBucketDomainCertificateAttachment(),
			"tencentcloud_cos_bucket_inventory":                                                     cos.ResourceTencentCloudCosBucketInventory(),
			"tencentcloud_cos_batch":                                                                cos.ResourceTencentCloudCosBatch(),
			"tencentcloud_cos_object_abort_multipart_upload_operation":                              cos.ResourceTencentCloudCosObjectAbortMultipartUploadOperation(),
			"tencentcloud_cos_object_copy_operation":                                                cos.ResourceTencentCloudCosObjectCopyOperation(),
			"tencentcloud_cos_object_restore_operation":                                             cos.ResourceTencentCloudCosObjectRestoreOperation(),
			"tencentcloud_cos_bucket_generate_inventory_immediately_operation":                      cos.ResourceTencentCloudCosBucketGenerateInventoryImmediatelyOperation(),
			"tencentcloud_cos_object_download_operation":                                            cos.ResourceTencentCloudCosObjectDownloadOperation(),
			"tencentcloud_address_template":                                                         vpc.ResourceTencentCloudAddressTemplate(),
			"tencentcloud_address_template_group":                                                   vpc.ResourceTencentCloudAddressTemplateGroup(),
			"tencentcloud_protocol_template":                                                        vpc.ResourceTencentCloudProtocolTemplate(),
			"tencentcloud_protocol_template_group":                                                  vpc.ResourceTencentCloudProtocolTemplateGroup(),
			"tencentcloud_kms_key":                                                                  kms.ResourceTencentCloudKmsKey(),
			"tencentcloud_kms_external_key":                                                         kms.ResourceTencentCloudKmsExternalKey(),
			"tencentcloud_kms_white_box_key":                                                        kms.ResourceTencentCloudKmsWhiteBoxKey(),
			"tencentcloud_kms_cloud_resource_attachment":                                            kms.ResourceTencentCloudKmsCloudResourceAttachment(),
			"tencentcloud_kms_overwrite_white_box_device_fingerprints":                              kms.ResourceTencentCloudKmsOverwriteWhiteBoxDeviceFingerprints(),
			"tencentcloud_ssm_secret":                                                               ssm.ResourceTencentCloudSsmSecret(),
			"tencentcloud_ssm_ssh_key_pair_secret":                                                  ssm.ResourceTencentCloudSsmSshKeyPairSecret(),
			"tencentcloud_ssm_product_secret":                                                       ssm.ResourceTencentCloudSsmProductSecret(),
			"tencentcloud_ssm_secret_version":                                                       ssm.ResourceTencentCloudSsmSecretVersion(),
			"tencentcloud_ssm_rotate_product_secret":                                                ssm.ResourceTencentCloudSsmRotateProductSecret(),
			"tencentcloud_cdh_instance":                                                             cdh.ResourceTencentCloudCdhInstance(),
			"tencentcloud_dnspod_domain_instance":                                                   dnspod.ResourceTencentCloudDnspodDomainInstance(),
			"tencentcloud_dnspod_domain_alias":                                                      dnspod.ResourceTencentCloudDnspodDomainAlias(),
			"tencentcloud_dnspod_record":                                                            dnspod.ResourceTencentCloudDnspodRecord(),
			"tencentcloud_dnspod_record_group":                                                      dnspod.ResourceTencentCloudDnspodRecordGroup(),
			"tencentcloud_dnspod_modify_domain_owner_operation":                                     dnspod.ResourceTencentCloudDnspodModifyDomainOwnerOperation(),
			"tencentcloud_dnspod_modify_record_group_operation":                                     dnspod.ResourceTencentCloudDnspodModifyRecordGroupOperation(),
			"tencentcloud_dnspod_download_snapshot_operation":                                       dnspod.ResourceTencentCloudDnspodDownloadSnapshotOperation(),
			"tencentcloud_dnspod_custom_line":                                                       dnspod.ResourceTencentCloudDnspodCustomLine(),
			"tencentcloud_dnspod_snapshot_config":                                                   dnspod.ResourceTencentCloudDnspodSnapshotConfig(),
			"tencentcloud_dnspod_domain_lock":                                                       dnspod.ResourceTencentCloudDnspodDomainLock(),
			"tencentcloud_subdomain_validate_txt_value_operation":                                   dnspod.ResourceTencentCloudSubdomainValidateTxtValueOperation(),
			"tencentcloud_private_dns_zone":                                                         privatedns.ResourceTencentCloudPrivateDnsZone(),
			"tencentcloud_private_dns_record":                                                       privatedns.ResourceTencentCloudPrivateDnsRecord(),
			"tencentcloud_private_dns_zone_vpc_attachment":                                          privatedns.ResourceTencentCloudPrivateDnsZoneVpcAttachment(),
			"tencentcloud_subscribe_private_zone_service":                                           privatedns.ResourceTencentCloudSubscribePrivateZoneService(),
			"tencentcloud_private_dns_forward_rule":                                                 privatedns.ResourceTencentCloudPrivateDnsForwardRule(),
			"tencentcloud_private_dns_end_point":                                                    privatedns.ResourceTencentCloudPrivateDnsEndPoint(),
			"tencentcloud_cls_logset":                                                               cls.ResourceTencentCloudClsLogset(),
			"tencentcloud_cls_topic":                                                                cls.ResourceTencentCloudClsTopic(),
			"tencentcloud_cls_config":                                                               cls.ResourceTencentCloudClsConfig(),
			"tencentcloud_cls_config_extra":                                                         cls.ResourceTencentCloudClsConfigExtra(),
			"tencentcloud_cls_config_attachment":                                                    cls.ResourceTencentCloudClsConfigAttachment(),
			"tencentcloud_cls_machine_group":                                                        cls.ResourceTencentCloudClsMachineGroup(),
			"tencentcloud_cls_cos_shipper":                                                          cls.ResourceTencentCloudClsCosShipper(),
			"tencentcloud_cls_index":                                                                cls.ResourceTencentCloudClsIndex(),
			"tencentcloud_cls_alarm":                                                                cls.ResourceTencentCloudClsAlarm(),
			"tencentcloud_cls_alarm_notice":                                                         cls.ResourceTencentCloudClsAlarmNotice(),
			"tencentcloud_cls_ckafka_consumer":                                                      cls.ResourceTencentCloudClsCkafkaConsumer(),
			"tencentcloud_cls_cos_recharge":                                                         cls.ResourceTencentCloudClsCosRecharge(),
			"tencentcloud_cls_export":                                                               cls.ResourceTencentCloudClsExport(),
			"tencentcloud_cls_data_transform":                                                       cls.ResourceTencentCloudClsDataTransform(),
			"tencentcloud_cls_cloud_product_log_task":                                               cls.ResourceTencentCloudClsCloudProductLogTask(),
			"tencentcloud_cls_notice_content":                                                       cls.ResourceTencentCloudClsNoticeContent(),
			"tencentcloud_lighthouse_instance":                                                      lighthouse.ResourceTencentCloudLighthouseInstance(),
			"tencentcloud_lighthouse_firewall_template":                                             lighthouse.ResourceTencentCloudLighthouseFirewallTemplate(),
			"tencentcloud_tem_environment":                                                          tem.ResourceTencentCloudTemEnvironment(),
			"tencentcloud_tem_application":                                                          tem.ResourceTencentCloudTemApplication(),
			"tencentcloud_tem_workload":                                                             tem.ResourceTencentCloudTemWorkload(),
			"tencentcloud_tem_app_config":                                                           tem.ResourceTencentCloudTemAppConfig(),
			"tencentcloud_tem_log_config":                                                           tem.ResourceTencentCloudTemLogConfig(),
			"tencentcloud_tem_scale_rule":                                                           tem.ResourceTencentCloudTemScaleRule(),
			"tencentcloud_tem_gateway":                                                              tem.ResourceTencentCloudTemGateway(),
			"tencentcloud_tem_application_service":                                                  tem.ResourceTencentCloudTemApplicationService(),
			"tencentcloud_teo_zone":                                                                 teo.ResourceTencentCloudTeoZone(),
			"tencentcloud_teo_zone_setting":                                                         teo.ResourceTencentCloudTeoZoneSetting(),
			"tencentcloud_teo_origin_group":                                                         teo.ResourceTencentCloudTeoOriginGroup(),
			"tencentcloud_teo_l4_proxy":                                                             teo.ResourceTencentCloudTeoL4Proxy(),
			"tencentcloud_teo_rule_engine":                                                          teo.ResourceTencentCloudTeoRuleEngine(),
			"tencentcloud_teo_ownership_verify":                                                     teo.ResourceTencentCloudTeoOwnershipVerify(),
			"tencentcloud_teo_certificate_config":                                                   teo.ResourceTencentCloudTeoCertificateConfig(),
			"tencentcloud_teo_acceleration_domain":                                                  teo.ResourceTencentCloudTeoAccelerationDomain(),
			"tencentcloud_teo_application_proxy":                                                    teo.ResourceTencentCloudTeoApplicationProxy(),
			"tencentcloud_teo_application_proxy_rule":                                               teo.ResourceTencentCloudTeoApplicationProxyRule(),
			"tencentcloud_teo_realtime_log_delivery":                                                teo.ResourceTencentCloudTeoRealtimeLogDelivery(),
			"tencentcloud_teo_security_ip_group":                                                    teo.ResourceTencentCloudTeoSecurityIpGroup(),
			"tencentcloud_teo_function":                                                             teo.ResourceTencentCloudTeoFunction(),
			"tencentcloud_teo_function_rule":                                                        teo.ResourceTencentCloudTeoFunctionRule(),
			"tencentcloud_teo_function_rule_priority":                                               teo.ResourceTencentCloudTeoFunctionRulePriority(),
			"tencentcloud_teo_function_runtime_environment":                                         teo.ResourceTencentCloudTeoFunctionRuntimeEnvironment(),
			"tencentcloud_tcm_mesh":                                                                 tcm.ResourceTencentCloudTcmMesh(),
			"tencentcloud_tcm_cluster_attachment":                                                   tcm.ResourceTencentCloudTcmClusterAttachment(),
			"tencentcloud_tcm_prometheus_attachment":                                                tcm.ResourceTencentCloudTcmPrometheusAttachment(),
			"tencentcloud_tcm_tracing_config":                                                       tcm.ResourceTencentCloudTcmTracingConfig(),
			"tencentcloud_tcm_access_log_config":                                                    tcm.ResourceTencentCloudTcmAccessLogConfig(),
			"tencentcloud_ses_domain":                                                               ses.ResourceTencentCloudSesDomain(),
			"tencentcloud_ses_template":                                                             ses.ResourceTencentCloudSesTemplate(),
			"tencentcloud_ses_email_address":                                                        ses.ResourceTencentCloudSesEmailAddress(),
			"tencentcloud_ses_receiver":                                                             ses.ResourceTencentCloudSesReceiver(),
			"tencentcloud_ses_send_email":                                                           ses.ResourceTencentCloudSesSendEmail(),
			"tencentcloud_ses_batch_send_email":                                                     ses.ResourceTencentCloudSesBatchSendEmail(),
			"tencentcloud_ses_verify_domain":                                                        ses.ResourceTencentCloudSesVerifyDomain(),
			"tencentcloud_ses_black_list_delete":                                                    ses.ResourceTencentCloudSesBlackListDelete(),
			"tencentcloud_sms_sign":                                                                 sms.ResourceTencentCloudSmsSign(),
			"tencentcloud_sms_template":                                                             sms.ResourceTencentCloudSmsTemplate(),
			"tencentcloud_dcdb_account":                                                             dcdb.ResourceTencentCloudDcdbAccount(),
			"tencentcloud_dcdb_hourdb_instance":                                                     dcdb.ResourceTencentCloudDcdbHourdbInstance(),
			"tencentcloud_dcdb_security_group_attachment":                                           dcdb.ResourceTencentCloudDcdbSecurityGroupAttachment(),
			"tencentcloud_dcdb_db_instance":                                                         dcdb.ResourceTencentCloudDcdbDbInstance(),
			"tencentcloud_dcdb_account_privileges":                                                  dcdb.ResourceTencentCloudDcdbAccountPrivileges(),
			"tencentcloud_dcdb_db_parameters":                                                       dcdb.ResourceTencentCloudDcdbDbParameters(),
			"tencentcloud_dcdb_encrypt_attributes_config":                                           dcdb.ResourceTencentCloudDcdbEncryptAttributesConfig(),
			"tencentcloud_dcdb_db_sync_mode_config":                                                 dcdb.ResourceTencentCloudDcdbDbSyncModeConfig(),
			"tencentcloud_dcdb_instance_config":                                                     dcdb.ResourceTencentCloudDcdbInstanceConfig(),
			"tencentcloud_dcdb_activate_hour_instance_operation":                                    dcdb.ResourceTencentCloudDcdbActivateHourInstanceOperation(),
			"tencentcloud_dcdb_isolate_hour_instance_operation":                                     dcdb.ResourceTencentCloudDcdbIsolateHourInstanceOperation(),
			"tencentcloud_dcdb_cancel_dcn_job_operation":                                            dcdb.ResourceTencentCloudDcdbCancelDcnJobOperation(),
			"tencentcloud_dcdb_flush_binlog_operation":                                              dcdb.ResourceTencentCloudDcdbFlushBinlogOperation(),
			"tencentcloud_dcdb_switch_db_instance_ha_operation":                                     dcdb.ResourceTencentCloudDcdbSwitchDbInstanceHaOperation(),
			"tencentcloud_cat_task_set":                                                             cat.ResourceTencentCloudCatTaskSet(),
			"tencentcloud_mariadb_dedicatedcluster_db_instance":                                     mariadb.ResourceTencentCloudMariadbDedicatedclusterDbInstance(),
			"tencentcloud_mariadb_instance":                                                         mariadb.ResourceTencentCloudMariadbInstance(),
			"tencentcloud_mariadb_hour_db_instance":                                                 mariadb.ResourceTencentCloudMariadbHourDbInstance(),
			"tencentcloud_mariadb_account":                                                          mariadb.ResourceTencentCloudMariadbAccount(),
			"tencentcloud_mariadb_parameters":                                                       mariadb.ResourceTencentCloudMariadbParameters(),
			"tencentcloud_mariadb_log_file_retention_period":                                        mariadb.ResourceTencentCloudMariadbLogFileRetentionPeriod(),
			"tencentcloud_mariadb_security_groups":                                                  mariadb.ResourceTencentCloudMariadbSecurityGroups(),
			"tencentcloud_mariadb_encrypt_attributes":                                               mariadb.ResourceTencentCloudMariadbEncryptAttributes(),
			"tencentcloud_mariadb_account_privileges":                                               mariadb.ResourceTencentCloudMariadbAccountPrivileges(),
			"tencentcloud_mariadb_operate_hour_db_instance":                                         mariadb.ResourceTencentCloudMariadbOperateHourDbInstance(),
			"tencentcloud_mariadb_backup_time":                                                      mariadb.ResourceTencentCloudMariadbBackupTime(),
			"tencentcloud_mariadb_cancel_dcn_job":                                                   mariadb.ResourceTencentCloudMariadbCancelDcnJob(),
			"tencentcloud_mariadb_flush_binlog":                                                     mariadb.ResourceTencentCloudMariadbFlushBinlog(),
			"tencentcloud_mariadb_switch_ha":                                                        mariadb.ResourceTencentCloudMariadbSwitchHA(),
			"tencentcloud_mariadb_restart_instance":                                                 mariadb.ResourceTencentCloudMariadbRestartInstance(),
			"tencentcloud_mariadb_renew_instance":                                                   mariadb.ResourceTencentCloudMariadbRenewInstance(),
			"tencentcloud_mariadb_instance_config":                                                  mariadb.ResourceTencentCloudMariadbInstanceConfig(),
			"tencentcloud_tdcpg_cluster":                                                            tdcpg.ResourceTencentCloudTdcpgCluster(),
			"tencentcloud_tdcpg_instance":                                                           tdcpg.ResourceTencentCloudTdcpgInstance(),
			"tencentcloud_css_watermark":                                                            css.ResourceTencentCloudCssWatermark(),
			"tencentcloud_css_watermark_rule_attachment":                                            css.ResourceTencentCloudCssWatermarkRuleAttachment(),
			"tencentcloud_css_pull_stream_task":                                                     css.ResourceTencentCloudCssPullStreamTask(),
			"tencentcloud_css_live_transcode_template":                                              css.ResourceTencentCloudCssLiveTranscodeTemplate(),
			"tencentcloud_css_live_transcode_rule_attachment":                                       css.ResourceTencentCloudCssLiveTranscodeRuleAttachment(),
			"tencentcloud_css_domain":                                                               css.ResourceTencentCloudCssDomain(),
			"tencentcloud_css_authenticate_domain_owner_operation":                                  css.ResourceTencentCloudCssAuthenticateDomainOwnerOperation(),
			"tencentcloud_css_play_domain_cert_attachment":                                          css.ResourceTencentCloudCssPlayDomainCertAttachment(),
			"tencentcloud_css_play_auth_key_config":                                                 css.ResourceTencentCloudCssPlayAuthKeyConfig(),
			"tencentcloud_css_push_auth_key_config":                                                 css.ResourceTencentCloudCssPushAuthKeyConfig(),
			"tencentcloud_css_backup_stream":                                                        css.ResourceTencentCloudCssBackupStream(),
			"tencentcloud_css_callback_rule_attachment":                                             css.ResourceTencentCloudCssCallbackRuleAttachment(),
			"tencentcloud_css_callback_template":                                                    css.ResourceTencentCloudCssCallbackTemplate(),
			"tencentcloud_css_domain_referer":                                                       css.ResourceTencentCloudCssDomainReferer(),
			"tencentcloud_css_enable_optimal_switching":                                             css.ResourceTencentCloudCssEnableOptimalSwitching(),
			"tencentcloud_css_record_rule_attachment":                                               css.ResourceTencentCloudCssRecordRuleAttachment(),
			"tencentcloud_css_record_template":                                                      css.ResourceTencentCloudCssRecordTemplate(),
			"tencentcloud_css_snapshot_rule_attachment":                                             css.ResourceTencentCloudCssSnapshotRuleAttachment(),
			"tencentcloud_css_snapshot_template":                                                    css.ResourceTencentCloudCssSnapshotTemplate(),
			"tencentcloud_css_pad_template":                                                         css.ResourceTencentCloudCssPadTemplate(),
			"tencentcloud_css_pad_rule_attachment":                                                  css.ResourceTencentCloudCssPadRuleAttachment(),
			"tencentcloud_css_timeshift_template":                                                   css.ResourceTencentCloudCssTimeshiftTemplate(),
			"tencentcloud_css_timeshift_rule_attachment":                                            css.ResourceTencentCloudCssTimeshiftRuleAttachment(),
			"tencentcloud_css_stream_monitor":                                                       css.ResourceTencentCloudCssStreamMonitor(),
			"tencentcloud_css_start_stream_monitor":                                                 css.ResourceTencentCloudCssStartStreamMonitor(),
			"tencentcloud_css_pull_stream_task_restart":                                             css.ResourceTencentCloudCssPullStreamTaskRestart(),
			"tencentcloud_pts_project":                                                              pts.ResourceTencentCloudPtsProject(),
			"tencentcloud_pts_alert_channel":                                                        pts.ResourceTencentCloudPtsAlertChannel(),
			"tencentcloud_pts_scenario":                                                             pts.ResourceTencentCloudPtsScenario(),
			"tencentcloud_pts_file":                                                                 pts.ResourceTencentCloudPtsFile(),
			"tencentcloud_pts_job":                                                                  pts.ResourceTencentCloudPtsJob(),
			"tencentcloud_pts_cron_job":                                                             pts.ResourceTencentCloudPtsCronJob(),
			"tencentcloud_pts_tmp_key_generate":                                                     pts.ResourceTencentCloudPtsTmpKeyGenerate(),
			"tencentcloud_pts_cron_job_restart":                                                     pts.ResourceTencentCloudPtsCronJobRestart(),
			"tencentcloud_pts_job_abort":                                                            pts.ResourceTencentCloudPtsJobAbort(),
			"tencentcloud_pts_cron_job_abort":                                                       pts.ResourceTencentCloudPtsCronJobAbort(),
			"tencentcloud_tat_command":                                                              tat.ResourceTencentCloudTatCommand(),
			"tencentcloud_tat_invoker":                                                              tat.ResourceTencentCloudTatInvoker(),
			"tencentcloud_tat_invoker_config":                                                       tat.ResourceTencentCloudTatInvokerConfig(),
			"tencentcloud_tat_invocation_invoke_attachment":                                         tat.ResourceTencentCloudTatInvocationInvokeAttachment(),
			"tencentcloud_tat_invocation_command_attachment":                                        tat.ResourceTencentCloudTatInvocationCommandAttachment(),
			"tencentcloud_organization_org_node":                                                    tco.ResourceTencentCloudOrganizationOrgNode(),
			"tencentcloud_organization_org_member":                                                  tco.ResourceTencentCloudOrganizationOrgMember(),
			"tencentcloud_organization_org_identity":                                                tco.ResourceTencentCloudOrganizationOrgIdentity(),
			"tencentcloud_organization_org_member_email":                                            tco.ResourceTencentCloudOrganizationOrgMemberEmail(),
			"tencentcloud_organization_instance":                                                    tco.ResourceTencentCloudOrganizationOrganization(),
			"tencentcloud_organization_policy_sub_account_attachment":                               tco.ResourceTencentCloudOrganizationPolicySubAccountAttachment(),
			"tencentcloud_organization_org_member_auth_identity_attachment":                         tco.ResourceTencentCloudOrganizationOrgMemberAuthIdentityAttachment(),
			"tencentcloud_organization_org_member_policy_attachment":                                tco.ResourceTencentCloudOrganizationOrgMemberPolicyAttachment(),
			"tencentcloud_organization_org_manage_policy_config":                                    tco.ResourceTencentCloudOrganizationOrgManagePolicyConfig(),
			"tencentcloud_organization_org_manage_policy":                                           tco.ResourceTencentCloudOrganizationOrgManagePolicy(),
			"tencentcloud_organization_org_manage_policy_target":                                    tco.ResourceTencentCloudOrganizationOrgManagePolicyTarget(),
			"tencentcloud_organization_service_assign":                                              tco.ResourceTencentCloudOrganizationServiceAssign(),
			"tencentcloud_identity_center_user":                                                     tco.ResourceTencentCloudIdentityCenterUser(),
			"tencentcloud_identity_center_group":                                                    tco.ResourceTencentCloudIdentityCenterGroup(),
			"tencentcloud_identity_center_user_group_attachment":                                    tco.ResourceTencentCloudIdentityCenterUserGroupAttachment(),
			"tencentcloud_identity_center_external_saml_identity_provider":                          tco.ResourceTencentCloudIdentityCenterExternalSamlIdentityProvider(),
			"tencentcloud_identity_center_role_configuration":                                       tco.ResourceTencentCloudIdentityCenterRoleConfiguration(),
			"tencentcloud_identity_center_role_configuration_permission_policy_attachment":          tco.ResourceTencentCloudIdentityCenterRoleConfigurationPermissionPolicyAttachment(),
			"tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment":   tco.ResourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachment(),
			"tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment": tco.ResourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachment(),
			"tencentcloud_identity_center_user_sync_provisioning":                                   tco.ResourceTencentCloudIdentityCenterUserSyncProvisioning(),
			"tencentcloud_identity_center_role_assignment":                                          tco.ResourceTencentCloudIdentityCenterRoleAssignment(),
			"tencentcloud_invite_organization_member_operation":                                     tco.ResourceTencentCloudInviteOrganizationMemberOperation(),
			"tencentcloud_open_identity_center_operation":                                           tco.ResourceTencentCloudOpenIdentityCenterOperation(),
			"tencentcloud_identity_center_scim_credential_status":                                   tco.ResourceTencentCloudIdentityCenterScimCredentialStatus(),
			"tencentcloud_identity_center_scim_credential":                                          tco.ResourceTencentCloudIdentityCenterScimCredential(),
			"tencentcloud_identity_center_scim_synchronization_status":                              tco.ResourceTencentCloudIdentityCenterScimSynchronizationStatus(),
			"tencentcloud_provision_role_configuration_operation":                                   tco.ResourceTencentCloudProvisionRoleConfigurationOperation(),
			"tencentcloud_dbbrain_sql_filter":                                                       dbbrain.ResourceTencentCloudDbbrainSqlFilter(),
			"tencentcloud_dbbrain_security_audit_log_export_task":                                   dbbrain.ResourceTencentCloudDbbrainSecurityAuditLogExportTask(),
			"tencentcloud_dbbrain_db_diag_report_task":                                              dbbrain.ResourceTencentCloudDbbrainDbDiagReportTask(),
			"tencentcloud_dbbrain_modify_diag_db_instance_operation":                                dbbrain.ResourceTencentCloudDbbrainModifyDiagDbInstanceOperation(),
			"tencentcloud_dbbrain_tdsql_audit_log":                                                  dbbrain.ResourceTencentCloudDbbrainTdsqlAuditLog(),
			"tencentcloud_rum_project":                                                              rum.ResourceTencentCloudRumProject(),
			"tencentcloud_rum_taw_instance":                                                         rum.ResourceTencentCloudRumTawInstance(),
			"tencentcloud_rum_whitelist":                                                            rum.ResourceTencentCloudRumWhitelist(),
			"tencentcloud_rum_offline_log_config_attachment":                                        rum.ResourceTencentCloudRumOfflineLogConfigAttachment(),
			"tencentcloud_rum_instance_status_config":                                               rum.ResourceTencentCloudRumInstanceStatusConfig(),
			"tencentcloud_rum_project_status_config":                                                rum.ResourceTencentCloudRumProjectStatusConfig(),
			"tencentcloud_rum_release_file":                                                         rum.ResourceTencentCloudRumReleaseFile(),
			"tencentcloud_tdmq_rocketmq_cluster":                                                    trocket.ResourceTencentCloudTdmqRocketmqCluster(),
			"tencentcloud_tdmq_rocketmq_namespace":                                                  trocket.ResourceTencentCloudTdmqRocketmqNamespace(),
			"tencentcloud_tdmq_rocketmq_role":                                                       trocket.ResourceTencentCloudTdmqRocketmqRole(),
			"tencentcloud_tdmq_rocketmq_topic":                                                      trocket.ResourceTencentCloudTdmqRocketmqTopic(),
			"tencentcloud_tdmq_rocketmq_group":                                                      trocket.ResourceTencentCloudTdmqRocketmqGroup(),
			"tencentcloud_tdmq_rocketmq_environment_role":                                           trocket.ResourceTencentCloudTdmqRocketmqEnvironmentRole(),
			"tencentcloud_tdmq_rocketmq_vip_instance":                                               trocket.ResourceTencentCloudTdmqRocketmqVipInstance(),
			"tencentcloud_trocket_rocketmq_instance":                                                trocket.ResourceTencentCloudTrocketRocketmqInstance(),
			"tencentcloud_trocket_rocketmq_topic":                                                   trocket.ResourceTencentCloudTrocketRocketmqTopic(),
			"tencentcloud_trocket_rocketmq_consumer_group":                                          trocket.ResourceTencentCloudTrocketRocketmqConsumerGroup(),
			"tencentcloud_trocket_rocketmq_role":                                                    trocket.ResourceTencentCloudTrocketRocketmqRole(),
			"tencentcloud_dts_sync_job":                                                             dts.ResourceTencentCloudDtsSyncJob(),
			"tencentcloud_dts_sync_config":                                                          dts.ResourceTencentCloudDtsSyncConfig(),
			"tencentcloud_dts_sync_check_job_operation":                                             dts.ResourceTencentCloudDtsSyncCheckJobOperation(),
			"tencentcloud_dts_sync_job_resume_operation":                                            dts.ResourceTencentCloudDtsSyncJobResumeOperation(),
			"tencentcloud_dts_sync_job_start_operation":                                             dts.ResourceTencentCloudDtsSyncJobStartOperation(),
			"tencentcloud_dts_sync_job_stop_operation":                                              dts.ResourceTencentCloudDtsSyncJobStopOperation(),
			"tencentcloud_dts_sync_job_resize_operation":                                            dts.ResourceTencentCloudDtsSyncJobResizeOperation(),
			"tencentcloud_dts_sync_job_recover_operation":                                           dts.ResourceTencentCloudDtsSyncJobRecoverOperation(),
			"tencentcloud_dts_sync_job_isolate_operation":                                           dts.ResourceTencentCloudDtsSyncJobIsolateOperation(),
			"tencentcloud_dts_sync_job_continue_operation":                                          dts.ResourceTencentCloudDtsSyncJobContinueOperation(),
			"tencentcloud_dts_sync_job_pause_operation":                                             dts.ResourceTencentCloudDtsSyncJobPauseOperation(),
			"tencentcloud_dts_migrate_service":                                                      dts.ResourceTencentCloudDtsMigrateService(),
			"tencentcloud_dts_migrate_job":                                                          dts.ResourceTencentCloudDtsMigrateJob(),
			"tencentcloud_dts_migrate_job_config":                                                   dts.ResourceTencentCloudDtsMigrateJobConfig(),
			"tencentcloud_dts_migrate_job_start_operation":                                          dts.ResourceTencentCloudDtsMigrateJobStartOperation(),
			"tencentcloud_dts_migrate_job_resume_operation":                                         dts.ResourceTencentCloudDtsMigrateJobResumeOperation(),
			"tencentcloud_dts_compare_task_stop_operation":                                          dts.ResourceTencentCloudDtsCompareTaskStopOperation(),
			"tencentcloud_dts_compare_task":                                                         dts.ResourceTencentCloudDtsCompareTask(),
			"tencentcloud_cvm_hpc_cluster":                                                          cvm.ResourceTencentCloudCvmHpcCluster(),
			"tencentcloud_vpc_flow_log":                                                             fl.ResourceTencentCloudVpcFlowLog(),
			"tencentcloud_vpc_end_point_service":                                                    pls.ResourceTencentCloudVpcEndPointService(),
			"tencentcloud_vpc_end_point":                                                            pls.ResourceTencentCloudVpcEndPoint(),
			"tencentcloud_vpc_end_point_service_white_list":                                         pls.ResourceTencentCloudVpcEndPointServiceWhiteList(),
			"tencentcloud_vpc_enable_end_point_connect":                                             pls.ResourceTencentCloudVpcEnableEndPointConnect(),
			"tencentcloud_ci_bucket_attachment":                                                     ci.ResourceTencentCloudCiBucketAttachment(),
			"tencentcloud_tcmq_queue":                                                               tcmq.ResourceTencentCloudTcmqQueue(),
			"tencentcloud_tcmq_topic":                                                               tcmq.ResourceTencentCloudTcmqTopic(),
			"tencentcloud_tcmq_subscribe":                                                           tcmq.ResourceTencentCloudTcmqSubscribe(),
			"tencentcloud_ci_bucket_pic_style":                                                      ci.ResourceTencentCloudCiBucketPicStyle(),
			"tencentcloud_ci_hot_link":                                                              ci.ResourceTencentCloudCiHotLink(),
			"tencentcloud_ci_media_snapshot_template":                                               ci.ResourceTencentCloudCiMediaSnapshotTemplate(),
			"tencentcloud_ci_media_transcode_template":                                              ci.ResourceTencentCloudCiMediaTranscodeTemplate(),
			"tencentcloud_ci_media_animation_template":                                              ci.ResourceTencentCloudCiMediaAnimationTemplate(),
			"tencentcloud_ci_media_concat_template":                                                 ci.ResourceTencentCloudCiMediaConcatTemplate(),
			"tencentcloud_ci_media_video_process_template":                                          ci.ResourceTencentCloudCiMediaVideoProcessTemplate(),
			"tencentcloud_ci_media_video_montage_template":                                          ci.ResourceTencentCloudCiMediaVideoMontageTemplate(),
			"tencentcloud_ci_media_voice_separate_template":                                         ci.ResourceTencentCloudCiMediaVoiceSeparateTemplate(),
			"tencentcloud_ci_media_super_resolution_template":                                       ci.ResourceTencentCloudCiMediaSuperResolutionTemplate(),
			"tencentcloud_ci_media_pic_process_template":                                            ci.ResourceTencentCloudCiMediaPicProcessTemplate(),
			"tencentcloud_ci_media_watermark_template":                                              ci.ResourceTencentCloudCiMediaWatermarkTemplate(),
			"tencentcloud_ci_media_tts_template":                                                    ci.ResourceTencentCloudCiMediaTtsTemplate(),
			"tencentcloud_ci_media_transcode_pro_template":                                          ci.ResourceTencentCloudCiMediaTranscodeProTemplate(),
			"tencentcloud_ci_media_smart_cover_template":                                            ci.ResourceTencentCloudCiMediaSmartCoverTemplate(),
			"tencentcloud_ci_media_speech_recognition_template":                                     ci.ResourceTencentCloudCiMediaSpeechRecognitionTemplate(),
			"tencentcloud_ci_guetzli":                                                               ci.ResourceTencentCloudCIGuetzli(),
			"tencentcloud_ci_original_image_protection":                                             ci.ResourceTencentCloudCIOriginalImageProtection(),
			"tencentcloud_cynosdb_audit_log_file":                                                   cynosdb.ResourceTencentCloudCynosdbAuditLogFile(),
			"tencentcloud_cynosdb_security_group":                                                   cynosdb.ResourceTencentCloudCynosdbSecurityGroup(),
			"tencentcloud_dayu_ddos_ip_attachment_v2":                                               dayuv2.ResourceTencentCloudDayuDDosIpAttachmentV2(),
			"tencentcloud_antiddos_ddos_black_white_ip":                                             dayuv2.ResourceTencentCloudAntiddosDdosBlackWhiteIp(),
			"tencentcloud_antiddos_ddos_geo_ip_block_config":                                        dayuv2.ResourceTencentCloudAntiddosDdosGeoIpBlockConfig(),
			"tencentcloud_antiddos_ddos_speed_limit_config":                                         dayuv2.ResourceTencentCloudAntiddosDdosSpeedLimitConfig(),
			"tencentcloud_antiddos_default_alarm_threshold":                                         dayuv2.ResourceTencentCloudAntiddosDefaultAlarmThreshold(),
			"tencentcloud_antiddos_scheduling_domain_user_name":                                     dayuv2.ResourceTencentCloudAntiddosSchedulingDomainUserName(),
			"tencentcloud_antiddos_ip_alarm_threshold_config":                                       dayuv2.ResourceTencentCloudAntiddosIpAlarmThresholdConfig(),
			"tencentcloud_antiddos_packet_filter_config":                                            dayuv2.ResourceTencentCloudAntiddosPacketFilterConfig(),
			"tencentcloud_antiddos_port_acl_config":                                                 dayuv2.ResourceTencentCloudAntiddosPortAclConfig(),
			"tencentcloud_antiddos_cc_black_white_ip":                                               dayuv2.ResourceTencentCloudAntiddosCcBlackWhiteIp(),
			"tencentcloud_antiddos_cc_precision_policy":                                             dayuv2.ResourceTencentCloudAntiddosCcPrecisionPolicy(),
			"tencentcloud_tsf_microservice":                                                         tsf.ResourceTencentCloudTsfMicroservice(),
			"tencentcloud_tsf_application_config":                                                   tsf.ResourceTencentCloudTsfApplicationConfig(),
			"tencentcloud_cvm_launch_template":                                                      cvm.ResourceTencentCloudCvmLaunchTemplate(),
			"tencentcloud_tsf_cluster":                                                              tsf.ResourceTencentCloudTsfCluster(),
			"tencentcloud_tsf_api_group":                                                            tsf.ResourceTencentCloudTsfApiGroup(),
			"tencentcloud_tsf_namespace":                                                            tsf.ResourceTencentCloudTsfNamespace(),
			"tencentcloud_tsf_path_rewrite":                                                         tsf.ResourceTencentCloudTsfPathRewrite(),
			"tencentcloud_tsf_unit_rule":                                                            tsf.ResourceTencentCloudTsfUnitRule(),
			"tencentcloud_tsf_task":                                                                 tsf.ResourceTencentCloudTsfTask(),
			"tencentcloud_tsf_config_template":                                                      tsf.ResourceTencentCloudTsfConfigTemplate(),
			"tencentcloud_tsf_api_rate_limit_rule":                                                  tsf.ResourceTencentCloudTsfApiRateLimitRule(),
			"tencentcloud_tsf_application_release_config":                                           tsf.ResourceTencentCloudTsfApplicationReleaseConfig(),
			"tencentcloud_tsf_lane":                                                                 tsf.ResourceTencentCloudTsfLane(),
			"tencentcloud_tsf_lane_rule":                                                            tsf.ResourceTencentCloudTsfLaneRule(),
			"tencentcloud_tsf_group":                                                                tsf.ResourceTencentCloudTsfGroup(),
			"tencentcloud_tsf_repository":                                                           tsf.ResourceTencentCloudTsfRepository(),
			"tencentcloud_tsf_application":                                                          tsf.ResourceTencentCloudTsfApplication(),
			"tencentcloud_tsf_application_public_config_release":                                    tsf.ResourceTencentCloudTsfApplicationPublicConfigRelease(),
			"tencentcloud_tsf_application_public_config":                                            tsf.ResourceTencentCloudTsfApplicationPublicConfig(),
			"tencentcloud_tsf_application_file_config_release":                                      tsf.ResourceTencentCloudTsfApplicationFileConfigRelease(),
			"tencentcloud_tsf_instances_attachment":                                                 tsf.ResourceTencentCloudTsfInstancesAttachment(),
			"tencentcloud_tsf_bind_api_group":                                                       tsf.ResourceTencentCloudTsfBindApiGroup(),
			"tencentcloud_tsf_application_file_config":                                              tsf.ResourceTencentCloudTsfApplicationFileConfig(),
			"tencentcloud_tsf_enable_unit_rule":                                                     tsf.ResourceTencentCloudTsfEnableUnitRule(),
			"tencentcloud_tsf_deploy_container_group":                                               tsf.ResourceTencentCloudTsfDeployContainerGroup(),
			"tencentcloud_tsf_deploy_vm_group":                                                      tsf.ResourceTencentCloudTsfDeployVmGroup(),
			"tencentcloud_tsf_release_api_group":                                                    tsf.ResourceTencentCloudTsfReleaseApiGroup(),
			"tencentcloud_tsf_operate_container_group":                                              tsf.ResourceTencentCloudTsfOperateContainerGroup(),
			"tencentcloud_tsf_operate_group":                                                        tsf.ResourceTencentCloudTsfOperateGroup(),
			"tencentcloud_tsf_unit_namespace":                                                       tsf.ResourceTencentCloudTsfUnitNamespace(),
			"tencentcloud_mps_workflow":                                                             mps.ResourceTencentCloudMpsWorkflow(),
			"tencentcloud_mps_enable_workflow_config":                                               mps.ResourceTencentCloudMpsEnableWorkflowConfig(),
			"tencentcloud_mps_flow":                                                                 mps.ResourceTencentCloudMpsFlow(),
			"tencentcloud_mps_input":                                                                mps.ResourceTencentCloudMpsInput(),
			"tencentcloud_mps_output":                                                               mps.ResourceTencentCloudMpsOutput(),
			"tencentcloud_mps_content_review_template":                                              mps.ResourceTencentCloudMpsContentReviewTemplate(),
			"tencentcloud_mps_start_flow_operation":                                                 mps.ResourceTencentCloudMpsStartFlowOperation(),
			"tencentcloud_mps_event":                                                                mps.ResourceTencentCloudMpsEvent(),
			"tencentcloud_mps_execute_function_operation":                                           mps.ResourceTencentCloudMpsExecuteFunctionOperation(),
			"tencentcloud_mps_manage_task_operation":                                                mps.ResourceTencentCloudMpsManageTaskOperation(),
			"tencentcloud_mps_transcode_template":                                                   mps.ResourceTencentCloudMpsTranscodeTemplate(),
			"tencentcloud_mps_watermark_template":                                                   mps.ResourceTencentCloudMpsWatermarkTemplate(),
			"tencentcloud_mps_image_sprite_template":                                                mps.ResourceTencentCloudMpsImageSpriteTemplate(),
			"tencentcloud_mps_snapshot_by_timeoffset_template":                                      mps.ResourceTencentCloudMpsSnapshotByTimeoffsetTemplate(),
			"tencentcloud_mps_sample_snapshot_template":                                             mps.ResourceTencentCloudMpsSampleSnapshotTemplate(),
			"tencentcloud_mps_animated_graphics_template":                                           mps.ResourceTencentCloudMpsAnimatedGraphicsTemplate(),
			"tencentcloud_mps_ai_recognition_template":                                              mps.ResourceTencentCloudMpsAiRecognitionTemplate(),
			"tencentcloud_mps_ai_analysis_template":                                                 mps.ResourceTencentCloudMpsAiAnalysisTemplate(),
			"tencentcloud_mps_adaptive_dynamic_streaming_template":                                  mps.ResourceTencentCloudMpsAdaptiveDynamicStreamingTemplate(),
			"tencentcloud_mps_person_sample":                                                        mps.ResourceTencentCloudMpsPersonSample(),
			"tencentcloud_mps_withdraws_watermark_operation":                                        mps.ResourceTencentCloudMpsWithdrawsWatermarkOperation(),
			"tencentcloud_mps_process_live_stream_operation":                                        mps.ResourceTencentCloudMpsProcessLiveStreamOperation(),
			"tencentcloud_mps_edit_media_operation":                                                 mps.ResourceTencentCloudMpsEditMediaOperation(),
			"tencentcloud_mps_word_sample":                                                          mps.ResourceTencentCloudMpsWordSample(),
			"tencentcloud_mps_schedule":                                                             mps.ResourceTencentCloudMpsSchedule(),
			"tencentcloud_mps_enable_schedule_config":                                               mps.ResourceTencentCloudMpsEnableScheduleConfig(),
			"tencentcloud_mps_process_media_operation":                                              mps.ResourceTencentCloudMpsProcessMediaOperation(),
			"tencentcloud_cbs_disk_backup":                                                          cbs.ResourceTencentCloudCbsDiskBackup(),
			"tencentcloud_cbs_snapshot_share_permission":                                            cbs.ResourceTencentCloudCbsSnapshotSharePermission(),
			"tencentcloud_cbs_disk_backup_rollback_operation":                                       cbs.ResourceTencentCloudCbsDiskBackupRollbackOperation(),
			"tencentcloud_chdfs_access_group":                                                       chdfs.ResourceTencentCloudChdfsAccessGroup(),
			"tencentcloud_chdfs_access_rule":                                                        chdfs.ResourceTencentCloudChdfsAccessRule(),
			"tencentcloud_chdfs_file_system":                                                        chdfs.ResourceTencentCloudChdfsFileSystem(),
			"tencentcloud_chdfs_life_cycle_rule":                                                    chdfs.ResourceTencentCloudChdfsLifeCycleRule(),
			"tencentcloud_chdfs_mount_point":                                                        chdfs.ResourceTencentCloudChdfsMountPoint(),
			"tencentcloud_chdfs_mount_point_attachment":                                             chdfs.ResourceTencentCloudChdfsMountPointAttachment(),
			"tencentcloud_mdl_stream_live_input":                                                    mdl.ResourceTencentCloudMdlStreamLiveInput(),
			"tencentcloud_lighthouse_blueprint":                                                     lighthouse.ResourceTencentCloudLighthouseBlueprint(),
			"tencentcloud_cvm_launch_template_version":                                              cvm.ResourceTencentCloudCvmLaunchTemplateVersion(),
			"tencentcloud_apm_instance":                                                             apm.ResourceTencentCloudApmInstance(),
			"tencentcloud_cvm_launch_template_default_version":                                      cvm.ResourceTencentCloudCvmLaunchTemplateDefaultVersion(),
			"tencentcloud_lighthouse_firewall_rule":                                                 lighthouse.ResourceTencentCloudLighthouseFirewallRule(),
			"tencentcloud_cvm_security_group_attachment":                                            cvm.ResourceTencentCloudCvmSecurityGroupAttachment(),
			"tencentcloud_cvm_reboot_instance":                                                      cvm.ResourceTencentCloudCvmRebootInstance(),
			"tencentcloud_cvm_chc_config":                                                           cvm.ResourceTencentCloudCvmChcConfig(),
			"tencentcloud_cvm_sync_image":                                                           cvm.ResourceTencentCloudCvmSyncImage(),
			"tencentcloud_cvm_renew_instance":                                                       cvm.ResourceTencentCloudCvmRenewInstance(),
			"tencentcloud_cvm_export_images":                                                        cvm.ResourceTencentCloudCvmExportImages(),
			"tencentcloud_cvm_image_share_permission":                                               cvm.ResourceTencentCloudCvmImageSharePermission(),
			"tencentcloud_cvm_import_image":                                                         cvm.ResourceTencentCloudCvmImportImage(),
			"tencentcloud_cvm_renew_host":                                                           cvm.ResourceTencentCloudCvmRenewHost(),
			"tencentcloud_cvm_program_fpga_image":                                                   cvm.ResourceTencentCloudCvmProgramFpgaImage(),
			"tencentcloud_cvm_modify_instance_disk_type":                                            cvm.ResourceTencentCloudCvmModifyInstanceDiskType(),
			"tencentcloud_cvm_action_timer":                                                         cvm.ResourceTencentCloudCvmActionTimer(),
			"tencentcloud_lighthouse_disk_backup":                                                   lighthouse.ResourceTencentCloudLighthouseDiskBackup(),
			"tencentcloud_lighthouse_apply_disk_backup":                                             lighthouse.ResourceTencentCloudLighthouseApplyDiskBackup(),
			"tencentcloud_lighthouse_disk_attachment":                                               lighthouse.ResourceTencentCloudLighthouseDiskAttachment(),
			"tencentcloud_lighthouse_key_pair":                                                      lighthouse.ResourceTencentCloudLighthouseKeyPair(),
			"tencentcloud_lighthouse_snapshot":                                                      lighthouse.ResourceTencentCloudLighthouseSnapshot(),
			"tencentcloud_lighthouse_apply_instance_snapshot":                                       lighthouse.ResourceTencentCloudLighthouseApplyInstanceSnapshot(),
			"tencentcloud_lighthouse_start_instance":                                                lighthouse.ResourceTencentCloudLighthouseStartInstance(),
			"tencentcloud_lighthouse_stop_instance":                                                 lighthouse.ResourceTencentCloudLighthouseStopInstance(),
			"tencentcloud_lighthouse_reboot_instance":                                               lighthouse.ResourceTencentCloudLighthouseRebootInstance(),
			"tencentcloud_lighthouse_key_pair_attachment":                                           lighthouse.ResourceTencentCloudLighthouseKeyPairAttachment(),
			"tencentcloud_lighthouse_disk":                                                          lighthouse.ResourceTencentCloudLighthouseDisk(),
			"tencentcloud_lighthouse_renew_disk":                                                    lighthouse.ResourceTencentCloudLighthouseRenewDisk(),
			"tencentcloud_lighthouse_renew_instance":                                                lighthouse.ResourceTencentCloudLighthouseRenewInstance(),
			"tencentcloud_clickhouse_backup":                                                        cdwch.ResourceTencentCloudClickhouseBackup(),
			"tencentcloud_clickhouse_backup_strategy":                                               cdwch.ResourceTencentCloudClickhouseBackupStrategy(),
			"tencentcloud_clickhouse_recover_backup_job":                                            cdwch.ResourceTencentCloudClickhouseRecoverBackupJob(),
			"tencentcloud_clickhouse_delete_backup_data":                                            cdwch.ResourceTencentCloudClickhouseDeleteBackupData(),
			"tencentcloud_clickhouse_account":                                                       cdwch.ResourceTencentCloudClickhouseAccount(),
			"tencentcloud_clickhouse_account_permission":                                            cdwch.ResourceTencentCloudClickhouseAccountPermission(),
			"tencentcloud_api_gateway_api_doc":                                                      apigateway.ResourceTencentCloudAPIGatewayAPIDoc(),
			"tencentcloud_api_gateway_api_app":                                                      apigateway.ResourceTencentCloudAPIGatewayAPIApp(),
			"tencentcloud_api_gateway_update_api_app_key":                                           apigateway.ResourceTencentCloudApiGatewayUpdateApiAppKey(),
			"tencentcloud_api_gateway_import_open_api":                                              apigateway.ResourceTencentCloudApiGatewayImportOpenApi(),
			"tencentcloud_tse_instance":                                                             tse.ResourceTencentCloudTseInstance(),
			"tencentcloud_tse_cngw_gateway":                                                         tse.ResourceTencentCloudTseCngwGateway(),
			"tencentcloud_tse_cngw_group":                                                           tse.ResourceTencentCloudTseCngwGroup(),
			"tencentcloud_tse_cngw_service":                                                         tse.ResourceTencentCloudTseCngwService(),
			"tencentcloud_tse_cngw_service_rate_limit":                                              tse.ResourceTencentCloudTseCngwServiceRateLimit(),
			"tencentcloud_tse_cngw_route":                                                           tse.ResourceTencentCloudTseCngwRoute(),
			"tencentcloud_tse_cngw_route_rate_limit":                                                tse.ResourceTencentCloudTseCngwRouteRateLimit(),
			"tencentcloud_tse_cngw_canary_rule":                                                     tse.ResourceTencentCloudTseCngwCanaryRule(),
			"tencentcloud_tse_cngw_certificate":                                                     tse.ResourceTencentCloudTseCngwCertificate(),
			"tencentcloud_tse_waf_protection":                                                       tse.ResourceTencentCloudTseWafProtection(),
			"tencentcloud_tse_waf_domains":                                                          tse.ResourceTencentCloudTseWafDomains(),
			"tencentcloud_tse_cngw_network":                                                         tse.ResourceTencentCloudTseCngwNetwork(),
			"tencentcloud_tse_cngw_network_access_control":                                          tse.ResourceTencentCloudTseCngwNetworkAccessControl(),
			"tencentcloud_tse_cngw_strategy":                                                        tse.ResourceTencentCloudTseCngwStrategy(),
			"tencentcloud_tse_cngw_strategy_bind_group":                                             tse.ResourceTencentCloudTseCngwStrategyBindGroup(),
			"tencentcloud_clickhouse_instance":                                                      cdwch.ResourceTencentCloudClickhouseInstance(),
			"tencentcloud_cls_kafka_recharge":                                                       cls.ResourceTencentCloudClsKafkaRecharge(),
			"tencentcloud_cls_scheduled_sql":                                                        cls.ResourceTencentCloudClsScheduledSql(),
			"tencentcloud_eb_event_transform":                                                       eb.ResourceTencentCloudEbEventTransform(),
			"tencentcloud_eb_event_bus":                                                             eb.ResourceTencentCloudEbEventBus(),
			"tencentcloud_eb_event_rule":                                                            eb.ResourceTencentCloudEbEventRule(),
			"tencentcloud_eb_event_target":                                                          eb.ResourceTencentCloudEbEventTarget(),
			"tencentcloud_eb_put_events":                                                            eb.ResourceTencentCloudEbPutEvents(),
			"tencentcloud_eb_event_connector":                                                       eb.ResourceTencentCloudEbEventConnector(),
			"tencentcloud_dlc_user":                                                                 dlc.ResourceTencentCloudDlcUser(),
			"tencentcloud_dlc_work_group":                                                           dlc.ResourceTencentCloudDlcWorkGroup(),
			"tencentcloud_dlc_data_engine":                                                          dlc.ResourceTencentCloudDlcDataEngine(),
			"tencentcloud_dlc_suspend_resume_data_engine":                                           dlc.ResourceTencentCloudDlcSuspendResumeDataEngine(),
			"tencentcloud_dlc_rollback_data_engine_image_operation":                                 dlc.ResourceTencentCloudDlcRollbackDataEngineImageOperation(),
			"tencentcloud_dlc_add_users_to_work_group_attachment":                                   dlc.ResourceTencentCloudDlcAddUsersToWorkGroupAttachment(),
			"tencentcloud_dlc_store_location_config":                                                dlc.ResourceTencentCloudDlcStoreLocationConfig(),
			"tencentcloud_dlc_modify_data_engine_description_operation":                             dlc.ResourceTencentCloudDlcModifyDataEngineDescriptionOperation(),
			"tencentcloud_dlc_modify_user_typ_operation":                                            dlc.ResourceTencentCloudDlcModifyUserTypOperation(),
			"tencentcloud_dlc_renew_data_engine_operation":                                          dlc.ResourceTencentCloudDlcRenewDataEngineOperation(),
			"tencentcloud_dlc_restart_data_engine_operation":                                        dlc.ResourceTencentCloudDlcRestartDataEngineOperation(),
			"tencentcloud_dlc_attach_user_policy_operation":                                         dlc.ResourceTencentCloudDlcAttachUserPolicyOperation(),
			"tencentcloud_dlc_detach_user_policy_operation":                                         dlc.ResourceTencentCloudDlcDetachUserPolicyOperation(),
			"tencentcloud_dlc_attach_work_group_policy_operation":                                   dlc.ResourceTencentCloudDlcAttachWorkGroupPolicyOperation(),
			"tencentcloud_dlc_detach_work_group_policy_operation":                                   dlc.ResourceTencentCloudDlcDetachWorkGroupPolicyOperation(),
			"tencentcloud_dlc_switch_data_engine_image_operation":                                   dlc.ResourceTencentCloudDlcSwitchDataEngineImageOperation(),
			"tencentcloud_dlc_update_data_engine_config_operation":                                  dlc.ResourceTencentCloudDlcUpdateDataEngineConfigOperation(),
			"tencentcloud_dlc_upgrade_data_engine_image_operation":                                  dlc.ResourceTencentCloudDlcUpgradeDataEngineImageOperation(),
			"tencentcloud_dlc_bind_work_groups_to_user_attachment":                                  dlc.ResourceTencentCloudDlcBindWorkGroupsToUserAttachment(),
			"tencentcloud_dlc_update_row_filter_operation":                                          dlc.ResourceTencentCloudDlcUpdateRowFilterOperation(),
			"tencentcloud_dlc_user_data_engine_config":                                              dlc.ResourceTencentCloudDlcUserDataEngineConfig(),
			"tencentcloud_waf_custom_rule":                                                          waf.ResourceTencentCloudWafCustomRule(),
			"tencentcloud_waf_custom_white_rule":                                                    waf.ResourceTencentCloudWafCustomWhiteRule(),
			"tencentcloud_waf_clb_domain":                                                           waf.ResourceTencentCloudWafClbDomain(),
			"tencentcloud_waf_saas_domain":                                                          waf.ResourceTencentCloudWafSaasDomain(),
			"tencentcloud_waf_clb_instance":                                                         waf.ResourceTencentCloudWafClbInstance(),
			"tencentcloud_waf_saas_instance":                                                        waf.ResourceTencentCloudWafSaasInstance(),
			"tencentcloud_waf_anti_fake":                                                            waf.ResourceTencentCloudWafAntiFake(),
			"tencentcloud_waf_anti_info_leak":                                                       waf.ResourceTencentCloudWafAntiInfoLeak(),
			"tencentcloud_waf_auto_deny_rules":                                                      waf.ResourceTencentCloudWafAutoDenyRules(),
			"tencentcloud_waf_module_status":                                                        waf.ResourceTencentCloudWafModuleStatus(),
			"tencentcloud_waf_protection_mode":                                                      waf.ResourceTencentCloudWafProtectionMode(),
			"tencentcloud_waf_web_shell":                                                            waf.ResourceTencentCloudWafWebShell(),
			"tencentcloud_waf_cc":                                                                   waf.ResourceTencentCloudWafCc(),
			"tencentcloud_waf_cc_auto_status":                                                       waf.ResourceTencentCloudWafCcAutoStatus(),
			"tencentcloud_waf_cc_session":                                                           waf.ResourceTencentCloudWafCcSession(),
			"tencentcloud_waf_ip_access_control":                                                    waf.ResourceTencentCloudWafIpAccessControl(),
			"tencentcloud_waf_ip_access_control_v2":                                                 waf.ResourceTencentCloudWafIpAccessControlV2(),
			"tencentcloud_wedata_rule_template":                                                     wedata.ResourceTencentCloudWedataRuleTemplate(),
			"tencentcloud_wedata_datasource":                                                        wedata.ResourceTencentCloudWedataDatasource(),
			"tencentcloud_wedata_function":                                                          wedata.ResourceTencentCloudWedataFunction(),
			"tencentcloud_wedata_resource":                                                          wedata.ResourceTencentCloudWedataResource(),
			"tencentcloud_wedata_script":                                                            wedata.ResourceTencentCloudWedataScript(),
			"tencentcloud_wedata_dq_rule":                                                           wedata.ResourceTencentCloudWedataDqRule(),
			"tencentcloud_wedata_baseline":                                                          wedata.ResourceTencentCloudWedataBaseline(),
			"tencentcloud_wedata_integration_offline_task":                                          wedata.ResourceTencentCloudWedataIntegrationOfflineTask(),
			"tencentcloud_wedata_integration_realtime_task":                                         wedata.ResourceTencentCloudWedataIntegrationRealtimeTask(),
			"tencentcloud_wedata_integration_task_node":                                             wedata.ResourceTencentCloudWedataIntegrationTaskNode(),
			"tencentcloud_cfw_address_template":                                                     cfw.ResourceTencentCloudCfwAddressTemplate(),
			"tencentcloud_cfw_block_ignore":                                                         cfw.ResourceTencentCloudCfwBlockIgnore(),
			"tencentcloud_cfw_edge_policy":                                                          cfw.ResourceTencentCloudCfwEdgePolicy(),
			"tencentcloud_cfw_nat_instance":                                                         cfw.ResourceTencentCloudCfwNatInstance(),
			"tencentcloud_cfw_nat_policy":                                                           cfw.ResourceTencentCloudCfwNatPolicy(),
			"tencentcloud_cfw_vpc_instance":                                                         cfw.ResourceTencentCloudCfwVpcInstance(),
			"tencentcloud_cfw_vpc_policy":                                                           cfw.ResourceTencentCloudCfwVpcPolicy(),
			"tencentcloud_cfw_sync_asset":                                                           cfw.ResourceTencentCloudCfwSyncAsset(),
			"tencentcloud_cfw_sync_route":                                                           cfw.ResourceTencentCloudCfwSyncRoute(),
			"tencentcloud_cfw_nat_firewall_switch":                                                  cfw.ResourceTencentCloudCfwNatFirewallSwitch(),
			"tencentcloud_cfw_vpc_firewall_switch":                                                  cfw.ResourceTencentCloudCfwVpcFirewallSwitch(),
			"tencentcloud_cfw_edge_firewall_switch":                                                 cfw.ResourceTencentCloudCfwEdgeFirewallSwitch(),
			"tencentcloud_sg_rule":                                                                  cfw.ResourceTencentCloudSgRule(),
			"tencentcloud_dasb_acl":                                                                 bh.ResourceTencentCloudDasbAcl(),
			"tencentcloud_dasb_cmd_template":                                                        bh.ResourceTencentCloudDasbCmdTemplate(),
			"tencentcloud_dasb_device_group":                                                        bh.ResourceTencentCloudDasbDeviceGroup(),
			"tencentcloud_dasb_user":                                                                bh.ResourceTencentCloudDasbUser(),
			"tencentcloud_dasb_device_account":                                                      bh.ResourceTencentCloudDasbDeviceAccount(),
			"tencentcloud_dasb_device_group_members":                                                bh.ResourceTencentCloudDasbDeviceGroupMembers(),
			"tencentcloud_dasb_user_group_members":                                                  bh.ResourceTencentCloudDasbUserGroupMembers(),
			"tencentcloud_dasb_bind_device_resource":                                                bh.ResourceTencentCloudDasbBindDeviceResource(),
			"tencentcloud_dasb_resource":                                                            bh.ResourceTencentCloudDasbResource(),
			"tencentcloud_dasb_device":                                                              bh.ResourceTencentCloudDasbDevice(),
			"tencentcloud_dasb_user_group":                                                          bh.ResourceTencentCloudDasbUserGroup(),
			"tencentcloud_dasb_reset_user":                                                          bh.ResourceTencentCloudDasbResetUser(),
			"tencentcloud_dasb_bind_device_account_private_key":                                     bh.ResourceTencentCloudDasbBindDeviceAccountPrivateKey(),
			"tencentcloud_dasb_bind_device_account_password":                                        bh.ResourceTencentCloudDasbBindDeviceAccountPassword(),
			"tencentcloud_ssl_check_certificate_chain_operation":                                    ssl.ResourceTencentCloudSslCheckCertificateChainOperation(),
			"tencentcloud_ssl_complete_certificate_operation":                                       ssl.ResourceTencentCloudSslCompleteCertificateOperation(),
			"tencentcloud_ssl_deploy_certificate_instance_operation":                                ssl.ResourceTencentCloudSslDeployCertificateInstanceOperation(),
			"tencentcloud_ssl_deploy_certificate_record_retry_operation":                            ssl.ResourceTencentCloudSslDeployCertificateRecordRetryOperation(),
			"tencentcloud_ssl_deploy_certificate_record_rollback_operation":                         ssl.ResourceTencentCloudSslDeployCertificateRecordRollbackOperation(),
			"tencentcloud_ssl_download_certificate_operation":                                       ssl.ResourceTencentCloudSslDownloadCertificateOperation(),
			"tencentcloud_ssl_check_certificate_domain_verification_operation":                      ssl.ResourceTencentCloudSslCheckCertificateDomainVerificationOperation(),
			"tencentcloud_cwp_license_order":                                                        cwp.ResourceTencentCloudCwpLicenseOrder(),
			"tencentcloud_cwp_license_bind_attachment":                                              cwp.ResourceTencentCloudCwpLicenseBindAttachment(),
			"tencentcloud_ssl_replace_certificate_operation":                                        ssl.ResourceTencentCloudSslReplaceCertificateOperation(),
			"tencentcloud_ssl_revoke_certificate_operation":                                         ssl.ResourceTencentCloudSslRevokeCertificateOperation(),
			"tencentcloud_ssl_update_certificate_instance_operation":                                ssl.ResourceTencentCloudSslUpdateCertificateInstanceOperation(),
			"tencentcloud_ssl_update_certificate_record_retry_operation":                            ssl.ResourceTencentCloudSslUpdateCertificateRecordRetryOperation(),
			"tencentcloud_ssl_update_certificate_record_rollback_operation":                         ssl.ResourceTencentCloudSslUpdateCertificateRecordRollbackOperation(),
			"tencentcloud_ssl_upload_revoke_letter_operation":                                       ssl.ResourceTencentCloudSslUploadRevokeLetterOperation(),
			"tencentcloud_bi_project":                                                               bi.ResourceTencentCloudBiProject(),
			"tencentcloud_bi_user_role":                                                             bi.ResourceTencentCloudBiUserRole(),
			"tencentcloud_bi_project_user_role":                                                     bi.ResourceTencentCloudBiProjectUserRole(),
			"tencentcloud_bi_datasource":                                                            bi.ResourceTencentCloudBiDatasource(),
			"tencentcloud_bi_datasource_cloud":                                                      bi.ResourceTencentCloudBiDatasourceCloud(),
			"tencentcloud_bi_embed_token_apply":                                                     bi.ResourceTencentCloudBiEmbedTokenApply(),
			"tencentcloud_bi_embed_interval_apply":                                                  bi.ResourceTencentCloudBiEmbedIntervalApply(),
			"tencentcloud_cdwpg_instance":                                                           cdwpg.ResourceTencentCloudCdwpgInstance(),
			"tencentcloud_clickhouse_keyval_config":                                                 cdwch.ResourceTencentCloudClickhouseKeyvalConfig(),
			"tencentcloud_clickhouse_xml_config":                                                    cdwch.ResourceTencentCloudClickhouseXmlConfig(),
			"tencentcloud_clb_target_group_attachments":                                             clb.ResourceTencentCloudClbTargetGroupAttachments(),
			"tencentcloud_vpc_peer_connect_manager":                                                 vpc.ResourceTencentCloudVpcPeerConnectManager(),
			"tencentcloud_vpc_peer_connect_accept_operation":                                        vpc.ResourceTencentCloudVpcPeerConnectAcceptOperation(),
			"tencentcloud_vpc_peer_connect_reject_operation":                                        vpc.ResourceTencentCloudVpcPeerConnectRejectOperation(),
			"tencentcloud_reserve_ip_address":                                                       vpc.ResourceTencentCloudReserveIpAddress(),
			"tencentcloud_elastic_public_ipv6":                                                      vpc.ResourceTencentCloudElasticPublicIpv6(),
			"tencentcloud_classic_elastic_public_ipv6":                                              vpc.ResourceTencentCloudClassicElasticPublicIpv6(),
			"tencentcloud_elastic_public_ipv6_attachment":                                           vpc.ResourceTencentCloudElasticPublicIpv6Attachment(),
			"tencentcloud_csip_risk_center":                                                         csip.ResourceTencentCloudCsipRiskCenter(),
			"tencentcloud_organization_org_share_unit_member":                                       tco.ResourceTencentCloudOrganizationOrgShareUnitMember(),
			"tencentcloud_organization_org_share_unit":                                              tco.ResourceTencentCloudOrganizationOrgShareUnit(),
			"tencentcloud_organization_org_share_unit_resource":                                     tco.ResourceTencentCloudOrganizationOrgShareUnitResource(),
			"tencentcloud_accept_join_share_unit_invitation_operation":                              tco.ResourceTencentCloudAcceptJoinShareUnitInvitationOperation(),
			"tencentcloud_reject_join_share_unit_invitation_operation":                              tco.ResourceTencentCloudRejectJoinShareUnitInvitationOperation(),
			"tencentcloud_kubernetes_addon":                                                         tke.ResourceTencentCloudKubernetesAddon(),
			"tencentcloud_kubernetes_addon_config":                                                  tke.ResourceTencentCloudKubernetesAddonConfig(),
			"tencentcloud_kubernetes_native_node_pool":                                              tke.ResourceTencentCloudKubernetesNativeNodePool(),
			"tencentcloud_cdc_site":                                                                 cdc.ResourceTencentCloudCdcSite(),
			"tencentcloud_cdc_dedicated_cluster":                                                    cdc.ResourceTencentCloudCdcDedicatedCluster(),
			"tencentcloud_cdwdoris_instance":                                                        cdwdoris.ResourceTencentCloudCdwdorisInstance(),
			"tencentcloud_cdwdoris_workload_group":                                                  cdwdoris.ResourceTencentCloudCdwdorisWorkloadGroup(),
			//"tencentcloud_cdwdoris_user":                                       cdwdoris.ResourceTencentCloudCdwdorisUser(),
			"tencentcloud_batch_apply_account_baselines":  controlcenter.ResourceTencentCloudBatchApplyAccountBaselines(),
			"tencentcloud_thpc_workspaces":                thpc.ResourceTencentCloudThpcWorkspaces(),
			"tencentcloud_kubernetes_health_check_policy": tke.ResourceTencentCloudKubernetesHealthCheckPolicy(),
			"tencentcloud_lite_hbase_instance":            emr.ResourceTencentCloudLiteHbaseInstance(),
			"tencentcloud_serverless_hbase_instance":      emr.ResourceTencentCloudServerlessHbaseInstance(),
			"tencentcloud_kubernetes_log_config":          tke.ResourceTencentCloudKubernetesLogConfig(),
			"tencentcloud_tcss_image_registry":            tcss.ResourceTencentCloudTcssImageRegistry(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var getProviderConfig = func(key string) string {
		var str string
		value, err := getConfigFromProfile(d, key)
		if err == nil && value != nil {
			str = value.(string)
		}

		return str
	}

	var (
		secretId            string
		secretKey           string
		securityToken       string
		region              string
		protocol            string
		domain              string
		cosDomain           string
		camRoleName         string
		allowedAccountIds   []string
		forbiddenAccountIds []string
		needSecret          = true
		needAccountFilter   = false
		err                 error
	)

	if v, ok := d.GetOk("secret_id"); ok {
		secretId = v.(string)
	}

	if v, ok := d.GetOk("secret_key"); ok {
		secretKey = v.(string)
	}

	if v, ok := d.GetOk("security_token"); ok {
		securityToken = v.(string)
	}

	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	if secretId == "" && secretKey == "" && securityToken == "" {
		secretId = getProviderConfig("secretId")
		secretKey = getProviderConfig("secretKey")
		securityToken = getProviderConfig("token")
		if region == "" {
			region = getProviderConfig("region")
		}
	}

	if region == "" {
		region = DEFAULT_REGION
	}

	if v, ok := d.GetOk("protocol"); ok {
		protocol = v.(string)
	}

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	if v, ok := d.GetOk("cos_domain"); ok {
		cosDomain = v.(string)
	}

	if v, ok := d.GetOk("cam_role_name"); ok {
		camRoleName = v.(string)
	}

	// standard client
	var tcClient TencentCloudClient
	tcClient.apiV3Conn = &connectivity.TencentCloudClient{
		Credential: sdkcommon.NewTokenCredential(
			secretId,
			secretKey,
			securityToken,
		),
		Region:    region,
		Protocol:  protocol,
		Domain:    domain,
		CosDomain: cosDomain,
	}

	if v, ok := d.GetOk("allowed_account_ids"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			allowedAccountIds = append(allowedAccountIds, v.(string))
		}

		needAccountFilter = true
	}

	if v, ok := d.GetOk("forbidden_account_ids"); ok && v.(*schema.Set).Len() > 0 {
		for _, v := range v.(*schema.Set).List() {
			forbiddenAccountIds = append(forbiddenAccountIds, v.(string))
		}

		needAccountFilter = true
	}

	// get auth from CAM role name
	if camRoleName != "" {
		needSecret = false
		err = genClientWithCAM(&tcClient, camRoleName)
		if err != nil {
			return nil, fmt.Errorf("Get auth from CAM role name failed. Reason: %s", err.Error())
		}
	}

	var (
		assumeRoleArn             string
		assumeRoleSessionName     string
		assumeRoleSessionDuration int
		assumeRolePolicy          string
		assumeRoleExternalId      string
	)

	// get assume role from credential
	if v, ok := providerConfig["role-arn"].(string); ok && v != "" {
		assumeRoleArn = v
	}

	if v, ok := providerConfig["role-session-name"].(string); ok && v != "" {
		assumeRoleSessionName = v
	}

	if assumeRoleArn != "" && assumeRoleSessionName != "" {
		assumeRoleSessionDuration = 7200
		err = genClientWithSTS(&tcClient, assumeRoleArn, assumeRoleSessionName, assumeRoleSessionDuration, assumeRolePolicy, assumeRoleExternalId)
		if err != nil {
			return nil, fmt.Errorf("Get auth from assume role by credential failed. Reason: %s", err.Error())
		}
	}

	// get assume role from env
	envRoleArn := os.Getenv(PROVIDER_ASSUME_ROLE_ARN)
	envSessionName := os.Getenv(PROVIDER_ASSUME_ROLE_SESSION_NAME)
	if envRoleArn != "" && envSessionName != "" {
		if envSessionDuration := os.Getenv(PROVIDER_ASSUME_ROLE_SESSION_DURATION); envSessionDuration != "" {
			var err error
			assumeRoleSessionDuration, err = strconv.Atoi(envSessionDuration)
			if err != nil {
				return nil, err
			}
		}

		if assumeRoleSessionDuration == 0 {
			assumeRoleSessionDuration = 7200
		}

		assumeRoleExternalId = os.Getenv(PROVIDER_ASSUME_ROLE_EXTERNAL_ID)

		// get assume role with saml from env
		envSamlAssertion := os.Getenv(PROVIDER_ASSUME_ROLE_SAML_ASSERTION)
		envPrincipalArn := os.Getenv(PROVIDER_ASSUME_ROLE_PRINCIPAL_ARN)
		// get assume role with web identity from env
		envWebIdentityToken := os.Getenv(PROVIDER_ASSUME_ROLE_WEB_IDENTITY_TOKEN)

		if envSamlAssertion == "" && envPrincipalArn == "" && envWebIdentityToken == "" {
			// use assume role
			err = genClientWithSTS(&tcClient, envRoleArn, envSessionName, assumeRoleSessionDuration, "", assumeRoleExternalId)
			if err != nil {
				return nil, fmt.Errorf("Get auth from assume role by env failed. Reason: %s", err.Error())
			}
		} else if envSamlAssertion != "" && envPrincipalArn != "" && envWebIdentityToken != "" {
			return nil, fmt.Errorf("Can not set `TENCENTCLOUD_ASSUME_ROLE_SAML_ASSERTION`, `TENCENTCLOUD_ASSUME_ROLE_PRINCIPAL_ARN`, `TENCENTCLOUD_ASSUME_ROLE_WEB_IDENTITY_TOKEN` at the same time.\n")
		} else if envSamlAssertion != "" && envPrincipalArn != "" {
			// use assume role with saml
			err = genClientWithSamlSTS(&tcClient, envRoleArn, envSessionName, assumeRoleSessionDuration, envSamlAssertion, envPrincipalArn)
			if err != nil {
				return nil, fmt.Errorf("Get auth from assume role with SAML by env failed. Reason: %s", err.Error())
			}

			needSecret = false
		} else if envWebIdentityToken != "" {
			// use assume role with oidc
			err = genClientWithOidcSTS(&tcClient, envRoleArn, envSessionName, assumeRoleSessionDuration, envWebIdentityToken)
			if err != nil {
				return nil, fmt.Errorf("Get auth from assume role with OIDC by env failed. Reason: %s", err.Error())
			}

			needSecret = false
		} else {
			return nil, fmt.Errorf("Get `assume_role` from env error.\n")
		}
	}

	// get assume role from tf
	if v, ok := d.GetOk("assume_role"); ok {
		assumeRoleList := v.(*schema.Set).List()
		if len(assumeRoleList) == 1 {
			assumeRole := assumeRoleList[0].(map[string]interface{})
			assumeRoleArn = assumeRole["role_arn"].(string)
			assumeRoleSessionName = assumeRole["session_name"].(string)
			assumeRoleSessionDuration = assumeRole["session_duration"].(int)
			assumeRolePolicy = assumeRole["policy"].(string)
			assumeRoleExternalId = assumeRole["external_id"].(string)

			err = genClientWithSTS(&tcClient, assumeRoleArn, assumeRoleSessionName, assumeRoleSessionDuration, assumeRolePolicy, assumeRoleExternalId)
			if err != nil {
				return nil, fmt.Errorf("Get auth from assume role failed. Reason: %s", err.Error())
			}

			if camRoleName != "" {
				needSecret = false
			} else {
				needSecret = true
			}
		}
	}

	var (
		assumeRoleSamlAssertion    string
		assumeRolePrincipalArn     string
		assumeRoleWebIdentityToken string
	)

	// get assume role with saml from tf
	if v, ok := d.GetOk("assume_role_with_saml"); ok {
		assumeRoleWithSamlList := v.([]interface{})
		if len(assumeRoleWithSamlList) == 1 {
			assumeRoleWithSaml := assumeRoleWithSamlList[0].(map[string]interface{})
			assumeRoleSamlAssertion = assumeRoleWithSaml["saml_assertion"].(string)
			assumeRolePrincipalArn = assumeRoleWithSaml["principal_arn"].(string)
			assumeRoleArn = assumeRoleWithSaml["role_arn"].(string)
			assumeRoleSessionName = assumeRoleWithSaml["session_name"].(string)
			assumeRoleSessionDuration = assumeRoleWithSaml["session_duration"].(int)

			err = genClientWithSamlSTS(&tcClient, assumeRoleArn, assumeRoleSessionName, assumeRoleSessionDuration, assumeRoleSamlAssertion, assumeRolePrincipalArn)
			if err != nil {
				return nil, fmt.Errorf("Get auth from assume role with SAML failed. Reason: %s", err.Error())
			}

			needSecret = false
		}
	}

	// get assume role with web identity from tf
	if v, ok := d.GetOk("assume_role_with_web_identity"); ok {
		assumeRoleWithWebIdentityList := v.([]interface{})
		if len(assumeRoleWithWebIdentityList) == 1 {
			assumeRoleWithWebIdentity := assumeRoleWithWebIdentityList[0].(map[string]interface{})
			assumeRoleWebIdentityToken = assumeRoleWithWebIdentity["web_identity_token"].(string)
			assumeRoleArn = assumeRoleWithWebIdentity["role_arn"].(string)
			assumeRoleSessionName = assumeRoleWithWebIdentity["session_name"].(string)
			assumeRoleSessionDuration = assumeRoleWithWebIdentity["session_duration"].(int)

			err = genClientWithOidcSTS(&tcClient, assumeRoleArn, assumeRoleSessionName, assumeRoleSessionDuration, assumeRoleWebIdentityToken)
			if err != nil {
				return nil, fmt.Errorf("Get auth from assume role with OIDC failed. Reason: %s", err.Error())
			}

			needSecret = false
		}
	}

	if v, ok := d.GetOkExists("enable_pod_oidc"); ok && v.(bool) {
		if os.Getenv(POD_OIDC_TKE_REGION) != "" && os.Getenv(POD_OIDC_TKE_WEB_IDENTITY_TOKEN_FILE) != "" && os.Getenv(POD_OIDC_TKE_PROVIDER_ID) != "" && os.Getenv(POD_OIDC_TKE_ROLE_ARN) != "" {
			err := genClientWithPodOidc(&tcClient)
			if err != nil {
				return nil, fmt.Errorf("Get auth from enable pod OIDC failed. Reason: %s", err.Error())
			}

			needSecret = false
		} else {
			return nil, fmt.Errorf("Can not get `TKE_REGION`, `TKE_WEB_IDENTITY_TOKEN_FILE`, `TKE_PROVIDER_ID`, `TKE_ROLE_ARN`. Must config serviceAccountName for pod.\n")
		}
	}

	if needSecret && (secretId == "" || secretKey == "") {
		return nil, fmt.Errorf("Please set your `secret_id` and `secret_key`.\n")
	}

	if needAccountFilter {
		// get indentity
		indentity, err := getCallerIdentity(&tcClient)
		if err != nil {
			return nil, err
		}

		// account filter
		err = verifyAccountIDAllowed(indentity, allowedAccountIds, forbiddenAccountIds)
		if err != nil {
			return nil, err
		}
	}

	return &tcClient, nil
}

func genClientWithCAM(tcClient *TencentCloudClient, roleName string) error {
	var camResp *tccommon.CAMResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := tccommon.GetAuthFromCAM(roleName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("Get cam failed, Response is nil."))
		}

		camResp = result
		return nil
	})

	if err != nil {
		return err
	}

	// using STS credentials
	tcClient.apiV3Conn.Credential = sdkcommon.NewTokenCredential(
		camResp.TmpSecretId,
		camResp.TmpSecretKey,
		camResp.Token,
	)

	return nil
}

func genClientWithSTS(tcClient *TencentCloudClient, assumeRoleArn, assumeRoleSessionName string, assumeRoleSessionDuration int, assumeRolePolicy string, assumeRoleExternalId string) error {
	// applying STS credentials
	request := sdksts.NewAssumeRoleRequest()
	response := sdksts.NewAssumeRoleResponse()
	request.RoleArn = helper.String(assumeRoleArn)
	request.RoleSessionName = helper.String(assumeRoleSessionName)
	request.DurationSeconds = helper.IntUint64(assumeRoleSessionDuration)
	if assumeRolePolicy != "" {
		request.Policy = helper.String(url.QueryEscape(assumeRolePolicy))
	}

	if assumeRoleExternalId != "" {
		request.ExternalId = helper.String(assumeRoleExternalId)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := tcClient.apiV3Conn.UseStsClient().AssumeRole(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Credentials == nil {
			return resource.NonRetryableError(fmt.Errorf("Get Assume Role failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		return err
	}

	if response.Response.Credentials.TmpSecretId == nil || response.Response.Credentials.TmpSecretKey == nil || response.Response.Credentials.Token == nil {
		return fmt.Errorf("Get Assume Role failed, Credentials is nil.")
	}

	// using STS credentials
	tcClient.apiV3Conn.Credential = sdkcommon.NewTokenCredential(
		*response.Response.Credentials.TmpSecretId,
		*response.Response.Credentials.TmpSecretKey,
		*response.Response.Credentials.Token,
	)

	return nil
}

func genClientWithSamlSTS(tcClient *TencentCloudClient, assumeRoleArn, assumeRoleSessionName string, assumeRoleSessionDuration int, assumeRoleSamlAssertion, assumeRolePrincipalArn string) error {
	// applying STS credentials
	request := sdksts.NewAssumeRoleWithSAMLRequest()
	response := sdksts.NewAssumeRoleWithSAMLResponse()
	request.RoleArn = helper.String(assumeRoleArn)
	request.RoleSessionName = helper.String(assumeRoleSessionName)
	request.DurationSeconds = helper.IntUint64(assumeRoleSessionDuration)
	request.SAMLAssertion = helper.String(assumeRoleSamlAssertion)
	request.PrincipalArn = helper.String(assumeRolePrincipalArn)
	var stsExtInfo connectivity.StsExtInfo
	stsExtInfo.Authorization = "SKIP"
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := tcClient.apiV3Conn.UseStsClient(stsExtInfo).AssumeRoleWithSAML(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Credentials == nil {
			return resource.NonRetryableError(fmt.Errorf("Get Assume Role with SAML failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		return err
	}

	if response.Response.Credentials.TmpSecretId == nil || response.Response.Credentials.TmpSecretKey == nil || response.Response.Credentials.Token == nil {
		return fmt.Errorf("Get Assume Role failed, Credentials is nil.")
	}

	// using STS credentials
	tcClient.apiV3Conn.Credential = sdkcommon.NewTokenCredential(
		*response.Response.Credentials.TmpSecretId,
		*response.Response.Credentials.TmpSecretKey,
		*response.Response.Credentials.Token,
	)

	return nil
}

func genClientWithOidcSTS(tcClient *TencentCloudClient, assumeRoleArn, assumeRoleSessionName string, assumeRoleSessionDuration int, assumeRolePolicy string) error {
	// applying STS credentials
	request := sdksts.NewAssumeRoleWithWebIdentityRequest()
	response := sdksts.NewAssumeRoleWithWebIdentityResponse()
	request.ProviderId = helper.String("OIDC")
	request.RoleArn = helper.String(assumeRoleArn)
	request.RoleSessionName = helper.String(assumeRoleSessionName)
	request.DurationSeconds = helper.IntInt64(assumeRoleSessionDuration)
	request.WebIdentityToken = helper.String(assumeRolePolicy)
	var stsExtInfo connectivity.StsExtInfo
	stsExtInfo.Authorization = "SKIP"
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := tcClient.apiV3Conn.UseStsClient(stsExtInfo).AssumeRoleWithWebIdentity(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Credentials == nil {
			return resource.NonRetryableError(fmt.Errorf("Get Assume Role with OIDC failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		return err
	}

	if response.Response.Credentials.TmpSecretId == nil || response.Response.Credentials.TmpSecretKey == nil || response.Response.Credentials.Token == nil {
		return fmt.Errorf("Get Assume Role failed, Credentials is nil.")
	}

	// using STS credentials
	tcClient.apiV3Conn.Credential = sdkcommon.NewTokenCredential(
		*response.Response.Credentials.TmpSecretId,
		*response.Response.Credentials.TmpSecretKey,
		*response.Response.Credentials.Token,
	)

	return nil
}

var providerConfig map[string]interface{}

func getConfigFromProfile(d *schema.ResourceData, ProfileKey string) (interface{}, error) {
	if providerConfig == nil {
		var (
			profile              string
			sharedCredentialsDir string
			credentialPath       string
			configurePath        string
		)

		if v, ok := d.GetOk("profile"); ok {
			profile = v.(string)
		} else {
			profile = DEFAULT_PROFILE
		}

		if v, ok := d.GetOk("shared_credentials_dir"); ok {
			sharedCredentialsDir = v.(string)
		}

		tmpSharedCredentialsDir, err := homedir.Expand(sharedCredentialsDir)
		if err != nil {
			return nil, err
		}

		if tmpSharedCredentialsDir == "" {
			credentialPath = fmt.Sprintf("%s/.tccli/%s.credential", os.Getenv("HOME"), profile)
			configurePath = fmt.Sprintf("%s/.tccli/%s.configure", os.Getenv("HOME"), profile)
			if runtime.GOOS == "windows" {
				credentialPath = fmt.Sprintf("%s/.tccli/%s.credential", os.Getenv("USERPROFILE"), profile)
				configurePath = fmt.Sprintf("%s/.tccli/%s.configure", os.Getenv("USERPROFILE"), profile)
			}
		} else {
			credentialPath = fmt.Sprintf("%s/%s.credential", tmpSharedCredentialsDir, profile)
			configurePath = fmt.Sprintf("%s/%s.configure", tmpSharedCredentialsDir, profile)
		}

		providerConfig = make(map[string]interface{})
		_, err = os.Stat(credentialPath)
		if !os.IsNotExist(err) {
			data, err := os.ReadFile(credentialPath)
			if err != nil {
				return nil, err
			}

			config := map[string]interface{}{}
			err = json.Unmarshal(data, &config)
			if err != nil {
				return nil, err
			}

			for k, v := range config {
				if strValue, ok := v.(string); ok {
					providerConfig[k] = strings.TrimSpace(strValue)
				}
			}
		}

		_, err = os.Stat(configurePath)
		if !os.IsNotExist(err) {
			data, err := os.ReadFile(configurePath)
			if err != nil {
				return nil, err
			}

			config := map[string]interface{}{}
			err = json.Unmarshal(data, &config)
			if err != nil {
				return nil, err
			}

		outerLoop:
			for k, v := range config {
				if k == "_sys_param" {
					tmpMap := v.(map[string]interface{})
					for tmpK, tmpV := range tmpMap {
						if tmpK == "region" {
							providerConfig[tmpK] = strings.TrimSpace(tmpV.(string))
							break outerLoop
						}
					}
				}
			}
		}
	}

	return providerConfig[ProfileKey], nil
}

func genClientWithPodOidc(tcClient *TencentCloudClient) error {
	provider, err := sdkcommon.DefaultTkeOIDCRoleArnProvider()
	if err != nil {
		return err
	}

	assumeResp, err := provider.GetCredential()
	if err != nil {
		return err
	}

	tcClient.apiV3Conn.Credential = sdkcommon.NewTokenCredential(
		assumeResp.GetSecretId(),
		assumeResp.GetSecretKey(),
		assumeResp.GetToken(),
	)

	return nil
}

func getCallerIdentity(tcClient *TencentCloudClient) (indentity *sdksts.GetCallerIdentityResponseParams, err error) {
	ak := tcClient.apiV3Conn.Credential.SecretId
	sk := tcClient.apiV3Conn.Credential.SecretKey
	token := tcClient.apiV3Conn.Credential.Token
	region := tcClient.apiV3Conn.Region
	credential := sdkcommon.NewTokenCredential(ak, sk, token)
	cpf := sdkprofile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sts.tencentcloudapi.com"
	client, _ := sdksts.NewClient(credential, region, cpf)
	request := sdksts.NewGetCallerIdentityRequest()
	response := sdksts.NewGetCallerIdentityResponse()
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := client.GetCallerIdentity(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Get caller identity failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		return
	}

	indentity = response.Response
	return
}

func verifyAccountIDAllowed(indentity *sdksts.GetCallerIdentityResponseParams, allowedAccountIds, forbiddenAccountIds []string) error {
	var accountId string
	if indentity.AccountId != nil {
		accountId = *indentity.AccountId
	} else {
		return fmt.Errorf("Caller identity accountId is illegal, The value is nil.")
	}

	if len(allowedAccountIds) > 0 {
		found := false
		for _, allowedAccountID := range allowedAccountIds {
			if accountId == allowedAccountID {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("TencentCloud account ID not allowed: %s", accountId)
		}
	}

	if len(forbiddenAccountIds) > 0 {
		for _, forbiddenAccountID := range forbiddenAccountIds {
			if accountId == forbiddenAccountID {
				return fmt.Errorf("TencentCloud account ID not allowed: %s", accountId)
			}
		}
	}

	return nil
}
