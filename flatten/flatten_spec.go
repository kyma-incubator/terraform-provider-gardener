package flatten

import (
	"encoding/json"

	gcpAlpha1 "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"

	awsAlpha1 "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"
	azAlpha1 "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-incubator/terraform-provider-gardener/expand"

	//"github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	//"github.com/kyma-incubator/terraform-provider-gardener/expand"
)

// Flatteners
func FlattenShoot(in corev1beta1.ShootSpec, d *schema.ResourceData, specPrefix ...string) ([]interface{}, error) {
	att := make(map[string]interface{})
	prefix := ""
	if len(specPrefix) > 0 {
		prefix = specPrefix[0]
	}

	if len(in.CloudProfileName) > 0 {
		att["cloud_profile_name"] = in.CloudProfileName
	}
	if len(in.SecretBindingName) > 0 {
		att["secret_binding_name"] = in.SecretBindingName
	}
	if len(in.Region) > 0 {
		att["region"] = in.Region
	}
	if in.Purpose != nil {
		att["purpose"] = *in.Purpose
	}
	if in.Addons != nil {
		configAddons := d.Get(prefix + "spec.0.addons").([]interface{})
		flattenedAddons := flattenAddons(in.Addons)
		att["addons"] = expand.RemoveInternalKeysArraySpec([]interface{}{flattenedAddons}, configAddons) //expand.RemoveInternalKeysMetadata(flattenedAddons, configAddons)
	}
	configProvider := d.Get(prefix + "spec.0.provider").([]interface{})
	flattenedProvider := flattenProvider(in.Provider)
	att["provider"] = expand.RemoveInternalKeysArraySpec(flattenedProvider, configProvider)

	if in.DNS != nil {
		configDNS := d.Get(prefix + "spec.0.dns").([]interface{})
		flattenedDNS := flattenDNS(in.DNS)
		att["dns"] = expand.RemoveInternalKeysArraySpec(flattenedDNS, configDNS)
	}
	if in.Hibernation != nil {
		configHibernation := d.Get(prefix + "spec.0.hibernation").([]interface{})
		flattenedHibernation := flattenHibernation(in.Hibernation)
		att["hibernation"] = expand.RemoveInternalKeysArraySpec(flattenedHibernation, configHibernation)
	}
	configKubernetes := d.Get(prefix + "spec.0.kubernetes").([]interface{})
	flattenedKubernetes := flattenKubernetes(in.Kubernetes)
	att["kubernetes"] = expand.RemoveInternalKeysArraySpec(flattenedKubernetes, configKubernetes)
	if in.Maintenance != nil {
		configMaintenance := d.Get(prefix + "spec.0.maintenance").([]interface{})
		flattenedMaintenance := flattenMaintenance(in.Maintenance)
		att["maintenance"] = expand.RemoveInternalKeysArraySpec(flattenedMaintenance, configMaintenance)
	}
	configNetworking := d.Get(prefix + "spec.0.networking").([]interface{})
	flattenedNetworking := flattenNetworking(in.Networking)
	att["networking"] = expand.RemoveInternalKeysArraySpec(flattenedNetworking, configNetworking)
	if in.Monitoring != nil {
		configMonitoring := d.Get(prefix + "spec.0.monitoring").([]interface{})
		flattenedMonitoring := flattenMonitoring(in.Monitoring)
		att["monitoring"] = expand.RemoveInternalKeysArraySpec(flattenedMonitoring, configMonitoring)
	}

	return []interface{}{att}, nil
}

func flattenAddons(in *corev1beta1.Addons) map[string]interface{} {
	att := make(map[string]interface{})

	if in.KubernetesDashboard != nil {
		dashboard := make(map[string]interface{})
		dashboard["enabled"] = in.KubernetesDashboard.Enabled
		if in.KubernetesDashboard.AuthenticationMode != nil {
			dashboard["authentication_mode"] = *in.KubernetesDashboard.AuthenticationMode
		}
		att["kubernetes_dashboard"] = []interface{}{dashboard}
	}
	if in.NginxIngress != nil {
		ingress := make(map[string]interface{})
		ingress["enabled"] = in.NginxIngress.Enabled
		att["nginx_ingress"] = []interface{}{ingress}
	}
	//if in.ClusterAutoscaler != nil {
	//	autoscaler := make(map[string]interface{})
	//	autoscaler["enabled"] = in.ClusterAutoscaler.Enabled
	//	att["cluster_autoscaler"] = []interface{}{autoscaler}
	//}

	return att
}

