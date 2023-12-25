package tcr_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTCRVPCAttachmentsNameAll = "data.tencentcloud_tcr_vpc_attachments.id_test"

func TestAccTencentCloudTcrVPCAttachmentsData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRVPCAttachmentsBasic,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRVPCAttachmentExists("tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment"),
					resource.TestCheckResourceAttr(testDataTCRVPCAttachmentsNameAll, "vpc_attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataTCRVPCAttachmentsNameAll, "vpc_attachment_list.0.status"),
				),
			},
		},
	})
}

const DefaultTcrVpcSubnets = `

data "tencentcloud_vpc_subnets" "sh" {
  availability_zone = "ap-shanghai-1"
}

locals {
  vpc_id = data.tencentcloud_vpc_subnets.sh.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.sh.instance_list.0.subnet_id
}`

const testAccTencentCloudDataTCRVPCAttachmentsBasic = DefaultTcrVpcSubnets + `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "test-tcr-attach"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_vpc_attachment" "mytcr_vpc_attachment" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
}

data "tencentcloud_tcr_vpc_attachments" "id_test" {
  instance_id = tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment.instance_id
}
`
