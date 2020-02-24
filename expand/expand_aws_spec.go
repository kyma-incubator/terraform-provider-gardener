package expand

import (
	"encoding/json"
	awsAlpha1 "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

func awsConfig(aws []interface{}) *corev1beta1.ProviderConfig {
	awsConfigObj := awsAlpha1.InfrastructureConfig{}
	obj := corev1beta1.ProviderConfig{}
	if len(aws) == 0 && aws[0] == nil {
		return &obj
	}
	in := aws[0].(map[string]interface{})

	awsConfigObj.APIVersion = "aws.provider.extensions.gardener.cloud/v1alpha1"
	awsConfigObj.Kind = "InfrastructureConfig"
	if v, ok := in["enableecraccess"].(bool); ok {
		awsConfigObj.EnableECRAccess = &v
	}
	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
		awsConfigObj.Networks = awsNetworks(v)
	}
	obj.Raw, _ = json.Marshal(awsConfigObj)
	return &obj
}

func awsNetworks(networks []interface{}) awsAlpha1.Networks {
	obj := awsAlpha1.Networks{}
	if networks == nil {
		return obj
	}
	in := networks[0].(map[string]interface{})

	if v, ok := in["vpc"].([]interface{}); ok {
		obj.VPC = vpc(v)
	}

	if v, ok := in["zones"].(*schema.Set); ok {
		obj.Zones = expandAwsZones(v)
	}

	return obj
}

func expandAwsZones(set *schema.Set) []awsAlpha1.Zone {
	result := make([]awsAlpha1.Zone, set.Len())
	for i, k := range set.List() {
		z := awsAlpha1.Zone{}
		if v, ok := k.(map[string]interface{})["name"].(string); ok && len(v) > 0 {
			z.Name = v
		}
		if v, ok := k.(map[string]interface{})["internal"].(string); ok && len(v) > 0 {
			z.Internal = v
		}
		if v, ok := k.(map[string]interface{})["public"].(string); ok && len(v) > 0 {
			z.Public = v
		}
		if v, ok := k.(map[string]interface{})["workers"].(string); ok && len(v) > 0 {
			z.Workers = v
		}

		result[i] = z
	}
	return result
}

func vpc(vpc []interface{}) awsAlpha1.VPC {
	obj := awsAlpha1.VPC{}

	if len(vpc) == 0 && vpc[0] == nil {
		return obj
	}
	in := vpc[0].(map[string]interface{})

	if v, ok := in["id"].(string); ok && len(v) > 0 {
		obj.ID = &v
	}
	if v, ok := in["cidr"].(string); ok && len(v) > 0 {
		obj.CIDR = &v
	}
	return obj
}
