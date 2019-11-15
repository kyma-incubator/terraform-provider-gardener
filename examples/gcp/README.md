# Deploy a cluster on GCP

## Overview
Use the configuration example located in [main.tf](main.tf) to deploy a GCP cluster using Gardener.

Currently, values for the following CIDRs are obtained from the similar clusters deployed on Gardener. See the example:

```bash
   workers       = ["10.250.0.0/19"]
```

You can define multiple workers using the following configuration:

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
