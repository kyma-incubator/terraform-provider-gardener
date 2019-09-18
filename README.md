
Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.12 or higher

## Overview
Creating a new Terraform provider for Gardener using the Gardener client.This PoC contains a terraform provider that is meant to be executed by the terraform CLI.

Developing the Provider
---------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](https://github.com/kyma-incubator/terraform-provider-gardener#requirements) before proceeding).


Building the provider
---------------------------

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

Using the provider
---------------------------
After building the provider as stated previously, you should create your terraform config file and run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider.For installation methods, documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/aws/index.html). After the provider is installed you can apply your configurations usin the command [`terraform apply`](https://www.terraform.io/docs/commands/apply.html).

Examples of different configurations for different cluster providers could be found in the [examples](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples) folder.
1. Go to the provider example folder
    ```bash
    cd examples/<provider>
    ```
2. Edit the `main.tf` file providing the required gardener configuration:
    - Gardener project name
    - Gardener secret for the choosen cloud provider(s).
    - Path to the gardener kubeconfig.
    > NOTE: To obtain the gardener secret and kubeconfig visit the gardener dashboard.
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
