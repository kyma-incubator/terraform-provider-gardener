package helper

import (
	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

// Expanders
func ExpandShoot(shoot []interface{}) v1beta1.ShootSpec {
	obj := v1beta1.ShootSpec{}

	if len(shoot) == 0 || shoot[0] == nil {
		return obj
	}
	in := shoot[0].(map[string]interface{})

	if v, ok := in["addons"].([]interface{}); ok && len(v) > 0 {
		obj.Addons = expandAddons(v)
	}
	if v, ok := in["cloud"].([]interface{}); ok && len(v) > 0 {
		obj.Cloud = expandCloud(v)
	}
	if v, ok := in["dns"].([]interface{}); ok && len(v) > 0 {
		dns := expandDNS(v)
		if dns.Domain != nil {
			obj.DNS = dns
		}
	}
	if v, ok := in["hibernation"].([]interface{}); ok && len(v) > 0 {
		obj.Hibernation = expandHibernation(v)
	}
	if v, ok := in["kubernetes"].([]interface{}); ok && len(v) > 0 {
		obj.Kubernetes = expandKubernetes(v)
	}
	if v, ok := in["maintenance"].([]interface{}); ok && len(v) > 0 {
		obj.Maintenance = expandMaintenance(v)
	}

	return obj
}

func expandAddons(addon []interface{}) *v1beta1.Addons {
	obj := &v1beta1.Addons{}

	if len(addon) == 0 || addon[0] == nil {
		return obj
	}
	in := addon[0].(map[string]interface{})

	if v, ok := in["kubernetes_dashboard"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubernetesDashboard = &v1beta1.KubernetesDashboard{}

		if v, ok := v["enabled"].(bool); ok {
			obj.KubernetesDashboard.Enabled = v
		}
		if v, ok := v["authentication_mode"].(string); ok && len(v) > 0 {
			obj.KubernetesDashboard.AuthenticationMode = &v
		}
	}
	if v, ok := in["nginx_ingress"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.NginxIngress = &v1beta1.NginxIngress{}

		if v, ok := v["enabled"].(bool); ok {
			obj.NginxIngress.Enabled = v
		}
		if v, ok := v["load_balancer_source_ranges"].(*schema.Set); ok {
			obj.NginxIngress.LoadBalancerSourceRanges = expandSet(v)
		}
	}
	if v, ok := in["cluster_autoscaler"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.ClusterAutoscaler = &v1beta1.AddonClusterAutoscaler{}

		if v, ok := v["enabled"].(bool); ok {
			obj.ClusterAutoscaler.Enabled = v
		}
	}
	if v, ok := in["heapster"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.Heapster = &v1beta1.Heapster{}

		if v, ok := v["enabled"].(bool); ok {
			obj.Heapster.Enabled = v
		}
	}
	if v, ok := in["kube2iam"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.Kube2IAM = &v1beta1.Kube2IAM{}

		if v, ok := v["enabled"].(bool); ok {
			obj.Kube2IAM.Enabled = v
		}
		if roles, ok := v["roles"].([]interface{}); ok {
			for _, r := range roles {
				if r, ok := r.(map[string]interface{}); ok {
					roleObj := v1beta1.Kube2IAMRole{}

					if v, ok := r["name"].(string); ok && len(v) > 0 {
						roleObj.Name = v
					}
					if v, ok := r["description"].(string); ok && len(v) > 0 {
						roleObj.Description = v
					}
					if v, ok := r["policy"].(string); ok && len(v) > 0 {
						roleObj.Policy = v
					}

					obj.Kube2IAM.Roles = append(obj.Kube2IAM.Roles, roleObj)
				}
			}
		}
	}
	if v, ok := in["kube_lego"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubeLego = &v1beta1.KubeLego{}

		if v, ok := v["enabled"].(bool); ok {
			obj.KubeLego.Enabled = v
		}
		if v, ok := v["email"].(string); ok && len(v) > 0 {
			obj.KubeLego.Mail = v
		}
	}
	if v, ok := in["monocular"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.Monocular = &v1beta1.Monocular{}

		if v, ok := v["enabled"].(bool); ok {
			obj.Monocular.Enabled = v
		}
	}

	return obj
}

func expandCloud(cloud []interface{}) v1beta1.Cloud {
	obj := v1beta1.Cloud{}

	if len(cloud) == 0 || cloud[0] == nil {
		return obj
	}
	in := cloud[0].(map[string]interface{})

	if v, ok := in["profile"].(string); ok {
		obj.Profile = v
	}
	if v, ok := in["region"].(string); ok {
		obj.Region = v
	}
	if v, ok := in["secret_binding_ref"].([]interface{}); ok && len(v) > 0 {
		obj.SecretBindingRef = *expandLocalObjectReference(v)
	}
	if v, ok := in["seed"].(string); ok {
		obj.Seed = &v
	}
	if v, ok := in["aws"].([]interface{}); ok && len(v) > 0 {
		obj.AWS = expandCloudAWS(v)
	} else if v, ok := in["gcp"].([]interface{}); ok && len(v) > 0 {
		obj.GCP = expandCloudGCP(v)
	} else if v, ok := in["azure"].([]interface{}); ok && len(v) > 0 {
		obj.Azure = expandCloudAzure(v)
	}

	return obj
}

func expandCloudAWS(aws []interface{}) *v1beta1.AWSCloud {
	obj := &v1beta1.AWSCloud{}

	if len(aws) == 0 || aws[0] == nil {
		return obj
	}
	in := aws[0].(map[string]interface{})

	if v, ok := in["machine_image"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.MachineImage = &v1beta1.ShootMachineImage{}

		if v, ok := v["name"].(string); ok && len(v) > 0 {
			obj.MachineImage.Name = v
		}
		if v, ok := v["version"].(string); ok && len(v) > 0 {
			obj.MachineImage.Version = v
		}
	}
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.Networks = v1beta1.AWSNetworks{}

		if v, ok := v["nodes"].(string); ok && len(v) > 0 {
			obj.Networks.Nodes = expandCIDR(v)
		}
		if v, ok := v["pods"].(string); ok && len(v) > 0 {
			obj.Networks.Pods = expandCIDR(v)
		}
		if v, ok := v["services"].(string); ok && len(v) > 0 {
			obj.Networks.Services = expandCIDR(v)
		}
		if v, ok := v["vpc"].([]interface{}); ok && len(v) > 0 {
			v := v[0].(map[string]interface{})
			obj.Networks.VPC = v1beta1.AWSVPC{}

			if v, ok := v["id"].(string); ok && len(v) > 0 {
				obj.Networks.VPC.ID = &v
			}
			if v, ok := v["cidr"].(string); ok && len(v) > 0 {
				obj.Networks.VPC.CIDR = expandCIDR(v)
			}
		}
		if internal, ok := v["internal"].(*schema.Set); ok {
			for _, i := range internal.List() {
				obj.Networks.Internal = append(obj.Networks.Internal, *expandCIDR(i.(string)))
			}
		}
		if public, ok := v["public"].(*schema.Set); ok {
			for _, p := range public.List() {
				obj.Networks.Public = append(obj.Networks.Public, *expandCIDR(p.(string)))
			}
		}
		if workers, ok := v["workers"].(*schema.Set); ok {
			for _, w := range workers.List() {
				obj.Networks.Workers = append(obj.Networks.Workers, *expandCIDR(w.(string)))
			}
		}
	}
	if workers, ok := in["workers"].([]interface{}); ok && len(workers) > 0 {
		for _, w := range workers {
			if w, ok := w.(map[string]interface{}); ok {
				workerObj := v1beta1.AWSWorker{}

				if v, ok := w["name"].(string); ok && len(v) > 0 {
					workerObj.Name = v
				}
				if v, ok := w["machine_type"].(string); ok && len(v) > 0 {
					workerObj.MachineType = v
				}
				if v, ok := w["auto_scaler_min"].(int); ok {
					workerObj.AutoScalerMin = v
				}
				if v, ok := w["auto_scaler_max"].(int); ok {
					workerObj.AutoScalerMax = v
				}
				if v, ok := w["max_surge"].(string); ok && len(v) > 0 {
					surge := intstr.FromString(v)
					workerObj.MaxSurge = &surge
				}
				if v, ok := w["max_unavailable"].(string); ok && len(v) > 0 {
					surge := intstr.FromString(v)
					workerObj.MaxUnavailable = &surge
				}
				if v, ok := w["annotations"].(map[string]interface{}); ok {
					workerObj.Annotations = expandStringMap(v)
				}
				if v, ok := w["labels"].(map[string]interface{}); ok {
					workerObj.Labels = expandStringMap(v)
				}
				if taints, ok := w["taints"].([]interface{}); ok && len(taints) > 0 {
					for _, t := range taints {
						if t, ok := t.(map[string]interface{}); ok {
							taint := corev1.Taint{}

							if v, ok := t["key"].(string); ok && len(v) > 0 {
								taint.Key = v
							}
							// if v, ok := t["operator"].(string); ok && len(v) > 0 {
							// 	taint.Operator = v
							// }
							if v, ok := t["value"].(string); ok && len(v) > 0 {
								taint.Value = v
							}
							if v, ok := t["effect"].(string); ok && len(v) > 0 {
								taint.Effect = corev1.TaintEffect(v)
							}

							workerObj.Taints = append(workerObj.Taints, taint)
						}
					}
				}
				if v, ok := w["volume_type"].(string); ok && len(v) > 0 {
					workerObj.VolumeType = v
				}
				if v, ok := w["volume_size"].(string); ok && len(v) > 0 {
					workerObj.VolumeSize = v
				}

				obj.Workers = append(obj.Workers, workerObj)
			}
		}
	}
	if v, ok := in["zones"].(*schema.Set); ok {
		obj.Zones = expandSet(v)
	}

	return obj
}
func expandCloudGCP(aws []interface{}) *v1beta1.GCPCloud {
	obj := &v1beta1.GCPCloud{}

	if len(aws) == 0 || aws[0] == nil {
		return obj
	}
	in := aws[0].(map[string]interface{})

	if v, ok := in["machine_image"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.MachineImage = &v1beta1.ShootMachineImage{}

		if v, ok := v["name"].(string); ok && len(v) > 0 {
			obj.MachineImage.Name = v
		}
		if v, ok := v["version"].(string); ok && len(v) > 0 {
			obj.MachineImage.Version = v
		}
	}
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.Networks = v1beta1.GCPNetworks{}

		if v, ok := v["nodes"].(string); ok && len(v) > 0 {
			obj.Networks.Nodes = expandCIDR(v)
		}
		if v, ok := v["pods"].(string); ok && len(v) > 0 {
			obj.Networks.Pods = expandCIDR(v)
		}
		if v, ok := v["services"].(string); ok && len(v) > 0 {
			obj.Networks.Services = expandCIDR(v)
		}
		if workers, ok := v["workers"].(*schema.Set); ok {
			for _, w := range workers.List() {
				obj.Networks.Workers = append(obj.Networks.Workers, *expandCIDR(w.(string)))
			}
		}
	}
	if workers, ok := in["workers"].([]interface{}); ok && len(workers) > 0 {
		for _, w := range workers {
			if w, ok := w.(map[string]interface{}); ok {
				workerObj := v1beta1.GCPWorker{}

				if v, ok := w["name"].(string); ok && len(v) > 0 {
					workerObj.Name = v
				}
				if v, ok := w["machine_type"].(string); ok && len(v) > 0 {
					workerObj.MachineType = v
				}
				if v, ok := w["auto_scaler_min"].(int); ok {
					workerObj.AutoScalerMin = v
				}
				if v, ok := w["auto_scaler_max"].(int); ok {
					workerObj.AutoScalerMax = v
				}
				if v, ok := w["max_surge"].(string); ok && len(v) > 0 {
					surge := intstr.FromString(v)
					workerObj.MaxSurge = &surge
				}
				if v, ok := w["max_unavailable"].(string); ok && len(v) > 0 {
					surge := intstr.FromString(v)
					workerObj.MaxUnavailable = &surge
				}
				if v, ok := w["annotations"].(map[string]interface{}); ok {
					workerObj.Annotations = expandStringMap(v)
				}
				if v, ok := w["labels"].(map[string]interface{}); ok {
					workerObj.Labels = expandStringMap(v)
				}
				if taints, ok := w["taints"].([]interface{}); ok && len(taints) > 0 {
					for _, t := range taints {
						if t, ok := t.(map[string]interface{}); ok {
							taint := corev1.Taint{}

							if v, ok := t["key"].(string); ok && len(v) > 0 {
								taint.Key = v
							}
							if v, ok := t["value"].(string); ok && len(v) > 0 {
								taint.Value = v
							}
							if v, ok := t["effect"].(string); ok && len(v) > 0 {
								taint.Effect = corev1.TaintEffect(v)
							}

							workerObj.Taints = append(workerObj.Taints, taint)
						}
					}
				}
				if v, ok := w["volume_type"].(string); ok && len(v) > 0 {
					workerObj.VolumeType = v
				}
				if v, ok := w["volume_size"].(string); ok && len(v) > 0 {
					workerObj.VolumeSize = v
				}

				obj.Workers = append(obj.Workers, workerObj)
			}
		}
	}
	if v, ok := in["zones"].(*schema.Set); ok {
		obj.Zones = expandSet(v)
	}

	return obj
}
func expandCloudAzure(aws []interface{}) *v1beta1.AzureCloud {
	obj := &v1beta1.AzureCloud{}

	if len(aws) == 0 || aws[0] == nil {
		return obj
	}
	in := aws[0].(map[string]interface{})

	if v, ok := in["machine_image"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.MachineImage = &v1beta1.ShootMachineImage{}

		if v, ok := v["name"].(string); ok && len(v) > 0 {
			obj.MachineImage.Name = v
		}
		if v, ok := v["version"].(string); ok && len(v) > 0 {
			obj.MachineImage.Version = v
		}
	}
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.Networks = v1beta1.AzureNetworks{}

		if v, ok := v["nodes"].(string); ok && len(v) > 0 {
			obj.Networks.Nodes = expandCIDR(v)
		}
		if v, ok := v["pods"].(string); ok && len(v) > 0 {
			obj.Networks.Pods = expandCIDR(v)
		}
		if v, ok := v["services"].(string); ok && len(v) > 0 {
			obj.Networks.Services = expandCIDR(v)
		}
		if v, ok := v["vnet"].([]interface{}); ok && len(v) > 0 {
			v := v[0].(map[string]interface{})
			obj.Networks.VNet = v1beta1.AzureVNet{}

			if v, ok := v["name"].(string); ok && len(v) > 0 {
				obj.Networks.VNet.Name = &v
			}
			if v, ok := v["cidr"].(string); ok && len(v) > 0 {
				obj.Networks.VNet.CIDR = expandCIDR(v)
			}
		}
		if workers, ok := v["workers"].(string); ok {
			obj.Networks.Workers = *expandCIDR(workers)

		}
	}
	if workers, ok := in["workers"].([]interface{}); ok && len(workers) > 0 {
		for _, w := range workers {
			if w, ok := w.(map[string]interface{}); ok {
				workerObj := v1beta1.AzureWorker{}

				if v, ok := w["name"].(string); ok && len(v) > 0 {
					workerObj.Name = v
				}
				if v, ok := w["machine_type"].(string); ok && len(v) > 0 {
					workerObj.MachineType = v
				}
				if v, ok := w["auto_scaler_min"].(int); ok {
					workerObj.AutoScalerMin = v
				}
				if v, ok := w["auto_scaler_max"].(int); ok {
					workerObj.AutoScalerMax = v
				}
				if v, ok := w["max_surge"].(string); ok && len(v) > 0 {
					surge := intstr.FromString(v)
					workerObj.MaxSurge = &surge
				}
				if v, ok := w["max_unavailable"].(string); ok && len(v) > 0 {
					surge := intstr.FromString(v)
					workerObj.MaxUnavailable = &surge
				}
				if v, ok := w["annotations"].(map[string]interface{}); ok {
					workerObj.Annotations = expandStringMap(v)
				}
				if v, ok := w["labels"].(map[string]interface{}); ok {
					workerObj.Labels = expandStringMap(v)
				}
				if taints, ok := w["taints"].([]interface{}); ok && len(taints) > 0 {
					for _, t := range taints {
						if t, ok := t.(map[string]interface{}); ok {
							taint := corev1.Taint{}

							if v, ok := t["key"].(string); ok && len(v) > 0 {
								taint.Key = v
							}
							// if v, ok := t["operator"].(string); ok && len(v) > 0 {
							// 	taint.Operator = v
							// }
							if v, ok := t["value"].(string); ok && len(v) > 0 {
								taint.Value = v
							}
							if v, ok := t["effect"].(string); ok && len(v) > 0 {
								taint.Effect = corev1.TaintEffect(v)
							}

							workerObj.Taints = append(workerObj.Taints, taint)
						}
					}
				}
				if v, ok := w["volume_type"].(string); ok && len(v) > 0 {
					workerObj.VolumeType = v
				}
				if v, ok := w["volume_size"].(string); ok && len(v) > 0 {
					workerObj.VolumeSize = v
				}

				obj.Workers = append(obj.Workers, workerObj)
			}
		}
	}

	return obj
}
func expandDNS(dns []interface{}) v1beta1.DNS {
	obj := v1beta1.DNS{}

	if len(dns) == 0 || dns[0] == nil {
		return obj
	}
	in := dns[0].(map[string]interface{})

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = &v
	}
	if v, ok := in["domain"].(string); ok && len(v) > 0 {
		obj.Domain = &v
	}
	if v, ok := in["secret_name"].(string); ok && len(v) > 0 {
		obj.SecretName = &v
	}

	return obj
}

