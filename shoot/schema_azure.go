package shoot

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func azureResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:        schema.TypeList,
				Description: "Networks is the network configuration (VNet, subnets, etc.).",
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
