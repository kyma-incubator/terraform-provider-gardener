package shoot

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func workerKubernetes() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kubelet": {
				Type:        schema.TypeList,
				Description: "Kubelet contains configuration settings for the kubelet.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pod_pids_limit": {
							Type:        schema.TypeInt,
							Description: "PodPIDsLimit is the maximum number of process IDs per pod allowed by the kubelet.",
							Optional:    true,
						},
						"cpu_cfs_quota": {
							Type:        schema.TypeBool,
							Description: "CPUCFSQuota allows you to disable/enable CPU throttling for Pods.",
							Optional:    true,
						},
						"cpu_manager_policy": {
							Type:        schema.TypeString,
							Description: "CPUManagerPolicy allows to set alternative CPU management policies (default: none).",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func workerConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"annotations": {
				Type:        schema.TypeMap,
				Description: "Annotations is a map of key/value pairs for annotations for all the `Node` objects in this worker pool.",
				Optional:    true,
			},
			"cabundle": {
				Type:        schema.TypeString,
				Description: "caBundle configuration",
				Optional:    true,
			},
			"kubernetes": {
				Type:        schema.TypeList,
				Description: "Kubernetes contains configuration for Kubernetes components related to this worker pool.",
				Optional:    true,
				Elem:        workerKubernetes(),
			},
			"labels": {
				Type:        schema.TypeMap,
				Description: "Labels is a map of key/value pairs for labels for all the `Node` objects in this worker pool.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name is the name of the worker group.",
				Required:    true,
			},
			"machine": {
				Type:        schema.TypeList,
				Description: "MachineType is the machine type of the worker group.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "Type is the machine type of the worker group.",
							Required:    true,
						},
						"image": {
							Type:        schema.TypeList,
							Description: "Image holds information about the machine image to use for all nodes of this pool. It will default to the latest version of the first image stated in the referenced CloudProfile if no value has been provided.",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "VolumeSize is the size of the root volume.",
										Required:    true,
									},
									"version": {
										Type:        schema.TypeString,
										Description: "Version is the version of the shoot's image.",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"minimum": {
				Type:        schema.TypeInt,
				Description: "Minimum is the minimum number of VMs to create.",
				Required:    true,
			},
			"maximum": {
				Type:        schema.TypeInt,
				Description: "Maximum is the maximum number of VMs to create.",
				Required:    true,
			},
			"max_surge": {
				Type:        schema.TypeInt,
				Description: "MaxSurge is maximum number of VMs that are created during an update.",
				Optional:    true,
			},
			"max_unavailable": {
				Type:        schema.TypeInt,
				Description: "MaxUnavailable is the maximum number of VMs that can be unavailable during an update.",
				Optional:    true,
			},
			"taints": {
				Type:        schema.TypeList,
				Description: "Taints is a list of taints for all the `Node` objects in this worker pool.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"volume": {
				Type:        schema.TypeList,
				Description: "Volume contains information about the volume type and size.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "VolumeType is the type of the root volumes.",
							Required:    true,
						},
						"size": {
							Type:        schema.TypeString,
							Description: "VolumeSize is the size of the root volume.",
							Required:    true,
						},
					},
				},
			},
			"zones": {
				Type:        schema.TypeSet,
				Description: "Zones is a list of availability zones to deploy the Shoot cluster to.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
		},
	}
}

func providerResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "Type is the type of the provider.",
				Required:    true,
			},
			"control_plane_config": {
				Type:        schema.TypeList,
				Description: "ControlPlaneConfig contains the provider-specific control plane config blob.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gcp": {
							Type:        schema.TypeList,
							Description: "GCP contains the Shoot specification for Google Cloud Platform.",
							Optional:    true,
							MaxItems:    1,
							Elem:        gcpControlPlaneResource(),
						},
					},
				},
			},
			"infrastructure_config": {
				Type:        schema.TypeList,
				Description: "InfrastructureConfig contains the provider-specific infrastructure config blob.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"azure": {
							Type:        schema.TypeList,
							Description: "Azure contains the Shoot specification for the Azure Cloud Platform.",
							Optional:    true,
							MaxItems:    1,
							Elem:        azureResource(),
						},
						"aws": {
							Type:        schema.TypeList,
							Description: "AWS contains the Shoot specification for the AWS Cloud Platform.",
							Optional:    true,
							MaxItems:    1,
							Elem:        awsResource(),
						},
						"gcp": {
							Type:        schema.TypeList,
							Description: "GCP contains the Shoot specification for Google Cloud Platform.",
							Optional:    true,
							MaxItems:    1,
							Elem:        gcpResource(),
						},
					},
				},
			},
			"worker": {
				Type:             schema.TypeList,
				Description:      "Workers is a list of worker groups.",
				Required:         true,
				Elem:             workerConfig(),
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
			},
		},
	}
}

func gcpControlPlaneResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Description: "Zone is the GCP zone.",
				Required:    true,
			},
		},
	}
}

func gcpResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:        schema.TypeList,
				Description: "Networks is the network configuration (VPC, subnets, etc.)",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:        schema.TypeList,
							Description: "VPC indicates whether to use an existing VPC or create a new one.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "Name is the VPC name.",
										Optional:    true,
									},
									"cloud_router": {
										Type:        schema.TypeList,
										Description: "CloudRouter indicates whether to use an existing CloudRouter or create a new one.",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Description: "Name is the CloudRouter name.",
													Optional:    true,
												},
											},
										},
									},
								},
							},
						},
						"cloud_nat": {
							Type:        schema.TypeList,
							Description: "CloudNAT contains configuration about the the CloudNAT configuration",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_ports_per_vm": {
										Type:        schema.TypeInt,
										Description: "MinPortsPerVM is the minimum number of ports allocated to a VM in the NAT config. The default value is 2048 ports.",
										Optional:    true,
									},
								},
							},
						},
						"internal": {
							Type:        schema.TypeString,
							Description: "Internal is a private subnet (used for internal load balancers).",
							Optional:    true,
						},
						"workers": {
							Type:        schema.TypeString,
							Description: "Workers is the worker subnet range to create (used for the VMs).",
							Required:    true,
						},
						"flow_logs": {
							Type:        schema.TypeList,
							Description: "FlowLogs contains the flow log configuration for the subnet.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"aggregation_interval": {
										Type:        schema.TypeString,
										Description: "AggregationInterval for collecting flow logs.",
										Optional:    true,
									},
									"flow_sampling": {
										Type:        schema.TypeFloat,
										Description: "FlowSampling sets the sampling rate of VPC flow logs within the subnetwork where 1.0 means all collected logs are reported and 0.0 means no logs are reported.",
										Optional:    true,
									},
									"metadata": {
										Type:        schema.TypeString,
										Description: "Metadata configures whether metadata fields should be added to the reported VPC flow logs.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func azureResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:        schema.TypeList,
				Description: "NetworkConfig holds information about the Kubernetes and infrastructure networks.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workers": {
							Type:        schema.TypeString,
							Description: "Workers is the worker subnet range to create (used for the VMs).",
							Required:    true,
						},
						"service_endpoints": {
							Type:        schema.TypeSet,
							Description: "List of Azure service endpoints connect to the created VNet.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
						"vnet": {
							Type:        schema.TypeList,
							Description: "VNet indicates whether to use an existing VNet or create a new one.",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "Name is the VNet name.",
										Optional:    true,
									},
									"cidr": {
										Type:        schema.TypeString,
										Description: "CIDR is the VNet CIDR.",
										Optional:    true,
									},
									"resource_group": {
										Type:        schema.TypeString,
										Description: "ResourceGroup is the resource group where the existing vNet belongs to.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func awsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enableecraccess": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"networks": {
				Type:        schema.TypeList,
				Description: "NetworkConfig holds information about the Kubernetes and infrastructure networks.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc": {
							Type:        schema.TypeList,
							Description: "VPC ID or CIDR for aws",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "ID of the VPC.",
										Optional:    true,
									},
									"cidr": {
										Type:        schema.TypeString,
										Description: "CIDR is the VPC CIDR.",
										Optional:    true,
									},
								},
							},
						},
						"zones": {
							Type:        schema.TypeSet,
							Description: "List of zones.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Description: "Name is the zone name.",
										Optional:    true,
									},
									"internal": {
										Type:        schema.TypeString,
										Description: "internal CIDR",
										Optional:    true,
									},
									"public": {
										Type:        schema.TypeString,
										Description: "public cidr",
										Optional:    true,
									},
									"workers": {
										Type:        schema.TypeString,
										Description: "worker cidr",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func addonsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kubernetes_dashboard": {
				Type:        schema.TypeList,
				Description: "Kubernetes dashboard holds configuration settings for the kubernetes dashboard addon.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"authentication_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"nginx_ingress": {
				Type:        schema.TypeList,
				Description: "NginxIngress holds configuration settings for the nginx-ingress addon.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dNSResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// "provider": {
			// 	Type:        schema.TypeString,
			// 	Description: "Provider is the DNS provider type for the Shoot.",
			// 	Optional:    true,
			// },
			"domain": {
				Type:        schema.TypeString,
				Description: "Domain is the external available domain of the Shoot cluster.",
				Optional:    true,
			},
			// "secret_name": {
			// 	Type:        schema.TypeString,
			// 	Description: "SecretName is a name of a secret containing credentials for the stated domain and the provider.",
			// 	Optional:    true,
			// },
		},
	}
}

func hibernationResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Enabled is true if Shoot is hibernated, false otherwise.",
				Required:    true,
			},
			"schedules": {
				Type:        schema.TypeList,
				Description: "Schedules determine the hibernation schedules.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Description: "Start is a Cron spec at which time a Shoot will be hibernated.",
							Optional:    true,
						},
						"end": {
							Type:        schema.TypeString,
							Description: "End is a Cron spec at which time a Shoot will be woken up.",
							Optional:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "Location is the time location in which both start and and shall be evaluated.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func kubernetesResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"allow_privileged_containers": {
				Type:        schema.TypeBool,
				Description: "AllowPrivilegedContainers indicates whether privileged containers are allowed in the Shoot.",
				Optional:    true,
			},
			"kube_api_server": {
				Type:        schema.TypeList,
				Description: "KubeAPIServer contains configuration settings for the kube-apiserver.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_basic_authentication": {
							Type:        schema.TypeBool,
							Description: "enable basic authentication flag.",
							Optional:    true,
						},
						"oidc_config": {
							Type:             schema.TypeList,
							Description:      "interface for adding oidc_config in kube api server section",
							MaxItems:         1,
							Optional:         true,
							DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ca_bundle": {
										Type:        schema.TypeString,
										Description: "ca_bundle for oidc config in kube api server section",
										Optional:    true,
									},
									"client_id": {
										Type:        schema.TypeString,
										Description: "client_id for oidc config in kube api server section",
										Optional:    true,
									},
									"groups_claim": {
										Type:        schema.TypeString,
										Description: "groups_claim for oidc config in kube api server section",
										Optional:    true,
									},
									"groups_prefix": {
										Type:        schema.TypeString,
										Description: "groups_prefix for oidc config in kube api server section",
										Optional:    true,
									},
									"issuer_url": {
										Type:        schema.TypeString,
										Description: "issuer_url for oidc config in kube api server section",
										Optional:    true,
									},
									"required_claims": {
										Type:        schema.TypeString,
										Description: "required_claims for oidc config in kube api server section",
										Optional:    true,
									},
									"signing_algs": {
										Type:        schema.TypeString,
										Description: "signing_algs for oidc config in kube api server section",
										Optional:    true,
									},
									"username_claim": {
										Type:        schema.TypeString,
										Description: "username_claim for oidc config in kube api server section",
										Optional:    true,
									},
									"username_prefix": {
										Type:        schema.TypeString,
										Description: "username_prefix for oidc config in kube api server section",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
			"kube_controller_manager": {
				Type:        schema.TypeList,
				Description: "KubeControllerManager contains configuration settings for the kube-controller-manager.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_cidr_mask_size": {
							Type:        schema.TypeInt,
							Description: "Size of the node mask.",
							Optional:    true,
						},
					},
				},
			},
			"kube_proxy": {
				Type:        schema.TypeList,
				Description: "KubeProxy contains configuration settings for the kube-proxy.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Description: "Mode specifies which proxy mode to use. defaults to IPTables.",
							Optional:    true,

							Default: "IPTables",
						},
					},
				},
			},
			"kubelet": {
				Type:        schema.TypeList,
				Description: "Kubelet contains configuration settings for the kubelet.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"feature_gates": {
							Type:        schema.TypeMap,
							Description: "FeatureGates contains information about enabled feature gates.",
							Optional:    true,
						},
						"pod_pids_limit": {
							Type:        schema.TypeInt,
							Description: "PodPIDsLimit is the maximum number of process IDs per pod allowed by the kubelet.",
							Optional:    true,
						},
						"cpu_cfs_quota": {
							Type:        schema.TypeBool,
							Description: "CPUCFSQuota allows you to disable/enable CPU throttling for Pods.",
							Optional:    true,
						},
						"cpu_manager_policy": {
							Type:        schema.TypeString,
							Description: "CPUManagerPolicy allows to set alternative CPU management policies (default: none).",
							Optional:    true,
						},
					},
				},
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Version is the semantic Kubernetes version to use for the Shoot cluster.",
				Required:    true,
			},
			"cluster_autoscaler": {
				Type:        schema.TypeList,
				Description: "ClusterAutoscaler contains the configration flags for the Kubernetes cluster autoscaler.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scale_down_utilization_threshold": {
							Type:        schema.TypeFloat,
							Description: "ScaleDownUtilizationThreshold defines the threshold in % under which a node is being removed.",
							Optional:    true,
						},
						"scale_down_unneeded_time": {
							Type:        schema.TypeString,
							Description: "ScaleDownUnneededTime defines how long a node should be unneeded before it is eligible for scale down (default: 10 mins).",
							Optional:    true,
						},
						"scale_down_delay_after_add": {
							Type:        schema.TypeString,
							Description: "ScaleDownDelayAfterAdd defines how long after scale up that scale down evaluation resumes (default: 10 mins).",
							Optional:    true,
						},
						"scale_down_delay_after_failure": {
							Type:        schema.TypeString,
							Description: "ScaleDownDelayAfterFailure how long after scale down failure that scale down evaluation resumes (default: 3 mins).",
							Optional:    true,
						},
						"scale_down_delay_after_delete": {
							Type:        schema.TypeString,
							Description: "ScaleDownDelayAfterDelete how long after node deletion that scale down evaluation resumes, defaults to scanInterval (defaults to ScanInterval).",
							Optional:    true,
						},
						"scan_interval": {
							Type:        schema.TypeString,
							Description: "ScanInterval how often cluster is reevaluated for scale up or down (default: 10 secs).",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func maintenanceResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_update": {
				Type:             schema.TypeList,
				Description:      "AutoUpdate contains information about which constraints should be automatically updated.",
				MaxItems:         1,
				Optional:         true,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kubernetes_version": {
							Type:        schema.TypeBool,
							Description: "KubernetesVersion indicates whether the patch Kubernetes version may be automatically updated.",
							Optional:    true,
							Default:     true,
						},
						"machine_image_version": {
							Type:        schema.TypeBool,
							Description: "machineImageVersion indicates whether the machine image version may be automatically updated.",
							Optional:    true,
							Default:     true,
						},
					},
				},
			},
			"time_window": {
				Type:             schema.TypeList,
				Description:      "TimeWindow contains information about the time window for maintenance operations.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin": {
							Type:             schema.TypeString,
							Description:      "Begin is the beginning of the time window in the format HHMMSS+ZONE, e.g. '220000+0100'.",
							Optional:         true,
							DiffSuppressFunc: suppressEmptyNewValue,
						},
						"end": {
							Type:             schema.TypeString,
							Description:      "End is the end of the time window in the format HHMMSS+ZONE, e.g. '220000+0100'.",
							Optional:         true,
							DiffSuppressFunc: suppressEmptyNewValue,
						},
					},
				},
			},
		},
	}
}

func monitoringResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alerting": {
				Type:             schema.TypeList,
				Description:      "Alert configuration to send notification to email lists.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"emailreceivers": {
							Type:             schema.TypeSet,
							Description:      "List of people who receiving alerts for this shoots",
							Optional:         true,
							Elem:             &schema.Schema{Type: schema.TypeString},
							Set:              schema.HashString,
							DiffSuppressFunc: suppressEmptyNewValue,
						},
					},
				},
			},
		},
	}
}

func shootSpecSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "ShootSpec is the specification of a Shoot. (see https://github.com/gardener/gardener/blob/master/pkg/apis/garden/v1beta1/types.go#L609)",
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"addons": {
					Type:        schema.TypeList,
					Description: "Addons contains information about enabled/disabled addons and their configuration.",
					Optional:    true,
					MaxItems:    1,
					Elem:        addonsResource(),
				},
				"cloud_profile_name": {
					Type:        schema.TypeString,
					Description: "Profile is a name of a CloudProfile object.",
					Required:    true,
				},
				"dns": {
					Type:        schema.TypeList,
					Description: "DNS contains information about the DNS settings of the Shoot.",
					Optional:    true,
					MaxItems:    1,
					Elem:        dNSResource(),
				},
				//"extensions": {
				//	Type: schema.TypeList,
				//	Description: "Extensions contain type and provider information for Shoot extensions.",
				//	Optional: true,
				//},
				"hibernation": {
					Type:        schema.TypeList,
					Description: "Hibernation contains information whether the Shoot is suspended or not.",
					Optional:    true,
					MaxItems:    1,
					Elem:        hibernationResource(),
				},
				"kubernetes": {
					Type:        schema.TypeList,
					Description: "Kubernetes contains the version and configuration settings of the control plane components.",
					Required:    true,
					MaxItems:    1,
					Elem:        kubernetesResource(),
				},
				"networking": {
					Type:        schema.TypeList,
					Description: "Networking contains information about cluster networking such as CNI Plugin type, CIDRs, ...etc.",
					Required:    true,
					Elem:        networkingResource(),
				},
				"maintenance": {
					Type:             schema.TypeList,
					Description:      "Maintenance contains information about the time window for maintenance operations and which operations should be performed.",
					Optional:         true,
					MaxItems:         1,
					Elem:             maintenanceResource(),
					DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				},
				"monitoring": {
					Type:             schema.TypeList,
					Description:      "Alert configuration to send notification to email lists.",
					Optional:         true,
					MaxItems:         1,
					Elem:             monitoringResource(),
					DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				},
				"provider": {
					Type:             schema.TypeList,
					Description:      "Provider contains all provider-specific and provider-relevant information.",
					Required:         true,
					Elem:             providerResource(),
					DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				},
				"purpose": {
					Type:        schema.TypeString,
					Description: "Purpose is the purpose class for this cluster.",
					Optional:    true,
				},
				"region": {
					Type:        schema.TypeString,
					Description: "Region is a name of a cloud provider region.",
					Required:    true,
				},
				"secret_binding_name": {
					Type:        schema.TypeString,
					Description: "Secret binding name is the name of the a SecretBinding that has a reference to the provider secret. The credentials inside the provider secret will be used to create the shoot in the respective account",
					Required:    true,
				},
				"seed_name": {
					Type:        schema.TypeString,
					Description: "Seed name is the name of the seed cluster that runs the control plane of the Shoot.",
					Optional:    true,
				},
			},
		},
	}
}

func networkingResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "Type identifies the type of the networking plugin.",
				Required:    true,
			},
			"pods": {
				Type:        schema.TypeString,
				Description: "Pods is the CIDR of the pod network.",
				Optional:    true,
			},
			"nodes": {
				Type:        schema.TypeString,
				Description: "Nodes is the CIDR of the entire node network.",
				Optional:    true,
			},
			"services": {
				Type:        schema.TypeString,
				Description: "Services is the CIDR of the service network.",
				Optional:    true,
			},
		},
	}
}
