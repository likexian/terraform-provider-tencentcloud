package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsTopicCreate,
		Read:   resourceTencentCloudClsTopicRead,
		Delete: resourceTencentCloudClsTopicDelete,
		Update: resourceTencentCloudClsTopicUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"logset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Logset ID.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Log topic name.",
			},
			"partition_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of log topic partitions. Default value: 1. Maximum value: 10.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list. Up to 10 tag key-value pairs are supported and must be unique.",
			},
			"auto_split": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable automatic split. Default value: true.",
			},
			"max_split_partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Maximum number of partitions to split into for this topic if" +
					" automatic split is enabled. Default value: 50.",
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Log topic storage class. Valid values: hot: real-time storage; cold: offline storage. Default value: hot. If cold is passed in, " +
					"please contact the customer service to add the log topic to the allowlist first.",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Lifecycle in days. Value range: 1~366. Default value: 30.",
			},
			"hot_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "0: Turn off log sinking. Non 0: The number of days of standard storage after enabling log settling. HotPeriod needs to be greater than or equal to 7 and less than Period. Only effective when StorageType is hot.",
			},
			"describes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log Topic Description.",
			},
		},
	}
}

func resourceTencentCloudClsTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_topic.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = cls.NewCreateTopicRequest()
		response *cls.CreateTopicResponse
	)

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("partition_count"); ok {
		request.PartitionCount = helper.IntInt64(v.(int))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &key,
				Value: &value,
			})
		}
	}

	if v, ok := d.GetOkExists("auto_split"); ok {
		request.AutoSplit = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("max_split_partitions"); ok {
		request.MaxSplitPartitions = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request.StorageType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("hot_period"); ok {
		request.HotPeriod = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("describes"); ok {
		request.Describes = helper.String(v.(string))
	} else {
		request.Describes = helper.String("")
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create cls topic failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls topic failed, reason:%+v", logId, err)
		return err
	}

	id := *response.Response.TopicId
	d.SetId(id)
	return resourceTencentCloudClsTopicRead(d, meta)
}

func resourceTencentCloudClsTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	topic, err := service.DescribeClsTopicById(ctx, id)
	if err != nil {
		return err
	}

	if topic == nil {
		d.SetId("")
		return fmt.Errorf("resource `Topic` %s does not exist", id)
	}

	_ = d.Set("logset_id", topic.LogsetId)
	_ = d.Set("topic_name", topic.TopicName)
	_ = d.Set("partition_count", topic.PartitionCount)

	tags := make(map[string]string, len(topic.Tags))
	for _, tag := range topic.Tags {
		tags[*tag.Key] = *tag.Value
	}

	_ = d.Set("tags", tags)
	_ = d.Set("auto_split", topic.AutoSplit)
	_ = d.Set("max_split_partitions", topic.MaxSplitPartitions)
	_ = d.Set("storage_type", topic.StorageType)
	_ = d.Set("period", topic.Period)
	_ = d.Set("hot_period", topic.HotPeriod)
	_ = d.Set("describes", topic.Describes)

	return nil
}

func resourceTencentCloudClsTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_topic.update")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = cls.NewModifyTopicRequest()
		id      = d.Id()
	)

	immutableArgs := []string{"partition_count", "storage_type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.TopicId = helper.String(id)

	if d.HasChange("topic_name") {
		request.TopicName = helper.String(d.Get("topic_name").(string))
	}

	if d.HasChange("tags") {
		tags := d.Get("tags").(map[string]interface{})
		request.Tags = make([]*cls.Tag, 0, len(tags))
		for k, v := range tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &cls.Tag{
				Key:   &key,
				Value: helper.String(value.(string)),
			})
		}
	}

	if d.HasChange("auto_split") {
		request.AutoSplit = helper.Bool(d.Get("auto_split").(bool))
	}

	if d.HasChange("max_split_partitions") {
		request.MaxSplitPartitions = helper.IntInt64(d.Get("max_split_partitions").(int))
	}

	if d.HasChange("period") {
		request.Period = helper.IntInt64(d.Get("period").(int))
	}

	if d.HasChange("hot_period") {
		request.HotPeriod = helper.IntUint64(d.Get("hot_period").(int))
	}

	if d.HasChange("describes") {
		request.Describes = helper.String(d.Get("describes").(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyTopic(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsTopicRead(d, meta)
}

func resourceTencentCloudClsTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_topic.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	if err := service.DeleteClsTopic(ctx, id); err != nil {
		return err
	}

	return nil
}
