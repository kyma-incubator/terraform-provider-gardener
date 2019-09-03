provider "gardener" {
  profile            = "<my-gardener-project>"
  gcp_secret_binding = "<my-gardener-gcp-secret>"
  kube_path          = "<my-gardener-service-account-kubeconfig>"
}

resource "gardener_gcp_shoot" "my-server" {
  name              = "tf-gardener-gcp"
  region            = "europe-west3"
  zones             = ["europe-west3-b"]
  workerscidr       = ["10.250.0.0/19"]
  kubernetesversion = "1.15.2"
  worker {
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
