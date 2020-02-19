#  Deploy a cluster on Azure

## Overview
Use the configuration example located in [main.tf](main.tf) to deploy an Azure cluster using Gardener.

Currently, values for the following CIDRs are obtained from the similar clusters deployed on Gardener. See the example:
```bash
   vnet{
     cidr          = "10.250.0.0/16"
   }
  workers        = "10.250.0.0/22"
```

You can define multiple workers using the following configuration:

```hcl-terraform
worker {
    name           = "cpu-worker1"
    minimum  = 2
    maximum  = 4
    max_surge       = 1
    max_unavailable = 0
    volume {
      size = "50Gi"
      type = "Standard_LRS"
    }
    machine {
      image  {
        name = "coreos"
        version = "2303.3.0"
      }
      type = "Standard_A4_v2"
    }
 }
 worker {
    name           = "cpu-worker2"
    minimum  = 2
    maximum  = 4
    max_surge       = 1
    max_unavailable = 0
    volume {
          size = "50Gi"
          type = "Standard_LRS"
        }
    machine {
      image  {
        name = "coreos"
        version = "2303.3.0"
      }
      type = "Standard_A4_v2"
    }
  }
```

The infrastructure config allows to pass extra information to the cloud provider. It can be done in following way:
```hcl-terraform
infrastructure_config {
  azure {
     networks {
       vnet {
         cidr = "10.250.0.0/16"
       }
       workers = "10.250.0.0/19"
       service_endpoints = ["microsoft.test"]
     }
  }
}
```
