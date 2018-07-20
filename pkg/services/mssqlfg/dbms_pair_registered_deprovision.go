package mssqlfg

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsPairRegisteredManager) GetDeprovisioner(
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
	)
}

func (d *dbmsPairRegisteredManager) deletePriARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsPairInstanceDetails)
	if err := deleteARMDeployment(
		&d.armDeployer,
		dt.PriARMDeploymentName,
		instance.ProvisioningParameters.GetString("primaryResourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}

func (d *dbmsPairRegisteredManager) deleteSecARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsPairInstanceDetails)
	if err := deleteARMDeployment(
		&d.armDeployer,
		dt.SecARMDeploymentName,
		instance.ProvisioningParameters.GetString("secondaryResourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return dt, nil
}
