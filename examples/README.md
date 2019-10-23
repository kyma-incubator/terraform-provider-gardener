# Terraform Provider for Gardener Examples

## Overview

This folder contains a set of examples which use Gardener services to deploy [AWS](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/aws), [GCP](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/gcp) and [Azure](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/azure) clusters. All provided examples have the following section in common:

```bash
provider "gardener" {
  kube_path          = "<my-gardener-service-account-kubeconfig>"
}
```
This section includes the following parameters:
* **kube_path** - the path for the kubeconfig file of the service bot of the profile. 

## Installation

Follow these steps to run an example:

1. Clone the `terraform-provider-gardener` repository.
2. Go to `terraform-provider-gardener/examples/{example_name}`.
3. Run `terraform apply`.

