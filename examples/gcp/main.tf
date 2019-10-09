provider "gardener" {
  kube_path = "<kube_path>"
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

      secret_binding_ref {
        name = "<secret_binding>"
      }

      gcp {
        networks {
          workers = ["10.250.0.0/19"]
        }

        workers {
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
