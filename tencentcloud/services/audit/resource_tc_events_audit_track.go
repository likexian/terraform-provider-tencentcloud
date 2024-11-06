package audit

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cloudaudit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEventsAuditTrack() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEventsAuditTrackCreate,
		Read:   resourceTencentCloudEventsAuditTrackRead,
		Update: resourceTencentCloudEventsAuditTrackUpdate,
		Delete: resourceTencentCloudEventsAuditTrackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Tracking set name, which can only contain 3-48 letters, digits, hyphens, and underscores.",
			},

			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Tracking set status (0: Not enabled; 1: Enabled).",
			},

			"storage": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Storage type of shipped data. Valid values: `cos`, `cls`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage type (Valid values: cos, cls).",
						},
						"storage_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "StorageRegion *string `json:'StorageRegion,omitnil,omitempty' name: 'StorageRegion'`.",
						},
						"storage_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage name. For COS, the storage name is the custom bucket name, which can contain up to 50 lowercase letters, digits, and hyphens. It cannot contain \"-APPID\" and cannot start or end with a hyphen. For CLS, the storage name is the log topic ID, which can contain 1-50 characters.",
						},
						"storage_prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage directory prefix. The COS log file prefix can only contain 3-40 letters and digits.",
						},
						"storage_account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Designated to store user ID.",
						},
						"storage_app_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Designated to store user app ID.",
						},
					},
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Data filtering criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_fields": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Resource filtering conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The product to which the tracking set event belongs. The value can be a single product such as `cos`, or `*` that indicates all products.",
									},
									"action_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tracking set event type (`Read`: Read; `Write`: Write; `*`: All).",
									},
									"event_names": {
										Type:        schema.TypeSet,
										Required:    true,
										Description: "The list of API names of tracking set events. When `ResourceType` is `*`, the value of `EventNames` must be `*`. When `ResourceType` is a specified product, the value of `EventNames` can be `*`. When `ResourceType` is `cos` or `cls`, up to 10 APIs are supported.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"track_for_all_members": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to enable the feature of shipping organization members operation logs to the organization admin account or the trusted service admin account (0: Not enabled; 1: Enabled. This feature can only be enabled by the organization admin account or the trusted service admin account).",
			},

			"track_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the log list has come to an end. `true`: Yes. Pagination is not required.",
			},
		},
	}
}

func resourceTencentCloudEventsAuditTrackCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_events_audit_track.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		trackId string
	)
	var (
		request  = cloudaudit.NewCreateEventsAuditTrackRequest()
		response = cloudaudit.NewCreateEventsAuditTrackResponse()
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	if storageMap, ok := helper.InterfacesHeadMap(d, "storage"); ok {
		storage := cloudaudit.Storage{}
		if v, ok := storageMap["storage_type"]; ok {
			storage.StorageType = helper.String(v.(string))
		}
		if v, ok := storageMap["storage_region"]; ok {
			storage.StorageRegion = helper.String(v.(string))
		}
		if v, ok := storageMap["storage_name"]; ok {
			storage.StorageName = helper.String(v.(string))
		}
		if v, ok := storageMap["storage_prefix"]; ok {
			storage.StoragePrefix = helper.String(v.(string))
		}
		if v, ok := storageMap["storage_account_id"]; ok {
			storage.StorageAccountId = helper.String(v.(string))
		}
		if v, ok := storageMap["storage_app_id"]; ok {
			storage.StorageAppId = helper.String(v.(string))
		}
		request.Storage = &storage
	}

	if filtersMap, ok := helper.InterfacesHeadMap(d, "filters"); ok {
		filter := cloudaudit.Filter{}
		if v, ok := filtersMap["resource_fields"]; ok {
			for _, item := range v.([]interface{}) {
				resourceFieldsMap := item.(map[string]interface{})
				resourceField := cloudaudit.ResourceField{}
				if v, ok := resourceFieldsMap["resource_type"]; ok {
					resourceField.ResourceType = helper.String(v.(string))
				}
				if v, ok := resourceFieldsMap["action_type"]; ok {
					resourceField.ActionType = helper.String(v.(string))
				}
				if v, ok := resourceFieldsMap["event_names"]; ok {
					eventNamesSet := v.(*schema.Set).List()
					for i := range eventNamesSet {
						eventNames := eventNamesSet[i].(string)
						resourceField.EventNames = append(resourceField.EventNames, helper.String(eventNames))
					}
				}
				filter.ResourceFields = append(filter.ResourceFields, &resourceField)
			}
		}
		request.Filters = &filter
	}

	if v, ok := d.GetOkExists("track_for_all_members"); ok {
		request.TrackForAllMembers = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAuditClient().CreateEventsAuditTrackWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create events audit track failed, reason:%+v", logId, err)
		return err
	}

	trackId = helper.UInt64ToStr(*response.Response.TrackId)

	d.SetId(trackId)

	return resourceTencentCloudEventsAuditTrackRead(d, meta)
}

func resourceTencentCloudEventsAuditTrackRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_events_audit_track.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := AuditService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	trackId := d.Id()

	respData, err := service.DescribeEventsAuditTrackById(ctx, trackId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `events_audit_track` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("track_id", helper.StrToInt(trackId))

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	storageMap := map[string]interface{}{}

	if respData.Storage != nil {
		if respData.Storage.StorageType != nil {
			storageMap["storage_type"] = respData.Storage.StorageType
		}

		if respData.Storage.StorageRegion != nil {
			storageMap["storage_region"] = respData.Storage.StorageRegion
		}

		if respData.Storage.StorageName != nil {
			storageMap["storage_name"] = respData.Storage.StorageName
		}

		if respData.Storage.StoragePrefix != nil {
			storageMap["storage_prefix"] = respData.Storage.StoragePrefix
		}

		_ = d.Set("storage", []interface{}{storageMap})
	}

	filtersMap := map[string]interface{}{}
	if respData.Filters != nil {
		resourceFieldsList := make([]map[string]interface{}, 0, len(respData.Filters.ResourceFields))
		if respData.Filters.ResourceFields != nil {
			for _, resourceFields := range respData.Filters.ResourceFields {
				resourceFieldsMap := map[string]interface{}{}

				if resourceFields.ResourceType != nil {
					resourceFieldsMap["resource_type"] = resourceFields.ResourceType
				}

				if resourceFields.ActionType != nil {
					resourceFieldsMap["action_type"] = resourceFields.ActionType
				}

				if resourceFields.EventNames != nil {
					resourceFieldsMap["event_names"] = resourceFields.EventNames
				}

				resourceFieldsList = append(resourceFieldsList, resourceFieldsMap)
			}

			filtersMap["resource_fields"] = resourceFieldsList
		}
		_ = d.Set("filters", []interface{}{filtersMap})
	}

	if respData.TrackForAllMembers != nil {
		_ = d.Set("track_for_all_members", respData.TrackForAllMembers)
	}

	_ = trackId
	return nil
}

func resourceTencentCloudEventsAuditTrackUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_events_audit_track.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	trackId := d.Id()

	needChange := false
	mutableArgs := []string{"status", "storage", "track_for_all_members", "filters"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := cloudaudit.NewModifyEventsAuditTrackRequest()

		request.TrackId = helper.StrToUint64Point(trackId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}

		if storageMap, ok := helper.InterfacesHeadMap(d, "storage"); ok {
			storage := cloudaudit.Storage{}
			if v, ok := storageMap["storage_type"]; ok {
				storage.StorageType = helper.String(v.(string))
			}
			if v, ok := storageMap["storage_region"]; ok {
				storage.StorageRegion = helper.String(v.(string))
			}
			if v, ok := storageMap["storage_name"]; ok {
				storage.StorageName = helper.String(v.(string))
			}
			if v, ok := storageMap["storage_prefix"]; ok {
				storage.StoragePrefix = helper.String(v.(string))
			}
			if v, ok := storageMap["storage_account_id"]; ok {
				storage.StorageAccountId = helper.String(v.(string))
			}
			if v, ok := storageMap["storage_app_id"]; ok {
				storage.StorageAppId = helper.String(v.(string))
			}
			request.Storage = &storage
		}

		if v, ok := d.GetOkExists("track_for_all_members"); ok {
			request.TrackForAllMembers = helper.IntUint64(v.(int))
		}

		if filtersMap, ok := helper.InterfacesHeadMap(d, "filters"); ok {
			filter := cloudaudit.Filter{}
			if v, ok := filtersMap["resource_fields"]; ok {
				for _, item := range v.([]interface{}) {
					resourceFieldsMap := item.(map[string]interface{})
					resourceField := cloudaudit.ResourceField{}
					if v, ok := resourceFieldsMap["resource_type"]; ok {
						resourceField.ResourceType = helper.String(v.(string))
					}
					if v, ok := resourceFieldsMap["action_type"]; ok {
						resourceField.ActionType = helper.String(v.(string))
					}
					if v, ok := resourceFieldsMap["event_names"]; ok {
						eventNamesSet := v.(*schema.Set).List()
						for i := range eventNamesSet {
							eventNames := eventNamesSet[i].(string)
							resourceField.EventNames = append(resourceField.EventNames, helper.String(eventNames))
						}
					}
					filter.ResourceFields = append(filter.ResourceFields, &resourceField)
				}
			}
			request.Filters = &filter
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAuditClient().ModifyEventsAuditTrackWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update events audit track failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = trackId
	return resourceTencentCloudEventsAuditTrackRead(d, meta)
}

func resourceTencentCloudEventsAuditTrackDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_events_audit_track.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	trackId := d.Id()

	var (
		request  = cloudaudit.NewDeleteAuditTrackRequest()
		response = cloudaudit.NewDeleteAuditTrackResponse()
	)

	request.TrackId = helper.StrToUint64Point(trackId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAuditClient().DeleteAuditTrackWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete events audit track failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	_ = trackId
	return nil
}
