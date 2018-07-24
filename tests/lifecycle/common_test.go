// +build !unit

package lifecycle

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
)

func getAzureConfigAndAuthorizer() (
	*azure.Config,
	*autorest.BearerAuthorizer,
	error,
) {
	azureConfig, err := azure.GetConfigFromEnvironment()
	if err != nil {
		return nil, nil, err
	}
	authorizer, err := azure.GetBearerTokenAuthorizer(
		azureConfig.Environment,
		azureConfig.TenantID,
		azureConfig.ClientID,
		azureConfig.ClientSecret,
	)
	if err != nil {
		return nil, nil, err
	}

	return &azureConfig, authorizer, nil
}
