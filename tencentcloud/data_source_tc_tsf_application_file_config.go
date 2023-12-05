package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfApplicationFileConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApplicationFileConfigRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration ID.",
			},

			"config_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of configuration item ID.",
			},

			"config_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration item name.",
			},

			"application_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Application ID.",
			},

			"config_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration item version.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "File configuration item list. Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "total count.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "File configuration array. Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Config ID. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item name. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration version. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_version_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item version description. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_file_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item file name. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_file_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration file content. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_file_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration file code. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"creation_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CreationTime. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "application Id. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "application name. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"delete_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "delete flag, true: allow delete; false: delete prohibit.",
									},
									"config_version_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "config version count.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "last update time.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_file_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "file config path. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_post_cmd": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "last update time.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"config_file_value_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "config item content length.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
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

func dataSourceTencentCloudTsfApplicationFileConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_application_file_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("config_id"); ok {
		paramMap["ConfigId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_id_list"); ok {
		configIdListSet := v.(*schema.Set).List()
		paramMap["ConfigIdList"] = helper.InterfacesStringsPoint(configIdListSet)
	}

	if v, ok := d.GetOk("config_name"); ok {
		paramMap["ConfigName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_version"); ok {
		paramMap["ConfigVersion"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var config *tsf.TsfPageFileConfig
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApplicationFileConfigByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		config = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(config.Content))
	tsfPageFileConfigMap := map[string]interface{}{}
	if config != nil {
		if config.TotalCount != nil {
			tsfPageFileConfigMap["total_count"] = config.TotalCount
		}

		if config.Content != nil {
			contentList := []interface{}{}
			for _, content := range config.Content {
				contentMap := map[string]interface{}{}

				if content.ConfigId != nil {
					contentMap["config_id"] = content.ConfigId
				}

				if content.ConfigName != nil {
					contentMap["config_name"] = content.ConfigName
				}

				if content.ConfigVersion != nil {
					contentMap["config_version"] = content.ConfigVersion
				}

				if content.ConfigVersionDesc != nil {
					contentMap["config_version_desc"] = content.ConfigVersionDesc
				}

				if content.ConfigFileName != nil {
					contentMap["config_file_name"] = content.ConfigFileName
				}

				if content.ConfigFileValue != nil {
					contentMap["config_file_value"] = content.ConfigFileValue
				}

				if content.ConfigFileCode != nil {
					contentMap["config_file_code"] = content.ConfigFileCode
				}

				if content.CreationTime != nil {
					contentMap["creation_time"] = content.CreationTime
				}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.DeleteFlag != nil {
					contentMap["delete_flag"] = content.DeleteFlag
				}

				if content.ConfigVersionCount != nil {
					contentMap["config_version_count"] = content.ConfigVersionCount
				}

				if content.LastUpdateTime != nil {
					contentMap["last_update_time"] = content.LastUpdateTime
				}

				if content.ConfigFilePath != nil {
					contentMap["config_file_path"] = content.ConfigFilePath
				}

				if content.ConfigPostCmd != nil {
					contentMap["config_post_cmd"] = content.ConfigPostCmd
				}

				if content.ConfigFileValueLength != nil {
					contentMap["config_file_value_length"] = content.ConfigFileValueLength
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.ConfigId)
			}

			tsfPageFileConfigMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageFileConfigMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageFileConfigMap); e != nil {
			return e
		}
	}
	return nil
}