func flattenDNS(in *corev1beta1.DNS) []interface{} {
	att := make(map[string]interface{})

	if len(in.Providers) > 0 {
		providers := make([]interface{}, len(in.Providers))
		for i, v := range in.Providers {
			m := map[string]interface{}{}
			if v.Domains != nil {
				m["domains"] = v.Domains
			}
			if v.SecretName != nil {
				m["secret_name"] = v.SecretName
			}
			if v.Type != nil {
				m["type"] = v.Type
			}

			if v.Zones != nil {
				m["zones"] = v.Zones
			}
			providers[i] = m
		}
		att["providers"] = providers
	}
	if in.Domain != nil {
		att["domain"] = *in.Domain
	}

	return []interface{}{att}
}

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

func flattenHibernation(in *corev1beta1.Hibernation) []interface{} {
	att := make(map[string]interface{})

	if in.Enabled != nil {
		att["enabled"] = *in.Enabled
	}
	if len(in.Schedules) > 0 {
		schedules := make([]interface{}, len(in.Schedules))
		for i, v := range in.Schedules {
			m := map[string]interface{}{}

			if v.Start != nil {
				m["start"] = *v.Start
			}
			if v.End != nil {
				m["end"] = *v.End
			}
			if v.Location != nil {
				m["location"] = *v.Location
			}
			schedules[i] = m
		}
		att["schedules"] = schedules
	}

	return []interface{}{att}
}

