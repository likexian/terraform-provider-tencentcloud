package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlDeployGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlDeployGroupCreate,
		Read:   resourceTencentCloudMysqlDeployGroupRead,
		Update: resourceTencentCloudMysqlDeployGroupUpdate,
		Delete: resourceTencentCloudMysqlDeployGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"deploy_group_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of deploy group. the maximum length cannot exceed 60 characters.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The description of deploy group. the maximum length cannot exceed 200 characters.",
			},

			"limit_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The limit on the number of instances on the same physical machine in deploy group affinity policy 1.",
			},

			"dev_class": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The device class of deploy group. optional value is SH12+SH02, TS85, etc.",
			},
		},
	}
}

func resourceTencentCloudMysqlDeployGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_deploy_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = mysql.NewCreateDeployGroupRequest()
		response      = mysql.NewCreateDeployGroupResponse()
		deployGroupId string
	)
	if v, ok := d.GetOk("deploy_group_name"); ok {
		request.DeployGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, _ := d.GetOk("limit_num"); v != nil {
		request.LimitNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("dev_class"); ok {
		devClassSet := v.(*schema.Set).List()
		for i := range devClassSet {
			devClass := devClassSet[i].(string)
			request.DevClass = append(request.DevClass, &devClass)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateDeployGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql deployGroup failed, reason:%+v", logId, err)
		return err
	}

	deployGroupId = *response.Response.DeployGroupId
	d.SetId(deployGroupId)

	return resourceTencentCloudMysqlDeployGroupRead(d, meta)
}

func resourceTencentCloudMysqlDeployGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_deploy_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	deployGroupId := d.Id()

	deployGroup, err := service.DescribeMysqlDeployGroupById(ctx, deployGroupId)
	if err != nil {
		return err
	}

	if deployGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_mysql_deploy_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil

	}
	if deployGroup.DeployGroupName != nil {
		_ = d.Set("deploy_group_name", deployGroup.DeployGroupName)
	}

	if deployGroup.Description != nil {
		_ = d.Set("description", deployGroup.Description)
	}

	if deployGroup.LimitNum != nil {
		_ = d.Set("limit_num", deployGroup.LimitNum)
	}

	if deployGroup.DevClass != nil {
		_ = d.Set("dev_class", []*string{deployGroup.DevClass})
	}

	return nil
}

func resourceTencentCloudMysqlDeployGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_deploy_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mysql.NewModifyNameOrDescByDpIdRequest()

	deployGroupId := d.Id()

	request.DeployGroupId = &deployGroupId

	immutableArgs := []string{"limit_num", "dev_class"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("deploy_group_name") {
		if v, ok := d.GetOk("deploy_group_name"); ok {
			request.DeployGroupName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyNameOrDescByDpId(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql deployGroup failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlDeployGroupRead(d, meta)
}

func resourceTencentCloudMysqlDeployGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_deploy_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	deployGroupId := d.Id()

	if err := service.DeleteMysqlDeployGroupById(ctx, deployGroupId); err != nil {
		return err
	}

	return nil
}
