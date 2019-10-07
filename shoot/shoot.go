package shoot

import (
	"fmt"
	"log"

	//"strings"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	gardener_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	gardener_apis "github.com/gardener/gardener/pkg/client/garden/clientset/versioned/typed/garden/v1beta1"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kyma-incubator/terraform-provider-gardener/client"
	"github.com/kyma-incubator/terraform-provider-gardener/helper"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//pkgApi "k8s.io/apimachinery/pkg/types"
	//util "k8s.io/apimachinery/pkg/util/intstr"
)

func resourceServerCreate(d *schema.ResourceData, m interface{}, provider string) error {
	client := m.(*client.Client)
	metadata := helper.ExpandMetadata(d.Get("metadata").([]interface{}))
	spec := helper.ExpandShoot(d.Get("spec").([]interface{}))

	shootCRD := gardener_types.Shoot{
		ObjectMeta: metadata,
		Spec:       spec,
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Shoot",
			APIVersion: "garden.sapcloud.io/v1beta1",
		},
	}

	log.Printf("[INFO] Creating new shoot: %#v", shootCRD)
	shootsClient := client.GardenerClientSet.Shoots(metadata.Namespace)
	shoot, err := shootsClient.Create(&shootCRD)
	if err != nil {
		return err
	}
	d.SetId(helper.BuildID(shoot.ObjectMeta))
	resource.Retry(d.Timeout(schema.TimeoutCreate),
		waitForShootFunc(shootsClient, metadata.Name))
	if err != nil {
		d.SetId("")
		return err
	}
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	namespace, name, err := helper.IdParts(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading Shoot %s", name)
	shootsClient := client.GardenerClientSet.Shoots(namespace)
	shoot, err := shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		d.SetId("")
		return err
	}

	err = d.Set("metadata", helper.FlattenMetadata(shoot.ObjectMeta, d))
	if err != nil {
		return err
	}

	spec, err := helper.FlattenShoot(shoot.Spec, d)
	if err != nil {
		return err
	}

	err = d.Set("spec", spec)
	if err != nil {
		return err
	}
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}, provider string) error {
	client := m.(*client.Client)
	namespace, name, err := helper.IdParts(d.Id())
	if err != nil {
		return err
	}
	shootsClient := client.GardenerClientSet.Shoots(namespace)
	shoot,err := shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Failed to get shoot: %s", err)
	}
	if d.HasChange("metadata") {
		shoot.ObjectMeta = helper.ExpandMetadata(d.Get("metadata").([]interface{}))
	}
	if d.HasChange("spec") {
		shoot.Spec = helper.ExpandShoot(d.Get("spec").([]interface{}))
	}
	_,err = shootsClient.Update(shoot)
	if err != nil {
		log.Printf("[INFO] Error while updating shoot cluster: %#v", err)
		d.SetId("")
		return err
	}
	resource.Retry(d.Timeout(schema.TimeoutCreate),
		waitForShootFunc(shootsClient, name))
	log.Printf("[INFO] Submitted updated Shoot: %#v", d.Id())
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	namespace, name, err := helper.IdParts(d.Id())
	shootsClient := client.GardenerClientSet.Shoots(namespace)
	err = shootsClient.Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}
	resource.Retry(d.Timeout(schema.TimeoutDelete),
		waitForShootDeleteFunc(shootsClient, name))
	d.SetId("")
	return nil
}

func resourceServerExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*client.Client)
	namespace, name, err := helper.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	log.Printf("[INFO] Checking Shoot %s", name)
	shootsClient := client.GardenerClientSet.Shoots(namespace)
	_, err = shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
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
