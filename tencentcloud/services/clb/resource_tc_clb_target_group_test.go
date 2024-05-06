package clb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"

	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_clb_target_group
	resource.AddTestSweepers("tencentcloud_clb_target_group", &resource.Sweeper{
		Name: "tencentcloud_clb_target_group",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := localclb.NewClbService(client)

			tgs, err := service.DescribeTargetGroups(ctx, "", nil)
			if err != nil {
				return err
			}

			// add scanning resources
			var resources, nonKeepResources []*tccommon.ResourceInstance
			for _, v := range tgs {
				if !tccommon.CheckResourcePersist(*v.TargetGroupName, *v.CreatedTime) {
					nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
						Id:   *v.TargetGroupId,
						Name: *v.TargetGroupName,
					})
				}
				resources = append(resources, &tccommon.ResourceInstance{
					Id:         *v.TargetGroupId,
					Name:       *v.TargetGroupName,
					CreateTime: *v.CreatedTime,
				})
			}
			tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateTargetGroup")

			for i := range tgs {
				tg := tgs[i]
				created := tccommon.ParseTimeFromCommonLayout(tg.CreatedTime)
				if tcacctest.IsResourcePersist(*tg.TargetGroupName, &created) {
					continue
				}
				log.Printf("%s will be remvoed", *tg.TargetGroupName)
				err = service.DeleteTarget(ctx, *tg.TargetGroupId)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudClbTargetGroup_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbTargetGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.test"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group.test", "target_group_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group.test", "vpc_id"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstanceTargetGroup(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstanceTargetGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.target_group"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_name", "tgt_grp_test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "port", "33"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.bind_ip", "10.0.0.4"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.port", "33"),
				),
			},
			{
				Config: testAccClbInstanceTargetGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbTargetGroupExists("tencentcloud_clb_target_group.target_group"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_name", "tgt_grp_test"),
					resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "port", "44"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.bind_ip", "10.0.0.4"),
					//resource.TestCheckResourceAttr("tencentcloud_clb_target_group.target_group", "target_group_instances.port", "44"),
				),
			},
		},
	})
}

func testAccCheckClbTargetGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_target_group" {
			continue
		}
		time.Sleep(5 * time.Second)
		filters := map[string]string{}
		targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, rs.Primary.ID, filters)
		if len(targetGroupInfos) > 0 && err == nil {
			return fmt.Errorf("[CHECK][CLB target group][Destroy] check: CLB target group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbTargetGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB target group][Exists] check: CLB target group %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB target group][Exists] check: CLB target group id is not set")
		}
		clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		filters := map[string]string{}
		targetGroupInfos, err := clbService.DescribeTargetGroups(ctx, rs.Primary.ID, filters)
		if err != nil {
			return err
		}
		if len(targetGroupInfos) == 0 {
			return fmt.Errorf("[CHECK][CLB target group][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClbTargetGroup_basic = `
resource "tencentcloud_clb_target_group" "test"{
    target_group_name = "qwe"
}
`

const testAccClbInstanceTargetGroup = `
resource "tencentcloud_clb_target_group" "target_group" {
    target_group_name = "tgt_grp_test"
    port              = 33
    vpc_id            = "vpc-4owdpnwr"
    target_group_instances {
      bind_ip = "172.16.16.95"
      port = 18800
    }
}
`

const testAccClbInstanceTargetGroupUpdate = `
resource "tencentcloud_clb_target_group" "target_group" {
    target_group_name = "tgt_grp_test"
    port              = 44
	vpc_id            = "vpc-4owdpnwr"
    target_group_instances {
      bind_ip = "172.16.16.95"
      port = 18800
    }
}
`
