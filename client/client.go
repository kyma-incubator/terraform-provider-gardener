package client

import (
	gardener_apis "github.com/gardener/gardener/pkg/client/core/clientset/versioned/typed/core/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	GardenerClientSet *gardener_apis.CoreV1beta1Client
}

// Client configures and returns a fully initialized GardenerClient
func New(c *Config) (interface{}, error) {
	kubeBytes := []byte(c.KubeFile)
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeBytes)
	if err != nil {
		return nil, err
	}
	clientset, err := gardener_apis.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client := &Client{
		GardenerClientSet: clientset,
	}
	return client, nil
}

type Config struct {
	KubeFile string
}
