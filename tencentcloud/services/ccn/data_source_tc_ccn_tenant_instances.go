package ccn

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCcnTenantInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTenantCcnRead,
		Schema: map[string]*schema.Schema{
			"ccn_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "filter by ccn ids, like: ['ccn-12345678'].",
			},

			"user_account_id": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "filter by ccn ids, like: ['12345678'].",
			},

			"is_security_lock": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "filter by locked, like ['true'].",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTenantCcnRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tenant_ccn.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("ccn_ids"); ok {
		valuesSet := v.(*schema.Set).List()
		paramMap["ccn-ids"] = helper.InterfacesStringsPoint(valuesSet)
	}

	if v, ok := d.GetOk("user_account_id"); ok {
		valuesSet := v.(*schema.Set).List()
		paramMap["user-account-id"] = helper.InterfacesStringsPoint(valuesSet)
	}

	if v, ok := d.GetOk("is_security_lock"); ok {
		valuesSet := v.(*schema.Set).List()
		paramMap["is-security-lock"] = helper.InterfacesStringsPoint(valuesSet)
	}
	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var ccnSet []*vpc.CcnInstanceInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTenantCcnByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		ccnSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ccnSet))
	tmpList := make([]map[string]interface{}, 0, len(ccnSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
