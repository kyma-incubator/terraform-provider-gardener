package expand

import (
	"encoding/json"
	gcpAlpha1 "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
)

func getGCPControlPlaneConfig(gcp []interface{}) *corev1beta1.ProviderConfig {
	gcpConfigObj := gcpAlpha1.ControlPlaneConfig{}
	obj := corev1beta1.ProviderConfig{}
	if len(gcp) == 0 && gcp[0] == nil {
		return &obj
	}
	in := gcp[0].(map[string]interface{})

	gcpConfigObj.APIVersion = "gcp.provider.extensions.gardener.cloud/v1alpha1"
	gcpConfigObj.Kind = "ControlPlaneConfig"

	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		gcpConfigObj.Zone = v
	}

	obj.Raw, _ = json.Marshal(gcpConfigObj)
	return &obj
}

func getGCPConfig(gcp []interface{}) *corev1beta1.ProviderConfig {
	gcpConfigObj := gcpAlpha1.InfrastructureConfig{}
	obj := corev1beta1.ProviderConfig{}

	if len(gcp) == 0 && gcp[0] == nil {
		return &obj
	}
	in := gcp[0].(map[string]interface{})

	gcpConfigObj.APIVersion = "gcp.provider.extensions.gardener.cloud/v1alpha1"
	gcpConfigObj.Kind = "InfrastructureConfig"

	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		gcpConfigObj.Networks = getGCPNetworks(v)
	}
	obj.Raw, _ = json.Marshal(gcpConfigObj)
	return &obj
}

func getGCPNetworks(networks []interface{}) gcpAlpha1.NetworkConfig {
	obj := gcpAlpha1.NetworkConfig{}
	if networks == nil {
		return obj
	}
	in := networks[0].(map[string]interface{})
	if v, ok := in["vpc"].([]interface{}); ok && len(v) > 0 {
		obj.VPC = getGCPvpc(v)
	}
	if v, ok := in["workers"].(string); ok && len(v) > 0 {
		obj.Workers = v
	}
	if v, ok := in["internal"].(string); ok && len(v) > 0 {
		obj.Internal = &v
	}
	if v, ok := in["cloud_nat"].([]interface{}); ok && len(v) > 0 {
		obj.CloudNAT = getGCPCloudNat(v)
	}
	if v, ok := in["flow_logs"].([]interface{}); ok && len(v) > 0 {
		obj.FlowLogs = getGCPFlowLogs(v)
	}

	return obj
}

func getGCPFlowLogs(fl []interface{}) *gcpAlpha1.FlowLogs {
	obj := gcpAlpha1.FlowLogs{}
	if len(fl) == 0 && fl[0] == nil {
		return &obj
	}
	in := fl[0].(map[string]interface{})

	if v, ok := in["aggregation_interval"].(string); ok && len(v) > 0 {
		obj.AggregationInterval = &v
	}
	if v, ok := in["flow_sampling"].(float32); ok {
		f := float32(v)
		obj.FlowSampling = &f
	}
	if v, ok := in["metadata"].(string); ok {
		obj.Metadata = &v
	}
	return &obj

}

func getGCPCloudNat(cn []interface{}) *gcpAlpha1.CloudNAT {
	obj := gcpAlpha1.CloudNAT{}
	if len(cn) == 0 && cn[0] == nil {
		return &obj
	}

	in := cn[0].(map[string]interface{})

	if v, ok := in["min_ports_per_vm"].(int); ok {
		f := int32(v)
		obj.MinPortsPerVM = &f
	}
	return &obj
}

func getGCPvpc(vpc []interface{}) *gcpAlpha1.VPC {
	obj := gcpAlpha1.VPC{}
	if len(vpc) == 0 && vpc[0] == nil {
		return &obj
	}
	in := vpc[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}
	if v, ok := in["cloud_router"].([]interface{}); ok && len(v) > 0 {
		obj.CloudRouter = getGCPCloudRouter(v)
	}
	return &obj
}

func getGCPCloudRouter(cr []interface{}) *gcpAlpha1.CloudRouter {

	obj := gcpAlpha1.CloudRouter{}
	if len(cr) == 0 && cr[0] == nil {
		return &obj
	}
	in := cr[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}
	return &obj
}
