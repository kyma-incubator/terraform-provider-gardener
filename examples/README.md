# Terraform Provider for Gardener Examples

## Overview

This folder contains a set of examples which use Gardener services to deploy [AWS](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/aws), [GCP](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/gcp) and [Azure](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/azure) clusters. All provided examples have the following section in common:

```bash
provider "gardener" {
  profile            = "<my-gardener-project>"
  <provider>_secret_binding = "<my-gardener-<provider>-secret>"
  kube_path          = "<my-gardener-service-account-kubeconfig>"
}
```
Here you can configure provider-specific parameters:

* **profile** - the name of your Gardener project, you want to deploy the cluster to.
* **{provider}_secret_binding** - the name of the secret containing the service account key. There may be one secret defined for a provider.
* **kube_path** - the path for the kubeconfig file of the service account. 

## Installation

Follow these steps to run an example:

1. Clone the `terraform-provider-gardener` repository.
2. Go to `terraform-provider-gardener/examples/{example_name}`.
3. Run  `terraform apply` to deploy the example.
