package flatten

import azAlpha1 "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"

func flattenAzure(in azAlpha1.InfrastructureConfig) []interface{} {
	att := make(map[string]interface{})

	net := make(map[string]interface{})
	if len(in.Networks.Workers) > 0 {
		net["workers"] = in.Networks.Workers
	}
	if len(in.Networks.ServiceEndpoints) > 0 {
		net["service_endpoints"] = in.Networks.ServiceEndpoints
	}
	vnet := make(map[string]interface{})
	if in.Networks.VNet.CIDR != nil {
		vnet["cidr"] = *in.Networks.VNet.CIDR
	}
	if in.Networks.VNet.Name != nil {
		vnet["name"] = *in.Networks.VNet.Name
	}
	if in.Networks.VNet.ResourceGroup != nil {
		vnet["resource_group"] = *in.Networks.VNet.ResourceGroup
	}
	net["vnet"] = []interface{}{vnet}
	att["networks"] = []interface{}{net}

	return []interface{}{att}
}
