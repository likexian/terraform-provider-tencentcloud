The TencentCloud provider is used to interact with many resources supported by [TencentCloud](https://intl.cloud.tencent.com).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** From version 1.9.0 (June 18, 2019), the provider start to support Terraform 0.12.x.

Example Usage

```hcl
terraform {
  required_providers {
    tencentcloud = {
      source = "tencentcloudstack/tencentcloud"
    }
  }
}

# Configure the TencentCloud Provider
provider "tencentcloud" {
  secret_id  = var.secret_id
  secret_key = var.secret_key
  region     = var.region
}

#Configure the TencentCloud Provider with STS
provider "tencentcloud" {
  secret_id  = var.secret_id
  secret_key = var.secret_key
  region     = var.region
  assume_role {
    role_arn         = var.assume_role_arn
    session_name     = var.session_name
    session_duration = var.session_duration
    policy           = var.policy
  }
}
```

Resources List

Provider Data Sources
  tencentcloud_availability_regions
  tencentcloud_availability_zones_by_product
  tencentcloud_availability_zones

Project
  Data Source
    tencentcloud_projects

  Resource
    tencentcloud_project

Anti-DDoS(antiddos)
  Data Source
    tencentcloud_antiddos_basic_device_status
    tencentcloud_antiddos_bgp_biz_trend
    tencentcloud_antiddos_list_listener
    tencentcloud_antiddos_overview_attack_trend

Anti-DDoS(DayuV2)
  Data Source
    tencentcloud_dayu_eip
    tencentcloud_dayu_l4_rules_v2
    tencentcloud_dayu_l7_rules_v2
    tencentcloud_antiddos_pending_risk_info
    tencentcloud_antiddos_overview_index
    tencentcloud_antiddos_overview_ddos_trend
    tencentcloud_antiddos_overview_ddos_event_list
    tencentcloud_antiddos_overview_cc_trend

  Resource
    tencentcloud_dayu_eip
    tencentcloud_dayu_l4_rule
    tencentcloud_dayu_l7_rule_v2
    tencentcloud_dayu_ddos_policy_v2
    tencentcloud_dayu_cc_policy_v2
    tencentcloud_dayu_ddos_ip_attachment_v2
    tencentcloud_antiddos_ddos_black_white_ip
    tencentcloud_antiddos_ddos_geo_ip_block_config
    tencentcloud_antiddos_ddos_speed_limit_config
    tencentcloud_antiddos_default_alarm_threshold
    tencentcloud_antiddos_scheduling_domain_user_name
    tencentcloud_antiddos_ip_alarm_threshold_config
    tencentcloud_antiddos_packet_filter_config
    tencentcloud_antiddos_port_acl_config
    tencentcloud_antiddos_cc_black_white_ip
    tencentcloud_antiddos_cc_precision_policy

Anti-DDoS(Dayu)
  Data Source
    tencentcloud_dayu_cc_http_policies
    tencentcloud_dayu_cc_https_policies
    tencentcloud_dayu_ddos_policies
    tencentcloud_dayu_ddos_policy_attachments
    tencentcloud_dayu_ddos_policy_cases
    tencentcloud_dayu_l4_rules
    tencentcloud_dayu_l7_rules

  Resource
    tencentcloud_dayu_cc_http_policy
    tencentcloud_dayu_cc_https_policy
    tencentcloud_dayu_ddos_policy
    tencentcloud_dayu_ddos_policy_attachment
    tencentcloud_dayu_ddos_policy_case
    tencentcloud_dayu_l4_rule
    tencentcloud_dayu_l7_rule

API GateWay(apigateway)
  Data Source
    tencentcloud_api_gateway_apis
    tencentcloud_api_gateway_services
    tencentcloud_api_gateway_throttling_services
    tencentcloud_api_gateway_throttling_apis
    tencentcloud_api_gateway_usage_plans
    tencentcloud_api_gateway_ip_strategies
    tencentcloud_api_gateway_customer_domains
    tencentcloud_api_gateway_usage_plan_environments
    tencentcloud_api_gateway_api_keys
    tencentcloud_api_gateway_api_docs
    tencentcloud_api_gateway_api_apps
    tencentcloud_api_gateway_plugins
    tencentcloud_api_gateway_upstreams
    tencentcloud_api_gateway_api_usage_plans
    tencentcloud_api_gateway_api_app_service
    tencentcloud_api_gateway_bind_api_apps_status
    tencentcloud_api_gateway_api_app_api
    tencentcloud_api_gateway_api_plugins
    tencentcloud_api_gateway_service_release_versions
    tencentcloud_api_gateway_service_environment_list

  Resource
    tencentcloud_api_gateway_api
    tencentcloud_api_gateway_service
    tencentcloud_api_gateway_custom_domain
    tencentcloud_api_gateway_usage_plan
    tencentcloud_api_gateway_usage_plan_attachment
    tencentcloud_api_gateway_ip_strategy
    tencentcloud_api_gateway_strategy_attachment
    tencentcloud_api_gateway_api_key
    tencentcloud_api_gateway_api_key_attachment
    tencentcloud_api_gateway_service_release
    tencentcloud_api_gateway_plugin
    tencentcloud_api_gateway_plugin_attachment
    tencentcloud_api_gateway_api_doc
    tencentcloud_api_gateway_api_app
    tencentcloud_api_gateway_upstream
    tencentcloud_api_gateway_api_app_attachment
    tencentcloud_api_gateway_update_api_app_key
    tencentcloud_api_gateway_import_open_api 
    tencentcloud_api_gateway_update_service

Cloud Audit(Audit)
  Data Source
    tencentcloud_audit_cos_regions
    tencentcloud_audit_key_alias
    tencentcloud_audits
    tencentcloud_audit_events

  Resource
    tencentcloud_audit_track
    tencentcloud_events_audit_track

Auto Scaling(AS)
  Data Source
    tencentcloud_as_scaling_configs
    tencentcloud_as_scaling_groups
    tencentcloud_as_scaling_policies
    tencentcloud_as_instances
    tencentcloud_as_advices
    tencentcloud_as_limits
    tencentcloud_as_last_activity

  Resource
    tencentcloud_as_scaling_config
    tencentcloud_as_scaling_group
    tencentcloud_as_scaling_group_status
    tencentcloud_as_attachment
    tencentcloud_as_scaling_policy
    tencentcloud_as_schedule
    tencentcloud_as_lifecycle_hook
    tencentcloud_as_notification
    tencentcloud_as_remove_instances
    tencentcloud_as_protect_instances
    tencentcloud_as_start_instances
    tencentcloud_as_stop_instances
    tencentcloud_as_scale_in_instances
    tencentcloud_as_scale_out_instances
    tencentcloud_as_execute_scaling_policy
    tencentcloud_as_complete_lifecycle
    tencentcloud_as_start_instance_refresh

Content Delivery Network(CDN)
  Data Source
    tencentcloud_cdn_domains
    tencentcloud_cdn_domain_verifier

  Resource
    tencentcloud_cdn_domain
    tencentcloud_cdn_url_push
    tencentcloud_cdn_url_purge

Cloud Kafka(ckafka)
  Data Source
    tencentcloud_ckafka_users
    tencentcloud_ckafka_acls
    tencentcloud_ckafka_topics
    tencentcloud_ckafka_instances
    tencentcloud_ckafka_connect_resource
    tencentcloud_ckafka_region
    tencentcloud_ckafka_datahub_topic
    tencentcloud_ckafka_datahub_group_offsets
    tencentcloud_ckafka_datahub_task
    tencentcloud_ckafka_group
    tencentcloud_ckafka_group_offsets
    tencentcloud_ckafka_group_info
    tencentcloud_ckafka_task_status
    tencentcloud_ckafka_topic_flow_ranking
    tencentcloud_ckafka_topic_produce_connection
    tencentcloud_ckafka_topic_subscribe_group
    tencentcloud_ckafka_topic_sync_replica
    tencentcloud_ckafka_zone

  Resource
    tencentcloud_ckafka_instance
    tencentcloud_ckafka_user
    tencentcloud_ckafka_acl
    tencentcloud_ckafka_topic
    tencentcloud_ckafka_datahub_topic
    tencentcloud_ckafka_connect_resource
    tencentcloud_ckafka_renew_instance
    tencentcloud_ckafka_acl_rule
    tencentcloud_ckafka_consumer_group
    tencentcloud_ckafka_consumer_group_modify_offset
    tencentcloud_ckafka_datahub_task
    tencentcloud_ckafka_route

Cloud Access Management(CAM)
  Data Source
    tencentcloud_cam_group_memberships
    tencentcloud_cam_group_policy_attachments
    tencentcloud_cam_groups
    tencentcloud_cam_policies
    tencentcloud_cam_role_policy_attachments
    tencentcloud_cam_roles
    tencentcloud_cam_saml_providers
    tencentcloud_cam_user_policy_attachments
    tencentcloud_cam_users
    tencentcloud_user_info
    tencentcloud_cam_list_entities_for_policy
    tencentcloud_cam_secret_last_used_time
    tencentcloud_cam_account_summary
    tencentcloud_cam_policy_granting_service_access
    tencentcloud_cam_oidc_config
    tencentcloud_cam_group_user_account
    tencentcloud_cam_sub_accounts
    tencentcloud_cam_role_detail

  Resource
    tencentcloud_cam_role
    tencentcloud_cam_role_by_name
    tencentcloud_cam_role_policy_attachment
    tencentcloud_cam_role_policy_attachment_by_name
    tencentcloud_cam_policy
    tencentcloud_cam_policy_by_name
    tencentcloud_cam_user
    tencentcloud_cam_user_policy_attachment
    tencentcloud_cam_group
    tencentcloud_cam_group_policy_attachment
    tencentcloud_cam_group_membership
    tencentcloud_cam_saml_provider
    tencentcloud_cam_oidc_sso
    tencentcloud_cam_role_sso
    tencentcloud_cam_service_linked_role
    tencentcloud_cam_mfa_flag
    tencentcloud_cam_access_key
    tencentcloud_cam_user_saml_config
    tencentcloud_cam_tag_role_attachment
    tencentcloud_cam_policy_version
    tencentcloud_cam_set_policy_version_config
    tencentcloud_cam_user_permission_boundary_attachment
    tencentcloud_cam_role_permission_boundary_attachment

Customer Identity and Access Management(CIAM)
  Resource
    tencentcloud_ciam_user_store
    tencentcloud_ciam_user_group

Cloud Block Storage(CBS)
  Data Source
    tencentcloud_cbs_snapshots
    tencentcloud_cbs_storages
    tencentcloud_cbs_storages_set
    tencentcloud_cbs_snapshot_policies

  Resource
    tencentcloud_cbs_storage
    tencentcloud_cbs_storage_set
    tencentcloud_cbs_storage_attachment
    tencentcloud_cbs_storage_set_attachment
    tencentcloud_cbs_snapshot
    tencentcloud_cbs_snapshot_policy
    tencentcloud_cbs_snapshot_policy_attachment
    tencentcloud_cbs_snapshot_share_permission
    tencentcloud_cbs_disk_backup
    tencentcloud_cbs_disk_backup_rollback_operation

Cloud Connect Network(CCN)
  Data Source
    tencentcloud_ccn_bandwidth_limits
    tencentcloud_ccn_instances
    tencentcloud_ccn_cross_border_compliance
    tencentcloud_ccn_tenant_instances
    tencentcloud_ccn_cross_border_flow_monitor
    tencentcloud_ccn_cross_border_region_bandwidth_limits
    tencentcloud_ccn_routes
    tencentcloud_ccn_route_table_input_policies

  Resource
    tencentcloud_ccn
    tencentcloud_ccn_attachment
    tencentcloud_ccn_bandwidth_limit
    tencentcloud_ccn_routes
    tencentcloud_ccn_instances_accept_attach
    tencentcloud_ccn_instances_reject_attach
    tencentcloud_ccn_instances_reset_attach
    tencentcloud_ccn_route_table
    tencentcloud_ccn_route_table_input_policies
    tencentcloud_ccn_route_table_broadcast_policies
    tencentcloud_ccn_route_table_selection_policies
    tencentcloud_ccn_route_table_associate_instance_config

CVM Dedicated Host(CDH)
  Data Source
    tencentcloud_cdh_instances

  Resource
    tencentcloud_cdh_instance

Cloud File Storage(CFS)
  Data Source
    tencentcloud_cfs_access_groups
    tencentcloud_cfs_access_rules
    tencentcloud_cfs_file_systems
    tencentcloud_cfs_mount_targets
    tencentcloud_cfs_file_system_clients
    tencentcloud_cfs_available_zone

  Resource
    tencentcloud_cfs_file_system
    tencentcloud_cfs_access_group
    tencentcloud_cfs_access_rule
    tencentcloud_cfs_auto_snapshot_policy
    tencentcloud_cfs_auto_snapshot_policy_attachment
    tencentcloud_cfs_snapshot
    tencentcloud_cfs_sign_up_cfs_service

Container Cluster(tke)
  Data Source
    tencentcloud_container_cluster_instances
    tencentcloud_container_clusters

  Resource
    tencentcloud_container_cluster
    tencentcloud_container_cluster_instance

Cloud Load Balancer(CLB)
  Data Source
    tencentcloud_clb_attachments
    tencentcloud_clb_instances
    tencentcloud_clb_listener_rules
    tencentcloud_clb_listeners
    tencentcloud_clb_redirections
    tencentcloud_clb_target_groups
    tencentcloud_clb_cluster_resources
    tencentcloud_clb_cross_targets
    tencentcloud_clb_exclusive_clusters
    tencentcloud_clb_idle_instances
    tencentcloud_clb_listeners_by_targets
    tencentcloud_clb_instance_by_cert_id
    tencentcloud_clb_instance_traffic
    tencentcloud_clb_instance_detail
    tencentcloud_clb_resources
    tencentcloud_clb_target_group_list
    tencentcloud_clb_target_health

  Resource
    tencentcloud_clb_instance
    tencentcloud_clb_listener
    tencentcloud_clb_listener_rule
    tencentcloud_clb_listener_default_domain
    tencentcloud_clb_attachment
    tencentcloud_clb_redirection
    tencentcloud_lb
    tencentcloud_alb_server_attachment
    tencentcloud_clb_target_group
    tencentcloud_clb_target_group_instance_attachment
    tencentcloud_clb_target_group_attachment
    tencentcloud_clb_target_group_attachments
    tencentcloud_clb_log_set
    tencentcloud_clb_log_topic
    tencentcloud_clb_customized_config
    tencentcloud_clb_snat_ip
    tencentcloud_clb_function_targets_attachment
    tencentcloud_clb_instance_sla_config
    tencentcloud_clb_instance_mix_ip_target_config
    tencentcloud_clb_replace_cert_for_lbs
    tencentcloud_clb_security_group_attachment

Cloud Object Storage(COS)
  Data Source
    tencentcloud_cos_bucket_object
    tencentcloud_cos_buckets
    tencentcloud_cos_batchs
    tencentcloud_cos_bucket_inventorys
    tencentcloud_cos_bucket_multipart_uploads

  Resource
    tencentcloud_cos_bucket
    tencentcloud_cos_bucket_object
    tencentcloud_cos_bucket_policy
    tencentcloud_cos_bucket_referer
    tencentcloud_cos_bucket_version
    tencentcloud_cos_bucket_domain_certificate_attachment
    tencentcloud_cos_bucket_inventory
    tencentcloud_cos_batch
    tencentcloud_cos_object_abort_multipart_upload_operation
    tencentcloud_cos_object_copy_operation
    tencentcloud_cos_object_restore_operation
    tencentcloud_cos_bucket_generate_inventory_immediately_operation
    tencentcloud_cos_object_download_operation

Cloud Virtual Machine(CVM)
  Data Source
    tencentcloud_image
    tencentcloud_image_from_family
    tencentcloud_images
    tencentcloud_instance_types
    tencentcloud_instances
    tencentcloud_instances_set
    tencentcloud_key_pairs
    tencentcloud_eip
    tencentcloud_eips
    tencentcloud_eip_address_quota
    tencentcloud_eip_network_account_type
    tencentcloud_placement_groups
    tencentcloud_reserved_instance_configs
    tencentcloud_reserved_instances
    tencentcloud_cvm_instances_modification
    tencentcloud_cvm_instance_vnc_url
    tencentcloud_cvm_disaster_recover_group_quota
    tencentcloud_cvm_chc_hosts
    tencentcloud_cvm_chc_denied_actions
    tencentcloud_cvm_image_quota
    tencentcloud_cvm_image_share_permission
    tencentcloud_cvm_import_image_os

  Resource
    tencentcloud_instance
    tencentcloud_instance_set
    tencentcloud_eip
    tencentcloud_eip_association
    tencentcloud_eip_address_transform
    tencentcloud_eip_public_address_adjust
    tencentcloud_eip_normal_address_return
    tencentcloud_key_pair
    tencentcloud_placement_group
    tencentcloud_reserved_instance
    tencentcloud_image
    tencentcloud_cvm_hpc_cluster
    tencentcloud_cvm_launch_template
    tencentcloud_cvm_launch_template_version
    tencentcloud_cvm_launch_template_default_version
    tencentcloud_cvm_security_group_attachment
    tencentcloud_cvm_reboot_instance
    tencentcloud_cvm_chc_config
    tencentcloud_cvm_renew_instance
    tencentcloud_cvm_sync_image
    tencentcloud_cvm_export_images
    tencentcloud_cvm_image_share_permission
    tencentcloud_cvm_action_timer

TDSQL-C MySQL(CynosDB)
  Data Source
    tencentcloud_cynosdb_clusters
    tencentcloud_cynosdb_instances
    tencentcloud_cynosdb_zone_config
    tencentcloud_cynosdb_accounts
    tencentcloud_cynosdb_cluster_instance_groups
    tencentcloud_cynosdb_cluster_params
    tencentcloud_cynosdb_param_templates
    tencentcloud_cynosdb_audit_logs
    tencentcloud_cynosdb_binlog_download_url
    tencentcloud_cynosdb_cluster_detail_databases
    tencentcloud_cynosdb_cluster_param_logs
    tencentcloud_cynosdb_cluster
    tencentcloud_cynosdb_describe_instance_slow_queries
    tencentcloud_cynosdb_describe_instance_error_logs
    tencentcloud_cynosdb_account_all_grant_privileges
    tencentcloud_cynosdb_resource_package_list
    tencentcloud_cynosdb_project_security_groups
    tencentcloud_cynosdb_resource_package_sale_specs
    tencentcloud_cynosdb_rollback_time_range
    tencentcloud_cynosdb_zone
    tencentcloud_cynosdb_instance_slow_queries
    tencentcloud_cynosdb_proxy_node
    tencentcloud_cynosdb_proxy_version

  Resource
    tencentcloud_cynosdb_cluster_resource_packages_attachment
    tencentcloud_cynosdb_cluster
    tencentcloud_cynosdb_readonly_instance
    tencentcloud_cynosdb_security_group
    tencentcloud_cynosdb_audit_log_file
    tencentcloud_cynosdb_cluster_password_complexity
    tencentcloud_cynosdb_export_instance_error_logs
    tencentcloud_cynosdb_export_instance_slow_queries
    tencentcloud_cynosdb_account_privileges
    tencentcloud_cynosdb_account
    tencentcloud_cynosdb_binlog_save_days
    tencentcloud_cynosdb_cluster_databases
    tencentcloud_cynosdb_instance_param
    tencentcloud_cynosdb_isolate_instance
    tencentcloud_cynosdb_param_template
    tencentcloud_cynosdb_restart_instance
    tencentcloud_cynosdb_roll_back_cluster
    tencentcloud_cynosdb_wan
    tencentcloud_cynosdb_proxy
    tencentcloud_cynosdb_reload_proxy_node
    tencentcloud_cynosdb_cluster_slave_zone
    tencentcloud_cynosdb_read_only_instance_exclusive_access
    tencentcloud_cynosdb_proxy_end_point
    tencentcloud_cynosdb_upgrade_proxy_version

Direct Connect(DC)
  Data Source
    tencentcloud_dc_instances
    tencentcloud_dc_access_points
    tencentcloud_dcx_instances
    tencentcloud_dc_internet_address_quota
    tencentcloud_dc_internet_address_statistics
    tencentcloud_dc_public_direct_connect_tunnel_routes

  Resource
    tencentcloud_dc_instance
    tencentcloud_dcx
    tencentcloud_dcx_extra_config
    tencentcloud_dc_share_dcx_config
    tencentcloud_dc_internet_address
    tencentcloud_dc_internet_address_config

Direct Connect Gateway(DCG)
  Data Source
    tencentcloud_dc_gateway_ccn_routes
    tencentcloud_dc_gateway_instances

  Resource
    tencentcloud_dc_gateway
    tencentcloud_dc_gateway_ccn_route
    tencentcloud_dc_gateway_attachment

Domain
  Data Source
    tencentcloud_domains

Elasticsearch Service(ES)
  Data Source
    tencentcloud_elasticsearch_describe_index_list
    tencentcloud_elasticsearch_instances
    tencentcloud_elasticsearch_instance_logs
    tencentcloud_elasticsearch_instance_operations
    tencentcloud_elasticsearch_logstash_instance_logs
    tencentcloud_elasticsearch_logstash_instance_operations
    tencentcloud_elasticsearch_views
    tencentcloud_elasticsearch_diagnose
    tencentcloud_elasticsearch_instance_plugin_list

  Resource
    tencentcloud_elasticsearch_instance
    tencentcloud_elasticsearch_security_group
    tencentcloud_elasticsearch_logstash
    tencentcloud_elasticsearch_logstash_pipeline
    tencentcloud_elasticsearch_restart_logstash_instance_operation
    tencentcloud_elasticsearch_start_logstash_pipeline_operation
    tencentcloud_elasticsearch_stop_logstash_pipeline_operation
    tencentcloud_elasticsearch_index
    tencentcloud_elasticsearch_restart_instance_operation
    tencentcloud_elasticsearch_restart_nodes_operation
    tencentcloud_elasticsearch_restart_kibana_operation
    tencentcloud_elasticsearch_diagnose
    tencentcloud_elasticsearch_diagnose_instance
    tencentcloud_elasticsearch_update_plugins_operation

Global Application Acceleration(GAAP)
  Data Source
    tencentcloud_gaap_certificates
    tencentcloud_gaap_http_domains
    tencentcloud_gaap_http_rules
    tencentcloud_gaap_layer4_listeners
    tencentcloud_gaap_layer7_listeners
    tencentcloud_gaap_proxies
    tencentcloud_gaap_realservers
    tencentcloud_gaap_security_policies
    tencentcloud_gaap_security_rules
    tencentcloud_gaap_domain_error_pages
    tencentcloud_gaap_access_regions
    tencentcloud_gaap_access_regions_by_dest_region
    tencentcloud_gaap_black_header
    tencentcloud_gaap_country_area_mapping
    tencentcloud_gaap_custom_header
    tencentcloud_gaap_dest_regions
    tencentcloud_gaap_proxy_detail
    tencentcloud_gaap_proxy_groups
    tencentcloud_gaap_proxy_statistics
    tencentcloud_gaap_proxy_group_statistics
    tencentcloud_gaap_real_servers_status
    tencentcloud_gaap_rule_real_servers
    tencentcloud_gaap_resources_by_tag
    tencentcloud_gaap_region_and_price
    tencentcloud_gaap_proxy_and_statistics_listeners
    tencentcloud_gaap_proxies_status
    tencentcloud_gaap_listener_statistics
    tencentcloud_gaap_listener_real_servers
    tencentcloud_gaap_group_and_statistics_proxy
    tencentcloud_gaap_domain_error_page_infos
    tencentcloud_gaap_check_proxy_create

  Resource
    tencentcloud_gaap_proxy
    tencentcloud_gaap_proxy_group
    tencentcloud_gaap_realserver
    tencentcloud_gaap_layer4_listener
    tencentcloud_gaap_layer7_listener
    tencentcloud_gaap_http_domain
    tencentcloud_gaap_http_rule
    tencentcloud_gaap_certificate
    tencentcloud_gaap_security_policy
    tencentcloud_gaap_security_rule
    tencentcloud_gaap_domain_error_page
    tencentcloud_gaap_global_domain_dns
    tencentcloud_gaap_global_domain
    tencentcloud_gaap_custom_header

Key Management Service(KMS)
  Data Source
    tencentcloud_kms_keys
    tencentcloud_kms_public_key
    tencentcloud_kms_get_parameters_for_import
    tencentcloud_kms_describe_keys
    tencentcloud_kms_white_box_key_details
    tencentcloud_kms_list_keys
    tencentcloud_kms_white_box_decrypt_key
    tencentcloud_kms_white_box_device_fingerprints
    tencentcloud_kms_list_algorithms

  Resource
    tencentcloud_kms_key
    tencentcloud_kms_external_key
    tencentcloud_kms_white_box_key
    tencentcloud_kms_cloud_resource_attachment
    tencentcloud_kms_overwrite_white_box_device_fingerprints

Tencent Kubernetes Engine(TKE)
  Data Source
    tencentcloud_kubernetes_clusters
    tencentcloud_kubernetes_cluster_levels
    tencentcloud_kubernetes_charts
    tencentcloud_kubernetes_cluster_common_names
    tencentcloud_kubernetes_available_cluster_versions
    tencentcloud_kubernetes_cluster_authentication_options
    tencentcloud_kubernetes_cluster_node_pools
    tencentcloud_kubernetes_cluster_instances
    tencentcloud_kubernetes_cluster_node_pools
    tencentcloud_kubernetes_cluster_native_node_pools

  Resource
    tencentcloud_kubernetes_cluster
    tencentcloud_kubernetes_scale_worker
    tencentcloud_kubernetes_cluster_attachment
    tencentcloud_kubernetes_node_pool
    tencentcloud_kubernetes_serverless_node_pool
    tencentcloud_kubernetes_backup_storage_location
    tencentcloud_kubernetes_encryption_protection
    tencentcloud_kubernetes_auth_attachment
    tencentcloud_kubernetes_addon_attachment
    tencentcloud_kubernetes_cluster_endpoint
    tencentcloud_kubernetes_addon
    tencentcloud_kubernetes_addon_config
    tencentcloud_kubernetes_native_node_pool
    tencentcloud_kubernetes_health_check_policy
    tencentcloud_kubernetes_log_config
    tencentcloud_kubernetes_cluster_master_attachment

TDMQ for Pulsar(tpulsar)
  Data Source
    tencentcloud_tdmq_environment_attributes
    tencentcloud_tdmq_publisher_summary
    tencentcloud_tdmq_publishers
    tencentcloud_tdmq_pro_instances
    tencentcloud_tdmq_pro_instance_detail

  Resource
    tencentcloud_tdmq_instance
    tencentcloud_tdmq_professional_cluster
    tencentcloud_tdmq_namespace
    tencentcloud_tdmq_topic
    tencentcloud_tdmq_topic_with_full_id
    tencentcloud_tdmq_role
    tencentcloud_tdmq_namespace_role_attachment
    tencentcloud_tdmq_subscription

TencentDB for MongoDB(mongodb)
  Data Source
    tencentcloud_mongodb_instances
    tencentcloud_mongodb_zone_config
    tencentcloud_mongodb_instance_backups
    tencentcloud_mongodb_instance_connections
    tencentcloud_mongodb_instance_current_op
    tencentcloud_mongodb_instance_params
    tencentcloud_mongodb_instance_slow_log

  Resource
    tencentcloud_mongodb_instance
    tencentcloud_mongodb_sharding_instance
    tencentcloud_mongodb_standby_instance
    tencentcloud_mongodb_instance_account
    tencentcloud_mongodb_instance_backup
    tencentcloud_mongodb_instance_transparent_data_encryption
    tencentcloud_mongodb_instance_backup_rule

TencentDB for MySQL(cdb)
  Data Source
    tencentcloud_mysql_backup_list
    tencentcloud_mysql_instance
    tencentcloud_mysql_parameter_list
    tencentcloud_mysql_default_params
    tencentcloud_mysql_zone_config
    tencentcloud_mysql_backup_overview
    tencentcloud_mysql_backup_summaries
    tencentcloud_mysql_bin_log
    tencentcloud_mysql_binlog_backup_overview
    tencentcloud_mysql_clone_list
    tencentcloud_mysql_data_backup_overview
    tencentcloud_mysql_db_features
    tencentcloud_mysql_inst_tables
    tencentcloud_mysql_instance_charset
    tencentcloud_mysql_instance_info
    tencentcloud_mysql_instance_param_record
    tencentcloud_mysql_instance_reboot_time
    tencentcloud_mysql_rollback_range_time
    tencentcloud_mysql_slow_log
    tencentcloud_mysql_slow_log_data
    tencentcloud_mysql_supported_privileges
    tencentcloud_mysql_switch_record
    tencentcloud_mysql_user_task
    tencentcloud_mysql_databases
    tencentcloud_mysql_error_log
    tencentcloud_mysql_project_security_group
    tencentcloud_mysql_ro_min_scale

  Resource
    tencentcloud_mysql_instance
    tencentcloud_mysql_database
    tencentcloud_mysql_readonly_instance
    tencentcloud_mysql_account
    tencentcloud_mysql_privilege
    tencentcloud_mysql_account_privilege
    tencentcloud_mysql_backup_policy
    tencentcloud_mysql_time_window
    tencentcloud_mysql_param_template
    tencentcloud_mysql_deploy_group
    tencentcloud_mysql_security_groups_attachment
    tencentcloud_mysql_local_binlog_config
    tencentcloud_mysql_audit_log_file
    tencentcloud_mysql_backup_download_restriction
    tencentcloud_mysql_renew_db_instance_operation
    tencentcloud_mysql_backup_encryption_status
    tencentcloud_mysql_dr_instance_to_mater
    tencentcloud_mysql_instance_encryption_operation
    tencentcloud_mysql_password_complexity
    tencentcloud_mysql_remote_backup_config
    tencentcloud_mysql_restart_db_instances_operation
    tencentcloud_mysql_switch_for_upgrade
    tencentcloud_mysql_rollback
    tencentcloud_mysql_rollback_stop
    tencentcloud_mysql_ro_group
    tencentcloud_mysql_ro_instance_ip
    tencentcloud_mysql_ro_group_load_operation
    tencentcloud_mysql_switch_master_slave_operation
    tencentcloud_mysql_proxy
    tencentcloud_mysql_reset_root_account
    tencentcloud_mysql_verify_root_account
    tencentcloud_mysql_reload_balance_proxy_node
    tencentcloud_mysql_ro_start_replication
    tencentcloud_mysql_ro_stop_replication
    tencentcloud_mysql_isolate_instance
    tencentcloud_mysql_dr_instance
    tencentcloud_mysql_ssl
    tencentcloud_mysql_cls_log_attachment

Cloud Monitor(Monitor)
  Data Source
    tencentcloud_monitor_policy_conditions
    tencentcloud_monitor_data
    tencentcloud_monitor_product_event
    tencentcloud_monitor_binding_objects
    tencentcloud_monitor_policy_groups
    tencentcloud_monitor_product_namespace
    tencentcloud_monitor_alarm_notices
    tencentcloud_monitor_alarm_history
    tencentcloud_monitor_alarm_metric
    tencentcloud_monitor_alarm_policy
    tencentcloud_monitor_alarm_basic_alarms
    tencentcloud_monitor_alarm_basic_metric
    tencentcloud_monitor_alarm_conditions_template
    tencentcloud_monitor_alarm_notice_callbacks
    tencentcloud_monitor_alarm_all_namespaces
    tencentcloud_monitor_alarm_monitor_type

  Resource
    tencentcloud_monitor_binding_object
    tencentcloud_monitor_policy_binding_object
    tencentcloud_monitor_binding_receiver
    tencentcloud_monitor_alarm_policy
    tencentcloud_monitor_alarm_notice
    tencentcloud_monitor_alarm_policy_set_default


Managed Service for Prometheus(TMP)
  Data Source
    tencentcloud_monitor_tmp_regions
    tencentcloud_monitor_tmp_instances

  Resource
    tencentcloud_monitor_tmp_instance
    tencentcloud_monitor_tmp_alert_rule
    tencentcloud_monitor_tmp_exporter_integration
    tencentcloud_monitor_tmp_cvm_agent
    tencentcloud_monitor_tmp_scrape_job
    tencentcloud_monitor_tmp_recording_rule
    tencentcloud_monitor_tmp_manage_grafana_attachment
    tencentcloud_monitor_tmp_tke_template
    tencentcloud_monitor_tmp_tke_template_attachment
    tencentcloud_monitor_tmp_tke_alert_policy
    tencentcloud_monitor_tmp_tke_config
    tencentcloud_monitor_tmp_tke_record_rule_yaml
    tencentcloud_monitor_tmp_tke_global_notification
    tencentcloud_monitor_tmp_tke_cluster_agent
    tencentcloud_monitor_tmp_tke_basic_config
    tencentcloud_monitor_tmp_alert_group

TencentCloud Managed Service for Grafana(TCMG)
  Data Source
    tencentcloud_monitor_grafana_plugin_overviews

  Resource
    tencentcloud_monitor_grafana_instance
    tencentcloud_monitor_grafana_integration
    tencentcloud_monitor_grafana_notification_channel
    tencentcloud_monitor_grafana_plugin
    tencentcloud_monitor_grafana_sso_account
    tencentcloud_monitor_tmp_grafana_config
    tencentcloud_monitor_grafana_dns_config
    tencentcloud_monitor_grafana_env_config
    tencentcloud_monitor_grafana_whitelist_config
    tencentcloud_monitor_grafana_sso_cam_config
    tencentcloud_monitor_grafana_sso_config
    tencentcloud_monitor_grafana_version_upgrade

TencentDB for PostgreSQL(PostgreSQL)
  Data Source
    tencentcloud_postgresql_instances
    tencentcloud_postgresql_specinfos
    tencentcloud_postgresql_xlogs
    tencentcloud_postgresql_parameter_templates
    tencentcloud_postgresql_readonly_groups
    tencentcloud_postgresql_base_backups
    tencentcloud_postgresql_log_backups
    tencentcloud_postgresql_backup_download_urls
    tencentcloud_postgresql_db_instance_classes
    tencentcloud_postgresql_default_parameters
    tencentcloud_postgresql_recovery_time
    tencentcloud_postgresql_regions
    tencentcloud_postgresql_db_instance_versions
    tencentcloud_postgresql_zones
    tencentcloud_postgresql_account_privileges
    tencentcloud_postgresql_dedicated_clusters

  Resource
    tencentcloud_postgresql_instance
    tencentcloud_postgresql_readonly_instance
    tencentcloud_postgresql_readonly_group
    tencentcloud_postgresql_readonly_attachment
    tencentcloud_postgresql_parameter_template
    tencentcloud_postgresql_backup_plan_config
    tencentcloud_postgresql_security_group_config
    tencentcloud_postgresql_backup_download_restriction_config
    tencentcloud_postgresql_restart_db_instance_operation
    tencentcloud_postgresql_renew_db_instance_operation
    tencentcloud_postgresql_isolate_db_instance_operation
    tencentcloud_postgresql_disisolate_db_instance_operation
    tencentcloud_postgresql_rebalance_readonly_group_operation
    tencentcloud_postgresql_delete_log_backup_operation
    tencentcloud_postgresql_modify_account_remark_operation
    tencentcloud_postgresql_modify_switch_time_period_operation
    tencentcloud_postgresql_base_backup
    tencentcloud_postgresql_instance_ha_config
    tencentcloud_postgresql_account
    tencentcloud_postgresql_apply_parameter_template_operation
    tencentcloud_postgresql_clone_db_instance
    tencentcloud_postgresql_instance_network_access

TencentDB for Redis(crs)
  Data Source
    tencentcloud_redis_zone_config
    tencentcloud_redis_instances
    tencentcloud_redis_backup
    tencentcloud_redis_backup_download_info
    tencentcloud_redis_param_records
    tencentcloud_redis_instance_shards
    tencentcloud_redis_instance_zone_info
    tencentcloud_redis_instance_task_list
    tencentcloud_redis_instance_node_info
    tencentcloud_redis_clusters

  Resource
    tencentcloud_redis_instance
    tencentcloud_redis_backup_config
    tencentcloud_redis_param_template
    tencentcloud_redis_account
    tencentcloud_redis_read_only
    tencentcloud_redis_ssl
    tencentcloud_redis_backup_download_restriction
    tencentcloud_redis_clear_instance_operation
    tencentcloud_redis_renew_instance_operation
    tencentcloud_redis_startup_instance_operation
    tencentcloud_redis_upgrade_cache_version_operation
    tencentcloud_redis_upgrade_multi_zone_operation
    tencentcloud_redis_upgrade_proxy_version_operation
    tencentcloud_redis_maintenance_window
    tencentcloud_redis_replica_readonly
    tencentcloud_redis_switch_master
    tencentcloud_redis_replicate_attachment
    tencentcloud_redis_backup_operation
    tencentcloud_redis_security_group_attachment
    tencentcloud_redis_connection_config
    tencentcloud_redis_log_delivery

Serverless Cloud Function(SCF)
  Data Source
    tencentcloud_scf_functions
    tencentcloud_scf_logs
    tencentcloud_scf_namespaces
    tencentcloud_scf_account_info
    tencentcloud_scf_async_event_management
    tencentcloud_scf_triggers
    tencentcloud_scf_async_event_status
    tencentcloud_scf_function_address
    tencentcloud_scf_request_status
    tencentcloud_scf_function_aliases
    tencentcloud_scf_layer_versions
    tencentcloud_scf_layers
    tencentcloud_scf_function_versions

  Resource
    tencentcloud_scf_function
    tencentcloud_scf_function_version
    tencentcloud_scf_function_event_invoke_config
    tencentcloud_scf_reserved_concurrency_config
    tencentcloud_scf_provisioned_concurrency_config
    tencentcloud_scf_invoke_function
    tencentcloud_scf_sync_invoke_function
    tencentcloud_scf_terminate_async_event
    tencentcloud_scf_namespace
    tencentcloud_scf_layer
    tencentcloud_scf_function_alias
    tencentcloud_scf_trigger_config
    tencentcloud_scf_custom_domain

SQLServer
  Data Source
    tencentcloud_sqlserver_zone_config
    tencentcloud_sqlserver_instances
    tencentcloud_sqlserver_dbs
    tencentcloud_sqlserver_accounts
    tencentcloud_sqlserver_account_db_attachments
    tencentcloud_sqlserver_readonly_groups
    tencentcloud_sqlserver_publish_subscribes
    tencentcloud_sqlserver_basic_instances
    tencentcloud_sqlserver_backup_commands
    tencentcloud_sqlserver_backup_by_flow_id
    tencentcloud_sqlserver_backup_upload_size
    tencentcloud_sqlserver_cross_region_zone
    tencentcloud_sqlserver_db_charsets
    tencentcloud_sqlserver_instance_param_records
    tencentcloud_sqlserver_project_security_groups
    tencentcloud_sqlserver_regions
    tencentcloud_sqlserver_rollback_time
    tencentcloud_sqlserver_slowlogs
    tencentcloud_sqlserver_upload_backup_info
    tencentcloud_sqlserver_upload_incremental_info
    tencentcloud_sqlserver_query_xevent
    tencentcloud_sqlserver_ins_attribute

  Resource
    tencentcloud_sqlserver_instance
    tencentcloud_sqlserver_readonly_instance
    tencentcloud_sqlserver_db
    tencentcloud_sqlserver_account
    tencentcloud_sqlserver_account_db_attachment
    tencentcloud_sqlserver_publish_subscribe
    tencentcloud_sqlserver_basic_instance
    tencentcloud_sqlserver_migration
    tencentcloud_sqlserver_config_backup_strategy
    tencentcloud_sqlserver_general_backup
    tencentcloud_sqlserver_general_clone
    tencentcloud_sqlserver_full_backup_migration
    tencentcloud_sqlserver_incre_backup_migration
    tencentcloud_sqlserver_business_intelligence_file
    tencentcloud_sqlserver_business_intelligence_instance
    tencentcloud_sqlserver_general_communication
    tencentcloud_sqlserver_general_cloud_instance
    tencentcloud_sqlserver_complete_expansion
    tencentcloud_sqlserver_config_database_cdc
    tencentcloud_sqlserver_config_database_ct
    tencentcloud_sqlserver_config_database_mdf
    tencentcloud_sqlserver_config_instance_param
    tencentcloud_sqlserver_config_instance_ro_group
    tencentcloud_sqlserver_renew_db_instance
    tencentcloud_sqlserver_renew_postpaid_db_instance
    tencentcloud_sqlserver_restart_db_instance
    tencentcloud_sqlserver_config_terminate_db_instance
    tencentcloud_sqlserver_restore_instance
    tencentcloud_sqlserver_rollback_instance
    tencentcloud_sqlserver_start_backup_full_migration
    tencentcloud_sqlserver_start_backup_incremental_migration
    tencentcloud_sqlserver_start_xevent
    tencentcloud_sqlserver_instance_tde
    tencentcloud_sqlserver_database_tde
    tencentcloud_sqlserver_general_cloud_ro_instance

SSL Certificates(ssl)
  Data Source
    tencentcloud_ssl_certificates
    tencentcloud_ssl_describe_certificate
    tencentcloud_ssl_describe_companies
    tencentcloud_ssl_describe_host_api_gateway_instance_list
    tencentcloud_ssl_describe_host_cdn_instance_list
    tencentcloud_ssl_describe_host_clb_instance_list
    tencentcloud_ssl_describe_host_cos_instance_list
    tencentcloud_ssl_describe_host_ddos_instance_list
    tencentcloud_ssl_describe_host_lighthouse_instance_list
    tencentcloud_ssl_describe_host_live_instance_list
    tencentcloud_ssl_describe_host_teo_instance_list
    tencentcloud_ssl_describe_host_tke_instance_list
    tencentcloud_ssl_describe_host_vod_instance_list
    tencentcloud_ssl_describe_host_waf_instance_list
    tencentcloud_ssl_describe_host_deploy_record
    tencentcloud_ssl_describe_host_deploy_record_detail
    tencentcloud_ssl_describe_host_update_record
    tencentcloud_ssl_describe_host_update_record_detail
    tencentcloud_ssl_describe_managers
    tencentcloud_ssl_describe_manager_detail

  Resource
    tencentcloud_ssl_certificate
    tencentcloud_ssl_pay_certificate
    tencentcloud_ssl_free_certificate
    tencentcloud_ssl_replace_certificate_operation
    tencentcloud_ssl_revoke_certificate_operation
    tencentcloud_ssl_update_certificate_instance_operation
    tencentcloud_ssl_update_certificate_record_retry_operation
    tencentcloud_ssl_update_certificate_record_rollback_operation
    tencentcloud_ssl_upload_revoke_letter_operation
    tencentcloud_ssl_complete_certificate_operation
    tencentcloud_ssl_check_certificate_chain_operation
    tencentcloud_ssl_deploy_certificate_instance_operation
    tencentcloud_ssl_deploy_certificate_record_retry_operation
    tencentcloud_ssl_deploy_certificate_record_rollback_operation
    tencentcloud_ssl_download_certificate_operation
    tencentcloud_ssl_check_certificate_domain_verification_operation

Secrets Manager(SSM)
  Data Source
    tencentcloud_ssm_products
    tencentcloud_ssm_secrets
    tencentcloud_ssm_secret_versions
    tencentcloud_ssm_rotation_detail
    tencentcloud_ssm_rotation_history
    tencentcloud_ssm_service_status
    tencentcloud_ssm_ssh_key_pair_value

  Resource
    tencentcloud_ssm_secret
    tencentcloud_ssm_secret_version
    tencentcloud_ssm_product_secret
    tencentcloud_ssm_ssh_key_pair_secret
    tencentcloud_ssm_rotate_product_secret

TcaplusDB
  Data Source
    tencentcloud_tcaplus_clusters
    tencentcloud_tcaplus_idls
    tencentcloud_tcaplus_tables
    tencentcloud_tcaplus_tablegroups

  Resource
    tencentcloud_tcaplus_cluster
    tencentcloud_tcaplus_tablegroup
    tencentcloud_tcaplus_idl
    tencentcloud_tcaplus_table

Tencent Container Registry(TCR)
  Data Source
    tencentcloud_tcr_instances
    tencentcloud_tcr_namespaces
    tencentcloud_tcr_repositories
    tencentcloud_tcr_tokens
    tencentcloud_tcr_vpc_attachments
    tencentcloud_tcr_webhook_trigger_logs
    tencentcloud_tcr_images
    tencentcloud_tcr_image_manifests
    tencentcloud_tcr_tag_retention_execution_tasks
    tencentcloud_tcr_tag_retention_executions
    tencentcloud_tcr_replication_instance_create_tasks
    tencentcloud_tcr_replication_instance_sync_status

  Resource
    tencentcloud_tcr_instance
    tencentcloud_tcr_namespace
    tencentcloud_tcr_repository
    tencentcloud_tcr_token
    tencentcloud_tcr_vpc_attachment
    tencentcloud_tcr_tag_retention_rule
    tencentcloud_tcr_webhook_trigger
    tencentcloud_tcr_manage_replication_operation
    tencentcloud_tcr_customized_domain
    tencentcloud_tcr_immutable_tag_rule
    tencentcloud_tcr_delete_image_operation
    tencentcloud_tcr_create_image_signature_operation
    tencentcloud_tcr_tag_retention_execution_config
    tencentcloud_tcr_service_account

Video on Demand(VOD)
  Data Source
    tencentcloud_vod_adaptive_dynamic_streaming_templates
    tencentcloud_vod_snapshot_by_time_offset_templates
    tencentcloud_vod_super_player_configs
    tencentcloud_vod_image_sprite_templates
    tencentcloud_vod_procedure_templates


  Resource
    tencentcloud_vod_adaptive_dynamic_streaming_template
    tencentcloud_vod_procedure_template
    tencentcloud_vod_snapshot_by_time_offset_template
    tencentcloud_vod_image_sprite_template
    tencentcloud_vod_sub_application
    tencentcloud_vod_sample_snapshot_template
    tencentcloud_vod_transcode_template
    tencentcloud_vod_watermark_template
    tencentcloud_vod_event_config

Oceanus
  Data Source
    tencentcloud_oceanus_resource_related_job
    tencentcloud_oceanus_savepoint_list
    tencentcloud_oceanus_system_resource
    tencentcloud_oceanus_work_spaces
    tencentcloud_oceanus_clusters
    tencentcloud_oceanus_tree_jobs
    tencentcloud_oceanus_tree_resources
    tencentcloud_oceanus_job_submission_log
    tencentcloud_oceanus_check_savepoint
    tencentcloud_oceanus_job_events
    tencentcloud_oceanus_meta_table

  Resource
    tencentcloud_oceanus_folder
    tencentcloud_oceanus_job
    tencentcloud_oceanus_job_config
    tencentcloud_oceanus_job_copy
    tencentcloud_oceanus_run_job
    tencentcloud_oceanus_stop_job
    tencentcloud_oceanus_trigger_job_savepoint
    tencentcloud_oceanus_resource
    tencentcloud_oceanus_resource_config
    tencentcloud_oceanus_work_space

Virtual Private Cloud(VPC)
  Data Source
    tencentcloud_route_table
    tencentcloud_security_group
    tencentcloud_security_groups
    tencentcloud_address_templates
    tencentcloud_address_template_groups
    tencentcloud_protocol_templates
    tencentcloud_protocol_template_groups
    tencentcloud_subnet
    tencentcloud_vpc
    tencentcloud_vpc_acls
    tencentcloud_vpc_account_attributes
    tencentcloud_vpc_classic_link_instances
    tencentcloud_vpc_gateway_flow_monitor_detail
    tencentcloud_vpc_gateway_flow_qos
    tencentcloud_vpc_cvm_instances
    tencentcloud_vpc_net_detect_states
    tencentcloud_vpc_net_detect_state_check
    tencentcloud_vpc_network_interface_limit
    tencentcloud_vpc_private_ip_addresses
    tencentcloud_vpc_product_quota
    tencentcloud_vpc_resource_dashboard
    tencentcloud_vpc_route_conflicts
    tencentcloud_vpc_security_group_limits
    tencentcloud_vpc_security_group_references
    tencentcloud_vpc_sg_snapshot_file_content
    tencentcloud_vpc_snapshot_files
    tencentcloud_vpc_subnet_resource_dashboard
    tencentcloud_vpc_template_limits
    tencentcloud_vpc_used_ip_address
    tencentcloud_vpc_limits
    tencentcloud_vpc_instances
    tencentcloud_vpc_route_tables
    tencentcloud_vpc_subnets
    tencentcloud_dnats
    tencentcloud_enis
    tencentcloud_ha_vip_eip_attachments
    tencentcloud_ha_vips
    tencentcloud_nat_gateways
    tencentcloud_nat_gateway_snats
    tencentcloud_nats
    tencentcloud_nat_dc_route
    tencentcloud_vpc_bandwidth_package_quota
    tencentcloud_vpc_bandwidth_package_bill_usage

  Resource
    tencentcloud_eni
    tencentcloud_eni_attachment
    tencentcloud_eni_sg_attachment
    tencentcloud_eni_ipv4_address
    tencentcloud_eni_ipv6_address
    tencentcloud_vpc
    tencentcloud_vpc_acl
    tencentcloud_vpc_acl_attachment
    tencentcloud_vpc_traffic_package
    tencentcloud_vpc_snapshot_policy
    tencentcloud_vpc_snapshot_policy_attachment
    tencentcloud_vpc_snapshot_policy_config
    tencentcloud_vpc_net_detect
    tencentcloud_vpc_dhcp_ip
    tencentcloud_vpc_ipv6_cidr_block
    tencentcloud_vpc_ipv6_subnet_cidr_block
    tencentcloud_vpc_ipv6_eni_address
    tencentcloud_vpc_local_gateway
    tencentcloud_vpc_notify_routes
    tencentcloud_vpc_resume_snapshot_instance 
    tencentcloud_vpc_peer_connect_manager
    tencentcloud_vpc_peer_connect_accept_operation
    tencentcloud_vpc_peer_connect_reject_operation
    tencentcloud_subnet
    tencentcloud_security_group
    tencentcloud_security_group_rule_set
    tencentcloud_security_group_lite_rule
    tencentcloud_address_template
    tencentcloud_address_template_group
    tencentcloud_protocol_template
    tencentcloud_protocol_template_group
    tencentcloud_route_table
    tencentcloud_vpc_private_nat_gateway
    tencentcloud_route_table_association
    tencentcloud_route_entry
    tencentcloud_route_table_entry
    tencentcloud_dnat
    tencentcloud_nat_gateway
    tencentcloud_nat_gateway_snat
    tencentcloud_nat_refresh_nat_dc_route
    tencentcloud_ha_vip
    tencentcloud_ha_vip_eip_attachment
    tencentcloud_vpc_bandwidth_package
    tencentcloud_vpc_bandwidth_package_attachment
    tencentcloud_ipv6_address_bandwidth
    tencentcloud_reserve_ip_address

Private Link(PLS)
  Resource
    tencentcloud_vpc_end_point_service
    tencentcloud_vpc_end_point
    tencentcloud_vpc_enable_end_point_connect
    tencentcloud_vpc_end_point_service_white_list

Flow Logs(FL)
  Resource
     tencentcloud_vpc_flow_log
    tencentcloud_vpc_flow_log_config

VPN Connections(VPN)
  Data Source
    tencentcloud_vpn_connections
    tencentcloud_vpn_customer_gateways
    tencentcloud_vpn_gateways
    tencentcloud_vpn_gateway_routes
    tencentcloud_vpn_customer_gateway_vendors
    tencentcloud_vpn_default_health_check_ip

  Resource
    tencentcloud_vpn_customer_gateway
    tencentcloud_vpn_gateway
    tencentcloud_vpn_gateway_route
    tencentcloud_vpn_connection
    tencentcloud_vpn_ssl_server
    tencentcloud_vpn_ssl_client
    tencentcloud_vpn_connection_reset
    tencentcloud_vpn_customer_gateway_configuration_download
    tencentcloud_vpn_gateway_ssl_client_cert
    tencentcloud_vpn_gateway_ccn_routes

MapReduce(EMR)
  Data Source
    tencentcloud_emr
    tencentcloud_emr_auto_scale_records
    tencentcloud_emr_nodes
    tencentcloud_emr_cvm_quota
    tencentcloud_lite_hbase_instances

  Resource
    tencentcloud_emr_cluster
    tencentcloud_emr_user_manager
    tencentcloud_lite_hbase_instance

DNSPOD
  Resource
    tencentcloud_dnspod_domain_instance
    tencentcloud_dnspod_domain_alias
    tencentcloud_dnspod_record
    tencentcloud_dnspod_record_group
    tencentcloud_dnspod_modify_record_group_operation
    tencentcloud_dnspod_modify_domain_owner_operation
    tencentcloud_dnspod_download_snapshot_operation
    tencentcloud_dnspod_custom_line
    tencentcloud_dnspod_snapshot_config
    tencentcloud_dnspod_domain_lock

  Data Source
    tencentcloud_dnspod_records
    tencentcloud_dnspod_domain_list
    tencentcloud_dnspod_domain_analytics
    tencentcloud_dnspod_domain_log_list
    tencentcloud_dnspod_record_analytics
    tencentcloud_dnspod_record_line_list
    tencentcloud_dnspod_record_list
    tencentcloud_dnspod_record_type

PrivateDNS
  Resource
    tencentcloud_private_dns_zone
    tencentcloud_private_dns_record
    tencentcloud_private_dns_zone_vpc_attachment
    tencentcloud_subscribe_private_zone_service
    tencentcloud_private_dns_forward_rule
    tencentcloud_private_dns_end_point
  Data Source
    tencentcloud_private_dns_records
    tencentcloud_private_dns_private_zone_list
    tencentcloud_private_dns_forward_rules
    tencentcloud_private_dns_end_points

Cloud Log Service(CLS)
  Resource
    tencentcloud_cls_logset
    tencentcloud_cls_topic
    tencentcloud_cls_config
    tencentcloud_cls_config_extra
    tencentcloud_cls_config_attachment
    tencentcloud_cls_machine_group
    tencentcloud_cls_cos_shipper
    tencentcloud_cls_index
    tencentcloud_cls_alarm
    tencentcloud_cls_alarm_notice
    tencentcloud_cls_ckafka_consumer
    tencentcloud_cls_kafka_recharge
    tencentcloud_cls_cos_recharge
    tencentcloud_cls_export
    tencentcloud_cls_scheduled_sql
    tencentcloud_cls_data_transform
    tencentcloud_cls_cloud_product_log_task
    tencentcloud_cls_notice_content

  Data Source
    tencentcloud_cls_shipper_tasks
    tencentcloud_cls_machines
    tencentcloud_cls_machine_group_configs
    tencentcloud_cls_logsets

TencentCloud Lighthouse(Lighthouse)
  Resource
    tencentcloud_lighthouse_instance
    tencentcloud_lighthouse_blueprint
    tencentcloud_lighthouse_firewall_rule
    tencentcloud_lighthouse_disk_backup
    tencentcloud_lighthouse_apply_disk_backup
    tencentcloud_lighthouse_disk_attachment
    tencentcloud_lighthouse_key_pair
    tencentcloud_lighthouse_snapshot
    tencentcloud_lighthouse_apply_instance_snapshot
    tencentcloud_lighthouse_start_instance
    tencentcloud_lighthouse_stop_instance
    tencentcloud_lighthouse_reboot_instance
    tencentcloud_lighthouse_key_pair_attachment
    tencentcloud_lighthouse_disk
    tencentcloud_lighthouse_renew_disk
    tencentcloud_lighthouse_renew_instance
    tencentcloud_lighthouse_firewall_template

  Data Source
    tencentcloud_lighthouse_firewall_rules_template
    tencentcloud_lighthouse_bundle
    tencentcloud_lighthouse_zone
    tencentcloud_lighthouse_scene
    tencentcloud_lighthouse_reset_instance_blueprint
    tencentcloud_lighthouse_region
    tencentcloud_lighthouse_instance_vnc_url
    tencentcloud_lighthouse_instance_traffic_package
    tencentcloud_lighthouse_instance_disk_num
    tencentcloud_lighthouse_instance_blueprint
    tencentcloud_lighthouse_disk_config
    tencentcloud_lighthouse_all_scene
    tencentcloud_lighthouse_modify_instance_bundle
    tencentcloud_lighthouse_disks

TencentCloud Elastic Microservice(TEM)
  Resource
    tencentcloud_tem_environment
    tencentcloud_tem_application
    tencentcloud_tem_workload
    tencentcloud_tem_app_config
    tencentcloud_tem_log_config
    tencentcloud_tem_scale_rule
    tencentcloud_tem_gateway
    tencentcloud_tem_application_service

TencentCloud EdgeOne(TEO)
  Data Source
    tencentcloud_teo_zone_available_plans
    tencentcloud_teo_rule_engine_settings

  Resource
    tencentcloud_teo_zone
    tencentcloud_teo_zone_setting
    tencentcloud_teo_origin_group
    tencentcloud_teo_rule_engine
    tencentcloud_teo_application_proxy_rule
    tencentcloud_teo_ownership_verify
    tencentcloud_teo_certificate_config
    tencentcloud_teo_acceleration_domain
    tencentcloud_teo_l4_proxy
    tencentcloud_teo_realtime_log_delivery
    tencentcloud_teo_function
    tencentcloud_teo_function_rule
    tencentcloud_teo_function_rule_priority
    tencentcloud_teo_function_runtime_environment

TencentCloud ServiceMesh(TCM)
  Data Source
    tencentcloud_tcm_mesh
  Resource
    tencentcloud_tcm_mesh
    tencentcloud_tcm_cluster_attachment
    tencentcloud_tcm_prometheus_attachment
    tencentcloud_tcm_tracing_config
    tencentcloud_tcm_access_log_config

Simple Email Service(SES)
  Data Source
    tencentcloud_ses_receivers
    tencentcloud_ses_send_tasks
    tencentcloud_ses_email_identities
    tencentcloud_ses_black_email_address
    tencentcloud_ses_statistics_report
    tencentcloud_ses_send_email_status

  Resource
    tencentcloud_ses_domain
    tencentcloud_ses_template
    tencentcloud_ses_email_address
    tencentcloud_ses_receiver
    tencentcloud_ses_send_email
    tencentcloud_ses_batch_send_email
    tencentcloud_ses_verify_domain
    tencentcloud_ses_black_list_delete

Security Token Service(STS)
  Data Source
    tencentcloud_sts_caller_identity

TDSQL for MySQL(DCDB)
  Data Source
    tencentcloud_dcdb_instances
    tencentcloud_dcdb_accounts
    tencentcloud_dcdb_databases
    tencentcloud_dcdb_parameters
    tencentcloud_dcdb_shards
    tencentcloud_dcdb_security_groups
    tencentcloud_dcdb_database_objects
    tencentcloud_dcdb_database_tables
    tencentcloud_dcdb_file_download_url
    tencentcloud_dcdb_log_files
    tencentcloud_dcdb_instance_node_info
    tencentcloud_dcdb_orders
    tencentcloud_dcdb_price
    tencentcloud_dcdb_project_security_groups
    tencentcloud_dcdb_projects
    tencentcloud_dcdb_renewal_price
    tencentcloud_dcdb_sale_info
    tencentcloud_dcdb_shard_spec
    tencentcloud_dcdb_slow_logs
    tencentcloud_dcdb_upgrade_price

  Resource
    tencentcloud_dcdb_account
    tencentcloud_dcdb_hourdb_instance
    tencentcloud_dcdb_security_group_attachment
    tencentcloud_dcdb_account_privileges
    tencentcloud_dcdb_db_parameters
    tencentcloud_dcdb_db_sync_mode_config
    tencentcloud_dcdb_encrypt_attributes_config
    tencentcloud_dcdb_instance_config
    tencentcloud_dcdb_cancel_dcn_job_operation
    tencentcloud_dcdb_activate_hour_instance_operation
    tencentcloud_dcdb_isolate_hour_instance_operation
    tencentcloud_dcdb_flush_binlog_operation
    tencentcloud_dcdb_switch_db_instance_ha_operation

Short Message Service(SMS)
  Resource
    tencentcloud_sms_sign
    tencentcloud_sms_template

Cloud Automated Testing(CAT)
  Data Source
    tencentcloud_cat_probe_data
    tencentcloud_cat_node
    tencentcloud_cat_metric_data

  Resource
    tencentcloud_cat_task_set

TencentDB for MariaDB(MariaDB)
  Data Source
    tencentcloud_mariadb_db_instances
    tencentcloud_mariadb_accounts
    tencentcloud_mariadb_security_groups
    tencentcloud_mariadb_database_objects
    tencentcloud_mariadb_databases
    tencentcloud_mariadb_database_table
    tencentcloud_mariadb_dcn_detail
    tencentcloud_mariadb_file_download_url
    tencentcloud_mariadb_flow
    tencentcloud_mariadb_instance_specs
    tencentcloud_mariadb_log_files
    tencentcloud_mariadb_orders
    tencentcloud_mariadb_price
    tencentcloud_mariadb_project_security_groups
    tencentcloud_mariadb_renewal_price
    tencentcloud_mariadb_sale_info
    tencentcloud_mariadb_slow_logs
    tencentcloud_mariadb_upgrade_price

  Resource
    tencentcloud_mariadb_dedicatedcluster_db_instance
    tencentcloud_mariadb_instance
    tencentcloud_mariadb_hour_db_instance
    tencentcloud_mariadb_account
    tencentcloud_mariadb_parameters
    tencentcloud_mariadb_log_file_retention_period
    tencentcloud_mariadb_security_groups
    tencentcloud_mariadb_account_privileges
    tencentcloud_mariadb_operate_hour_db_instance
    tencentcloud_mariadb_backup_time
    tencentcloud_mariadb_cancel_dcn_job
    tencentcloud_mariadb_flush_binlog
    tencentcloud_mariadb_switch_ha
    tencentcloud_mariadb_restart_instance
    tencentcloud_mariadb_renew_instance
    tencentcloud_mariadb_instance_config

Real User Monitoring(RUM)
  Data Source
    tencentcloud_rum_project
    tencentcloud_rum_offline_log_config
    tencentcloud_rum_whitelist
    tencentcloud_rum_taw_instance
    tencentcloud_rum_custom_url
    tencentcloud_rum_event_url
    tencentcloud_rum_fetch_url_info
    tencentcloud_rum_fetch_url
    tencentcloud_rum_group_log
    tencentcloud_rum_log_url_statistics
    tencentcloud_rum_performance_page
    tencentcloud_rum_pv_url_info
    tencentcloud_rum_pv_url_statistics
    tencentcloud_rum_report_count
    tencentcloud_rum_scores
    tencentcloud_rum_set_url_statistics
    tencentcloud_rum_sign
    tencentcloud_rum_static_project
    tencentcloud_rum_static_resource
    tencentcloud_rum_static_url
    tencentcloud_rum_web_vitals_page
    tencentcloud_rum_log_export_list

  Resource
    tencentcloud_rum_project
    tencentcloud_rum_taw_instance
    tencentcloud_rum_whitelist
    tencentcloud_rum_offline_log_config_attachment
    tencentcloud_rum_instance_status_config
    tencentcloud_rum_project_status_config

Cloud Streaming Services(CSS)
  Resource
    tencentcloud_css_watermark
    tencentcloud_css_watermark_rule_attachment
    tencentcloud_css_pull_stream_task
    tencentcloud_css_live_transcode_template
    tencentcloud_css_live_transcode_rule_attachment
    tencentcloud_css_domain
    tencentcloud_css_authenticate_domain_owner_operation
    tencentcloud_css_play_domain_cert_attachment
    tencentcloud_css_play_auth_key_config
    tencentcloud_css_push_auth_key_config
    tencentcloud_css_backup_stream
    tencentcloud_css_callback_rule_attachment
    tencentcloud_css_callback_template
    tencentcloud_css_domain_referer
    tencentcloud_css_enable_optimal_switching
    tencentcloud_css_record_rule_attachment
    tencentcloud_css_snapshot_rule_attachment
    tencentcloud_css_snapshot_template
    tencentcloud_css_pad_template
    tencentcloud_css_pad_rule_attachment
    tencentcloud_css_timeshift_template
    tencentcloud_css_timeshift_rule_attachment
    tencentcloud_css_stream_monitor
    tencentcloud_css_start_stream_monitor
    tencentcloud_css_pull_stream_task_restart

  Data Source
    tencentcloud_css_domains
    tencentcloud_css_backup_stream
    tencentcloud_css_monitor_report
    tencentcloud_css_pad_templates
    tencentcloud_css_pull_stream_task_status
    tencentcloud_css_stream_monitor_list
    tencentcloud_css_time_shift_record_detail
    tencentcloud_css_time_shift_stream_list
    tencentcloud_css_watermarks
    tencentcloud_css_xp2p_detail_info_list

Performance Testing Service(PTS)
  Data Source
    tencentcloud_pts_scenario_with_jobs

  Resource
    tencentcloud_pts_project
    tencentcloud_pts_alert_channel
    tencentcloud_pts_scenario
    tencentcloud_pts_file
    tencentcloud_pts_job
    tencentcloud_pts_cron_job
    tencentcloud_pts_tmp_key_generate
    tencentcloud_pts_cron_job_restart
    tencentcloud_pts_job_abort
    tencentcloud_pts_cron_job_abort

TencentCloud Automation Tools(TAT)
  Data Source
    tencentcloud_tat_command
    tencentcloud_tat_invoker
    tencentcloud_tat_invoker_records
    tencentcloud_tat_agent
    tencentcloud_tat_invocation_task
  Resource
    tencentcloud_tat_command
    tencentcloud_tat_invoker
    tencentcloud_tat_invoker_config
    tencentcloud_tat_invocation_invoke_attachment
    tencentcloud_tat_invocation_command_attachment

Tencent Cloud Organization (TCO)
  Data Source
    tencentcloud_organization_members
    tencentcloud_organization_org_auth_node
    tencentcloud_organization_org_financial_by_member
    tencentcloud_organization_org_financial_by_month
    tencentcloud_organization_org_financial_by_product
    tencentcloud_organization_org_share_area
    tencentcloud_organization_services
    tencentcloud_identity_center_groups
    tencentcloud_identity_center_role_configurations
    tencentcloud_identity_center_users
    tencentcloud_organization_nodes
  Resource
    tencentcloud_organization_instance
    tencentcloud_organization_org_node
    tencentcloud_organization_org_member
    tencentcloud_organization_org_identity
    tencentcloud_organization_org_member_email
    tencentcloud_organization_org_member_auth_identity_attachment
    tencentcloud_organization_org_member_policy_attachment
    tencentcloud_organization_org_manage_policy
    tencentcloud_organization_org_manage_policy_config
    tencentcloud_organization_org_manage_policy_target
    tencentcloud_organization_policy_sub_account_attachment
    tencentcloud_organization_quit_organization_operation
    tencentcloud_organization_org_share_unit
    tencentcloud_organization_org_share_unit_member
    tencentcloud_organization_service_assign
    tencentcloud_identity_center_user
    tencentcloud_identity_center_group
    tencentcloud_identity_center_user_group_attachment
    tencentcloud_identity_center_external_saml_identity_provider
    tencentcloud_identity_center_role_configuration
    tencentcloud_identity_center_role_configuration_permission_policy_attachment
    tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment
    tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment
    tencentcloud_identity_center_user_sync_provisioning
    tencentcloud_identity_center_role_assignment
    tencentcloud_invite_organization_member_operation
    tencentcloud_open_identity_center_operation
    tencentcloud_identity_center_scim_credential_status
    tencentcloud_identity_center_scim_credential
    tencentcloud_identity_center_scim_synchronization_status

TDSQL-C for PostgreSQL(TDCPG)
  Data Source
    tencentcloud_tdcpg_clusters
    tencentcloud_tdcpg_instances
  Resource
    tencentcloud_tdcpg_cluster
    tencentcloud_tdcpg_instance

TencentDB for DBbrain(dbbrain)
  Data Source
    tencentcloud_dbbrain_sql_filters
    tencentcloud_dbbrain_security_audit_log_export_tasks
    tencentcloud_dbbrain_diag_event
    tencentcloud_dbbrain_diag_events
    tencentcloud_dbbrain_diag_history
    tencentcloud_dbbrain_security_audit_log_download_urls
    tencentcloud_dbbrain_slow_log_time_series_stats
    tencentcloud_dbbrain_slow_log_top_sqls
    tencentcloud_dbbrain_slow_log_user_host_stats
    tencentcloud_dbbrain_slow_log_user_sql_advice
    tencentcloud_dbbrain_slow_logs
    tencentcloud_dbbrain_health_scores
    tencentcloud_dbbrain_sql_templates
    tencentcloud_dbbrain_db_space_status
    tencentcloud_dbbrain_top_space_schemas
    tencentcloud_dbbrain_top_space_tables
    tencentcloud_dbbrain_top_space_schema_time_series
    tencentcloud_dbbrain_top_space_table_time_series
    tencentcloud_dbbrain_diag_db_instances
    tencentcloud_dbbrain_mysql_process_list
    tencentcloud_dbbrain_no_primary_key_tables
    tencentcloud_dbbrain_redis_top_big_keys
    tencentcloud_dbbrain_redis_top_key_prefix_list

  Resource
    tencentcloud_dbbrain_sql_filter
    tencentcloud_dbbrain_security_audit_log_export_task
    tencentcloud_dbbrain_db_diag_report_task
    tencentcloud_dbbrain_modify_diag_db_instance_operation
    tencentcloud_dbbrain_tdsql_audit_log

Data Transmission Service(DTS)
  Data Source
    tencentcloud_dts_sync_jobs
    tencentcloud_dts_migrate_jobs
    tencentcloud_dts_compare_tasks
    tencentcloud_dts_migrate_db_instances

  Resource
    tencentcloud_dts_sync_job
    tencentcloud_dts_sync_config
    tencentcloud_dts_sync_check_job_operation
    tencentcloud_dts_sync_job_resume_operation
    tencentcloud_dts_sync_job_start_operation
    tencentcloud_dts_sync_job_stop_operation
    tencentcloud_dts_sync_job_resize_operation
    tencentcloud_dts_sync_job_recover_operation
    tencentcloud_dts_sync_job_isolate_operation
    tencentcloud_dts_sync_job_continue_operation
    tencentcloud_dts_sync_job_pause_operation
    tencentcloud_dts_migrate_service
    tencentcloud_dts_migrate_job
    tencentcloud_dts_migrate_job_config
    tencentcloud_dts_migrate_job_start_operation
    tencentcloud_dts_migrate_job_resume_operation
    tencentcloud_dts_compare_task_stop_operation
    tencentcloud_dts_compare_task

TDMQ for RocketMQ(trocket)
  Data Source
    tencentcloud_tdmq_rocketmq_cluster
    tencentcloud_tdmq_rocketmq_namespace
    tencentcloud_tdmq_rocketmq_topic
    tencentcloud_tdmq_rocketmq_role
    tencentcloud_tdmq_rocketmq_group
    tencentcloud_tdmq_rocketmq_messages

  Resource
    tencentcloud_tdmq_rocketmq_cluster
    tencentcloud_tdmq_rocketmq_namespace
    tencentcloud_tdmq_rocketmq_role
    tencentcloud_tdmq_rocketmq_topic
    tencentcloud_tdmq_rocketmq_group
    tencentcloud_tdmq_rocketmq_environment_role
    tencentcloud_tdmq_send_rocketmq_message
    tencentcloud_tdmq_rocketmq_vip_instance
    tencentcloud_trocket_rocketmq_instance
    tencentcloud_trocket_rocketmq_topic
    tencentcloud_trocket_rocketmq_consumer_group
    tencentcloud_trocket_rocketmq_role

TDMQ for RabbitMQ(trabbit)
  Resource
    tencentcloud_tdmq_rabbitmq_user
    tencentcloud_tdmq_rabbitmq_virtual_host
    tencentcloud_tdmq_rabbitmq_vip_instance


Cloud Infinite(CI)
  Resource
    tencentcloud_ci_bucket_attachment
    tencentcloud_ci_bucket_pic_style
    tencentcloud_ci_hot_link
    tencentcloud_ci_media_snapshot_template
    tencentcloud_ci_media_transcode_template
    tencentcloud_ci_media_animation_template
    tencentcloud_ci_media_concat_template
    tencentcloud_ci_media_video_process_template
    tencentcloud_ci_media_video_montage_template
    tencentcloud_ci_media_voice_separate_template
    tencentcloud_ci_media_super_resolution_template
    tencentcloud_ci_media_pic_process_template
    tencentcloud_ci_media_watermark_template
    tencentcloud_ci_media_tts_template
    tencentcloud_ci_media_transcode_pro_template
    tencentcloud_ci_media_smart_cover_template
    tencentcloud_ci_media_speech_recognition_template
    tencentcloud_ci_guetzli
    tencentcloud_ci_original_image_protection

TDMQ for CMQ(tcmq)
  Data Source
    tencentcloud_tcmq_queue
    tencentcloud_tcmq_topic
    tencentcloud_tcmq_subscribe

  Resource
    tencentcloud_tcmq_queue
    tencentcloud_tcmq_topic
    tencentcloud_tcmq_subscribe

Tencent Service Framework(TSF)
  Data Source
    tencentcloud_tsf_application
    tencentcloud_tsf_application_config
    tencentcloud_tsf_application_file_config
    tencentcloud_tsf_application_public_config
    tencentcloud_tsf_cluster
    tencentcloud_tsf_microservice
    tencentcloud_tsf_unit_rules
    tencentcloud_tsf_config_summary
    tencentcloud_tsf_delivery_config_by_group_id
    tencentcloud_tsf_delivery_configs
    tencentcloud_tsf_public_config_summary
    tencentcloud_tsf_api_group
    tencentcloud_tsf_application_attribute
    tencentcloud_tsf_business_log_configs
    tencentcloud_tsf_api_detail
    tencentcloud_tsf_microservice_api_version
    tencentcloud_tsf_repository
    tencentcloud_tsf_pod_instances
    tencentcloud_tsf_gateway_all_group_apis
    tencentcloud_tsf_group_gateways
    tencentcloud_tsf_usable_unit_namespaces
    tencentcloud_tsf_group_instances
    tencentcloud_tsf_group_config_release
    tencentcloud_tsf_container_group
    tencentcloud_tsf_groups
    tencentcloud_tsf_ms_api_list

  Resource
    tencentcloud_tsf_cluster
    tencentcloud_tsf_microservice
    tencentcloud_tsf_application_config
    tencentcloud_tsf_api_group
    tencentcloud_tsf_namespace
    tencentcloud_tsf_path_rewrite
    tencentcloud_tsf_unit_rule
    tencentcloud_tsf_task
    tencentcloud_tsf_config_template
    tencentcloud_tsf_api_rate_limit_rule
    tencentcloud_tsf_application_release_config
    tencentcloud_tsf_lane
    tencentcloud_tsf_lane_rule
    tencentcloud_tsf_group
    tencentcloud_tsf_application
    tencentcloud_tsf_application_public_config_release
    tencentcloud_tsf_application_public_config
    tencentcloud_tsf_application_file_config_release
    tencentcloud_tsf_instances_attachment
    tencentcloud_tsf_bind_api_group
    tencentcloud_tsf_application_file_config
    tencentcloud_tsf_enable_unit_rule
    tencentcloud_tsf_deploy_container_group
    tencentcloud_tsf_deploy_vm_group
    tencentcloud_tsf_release_api_group
    tencentcloud_tsf_operate_container_group
    tencentcloud_tsf_operate_group
    tencentcloud_tsf_unit_namespace

Media Processing Service(MPS)
  Data Source
    tencentcloud_mps_schedules
    tencentcloud_mps_tasks
    tencentcloud_mps_parse_live_stream_process_notification
    tencentcloud_mps_parse_notification
    tencentcloud_mps_media_meta_data

  Resource
    tencentcloud_mps_workflow
    tencentcloud_mps_enable_workflow_config
    tencentcloud_mps_transcode_template
    tencentcloud_mps_watermark_template
    tencentcloud_mps_image_sprite_template
    tencentcloud_mps_snapshot_by_timeoffset_template
    tencentcloud_mps_sample_snapshot_template
    tencentcloud_mps_animated_graphics_template
    tencentcloud_mps_ai_recognition_template
    tencentcloud_mps_ai_analysis_template
    tencentcloud_mps_adaptive_dynamic_streaming_template
    tencentcloud_mps_person_sample
    tencentcloud_mps_withdraws_watermark_operation
    tencentcloud_mps_process_live_stream_operation
    tencentcloud_mps_edit_media_operation
    tencentcloud_mps_word_sample
    tencentcloud_mps_schedule
    tencentcloud_mps_enable_schedule_config
    tencentcloud_mps_flow
    tencentcloud_mps_input
    tencentcloud_mps_output
    tencentcloud_mps_content_review_template
    tencentcloud_mps_start_flow_operation
    tencentcloud_mps_event
    tencentcloud_mps_manage_task_operation
    tencentcloud_mps_execute_function_operation
    tencentcloud_mps_process_media_operation

Cloud HDFS(CHDFS)
  Data Source
    tencentcloud_chdfs_access_groups
    tencentcloud_chdfs_mount_points
    tencentcloud_chdfs_file_systems

  Resource
    tencentcloud_chdfs_access_group
    tencentcloud_chdfs_access_rule
    tencentcloud_chdfs_file_system
    tencentcloud_chdfs_life_cycle_rule
    tencentcloud_chdfs_mount_point
    tencentcloud_chdfs_mount_point_attachment

StreamLive(MDL)
  Resource
    tencentcloud_mdl_stream_live_input

Application Performance Management(APM)
  Resource
    tencentcloud_apm_instance

Tencent Cloud Service Engine(TSE)
  Data Source
    tencentcloud_tse_access_address
    tencentcloud_tse_nacos_replicas
    tencentcloud_tse_zookeeper_replicas
    tencentcloud_tse_zookeeper_server_interfaces
    tencentcloud_tse_nacos_server_interfaces
    tencentcloud_tse_groups
    tencentcloud_tse_gateways
    tencentcloud_tse_gateway_nodes
    tencentcloud_tse_gateway_routes
    tencentcloud_tse_gateway_canary_rules
    tencentcloud_tse_gateway_services
    tencentcloud_tse_gateway_certificates

  Resource
    tencentcloud_tse_instance
    tencentcloud_tse_cngw_service
    tencentcloud_tse_cngw_canary_rule
    tencentcloud_tse_cngw_gateway
    tencentcloud_tse_cngw_group
    tencentcloud_tse_cngw_service_rate_limit
    tencentcloud_tse_cngw_route
    tencentcloud_tse_cngw_route_rate_limit
    tencentcloud_tse_cngw_certificate
    tencentcloud_tse_waf_protection
    tencentcloud_tse_waf_domains
    tencentcloud_tse_cngw_network
    tencentcloud_tse_cngw_strategy
    tencentcloud_tse_cngw_strategy_bind_group
    tencentcloud_tse_cngw_network_access_control

ClickHouse(CDWCH)
  Data Source
    tencentcloud_clickhouse_backup_jobs
    tencentcloud_clickhouse_backup_job_detail
    tencentcloud_clickhouse_backup_tables
    tencentcloud_clickhouse_spec
    tencentcloud_clickhouse_instance_shards
    tencentcloud_clickhouse_instance_nodes

  Resource
    tencentcloud_clickhouse_instance
    tencentcloud_clickhouse_backup
    tencentcloud_clickhouse_backup_strategy
    tencentcloud_clickhouse_recover_backup_job
    tencentcloud_clickhouse_delete_backup_data
    tencentcloud_clickhouse_account
    tencentcloud_clickhouse_account_permission
    tencentcloud_clickhouse_keyval_config 
    tencentcloud_clickhouse_xml_config

Tag
  Resource
    tencentcloud_tag
    tencentcloud_tag_attachment

EventBridge(EB)
  Data Source
    tencentcloud_eb_bus
    tencentcloud_eb_event_rules
    tencentcloud_eb_platform_event_names
    tencentcloud_eb_platform_event_patterns
    tencentcloud_eb_platform_products
    tencentcloud_eb_plateform_event_template

  Resource
    tencentcloud_eb_event_transform
    tencentcloud_eb_event_bus
    tencentcloud_eb_event_rule
    tencentcloud_eb_event_target
    tencentcloud_eb_put_events
    tencentcloud_eb_event_connector

Data Lake Compute(DLC)
  Data Source
    tencentcloud_dlc_describe_user_type
    tencentcloud_dlc_describe_user_info
    tencentcloud_dlc_describe_user_roles
    tencentcloud_dlc_describe_data_engine
    tencentcloud_dlc_describe_data_engine_image_versions
    tencentcloud_dlc_describe_data_engine_python_spark_images
    tencentcloud_dlc_describe_engine_usage_info
    tencentcloud_dlc_describe_work_group_info
    tencentcloud_dlc_check_data_engine_image_can_be_rollback
    tencentcloud_dlc_check_data_engine_image_can_be_upgrade
    tencentcloud_dlc_check_data_engine_config_pairs_validity
    tencentcloud_dlc_describe_updatable_data_engines
    tencentcloud_dlc_describe_data_engine_events

  Resource
    tencentcloud_dlc_work_group
    tencentcloud_dlc_user
    tencentcloud_dlc_data_engine
    tencentcloud_dlc_rollback_data_engine_image_operation
    tencentcloud_dlc_add_users_to_work_group_attachment
    tencentcloud_dlc_store_location_config
    tencentcloud_dlc_suspend_resume_data_engine
    tencentcloud_dlc_modify_data_engine_description_operation
    tencentcloud_dlc_modify_user_typ_operation
    tencentcloud_dlc_renew_data_engine_operation
    tencentcloud_dlc_restart_data_engine_operation
    tencentcloud_dlc_switch_data_engine_image_operation
    tencentcloud_dlc_update_data_engine_config_operation
    tencentcloud_dlc_upgrade_data_engine_image_operation
    tencentcloud_dlc_user_data_engine_config
    tencentcloud_dlc_update_row_filter_operation
    tencentcloud_dlc_bind_work_groups_to_user_attachment

Web Application Firewall(WAF)
  Data Source
    tencentcloud_waf_ciphers
    tencentcloud_waf_tls_versions
    tencentcloud_waf_domains
    tencentcloud_waf_find_domains
    tencentcloud_waf_ports
    tencentcloud_waf_user_domains
    tencentcloud_waf_attack_log_histogram
    tencentcloud_waf_attack_log_list
    tencentcloud_waf_attack_overview
    tencentcloud_waf_attack_total_count
    tencentcloud_waf_peak_points
    tencentcloud_waf_instance_qps_limit
    tencentcloud_waf_user_clb_regions

  Resource
    tencentcloud_waf_custom_rule
    tencentcloud_waf_custom_white_rule
    tencentcloud_waf_clb_domain
    tencentcloud_waf_saas_domain
    tencentcloud_waf_clb_instance
    tencentcloud_waf_saas_instance
    tencentcloud_waf_anti_fake
    tencentcloud_waf_anti_info_leak
    tencentcloud_waf_auto_deny_rules
    tencentcloud_waf_module_status
    tencentcloud_waf_protection_mode
    tencentcloud_waf_web_shell
    tencentcloud_waf_cc
    tencentcloud_waf_cc_auto_status
    tencentcloud_waf_cc_session
    tencentcloud_waf_ip_access_control
    tencentcloud_waf_ip_access_control_v2

Wedata
  Data Source
    tencentcloud_wedata_rule_templates
    tencentcloud_wedata_data_source_list
    tencentcloud_wedata_data_source_without_info

  Resource
    tencentcloud_wedata_datasource
    tencentcloud_wedata_function
    tencentcloud_wedata_resource
    tencentcloud_wedata_script
    tencentcloud_wedata_dq_rule
    tencentcloud_wedata_rule_template
    tencentcloud_wedata_baseline
    tencentcloud_wedata_integration_offline_task
    tencentcloud_wedata_integration_realtime_task
    tencentcloud_wedata_integration_task_node

Cloud Firewall(CFW)
  Data Source
    tencentcloud_cfw_nat_fw_switches
    tencentcloud_cfw_vpc_fw_switches
    tencentcloud_cfw_edge_fw_switches

  Resource
    tencentcloud_cfw_address_template
    tencentcloud_cfw_block_ignore
    tencentcloud_cfw_edge_policy
    tencentcloud_cfw_nat_instance
    tencentcloud_cfw_nat_policy
    tencentcloud_cfw_vpc_instance
    tencentcloud_cfw_vpc_policy
    tencentcloud_cfw_sync_asset
    tencentcloud_cfw_sync_route
    tencentcloud_cfw_nat_firewall_switch
    tencentcloud_cfw_vpc_firewall_switch
    tencentcloud_cfw_edge_firewall_switch
    tencentcloud_sg_rule

Bastion Host(BH)
  Resource
    tencentcloud_dasb_resource
    tencentcloud_dasb_acl
    tencentcloud_dasb_cmd_template
    tencentcloud_dasb_device_group
    tencentcloud_dasb_user
    tencentcloud_dasb_device_account
    tencentcloud_dasb_device_group_members
    tencentcloud_dasb_user_group_members
    tencentcloud_dasb_bind_device_resource
    tencentcloud_dasb_device
    tencentcloud_dasb_user_group
    tencentcloud_dasb_reset_user
    tencentcloud_dasb_bind_device_account_private_key
    tencentcloud_dasb_bind_device_account_password

Cwp
  Data Source
    tencentcloud_cwp_machines_simple

  Resource
    tencentcloud_cwp_license_order
    tencentcloud_cwp_license_bind_attachment

Business Intelligence(BI)
  Data Source
    tencentcloud_bi_project
    tencentcloud_bi_user_project

  Resource
    tencentcloud_bi_project
    tencentcloud_bi_user_role
    tencentcloud_bi_project_user_role
    tencentcloud_bi_datasource
    tencentcloud_bi_datasource_cloud
    tencentcloud_bi_embed_token_apply
    tencentcloud_bi_embed_interval_apply

CDWPG
  Resource
    tencentcloud_cdwpg_instance

CSIP
  Resource
    tencentcloud_csip_risk_center

CDC
  Data Source
    tencentcloud_cdc_dedicated_cluster_hosts
    tencentcloud_cdc_dedicated_cluster_instance_types
    tencentcloud_cdc_dedicated_cluster_orders

  Resource
    tencentcloud_cdc_site
    tencentcloud_cdc_dedicated_cluster
  
CdwDoris
  Data Source
    tencentcloud_cdwdoris_instances
  Resource
    tencentcloud_cdwdoris_instance
    tencentcloud_cdwdoris_workload_group

ControlCenter
  Resource
    tencentcloud_batch_apply_account_baselines

THPC
  Resource
    tencentcloud_thpc_workspaces

Tencent Container Security Service(TCSS)
  Resource
    tencentcloud_tcss_image_registry
