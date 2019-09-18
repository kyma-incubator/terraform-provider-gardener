# Gardener Aws Example

This is a configuration example to deploy an aws cluster using gardener. You need to make sure that the secret binding parameter is "aws_secret_binding" for the deployment to work properly and that the resource is "gardener_aws_shoot".

For now the following cidrs are obtrained from the similar clusters deployed on gardener, in the future these values should be obtained programatically.
```bash
  workerscidr       = ["10.250.0.0/17"]
  internalscidr     = ["10.250.112.0/22"]
  publicscidr       = ["10.250.96.0/22"]
  vpccidr = "10.250.0.0/16"
```

multiple workers can be defined as follows
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


