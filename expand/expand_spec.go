package expand

import (
	"encoding/json"
	azAlpha1 "github.com/gardener/gardener-extensions/controllers/provider-azure/pkg/apis/azure/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	//v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Expanders
func ExpandShoot(shoot []interface{}) corev1beta1.ShootSpec {
	//obj := v1beta1.ShootSpec{}
	obj := corev1beta1.ShootSpec{}

	if len(shoot) == 0 || shoot[0] == nil {
		return obj
	}
	in := shoot[0].(map[string]interface{})

	if v, ok := in["addons"].([]interface{}); ok && len(v) > 0 {
		obj.Addons = expandAddons(v)
	}
	if v, ok := in["cloud_profile_name"].(string); ok && len(v) > 0 {
		obj.CloudProfileName = v
	}
	//if v, ok := in["cloud"].([]interface{}); ok && len(v) > 0 {
	//	obj.Cloud = expandCloud(v)
	//}
	if v, ok := in["provider"].([]interface{}); ok && len(v) > 0 {
		obj.Provider = expandProvider(v)
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

	if v, ok := in["networking"].([]interface{}); ok && len(v) > 0 {
		obj.Networking = expandNetwokring(v)
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["secret_binding_name"].(string); ok && len(v) > 0 {
		obj.SecretBindingName = v
	}

	if v, ok := in["seed_name"].(string); ok && len(v) > 0 {
		obj.SeedName = &v
	}

	return obj
}

func expandNetwokring(networking []interface{}) corev1beta1.Networking {
	obj := corev1beta1.Networking{}
	if len(networking) == 0 || networking[0] == nil {
		return obj
	}

	in := networking[0].(map[string]interface{})
	if v, ok := in["type"].(string); ok && len(v) > 0 {
		obj.Type = v
	}

	if v, ok := in["pods"].(string); ok && len(v) > 0 {
		obj.Pods = &v
	}

	if v, ok := in["nodes"].(string); ok && len(v) > 0 {
		obj.Nodes = &v
	}

	if v, ok := in["services"].(string); ok && len(v) > 0 {
		obj.Services = &v
	}

	return  obj

}

func expandProvider(provider []interface{}) corev1beta1.Provider{
	obj := corev1beta1.Provider{}
	if  len(provider) ==  0 || provider[0]== nil {
		return obj
	}
	in := provider[0].(map[string]interface{})


	if v, ok := in["type"].(string); ok && len(v) > 0 {
		obj.Type = v
	}

	if v, ok := in["infrastructure_config"].([]interface{}); ok && len(v) > 0 {
		cloud := v[0].(map[string]interface{})
		if az, ok := cloud["azure"].([]interface{}); ok && len(az) > 0 {
			//obj.ControlPlaneConfig = getAzControlPlaneConfig()
			obj.InfrastructureConfig  = getAzureConfig(az)
		}
	}
	if workers, ok := in["worker"].([]interface{}); ok && len(workers) > 0 {
		for _, w := range workers {
			if w, ok := w.(map[string]interface{}); ok {
				workerObj := expandWorker(w)
				obj.Workers = append(obj.Workers, workerObj)
			}
		}
	}

	return  obj
}

func getAzControlPlaneConfig() *corev1beta1.ProviderConfig {
	azConfig := `
      apiVersion: azure.provider.extensions.gardener.cloud/v1alpha1
      kind: ControlPlaneConfig`
	obj := corev1beta1.ProviderConfig{}
	obj.Raw = []byte(azConfig)
	return  &obj
}



func getAzureConfig(az []interface{}) *corev1beta1.ProviderConfig {
	azConfigObj := azAlpha1.InfrastructureConfig{}
	obj := corev1beta1.ProviderConfig{}
	if len(az) == 0 && az[0] == nil {
		return &obj
	}
	in := az[0].(map[string]interface{})

	azConfigObj.APIVersion  = "azure.provider.extensions.gardener.cloud/v1alpha1"
	azConfigObj.Kind = "InfrastructureConfig"
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		azConfigObj.Networks = getNetworks(v)
	}
	obj.Raw, _ = json.Marshal(azConfigObj)
	return  &obj
}

func getNetworks(networks []interface{}) azAlpha1.NetworkConfig {
	obj := azAlpha1.NetworkConfig{}
	if networks == nil {
		return  obj
	}
	in := networks[0].(map[string]interface{})

	if v, ok := in["vnet"].([]interface{}); ok {
		obj.VNet = getVNET(v)
	}
	if v, ok := in["workers"].(string); ok {
		obj.Workers = v
	}

	return obj
}

func getVNET(vnet []interface{}) azAlpha1.VNet {
	obj := azAlpha1.VNet{}


	if vnet == nil {
		return  obj
	}

	in := vnet[0].(map[string]interface{})
	return obj
	if v, ok := in["cidr"].(string); ok && len(v) >0 {
		obj.CIDR = &v
	}

	if v, ok := in["name"].(string); ok && len(v) >0 {
		obj.Name = &v
	}

	if v, ok := in["resource_group"].(string); ok && len(v) >0 {
		obj.ResourceGroup = &v
	}

return obj
}

func expandWorker(w interface{}) corev1beta1.Worker {
obj := corev1beta1.Worker{}
if  w == nil {
	return obj
}
	in := w.(map[string]interface{})

	if v, ok := in["annotations"].(map[string]interface{}); ok {
		obj.Annotations = expandStringMap(v)
	}

	if v, ok := in["ca_bundle"].(string); ok && len(v) > 0 {
		obj.CABundle = &v
	}

	if v, ok := in["worker_kubernetes"].([]interface{}); ok && len(v) > 0 {
		obj.Kubernetes = expand_worker_kubernetes(v)
	}

	if v, ok := in["labels"].(map[string]interface{}); ok {
		obj.Labels = expandStringMap(v)
	}

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in["machine"].([]interface{}); ok && len(v) > 0 {
		obj.Machine = expandMachine(v)
	}

	if v, ok := in["maximum"].(int); ok   {
		obj.Maximum = int32(v)
	}

	if v, ok := in["minimum"].(int); ok   {
		obj.Minimum = int32(v)
	}

	if v, ok := in["max_surge"].(string); ok && len(v) > 0  {
		surge := intstr.FromString(v)
		obj.MaxSurge = &surge
	}

	if v, ok := in["max_unavailable"].(string); ok && len(v) > 0  {
		unavailable := intstr.FromString(v)
		obj.MaxUnavailable = &unavailable
	}

	if taints, ok := in["taints"].([]interface{}); ok && len(taints) > 0 {
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

				obj.Taints = append(obj.Taints, taint)
			}
		}
	}

	if v, ok := in["volume"].([]interface{}); ok {
		obj.Volume = expandVolume(v)
	}

	if zones, ok := in["zones"].([]string); ok && len(zones) > 0 {
		obj.Zones = zones
	}
return obj
}

