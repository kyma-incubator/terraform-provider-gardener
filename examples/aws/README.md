#  Deploy a cluster on AWS

## Overview
Use the configuration example located in [main.tf](main.tf) to deploy an AWS cluster using Gardener.

Currently, values for the following CIDRs are obtained from the similar clusters deployed on Gardener. The CIDR for each zone is calculated from the VPC CIDR.
See the example:
```bash
  networks {
    vpc {
      cidr = "10.250.0.0/16"
    }
    zones {
      name     = "eu-central-1a"
      workers  = "10.250.0.0/19"
      public   = "10.250.32.0/20"
      internal = "10.250.48.0/20"
    }
  }
```

You can define multiple workers using the following configuration:

```bash
  worker {
    name            = "cpu-worker1"
    zones           = ["eu-central-1a"]
    max_surge       = "3"
    max_unavailable = "1"
    maximum         = "4"
    minimum         = "2"
    volume {
      size = "30Gi"
      type = "gp2"
    }
    machine {
      image {
        name    = "gardenlinux"
        version = "27.1.0"
      }
      type = "m5.xlarge"
    }
  }
  worker {
    name            = "cpu-worker2"
    zones           = ["eu-central-1a"]
    max_surge       = "3"
    max_unavailable = "1"
    maximum         = "4"
    minimum         = "2"
    volume {
      size = "30Gi"
      type = "gp2"
    }
    machine {
      image {
        name    = "gardenlinux"
        version = "27.1.0"
      }
      type = "m5.xlarge"
    }
  }
```
