package mssql

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databaseFeManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("getDatabase", d.getDatabase),
		service.NewProvisioningStep("testConnection", d.testConnection),
		service.NewProvisioningStep("deployARMTemplate", d.deployARMTemplate),
	)
}

func (d *databaseFeManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	return &databaseInstanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		DatabaseName:      instance.ProvisioningParameters.GetString("database"),
	}, nil
}

func (d *databaseFeManager) getDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt := instance.Details.(*databaseInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	resourceGroup :=
		instance.Parent.ProvisioningParameters.GetString("resourceGroup")
	result, err := d.databasesClient.Get(
		ctx,
		resourceGroup,
		pdt.ServerName,
		dt.DatabaseName,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf("error getting sql database: %s", err)
	}
	if result.Name == nil {
		err = fmt.Errorf(
			"can't find sql database %s in server %s in the resource group %s",
			dt.DatabaseName,
			pdt.ServerName,
			resourceGroup,
		)
		return nil, err
	}
	return instance.Details, nil
}

func (d *databaseFeManager) testConnection(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databaseInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	// connect to master database to create login
	db, err := getDBConnection(
		pdt.AdministratorLogin,
		string(pdt.AdministratorLoginPassword),
		pdt.FullyQualifiedDomainName,
		dt.DatabaseName,
	)
	if err != nil {
		return nil, err
	}
	defer db.Close() // nolint: errcheck

	// Is there a better approach to verify if it is a sys admin?
	rows, err := db.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='ALTER'") // nolint: lll
	if err != nil {
		return nil, fmt.Errorf(
			`error querying SELECT from table fn_my_permissions: %s`,
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return nil, fmt.Errorf(
			`error parent server user doesn't have permission 'ALTER'`,
		)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			`error iterating rows`,
		)
	}

	return instance.Details, nil
}

func (d *databaseFeManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*databaseInstanceDetails)
	pdt := instance.Parent.Details.(*dbmsInstanceDetails)
	goTemplateParams := map[string]interface{}{}
	goTemplateParams["serverName"] = pdt.ServerName
	goTemplateParams["location"] =
		instance.Parent.ProvisioningParameters.GetString("location")
	goTemplateParams["databaseName"] = dt.DatabaseName
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	// No output, so ignore the output
	_, err := d.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.Parent.ProvisioningParameters.GetString("resourceGroup"),
		instance.Parent.ProvisioningParameters.GetString("location"),
		databaseFeARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}
