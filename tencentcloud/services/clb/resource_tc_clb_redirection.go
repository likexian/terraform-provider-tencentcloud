package clb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbRedirection() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbRedirectionCreate,
		Read:   resourceTencentCloudClbRedirectionRead,
		Update: resourceTencentCloudClbRedirectionUpdate,
		Delete: resourceTencentCloudClbRedirectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of CLB instance.",
			},
			"source_listener_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,

				Description: "ID of source listener.",
			},
			"target_listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of source listener.",
			},
			"source_rule_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,

				Description: "Rule ID of source listener.",
			},
			"target_rule_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Rule ID of target listener.",
			},
			"is_auto_rewrite": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "Indicates whether automatic forwarding is enable, default is `false`. If enabled, the source listener and location should be empty, the target listener must be https protocol and port is 443.",
			},
			"delete_all_auto_rewrite": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Indicates whether delete all auto redirection. Default is `false`. It will take effect only when this redirection is auto-rewrite and this auto-rewrite auto redirected more than one rules. All the auto-rewrite relations will be deleted when this parameter set true.",
			},
		},
	}
}

func resourceTencentCloudClbRedirectionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_redirection.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbId := d.Get("clb_id").(string)
	targetListenerId := d.Get("target_listener_id").(string)
	checkErr := ListenerIdCheck(targetListenerId)
	if checkErr != nil {
		return checkErr
	}
	targetLocId := d.Get("target_rule_id").(string)
	checkErr = RuleIdCheck(targetLocId)
	if checkErr != nil {
		return checkErr
	}
	sourceListenerId := ""
	sourceLocId := ""
	if v, ok := d.GetOk("source_listener_id"); ok {
		sourceListenerId = v.(string)
		checkErr := ListenerIdCheck(sourceListenerId)
		if checkErr != nil {
			return checkErr
		}
	}
	if v, ok := d.GetOk("source_rule_id"); ok {
		sourceLocId = v.(string)
		checkErr = RuleIdCheck(sourceLocId)
		if checkErr != nil {
			return checkErr
		}
	}

	//check is auto forwarding or not
	isAutoRewrite := false
	if v, ok := d.GetOkExists("is_auto_rewrite"); ok {
		isAutoRewrite = v.(bool)
	}

	if isAutoRewrite {
		request := clb.NewAutoRewriteRequest()

		request.LoadBalancerId = helper.String(clbId)
		request.ListenerId = helper.String(targetListenerId)

		//check target listener is https:443
		clbService := ClbService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		protocol := ""
		port := -1
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, e := clbService.DescribeListenerById(ctx, targetListenerId, clbId)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if instance == nil {
				return resource.NonRetryableError(fmt.Errorf("[CLB redirection][Create] the queried instance is empty [DescribeListenerById]"))
			}
			if instance.Protocol == nil || instance.Port == nil {
				return resource.NonRetryableError(fmt.Errorf("[CLB redirection][Create] protocol or port is nil, get protocol and port fail [DescribeListenerById]"))
			}
			protocol = *(instance.Protocol)
			port = int(*(instance.Port))
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s get CLB listener failed, reason:%+v", logId, err)
			return err
		}

		if protocol == CLB_LISTENER_PROTOCOL_HTTPS && port != AUTO_TARGET_PORT {
			return fmt.Errorf("[CHECK][CLB redirection][Create] check: The target listener must be https:443 when applying auto rewrite")
		}

		//get host array from location
		filter := map[string]string{"rule_id": targetLocId, "listener_id": targetListenerId, "clb_id": clbId}
		var instances []*clb.RuleOutput
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			results, e := clbService.DescribeRulesByFilter(ctx, filter)
			if e != nil {
				return tccommon.RetryError(e)
			}
			instances = results
			return nil

		})
		if err != nil {
			log.Printf("[CRITAL]%s read CLB listener rule failed, reason:%+v", logId, err)
			return err
		}

		if len(instances) == 0 {
			return fmt.Errorf("[CHECK][CLB redirection][Create] check: rule %s not found!", targetLocId)
		}
		instance := instances[0]
		domain := instance.Domain
		url := instance.Url
		request.Domains = []*string{domain}
		//check source listener is null
		if sourceListenerId != "" || sourceLocId != "" {
			return fmt.Errorf("[CHECK][CLB redirection][Create] check: auto rewrite cannot specify source")
		}
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().AutoRewrite(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create CLB redirection failed, reason:%+v", logId, err)
			return err
		}

		params := make(map[string]interface{})
		params["clb_id"] = clbId
		params["port"] = AUTO_SOURCE_PORT
		params["protocol"] = CLB_LISTENER_PROTOCOL_HTTP
		var listeners []*clb.Listener
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			results, e := clbService.DescribeListenersByFilter(ctx, params)
			if e != nil {
				return tccommon.RetryError(e)
			}
			listeners = results
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read CLB listeners failed, reason:%+v", logId, err)
			return err
		}
		if len(listeners) == 0 {
			return fmt.Errorf("[CHECK][CLB redirection][Create] check: listener not found!")
		}
		listener := listeners[0]
		sourceListenerId = *listener.ListenerId
		rparams := make(map[string]string)
		rparams["clb_id"] = clbId
		rparams["domain"] = *domain
		rparams["url"] = *url
		rparams["listener_id"] = sourceListenerId
		rparams["url"] = *url
		var rules []*clb.RuleOutput
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			results, e := clbService.DescribeRulesByFilter(ctx, rparams)
			if e != nil {
				return tccommon.RetryError(e)
			}
			rules = results
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read CLB listener rules failed, reason:%+v", logId, err)
			return err
		}
		if len(rules) == 0 {
			return fmt.Errorf("[CHECK][CLB redirection][Create] check: rule not found!")
		}

		rule := rules[0]
		sourceLocId = *rule.LocationId

	} else {
		request := clb.NewManualRewriteRequest()

		request.LoadBalancerId = helper.String(clbId)
		request.SourceListenerId = helper.String(sourceListenerId)
		request.TargetListenerId = helper.String(targetListenerId)

		var rewriteInfo clb.RewriteLocationMap
		rewriteInfo.SourceLocationId = helper.String(sourceLocId)
		rewriteInfo.TargetLocationId = helper.String(targetLocId)
		request.RewriteInfos = []*clb.RewriteLocationMap{&rewriteInfo}
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ManualRewrite(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create CLB redirection failed, reason:%+v", logId, err)
			return err
		}

	}

	d.SetId(sourceLocId + "#" + targetLocId + "#" + sourceListenerId + "#" + targetListenerId + "#" + clbId)

	return resourceTencentCloudClbRedirectionRead(d, meta)
}

func resourceTencentCloudClbRedirectionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_redirection.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	rewriteId := d.Id()
	isAutoRewrite := false
	if v, ok := d.GetOkExists("is_auto_rewrite"); ok {
		isAutoRewrite = v.(bool)
		_ = d.Set("is_auto_rewrite", isAutoRewrite)
	}
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *map[string]string
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeRedirectionById(ctx, rewriteId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB redirection failed, reason:%+v", logId, err)
		return err
	}

	if instance == nil || len(*instance) == 0 {
		d.SetId("")
		return nil
	}

	_ = d.Set("clb_id", (*instance)["clb_id"])
	_ = d.Set("source_listener_id", (*instance)["source_listener_id"])
	_ = d.Set("target_listener_id", (*instance)["target_listener_id"])
	_ = d.Set("source_rule_id", (*instance)["source_rule_id"])
	_ = d.Set("target_rule_id", (*instance)["target_rule_id"])

	return nil
}

func resourceTencentCloudClbRedirectionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_redirection.update")()
	defer tccommon.InconsistentCheck(d, meta)()
	// this nil update method works for the only filed `delete_all_auto_rewrite`
	return resourceTencentCloudClbRedirectionRead(d, meta)
}

func resourceTencentCloudClbRedirectionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_redirection.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	deleteAll := d.Get("delete_all_auto_rewrite").(bool)
	isAutoRewrite := d.Get("is_auto_rewrite").(bool)
	if deleteAll && isAutoRewrite {
		//delete all the auto rewrite
		var rewrites []*map[string]string
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, inErr := clbService.DescribeAllAutoRedirections(ctx, id)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			rewrites = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s delete CLB redirection failed, reason:%+v", logId, err)
			return err
		}

		for _, rewrite := range rewrites {
			if rewrite == nil {
				continue
			}
			rewriteId := (*rewrite)["source_rule_id"] + tccommon.FILED_SP + (*rewrite)["target_rule_id"] + tccommon.FILED_SP + (*rewrite)["source_listener_id"] + tccommon.FILED_SP + (*rewrite)["target_listener_id"] + tccommon.FILED_SP + (*rewrite)["clb_id"]
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				e := clbService.DeleteRedirectionById(ctx, rewriteId)
				if e != nil {
					return tccommon.RetryError(e)
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s delete CLB redirection failed, reason:%+v", logId, err)
				return err
			}
		}
	} else {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := clbService.DeleteRedirectionById(ctx, id)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s delete CLB redirection failed, reason:%+v", logId, err)
			return err
		}
	}
	return nil
}
