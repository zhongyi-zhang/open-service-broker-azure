package mssqlfg

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
)

func deleteARMDeployment(
	armDeployer *arm.Deployer,
	armDeploymentName string,
	resourceGroup string,
) error {
	return (*armDeployer).Delete(
		armDeploymentName,
		resourceGroup,
	)
}
