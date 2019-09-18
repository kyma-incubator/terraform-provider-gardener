# Gardener Provider Examples

This directory contains a set of examples of using gardener services to deploy clusters of 3 different providers [aws](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/aws), [gcp](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/gcp) and [azure](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/azure).The examples have their own README, but all of them have a common provider section of the configuration which is as follows:

```bash
provider "gardener" {
  profile            = "<my-gardener-project>"
  <provider>_secret_binding = "<my-gardener-<provider>-secret>"
  kube_path          = "<my-gardener-service-account-kubeconfig>"
}
```
The profile is the profile you want to deploy to in gardener. The <provider>_secret_binding is the provider secret binding defined for the profile that you want to use (there could be more than one secret per provider for a profile).kube_path is the path for the kubeConfig file of the service bot of the profile. 
To run any example, clone the repository and run `terraform apply` within the example's own directory.
