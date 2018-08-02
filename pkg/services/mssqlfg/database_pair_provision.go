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
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep(
			"checkNameAvailability",
			d.checkNameAvailability,
		),
		service.NewProvisioningStep("deployPriARMTemplate", d.deployPriARMTemplate),
		service.NewProvisioningStep("deployFgARMTemplate", d.deployFgARMTemplate),
		service.NewProvisioningStep("deploySecARMTemplate", d.deploySecARMTemplate),
	)
}

func (d *databasePairManager) preProvision(
	_ context.Context,
	_ service.Instance,
) (service.InstanceDetails, error) {
	return &databasePairInstanceDetails{
		PriARMDeploymentName: uuid.NewV4().String(),
		SecARMDeploymentName: uuid.NewV4().String(),
		FgARMDeploymentName:  uuid.NewV4().String(),
	}, nil
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
	); err != nil {
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Primary database with the name is already existed")
	}
	if err := getDatabase(
		ctx,
		&d.databasesClient,
		ppp.GetString("secondaryResourceGroup"),
		pdt.SecServerName,
		pp.GetString("database"),
	); err != nil {
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Secondary database with the name is already " +
			"existed")
	}
	if err := getFailoverGroup(
		ctx,
		&d.failoverGroupsClient,
		ppp.GetString("primaryResourceGroup"),
		pdt.SecServerName,
		pp.GetString("failoverGroup"),
	); err != nil {
		if !strings.Contains(err.Error(), "ResourceNotFound") {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Failover group with the name is already existed")
	}
	dt := instance.Details.(*databasePairInstanceDetails)
	dt.FailoverGroupName = pp.GetString("failoverGroup")
	dt.DatabaseName = pp.GetString("database")
	return dt, nil
}

func (d *databasePairManager) deployPriARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databasePairInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pd := instance.Plan.GetProperties().Extended["tierDetails"].(planDetails)
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	if err := deployDatabaseARMTemplate(
		&d.armDeployer,
		dt.PriARMDeploymentName,
		ppp.GetString("primaryResourceGroup"),
		ppp.GetString("primaryLocation"),
		pdt.PriServerName,
		dt.DatabaseName,
		*pp,
		pd,
		tags,
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (d *databasePairManager) deployFgARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := deployFailoverGroupARMTemplate(
		&d.armDeployer,
		instance,
	); err != nil {
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
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	if err := deployDatabaseFeARMTemplate(
		&d.armDeployer,
		dt.SecARMDeploymentName,
		ppp.GetString("secondaryResourceGroup"),
		ppp.GetString("secondaryLocation"),
		pdt.SecServerName,
		pp.GetString("database"),
		tags,
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}
