package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func resourceTencentCloudAPIGatewayAPIKeyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayAPIKeyAttachmentCreate,
		Read:   resourceTencentCloudAPIGatewayAPIKeyAttachmentRead,
		Delete: resourceTencentCloudAPIGatewayAPIKeyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"api_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of API key.",
			},
			"usage_plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the usage plan.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayAPIKeyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_key_attachment.create")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		apiKeyId          = d.Get("api_key_id").(string)
		usagePlanId       = d.Get("usage_plan_id").(string)
		has               bool
		err               error
	)

	//check usage plan is exist
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, err = apiGatewayService.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("usage plan %s is not exist", usagePlanId)
	}

	//check API key is exist
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, has, err = apiGatewayService.DescribeApiKey(ctx, apiKeyId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("API key %s is not exist", apiKeyId)
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if err = apiGatewayService.BindSecretId(ctx, usagePlanId, apiKeyId); err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	//waiting bind success
	var info apigateway.UsagePlanInfo
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, err = apiGatewayService.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			return retryError(err, InternalError)
		}
		if !has {
			return nil
		}
		for _, v := range info.BindSecretIds {
			if *v == apiKeyId {
				return nil
			}
		}
		return resource.RetryableError(
			fmt.Errorf("API key  %s attach to usage plan %s still is doing",
				apiKeyId, usagePlanId))

	}); err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("usage plan %s has been deleted", usagePlanId)
	}
	d.SetId(strings.Join([]string{apiKeyId, usagePlanId}, FILED_SP))

	return resourceTencentCloudAPIGatewayAPIKeyAttachmentRead(d, meta)
}

func resourceTencentCloudAPIGatewayAPIKeyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_key_attachment.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		info              apigateway.UsagePlanInfo
		err               error
		has               bool
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	apiKeyId := idSplit[0]
	usagePlanId := idSplit[1]
	if apiKeyId == "" || usagePlanId == "" {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, err = apiGatewayService.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	for _, v := range info.BindSecretIds {
		if *v == apiKeyId {
			_ = d.Set("api_key_id", apiKeyId)
			_ = d.Set("usage_plan_id", usagePlanId)
			break
		}
	}
	return nil
}

func resourceTencentCloudAPIGatewayAPIKeyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_api_gateway_api_key_attachment.delete")()

	var (
		logId             = getLogId(contextNil)
		ctx               = context.WithValue(context.TODO(), logIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(*TencentCloudClient).apiV3Conn}
		info              apigateway.UsagePlanInfo
		err               error
		has               bool
	)
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	apiKeyId := idSplit[0]
	usagePlanId := idSplit[1]
	if apiKeyId == "" || usagePlanId == "" {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err = apiGatewayService.UnBindSecretId(ctx, usagePlanId, apiKeyId)
		if err != nil {
			return retryError(err)
		}
		return nil
	}); err != nil {
		return err
	}

	//waiting delete ok
	if err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		info, has, err = apiGatewayService.DescribeUsagePlan(ctx, usagePlanId)
		if err != nil {
			return retryError(err, InternalError)
		}
		if !has {
			return nil
		}
		for _, v := range info.BindSecretIds {
			if *v == apiKeyId {
				return resource.RetryableError(
					fmt.Errorf("API key  %s attach to usage plan %s still is deleting.",
						apiKeyId, usagePlanId))
			}
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
