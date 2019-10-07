package shoot

import (
	//"strconv"

	//gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	//gardener_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

const azure string = "az"

func AzureShoot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureServerCreate,
		Read:   resourceServerRead,
		Exists: resourceServerExists,
		Update: resourceAzureServerUpdate,
		Delete: resourceServerDelete,
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("shoot", false),
			"spec":     shootSpecSchema(),
		},
	}
}
func resourceAzureServerCreate(d *schema.ResourceData, m interface{}) error {
	return resourceServerCreate(d, m, azure)
}

func resourceAzureServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerUpdate(d, m, azure)
}

