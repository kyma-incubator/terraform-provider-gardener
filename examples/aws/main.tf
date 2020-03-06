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
    cloud_profile_name  = "aws"
    region              = "eu-central-1"
    secret_binding_name = "<secret_binding>"

    networking {
      type     = "calico"
      nodes    = "10.250.0.0/19"
      pods     = "100.96.0.0/11"
      services = "100.64.0.0/13"
    }

    provider {
      type = "aws"
      infrastructure_config {
        aws {
          enableecraccess = true
          networks {
            vpc {
              cidr = "10.50.0.0/16"
            }
            zones {
              name     = "eu-central-1"
              internal = "10.50.112.0/22"
              public   = "10.50.96.0/22"
              workers  = "10.50.0.0/19"
            }
          }
        }
      }
      worker {
        name            = "cpu-worker"
        max_surge       = 1
        max_unavailable = 0
        maximum         = 2
        minimum         = 2
        volume {
          type = "gp2"
          size = "30Gi"
        }
        machine {
          type = "t3.medium"
          image {
            name    = "coreos"
            version = "ami-d0dcef3"
          }
        }
        
        zones = ["eu-central-1a", "eu-central-1b", "eu-central-1c"]
      }
    }

    kubernetes {
      version = "1.17.3"
    }
  }
}
