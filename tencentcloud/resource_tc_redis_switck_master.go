/*
Provides a resource to create a redis switck_master

Example Usage

```hcl
resource "tencentcloud_redis_switck_master" "switck_master" {
  instance_id = "crs-kfdkirid"
  group_id = 29369
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRedisSwitckMaster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisSwitckMasterCreate,
		Read:   resourceTencentCloudRedisSwitckMasterRead,
		Update: resourceTencentCloudRedisSwitckMasterUpdate,
		Delete: resourceTencentCloudRedisSwitckMasterDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Replication group ID, required for multi-AZ instances.",
			},
		},
	}
}

func resourceTencentCloudRedisSwitckMasterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_switck_master.create")()
	defer inconsistentCheck(d, meta)()

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudRedisSwitckMasterUpdate(d, meta)
}

func resourceTencentCloudRedisSwitckMasterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_switck_master.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()
	paramMap := make(map[string]interface{})
	paramMap["InstanceId"] = &instanceId

	switckMaster, err := service.DescribeRedisInstanceZoneInfoByFilter(ctx, paramMap)
	if err != nil {
		return err
	}

	if switckMaster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `RedisSwitckMaster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if len(switckMaster) > 1 {
		for _, v := range switckMaster {
			if *v.Role == "master" {
				_ = d.Set("group_id", v.GroupId)
				break
			}
		}
	}

	return nil
}

func resourceTencentCloudRedisSwitckMasterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_switck_master.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := redis.NewChangeReplicaToMasterRequest()
	response := redis.NewChangeReplicaToMasterResponse()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRedisClient().ChangeReplicaToMaster(request)
		if e != nil {
			if _, ok := e.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(e)
			} else {
				return resource.NonRetryableError(e)
			}
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update redis switckMaster failed, reason:%+v", logId, err)
		return err
	}

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	taskId := *response.Response.TaskId
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeTaskInfo(ctx, instanceId, taskId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("update redis switckMaster is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s update redis switckMaster fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudRedisSwitckMasterRead(d, meta)
}

func resourceTencentCloudRedisSwitckMasterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_redis_switck_master.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