func expandHibernation(hibernation []interface{}) *v1beta1.Hibernation {
	obj := &v1beta1.Hibernation{}

	if len(hibernation) == 0 || hibernation[0] == nil {
		return obj
	}
	in := hibernation[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = &v
	}
	if schedules, ok := in["schedules"].([]interface{}); ok {
		for _, s := range schedules {
			if s, ok := s.(map[string]interface{}); ok {
				scheduleObj := v1beta1.HibernationSchedule{}

				if v, ok := s["start"].(string); ok && len(v) > 0 {
					scheduleObj.Start = &v
				}
				if v, ok := s["end"].(string); ok && len(v) > 0 {
					scheduleObj.End = &v
				}
				if v, ok := s["location"].(string); ok && len(v) > 0 {
					scheduleObj.Location = &v
				}

				obj.Schedules = append(obj.Schedules, scheduleObj)
			}
		}
	}

	return obj
}

func expandKubernetes(kubernetes []interface{}) v1beta1.Kubernetes {
	obj := v1beta1.Kubernetes{}

	if len(kubernetes) == 0 || kubernetes[0] == nil {
		return obj
	}
	in := kubernetes[0].(map[string]interface{})

	if v, ok := in["allow_privileged_containers"].(bool); ok {
		obj.AllowPrivilegedContainers = &v
	}
	if v, ok := in["kube_api_server"].([]interface{}); ok && len(v) > 0 {
		obj.KubeAPIServer = expandKubernetesAPIServer(v)
	}
	if v, ok := in["cloud_controller_manager"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.CloudControllerManager = &v1beta1.CloudControllerManagerConfig{}

		if v, ok := v["feature_gates"].(map[string]interface{}); ok {
			obj.CloudControllerManager.FeatureGates = expandBoolMap(v)
		}
	}
	if v, ok := in["kube_controller_manager"].([]interface{}); ok && len(v) > 0 {
		obj.KubeControllerManager = expandKubernetesControllerManager(v)
	}
	if v, ok := in["kube_scheduler"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubeScheduler = &v1beta1.KubeSchedulerConfig{}

		if v, ok := v["feature_gates"].(map[string]interface{}); ok {
			obj.KubeScheduler.FeatureGates = expandBoolMap(v)
		}
	}
	if v, ok := in["kube_proxy"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubeProxy = &v1beta1.KubeProxyConfig{}

		if v, ok := v["feature_gates"].(map[string]interface{}); ok {
			obj.KubeProxy.FeatureGates = expandBoolMap(v)
		}
		if v, ok := v["mode"].(string); ok && len(v) > 0 {
			mode := v1beta1.ProxyModeIPTables
			obj.KubeProxy.Mode = &mode
		}
	}
	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.Kubelet = &v1beta1.KubeletConfig{}

		if v, ok := v["feature_gates"].(map[string]interface{}); ok {
			obj.Kubelet.FeatureGates = expandBoolMap(v)
		}
		if v, ok := v["pod_pids_limit"].(*int64); ok {
			obj.Kubelet.PodPIDsLimit = v
		}
		if v, ok := v["cpu_cfs_quota"].(*bool); ok {
			obj.Kubelet.CPUCFSQuota = v
		}
	}
	if v, ok := in["version"].(string); ok {
		obj.Version = v
	}
	if v, ok := in["cluster_autoscaler"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.ClusterAutoscaler = &v1beta1.ClusterAutoscaler{}

		if v, ok := v["scale_down_utilization_threshold"].(float64); ok {
			obj.ClusterAutoscaler.ScaleDownUtilizationThreshold = &v
		}
	}

	return obj
}

