package shoot

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func awsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enableecraccess": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"networks": {
				Type:        schema.TypeList,
				Description: "Networks is the network configuration (VNet, subnets, etc.).",
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
