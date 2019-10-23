provider "gardener" {
  kube_path = "<my-gardener-service-account-kubeconfig>"
}

resource "gardener_shoot" "test_cluster" {
  metadata {
    name      = "test-cluster"
    namespace = "garden-<profile>"

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
