package vpc_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_security_group", &resource.Sweeper{
		Name: "tencentcloud_security_group",
		F:    testSweepSecurityGroups,
	})
}

func testSweepSecurityGroups(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcvpc.NewVpcService(client)

	sgs, err := service.DescribeSecurityGroups(ctx, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("DescribeSecurityGroups error: %s", err.Error())
	}

	// add scanning resources
	var resources, nonKeepResources []*tccommon.ResourceInstance
	for _, v := range sgs {
		if !tccommon.CheckResourcePersist(*v.SecurityGroupName, *v.CreatedTime) {
			nonKeepResources = append(nonKeepResources, &tccommon.ResourceInstance{
				Id:   *v.SecurityGroupId,
				Name: *v.SecurityGroupName,
			})
		}
		resources = append(resources, &tccommon.ResourceInstance{
			Id:        *v.SecurityGroupId,
			Name:      *v.SecurityGroupName,
			CreatTime: *v.CreatedTime,
		})
	}
	tccommon.ProcessScanCloudResources(client, resources, nonKeepResources, "CreateSecurityGroup")

	for _, v := range sgs {
		name := *v.SecurityGroupName
		id := *v.SecurityGroupId

		now := time.Now()
		createTime := tccommon.StringToTime(*v.CreatedTime)
		interval := now.Sub(createTime).Minutes()

		if strings.HasPrefix(name, tcacctest.KeepResource) || strings.HasPrefix(name, tcacctest.DefaultResource) {
			continue
		}

		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err := service.DeleteSecurityGroup(ctx, id); err != nil {
			log.Printf("[ERROR] sweep security group %s error: %s", id, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudSecurityGroup_basic(t *testing.T) {
	t.Parallel()
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("tencentcloud_security_group.foo", &sgId),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "name", "ci-temp-test-sg"),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "project_id", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_security_group.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudSecurityGroup_update(t *testing.T) {
	t.Parallel()
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("tencentcloud_security_group.foo", &sgId),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "name", "ci-temp-test-sg"),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "description", "ci-temp-test-sg-desc"),
				),
			},
			{
				Config: testAccSecurityGroupConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("tencentcloud_security_group.foo", &sgId),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "name", "ci-temp-test-sg-updated"),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "description", "ci-temp-test-sg-desc-updated"),
				),
			},
		},
	})
}

func TestAccTencentCloudSecurityGroup_tags(t *testing.T) {
	t.Parallel()
	var sgId string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy(&sgId),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfigTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("tencentcloud_security_group.foo", &sgId),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "description", ""),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "tags.test", "test"),
				),
			},
			{
				Config: testAccSecurityGroupConfigTagsUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("tencentcloud_security_group.foo", &sgId),
					resource.TestCheckNoResourceAttr("tencentcloud_security_group.foo", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_security_group.foo", "tags.abc", "abc"),
				),
			},
		},
	})
}

func testAccCheckSecurityGroupDestroy(id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		service := svcvpc.NewVpcService(client)

		sg, err := service.DescribeSecurityGroup(context.TODO(), *id)
		if err != nil {
			return err
		}

		if sg != nil {
			return fmt.Errorf("security group still exists")
		}

		return nil
	}
}

func testAccCheckSecurityGroupExists(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no security group ID is set")
		}

		service := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		sg, err := service.DescribeSecurityGroup(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if sg == nil {
			return fmt.Errorf("security group not found: %s", rs.Primary.ID)
		}

		*id = rs.Primary.ID

		return nil
	}
}

const testAccSecurityGroupConfigBasic = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}
`

const testAccSecurityGroupConfig = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg"
  description = "ci-temp-test-sg-desc"
}
`
const testAccSecurityGroupConfigUpdate = `
resource "tencentcloud_security_group" "foo" {
  name        = "ci-temp-test-sg-updated"
  description = "ci-temp-test-sg-desc-updated"
}
`

const testAccSecurityGroupConfigTags = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"

  tags = {
    "test" = "test"
  }
}
`

const testAccSecurityGroupConfigTagsUpdate = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"

  tags = {
    "abc" = "abc"
  }
}
`
