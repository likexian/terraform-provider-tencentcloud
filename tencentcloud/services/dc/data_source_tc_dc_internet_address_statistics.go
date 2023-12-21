package dc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcInternetAddressStatistics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcInternetAddressStatisticsRead,
		Schema: map[string]*schema.Schema{
			"internet_address_statistics": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Statistical Information List of Internet Public Network Addresses.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"subnet_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of Internet public network addresses.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDcInternetAddressStatisticsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dc_internet_address_statistics.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var internetAddressStatistics []*dc.InternetAddressStatistics

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcInternetAddressStatistics(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		internetAddressStatistics = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(internetAddressStatistics))
	tmpList := make([]map[string]interface{}, 0, len(internetAddressStatistics))

	if internetAddressStatistics != nil {
		for _, internetAddressStatistics := range internetAddressStatistics {
			internetAddressStatisticsMap := map[string]interface{}{}

			if internetAddressStatistics.Region != nil {
				internetAddressStatisticsMap["region"] = internetAddressStatistics.Region
			}

			if internetAddressStatistics.SubnetNum != nil {
				internetAddressStatisticsMap["subnet_num"] = internetAddressStatistics.SubnetNum
			}

			ids = append(ids, *internetAddressStatistics.Region)
			tmpList = append(tmpList, internetAddressStatisticsMap)
		}

		_ = d.Set("internet_address_statistics", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
