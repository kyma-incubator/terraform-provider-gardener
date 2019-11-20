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

```bash
worker {
    name           = "cpu-worker1"
    machinetype    = "Standard_D2_v3"
    autoscalermin  = 2
    autoscalermax  = 4
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "35Gi"
    volumetype     = "standard"
  }
 worker {
    name           = "cpu-worker2"
    machinetype    = "Standard_D2_v3"
    autoscalermin  = 2
    autoscalermax  = 4
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "20Gi"
    volumetype     = "standard"
  }
```
