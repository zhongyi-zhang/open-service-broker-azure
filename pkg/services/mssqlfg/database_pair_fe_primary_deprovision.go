package mssqlfg

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databasePairFePrimaryManager) GetDeprovisioner(
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
		service.NewDeprovisioningStep(
			"deletePriDatabase",
			d.deletePriDatabase,
		),
		service.NewDeprovisioningStep(
			"deleteSecDatabase",
			d.deleteSecDatabase,
		),
		service.NewDeprovisioningStep(
			"deleteFailoverGroup",
			d.deleteFailoverGroup,
		),
	)
}

func (d *databasePairFePrimaryManager) deletePriARMDeployment(
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

func (d *databasePairFePrimaryManager) deleteSecARMDeployment(
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

func (d *databasePairFePrimaryManager) deleteFgARMDeployment(
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

func (d *databasePairFePrimaryManager) deletePriDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if _, err := d.databasesClient.Delete(
		ctx,
		instance.Parent.ProvisioningParameters.GetString("primaryResourceGroup"),
		pdt.PriServerName,
		dt.DatabaseName,
	); err != nil {
		return nil, fmt.Errorf("error deleting sql database: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairFePrimaryManager) deleteSecDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if _, err := d.databasesClient.Delete(
		ctx,
		instance.Parent.ProvisioningParameters.GetString("secResourceGroup"),
		pdt.SecServerName,
		dt.DatabaseName,
	); err != nil {
		return nil, fmt.Errorf("error deleting sql database: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairFePrimaryManager) deleteFailoverGroup(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if _, err := d.databasesClient.Delete(
		ctx,
		instance.Parent.ProvisioningParameters.GetString("primaryResourceGroup"),
		pdt.PriServerName,
		dt.FailoverGroupName,
	); err != nil {
		return nil, fmt.Errorf("error deleting failover group: %s", err)
	}
	return instance.Details, nil
}
