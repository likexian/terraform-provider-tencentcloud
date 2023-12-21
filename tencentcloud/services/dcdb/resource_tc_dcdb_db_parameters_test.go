package dcdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdcdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dcdb"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudDCDBDbParametersResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDCDBDbParameters_basic, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCDBDbParametersExists("tencentcloud_dcdb_db_parameters.db_parameters"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_parameters.db_parameters", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_parameters.db_parameters", "params.#"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_parameters.db_parameters", "params.0.param", "max_connections"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_parameters.db_parameters", "params.0.value", "9999"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDCDBDbParameters_update, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDCDBDbParametersExists("tencentcloud_dcdb_db_parameters.db_parameters"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_parameters.db_parameters", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_parameters.db_parameters", "params.#"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_parameters.db_parameters", "params.0.param", "max_connections"),
					resource.TestCheckResourceAttr("tencentcloud_dcdb_db_parameters.db_parameters", "params.0.value", "10001"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_db_parameters.db_parameters",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDCDBDbParametersExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("dcdb db parameters  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("dcdb db parameters id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}

		instanceId := idSplit[0]

		dcdbService := svcdcdb.NewDcdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		ret, err := dcdbService.DescribeDcdbDbParametersById(ctx, instanceId)
		if err != nil {
			return err
		}

		if ret.InstanceId == nil {
			return fmt.Errorf("dcdb account privileges not found, instanceId: %v", rs.Primary.ID)
		}

		return nil
	}
}

const testAccDCDBDbParameters_basic = `

resource "tencentcloud_dcdb_db_parameters" "db_parameters" {
  instance_id = "%s"
  params {
	param = "max_connections"
	value = "9999"
  }
}

`

const testAccDCDBDbParameters_update = `

resource "tencentcloud_dcdb_db_parameters" "db_parameters" {
  instance_id = "%s"
  params {
	param = "max_connections"
	value = "10001"
  }
}

`
