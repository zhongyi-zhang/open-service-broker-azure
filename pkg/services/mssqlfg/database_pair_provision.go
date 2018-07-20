package mssqlfg

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databasePairManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep(
			"checkNameAvailability",
			d.checkNameAvailability,
		),
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("deployPriARMTemplate", d.deployPriARMTemplate),
		service.NewProvisioningStep("deploySecARMTemplate", d.deploySecARMTemplate),
		service.NewProvisioningStep("deployFgARMTemplate", d.deployFgARMTemplate),
	)
}

func (d *databasePairManager) checkNameAvailability(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if err := getDatabase(
		ctx,
		&d.databasesClient,
		ppp.GetString("primaryResourceGroup"),
		pdt.PriServerName,
		pp.GetString("database"),
	); !strings.HasPrefix(err.Error(), "can't find") {
		return nil, fmt.Errorf("Primary database with the name is already existed")
	}
	if err := getDatabase(
		ctx,
		&d.databasesClient,
		ppp.GetString("secondaryResourceGroup"),
		pdt.SecServerName,
		pp.GetString("database"),
	); !strings.HasPrefix(err.Error(), "can't find") {
		return nil, fmt.Errorf("Secondary database with the name " +
			"is already existed",
		)
	}
	if err := getFailoverGroup(
		ctx,
		&d.failoverGroupsClient,
		ppp.GetString("primaryResourceGroup"),
		pdt.SecServerName,
		pp.GetString("failoverGroup"),
	); !strings.HasPrefix(err.Error(), "can't find") {
		return nil, fmt.Errorf("Failover group with the name is already existed")
	}
	return instance.Details, nil
}

func (d *databasePairManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	return &databasePairInstanceDetails{
		PriARMDeploymentName: uuid.NewV4().String(),
		SecARMDeploymentName: uuid.NewV4().String(),
		FgARMDeploymentName:  uuid.NewV4().String(),
		FailoverGroupName:    pp.GetString("failoverGroup"),
		DatabaseName:         pp.GetString("database"),
	}, nil
}

func (d *databasePairManager) deployPriARMTemplate(
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
	err := deployDatabaseARMTemplate(
		&d.armDeployer,
		dt.PriARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("primaryResourceGroup"),
		instance.Parent.ProvisioningParameters.GetString("primaryLocation"),
		pdt.PriServerName,
		dt.DatabaseName,
		*instance.ProvisioningParameters,
		pd,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairManager) deploySecARMTemplate(
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
	err := deployDatabaseARMTemplate(
		&d.armDeployer,
		dt.SecARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("secondaryResourceGroup"),
		instance.Parent.ProvisioningParameters.GetString("secondaryLocation"),
		pdt.SecServerName,
		dt.DatabaseName,
		*instance.ProvisioningParameters,
		pd,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairManager) deployFgARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	err := deployFailoverGroupARMTemplate(
		&d.armDeployer,
		instance,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}
