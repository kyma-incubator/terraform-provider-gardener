package expand

//
//import (
//	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
//	"github.com/hashicorp/terraform/helper/schema"
//	corev1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/util/intstr"
//)
//
//func expandCloud(cloud []interface{}) v1beta1.Cloud {
//	obj := v1beta1.Cloud{}
//
//	if len(cloud) == 0 || cloud[0] == nil {
//		return obj
//	}
//	in := cloud[0].(map[string]interface{})
//
//	if v, ok := in["profile"].(string); ok {
//		obj.Profile = v
//	}
//	if v, ok := in["region"].(string); ok {
//		obj.Region = v
//	}
//	if v, ok := in["secret_binding_ref"].([]interface{}); ok && len(v) > 0 {
//		obj.SecretBindingRef = *(expandLocalObjectReference(v))
//	}
//	if v, ok := in["seed"].(string); ok {
//		obj.Seed = &v
//	}
//	if v, ok := in["aws"].([]interface{}); ok && len(v) > 0 {
//		obj.AWS = expandCloudAWS(v)
//	} else if v, ok := in["gcp"].([]interface{}); ok && len(v) > 0 {
//		obj.GCP = expandCloudGCP(v)
//	} else if v, ok := in["azure"].([]interface{}); ok && len(v) > 0 {
//		obj.Azure = expandCloudAzure(v)
//	}
//
//	return obj
//}
//
//func expandCloudAWS(aws []interface{}) *v1beta1.AWSCloud {
//	obj := &v1beta1.AWSCloud{}
//
//	if len(aws) == 0 || aws[0] == nil {
//		return obj
//	}
//	in := aws[0].(map[string]interface{})
//
//	if v, ok := in["machine_image"].([]interface{}); ok && len(v) > 0 {
//		v := v[0].(map[string]interface{})
//		obj.MachineImage = &v1beta1.ShootMachineImage{}
//
//		if v, ok := v["name"].(string); ok && len(v) > 0 {
//			obj.MachineImage.Name = v
//		}
//		if v, ok := v["version"].(string); ok && len(v) > 0 {
//			obj.MachineImage.Version = v
//		}
//	}
//	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
//		v := v[0].(map[string]interface{})
//		obj.Networks = v1beta1.AWSNetworks{}
//
//		if v, ok := v["vpc"].([]interface{}); ok && len(v) > 0 {
//			v := v[0].(map[string]interface{})
//			obj.Networks.VPC = v1beta1.AWSVPC{}
//
//			if v, ok := v["id"].(string); ok && len(v) > 0 {
//				obj.Networks.VPC.ID = &v
//			}
//			if v, ok := v["cidr"].(string); ok && len(v) > 0 {
//				obj.Networks.VPC.CIDR = &v
//			}
//		}
//		if internal, ok := v["internal"].(*schema.Set); ok {
//			for _, i := range internal.List() {
//				obj.Networks.Internal = append(obj.Networks.Internal, i.(string))
//			}
//		}
//		if public, ok := v["public"].(*schema.Set); ok {
//			for _, p := range public.List() {
//				obj.Networks.Public = append(obj.Networks.Public, p.(string))
//			}
//		}
//		if workers, ok := v["workers"].(*schema.Set); ok {
//			for _, w := range workers.List() {
//				obj.Networks.Workers = append(obj.Networks.Workers, w.(string))
//			}
//		}
//	}
//	if workers, ok := in["worker"].([]interface{}); ok && len(workers) > 0 {
//		for _, w := range workers {
//			if w, ok := w.(map[string]interface{}); ok {
//				workerObj := v1beta1.AWSWorker{}
//
//				if v, ok := w["name"].(string); ok && len(v) > 0 {
//					workerObj.Name = v
//				}
//				if v, ok := w["machine_type"].(string); ok && len(v) > 0 {
//					workerObj.MachineType = v
//				}
//				if v, ok := w["auto_scaler_min"].(int); ok {
//					workerObj.AutoScalerMin = v
//				}
//				if v, ok := w["auto_scaler_max"].(int); ok {
//					workerObj.AutoScalerMax = v
//				}
//				if v, ok := w["max_surge"].(string); ok && len(v) > 0 {
//					surge := intstr.FromString(v)
//					workerObj.MaxSurge = &surge
//				}
//				if v, ok := w["max_unavailable"].(string); ok && len(v) > 0 {
//					surge := intstr.FromString(v)
//					workerObj.MaxUnavailable = &surge
//				}
//				if v, ok := w["annotations"].(map[string]interface{}); ok {
//					workerObj.Annotations = expandStringMap(v)
//				}
//				if v, ok := w["labels"].(map[string]interface{}); ok {
//					workerObj.Labels = expandStringMap(v)
//				}
//				if taints, ok := w["taints"].([]interface{}); ok && len(taints) > 0 {
//					for _, t := range taints {
//						if t, ok := t.(map[string]interface{}); ok {
//							taint := corev1.Taint{}
//
//							if v, ok := t["key"].(string); ok && len(v) > 0 {
//								taint.Key = v
//							}
//							// if v, ok := t["operator"].(string); ok && len(v) > 0 {
//							// 	taint.Operator = v
//							// }
//							if v, ok := t["value"].(string); ok && len(v) > 0 {
//								taint.Value = v
//							}
//							if v, ok := t["effect"].(string); ok && len(v) > 0 {
//								taint.Effect = corev1.TaintEffect(v)
//							}
//
//							workerObj.Taints = append(workerObj.Taints, taint)
//						}
//					}
//				}
//				if v, ok := w["volume_type"].(string); ok && len(v) > 0 {
//					workerObj.VolumeType = v
//				}
//				if v, ok := w["volume_size"].(string); ok && len(v) > 0 {
//					workerObj.VolumeSize = v
//				}
//
//				obj.Workers = append(obj.Workers, workerObj)
//			}
//		}
//	}
//	if v, ok := in["zones"].(*schema.Set); ok {
//		obj.Zones = expandSet(v)
//	}
//
//	return obj
//}
//func expandCloudGCP(aws []interface{}) *v1beta1.GCPCloud {
//	obj := &v1beta1.GCPCloud{}
//
//	if len(aws) == 0 || aws[0] == nil {
//		return obj
//	}
//	in := aws[0].(map[string]interface{})
//
//	if v, ok := in["machine_image"].([]interface{}); ok && len(v) > 0 {
//		v := v[0].(map[string]interface{})
//		obj.MachineImage = &v1beta1.ShootMachineImage{}
//
//		if v, ok := v["name"].(string); ok && len(v) > 0 {
//			obj.MachineImage.Name = v
//		}
//		if v, ok := v["version"].(string); ok && len(v) > 0 {
//			obj.MachineImage.Version = v
//		}
//	}
//	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
//		v := v[0].(map[string]interface{})
//		obj.Networks = v1beta1.GCPNetworks{}
//
//		if workers, ok := v["workers"].(*schema.Set); ok {
//			for _, w := range workers.List() {
//				obj.Networks.Workers = append(obj.Networks.Workers, w.(string))
//			}
//		}
//
//		if v, ok := v["internal"].(string); ok && len(v) > 0 {
//			obj.Networks.Internal = &v
//		}
//	}
//	if workers, ok := in["worker"].([]interface{}); ok && len(workers) > 0 {
//		for _, w := range workers {
//			if w, ok := w.(map[string]interface{}); ok {
//				workerObj := v1beta1.GCPWorker{}
//
//				if v, ok := w["name"].(string); ok && len(v) > 0 {
//					workerObj.Name = v
//				}
//				if v, ok := w["machine_type"].(string); ok && len(v) > 0 {
//					workerObj.MachineType = v
//				}
//				if v, ok := w["auto_scaler_min"].(int); ok {
//					workerObj.AutoScalerMin = v
//				}
//				if v, ok := w["auto_scaler_max"].(int); ok {
//					workerObj.AutoScalerMax = v
//				}
//				if v, ok := w["max_surge"].(string); ok && len(v) > 0 {
//					surge := intstr.FromString(v)
//					workerObj.MaxSurge = &surge
//				}
//				if v, ok := w["max_unavailable"].(string); ok && len(v) > 0 {
//					surge := intstr.FromString(v)
//					workerObj.MaxUnavailable = &surge
//				}
//				if v, ok := w["annotations"].(map[string]interface{}); ok {
//					workerObj.Annotations = expandStringMap(v)
//				}
//				if v, ok := w["labels"].(map[string]interface{}); ok {
//					workerObj.Labels = expandStringMap(v)
//				}
//				if taints, ok := w["taints"].([]interface{}); ok && len(taints) > 0 {
//					for _, t := range taints {
//						if t, ok := t.(map[string]interface{}); ok {
//							taint := corev1.Taint{}
//
//							if v, ok := t["key"].(string); ok && len(v) > 0 {
//								taint.Key = v
//							}
//							if v, ok := t["value"].(string); ok && len(v) > 0 {
//								taint.Value = v
//							}
//							if v, ok := t["effect"].(string); ok && len(v) > 0 {
//								taint.Effect = corev1.TaintEffect(v)
//							}
//
//							workerObj.Taints = append(workerObj.Taints, taint)
//						}
//					}
//				}
//				if v, ok := w["volume_type"].(string); ok && len(v) > 0 {
//					workerObj.VolumeType = v
//				}
//				if v, ok := w["volume_size"].(string); ok && len(v) > 0 {
//					workerObj.VolumeSize = v
//				}
//
//				obj.Workers = append(obj.Workers, workerObj)
//			}
//		}
//	}
//	if v, ok := in["zones"].(*schema.Set); ok {
//		obj.Zones = expandSet(v)
//	}
//
//	return obj
//}
//func expandCloudAzure(aws []interface{}) *v1beta1.AzureCloud {
//	obj := &v1beta1.AzureCloud{}
//
//	if len(aws) == 0 || aws[0] == nil {
//		return obj
//	}
//	in := aws[0].(map[string]interface{})
//
//	if v, ok := in["machine_image"].([]interface{}); ok && len(v) > 0 {
//		v := v[0].(map[string]interface{})
//		obj.MachineImage = &v1beta1.ShootMachineImage{}
//
//		if v, ok := v["name"].(string); ok && len(v) > 0 {
//			obj.MachineImage.Name = v
//		}
//		if v, ok := v["version"].(string); ok && len(v) > 0 {
//			obj.MachineImage.Version = v
//		}
//	}
//	if v, ok := in["networks"].([]interface{}); ok && len(v) > 0 {
//		v := v[0].(map[string]interface{})
//		obj.Networks = v1beta1.AzureNetworks{}
//
//		if v, ok := v["vnet"].([]interface{}); ok && len(v) > 0 {
//			v := v[0].(map[string]interface{})
//			obj.Networks.VNet = v1beta1.AzureVNet{}
//
//			if v, ok := v["name"].(string); ok && len(v) > 0 {
//				obj.Networks.VNet.Name = &v
//			}
//			if v, ok := v["cidr"].(string); ok && len(v) > 0 {
//				obj.Networks.VNet.CIDR = &v
//			}
//		}
//		if workers, ok := v["workers"].(string); ok {
//			obj.Networks.Workers = workers
//
//		}
//	}
//	if workers, ok := in["worker"].([]interface{}); ok && len(workers) > 0 {
//		for _, w := range workers {
//			if w, ok := w.(map[string]interface{}); ok {
//				workerObj := v1beta1.AzureWorker{}
//
//				if v, ok := w["name"].(string); ok && len(v) > 0 {
//					workerObj.Name = v
//				}
//				if v, ok := w["machine_type"].(string); ok && len(v) > 0 {
//					workerObj.MachineType = v
//				}
//				if v, ok := w["auto_scaler_min"].(int); ok {
//					workerObj.AutoScalerMin = v
//				}
//				if v, ok := w["auto_scaler_max"].(int); ok {
//					workerObj.AutoScalerMax = v
//				}
//				if v, ok := w["max_surge"].(string); ok && len(v) > 0 {
//					surge := intstr.FromString(v)
//					workerObj.MaxSurge = &surge
//				}
//				if v, ok := w["max_unavailable"].(string); ok && len(v) > 0 {
//					surge := intstr.FromString(v)
//					workerObj.MaxUnavailable = &surge
//				}
//				if v, ok := w["annotations"].(map[string]interface{}); ok {
//					workerObj.Annotations = expandStringMap(v)
//				}
//				if v, ok := w["labels"].(map[string]interface{}); ok {
//					workerObj.Labels = expandStringMap(v)
//				}
//				if taints, ok := w["taints"].([]interface{}); ok && len(taints) > 0 {
//					for _, t := range taints {
//						if t, ok := t.(map[string]interface{}); ok {
//							taint := corev1.Taint{}
//
//							if v, ok := t["key"].(string); ok && len(v) > 0 {
//								taint.Key = v
//							}
//							// if v, ok := t["operator"].(string); ok && len(v) > 0 {
//							// 	taint.Operator = v
//							// }
//							if v, ok := t["value"].(string); ok && len(v) > 0 {
//								taint.Value = v
//							}
//							if v, ok := t["effect"].(string); ok && len(v) > 0 {
//								taint.Effect = corev1.TaintEffect(v)
//							}
//
//							workerObj.Taints = append(workerObj.Taints, taint)
//						}
//					}
//				}
//				if v, ok := w["volume_type"].(string); ok && len(v) > 0 {
//					workerObj.VolumeType = v
//				}
//				if v, ok := w["volume_size"].(string); ok && len(v) > 0 {
//					workerObj.VolumeSize = v
//				}
//
//				obj.Workers = append(obj.Workers, workerObj)
//			}
//		}
//	}
//
//	return obj
//}
