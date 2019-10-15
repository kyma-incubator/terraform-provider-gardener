package shoot

import (
	"fmt"
	"log"
	"strings"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	gardener_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	gardener_apis "github.com/gardener/gardener/pkg/client/garden/clientset/versioned/typed/garden/v1beta1"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kyma-incubator/terraform-provider-gardener/client"

	"github.com/hashicorp/terraform/helper/resource"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	util "k8s.io/apimachinery/pkg/util/intstr"
)

func resourceServerCreate(d *schema.ResourceData, m interface{}, provider string) error {
	client := m.(*client.Client)
	name := d.Get("name").(string)
	d.SetId(name)
	shootsClient := client.GardenerClientSet.Shoots(client.NameSpace)
	_, err := shootsClient.Create(createCRD(d, client, provider))
	if err != nil {
		return err
	}
	resource.Retry(d.Timeout(schema.TimeoutCreate),
		waitForShootFunc(shootsClient, name))
	if err != nil {
		d.SetId("")
		return err
	}
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	name := d.Get("name").(string)
	shootsClient := client.GardenerClientSet.Shoots(client.NameSpace)
	shoot, err := shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		d.SetId("")
		return err
	}
	flattenShoot(d, shoot)
	d.SetId(name)
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}, provider string) error {
	client := m.(*client.Client)
	name := d.Get("name").(string)
	d.SetId(name)
	shootsClient := client.GardenerClientSet.Shoots(client.NameSpace)
	shoot, err := shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}
	err = updateShootSpecFromConfig(d, shoot)
	if err != nil {
		return err
	}
	_, err = shootsClient.Update(shoot)
	if err != nil {
		d.SetId("")
		return err
	}
	resource.Retry(d.Timeout(schema.TimeoutCreate),
		waitForShootFunc(shootsClient, name))
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	name := d.Get("name").(string)
	shootsClient := client.GardenerClientSet.Shoots(client.NameSpace)
	err := shootsClient.Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}
	resource.Retry(d.Timeout(schema.TimeoutDelete),
		waitForShootDeleteFunc(shootsClient, name))
	d.SetId("")
	return nil
}

//createCRD return deployment structure
func createCRD(d *schema.ResourceData, client *client.Client, provider string) *gardener_types.Shoot {

	name := d.Get("name").(string)
	allowPrivilegedContainers := true
	spec := gardener_types.ShootSpec{
		Cloud: gardener_types.Cloud{
			Profile: provider,
			Region:  d.Get("region").(string),
		},
		Kubernetes: gardener_types.Kubernetes{
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
	return &gardener_types.Shoot{
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
func createGardenWorker(workerindex string, d *schema.ResourceData) gardener_types.Worker {
	return gardener_types.Worker{
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
func flattenShoot(d *schema.ResourceData, shoot *gardener_types.Shoot) error {
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
func flattenSpec(d *schema.ResourceData, spec *gardener_types.ShootSpec) {
	cloud := spec.Cloud
	d.Set("region", cloud.Region)
	d.Set("kubernetesversion ", spec.Kubernetes.Version)
	if cloud.GCP != nil {
		d.Set("gcp_secret_binding", cloud.SecretBindingRef.Name)
		d.Set("zones", cloud.GCP.Zones)
		d.Set("workerscidr", cloud.GCP.Networks.Workers)
		flattenGCPWorkers(d, cloud.GCP.Workers)
	}
	if cloud.AWS != nil {
		d.Set("aws_secret_binding", cloud.SecretBindingRef.Name)
		d.Set("zones", cloud.AWS.Zones)
		d.Set("workerscidr", cloud.AWS.Networks.Workers)
		d.Set("internalscidr", cloud.AWS.Networks.Internal)
		d.Set("publicscidr", cloud.AWS.Networks.Public)
		d.Set("vpccidr", cloud.AWS.Networks.VPC.CIDR)
		flattenAWSWorkers(d, cloud.AWS.Workers)
	}
	if cloud.Azure != nil {
		d.Set("azure_secret_binding", cloud.SecretBindingRef.Name)
		d.Set("workerscidr", cloud.Azure.Networks.Workers)
		d.Set("vnetcidr", cloud.Azure.Networks.VNet.CIDR)
		flattenAzureWorkers(d, cloud.Azure.Workers)
	}

}

func resourceServerExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*client.Client)
	name := d.Id()
	if name == "" {
		return false, fmt.Errorf("name is not set")
	}
	log.Printf("[INFO] Checking Shoot %s", name)
	shootsClient := client.GardenerClientSet.Shoots(client.NameSpace)
	_, err := shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}

func updateShootSpecFromConfig(d *schema.ResourceData, shoot *gardener_types.Shoot) error {
	if d.HasChange("name") {
		return fmt.Errorf("Can not change the name")
	}
	if d.HasChange("region") {
		shoot.Spec.Cloud.Region = d.Get("region").(string)
	}
	if d.HasChange("kubernetesversion") {
		shoot.Spec.Kubernetes.Version = d.Get("kubernetesversion").(string)
	}
	if shoot.Spec.Cloud.GCP != nil {
		updateGCPSpec(d, shoot.Spec.Cloud.GCP)
	}
	if shoot.Spec.Cloud.Azure != nil {
		updateAzureSpec(d, shoot.Spec.Cloud.Azure)
	}
	if shoot.Spec.Cloud.AWS != nil {
		updateAWSSpec(d, shoot.Spec.Cloud.AWS)
	}
	return nil
}

func waitForShootFunc(shootsClient gardener_apis.ShootInterface, name string) resource.RetryFunc {
	return func() *resource.RetryError {
		// Query the shoot to get a status update.
		shoot, err := shootsClient.Get(name, meta_v1.GetOptions{})
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if shoot.Generation <= shoot.Status.ObservedGeneration {
			for _, condition := range shoot.Status.Conditions {
				if condition.Status == gardencorev1alpha1.ConditionProgressing {
					return resource.RetryableError(fmt.Errorf("Waiting for shoot condition to finish: %s", condition.Type))
				}
				if condition.Status == gardencorev1alpha1.ConditionFalse {
					return resource.NonRetryableError(fmt.Errorf("Shoot condition failed: %s", condition.Message))
				}
			}

			if shoot.Status.LastOperation.State == gardencorev1alpha1.LastOperationStatePending || shoot.Status.LastOperation.State == gardencorev1alpha1.LastOperationStateProcessing {
				return resource.RetryableError(fmt.Errorf("Waiting for last operation to finish: %s", shoot.Status.LastOperation.Description))
			}
			if shoot.Status.LastOperation.State == gardencorev1alpha1.LastOperationStateAborted || shoot.Status.LastOperation.State == gardencorev1alpha1.LastOperationStateError || shoot.Status.LastOperation.State == gardencorev1alpha1.LastOperationStateFailed {
				return resource.NonRetryableError(fmt.Errorf("Shoot operation failed: %s", shoot.Status.LastOperation.Description))
			}
		} else {
			return resource.RetryableError(fmt.Errorf("Waiting for rollout to start"))
		}

		return nil
	}
}

func waitForShootDeleteFunc(shootsClient gardener_apis.ShootInterface, name string) resource.RetryFunc {
	return func() *resource.RetryError {
		// Query the shoot to get a status update.
		_, err := shootsClient.Get(name, meta_v1.GetOptions{})
		if err != nil {
			if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Cloud not get shoot state: %#v", err))
		}
		return resource.RetryableError(fmt.Errorf("Waiting for shoot to be deleted: %s", name))
	}
}
