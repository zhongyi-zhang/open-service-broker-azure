package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *dbmsFeManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("getServer", d.getServer),
		service.NewProvisioningStep("testConnection", d.testConnection),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (d *dbmsFeManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	return &dbmsInstanceDetails{
		ARMDeploymentName:          uuid.NewV4().String(),
		ServerName:                 pp.GetString("server"),
		AdministratorLogin:         pp.GetString("administratorLogin"),
		AdministratorLoginPassword: service.SecureString(pp.GetString("administratorLoginPassword")), // nolint: lll
	}, nil
}

func (d *dbmsFeManager) getServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pp := instance.ProvisioningParameters
	dt := instance.Details.(*dbmsInstanceDetails)
	resourceGroup := pp.GetString("resourceGroup")
	result, err := d.serversClient.Get(
		ctx,
		resourceGroup,
		dt.ServerName,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting sql server: %s", err)
	}
	if result.Name == nil {
		err = fmt.Errorf(
			"can't find sql server %s in the resource group %s",
			dt.ServerName,
			resourceGroup,
		)
		return nil, err
	}
	expectedVersion :=
		instance.Service.GetProperties().Extended["version"].(string)
	if *result.Version != expectedVersion {
		return nil, fmt.Errorf(
			"sql server version validation failed, "+
				"expected version: %s, current version: %s",
			expectedVersion,
			result.Version,
		)
	}
	return instance.Details, nil
}

func (d *dbmsFeManager) testConnection(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	// connect to master database to create login
	masterDb, err := getDBConnection(
		dt.AdministratorLogin,
		string(dt.AdministratorLoginPassword),
		fmt.Sprintf("%s.%s", dt.ServerName, d.sqlDatabaseDNSSuffix),
		"master",
	)
	if err != nil {
		return nil, err
	}
	defer masterDb.Close() // nolint: errcheck

	// Is there a better approach to verify if it is a sys admin?
	rows, err := masterDb.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='ALTER ANY USER'") // nolint: lll
	if err != nil {
		return nil, fmt.Errorf(
			`error querying SELECT from table fn_my_permissions: %s`,
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return nil, fmt.Errorf(
			`error user doesn't have permission 'ALTER ANY USER'`,
		)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			`error iterating rows`,
		)
	}

	return instance.Details, nil
}

func (d *dbmsFeManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*dbmsInstanceDetails)
	version := instance.Service.GetProperties().Extended["version"].(string)
	goTemplateParams, err := buildDBMSGoTemplateParameters(
		dt,
		*instance.ProvisioningParameters,
		version,
	)
	if err != nil {
		return nil, err
	}
	goTemplateParams["location"] =
		instance.ProvisioningParameters.GetString("location")
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	outputs, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		dbmsARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	var ok bool
	dt.FullyQualifiedDomainName, ok = outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	return dt, err
}
