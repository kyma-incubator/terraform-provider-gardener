// This file needs to be placed in this package, because of circular dependency to "flatten" package
package shoot

import (
	"encoding/json"
	"testing"

	gcpAlpha1 "github.com/gardener/gardener-extension-provider-gcp/pkg/apis/gcp/v1alpha1"

	awsAlpha1 "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"
	azAlpha1 "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/v1alpha1"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	corev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cmp "github.com/google/go-cmp/cmp"
	"github.com/kyma-incubator/terraform-provider-gardener/flatten"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestFlattenShoot(t *testing.T) {
	purpose := v1beta1.ShootPurposeEvaluation
	nodes := "10.250.0.0/19"
	pods := "100.96.0.0/11"
	services := "100.64.0.0/13"
	volumeType := "Standard_LRS"
	caBundle := "caBundle"
	policy := "policy"
	quota := true
	pdsLimit := int64(4)
	domain := "domain"
	authMode := "auth_mode"
	location := "Pacific/Auckland"
	hibernationStart := "00 17 * * 1"
	hibernationEnd := "00 00 * * 1"
	hibernationEnabled := true
	allowPrivilegedContainers := true
	enableBasicAuthentication := true
	clientID := "ClientID"
	groupsClaim := "GroupsClaim"
	groupsPrefix := "GroupsPrefix"
	issuerURL := "IssuerURL"
	usernameClaim := "UsernameClaim"
	usernamePrefix := "UsernamePrefix"

	d := ResourceShoot().TestResourceData()
	shoot := corev1beta1.ShootSpec{
		CloudProfileName:  "az",
		SecretBindingName: "test-secret",
		Purpose:           &purpose,
		Networking: corev1beta1.Networking{
			Nodes:    &nodes,
			Pods:     &pods,
			Services: &services,
			Type:     "calico",
		},
		Maintenance: &corev1beta1.Maintenance{
			AutoUpdate: &corev1beta1.MaintenanceAutoUpdate{
				KubernetesVersion:   true,
				MachineImageVersion: true,
			},
			TimeWindow: &corev1beta1.MaintenanceTimeWindow{
				Begin: "030000+0000",
				End:   "040000+0000",
			},
		},
		Provider: corev1beta1.Provider{
			Workers: []corev1beta1.Worker{
				corev1beta1.Worker{
					MaxSurge: &intstr.IntOrString{
						IntVal: 1,
					},
					MaxUnavailable: &intstr.IntOrString{
						IntVal: 0,
					},
					Maximum: 2,
					Minimum: 1,
					Volume: &corev1beta1.Volume{
						Size: "50Gi",
						Type: &volumeType,
					},
					Name: "cpu-worker",
					Machine: corev1beta1.Machine{
						Image: &corev1beta1.ShootMachineImage{
							Name:    "coreos",
							Version: "2303.3.0",
						},
						Type: "Standard_A4_v2",
					},
					Annotations: map[string]string{
						"test-key-annotation": "test-value-annotation",
					},
					Labels: map[string]string{
						"test-key-label": "test-value-label",
					},
					CABundle: &caBundle,
					Zones:    []string{"foo", "bar"},
					Taints: []corev1.Taint{
						corev1.Taint{
							Key:    "key",
							Value:  "value",
							Effect: corev1.TaintEffectNoExecute,
						},
					},
					Kubernetes: &corev1beta1.WorkerKubernetes{
						Kubelet: &corev1beta1.KubeletConfig{
							PodPIDsLimit:     &pdsLimit,
							CPUCFSQuota:      &quota,
							CPUManagerPolicy: &policy,
						},
					},
				},
			},
		},
		Region: "westeurope",
		Kubernetes: corev1beta1.Kubernetes{
			Version:                   "1.15.4",
			AllowPrivilegedContainers: &allowPrivilegedContainers,
			KubeAPIServer: &corev1beta1.KubeAPIServerConfig{
				EnableBasicAuthentication: &enableBasicAuthentication,
				OIDCConfig: &corev1beta1.OIDCConfig{
					CABundle:       &caBundle,
					ClientID:       &clientID,
					GroupsClaim:    &groupsClaim,
					GroupsPrefix:   &groupsPrefix,
					IssuerURL:      &issuerURL,
					RequiredClaims: map[string]string{"key": "value"},
					SigningAlgs:    []string{"bar", "foo"},
					UsernameClaim:  &usernameClaim,
					UsernamePrefix: &usernamePrefix,
				},
				AuditConfig: &corev1beta1.AuditConfig{
					AuditPolicy: &corev1beta1.AuditPolicy{
						ConfigMapRef: &corev1.ObjectReference{
							Name: "audit-policy",
						},
					},
				},
			},
		},
		DNS: &corev1beta1.DNS{
			Domain: &domain,
		},
		Addons: &corev1beta1.Addons{
			KubernetesDashboard: &corev1beta1.KubernetesDashboard{
				Addon: corev1beta1.Addon{
					Enabled: true,
				},
				AuthenticationMode: &authMode,
			},
			NginxIngress: &corev1beta1.NginxIngress{
				Addon: corev1beta1.Addon{
					Enabled: true,
				},
			},
		},
		Hibernation: &corev1beta1.Hibernation{
			Enabled: &hibernationEnabled,
			Schedules: []corev1beta1.HibernationSchedule{
				corev1beta1.HibernationSchedule{
					Start:    &hibernationStart,
					End:      &hibernationEnd,
					Location: &location,
				},
			},
		},
		Monitoring: &corev1beta1.Monitoring{
			Alerting: &corev1beta1.Alerting{
				EmailReceivers: []string{"receiver1", "receiver2"},
			},
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			"cloud_profile_name":  "az",
			"secret_binding_name": "test-secret",
			"purpose":             purpose,
			"region":              "westeurope",
			"networking": []interface{}{
				map[string]interface{}{
					"nodes":    nodes,
					"pods":     pods,
					"services": services,
					"type":     "calico",
				},
			},
			"kubernetes": []interface{}{
				map[string]interface{}{
					"version":                     "1.15.4",
					"allow_privileged_containers": true,
					"kube_api_server": []interface{}{
						map[string]interface{}{
							"enable_basic_authentication": true,
							"oidc_config": []interface{}{
								map[string]interface{}{
									"ca_bundle":       caBundle,
									"client_id":       clientID,
									"groups_claim":    groupsClaim,
									"groups_prefix":   groupsPrefix,
									"issuer_url":      issuerURL,
									"required_claims": map[string]string{"key": "value"},
									"signing_algs":    []string{"bar", "foo"},
									"username_claim":  usernameClaim,
									"username_prefix": usernamePrefix,
								},
							},
							"audit_config": []interface{}{
								map[string]interface{}{
									"audit_policy": []interface{}{
										map[string]interface{}{
											"config_map_ref": []interface{}{
												map[string]interface{}{
													"name": "audit-policy",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"maintenance": []interface{}{
				map[string]interface{}{
					"auto_update": []interface{}{
						map[string]interface{}{
							"kubernetes_version":    true,
							"machine_image_version": true,
						},
					},
					"time_window": []interface{}{
						map[string]interface{}{
							"begin": "030000+0000",
							"end":   "040000+0000",
						},
					},
				},
			},
			"provider": []interface{}{
				map[string]interface{}{
					"worker": []interface{}{
						map[string]interface{}{
							"name":            "cpu-worker",
							"max_surge":       1,
							"max_unavailable": 0,
							"maximum":         int32(2),
							"minimum":         int32(1),
							"volume": []interface{}{
								map[string]interface{}{
									"size": "50Gi",
									"type": "Standard_LRS",
								},
							},
							"machine": []interface{}{
								map[string]interface{}{
									"type": "Standard_A4_v2",
									"image": []interface{}{
										map[string]interface{}{
											"name":    "coreos",
											"version": "2303.3.0",
										},
									},
								},
							},
							"annotations": map[string]string{
								"test-key-annotation": "test-value-annotation",
							},
							"labels": map[string]string{
								"test-key-label": "test-value-label",
							},
							"zones":    []string{"foo", "bar"},
							"cabundle": caBundle,
							"taints": []interface{}{
								map[string]interface{}{
									"key":    "key",
									"value":  "value",
									"effect": corev1.TaintEffect("NoExecute"),
								},
							},
							"kubernetes": []interface{}{
								map[string]interface{}{
									"kubelet": []interface{}{
										map[string]interface{}{
											"pod_pids_limit":     int64(4),
											"cpu_cfs_quota":      true,
											"cpu_manager_policy": "policy",
										},
									},
								},
							},
						},
					},
				},
			},
			"dns": []interface{}{
				map[string]interface{}{
					"domain": domain,
				},
			},
			"addons": []interface{}{
				map[string]interface{}{
					"kubernetes_dashboard": []interface{}{
						map[string]interface{}{
							"enabled":             true,
							"authentication_mode": authMode,
						},
					},
					"nginx_ingress": []interface{}{
						map[string]interface{}{
							"enabled": true,
						},
					},
				},
			},
			"hibernation": []interface{}{
				map[string]interface{}{
					"enabled": hibernationEnabled,
					"schedules": []interface{}{
						map[string]interface{}{
							"start":    hibernationStart,
							"end":      hibernationEnd,
							"location": location,
						},
					},
				},
			},
			"monitoring": []interface{}{
				map[string]interface{}{
					"alerting": []interface{}{
						map[string]interface{}{
							"emailreceivers": []string{"receiver1", "receiver2"},
						},
					},
				},
			},
		},
	}
	d.Set("spec", expected)

	out, _ := flatten.FlattenShoot(shoot, d, "")
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenShootAzure(t *testing.T) {
	vnetCIDR := "10.250.0.0/16"
	vnetName := "test"
	resGroup := "test"
	azConfig, _ := json.Marshal(azAlpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: "azure.provider.extensions.gardener.cloud/v1alpha1",
			Kind:       "InfrastructureConfig",
		},
		Networks: azAlpha1.NetworkConfig{
			VNet: azAlpha1.VNet{
				CIDR:          &vnetCIDR,
				Name:          &vnetName,
				ResourceGroup: &resGroup,
			},
			Workers:          "10.250.0.0/19",
			ServiceEndpoints: []string{"microsoft.test"},
		},
	})

	d := ResourceShoot().TestResourceData()
	shoot := corev1beta1.ShootSpec{
		Provider: corev1beta1.Provider{
			Type: "azure",
			InfrastructureConfig: &corev1beta1.ProviderConfig{
				RawExtension: runtime.RawExtension{
					Raw: azConfig,
				},
			},
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			"kubernetes": []interface{}{},
			"networking": []interface{}{},
			"provider": []interface{}{
				map[string]interface{}{
					"type": "azure",
					"infrastructure_config": []interface{}{
						map[string]interface{}{
							"azure": []interface{}{
								map[string]interface{}{
									"networks": []interface{}{
										map[string]interface{}{
											"vnet": []interface{}{
												map[string]interface{}{
													"cidr":           "10.250.0.0/16",
													"name":           "test",
													"resource_group": "test",
												},
											},
											"workers":           "10.250.0.0/19",
											"service_endpoints": []string{"microsoft.test"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := d.Set("spec", expected)
	if err != nil {
		t.Fatalf("Unable to set the spec: %v\n", err)
	}

	out, _ := flatten.FlattenShoot(shoot, d, "")
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenShootAws(t *testing.T) {
	vpcCIDR := "10.250.0.0/16"
	EnableECRAccess := false
	awsConfig, _ := json.Marshal(awsAlpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: "aws.provider.extensions.gardener.cloud/v1alpha1",
			Kind:       "InfrastructureConfig",
		},
		EnableECRAccess: &EnableECRAccess,
		Networks: awsAlpha1.Networks{
			VPC: awsAlpha1.VPC{
				CIDR: &vpcCIDR,
			},
			Zones: []awsAlpha1.Zone{
				awsAlpha1.Zone{
					Name:     "eu-central-1a",
					Internal: vpcCIDR,
					Public:   vpcCIDR,
					Workers:  vpcCIDR,
				},
			},
		},
	})

	d := ResourceShoot().TestResourceData()
	shoot := corev1beta1.ShootSpec{
		Provider: corev1beta1.Provider{
			Type: "aws",
			InfrastructureConfig: &corev1beta1.ProviderConfig{
				RawExtension: runtime.RawExtension{
					Raw: awsConfig,
				},
			},
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			"kubernetes": []interface{}{},
			"networking": []interface{}{},
			"provider": []interface{}{
				map[string]interface{}{
					"type": "aws",
					"infrastructure_config": []interface{}{
						map[string]interface{}{
							"aws": []interface{}{
								map[string]interface{}{
									"enableecraccess": &EnableECRAccess,
									"networks": []interface{}{
										map[string]interface{}{
											"vpc": []interface{}{
												map[string]interface{}{
													"cidr": &vpcCIDR,
												},
											},
											"zones": []map[string]interface{}{
												{
													"name":     "eu-central-1a",
													"internal": "10.250.0.0/16",
													"public":   "10.250.0.0/16",
													"workers":  "10.250.0.0/16",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := d.Set("spec", expected)

	if err != nil {
		t.Fatalf("Error setting expected with spec: \n%s", err)
	}

	out, _ := flatten.FlattenShoot(shoot, d, "")
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}

func TestFlattenShootGCP(t *testing.T) {
	minPorts := int32(2)
	aggregationInterval := "test"
	metadata := "test"
	fs := float32(2)
	internal := "test"
	gcpControlPlaneConfig, _ := json.Marshal(gcpAlpha1.ControlPlaneConfig{
		Zone: "zone1",
	})
	gcpConfig, _ := json.Marshal(gcpAlpha1.InfrastructureConfig{
		TypeMeta: v1.TypeMeta{
			APIVersion: "gcp.provider.extensions.gardener.cloud/v1alpha1",
			Kind:       "InfrastructureConfig",
		},
		Networks: gcpAlpha1.NetworkConfig{
			VPC: &gcpAlpha1.VPC{
				CloudRouter: &gcpAlpha1.CloudRouter{
					Name: "bar",
				},
				Name: "foo",
			},
			Workers: "10.250.0.0/19",
			CloudNAT: &gcpAlpha1.CloudNAT{
				MinPortsPerVM: &minPorts,
			},
			Internal: &internal,
			FlowLogs: &gcpAlpha1.FlowLogs{
				AggregationInterval: &aggregationInterval,
				FlowSampling:        &fs,
				Metadata:            &metadata,
			},
		},
	})

	d := ResourceShoot().TestResourceData()
	shoot := corev1beta1.ShootSpec{
		Provider: corev1beta1.Provider{
			Type: "gcp",
			ControlPlaneConfig: &corev1beta1.ProviderConfig{
				RawExtension: runtime.RawExtension{
					Raw: gcpControlPlaneConfig,
				},
			},
			InfrastructureConfig: &corev1beta1.ProviderConfig{
				RawExtension: runtime.RawExtension{
					Raw: gcpConfig,
				},
			},
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			"kubernetes": []interface{}{},
			"networking": []interface{}{},
			"provider": []interface{}{
				map[string]interface{}{
					"type": "gcp",
					"control_plane_config": []interface{}{
						map[string]interface{}{
							"gcp": []interface{}{
								map[string]interface{}{
									"zone": "zone1",
								},
							},
						},
					},
					"infrastructure_config": []interface{}{
						map[string]interface{}{
							"gcp": []interface{}{
								map[string]interface{}{
									"networks": []interface{}{
										map[string]interface{}{
											"vpc": []interface{}{
												map[string]interface{}{
													"name": "foo",
													"cloud_router": []interface{}{
														map[string]interface{}{
															"name": "bar",
														},
													},
												},
											},
											"workers": "10.250.0.0/19",
											"cloud_nat": []interface{}{
												map[string]interface{}{
													"min_ports_per_vm": int32(2),
												},
											},
											"internal": "test",
											"flow_logs": []interface{}{
												map[string]interface{}{
													"aggregation_interval": "test",
													"flow_sampling":        float32(2),
													"metadata":             "test",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := d.Set("spec", expected)
	if err != nil {
		t.Fatalf("Unable to set the spec: %v\n", err)
	}

	out, _ := flatten.FlattenShoot(shoot, d, "")
	if diff := cmp.Diff(expected, out); diff != "" {
		t.Fatalf("Error matching output and expected: \n%s", diff)
	}
}
