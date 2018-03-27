package azure

import (
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/kelseyhightower/envconfig"
)

const (
	envconfigPrefix = "AZURE"
	azureStackCloud = "AzureStackCloud"
)

// Config represents details necessary for the broker to interact with
// an Azure subscription
type Config struct {
	Environment    azure.Environment
	SubscriptionID string `envconfig:"SUBSCRIPTION_ID" required:"true"`
	TenantID       string `envconfig:"TENANT_ID" required:"true"`
	ClientID       string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret   string `envconfig:"CLIENT_SECRET" required:"true"`
}

type tempConfig struct {
	Config
	EnvironmentStr		string `envconfig:"ENVIRONMENT" default:"AzurePublicCloud"` // nolint: lll
	ResourceManagerEndpoint string `envconfig:"RESOURCE_MANAGER_ENDPOINT"` // nolint: lll
}

// NewConfigWithDefaults returns a Config object with default values already
// applied. Callers are then free to set custom values for the remaining fields
// and/or override default values.
func NewConfigWithDefaults() Config {
	return Config{}
}

// GetConfigFromEnvironment returns Azure-related configuration derived from
// environment variables
func GetConfigFromEnvironment() (Config, error) {
	c := tempConfig{
		Config: NewConfigWithDefaults(),
	}
	err := envconfig.Process(envconfigPrefix, &c)
	if err != nil {
		return c.Config, err
	}
	if c.EnvironmentStr != azureStackCloud {
		c.Environment, err = azure.EnvironmentFromName(c.EnvironmentStr)
        } else {
		properties := azure.OverrideProperty{
			Key: azure.EnvironmentName,
			Value: azureStackCloud,
		}
		c.Environment, err = azure.EnvironmentFromURL(c.ResourceManagerEndpoint, properties)
        }
	return c.Config, err
}

func IsAzureStackCloud() bool {
	c := tempConfig{
		Config: NewConfigWithDefaults(),
	}
	err := envconfig.Process(envconfigPrefix, &c)
	if err != nil {
		return false
	}
	return c.EnvironmentStr == azureStackCloud
}
