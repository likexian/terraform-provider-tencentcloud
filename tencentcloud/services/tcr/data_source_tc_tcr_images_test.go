package tcr_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testObjectName = "data.tencentcloud_tcr_images.images"

func TestAccTencentCloudTcrImagesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrImagesDataSource_id, tcacctest.DefaultTCRInstanceId, tcacctest.DefaultTCRNamespace, tcacctest.DefaultTCRRepoName),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttr(testObjectName, "registry_id", tcacctest.DefaultTCRInstanceId),
					resource.TestCheckResourceAttr(testObjectName, "namespace_name", tcacctest.DefaultTCRNamespace),
					resource.TestCheckResourceAttr(testObjectName, "repository_name", tcacctest.DefaultTCRRepoName),
					resource.TestCheckResourceAttrSet(testObjectName, "image_info_list.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudTcrImagesDataSource_exact(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrImagesDataSource_exact, tcacctest.DefaultTCRInstanceId, tcacctest.DefaultTCRNamespace, tcacctest.DefaultTCRRepoName),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttr(testObjectName, "registry_id", tcacctest.DefaultTCRInstanceId),
					resource.TestCheckResourceAttr(testObjectName, "namespace_name", tcacctest.DefaultTCRNamespace),
					resource.TestCheckResourceAttr(testObjectName, "repository_name", tcacctest.DefaultTCRRepoName),
					resource.TestCheckResourceAttr(testObjectName, "exact_match", "true"),
					resource.TestCheckResourceAttrSet(testObjectName, "image_info_list.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudTcrImagesDataSource_exact_version(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrImagesDataSource_exact_version, tcacctest.DefaultTCRInstanceId, tcacctest.DefaultTCRNamespace, tcacctest.DefaultTCRRepoName),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testObjectName, "id"),
					resource.TestCheckResourceAttr(testObjectName, "registry_id", tcacctest.DefaultTCRInstanceId),
					resource.TestCheckResourceAttr(testObjectName, "namespace_name", tcacctest.DefaultTCRNamespace),
					resource.TestCheckResourceAttr(testObjectName, "repository_name", tcacctest.DefaultTCRRepoName),
					resource.TestCheckResourceAttr(testObjectName, "image_version", "v1"),
					resource.TestCheckResourceAttr(testObjectName, "exact_match", "true"),
					resource.TestCheckResourceAttrSet(testObjectName, "image_info_list.#"),
				),
			},
		},
	})
}

const testAccTcrImagesDataSource_id = `

data "tencentcloud_tcr_images" "images" {
  registry_id = "%s"
  namespace_name = "%s" 
  repository_name = "%s"
  }

`

const testAccTcrImagesDataSource_exact = `

data "tencentcloud_tcr_images" "images" {
  registry_id = "%s"
  namespace_name = "%s" 
  repository_name = "%s"
  exact_match = true
  }

`

const testAccTcrImagesDataSource_exact_version = `

data "tencentcloud_tcr_images" "images" {
  registry_id = "%s"
  namespace_name = "%s" 
  repository_name = "%s"
  image_version = "v1"
  exact_match = true
  }

`
