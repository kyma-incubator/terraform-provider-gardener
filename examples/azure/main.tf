provider "gardener" {
  profile              = "<my-gardener-project>"
  azure_secret_binding = "<my-gardener-gcp-secret>"
  kube_path            = "<my-gardener-service-account-kubeconfig>"
}

resource "gardener_azure_shoot" "my-server" {
  name              = "tf-gardener-azure"
  region            = "westeurope"
  kubernetesversion = "1.15.2"
  vnetcidr          = "10.250.112.0/16"
  workercidr        = "10.250.112.0/22"
  worker {
    name           = "cpu-worker"
    machinetype    = "Standard_D2_v3"
    autoscalermin  = 2
    autoscalermax  = 2
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "35Gi"
    volumetype     = "standard"
  }
}
