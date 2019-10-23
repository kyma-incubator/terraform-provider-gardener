package expand

import (
	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

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
