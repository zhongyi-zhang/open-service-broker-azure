package mssqlfg

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *dbmsPairRegisteredManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("getPriServer", d.getPriServer),
		service.NewProvisioningStep("getSecServer", d.getSecServer),
		service.NewProvisioningStep("testPriConnection", d.testPriConnection),
		service.NewProvisioningStep("testSecConnection", d.testSecConnection),
		service.NewProvisioningStep("deployPriARMTemplate", d.deployPriARMTemplate),
		service.NewProvisioningStep("deploySecARMTemplate", d.deploySecARMTemplate),
	)
}

func (d *dbmsPairRegisteredManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	return &dbmsPairInstanceDetails{
		PriARMDeploymentName:          uuid.NewV4().String(),
		PriServerName:                 pp.GetString("primaryServer"),
		PriAdministratorLogin:         pp.GetString("primaryAdministratorLogin"),
		PriAdministratorLoginPassword: service.SecureString(pp.GetString("primaryAdministratorLoginPassword")), // nolint: lll
		SecARMDeploymentName:          uuid.NewV4().String(),
		SecServerName:                 pp.GetString("secondaryServer"),
		SecAdministratorLogin:         pp.GetString("secondaryAdministratorLogin"),
		SecAdministratorLoginPassword: service.SecureString(pp.GetString("secondaryAdministratorLoginPassword")), // nolint: lll
	}, nil
}

func (d *dbmsPairRegisteredManager) getPriServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	dt := instance.Details.(*dbmsPairInstanceDetails)
	if err := getServer(
		ctx,
		&d.serversClient,
		pp.GetString("primaryResourceGroup"),
		dt.PriServerName,
		instance.Service.GetProperties().Extended["version"].(string),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (d *dbmsPairRegisteredManager) getSecServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	dt := instance.Details.(*dbmsPairInstanceDetails)
	if err := getServer(
		ctx,
		&d.serversClient,
		pp.GetString("secondaryResourceGroup"),
		dt.SecServerName,
		instance.Service.GetProperties().Extended["version"].(string),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (d *dbmsPairRegisteredManager) testPriConnection(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsPairInstanceDetails)
	if err := testConnection(
		fmt.Sprintf("%s.%s", dt.PriServerName, d.sqlDatabaseDNSSuffix),
		dt.PriAdministratorLogin,
		string(dt.PriAdministratorLoginPassword),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (d *dbmsPairRegisteredManager) testSecConnection(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsPairInstanceDetails)
	if err := testConnection(
		fmt.Sprintf("%s.%s", dt.SecServerName, d.sqlDatabaseDNSSuffix),
		dt.SecAdministratorLogin,
		string(dt.SecAdministratorLoginPassword),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (d *dbmsPairRegisteredManager) deployPriARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsPairInstanceDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := deployDbmsFeARMTemplate(
		&d.armDeployer,
		dt.PriARMDeploymentName,
		instance.ProvisioningParameters.GetString("primaryResourceGroup"),
		instance.ProvisioningParameters.GetString("primaryLocation"),
		dt.PriServerName,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	var ok bool
	dt.PriFullyQualifiedDomainName, ok =
		outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	return dt, err
}

func (d *dbmsPairRegisteredManager) deploySecARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsPairInstanceDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := deployDbmsFeARMTemplate(
		&d.armDeployer,
		dt.SecARMDeploymentName,
		instance.ProvisioningParameters.GetString("secondaryResourceGroup"),
		instance.ProvisioningParameters.GetString("secondaryLocation"),
		dt.SecServerName,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	var ok bool
	dt.SecFullyQualifiedDomainName, ok =
		outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	return dt, err
}
