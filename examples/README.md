# Terraform Provider for Gardener Examples

## Overview

This folder contains a set of examples which use Gardener services to deploy [AWS](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/aws), [GCP](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/gcp) and [Azure](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples/azure) clusters. All provided examples have the following section in common:

```bash
provider "gardener" {
  kube_file          = "${file("<my-gardener-service-account-kubeconfig>")}"
}
```
You can pass the kube_file using the raw text alone as follows:
```bash
kube_file          =<<-EOT
    kind: Config
    clusters:
      - cluster:
          certificate-authority-data: >-
            <certificate-authority-data>
          server: "https://gardener.garden.canary.k8s.ondemand.com"
        name: garden
    users:
      - user:
          token: >-
            <token>
        name: robot
    contexts:
      - context:
          cluster: garden
          user: robot
          namespace: garden-<profile>
        name: garden-<profile>-robot
    current-context: garden-<profile>-robot

    EOT
```
This section includes the following parameters:
* **kube_file** - the raw string of the kube config file.

## Installation

Follow these steps to run an example:

1. Clone the `terraform-provider-gardener` repository.
2. Go to `terraform-provider-gardener/examples/{example_name}`.
3. Run `terraform apply`.