func expandKubernetesControllerManager(controller []interface{}) *v1beta1.KubeControllerManagerConfig {
	obj := &v1beta1.KubeControllerManagerConfig{}

	if len(controller) == 0 || controller[0] == nil {
		return obj
	}
	in := controller[0].(map[string]interface{})

	if v, ok := in["feature_gates"].(map[string]interface{}); ok {
		obj.FeatureGates = expandBoolMap(v)
	}
	if v, ok := in["horizontal_pod_autoscaler"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.HorizontalPodAutoscalerConfig = &v1beta1.HorizontalPodAutoscalerConfig{}

		if v, ok := v["downscale_delay"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.DownscaleDelay = expandDuration(v)
		}
		if v, ok := v["sync_period"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.SyncPeriod = expandDuration(v)
		}
		if v, ok := v["tolerance"].(*float64); ok {
			obj.HorizontalPodAutoscalerConfig.Tolerance = v
		}
		if v, ok := v["upscale_delay"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.UpscaleDelay = expandDuration(v)
		}
		if v, ok := v["downscale_stabilization"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.DownscaleStabilization = expandDuration(v)
		}
		if v, ok := v["initial_readiness_delay"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.InitialReadinessDelay = expandDuration(v)
		}
		if v, ok := v["cpu_initialization_period"].(string); ok && len(v) > 0 {
			obj.HorizontalPodAutoscalerConfig.CPUInitializationPeriod = expandDuration(v)
		}
	}

	return obj
}

