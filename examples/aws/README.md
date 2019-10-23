#  Deploy a cluster on AWS

## Overview
Use this configuration example to deploy an AWS cluster using Gardener.

Currently, values for the following CIDRs are obtained from the similar clusters deployed on Gardener. See the example:
```bash
  workerscidr       = ["10.250.0.0/17"]
  internalscidr     = ["10.250.112.0/22"]
  publicscidr       = ["10.250.96.0/22"]
  vpc{
    cidr = "10.250.0.0/16"
  }
```

You can define multiple workers using the following configuration:

```bash
worker {
    name           = "cpu-worker1"
    machinetype    = "m5.large"
    autoscalermin  = 2
    autoscalermax  = 2
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "20Gi"
    volumetype     = "gp2"
  }
  worker {
    name           = "cpu-worker2"
    machinetype    = "m5.large"
    autoscalermin  = 2
    autoscalermax  = 2
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "30Gi"
    volumetype     = "gp2"
  }
```
