package ssm

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSsmSshKeyPairSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmSshKeyPairSecretCreate,
		Read:   resourceTencentCloudSsmSshKeyPairSecretRead,
		Update: resourceTencentCloudSsmSshKeyPairSecretUpdate,
		Delete: resourceTencentCloudSsmSshKeyPairSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Secret name, which must be unique in the same region. It can contain 128 bytes of letters, digits, hyphens and underscores and must begin with a letter or digit.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description, such as what it is used for. It contains up to 2,048 bytes.",
			},
			"kms_key_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Specifies a KMS CMK to encrypt the secret.If this parameter is left empty, the CMK created by Secrets Manager by default will be used for encryption.You can also specify a custom KMS CMK created in the same region for encryption.",
			},
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "ID of the project to which the created SSH key belongs.",
			},
			"ssh_key_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Name of the SSH key pair, which only contains digits, letters and underscores and must start with a digit or letter. The maximum length is 25 characters.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of secret.",
			},
			"status": {
				Optional:     true,
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"Enabled", "Disabled"}),
				Description:  "Enable or Disable Secret. Valid values is `Enabled` or `Disabled`. Default is `Enabled`.",
			},
			"clean_ssh_key": {
				Optional: true,
				Type:     schema.TypeBool,
				Description: "Specifies whether to delete the SSH key from both the secret and the SSH key list in the CVM console. This field is only take effect when delete SSH key secrets. Valid values: " +
					"`True`: deletes SSH key from both the secret and SSH key list in the CVM console. Note that the deletion will fail if the SSH key is already bound to a CVM instance." +
					"`False`: only deletes the SSH key information in the secret.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Credential creation time in UNIX timestamp format.",
			},
			"secret_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "`0`: user-defined secret. `1`: Tencent Cloud services secret. `2`: SSH key secret. `3`: Tencent Cloud API key secret. Note: this field may return `null`, indicating that no valid values can be obtained.",
			},
		},
	}
}

func resourceTencentCloudSsmSshKeyPairSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_ssh_key_pair_secret.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = ssm.NewCreateSSHKeyPairSecretRequest()
		response   = ssm.NewCreateSSHKeyPairSecretResponse()
		secretInfo *SecretInfo
		secretName string
	)

	if v, ok := d.GetOk("secret_name"); ok {
		request.SecretName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("ssh_key_name"); ok {
		request.SSHKeyName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().CreateSSHKeyPairSecret(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ssm sshKeyPairSecret failed, reason:%+v", logId, err)
		return err
	}

	secretName = *response.Response.SecretName
	d.SetId(secretName)

	// update status if disabled
	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		if status == "Disabled" {
			err = ssmService.DisableSecret(ctx, secretName)
			if err != nil {
				return err
			}
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			secretInfo, err = ssmService.DescribeSecretByName(ctx, secretName)
			if err != nil {
				return tccommon.RetryError(err)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err = tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudSsmSshKeyPairSecretRead(d, meta)
}

func resourceTencentCloudSsmSshKeyPairSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_ssh_key_pair_secret.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretInfo *SecretInfo
		secretName = d.Id()
	)

	sshKeyPairSecret, err := service.DescribeSecretById(ctx, secretName, 2)
	if err != nil {
		return err
	}

	if sshKeyPairSecret == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SsmSshKeyPairSecret` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if sshKeyPairSecret.SecretName != nil {
		_ = d.Set("secret_name", sshKeyPairSecret.SecretName)
	}

	if sshKeyPairSecret.ProjectID != nil {
		_ = d.Set("project_id", sshKeyPairSecret.ProjectID)
	}

	if sshKeyPairSecret.Description != nil {
		_ = d.Set("description", sshKeyPairSecret.Description)
	}

	if sshKeyPairSecret.KmsKeyId != nil {
		_ = d.Set("kms_key_id", sshKeyPairSecret.KmsKeyId)
	}

	if sshKeyPairSecret.ResourceName != nil {
		_ = d.Set("ssh_key_name", sshKeyPairSecret.ResourceName)
	}

	if sshKeyPairSecret.Status != nil {
		_ = d.Set("status", sshKeyPairSecret.Status)
	}

	if sshKeyPairSecret.CreateTime != nil {
		_ = d.Set("create_time", sshKeyPairSecret.CreateTime)
	}

	if sshKeyPairSecret.SecretType != nil {
		_ = d.Set("secret_type", sshKeyPairSecret.SecretType)
	}

	outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		secretInfo, err = service.DescribeSecretByName(ctx, secretName)
		if err != nil {
			return tccommon.RetryError(err)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "ssm", "secret", tcClient.Region, secretInfo.resourceId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudSsmSshKeyPairSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_ssh_key_pair_secret.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretName = d.Id()
	)

	immutableArgs := []string{
		"project_id",
		"kms_key_id",
		"ssh_key_name",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		request := ssm.NewUpdateDescriptionRequest()
		request.SecretName = &secretName

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().UpdateDescription(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update ssm sshKeyPairSecret failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("status") {
		service := SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			if status == "Disabled" {
				err := service.DisableSecret(ctx, secretName)
				if err != nil {
					return err
				}
			} else {
				err := service.EnableSecret(ctx, secretName)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		secretInfo, err := ssmService.DescribeSecretByName(ctx, secretName)
		if err != nil {
			return err
		}

		resourceName := tccommon.BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err = tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	return resourceTencentCloudSsmSshKeyPairSecretRead(d, meta)
}

func resourceTencentCloudSsmSshKeyPairSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_ssh_key_pair_secret.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretName = d.Id()
	)

	// disable before destroy
	err := service.DisableSecret(ctx, secretName)
	if err != nil {
		return err
	}

	var cleanSSHKey *bool

	if v, ok := d.GetOkExists("clean_ssh_key"); ok {
		cleanSSHKey = helper.Bool(v.(bool))
	}

	if err = service.DeleteSsmSshKeyPairSecretById(ctx, secretName, cleanSSHKey); err != nil {
		return err
	}

	return nil
}
