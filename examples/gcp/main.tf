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
    namespace = "garden-berlin"
  }
  spec {
    cloud_profile_name = "gcp"
    secret_binding_name = "<secret_binding>"
    region  = "europe-west3"

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
    networking{
      type = "calico"
      nodes ="10.250.0.0/16"
    }

    provider {
      type = "gcp"
      control_plane_config {
        gcp {
          zone = "europe-west3-a"
        }

      }
      infrastructure_config {
        gcp {
          networks {
            workers = "10.250.0.0/16"
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
          type = "pd-standard"
        }
        name = "cpu-worker"
        machine {
          image  {
            name = "coreos"
            version = "2303.3.0"
          }
          type = "n1-standard-4"
        }
        zones = ["europe-west3-a"]
      }
    }

    kubernetes {
      version = "1.17.2"
    }
  }
}
