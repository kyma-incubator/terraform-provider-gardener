package gardener

import (
	"github.com/hashicorp/terraform/helper/schema"

	//"log"
	"fmt"

	//  appsv1 "k8s.io/api/apps/v1"
	//  apiv1 "k8s.io/api/core/v1"
	//"github.com/aws/aws-sdk-go/aws/awsutil"
	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"

	//apierrors "k8s.io/apimachinery/pkg/api/errors"
	gardner_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "k8s.io/apimachinery/pkg/util/intstr"
)

// type command struct {
// 	opts *Options
// 	core.Command
// }

func resourceGCPShoot() *schema.Resource {
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
			"region": &schema.Schema{
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
			"worker": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
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
			},
		},
	}
}
func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*GardenerClient)
	name := d.Get("name").(string)
	d.SetId(name)
	shoots := client.GardenerClientSet.Shoots(client.NameSpace)
	shoot, err := shoots.Create(createCRD(d, client))
	if err != nil {
		panic(err)
	}
	fmt.Println(shoot)
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*GardenerClient)
	name := d.Get("name").(string)
	shoots := client.GardenerClientSet.Shoots(client.NameSpace)
	obj, err := shoots.Get(name, meta_v1.GetOptions{})
	if err != nil {
		d.SetId("")
		panic(err)
	}

	fmt.Println("obj: ", obj)
	d.SetId(name)
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*GardenerClient)
	name := d.Get("name").(string)
	d.SetId(name)
	shoots := client.GardenerClientSet.Shoots(client.NameSpace)
	shoot, err := shoots.Update(createCRD(d, client))
	if err != nil {
		panic(err)
	}
	fmt.Println(shoot)
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*GardenerClient)
	name := d.Get("name").(string)
	shoots := client.GardenerClientSet.Shoots(client.NameSpace)
	shoots.Delete(name, &meta_v1.DeleteOptions{})
	d.SetId("")
	return nil
}

//createCRD return deployment structure
func createCRD(d *schema.ResourceData, client *GardenerClient) *gardner_types.Shoot {
	var internal gardencorev1alpha1.CIDR = "10.250.112.0/22" // TODO replace hardcoded
	name := d.Get("name").(string)
	domain := name + "." + client.DNSBase
	region := d.Get("region").(string)
	allowPrivilegedContainers := true
	return &gardner_types.Shoot{
		TypeMeta:   meta_v1.TypeMeta{Kind: "Shoot", APIVersion: "garden.sapcloud.io/v1beta1"},
		ObjectMeta: meta_v1.ObjectMeta{Name: name, Namespace: client.NameSpace},
		Spec: gardner_types.ShootSpec{
			Cloud: gardner_types.Cloud{
				Profile: "gcp",
				Region:  region,
				SecretBindingRef: corev1.LocalObjectReference{
					Name: client.SecretBindings.GcpSecretBinding,
				},
				GCP: &gardner_types.GCPCloud{
					Networks: gardner_types.GCPNetworks{
						Internal: &internal,
						Workers:  []gardencorev1alpha1.CIDR{"10.250.0.0/19"}, // TODO replace hardcoded
					},
					Workers: []gardner_types.GCPWorker{ // TODO iterate on multiple workers
						gardner_types.GCPWorker{
							Worker: gardner_types.Worker{
								Name:          d.Get("worker.0.name").(string),
								MachineType:   d.Get("worker.0.machinetype").(string),
								AutoScalerMin: 2,
								AutoScalerMax: 2,
								MaxSurge: &util.IntOrString{
									IntVal: 1,
								},
								MaxUnavailable: &util.IntOrString{
									IntVal: 0,
								},
							},
							VolumeSize: d.Get("worker.0.volumesize").(string),
							VolumeType: d.Get("worker.0.volumetype").(string),
						},
					},
					Zones: []string{"europe-west3-b"}, //d.Get("zones").([]string), FIX THIS
				},
			},
			Kubernetes: gardner_types.Kubernetes{
				Version:                   "1.15.2", // TODO replace hardcoded
				AllowPrivilegedContainers: &allowPrivilegedContainers,
			},
			DNS: gardner_types.DNS{
				Domain: &domain,
			},
		},
	}

}
