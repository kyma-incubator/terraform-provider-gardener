package gardener

import (
	"fmt"
	"strconv"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	gardner_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "k8s.io/apimachinery/pkg/util/intstr"
)

func resourceServerCreate(d *schema.ResourceData, m interface{}, provider string) error {
	client := m.(*GardenerClient)
	name := d.Get("name").(string)
	d.SetId(name)
	shoots := client.GardenerClientSet.Shoots(client.NameSpace)
	shoot, err := shoots.Create(createCRD(d, client, provider))
	if err != nil {
		d.SetId("")
		return err
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

func resourceServerUpdate(d *schema.ResourceData, m interface{}, provider string) error {
	client := m.(*GardenerClient)
	name := d.Get("name").(string)
	d.SetId(name)
	shoots := client.GardenerClientSet.Shoots(client.NameSpace)
	shoot, err := shoots.Update(createCRD(d, client, provider))
	if err != nil {
		d.SetId("")
		return err
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
func createCRD(d *schema.ResourceData, client *GardenerClient, provider string) *gardner_types.Shoot {
	var internal gardencorev1alpha1.CIDR = "10.250.112.0/22" // TODO replace hardcoded
	var cidr gardencorev1alpha1.CIDR = "10.250.112.0/16"
	name := d.Get("name").(string)
	domain := name + "." + client.DNSBase
	allowPrivilegedContainers := true
	spec := gardner_types.ShootSpec{
		Cloud: gardner_types.Cloud{
			Profile: provider,
			Region:  d.Get("region").(string),
		},
		Kubernetes: gardner_types.Kubernetes{
			Version:                   d.Get("kubernetesversion").(string),
			AllowPrivilegedContainers: &allowPrivilegedContainers,
		},
		DNS: gardner_types.DNS{
			Domain: &domain,
		},
	}
	switch provider {
	case gcp:
		spec = createGCPSpec(spec, d, client.SecretBindings.GcpSecretBinding)
		spec.Cloud.GCP.Networks.Internal = &internal
	case aws:
		spec = createAWSSpec(spec, d, client.SecretBindings.AwsSecretBinding)
		spec.Cloud.AWS.Networks.VPC.CIDR = &cidr
	}
	//TODO check if secret binding is empty then return error
	return &gardner_types.Shoot{
		TypeMeta:   meta_v1.TypeMeta{Kind: "Shoot", APIVersion: "garden.sapcloud.io/v1beta1"},
		ObjectMeta: meta_v1.ObjectMeta{Name: name, Namespace: client.NameSpace},
		Spec:       spec,
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
func getGCPWorkers(d *schema.ResourceData) []gardner_types.GCPWorker {
	numWorkers := d.Get("worker.#").(int)
	resultWorkers := make([]gardner_types.GCPWorker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		var worker = "worker." + strconv.Itoa(i)
		resultWorkers[i] = gardner_types.GCPWorker{
			Worker:     createGardenWorker(worker, d),
			VolumeSize: d.Get(worker + ".volumesize").(string),
			VolumeType: d.Get(worker + ".volumetype").(string),
		}
	}
	return resultWorkers
}

func getAWSWorkers(d *schema.ResourceData) []gardner_types.AWSWorker {
	numWorkers := d.Get("worker.#").(int)
	resultWorkers := make([]gardner_types.AWSWorker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		var worker = "worker." + strconv.Itoa(i)
		resultWorkers[i] = gardner_types.AWSWorker{
			Worker:     createGardenWorker(worker, d),
			VolumeSize: d.Get(worker + ".volumesize").(string),
			VolumeType: d.Get(worker + ".volumetype").(string),
		}
	}
	return resultWorkers
}
func createGardenWorker(workerindex string, d *schema.ResourceData) gardner_types.Worker {
	return gardner_types.Worker{
		Name:          d.Get(workerindex + ".name").(string),
		MachineType:   d.Get(workerindex + ".machinetype").(string),
		AutoScalerMin: d.Get(workerindex + ".autoscalermin").(int),
		AutoScalerMax: d.Get(workerindex + ".autoscalermax").(int),
		MaxSurge: &util.IntOrString{
			IntVal: int32(d.Get(workerindex + ".maxsurge").(int)),
		},
		MaxUnavailable: &util.IntOrString{
			IntVal: int32(d.Get(workerindex + ".maxunavailable").(int)),
		},
	}
}
func createGCPSpec(spec gardner_types.ShootSpec, d *schema.ResourceData, secretBinding string) gardner_types.ShootSpec {
	spec.Cloud.SecretBindingRef.Name = secretBinding
	spec.Cloud.GCP = &gardner_types.GCPCloud{
		Networks: gardner_types.GCPNetworks{
			Workers: []gardencorev1alpha1.CIDR{"10.250.0.0/19"}, // TODO replace hardcoded
		},
		Workers: getGCPWorkers(d),
		Zones:   getZones(d),
	}
	return spec
}

func createAWSSpec(spec gardner_types.ShootSpec, d *schema.ResourceData, secretBinding string) gardner_types.ShootSpec {
	spec.Cloud.SecretBindingRef.Name = secretBinding
	spec.Cloud.AWS = &gardner_types.AWSCloud{
		Networks: gardner_types.AWSNetworks{
			Workers:  []gardencorev1alpha1.CIDR{"10.250.0.0/19"}, // TODO replace hardcoded
			Internal: []gardencorev1alpha1.CIDR{"10.250.112.0/22"},
			Public:   []gardencorev1alpha1.CIDR{"10.250.96.0/22"},
		},
		Workers: getAWSWorkers(d),
		Zones:   getZones(d),
	}
	return spec
}
