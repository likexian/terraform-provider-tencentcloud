package sqlserver

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverConfigBackupStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigBackupStrategyCreate,
		Read:   resourceTencentCloudSqlserverConfigBackupStrategyRead,
		Update: resourceTencentCloudSqlserverConfigBackupStrategyUpdate,
		Delete: resourceTencentCloudSqlserverConfigBackupStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"backup_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backup type. Valid values: weekly (when length(BackupDay) <=7 && length(BackupDay) >=2), daily (when length(BackupDay)=1). Default value: daily.",
			},

			"backup_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup time. Value range: an integer from 0 to 23.",
			},

			"backup_day": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup interval in days when the BackupType is daily. The current value can only be 1.",
			},

			"backup_model": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backup mode. Valid values: master_pkg (archive the backup files of the primary node), master_no_pkg (do not archive the backup files of the primary node), slave_pkg (archive the backup files of the replica node), slave_no_pkg (do not archive the backup files of the replica node). Backup files of the replica node are supported only when Always On disaster recovery is enabled.",
			},

			"backup_cycle": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The days of the week on which backup will be performed when `BackupType` is weekly. If data backup retention period is less than 7 days, the values will be 1-7, indicating that backup will be performed everyday by default; if data backup retention period is greater than or equal to 7 days, the values will be at least any two days, indicating that backup will be performed at least twice in a week by default.",
			},

			"backup_save_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Data (log) backup retention period. Value range: 3-1830 days, default value: 7 days.",
			},

			"regular_backup_enable": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Archive backup status. Valid values: enable (enabled); disable (disabled). Default value: disable.",
			},

			"regular_backup_save_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Archive backup retention days. Value range: 90-3650 days. Default value: 365 days.",
			},

			"regular_backup_strategy": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Archive backup policy. Valid values: years (yearly); quarters (quarterly); months(monthly); Default value: `months`.",
			},

			"regular_backup_counts": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of retained archive backups. Default value: 1.",
			},

			"regular_backup_start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Archive backup start date in YYYY-MM-DD format, which is the current time by default.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigBackupStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_backup_strategy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigBackupStrategyUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigBackupStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_backup_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	configBackupStrategy, err := service.DescribeSqlserverConfigBackupStrategyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configBackupStrategy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigBackupStrategy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configBackupStrategy.InstanceId != nil {
		_ = d.Set("instance_id", configBackupStrategy.InstanceId)
	}

	if configBackupStrategy.BackupCycleType != nil {
		_ = d.Set("backup_type", configBackupStrategy.BackupCycleType)
		if configBackupStrategy.BackupCycleType == helper.String(SQLSERVER_BACKUP_CYCLETYPE_DAILY) {
			//Backup interval in days. When the BackupType is daily, valid value is 1.
			_ = d.Set("backup_day", 1)
		}
	}

	if configBackupStrategy.BackupTime != nil {
		_ = d.Set("backup_time", helper.StrToInt(*configBackupStrategy.BackupTime))
	}

	if configBackupStrategy.BackupModel != nil {
		_ = d.Set("backup_model", configBackupStrategy.BackupModel)
	}

	if configBackupStrategy.BackupCycle != nil {
		_ = d.Set("backup_cycle", configBackupStrategy.BackupCycle)
	}

	if configBackupStrategy.BackupSaveDays != nil {
		_ = d.Set("backup_save_days", configBackupStrategy.BackupSaveDays)
	}

	// if configBackupStrategy.RegularBackupEnable != nil {
	// 	_ = d.Set("regular_backup_enable", configBackupStrategy.RegularBackupEnable)
	// }

	// if configBackupStrategy.RegularBackupSaveDays != nil {
	// 	_ = d.Set("regular_backup_save_days", configBackupStrategy.RegularBackupSaveDays)
	// }

	// if configBackupStrategy.RegularBackupStrategy != nil {
	// 	_ = d.Set("regular_backup_strategy", configBackupStrategy.RegularBackupStrategy)
	// }

	// if configBackupStrategy.RegularBackupCounts != nil {
	// 	_ = d.Set("regular_backup_counts", configBackupStrategy.RegularBackupCounts)
	// }

	// if configBackupStrategy.RegularBackupStartTime != nil {
	// 	_ = d.Set("regular_backup_start_time", configBackupStrategy.RegularBackupStartTime)
	// }

	return nil
}

func resourceTencentCloudSqlserverConfigBackupStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_backup_strategy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := sqlserver.NewModifyBackupStrategyRequest()

	needChange := false

	request.InstanceId = helper.String(d.Id())

	mutableArgs := []string{"backup_type", "backup_time", "backup_day", "backup_model", "backup_cycle", "backup_save_days", "regular_backup_enable", "regular_backup_save_days", "regular_backup_strategy", "regular_backup_counts", "regular_backup_start_time"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("backup_type"); ok {
			request.BackupType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("backup_model"); ok {
			request.BackupModel = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("backup_time"); ok {
			request.BackupTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("backup_day"); ok {
			request.BackupDay = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("backup_cycle"); ok {
			backupCycleSet := v.(*schema.Set).List()
			for i := range backupCycleSet {
				backupCycle := backupCycleSet[i].(int)
				request.BackupCycle = append(request.BackupCycle, helper.IntUint64(backupCycle))
			}
		}

		if v, ok := d.GetOkExists("backup_save_days"); ok {
			request.BackupSaveDays = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("regular_backup_enable"); ok {
			request.RegularBackupEnable = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("regular_backup_save_days"); ok {
			request.RegularBackupSaveDays = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("regular_backup_strategy"); ok {
			request.RegularBackupStrategy = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("regular_backup_counts"); ok {
			request.RegularBackupCounts = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("regular_backup_start_time"); ok {
			request.RegularBackupStartTime = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyBackupStrategy(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update sqlserver configBackupStrategy failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudSqlserverConfigBackupStrategyRead(d, meta)
}

func resourceTencentCloudSqlserverConfigBackupStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_backup_strategy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
