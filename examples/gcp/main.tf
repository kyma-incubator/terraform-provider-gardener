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
    cloud {
      profile = "gcp"
      region  = "europe-west3"
      seed    = "gcp-eu1"
      secret_binding_ref {
        name = "<secret_binding>"
      }

      gcp {
        networks {
          workers = ["10.250.0.0/19"]
        }

        worker {
          name            = "cpu-worker"
          machine_type    = "n1-standard-4"
          auto_scaler_min = 3
          auto_scaler_max = 3
          max_surge       = 1
          max_unavailable = 0
          volume_type     = "pd-standard"
          volume_size     = "50Gi"
        }

        zones = ["europe-west3-b"]
      }
    }

    kubernetes {
      version = "1.15.4"
    }
  }
}
