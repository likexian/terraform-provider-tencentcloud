package pls

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcEndPointServiceWhiteList() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcEndPointServiceWhiteListCreate,
		Read:   resourceTencentCloudVpcEndPointServiceWhiteListRead,
		Update: resourceTencentCloudVpcEndPointServiceWhiteListUpdate,
		Delete: resourceTencentCloudVpcEndPointServiceWhiteListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_uin": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "UIN.",
			},

			"end_point_service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of endpoint service.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description of white list.",
			},

			"owner": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "APPID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create Time.",
			},
		},
	}
}

func resourceTencentCloudVpcEndPointServiceWhiteListCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point_service_white_list.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request           = vpc.NewCreateVpcEndPointServiceWhiteListRequest()
		userUin           string
		endPointServiceId string
	)
	if v, ok := d.GetOk("user_uin"); ok {
		userUin = v.(string)
		request.UserUin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_point_service_id"); ok {
		endPointServiceId = v.(string)
		request.EndPointServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateVpcEndPointServiceWhiteList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc endPointServiceWhiteList failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userUin + tccommon.FILED_SP + endPointServiceId)

	return resourceTencentCloudVpcEndPointServiceWhiteListRead(d, meta)
}

func resourceTencentCloudVpcEndPointServiceWhiteListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point_service_white_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	userUin := idSplit[0]
	endPointServiceId := idSplit[1]

	endPointServiceWhiteList, err := service.DescribeVpcEndPointServiceWhiteListById(ctx, userUin, endPointServiceId)
	if err != nil {
		return err
	}

	if endPointServiceWhiteList == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if endPointServiceWhiteList.UserUin != nil {
		_ = d.Set("user_uin", endPointServiceWhiteList.UserUin)
	}

	if endPointServiceWhiteList.EndPointServiceId != nil {
		_ = d.Set("end_point_service_id", endPointServiceWhiteList.EndPointServiceId)
	}

	if endPointServiceWhiteList.Description != nil {
		_ = d.Set("description", endPointServiceWhiteList.Description)
	}

	if endPointServiceWhiteList.Owner != nil {
		_ = d.Set("owner", helper.UInt64ToStr(*endPointServiceWhiteList.Owner))
	}

	if endPointServiceWhiteList.CreateTime != nil {
		_ = d.Set("create_time", endPointServiceWhiteList.CreateTime)
	}

	return nil
}

func resourceTencentCloudVpcEndPointServiceWhiteListUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point_service_white_list.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewModifyVpcEndPointServiceWhiteListRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	userUin := idSplit[0]
	endPointServiceId := idSplit[1]

	request.UserUin = &userUin
	request.EndPointServiceId = &endPointServiceId

	unsupportedUpdateFields := []string{
		"user_uin",
		"end_point_service_id",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_vpc_end_point_service_white_list update on %s is not support yet", field)
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpcEndPointServiceWhiteList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vpc endPointServiceWhiteList failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpcEndPointServiceWhiteListRead(d, meta)
}

func resourceTencentCloudVpcEndPointServiceWhiteListDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_end_point_service_white_list.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	userUin := idSplit[0]
	endPointServiceId := idSplit[1]

	if err := service.DeleteVpcEndPointServiceWhiteListById(ctx, userUin, endPointServiceId); err != nil {
		return nil
	}

	return nil
}
