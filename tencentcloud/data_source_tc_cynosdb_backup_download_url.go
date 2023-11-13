/*
Use this data source to query detailed information of cynosdb backup_download_url

Example Usage

```hcl
data "tencentcloud_cynosdb_backup_download_url" "backup_download_url" {
  cluster_id = "cynosdbmysql-123"
  backup_id = 100
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbBackupDownloadUrl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbBackupDownloadUrlRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"backup_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Backup ID.",
			},

			"download_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup download address.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCynosdbBackupDownloadUrlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_backup_download_url.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("backup_id"); v != nil {
		paramMap["BackupId"] = helper.IntInt64(v.(int))
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbBackupDownloadUrlByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		downloadUrl = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(downloadUrl))
	if downloadUrl != nil {
		_ = d.Set("download_url", downloadUrl)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
