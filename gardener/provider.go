package gardener

import (
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
)

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROFILE", ""),
			},
			"kube_path": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KUBECONFIG", ""),
			},
			"aws_secret_binding": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"azure_secret_binding": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"gcp_secret_binding": {
				Type:         schema.TypeString,
				Optional:     true,
				InputDefault: "",
			},
			"openstack_secret_binding": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"alicloud_secret_binding": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gardener_gcp_shoot": resourceGCPShoot(),
		},
		ConfigureFunc: providerConfigure,
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
