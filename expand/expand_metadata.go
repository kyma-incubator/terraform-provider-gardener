package expand

import (
	"time"

	v1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ExpandMetadata(in []interface{}) metav1.ObjectMeta {
	meta := metav1.ObjectMeta{}
	if len(in) < 1 {
		return meta
	}
	m := in[0].(map[string]interface{})

	if v, ok := m["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Annotations = expandStringMap(m["annotations"].(map[string]interface{}))
	}
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string)
	}
	meta.Annotations["confirmation.garden.sapcloud.io/deletion"] = "true"
	if v, ok := m["labels"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Labels = expandStringMap(m["labels"].(map[string]interface{}))
	}

	if v, ok := m["generate_name"]; ok {
		meta.GenerateName = v.(string)
	}
	if v, ok := m["name"]; ok {
		meta.Name = v.(string)
	}
	if v, ok := m["namespace"]; ok {
		meta.Namespace = v.(string)
	}

	return meta
}

func expandLocalObjectReference(l []interface{}) *v1.LocalObjectReference {
	if len(l) == 0 || l[0] == nil {
		return &v1.LocalObjectReference{}
	}
	in := l[0].(map[string]interface{})
	obj := &v1.LocalObjectReference{}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	return obj
}

func expandDuration(v string) *v1beta1.GardenerDuration {
	d, err := time.ParseDuration(v)
	if err != nil {
		return &v1beta1.GardenerDuration{
			Duration: d,
		}
	}

	return nil
}

func expandCIDR(v string) *v1alpha1.CIDR {
	cidr := v1alpha1.CIDR(v)
	return &cidr
}

func expandSet(set *schema.Set) []string {
	result := make([]string, set.Len())
	for i, k := range set.List() {
		result[i] = k.(string)
	}

	return result
}

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		if v, ok := v.(string); ok {
			result[k] = v
		}
	}
	return result
}

func expandBoolMap(m map[string]interface{}) map[string]bool {
	result := make(map[string]bool)
	for k, v := range m {
		if v, ok := v.(bool); ok {
			result[k] = v
		}
	}
	return result
}

func ParseMap(aMap map[string]string, bMap map[string]interface{}) map[string]string {
	for key, _ := range aMap {
		if _, ok := bMap[key]; !ok {
			delete(aMap, key)
		}
	}
	return aMap
}
