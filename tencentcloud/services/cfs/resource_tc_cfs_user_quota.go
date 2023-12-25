package cfs

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfsUserQuota() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsUserQuotaCreate,
		Read:   resourceTencentCloudCfsUserQuotaRead,
		Update: resourceTencentCloudCfsUserQuotaUpdate,
		Delete: resourceTencentCloudCfsUserQuotaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "File system ID.",
			},

			"user_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Quota type. Valid value: `Uid`, `Gid`.",
			},

			"user_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Info of UID/GID.",
			},

			"capacity_hard_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Capacity Limit(GB).",
			},

			"file_hard_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "File limit.",
			},
		},
	}
}

func resourceTencentCloudCfsUserQuotaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_user_quota.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = cfs.NewSetUserQuotaRequest()
		fileSystemId string
		userType     string
		userId       string
	)
	if v, ok := d.GetOk("file_system_id"); ok {
		fileSystemId = v.(string)
		request.FileSystemId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_type"); ok {
		userType = v.(string)
		request.UserType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		request.UserId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("capacity_hard_limit"); v != nil {
		request.CapacityHardLimit = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("file_hard_limit"); v != nil {
		request.FileHardLimit = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().SetUserQuota(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cfs userQuota failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(fileSystemId + tccommon.FILED_SP + userType + tccommon.FILED_SP + userId)

	return resourceTencentCloudCfsUserQuotaRead(d, meta)
}

func resourceTencentCloudCfsUserQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_user_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	fileSystemId := idSplit[0]
	userType := idSplit[1]
	userId := idSplit[2]

	userQuota, err := service.DescribeCfsUserQuotaById(ctx, fileSystemId, userType, userId)
	if err != nil {
		return err
	}

	if userQuota == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CfsUserQuota` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if userQuota.FileSystemId != nil {
		_ = d.Set("file_system_id", userQuota.FileSystemId)
	}

	if userQuota.UserType != nil {
		_ = d.Set("user_type", userQuota.UserType)
	}

	if userQuota.UserId != nil {
		_ = d.Set("user_id", userQuota.UserId)
	}

	if userQuota.CapacityHardLimit != nil {
		_ = d.Set("capacity_hard_limit", userQuota.CapacityHardLimit)
	}

	if userQuota.FileHardLimit != nil {
		_ = d.Set("file_hard_limit", userQuota.FileHardLimit)
	}

	return nil
}

func resourceTencentCloudCfsUserQuotaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_user_quota.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"file_system_id", "user_type", "user_id", "capacity_hard_limit", "file_hard_limit"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudCfsUserQuotaRead(d, meta)
}

func resourceTencentCloudCfsUserQuotaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_user_quota.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	fileSystemId := idSplit[0]
	userType := idSplit[1]
	userId := idSplit[2]

	if err := service.DeleteCfsUserQuotaById(ctx, fileSystemId, userType, userId); err != nil {
		return err
	}

	return nil
}
