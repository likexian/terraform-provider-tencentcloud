package cos

import (
	"context"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCosObjectSignedUrl() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceTencentCloudCosObjectSignedUrlRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the bucket.",
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The full path to the object inside the bucket.",
			},
			"method": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "GET",
				ValidateFunc: validation.StringInSlice([]string{"GET", "PUT"}, true),
				Description:  "Method, GET or PUT. Default value is GET.",
			},
			"duration": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1m",
				Description: "Duration of signed url. Default value is 1m.",
			},
			"headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request headers.",
			},
			"queries": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Request query parameters.",
			},
			"signed_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Signed URL.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func DataSourceTencentCloudCosObjectSignedUrlRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cos_object_signed_url.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Get("bucket").(string)
	path := d.Get("path").(string)
	method := "GET"
	durationString := "1m"
	opt := &cos.PresignedURLOptions{}
	signHost := true

	if v, ok := d.GetOk("method"); ok {
		method = v.(string)
	}

	if v, ok := d.GetOk("duration"); ok {
		durationString = v.(string)
	}

	duration, err := time.ParseDuration(durationString)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("headers"); ok {
		for key, value := range v.(map[string]string) {
			opt.Header.Set(key, value)
		}
	}

	if v, ok := d.GetOk("queries"); ok {
		for key, value := range v.(map[string]string) {
			opt.Query.Set(key, value)
		}
	}

	result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Object.GetPresignedURL2(ctx, method, path, duration, opt, signHost)
	if err != nil {
		return err
	}

	signedUrl := result.String()

	d.SetId(helper.DataResourceIdHash(path))
	_ = d.Set("signed_url", signedUrl)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), signedUrl); err != nil {
			return err
		}
	}

	return nil
}
