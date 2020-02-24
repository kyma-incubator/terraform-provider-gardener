package expand

import (
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func expandProvider(provider []interface{}) corev1beta1.Provider {
	obj := corev1beta1.Provider{}
	if len(provider) == 0 || provider[0] == nil {
		return obj
	}
	in := provider[0].(map[string]interface{})

	if v, ok := in["type"].(string); ok && len(v) > 0 {
		obj.Type = v
	}

	if v, ok := in["control_plane_config"].([]interface{}); ok && len(v) > 0 {
		cloud := v[0].(map[string]interface{})
		if az, ok := cloud["azure"].([]interface{}); ok && len(az) > 0 {
			obj.ControlPlaneConfig = azControlPlaneConfig()
		}

		if gcp, ok := cloud["gcp"].([]interface{}); ok && len(gcp) > 0 {
			obj.ControlPlaneConfig = gcpControlPlaneConfig(gcp)
		}
	}

	if v, ok := in["infrastructure_config"].([]interface{}); ok && len(v) > 0 {
		cloud := v[0].(map[string]interface{})
		if az, ok := cloud["azure"].([]interface{}); ok && len(az) > 0 {
			obj.InfrastructureConfig = azureConfig(az)
		}
		if aws, ok := cloud["aws"].([]interface{}); ok && len(aws) > 0 {
			obj.InfrastructureConfig = awsConfig(aws)
		}
		if gcp, ok := cloud["gcp"].([]interface{}); ok && len(gcp) > 0 {
			obj.InfrastructureConfig = gcpConfig(gcp)
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

	return obj
}

func expandWorker(w interface{}) corev1beta1.Worker {
	obj := corev1beta1.Worker{}
	if w == nil {
		return obj
	}
	in := w.(map[string]interface{})

	if v, ok := in["annotations"].(map[string]interface{}); ok {
		obj.Annotations = expandStringMap(v)
	}

	if v, ok := in["cabundle"].(string); ok && len(v) > 0 {
		obj.CABundle = &v
	}

	if v, ok := in["kubernetes"].([]interface{}); ok && len(v) > 0 {
		obj.Kubernetes = expandWorkerKubernetes(v)
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

	if v, ok := in["maximum"].(int); ok {
		obj.Maximum = int32(v)
	}

	if v, ok := in["minimum"].(int); ok {
		obj.Minimum = int32(v)
	}

	if v, ok := in["max_surge"].(int); ok {
		surge := intstr.FromInt(v)
		obj.MaxSurge = &surge
	}

	if v, ok := in["max_unavailable"].(int); ok {
		unavailable := intstr.FromInt(v)
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

	if v, ok := in["zones"].(*schema.Set); ok {
		obj.Zones = expandSet(v)
	}

	return obj
}

func expandWorkerKubernetes(wk []interface{}) *corev1beta1.WorkerKubernetes {
	obj := &corev1beta1.WorkerKubernetes{}
	if len(wk) == 0 && wk[0] == nil {
		return obj
	}
	in := wk[0].(map[string]interface{})

	if v, ok := in["kubelet"].([]interface{}); ok {
		obj.Kubelet = expandKubelet(v)
	}

	return obj
}

func expandMachine(m []interface{}) corev1beta1.Machine {
	obj := corev1beta1.Machine{}

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

func expandVolume(v []interface{}) *corev1beta1.Volume {
	obj := &corev1beta1.Volume{}

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
	return obj
}

func expandKubelet(kubelet []interface{}) *corev1beta1.KubeletConfig {
	obj := &corev1beta1.KubeletConfig{}

	if len(kubelet) == 0 && kubelet[0] == nil {
		return obj
	}
	in := kubelet[0].(map[string]interface{})

	if v, ok := in["pod_pids_limit"].(int); ok && v > 0 {
		limit := int64(v)
		obj.PodPIDsLimit = &limit
	}
	if v, ok := in["cpu_cfs_quota"].(bool); ok {
		obj.CPUCFSQuota = &v
	}
	if v, ok := in["cpu_manager_policy"].(string); ok && len(v) > 0 {
		obj.CPUManagerPolicy = &v
	}

	return obj
}

func expandShootImage(si []interface{}) *corev1beta1.ShootMachineImage {
	obj := &corev1beta1.ShootMachineImage{}

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
