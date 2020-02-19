package shoot

import "github.com/hashicorp/terraform/helper/schema"

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
