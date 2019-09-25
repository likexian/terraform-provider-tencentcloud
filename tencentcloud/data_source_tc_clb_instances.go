/*
Use this data source to query detailed information of CLB

Example Usage

```hcl
data "tencentcloud_clb_instances" "foo" {
    clb_id             = "lb-k2zjp9lv"
    network_type       = "OPEN"
    clb_name           = "myclb"
    project_id         = 0
    result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func dataSourceTencentCloudClbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstancesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the CLB to be queried.",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_NETWORK_TYPE),
				Description:  "Type of CLB instance, and available values include 'OPEN' and 'INTERNAL'.",
			},
			"clb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CLB to be queried.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project id of the CLB.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"clb_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of cloud load balancers. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of CLB.",
						},
						"clb_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CLB.",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of CLB.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Id of the project.",
						},
						"clb_vips": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The virtual service address table of the CLB.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of CLB.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the CLB.",
						},
						"status_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest state transition time of CLB.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the subnet.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Id of the security groups.",
						},
						"target_region_info_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region information of backend service are attached the CLB.",
						},
						"target_region_info_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VpcId information of backend service are attached the CLB.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this CLB.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_clb_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("clb_id"); ok {
		params["clb_id"] = v.(string)
	}
	if v, ok := d.GetOk("clb_name"); ok {
		params["clb_name"] = v.(string)
	}
	if v, ok := d.GetOk("project_id"); ok {
		params["project_id"] = v.(int)
	}
	if v, ok := d.GetOk("network_type"); ok {
		params["network_type"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var clbs []*clb.LoadBalancer
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		clbs = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read clb instances failed, reason:%s\n ", logId, err.Error())
		return err
	}
	clbList := make([]map[string]interface{}, 0, len(clbs))
	ids := make([]string, 0, len(clbs))
	for _, clb := range clbs {
		mapping := map[string]interface{}{
			"clb_id":                    *clb.LoadBalancerId,
			"clb_name":                  *clb.LoadBalancerName,
			"network_type":              *clb.LoadBalancerType,
			"status":                    *clb.Status,
			"create_time":               *clb.CreateTime,
			"status_time":               *clb.StatusTime,
			"project_id":                *clb.ProjectId,
			"vpc_id":                    *clb.VpcId,
			"subnet_id":                 *clb.SubnetId,
			"clb_vips":                  flattenStringList(clb.LoadBalancerVips),
			"target_region_info_region": *(clb.TargetRegionInfo.Region),
			"target_region_info_vpc_id": *(clb.TargetRegionInfo.VpcId),
			"security_groups":           flattenStringList(clb.SecureGroups),
		}
		if clb.Tags != nil {
			tags := make(map[string]interface{}, len(clb.Tags))
			for _, t := range clb.Tags {
				tags[*t.TagKey] = *t.TagValue
			}
			mapping["tags"] = tags
		}
		clbList = append(clbList, mapping)
		ids = append(ids, *clb.LoadBalancerId)
	}

	d.SetId(dataResourceIdsHash(ids))
	if e := d.Set("clb_list", clbList); e != nil {
		log.Printf("[CRITAL]%s provider set clb list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), clbList); e != nil {
			return e
		}
	}

	return nil
}
