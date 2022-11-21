/*
Provides a resource to create a organization policy_sub_account_attachment

Example Usage

```hcl
resource "tencentcloud_organization_policy_sub_account_attachment" "policy_sub_account_attachment" {
  member_uin               = 100028582828
  org_sub_account_uin      = 100028223737
  policy_id                = 144256499
}
```
Import

organization policy_sub_account_attachment can be imported using the id, e.g.
```
$ terraform import tencentcloud_organization_policy_sub_account_attachment.policy_sub_account_attachment policyId#memberUin#orgSubAccountUin
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOrganizationPolicySubAccountAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudOrganizationPolicySubAccountAttachmentRead,
		Create: resourceTencentCloudOrganizationPolicySubAccountAttachmentCreate,
		Delete: resourceTencentCloudOrganizationPolicySubAccountAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Policy ID.",
			},

			"org_sub_account_uin": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Organization administrator sub account uin list.",
			},

			"member_uin": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Organization member uin.",
			},

			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy name.",
			},

			"identity_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Manage Identity ID.",
			},

			"identity_role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity role name.",
			},

			"identity_role_alias_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identity role alias name.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},

			"org_sub_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Organization administrator sub account name.",
			},
		},
	}
}

func resourceTencentCloudOrganizationPolicySubAccountAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_policy_sub_account_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = organization.NewBindOrganizationMemberAuthAccountRequest()
		response         *organization.BindOrganizationMemberAuthAccountResponse
		policyId         int
		memberUin        int
		orgSubAccountUin int
	)

	if v, ok := d.GetOk("policy_id"); ok {
		policyId = v.(int)
		request.PolicyId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("org_sub_account_uin"); ok {
		orgSubAccountUin = v.(int)
		request.OrgSubAccountUins = []*int64{helper.IntInt64(v.(int))}
	}

	if v, ok := d.GetOk("member_uin"); ok {
		memberUin = v.(int)
		request.MemberUin = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().BindOrganizationMemberAuthAccount(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create organization policySubAccountAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strconv.Itoa(policyId) + FILED_SP + strconv.Itoa(memberUin) + FILED_SP + strconv.Itoa(orgSubAccountUin))
	return resourceTencentCloudOrganizationPolicySubAccountAttachmentRead(d, meta)
}

func resourceTencentCloudOrganizationPolicySubAccountAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_policy_sub_account_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	memberUin := idSplit[1]

	policySubAccountAttachment, err := service.DescribeOrganizationPolicySubAccountAttachment(ctx, policyId, memberUin)

	if err != nil {
		return err
	}

	if policySubAccountAttachment == nil {
		d.SetId("")
		return fmt.Errorf("resource `policySubAccountAttachment` %s does not exist", d.Id())
	}

	if policySubAccountAttachment.PolicyId != nil {
		_ = d.Set("policy_id", policySubAccountAttachment.PolicyId)
	}

	if policySubAccountAttachment.OrgSubAccountUin != nil {
		_ = d.Set("org_sub_account_uin", policySubAccountAttachment.OrgSubAccountUin)
	}

	_ = d.Set("member_uin", helper.StrToInt64(memberUin))

	if policySubAccountAttachment.PolicyName != nil {
		_ = d.Set("policy_name", policySubAccountAttachment.PolicyName)
	}

	if policySubAccountAttachment.IdentityId != nil {
		_ = d.Set("identity_id", policySubAccountAttachment.IdentityId)
	}

	if policySubAccountAttachment.IdentityRoleName != nil {
		_ = d.Set("identity_role_name", policySubAccountAttachment.IdentityRoleName)
	}

	if policySubAccountAttachment.IdentityRoleAliasName != nil {
		_ = d.Set("identity_role_alias_name", policySubAccountAttachment.IdentityRoleAliasName)
	}

	if policySubAccountAttachment.CreateTime != nil {
		_ = d.Set("create_time", policySubAccountAttachment.CreateTime)
	}

	if policySubAccountAttachment.UpdateTime != nil {
		_ = d.Set("update_time", policySubAccountAttachment.UpdateTime)
	}

	if policySubAccountAttachment.OrgSubAccountName != nil {
		_ = d.Set("org_sub_account_name", policySubAccountAttachment.OrgSubAccountName)
	}

	return nil
}

func resourceTencentCloudOrganizationPolicySubAccountAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_policy_sub_account_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	memberUin := idSplit[1]
	orgSubAccountUin := idSplit[2]

	if err := service.DeleteOrganizationPolicySubAccountAttachmentById(ctx, policyId, memberUin, orgSubAccountUin); err != nil {
		return err
	}

	return nil
}
