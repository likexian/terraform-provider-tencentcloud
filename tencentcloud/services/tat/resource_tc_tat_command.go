package tat

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTatCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatCommandCreate,
		Read:   resourceTencentCloudTatCommandRead,
		Update: resourceTencentCloudTatCommandUpdate,
		Delete: resourceTencentCloudTatCommandDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"command_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Command name. The name can be up to 60 bytes, and contain [a-z], [A-Z], [0-9] and [_-.].",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Base64-encoded command. The maximum length is 64 KB.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Command description. The maximum length is 120 characters.",
			},

			"command_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Command type. `SHELL` and `POWERSHELL` are supported. The default value is `SHELL`.",
			},

			"working_directory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Command execution path. The default value is /root for `SHELL` commands and C:/Program Files/qcloudtat_agent/workdir for `POWERSHELL` commands.",
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Command timeout period. Default value: 60 seconds. Value range: [1, 86400].",
			},

			"enable_parameter": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable the custom parameter feature.This cannot be modified once created.Default value: `false`.",
			},

			"default_parameters": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default value of the custom parameter value when it is enabled. The field type is JSON encoded string. For example, {\"varA\": \"222\"}.`key` is the name of the custom parameter and value is the default value. Both `key` and `value` are strings.If no parameter value is provided in the `InvokeCommand` API, the default value is used.Up to 20 custom parameters are supported.The name of the custom parameter cannot exceed 64 characters and can contain [a-z], [A-Z], [0-9] and [-_].",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags bound to the command. At most 10 tags are allowed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateNotEmpty,
							Description:  "Tag key.",
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateNotEmpty,
							Description:  "Tag value.",
						},
					},
				},
			},

			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username used to execute the command on the CVM or Lighthouse instance.The principle of least privilege is the best practice for permission management. We recommend you execute TAT commands as a general user. By default, the root user is used to execute commands on Linux and the System user is used on Windows.",
			},

			"output_cos_bucket_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The COS bucket URL for uploading logs. The URL must start with `https`, such as `https://BucketName-123454321.cos.ap-beijing.myqcloud.com`.",
			},

			"output_cos_key_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The COS bucket directory where the logs are saved. Check below for the rules of the directory name.1. It must be a combination of number, letters, and visible characters. Up to 60 characters are allowed.2. Use a slash (/) to create a subdirectory.3. Consecutive dots (.) and slashes (/) are not allowed. It can not start with a slash (/).",
			},

			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Command creation time.",
			},

			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Command update time.",
			},

			"formatted_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Formatted description of the command. This parameter is an empty string for user commands and contains values for public commands.",
			},

			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Command creator. `TAT` indicates a public command and `USER` indicates a personal command.",
			},
		},
	}
}

func resourceTencentCloudTatCommandCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_command.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = tat.NewCreateCommandRequest()
		response  *tat.CreateCommandResponse
		commandId string
	)

	if v, ok := d.GetOk("command_name"); ok {
		request.CommandName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(tccommon.StringToBase64(v.(string)))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command_type"); ok {
		request.CommandType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("working_directory"); ok {
		request.WorkingDirectory = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		request.Timeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_parameter"); ok {
		request.EnableParameter = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("default_parameters"); ok {
		request.DefaultParameters = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tag := tat.Tag{}
			if v, ok := dMap["key"]; ok {
				tag.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				tag.Value = helper.String(v.(string))
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	if v, ok := d.GetOk("username"); ok {
		request.Username = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output_cos_bucket_url"); ok {
		request.OutputCOSBucketUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output_cos_key_prefix"); ok {
		request.OutputCOSKeyPrefix = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTatClient().CreateCommand(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tat command failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tat command failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.CommandId == nil {
		return fmt.Errorf("CommandId is nil.")
	}

	commandId = *response.Response.CommandId

	d.SetId(commandId)
	return resourceTencentCloudTatCommandRead(d, meta)
}

func resourceTencentCloudTatCommandRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_command.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	commandId := d.Id()

	command, err := service.DescribeTatCommand(ctx, commandId)

	if err != nil {
		return err
	}

	if command == nil {
		d.SetId("")
		return fmt.Errorf("resource `command` %s does not exist", commandId)
	}

	if command.CommandName != nil {
		_ = d.Set("command_name", command.CommandName)
	}

	if command.Content != nil {
		content, err := tccommon.Base64ToString(*command.Content)
		if err != nil {
			return fmt.Errorf("`Content` [%v] base64 to string failed, err: %v.", *command.Content, err)
		}
		_ = d.Set("content", content)
	}

	if command.Description != nil {
		_ = d.Set("description", command.Description)
	}

	if command.CommandType != nil {
		_ = d.Set("command_type", command.CommandType)
	}

	if command.WorkingDirectory != nil {
		_ = d.Set("working_directory", command.WorkingDirectory)
	}

	if command.Timeout != nil {
		_ = d.Set("timeout", command.Timeout)
	}

	if command.EnableParameter != nil {
		_ = d.Set("enable_parameter", command.EnableParameter)
	}

	if command.DefaultParameters != nil {
		_ = d.Set("default_parameters", command.DefaultParameters)
	}

	if command.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range command.Tags {
			tagsMap := map[string]interface{}{}
			if tags.Key != nil {
				tagsMap["key"] = tags.Key
			}
			if tags.Value != nil {
				tagsMap["value"] = tags.Value
			}

			tagsList = append(tagsList, tagsMap)
		}
		_ = d.Set("tags", tagsList)
	}

	if command.Username != nil {
		_ = d.Set("username", command.Username)
	}

	if command.OutputCOSBucketUrl != nil {
		_ = d.Set("output_cos_bucket_url", command.OutputCOSBucketUrl)
	}

	if command.OutputCOSKeyPrefix != nil {
		_ = d.Set("output_cos_key_prefix", command.OutputCOSKeyPrefix)
	}

	if command.CreatedTime != nil {
		_ = d.Set("created_time", command.CreatedTime)
	}

	if command.UpdatedTime != nil {
		_ = d.Set("updated_time", command.UpdatedTime)
	}

	if command.FormattedDescription != nil {
		_ = d.Set("formatted_description", command.FormattedDescription)
	}

	if command.CreatedBy != nil {
		_ = d.Set("created_by", command.CreatedBy)
	}

	return nil
}

func resourceTencentCloudTatCommandUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_command.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := tat.NewModifyCommandRequest()

	if d.HasChange("enable_parameter") {
		return fmt.Errorf("`enable_parameter` do not support change now.")
	}

	if d.HasChange("tags") {
		return fmt.Errorf("`tags` do not support change now.")
	}

	commandId := d.Id()
	request.CommandId = &commandId

	if d.HasChange("command_name") {
		if v, ok := d.GetOk("command_name"); ok {
			request.CommandName = helper.String(v.(string))
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(tccommon.StringToBase64(v.(string)))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("command_type") {
		if v, ok := d.GetOk("command_type"); ok {
			request.CommandType = helper.String(v.(string))
		}
	}

	if d.HasChange("working_directory") {
		if v, ok := d.GetOk("working_directory"); ok {
			request.WorkingDirectory = helper.String(v.(string))
		}
	}

	if d.HasChange("timeout") {
		if v, ok := d.GetOk("timeout"); ok {
			request.Timeout = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("default_parameters") {
		if v, ok := d.GetOk("default_parameters"); ok {
			request.DefaultParameters = helper.String(v.(string))
		}
	}

	if d.HasChange("username") {
		if v, ok := d.GetOk("username"); ok {
			request.Username = helper.String(v.(string))
		}
	}

	if d.HasChange("output_cos_bucket_url") {
		if v, ok := d.GetOk("output_cos_bucket_url"); ok {
			request.OutputCOSBucketUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("output_cos_key_prefix") {
		if v, ok := d.GetOk("output_cos_key_prefix"); ok {
			request.OutputCOSKeyPrefix = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTatClient().ModifyCommand(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tat command failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTatCommandRead(d, meta)
}

func resourceTencentCloudTatCommandDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_command.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	commandId := d.Id()

	if err := service.DeleteTatCommandById(ctx, commandId); err != nil {
		return err
	}

	return nil
}
