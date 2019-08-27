provider "gardener" {
  profile            = "<my-gardener-project>"
  aws_secret_binding = "<my-gardener-gcp-secret>"
  kube_path          = "<my-gardener-service-account-kubeconfig>"
}

resource "gardener_aws_shoot" "my-server" {
  name              = "tf-gardener-aws"
  region            = "eu-west-1"
  zones             = ["eu-west-1a"]
  kubernetesversion = "1.15.2"
  worker {
    name           = "cpu-worker"
    machinetype    = "m5.large"
    autoscalermin  = 2
    autoscalermax  = 2
    maxsurge       = 1
    maxunavailable = 0
    volumesize     = "20Gi"
    volumetype     = "gp2"
  }
}
