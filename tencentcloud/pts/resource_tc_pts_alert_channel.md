Provides a resource to create a pts alert_channel

~> **NOTE:** Modification is not currently supported, please go to the console to modify.

Example Usage

```hcl
resource "tencentcloud_monitor_alarm_notice" "example" {
  name                  = "test_alarm_notice_1"
  notice_type           = "ALL"
  notice_language       = "zh-CN"

  user_notices    {
      receiver_type              = "USER"
      start_time                 = 0
      end_time                   = 1
      notice_way                 = ["EMAIL", "SMS", "WECHAT"]
      user_ids                   = [10001]
      group_ids                  = []
      phone_order                = [10001]
      phone_circle_times         = 2
      phone_circle_interval      = 50
      phone_inner_interval       = 60
      need_phone_arrive_notice   = 1
      phone_call_type            = "CIRCLE"
      weekday                    =[1,2,3,4,5,6,7]
  }

  url_notices {
      url    = "https://www.mytest.com/validate"
      end_time =  0
      start_time = 1
      weekday = [1,2,3,4,5,6,7]
  }

}

resource "tencentcloud_pts_project" "project" {
  name = "ptsObjectName"
  description = "desc"
  tags {
    tag_key = "createdBy"
    tag_value = "terraform"
  }
}

resource "tencentcloud_pts_alert_channel" "alert_channel" {
  notice_id = tencentcloud_monitor_alarm_notice.example.id
  project_id = tencentcloud_pts_project.project.id
  amp_consumer_id = "Consumer-vvy1xxxxxx"
}

```
Import

pts alert_channel can be imported using the project_id#notice_id, e.g.
```
$ terraform import tencentcloud_pts_alert_channel.alert_channel project-kww5v8se#notice-kl66t6y9
```