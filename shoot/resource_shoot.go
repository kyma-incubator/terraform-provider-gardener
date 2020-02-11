package shoot

import (
	"fmt"
	"time"

	//"encoding/json"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardener_apis "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	"github.com/hashicorp/terraform/helper/mutexkv"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kyma-incubator/terraform-provider-gardener/client"
	"github.com/kyma-incubator/terraform-provider-gardener/expand"
	"github.com/kyma-incubator/terraform-provider-gardener/flatten"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultCreateTimeout = time.Minute * 30
	defaultUpdateTimeout = time.Minute * 30
	defaultDeleteTimeout = time.Minute * 20
)

// Shoot mutex prevents concurrent writes to the CRD
var shootMutex = mutexkv.NewMutexKV()

func ResourceShoot() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Exists: resourceServerExists,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServerImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultCreateTimeout),
			Update: schema.DefaultTimeout(defaultUpdateTimeout),
			Delete: schema.DefaultTimeout(defaultDeleteTimeout),
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("shoot", false),
			"spec":     shootSpecSchema(),
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	metadata := expand.ExpandMetadata(d.Get("metadata").([]interface{}))
	spec := expand.ExpandShoot(d.Get("spec").([]interface{}))

	mutex_key := fmt.Sprintf(`namespace-%s`, metadata.Namespace)
	shootMutex.Lock(mutex_key)
	defer shootMutex.Unlock(mutex_key)
	shootCRD := gardencorev1beta1.Shoot{
		ObjectMeta: metadata,
		Spec:       spec,
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Shoot",
			APIVersion: "core.gardener.cloud/v1beta1",
		},
	}
	//foo, _:= json.Marshal(&shootCRD)
	//return fmt.Errorf("fffff: %v", string(foo))

	shootsClient := client.GardenerClientSet.Shoots(metadata.Namespace)
	shoot, err := shootsClient.Create(&shootCRD)

	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId(flatten.BuildID(shoot.ObjectMeta))

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), waitForShootFunc(shootsClient, metadata.Name))
	if err != nil {
		d.SetId("")
		return err
	}
	return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	namespace, name, err := flatten.IdParts(d.Id())
	if err != nil {
		return err
	}

	shootsClient := client.GardenerClientSet.Shoots(namespace)
	shoot, err := shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		d.SetId("")
		return err
	}
	shoot.ObjectMeta.Annotations["confirmation.garden.sapcloud.io/deletion"] = "true"
	err = d.Set("metadata", flatten.FlattenMetadata(shoot.ObjectMeta, d))
	if err != nil {
		return err
	}

	spec, err := flatten.FlattenShoot(shoot.Spec, d)
	if err != nil {
		return err
	}

	err = d.Set("spec", spec)
	if err != nil {
		return err
	}
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	namespace, name, err := flatten.IdParts(d.Id())
	if err != nil {
		return err
	}

	mutex_key := fmt.Sprintf(`namespace-%s`, namespace)
	shootMutex.Lock(mutex_key)
	defer shootMutex.Unlock(mutex_key)

	shootsClient := client.GardenerClientSet.Shoots(namespace)
	shoot, err := shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Failed to get shoot: %s", err)
	}

	newShoot := gardencorev1beta1.Shoot{
		ObjectMeta: expand.ExpandMetadata(d.Get("metadata").([]interface{})),
		Spec:       expand.ExpandShoot(d.Get("spec").([]interface{})),
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Shoot",
			APIVersion: "core.gardener.cloud/v1beta1",
		},
	}
	expand.AddMissingDataForUpdate(shoot, &newShoot)

	_, err = shootsClient.Update(&newShoot)
	if err != nil {
		d.SetId("")
		return err
	}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), waitForShootFunc(shootsClient, name))
	if err != nil {
		return err
	}

	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)
	namespace, name, err := flatten.IdParts(d.Id())
	if err != nil {
		return err
	}
	mutex_key := fmt.Sprintf(`namespace-%s`, namespace)
	shootMutex.Lock(mutex_key)
	defer shootMutex.Unlock(mutex_key)
	shootsClient := client.GardenerClientSet.Shoots(namespace)
	err = shootsClient.Delete(name, &meta_v1.DeleteOptions{})
	if err != nil {
		return err
	}

	err = resource.Retry(d.Timeout(schema.TimeoutDelete), waitForShootDeleteFunc(shootsClient, name))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceServerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	client := m.(*client.Client)
	namespace, name, err := flatten.IdParts(d.Id())
	if err != nil {
		return nil, err
	}
	mutex_key := fmt.Sprintf(`namespace-%s`, namespace)
	shootMutex.Lock(mutex_key)
	defer shootMutex.Unlock(mutex_key)
	shootsClient := client.GardenerClientSet.Shoots(namespace)

	// Wait for cluster if it is still not ready
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), waitForShootFunc(shootsClient, name))
	if err != nil {
		d.SetId("")
		return nil, err
	}
	// Set ID
	shoot, err := shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	d.SetId(flatten.BuildID(shoot.ObjectMeta))

	return []*schema.ResourceData{d}, nil
}

func resourceServerExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*client.Client)
	namespace, name, err := flatten.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	shootsClient := client.GardenerClientSet.Shoots(namespace)
	_, err = shootsClient.Get(name, meta_v1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, err
		}
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
				if condition.Status == gardencorev1beta1.ConditionProgressing {
					return resource.RetryableError(fmt.Errorf("Waiting for shoot condition to finish: %s", condition.Type))
				}
				if condition.Status == gardencorev1beta1.ConditionFalse {
					return resource.RetryableError(fmt.Errorf("Shoot condition failed: %s", condition.Message))
				}
			}

			if shoot.Status.LastOperation.State == gardencorev1beta1.LastOperationStatePending || shoot.Status.LastOperation.State == gardencorev1beta1.LastOperationStateProcessing {
				return resource.RetryableError(fmt.Errorf("Waiting for last operation to finish: %s", shoot.Status.LastOperation.Description))
			}
			if shoot.Status.LastOperation.State == gardencorev1beta1.LastOperationStateAborted || shoot.Status.LastOperation.State == gardencorev1beta1.LastOperationStateError || shoot.Status.LastOperation.State == gardencorev1beta1.LastOperationStateFailed {
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
