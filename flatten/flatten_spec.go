package flatten

import (
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
		att["allow_privileged_containers"] = *in.AllowPrivilegedContainers
	}
	if in.KubeAPIServer != nil {
		att["kube_api_server"] = flattenKubeAPIServer(in.KubeAPIServer)
	}
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

func flattenKubeAPIServer(in *corev1beta1.KubeAPIServerConfig) []interface{} {
	att := make(map[string]interface{})

	if in.EnableBasicAuthentication != nil {
		att["enable_basic_authentication"] = *in.EnableBasicAuthentication
	}

	if in.OIDCConfig != nil {
		config := make(map[string]interface{})

		if in.OIDCConfig.CABundle != nil {
			config["ca_bundle"] = *in.OIDCConfig.CABundle
		}
		if in.OIDCConfig.ClientID != nil {
			config["client_id"] = *in.OIDCConfig.ClientID
		}
		if in.OIDCConfig.GroupsClaim != nil {
			config["groups_claim"] = *in.OIDCConfig.GroupsClaim
		}
		if in.OIDCConfig.GroupsPrefix != nil {
			config["groups_prefix"] = *in.OIDCConfig.GroupsPrefix
		}
		if in.OIDCConfig.IssuerURL != nil {
			config["issuer_url"] = *in.OIDCConfig.IssuerURL
		}
		if in.OIDCConfig.RequiredClaims != nil {
			config["required_claims"] = in.OIDCConfig.RequiredClaims
		}
		if len(in.OIDCConfig.SigningAlgs) > 0 {
			config["signing_algs"] = in.OIDCConfig.SigningAlgs
		}
		if in.OIDCConfig.UsernameClaim != nil {
			config["username_claim"] = *in.OIDCConfig.UsernameClaim
		}
		if in.OIDCConfig.UsernamePrefix != nil {
			config["username_prefix"] = *in.OIDCConfig.UsernamePrefix
		}

		att["oidc_config"] = []interface{}{config}
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
