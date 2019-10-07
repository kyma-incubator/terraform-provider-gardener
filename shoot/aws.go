package shoot

import (
	//"strconv"

	//gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	//gardener_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

const aws string = "aws"

func AWSShoot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAWSServerCreate,
		Read:   resourceServerRead,
		Exists: resourceServerExists,
		Update: resourceAWSServerUpdate,
		Delete: resourceServerDelete,
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("shoot", false),
			"spec":     shootSpecSchema(),
		},
	}
}
func resourceAWSServerCreate(d *schema.ResourceData, m interface{}) error {
	return resourceServerCreate(d, m, aws)
}

func resourceAWSServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerUpdate(d, m, aws)
}

