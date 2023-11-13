/*
Use this data source to query detailed information of mongodb instance_current_op

Example Usage

```hcl
data "tencentcloud_mongodb_instance_current_op" "instance_current_op" {
  instance_id = "cmgo-9d0p6umb"
  ns = ""
  millisecond_running = 10
  op = "update"
  replica_set_name = ""
  state = "secondary"
  limit = 10
  offset = 0
  order_by = ""
  order_by_type = "desc"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMongodbInstanceCurrentOp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceCurrentOpRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.",
			},

			"ns": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter condition, the namespace namespace to which the operation belongs, in the format of db.collection.",
			},

			"millisecond_running": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Filter condition, the time that the operation has been executed (unit: millisecond),the result will return the operation that exceeds the set time, the default value is 0,and the value range is [0, 3600000].",
			},

			"op": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter condition, operation type, possible values: none, update, insert, query, command, getmore,remove and killcursors.",
			},

			"replica_set_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter condition, shard name.",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter condition, node status, possible value: primary, secondary.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number returned by a single request, the default value is 100, and the value range is [0,100].",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset, the default value is 0, and the value range is [0,10000].",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Returns the sorted field of the result set, currently supports: MicrosecsRunning/microsecsrunning,the default is ascending sort.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Returns the sorting method of the result set, possible values: ASC/asc or DESC/desc.",
			},

			"current_ops": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Current operation list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"op_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Operation id.",
						},
						"ns": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operation namespace.",
						},
						"query": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operation query.",
						},
						"op": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operation value.",
						},
						"replica_set_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Replication name.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operation state.",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operation info.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name.",
						},
						"microsecs_running": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Running time(ms).",
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

func dataSourceTencentCloudMongodbInstanceCurrentOpRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mongodb_instance_current_op.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ns"); ok {
		paramMap["Ns"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("millisecond_running"); v != nil {
		paramMap["MillisecondRunning"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("op"); ok {
		paramMap["Op"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("replica_set_name"); ok {
		paramMap["ReplicaSetName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		paramMap["State"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	service := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var currentOps []*mongodb.CurrentOp

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceCurrentOpByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		currentOps = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(currentOps))
	tmpList := make([]map[string]interface{}, 0, len(currentOps))

	if currentOps != nil {
		for _, currentOp := range currentOps {
			currentOpMap := map[string]interface{}{}

			if currentOp.OpId != nil {
				currentOpMap["op_id"] = currentOp.OpId
			}

			if currentOp.Ns != nil {
				currentOpMap["ns"] = currentOp.Ns
			}

			if currentOp.Query != nil {
				currentOpMap["query"] = currentOp.Query
			}

			if currentOp.Op != nil {
				currentOpMap["op"] = currentOp.Op
			}

			if currentOp.ReplicaSetName != nil {
				currentOpMap["replica_set_name"] = currentOp.ReplicaSetName
			}

			if currentOp.State != nil {
				currentOpMap["state"] = currentOp.State
			}

			if currentOp.Operation != nil {
				currentOpMap["operation"] = currentOp.Operation
			}

			if currentOp.NodeName != nil {
				currentOpMap["node_name"] = currentOp.NodeName
			}

			if currentOp.MicrosecsRunning != nil {
				currentOpMap["microsecs_running"] = currentOp.MicrosecsRunning
			}

			ids = append(ids, *currentOp.InstanceId)
			tmpList = append(tmpList, currentOpMap)
		}

		_ = d.Set("current_ops", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
