package types

type ProviderType string
type ProfileType string
type SeedType string

const (
	// GCPProvider stands for the Google Cloud Platform provider.
	GCPProvider ProviderType = "gcp"
	// AzureProvider stands for the Microsoft Azure Cloud Computing Platform provider.
	AzureProvider ProviderType = "azure"
	// AWSProvider stands for Amazon Web Services provider.
	AWSProvider ProviderType = "aws"
	// GCPProfile stands for the Google Cloud Platform profile.
	GCPProfile ProfileType = "gcp"
	// AzureProfile stands for the Microsoft Azure Cloud Computing Platform profile.
	AzureProfile ProfileType = "az"
	// AWSProfile stands for Amazon Web Services profile.
	AWSProfile ProfileType = "aws"
	// GCPSeed stands for the Google Cloud Platform seed
	GCPSeed SeedType = "gcp-eu1"
	// AzureSeed stands for the Microsoft Azure Cloud Computing Platform seed.
	AzureSeed SeedType = "az-eu1"
	// AWSProfile stands for Amazon Web Services seed.
	AWSSeed SeedType = "aws-eu1"
)
