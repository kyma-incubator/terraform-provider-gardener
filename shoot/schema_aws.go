package shoot

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func cloudResourceAWS() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"machine_image": {
				Type:             schema.TypeList,
				Description:      "MachineImage holds information about the machine image to use for all workers.",
				Optional:         true,
				MaxItems:         1,
				DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name is the name of the image.",
							Optional:    true,
							Default:     "coreos",
						},
						"version": {
							Type:             schema.TypeString,
							Description:      "Version is the version of the image.",
							Optional:         true,
							DiffSuppressFunc: suppressEmptyNewValue,
						},
					},
				},
			},
			"networks": {
				Type:        schema.TypeList,
				Description: "Networks holds information about the Kubernetes and infrastructure networks.",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes": {
							Type:        schema.TypeString,
							Description: "Nodes is the CIDR of the node network.",
							Optional:    true,
						},
						"pods": {
							Type:        schema.TypeString,
							Description: "Pods is the CIDR of the pod network.",
							Optional:    true,
						},
						"services": {
							Type:        schema.TypeString,
							Description: "Services is the CIDR of the service network.",
							Optional:    true,
						},
						"vpc": {
							Type:        schema.TypeList,
							Description: "VPC indicates whether to use an existing VPC or create a new one.",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "ID is the AWS VPC id of an existing VPC.",
										Optional:    true,
									},
									"cidr": {
										Type:        schema.TypeString,
										Description: "CIDR is a CIDR range for a new VPC.",
										Optional:    true,
									},
								},
							},
						},
						"internal": {
							Type:        schema.TypeSet,
							Description: "Internal is a list of private subnets to create (used for internal load balancers).",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
						"public": {
							Type:        schema.TypeSet,
							Description: "Public is a list of public subnets to create (used for bastion and load balancers).",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
						"workers": {
							Type:        schema.TypeSet,
							Description: "Workers is a list of worker subnets (private) to create (used for the VMs).",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         schema.HashString,
						},
					},
				},
			},
			"worker": {
				Type:        schema.TypeList,
				Description: "Workers is a list of worker groups.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Name is the name of the worker group.",
							Required:    true,
						},
						"machine_type": {
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
						"annotations": {
							Type:        schema.TypeMap,
							Description: "Annotations is a map of key/value pairs for annotations for all the `Node` objects in this worker pool.",
							Optional:    true,
						},
						"labels": {
							Type:        schema.TypeMap,
							Description: "Labels is a map of key/value pairs for labels for all the `Node` objects in this worker pool.",
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
					},
				},
			},
			"zones": {
				Type:        schema.TypeSet,
				Description: "Zones is a list of availability zones to deploy the Shoot cluster to.",
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
			},
			"cabundle": {
				Type:        schema.TypeString,
				Description: "caBundle configuration",
				Optional:    true,
			},
		},
	}
}
