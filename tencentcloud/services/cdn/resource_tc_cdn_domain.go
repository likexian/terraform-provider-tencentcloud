package cdn

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudCdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdnDomainCreate,
		Read:   resourceTencentCloudCdnDomainRead,
		Update: resourceTencentCloudCdnDomainUpdate,
		Delete: resourceTencentCloudCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			//State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
			//	getDefaultSwitchOffMap := func() []interface{} {
			//		return []interface{}{
			//			map[string]interface{}{"switch": "off"},
			//		}
			//	}
			//	_ = d.Set("authentication", getDefaultSwitchOffMap())
			//	return []*schema.ResourceData{d}, nil
			//},
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"authentication": []interface{}{map[string]interface{}{
					"switch": "off",
				}},
				"cache_key": []interface{}{map[string]interface{}{
					"full_url_cache": "on",
				}},
				"full_url_cache": true,
			}),
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the acceleration domain.",
			},
			"service_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SERVICE_TYPE),
				Description:  "Acceleration domain name service type. `web`: static acceleration, `download`: download acceleration, `media`: streaming media VOD acceleration, `hybrid`: hybrid acceleration, `dynamic`: dynamic acceleration.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The project CDN belongs to, default to 0.",
			},
			"area": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_AREA),
				Description:  "Domain name acceleration region. `mainland`: acceleration inside mainland China, `overseas`: acceleration outside mainland China, `global`: global acceleration. Overseas acceleration service must be enabled to use overseas acceleration and global acceleration.",
			},
			"full_url_cache": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       true,
				ConflictsWith: []string{"cache_key"},
				Deprecated:    "Use `cache_key` -> `full_url_cache` instead.",
				Description:   "Whether to enable full-path cache. Default value is `true`.",
			},
			"origin": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Origin server configuration. It's a list and consist of at most one item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_ORIGIN_TYPE),
							Description: "Master origin server type. The following types are supported: `domain`: Domain name, `domainv6`: IPv6 domain name, `cos`: COS bucket address, `third_party`: Third-party object storage origin, " +
								"`igtm`: IGTM origin, `ip`: IP address, `ipv6`: One IPv6 address, `ip_ipv6`: Multiple IPv4 addresses and one IPv6 address, " +
								"`ip_domain`: IP addresses and domain names (only available to beta users), `ip_domainv6`: Multiple IPv4 addresses and one IPv6 domain name, `ipv6_domain`: Multiple IPv6 addresses and one domain name, `ipv6_domainv6`: Multiple IPv6 addresses and one IPv6 domain name, " +
								"`domain_domainv6`: Multiple IPv4 domain names and one IPv6 domain name, `ip_ipv6_domain`: Multiple IPv4 and IPv6 addresses and one domain name, `ip_ipv6_domainv6`: Multiple IPv4 and IPv6 addresses and one IPv6 domain name, `ip_domain_domainv6`: Multiple IPv4 addresses and IPv4 domain names and one IPv6 domain name, " +
								"`ipv6_domain_domainv6`: Multiple IPv4 domain names and IPv6 addresses and one IPv6 domain name, `ip_ipv6_domain_domainv6`: Multiple IPv4 and IPv6 addresses and IPv4 domain names and one IPv6 domain name.",
						},
						"origin_list": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Master origin server list. Valid values can be ip or domain name. When modifying the origin server, you need to enter the corresponding `origin_type`.",
						},
						"server_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Host header used when accessing the master origin server. If left empty, the acceleration domain name will be used by default.",
						},
						"cos_private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     CDN_SWITCH_OFF,
							Description: "When OriginType is COS, you can specify if access to private buckets is allowed. Valid values are `on` and `off`. and default value is `off`.",
						},
						"origin_pull_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_ORIGIN_PULL_PROTOCOL_HTTP,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_ORIGIN_PULL_PROTOCOL),
							Description:  "Origin-pull protocol configuration. `http`: forced HTTP origin-pull, `follow`: protocol follow origin-pull, `https`: forced HTTPS origin-pull. This only supports origin server port 443 for origin-pull.",
						},
						"backup_origin_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_BACKUP_ORIGIN_TYPE),
							Description:  "Backup origin server type, which supports the following types: `domain`: domain name type, `ip`: IP list used as origin server, `ipv6_domain`: Multiple IPv6 addresses and one domain name, `ip_ipv6`: Multiple IPv4 addresses and one IPv6 address, `ip_ipv6_domain`: Multiple IPv4 and IPv6 addresses and one domain name.",
						},
						"backup_origin_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Backup origin server list. Valid values can be ip or domain name. When modifying the backup origin server, you need to enter the corresponding `backup_origin_type`.",
						},
						"backup_server_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Host header used when accessing the backup origin server. If left empty, the ServerName of master origin server will be used by default.",
						},
						"origin_company": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object storage back to the source vendor. Required when the source station type is a third-party storage source station (third_party). Optional values include the following: `aws_s3`: AWS S3; `ali_oss`: Alibaba Cloud OSS; `hw_obs`: Huawei OBS; `qiniu_kodo`: Qiniu Cloud kodo; `others`: other vendors' object storage, only supports object storage compatible with AWS signature algorithm, such as Tencent Cloud Financial Zone COS. Example value: `hw_obs`.",
						},
					},
				},
			},
			"https_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "HTTPS acceleration configuration. It's a list and consist of at most one item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"https_switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "HTTPS configuration switch. Valid values are `on` and `off`.",
						},
						"http2_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "HTTP2 configuration switch. Valid values are `on` and `off`. and default value is `off`.",
						},
						"ocsp_stapling_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "OCSP configuration switch. Valid values are `on` and `off`. and default value is `off`.",
						},
						"spdy_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Spdy configuration switch. Valid values are `on` and `off`. and default value is `off`. This parameter is for white-list customer.",
						},
						"verify_client": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Client certificate authentication feature. Valid values are `on` and `off`. and default value is `off`.",
						},
						"server_certificate_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Server certificate configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server certificate ID.",
									},
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Server certificate name.",
									},
									"certificate_content": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server certificate information. This is required when uploading an external certificate, which should contain the complete certificate chain.",
									},
									"private_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server key information. This is required when uploading an external certificate.",
									},
									"message": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Certificate remarks.",
									},
									"deploy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy time of server certificate.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire time of server certificate.",
									},
								},
							},
						},
						"client_certificate_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Client certificate configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client certificate name.",
									},
									"certificate_content": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Client Certificate PEM format, requires Base64 encoding.",
									},
									"deploy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy time of client certificate.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire time of client certificate.",
									},
								},
							},
						},
						"force_redirect": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Configuration of forced HTTP or HTTPS redirects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      CDN_SWITCH_OFF,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
										Description:  "Forced redirect configuration switch. Valid values are `on` and `off`. Default value is `off`.",
									},
									"redirect_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      CDN_ORIGIN_PULL_PROTOCOL_HTTP,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_FORCE_REDIRECT_TYPE),
										Description: "Forced redirect type. Valid values are `http` and `https`. `http` means a forced redirect from HTTPS to HTTP, `https` means a forced redirect from HTTP to HTTPS. " +
											"When `switch` setting `off`, this property does not need to be set or set to `http`. Default value is `http`.",
									},
									"redirect_status_code": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      302,
										ValidateFunc: tccommon.ValidateAllowedIntValue([]int{301, 302}),
										Description: "Forced redirect status code. Valid values are `301` and `302`. " +
											"When `switch` setting `off`, this property does not need to be set or set to `302`. Default value is `302`.",
									},
									"carry_headers": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "off",
										Description: "Whether to return the newly added header during force redirection. Values: `on`, `off`.",
									},
								},
							},
						},
						"tls_versions": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Tls version settings, only support some Advanced domain names, support settings TLSv1, TLSv1.1, TLSv1.2, TLSv1.3, when modifying must open consecutive versions.",
						},
					},
				},
			},
			"range_origin_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_ON,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
				Description:  "Sharding back to source configuration switch. Valid values are `on` and `off`. Default value is `on`.",
			},
			"ipv6_access_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_OFF,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
				Description:  "ipv6 access configuration switch. Only available when area set to `mainland`. Valid values are `on` and `off`. Default value is `off`.",
			},
			"follow_redirect_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_OFF,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
				Description:  "301/302 redirect following switch, available values: `on`, `off` (default).",
			},
			"authentication": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specify timestamp hotlink protection configuration, NOTE: only one type can choose for the sub elements.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Authentication switching, available values: `on`, `off`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
						},
						"type_a": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp hotlink protection mode A configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.",
									},
									"sign_param": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Signature parameter name. Only upper and lower-case letters, digits, and underscores (_) are allowed. It cannot start with a digit. Length limit: 1-100 characters.",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Signature expiration time in second. The maximum value is 630720000.",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.",
									},
								},
							},
						},
						"type_b": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp hotlink protection mode B configuration. NOTE: according to upgrading of TencentCloud Platform, TypeB is unavailable for now.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Signature expiration time in second. The maximum value is 630720000.",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.",
									},
								},
							},
						},
						"type_c": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp hotlink protection mode C configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Signature expiration time in second. The maximum value is 630720000.",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.",
									},
									"time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Timestamp formation, available values: `dec`, `hex`.",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.",
									},
								},
							},
						},
						"type_d": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp hotlink protection mode D configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The key for signature calculation. Only digits, upper and lower-case letters are allowed. Length limit: 6-32 characters.",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Signature expiration time in second. The maximum value is 630720000.",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension list settings determining if authentication should be performed. NOTE: If it contains an asterisk (*), this indicates all files.",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Available values: `whitelist` - all types apart from `file_extensions` are authenticated, `blacklist`: - only the types in the `file_extensions` are authenticated.",
									},
									"time_param": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Timestamp parameter name. Only upper and lower-case letters, digits, and underscores (_) are allowed. It cannot start with a digit. Length limit: 1-100 characters.",
									},
									"time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Timestamp formation, available values: `dec`, `hex`.",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Used for calculate a signature. 6-32 characters. Only digits and letters are allowed.",
									},
								},
							},
						},
					},
				},
			},
			"rule_cache": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Advanced path cache configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Matching content under the corresponding type of CacheType: `all`: fill *, `file`: fill in the suffix name, such as jpg, txt, " +
								"`directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html, `index`: fill /.",
						},
						"rule_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_RULE_TYPE_DEFAULT,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_RULE_TYPE),
							Description: "Rule type. The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, " +
								"`directory`: the specified path takes effect, `path`: specify the absolute path to take effect, `index`: home page.",
						},
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Cache configuration switch. Valid values are `on` and `off`.",
						},
						"cache_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Cache expiration time setting, the unit is second, the maximum can be set to 365 days.",
						},
						"compare_max_age": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description: "Advanced cache expiration configuration. When it is turned on, it will compare the max-age value returned by the origin site with the cache expiration time set in CacheRules, " +
								"and take the minimum value to cache at the node. Valid values are `on` and `off`. Default value is `off`.",
						},
						"ignore_cache_control": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description: "Force caching. After opening, the no-store and no-cache resources returned by the origin site will also be cached in accordance with the CacheRules " +
								"rules. Valid values are `on` and `off`. Default value is `off`.",
						},
						"ignore_set_cookie": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Ignore the Set-Cookie header of the origin site. Valid values are `on` and `off`. Default value is `off`. This parameter is for white-list customer.",
						},
						"no_cache_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Cache configuration switch. Valid values are `on` and `off`.",
						},
						"re_validate": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Always check back to origin. Valid values are `on` and `off`. Default value is `off`.",
						},
						"follow_origin_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Follow the source station configuration switch. Valid values are `on` and `off`.",
						},
						"heuristic_cache_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Specify whether to enable heuristic cache, only available while `follow_origin_switch` enabled, values: `on`, `off` (Default).",
						},
						"heuristic_cache_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specify heuristic cache time in second, only available while `follow_origin_switch` and `heuristic_cache_switch` enabled.",
						},
					},
				},
			},
			"request_header": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Request header configuration. It's a list and consist of at most one item.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Custom request header configuration switch. Valid values are `on` and `off`. and default value is `off`.",
						},
						"header_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom request header configuration rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Http header setting method. The following types are supported: `add`: add a head, if a head already exists, there will be a duplicate head, `del`: delete the head.",
									},
									"header_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateStringLengthInRange(1, 100),
										Description:  "Http header name.",
									},
									"header_value": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateStringLengthInRange(1, 1000),
										Description:  "Http header value, optional when Mode is `del`, Required when Mode is `add`/`set`.",
									},
									"rule_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_HEADER_RULE),
										Description: "Rule type. The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, " +
											"`directory`: the specified path takes effect, `path`: specify the absolute path to take effect.",
									},
									"rule_paths": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Description: "Matching content under the corresponding type of CacheType: `all`: fill *, `file`: fill in the suffix name, such as jpg, txt, " +
											"`directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html.",
									},
								},
							},
						},
					},
				},
			},
			// extensions
			"ip_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify Ip filter configurations.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"filter_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP `blacklist`/`whitelist` type.",
						},
						"filters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Ip filter list, Supports IPs in X.X.X.X format, or /8, /16, /24 format IP ranges. Up to 50 allowlists or blocklists can be entered.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"filter_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Ip filter rules, This feature is only available to selected beta customers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Ip filter `blacklist`/`whitelist` type of filter rules.",
									},
									"filters": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Ip filter rule list, supports IPs in X.X.X.X format, or /8, /16, /24 format IP ranges. Up to 50 allowlists or blocklists can be entered.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Ip filter rule type of filter rules, available: `all`, `file`, `directory`, `path`.",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Content list for each `rule_type`: `*` for `all`, file ext like `jpg` for `file`, `/dir/like/` for `directory` and `/path/index.html` for `path`.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"return_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Return code, available values: 400-499.",
						},
					},
				},
			},
			"ip_freq_limit": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Specify Ip frequency limit configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"qps": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Sets the limited number of requests per second, 514 will be returned for requests that exceed the limit.",
						},
					},
				},
			},
			"status_code_cache": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Status code cache configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"cache_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of cache rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status_code": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Code of status cache. available values: `403`, `404`.",
									},
									"cache_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Status code cache expiration time (in seconds).",
									},
								},
							},
						},
					},
				},
			},
			"compression": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Smart compression configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"compression_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of compression rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compress": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Must be set as true, enables compression.",
									},
									"min_length": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The minimum file size to trigger compression (in bytes).",
									},
									"max_length": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The maximum file size to trigger compression (in bytes).",
									},
									"algorithms": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of algorithms, available: `gzip` and `brotli`.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of file extensions like `jpg`, `txt`.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"rule_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Rule type, available: `all`, `file`, `directory`, `path`, `contentType`.",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of rule paths for each `rule_type`: `*` for `all`, file ext like `jpg` for `file`, `/dir/like/` for `directory` and `/path/index.html` for `path`.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"band_width_alert": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Bandwidth cap configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"bps_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "threshold of bps.",
						},
						"counter_measure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Counter measure.",
						},
						"last_trigger_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last trigger time.",
						},
						"alert_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Switch alert.",
						},
						"alert_percentage": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alert percentage.",
						},
						"last_trigger_time_overseas": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last trigger time of overseas.",
						},
						"metric": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Metric.",
						},
						"statistic_item": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Specify statistic item configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration switch, available values: `on`, `off` (default).",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Type of statistic item.",
									},
									"unblock_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Time of auto unblock.",
									},
									"bps_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "threshold of bps.",
									},
									"counter_measure": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Counter measure, values: `RETURN_404`, `RESOLVE_DNS_TO_ORIGIN`.",
									},
									"alert_switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Switch alert.",
									},
									"alert_percentage": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Alert percentage.",
									},
									"metric": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Metric.",
									},
									"cycle": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Cycle of checking in minutes, values `60`, `1440`.",
									},
								},
							},
						},
					},
				},
			},
			"error_page": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Error page configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"page_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of error page rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status_code": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Status code of error page rules.",
									},
									"redirect_code": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Redirect code of error page rules.",
									},
									"redirect_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Redirect url of error page rules.",
									},
								},
							},
						},
					},
				},
			},
			"response_header": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Response header configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"header_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of response header rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Response header mode.",
									},
									"header_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "response header name of rule.",
									},
									"header_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "response header value of rule.",
									},
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "response rule type of rule.",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "response rule paths of rule.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"downstream_capping": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Downstream capping configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"capping_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of capping rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Capping rule type.",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of capping rule path.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"kbps_threshold": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Capping rule kbps threshold.",
									},
								},
							},
						},
					},
				},
			},
			"response_header_cache_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Response header cache switch, available values: `on`, `off` (default).",
			},
			"origin_pull_optimization": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Cross-border linkage optimization configuration. (This feature is in beta and not generally available yet).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"optimization_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optimization type, values: `OVToCN` - Overseas to CN, `CNToOV` CN to Overseas.",
						},
					},
				},
			},
			"seo_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SEO switch, available values: `on`, `off` (default).",
			},
			"referer": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Referer configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"referer_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of referer rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Referer rule type.",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Referer rule path list.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"referer_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Referer type.",
									},
									"referers": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Referer list.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"allow_empty": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether to allow emptpy.",
									},
								},
							},
						},
					},
				},
			},
			"video_seek_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Video seek switch, available values: `on`, `off` (default).",
			},
			"max_age": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Browser cache configuration. (This feature is in beta and not generally available yet).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"max_age_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Max Age rule configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_age_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The following types are supported: `all`: all documents take effect, `file`: the specified file suffix takes effect, `directory`: the specified path takes effect, `path`: specify the absolute path to take effect, `index`: home page.",
									},
									"max_age_contents": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of rule paths for each `max_age_type`: `*` for `all`, file ext like `jpg` for `file`, `/dir/like/` for `directory` and `/path/index.html` for `path`.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"max_age_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Max Age time in seconds, this can set to `0` that stands for no cache.",
									},
									"follow_origin": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to follow origin, values: `on`/`off`, if set to `on`, the `max_age_time` will be ignored.",
									},
								},
							},
						},
					},
				},
			},
			"specific_config_mainland": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specific configuration for mainland, NOTE: Both specifying full schema or using it is superfluous, please use cloud api parameters json passthroughs, check the [Data Types](https://www.tencentcloud.com/document/api/228/31739#MainlandConfig) for more details.",
				DiffSuppressFunc: helper.DiffSupressJSON,
			},
			"specific_config_overseas": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specific configuration for oversea, NOTE: Both specifying full schema or using it is superfluous, please use cloud api parameters json passthroughs, check the [Data Types](https://www.tencentcloud.com/document/api/228/31739#OverseaConfig) for more details.",
				DiffSuppressFunc: helper.DiffSupressJSON,
			},
			"origin_pull_timeout": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Cross-border linkage optimization configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_timeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The origin-pull connection timeout (in seconds). Valid range: 5-60.",
						},
						"receive_timeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The origin-pull receipt timeout (in seconds). Valid range: 10-60.",
						},
					},
				},
			},
			"offline_cache_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Offline cache switch, available values: `on`, `off` (default).",
			},
			"post_max_size": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Maximum post size configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"max_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum size in MB, value range is `[1, 200]`.",
						},
					},
				},
			},
			"quic_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "QUIC switch, available values: `on`, `off` (default).",
			},
			"cache_key": {
				Optional:      true,
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"full_url_cache"},
				Description:   "Cache key configuration (Ignore Query String configuration). NOTE: All of `full_url_cache` default value is `on`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_url_cache": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     CDN_SWITCH_ON,
							Description: "Whether to enable full-path cache, values `on` (DEFAULT ON), `off`.",
						},
						"ignore_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     CDN_SWITCH_OFF,
							Description: "Specifies whether the cache key is case sensitive.",
						},
						"query_string": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Request parameter contained in CacheKey.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     CDN_SWITCH_OFF,
										Description: "Whether to use QueryString as part of CacheKey, values `on`, `off` (Default).",
									},
									"reorder": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     CDN_SWITCH_OFF,
										Description: "Whether to sort again, values `on`, `off` (Default).",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Include/exclude query parameters. Values: `includeAll` (Default), `excludeAll`, `includeCustom`, `excludeCustom`.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Array of included/excluded query strings (separated by `;`).",
									},
								},
							},
						},
						"key_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Path-specific cache key configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_paths": {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Required:    true,
										Description: "List of rule paths for each `key_rules`: `/` for `index`, file ext like `jpg` for `file`, `/dir/like/` for `directory` and `/path/index.html` for `path`.",
									},
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Rule type, available: `file`, `directory`, `path`, `index`.",
									},
									"full_url_cache": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     CDN_SWITCH_ON,
										Description: "Whether to enable full-path cache, values `on` (DEFAULT ON), `off`.",
									},
									"ignore_case": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     CDN_SWITCH_OFF,
										Description: "Whether caches are case insensitive.",
									},
									"query_string": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "Request parameter contained in CacheKey.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     CDN_SWITCH_OFF,
													Description: "Whether to use QueryString as part of CacheKey, values `on`, `off` (Default).",
												},
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specify key rule QS action, values: `includeCustom`, `excludeCustom`.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "",
													Description: "Array of included/excluded query strings (separated by `;`).",
												},
											},
										},
									},
									"rule_tag": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specify rule tag, default value is `user`.",
									},
								},
							},
						},
					},
				},
			},
			"aws_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access authentication for S3 origin.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID.",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key.",
							Sensitive:   true,
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region.",
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bucket.",
						},
					},
				},
			},
			"oss_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access authentication for OSS origin.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID.",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key.",
							Sensitive:   true,
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region.",
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bucket.",
						},
					},
				},
			},
			"hw_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access authentication for OBS origin.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID.",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key.",
							Sensitive:   true,
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bucket.",
						},
					},
				},
			},
			"qn_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access authentication for OBS origin.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID.",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key.",
							Sensitive:   true,
						},
					},
				},
			},
			"others_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Object storage back-to-source authentication of other vendors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch, available values: `on`, `off` (default).",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID.",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key.",
							Sensitive:   true,
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Region.",
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bucket.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of cdn domain.",
			},

			// computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration service status.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME address of domain name.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of domain name.",
			},
			"explicit_using_dry_run": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Used for validate only by store arguments to request json string as expected, WARNING: if set to `true`, NO Cloud Api will be invoked but store as local data, do not use this argument unless you really know what you are doing.",
			},
			"dry_run_create_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used for store `dry_run` request json.",
			},
			"dry_run_update_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used for store `dry_run` update request json.",
			},
		},
	}
}

func resourceTencentCloudCdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cdnService := CdnService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	request := cdn.NewAddCdnDomainRequest()
	domain := d.Get("domain").(string)
	request.Domain = &domain
	request.ServiceType = helper.String(d.Get("service_type").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}
	// Range Origin Pull
	request.RangeOriginPull = &cdn.RangeOriginPull{}
	request.RangeOriginPull.Switch = helper.String(d.Get("range_origin_switch").(string))

	if v, ok := d.GetOk("ipv6_access_switch"); ok {
		request.Ipv6Access = &cdn.Ipv6Access{
			Switch: helper.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("follow_redirect_switch"); ok {
		request.FollowRedirect = &cdn.FollowRedirect{
			Switch: helper.String(v.(string)),
		}
	}

	if v, ok := helper.InterfacesHeadMap(d, "authentication"); ok {
		switchOn := v["switch"].(string)
		request.Authentication = &cdn.Authentication{
			Switch: &switchOn,
		}

		if v, ok := v["type_a"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				signParam      = item["sign_param"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeA = &cdn.AuthenticationTypeA{
				SecretKey:      &secretKey,
				SignParam:      &signParam,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeA.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_b"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeB = &cdn.AuthenticationTypeB{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeB.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_c"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeC = &cdn.AuthenticationTypeC{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if timeFormat, ok := item["time_format"].(string); ok {
				request.Authentication.TypeC.TimeFormat = &timeFormat
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeC.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_d"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
				timeParam      = item["time_param"].(string)
			)

			request.Authentication.TypeD = &cdn.AuthenticationTypeD{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
				TimeParam:      &timeParam,
			}

			if timeFormat, ok := item["time_format"].(string); ok {
				request.Authentication.TypeD.TimeFormat = &timeFormat
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeD.BackupSecretKey = &backupSecretKey
			}
		}
	}

	// rule_cache
	if v, ok := d.GetOk("rule_cache"); ok {
		ruleCache := v.([]interface{})
		var ruleCaches []*cdn.RuleCache
		for _, v := range ruleCache {
			re := &cdn.RuleCache{}
			ruleCacheMap := v.(map[string]interface{})
			rulePaths := ruleCacheMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			ruleType := ruleCacheMap["rule_type"].(string)
			if ruleType == CDN_RULE_TYPE_DEFAULT {
				rulePathList = append(rulePathList, helper.String(CDN_RULE_PATH))
			} else {
				for _, value := range rulePaths {
					rulePathList = append(rulePathList, helper.String(value.(string)))
				}
			}
			switchFlag := ruleCacheMap["switch"].(string)
			cacheTime := ruleCacheMap["cache_time"].(int)
			compareMaxAge := ruleCacheMap["compare_max_age"].(string)
			ignoreCacheControl := ruleCacheMap["ignore_cache_control"].(string)
			ignoreSetCookie := ruleCacheMap["ignore_set_cookie"].(string)
			noCacheSwitch := ruleCacheMap["no_cache_switch"].(string)
			reValidate := ruleCacheMap["re_validate"].(string)
			followOriginSwitch := ruleCacheMap["follow_origin_switch"].(string)
			ruleCacheConfig := &cdn.RuleCacheConfig{}
			cache := &cdn.CacheConfigCache{}
			noCache := &cdn.CacheConfigNoCache{}
			followOrigin := &cdn.CacheConfigFollowOrigin{}
			ruleCacheConfig.Cache = cache
			ruleCacheConfig.NoCache = noCache
			ruleCacheConfig.FollowOrigin = followOrigin
			re.CacheConfig = ruleCacheConfig
			re.RulePaths = rulePathList
			re.RuleType = &ruleType
			re.CacheConfig.Cache.Switch = &switchFlag
			re.CacheConfig.Cache.CacheTime = helper.IntInt64(cacheTime)
			re.CacheConfig.Cache.CompareMaxAge = &compareMaxAge
			re.CacheConfig.Cache.IgnoreCacheControl = &ignoreCacheControl
			re.CacheConfig.Cache.IgnoreSetCookie = &ignoreSetCookie
			re.CacheConfig.NoCache.Switch = &noCacheSwitch
			re.CacheConfig.NoCache.Revalidate = &reValidate
			re.CacheConfig.FollowOrigin.Switch = &followOriginSwitch
			heuristicCacheSwitch := ruleCacheMap["heuristic_cache_switch"].(string)
			heuristicCacheTime := ruleCacheMap["heuristic_cache_time"].(int)
			if heuristicCacheSwitch != "" {
				re.CacheConfig.FollowOrigin.HeuristicCache = &cdn.HeuristicCache{
					Switch:      &heuristicCacheSwitch,
					CacheConfig: &cdn.CacheConfig{},
				}
				if heuristicCacheTime > 0 {
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTimeSwitch = helper.String("on")
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTime = helper.IntInt64(heuristicCacheTime)
				} else {
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTimeSwitch = helper.String("off")
				}
			}
			ruleCaches = append(ruleCaches, re)
		}
		request.Cache = &cdn.Cache{}
		request.Cache.RuleCache = ruleCaches
	}

	if v, ok := d.GetOk("request_header"); ok {
		requestHeaders := v.([]interface{})
		requestHeader := requestHeaders[0].(map[string]interface{})
		headerRule := requestHeader["header_rules"].([]interface{})
		var headerRules []*cdn.HttpHeaderPathRule
		for _, value := range headerRule {
			hr := &cdn.HttpHeaderPathRule{}
			headerRuleMap := value.(map[string]interface{})
			headerMode := headerRuleMap["header_mode"].(string)
			headerName := headerRuleMap["header_name"].(string)
			headerValue := headerRuleMap["header_value"].(string)
			ruleType := headerRuleMap["rule_type"].(string)
			rulePaths := headerRuleMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			for _, value := range rulePaths {
				rulePathList = append(rulePathList, helper.String(value.(string)))
			}
			hr.HeaderMode = &headerMode
			hr.HeaderName = &headerName
			hr.HeaderValue = &headerValue
			hr.RuleType = &ruleType
			hr.RulePaths = rulePathList
			headerRules = append(headerRules, hr)
		}
		request.RequestHeader = &cdn.RequestHeader{}
		request.RequestHeader.Switch = helper.String(requestHeader["switch"].(string))
		request.RequestHeader.HeaderRules = headerRules
	}

	// origin
	origins := d.Get("origin").([]interface{})
	if len(origins) < 1 {
		return fmt.Errorf("origin is required")
	}
	origin := origins[0].(map[string]interface{})
	request.Origin = &cdn.Origin{}
	request.Origin.OriginType = helper.String(origin["origin_type"].(string))
	originList := origin["origin_list"].([]interface{})
	request.Origin.Origins = make([]*string, 0, len(originList))
	for _, item := range originList {
		request.Origin.Origins = append(request.Origin.Origins, helper.String(item.(string)))
	}
	if v := origin["server_name"]; v.(string) != "" {
		request.Origin.ServerName = helper.String(v.(string))
	}
	if v := origin["cos_private_access"]; v.(string) != "" {
		request.Origin.CosPrivateAccess = helper.String(v.(string))
	}
	if v := origin["origin_pull_protocol"]; v.(string) != "" {
		request.Origin.OriginPullProtocol = helper.String(v.(string))
	}
	if v := origin["backup_origin_type"]; v.(string) != "" {
		request.Origin.BackupOriginType = helper.String(v.(string))
	}
	if v := origin["backup_server_name"]; v.(string) != "" {
		request.Origin.BackupServerName = helper.String(v.(string))
	}
	if v := origin["backup_origin_list"]; len(v.([]interface{})) > 0 {
		backupOriginList := v.([]interface{})
		request.Origin.BackupOrigins = make([]*string, 0, len(backupOriginList))
		for _, item := range backupOriginList {
			request.Origin.BackupOrigins = append(request.Origin.BackupOrigins, helper.String(item.(string)))
		}
	}
	if v := origin["origin_company"]; v.(string) != "" {
		request.Origin.OriginCompany = helper.String(v.(string))
	}

	// https config
	if v, ok := d.GetOk("https_config"); ok {
		httpsConfigs := v.([]interface{})
		if len(httpsConfigs) > 0 {
			config := httpsConfigs[0].(map[string]interface{})
			request.Https = &cdn.Https{}
			request.Https.Switch = helper.String(config["https_switch"].(string))
			if v := config["http2_switch"]; v.(string) != "" {
				request.Https.Http2 = helper.String(v.(string))
			}
			request.Https.OcspStapling = helper.String(config["ocsp_stapling_switch"].(string))
			request.Https.Spdy = helper.String(config["spdy_switch"].(string))
			request.Https.VerifyClient = helper.String(config["verify_client"].(string))
			if v := config["server_certificate_config"]; len(v.([]interface{})) > 0 {
				serverCerts := v.([]interface{})
				if len(serverCerts) > 0 && serverCerts[0] != nil {
					serverCert := serverCerts[0].(map[string]interface{})
					request.Https.CertInfo = &cdn.ServerCert{}
					if v := serverCert["certificate_id"]; v.(string) != "" {
						request.Https.CertInfo.CertId = helper.String(v.(string))
					}
					if v := serverCert["certificate_content"]; v.(string) != "" {
						request.Https.CertInfo.Certificate = helper.String(v.(string))
					}
					if v := serverCert["private_key"]; v.(string) != "" {
						request.Https.CertInfo.PrivateKey = helper.String(v.(string))
					}
					if v := serverCert["message"]; v.(string) != "" {
						request.Https.CertInfo.Message = helper.String(v.(string))
					}
				}
			}
			if v := config["client_certificate_config"]; len(v.([]interface{})) > 0 {
				clientCerts := v.([]interface{})
				if len(clientCerts) > 0 && clientCerts[0] != nil {
					clientCert := clientCerts[0].(map[string]interface{})
					request.Https.ClientCertInfo = &cdn.ClientCert{}
					if v := clientCert["certificate_content"]; v.(string) != "" {
						request.Https.ClientCertInfo.Certificate = helper.String(v.(string))
					}
				}
			}
			if v, ok := config["force_redirect"]; ok {
				forceRedirect := v.([]interface{})
				if len(forceRedirect) > 0 && forceRedirect[0] != nil {
					var redirect cdn.ForceRedirect
					redirectMap := forceRedirect[0].(map[string]interface{})
					if sw := redirectMap["switch"]; sw.(string) != "" {
						redirect.Switch = helper.String(sw.(string))
					}
					if rt := redirectMap["redirect_type"]; rt.(string) != "" {
						redirect.RedirectType = helper.String(rt.(string))
					}
					if rsc := redirectMap["redirect_status_code"]; rsc.(int) != 0 {
						redirect.RedirectStatusCode = helper.Int64(int64(rsc.(int)))
					}
					if ch := redirectMap["carry_headers"]; ch.(string) != "" {
						redirect.CarryHeaders = helper.String(ch.(string))
					}
					request.ForceRedirect = &redirect
				}
			}
			if v, ok := config["tls_versions"]; ok {
				request.Https.TlsVersion = helper.InterfacesStringsPoint(v.([]interface{}))
			}
		}
	}

	// more added
	if v, ok := helper.InterfacesHeadMap(d, "ip_filter"); ok {
		vSwitch := v["switch"].(string)
		request.IpFilter = &cdn.IpFilter{
			Switch: &vSwitch,
		}
		if vv, ok := v["filter_type"].(string); ok {
			request.IpFilter.FilterType = &vv
		}
		if vv, ok := v["filters"].([]interface{}); ok {
			request.IpFilter.Filters = helper.InterfacesStringsPoint(vv)
		}

		//need white list func
		if vv, ok := v["filter_rules"].([]interface{}); ok && len(vv) > 0 {
			filterRules := make([]*cdn.IpFilterPathRule, 0)
			for i := range vv {
				item := vv[i].(map[string]interface{})
				rule := &cdn.IpFilterPathRule{}
				if rv, ok := item["filter_type"].(string); ok && rv != "" {
					rule.FilterType = &rv
				}
				if rv, ok := item["filters"].([]interface{}); ok && len(rv) > 0 {
					rule.Filters = helper.InterfacesStringsPoint(rv)
				}
				if rv, ok := item["rule_type"].(string); ok && rv != "" {
					rule.RuleType = &rv
				}
				if rv, ok := item["rule_paths"].([]interface{}); ok && len(rv) > 0 {
					rule.RulePaths = helper.InterfacesStringsPoint(rv)
				}
				filterRules = append(filterRules, rule)
			}
			request.IpFilter.FilterRules = filterRules
		}

		if vv, ok := v["return_code"].(int); ok {
			request.IpFilter.ReturnCode = helper.IntInt64(vv)
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "ip_freq_limit"); ok {
		vSwitch := v["switch"].(string)
		qps := v["qps"].(int)
		request.IpFreqLimit = &cdn.IpFreqLimit{
			Switch: &vSwitch,
			Qps:    helper.IntInt64(qps),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "status_code_cache"); ok {
		vSwitch := v["switch"].(string)
		request.StatusCodeCache = &cdn.StatusCodeCache{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.StatusCodeCacheRule, 0)
		rules := v["cache_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.StatusCodeCacheRule{
				StatusCode: helper.String(item["status_code"].(string)),
			}
			if v, ok := item["cache_time"].(int); ok && v > 0 {
				rule.CacheTime = helper.IntInt64(v)
			}
			requestRules = append(requestRules, rule)
		}
		request.StatusCodeCache.CacheRules = requestRules
	}
	if v, ok := helper.InterfacesHeadMap(d, "compression"); ok {
		vSwitch := v["switch"].(string)
		request.Compression = &cdn.Compression{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.CompressionRule, 0)
		rules := v["compression_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			var (
				compress = item["compress"].(bool)
			)
			rule := &cdn.CompressionRule{
				Compress: &compress,
			}
			if v, ok := item["min_length"].(int); ok && v > 0 {
				rule.MinLength = helper.IntInt64(v)
			}
			if v, ok := item["max_length"].(int); ok && v > 0 {
				rule.MaxLength = helper.IntInt64(v)
			}
			if v, ok := item["algorithms"].([]interface{}); ok && len(v) > 0 {
				rule.Algorithms = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["file_extensions"].([]interface{}); ok && len(v) > 0 {
				rule.FileExtensions = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}

			requestRules = append(requestRules, rule)
		}
		request.Compression.CompressionRules = requestRules

	}
	if v, ok := helper.InterfacesHeadMap(d, "band_width_alert"); ok {
		vSwitch := v["switch"].(string)
		request.BandwidthAlert = &cdn.BandwidthAlert{
			Switch: &vSwitch,
		}
		if v, ok := v["bps_threshold"].(int); ok && v > 0 {
			request.BandwidthAlert.BpsThreshold = helper.IntInt64(v)
		}
		if v, ok := v["counter_measure"].(string); ok && v != "" {
			request.BandwidthAlert.CounterMeasure = &v
		}
		//if v, ok := v["last_trigger_time"].(string); ok && v != "" {
		//	request.BandwidthAlert.LastTriggerTime = &v
		//}
		if v, ok := v["alert_switch"].(string); ok && v != "" {
			request.BandwidthAlert.AlertSwitch = &v
		}
		if v, ok := v["alert_percentage"].(int); ok && v > 0 {
			request.BandwidthAlert.AlertPercentage = helper.IntInt64(v)
		}
		//if v, ok := v["last_trigger_time_overseas"].(string); ok && v != "" {
		//	request.BandwidthAlert.LastTriggerTimeOverseas = &v
		//}
		if v, ok := v["metric"].(string); ok && v != "" {
			request.BandwidthAlert.Metric = &v
		}
		if statistic, ok := v["statistic_item"].([]interface{}); ok && len(statistic) > 0 {
			for i := range statistic {
				item := statistic[i].(map[string]interface{})
				vSwitch := item["switch"].(string)
				sItem := &cdn.StatisticItem{
					Switch: &vSwitch,
				}
				if vv, ok := item["type"].(string); ok && vv != "" {
					sItem.Type = &vv
				}
				if vv, ok := item["unblock_time"].(int); ok && vv != 0 {
					sItem.UnBlockTime = helper.IntUint64(vv)
				}
				if vv, ok := item["bps_threshold"].(int); ok && vv != 0 {
					sItem.BpsThreshold = helper.IntUint64(vv)
				}
				if vv, ok := item["counter_measure"].(string); ok && vv != "" {
					sItem.CounterMeasure = &vv
				}
				if vv, ok := item["alert_switch"].(string); ok && vv != "" {
					sItem.AlertSwitch = &vv
				}
				if vv, ok := item["alert_percentage"].(int); ok && vv != 0 {
					sItem.AlertPercentage = helper.IntUint64(vv)
				}
				if vv, ok := item["metric"].(string); ok && vv != "" {
					sItem.Metric = &vv
				}
				if vv, ok := item["cycle"].(int); ok && vv != 0 {
					sItem.BpsThreshold = helper.IntUint64(vv)
				}
				request.BandwidthAlert.StatisticItems = append(request.BandwidthAlert.StatisticItems, sItem)
			}

		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "error_page"); ok {
		vSwitch := v["switch"].(string)
		request.ErrorPage = &cdn.ErrorPage{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.ErrorPageRule, 0)
		rules := v["page_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.ErrorPageRule{}
			if v, ok := item["status_code"].(int); ok && v != 0 {
				rule.StatusCode = helper.IntInt64(v)
			}
			if v, ok := item["redirect_code"].(int); ok && v != 0 {
				rule.RedirectCode = helper.IntInt64(v)
			}
			if v, ok := item["redirect_url"].(string); ok && v != "" {
				rule.RedirectUrl = &v
			}
			requestRules = append(requestRules, rule)
		}
		request.ErrorPage.PageRules = requestRules
	}
	if v, ok := helper.InterfacesHeadMap(d, "response_header"); ok {
		vSwitch := v["switch"].(string)
		request.ResponseHeader = &cdn.ResponseHeader{
			Switch: &vSwitch,
		}
		responseRules := make([]*cdn.HttpHeaderPathRule, 0)
		rules := v["header_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.HttpHeaderPathRule{}
			if v, ok := item["header_mode"].(string); ok && v != "" {
				rule.HeaderMode = &v
			}
			if v, ok := item["header_name"].(string); ok && v != "" {
				rule.HeaderName = &v
			}
			if v, ok := item["header_value"].(string); ok && v != "" {
				rule.HeaderValue = &v
			}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			responseRules = append(responseRules, rule)
		}
		request.ResponseHeader.HeaderRules = responseRules
	}
	if v, ok := helper.InterfacesHeadMap(d, "downstream_capping"); ok {
		vSwitch := v["switch"].(string)
		request.DownstreamCapping = &cdn.DownstreamCapping{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.CappingRule, 0)
		rules := v["capping_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.CappingRule{}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["kbps_threshold"].(int); ok && v > 0 {
				rule.KBpsThreshold = helper.IntInt64(v)
			}
			requestRules = append(requestRules, rule)
		}
		request.DownstreamCapping.CappingRules = requestRules
	}
	if v, ok := d.GetOk("response_header_cache_switch"); ok {
		request.ResponseHeaderCache = &cdn.ResponseHeaderCache{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "origin_pull_optimization"); ok {
		vSwitch := v["switch"].(string)
		request.OriginPullOptimization = &cdn.OriginPullOptimization{
			Switch: &vSwitch,
		}
		if v, ok := v["optimization_type"].(string); ok && v != "" {
			request.OriginPullOptimization.OptimizationType = &v
		}
	}
	if v, ok := d.GetOk("seo_switch"); ok {
		request.Seo = &cdn.Seo{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "referer"); ok {
		vSwitch := v["switch"].(string)
		request.Referer = &cdn.Referer{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.RefererRule, 0)
		rules := v["referer_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.RefererRule{}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["referer_type"].(string); ok && v != "" {
				rule.RefererType = &v
			}
			if v, ok := item["referers"].([]interface{}); ok && len(v) > 0 {
				rule.Referers = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["allow_empty"].(bool); ok {
				rule.AllowEmpty = &v
			}
			requestRules = append(requestRules, rule)
		}
		request.Referer.RefererRules = requestRules
	}
	if v, ok := d.GetOk("video_seek_switch"); ok {
		request.VideoSeek = &cdn.VideoSeek{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "max_age"); ok {
		vSwitch := v["switch"].(string)
		request.MaxAge = &cdn.MaxAge{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.MaxAgeRule, 0)
		rules := v["max_age_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.MaxAgeRule{}

			if v, ok := item["max_age_type"].(string); ok && v != "" {
				rule.MaxAgeType = &v
			}
			if v, ok := item["max_age_contents"].([]interface{}); ok && len(v) > 0 {
				rule.MaxAgeContents = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["max_age_time"].(int); ok && v > 0 {
				rule.MaxAgeTime = helper.IntInt64(v)
			}
			if v, ok := item["follow_origin"].(string); ok && v != "" {
				rule.FollowOrigin = &v
			}

			requestRules = append(requestRules, rule)
		}
		request.MaxAge.MaxAgeRules = requestRules
	}
	if v, ok := d.GetOk("specific_config_mainland"); ok && v.(string) != "" {
		request.SpecificConfig = &cdn.SpecificConfig{}
		configStruct := cdn.MainlandConfig{}
		err := json.Unmarshal([]byte(v.(string)), &configStruct)
		if err != nil {
			return fmt.Errorf("unmarshal specific_config_mainland fail: %s", err.Error())
		}
		request.SpecificConfig.Mainland = &configStruct
	}
	if v, ok := d.GetOk("specific_config_overseas"); ok && v.(string) != "" {
		if request.SpecificConfig == nil {
			request.SpecificConfig = &cdn.SpecificConfig{}
		}
		configStruct := cdn.OverseaConfig{}
		err := json.Unmarshal([]byte(v.(string)), &configStruct)
		if err != nil {
			return fmt.Errorf("unmarshal specific_config_overseas fail: %s", err.Error())
		}
		request.SpecificConfig.Overseas = &configStruct
	}
	if v, ok := helper.InterfacesHeadMap(d, "origin_pull_timeout"); ok {
		request.OriginPullTimeout = &cdn.OriginPullTimeout{}
		if vv, ok := v["connect_timeout"].(int); ok && vv > 0 {
			request.OriginPullTimeout.ConnectTimeout = helper.IntUint64(vv)
		}
		if vv, ok := v["receive_timeout"].(int); ok && vv > 0 {
			request.OriginPullTimeout.ReceiveTimeout = helper.IntUint64(vv)
		}
	}
	if v, ok := d.GetOk("offline_cache_switch"); ok {
		request.OfflineCache = &cdn.OfflineCache{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := d.GetOk("quic_switch"); ok {
		request.Quic = &cdn.Quic{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "cache_key"); ok {
		request.CacheKey = &cdn.CacheKey{}
		if fuc := v["full_url_cache"].(string); fuc != "" {
			request.CacheKey.FullUrlCache = &fuc
		}
		if ic := v["ignore_case"].(string); ic != "" {
			request.CacheKey.IgnoreCase = &ic
		}
		if qs, ok := v["query_string"].([]interface{}); ok && len(qs) > 0 {
			if dMap, ok := qs[0].(map[string]interface{}); ok {
				qSwitch := dMap["switch"].(string)
				reorder := dMap["reorder"].(string)
				action := dMap["action"].(string)
				value := dMap["value"].(string)
				request.CacheKey.QueryString = &cdn.QueryStringKey{
					Switch:  &qSwitch,
					Reorder: &reorder,
					Action:  &action,
					Value:   &value,
				}
			}
		}
		if kr, ok := v["key_rules"].([]interface{}); ok {
			for i := range kr {
				rule, ok := kr[i].(map[string]interface{})
				if !ok {
					continue
				}
				ruleType := rule["rule_type"].(string)
				keyRule := &cdn.KeyRule{
					RuleType: &ruleType,
				}
				if vv := rule["full_url_cache"].(string); vv != "" {
					keyRule.FullUrlCache = &vv
				}
				if vv := rule["ignore_case"].(string); vv != "" {
					keyRule.IgnoreCase = &vv
				}
				if vv := rule["rule_tag"].(string); vv != "" {
					keyRule.RuleTag = &vv
				}
				if rp, ok := rule["rule_paths"].([]interface{}); ok {
					keyRule.RulePaths = helper.InterfacesStringsPoint(rp)
				}
				if qs, ok := rule["query_string"].([]interface{}); ok && len(qs) > 0 {
					if dMap, ok := qs[0].(map[string]interface{}); ok {
						vSwitch := dMap["switch"].(string)
						keyRule.QueryString = &cdn.RuleQueryString{
							Switch: &vSwitch,
						}
						if v := dMap["action"].(string); v != "" && vSwitch == "on" {
							keyRule.QueryString.Action = &v
						}
						if v := dMap["value"].(string); v != "" {
							keyRule.QueryString.Value = &v
						}
					}
				}
				request.CacheKey.KeyRules = append(request.CacheKey.KeyRules, keyRule)
			}
		}
	} else {
		fullUrlCache := d.Get("full_url_cache").(bool)
		request.CacheKey = &cdn.CacheKey{}
		if fullUrlCache {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_ON)
		} else {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_OFF)
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "aws_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.AwsPrivateAccess = &cdn.AwsPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.AwsPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.AwsPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.AwsPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.AwsPrivateAccess.Bucket = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "oss_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.OssPrivateAccess = &cdn.OssPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.OssPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.OssPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.OssPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.OssPrivateAccess.Bucket = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "hw_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.HwPrivateAccess = &cdn.HwPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.HwPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.HwPrivateAccess.SecretKey = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.HwPrivateAccess.Bucket = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "qn_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.QnPrivateAccess = &cdn.QnPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.QnPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.QnPrivateAccess.SecretKey = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "others_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.OthersPrivateAccess = &cdn.OthersPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.OthersPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.OthersPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.OthersPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.OthersPrivateAccess.Bucket = &v
		}
	}

	if v := d.Get("explicit_using_dry_run").(bool); v {
		d.SetId(domain)
		_ = d.Set("dry_run_create_result", request.ToJsonString())
		return nil
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().AddCdnDomain(request)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROR || sdkErr.Code == CDN_HOST_EXISTS {
					return resource.NonRetryableError(err)
				}
			}
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(domain)

	time.Sleep(1 * time.Second)
	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
			return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := updateCdnModifyOnlyParams(d, meta, ctx); err != nil {
		return err
	}

	// tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		region := client.Region
		resourceName := tccommon.BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		err := tagService.ModifyTags(ctx, resourceName, tags, nil)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudCdnDomainRead(d, meta)
}

func resourceTencentCloudCdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	region := client.Region
	cdnService := CdnService{client: client}
	tagService := svctag.NewTagService(client)

	domain := d.Id()

	if v, ok := d.Get("explicit_using_dry_run").(bool); ok && v {
		d.SetId(domain)
		return nil
	}

	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if domainConfig == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("service_type", domainConfig.ServiceType)
	_ = d.Set("project_id", domainConfig.ProjectId)
	_ = d.Set("area", domainConfig.Area)
	_ = d.Set("status", domainConfig.Status)
	_ = d.Set("create_time", domainConfig.CreateTime)
	_ = d.Set("cname", domainConfig.Cname)
	_ = d.Set("range_origin_switch", domainConfig.RangeOriginPull.Switch)

	if domainConfig.Ipv6Access != nil {
		_ = d.Set("ipv6_access_switch", domainConfig.Ipv6Access.Switch)
	}
	if domainConfig.FollowRedirect != nil {
		_ = d.Set("follow_redirect_switch", domainConfig.FollowRedirect.Switch)
	}

	origins := make([]map[string]interface{}, 0, 1)
	origin := make(map[string]interface{}, 8)
	origin["origin_type"] = domainConfig.Origin.OriginType
	origin["origin_list"] = domainConfig.Origin.Origins
	origin["server_name"] = domainConfig.Origin.ServerName
	origin["cos_private_access"] = domainConfig.Origin.CosPrivateAccess
	origin["origin_pull_protocol"] = domainConfig.Origin.OriginPullProtocol
	origin["backup_origin_type"] = domainConfig.Origin.BackupOriginType
	origin["backup_origin_list"] = domainConfig.Origin.BackupOrigins
	origin["backup_server_name"] = domainConfig.Origin.BackupServerName
	origin["origin_company"] = domainConfig.Origin.OriginCompany
	origins = append(origins, origin)
	_ = d.Set("origin", origins)

	if len(domainConfig.Cache.RuleCache) > 0 {
		ruleCaches := make([]map[string]interface{}, len(domainConfig.Cache.RuleCache))
		for index, value := range domainConfig.Cache.RuleCache {
			ruleCache := make(map[string]interface{})
			ruleCache["rule_paths"] = value.RulePaths
			ruleCache["rule_type"] = value.RuleType
			if value.CacheConfig == nil {
				continue
			}
			if value.CacheConfig.Cache != nil {
				ruleCache["switch"] = value.CacheConfig.Cache.Switch
				ruleCache["cache_time"] = value.CacheConfig.Cache.CacheTime
				ruleCache["compare_max_age"] = value.CacheConfig.Cache.CompareMaxAge
				ruleCache["ignore_cache_control"] = value.CacheConfig.Cache.IgnoreCacheControl
				ruleCache["ignore_set_cookie"] = value.CacheConfig.Cache.IgnoreSetCookie
			}
			if value.CacheConfig.NoCache != nil {
				ruleCache["no_cache_switch"] = value.CacheConfig.NoCache.Switch
				ruleCache["re_validate"] = value.CacheConfig.NoCache.Revalidate
			}
			if value.CacheConfig.FollowOrigin != nil {
				ruleCache["follow_origin_switch"] = value.CacheConfig.FollowOrigin.Switch
				if htc := value.CacheConfig.FollowOrigin.HeuristicCache; htc != nil {
					ruleCache["heuristic_cache_switch"] = htc.Switch
					if htc.CacheConfig != nil {
						ruleCache["heuristic_cache_time"] = htc.CacheConfig.HeuristicCacheTime
					}
				}
			}
			ruleCaches[index] = ruleCache
		}
		_ = d.Set("rule_cache", ruleCaches)
	}

	requestHeaders := make([]map[string]interface{}, 1)
	requestHeader := make(map[string]interface{})
	if domainConfig.RequestHeader != nil {
		requestHeader["switch"] = domainConfig.RequestHeader.Switch
		if len(domainConfig.RequestHeader.HeaderRules) > 0 {
			headerRules := make([]map[string]interface{}, len(domainConfig.RequestHeader.HeaderRules))
			headerRuleList := domainConfig.RequestHeader.HeaderRules
			for index, value := range headerRuleList {
				headerRule := make(map[string]interface{})
				headerRule["header_mode"] = value.HeaderMode
				headerRule["header_name"] = value.HeaderName
				headerRule["header_value"] = value.HeaderValue
				headerRule["rule_type"] = value.RuleType
				headerRule["rule_paths"] = value.RulePaths
				headerRules[index] = headerRule
			}
			requestHeader["header_rules"] = headerRules
		}
		requestHeaders[0] = requestHeader
		_ = d.Set("request_header", requestHeaders)
	}

	httpsConfigs := make([]map[string]interface{}, 0, 1)
	httpsConfig := make(map[string]interface{}, 8)
	httpsConfig["https_switch"] = domainConfig.Https.Switch
	httpsConfig["http2_switch"] = domainConfig.Https.Http2
	httpsConfig["ocsp_stapling_switch"] = domainConfig.Https.OcspStapling
	httpsConfig["spdy_switch"] = domainConfig.Https.Spdy
	httpsConfig["verify_client"] = domainConfig.Https.VerifyClient

	oldHttpsConfigs := make([]interface{}, 0)
	if _, ok := d.GetOk("https_config"); ok {
		oldHttpsConfigs = d.Get("https_config").([]interface{})
	}
	oldHttpsConfig := make(map[string]interface{})
	if len(oldHttpsConfigs) > 0 && oldHttpsConfigs[0] != nil {
		oldHttpsConfig = oldHttpsConfigs[0].(map[string]interface{})
	}
	oldServerConfigs := make([]interface{}, 0)
	if _, ok := oldHttpsConfig["server_certificate_config"]; ok {
		oldServerConfigs = oldHttpsConfig["server_certificate_config"].([]interface{})
	}
	oldServerConfig := make(map[string]interface{})
	if len(oldServerConfigs) > 0 && oldServerConfigs[0] != nil {
		oldServerConfig = oldServerConfigs[0].(map[string]interface{})
	}
	oldClientConfigs := make([]interface{}, 0)
	if _, ok := oldHttpsConfig["client_certificate_config"]; ok {
		oldClientConfigs = oldHttpsConfig["client_certificate_config"].([]interface{})
	}
	oldClientConfig := make(map[string]interface{})
	if len(oldClientConfigs) > 0 && oldClientConfigs[0] != nil {
		oldClientConfig = oldClientConfigs[0].(map[string]interface{})
	}

	if domainConfig.Https.CertInfo != nil && domainConfig.Https.CertInfo.CertName != nil {
		serverCertConfigs := make([]map[string]interface{}, 0, 1)
		serverCertConfig := make(map[string]interface{}, 5)
		serverCertConfig["certificate_id"] = domainConfig.Https.CertInfo.CertId
		serverCertConfig["certificate_name"] = domainConfig.Https.CertInfo.CertName
		serverCertConfig["certificate_content"] = oldServerConfig["certificate_content"]
		serverCertConfig["private_key"] = oldServerConfig["private_key"]
		serverCertConfig["message"] = domainConfig.Https.CertInfo.Message
		serverCertConfig["deploy_time"] = domainConfig.Https.CertInfo.DeployTime
		serverCertConfig["expire_time"] = domainConfig.Https.CertInfo.ExpireTime
		serverCertConfigs = append(serverCertConfigs, serverCertConfig)
		httpsConfig["server_certificate_config"] = serverCertConfigs
	}
	if domainConfig.Https.ClientCertInfo != nil && domainConfig.Https.ClientCertInfo.CertName != nil {
		clientCertConfigs := make([]map[string]interface{}, 0, 1)
		clientCertConfig := make(map[string]interface{}, 2)
		clientCertConfig["certificate_content"] = oldClientConfig["certificate_content"]
		clientCertConfig["certificate_name"] = domainConfig.Https.ClientCertInfo.CertName
		clientCertConfig["deploy_time"] = domainConfig.Https.ClientCertInfo.DeployTime
		clientCertConfig["expire_time"] = domainConfig.Https.ClientCertInfo.ExpireTime
		clientCertConfigs = append(clientCertConfigs, clientCertConfig)
		httpsConfig["client_certificate_config"] = clientCertConfigs
	}
	if domainConfig.ForceRedirect != nil {
		httpsConfig["force_redirect"] = []map[string]interface{}{
			{
				"switch":               domainConfig.ForceRedirect.Switch,
				"redirect_type":        domainConfig.ForceRedirect.RedirectType,
				"redirect_status_code": domainConfig.ForceRedirect.RedirectStatusCode,
				"carry_headers":        domainConfig.ForceRedirect.CarryHeaders,
			},
		}
	}
	if len(domainConfig.Https.TlsVersion) > 0 {
		tlsVersions := make([]string, 0)
		for _, tlsVersionItem := range domainConfig.Https.TlsVersion {
			tlsVersions = append(tlsVersions, *tlsVersionItem)
		}
		httpsConfig["tls_versions"] = tlsVersions
	}
	httpsConfigs = append(httpsConfigs, httpsConfig)
	_ = d.Set("https_config", httpsConfigs)

	authRaw := d.Get("authentication").([]interface{})
	if authentication := domainConfig.Authentication; authentication != nil && len(authRaw) > 0 {
		auth := make(map[string]interface{})
		auth["switch"] = authentication.Switch
		if authType := authentication.TypeA; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"sign_param":        authType.SignParam,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_a"] = []interface{}{dMap}
		}
		if authType := authentication.TypeB; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_b"] = []interface{}{dMap}
		}
		if authType := authentication.TypeC; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"time_format":       authType.TimeFormat,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_c"] = []interface{}{dMap}
		}
		if authType := authentication.TypeD; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"time_param":        authType.TimeParam,
				"time_format":       authType.TimeFormat,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_d"] = []interface{}{dMap}
		}
		_ = d.Set("authentication", []interface{}{auth})
	}

	dc := domainConfig

	if ok := checkCdnInfoWritable(d, "ip_filter", dc.IpFilter); ok {
		dMap := map[string]interface{}{
			"switch":      dc.IpFilter.Switch,
			"filter_type": dc.IpFilter.FilterType,
			"filters":     dc.IpFilter.Filters,
			"return_code": dc.IpFilter.ReturnCode,
		}
		if rules := dc.IpFilter.FilterRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"filter_type": item.FilterType,
					"filters":     item.Filters,
					"rule_type":   item.RuleType,
					"rule_paths":  item.RulePaths,
				}
				list = append(list, rule)
			}
			dMap["filter_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "ip_filter", dMap)
	}
	if ok := checkCdnInfoWritable(d, "ip_freq_limit", dc.IpFreqLimit); ok {
		dMap := map[string]interface{}{
			"switch": dc.IpFreqLimit.Switch,
			"qps":    dc.IpFreqLimit.Qps,
		}
		_ = helper.SetMapInterfaces(d, "ip_freq_limit", dMap)
	}
	if ok := checkCdnInfoWritable(d, "status_code_cache", dc.StatusCodeCache); ok {
		dMap := map[string]interface{}{
			"switch": dc.StatusCodeCache.Switch,
		}
		if rules := dc.StatusCodeCache.CacheRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"status_code": item.StatusCode,
					"cache_time":  item.CacheTime,
				}
				list = append(list, rule)
			}
			dMap["cache_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "status_code_cache", dMap)
	}
	if ok := checkCdnInfoWritable(d, "compression", dc.Compression); ok {
		dMap := map[string]interface{}{
			"switch": dc.Compression.Switch,
		}
		if rules := dc.Compression.CompressionRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"compress":        item.Compress,
					"min_length":      item.MinLength,
					"max_length":      item.MaxLength,
					"algorithms":      item.Algorithms,
					"file_extensions": item.FileExtensions,
					"rule_type":       item.RuleType,
					"rule_paths":      item.RulePaths,
				}
				list = append(list, rule)
			}
			dMap["compression_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "compression", dMap)
	}
	if ok := checkCdnInfoWritable(d, "band_width_alert", dc.BandwidthAlert); ok {
		dMap := map[string]interface{}{
			"switch":                     dc.BandwidthAlert.Switch,
			"bps_threshold":              dc.BandwidthAlert.BpsThreshold,
			"counter_measure":            dc.BandwidthAlert.CounterMeasure,
			"last_trigger_time":          dc.BandwidthAlert.LastTriggerTime,
			"alert_switch":               dc.BandwidthAlert.AlertSwitch,
			"alert_percentage":           dc.BandwidthAlert.AlertPercentage,
			"last_trigger_time_overseas": dc.BandwidthAlert.LastTriggerTimeOverseas,
			"metric":                     dc.BandwidthAlert.Metric,
		}
		if si := dc.BandwidthAlert.StatisticItems; len(si) > 0 {
			rules := make([]interface{}, 0)
			for i := range si {
				item := si[i]
				rule := map[string]interface{}{
					"switch":           item.Switch,
					"type":             item.Type,
					"unblock_time":     item.UnBlockTime,
					"bps_threshold":    item.BpsThreshold,
					"counter_measure":  item.CounterMeasure,
					"alert_switch":     item.AlertSwitch,
					"alert_percentage": item.AlertPercentage,
					"metric":           item.Metric,
					"cycle":            item.Cycle,
				}
				rules = append(rules, rule)
			}
			dMap["statistic_item"] = rules
		}
		_ = helper.SetMapInterfaces(d, "band_width_alert", dMap)
	}
	if ok := checkCdnInfoWritable(d, "error_page", dc.ErrorPage); ok {
		dMap := map[string]interface{}{
			"switch": dc.ErrorPage.Switch,
		}
		if rules := dc.ErrorPage.PageRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"status_code":   item.StatusCode,
					"redirect_code": item.RedirectCode,
					"redirect_url":  item.RedirectUrl,
				}
				list = append(list, rule)
			}
			dMap["page_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "error_page", dMap)
	}
	if ok := checkCdnInfoWritable(d, "response_header", dc.ResponseHeader); ok {
		dMap := map[string]interface{}{
			"switch": dc.ResponseHeader.Switch,
		}
		if rules := dc.ResponseHeader.HeaderRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"header_mode":  item.HeaderMode,
					"header_name":  item.HeaderName,
					"header_value": item.HeaderValue,
					"rule_type":    item.RuleType,
					"rule_paths":   item.RulePaths,
				}
				list = append(list, rule)
			}
			dMap["header_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "response_header", dMap)
	}
	if ok := checkCdnInfoWritable(d, "downstream_capping", dc.DownstreamCapping); ok {
		dMap := map[string]interface{}{
			"switch": dc.DownstreamCapping.Switch,
		}
		if rules := dc.DownstreamCapping.CappingRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"rule_type":      item.RuleType,
					"rule_paths":     item.RulePaths,
					"kbps_threshold": item.KBpsThreshold,
				}
				list = append(list, rule)
			}
			dMap["capping_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "downstream_capping", dMap)
	}
	if _, ok := d.GetOk("response_header_cache_switch"); ok && dc.ResponseHeaderCache != nil {
		_ = d.Set("response_header_cache_switch", dc.ResponseHeaderCache.Switch)
	}
	if ok := checkCdnInfoWritable(d, "origin_pull_optimization", dc.OriginPullOptimization); ok {
		dMap := map[string]interface{}{
			"switch":            dc.OriginPullOptimization.Switch,
			"optimization_type": dc.OriginPullOptimization.OptimizationType,
		}
		_ = helper.SetMapInterfaces(d, "origin_pull_optimization", dMap)
	}
	if _, ok := d.GetOk("seo_switch"); ok && dc.Seo != nil {
		_ = d.Set("seo_switch", dc.Seo.Switch)
	}
	if ok := checkCdnInfoWritable(d, "referer", dc.Referer); ok {
		dMap := map[string]interface{}{
			"switch": dc.Referer.Switch,
		}
		if rules := dc.Referer.RefererRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"rule_type":    item.RuleType,
					"rule_paths":   item.RulePaths,
					"referer_type": item.RefererType,
					"referers":     item.Referers,
					"allow_empty":  item.AllowEmpty,
				}
				list = append(list, rule)
			}
			dMap["referer_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "referer", dMap)
	}
	if _, ok := d.GetOk("video_seek_switch"); ok && dc.VideoSeek != nil {
		_ = d.Set("video_seek_switch", dc.VideoSeek.Switch)
	}
	if ok := checkCdnInfoWritable(d, "max_age", dc.MaxAge); ok {
		dMap := map[string]interface{}{
			"switch": dc.MaxAge.Switch,
		}
		if rules := dc.MaxAge.MaxAgeRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"follow_origin":    item.FollowOrigin,
					"max_age_contents": item.MaxAgeContents,
					"max_age_type":     item.MaxAgeType,
					"max_age_time":     item.MaxAgeTime,
				}
				list = append(list, rule)
			}
			dMap["max_age_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "max_age", dMap)
	}
	if ok := checkCdnInfoWritable(d, "specific_config_mainland", dc.SpecificConfig); ok {
		specConfig, err := json.Marshal(dc.SpecificConfig.Mainland)
		if err == nil {
			_ = d.Set("specific_config_mainland", string(specConfig))
		}
	}
	if ok := checkCdnInfoWritable(d, "specific_config_overseas", dc.SpecificConfig); ok {
		specConfig, err := json.Marshal(dc.SpecificConfig.Overseas)
		if err == nil {
			_ = d.Set("specific_config_overseas", string(specConfig))
		}
	}
	if ok := checkCdnInfoWritable(d, "origin_pull_timeout", dc.OriginPullTimeout); ok {
		_ = helper.SetMapInterfaces(d, "origin_pull_timeout", map[string]interface{}{
			"connect_timeout": dc.OriginPullTimeout.ConnectTimeout,
			"receive_timeout": dc.OriginPullTimeout.ReceiveTimeout,
		})
	}
	if ok := checkCdnInfoWritable(d, "post_max_size", dc.PostMaxSize); ok {
		dMap := map[string]interface{}{
			"switch":   dc.PostMaxSize.Switch,
			"max_size": *dc.PostMaxSize.MaxSize / 1024 / 1024,
		}
		_ = helper.SetMapInterfaces(d, "post_max_size", dMap)
	}
	if ok := checkCdnInfoWritable(d, "cache_key", dc.CacheKey); ok {
		dMap := map[string]interface{}{
			"full_url_cache": dc.CacheKey.FullUrlCache,
			"ignore_case":    dc.CacheKey.IgnoreCase,
		}
		if qs := dc.CacheKey.QueryString; qs != nil {
			dMap["query_string"] = []interface{}{
				map[string]interface{}{
					"switch":  qs.Switch,
					"action":  qs.Action,
					"value":   qs.Value,
					"reorder": qs.Reorder,
				},
			}
		}
		if krs := dc.CacheKey.KeyRules; len(krs) > 0 {
			dMaps := make([]interface{}, 0)
			for i := range krs {
				kr := krs[i]
				dMap := map[string]interface{}{
					"full_url_cache": kr.FullUrlCache,
					"ignore_case":    kr.IgnoreCase,
				}
				if kr.RuleType != nil {
					dMap["rule_type"] = kr.RuleType
				}
				if len(kr.RulePaths) > 0 {
					dMap["rule_paths"] = helper.StringsInterfaces(kr.RulePaths)
				}
				if krqs := kr.QueryString; krqs != nil {
					dMap["query_string"] = []interface{}{
						map[string]interface{}{
							"value":  krqs.Value,
							"switch": krqs.Switch,
							"action": krqs.Action,
						},
					}
				}
				dMaps = append(dMaps, dMap)
			}
			dMap["key_rules"] = dMaps
		}
		_ = helper.SetMapInterfaces(d, "cache_key", dMap)
	} else if dc.CacheKey != nil && dc.CacheKey.FullUrlCache != nil {
		fullUrlCache := *dc.CacheKey.FullUrlCache == CDN_SWITCH_ON
		_ = d.Set("full_url_cache", fullUrlCache)
	}
	if _, ok := d.GetOk("offline_cache_switch"); ok && dc.OfflineCache != nil {
		_ = d.Set("offline_cache_switch", dc.OfflineCache.Switch)
	}
	if _, ok := d.GetOk("quic_switch"); ok && dc.Quic != nil {
		_ = d.Set("quic_switch", dc.Quic.Switch)
	}
	if ok := checkCdnInfoWritable(d, "aws_private_access", dc.AwsPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "aws_private_access", map[string]interface{}{
			"switch":     dc.AwsPrivateAccess.Switch,
			"access_key": dc.AwsPrivateAccess.AccessKey,
			"secret_key": dc.AwsPrivateAccess.SecretKey,
			"bucket":     dc.AwsPrivateAccess.Bucket,
			"region":     dc.AwsPrivateAccess.Region,
		})
	}
	if ok := checkCdnInfoWritable(d, "oss_private_access", dc.OssPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "oss_private_access", map[string]interface{}{
			"switch":     dc.OssPrivateAccess.Switch,
			"access_key": dc.OssPrivateAccess.AccessKey,
			"secret_key": dc.OssPrivateAccess.SecretKey,
			"bucket":     dc.OssPrivateAccess.Bucket,
			"region":     dc.OssPrivateAccess.Region,
		})
	}
	if ok := checkCdnInfoWritable(d, "hw_private_access", dc.HwPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "hw_private_access", map[string]interface{}{
			"switch":     dc.HwPrivateAccess.Switch,
			"access_key": dc.HwPrivateAccess.AccessKey,
			"secret_key": dc.HwPrivateAccess.SecretKey,
			"bucket":     dc.HwPrivateAccess.Bucket,
		})
	}
	if ok := checkCdnInfoWritable(d, "qn_private_access", dc.QnPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "qn_private_access", map[string]interface{}{
			"switch":     dc.QnPrivateAccess.Switch,
			"access_key": dc.QnPrivateAccess.AccessKey,
			"secret_key": dc.QnPrivateAccess.SecretKey,
		})
	}
	if ok := checkCdnInfoWritable(d, "others_private_access", dc.OthersPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "others_private_access", map[string]interface{}{
			"switch":     dc.OthersPrivateAccess.Switch,
			"access_key": dc.OthersPrivateAccess.AccessKey,
			"secret_key": dc.OthersPrivateAccess.SecretKey,
			"bucket":     dc.OthersPrivateAccess.Bucket,
			"region":     dc.OthersPrivateAccess.Region,
		})
	}

	tags, errRet := tagService.DescribeResourceTags(ctx, CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
	if errRet != nil {
		return errRet
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cdnService := CdnService{client: client}

	d.Partial(true)
	updateAttrs := make([]string, 0)

	domain := d.Id()
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = &domain

	if d.HasChange("service_type") {
		request.ServiceType = helper.String(d.Get("service_type").(string))
		updateAttrs = append(updateAttrs, "service_type")
	}
	if d.HasChange("project_id") {
		request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
		updateAttrs = append(updateAttrs, "project_id")
	}
	if d.HasChange("area") {
		request.Area = helper.String(d.Get("area").(string))
		updateAttrs = append(updateAttrs, "area")
	}
	if d.HasChange("range_origin_switch") {
		request.RangeOriginPull = &cdn.RangeOriginPull{}
		request.RangeOriginPull.Switch = helper.String(d.Get("range_origin_switch").(string))
		updateAttrs = append(updateAttrs, "range_origin_switch")
	}
	if d.HasChange("ipv6_access_switch") {
		request.Ipv6Access = &cdn.Ipv6Access{}
		request.Ipv6Access.Switch = helper.String(d.Get("ipv6_access_switch").(string))
		updateAttrs = append(updateAttrs, "ipv6_access_switch")
	}
	if d.HasChange("follow_redirect_switch") {
		request.FollowRedirect = &cdn.FollowRedirect{}
		request.FollowRedirect.Switch = helper.String(d.Get("follow_redirect_switch").(string))
		updateAttrs = append(updateAttrs, "follow_redirect_switch")
	}
	if d.HasChange("origin") {
		updateAttrs = append(updateAttrs, "origin")
		origins := d.Get("origin").([]interface{})
		if len(origins) < 1 {
			return fmt.Errorf("origin is required")
		}
		origin := origins[0].(map[string]interface{})
		request.Origin = &cdn.Origin{}
		request.Origin.OriginType = helper.String(origin["origin_type"].(string))
		originList := origin["origin_list"].([]interface{})
		request.Origin.Origins = make([]*string, 0, len(originList))
		for _, item := range originList {
			request.Origin.Origins = append(request.Origin.Origins, helper.String(item.(string)))
		}
		if v := origin["server_name"]; v.(string) != "" {
			request.Origin.ServerName = helper.String(v.(string))
		}
		if v := origin["cos_private_access"]; v.(string) != "" {
			request.Origin.CosPrivateAccess = helper.String(v.(string))
		}
		if v := origin["origin_pull_protocol"]; v.(string) != "" {
			request.Origin.OriginPullProtocol = helper.String(v.(string))
		}
		if v := origin["backup_origin_type"]; v.(string) != "" {
			request.Origin.BackupOriginType = helper.String(v.(string))
		}
		if v := origin["backup_server_name"]; v.(string) != "" {
			request.Origin.BackupServerName = helper.String(v.(string))
		}
		if v := origin["backup_origin_list"]; len(v.([]interface{})) > 0 {
			backupOriginList := v.([]interface{})
			request.Origin.BackupOrigins = make([]*string, 0, len(backupOriginList))
			for _, item := range backupOriginList {
				request.Origin.BackupOrigins = append(request.Origin.BackupOrigins, helper.String(item.(string)))
			}
		}
		if v := origin["origin_company"]; v.(string) != "" {
			request.Origin.OriginCompany = helper.String(v.(string))
		}
	}
	if d.HasChange("request_header") {
		updateAttrs = append(updateAttrs, "request_header")
		requestHeaders := d.Get("request_header").([]interface{})
		requestHeader := requestHeaders[0].(map[string]interface{})
		headerRule := requestHeader["header_rules"].([]interface{})
		var headerRules []*cdn.HttpHeaderPathRule
		for _, value := range headerRule {
			hr := &cdn.HttpHeaderPathRule{}
			headerRuleMap := value.(map[string]interface{})
			headerMode := headerRuleMap["header_mode"].(string)
			headerName := headerRuleMap["header_name"].(string)
			headerValue := headerRuleMap["header_value"].(string)
			ruleType := headerRuleMap["rule_type"].(string)
			rulePaths := headerRuleMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			for _, value := range rulePaths {
				rulePathList = append(rulePathList, helper.String(value.(string)))
			}
			hr.HeaderMode = &headerMode
			hr.HeaderName = &headerName
			hr.HeaderValue = &headerValue
			hr.RuleType = &ruleType
			hr.RulePaths = rulePathList
			headerRules = append(headerRules, hr)
		}
		request.RequestHeader = &cdn.RequestHeader{}
		request.RequestHeader.Switch = helper.String(requestHeader["switch"].(string))
		request.RequestHeader.HeaderRules = headerRules
	}
	if d.HasChange("rule_cache") {
		updateAttrs = append(updateAttrs, "rule_cache")
		ruleCache := d.Get("rule_cache").([]interface{})
		var ruleCaches []*cdn.RuleCache
		for _, v := range ruleCache {
			re := &cdn.RuleCache{}
			ruleCacheMap := v.(map[string]interface{})
			rulePaths := ruleCacheMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			ruleType := ruleCacheMap["rule_type"].(string)
			if ruleType == CDN_RULE_TYPE_DEFAULT {
				rulePathList = append(rulePathList, helper.String(CDN_RULE_PATH))
			} else {
				for _, value := range rulePaths {
					rulePathList = append(rulePathList, helper.String(value.(string)))
				}
			}
			switchFlag := ruleCacheMap["switch"].(string)
			cacheTime := ruleCacheMap["cache_time"].(int)
			compareMaxAge := ruleCacheMap["compare_max_age"].(string)
			ignoreCacheControl := ruleCacheMap["ignore_cache_control"].(string)
			ignoreSetCookie := ruleCacheMap["ignore_set_cookie"].(string)
			noCacheSwitch := ruleCacheMap["no_cache_switch"].(string)
			reValidate := ruleCacheMap["re_validate"].(string)
			followOriginSwitch := ruleCacheMap["follow_origin_switch"].(string)
			ruleCacheConfig := &cdn.RuleCacheConfig{}
			cache := &cdn.CacheConfigCache{}
			noCache := &cdn.CacheConfigNoCache{}
			followOrigin := &cdn.CacheConfigFollowOrigin{}
			ruleCacheConfig.Cache = cache
			ruleCacheConfig.NoCache = noCache
			ruleCacheConfig.FollowOrigin = followOrigin
			re.CacheConfig = ruleCacheConfig
			re.RulePaths = rulePathList
			re.RuleType = &ruleType
			re.CacheConfig.Cache.Switch = &switchFlag
			re.CacheConfig.Cache.CacheTime = helper.IntInt64(cacheTime)
			re.CacheConfig.Cache.CompareMaxAge = &compareMaxAge
			re.CacheConfig.Cache.IgnoreCacheControl = &ignoreCacheControl
			re.CacheConfig.Cache.IgnoreSetCookie = &ignoreSetCookie
			re.CacheConfig.NoCache.Switch = &noCacheSwitch
			re.CacheConfig.NoCache.Revalidate = &reValidate
			re.CacheConfig.FollowOrigin.Switch = &followOriginSwitch
			heuristicCacheSwitch := ruleCacheMap["heuristic_cache_switch"].(string)
			heuristicCacheTime := ruleCacheMap["heuristic_cache_time"].(int)
			if heuristicCacheSwitch != "" {
				re.CacheConfig.FollowOrigin.HeuristicCache = &cdn.HeuristicCache{
					Switch:      &heuristicCacheSwitch,
					CacheConfig: &cdn.CacheConfig{},
				}
				if heuristicCacheTime > 0 {
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTimeSwitch = helper.String("on")
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTime = helper.IntInt64(heuristicCacheTime)
				} else {
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTimeSwitch = helper.String("off")
				}
			}
			ruleCaches = append(ruleCaches, re)
		}
		request.Cache = &cdn.Cache{}
		request.Cache.RuleCache = ruleCaches
	}
	if d.HasChange("https_config") {
		updateAttrs = append(updateAttrs, "https_config")
		httpsConfigs := d.Get("https_config").([]interface{})
		if len(httpsConfigs) > 0 {
			config := httpsConfigs[0].(map[string]interface{})
			request.Https = &cdn.Https{}
			request.Https.Switch = helper.String(config["https_switch"].(string))
			if v := config["http2_switch"]; v.(string) != "" {
				request.Https.Http2 = helper.String(v.(string))
			}
			request.Https.OcspStapling = helper.String(config["ocsp_stapling_switch"].(string))
			request.Https.Spdy = helper.String(config["spdy_switch"].(string))
			request.Https.VerifyClient = helper.String(config["verify_client"].(string))
			if v := config["server_certificate_config"]; len(v.([]interface{})) > 0 {
				serverCerts := v.([]interface{})
				if len(serverCerts) > 0 && serverCerts[0] != nil {
					serverCert := serverCerts[0].(map[string]interface{})
					request.Https.CertInfo = &cdn.ServerCert{}
					if v := serverCert["certificate_id"]; v.(string) != "" {
						request.Https.CertInfo.CertId = helper.String(v.(string))
					}
					if v := serverCert["certificate_content"]; v.(string) != "" {
						request.Https.CertInfo.Certificate = helper.String(v.(string))
					}
					if v := serverCert["private_key"]; v.(string) != "" {
						request.Https.CertInfo.PrivateKey = helper.String(v.(string))
					}
					if v := serverCert["message"]; v.(string) != "" {
						request.Https.CertInfo.Message = helper.String(v.(string))
					}
				}
			}
			if v := config["client_certificate_config"]; len(v.([]interface{})) > 0 {
				clientCerts := v.([]interface{})
				if len(clientCerts) > 0 && clientCerts[0] != nil {
					clientCert := clientCerts[0].(map[string]interface{})
					request.Https.ClientCertInfo = &cdn.ClientCert{}
					if v := clientCert["certificate_content"]; v.(string) != "" {
						request.Https.ClientCertInfo.Certificate = helper.String(v.(string))
					}
				}
			}
			if v, ok := config["force_redirect"]; ok {
				forceRedirect := v.([]interface{})
				if len(forceRedirect) > 0 && forceRedirect[0] != nil {
					var redirect cdn.ForceRedirect
					redirectMap := forceRedirect[0].(map[string]interface{})
					if sw := redirectMap["switch"]; sw.(string) != "" {
						redirect.Switch = helper.String(sw.(string))
					}
					if rt := redirectMap["redirect_type"]; rt.(string) != "" {
						redirect.RedirectType = helper.String(rt.(string))
					}
					if rsc := redirectMap["redirect_status_code"]; rsc.(int) != 0 {
						redirect.RedirectStatusCode = helper.Int64(int64(rsc.(int)))
					}
					if ch := redirectMap["carry_headers"]; ch.(string) != "" {
						redirect.CarryHeaders = helper.String(ch.(string))
					}
					request.ForceRedirect = &redirect
				}
			}
			if v, ok := config["tls_versions"]; ok {
				request.Https.TlsVersion = helper.InterfacesStringsPoint(v.([]interface{}))
			}
		}
	}

	if d.HasChange("authentication") {
		updateAttrs = append(updateAttrs, "authentication")
		if v, ok := helper.InterfacesHeadMap(d, "authentication"); ok {
			switchOn := v["switch"].(string)
			request.Authentication = &cdn.Authentication{
				Switch: &switchOn,
			}

			if v, ok := v["type_a"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					signParam      = item["sign_param"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeA = &cdn.AuthenticationTypeA{
					SecretKey:      &secretKey,
					SignParam:      &signParam,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeA.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_b"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeB = &cdn.AuthenticationTypeB{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeB.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_c"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeC = &cdn.AuthenticationTypeC{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if timeFormat, ok := item["time_format"].(string); ok {
					request.Authentication.TypeC.TimeFormat = &timeFormat
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeC.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_d"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
					timeParam      = item["time_param"].(string)
				)

				request.Authentication.TypeD = &cdn.AuthenticationTypeD{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
					TimeParam:      &timeParam,
				}

				if timeFormat, ok := item["time_format"].(string); ok {
					request.Authentication.TypeD.TimeFormat = &timeFormat
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeD.BackupSecretKey = &backupSecretKey
				}
			}
		}
	}

	// more added
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "ip_filter"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "ip_filter")
		vSwitch := v["switch"].(string)
		request.IpFilter = &cdn.IpFilter{
			Switch: &vSwitch,
		}
		if vv, ok := v["filter_type"].(string); ok {
			request.IpFilter.FilterType = &vv
		}
		if vv, ok := v["filters"].([]interface{}); ok {
			request.IpFilter.Filters = helper.InterfacesStringsPoint(vv)
		}

		//need white list func
		if vv, ok := v["filter_rules"].([]interface{}); ok && len(vv) > 0 {
			filterRules := make([]*cdn.IpFilterPathRule, 0)
			for i := range vv {
				item := vv[i].(map[string]interface{})
				rule := &cdn.IpFilterPathRule{}
				if rv, ok := item["filter_type"].(string); ok && rv != "" {
					rule.FilterType = &rv
				}
				if rv, ok := item["filters"].([]interface{}); ok && len(rv) > 0 {
					rule.Filters = helper.InterfacesStringsPoint(rv)
				}
				if rv, ok := item["rule_type"].(string); ok && rv != "" {
					rule.RuleType = &rv
				}
				if rv, ok := item["rule_paths"].([]interface{}); ok && len(rv) > 0 {
					rule.RulePaths = helper.InterfacesStringsPoint(rv)
				}
				filterRules = append(filterRules, rule)
			}
			request.IpFilter.FilterRules = filterRules
		}

		if vv, ok := v["return_code"].(int); ok {
			request.IpFilter.ReturnCode = helper.IntInt64(vv)
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "ip_freq_limit"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "ip_freq_limit")
		vSwitch := v["switch"].(string)
		qps := v["qps"].(int)
		request.IpFreqLimit = &cdn.IpFreqLimit{
			Switch: &vSwitch,
			Qps:    helper.IntInt64(qps),
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "status_code_cache"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "status_code_cache")
		vSwitch := v["switch"].(string)
		request.StatusCodeCache = &cdn.StatusCodeCache{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.StatusCodeCacheRule, 0)
		rules := v["cache_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.StatusCodeCacheRule{
				StatusCode: helper.String(item["status_code"].(string)),
			}
			if v, ok := item["cache_time"].(int); ok && v > 0 {
				rule.CacheTime = helper.IntInt64(v)
			}
			requestRules = append(requestRules, rule)
		}
		request.StatusCodeCache.CacheRules = requestRules
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "compression"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "compression")
		vSwitch := v["switch"].(string)
		request.Compression = &cdn.Compression{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.CompressionRule, 0)
		rules := v["compression_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			var (
				compress = item["compress"].(bool)
			)
			rule := &cdn.CompressionRule{
				Compress: &compress,
			}
			if v, ok := item["min_length"].(int); ok && v > 0 {
				rule.MinLength = helper.IntInt64(v)
			}
			if v, ok := item["max_length"].(int); ok && v > 0 {
				rule.MaxLength = helper.IntInt64(v)
			}
			if v, ok := item["algorithms"].([]interface{}); ok && len(v) > 0 {
				rule.Algorithms = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["file_extensions"].([]interface{}); ok && len(v) > 0 {
				rule.FileExtensions = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}

			requestRules = append(requestRules, rule)
		}
		request.Compression.CompressionRules = requestRules

	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "band_width_alert"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "band_width_alert")
		vSwitch := v["switch"].(string)
		request.BandwidthAlert = &cdn.BandwidthAlert{
			Switch: &vSwitch,
		}
		if v, ok := v["bps_threshold"].(int); ok && v > 0 {
			request.BandwidthAlert.BpsThreshold = helper.IntInt64(v)
		}
		if v, ok := v["counter_measure"].(string); ok && v != "" {
			request.BandwidthAlert.CounterMeasure = &v
		}
		//if v, ok := v["last_trigger_time"].(string); ok && v != "" {
		//	request.BandwidthAlert.LastTriggerTime = &v
		//}
		if v, ok := v["alert_switch"].(string); ok && v != "" {
			request.BandwidthAlert.AlertSwitch = &v
		}
		if v, ok := v["alert_percentage"].(int); ok && v > 0 {
			request.BandwidthAlert.AlertPercentage = helper.IntInt64(v)
		}
		//if v, ok := v["last_trigger_time_overseas"].(string); ok && v != "" {
		//	request.BandwidthAlert.LastTriggerTimeOverseas = &v
		//}
		if v, ok := v["metric"].(string); ok && v != "" {
			request.BandwidthAlert.Metric = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "error_page"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "error_page")
		vSwitch := v["switch"].(string)
		request.ErrorPage = &cdn.ErrorPage{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.ErrorPageRule, 0)
		rules := v["page_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.ErrorPageRule{}
			if v, ok := item["status_code"].(int); ok && v != 0 {
				rule.StatusCode = helper.IntInt64(v)
			}
			if v, ok := item["redirect_code"].(int); ok && v != 0 {
				rule.RedirectCode = helper.IntInt64(v)
			}
			if v, ok := item["redirect_url"].(string); ok && v != "" {
				rule.RedirectUrl = &v
			}
			requestRules = append(requestRules, rule)
		}
		request.ErrorPage.PageRules = requestRules
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "response_header"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "response_header")
		vSwitch := v["switch"].(string)
		request.ResponseHeader = &cdn.ResponseHeader{
			Switch: &vSwitch,
		}
		responseRules := make([]*cdn.HttpHeaderPathRule, 0)
		rules := v["header_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.HttpHeaderPathRule{}
			if v, ok := item["header_mode"].(string); ok && v != "" {
				rule.HeaderMode = &v
			}
			if v, ok := item["header_name"].(string); ok && v != "" {
				rule.HeaderName = &v
			}
			if v, ok := item["header_value"].(string); ok && v != "" {
				rule.HeaderValue = &v
			}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			responseRules = append(responseRules, rule)
		}
		request.ResponseHeader.HeaderRules = responseRules
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "downstream_capping"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "downstream_capping")
		vSwitch := v["switch"].(string)
		request.DownstreamCapping = &cdn.DownstreamCapping{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.CappingRule, 0)
		rules := v["capping_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.CappingRule{}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["kbps_threshold"].(int); ok && v > 0 {
				rule.KBpsThreshold = helper.IntInt64(v)
			}
			requestRules = append(requestRules, rule)
		}
		request.DownstreamCapping.CappingRules = requestRules
	}
	if d.HasChange("response_header_cache_switch") {
		updateAttrs = append(updateAttrs, "response_header_cache_switch")
		v, ok := d.Get("response_header_cache_switch").(string)
		if !ok || v == "" {
			v = "off"
		}
		request.ResponseHeaderCache = &cdn.ResponseHeaderCache{
			Switch: &v,
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "origin_pull_optimization"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "origin_pull_optimization")
		vSwitch := v["switch"].(string)
		request.OriginPullOptimization = &cdn.OriginPullOptimization{
			Switch: &vSwitch,
		}
		if v, ok := v["optimization_type"].(string); ok && v != "" {
			request.OriginPullOptimization.OptimizationType = &v
		}
	}
	if d.HasChange("seo_switch") {
		updateAttrs = append(updateAttrs, "seo_switch")
		v, ok := d.Get("seo_switch").(string)
		if !ok || v == "" {
			v = "off"
		}
		request.Seo = &cdn.Seo{
			Switch: &v,
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "referer"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "referer")
		vSwitch := v["switch"].(string)
		request.Referer = &cdn.Referer{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.RefererRule, 0)
		rules := v["referer_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.RefererRule{}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["referer_type"].(string); ok && v != "" {
				rule.RefererType = &v
			}
			if v, ok := item["referers"].([]interface{}); ok && len(v) > 0 {
				rule.Referers = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["allow_empty"].(bool); ok {
				rule.AllowEmpty = &v
			}
			requestRules = append(requestRules, rule)
		}
		request.Referer.RefererRules = requestRules
	}
	if d.HasChange("video_seek_switch") {
		updateAttrs = append(updateAttrs, "video_seek_switch")
		v, ok := d.Get("video_seek_switch").(string)
		if !ok || v == "" {
			v = "off"
		}
		request.VideoSeek = &cdn.VideoSeek{
			Switch: &v,
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "max_age"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "max_age")
		vSwitch := v["switch"].(string)
		request.MaxAge = &cdn.MaxAge{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.MaxAgeRule, 0)
		rules := v["max_age_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.MaxAgeRule{}

			if v, ok := item["max_age_type"].(string); ok && v != "" {
				rule.MaxAgeType = &v
			}
			if v, ok := item["max_age_contents"].([]interface{}); ok && len(v) > 0 {
				rule.MaxAgeContents = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["max_age_time"].(int); ok && v > 0 {
				rule.MaxAgeTime = helper.IntInt64(v)
			}
			if v, ok := item["follow_origin"].(string); ok && v != "" {
				rule.FollowOrigin = &v
			}

			requestRules = append(requestRules, rule)
		}
		request.MaxAge.MaxAgeRules = requestRules
	}
	if v, ok := d.GetOk("specific_config_mainland"); ok && v.(string) != "" {
		updateAttrs = append(updateAttrs, "specific_config_mainland")
		request.SpecificConfig = &cdn.SpecificConfig{}
		configStruct := cdn.MainlandConfig{}
		err := json.Unmarshal([]byte(v.(string)), &configStruct)
		if err != nil {
			return fmt.Errorf("unmarshal specific_config_mainland fail: %s", err.Error())
		}
		request.SpecificConfig.Mainland = &configStruct
	}
	if v, ok := d.GetOk("specific_config_overseas"); ok && v.(string) != "" {
		updateAttrs = append(updateAttrs, "specific_config_overseas")
		if request.SpecificConfig == nil {
			request.SpecificConfig = &cdn.SpecificConfig{}
		}
		configStruct := cdn.OverseaConfig{}
		err := json.Unmarshal([]byte(v.(string)), &configStruct)
		if err != nil {
			return fmt.Errorf("unmarshal specific_config_overseas fail: %s", err.Error())
		}
		request.SpecificConfig.Overseas = &configStruct
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "origin_pull_timeout"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "origin_pull_timeout")
		request.OriginPullTimeout = &cdn.OriginPullTimeout{}
		if vv, ok := v["connect_timeout"].(int); ok && vv > 0 {
			request.OriginPullTimeout.ConnectTimeout = helper.IntUint64(vv)
		}
		if vv, ok := v["receive_timeout"].(int); ok && vv > 0 {
			request.OriginPullTimeout.ReceiveTimeout = helper.IntUint64(vv)
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "post_max_size"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "post_max_size")
		vSwitch := v["switch"].(string)
		maxSize := v["max_size"].(int)
		request.PostMaxSize = &cdn.PostSize{
			Switch: &vSwitch,
		}
		if maxSize > 0 {
			request.PostMaxSize.MaxSize = helper.IntInt64(maxSize * 1024 * 1024)
		}
	}
	if v, ok := d.GetOk("offline_cache_switch"); ok {
		request.OfflineCache = &cdn.OfflineCache{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := d.GetOk("quic_switch"); ok {
		request.Quic = &cdn.Quic{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "cache_key"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "cache_key")
		request.CacheKey = &cdn.CacheKey{}
		if fuc := v["full_url_cache"].(string); fuc != "" {
			request.CacheKey.FullUrlCache = &fuc
		}
		if ic := v["ignore_case"].(string); ic != "" {
			request.CacheKey.IgnoreCase = &ic
		}
		if qs, ok := v["query_string"].([]interface{}); ok && len(qs) > 0 {
			if dMap, ok := qs[0].(map[string]interface{}); ok {
				qSwitch := dMap["switch"].(string)
				reorder := dMap["reorder"].(string)
				action := dMap["action"].(string)
				value := dMap["value"].(string)
				request.CacheKey.QueryString = &cdn.QueryStringKey{
					Switch:  &qSwitch,
					Reorder: &reorder,
					Action:  &action,
					Value:   &value,
				}
			}
		}
		if kr, ok := v["key_rules"].([]interface{}); ok {
			for i := range kr {
				rule, ok := kr[i].(map[string]interface{})
				if !ok {
					continue
				}
				ruleType := rule["rule_type"].(string)
				keyRule := &cdn.KeyRule{
					RuleType: &ruleType,
				}
				if vv := rule["full_url_cache"].(string); vv != "" {
					keyRule.FullUrlCache = &vv
				}
				if vv := rule["ignore_case"].(string); vv != "" {
					keyRule.IgnoreCase = &vv
				}
				if vv := rule["rule_tag"].(string); vv != "" {
					keyRule.RuleTag = &vv
				}
				if rp, ok := rule["rule_paths"].([]interface{}); ok {
					keyRule.RulePaths = helper.InterfacesStringsPoint(rp)
				}
				if qs, ok := rule["query_string"].([]interface{}); ok && len(qs) > 0 {
					if dMap, ok := qs[0].(map[string]interface{}); ok {
						vSwitch := dMap["switch"].(string)
						keyRule.QueryString = &cdn.RuleQueryString{
							Switch: &vSwitch,
						}
						if v := dMap["action"].(string); v != "" && vSwitch == "on" {
							keyRule.QueryString.Action = &v
						}
						if v := dMap["value"].(string); v != "" {
							keyRule.QueryString.Value = &v
						}
					}
				}
				request.CacheKey.KeyRules = append(request.CacheKey.KeyRules, keyRule)
			}
		}
	} else if d.HasChange("full_url_cache") {
		fullUrlCache := d.Get("full_url_cache").(bool)
		request.CacheKey = &cdn.CacheKey{}
		if fullUrlCache {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_ON)
		} else {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_OFF)
		}
		updateAttrs = append(updateAttrs, "full_url_cache")
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "aws_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.AwsPrivateAccess = &cdn.AwsPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.AwsPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.AwsPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.AwsPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.AwsPrivateAccess.Bucket = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "oss_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.OssPrivateAccess = &cdn.OssPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.OssPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.OssPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.OssPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.OssPrivateAccess.Bucket = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "hw_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.HwPrivateAccess = &cdn.HwPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.HwPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.HwPrivateAccess.SecretKey = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.HwPrivateAccess.Bucket = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "qn_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.QnPrivateAccess = &cdn.QnPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.QnPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.QnPrivateAccess.SecretKey = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "others_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.OthersPrivateAccess = &cdn.OthersPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.OthersPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.OthersPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.OthersPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.OthersPrivateAccess.Bucket = &v
		}
	}

	if v := d.Get("explicit_using_dry_run").(bool); v {
		_ = d.Set("dry_run_update_result", request.ToJsonString())
		return resourceTencentCloudCdnDomainRead(d, meta)
	}

	if len(updateAttrs) > 0 {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().UpdateDomainConfig(request)
			if err != nil {
				if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROR {
						return resource.NonRetryableError(err)
					}
				}
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := svctag.NewTagService(client)
		region := client.Region
		resourceName := tccommon.BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudCdnDomainRead(d, meta)
}

func resourceTencentCloudCdnDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cdnService := CdnService{client: client}

	domain := d.Id()

	if v, ok := d.Get("explicit_using_dry_run").(bool); ok && v {
		d.SetId("")
		return nil
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(client)
		region := client.Region
		resourceName := tccommon.BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		deleteTags := make([]string, 0, len(tags))
		for key := range tags {
			deleteTags = append(deleteTags, key)
		}
		err := tagService.ModifyTags(ctx, resourceName, nil, deleteTags)
		if err != nil {
			return err
		}
	}

	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if domainConfig == nil {
		return nil
	}

	if *domainConfig.Status == CDN_DOMAIN_STATUS_ONLINE {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet = cdnService.StopDomain(ctx, domain)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet = cdnService.DeleteDomain(ctx, domain)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func updateCdnModifyOnlyParams(d *schema.ResourceData, meta interface{}, ctx context.Context) error {
	if !d.HasChanges("post_max_size") {
		return nil
	}

	domain := d.Id()
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := CdnService{client}
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = &domain

	if v, ok := helper.InterfacesHeadMap(d, "post_max_size"); ok {
		vSwitch := v["switch"].(string)
		maxSize := v["max_size"].(int)
		request.PostMaxSize = &cdn.PostSize{
			Switch: &vSwitch,
		}
		if maxSize > 0 {
			request.PostMaxSize.MaxSize = helper.IntInt64(maxSize * 1024 * 1024)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		err := service.UpdateDomainConfig(ctx, request)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROR {
					return resource.NonRetryableError(err)
				}
			}
			return tccommon.RetryError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func checkCdnHeadMapOkAndChanged(d *schema.ResourceData, key string) (v map[string]interface{}, ok bool, changed bool) {
	changed = d.HasChange(key)
	v, ok = helper.InterfacesHeadMap(d, key)
	return
}

func checkCdnInfoWritable(d *schema.ResourceData, key string, val interface{}) bool {
	_, ok := helper.InterfacesHeadMap(d, key)
	return val != nil && ok
}
