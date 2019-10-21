package flatters

import (
	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func flattenCloud(in v1beta1.Cloud) []interface{} {
	att := make(map[string]interface{})

	if len(in.Profile) > 0 {
		att["profile"] = in.Profile
	}
	if len(in.Region) > 0 {
		att["region"] = in.Region
	}
	att["secret_binding_ref"] = flattenLocalObjectReference(&in.SecretBindingRef)
	if in.Seed != nil {
		att["seed"] = in.Seed
	}
	if in.AWS != nil {
		att["aws"] = flattenCloudAWS(in.AWS)
	}
	if in.GCP != nil {
		att["gcp"] = flattenCloudGCP(in.GCP)
	}
	if in.Azure != nil {
		att["azure"] = flattenCloudAzure(in.Azure)
	}
	return []interface{}{att}
}

func flattenCloudAWS(in *v1beta1.AWSCloud) []interface{} {
	att := make(map[string]interface{})

	if in.MachineImage != nil {
		image := make(map[string]interface{})
		if len(in.MachineImage.Name) > 0 {
			image["name"] = in.MachineImage.Name
		}
		if len(in.MachineImage.Version) > 0 {
			image["version"] = in.MachineImage.Version
		}
		att["machine_image"] = []interface{}{image}
	}
	networks := make(map[string]interface{})
	if in.Networks.Nodes != nil {
		networks["nodes"] = in.Networks.Nodes
	}
	if in.Networks.Pods != nil {
		networks["pods"] = in.Networks.Pods
	}
	if in.Networks.Services != nil {
		networks["services"] = in.Networks.Services
	}
	vpc := make(map[string]interface{})
	if in.Networks.VPC.ID != nil {
		vpc["id"] = in.Networks.VPC.ID
	}
	if in.Networks.VPC.CIDR != nil {
		vpc["cidr"] = in.Networks.VPC.CIDR
	}
	networks["vpc"] = []interface{}{vpc}
	if in.Networks.Internal != nil {
		networks["internal"] = newCIDRSet(schema.HashString, in.Networks.Internal)
	}
	if in.Networks.Public != nil {
		networks["public"] = newCIDRSet(schema.HashString, in.Networks.Public)
	}
	if in.Networks.Workers != nil {
		networks["workers"] = newCIDRSet(schema.HashString, in.Networks.Workers)
	}
	att["networks"] = []interface{}{networks}
	if len(in.Workers) > 0 {
		workers := make([]interface{}, len(in.Workers))
		for i, v := range in.Workers {
			m := map[string]interface{}{}

			if v.Name != "" {
				m["name"] = v.Name
			}
			if v.MachineType != "" {
				m["machine_type"] = v.MachineType
			}
			if v.AutoScalerMin != 0 {
				m["auto_scaler_min"] = v.AutoScalerMin
			}
			if v.AutoScalerMax != 0 {
				m["auto_scaler_max"] = v.AutoScalerMax
			}
			if v.MaxSurge != nil {
				m["max_surge"] = v.MaxSurge.IntValue()
			}
			if v.MaxUnavailable != nil {
				m["max_unavailable"] = v.MaxUnavailable.IntValue()
			}
			if len(v.Annotations) > 0 {
				m["annotations"] = v.Annotations
			}
			if len(v.Labels) > 0 {
				m["labels"] = v.Labels
			}
			if len(v.Taints) > 0 {
				taints := make([]interface{}, len(v.Taints))
				for i, v := range v.Taints {
					m := map[string]interface{}{}

					if v.Key != "" {
						m["key"] = v.Key
					}
					// if v.Operator != "" {
					// 	m["operator"] = v.Operator
					// }
					if v.Value != "" {
						m["value"] = v.Value
					}
					if v.Effect != "" {
						m["effect"] = v.Effect
					}
					taints[i] = m
				}
				m["taints"] = taints
			}
			if len(v.VolumeType) > 0 {
				m["volume_type"] = v.VolumeType
			}
			if len(v.VolumeSize) > 0 {
				m["volume_size"] = v.VolumeSize
			}
			workers[i] = m
		}
		att["workers"] = workers
		if in.Zones != nil {
			att["zones"] = newStringSet(schema.HashString, in.Zones)
		}
	}

	return []interface{}{att}
}

func flattenCloudGCP(in *v1beta1.GCPCloud) []interface{} {
	att := make(map[string]interface{})

	if in.MachineImage != nil {
		image := make(map[string]interface{})
		if len(in.MachineImage.Name) > 0 {
			image["name"] = in.MachineImage.Name
		}
		if len(in.MachineImage.Version) > 0 {
			image["version"] = in.MachineImage.Version
		}
		att["machine_image"] = []interface{}{image}
	}
	networks := make(map[string]interface{})
	if in.Networks.Nodes != nil {
		networks["nodes"] = in.Networks.Nodes
	}
	if in.Networks.Pods != nil {
		networks["pods"] = in.Networks.Pods
	}
	if in.Networks.Services != nil {
		networks["services"] = in.Networks.Services
	}
	if in.Networks.Workers != nil {
		networks["workers"] = newCIDRSet(schema.HashString, in.Networks.Workers)
	}
	att["networks"] = []interface{}{networks}
	if len(in.Workers) > 0 {
		workers := make([]interface{}, len(in.Workers))
		for i, v := range in.Workers {
			m := map[string]interface{}{}

			if v.Name != "" {
				m["name"] = v.Name
			}
			if v.MachineType != "" {
				m["machine_type"] = v.MachineType
			}
			if v.AutoScalerMin != 0 {
				m["auto_scaler_min"] = v.AutoScalerMin
			}
			if v.AutoScalerMax != 0 {
				m["auto_scaler_max"] = v.AutoScalerMax
			}
			if v.MaxSurge != nil {
				m["max_surge"] = v.MaxSurge.IntValue()
			}
			if v.MaxUnavailable != nil {
				m["max_unavailable"] = v.MaxUnavailable.IntValue()
			}
			if len(v.Annotations) > 0 {
				m["annotations"] = v.Annotations
			}
			if len(v.Labels) > 0 {
				m["labels"] = v.Labels
			}
			if len(v.Taints) > 0 {
				taints := make([]interface{}, len(v.Taints))
				for i, v := range v.Taints {
					m := map[string]interface{}{}

					if v.Key != "" {
						m["key"] = v.Key
					}
					// if v.Operator != "" {
					// 	m["operator"] = v.Operator
					// }
					if v.Value != "" {
						m["value"] = v.Value
					}
					if v.Effect != "" {
						m["effect"] = v.Effect
					}
					taints[i] = m
				}
				m["taints"] = taints
			}
			if len(v.VolumeType) > 0 {
				m["volume_type"] = v.VolumeType
			}
			if len(v.VolumeSize) > 0 {
				m["volume_size"] = v.VolumeSize
			}
			workers[i] = m
		}
		att["workers"] = workers
		if in.Zones != nil {
			att["zones"] = newStringSet(schema.HashString, in.Zones)
		}
	}

	return []interface{}{att}
}

func flattenCloudAzure(in *v1beta1.AzureCloud) []interface{} {
	att := make(map[string]interface{})

	if in.MachineImage != nil {
		image := make(map[string]interface{})
		if len(in.MachineImage.Name) > 0 {
			image["name"] = in.MachineImage.Name
		}
		if len(in.MachineImage.Version) > 0 {
			image["version"] = in.MachineImage.Version
		}
		att["machine_image"] = []interface{}{image}
	}
	networks := make(map[string]interface{})
	if in.Networks.Nodes != nil {
		networks["nodes"] = in.Networks.Nodes
	}
	if in.Networks.Pods != nil {
		networks["pods"] = in.Networks.Pods
	}
	if in.Networks.Services != nil {
		networks["services"] = in.Networks.Services
	}
	vnet := make(map[string]interface{})
	if in.Networks.VNet.Name != nil {
		vnet["id"] = in.Networks.VNet.Name
	}
	if in.Networks.VNet.CIDR != nil {
		vnet["cidr"] = in.Networks.VNet.CIDR
	}
	networks["vnet"] = []interface{}{vnet}
	if in.Networks.Workers != "" {
		networks["workers"] = in.Networks.Workers
	}
	att["networks"] = []interface{}{networks}
	if len(in.Workers) > 0 {
		workers := make([]interface{}, len(in.Workers))
		for i, v := range in.Workers {
			m := map[string]interface{}{}

			if v.Name != "" {
				m["name"] = v.Name
			}
			if v.MachineType != "" {
				m["machine_type"] = v.MachineType
			}
			if v.AutoScalerMin != 0 {
				m["auto_scaler_min"] = v.AutoScalerMin
			}
			if v.AutoScalerMax != 0 {
				m["auto_scaler_max"] = v.AutoScalerMax
			}
			if v.MaxSurge != nil {
				m["max_surge"] = v.MaxSurge.IntValue()
			}
			if v.MaxUnavailable != nil {
				m["max_unavailable"] = v.MaxUnavailable.IntValue()
			}
			if len(v.Annotations) > 0 {
				m["annotations"] = v.Annotations
			}
			if len(v.Labels) > 0 {
				m["labels"] = v.Labels
			}
			if len(v.Taints) > 0 {
				taints := make([]interface{}, len(v.Taints))
				for i, v := range v.Taints {
					m := map[string]interface{}{}

					if v.Key != "" {
						m["key"] = v.Key
					}
					// if v.Operator != "" {
					// 	m["operator"] = v.Operator
					// }
					if v.Value != "" {
						m["value"] = v.Value
					}
					if v.Effect != "" {
						m["effect"] = v.Effect
					}
					taints[i] = m
				}
				m["taints"] = taints
			}
			if len(v.VolumeType) > 0 {
				m["volume_type"] = v.VolumeType
			}
			if len(v.VolumeSize) > 0 {
				m["volume_size"] = v.VolumeSize
			}
			workers[i] = m
		}
		att["workers"] = workers
	}

	return []interface{}{att}
}
