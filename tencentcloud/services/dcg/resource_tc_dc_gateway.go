package dcg

import (
	"context"
	"fmt"
	"log"
	"strings"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudDcGatewayInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcGatewayCreate,
		Read:   resourceTencentCloudDcGatewayRead,
		Update: resourceTencentCloudDcGatewayUpdate,
		Delete: resourceTencentCloudDcGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the DCG.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DCG_NETWORK_TYPES),
				Description:  "Type of associated network. Valid value: `VPC` and `CCN`.",
			},
			"network_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "If the `network_type` value is `VPC`, the available value is VPC ID. But when the `network_type` value is `CCN`, the available value is CCN instance ID.",
			},
			"gateway_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DCG_GATEWAY_TYPE_NORMAL,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DCG_GATEWAY_TYPES),
				Description:  "Type of the gateway. Valid value: `NORMAL` and `NAT`. Default is `NORMAL`. NOTES: CCN only supports `NORMAL` and a VPC can create two DCGs, the one is NAT type and the other is non-NAT type.",
			},

			//compute
			"cnn_route_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of CCN route. Valid value: `BGP` and `STATIC`. The property is available when the DCG type is CCN gateway and BGP enabled.",
			},
			"enable_bgp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the BGP is enabled.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of resource.",
			},
		},
	}
}

func resourceTencentCloudDcGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway.create")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		request           = vpc.NewCreateDirectConnectGatewayRequest()
		response          = vpc.NewCreateDirectConnectGatewayResponse()
		networkType       string
		networkInstanceId string
		gatewayType       string
	)

	if v, ok := d.GetOk("name"); ok {
		request.DirectConnectGatewayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok {
		request.NetworkType = helper.String(v.(string))
		networkType = v.(string)
	}

	if v, ok := d.GetOk("network_instance_id"); ok {
		request.NetworkInstanceId = helper.String(v.(string))
		networkInstanceId = v.(string)
	}

	if v, ok := d.GetOk("gateway_type"); ok {
		request.GatewayType = helper.String(v.(string))
		gatewayType = v.(string)
	}

	if networkType == DCG_NETWORK_TYPE_VPC && !strings.HasPrefix(networkInstanceId, "vpc") {
		return fmt.Errorf("if `network_type` is '%s', the field `network_instance_id` must be a VPC resource", DCG_NETWORK_TYPE_VPC)
	}

	if networkType == DCG_NETWORK_TYPE_CCN && !strings.HasPrefix(networkInstanceId, "ccn") {
		return fmt.Errorf("if `network_type` is '%s', the field `network_instance_id` must be a CCN resource", DCG_NETWORK_TYPE_CCN)
	}

	if networkType == DCG_NETWORK_TYPE_CCN && gatewayType != DCG_GATEWAY_TYPE_NORMAL {
		return fmt.Errorf("if `network_type` is '%s', the field `gateway_type` must be '%s'", DCG_NETWORK_TYPE_CCN, DCG_GATEWAY_TYPE_NORMAL)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateDirectConnectGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.DirectConnectGateway == nil {
			return resource.NonRetryableError(fmt.Errorf("Create direct connect gateway failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create direct connect gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	if response.Response.DirectConnectGateway.DirectConnectGatewayId == nil {
		return fmt.Errorf("DirectConnectGatewayId is nil.")
	}

	d.SetId(*response.Response.DirectConnectGateway.DirectConnectGatewayId)

	return resourceTencentCloudDcGatewayRead(d, meta)
}

func resourceTencentCloudDcGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		info, has, e := service.DescribeDirectConnectGateway(ctx, d.Id())
		if e != nil {
			return tccommon.RetryError(e)
		}

		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("name", info.name)
		_ = d.Set("network_type", info.networkType)
		_ = d.Set("network_instance_id", info.networkInstanceId)
		_ = d.Set("gateway_type", info.gatewayType)
		_ = d.Set("cnn_route_type", info.cnnRouteType)
		_ = d.Set("enable_bgp", info.enableBGP)
		_ = d.Set("create_time", info.createTime)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudDcGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway.update")()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		dcgId = d.Id()
	)

	if d.HasChange("name") {
		request := vpc.NewModifyDirectConnectGatewayAttributeRequest()
		request.DirectConnectGatewayId = helper.String(dcgId)
		if v, ok := d.GetOk("name"); ok {
			request.DirectConnectGatewayName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyDirectConnectGatewayAttribute(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update direct connect gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudDcGatewayRead(d, meta)
}

func resourceTencentCloudDcGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dc_gateway.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vpc.NewDeleteDirectConnectGatewayRequest()
		dcgId   = d.Id()
	)

	request.DirectConnectGatewayId = helper.String(dcgId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteDirectConnectGateway(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete direct connect gateway failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
