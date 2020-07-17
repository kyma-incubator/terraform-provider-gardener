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

resource "gardener_shoot" "gardener_cluster" {
  metadata {
    name      = "test-cluster"
    namespace = "garden-<profile>"
  }
  timeouts {
    create = "30m0s"
    update = "30m0s"
    delete = "20m0s"
  }
  spec {
    cloud_profile_name  = "aws"
    region              = "eu-central-1"
    secret_binding_name = "<secret_binding>"
    networking {
      nodes    = "10.250.0.0/16"
      pods     = ""
      services = ""
      type     = "calico"
    }
    maintenance {
      auto_update {
        kubernetes_version    = "true"
        machine_image_version = "true"
      }
      time_window {
        begin = "030000+0000"
        end   = "040000+0000"
      }
    }
    provider {
      type = "aws"
      infrastructure_config {
        aws {
          networks {
            vpc {
              cidr = "10.250.0.0/16"
            }
            zones {
              name     = "eu-central-1a"
              workers  = "10.250.0.0/19"
              public   = "10.250.32.0/20"
              internal = "10.250.48.0/20"
            }
          }
        }
      }
      worker {
        name            = "cpu-worker"
        zones           = ["eu-central-1a"]
        max_surge       = "3"
        max_unavailable = "1"
        maximum         = "4"
        minimum         = "2"
        volume {
          size = "30Gi"
          type = "gp2"
        }
        machine {
          image {
            name    = "gardenlinux"
            version = "27.1.0"
          }
          type = "m5.xlarge"
        }
      }
    }
    kubernetes {
      allow_privileged_containers = true
      version                     = "1.18"
    }
  }
}
