provider "gardener" {
  kube_file = "${file("my-gardener-service-account-kubeconfig")}"
  /*kube_file          =<<-EOT
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

    EOT*/
}

resource "gardener_shoot" "test_cluster" {
  metadata {
    name      = "test-cluster"
    namespace = "garden-<profile>"

  }

  spec {
    cloud_profile_name = "az"
    secret_binding_name = "<secret_binding>"
    purpose = "evaluation"
    networking {
      nodes = "10.250.0.0/19"
      pods= "100.96.0.0/11"
      services= "100.64.0.0/13"
      type= "calico"
    }
    maintenance {
      auto_update{
        kubernetes_version = true
        machine_image_version = true
      }

      time_window {
        begin = "030000+0000"
        end = "040000+0000"
      }
    }
    provider {
      type = "azure"
      infrastructure_config {
        azure {
          networks {
            vnet {
              cidr = "10.250.0.0/16"
            }
            workers = "10.250.0.0/19"
           # service_endpoints = ["microsoft.test"]
          }
        }
      }
      worker {
        max_surge = 1
        max_unavailable = 0
        maximum = 2
        minimum = 2
        volume {
          size = "50Gi"
          type = "Standard_LRS"
        }
        name = "cpu-worker"
        machine {
          image  {
            name = "coreos"
            version = "2303.3.0"
          }
          type = "Standard_A4_v2"
        }
      }
    }
    region =  "westeurope"
    kubernetes {
      version = "1.15.4"
    }
  }
}
