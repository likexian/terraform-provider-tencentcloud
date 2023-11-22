/*
Provides a resource to create a antiddos ddos black white ip

Example Usage

```hcl
resource "tencentcloud_antiddos_ddos_black_white_ip" "ddos_black_white_ip" {
  instance_id = "bgp-xxxxxx"
  ip = "1.2.3.5"
  mask = 0
  type = "black"
}
```

Import

antiddos ddos_black_white_ip can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ddos_black_white_ip.ddos_black_white_ip ${instanceId}#${ip}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAntiddosDdosBlackWhiteIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosDdosBlackWhiteIpCreate,
		Read:   resourceTencentCloudAntiddosDdosBlackWhiteIpRead,
		Update: resourceTencentCloudAntiddosDdosBlackWhiteIpUpdate,
		Delete: resourceTencentCloudAntiddosDdosBlackWhiteIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ip list.",
			},
			"mask": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ip mask.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ip type, black: black ip list, white: white ip list.",
			},
		},
	}
}

func resourceTencentCloudAntiddosDdosBlackWhiteIpCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_ddos_black_white_ip.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = antiddos.NewCreateDDoSBlackWhiteIpListRequest()
		instanceId string
		ipType     string
		ip         string
		ipMask     int
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	ipSegment := antiddos.IpSegment{}
	if v, ok := d.GetOk("ip"); ok {
		ip = v.(string)
		ipSegment.Ip = helper.String(ip)
	}
	if v, ok := d.GetOkExists("mask"); ok {
		ipMask = v.(int)
		ipSegment.Mask = helper.IntUint64(ipMask)
	}
	request.IpList = append(request.IpList, &ipSegment)

	if v, ok := d.GetOk("type"); ok {
		ipType = v.(string)
		request.Type = helper.String(ipType)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAntiddosClient().CreateDDoSBlackWhiteIpList(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create antiddos ddosBlackWhiteIp failed, reason:%+v", logId, err)
		return err
	}

	idItems := []string{instanceId, ip}
	d.SetId(strings.Join(idItems, FILED_SP))

	return resourceTencentCloudAntiddosDdosBlackWhiteIpRead(d, meta)
}

func resourceTencentCloudAntiddosDdosBlackWhiteIpRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_ddos_black_white_ip.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := idSplit[0]
	ip := idSplit[1]
	ddosBlackWhiteIpListResponseParams, err := service.DescribeAntiddosDdosBlackWhiteIpListById(ctx, instanceId)
	if err != nil {
		return err
	}
	_ = d.Set("instance_id", instanceId)

	if len(ddosBlackWhiteIpListResponseParams.BlackIpList) > 0 {
		var targetIpSegment *antiddos.IpSegment
		ipList := ddosBlackWhiteIpListResponseParams.BlackIpList
		for i := range ipList {
			if *(ipList[i].Ip) == ip {
				targetIpSegment = ipList[i]
				break
			}
		}
		if targetIpSegment != nil {
			_ = d.Set("type", IP_TYPE_BLACK)
			_ = d.Set("ip", targetIpSegment.Ip)
			_ = d.Set("mask", targetIpSegment.Mask)
		}
	}

	if len(ddosBlackWhiteIpListResponseParams.WhiteIpList) > 0 {
		var targetIpSegment *antiddos.IpSegment
		ipList := ddosBlackWhiteIpListResponseParams.WhiteIpList
		for i := range ipList {
			if *(ipList[i].Ip) == ip {
				targetIpSegment = ipList[i]
				break
			}
		}
		if targetIpSegment != nil {
			_ = d.Set("type", IP_TYPE_WHITE)
			_ = d.Set("ip", targetIpSegment.Ip)
			_ = d.Set("mask", targetIpSegment.Mask)
		}
	}

	return nil
}

func resourceTencentCloudAntiddosDdosBlackWhiteIpUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_ddos_black_white_ip.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := antiddos.NewModifyDDoSBlackWhiteIpListRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "ip"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	isChange := false
	if d.HasChange("type") || d.HasChange("mask") {
		isChange = true
		oldType, newType := d.GetChange("type")
		ip := d.Get("ip")
		oldMask, newMask := d.GetChange("mask")

		request.OldIpType = helper.String(oldType.(string))
		request.NewIpType = helper.String(newType.(string))

		oldIpSegment := antiddos.IpSegment{
			Ip:   helper.String(ip.(string)),
			Mask: helper.IntUint64(oldMask.(int)),
		}
		newIpSegment := antiddos.IpSegment{
			Ip:   helper.String(ip.(string)),
			Mask: helper.IntUint64(newMask.(int)),
		}
		request.OldIp = &oldIpSegment
		request.NewIp = &newIpSegment
	}

	if isChange {
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseAntiddosClient().ModifyDDoSBlackWhiteIpList(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update antiddos ddosBlackWhiteIp failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAntiddosDdosBlackWhiteIpRead(d, meta)
}

func resourceTencentCloudAntiddosDdosBlackWhiteIpDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_antiddos_ddos_black_white_ip.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	params := make(map[string]interface{})
	params["instanceId"] = idSplit[0]
	params["ip"] = idSplit[1]
	params["ipType"] = d.Get("type").(string)
	params["ipMask"] = d.Get("mask").(int)
	if err := service.DeleteAntiddosDdosBlackWhiteIpListById(ctx, params); err != nil {
		return err
	}

	return nil
}
