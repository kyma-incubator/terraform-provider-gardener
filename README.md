
# Terraform Provider for Gardener


## Overview

The Terraform Provider for Gardener enables [Terraform](https://www.terraform.io) to provision [Kubernetes](https://kubernetes.io) clusters on chosen cloud providers using [Gardener](https://gardener.cloud/). Currently, it supports AWS, Azure, and GCP.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.12 or higher
- Gardener project with kubeconfig access and configured cloud provider secrets

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

At the moment the resource type [`gardener_shoot`](shoot/schema_shoot.go) is supported corresponding to the [`shoot`](https://github.com/gardener/gardener/blob/master/pkg/apis/core/v1alpha1/types_shoot.go) resource of gardener. The cloud provider specific configuration and a related example is listed in the table below:

|Cloud Provider|Specification|Example|
|---|---|---|
|Amazon Web Services|[schema.go](shoot/schema_aws.go)|[main.tf](examples/aws/main.tf)|
|Goocle Cloud Platform|[schema.go](shoot/schema_gcp.go)|[main.tf](examples/gcp/main.tf)|
|Azure|[schema.go](shoot/schema_azure.go)|[main.tf](examples/azure/main.tf)|

-Perform the following steps to use the provider:
 
1. Navigate to one of the main.tf files as listed above

2. Edit the `main.tf` file to provide the missing parameters

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
