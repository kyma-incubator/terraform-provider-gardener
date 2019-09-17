provider "gardener" {
  profile              = "icke"
  azure_secret_binding = "my-azure-secret"
  kube_path            = "C:\\Users\\D074188\\.kube\\config"
}

resource "gardener_azure_shoot" "my-server" {
  name              = "tf-gardener-azure"
  region            = "westeurope"
  kubernetesversion = "1.15.2"
  vnetcidr          = "10.250.0.0/16"
  workercidr        = "10.250.0.0/22"
  worker {
    name           = "cpu-worker"
    machinetype    = "Standard_D2_v3"
    autoscalermin  = 2
    autoscalermax  = 4
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "35Gi"
    volumetype     = "standard"
  }
}
