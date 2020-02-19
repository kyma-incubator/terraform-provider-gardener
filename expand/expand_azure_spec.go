package expand

import (
	"encoding/json"
	azAlpha1 "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func getAzControlPlaneConfig() *corev1beta1.ProviderConfig {
	azConfig := `
      apiVersion: azure.provider.extensions.gardener.cloud/v1alpha1
      kind: ControlPlaneConfig`
	obj := corev1beta1.ProviderConfig{}
	obj.Raw = []byte(azConfig)
	return &obj
}

func getAzureConfig(az []interface{}) *corev1beta1.ProviderConfig {
	azConfigObj := azAlpha1.InfrastructureConfig{}
	obj := corev1beta1.ProviderConfig{}
	if len(az) == 0 && az[0] == nil {
		return &obj
	}
	in := az[0].(map[string]interface{})

	azConfigObj.APIVersion = "azure.provider.extensions.gardener.cloud/v1alpha1"
	azConfigObj.Kind = "InfrastructureConfig"
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		azConfigObj.Networks = getNetworks(v)
	}
	obj.Raw, _ = json.Marshal(azConfigObj)
	return &obj
}

func getNetworks(networks []interface{}) azAlpha1.NetworkConfig {
	obj := azAlpha1.NetworkConfig{}
	if networks == nil {
		return obj
	}
	in := networks[0].(map[string]interface{})

	if v, ok := in["vnet"].([]interface{}); ok {
		obj.VNet = getVNET(v)
	}
	if v, ok := in["workers"].(string); ok {
		obj.Workers = v
	}
	if v, ok := in["service_endpoints"].(*schema.Set); ok {
		obj.ServiceEndpoints = expandSet(v)
	}

	return obj
}

func getVNET(vnet []interface{}) azAlpha1.VNet {
	obj := azAlpha1.VNet{}

	if len(vnet) == 0 && vnet[0] == nil {
		return obj
	}
	in := vnet[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = &v
	}
	if v, ok := in["resource_group"].(string); ok && len(v) > 0 {
		obj.ResourceGroup = &v
	}

	if v, ok := in["cidr"].(string); ok && len(v) > 0 {
		obj.CIDR = &v
	}
	return obj
}
