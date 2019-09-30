
# Terraform Provider


## Overview
Creating a new Terraform provider for Gardener using the Gardener client.This PoC contains a terraform provider that is meant to be executed by the terraform CLI. The provider supports provisioning of three infrastructures (aws,azure,gcp) for now and in the future we will add the other infrastructures supported by Gardener.
### Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.12 or higher

## Development

1. Resolve dependencies:
    ```bash
    dep ensure
    ```
2. Build the provider:
    ```bash
    go build -o terraform-provider-gardener
    ```
3. Move the gardener provider binary into the terraform plugins folder*.
    >\* See https://www.terraform.io/docs/plugins/basics.html#installing-plugins for more information.

## Usage

After building the provider as stated previously, you should create your terraform config file and run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider.For installation methods. After the provider is installed you can apply your configurations usin the command [`terraform apply`](https://www.terraform.io/docs/commands/apply.html).

1. Go to the provider [example](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples) folder
    ```bash
    cd examples/<provider>
    ```
2. Edit the `main.tf` file providing the required gardener configuration:
    - Gardener project name
    - Gardener secret for the choosen cloud provider(s).
    - Path to the gardener kubeconfig.
    > **NOTE:** To obtain the gardener secret and kubeconfig visit the [gardener dashboard](https://dashboard.garden.canary.k8s.ondemand.com/login).
    ```bash
    provider "gardener" {
        profile            = "<my-gardener-project>"
        gcp_secret_binding = "<my-gardener-gcp-secret>"
        kube_path          = "<my-gardener-service-account-kubeconfig>"
    }
    ```
3. Initialise terraform:
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
Example config files can be found in the [examples](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples) folder for the three providers that we currently support.
