package flatten

import awsAlpha1 "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"

func flattenAws(in awsAlpha1.InfrastructureConfig) []interface{} {
	att := make(map[string]interface{})
	net := make(map[string]interface{})
	vpc := make(map[string]interface{})

	if in.EnableECRAccess != nil {
		att["enableecraccess"] = in.EnableECRAccess
	}
	if in.Networks.VPC.ID != nil {
		vpc["id"] = in.Networks.VPC.ID
	}
	if in.Networks.VPC.CIDR != nil {
		vpc["cidr"] = in.Networks.VPC.CIDR
	}
	net["vpc"] = []interface{}{vpc}

	if len(in.Networks.Zones) > 0 {
		zones := make([]map[string]interface{}, len(in.Networks.Zones))
		for i, v := range in.Networks.Zones {
			zone := map[string]interface{}{}
			if len(v.Name) > 0 {
				zone["name"] = v.Name
			}
			if len(v.Internal) > 0 {
				zone["internal"] = v.Internal
			}
			if len(v.Public) > 0 {
				zone["public"] = v.Public
			}
			if len(v.Workers) > 0 {
				zone["workers"] = v.Workers
			}
			zones[i] = zone
		}
		net["zones"] = zones
	}
	att["networks"] = []interface{}{net}

	return []interface{}{att}
}
