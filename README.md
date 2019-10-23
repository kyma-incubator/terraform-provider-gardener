
# Terraform Provider for Gardener


## Overview

The Terraform Provider for Gardener enables [Terraform](https://www.terraform.io) to provision [Kubernetes](https://kubernetes.io) clusters on chosen cloud providers using [Gardener](https://gardener.cloud/). Currently, it supports AWS, Azure, and GCP.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
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

Perform the following steps to use the provider:

1. Go to the provider [example](https://github.com/kyma-incubator/terraform-provider-gardener/tree/master/examples) folder:

    ```bash
    cd examples/<provider>
    ```
2. Edit the `main.tf` file to provide the following parameters:

    - Path to the Gardener kubeconfig
    - Shoot specification

    > **NOTE:** To obtain the gardener secret and kubeconfig go to the [Gardener dashboard](https://dashboard.garden.canary.k8s.ondemand.com/login).
    ```bash
    provider "gardener" {
        kube_path          = "<my-gardener-service-account-kubeconfig>"
    }
    resource "gardener_shoot" "<Name>" {
        metadata {
            name      = "<name-to-be-shown-in-gardener>"
            namespace = "<gardener-profile-namespace>"

        }
        spec {
            cloud {
            profile = "az"
            region  = "westeurope"
            seed    = "az-eu1"

            secret_binding_ref {
                name = "<secret_binding>"
            }

            azure {
                    networks {
                        vnet {
                            cidr = "10.250.0.0/16"
                        }
                        workers = "10.250.0.0/19"
                    }

                    worker {
                        name            = "cpu-worker"
                        machine_type    = "Standard_D2_v3"
                        auto_scaler_min = 3
                        auto_scaler_max = 3
                        max_surge       = 1
                        max_unavailable = 0
                        volume_type     = "standard"
                        volume_size     = "50Gi"
                    }
                }
            }

            kubernetes {
                version = "1.15.4"
            }
        }
    }
    ```
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
