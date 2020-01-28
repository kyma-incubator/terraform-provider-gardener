package flatten

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	//"github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	//"github.com/kyma-incubator/terraform-provider-gardener/expand"
)

// Flatteners
func FlattenShoot(in corev1beta1.ShootSpec, d *schema.ResourceData, specPrefix ...string) ([]interface{}, error) {
	att := make(map[string]interface{})
	//prefix := ""
	//if len(specPrefix) > 0 {
	//	prefix = specPrefix[0]
	//}
	//if in.Addons != nil {
	//	configAddons := d.Get(prefix + "spec.0.addons").([]interface{})
	//	flattenedAddons := flattenAddons(in.Addons)
	//	att["addons"] = expand.RemoveInternalKeysArraySpec([]interface{}{flattenedAddons}, configAddons) //expand.RemoveInternalKeysMetadata(flattenedAddons, configAddons)
	//}
	//configCloud := d.Get(prefix + "spec.0.cloud").([]interface{})
	//flattenedCloud := flattenCloud(in.Cloud)
	//att["cloud"] = expand.RemoveInternalKeysArraySpec(flattenedCloud, configCloud)
	//configDNS := d.Get(prefix + "spec.0.dns").([]interface{})
	//flattenedDNS := flattenDNS(in.DNS)
	//att["dns"] = expand.RemoveInternalKeysArraySpec(flattenedDNS, configDNS)
	//if in.Hibernation != nil {
	//	configHibernation := d.Get(prefix + "spec.0.hibernation").([]interface{})
	//	flattenedHibernation := flattenHibernation(in.Hibernation)
	//	att["hibernation"] = expand.RemoveInternalKeysArraySpec(flattenedHibernation, configHibernation)
	//}
	//configKubernetes := d.Get(prefix + "spec.0.kubernetes").([]interface{})
	//flattenedKubernetes := flattenKubernetes(in.Kubernetes)
	//att["kubernetes"] = expand.RemoveInternalKeysArraySpec(flattenedKubernetes, configKubernetes)
	//if in.Maintenance != nil {
	//	configMaintenance := d.Get(prefix + "spec.0.maintenance").([]interface{})
	//	flattenedMaintenance := flattenMaintenance(in.Maintenance)
	//	att["maintenance"] = expand.RemoveInternalKeysArraySpec(flattenedMaintenance, configMaintenance)
	//}

	return []interface{}{att}, nil
}

