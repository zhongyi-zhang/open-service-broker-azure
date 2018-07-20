package mssqlfg

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databasePairFeManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	td := instance.Plan.GetProperties().Extended["tierDetails"]
	details := td.(planDetails)
	return details.validateUpdateParameters(instance)
}

func (d *databasePairFeManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	// There isn't a need to do any "pre-provision here. just the update step"
	return service.NewUpdater(
		service.NewUpdatingStep("updatePriARMTemplate", d.updatePriARMTemplate),
		service.NewUpdatingStep("updateSecARMTemplate", d.updateSecARMTemplate),
	)
}

func (d *databasePairFeManager) updatePriARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	err := updateDatabaseARMTemplate(
		&d.armDeployer,
		dt.PriARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("primaryResourceGroup"),
		instance.Parent.ProvisioningParameters.GetString("primaryLocation"),
		pdt.PriServerName,
		dt.DatabaseName,
		*instance.UpdatingParameters,
		pd,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairFeManager) updateSecARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	err := updateDatabaseARMTemplate(
		&d.armDeployer,
		dt.SecARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("secondaryResourceGroup"),
		instance.Parent.ProvisioningParameters.GetString("secondaryLocation"),
		pdt.SecServerName,
		dt.DatabaseName,
		*instance.UpdatingParameters,
		pd,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}
