package flatters

import (
	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

// Flatteners
func FlattenShoot(in v1beta1.ShootSpec, d *schema.ResourceData) ([]interface{}, error) {
	att := make(map[string]interface{})

	if in.Addons != nil {
		att["addons"] = flattenAddons(in.Addons)
	}
	att["cloud"] = flattenCloud(in.Cloud)
	att["dns"] = flattenDNS(in.DNS)
	if in.Hibernation != nil {
		att["hibernation"] = flattenHibernation(in.Hibernation)
	}
	att["kubernetes"] = flattenKubernetes(in.Kubernetes)
	if in.Maintenance != nil {
		att["maintenance"] = flattenMaintenance(in.Maintenance)
	}

	return []interface{}{att}, nil
}

func flattenAddons(in *v1beta1.Addons) []interface{} {
	att := make(map[string]interface{})

	if in.KubernetesDashboard != nil {
		dashboard := make(map[string]interface{})
		dashboard["enabled"] = in.KubernetesDashboard.Enabled
		if in.KubernetesDashboard.AuthenticationMode != nil {
			dashboard["authentication_mode"] = in.KubernetesDashboard.AuthenticationMode
		}
		att["kubernetes_dashboard"] = []interface{}{dashboard}
	}
	if in.NginxIngress != nil {
		ingress := make(map[string]interface{})
		ingress["enabled"] = in.NginxIngress.Enabled
		att["nginx_ingress"] = []interface{}{ingress}
	}
	if in.ClusterAutoscaler != nil {
		autoscaler := make(map[string]interface{})
		autoscaler["enabled"] = in.ClusterAutoscaler.Enabled
		att["cluster_autoscaler"] = []interface{}{autoscaler}
	}

	return []interface{}{att}
}

func flattenDNS(in v1beta1.DNS) []interface{} {
	att := make(map[string]interface{})

	if in.Provider != nil {
		att["provider"] = in.Provider
	}
	if in.Domain != nil {
		att["domain"] = in.Domain
	}
	if in.SecretName != nil {
		att["secret_name"] = in.SecretName
	}

	return []interface{}{att}
}

func flattenHibernation(in *v1beta1.Hibernation) []interface{} {
	att := make(map[string]interface{})

	att["enabled"] = in.Enabled
	if len(in.Schedules) > 0 {
		schedules := make([]interface{}, len(in.Schedules))
		for i, v := range in.Schedules {
			m := map[string]interface{}{}

			if v.Start != nil {
				m["start"] = v.Start
			}
			if v.End != nil {
				m["end"] = v.End
			}
			if v.Location != nil {
				m["location"] = v.Location
			}
			schedules[i] = m
		}
		att["schedules"] = schedules
	}

	return []interface{}{att}
}

func flattenKubernetes(in v1beta1.Kubernetes) []interface{} {
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
	if in.CloudControllerManager != nil {
		manager := make(map[string]interface{})
		if in.CloudControllerManager.FeatureGates != nil {
			manager["feature_gates"] = in.CloudControllerManager.FeatureGates
		}
		att["cloud_controller_manager"] = []interface{}{manager}
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

func flattenMaintenance(in *v1beta1.Maintenance) []interface{} {
	att := make(map[string]interface{})

	if in.AutoUpdate != nil {
		update := make(map[string]interface{})
		update["kubernetes_version"] = in.AutoUpdate.KubernetesVersion
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
