package cos

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/beevik/etree"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func originPullRules() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority of origin-pull rules, do not set the same value for multiple rules.",
			},
			"sync_back_to_source": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If `true`, COS will not return 3XX status code when pulling data from an origin server. Current available zone: ap-beijing, ap-shanghai, ap-singapore, ap-mumbai.",
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Triggers the origin-pull rule when the requested file name matches this prefix.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "the protocol used for COS to access the specified origin server. The available value include `HTTP`, `HTTPS` and `FOLLOW`.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Allows only a domain name or IP address. You can optionally append a port number to the address.",
			},
			"follow_query_string": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specifies whether to pass through COS request query string when accessing the origin server.",
			},
			"follow_redirection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specifies whether to follow 3XX redirect to another origin server to pull data from.",
			},
			//"copy_origin_data": {
			//	Type:		 schema.TypeBool,
			//	Optional: 	 true,
			//	Default:	 true,
			//	Description: "",
			//},
			"follow_http_headers": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(i interface{}) int {
					return helper.HashString(i.(string))
				},
				Description: "Specifies the pass through headers when accessing the origin server.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"custom_http_headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Specifies the custom headers that you can add for COS to access your origin server.",
			},
			//"redirect_prefix": {
			//	Type:		schema.TypeString,
			//	Optional:   true,
			//	Description: "Prefix for the file to which a request is redirected when the origin-pull rule is triggered.",
			//},
			//"redirect_suffix": {
			//	Type:		schema.TypeString,
			//	Optional:   true,
			//	Description: "Suffix for the file to which a request is redirected when the origin-pull rule is triggered.",
			//},
		},
	}
}

// x-cos-grant-* headers may conflict with xml acl body, we don't open up for now.
//func aclGrantHeaders() *schema.Schema {
//	return &schema.Schema{
//		Type:        schema.TypeMap,
//		Optional:    true,
//		Description: "ACL x-cos-grant-* headers for multiple grand info",
//		Elem: &schema.Resource{
//			Schema: map[string]*schema.Schema{
//				"grant_read": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Allows grantee to read the bucket; format: `id=\"[OwnerUin]\"`.Use comma (,) to separate multiple users, e.g `id=\"100000000001\",id=\"100000000002\"`",
//				},
//				"grant_write": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Allows grantee to write to the bucket; format: `id=\"[OwnerUin]\"`.Use comma (,) to separate multiple users, e.g `id=\"100000000001\",id=\"100000000002\"`",
//				},
//				"grant_read_acp": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Allows grantee to read the ACL of the bucket; format: `id=\"[OwnerUin]\"`.Use comma (,) to separate multiple users, e.g `id=\"100000000001\",id=\"100000000002\"`",
//				},
//				"grant_write_acp": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Allows grantee to write the ACL of the bucket; format: `id=\"[OwnerUin]\"`.Use comma (,) to separate multiple users, e.g `id=\"100000000001\",id=\"100000000002\"`",
//				},
//				"grant_full_control": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Grants a user full permission to perform operations on the bucket; format: `id=\"[OwnerUin]\"`.Use comma (,) to separate multiple users, e.g `id=\"100000000001\",id=\"100000000002\"`",
//				},
//			},
//		},
//	}
//}

func ResourceTencentCloudCosBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketCreate,
		Read:   resourceTencentCloudCosBucketRead,
		Update: resourceTencentCloudCosBucketUpdate,
		Delete: resourceTencentCloudCosBucketDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"force_clean": false,
			}),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateCosBucketName,
				Description:  "The name of a bucket to be created. Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.",
			},
			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  s3.ObjectCannedACLPrivate,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{
					s3.ObjectCannedACLPrivate,
					s3.ObjectCannedACLPublicRead,
					s3.ObjectCannedACLPublicReadWrite,
				}),
				Description: "The canned ACL to apply. Valid values: private, public-read, and public-read-write. Defaults to private.",
			},
			"acl_body": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
					return ACLBodyDiffFunc(olds, news, d)
				},
				DiffSuppressOnRefresh: true,
				ValidateFunc:          validateACLBody,
				Description:           "ACL XML body for multiple grant info. NOTE: this argument will overwrite `acl`. Check https://intl.cloud.tencent.com/document/product/436/7737 for more detail.",
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The server-side encryption algorithm to use. Valid values are `AES256`, `KMS` and `cos/kms`, `cos/kms` is for cdc cos scenario.",
			},
			"kms_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The KMS Master Key ID. This value is valid only when `encryption_algorithm` is set to KMS or cos/kms. Set kms id to the specified value. If not specified, the default kms id is used.",
			},
			"versioning_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable bucket versioning. NOTE: The `multi_az` feature is true for the current bucket, cannot disable version control.",
			},
			"acceleration_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable bucket acceleration.",
			},
			"force_clean": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force cleanup all objects before delete bucket.",
			},
			"replica_role": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"replica_rules", "versioning_enable"},
				Description:  "Request initiator identifier, format: `qcs::cam::uin/<owneruin>:uin/<subuin>`. NOTE: only `versioning_enable` is true can configure this argument.",
			},
			"replica_rules": {
				Type:         schema.TypeList,
				Optional:     true,
				Description:  "List of replica rule. NOTE: only `versioning_enable` is true and `replica_role` set can configure this argument.",
				RequiredWith: []string{"replica_role", "versioning_enable"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of a specific rule.",
						},
						"status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Status identifier, available values: `Enabled`, `Disabled`.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Prefix matching policy. Policies cannot overlap; otherwise, an error will be returned. To match the root directory, leave this parameter empty.",
						},
						"destination_bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Destination bucket identifier, format: `qcs::cos:<region>::<bucketname-appid>`. NOTE: destination bucket must enable versioning.",
						},
						"destination_storage_class": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Storage class of destination, available values: `STANDARD`, `INTELLIGENT_TIERING`, `STANDARD_IA`. default is following current class of destination.",
						},
					},
				},
			},
			"cors_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A rule of Cross-Origin Resource Sharing (documented below).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies which origins are allowed.",
						},
						"allowed_methods": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies which methods are allowed. Can be `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.",
						},
						"allowed_headers": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies which headers are allowed.",
						},
						"max_age_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies time in seconds that browser can cache the response for a preflight request.",
						},
						"expose_headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Specifies expose header in the response.",
						},
					},
				},
			},
			"origin_pull_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Bucket Origin-Pull settings.",
				Elem:        originPullRules(),
			},
			"origin_domain_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Bucket Origin Domain settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specify domain host.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "REST",
							Description: "Specify origin domain type, available values: `REST`, `WEBSITE`, `ACCELERATE`, default: `REST`.",
						},
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "ENABLED",
							Description:  "Domain status, default: `ENABLED`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ENABLED", "DISABLED"}),
						},
						//"force_replacement": {
						//	Type:		 schema.TypeString,
						//	Optional: 	 true,
						//	Description: "Specify type to replace exist domain resolve record.",
						//},
					},
				},
			},
			"lifecycle_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A configuration of object lifecycle management (documented below).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "A unique identifier for the rule. It can be up to 255 characters.",
						},
						"filter_prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Object key prefix identifying one or more objects to which the rule applies.",
						},
						"transition": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         transitionHash,
							Description: "Specifies a period in the object's transitions (documented below).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: tccommon.ValidateCosBucketLifecycleTimestamp,
										Description:  "Specifies the date after which you want the corresponding action to take effect.",
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(0),
										Description:  "Specifies the number of days after object creation when the specific rule action takes effect.",
									},
									"storage_class": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the storage class to which you want the object to transition. Available values include `STANDARD_IA`, `MAZ_STANDARD_IA`, `INTELLIGENT_TIERING`, `MAZ_INTELLIGENT_TIERING`, `ARCHIVE`, `DEEP_ARCHIVE`. For more information, please refer to: https://cloud.tencent.com/document/product/436/33417.",
									},
								},
							},
						},
						"expiration": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         expirationHash,
							MaxItems:    1,
							Description: "Specifies a period in the object's expire (documented below).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: tccommon.ValidateCosBucketLifecycleTimestamp,
										Description:  "Specifies the date after which you want the corresponding action to take effect.",
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(0),
										Description:  "Specifies the number of days after object creation when the specific rule action takes effect.",
									},
									"delete_marker": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Indicates whether the delete marker of an expired object will be removed.",
									},
								},
							},
						},
						"non_current_transition": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         nonCurrentTransitionHash,
							Description: "Specifies a period in the non current object's transitions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"non_current_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(0),
										Description:  "Number of days after non current object creation when the specific rule action takes effect.",
									},
									"storage_class": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the storage class to which you want the non current object to transition. Available values include `STANDARD_IA`, `MAZ_STANDARD_IA`, `INTELLIGENT_TIERING`, `MAZ_INTELLIGENT_TIERING`, `ARCHIVE`, `DEEP_ARCHIVE`. For more information, please refer to: https://cloud.tencent.com/document/product/436/33417.",
									},
								},
							},
						},
						"non_current_expiration": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         nonCurrentExpirationHash,
							MaxItems:    1,
							Description: "Specifies when non current object versions shall expire.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"non_current_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(0),
										Description:  "Number of days after non current object creation when the specific rule action takes effect. The maximum value is 3650.",
									},
								},
							},
						},
						"abort_incomplete_multipart_upload": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         abortIncompleteMultipartUploadHash,
							MaxItems:    1,
							Description: "Set the maximum time a multipart upload is allowed to remain running.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_after_initiation": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(1),
										Description:  "Specifies the number of days after the multipart upload starts that the upload must be completed. The maximum value is 3650.",
									},
								},
							},
						},
					},
				},
			},
			"website": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A website object(documented below).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS returns this index document when requests are made to the root domain or any of the subfolders.",
						},
						"error_document": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An absolute path to the document to return in case of a 4XX error.",
						},
						"redirect_all_requests_to": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"http", "https"}),
							Description:  "Redirects all request configurations. Valid values: http, https. Default is `http`.",
						},
						"endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`Endpoint` of the static website.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of a bucket.",
			},
			"log_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate the access log of this bucket to be saved or not. Default is `false`. If set `true`, the access log will be saved with `log_target_bucket`. To enable log, the full access of log service must be granted. [Full Access Role Policy](https://intl.cloud.tencent.com/document/product/436/16920).",
			},
			"log_target_bucket": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The target bucket name which saves the access log of this bucket per 5 minutes. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`. User must have full access on this bucket.",
			},
			"log_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The prefix log name which saves the access log of this bucket per 5 minutes. Eg. `MyLogPrefix/`. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`.",
			},
			"multi_az": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates whether to create a bucket of multi available zone.",
			},
			"enable_intelligent_tiering": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable intelligent tiering. NOTE: When intelligent tiering configuration is enabled, it cannot be turned off or modified.",
			},
			"intelligent_tiering_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the limit of days for standard-tier data to low-frequency data in an intelligent tiered storage configuration, with optional days of 30, 60, 90. Default value is 30.",
			},
			"intelligent_tiering_request_frequent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specify the access limit for converting standard layer data into low-frequency layer data in the configuration. The default value is once, which can be used in combination with the number of days to achieve the conversion effect. For example, if the parameter is set to 1 and the number of access days is 30, it means that objects with less than one visit in 30 consecutive days will be reduced from the standard layer to the low frequency layer.",
			},
			"cdc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "CDC cluster ID.",
			},
			//computed
			"cos_bucket_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of this cos bucket.",
			},
		},
	}
}

func resourceTencentCloudCosBucketCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.create")()

	var err error

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	role, roleOk := d.GetOk("replica_role")
	rule, ruleOk := d.GetOk("replica_rules")
	versioning := d.Get("versioning_enable").(bool)
	cdcId := d.Get("cdc_id").(string)

	if !versioning {
		if roleOk || role.(string) != "" {
			return fmt.Errorf("cannot configure role unless versioning enable")
		} else if ruleOk || len(rule.([]interface{})) > 0 {
			return fmt.Errorf("cannot configure replica rule unless versioning enable")
		}
	}

	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	useCosService, createOptions := getBucketPutOptions(d)

	if useCosService {
		// tencent
		err = cosService.TencentCosPutBucket(ctx, bucket, createOptions, cdcId)
	} else {
		// s3
		err = cosService.PutBucket(ctx, bucket, acl, cdcId)
	}
	if err != nil {
		return err
	}

	d.SetId(bucket)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		if err := cosService.SetBucketTags(ctx, bucket, tags, cdcId); err != nil {
			return err
		}
	}

	return resourceTencentCloudCosBucketUpdate(d, meta)
}

func resourceTencentCloudCosBucketRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()
	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cdcId := d.Get("cdc_id").(string)
	code, header, err := cosService.TencentcloudHeadBucket(ctx, bucket, cdcId)
	if err != nil {
		if code == 404 {
			log.Printf("[WARN]%s bucket (%s) not found, error code (404)", logId, bucket)
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	if header != nil && len(header["X-Cos-Bucket-Az-Type"]) > 0 && header["X-Cos-Bucket-Az-Type"][0] == "MAZ" {
		_ = d.Set("multi_az", true)
	}

	cosDomain := meta.(tccommon.ProviderMeta).GetAPIV3Conn().CosDomain
	var cosBucketUrl string
	if cdcId == "" && cosDomain == "" {
		cosBucketUrl = fmt.Sprintf("%s.cos.%s.myqcloud.com", d.Id(), meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region)
	} else if cosDomain != "" {
		parsedURL, _ := url.Parse(cosDomain)
		parsedURL.Host = bucket + "." + parsedURL.Host
		cosBucketUrl = parsedURL.String()
	} else {
		cosBucketUrl = fmt.Sprintf("https://%s.%s.cos-cdc.%s.myqcloud.com", bucket, cdcId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region)
	}

	_ = d.Set("cos_bucket_url", cosBucketUrl)
	// set bucket in the import case
	if _, ok := d.GetOk("bucket"); !ok {
		_ = d.Set("bucket", d.Id())
	}

	// acl
	aclResult, err := cosService.GetBucketACL(ctx, bucket, cdcId)

	if err != nil {
		return err
	}

	aclBody, err := xml.Marshal(aclResult)
	if err != nil {
		return err
	}

	_ = d.Set("acl_body", string(aclBody))

	acl := GetBucketPublicACL(aclResult)

	_ = d.Set("acl", acl)

	// read the cors
	corsRules, err := cosService.GetBucketCors(ctx, bucket, cdcId)
	if err != nil {
		return err
	}
	if err = d.Set("cors_rules", corsRules); err != nil {
		return fmt.Errorf("setting cors_rules error: %v", err)
	}

	if cdcId == "" && cosDomain == "" {
		originPullRules, err := cosService.GetBucketPullOrigin(ctx, bucket)
		if err != nil {
			return err
		}

		if err = d.Set("origin_pull_rules", originPullRules); err != nil {
			return fmt.Errorf("setting origin_pull_rules error: %v", err)
		}

		originDomainRules, err := cosService.GetBucketOriginDomain(ctx, bucket)
		if err != nil {
			return err
		}

		if err = d.Set("origin_domain_rules", originDomainRules); err != nil {
			return fmt.Errorf("setting origin_domain_rules error: %v", err)
		}

		replicaResult, err := cosService.GetBucketReplication(ctx, bucket, cdcId)
		if err != nil {
			return err
		}

		if replicaResult != nil {
			err := setBucketReplication(d, *replicaResult)
			if err != nil {
				return err
			}
		}
	}

	// read the lifecycle
	lifecycleRules, err := cosService.GetBucketLifecycle(ctx, bucket, cdcId)
	if err != nil {
		return err
	}
	if err = d.Set("lifecycle_rules", lifecycleRules); err != nil {
		return fmt.Errorf("setting lifecycle_rules error: %v", err)
	}

	// read the website
	website, err := cosService.GetBucketWebsite(ctx, bucket, cdcId)
	if err != nil {
		return err
	}
	if len(website) > 0 && cosDomain == "" {
		// {bucket}.cos-website.{region}.myqcloud.com
		endPointUrl := fmt.Sprintf("%s.cos-website.%s.myqcloud.com", d.Id(), meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region)
		website[0]["endpoint"] = endPointUrl
	}
	if err = d.Set("website", website); err != nil {
		return fmt.Errorf("setting website error: %v", err)
	}

	// read the encryption algorithm
	encryption, kmsId, err := cosService.GetBucketEncryption(ctx, bucket, cdcId)
	if err != nil {
		return err
	}
	if err = d.Set("encryption_algorithm", encryption); err != nil {
		return fmt.Errorf("setting encryption error: %v", err)
	}
	if err = d.Set("kms_id", kmsId); err != nil {
		return fmt.Errorf("setting kms_id error: %v", err)
	}

	// read the versioning
	versioning, err := cosService.GetBucketVersioning(ctx, bucket, cdcId)
	if err != nil {
		return err
	}
	if err = d.Set("versioning_enable", versioning); err != nil {
		return fmt.Errorf("setting versioning_enable error: %v", err)
	}

	// read the acceleration
	acceleration, err := cosService.GetBucketAccleration(ctx, bucket, cdcId)
	if err != nil {
		return err
	}
	if err = d.Set("acceleration_enable", acceleration); err != nil {
		return fmt.Errorf("setting acceleration_enable error: %v", err)
	}

	//read the log
	logEnable, logTargetBucket, logPrefix, err := cosService.GetBucketLogStatus(ctx, bucket, cdcId)
	if err != nil {
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() != "UnSupportedLoggingRegion" {
				return err
			}
		}
	} else {
		_ = d.Set("log_enable", logEnable)
		_ = d.Set("log_target_bucket", logTargetBucket)
		_ = d.Set("log_prefix", logPrefix)
	}

	// read the tags
	tags, err := cosService.GetBucketTags(ctx, bucket, cdcId)
	if err != nil {
		return fmt.Errorf("get tags failed: %v", err)
	}
	if len(tags) > 0 {
		_ = d.Set("tags", tags)
	}

	//read intelligent tiering
	result, err := cosService.BucketGetIntelligentTiering(ctx, bucket, cdcId)
	if err != nil {
		return fmt.Errorf("get intelligent tiering failed: %v", err)
	}
	if result != nil {
		if result.Status == "Enabled" {
			_ = d.Set("enable_intelligent_tiering", true)
		} else {
			_ = d.Set("enable_intelligent_tiering", false)
		}
	}
	if result != nil && result.Transition != nil {
		_ = d.Set("intelligent_tiering_days", result.Transition.Days)
		_ = d.Set("intelligent_tiering_request_frequent", result.Transition.RequestFrequent)
	}
	return nil
}

func resourceTencentCloudCosBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	cdcId := d.Get("cdc_id").(string)
	if d.HasChange("enable_intelligent_tiering") || d.HasChange("intelligent_tiering_days") || d.HasChange("intelligent_tiering_request_frequent") {
		old, new := d.GetChange("enable_intelligent_tiering")
		if old.(bool) && !new.(bool) {
			return fmt.Errorf("enable_intelligent_tiering, intelligent_tiering_days and intelligent_tiering_request_frequent not support change!")
		}
		var transition cos.BucketIntelligentTieringTransition
		if v, ok := d.GetOk("intelligent_tiering_days"); ok {
			transition.Days = v.(int)
		} else {
			transition.Days = 30
		}
		if v, ok := d.GetOk("intelligent_tiering_request_frequent"); ok {
			transition.RequestFrequent = v.(int)
		} else {
			transition.RequestFrequent = 1
		}

		if v, ok := d.GetOk("enable_intelligent_tiering"); ok && v.(bool) {
			opt := &cos.BucketPutIntelligentTieringOptions{
				Status:     "Enabled",
				Transition: &transition,
			}
			err := cosService.BucketPutIntelligentTiering(ctx, d.Id(), opt, cdcId)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("acl") {
		bucket := d.Get("bucket").(string)
		err := waitAclEnable(ctx, meta, bucket, cdcId)
		if err != nil {
			return err
		}

		err = resourceTencentCloudCosBucketAclUpdate(ctx, meta, d)
		if err != nil {
			return err
		}
	}

	if d.HasChange("acl_body") {
		body := d.Get("acl_body")
		bucket := d.Get("bucket").(string)
		err := waitAclEnable(ctx, meta, bucket, cdcId)
		if err != nil {
			return err
		}

		if err := resourceTencentCloudCosBucketOriginACLBodyUpdate(ctx, cosService, d); err != nil {
			return err
		}
		_ = d.Set("acl_body", body)
	}

	if d.HasChange("cors_rules") {
		err := resourceTencentCloudCosBucketCorsUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("origin_pull_rules") {
		rules := d.Get("origin_pull_rules")
		err := resourceTencentCloudCosBucketOriginPullUpdate(ctx, cosService, d)
		if err != nil {
			return err
		}
		_ = d.Set("origin_pull_rules", rules)
	}

	if d.HasChange("origin_domain_rules") {
		rules := d.Get("origin_domain_rules")
		if err := resourceTencentCloudCosBucketOriginDomainUpdate(ctx, cosService, d); err != nil {
			return err
		}
		_ = d.Set("origin_domain_rules", rules)
	}

	if d.HasChange("lifecycle_rules") {
		err := resourceTencentCloudCosBucketLifecycleUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("website") {
		err := resourceTencentCloudCosBucketWebsiteUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("encryption_algorithm") || d.HasChange("kms_id") {
		err := resourceTencentCloudCosBucketEncryptionUpdate(ctx, meta, d)
		if err != nil {
			return err
		}
	}

	if d.HasChange("versioning_enable") {
		err := resourceTencentCloudCosBucketVersioningUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("acceleration_enable") {
		err := resourceTencentCloudCosBucketAccelerationUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("replica_role") || d.HasChange("replica_rules") {
		err := resourceTencentCloudCosBucketReplicaUpdate(ctx, cosService, d)

		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		bucket := d.Id()

		cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		if err := cosService.SetBucketTags(ctx, bucket, helper.GetTags(d, "tags"), cdcId); err != nil {
			return err
		}

	}

	if d.HasChange("log_enable") || d.HasChange("log_target_bucket") || d.HasChange("log_prefix") {
		err := resourceTencentCloudCosBucketLogStatusUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	d.Partial(false)

	// wait for update cache
	// if not, the data may be outdated.
	time.Sleep(3 * time.Second)

	return resourceTencentCloudCosBucketRead(d, meta)
}

func waitAclEnable(ctx context.Context, meta interface{}, bucket string, cdcId string) error {
	logId := tccommon.GetLogId(ctx)
	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		aclResult, e := cosService.GetBucketACL(ctx, bucket, cdcId)
		if e != nil {
			if strings.Contains(e.Error(), "NoSuchBucket") {
				log.Printf("[CRITAL][retry]%s api[%s] because of bucket[%s] still on creating, need try again.\n", logId, "GetBucketACL", bucket)
				return resource.RetryableError(fmt.Errorf("[CRITAL][retry]%s api[%s] it still on creating, need try again.\n", logId, "GetBucketACL"))
			}
			log.Printf("[CRITAL]%s api[%s] fail when try to update acl, reason[%s]\n", logId, "GetBucketACL", e.Error())
			return resource.NonRetryableError(e)
		}

		if aclResult == nil {
			return resource.RetryableError(fmt.Errorf("[CRITAL][retry]%s api[%s] it still on creating, need try again.\n", logId, "GetBucketACL"))
		}
		return nil
	})
	return err
}

func resourceTencentCloudCosBucketDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()
	forced := d.Get("force_clean").(bool)
	versioned := d.Get("versioning_enable").(bool)
	cosService := CosService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	cdcId := d.Get("cdc_id").(string)
	err := cosService.DeleteBucket(ctx, bucket, forced, versioned, cdcId)
	if err != nil {
		return err
	}

	// wait for update cache
	// if not, head bucket may be successful
	time.Sleep(3 * time.Second)

	return nil
}

func resourceTencentCloudCosBucketEncryptionUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	encryption := d.Get("encryption_algorithm").(string)
	kmsId := d.Get("kms_id").(string)
	cdcId := d.Get("cdc_id").(string)
	if encryption == "" {
		request := s3.DeleteBucketEncryptionInput{
			Bucket: aws.String(bucket),
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).DeleteBucketEncryption(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket encryption", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket encryption", request.String(), response.String())

		return nil
	}

	request := s3.PutBucketEncryptionInput{
		Bucket: aws.String(bucket),
	}
	request.ServerSideEncryptionConfiguration = &s3.ServerSideEncryptionConfiguration{}
	rules := make([]*s3.ServerSideEncryptionRule, 0)
	defaultRule := &s3.ServerSideEncryptionByDefault{
		SSEAlgorithm:   aws.String(encryption),
		KMSMasterKeyID: aws.String(kmsId),
	}
	rule := &s3.ServerSideEncryptionRule{
		ApplyServerSideEncryptionByDefault: defaultRule,
	}
	rules = append(rules, rule)
	request.ServerSideEncryptionConfiguration.Rules = rules
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketEncryption(&request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket encryption", request.String(), err.Error())
		return fmt.Errorf("cos put bucket encryption error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket encryption", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketVersioningUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	versioning := d.Get("versioning_enable").(bool)
	cdcId := d.Get("cdc_id").(string)
	status := "Suspended"
	if versioning {
		status = "Enabled"
	}
	request := s3.PutBucketVersioningInput{
		Bucket: aws.String(bucket),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String(status),
		},
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketVersioning(&request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket encryption", request.String(), err.Error())
		return fmt.Errorf("cos put bucket encryption error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket encryption", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketAccelerationUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	enabled := d.Get("acceleration_enable").(bool)
	cdcId := d.Get("cdc_id").(string)
	status := "Suspended"
	if enabled {
		status = "Enabled"
	}

	opt := &cos.BucketPutAccelerateOptions{
		Status: status,
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClientNew(bucket, cdcId).Bucket.PutAccelerate(ctx, opt)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, status [%s], reason[%s]\n",
			logId, "put bucket acceleration", opt.Status, err.Error())
		return fmt.Errorf("cos put bucket acceleration error: %s, bucket: %s", err.Error(), bucket)
	}
	rb, _ := ioutil.ReadAll(response.Body)
	body, _ := json.Marshal(rb)
	log.Printf("[DEBUG]%s api[%s] success, status [%s], response body [%s]\n",
		logId, "put bucket acceleration", opt.Status, string(body))

	return err
}

func resourceTencentCloudCosBucketReplicaUpdate(ctx context.Context, service CosService, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	oldRole, newRole := d.GetChange("replica_role")
	oldRules, newRules := d.GetChange("replica_rules")
	cdcId := d.Get("cdc_id").(string)
	oldRuleLength := len(oldRules.([]interface{}))
	newRuleLength := len(newRules.([]interface{}))

	// check if remove
	if oldRole.(string) != "" && newRole.(string) == "" || oldRuleLength > 0 && newRuleLength == 0 {
		result, err := service.GetBucketReplication(ctx, bucket, cdcId)
		if err != nil {
			return err
		}

		if result != nil {
			err := service.DeleteBucketReplication(ctx, d.Get("bucket").(string), cdcId)
			if err != nil {
				return err
			}
		}
	} else if newRole.(string) != "" || newRuleLength > 0 {
		role, rules, _ := getBucketReplications(d)
		err := service.PutBucketReplication(ctx, d.Get("bucket").(string), role, rules, cdcId)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudCosBucketAclUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	cdcId := d.Get("cdc_id").(string)
	request := s3.PutBucketAclInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketAcl(&request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket acl", request.String(), err.Error())
		return fmt.Errorf("cos put bucket error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket acl", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketCorsUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	cors := d.Get("cors_rules").([]interface{})
	cdcId := d.Get("cdc_id").(string)

	if len(cors) == 0 {
		request := s3.DeleteBucketCorsInput{
			Bucket: aws.String(bucket),
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).DeleteBucketCors(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket cors", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket cors error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket cors", request.String(), response.String())
	} else {
		rules := make([]*s3.CORSRule, 0, len(cors))
		for _, item := range cors {
			corsMap := item.(map[string]interface{})
			rule := &s3.CORSRule{}
			for k, v := range corsMap {
				if k == "max_age_seconds" {
					rule.MaxAgeSeconds = aws.Int64(int64(v.(int)))
				} else {
					vMap := make([]*string, len(v.([]interface{})))
					for i, value := range v.([]interface{}) {
						if str, ok := value.(string); ok {
							vMap[i] = aws.String(str)
						}
					}
					switch k {
					case "allowed_origins":
						rule.AllowedOrigins = vMap
					case "allowed_methods":
						rule.AllowedMethods = vMap
					case "allowed_headers":
						rule.AllowedHeaders = vMap
					case "expose_headers":
						rule.ExposeHeaders = vMap
					}
				}
			}
			rules = append(rules, rule)
		}
		request := s3.PutBucketCorsInput{
			Bucket: aws.String(bucket),
			CORSConfiguration: &s3.CORSConfiguration{
				CORSRules: rules,
			},
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketCors(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket cors", request.String(), err.Error())
			return fmt.Errorf("cos put bucket cors error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "put bucket cors", request.String(), response.String())
	}
	return nil
}

func resourceTencentCloudCosBucketLifecycleUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	lifecycleRules := d.Get("lifecycle_rules").([]interface{})
	cdcId := d.Get("cdc_id").(string)
	if len(lifecycleRules) == 0 {
		request := s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucket),
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).DeleteBucketLifecycle(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket lifecycle", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket lifecycle error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket lifecycle", request.String(), response.String())
	} else {
		rules := make([]*s3.LifecycleRule, 0, len(lifecycleRules))
		for i, lifecycleRule := range lifecycleRules {
			r := lifecycleRule.(map[string]interface{})
			rule := &s3.LifecycleRule{}
			id, ok := r["id"].(string)
			if ok {
				rule.ID = &id
			}
			rule.Status = helper.String(s3.ExpirationStatusEnabled)
			prefix := r["filter_prefix"].(string)
			rule.Filter = &s3.LifecycleRuleFilter{
				Prefix: &prefix,
			}

			// Transitions
			transitions := d.Get(fmt.Sprintf("lifecycle_rules.%d.transition", i)).(*schema.Set).List()
			if len(transitions) > 0 {
				rule.Transitions = make([]*s3.Transition, 0, len(transitions))
				for _, transition := range transitions {
					transitionValue := transition.(map[string]interface{})
					t := &s3.Transition{}
					if val, ok := transitionValue["date"].(string); ok && val != "" {
						date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
						if err != nil {
							return fmt.Errorf("parsing cos bucket lifecycle transition date(%s) error: %s", val, err.Error())
						}
						t.Date = aws.Time(date)
					} else if val, ok := transitionValue["days"].(int); ok && val >= 0 {
						t.Days = aws.Int64(int64(val))
					}
					if val, ok := transitionValue["storage_class"].(string); ok && val != "" {
						t.StorageClass = aws.String(val)
					}

					rule.Transitions = append(rule.Transitions, t)
				}
			}

			// Expiration
			expirations := d.Get(fmt.Sprintf("lifecycle_rules.%d.expiration", i)).(*schema.Set).List()
			if len(expirations) > 0 {
				expiration := expirations[0].(map[string]interface{})
				e := &s3.LifecycleExpiration{}

				if val, ok := expiration["date"].(string); ok && val != "" {
					date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
					if err != nil {
						return fmt.Errorf("parsing cos bucket lifecycle expiration data(%s) error: %s", val, err.Error())
					}
					e.Date = aws.Time(date)
				} else if val, ok := expiration["days"].(int); ok && val > 0 {
					e.Days = aws.Int64(int64(val))
				}

				if val, ok := expiration["delete_marker"].(bool); ok && val {
					e.ExpiredObjectDeleteMarker = helper.Bool(true)
				}

				rule.Expiration = e
			}

			// Non Current Transitions
			nonCurrentTransitions := d.Get(fmt.Sprintf("lifecycle_rules.%d.non_current_transition", i)).(*schema.Set).List()
			if len(nonCurrentTransitions) > 0 {
				rule.NoncurrentVersionTransitions = make([]*s3.NoncurrentVersionTransition, 0, len(transitions))
				for _, transition := range nonCurrentTransitions {
					transitionValue := transition.(map[string]interface{})
					t := &s3.NoncurrentVersionTransition{}
					if val, ok := transitionValue["non_current_days"].(int); ok && val >= 0 {
						t.NoncurrentDays = aws.Int64(int64(val))
					}
					if val, ok := transitionValue["storage_class"].(string); ok && val != "" {
						t.StorageClass = aws.String(val)
					}

					rule.NoncurrentVersionTransitions = append(rule.NoncurrentVersionTransitions, t)
				}
			}

			// Non Current Expiration
			nonCurrentExpirations := d.Get(fmt.Sprintf("lifecycle_rules.%d.non_current_expiration", i)).(*schema.Set).List()
			if len(nonCurrentExpirations) > 0 {
				nonCurrentExpiration := nonCurrentExpirations[0].(map[string]interface{})
				e := &s3.NoncurrentVersionExpiration{}

				if val, ok := nonCurrentExpiration["non_current_days"].(int); ok && val > 0 {
					e.NoncurrentDays = aws.Int64(int64(val))
				}

				rule.NoncurrentVersionExpiration = e
			}

			// AbortIncompleteMultipartUpload
			abortIncompleteMultipartUploads := d.Get(fmt.Sprintf("lifecycle_rules.%d.abort_incomplete_multipart_upload", i)).(*schema.Set).List()
			if len(abortIncompleteMultipartUploads) > 0 {
				abortIncompleteMultipartUpload := abortIncompleteMultipartUploads[0].(map[string]interface{})
				e := &s3.AbortIncompleteMultipartUpload{}

				if val, ok := abortIncompleteMultipartUpload["days_after_initiation"].(int); ok && val > 0 {
					e.DaysAfterInitiation = aws.Int64(int64(val))
				}

				rule.AbortIncompleteMultipartUpload = e
			}
			rules = append(rules, rule)
		}

		request := s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucket),
			LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
				Rules: rules,
			},
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketLifecycleConfiguration(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket lifecycle", request.String(), err.Error())
			return fmt.Errorf("cos put bucket lifecycle error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "put bucket lifecycle", request.String(), response.String())
	}

	return nil
}

func resourceTencentCloudCosBucketWebsiteUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	website := d.Get("website").([]interface{})
	cdcId := d.Get("cdc_id").(string)

	if len(website) == 0 {
		request := s3.DeleteBucketWebsiteInput{
			Bucket: aws.String(bucket),
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).DeleteBucketWebsite(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket website", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket website error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket website", request.String(), response.String())
	} else {
		if cdcId != "" {
			return fmt.Errorf("cdc cos not support set website.\n")
		}

		var w map[string]interface{}
		if website[0] != nil {
			w = website[0].(map[string]interface{})
		} else {
			w = make(map[string]interface{})
		}
		var indexDocument, errorDocument string
		var redirectAllRequestsTo = "http"
		if v, ok := w["index_document"]; ok {
			indexDocument = v.(string)
		}
		if v, ok := w["error_document"]; ok {
			errorDocument = v.(string)
		}
		if v, ok := w["redirect_all_requests_to"]; ok && v != "" {
			redirectAllRequestsTo = v.(string)
		}
		endPointUrl := fmt.Sprintf("%s.cos-website.%s.myqcloud.com", d.Id(), meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region)
		request := s3.PutBucketWebsiteInput{
			Bucket: aws.String(bucket),
			WebsiteConfiguration: &s3.WebsiteConfiguration{
				IndexDocument: &s3.IndexDocument{
					Suffix: aws.String(indexDocument),
				},
				ErrorDocument: &s3.ErrorDocument{
					Key: aws.String(errorDocument),
				},
				RedirectAllRequestsTo: &s3.RedirectAllRequestsTo{
					HostName: aws.String(endPointUrl),
					Protocol: aws.String(redirectAllRequestsTo),
				},
			},
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketWebsite(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket website", request.String(), err.Error())
			return fmt.Errorf("cos put bucket website error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "put bucket website", request.String(), response.String())
	}

	return nil
}

func resourceTencentCloudCosBucketLogStatusUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Id()

	logSwitch := d.Get("log_enable").(bool)
	cdcId := d.Get("cdc_id").(string)
	if logSwitch {
		if d.HasChange("log_target_bucket") || d.HasChange("log_prefix") {
			targetBucket := d.Get("log_target_bucket").(string)
			logPrefix := d.Get("log_prefix").(string)
			//check
			if targetBucket == "" || logPrefix == "" {
				return fmt.Errorf("log_target_bucket and log_prefix should set valid value when log_enable is true")
			}

			//set log target bucket and prefix
			//grant are solved by the tencentcloud_cam_role_attachment resource
			request := &s3.PutBucketLoggingInput{
				Bucket: aws.String(bucket),
				BucketLoggingStatus: &s3.BucketLoggingStatus{
					LoggingEnabled: &s3.LoggingEnabled{
						TargetBucket: aws.String(targetBucket),
						TargetPrefix: aws.String(logPrefix),
					},
				},
			}

			resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketLogging(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, "cos enable log error", request.String(), err.Error())
				return fmt.Errorf("cos enable log error: %s, bucket: %s", err.Error(), bucket)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], resp[%s]\n",
				logId, "cos enable log success", request.String(), resp.String())
		}
	} else {
		targetBucket := d.Get("log_target_bucket").(string)
		logPrefix := d.Get("log_prefix").(string)
		//check
		if targetBucket != "" || logPrefix != "" {
			return fmt.Errorf("log_target_bucket and log_prefix should set null when log_enable is false")
		}
		// set disabled, put empty request
		request := &s3.PutBucketLoggingInput{
			Bucket:              aws.String(bucket),
			BucketLoggingStatus: &s3.BucketLoggingStatus{},
		}

		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketLogging(request)
		if err != nil {
			return fmt.Errorf("cos disable log error: %s, bucket: %s", err.Error(), bucket)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], resp[%s]\n",
			logId, "cos enable log success", request.String(), resp.String())
	}

	return nil
}

func resourceTencentCloudCosBucketOriginACLBodyUpdate(ctx context.Context, service CosService, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)
	aclHeader := ""
	aclBody := ""
	body, bodyOk := d.GetOk("acl_body")
	header, headerOk := d.GetOk("acl")
	bucket := d.Get("bucket").(string)
	cdcId := d.Get("cdc_id").(string)
	// If ACLXML update to empty, this will pass default header to delete verbose acl info
	if bodyOk {
		aclBody = body.(string)
	} else if headerOk {
		aclHeader = header.(string)
	} else {
		aclHeader = "private"
	}

	aclBodyOrderly, err := service.transACLBodyOrderly(ctx, aclBody)
	if err != nil {
		return fmt.Errorf("transfer ACL Body failed, reason:%v", err.Error())
	}

	log.Printf("[DEBUG]%s transACLBodyOrderly success, before:[\n%s\n], after:[\n%s\n]\n", logId, aclBody, aclBodyOrderly)

	if err = service.TencentCosPutBucketACLBody(ctx, bucket, aclBodyOrderly, aclHeader, cdcId); err != nil {
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, bucket:[%s]\n", logId, "put bucket acl body", bucket)

	return nil
}

func resourceTencentCloudCosBucketOriginPullUpdate(ctx context.Context, service CosService, d *schema.ResourceData) error {
	var rules []cos.BucketOriginRule
	v, ok := d.GetOk("origin_pull_rules")
	bucket := d.Get("bucket").(string)
	cdcId := d.Get("cdc_id").(string)
	if !ok {
		if err := service.DeleteBucketPullOrigin(ctx, bucket, cdcId); err != nil {
			return err
		}
		return nil
	}
	rulesRaw := v.([]interface{})
	for _, i := range rulesRaw {
		var (
			dMap = i.(map[string]interface{})
			item = &cos.BucketOriginRule{
				OriginCondition: &cos.BucketOriginCondition{
					HTTPStatusCode: "404",
				},
				OriginParameter: &cos.BucketOriginParameter{
					CopyOriginData: true,
					HttpHeader:     &cos.BucketOriginHttpHeader{},
				},
			}
		)

		if v := dMap["sync_back_to_source"]; v.(bool) {
			item.OriginType = "Mirror"
		} else {
			item.OriginType = "Proxy"
		}

		if v, ok := dMap["priority"]; ok {
			item.RulePriority = v.(int)
		}
		if v, ok := dMap["prefix"]; ok {
			item.OriginCondition.Prefix = v.(string)
		}
		if v, ok := dMap["protocol"]; ok {
			item.OriginParameter.Protocol = v.(string)
		}
		if v, ok := dMap["host"]; ok {
			item.OriginInfo.HostInfo.HostName = v.(string)
		}
		if v, ok := dMap["follow_query_string"]; ok {
			item.OriginParameter.FollowQueryString = v.(bool)
		}
		if v, ok := dMap["follow_redirection"]; ok {
			item.OriginParameter.FollowRedirection = v.(bool)
		}
		//if v, ok := dMap["copy_origin_data"]; ok {
		//	item.OriginParameter.CopyOriginData = v.(bool)
		//}
		if v, ok := dMap["redirect_prefix"]; ok {
			value := v.(string)
			item.OriginInfo.FileInfo.PrefixConfiguration.Prefix = value
		}
		if v, ok := dMap["redirect_suffix"]; ok {
			value := v.(string)
			item.OriginInfo.FileInfo.SuffixConfiguration.Suffix = value
		}
		if v, ok := dMap["custom_http_headers"]; ok {
			var customHeaders []cos.OriginHttpHeader
			for key, val := range v.(map[string]interface{}) {
				customHeaders = append(customHeaders, cos.OriginHttpHeader{
					Key:   key,
					Value: val.(string),
				})
			}
			item.OriginParameter.HttpHeader.NewHttpHeaders = customHeaders
		}
		if v, ok := dMap["follow_http_headers"]; ok {
			var followHeaders []cos.OriginHttpHeader
			for _, item := range v.(*schema.Set).List() {
				header := cos.OriginHttpHeader{
					Key:   item.(string),
					Value: "",
				}
				followHeaders = append(followHeaders, header)
			}
			item.OriginParameter.HttpHeader.FollowHttpHeaders = followHeaders
		}
		rules = append(rules, *item)
	}

	if err := service.PutBucketPullOrigin(ctx, bucket, rules, cdcId); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudCosBucketOriginDomainUpdate(ctx context.Context, service CosService, d *schema.ResourceData) error {
	v, ok := d.GetOk("origin_domain_rules")
	bucket := d.Get("bucket").(string)
	cdcId := d.Get("cdc_id").(string)
	if !ok {
		if err := service.DeleteBucketOriginDomain(ctx, bucket, cdcId); err != nil {
			return err
		}
		return nil
	}
	rules := v.([]interface{})
	domainRules := make([]cos.BucketDomainRule, 0)

	for _, rule := range rules {
		dMap := rule.(map[string]interface{})
		item := cos.BucketDomainRule{}
		if name, ok := dMap["domain"]; ok {
			item.Name = name.(string)
		}
		if status, ok := dMap["status"]; ok {
			item.Status = status.(string)
		}
		if domainType, ok := dMap["type"]; ok {
			item.Type = domainType.(string)
		}
		domainRules = append(domainRules, item)
	}

	if err := service.PutBucketOriginDomain(ctx, bucket, domainRules, cdcId); err != nil {
		return err
	}
	return nil
}

func getBucketPutOptions(d *schema.ResourceData) (useCosService bool, options *cos.BucketPutOptions) {
	opt := &cos.BucketPutOptions{
		XCosACL:              d.Get("acl").(string),
		XCosGrantRead:        "",
		XCosGrantWrite:       "",
		XCosGrantReadACP:     "",
		XCosGrantWriteACP:    "",
		XCosGrantFullControl: "",
	}
	grants, hasGrantHeaders := d.GetOk("grant_headers")
	maz, hasMAZ := d.GetOk("multi_az")

	if !hasGrantHeaders && !hasMAZ {
		return false, opt
	}

	if hasGrantHeaders {
		headers := grants.(map[string]interface{})
		if v, ok := headers["grant_read"]; ok {
			opt.XCosGrantRead = v.(string)
		}
		if v, ok := headers["grant_write"]; ok {
			opt.XCosGrantWrite = v.(string)
		}
		if v, ok := headers["grant_read_acp"]; ok {
			opt.XCosGrantReadACP = v.(string)
		}
		if v, ok := headers["grant_write_acp"]; ok {
			opt.XCosGrantWriteACP = v.(string)
		}
		if v, ok := headers["grant_full_control"]; ok {
			opt.XCosGrantFullControl = v.(string)
		}
	}

	if hasMAZ {
		if maz.(bool) {
			opt.CreateBucketConfiguration = &cos.CreateBucketConfiguration{
				BucketAZConfig: "MAZ",
			}
		}
	}

	return true, opt
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return helper.HashString(buf.String())
}

func nonCurrentExpirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["non_current_days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return helper.HashString(buf.String())
}

func abortIncompleteMultipartUploadHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["days_after_initiation"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return helper.HashString(buf.String())
}

func transitionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return helper.HashString(buf.String())
}

func nonCurrentTransitionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["non_current_days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return helper.HashString(buf.String())
}

func getBucketReplications(d *schema.ResourceData) (role string, rules []cos.BucketReplicationRule, err error) {
	role = d.Get("replica_role").(string)
	replicaRules := d.Get("replica_rules").([]interface{})
	for i := range replicaRules {
		item := replicaRules[i].(map[string]interface{})
		rule := cos.BucketReplicationRule{
			Status: item["status"].(string),
			Destination: &cos.ReplicationDestination{
				Bucket: item["destination_bucket"].(string),
			},
		}
		if v, ok := item["prefix"].(string); ok {
			rule.Prefix = v
		}
		if v, ok := item["id"].(string); ok {
			rule.ID = v
		}
		if v, ok := item["destination_storage_class"].(string); ok {
			rule.Destination.StorageClass = v
		}
		rules = append(rules, rule)
	}
	return
}

func setBucketReplication(d *schema.ResourceData, result cos.GetBucketReplicationResult) (err error) {
	if result.Role != "" {
		_ = d.Set("replica_role", result.Role)
	}
	rules := make([]map[string]interface{}, 0)
	if len(result.Rule) > 0 {
		for i := range result.Rule {
			item := result.Rule[i]
			rule := map[string]interface{}{
				"status":                    item.Status,
				"destination_bucket":        item.Destination.Bucket,
				"destination_storage_class": item.Destination.StorageClass,
			}
			if item.ID != "" {
				rule["id"] = item.ID
			}
			if item.Prefix != "" {
				rule["prefix"] = item.Prefix
			}
			rules = append(rules, rule)
		}
	}
	err = d.Set("replica_rules", rules)
	return
}

func ACLBodyDiffFunc(olds, news string, d *schema.ResourceData) (result bool) {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.ACLBodyDiffFunc")()
	log.Printf("[DEBUG] ACLBodyDiffFunc called, before:[\n%s\n], after:[\n%s\n]\n", olds, news)

	oldDoc := etree.NewDocument()
	newDoc := etree.NewDocument()

	if err := oldDoc.ReadFromString(olds); err != nil {
		log.Printf("[CRITAL]read old xml from string error: %v", err)
		return false
	}

	if err := newDoc.ReadFromString(news); err != nil {
		log.Printf("[CRITAL]read new xml from string error: %v", err)
		return false
	}

	oldRoot := oldDoc.SelectElement("AccessControlPolicy")
	newRoot := newDoc.SelectElement("AccessControlPolicy")

	if oldRoot == nil || newRoot == nil {
		log.Println("[CRITAL]oldRoot or newRoot is nil: return false.")
		return false
	}

	oldOwner := oldRoot.SelectElement("Owner")
	newOwner := newRoot.SelectElement("Owner")

	if oldOwner == nil || newOwner == nil {
		log.Println("[CRITAL]oldOwner or newOwner is nil: return false.")
		return false
	}

	oldOwnerId := oldOwner.SelectElement("ID")
	oldOwnerName := oldOwner.SelectElement("DisplayName")
	newOwnerId := newOwner.SelectElement("ID")
	newOwnerName := newOwner.SelectElement("DisplayName")

	// diff: Owner element
	if oldOwnerId.Text() != newOwnerId.Text() || oldOwnerName.Text() != newOwnerName.Text() {
		log.Printf("[CRITAL]OwnerId[old:%s, new:%s] or OwnerName[old:%s, new:%s] not equal: return false.\n", oldOwnerId.Text(), newOwnerId.Text(), oldOwnerName.Text(), newOwnerName.Text())
		return false
	}

	// diff check: owner display name(if have)
	if oldOwnerName != nil {
		if newOwnerName == nil {
			log.Println("[CRITAL]newOwnerName is nil: return false.")
			return false
		}
		if oldOwnerName.Text() != newOwnerName.Text() {
			log.Printf("[CRITAL]OwnerName[old:%s, new:%s] not equal: return false.\n", oldOwnerName.Text(), newOwnerName.Text())
			return false
		}
	}

	// diff: ACL element
	oldGrantees := oldRoot.FindElements("//Grantee")
	newGrantees := newRoot.FindElements("//Grantee")
	// check count
	if len(oldGrantees) != len(newGrantees) {
		return false
	}
	// check content
	for _, oldGrantee := range oldGrantees {
		for _, attr := range oldGrantee.Attr {
			if attr.Key != "type" {
				// only need to handle the type attribute
				continue
			}
			// anonymous or real user
			oldGranteeType := attr.Value

			oldGranteeID := oldGrantee.SelectElement("ID")
			oldGranteeURI := oldGrantee.SelectElement("URI")
			oldGranteeDisplayName := oldGrantee.SelectElement("DisplayName")
			oldGrant := oldGrantee.Parent()
			oldGrantPermission := oldGrant.SelectElement("Permission")

			// find the new grant permission by specified grantee type
			result = false
			for _, newGrantee := range newRoot.FindElements(fmt.Sprintf("//Grantee[@type='%s']", oldGranteeType)) {
				newGranteeID := newGrantee.SelectElement("ID")
				newGranteeURI := newGrantee.SelectElement("URI")
				newGranteeDisplayName := newGrantee.SelectElement("DisplayName")
				newGrant := newGrantee.Parent()
				newGrantPermission := newGrant.SelectElement("Permission")

				// diff check: grantee id and name for real user
				if oldGranteeType == COS_ACL_GRANTEE_TYPE_USER {
					if oldGranteeID == nil || newGranteeID == nil {
						continue
					}
					if oldGranteeID.Text() != newGranteeID.Text() {
						continue
					}

					// diff check: grantee display name(if have)
					if oldGranteeDisplayName != nil {
						if newGranteeDisplayName == nil {
							continue
						}
						if oldGranteeDisplayName.Text() != newGranteeDisplayName.Text() {
							continue
						}
					}
				}

				// diff check: grantee uri for anonymous
				if oldGranteeType == COS_ACL_GRANTEE_TYPE_ANONYMOUS {
					if oldGranteeURI == nil || newGranteeURI == nil {
						continue
					}
					if oldGranteeURI.Text() != newGranteeURI.Text() {
						continue
					}
				}

				// diff check: permission
				if oldGrantPermission == nil || newGrantPermission == nil {
					continue
				}
				if oldGrantPermission.Text() != newGrantPermission.Text() {
					continue
				}

				// congrats! passed all diff checks for this grant.
				result = true

				var uid string
				if oldGranteeType == COS_ACL_GRANTEE_TYPE_USER {
					uid = oldGranteeID.Text()
				} else {
					uid = oldGranteeURI.Text()
				}
				log.Printf("[DEBUG] diff verification passed for grantee:[%s:%s]\n", oldGranteeType, uid)
			}
			if !result {
				return false
			}
		}
	}
	log.Printf("[DEBUG] Owner:%s's final equation result between old and new ACL is:[%v]\n", oldOwnerId.Text(), result)
	return result
}
