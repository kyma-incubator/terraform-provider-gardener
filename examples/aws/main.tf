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
      profile = "aws"
      region  = "eu-central-1"
      seed    = "aws-eu1"

      secret_binding_ref {
        name = "<secret_binding>"
      }

      aws {
        networks {
          vpc {
            cidr = "10.250.0.0/16"
          }

          internal = ["10.250.112.0/22"]
          public   = ["10.250.96.0/22"]
          workers  = ["10.250.0.0/19"]
        }

        workers {
          name            = "cpu-worker"
          machine_type    = "m5.large"
          auto_scaler_min = 3
          auto_scaler_max = 4
          max_surge       = 1
          max_unavailable = 0
          volume_type     = "gp2"
          volume_size     = "50Gi"
        }

        zones = ["eu-central-1a"]
      }
    }

    kubernetes {
      version = "1.15.4"
    }
  }
}
