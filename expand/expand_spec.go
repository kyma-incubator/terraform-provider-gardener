package expand

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"

	"github.com/hashicorp/terraform/helper/schema"
)

// Expanders
func ExpandShoot(shoot []interface{}) corev1beta1.ShootSpec {
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
	if v, ok := in["purpose"].(string); ok && len(v) > 0 {
		purpose := corev1beta1.ShootPurpose(v)
		obj.Purpose = &purpose
	}

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

	if v, ok := in["monitoring"].([]interface{}); ok && len(v) > 0 {
		obj.Monitoring = expandMonitoring(v)
	}

	if v, ok := in["networking"].([]interface{}); ok && len(v) > 0 {
		obj.Networking = expandNetworking(v)
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

func expandMonitoring(monitoring []interface{}) *corev1beta1.Monitoring {
	obj := corev1beta1.Monitoring{}
	if len(monitoring) == 0 || monitoring[0] == nil {
		return &obj
	}

	in := monitoring[0].(map[string]interface{})
	if v, ok := in["alerting"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		alert := corev1beta1.Alerting{}

		in := v[0].(map[string]interface{})
		if v, ok := in["emailreceivers"].(*schema.Set); ok {
			alert.EmailReceivers = expandSet(v)
		}

		obj.Alerting = &alert
	}

	return &obj
}

func expandNetworking(networking []interface{}) corev1beta1.Networking {
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

	return obj
}

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
	return obj
}

func expandDNS(dns []interface{}) *corev1beta1.DNS {
	obj := corev1beta1.DNS{}

	if len(dns) == 0 || dns[0] == nil {
		return &obj
	}
	in := dns[0].(map[string]interface{})

	// if v, ok := in["providers"].([]corev1beta1.DNSProvider); ok && len(v) > 0 {
	// 	obj.Providers = v
	// }
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

	if v, ok := in["enable_basic_authentication"].(bool); ok {
		obj.EnableBasicAuthentication = &v
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
		if v, ok := v["machine_image_version"].(bool); ok {
			obj.AutoUpdate.MachineImageVersion = v
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
	if new.Spec.DNS == nil {
		new.Spec.DNS = &corev1beta1.DNS{}
	}
	if new.Spec.DNS.Domain == nil || *new.Spec.DNS.Domain == "" {
		new.Spec.DNS.Domain = old.Spec.DNS.Domain
	}
	new.Spec.SeedName = old.Spec.SeedName
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
