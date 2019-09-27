/*
Provides a resource for an AS (Auto scaling) schedule.

Example Usage

```hcl
resource "tencentcloud_as_schedule" "schedule" {
	scaling_group_id = "sg-12af45"
	schedule_action_name = "tf-as-schedule"
	max_size = 10
	min_size = 0
	desired_capacity = 0
	start_time = "2019-01-01T00:00:00+08:00"
	end_time = "2019-12-01T00:00:00+08:00"
	recurrence = "0 0 * * *"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func resourceTencentCloudAsSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScheduleCreate,
		Read:   resourceTencentCloudAsScheduleRead,
		Update: resourceTencentCloudAsScheduleUpdate,
		Delete: resourceTencentCloudAsScheduleDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a scaling group.",
			},
			"schedule_action_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "The name of this scaling action.",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The maximum size for the Auto Scaling group.",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The minimum size for the Auto Scaling group.",
			},
			"desired_capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The desired number of CVM instances that should be running in the group.",
			},
			"start_time": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAsScheduleTimestamp,
				Description:  "The time for this action to start, in \"YYYY-MM-DDThh:mm:ss+08:00\" format (UTC+8).",
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAsScheduleTimestamp,
				Description:  "The time for this action to end, in \"YYYY-MM-DDThh:mm:ss+08:00\" format (UTC+8).",
			},
			"recurrence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The time when recurring future actions will start. Start time is specified by the user following the Unix cron syntax format. And this argument should be set with end_time together.",
			},
		},
	}
}

func resourceTencentCloudAsScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_schedule.create")()

	logId := getLogId(contextNil)

	request := as.NewCreateScheduledActionRequest()
	request.AutoScalingGroupId = stringToPointer(d.Get("scaling_group_id").(string))
	request.ScheduledActionName = stringToPointer(d.Get("schedule_action_name").(string))
	request.MaxSize = intToPointer(d.Get("max_size").(int))
	request.MinSize = intToPointer(d.Get("min_size").(int))
	request.DesiredCapacity = intToPointer(d.Get("desired_capacity").(int))
	request.StartTime = stringToPointer(d.Get("start_time").(string))

	// end_time and recurrence must be specified at the same time
	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = stringToPointer(v.(string))
		if vv, ok := d.GetOk("recurrence"); ok {
			request.Recurrence = stringToPointer(vv.(string))
		} else {
			return fmt.Errorf("end_time and recurrence must be specified at the same time.")
		}
	} else {
		if _, ok := d.GetOk("recurrence"); ok {
			return fmt.Errorf("end_time and recurrence must be specified at the same time.")
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().CreateScheduledAction(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ScheduledActionId == nil {
		return fmt.Errorf("schedule action id is nil")
	}
	d.SetId(*response.Response.ScheduledActionId)

	return resourceTencentCloudAsScheduleRead(d, meta)
}

func resourceTencentCloudAsScheduleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_schedule.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scheduledActionId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		scheduledAction, e := asService.DescribeScheduledActionById(ctx, scheduledActionId)
		if e != nil {
			if sdkErr, ok := e.(*errors.TencentCloudSDKError); ok {
				if sdkErr.Code == AsScheduleNotFound {
					d.SetId("")
					return nil
				}
			}
			return retryError(e)
		}

		d.Set("scaling_group_id", *scheduledAction.AutoScalingGroupId)
		d.Set("schedule_action_name", *scheduledAction.ScheduledActionName)
		d.Set("max_size", *scheduledAction.MaxSize)
		d.Set("min_size", *scheduledAction.MinSize)
		d.Set("desired_capacity", *scheduledAction.DesiredCapacity)
		d.Set("start_time", *scheduledAction.StartTime)

		if scheduledAction.EndTime != nil {
			d.Set("end_time", *scheduledAction.EndTime)
		}
		if scheduledAction.Recurrence != nil {
			d.Set("recurrence", *scheduledAction.Recurrence)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_schedule.update")()

	logId := getLogId(contextNil)

	request := as.NewModifyScheduledActionRequest()
	scheduledActionId := d.Id()
	request.ScheduledActionId = &scheduledActionId
	if d.HasChange("schedule_action_name") {
		request.ScheduledActionName = stringToPointer(d.Get("schedule_action_name").(string))
	}
	if d.HasChange("max_size") {
		request.MaxSize = intToPointer(d.Get("max_size").(int))
	}
	if d.HasChange("min_size") {
		request.MinSize = intToPointer(d.Get("min_size").(int))
	}
	if d.HasChange("desired_capacity") {
		request.DesiredCapacity = intToPointer(d.Get("desired_capacity").(int))
	}
	if d.HasChange("start_time") {
		request.StartTime = stringToPointer(d.Get("start_time").(string))
	}
	if d.HasChange("end_time") {
		request.EndTime = stringToPointer(d.Get("end_time").(string))
		request.Recurrence = stringToPointer(d.Get("recurrence").(string))
	}
	if d.HasChange("recurrence") {
		request.Recurrence = stringToPointer(d.Get("recurrence").(string))
		request.EndTime = stringToPointer(d.Get("end_time").(string))
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().ModifyScheduledAction(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func resourceTencentCloudAsScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_schedule.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	scheduledActionId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := asService.DeleteScheduledAction(ctx, scheduledActionId)
	if err != nil {
		return err
	}

	return nil
}
