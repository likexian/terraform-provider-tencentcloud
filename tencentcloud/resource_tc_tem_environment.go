/*
Provides a resource to create a tem environment

Example Usage

```hcl
resource "tencentcloud_tem_environment" "environment" {
  environment_name = "demo"
  description      = "demo for test"
  vpc              = "vpc-2hfyray3"
  subnet_ids       = ["subnet-rdkj0agk", "subnet-r1c4pn5m", "subnet-02hcj95c"]
  tag {
    tag_key = "createdBy"
	tag_value = "terraform"
  }
}

```
Import

tem environment can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_environment.environment environment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTemEnvironment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemEnvironmentRead,
		Create: resourceTencentCloudTemEnvironmentCreate,
		Update: resourceTencentCloudTemEnvironmentUpdate,
		Delete: resourceTencentCloudTemEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "environment name.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "environment description.",
			},

			"vpc": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vpc ID.",
			},

			"subnet_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "subnet IDs.",
			},
			"tag": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "environment tag list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_environment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tem.NewCreateEnvironmentRequest()
		response *tem.CreateEnvironmentResponse
	)

	if v, ok := d.GetOk("environment_name"); ok {
		request.EnvironmentName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc"); ok {
		request.Vpc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIdsSet := v.(*schema.Set).List()
		for i := range subnetIdsSet {
			subnetIds := subnetIdsSet[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetIds)
		}
	}

	if v, ok := d.GetOk("tag"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tag := tem.Tag{}
			if v, ok := dMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateEnvironment(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tem environment failed, reason:%+v", logId, err)
		return err
	}

	environmentId := *response.Response.Result

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTemEnvironmentStatus(ctx, environmentId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.ClusterStatus == "NORMAL" {
			return nil
		}
		if *instance.ClusterStatus == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("environment status is %v, operate failed.", *instance.ClusterStatus))
		}
		return resource.RetryableError(fmt.Errorf("environment status is %v, retry...", *instance.ClusterStatus))
	})
	if err != nil {
		return err
	}

	d.SetId(environmentId)
	return resourceTencentCloudTemEnvironmentRead(d, meta)
}

func resourceTencentCloudTemEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_environment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	environmentId := d.Id()

	environments, err := service.DescribeTemEnvironment(ctx, environmentId)

	if err != nil {
		return err
	}
	environment := environments.Result
	if environment == nil {
		d.SetId("")
		return fmt.Errorf("resource `environment` %s does not exist", environmentId)
	}

	if environment.EnvironmentName != nil {
		_ = d.Set("environment_name", environment.EnvironmentName)
	}

	if environment.Description != nil {
		_ = d.Set("description", environment.Description)
	}

	if environment.VpcId != nil {
		_ = d.Set("vpc", environment.VpcId)
	}

	if environment.SubnetIds != nil {
		_ = d.Set("subnet_ids", environment.SubnetIds)
	}

	if environment.Tags != nil {
		tagList := []interface{}{}
		for _, tag := range environment.Tags {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = tag.TagValue
			}
			tagList = append(tagList, tagMap)
		}
		_ = d.Set("tag", tagList)
	}

	return nil
}

func resourceTencentCloudTemEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_environment.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyEnvironmentRequest()

	request.EnvironmentId = helper.String(d.Id())

	if d.HasChange("environment_name") {
		if v, ok := d.GetOk("environment_name"); ok {
			request.EnvironmentName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("vpc") {
		return fmt.Errorf("`vpc` do not support change now.")
	}

	if d.HasChange("subnet_ids") {
		if v, ok := d.GetOk("subnet_ids"); ok {
			subnetIdsSet := v.(*schema.Set).List()
			for i := range subnetIdsSet {
				subnetIds := subnetIdsSet[i].(string)
				request.SubnetIds = append(request.SubnetIds, &subnetIds)
			}
		}
	}

	if d.HasChange("tag") {
		return fmt.Errorf("`tag` do not support change now.")
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyEnvironment(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTemEnvironmentRead(d, meta)
}

func resourceTencentCloudTemEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_environment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	environmentId := d.Id()

	if err := service.DeleteTemEnvironmentById(ctx, environmentId); err != nil {
		return err
	}

	return nil
}
