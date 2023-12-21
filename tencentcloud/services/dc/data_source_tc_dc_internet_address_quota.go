package dc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcInternetAddressQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcInternetAddressQuotaRead,
		Schema: map[string]*schema.Schema{
			"ipv6_prefix_len": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The minimum prefix length allowed on the IPv6 Internet public network.",
			},

			"ipv4_bgp_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "BGP type IPv4 Internet address quota.",
			},

			"ipv4_other_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Non-BGP type IPv4 Internet address quota.",
			},

			"ipv4_bgp_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of used BGP type IPv4 Internet addresses.",
			},

			"ipv4_other_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of non-BGP Internet addresses used.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDcInternetAddressQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dc_internet_address_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var quota *dc.DescribeInternetAddressQuotaResponse

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcInternetAddressQuota(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		quota = result
		return nil
	})
	if err != nil {
		return err
	}

	if quota.Response.Ipv6PrefixLen != nil {
		_ = d.Set("ipv6_prefix_len", quota.Response.Ipv6PrefixLen)
	}

	if quota.Response.Ipv4BgpQuota != nil {
		_ = d.Set("ipv4_bgp_quota", quota.Response.Ipv4BgpQuota)
	}

	if quota.Response.Ipv4OtherQuota != nil {
		_ = d.Set("ipv4_other_quota", quota.Response.Ipv4OtherQuota)
	}

	if quota.Response.Ipv4BgpNum != nil {
		_ = d.Set("ipv4_bgp_num", quota.Response.Ipv4BgpNum)
	}

	if quota.Response.Ipv4OtherNum != nil {
		_ = d.Set("ipv4_other_num", quota.Response.Ipv4OtherNum)
	}

	tmpList := []map[string]interface{}{
		{
			"ipv6_prefix_len":  quota.Response.Ipv6PrefixLen,
			"ipv4_bgp_quota":   quota.Response.Ipv4BgpQuota,
			"ipv4_other_quota": quota.Response.Ipv4OtherQuota,
			"ipv4_bgp_num":     quota.Response.Ipv4BgpNum,
			"ipv4_other_num":   quota.Response.Ipv4OtherNum,
		},
	}

	d.SetId(helper.Int64ToStr(*quota.Response.Ipv4BgpQuota))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
