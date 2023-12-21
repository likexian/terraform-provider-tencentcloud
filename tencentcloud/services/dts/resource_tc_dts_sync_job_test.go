package dts_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdts "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dts"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dts_sync_job", &resource.Sweeper{
		Name: "tencentcloud_dts_sync_job",
		F:    testSweepDtsSyncJob,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dts_sync_job
func testSweepDtsSyncJob(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	dtsService := svcdts.NewDtsService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())
	param := map[string]interface{}{}

	ret, err := dtsService.DescribeDtsSyncJobsByFilter(ctx, param)
	if err != nil {
		return err
	}

	for _, v := range ret {
		delId := *v.JobId

		if strings.HasPrefix(*v.JobName, tcacctest.KeepResource) || strings.HasPrefix(*v.JobName, tcacctest.DefaultResource) {
			continue
		}

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := dtsService.DeleteDtsSyncJobById(ctx, delId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("[ERROR] sweeper tencentcloud_dts_sync_job:[%v] failed! reason:[%s]", delId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudDtsSyncJobResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDtsSyncJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsSyncJob,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDtsSyncJobExists("tencentcloud_dts_sync_job.sync_job"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job.sync_job", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_job.sync_job", "pay_mode", "PostPay"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_job.sync_job", "src_database_type", "mysql"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_job.sync_job", "src_region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_job.sync_job", "dst_database_type", "cynosdbmysql"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_job.sync_job", "dst_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_sync_job.sync_job", "tags.#"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_job.sync_job", "auto_renew", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dts_sync_job.sync_job", "instance_class", "micro"),
				),
			},
		},
	})
}

func testAccCheckDtsSyncJobDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dtsService := svcdts.NewDtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dts_sync_job" {
			continue
		}

		job, err := dtsService.DescribeDtsSyncJob(ctx, helper.String(rs.Primary.ID))
		if err != nil {
			return err
		}

		if job != nil {
			status := *job.Status
			if status != "UnInitialized" {
				return fmt.Errorf("DTS sync job still exist, Id: %v, status:%s", rs.Primary.ID, status)
			}
		}
	}
	return nil
}

func testAccCheckDtsSyncJobExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		dtsService := svcdts.NewDtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("DTS sync job %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DTS sync job id is not set")
		}

		job, err := dtsService.DescribeDtsSyncJob(ctx, helper.String(rs.Primary.ID))
		if err != nil {
			return err
		}

		if job == nil {
			return fmt.Errorf("DTS sync job not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccDtsSyncJob = `

resource "tencentcloud_dts_sync_job" "sync_job" {
  pay_mode = "PostPay"
  src_database_type = "mysql"
  src_region = "ap-guangzhou"
  dst_database_type = "cynosdbmysql"
  dst_region = "ap-guangzhou"
  tags {
	tag_key = "aaa"
	tag_value = "bbb"
  }
  auto_renew = 0
  instance_class = "micro"
}

`
