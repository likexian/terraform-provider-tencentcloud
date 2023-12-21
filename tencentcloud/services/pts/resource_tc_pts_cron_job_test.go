package pts_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcpts "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/pts"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudPtsCronJobResource_basic -v
func TestAccTencentCloudPtsCronJobResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckPtsCronJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsCronJob,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPtsCronJobExists("tencentcloud_pts_cron_job.cron_job"),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "name", "iac-cron_job-update"),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "project_id", tcacctest.DefaultPtsProjectId),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "scenario_id", tcacctest.DefaultScenarioId),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "scenario_name", "keep-pts-js"),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "frequency_type", "2"),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "cron_expression", "* 1 * * *"),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "job_owner", "userName"),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "notice_id", tcacctest.DefaultPtsNoticeId),
					resource.TestCheckResourceAttr("tencentcloud_pts_cron_job.cron_job", "note", "desc"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_cron_job.cron_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPtsCronJobDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcpts.NewPtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_pts_project" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectId := idSplit[0]
		cronJobId := idSplit[1]

		cronJob, err := service.DescribePtsCronJob(ctx, cronJobId, projectId)
		if cronJob != nil {
			return fmt.Errorf("pts cronJob %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckPtsCronJobExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectId := idSplit[0]
		cronJobId := idSplit[1]

		service := svcpts.NewPtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		cronJob, err := service.DescribePtsCronJob(ctx, cronJobId, projectId)
		if cronJob == nil {
			return fmt.Errorf("pts cronJob %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccPtsCronJobVar = `
variable "project_id" {
  default = "` + tcacctest.DefaultPtsProjectId + `"
}
variable "scenario_id" {
	default = "` + tcacctest.DefaultScenarioId + `"
}
variable "notice_id" {
	default = "` + tcacctest.DefaultPtsNoticeId + `"
}
  
`

const testAccPtsCronJob = testAccPtsCronJobVar + `

resource "tencentcloud_pts_cron_job" "cron_job" {
	name = "iac-cron_job-update"
	project_id = var.project_id
	scenario_id = var.scenario_id
	scenario_name = "keep-pts-js"
	frequency_type = 2
	cron_expression = "* 1 * * *"
	job_owner = "userName"
	# end_time = ""
	notice_id = var.notice_id
	note = "desc"
  }

`
