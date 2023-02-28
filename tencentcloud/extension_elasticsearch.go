package tencentcloud

const (
	ES_CHARGE_TYPE_PREPAID          = "PREPAID"
	ES_CHARGE_TYPE_POSTPAID_BY_HOUR = "POSTPAID_BY_HOUR"

	ES_DEPLOY_MODE_SINGLE_REGION = 0
	ES_DEPLOY_MODE_MULTI_REGION  = 1

	ES_LICENSE_TYPE_OSS      = "oss"
	ES_LICENSE_TYPE_BASIC    = "basic"
	ES_LICENSE_TYPE_PLATINUM = "platinum"

	ES_BASIC_SECURITY_TYPE_ON  = 2
	ES_BASIC_SECURITY_TYPE_OFF = 1

	ES_NODE_TYPE_HOT_DATA        = "hotData"
	ES_NODE_TYPE_WARM_DATA       = "warmData"
	ES_NODE_TYPE_DEDICATED_MATER = "dedicatedMaster"

	ES_RENEW_FLAG_AUTO   = "RENEW_FLAG_AUTO"
	ES_RENEW_FLAG_MANUAL = "RENEW_FLAG_MANUAL"

	ES_INSTANCE_STATUS_PROCESSING = 0
	ES_INSTANCE_STATUS_NORMAL     = 1
	ES_INSTANCE_STATUS_CREATING   = 2
	ES_INSTANCE_STATUS_STOP       = -1
	ES_INSTANCE_STATUS_DESTROYING = -2
	ES_INSTANCE_STATUS_DESTROYED  = -3
)

var ES_CHARGE_TYPE = []string{
	ES_CHARGE_TYPE_POSTPAID_BY_HOUR,
	ES_CHARGE_TYPE_PREPAID,
}

var ES_DEPLOY_MODE = []int{
	ES_DEPLOY_MODE_SINGLE_REGION,
	ES_DEPLOY_MODE_MULTI_REGION,
}

var ES_LICENSE_TYPE = []string{
	ES_LICENSE_TYPE_BASIC,
	ES_LICENSE_TYPE_PLATINUM,
}

var ES_BASIC_SECURITY_TYPE = []int{
	ES_BASIC_SECURITY_TYPE_ON,
	ES_BASIC_SECURITY_TYPE_OFF,
}

var ES_NODE_TYPE = []string{
	ES_NODE_TYPE_HOT_DATA,
	ES_NODE_TYPE_WARM_DATA,
	ES_NODE_TYPE_DEDICATED_MATER,
}

var ES_NODE_DISK_TYPE = []string{
	CVM_DISK_TYPE_CLOUD_SSD,
	CVM_DISK_TYPE_CLOUD_PREMIUM,
}

var ES_RENEW_FLAG = []string{
	ES_RENEW_FLAG_AUTO,
	ES_RENEW_FLAG_MANUAL,
}
