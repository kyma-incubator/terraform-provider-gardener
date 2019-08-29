package shoot

import (
	"strconv"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	gardener_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

const aws string = "aws"

func AWSShoot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAWSServerCreate,
		Read:   resourceServerRead,
		Update: resourceAWSServerUpdate,
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
func resourceAWSServerCreate(d *schema.ResourceData, m interface{}) error {
	return resourceServerCreate(d, m, aws)
}

func resourceAWSServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerUpdate(d, m, aws)
}

func createAWSSpec(spec gardener_types.ShootSpec, d *schema.ResourceData, secretBinding string) gardener_types.ShootSpec {
	spec.Cloud.SecretBindingRef.Name = secretBinding
	spec.Cloud.AWS = &gardener_types.AWSCloud{
		Networks: gardener_types.AWSNetworks{
			Workers:  []gardencorev1alpha1.CIDR{"10.250.0.0/19"}, // TODO replace hardcoded
			Internal: []gardencorev1alpha1.CIDR{"10.250.112.0/22"},
			Public:   []gardencorev1alpha1.CIDR{"10.250.96.0/22"},
		},
		Workers: getAWSWorkers(d),
		Zones:   getZones(d),
	}
	return spec
}

func getAWSWorkers(d *schema.ResourceData) []gardener_types.AWSWorker {
	numWorkers := d.Get("worker.#").(int)
	resultWorkers := make([]gardener_types.AWSWorker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		var worker = "worker." + strconv.Itoa(i)
		resultWorkers[i] = gardener_types.AWSWorker{
			Worker:     createGardenWorker(worker, d),
			VolumeSize: d.Get(worker + ".volumesize").(string),
			VolumeType: d.Get(worker + ".volumetype").(string),
		}
	}
	return resultWorkers
}
