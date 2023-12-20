package ckafka

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCkafkaConsumerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaConsumerGroupCreate,
		Read:   resourceTencentCloudCkafkaConsumerGroupRead,
		Delete: resourceTencentCloudCkafkaConsumerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "InstanceId.",
			},

			"group_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "GroupName.",
			},

			"topic_name_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "array of topic names.",
			},
		},
	}
}

func resourceTencentCloudCkafkaConsumerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = ckafka.NewCreateConsumerRequest()
		instanceId string
		groupName  string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("group_name"); ok {
		groupName = v.(string)
		request.GroupName = helper.String(groupName)
	}

	if v, ok := d.GetOk("topic_name_list"); ok {
		topicNameListSet := v.(*schema.Set).List()
		for i := range topicNameListSet {
			topicNameList := topicNameListSet[i].(string)
			request.TopicNameList = append(request.TopicNameList, &topicNameList)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().CreateConsumer(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka consumerGroup failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + groupName)

	return resourceTencentCloudCkafkaConsumerGroupRead(d, meta)
}

func resourceTencentCloudCkafkaConsumerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupName := idSplit[1]

	consumerGroup, err := service.DescribeCkafkaConsumerGroupById(ctx, instanceId, groupName)
	if err != nil {
		return err
	}

	if consumerGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CkafkaConsumerGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("group_name", groupName)

	if consumerGroup.TopicList != nil {
		topicNameList := make([]string, 0)
		for _, v := range consumerGroup.TopicList {
			topicNameList = append(topicNameList, *v.TopicName)
		}
		_ = d.Set("topic_name_list", topicNameList)
	}

	return nil
}

func resourceTencentCloudCkafkaConsumerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	groupName := idSplit[1]

	if err := service.DeleteCkafkaConsumerGroupById(ctx, instanceId, groupName); err != nil {
		return err
	}

	return nil
}