//func flattenAddons(in *corev1beta1.Addons) map[string]interface{} {
//	att := make(map[string]interface{})
//
//	if in.KubernetesDashboard != nil {
//		dashboard := make(map[string]interface{})
//		dashboard["enabled"] = in.KubernetesDashboard.Enabled
//		if in.KubernetesDashboard.AuthenticationMode != nil {
//			dashboard["authentication_mode"] = in.KubernetesDashboard.AuthenticationMode
//		}
//		att["kubernetes_dashboard"] = []interface{}{dashboard}
//	}
//	if in.NginxIngress != nil {
//		ingress := make(map[string]interface{})
//		ingress["enabled"] = in.NginxIngress.Enabled
//		att["nginx_ingress"] = []interface{}{ingress}
//	}
//	if in.ClusterAutoscaler != nil {
//		autoscaler := make(map[string]interface{})
//		autoscaler["enabled"] = in.ClusterAutoscaler.Enabled
//		att["cluster_autoscaler"] = []interface{}{autoscaler}
//	}
//
//	return att
//}
//
//func flattenDNS(in corev1beta1.DNS) []interface{} {
//	att := make(map[string]interface{})
//
//	if in.Provider != nil {
//		att["provider"] = in.Provider
//	}
//	if in.Domain != nil {
//		att["domain"] = in.Domain
//	}
//	if in.SecretName != nil {
//		att["secret_name"] = in.SecretName
//	}
//
//	return []interface{}{att}
//}
//
//func flattenHibernation(in *corev1beta1.Hibernation) []interface{} {
//	att := make(map[string]interface{})
//
//	att["enabled"] = in.Enabled
//	if len(in.Schedules) > 0 {
//		schedules := make([]interface{}, len(in.Schedules))
//		for i, v := range in.Schedules {
//			m := map[string]interface{}{}
//
//			if v.Start != nil {
//				m["start"] = v.Start
//			}
//			if v.End != nil {
//				m["end"] = v.End
//			}
//			if v.Location != nil {
//				m["location"] = v.Location
//			}
//			schedules[i] = m
//		}
//		att["schedules"] = schedules
//	}
//
//	return []interface{}{att}
//}
//
//func flattenKubernetes(in corev1beta1.Kubernetes) []interface{} {
//	att := make(map[string]interface{})
//
//	if in.AllowPrivilegedContainers != nil {
//		att["allow_privileged_containers"] = in.AllowPrivilegedContainers
//	}
//	if in.KubeAPIServer != nil {
//		server := make(map[string]interface{})
//		if in.KubeAPIServer.FeatureGates != nil {
//			server["enable_basic_authentication"] = in.KubeAPIServer.EnableBasicAuthentication
//		}
//		att["kube_api_server"] = []interface{}{server}
//	}
//	if in.Kubelet != nil {
//		kubelet := make(map[string]interface{})
//		if in.Kubelet.FeatureGates != nil {
//			kubelet["feature_gates"] = in.Kubelet.FeatureGates
//		}
//		if in.Kubelet.PodPIDsLimit != nil {
//			kubelet["pod_pids_limit"] = in.Kubelet.PodPIDsLimit
//		}
//		if in.Kubelet.CPUCFSQuota != nil {
//			kubelet["cpu_cfs_quota"] = in.Kubelet.CPUCFSQuota
//		}
//		att["kubelet"] = []interface{}{kubelet}
//	}
//	if len(in.Version) > 0 {
//		att["version"] = in.Version
//	}
//	if in.CloudControllerManager != nil {
//		manager := make(map[string]interface{})
//		if in.CloudControllerManager.FeatureGates != nil {
//			manager["feature_gates"] = in.CloudControllerManager.FeatureGates
//		}
//		att["cloud_controller_manager"] = []interface{}{manager}
//	}
//	if in.KubeControllerManager != nil {
//		manager := make(map[string]interface{})
//		if in.KubeControllerManager.FeatureGates != nil {
//			manager["node_cidr_mask_size"] = in.KubeControllerManager.NodeCIDRMaskSize
//		}
//		att["kube_controller_manager"] = []interface{}{manager}
//	}
//	if in.KubeProxy != nil {
//		proxy := make(map[string]interface{})
//		if in.KubeProxy.FeatureGates != nil {
//			proxy["feature_gates"] = in.KubeProxy.FeatureGates
//		}
//		if in.KubeProxy.Mode != nil {
//			proxy["mode"] = in.KubeProxy.Mode
//		}
//		att["kube_proxy"] = []interface{}{proxy}
//	}
//	if in.Kubelet != nil {
//		kubelet := make(map[string]interface{})
//		if in.Kubelet.FeatureGates != nil {
//			kubelet["feature_gates"] = in.Kubelet.FeatureGates
//		}
//		if in.Kubelet.PodPIDsLimit != nil {
//			kubelet["pod_pids_limit"] = in.Kubelet.PodPIDsLimit
//		}
//		if in.Kubelet.CPUCFSQuota != nil {
//			kubelet["cpu_cfs_quota"] = in.Kubelet.CPUCFSQuota
//		}
//		att["kubelet"] = []interface{}{kubelet}
//	}
//	if len(in.Version) > 0 {
//		att["version"] = in.Version
//	}
//	if in.ClusterAutoscaler != nil {
//		scaler := make(map[string]interface{})
//		if in.ClusterAutoscaler.ScaleDownUtilizationThreshold != nil {
//			scaler["scale_down_utilization_threshold"] = in.ClusterAutoscaler.ScaleDownUtilizationThreshold
//		}
//		att["cluster_autoscaler"] = []interface{}{scaler}
//	}
//
//	return []interface{}{att}
//}
//
//func flattenMaintenance(in *corev1beta1.Maintenance) []interface{} {
//	att := make(map[string]interface{})
//
//	if in.AutoUpdate != nil {
//		update := make(map[string]interface{})
//		update["kubernetes_version"] = in.AutoUpdate.KubernetesVersion
//		att["auto_update"] = []interface{}{update}
//	}
//	if in.TimeWindow != nil {
//		window := make(map[string]interface{})
//		if len(in.TimeWindow.Begin) > 0 {
//			window["begin"] = in.TimeWindow.Begin
//		}
//		if len(in.TimeWindow.End) > 0 {
//			window["end"] = in.TimeWindow.End
//		}
//		att["time_window"] = []interface{}{window}
//	}
//
//	return []interface{}{att}
//}
