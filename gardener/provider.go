package gardener

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"kube_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_secret_binding": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"azure_secret_binding": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"gcp_secret_binding": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"openstack_secret_binding": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"alicloud_secret_binding": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gcp_shoot": resourceGCPShoot(),
		},
	}
}
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	SecretBindings := &Bindings{
		AwsSecretBinding:       d.Get("aws_secret_binding").(string),
		GcpSecretBinding:       d.Get("gcp_secret_binding").(string),
		AzureSecretBinding:     d.Get("azure_secret_binding").(string),
		OpenStackSecretBinding: d.Get("openstack_secret_binding").(string),
		AliCloudSecretBinding:  d.Get("alicloud_secret_binding").(string),
	}
	config := Config{
		Profile:        d.Get("profile").(string),
		KubePath:       d.Get("kube_path").(string),
		SecretBindings: SecretBindings,
	}
	return config.Client()
}
