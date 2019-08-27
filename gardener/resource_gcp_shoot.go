package gardener

import (
	"github.com/hashicorp/terraform/helper/schema"

	//"log"
	"fmt"
	"strconv"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"

	gardner_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "k8s.io/apimachinery/pkg/util/intstr"
)

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
			"kubernetesversion": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zones": &schema.Schema{
				Type: schema.TypeSet,
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
	_, err := shoots.Get(name, meta_v1.GetOptions{})
	if err != nil {
		d.SetId("")
		return err
	}
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
	allowPrivilegedContainers := true
	return &gardner_types.Shoot{
		TypeMeta:   meta_v1.TypeMeta{Kind: "Shoot", APIVersion: "garden.sapcloud.io/v1beta1"},
		ObjectMeta: meta_v1.ObjectMeta{Name: name, Namespace: client.NameSpace},
		Spec: gardner_types.ShootSpec{
			Cloud: gardner_types.Cloud{
				Profile: "gcp",
				Region:  d.Get("region").(string),
				SecretBindingRef: corev1.LocalObjectReference{
					Name: client.SecretBindings.GcpSecretBinding,
				},
				GCP: &gardner_types.GCPCloud{
					Networks: gardner_types.GCPNetworks{
						Internal: &internal,
						Workers:  []gardencorev1alpha1.CIDR{"10.250.0.0/19"}, // TODO replace hardcoded
					},
					Workers: getWorkers(d),
					Zones:   getZones(d),
				},
			},
			Kubernetes: gardner_types.Kubernetes{
				Version:                   d.Get("kubernetesversion").(string),
				AllowPrivilegedContainers: &allowPrivilegedContainers,
			},
			DNS: gardner_types.DNS{
				Domain: &domain,
			},
		},
	}

}

func getZones(d *schema.ResourceData) []string {
	zonesSet := d.Get("zones").(*schema.Set)
	zonesInterface := zonesSet.List()
	zones := make([]string, len(zonesInterface))
	for i, v := range zonesInterface {
		zones[i] = v.(string)
	}
	return zones
}
func getWorkers(d *schema.ResourceData) []gardner_types.GCPWorker {
	numWorkers := d.Get("worker.#").(int)
	resultWorkers := make([]gardner_types.GCPWorker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		var worker = "worker." + strconv.Itoa(i)
		resultWorkers[i] = gardner_types.GCPWorker{
			Worker: gardner_types.Worker{
				Name:          d.Get(worker + ".name").(string),
				MachineType:   d.Get(worker + ".machinetype").(string),
				AutoScalerMin: d.Get(worker + ".autoscalermin").(int),
				AutoScalerMax: d.Get(worker + ".autoscalermax").(int),
				MaxSurge: &util.IntOrString{
					IntVal: int32(d.Get(worker + ".maxsurge").(int)),
				},
				MaxUnavailable: &util.IntOrString{
					IntVal: int32(d.Get(worker + ".maxunavailable").(int)),
				},
			},
			VolumeSize: d.Get(worker + ".volumesize").(string),
			VolumeType: d.Get(worker + ".volumetype").(string),
		}
	}
	return resultWorkers
}
