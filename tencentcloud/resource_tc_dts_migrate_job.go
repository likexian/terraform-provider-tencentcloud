/*
Provides a resource to create a dts migrate_job

Example Usage

```hcl
resource "tencentcloud_dts_migrate_service" "service" {
	src_database_type = "mysql"
	dst_database_type = "cynosdbmysql"
	src_region = "ap-guangzhou"
	dst_region = "ap-guangzhou"
	instance_class = "small"
	job_name = "tf_test_migration_service_1"
	tags {
	  tag_key = "aaa"
	  tag_value = "bbb"
	}
  }

resource "tencentcloud_dts_migrate_job" "job" {
  	service_id = tencentcloud_dts_migrate_service.service.id
	run_mode = "immediate"
	migrate_option {
		database_table {
			object_mode = "partial"
			databases {
				db_name = "tf_ci_test"
				db_mode = "partial"
				table_mode = "partial"
				// view_mode = "partial"
				tables {
					table_name = "test"
					new_table_name = "test_xxx"
					table_edit_mode = "rename"
				}
			}
		}
	}
	src_info {
			region = "ap-guangzhou"
			access_type = "cdb"
			database_type = "mysql"
			node_type = "simple"
			info {
				user = "user"
				password = "password"
				instance_id = "cdb-xxx"
			}

	}
	dst_info {
			region = "ap-guangzhou"
			access_type = "cdb"
			database_type = "cynosdbmysql"
			node_type = "simple"
			info {
				user = "user"
				password = "password"
				instance_id = "cynosdbmysql-xxx"
			}
	}
	job_name = "tf_migrate_job_config_test"
	auto_retry_time_range_minutes = 0
}

resource "tencentcloud_dts_migrate_job_start_operation" "start"{
	job_id = tencentcloud_dts_migrate_job.job.id
}
```

Import

dts migrate_job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job.migrate_job migrate_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDtsMigrateJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobCreate,
		Read:   resourceTencentCloudDtsMigrateJobRead,
		// Update: resourceTencentCloudDtsMigrateJobUpdate,
		Delete: resourceTencentCloudDtsMigrateJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
				Description: "Migrate service Id from `tencentcloud_dts_migrate_service`.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Migrate job status.",
			},

			"resume_option": { // for resume operation
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The mode of the recovery task, the valid values: `clearData`: clears the target instance data. `overwrite`: executes the task in an overwriting way. `normal`: the normal process, no additional action is performed.",
			},

			"complete_mode": { // for complete operation
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The way to complete the task, only support the old version of MySQL migration task, the valid values: waitForSync,immediately.",
			},

			// for modify operation
			"run_mode": {
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeString,
				Description: "Run Mode. eg:immediate,timed.",
			},

			"migrate_option": {
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Migration job configuration options, used to describe how the task performs migration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_table": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Migration object option, you need to tell the migration service which library table objects to migrate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Object mode. eg:all,partial.",
									},
									"databases": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The database list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "database name.",
												},
												"new_db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "New database name.",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "schema name.",
												},
												"new_schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "schema name after migration or synchronization.",
												},
												"db_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "DB selection mode:all (for all objects under the current object), partial (partial objects), when the ObjectMode is partial, this item is required.",
												},
												"schema_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "schema mode: all,partial.",
												},
												"table_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "table mode: all,partial.",
												},
												"tables": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "tables list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "table name.",
															},
															"new_table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "new table name.",
															},
															"tmp_tables": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Optional:    true,
																Description: "temporary tables.",
															},
															"table_edit_mode": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "table edit mode.",
															},
														},
													},
												},
												"view_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "ViewMode.",
												},
												"views": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Views.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "ViewName.",
															},
															"new_view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "NewViewName.",
															},
														},
													},
												},
												"role_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "RoleMode.",
												},
												"roles": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Roles.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"role_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "RoleName.",
															},
															"new_role_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "NewRoleName.",
															},
														},
													},
												},
												"function_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "FunctionMode.",
												},
												"trigger_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "TriggerMode.",
												},
												"event_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "EventMode.",
												},
												"procedure_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "ProcedureMode.",
												},
												"functions": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Functions.",
												},
												"procedures": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Procedures.",
												},
												"events": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Events.",
												},
												"triggers": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "Triggers.",
												},
											},
										},
									},
									"advanced_objects": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "AdvancedObjects.",
									},
								},
							},
						},
						"migrate_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "MigrateType.",
						},
						"consistency": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Consistency.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ConsistencyOption.",
									},
								},
							},
						},
						"is_migrate_account": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "IsMigrateAccount.",
						},
						"is_override_root": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "IsOverrideRoot.",
						},
						"is_dst_read_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "IsDstReadOnly.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ExtraAttr.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"src_info": {
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "SrcInfo.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AccessType.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DatabaseType.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NodeType.",
						},
						"info": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Role.",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "DbKernel.",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Host.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Port.",
									},
									"user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "User.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "Password.",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CvmInstanceId.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UniqVpnGwId.",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UniqDcgId.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "InstanceId.",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CcnGwId.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VpcId.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SubnetId.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "EngineVersion.",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account.",
									},
									"account_role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "AccountRole.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "AccountMode.",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TmpSecretId.",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TmpSecretKey.",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "TmpToken.",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Supplier.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ExtraAttr.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"dst_info": {
				ForceNew:    true,
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "DstInfo.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AccessType.",
						},
						"database_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DatabaseType.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NodeType.",
						},
						"info": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Role.",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "DbKernel.",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Host.",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Port.",
									},
									"user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "User.",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "Password.",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CvmInstanceId.",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UniqVpnGwId.",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "UniqDcgId.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "InstanceId.",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CcnGwId.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VpcId.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SubnetId.",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Engine Version.",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account.",
									},
									"account_role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account Role.",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Account Mode.",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tmp SecretId.",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tmp SecretKey.",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Tmp Token.",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Supplier.",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ExtraAttr.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Value.",
									},
								},
							},
						},
					},
				},
			},

			"job_name": {
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
				Description: "JobName.",
			},

			"expect_run_time": {
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ExpectRunTime.",
			},

			"auto_retry_time_range_minutes": {
				ForceNew:    true,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "AutoRetryTimeRangeMinutes.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	// ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		tcClient  = meta.(*TencentCloudClient).apiV3Conn
		service   = DtsService{client: tcClient}
		serviceId string
		conf      *resource.StateChangeConf
	)

	if v, ok := d.GetOk("service_id"); ok {
		serviceId = v.(string)
	}

	// case "modify":
	err := handleModifyMigrate(d, tcClient, logId, serviceId)
	if err != nil {
		return err
	}

	conf = BuildStateChangeConf([]string{}, []string{"created"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(serviceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	// case "check":
	err = handleCheckMigrate(d, tcClient, logId, serviceId)
	if err != nil {
		return err
	}

	conf = BuildStateChangeConf([]string{}, []string{"checkPass", "checkNotPass"}, 3*readRetryTimeout, time.Second, service.DtsMigrateCheckConfigStateRefreshFunc(serviceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	// // case "resume":
	// err = handleResumeMigrate(d, tcClient, logId, jobId)
	// if err != nil {
	// 	return err
	// }

	// // case "compare":
	// err = handleCompareMigrate(d, tcClient, logId, jobId)
	// if err != nil {
	// 	return err
	// }

	// // case "complete":
	// err = handleCompleteMigrate(d, tcClient, logId, jobId)
	// if err != nil {
	// 	return err
	// }

	// // case "stop":
	// err = handleStopMigrate(d, tcClient, logId, jobId)
	// if err != nil {
	// 	return err
	// }

	d.SetId(serviceId)
	return resourceTencentCloudDtsMigrateJobRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	jobId := d.Id()

	migrateJob, err := service.DescribeDtsMigrateJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if migrateJob.JobId != nil {
		_ = d.Set("service_id", migrateJob.JobId)
	}

	if migrateJob.Status != nil {
		_ = d.Set("status", migrateJob.Status)
	}

	// for modify operation
	if migrateJob.RunMode != nil {
		_ = d.Set("run_mode", migrateJob.RunMode)
	}

	if migrateJob.MigrateOption != nil {
		migrateOptionMap := map[string]interface{}{}

		if migrateJob.MigrateOption.DatabaseTable != nil {
			databaseTableMap := map[string]interface{}{}

			if migrateJob.MigrateOption.DatabaseTable.ObjectMode != nil {
				databaseTableMap["object_mode"] = migrateJob.MigrateOption.DatabaseTable.ObjectMode
			}

			if migrateJob.MigrateOption.DatabaseTable.Databases != nil {
				databasesList := []interface{}{}
				for _, databases := range migrateJob.MigrateOption.DatabaseTable.Databases {
					databasesMap := map[string]interface{}{}

					if databases.DbName != nil {
						databasesMap["db_name"] = databases.DbName
					}

					if databases.NewDbName != nil {
						databasesMap["new_db_name"] = databases.NewDbName
					}

					if databases.SchemaName != nil {
						databasesMap["schema_name"] = databases.SchemaName
					}

					if databases.NewSchemaName != nil {
						databasesMap["new_schema_name"] = databases.NewSchemaName
					}

					if databases.DBMode != nil {
						databasesMap["db_mode"] = databases.DBMode
					}

					if databases.SchemaMode != nil {
						databasesMap["schema_mode"] = databases.SchemaMode
					}

					if databases.TableMode != nil {
						databasesMap["table_mode"] = databases.TableMode
					}

					if databases.Tables != nil {
						tablesList := []interface{}{}
						for _, tables := range databases.Tables {
							tablesMap := map[string]interface{}{}

							if tables.TableName != nil {
								tablesMap["table_name"] = tables.TableName
							}

							if tables.NewTableName != nil {
								tablesMap["new_table_name"] = tables.NewTableName
							}

							if tables.TmpTables != nil {
								tablesMap["tmp_tables"] = tables.TmpTables
							}

							if tables.TableEditMode != nil {
								tablesMap["table_edit_mode"] = tables.TableEditMode
							}

							tablesList = append(tablesList, tablesMap)
						}

						databasesMap["tables"] = []interface{}{tablesList}
					}

					if databases.ViewMode != nil {
						databasesMap["view_mode"] = databases.ViewMode
					}

					if databases.Views != nil {
						viewsList := []interface{}{}
						for _, views := range databases.Views {
							viewsMap := map[string]interface{}{}

							if views.ViewName != nil {
								viewsMap["view_name"] = views.ViewName
							}

							if views.NewViewName != nil {
								viewsMap["new_view_name"] = views.NewViewName
							}

							viewsList = append(viewsList, viewsMap)
						}

						databasesMap["views"] = []interface{}{viewsList}
					}

					if databases.RoleMode != nil {
						databasesMap["role_mode"] = databases.RoleMode
					}

					if databases.Roles != nil {
						rolesList := []interface{}{}
						for _, roles := range databases.Roles {
							rolesMap := map[string]interface{}{}

							if roles.RoleName != nil {
								rolesMap["role_name"] = roles.RoleName
							}

							if roles.NewRoleName != nil {
								rolesMap["new_role_name"] = roles.NewRoleName
							}

							rolesList = append(rolesList, rolesMap)
						}

						databasesMap["roles"] = []interface{}{rolesList}
					}

					if databases.FunctionMode != nil {
						databasesMap["function_mode"] = databases.FunctionMode
					}

					if databases.TriggerMode != nil {
						databasesMap["trigger_mode"] = databases.TriggerMode
					}

					if databases.EventMode != nil {
						databasesMap["event_mode"] = databases.EventMode
					}

					if databases.ProcedureMode != nil {
						databasesMap["procedure_mode"] = databases.ProcedureMode
					}

					if databases.Functions != nil {
						databasesMap["functions"] = databases.Functions
					}

					if databases.Procedures != nil {
						databasesMap["procedures"] = databases.Procedures
					}

					if databases.Events != nil {
						databasesMap["events"] = databases.Events
					}

					if databases.Triggers != nil {
						databasesMap["triggers"] = databases.Triggers
					}

					databasesList = append(databasesList, databasesMap)
				}

				databaseTableMap["databases"] = []interface{}{databasesList}
			}

			if migrateJob.MigrateOption.DatabaseTable.AdvancedObjects != nil {
				databaseTableMap["advanced_objects"] = migrateJob.MigrateOption.DatabaseTable.AdvancedObjects
			}

			migrateOptionMap["database_table"] = []interface{}{databaseTableMap}
		}

		if migrateJob.MigrateOption.MigrateType != nil {
			migrateOptionMap["migrate_type"] = migrateJob.MigrateOption.MigrateType
		}

		if migrateJob.MigrateOption.Consistency != nil {
			consistencyMap := map[string]interface{}{}

			if migrateJob.MigrateOption.Consistency.Mode != nil {
				consistencyMap["mode"] = migrateJob.MigrateOption.Consistency.Mode
			}

			migrateOptionMap["consistency"] = []interface{}{consistencyMap}
		}

		if migrateJob.MigrateOption.IsMigrateAccount != nil {
			migrateOptionMap["is_migrate_account"] = migrateJob.MigrateOption.IsMigrateAccount
		}

		if migrateJob.MigrateOption.IsOverrideRoot != nil {
			migrateOptionMap["is_override_root"] = migrateJob.MigrateOption.IsOverrideRoot
		}

		if migrateJob.MigrateOption.IsDstReadOnly != nil {
			migrateOptionMap["is_dst_read_only"] = migrateJob.MigrateOption.IsDstReadOnly
		}

		if migrateJob.MigrateOption.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateJob.MigrateOption.ExtraAttr {
				extraAttrMap := map[string]interface{}{}

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			migrateOptionMap["extra_attr"] = []interface{}{extraAttrList}
		}

		_ = d.Set("migrate_option", []interface{}{migrateOptionMap})
	}

	if migrateJob.SrcInfo != nil {
		srcInfoMap := map[string]interface{}{}

		if migrateJob.SrcInfo.Region != nil {
			srcInfoMap["region"] = migrateJob.SrcInfo.Region
		}

		if migrateJob.SrcInfo.AccessType != nil {
			srcInfoMap["access_type"] = migrateJob.SrcInfo.AccessType
		}

		if migrateJob.SrcInfo.DatabaseType != nil {
			srcInfoMap["database_type"] = migrateJob.SrcInfo.DatabaseType
		}

		if migrateJob.SrcInfo.NodeType != nil {
			srcInfoMap["node_type"] = migrateJob.SrcInfo.NodeType
		}

		if migrateJob.SrcInfo.Info != nil {
			infoList := []interface{}{}
			for _, info := range migrateJob.SrcInfo.Info {
				infoMap := map[string]interface{}{}

				if info.Role != nil {
					infoMap["role"] = info.Role
				}

				if info.DbKernel != nil {
					infoMap["db_kernel"] = info.DbKernel
				}

				if info.Host != nil {
					infoMap["host"] = info.Host
				}

				if info.Port != nil {
					infoMap["port"] = info.Port
				}

				if info.User != nil {
					infoMap["user"] = info.User
				}

				if info.Password != nil {
					infoMap["password"] = info.Password
				}

				if info.CvmInstanceId != nil {
					infoMap["cvm_instance_id"] = info.CvmInstanceId
				}

				if info.UniqVpnGwId != nil {
					infoMap["uniq_vpn_gw_id"] = info.UniqVpnGwId
				}

				if info.UniqDcgId != nil {
					infoMap["uniq_dcg_id"] = info.UniqDcgId
				}

				if info.InstanceId != nil {
					infoMap["instance_id"] = info.InstanceId
				}

				if info.CcnGwId != nil {
					infoMap["ccn_gw_id"] = info.CcnGwId
				}

				if info.VpcId != nil {
					infoMap["vpc_id"] = info.VpcId
				}

				if info.SubnetId != nil {
					infoMap["subnet_id"] = info.SubnetId
				}

				if info.EngineVersion != nil {
					infoMap["engine_version"] = info.EngineVersion
				}

				if info.Account != nil {
					infoMap["account"] = info.Account
				}

				if info.AccountRole != nil {
					infoMap["account_role"] = info.AccountRole
				}

				if info.AccountMode != nil {
					infoMap["account_mode"] = info.AccountMode
				}

				if info.TmpSecretId != nil {
					infoMap["tmp_secret_id"] = info.TmpSecretId
				}

				if info.TmpSecretKey != nil {
					infoMap["tmp_secret_key"] = info.TmpSecretKey
				}

				if info.TmpToken != nil {
					infoMap["tmp_token"] = info.TmpToken
				}

				infoList = append(infoList, infoMap)
			}

			srcInfoMap["info"] = []interface{}{infoList}
		}

		if migrateJob.SrcInfo.Supplier != nil {
			srcInfoMap["supplier"] = migrateJob.SrcInfo.Supplier
		}

		if migrateJob.SrcInfo.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateJob.SrcInfo.ExtraAttr {
				extraAttrMap := map[string]interface{}{}

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			srcInfoMap["extra_attr"] = []interface{}{extraAttrList}
		}

		_ = d.Set("src_info", []interface{}{srcInfoMap})
	}

	if migrateJob.DstInfo != nil {
		dstInfoMap := map[string]interface{}{}

		if migrateJob.DstInfo.Region != nil {
			dstInfoMap["region"] = migrateJob.DstInfo.Region
		}

		if migrateJob.DstInfo.AccessType != nil {
			dstInfoMap["access_type"] = migrateJob.DstInfo.AccessType
		}

		if migrateJob.DstInfo.DatabaseType != nil {
			dstInfoMap["database_type"] = migrateJob.DstInfo.DatabaseType
		}

		if migrateJob.DstInfo.NodeType != nil {
			dstInfoMap["node_type"] = migrateJob.DstInfo.NodeType
		}

		if migrateJob.DstInfo.Info != nil {
			infoList := []interface{}{}
			for _, info := range migrateJob.DstInfo.Info {
				infoMap := map[string]interface{}{}

				if info.Role != nil {
					infoMap["role"] = info.Role
				}

				if info.DbKernel != nil {
					infoMap["db_kernel"] = info.DbKernel
				}

				if info.Host != nil {
					infoMap["host"] = info.Host
				}

				if info.Port != nil {
					infoMap["port"] = info.Port
				}

				if info.User != nil {
					infoMap["user"] = info.User
				}

				if info.Password != nil {
					infoMap["password"] = info.Password
				}

				if info.CvmInstanceId != nil {
					infoMap["cvm_instance_id"] = info.CvmInstanceId
				}

				if info.UniqVpnGwId != nil {
					infoMap["uniq_vpn_gw_id"] = info.UniqVpnGwId
				}

				if info.UniqDcgId != nil {
					infoMap["uniq_dcg_id"] = info.UniqDcgId
				}

				if info.InstanceId != nil {
					infoMap["instance_id"] = info.InstanceId
				}

				if info.CcnGwId != nil {
					infoMap["ccn_gw_id"] = info.CcnGwId
				}

				if info.VpcId != nil {
					infoMap["vpc_id"] = info.VpcId
				}

				if info.SubnetId != nil {
					infoMap["subnet_id"] = info.SubnetId
				}

				if info.EngineVersion != nil {
					infoMap["engine_version"] = info.EngineVersion
				}

				if info.Account != nil {
					infoMap["account"] = info.Account
				}

				if info.AccountRole != nil {
					infoMap["account_role"] = info.AccountRole
				}

				if info.AccountMode != nil {
					infoMap["account_mode"] = info.AccountMode
				}

				if info.TmpSecretId != nil {
					infoMap["tmp_secret_id"] = info.TmpSecretId
				}

				if info.TmpSecretKey != nil {
					infoMap["tmp_secret_key"] = info.TmpSecretKey
				}

				if info.TmpToken != nil {
					infoMap["tmp_token"] = info.TmpToken
				}

				infoList = append(infoList, infoMap)
			}

			dstInfoMap["info"] = []interface{}{infoList}
		}

		if migrateJob.DstInfo.Supplier != nil {
			dstInfoMap["supplier"] = migrateJob.DstInfo.Supplier
		}

		if migrateJob.DstInfo.ExtraAttr != nil {
			extraAttrList := []interface{}{}
			for _, extraAttr := range migrateJob.DstInfo.ExtraAttr {
				extraAttrMap := map[string]interface{}{}

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			dstInfoMap["extra_attr"] = []interface{}{extraAttrList}
		}

		_ = d.Set("dst_info", []interface{}{dstInfoMap})
	}

	if migrateJob.JobName != nil {
		_ = d.Set("job_name", migrateJob.JobName)
	}

	if migrateJob.ExpectRunTime != nil {
		_ = d.Set("expect_run_time", migrateJob.ExpectRunTime)
	}

	return nil
}

func handleModifyMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	configMigrationJobRequest := dts.NewModifyMigrationJobRequest()
	configMigrationJobRequest.JobId = helper.String(jobId)

	if v, ok := d.GetOk("run_mode"); ok {
		configMigrationJobRequest.RunMode = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "migrate_option"); ok {
		migrateOption := dts.MigrateOption{}
		if databaseTableMap, ok := helper.InterfaceToMap(dMap, "database_table"); ok {
			databaseTableObject := dts.DatabaseTableObject{}
			if v, ok := databaseTableMap["object_mode"]; ok && v.(string) != "" {
				databaseTableObject.ObjectMode = helper.String(v.(string))
			}
			if v, ok := databaseTableMap["databases"]; ok {
				for _, item := range v.([]interface{}) {
					databasesMap := item.(map[string]interface{})
					dBItem := dts.DBItem{}
					if v, ok := databasesMap["db_name"]; ok && v.(string) != "" {
						dBItem.DbName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_db_name"]; ok && v.(string) != "" {
						dBItem.NewDbName = helper.String(v.(string))
					}
					if v, ok := databasesMap["schema_name"]; ok && v.(string) != "" {
						dBItem.SchemaName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_schema_name"]; ok && v.(string) != "" {
						dBItem.NewSchemaName = helper.String(v.(string))
					}
					if v, ok := databasesMap["db_mode"]; ok && v.(string) != "" {
						dBItem.DBMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["schema_mode"]; ok && v.(string) != "" {
						dBItem.SchemaMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["table_mode"]; ok && v.(string) != "" {
						dBItem.TableMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["tables"]; ok {
						for _, item := range v.([]interface{}) {
							tablesMap := item.(map[string]interface{})
							tableItem := dts.TableItem{}
							if v, ok := tablesMap["table_name"]; ok && v.(string) != "" {
								tableItem.TableName = helper.String(v.(string))
							}
							if v, ok := tablesMap["new_table_name"]; ok && v.(string) != "" {
								tableItem.NewTableName = helper.String(v.(string))
							}
							if v, ok := tablesMap["tmp_tables"]; ok {
								tmpTablesSet := v.(*schema.Set).List()
								for i := range tmpTablesSet {
									tmpTables := tmpTablesSet[i].(string)
									tableItem.TmpTables = append(tableItem.TmpTables, &tmpTables)
								}
							}
							if v, ok := tablesMap["table_edit_mode"]; ok && v.(string) != "" {
								tableItem.TableEditMode = helper.String(v.(string))
							}
							dBItem.Tables = append(dBItem.Tables, &tableItem)
						}
					}
					if v, ok := databasesMap["view_mode"]; ok && v.(string) != "" {
						dBItem.ViewMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["views"]; ok {
						for _, item := range v.([]interface{}) {
							viewsMap := item.(map[string]interface{})
							viewItem := dts.ViewItem{}
							if v, ok := viewsMap["view_name"]; ok && v.(string) != "" {
								viewItem.ViewName = helper.String(v.(string))
							}
							if v, ok := viewsMap["new_view_name"]; ok && v.(string) != "" {
								viewItem.NewViewName = helper.String(v.(string))
							}
							dBItem.Views = append(dBItem.Views, &viewItem)
						}
					}
					if v, ok := databasesMap["role_mode"]; ok && v.(string) != "" {
						dBItem.RoleMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["roles"]; ok {
						for _, item := range v.([]interface{}) {
							rolesMap := item.(map[string]interface{})
							roleItem := dts.RoleItem{}
							if v, ok := rolesMap["role_name"]; ok && v.(string) != "" {
								roleItem.RoleName = helper.String(v.(string))
							}
							if v, ok := rolesMap["new_role_name"]; ok && v.(string) != "" {
								roleItem.NewRoleName = helper.String(v.(string))
							}
							dBItem.Roles = append(dBItem.Roles, &roleItem)
						}
					}
					if v, ok := databasesMap["function_mode"]; ok && v.(string) != "" {
						dBItem.FunctionMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["trigger_mode"]; ok && v.(string) != "" {
						dBItem.TriggerMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["event_mode"]; ok && v.(string) != "" {
						dBItem.EventMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["procedure_mode"]; ok && v.(string) != "" {
						dBItem.ProcedureMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["functions"]; ok {
						functionsSet := v.(*schema.Set).List()
						for i := range functionsSet {
							functions := functionsSet[i].(string)
							dBItem.Functions = append(dBItem.Functions, &functions)
						}
					}
					if v, ok := databasesMap["procedures"]; ok {
						proceduresSet := v.(*schema.Set).List()
						for i := range proceduresSet {
							procedures := proceduresSet[i].(string)
							dBItem.Procedures = append(dBItem.Procedures, &procedures)
						}
					}
					if v, ok := databasesMap["events"]; ok {
						eventsSet := v.(*schema.Set).List()
						for i := range eventsSet {
							events := eventsSet[i].(string)
							dBItem.Events = append(dBItem.Events, &events)
						}
					}
					if v, ok := databasesMap["triggers"]; ok {
						triggersSet := v.(*schema.Set).List()
						for i := range triggersSet {
							triggers := triggersSet[i].(string)
							dBItem.Triggers = append(dBItem.Triggers, &triggers)
						}
					}
					databaseTableObject.Databases = append(databaseTableObject.Databases, &dBItem)
				}
			}
			if v, ok := databaseTableMap["advanced_objects"]; ok {
				advancedObjectsSet := v.(*schema.Set).List()
				for i := range advancedObjectsSet {
					advancedObjects := advancedObjectsSet[i].(string)
					databaseTableObject.AdvancedObjects = append(databaseTableObject.AdvancedObjects, &advancedObjects)
				}
			}
			migrateOption.DatabaseTable = &databaseTableObject
		}
		if v, ok := dMap["migrate_type"]; ok && v.(string) != "" {
			migrateOption.MigrateType = helper.String(v.(string))
		}
		if consistencyMap, ok := helper.InterfaceToMap(dMap, "consistency"); ok {
			consistencyOption := dts.ConsistencyOption{}
			if v, ok := consistencyMap["mode"]; ok && v.(string) != "" {
				consistencyOption.Mode = helper.String(v.(string))
			}
			migrateOption.Consistency = &consistencyOption
		}
		if v, ok := dMap["is_migrate_account"]; ok {
			migrateOption.IsMigrateAccount = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_override_root"]; ok {
			migrateOption.IsOverrideRoot = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_dst_read_only"]; ok {
			migrateOption.IsDstReadOnly = helper.Bool(v.(bool))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				migrateOption.ExtraAttr = append(migrateOption.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.MigrateOption = &migrateOption
	}

	// if d.HasChange("src_info") {
	if dMap, ok := helper.InterfacesHeadMap(d, "src_info"); ok {
		dBEndpointInfo := dts.DBEndpointInfo{}
		if v, ok := dMap["region"]; ok && v.(string) != "" {
			dBEndpointInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["access_type"]; ok && v.(string) != "" {
			dBEndpointInfo.AccessType = helper.String(v.(string))
		}
		if v, ok := dMap["database_type"]; ok && v.(string) != "" {
			dBEndpointInfo.DatabaseType = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok && v.(string) != "" {
			dBEndpointInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["info"]; ok {
			for _, item := range v.([]interface{}) {
				infoMap := item.(map[string]interface{})
				dBInfo := dts.DBInfo{}
				if v, ok := infoMap["role"]; ok && v.(string) != "" {
					dBInfo.Role = helper.String(v.(string))
				}
				if v, ok := infoMap["db_kernel"]; ok && v.(string) != "" {
					dBInfo.DbKernel = helper.String(v.(string))
				}
				if v, ok := infoMap["host"]; ok && v.(string) != "" {
					dBInfo.Host = helper.String(v.(string))
				}
				if v, ok := infoMap["port"]; ok {
					dBInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := infoMap["user"]; ok && v.(string) != "" {
					dBInfo.User = helper.String(v.(string))
				}
				if v, ok := infoMap["password"]; ok && v.(string) != "" {
					dBInfo.Password = helper.String(v.(string))
				}
				if v, ok := infoMap["cvm_instance_id"]; ok && v.(string) != "" {
					dBInfo.CvmInstanceId = helper.String(v.(string))
				}
				if v, ok := infoMap["uniq_vpn_gw_id"]; ok && v.(string) != "" {
					dBInfo.UniqVpnGwId = helper.String(v.(string))
				}
				if v, ok := infoMap["uniq_dcg_id"]; ok && v.(string) != "" {
					dBInfo.UniqDcgId = helper.String(v.(string))
				}
				if v, ok := infoMap["instance_id"]; ok && v.(string) != "" {
					dBInfo.InstanceId = helper.String(v.(string))
				}
				if v, ok := infoMap["ccn_gw_id"]; ok && v.(string) != "" {
					dBInfo.CcnGwId = helper.String(v.(string))
				}
				if v, ok := infoMap["vpc_id"]; ok && v.(string) != "" {
					dBInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := infoMap["subnet_id"]; ok && v.(string) != "" {
					dBInfo.SubnetId = helper.String(v.(string))
				}
				if v, ok := infoMap["engine_version"]; ok && v.(string) != "" {
					dBInfo.EngineVersion = helper.String(v.(string))
				}
				if v, ok := infoMap["account"]; ok && v.(string) != "" {
					dBInfo.Account = helper.String(v.(string))
				}
				if v, ok := infoMap["account_role"]; ok && v.(string) != "" {
					dBInfo.AccountRole = helper.String(v.(string))
				}
				if v, ok := infoMap["account_mode"]; ok && v.(string) != "" {
					dBInfo.AccountMode = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_secret_id"]; ok && v.(string) != "" {
					dBInfo.TmpSecretId = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_secret_key"]; ok && v.(string) != "" {
					dBInfo.TmpSecretKey = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_token"]; ok && v.(string) != "" {
					dBInfo.TmpToken = helper.String(v.(string))
				}
				dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
			}
		}
		if v, ok := dMap["supplier"]; ok && v.(string) != "" {
			dBEndpointInfo.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.SrcInfo = &dBEndpointInfo
	}
	// }

	// if d.HasChange("dst_info") {
	if dMap, ok := helper.InterfacesHeadMap(d, "dst_info"); ok {
		dBEndpointInfo := dts.DBEndpointInfo{}
		if v, ok := dMap["region"]; ok && v.(string) != "" {
			dBEndpointInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["access_type"]; ok && v.(string) != "" {
			dBEndpointInfo.AccessType = helper.String(v.(string))
		}
		if v, ok := dMap["database_type"]; ok && v.(string) != "" {
			dBEndpointInfo.DatabaseType = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok && v.(string) != "" {
			dBEndpointInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["info"]; ok {
			for _, item := range v.([]interface{}) {
				infoMap := item.(map[string]interface{})
				dBInfo := dts.DBInfo{}
				if v, ok := infoMap["role"]; ok && v.(string) != "" {
					dBInfo.Role = helper.String(v.(string))
				}
				if v, ok := infoMap["db_kernel"]; ok && v.(string) != "" {
					dBInfo.DbKernel = helper.String(v.(string))
				}
				if v, ok := infoMap["host"]; ok && v.(string) != "" {
					dBInfo.Host = helper.String(v.(string))
				}
				if v, ok := infoMap["port"]; ok {
					dBInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := infoMap["user"]; ok && v.(string) != "" {
					dBInfo.User = helper.String(v.(string))
				}
				if v, ok := infoMap["password"]; ok && v.(string) != "" {
					dBInfo.Password = helper.String(v.(string))
				}
				if v, ok := infoMap["cvm_instance_id"]; ok && v.(string) != "" {
					dBInfo.CvmInstanceId = helper.String(v.(string))
				}
				if v, ok := infoMap["uniq_vpn_gw_id"]; ok && v.(string) != "" {
					dBInfo.UniqVpnGwId = helper.String(v.(string))
				}
				if v, ok := infoMap["uniq_dcg_id"]; ok && v.(string) != "" {
					dBInfo.UniqDcgId = helper.String(v.(string))
				}
				if v, ok := infoMap["instance_id"]; ok && v.(string) != "" {
					dBInfo.InstanceId = helper.String(v.(string))
				}
				if v, ok := infoMap["ccn_gw_id"]; ok && v.(string) != "" {
					dBInfo.CcnGwId = helper.String(v.(string))
				}
				if v, ok := infoMap["vpc_id"]; ok && v.(string) != "" {
					dBInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := infoMap["subnet_id"]; ok && v.(string) != "" {
					dBInfo.SubnetId = helper.String(v.(string))
				}
				if v, ok := infoMap["engine_version"]; ok && v.(string) != "" {
					dBInfo.EngineVersion = helper.String(v.(string))
				}
				if v, ok := infoMap["account"]; ok && v.(string) != "" {
					dBInfo.Account = helper.String(v.(string))
				}
				if v, ok := infoMap["account_role"]; ok && v.(string) != "" {
					dBInfo.AccountRole = helper.String(v.(string))
				}
				if v, ok := infoMap["account_mode"]; ok && v.(string) != "" {
					dBInfo.AccountMode = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_secret_id"]; ok && v.(string) != "" {
					dBInfo.TmpSecretId = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_secret_key"]; ok && v.(string) != "" {
					dBInfo.TmpSecretKey = helper.String(v.(string))
				}
				if v, ok := infoMap["tmp_token"]; ok && v.(string) != "" {
					dBInfo.TmpToken = helper.String(v.(string))
				}
				dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
			}
		}
		if v, ok := dMap["supplier"]; ok && v.(string) != "" {
			dBEndpointInfo.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.DstInfo = &dBEndpointInfo
	}
	// }

	// if d.HasChange("job_name") {
	if v, ok := d.GetOk("job_name"); ok && v.(string) != "" {
		configMigrationJobRequest.JobName = helper.String(v.(string))
	}
	// }

	// if d.HasChange("expect_run_time") {
	if v, ok := d.GetOk("expect_run_time"); ok && v.(string) != "" {
		configMigrationJobRequest.ExpectRunTime = helper.String(v.(string))
	}
	// }

	// if d.HasChange("auto_retry_time_range_minutes") {
	if v, _ := d.GetOk("auto_retry_time_range_minutes"); v != nil {
		configMigrationJobRequest.AutoRetryTimeRangeMinutes = helper.IntInt64(v.(int))
	}
	// }

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().ModifyMigrationJob(configMigrationJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, configMigrationJobRequest.GetAction(), configMigrationJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}
	return nil
}

func handleCheckMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	checkMigrateJobRequest := dts.NewCreateMigrateCheckJobRequest()
	checkMigrateJobRequest.JobId = helper.String(jobId)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().CreateMigrateCheckJob(checkMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, checkMigrateJobRequest.GetAction(), checkMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s check dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func handleResumeMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	resumeMigrateJobRequest := dts.NewResumeMigrateJobRequest()
	resumeMigrateJobRequest.JobId = helper.String(jobId)
	service := DtsService{client: tcClient}

	if d.HasChange("resume_option") {
		if v, ok := d.GetOk("resume_option"); ok {
			resumeMigrateJobRequest.ResumeOption = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().ResumeMigrateJob(resumeMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, resumeMigrateJobRequest.GetAction(), resumeMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s resume dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"readyComplete", "success", "failed"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func handleCompleteMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	completeMigrateJobRequest := dts.NewCompleteMigrateJobRequest()
	completeMigrateJobRequest.JobId = helper.String(jobId)
	service := DtsService{client: tcClient}

	if d.HasChange("complete_mode") {
		if v, ok := d.GetOk("complete_mode"); ok {
			completeMigrateJobRequest.CompleteMode = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().CompleteMigrateJob(completeMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, completeMigrateJobRequest.GetAction(), completeMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s complete dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"success", "error", "failed"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func handleCompareMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	startCompareRequest := dts.NewStartCompareRequest()
	startCompareRequest.JobId = helper.String(jobId)

	if d.HasChange("compare_task_id") {
		if v, ok := d.GetOk("compare_task_id"); ok {
			startCompareRequest.CompareTaskId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().StartCompare(startCompareRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, startCompareRequest.GetAction(), startCompareRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s compare dts migrate job failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func handleStopMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	stopMigrateJobRequest := dts.NewStopMigrateJobRequest()
	stopMigrateJobRequest.JobId = helper.String(jobId)
	service := DtsService{client: tcClient}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().StopMigrateJob(stopMigrateJobRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, stopMigrateJobRequest.GetAction(), stopMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s stop dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	conf := BuildStateChangeConf([]string{}, []string{"canceled"}, 3*readRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}

func resourceTencentCloudDtsMigrateJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
