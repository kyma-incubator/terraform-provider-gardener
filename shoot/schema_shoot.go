package shoot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/runtime"
)

func workerKubernetes() schema.Resource{
 return  schema.Resource{}
}

func workerConfig()  *schema.Resource{
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
				Type: schema.TypeList,
				Description: "Kubernetes contains configuration for Kubernetes components related to this worker pool.",
				Optional: true,
				Elem:workerKubernetes(),
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
				Type:        schema.TypeString,
				Description: "MachineType is the machine type of the worker group.",
				Required:    true,
			},
			"auto_scaler_min": {
				Type:        schema.TypeInt,
				Description: "AutoScalerMin is the minimum number of VMs to create.",
				Required:    true,
			},
			"auto_scaler_max": {
				Type:        schema.TypeInt,
				Description: "AutoScalerMax is the maximum number of VMs to create.",
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
						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"volume_type": {
				Type:        schema.TypeString,
				Description: "VolumeType is the type of the root volumes.",
				Required:    true,
			},
			"volume_size": {
				Type:        schema.TypeString,
				Description: "VolumeSize is the size of the root volume.",
				Required:    true,
			},
			"zones": {
				Type:        schema.TypeSet,
				Description: "Zones is a list of availability zones to deploy the Shoot cluster to.",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
		},
	}
}

func providerResource() schema.Resource {
	return schema.Resource{
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
						"ProviderConfig": {
							Type:     schema.TypeList,
							Elem:runtime.RawExtension{},
							Required: true,
						},
					},
				},
			},
			//"infrastructure_config": {
			//	Type:        schema.TypeList,
			//	Description: "InfrastructureConfig contains the provider-specific infrastructure config blob.",
			//	Optional:    true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"ProviderConfig": {
			//				Type:     schema.TypeList,
			//				Elem:runtime.RawExtension{},
			//				Required: true,
			//			},
			//		},
			//	},
			//},
			"workers": {
				Type: schema.TypeList,
				Description: "Workers is a list of worker groups.",
				Required: true,
				Elem: workerConfig(),
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
			"provider": {
				Type:        schema.TypeString,
				Description: "Provider is the DNS provider type for the Shoot.",
				Optional:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "Domain is the external available domain of the Shoot cluster.",
				Optional:    true,
			},
			"secret_name": {
				Type:        schema.TypeString,
				Description: "SecretName is a name of a secret containing credentials for the stated domain and the provider.",
				Optional:    true,
			},
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
			"cloud_controller_manager": {
				Type:        schema.TypeList,
				Description: "CloudControllerManager contains configuration settings for the cloud-controller-manager.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"feature_gates": {
							Type:        schema.TypeMap,
							Description: "FeatureGates contains information about enabled feature gates.",
							Optional:    true,
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
							Type:        schema.TypeList,
							Description: "List of people who receiving alerts for this shoots",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
				"extensions": {
					Type: schema.TypeList,
					Description: "Extensions contain type and provider information for Shoot extensions.",
					Optional: true,
				},
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
					Type: schema.TypeList,
					Description: "Provider contains all provider-specific and provider-relevant information.",
					Required: true,
					Elem: providerResource(),
				},
				"purpose": {
					Type: schema.TypeString,
					Description: "Purpose is the purpose class for this cluster.",
					Optional: true,
				},
				"region": {
					Type:        schema.TypeString,
					Description: "Region is a name of a cloud provider region.",
					Required:    true,
				},
				"secretBindingName": {
					Type: schema.TypeString,
					Description: "SecretBindingName is the name of the a SecretBinding that has a reference to the provider secret. The credentials inside the provider secret will be used to create the shoot in the respective account",
					Required: true,
				},
				"seedName": {
					Type: schema.TypeString,
					Description: "SeedName is the name of the seed cluster that runs the control plane of the Shoot.",
					Optional: true,
				},
				//"cloud": {
				//	Type:        schema.TypeList,
				//	Description: "Cloud contains information about the cloud environment and their specific settings.",
				//	Required:    true,
				//	MaxItems:    1,
				//	Elem:        cloudResource(),
				//},
			},
		},
	}
}

func cloudResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			//"profile": {
			//	Type:        schema.TypeString,
			//	Description: "Profile is a name of a CloudProfile object.",
			//	Required:    true,
			//},
			//"region": {
			//	Type:        schema.TypeString,
			//	Description: "Region is a name of a cloud provider region.",
			//	Required:    true,
			//},
			"secret_binding_ref": {
				Type:        schema.TypeList,
				Description: "SecretBindingRef is a reference to a SecretBinding object.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name of the secret.",
							Required:    true,
						},
					},
				},
			},
			"aws": {
				Type:        schema.TypeList,
				Description: "AWS contains the Shoot specification for the Amazon Web Services cloud.",
				Optional:    true,
				MaxItems:    1,
				Elem:        cloudResourceAWS(),
			},
			"gcp": {
				Type:        schema.TypeList,
				Description: "GCP contains the Shoot specification for the Google Cloud Platform.",
				Optional:    true,
				MaxItems:    1,
				Elem:        cloudResourceGCP(),
			},
			"azure": {
				Type:        schema.TypeList,
				Description: "Azure contains the Shoot specification for the Azure Cloud Platform.",
				Optional:    true,
				MaxItems:    1,
				Elem:        cloudResourceAzure(),
			},
		},
	}
}