func expandKubernetesAPIServer(server []interface{}) *v1beta1.KubeAPIServerConfig {
	obj := &v1beta1.KubeAPIServerConfig{}

	if len(server) == 0 || server[0] == nil {
		return obj
	}
	in := server[0].(map[string]interface{})

	if v, ok := in["feature_gates"].(map[string]interface{}); ok {
		obj.FeatureGates = expandBoolMap(v)
	}
	if v, ok := in["runtime_config"].(map[string]interface{}); ok {
		obj.RuntimeConfig = expandBoolMap(v)
	}
	if v, ok := in["oidc_config"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.OIDCConfig = &v1beta1.OIDCConfig{}

		if v, ok := v["ca_bundle"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.CABundle = &v
		}
		if v, ok := v["client_id"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.ClientID = &v
		}
		if v, ok := v["groups_claim"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.GroupsClaim = &v
		}
		if v, ok := v["groups_prefix"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.GroupsPrefix = &v
		}
		if v, ok := v["issuer_url"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.IssuerURL = &v
		}
		if v, ok := v["required_claims"].(map[string]interface{}); ok {
			obj.OIDCConfig.RequiredClaims = expandStringMap(v)
		}
		if v, ok := v["signing_algs"].(*schema.Set); ok {
			obj.OIDCConfig.SigningAlgs = expandSet(v)
		}
		if v, ok := v["username_claim"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.UsernameClaim = &v
		}
		if v, ok := v["username_prefix"].(string); ok && len(v) > 0 {
			obj.OIDCConfig.UsernamePrefix = &v
		}
	}
	if plugins, ok := in["admission_plugins"].([]interface{}); ok && len(plugins) > 0 {
		for _, p := range plugins {
			if p, ok := p.(map[string]interface{}); ok {
				pluginObj := v1beta1.AdmissionPlugin{}

				if v, ok := p["name"].(string); ok && len(v) > 0 {
					pluginObj.Name = v
				}
				if v, ok := p["config"].(string); ok && len(v) > 0 {
					pluginObj.Config = &v
				}

				obj.AdmissionPlugins = append(obj.AdmissionPlugins, pluginObj)
			}
		}
	}
	if v, ok := in["audit_config"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.AuditConfig = &v1beta1.AuditConfig{}

		if v, ok := v["audit_policy"].([]interface{}); ok && len(v) > 0 {
			v := v[0].(map[string]interface{})
			obj.AuditConfig.AuditPolicy = &v1beta1.AuditPolicy{}

			if v, ok := v["config_map_ref"].([]interface{}); ok {
				obj.AuditConfig.AuditPolicy.ConfigMapRef = expandLocalObjectReference(v)
			}
		}
	}

	return obj
}

