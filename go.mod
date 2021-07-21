module github.com/kyma-incubator/terraform-provider-gardener

go 1.13

require (
	github.com/gardener/gardener v1.0.2
	github.com/gardener/gardener-extension-provider-aws v1.3.3
	github.com/gardener/gardener-extension-provider-azure v1.3.0
	github.com/gardener/gardener-extension-provider-gcp v1.3.0
	github.com/google/go-cmp v0.5.4
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/terraform v0.12.13
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0 // indirect
	google.golang.org/appengine v1.6.4 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90 // kubernetes-1.16.0
