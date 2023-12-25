package sqlserver

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverGeneralBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralBackupCreate,
		Read:   resourceTencentCloudSqlserverGeneralBackupRead,
		Update: resourceTencentCloudSqlserverGeneralBackupUpdate,
		Delete: resourceTencentCloudSqlserverGeneralBackupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"strategy": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Backup policy (0: instance backup, 1: multi-database backup).",
			},
			"db_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of names of databases to be backed up (required only for multi-database backup).",
			},
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of mssql-i1z41iwd.",
			},
			"backup_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup name. If this parameter is left empty, a backup name in the format of [Instance ID]_[Backup start timestamp] will be automatically generated.",
			},
			"flow_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "flow id.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralBackupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_backup.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = sqlserver.NewCreateBackupRequest()
		instanceId string
		flowId     string
		backupId   uint64
		startStr   string
		endStr     string
		fileName   string
		err        error
	)

	if v, ok := d.GetOk("strategy"); ok {
		request.Strategy = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("db_names"); ok {
		dBNamesSet := v.(*schema.Set).List()
		for i := range dBNamesSet {
			dBNames := dBNamesSet[i].(string)
			request.DBNames = append(request.DBNames, &dBNames)
		}
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = *helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_name"); ok {
		request.BackupName = helper.String(v.(string))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().CreateBackup(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			err = fmt.Errorf("sqlserver Backup %s not exists", instanceId)
			return resource.NonRetryableError(err)
		}

		flowId = strconv.FormatInt(*result.Response.FlowId, 10)
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver Backup failed, reason:%+v", logId, err)
		return err
	}

	// waiting for backup done.
	err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBackupByFlowId(ctx, instanceId, flowId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			err = fmt.Errorf("sqlserver Backup %s not exists", instanceId)
			return resource.NonRetryableError(err)
		}

		if *result.Response.Status == SQLSERVER_BACKUP_RUNNING {
			return resource.RetryableError(fmt.Errorf("create sqlserver Backup task status is running"))
		}

		if *result.Response.Status == SQLSERVER_BACKUP_SUCCESS {
			backupId = *result.Response.Id
			startStr = *result.Response.StartTime
			endStr = *result.Response.EndTime
			fileName = *result.Response.FileName
			return nil
		}

		if *result.Response.Status == SQLSERVER_BACKUP_FAIL {
			return resource.NonRetryableError(fmt.Errorf("create sqlserver Backup task status is failed"))
		}

		err = fmt.Errorf("create sqlserver Backup task status is %v, we won't wait for it finish", *result.Response.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver Backup task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{instanceId, strconv.Itoa(int(backupId)), flowId, startStr, endStr, fileName}, tccommon.FILED_SP))
	return resourceTencentCloudSqlserverGeneralBackupRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_backup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
		flowId     string
		startStr   string
		endStr     string
		backupId   uint64
		strategy   int64
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 6 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId = idSplit[0]
	tempD, _ := strconv.Atoi(idSplit[1])
	backupId = uint64(tempD)
	flowId = idSplit[2]
	startStr = idSplit[3]
	endStr = idSplit[4]

	backupList, err := service.DescribeSqlserverBackupByBackupId(ctx, instanceId, startStr, endStr, backupId)
	if err != nil {
		return err
	}

	if backupList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralBackups` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	backupDetail := backupList[0]

	if backupDetail.Strategy != nil {
		strategy = *backupDetail.Strategy
		_ = d.Set("strategy", backupDetail.Strategy)
	}

	if backupDetail.BackupName != nil {
		_ = d.Set("backup_name", backupDetail.BackupName)
	}

	if strategy == SQLSERVER_BACKUP_STRATEGY_SINGEL {
		if backupDetail.DBs != nil {
			_ = d.Set("db_names", backupDetail.DBs)
		}
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("flow_id", flowId)

	return nil
}

func resourceTencentCloudSqlserverGeneralBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_backup.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"strategy", "db_names", "instance_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = sqlserver.NewModifyBackupNameRequest()
		instanceId string
		backupId   uint64
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 6 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId = idSplit[0]
	tempD, _ := strconv.Atoi(idSplit[1])
	backupId = uint64(tempD)

	request.InstanceId = &instanceId
	request.BackupId = &backupId

	if d.HasChange("backup_name") {
		if v, ok := d.GetOk("backup_name"); ok {
			request.BackupName = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyBackupName(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver generalBackups failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverGeneralBackupRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralBackupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_backup.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
		fileName   string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 6 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId = idSplit[0]
	fileName = idSplit[5]

	if err := service.DeleteSqlserverGeneralBackupsById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}
