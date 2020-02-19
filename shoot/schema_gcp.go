package shoot

import "github.com/hashicorp/terraform/helper/schema"

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
