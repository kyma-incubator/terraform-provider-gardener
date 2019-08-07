package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/client-go/tools/clientcmd"

	//"log"
	"fmt"
	//  appsv1 "k8s.io/api/apps/v1"
	//  apiv1 "k8s.io/api/core/v1"
	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"

	//apierrors "k8s.io/apimachinery/pkg/api/errors"
	gardner_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	gardner_apis "github.com/gardener/gardener/pkg/client/garden/clientset/versioned/typed/garden/v1beta1"
	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "k8s.io/apimachinery/pkg/util/intstr"
)

// type command struct {
// 	opts *Options
// 	core.Command
// }

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"profile": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_binding": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zones": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"workers": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"machinetype": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"autoscalermin": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"autoscalermax": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"maxsurge": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"maxunavailable": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumesize": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"volumetype": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
		},
	}
}
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	d.SetId(name)
	profile := d.Get("profile").(string)
	nameSpace := "garden-" + profile
	kubeconfig := "C:\\Users\\D074188\\.kube\\config"
	fmt.Println(kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := gardner_apis.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	shoots := clientset.Shoots("garden-icke")
	shoot, err := shoots.Create(CreateCRD(name, profile, nameSpace))
	if err != nil {
		panic(err)
	}
	fmt.Println(shoot)
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	d.SetId(name)
	profile := d.Get("profile").(string)
	nameSpace := "garden-" + profile
	fmt.Println(d.Get("workers"))
	kubeconfig := "C:\\Users\\D074188\\.kube\\config"
	fmt.Println(kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := gardner_apis.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	shoots := clientset.Shoots(nameSpace)
	shoot, err := shoots.Update(CreateCRD(name, profile, nameSpace))
	if err != nil {
		panic(err)
	}
	fmt.Println(shoot)
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	d.SetId("")
	kubeconfig := "C:\\Users\\D074188\\.kube\\config"
	fmt.Println(kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := gardner_apis.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	shoots := clientset.Shoots("garden-icke")
	shoots.Delete(name, &meta_v1.DeleteOptions{})
	d.SetId("")
	return nil
}

//CreateCRD return deployment structure
func CreateCRD(name string, profile string, nameSpace string) *gardner_types.Shoot {
	var internal gardencorev1alpha1.CIDR = "10.250.112.0/22"
	domain := name + "." + profile + ".shoot.canary.k8s-hana.ondemand.com"
	//domain := "johndoe-gcp.garden-dev.exam"ple.com"
	allowPrivilegedContainers := true
	return &gardner_types.Shoot{
		TypeMeta:   meta_v1.TypeMeta{Kind: "Shoot", APIVersion: "garden.sapcloud.io/v1beta1"},
		ObjectMeta: meta_v1.ObjectMeta{Name: name, Namespace: nameSpace},
		Spec: gardner_types.ShootSpec{
			Cloud: gardner_types.Cloud{
				Profile: "gcp",
				Region:  "europe-west3",
				SecretBindingRef: corev1.LocalObjectReference{
					Name: "icke-architecture",
				},
				GCP: &gardner_types.GCPCloud{
					Networks: gardner_types.GCPNetworks{
						Internal: &internal,
						Workers:  []gardencorev1alpha1.CIDR{"10.250.0.0/19"},
					},
					Workers: []gardner_types.GCPWorker{
						gardner_types.GCPWorker{
							Worker: gardner_types.Worker{
								Name:          "cpu-worker",
								MachineType:   "n1-standard-4",
								AutoScalerMin: 2,
								AutoScalerMax: 2,
								MaxSurge: &util.IntOrString{
									IntVal: 1,
								},
								MaxUnavailable: &util.IntOrString{
									IntVal: 0,
								},
							},
							VolumeSize: "20Gi",
							VolumeType: "pd-standard",
						},
					},
					Zones: []string{"europe-west3-b"},
				},
			},
			Kubernetes: gardner_types.Kubernetes{
				Version:                   "1.15.2",
				AllowPrivilegedContainers: &allowPrivilegedContainers,
			},
			DNS: gardner_types.DNS{
				Domain: &domain,
			},
		},
	}

}
