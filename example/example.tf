provider "example" {
}

resource "example_server" "my-server" {
  name           = "demo2-gcp"
  profile        = "omar"
  region         = "europe-west3"
  secret_binding = "icke-architecture"
  zones          = ["europe-west3-b"]
  workers {
    name           = "cpu-worker"
    machinetype    = "n1-standard-4"
    autoscalermin  = 2
    autoscalermax  = 2
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "20Gi"
    volumetype     = "pd-standard"
  }
}
