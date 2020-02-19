provider "gardener" {
  kube_file = "${file("<my-gardener-service-account-kubeconfig>")}"
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
    cloud_profile_name  = "gcp"
    region              = "europe-west3"
    secret_binding_name = "<secret_binding>"

    networking {
      nodes    = "10.250.0.0/19"
      pods     = "100.96.0.0/11"
      services = "100.64.0.0/13"
      type     = "calico"
    }

    provider {
      type = "gcp"

      infrastructure_config {
        aws {
          enableecraccess = true
          networks {
            vpc {
              cidr = "10.50.0.0/16"
            }
            zones {
              name     = "eu-central-1a"
              internal = "10.50.112.0/22"
              public   = "10.50.96.0/22"
              workers  = "10.50.0.0/19"
            }
          }
        }
      }

      worker {
        max_surge       = 1
        max_unavailable = 0
        maximum         = 2
        minimum         = 2
        volume {
          size = "50Gi"
          type = "Standard_LRS"
        }
        name = "cpu-worker"

        machine {
          image {
            name    = "coreos"
            version = "2303.3.0"
          }
          type = "Standard_A4_v2"
        }
      }
    }

    kubernetes {
      version = "1.15.4"
    }
  }
}
