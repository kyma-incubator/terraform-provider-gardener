# Gardener Provider Examples

## Overview
This folder contains a set of examples which use Gardener services to deploy [AWS](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/aws), [GCP](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/gcp) and [Azure](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/azure) clusters. Although the configurations for these providers differ, they have the following common section:

```bash
provider "gardener" {
  profile            = "<my-gardener-project>"
  <provider>_secret_binding = "<my-gardener-<provider>-secret>"
  kube_path          = "<my-gardener-service-account-kubeconfig>"
}
```
This section includes the following parameters:
* **profile** - the profile you want to deploy to in gardener. 
* **<provider>_secret_binding** - the provider secret binding defined for the profile that you want to use. There might be more than one secret per provider for a profile.
* **kube_path** - the path for the kubeconfig file of the service bot of the profile. 

## Installation
Follow these steps to run an example:
1. Clone the `terraform-provider-gardener` repository.
2. Go to `terraform-provider-gardener/examples/{example_name}`.
3. Run  `terraform apply`. 
