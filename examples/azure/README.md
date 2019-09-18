# Gardener Azure Example

This is a configuration example to deploy an azure cluster using gardener. You need to make sure that the secret binding parameter is "azure_secret_binding" for the deployment to work properly and that the resource is "gardener_azure_shoot".

For now the following cidrs are obtrained from the similar clusters deployed on gardener, in the future these values should be obtained programatically.
```bash
   vnetcidr          = "10.250.0.0/16"
  workercidr        = "10.250.0.0/22"
```

multiple workers can be defined as follows
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


