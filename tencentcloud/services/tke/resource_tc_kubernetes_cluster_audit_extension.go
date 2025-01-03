package tke

import (
	"context"
	"fmt"
	"log"

	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func resourceTencentCloudKubernetesClusterAuditReadPreHandleResponse0(ctx context.Context, resp *tkev20180525.DescribeLogSwitchesResponseParams) error {
	logId := tccommon.GetLogId(ctx)
	d := tccommon.ResourceDataFromContext(ctx)
	if d == nil {
		return fmt.Errorf("resource data can not be nil")
	}

	if resp.SwitchSet == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster_audit` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	auditInfo := resp.SwitchSet[0].Audit
	if auditInfo.LogsetId != nil {
		_ = d.Set("logset_id", auditInfo.LogsetId)
	}

	if auditInfo.TopicId != nil {
		_ = d.Set("topic_id", auditInfo.TopicId)
	}

	if auditInfo.TopicRegion != nil {
		_ = d.Set("topic_region", auditInfo.TopicRegion)
	}

	return nil
}
