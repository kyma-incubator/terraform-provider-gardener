package shoot

import (
	"fmt"
	"strings"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	gardner_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"

	//pkgApi "k8s.io/apimachinery/pkg/types"

	//"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kyma-incubator/terraform-provider-gardener/client"

	//"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "k8s.io/apimachinery/pkg/util/intstr"
)

func resourceServerCreate(d *schema.ResourceData, m interface{}, provider string) error {
	client := m.(*client.Client)
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
	client := m.(*client.Client)
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
	client := m.(*client.Client)
	name := d.Get("name").(string)
	d.SetId(name)
	shoots := client.GardenerClientSet.Shoots(client.NameSpace)
	shoot, err := shoots.Get(name, meta_v1.GetOptions{})
	//spec2 := createGCPSpec(spec, d, client.SecretBindings.GcpSecretBinding)
	shoot, err = GetUpdatedSpec(d, shoot)
	shoot, err = shoots.Update(shoot)
	if err != nil {
		d.SetId("")
		return err
	}
	fmt.Println(shoot)
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	name := d.Get("name").(string)
	shoots := client.GardenerClientSet.Shoots(client.NameSpace)
	shoots.Delete(name, &meta_v1.DeleteOptions{})
	d.SetId("")
	return nil
}

//createCRD return deployment structure
func createCRD(d *schema.ResourceData, client *client.Client, provider string) *gardner_types.Shoot {

	name := d.Get("name").(string)
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
	}
	switch provider {
	case gcp:
		spec = createGCPSpec(spec, d, client.SecretBindings.GcpSecretBinding)
	case aws:
		spec = createAWSSpec(spec, d, client.SecretBindings.AwsSecretBinding)
		var cidr = gardencorev1alpha1.CIDR(d.Get("vpccidr").(string))
		spec.Cloud.AWS.Networks.VPC.CIDR = &cidr
	case azure:
		spec = createAzureSpec(spec, d, client.SecretBindings.AzureSecretBinding)
		var cidr = gardencorev1alpha1.CIDR(d.Get("vnetcidr").(string))
		spec.Cloud.Azure.Networks.VNet.CIDR = &cidr
	}
	annotations := make(map[string]string)
	annotations["confirmation.garden.sapcloud.io/deletion"] = "true"
	return &gardner_types.Shoot{
		TypeMeta:   meta_v1.TypeMeta{Kind: "Shoot", APIVersion: "garden.sapcloud.io/v1beta1"},
		ObjectMeta: meta_v1.ObjectMeta{Name: name, Namespace: client.NameSpace, Annotations: annotations},
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
func getCidrs(property string, d *schema.ResourceData) []gardencorev1alpha1.CIDR {
	cidrSet := d.Get(property).(*schema.Set)
	cidrInterface := cidrSet.List()
	cidr := make([]gardencorev1alpha1.CIDR, len(cidrInterface))
	for i, v := range cidrInterface {
		cidr[i] = gardencorev1alpha1.CIDR(v.(string))
	}
	return cidr
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
func flattenShoot(d *schema.ResourceData, shoot *gardner_types.Shoot) error {
	d.Set("name", shoot.Name)
	i := strings.Index(shoot.Namespace, "-")
	if i > -1 {
		profile := shoot.Namespace[i+1:]
		d.Set("profile", profile)
	} else {
		err := fmt.Errorf("Unexpected Namespace format")
		return err
	}
	flattenSpec(d, &shoot.Spec)
	return nil
}
func flattenSpec(d *schema.ResourceData, spec *gardner_types.ShootSpec) {
	cloud := spec.Cloud
	d.Set("region", cloud.Region)
	d.Set("kubernetesversion ", spec.Kubernetes.Version)
	if cloud.GCP != nil {
		d.Set("gcp_secret_binding", cloud.SecretBindingRef.Name)
		d.Set("zones", cloud.GCP.Zones)
		d.Set("workerscidr", cloud.GCP.Networks.Workers)
		SetGCPWorkersFromShoot(d, cloud.GCP.Workers)
	}
	if cloud.AWS != nil {
		d.Set("aws_secret_binding", cloud.SecretBindingRef.Name)
		d.Set("zones", cloud.AWS.Zones)
		d.Set("workerscidr", cloud.AWS.Networks.Workers)
		d.Set("internalscidr", cloud.AWS.Networks.Internal)
		d.Set("publicscidr", cloud.AWS.Networks.Public)
		d.Set("vpccidr", cloud.AWS.Networks.VPC.CIDR)
		SetAWSWorkersFromShoot(d, cloud.AWS.Workers)
	}
	if cloud.Azure != nil {
		d.Set("azure_secret_binding", cloud.SecretBindingRef.Name)
		d.Set("workerscidr", cloud.AWS.Networks.Workers)
		d.Set("vnetcidr", cloud.Azure.Networks.VNet.CIDR)
		SetAzureWorkersFromShoot(d, cloud.Azure.Workers)
	}

}
func GetUpdatedSpec(d *schema.ResourceData, shoot *gardner_types.Shoot) (*gardner_types.Shoot, error) {
	if d.HasChange("name") {
		return nil, fmt.Errorf("Can not change the name")
	}
	if d.HasChange("region") {
		shoot.Spec.Cloud.Region = d.Get("region").(string)
	}
	if d.HasChange("kubernetesversion") {
		shoot.Spec.Kubernetes.Version = d.Get("kubernetesversion").(string)
	}
	if shoot.Spec.Cloud.GCP != nil {
		shoot.Spec.Cloud.GCP = SetGCPChanges(d, shoot.Spec.Cloud.GCP)
	}
	if shoot.Spec.Cloud.Azure != nil {
		shoot.Spec.Cloud.Azure = SetAzureChanges(d, shoot.Spec.Cloud.Azure)
	}
	if shoot.Spec.Cloud.AWS != nil {
		shoot.Spec.Cloud.AWS = SetAWSChanges(d, shoot.Spec.Cloud.AWS)
	}
	return shoot, nil
}
