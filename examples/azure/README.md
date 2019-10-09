#  Deploy a cluster on Azure

## Overview
Use this configuration example to deploy a cluster on Azure using Gardener.
 >**NOTE**: For the example to work properly, make sure that the **secret_binding** parameter is set to `azure_secret_binding` and that the **resource** is `gardener_azure_shoot`.

Currently, values for the following CIDRs are obtained from the similar clusters deployed on Gardener. See the example:
```bash
   vnetcidr          = "10.250.0.0/16"
  workercidr        = "10.250.0.0/22"
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


