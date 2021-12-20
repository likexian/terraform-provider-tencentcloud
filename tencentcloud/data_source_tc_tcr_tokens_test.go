package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTCRTokensNameAll = "data.tencentcloud_tcr_tokens.id_test"

func TestAccTencentCloudDataTCRTokens(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRTokensBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRTokenExists("tencentcloud_tcr_token.mytcr_token"),
					resource.TestCheckResourceAttrSet(testDataTCRTokensNameAll, "token_list.0.token_id"),
					resource.TestCheckResourceAttrSet(testDataTCRTokensNameAll, "token_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataTCRTokensNameAll, "token_list.0.description"),
					resource.TestCheckResourceAttr(testDataTCRTokensNameAll, "token_list.0.enable", "true"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRTokensBasic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "standard"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_token" "mytcr_token" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  description       = "test"
  enable   = true
}

data "tencentcloud_tcr_tokens" "id_test" {
  instance_id = tencentcloud_tcr_token.mytcr_token.instance_id
}
`