func expandVolume(v []interface{}) *corev1beta1.Volume {
	obj  := &corev1beta1.Volume{}

	if len(v) == 0 && v[0] == nil {
		return obj
	}
	in := v[0].(map[string]interface{})

	if c, ok := in["type"].(string); ok && len(c) > 0 {
		obj.Type = &c
	}

	if c, ok := in["size"].(string); ok && len(c) > 0 {
		obj.Size = c
	}
	return  obj
}

func expandMachine(m []interface{}) corev1beta1.Machine {
	obj  := corev1beta1.Machine{}

	if len(m) == 0 && m[0] == nil {
		return obj
	}
	in := m[0].(map[string]interface{})

	if v, ok := in["type"].(string); ok && len(v) > 0 {
		obj.Type = v
	}

	if v, ok := in["image"].([]interface{}); ok && len(v) > 0 {
		obj.Image = expandShootImage(v)
	}

return obj
}

func expandShootImage(si []interface{}) *corev1beta1.ShootMachineImage {
	obj  := &corev1beta1.ShootMachineImage{}

	if len(si) == 0 && si[0] == nil {
		return obj
	}
	in := si[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in["version"].(string); ok && len(v) > 0 {
		obj.Version = v
	}
	return obj
}

func expand_worker_kubernetes(wk []interface{}) *corev1beta1.WorkerKubernetes {
	obj  := &corev1beta1.WorkerKubernetes{}
	if len(wk) == 0 && wk[0] == nil {
		return obj
	}
	in := wk[0].(map[string]interface{})

	if v, ok := in["kubelet"].(corev1beta1.KubeletConfig); ok {
		obj.Kubelet = &v
	}
	return obj

}

//func expandKubelet(kubelet []interface{}) *corev1beta1.KubeletConfig {
//	obj := &corev1beta1.KubeletConfig{}
//	if len(kubelet) == 0 && kubelet[0] == nil {
//		return  obj
//	}
//
//	in := kubelet[0].(map[string]interface{})
//	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0{
//		obj. = expandKubelet(v)
//	}
//
//	return  obj
//
//}

func expandAddons(addon []interface{}) *corev1beta1.Addons {
	obj := &corev1beta1.Addons{}

	if len(addon) == 0 || addon[0] == nil {
		return obj
	}
	in := addon[0].(map[string]interface{})

	if v, ok := in["kubernetes_dashboard"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubernetesDashboard = &corev1beta1.KubernetesDashboard{}

		if v, ok := v["enabled"].(bool); ok {
			obj.KubernetesDashboard.Enabled = v
		}
		if v, ok := v["authentication_mode"].(string); ok && len(v) > 0 {
			obj.KubernetesDashboard.AuthenticationMode = &v
		}
	}
	if v, ok := in["nginx_ingress"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.NginxIngress = &corev1beta1.NginxIngress{}

		if v, ok := v["enabled"].(bool); ok {
			obj.NginxIngress.Enabled = v
		}
		if v, ok := v["load_balancer_source_ranges"].(*schema.Set); ok {
			obj.NginxIngress.LoadBalancerSourceRanges = expandSet(v)
		}
	}
	//if v, ok := in["cluster_autoscaler"].([]interface{}); ok && len(v) > 0 {
	//	v := v[0].(map[string]interface{})
	//	obj.ClusterAutoscaler = &v1beta1.AddonClusterAutoscaler{}
	//
	//	if v, ok := v["enabled"].(bool); ok {
	//		obj.ClusterAutoscaler.Enabled = v
	//	}
	//}
	//if v, ok := in["heapster"].([]interface{}); ok && len(v) > 0 {
	//	v := v[0].(map[string]interface{})
	//	obj.Heapster = &v1beta1.Heapster{}
	//
	//	if v, ok := v["enabled"].(bool); ok {
	//		obj.Heapster.Enabled = v
	//	}
	//}
	//if v, ok := in["kube2iam"].([]interface{}); ok && len(v) > 0 {
	//	v := v[0].(map[string]interface{})
	//	obj.Kube2IAM = &v1beta1.Kube2IAM{}
	//
	//	if v, ok := v["enabled"].(bool); ok {
	//		obj.Kube2IAM.Enabled = v
	//	}
	//	if roles, ok := v["roles"].([]interface{}); ok {
	//		for _, r := range roles {
	//			if r, ok := r.(map[string]interface{}); ok {
	//				roleObj := v1beta1.Kube2IAMRole{}
	//
	//				if v, ok := r["name"].(string); ok && len(v) > 0 {
	//					roleObj.Name = v
	//				}
	//				if v, ok := r["description"].(string); ok && len(v) > 0 {
	//					roleObj.Description = v
	//				}
	//				if v, ok := r["policy"].(string); ok && len(v) > 0 {
	//					roleObj.Policy = v
	//				}
	//
	//				obj.Kube2IAM.Roles = append(obj.Kube2IAM.Roles, roleObj)
	//			}
	//		}
	//	}
	//}
	//if v, ok := in["kube_lego"].([]interface{}); ok && len(v) > 0 {
	//	v := v[0].(map[string]interface{})
	//	obj.KubeLego = &v1beta1.KubeLego{}
	//
	//	if v, ok := v["enabled"].(bool); ok {
	//		obj.KubeLego.Enabled = v
	//	}
	//	if v, ok := v["email"].(string); ok && len(v) > 0 {
	//		obj.KubeLego.Mail = v
	//	}
	//}
	//if v, ok := in["monocular"].([]interface{}); ok && len(v) > 0 {
	//	v := v[0].(map[string]interface{})
	//	obj.Monocular = &v1beta1.Monocular{}
	//
	//	if v, ok := v["enabled"].(bool); ok {
	//		obj.Monocular.Enabled = v
	//	}
	//}

	return obj
}

func expandDNS(dns []interface{}) *corev1beta1.DNS {
	obj := corev1beta1.DNS{}

	if len(dns) == 0 || dns[0] == nil {
		return &obj
	}
	in := dns[0].(map[string]interface{})

	if v, ok := in["providers"].([]corev1beta1.DNSProvider); ok && len(v) > 0 {
		obj.Providers = v
	}
	if v, ok := in["domain"].(string); ok && len(v) > 0 {
		obj.Domain = &v
	}
	//if v, ok := in["secret_name"].(string); ok && len(v) > 0 {
	//	obj.SecretName = &v
	//}

	return &obj
}

func expandHibernation(hibernation []interface{}) *corev1beta1.Hibernation {
	obj := &corev1beta1.Hibernation{}

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
				scheduleObj := corev1beta1.HibernationSchedule{}

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

func expandKubernetes(kubernetes []interface{}) corev1beta1.Kubernetes {
	obj := corev1beta1.Kubernetes{}

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
	//if v, ok := in["cloud_controller_manager"].([]interface{}); ok && len(v) > 0 {
	//	v := v[0].(map[string]interface{})
	//	obj.CloudControllerManager = &v1beta1.CloudControllerManagerConfig{}
	//
	//	if v, ok := v["feature_gates"].(map[string]interface{}); ok {
	//		obj.CloudControllerManager.FeatureGates = expandBoolMap(v)
	//	}
	//}
	if v, ok := in["kube_controller_manager"].([]interface{}); ok && len(v) > 0 {
		obj.KubeControllerManager = expandKubernetesControllerManager(v)
	}
	if v, ok := in["kube_scheduler"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubeScheduler = &corev1beta1.KubeSchedulerConfig{}

		if v, ok := v["feature_gates"].(map[string]interface{}); ok {
			obj.KubeScheduler.FeatureGates = expandBoolMap(v)
		}
	}
	if v, ok := in["kube_proxy"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.KubeProxy = &corev1beta1.KubeProxyConfig{}

		if v, ok := v["feature_gates"].(map[string]interface{}); ok {
			obj.KubeProxy.FeatureGates = expandBoolMap(v)
		}
		if v, ok := v["mode"].(string); ok && len(v) > 0 {
			mode := corev1beta1.ProxyModeIPTables
			obj.KubeProxy.Mode = &mode
		}
	}
	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.Kubelet = &corev1beta1.KubeletConfig{}

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
		obj.ClusterAutoscaler = &corev1beta1.ClusterAutoscaler{}

		if v, ok := v["scale_down_utilization_threshold"].(float64); ok {
			obj.ClusterAutoscaler.ScaleDownUtilizationThreshold = &v
		}
	}

	return obj
}

func expandKubernetesControllerManager(controller []interface{}) *corev1beta1.KubeControllerManagerConfig {
	obj := &corev1beta1.KubeControllerManagerConfig{}

	if len(controller) == 0 || controller[0] == nil {
		return obj
	}
	in := controller[0].(map[string]interface{})

	if v, ok := in["feature_gates"].(map[string]interface{}); ok {
		obj.FeatureGates = expandBoolMap(v)
	}
	if v, ok := in["horizontal_pod_autoscaler"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.HorizontalPodAutoscalerConfig = &corev1beta1.HorizontalPodAutoscalerConfig{}

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

func expandKubernetesAPIServer(server []interface{}) *corev1beta1.KubeAPIServerConfig {
	obj := &corev1beta1.KubeAPIServerConfig{}

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
		obj.OIDCConfig = &corev1beta1.OIDCConfig{}

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
	if v, ok := in["audit_config"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.AuditConfig = &corev1beta1.AuditConfig{}

		if v, ok := v["audit_policy"].([]interface{}); ok && len(v) > 0 {
			v := v[0].(map[string]interface{})
			obj.AuditConfig.AuditPolicy = &corev1beta1.AuditPolicy{}

			if v, ok := v["config_map_ref"].([]interface{}); ok {
				obj.AuditConfig.AuditPolicy.ConfigMapRef = expandObjectReference(v)
			}
		}
	}

	return obj
}

func expandMaintenance(maintenance []interface{}) *corev1beta1.Maintenance {
	obj := &corev1beta1.Maintenance{}

	if len(maintenance) == 0 || maintenance[0] == nil {
		return obj
	}
	in := maintenance[0].(map[string]interface{})

	if v, ok := in["auto_update"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.AutoUpdate = &corev1beta1.MaintenanceAutoUpdate{}

		if v, ok := v["kubernetes_version"].(bool); ok {
			obj.AutoUpdate.KubernetesVersion = v
		}
	}
	if v, ok := in["time_window"].([]interface{}); ok && len(v) > 0 {
		v := v[0].(map[string]interface{})
		obj.TimeWindow = &corev1beta1.MaintenanceTimeWindow{}

		if v, ok := v["begin"].(string); ok && len(v) > 0 {
			obj.TimeWindow.Begin = v
		}
		if v, ok := v["end"].(string); ok && len(v) > 0 {
			obj.TimeWindow.End = v
		}
	}

	return obj
}
func AddMissingDataForUpdate(old *corev1beta1.Shoot, new *corev1beta1.Shoot) {
	if new.Spec.DNS.Domain == nil || *new.Spec.DNS.Domain == "" {
		new.Spec.DNS.Domain = old.Spec.DNS.Domain
	}
	new.ObjectMeta.ResourceVersion = old.ObjectMeta.ResourceVersion
	new.ObjectMeta.Finalizers = old.ObjectMeta.Finalizers
}
func RemoveInternalKeysMapSpec(aMap map[string]interface{}, bMap map[string]interface{}) map[string]interface{} {
	for key, val := range aMap {
		switch val.(type) {
		case map[string]interface{}:
			if val2, ok := bMap[key]; !ok {
				delete(aMap, key)
			} else {
				aMap[key] = RemoveInternalKeysMapSpec(val.(map[string]interface{}), val2.(map[string]interface{}))
			}
		case []interface{}:
			if val2, ok := bMap[key]; !ok {
				delete(aMap, key)
			} else {
				aMap[key] = RemoveInternalKeysArraySpec(val.([]interface{}), val2.([]interface{}))
			}
		default:
			if val2, ok := bMap[key]; !ok || val2 == "" {
				delete(aMap, key)
			}
		}
	}
	return aMap
}
func RemoveInternalKeysArraySpec(ArrayA, ArrayB []interface{}) []interface{} {
	for i, val := range ArrayA {
		switch val.(type) {
		case map[string]interface{}:
			if i >= len(ArrayB) || ArrayB[i] == nil {
				ArrayA = remove(ArrayA, i)
			} else {
				ArrayA[i] = RemoveInternalKeysMapSpec(val.(map[string]interface{}), ArrayB[i].(map[string]interface{}))
			}
		case []interface{}:
			if i >= len(ArrayB) || ArrayB[i] == nil {
				ArrayA = remove(ArrayA, i)
			} else {
				ArrayA[i] = RemoveInternalKeysArraySpec(val.([]interface{}), ArrayB[i].([]interface{}))
			}
		default:
			if i >= len(ArrayB) || ArrayB[i] == nil || ArrayB[i] == "" {
				ArrayA = remove(ArrayA, i)
			}
		}
	}
	return ArrayA
}
func remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}
