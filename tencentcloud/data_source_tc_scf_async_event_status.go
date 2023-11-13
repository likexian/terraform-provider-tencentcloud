/*
Use this data source to query detailed information of scf async_event_status

Example Usage

```hcl
data "tencentcloud_scf_async_event_status" "async_event_status" {
  invoke_request_id = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfAsyncEventStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfAsyncEventStatusRead,
		Schema: map[string]*schema.Schema{
			"invoke_request_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the async execution request.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Async event status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Async event status. Values: `RUNNING` (running); `FINISHED` (invoked successfully); `ABORTED` (invocation ended); `FAILED` (invocation failed).",
						},
						"status_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Request status code.",
						},
						"invoke_request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Async execution request ID.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudScfAsyncEventStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_async_event_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("invoke_request_id"); ok {
		paramMap["InvokeRequestId"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*scf.AsyncEventStatus

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfAsyncEventStatusByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		asyncEventStatusMap := map[string]interface{}{}

		if result.Status != nil {
			asyncEventStatusMap["status"] = result.Status
		}

		if result.StatusCode != nil {
			asyncEventStatusMap["status_code"] = result.StatusCode
		}

		if result.InvokeRequestId != nil {
			asyncEventStatusMap["invoke_request_id"] = result.InvokeRequestId
		}

		ids = append(ids, *result.InvokeRequestId)
		_ = d.Set("result", asyncEventStatusMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), asyncEventStatusMap); e != nil {
			return e
		}
	}
	return nil
}
