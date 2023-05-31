/*
Provides a resource to create a tdmq subscription_attachment

Example Usage

```hcl
resource "tencentcloud_tdmq_subscription_attachment" "subscription_attachment" {
  environment_id    			= "keep-ns"
  topic_name        			= "keep-topic"
  subscription_name 			= "test-subcription"
  remark 						= "test"
  cluster_id 					= "pulsar-9n95ax58b9vn"
  auto_create_policy_topic 		= true
}
```

Import

tdmq subscription_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_subscription_attachment.subscription_attachment subscription_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTdmqSubscriptionAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqSubscriptionAttachmentCreate,
		Read:   resourceTencentCloudTdmqSubscriptionAttachmentRead,
		Delete: resourceTencentCloudTdmqSubscriptionAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Environment (namespace) name.",
			},
			"topic_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "topic name.",
			},
			"subscription_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subscriber name, no more than 128 characters.",
			},
			"remark": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Remarks, within 128 characters.",
			},
			"cluster_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the Pulsar cluster.",
			},
			"auto_create_policy_topic": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to automatically create dead letters and retry topics, True means to create, False means not to create, the default is to automatically create dead letters and retry topics.",
			},
		},
	}
}

func resourceTencentCloudTdmqSubscriptionAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscription_attachment.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                 = getLogId(contextNil)
		request               = tdmq.NewCreateSubscriptionRequest()
		environmentId         string
		Topic                 string
		subscriptionName      string
		clusterId             string
		autoCreatePolicyTopic bool
	)

	if v, ok := d.GetOk("environment_id"); ok {
		request.EnvironmentId = helper.String(v.(string))
		environmentId = v.(string)
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
		Topic = v.(string)
	}

	if v, ok := d.GetOk("subscription_name"); ok {
		request.SubscriptionName = helper.String(v.(string))
		subscriptionName = v.(string)
	}

	if v, ok := d.GetOk("is_idempotent"); ok {
		request.IsIdempotent = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("auto_create_policy_topic"); ok {
		request.AutoCreatePolicyTopic = helper.Bool(v.(bool))
		autoCreatePolicyTopic = v.(bool)
	}

	var isIdempotent = false
	request.IsIdempotent = &isIdempotent
	request.AutoCreatePolicyTopic = &autoCreatePolicyTopic

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTdmqClient().CreateSubscription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq subscriptionAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{environmentId, Topic, subscriptionName, clusterId, strconv.FormatBool(autoCreatePolicyTopic)}, FILED_SP))

	return resourceTencentCloudTdmqSubscriptionAttachmentRead(d, meta)
}

func resourceTencentCloudTdmqSubscriptionAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscription_attachment.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}
	environmentId := idSplit[0]
	Topic := idSplit[1]
	subscriptionName := idSplit[2]
	clusterId := idSplit[3]
	autoCreatePolicyTopicStr := idSplit[4]

	subscriptionAttachment, err := service.DescribeTdmqSubscriptionAttachmentById(ctx, environmentId, Topic, subscriptionName, clusterId)
	if err != nil {
		return err
	}

	if subscriptionAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqSubscriptionAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if subscriptionAttachment.EnvironmentId != nil {
		_ = d.Set("environment_id", subscriptionAttachment.EnvironmentId)
	}

	if subscriptionAttachment.TopicName != nil {
		_ = d.Set("topic_name", subscriptionAttachment.TopicName)
	}

	if subscriptionAttachment.SubscriptionName != nil {
		_ = d.Set("subscription_name", subscriptionAttachment.SubscriptionName)
	}

	if subscriptionAttachment.Remark != nil {
		_ = d.Set("remark", subscriptionAttachment.Remark)
	}

	_ = d.Set("cluster_id", clusterId)

	autoCreatePolicyTopic, _ := strconv.ParseBool(autoCreatePolicyTopicStr)
	_ = d.Set("auto_create_policy_topic", autoCreatePolicyTopic)

	return nil
}

func resourceTencentCloudTdmqSubscriptionAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tdmq_subscription_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 5 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	environmentId := idSplit[0]
	Topic := idSplit[1]
	subscriptionName := idSplit[2]
	clusterId := idSplit[3]
	autoCreatePolicyTopicStr := idSplit[4]
	autoCreatePolicyTopic, _ := strconv.ParseBool(autoCreatePolicyTopicStr)

	// Delete Subscription
	if err := service.DeleteTdmqSubscriptionAttachmentById(ctx, environmentId, Topic, subscriptionName, clusterId); err != nil {
		return err
	}

	if autoCreatePolicyTopic {
		// Delete Topics
		if err := service.DeleteTdmqTopicsAttachmentById(ctx, environmentId, Topic, subscriptionName, clusterId); err != nil {
			return err
		}
	}

	return nil
}
