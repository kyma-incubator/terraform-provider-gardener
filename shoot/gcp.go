package shoot

import (
	//"strconv"

	//gardener_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

const gcp string = "gcp"

func GCPShoot() *schema.Resource {
	return &schema.Resource{
		Create: resourceGCPServerCreate,
		Read:   resourceServerRead,
		Exists: resourceServerExists,
		Update: resourceGCPServerUpdate,
		Delete: resourceServerDelete,
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("shoot", false),
			"spec":     shootSpecSchema(),
		},
	}
}
func resourceGCPServerCreate(d *schema.ResourceData, m interface{}) error {
	return resourceServerCreate(d, m, gcp)
}

func resourceGCPServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerUpdate(d, m, gcp)
}

