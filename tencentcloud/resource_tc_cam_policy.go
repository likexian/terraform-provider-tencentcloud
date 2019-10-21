/*
Provides a resource to create a CAM policy.

Example Usage

```hcl
resource "tencentcloud_cam_policy" "foo" {
  name        = "cam-policy-test"
  document    = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"resource\":[\"*\"]}]}"
  description = "test"
}
```

Import

CAM policy can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_policy.foo 26655801
```
*/
package tencentcloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
)

func resourceTencentCloudCamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamPolicyCreate,
		Read:   resourceTencentCloudCamPolicyRead,
		Update: resourceTencentCloudCamPolicyUpdate,
		Delete: resourceTencentCloudCamPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of CAM policy.",
			},
			"document": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Document of the CAM policy. The syntax refers to https://intl.cloud.tencent.com/document/product/598/10604. There are some notes when using this para in terraform: 1. The elements in JSON claimed supporting two types as `string` and `array` only support type `array`; 2. Terraform does not support the `root` syntax, when it appears, it must be replaced with the uin it stands for.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the CAM policy.",
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type of the policy strategy. 1 means customer strategy and 2 means preset strategy.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the CAM policy.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update time of the CAM policy.",
			},
		},
	}
}

func resourceTencentCloudCamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy.create")()

	logId := getLogId(contextNil)

	name := d.Get("name").(string)
	document := d.Get("document").(string)

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	documentErr := camService.PolicyDocumentForceCheck(document)
	if documentErr != nil {
		return documentErr
	}
	request := cam.NewCreatePolicyRequest()
	request.PolicyName = &name
	request.PolicyDocument = &document
	if v, ok := d.GetOk("description"); ok {
		request.Description = stringToPointer(v.(string))
	}

	var response *cam.CreatePolicyResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().CreatePolicy(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	if response.Response.PolicyId == nil {
		return fmt.Errorf("CAM policy id is nil")
	}
	d.SetId(strconv.Itoa(int(*response.Response.PolicyId)))

	return resourceTencentCloudCamPolicyRead(d, meta)
}

func resourceTencentCloudCamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy.read")()

	logId := getLogId(contextNil)

	policyId := d.Id()
	request := cam.NewGetPolicyRequest()
	policyIdInt, e := strconv.Atoi(policyId)
	if e != nil {
		return e
	}
	policyIdInt64 := uint64(policyIdInt)
	request.PolicyId = &policyIdInt64
	var instance *cam.GetPolicyResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().GetPolicy(request)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}

	d.Set("name", *instance.Response.PolicyName)
	//document with special change rule, the `\\/` must be replaced with `/`
	d.Set("document", strings.Replace(*instance.Response.PolicyDocument, "\\/", "/", -1))
	d.Set("create_time", *instance.Response.AddTime)
	d.Set("update_time", *instance.Response.UpdateTime)
	d.Set("type", int(*instance.Response.Type))
	if instance.Response.Description != nil {
		d.Set("description", *instance.Response.Description)
	}
	return nil
}

func resourceTencentCloudCamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy.update")()

	logId := getLogId(contextNil)

	policyId := d.Id()
	policyIdInt, e := strconv.Atoi(policyId)
	if e != nil {
		return e
	}
	policyIdInt64 := uint64(policyIdInt)
	request := cam.NewUpdatePolicyRequest()
	request.PolicyId = &policyIdInt64
	changeFlag := false

	if d.HasChange("description") {
		request.Description = stringToPointer(d.Get("description").(string))
		changeFlag = true

	}
	if d.HasChange("name") {
		request.PolicyName = stringToPointer(d.Get("name").(string))
		changeFlag = true
	}

	if d.HasChange("document") {
		o, n := d.GetChange("document")
		flag, err := diffJson(o.(string), n.(string))
		if err != nil {
			log.Printf("[CRITAL]%s update CAM policy document failed, reason:%s\n", logId, err.Error())
			return err
		}
		if flag {
			document := d.Get("document").(string)
			camService := CamService{
				client: meta.(*TencentCloudClient).apiV3Conn,
			}
			documentErr := camService.PolicyDocumentForceCheck(document)
			if documentErr != nil {
				return documentErr
			}
			request.PolicyDocument = &document
			changeFlag = true
		}
	}
	if changeFlag {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().UpdatePolicy(request)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM policy description failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCamPolicyRead(d, meta)
}

func resourceTencentCloudCamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cam_policy.delete")()

	logId := getLogId(contextNil)

	policyId := d.Id()
	policyIdInt, e := strconv.Atoi(policyId)
	if e != nil {
		return e
	}
	policyIdInt64 := uint64(policyIdInt)
	request := cam.NewDeletePolicyRequest()
	request.PolicyId = []*uint64{&policyIdInt64}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, e := meta.(*TencentCloudClient).apiV3Conn.UseCamClient().DeletePolicy(request)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM policy failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}