func expandMaintenance(maintenance []interface{}) *v1beta1.Maintenance {
	obj := &v1beta1.Maintenance{}

	if len(maintenance) == 0 || maintenance[0] == nil {
		return obj
	}
	in := maintenance[0].(map[string]interface{})

	if v, ok := in["auto_update"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.AutoUpdate = &v1beta1.MaintenanceAutoUpdate{}

		if v, ok := v["kubernetes_version"].(bool); ok {
			obj.AutoUpdate.KubernetesVersion = v
		}
	}
	if v, ok := in["time_window"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.TimeWindow = &v1beta1.MaintenanceTimeWindow{}

		if v, ok := v["begin"].(string); ok && len(v) > 0 {
			obj.TimeWindow.Begin = v
		}
		if v, ok := v["end"].(string); ok && len(v) > 0 {
			obj.TimeWindow.End = v
		}
	}

	return obj
}
func AddMissingDataForUpdate(old *v1beta1.Shoot, new *v1beta1.Shoot) {
	if new.Spec.DNS.Domain == nil || *new.Spec.DNS.Domain == "" {
		new.Spec.DNS.Domain = old.Spec.DNS.Domain
	}
	new.ObjectMeta.ResourceVersion = old.ObjectMeta.ResourceVersion
	new.ObjectMeta.Finalizers = old.ObjectMeta.Finalizers
}
