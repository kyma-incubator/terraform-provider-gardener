package shoot

import (
	"strconv"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	gardener_types "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	"github.com/hashicorp/terraform/helper/schema"
)

const azure string = "az"

func AzureShoot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureServerCreate,
		Read:   resourceServerRead,
		Update: resourceAzureServerUpdate,
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
			"workercidr": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vnetcidr": &schema.Schema{
				Type:     schema.TypeString,
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
func resourceAzureServerCreate(d *schema.ResourceData, m interface{}) error {
	return resourceServerCreate(d, m, azure)
}

func resourceAzureServerUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceServerUpdate(d, m, azure)
}

func createAzureSpec(spec gardener_types.ShootSpec, d *schema.ResourceData, secretBinding string) gardener_types.ShootSpec {
	spec.Cloud.SecretBindingRef.Name = secretBinding
	spec.Cloud.Azure = &gardener_types.AzureCloud{
		Networks: gardener_types.AzureNetworks{
			Workers: gardencorev1alpha1.CIDR(d.Get("workercidr").(string)),
		},
		Workers: getAzureWorkers(d),
	}
	return spec
}

func getAzureWorkers(d *schema.ResourceData) []gardener_types.AzureWorker {
	numWorkers := d.Get("worker.#").(int)
	resultWorkers := make([]gardener_types.AzureWorker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		var worker = "worker." + strconv.Itoa(i)
		resultWorkers[i] = gardener_types.AzureWorker{
			Worker:     createGardenWorker(worker, d),
			VolumeSize: d.Get(worker + ".volumesize").(string),
			VolumeType: d.Get(worker + ".volumetype").(string),
		}
	}
	return resultWorkers
}
func SetAzureChanges(d *schema.ResourceData, azureSpec *gardener_types.AzureCloud) *gardener_types.AzureCloud {

	if d.HasChange("workerscidr") {
		azureSpec.Networks.Workers = gardencorev1alpha1.CIDR(d.Get("workercidr").(string))
	}
	var cidr = gardencorev1alpha1.CIDR(d.Get("vnetcidr").(string))
	azureSpec.Networks.VNet.CIDR = &cidr
	azureSpec.Workers = getAzureWorkers(d)

	return azureSpec
}
func SetAzureWorkersFromShoot(d *schema.ResourceData, workersarray []gardener_types.AzureWorker) {

	if len(workersarray) > 0 {
		workers := make([]interface{}, len(workersarray))
		for i, v := range workersarray {
			m := map[string]interface{}{}

			if v.Name != "" {
				m["name"] = v.Name
			}
			if v.MachineType != "" {
				m["machine_type"] = v.MachineType
			}
			if v.AutoScalerMin != 0 {
				m["auto_scaler_min"] = v.AutoScalerMin
			}
			if v.AutoScalerMax != 0 {
				m["auto_scaler_max"] = v.AutoScalerMax
			}
			if v.MaxSurge != nil {
				m["max_surge"] = v.MaxSurge.IntValue()
			}
			if v.MaxUnavailable != nil {
				m["max_unavailable"] = v.MaxUnavailable.IntValue()
			}
			if len(v.VolumeType) > 0 {
				m["volume_type"] = v.VolumeType
			}
			if len(v.VolumeSize) > 0 {
				m["volume_size"] = v.VolumeSize
			}
			workers[i] = m
		}
		d.Set("worker", workers)
	}
}
