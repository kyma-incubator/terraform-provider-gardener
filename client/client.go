package client

import (
	"errors"

	gardener_apis "github.com/gardener/gardener/pkg/client/garden/clientset/versioned/typed/garden/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	NameSpace         string
	DNSBase           string
	GardenerClientSet *gardener_apis.GardenV1beta1Client
	SecretBindings    *Bindings
}

// Client configures and returns a fully initialized GardenerClient
func New(c *Config) (interface{}, error) {
	if c.SecretBindings.AwsSecretBinding == "" && c.SecretBindings.GcpSecretBinding == "" && c.SecretBindings.AzureSecretBinding == "" &&
		c.SecretBindings.OpenStackSecretBinding == "" && c.SecretBindings.AliCloudSecretBinding == "" {
		return nil, errors.New("at least one binding needs to be defined")
	}

	config, err := clientcmd.BuildConfigFromFlags("", c.KubePath)
	if err != nil {
		return nil, err
	}
	clientset, err := gardener_apis.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client := &Client{
		NameSpace:         "garden-" + c.Profile,
		GardenerClientSet: clientset,
		SecretBindings:    c.SecretBindings,
	}
	return client, nil
}

type Config struct {
	Profile        string
	KubePath       string
	SecretBindings *Bindings
}

type Bindings struct {
	AwsSecretBinding       string
	GcpSecretBinding       string
	AzureSecretBinding     string
	OpenStackSecretBinding string
	AliCloudSecretBinding  string
}
