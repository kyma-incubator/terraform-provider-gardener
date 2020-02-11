module github.com/kyma-incubator/terraform-provider-gardener

go 1.13

require (
	github.com/gardener/gardener v1.0.2
	github.com/gardener/gardener-extension-provider-aws v1.3.3 // indirect
	github.com/gardener/gardener-extension-provider-azure v1.3.0
	github.com/gardener/gardener-extension-provider-gcp v1.3.0
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/golang/gddo v0.0.0-20190904175337-72a348e765d2 // indirect
	github.com/google/go-cmp v0.3.1
	github.com/gorilla/mux v1.7.0 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/terraform v0.12.13
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/karrick/godirwalk v1.10.12 // indirect
	github.com/mattn/go-ieproxy v0.0.0-20190805055040-f9202b1cfdeb // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0 // indirect
	google.golang.org/appengine v1.6.4 // indirect
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c // indirect
	google.golang.org/grpc v1.24.0 // indirect
	gopkg.in/ini.v1 v1.52.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

replace (
	github.com/census-instrumentation/opencensus-proto v0.1.0-0.20181214143942-ba49f56771b8 => github.com/census-instrumentation/opencensus-proto v0.0.3-0.20181214143942-ba49f56771b8
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655 // kubernetes-1.16.0
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90 // kubernetes-1.16.0

)
