/*
Provides a resource to create a mariadb instance_config

Example Usage

```hcl
resource "tencentcloud_mariadb_instance_config" "test" {
  instance_id        = "tdsql-9vqvls95"
  vpc_id             = "vpc-ii1jfbhl"
  subnet_id          = "subnet-3ku415by"
  rs_access_strategy = 1
  extranet_access    = 0
  project_id         = 0
  vip                = "127.0.0.1"
}
```

Import

mariadb instance_config can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_instance_config.test id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMariadbInstanceConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbInstanceConfigCreate,
		Read:   resourceTencentCloudMariadbInstanceConfigRead,
		Update: resourceTencentCloudMariadbInstanceConfigUpdate,
		Delete: resourceTencentCloudMariadbInstanceConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},
			"rs_access_strategy": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "RS proximity mode, 0- no strategy, 1- access to the nearest available zone.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Project ID, if not passed, indicates the default project.",
			},
			"extranet_access": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "External network status, 0-closed; 1- Opening; Default not enabled.",
			},
			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "vip.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "subnet id, it&amp;#39;s required when vpcId is set; If modifying VIP, this parameter is required.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "vpc id; If modifying VIP, this parameter is required.",
			},
		},
	}
}

func resourceTencentCloudMariadbInstanceConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMariadbInstanceConfigUpdate(d, meta)
}

func resourceTencentCloudMariadbInstanceConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance_config.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId = d.Id()
	)

	dbDetail, err := service.DescribeDBInstanceDetailById(ctx, instanceId)
	if err != nil {
		return err
	}

	if dbDetail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MariadbInstanceConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbDetail.InstanceId != nil {
		_ = d.Set("instance_id", dbDetail.InstanceId)
	}

	if dbDetail.ProjectId != nil {
		_ = d.Set("project_id", dbDetail.ProjectId)
	}

	if dbDetail.RsAccessStrategy != nil {
		_ = d.Set("rs_access_strategy", dbDetail.RsAccessStrategy)
	}

	if dbDetail.WanStatus != nil {
		_ = d.Set("extranet_access", dbDetail.WanStatus)
	}

	if dbDetail.Vip != nil {
		_ = d.Set("vip", dbDetail.Vip)
	}

	if dbDetail.SubnetId != nil {
		_ = d.Set("subnet_id", dbDetail.SubnetId)
	}

	if dbDetail.VpcId != nil {
		_ = d.Set("vpc_id", dbDetail.VpcId)
	}

	return nil
}

func resourceTencentCloudMariadbInstanceConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance_config.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request    = mariadb.NewModifyAccountPrivilegesRequest()
		instanceId = d.Id()
	)

	needChange := false

	mutableArgs := []string{"rs_access_strategy", "project_id", "extranet_access", "vip", "subnet_id", "vpc_id"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		// set rs_access_strategy
		if v, ok := d.GetOkExists("rs_access_strategy"); ok {
			rsAccessStrategy := int64(v.(int))
			if rsAccessStrategy == RSACCESSSTRATEGY_ENABLE {
				rsAccessStrategyRequest := mariadb.NewModifyRealServerAccessStrategyRequest()
				rsAccessStrategyRequest.InstanceId = &instanceId
				rsAccessStrategyRequest.RsAccessStrategy = &rsAccessStrategy

				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyRealServerAccessStrategy(rsAccessStrategyRequest)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb accessStrategy failed, reason:%+v", logId, err)
					return err
				}
			}
		}

		// set ExtranetAccess
		if v, ok := d.GetOkExists("extranet_access"); ok {
			extranetAccess := v.(int)
			var extranetAccessFlowId int64
			if extranetAccess == ExtranetAccess_ENABLE {
				extranetAccessRequest := mariadb.NewOpenDBExtranetAccessRequest()
				extranetAccessRequest.InstanceId = &instanceId
				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().OpenDBExtranetAccess(extranetAccessRequest)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					extranetAccessFlowId = *result.Response.FlowId
					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb openDBExtranetAccess failed, reason:%+v", logId, err)
					return err
				}

			} else {
				extranetAccessRequest := mariadb.NewCloseDBExtranetAccessRequest()
				extranetAccessRequest.InstanceId = &instanceId
				err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
					result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CloseDBExtranetAccess(extranetAccessRequest)
					if e != nil {
						return retryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					extranetAccessFlowId = *result.Response.FlowId
					return nil
				})
				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb closeDBExtranetAccess failed, reason:%+v", logId, err)
					return err
				}
			}

			// wait
			if extranetAccessFlowId != NONE_FLOW_TASK {
				err := resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
					result, e := service.DescribeFlowById(ctx, extranetAccessFlowId)
					if e != nil {
						return retryError(e)
					}

					if *result.Status == MARIADB_TASK_SUCCESS {
						return nil
					} else if *result.Status == MARIADB_TASK_RUNNING {
						return resource.RetryableError(fmt.Errorf("operate mariadb DBExtranetAccess status is running"))
					} else if *result.Status == MARIADB_TASK_FAIL {
						return resource.NonRetryableError(fmt.Errorf("operate mariadb DBExtranetAccess status is fail"))
					} else {
						e = fmt.Errorf("operate mariadb DBExtranetAccess status illegal")
						return resource.NonRetryableError(e)
					}
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb DBExtranetAccess task failed, reason:%+v", logId, err)
					return err
				}
			}
		}

		// set vip
		if v, ok := d.GetOk("vip"); ok {
			Vip := v.(string)
			var VipFlowId int64
			VipRequest := mariadb.NewModifyInstanceNetworkRequest()
			VipRequest.InstanceId = &Vip
			if v, ok := d.GetOk("vpc_id"); ok {
				VipRequest.VpcId = helper.String(v.(string))
			}

			if v, ok := d.GetOk("subnet_id"); ok {
				VipRequest.SubnetId = helper.String(v.(string))
			}

			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyInstanceNetwork(VipRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				VipFlowId = *result.Response.FlowId
				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate mariadb network failed, reason:%+v", logId, err)
				return err
			}

			// wait
			if VipFlowId != NONE_FLOW_TASK {
				err = resource.Retry(10*writeRetryTimeout, func() *resource.RetryError {
					result, e := service.DescribeFlowById(ctx, VipFlowId)
					if e != nil {
						return retryError(e)
					}

					if *result.Status == MARIADB_TASK_SUCCESS {
						return nil
					} else if *result.Status == MARIADB_TASK_RUNNING {
						return resource.RetryableError(fmt.Errorf("operate mariadb network status is running"))
					} else if *result.Status == MARIADB_TASK_FAIL {
						return resource.NonRetryableError(fmt.Errorf("operate mariadb network status is fail"))
					} else {
						e = fmt.Errorf("operate mariadb network status illegal")
						return resource.NonRetryableError(e)
					}
				})

				if err != nil {
					log.Printf("[CRITAL]%s operate mariadb network task failed, reason:%+v", logId, err)
					return err
				}
			}
		}

		// set project_id
		if v, ok := d.GetOk("project_id"); ok {
			projectId := int64(v.(int))

			MPRequest := mariadb.NewModifyDBInstancesProjectRequest()
			MPRequest.InstanceIds = common.StringPtrs([]string{instanceId})
			MPRequest.ProjectId = &projectId

			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyDBInstancesProject(MPRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s operate mariadb modifyInstanceProject failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudMariadbInstanceConfigRead(d, meta)
}

func resourceTencentCloudMariadbInstanceConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_instance_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
