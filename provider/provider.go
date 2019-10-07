package provider

import (
	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/kyma-incubator/terraform-provider-gardener/client"
	"github.com/kyma-incubator/terraform-provider-gardener/shoot"
)

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kube_path": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KUBECONFIG", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gardener_gcp_shoot":   shoot.GCPShoot(),
			"gardener_aws_shoot":   shoot.AWSShoot(),
			"gardener_azure_shoot": shoot.AzureShoot(),
		},
		ConfigureFunc: providerConfigure,
	}
}
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &client.Config{
		KubePath:       d.Get("kube_path").(string),
	}
	return client.New(config)
}
