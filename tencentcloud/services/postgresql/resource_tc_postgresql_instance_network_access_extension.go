package postgresql

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
)

func networkAccessCustomResourceImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// TODO: implement me
	panic("TODO: implement me")
	return []*schema.ResourceData{d}, nil
}

func resourceTencentCloudPostgresqlInstanceNetworkAccessReadPreHandleResponse0(ctx context.Context, resp *postgresv20170312.DBInstance) error {
	// TODO: implement me
	panic("TODO: implement me")
}
