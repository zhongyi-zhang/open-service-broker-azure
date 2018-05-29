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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pp := dbmsFeProvisioningParams{}
	if err := service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, nil, err
	}
	spp := secureDBMSFeProvisioningParams{}
	if err := service.GetStructFromMap(instance.SecureProvisioningParameters, &spp); err != nil {
		return nil, nil, err
	}
	dt := dbmsInstanceDetails{
		ARMDeploymentName:  uuid.NewV4().String(),
		ServerName:         pp.ServerName,
		AdministratorLogin: pp.AdministratorLogin,
	}
	sdt := secureDBMSInstanceDetails{
		AdministratorLoginPassword: spp.AdministratorLoginPassword,
	}
	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		return nil, nil, err
	}
	sdtMap, err := service.GetMapFromStruct(sdt)
	return dtMap, sdtMap, err
}

func (d *dbmsFeManager) getServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := dbmsInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	result, err := d.serversClient.Get(
		ctx,
		instance.ResourceGroup,
		dt.ServerName,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting sql server: %s", err)
	}
	if result.Name == nil {
		err = fmt.Errorf(
			"can't find sql server %s in the resource group %s",
			dt.ServerName,
			instance.ResourceGroup,
		)
		return nil, nil, err
	}
	expectedVersion := instance.Service.GetProperties().Extended["version"]
	if *result.Version != expectedVersion {
		return nil, nil, fmt.Errorf(
			"sql server version validation failed, "+
				"expected version: %s, current version: %s",
			expectedVersion,
			result.Version,
		)
	}
	return instance.Details, instance.SecureDetails, nil
}

func (d *dbmsFeManager) testConnection(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	pp := dbmsFeProvisioningParams{}
	if err := service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, nil, err
	}
	spp := secureDBMSFeProvisioningParams{}
	if err := service.GetStructFromMap(instance.SecureProvisioningParameters, &spp); err != nil {
		return nil, nil, err
	}

	// connect to master database to create login
	masterDb, err := getDBConnection(
		pp.AdministratorLogin,
		spp.AdministratorLoginPassword,
		fmt.Sprintf("%s.%s", pp.ServerName, d.sqlDatabaseDNSSuffix),
		"master",
	)
	if err != nil {
		return nil, nil, err
	}
	defer masterDb.Close() // nolint: errcheck

	// Is there a better approach to verify if it is a sys admin?
	rows, err := masterDb.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='ALTER ANY USER'") // nolint: lll
	if err != nil {
		return nil, nil, fmt.Errorf(
			`error querying SELECT from table fn_my_permissions: %s`,
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return nil, nil, fmt.Errorf(
			`error user doesn't have permission 'ALTER ANY USER'`,
		)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf(
			`error iterating rows`,
		)
	}

	return instance.Details, instance.SecureDetails, nil
}

func (d *dbmsFeManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := dbmsInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	goTemplateParams, err := buildDBMSGoTemplateParameters(instance)
	if err != nil {
		return nil, nil, err
	}
	outputs, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		dbmsARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	var ok bool
	dt.FullyQualifiedDomainName, ok = outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, instance.SecureDetails, err
}
