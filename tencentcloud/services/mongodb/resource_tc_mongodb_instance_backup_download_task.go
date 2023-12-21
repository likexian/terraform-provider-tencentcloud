package mongodb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMongodbInstanceBackupDownloadTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceBackupDownloadTaskCreate,
		Read:   resourceTencentCloudMongodbInstanceBackupDownloadTaskRead,
		Delete: resourceTencentCloudMongodbInstanceBackupDownloadTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.",
			},

			"backup_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the backup file to be downloaded can be obtained through the DescribeDBBackups interface.",
			},

			"backup_sets": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Specifies the node names of replica sets to download or a list of shard names for sharded clusters.For example, the replica set cmgo-p8vnipr5, example (fixed value): BackupSets.0=cmgo-p8vnipr5_0, the full amount of data can be downloaded.For example, the sharded cluster cmgo-p8vnipr5, for example: BackupSets.0=cmgo-p8vnipr5_0&amp;amp;BackupSets.1=cmgo-p8vnipr5_1, that is, to download the data of shard 0 and 1. If the sharded cluster needs to be downloaded in full, please pass in the example. Full slice name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replica_set_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Replication Id.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceBackupDownloadTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_download_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mongodb.NewCreateBackupDownloadTaskRequest()
		instanceId string
		backupName string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_name"); ok {
		backupName = v.(string)
		request.BackupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_sets"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			replicaSetInfo := mongodb.ReplicaSetInfo{}
			if v, ok := dMap["replica_set_id"]; ok {
				replicaSetInfo.ReplicaSetId = helper.String(v.(string))
			}
			request.BackupSets = append(request.BackupSets, &replicaSetInfo)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().CreateBackupDownloadTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mongodb instanceBackupDownloadTask failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + backupName)

	return resourceTencentCloudMongodbInstanceBackupDownloadTaskRead(d, meta)
}

func resourceTencentCloudMongodbInstanceBackupDownloadTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_download_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	backupName := idSplit[1]

	instanceBackupDownloadTask, err := service.DescribeMongodbInstanceBackupDownloadTaskById(ctx, instanceId, backupName)
	if err != nil {
		return err
	}

	if instanceBackupDownloadTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MongodbInstanceBackupDownloadTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("backup_name", backupName)

	if instanceBackupDownloadTask != nil {
		backupSetsList := []interface{}{}
		for _, backupSet := range instanceBackupDownloadTask {
			backupSetsMap := map[string]interface{}{}

			if backupSet.ReplicaSetId != nil {
				backupSetsMap["replica_set_id"] = backupSet.ReplicaSetId
			}
			backupSetsList = append(backupSetsList, backupSetsMap)
		}
		_ = d.Set("backup_sets", backupSetsList)
	}
	return nil
}

func resourceTencentCloudMongodbInstanceBackupDownloadTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_download_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
