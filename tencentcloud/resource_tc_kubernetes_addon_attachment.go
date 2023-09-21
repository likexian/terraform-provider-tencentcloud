/*
Provide a resource to configure kubernetes cluster app addons.

~> **NOTE**: Avoid to using legacy "1.0.0" version, leave the versions empty so we can fetch the latest while creating.

Example Usage

Install cbs addon by passing values

```hcl

resource "tencentcloud_kubernetes_addon_attachment" "addon_cbs" {
  cluster_id = "cls-xxxxxxxx"
  name = "cbs"
  # version = "1.0.5"
  values = [
    "rootdir=/var/lib/kubelet"
  ]
}
```

Install tcr addon by passing values

```hcl
resource "tencentcloud_kubernetes_addon_attachment" "addon_tcr" {
  cluster_id = "cls-xxxxxxxx" #specify your tke cluster id
  name       = "tcr"
  version    = "1.0.0"
  values = [
    # imagePullSecretsCrs is an array which can configure image pull
    "global.imagePullSecretsCrs[0].name=${local.tcr_id}-vpc",                              #specify a unique name, invalid format as: `${tcrId}-vpc`
    "global.imagePullSecretsCrs[0].namespaces=${local.ns_name}",                           #input the specified namespaces of the cluster, or input `*` for all.
    "global.imagePullSecretsCrs[0].serviceAccounts=*",                                     #input the specified service account of the cluster, or input `*` for all.
    "global.imagePullSecretsCrs[0].type=docker",                                           #only support docker now
    "global.imagePullSecretsCrs[0].dockerUsername=${local.user_name}",                     #input the access username, or you can create it from `tencentcloud_tcr_token`
    "global.imagePullSecretsCrs[0].dockerPassword=${local.token}",                         #input the access token, or you can create it from `tencentcloud_tcr_token`
    "global.imagePullSecretsCrs[0].dockerServer=${local.tcr_name}-vpc.tencentcloudcr.com", #invalid format as: `${tcr_name}-vpc.tencentcloudcr.com`
    "global.imagePullSecretsCrs[1].name=${local.tcr_id}-public",                           #specify a unique name, invalid format as: `${tcr_id}-public`
    "global.imagePullSecretsCrs[1].namespaces=${local.ns_name}",
    "global.imagePullSecretsCrs[1].serviceAccounts=*",
    "global.imagePullSecretsCrs[1].type=docker",
    "global.imagePullSecretsCrs[1].dockerUsername=${local.user_name}",                 #refer to previous description
    "global.imagePullSecretsCrs[1].dockerPassword=${local.token}",                     #refer to previous description
    "global.imagePullSecretsCrs[1].dockerServer=${local.tcr_name}.tencentcloudcr.com", #invalid format as: `${tcr_name}.tencentcloudcr.com`
    "global.cluster.region=gz",
    "global.cluster.longregion=ap-guangzhou",
    # Specify global hosts(optional), the numbers of hosts must be matched with the numbers of imagePullSecretsCrs
    "global.hosts[0].domain=${local.tcr_name}-vpc.tencentcloudcr.com",                 #Corresponds to the dockerServer in the imagePullSecretsCrs above
    "global.hosts[0].ip=${local.end_point}",                                           #input InternalEndpoint of tcr instance, you can get it from data source `tencentcloud_tcr_instances`
    "global.hosts[0].disabled=false",                                                  #disabled this host config or not
    "global.hosts[1].domain=${local.tcr_name}.tencentcloudcr.com",
    "global.hosts[1].ip=${local.end_point}",
    "global.hosts[1].disabled=false",
  ]
}

locals {
  tcr_id   = tencentcloud_tcr_instance.mytcr.id
  tcr_name = tencentcloud_tcr_instance.mytcr.name
  ns_name   = tencentcloud_tcr_namespace.my_ns.name
  user_name = tencentcloud_tcr_token.my_token.user_name
  token     = tencentcloud_tcr_token.my_token.token
  end_point = data.tencentcloud_tcr_instances.my_ins.instance_list.0.internal_end_point
}

resource "tencentcloud_tcr_token" "my_token" {
  instance_id = local.tcr_id
  description = "tcr token"
}

data "tencentcloud_tcr_instances" "my_ins" {
  instance_id = local.tcr_id
}

resource "tencentcloud_tcr_instance" "mytcr" {
  name          = "tf-test-tcr-addon"
  instance_type = "basic"
  delete_bucket = true

  tags = {
    test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "my_ns" {
  instance_id    = local.tcr_id
  name           = "tf_test_tcr_ns_addon"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}
```

Install new addon by passing spec json to req_body directly

```hcl
resource "tencentcloud_kubernetes_addon_attachment" "addon_cbs" {
  cluster_id = "cls-xxxxxxxx"
  request_body = <<EOF
  {
    "spec":{
        "chart":{
            "chartName":"cbs",
            "chartVersion":"1.0.5"
        },
        "values":{
            "rawValuesType":"yaml",
            "values":[
              "rootdir=/var/lib/kubelet"
            ]
        }
    }
  }
EOF
}
```

Import

Addon can be imported by using cluster_id#addon_name
```
$ terraform import tencentcloud_kubernetes_addon_attachment.addon_cos cls-xxxxxxxx#cos
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTkeAddonAttachment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of addon.",
			},
			"version": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Addon version, default latest version. Conflict with `request_body`.",
				ConflictsWith: []string{"request_body"},
			},
			"values": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Description:   "Values the addon passthroughs. Conflict with `request_body`.",
				ConflictsWith: []string{"request_body"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"request_body": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Serialized json string as request body of addon spec. If set, will ignore `version` and `values`.",
				ConflictsWith: []string{"version", "values"},
			},
			"response_body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Addon response body.",
			},
			"status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Addon current status.",
			},
		},
		Create: resourceTencentCloudTkeAddonAttachmentCreate,
		Update: resourceTencentCloudTkeAddonAttachmentUpdate,
		Read:   resourceTencentCloudTkeAddonAttachmentRead,
		Delete: resourceTencentCloudTkeAddonAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceTencentCloudTkeAddonAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_addon_attachment.create")()
	logId := getLogId(contextNil)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		clusterId = d.Get("cluster_id").(string)
		addonName = d.Get("name").(string)
		version   = d.Get("version").(string)
		values    = d.Get("values").([]interface{})
		reqBody   = d.Get("request_body").(string)
	)

	if version == "" {
		request := tke.NewGetTkeAppChartListRequest()
		chartList, err := service.GetTkeAppChartList(ctx, request)
		if err != nil {
			return fmt.Errorf("error while fetching latest chart versions, %s", err.Error())
		}
		for i := range chartList {
			chart := chartList[i]
			if *chart.Name == addonName {
				version = *chart.LatestVersion
				break
			}
		}
	}

	if reqBody == "" {
		var reqErr error
		v := helper.InterfacesStringsPoint(values)
		reqBody, reqErr = service.GetAddonReqBody(addonName, version, v)
		if reqErr != nil {
			return reqErr
		}
	}

	err := service.CreateExtensionAddon(ctx, clusterId, reqBody)

	if err != nil {
		return err
	}

	d.SetId(clusterId + FILED_SP + addonName)

	resData := &AddonResponseData{}
	reason := "unknown error"
	phase, has, _ := service.PollingAddonsPhase(ctx, clusterId, addonName, resData)

	if resData.Status != nil && resData.Status["reason"] != nil {
		reason = resData.Status["reason"].(string)
	}

	if !has {
		return fmt.Errorf("addon %s not exists", addonName)
	}

	if phase == "ChartFetchFailed" || phase == "Failed" || phase == "RollbackFailed" || phase == "SyncFailed" {
		msg := fmt.Sprintf("Unexpected chart phase `%s`, reason: %s", phase, reason)
		if err := resourceTencentCloudTkeAddonAttachmentDelete(d, meta); err != nil {
			return err
		}
		d.SetId("")
		return fmt.Errorf(msg)
	}

	return resourceTencentCloudTkeAddonAttachmentRead(d, meta)
}

func resourceTencentCloudTkeAddonAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_addon_attachment.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()
	has := false
	split := strings.Split(id, FILED_SP)
	if len(split) < 2 {
		return fmt.Errorf("id expected format: cluster_id#addon_name but no addon_name provided")
	}
	clusterId := split[0]
	addonName := split[1]

	var (
		err               error
		response          string
		addonResponseData = &AddonResponseData{}
	)

	_, has, err = service.PollingAddonsPhase(ctx, clusterId, addonName, addonResponseData)

	if err != nil || !has {
		d.SetId("")
		return err
	}

	response, _, err = service.DescribeExtensionAddon(ctx, clusterId, addonName)

	if err != nil {
		d.SetId("")
		return err
	}

	_ = d.Set("response_body", response)

	spec := addonResponseData.Spec
	statuses := addonResponseData.Status

	if spec != nil {
		_ = d.Set("cluster_id", clusterId)
		_ = d.Set("name", spec.Chart.ChartName)
		_ = d.Set("version", spec.Chart.ChartVersion)
		if spec.Values != nil && len(spec.Values.Values) > 0 {

			// Filter auto-filled values from addon creation
			filteredValues := getFilteredValues(d, spec.Values.Values)
			_ = d.Set("values", filteredValues)
		}
	}

	if statuses != nil || len(statuses) == 0 {
		strMap := helper.CovertInterfaceMapToStrPtr(statuses)
		err := d.Set("status", strMap)
		if err != nil {
			return err
		}
	}

	d.SetId(id)

	return nil
}

func resourceTencentCloudTkeAddonAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_addon_attachment.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		id        = d.Id()
		split     = strings.Split(id, FILED_SP)
		clusterId = split[0]
		addonName = split[1]
		version   = d.Get("version").(string)
		values    = d.Get("values").([]interface{})
		reqBody   = d.Get("request_body").(string)
		err       error
	)

	if d.HasChange("request_body") && reqBody == "" || d.HasChange("version") || d.HasChange("values") {
		reqBody, err = service.GetAddonReqBody(addonName, version, helper.InterfacesStringsPoint(values))
	}

	if err != nil {
		return err
	}

	err = service.UpdateExtensionAddon(ctx, clusterId, addonName, reqBody)

	if err != nil {
		return err
	}

	return resourceTencentCloudTkeAddonAttachmentRead(d, meta)
}

func resourceTencentCloudTkeAddonAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_addon_attachment.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		id        = d.Id()
		split     = strings.Split(id, FILED_SP)
		clusterId = split[0]
		addonName = split[1]
		has       bool
	)

	if err := service.DeleteExtensionAddon(ctx, clusterId, addonName); err != nil {
		return err
	}

	// check if addon terminating or still exists
	_, has, _ = service.PollingAddonsPhase(ctx, clusterId, addonName, nil)

	if has {
		return fmt.Errorf("addon %s still exists", addonName)
	}

	return nil
}

func getFilteredValues(d *schema.ResourceData, values []*string) []string {
	rawValues := helper.InterfacesStrings(d.Get("values").([]interface{}))

	for _, value := range values {
		kv := strings.Split(*value, "=")
		key := kv[0]

		if IsContains(TKE_ADDON_DEFAULT_VALUES_KEY, key) || IsContains(rawValues, *value) {
			continue
		}
		rawValues = append(rawValues, *value)
	}
	return rawValues
}
