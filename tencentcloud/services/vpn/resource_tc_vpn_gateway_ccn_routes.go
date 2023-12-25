package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpnGatewayCcnRoutes() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpnGatewayCcnRoutesCreate,
		Read:   resourceTencentCloudVpnGatewayCcnRoutesRead,
		Update: resourceTencentCloudVpnGatewayCcnRoutesUpdate,
		Delete: resourceTencentCloudVpnGatewayCcnRoutesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VPN GATEWAY INSTANCE ID.",
			},
			"route_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Route Id.",
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Whether routing information is enabled. `ENABLE`: Enable Route, `DISABLE`: Disable Route.",
			},
			"destination_cidr_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Routing CIDR.",
			},
		},
	}
}

func resourceTencentCloudVpnGatewayCcnRoutesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_ccn_routes.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		vpnGwId string
		routeId string
	)

	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		vpnGwId = v.(string)
	}

	if v, ok := d.GetOk("route_id"); ok {
		routeId = v.(string)
	}

	d.SetId(vpnGwId + tccommon.FILED_SP + routeId)

	return resourceTencentCloudVpnGatewayCcnRoutesUpdate(d, meta)
}

func resourceTencentCloudVpnGatewayCcnRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_ccn_routes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpnGatewayId := idSplit[0]
	routeId := idSplit[1]

	vpnGatewayCcnRoutes, err := service.DescribeVpcVpnGatewayCcnRoutesById(ctx, vpnGatewayId, routeId)
	if err != nil {
		return err
	}

	if vpnGatewayCcnRoutes == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcVpnGatewayCcnRoutes` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("vpn_gateway_id", vpnGatewayId)
	_ = d.Set("route_id", vpnGatewayCcnRoutes.RouteId)
	_ = d.Set("status", vpnGatewayCcnRoutes.Status)
	_ = d.Set("destination_cidr_block", vpnGatewayCcnRoutes.DestinationCidrBlock)

	return nil
}

func resourceTencentCloudVpnGatewayCcnRoutesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_ccn_routes.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewModifyVpnGatewayCcnRoutesRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	vpnGatewayId := idSplit[0]
	routeId := idSplit[1]

	request.VpnGatewayId = &vpnGatewayId
	route := vpc.VpngwCcnRoutes{}
	route.RouteId = &routeId
	route.Status = helper.String(d.Get("status").(string))
	route.DestinationCidrBlock = helper.String(d.Get("destination_cidr_block").(string))
	request.Routes = append(request.Routes, &route)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyVpnGatewayCcnRoutes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vpnGatewayCcnRoutes failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVpnGatewayCcnRoutesRead(d, meta)
}

func resourceTencentCloudVpnGatewayCcnRoutesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpn_gateway_ccn_routes.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