func flattenKubernetes(in corev1beta1.Kubernetes) []interface{} {
	att := make(map[string]interface{})

	if in.AllowPrivilegedContainers != nil {
		att["allow_privileged_containers"] = in.AllowPrivilegedContainers
	}
	if in.KubeAPIServer != nil {
		server := make(map[string]interface{})
		if in.KubeAPIServer.FeatureGates != nil {
			server["enable_basic_authentication"] = in.KubeAPIServer.EnableBasicAuthentication
		}
		att["kube_api_server"] = []interface{}{server}
	}
	//if in.CloudControllerManager != nil {
	//	manager := make(map[string]interface{})
	//	if in.CloudControllerManager.FeatureGates != nil {
	//		manager["feature_gates"] = in.CloudControllerManager.FeatureGates
	//	}
	//	att["cloud_controller_manager"] = []interface{}{manager}
	//}
	if in.KubeControllerManager != nil {
		manager := make(map[string]interface{})
		if in.KubeControllerManager.FeatureGates != nil {
			manager["node_cidr_mask_size"] = in.KubeControllerManager.NodeCIDRMaskSize
		}
		att["kube_controller_manager"] = []interface{}{manager}
	}
	if in.KubeProxy != nil {
		proxy := make(map[string]interface{})
		if in.KubeProxy.FeatureGates != nil {
			proxy["feature_gates"] = in.KubeProxy.FeatureGates
		}
		if in.KubeProxy.Mode != nil {
			proxy["mode"] = in.KubeProxy.Mode
		}
		att["kube_proxy"] = []interface{}{proxy}
	}
	if in.Kubelet != nil {
		kubelet := make(map[string]interface{})
		if in.Kubelet.FeatureGates != nil {
			kubelet["feature_gates"] = in.Kubelet.FeatureGates
		}
		if in.Kubelet.PodPIDsLimit != nil {
			kubelet["pod_pids_limit"] = in.Kubelet.PodPIDsLimit
		}
		if in.Kubelet.CPUCFSQuota != nil {
			kubelet["cpu_cfs_quota"] = in.Kubelet.CPUCFSQuota
		}
		att["kubelet"] = []interface{}{kubelet}
	}
	if len(in.Version) > 0 {
		att["version"] = in.Version
	}
	if in.ClusterAutoscaler != nil {
		scaler := make(map[string]interface{})
		if in.ClusterAutoscaler.ScaleDownUtilizationThreshold != nil {
			scaler["scale_down_utilization_threshold"] = in.ClusterAutoscaler.ScaleDownUtilizationThreshold
		}
		att["cluster_autoscaler"] = []interface{}{scaler}
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

func flattenMaintenance(in *corev1beta1.Maintenance) []interface{} {
	att := make(map[string]interface{})

	if in.AutoUpdate != nil {
		update := make(map[string]interface{})
		update["kubernetes_version"] = in.AutoUpdate.KubernetesVersion
		update["machine_image_version"] = in.AutoUpdate.MachineImageVersion
		att["auto_update"] = []interface{}{update}
	}
	if in.TimeWindow != nil {
		window := make(map[string]interface{})
		if len(in.TimeWindow.Begin) > 0 {
			window["begin"] = in.TimeWindow.Begin
		}
		if len(in.TimeWindow.End) > 0 {
			window["end"] = in.TimeWindow.End
		}
		att["time_window"] = []interface{}{window}
	}

	return []interface{}{att}
}

func flattenNetworking(in corev1beta1.Networking) []interface{} {
	att := make(map[string]interface{})

	if in.Nodes != nil {
		att["nodes"] = *in.Nodes
	}
	if in.Pods != nil {
		att["pods"] = *in.Pods
	}
	if in.Services != nil {
		att["services"] = *in.Services
	}
	if len(in.Type) > 0 {
		att["type"] = in.Type
	}

	return []interface{}{att}
}

func flattenMonitoring(in *corev1beta1.Monitoring) []interface{} {
	att := make(map[string]interface{})

	if in.Alerting != nil {
		alerting := make(map[string]interface{})

		if len(in.Alerting.EmailReceivers) > 0 {
			alerting["emailreceivers"] = in.Alerting.EmailReceivers
		}

		att["alerting"] = []interface{}{alerting}
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

func flattenGcpInfra(in gcpAlpha1.InfrastructureConfig) []interface{} {
	att := make(map[string]interface{})
	net := make(map[string]interface{})

	if len(in.Networks.Workers) > 0 {
		net["workers"] = in.Networks.Workers
	}

	if in.Networks.Internal != nil {
		net["internal"] = *in.Networks.Internal
	}

	vpc := make(map[string]interface{})

	if in.Networks.VPC != nil && len(in.Networks.VPC.Name) > 0 {
		vpc["name"] = in.Networks.VPC.Name
	}
	cr := make(map[string]interface{})
	if in.Networks.VPC != nil && len(in.Networks.VPC.CloudRouter.Name) > 0 {
		cr["name"] = in.Networks.VPC.CloudRouter.Name
		vpc["cloud_router"] = []interface{}{cr}
	}
	net["vpc"] = []interface{}{vpc}

	cn := make(map[string]interface{})
	if in.Networks.CloudNAT != nil && in.Networks.CloudNAT.MinPortsPerVM != nil {
		cn["min_ports_per_vm"] = *in.Networks.CloudNAT.MinPortsPerVM
	}
	net["cloud_nat"] = []interface{}{cn}

	fl := make(map[string]interface{})
	if in.Networks.FlowLogs != nil && in.Networks.FlowLogs.AggregationInterval != nil {
		fl["aggregation_interval"] = *in.Networks.FlowLogs.AggregationInterval
	}
	if in.Networks.FlowLogs != nil && in.Networks.FlowLogs.Metadata != nil {
		fl["metadata"] = *in.Networks.FlowLogs.Metadata
	}
	if in.Networks.FlowLogs != nil && in.Networks.FlowLogs.FlowSampling != nil {
		fl["flow_sampling"] = *in.Networks.FlowLogs.FlowSampling
	}

	net["flow_logs"] = []interface{}{fl}
	att["networks"] = []interface{}{net}

	return []interface{}{att}
}

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
		// zones := make([]interface{}, len(in.Networks.Zones))
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

func flattenGcpControlPlane(in gcpAlpha1.ControlPlaneConfig) []interface{} {
	att := make(map[string]interface{})

	if len(in.Zone) > 0 {
		att["zone"] = in.Zone
	}

	return []interface{}{att}
}
