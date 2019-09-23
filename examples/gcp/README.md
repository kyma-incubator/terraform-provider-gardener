# Gardener GCP

Use this configuration example to deploy a GCP cluster using Gardener.
 >**NOTE**: For the example to work properly, make sure that the **secret_binding** parameter is set to `gcp_secret_binding` and that the **resource** is `gardener_gcp_shoot`.

For now the following cidrs are obtrained from the similar clusters deployed on gardener, in the future these values should be obtained programatically.
```bash
   workerscidr       = ["10.250.0.0/19"]
```

You can also define multiple workers:
```bash
 worker {
    name           = "cpu-worker1"
    machinetype    = "n1-standard-4"
    autoscalermin  = 2
    autoscalermax  = 2
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "20Gi"
    volumetype     = "pd-standard"
  }
  worker {
    name           = "cpu-worker2"
    machinetype    = "n1-standard-4"
    autoscalermin  = 2
    autoscalermax  = 2
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "30Gi"
    volumetype     = "pd-standard"
  }
```


