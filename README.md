
# Terraform Provider for Gardener


## Overview

The Terraform Provider for Gardener enables [Terraform](https://www.terraform.io) to provision [Kubernetes](https://kubernetes.io) clusters on chosen cloud providers using [Gardener](https://gardener.cloud/). Currently, it supports AWS, Azure, and GCP.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.12 or higher
-  Gardener project with kubeconfig access and configured cloud provider secrets

## Development

Perform the following steps to build the providers:

1. Build the provider:
    ```bash
    go build -o terraform-provider-gardener
    ```
2. Move the gardener provider binary into the terraform plugins folder:

   On MacOS:
   ```bash
   mv terraform-provider-gardener ~/.terraform.d/plugins/terraform-provider-gardener
   ```
   For other systems and further details, see [this](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) document.

## Usage

Perform the following steps to use the provider:

1. Go to the provider [example](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples) folder:

    ```bash
    cd examples/<provider>
    ```
2. Edit the `main.tf` file to provide the following parameters:

    * Gardener project name
    * Gardener secret for the choosen cloud provider(s)
    * Path to the Gardener kubeconfig

     ```bash
     provider "gardener" {
        profile            = "<my-gardener-project>"
        gcp_secret_binding = "<my-gardener-gcp-secret>"
        kube_path          = "<my-gardener-service-account-kubeconfig>"
     }
     ```
      >**NOTE:** To obtain the gardener secret and kubeconfig go to the [Gardener dashboard](https://dashboard.garden.canary.k8s.ondemand.com/login).
      
3. Initialize Terraform:
    ```bash
    terraform init
    ```
4. Plan the provisioning:
    ```bash
    terraform plan
    ```
5. Deploy the cluster:
    ```bash
    terraform apply
    ```
## Examples

See the [examples](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples) to learn how to use Gardener to deploy clusters on AWS, GCP, and Azure.
