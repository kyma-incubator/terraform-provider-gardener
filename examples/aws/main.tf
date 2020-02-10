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
      worker {
        name            = "cpu-worker"
        max_surge       = 1
        max_unavailable = 0
        maximum         = 2
        minimum         = 2
        volume {
          size = "50Gi"
          type = "Standard_LRS"
        }
        machine {
          type    = "m5.large"
        }
      }

    }

    kubernetes {
      version = "1.15.4"
    }
  }
}

