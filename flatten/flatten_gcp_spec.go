package flatten

import gcpAlpha1 "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"

func flattenGcpControlPlane(in gcpAlpha1.ControlPlaneConfig) []interface{} {
	att := make(map[string]interface{})

	if len(in.Zone) > 0 {
		att["zone"] = in.Zone
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
