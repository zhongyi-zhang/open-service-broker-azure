package mssqlfg

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databasePairManager) GetDeprovisioner(
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

func (d *databasePairManager) deletePriARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	if err := d.armDeployer.Delete(
		dt.PriARMDeploymentName,
		ppp.GetString("primaryResourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairManager) deleteSecARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	if err := d.armDeployer.Delete(
		dt.SecARMDeploymentName,
		ppp.GetString("secondaryResourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairManager) deleteFgARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	if err := d.armDeployer.Delete(
		dt.FgARMDeploymentName,
		ppp.GetString("primaryResourceGroup"),
	); err != nil {
		return nil, fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairManager) deletePriDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	if dt.DatabaseName != "" {
		if _, err := d.databasesClient.Delete(
			ctx,
			ppp.GetString("primaryResourceGroup"),
			pdt.PriServerName,
			dt.DatabaseName,
		); err != nil {
			return nil, fmt.Errorf("error deleting sql database: %s", err)
		}
	}
	return instance.Details, nil
}

func (d *databasePairManager) deleteSecDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	if dt.DatabaseName != "" {
		if _, err := d.databasesClient.Delete(
			ctx,
			ppp.GetString("secondaryResourceGroup"),
			pdt.SecServerName,
			dt.DatabaseName,
		); err != nil {
			return nil, fmt.Errorf("error deleting sql database: %s", err)
		}
	}
	return instance.Details, nil
}

func (d *databasePairManager) deleteFailoverGroup(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	ppp := instance.Parent.ProvisioningParameters
	if dt.FailoverGroupName != "" {
		if _, err := d.failoverGroupsClient.Delete(
			ctx,
			ppp.GetString("primaryResourceGroup"),
			pdt.PriServerName,
			dt.FailoverGroupName,
		); err != nil {
			return nil, fmt.Errorf("error deleting failover group: %s", err)
		}
	}
	return instance.Details, nil
}
