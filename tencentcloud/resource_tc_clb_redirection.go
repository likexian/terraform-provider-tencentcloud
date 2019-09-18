/*
Provides a resource to create a CLB redirection.

Example Usage

Manual Rewrite

```hcl
resource "tencentcloud_clb_redirection" "foo" {
  clb_id             = "lb-p7olt9e5"
  source_listener_id = "lbl-jc1dx6ju"
  target_listener_id = "lbl-asj1hzuo"
  source_rule_id     = "loc-ft8fmngv"
  target_rule_id     = "loc-4xxr2cy7"
}
```

Auto Rewrite

```hcl
resource "tencentcloud_clb_redirection" "foo" {
  clb_id             = "lb-p7olt9e5"
  target_listener_id = "lbl-asj1hzuo"
  target_rule_id     = "loc-4xxr2cy7"
  is_auto_rewrite    = true
}
```

Import

CLB redirection can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_redirection.foo loc-ft8fmngv#loc-4xxr2cy7#lbl-jc1dx6ju#lbl-asj1hzuo#lb-p7olt9e5
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbRedirection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbRedirectionCreate,
		Read:   resourceTencentCloudClbRedirectionRead,
		Delete: resourceTencentCloudClbRedirectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of CLB instance.",
			},
			"source_listener_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,

				Description: "Id of source listener.",
			},
			"target_listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of source listener.",
			},
			"source_rule_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,

				Description: "Rule id of source listener.",
			},
			"target_rule_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Rule id of target listener.",
			},
			"is_auto_rewrite": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "Indicates whether automatic forwarding is enable, default is false. If enabled, the source listener and location should be empty, the target listener must be https protocol and port is 443.",
			},
		},
	}
}

func resourceTencentCloudClbRedirectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_redirection.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbId := d.Get("clb_id").(string)
	targetListenerId := d.Get("target_listener_id").(string)
	targetLocId := d.Get("target_rule_id").(string)
	sourceListenerId := ""
	sourceLocId := ""
	if v, ok := d.GetOk("source_listener_id"); ok {
		sourceListenerId = v.(string)
	}
	if v, ok := d.GetOk("source_rule_id"); ok {
		sourceLocId = v.(string)
	}

	//check is auto forwarding or not
	isAutoRewrite := false
	if v, ok := d.GetOk("is_auto_rewrite"); ok {
		isAutoRewrite = v.(bool)
	}

	if isAutoRewrite {
		request := clb.NewAutoRewriteRequest()

		request.LoadBalancerId = stringToPointer(clbId)
		request.ListenerId = stringToPointer(targetListenerId)

		//check target listener is https:443
		clbService := ClbService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		protocol := ""
		port := -1
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			instance, e := clbService.DescribeListenerById(ctx, targetListenerId, clbId)
			if e != nil {
				return retryError(e)
			}
			protocol = *(instance.Protocol)
			port = int(*(instance.Port))
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s get clb listener failed, reason:%s\n ", logId, err.Error())
			return err
		}

		if protocol == CLB_LISTENER_PROTOCOL_HTTPS && port != AUTO_TARGET_PORT {
			return fmt.Errorf("The target listener must be https:443 when applying auto rewrite")
		}

		//get host array from location
		filter := map[string]string{"rule_id": targetLocId, "listener_id": targetListenerId, "clb_id": clbId}
		var instances []*clb.RuleOutput
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			results, e := clbService.DescribeRulesByFilter(ctx, filter)
			if e != nil {
				return retryError(e)
			}
			instances = results
			return nil

		})
		if err != nil {
			log.Printf("[CRITAL]%s read clb listener rule failed, reason:%s\n ", logId, err.Error())
			return err
		}

		if len(instances) == 0 {
			return fmt.Errorf("rule not found!")
		}
		instance := instances[0]
		domain := instance.Domain
		request.Domains = []*string{domain}
		//check source listener is null
		if sourceListenerId != "" || sourceLocId != "" {
			return fmt.Errorf("auto rewrite cannot specify source")
		}
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().AutoRewrite(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(retryErr)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create clb redirection failed, reason:%s\n ", logId, err.Error())
			return err
		}

		params := make(map[string]interface{})
		params["clb_id"] = clbId
		params["port"] = AUTO_SOURCE_PORT
		params["protocol"] = CLB_LISTENER_PROTOCOL_HTTP
		var listeners []*clb.Listener
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			results, e := clbService.DescribeListenersByFilter(ctx, params)
			if e != nil {
				return retryError(e)
			}
			listeners = results
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read clb listeners failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(listeners) == 0 {
			return fmt.Errorf("listener not found!")
		}
		listener := listeners[0]
		sourceListenerId = *listener.ListenerId
		rparams := make(map[string]string)
		rparams["clb_id"] = clbId
		rparams["domain"] = *domain
		rparams["listener_id"] = sourceListenerId
		var rules []*clb.RuleOutput
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			results, e := clbService.DescribeRulesByFilter(ctx, rparams)
			if e != nil {
				return retryError(e)
			}
			rules = results
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read clb listener rules failed, reason:%s\n ", logId, err.Error())
			return err
		}
		if len(rules) == 0 {
			return fmt.Errorf("rule not found!")
		}
		rule := rules[0]
		sourceLocId = *rule.LocationId

	} else {
		request := clb.NewManualRewriteRequest()

		request.LoadBalancerId = stringToPointer(clbId)
		request.SourceListenerId = stringToPointer(sourceListenerId)
		request.TargetListenerId = stringToPointer(targetListenerId)

		var rewriteInfo clb.RewriteLocationMap
		rewriteInfo.SourceLocationId = stringToPointer(sourceLocId)
		rewriteInfo.TargetLocationId = stringToPointer(targetLocId)
		request.RewriteInfos = []*clb.RewriteLocationMap{&rewriteInfo}
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ManualRewrite(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(retryErr)
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create clb redirection failed, reason:%s\n ", logId, err.Error())
			return err
		}

	}

	d.SetId(sourceLocId + "#" + targetLocId + "#" + sourceListenerId + "#" + targetListenerId + "#" + clbId)

	return resourceTencentCloudClbRedirectionRead(d, meta)
}

func resourceTencentCloudClbRedirectionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_redirection.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	rewriteId := d.Id()
	isAutoRewrite := false
	if v, ok := d.GetOk("is_auto_rewrite"); ok {
		isAutoRewrite = v.(bool)
		d.Set("is_auto_rewrite", isAutoRewrite)
	}
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *map[string]string
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeRedirectionById(ctx, rewriteId)
		if e != nil {
			return retryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read clb redirection failed, reason:%s\n ", logId, err.Error())
		return err
	}
	d.Set("clb_id", (*instance)["clb_id"])
	d.Set("source_listener_id", (*instance)["source_listener_id"])
	d.Set("target_listener_id", (*instance)["target_listener_id"])
	d.Set("source_rule_id", (*instance)["source_rule_id"])
	d.Set("target_rule_id", (*instance)["target_rule_id"])

	return nil
}

func resourceTencentCloudClbRedirectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_redirection.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteRedirectionById(ctx, clbId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete clb redirection failed, reason:%s\n ", logId, err.Error())
		return err
	}
	return nil
}
