package flatten

import (
	"encoding/json"
	awsAlpha1 "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"
	azAlpha1 "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
	gcpAlpha1 "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
)

func flattenProvider(in corev1beta1.Provider) []interface{} {
	att := make(map[string]interface{})

	if len(in.Type) > 0 {
		att["type"] = in.Type
	}

	if len(in.Workers) > 0 {
		workers := make([]interface{}, len(in.Workers))
		for i, v := range in.Workers {
			m := map[string]interface{}{}

			if len(v.Name) > 0 {
				m["name"] = v.Name
			}
			if len(v.Zones) > 0 {
				m["zones"] = v.Zones
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
			if v.MaxSurge != nil {
				m["max_surge"] = v.MaxSurge.IntValue()
			}
			if v.MaxUnavailable != nil {
				m["max_unavailable"] = v.MaxUnavailable.IntValue()
			}
			if v.CABundle != nil {
				m["cabundle"] = *v.CABundle
			}

			if v.Minimum != 0 {
				m["minimum"] = v.Minimum
			}

			if v.Maximum != 0 {
				m["maximum"] = v.Maximum
			}

			if v.Kubernetes != nil {
				m["kubernetes"] = flattenWorkerKubernetes(v.Kubernetes)
			}

			if len(v.Annotations) > 0 {
				m["annotations"] = v.Annotations
			}
			if len(v.Labels) > 0 {
				m["labels"] = v.Labels
			}
			if v.Volume != nil {
				m["volume"] = flattenVolume(v.Volume)
			}
			m["machine"] = flattenMachine(v.Machine)

			workers[i] = m
		}
		att["worker"] = workers
	}

	if in.InfrastructureConfig != nil {
		att["infrastructure_config"] = flattenInfrastructureConfig(in.Type, in.InfrastructureConfig)
	}

	if in.ControlPlaneConfig != nil {
		att["control_plane_config"] = flattenControlPlaneConfig(in.Type, in.ControlPlaneConfig)
	}

	return []interface{}{att}
}

func flattenWorkerKubernetes(in *corev1beta1.WorkerKubernetes) []interface{} {
	att := make(map[string]interface{})

	if in.Kubelet != nil {
		kubelet := make(map[string]interface{})
		if in.Kubelet.PodPIDsLimit != nil {
			kubelet["pod_pids_limit"] = *in.Kubelet.PodPIDsLimit
		}
		if in.Kubelet.CPUManagerPolicy != nil {
			kubelet["cpu_manager_policy"] = *in.Kubelet.CPUManagerPolicy
		}
		if in.Kubelet.CPUCFSQuota != nil {
			kubelet["cpu_cfs_quota"] = *in.Kubelet.CPUCFSQuota
		}
		att["kubelet"] = []interface{}{kubelet}
	}

	return []interface{}{att}
}

func flattenVolume(in *corev1beta1.Volume) []interface{} {
	att := map[string]interface{}{}

	if len(in.Size) > 0 {
		att["size"] = in.Size
	}
	if in.Type != nil {
		att["type"] = *in.Type
	}

	return []interface{}{att}
}
func flattenMachine(in corev1beta1.Machine) []interface{} {
	att := map[string]interface{}{}

	if len(in.Type) > 0 {
		att["type"] = in.Type
	}
	if in.Image != nil {
		att["image"] = flattenMachineImage(in.Image)
	}

	return []interface{}{att}
}

func flattenMachineImage(in *corev1beta1.ShootMachineImage) []interface{} {
	att := map[string]interface{}{}

	if len(in.Name) > 0 {
		att["name"] = in.Name
	}
	if len(in.Version) > 0 {
		att["version"] = in.Version
	}

	return []interface{}{att}
}

func flattenInfrastructureConfig(providerType string, in *corev1beta1.ProviderConfig) []interface{} {
	att := map[string]interface{}{}

	if providerType == "azure" {
		azConfigObj := azAlpha1.InfrastructureConfig{}
		if err := json.Unmarshal(in.RawExtension.Raw, &azConfigObj); err == nil {
			att["azure"] = flattenAzure(azConfigObj)
		}
	}

	if providerType == "aws" {
		awsConfigObj := awsAlpha1.InfrastructureConfig{}
		if err := json.Unmarshal(in.RawExtension.Raw, &awsConfigObj); err == nil {
			att["aws"] = flattenAws(awsConfigObj)
		}
	}

	if providerType == "gcp" {
		gcpConfigObj := gcpAlpha1.InfrastructureConfig{}
		if err := json.Unmarshal(in.RawExtension.Raw, &gcpConfigObj); err == nil {
			att["gcp"] = flattenGcpInfra(gcpConfigObj)
		}
	}

	return []interface{}{att}
}

func flattenControlPlaneConfig(providerType string, in *corev1beta1.ProviderConfig) []interface{} {
	att := map[string]interface{}{}
	if providerType == "gcp" {
		gcpConfigObj := gcpAlpha1.ControlPlaneConfig{}
		if err := json.Unmarshal(in.RawExtension.Raw, &gcpConfigObj); err == nil {
			att["gcp"] = flattenGcpControlPlane(gcpConfigObj)
		}
	}
	return []interface{}{att}
}
