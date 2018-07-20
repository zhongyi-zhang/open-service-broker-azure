package mssqlfg

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databasePairRegisteredManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep(
			"deletePriARMDeployment",
			d.deletePriARMDeployment,
		),
		service.NewDeprovisioningStep(
			"deleteSecARMDeployment",
			d.deleteSecARMDeployment,
		),
		service.NewDeprovisioningStep(
			"deleteFgARMDeployment",
			d.deleteFgARMDeployment,
		),
	)
}

func (d *databasePairRegisteredManager) deletePriARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	err := d.armDeployer.Delete(
		dt.PriARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("primaryResourceGroup"),
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairRegisteredManager) deleteSecARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	err := d.armDeployer.Delete(
		dt.SecARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("secResourceGroup"),
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairRegisteredManager) deleteFgARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	err := d.armDeployer.Delete(
		dt.FgARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("primaryResourceGroup"),
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}
