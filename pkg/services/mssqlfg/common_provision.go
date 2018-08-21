package mssqlfg

import (
	"context"
	"fmt"
	"strings"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func validateServer(
	ctx context.Context,
	serversClient *sqlSDK.ServersClient,
	resourceGroup string,
	serverName string,
	expectedVersion string,
	expectedLocation string,
) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	result, err := serversClient.Get(
		ctx,
		resourceGroup,
		serverName,
	)
	if err != nil {
		return "", fmt.Errorf("error getting the sql server: %s", err)
	}
	if result.Name == nil {
		return "", fmt.Errorf(
			"can't find sql server %s in the resource group %s",
			serverName,
			resourceGroup,
		)
	}
	if *result.Version != expectedVersion {
		return "", fmt.Errorf(
			"sql server version validation failed, "+
				"expected version: %s, actual version: %s",
			expectedVersion,
			result.Version,
		)
	}
	expectedLocation = strings.Replace(
		strings.ToLower(expectedLocation),
		" ",
		"",
		-1,
	)
	actualLocation := strings.Replace(
		strings.ToLower(*result.Location),
		" ",
		"",
		-1,
	)
	if expectedLocation != actualLocation {
		return "", fmt.Errorf(
			"sql server location validation failed, "+
				"expected location: %s, actual location: %s",
			expectedLocation,
			actualLocation,
		)
	}
	if *result.FullyQualifiedDomainName == "" {
		return "", fmt.Errorf(
			"sql server details doesn't contain FQDN",
		)
	}
	return *result.FullyQualifiedDomainName, nil
}

func validateDatabase(
	ctx context.Context,
	databasesClient *sqlSDK.DatabasesClient,
	resourceGroup string,
	serverName string,
	databaseName string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, err := databasesClient.Get(
		ctx,
		resourceGroup,
		serverName,
		databaseName,
		"",
	)
	if err != nil {
		return fmt.Errorf("error getting the sql database: %s", err)
	}
	// TODO: add the plan as param and validate?
	return nil
}

func validateFailoverGroup(
	ctx context.Context,
	failoverGroupsClient *sqlSDK.FailoverGroupsClient,
	resourceGroup string,
	serverName string,
	failoverGroupName string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, err := failoverGroupsClient.Get(
		ctx,
		resourceGroup,
		serverName,
		failoverGroupName,
	)
	if err != nil {
		return fmt.Errorf("error getting the failover group: %s", err)
	}
	// TODO: check details to validate roles
	return nil
}

func testConnection(
	fqdn string,
	administratorLogin string,
	administratorLoginPassword string,
) error {
	masterDb, err := getDBConnection(
		administratorLogin,
		administratorLoginPassword,
		fqdn,
		"master",
	)
	if err != nil {
		return err
	}
	defer masterDb.Close() // nolint: errcheck
	// TODO: Is there a better approach to verify if it is a sys admin?
	rows, err := masterDb.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='ALTER ANY USER'") // nolint: lll
	if err != nil {
		return fmt.Errorf(
			`error querying SELECT from table fn_my_permissions: %s`,
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return fmt.Errorf(
			`error user doesn't have permission 'ALTER ANY USER'`,
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(`server error iterating rows`)
	}
	return nil
}

func buildDatabaseGoTemplateParameters(
	databaseName string,
	pp service.ProvisioningParameters,
	pd planDetails,
) (map[string]interface{}, error) {
	td, err := pd.getTierProvisionParameters(pp)
	if err != nil {
		return nil, err
	}
	p := map[string]interface{}{}
	p["databaseName"] = databaseName
	for key, value := range td {
		p[key] = value
	}
	return p, nil
}

func deployDatabaseARMTemplate(
	armDeployer *arm.Deployer,
	armDeploymentName string,
	resourceGroup string,
	location string,
	serverName string,
	databaseName string,
	pp service.ProvisioningParameters,
	pd planDetails,
	tags map[string]string,
) error {
	goTemplateParams, err := buildDatabaseGoTemplateParameters(
		databaseName,
		pp,
		pd,
	)
	if err != nil {
		return err
	}
	goTemplateParams["location"] = location
	goTemplateParams["serverName"] = serverName
	_, err = (*armDeployer).Deploy(
		armDeploymentName,
		resourceGroup,
		location,
		databaseARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	return err
}

func deployDatabaseFeARMTemplate(
	armDeployer *arm.Deployer,
	armDeploymentName string,
	resourceGroup string,
	location string,
	serverName string,
	databaseName string,
	tags map[string]string,
) error {
	goTemplateParams := map[string]interface{}{}
	goTemplateParams["location"] = location
	goTemplateParams["serverName"] = serverName
	goTemplateParams["databaseName"] = databaseName
	_, err := (*armDeployer).Deploy(
		armDeploymentName,
		resourceGroup,
		location,
		databaseFeARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	return err
}

func deployFailoverGroupARMTemplate(
	armDeployer *arm.Deployer,
	instance service.Instance,
) error {
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	dt := instance.Details.(*databasePairInstanceDetails)
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	goTemplateParams := map[string]interface{}{}
	goTemplateParams["priServerName"] = pdt.PriServerName
	goTemplateParams["secServerName"] = pdt.SecServerName
	goTemplateParams["failoverGroupName"] = pp.GetString("failoverGroup")
	goTemplateParams["databaseName"] = pp.GetString("database")
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	_, err := (*armDeployer).Deploy(
		dt.FgARMDeploymentName,
		ppp.GetString("primaryResourceGroup"),
		ppp.GetString("primaryLocation"),
		failoverGroupARMTemplateBytes,
		goTemplateParams,
		map[string]interface{}{}, // empty arm params
		tags,
	)
	return err
}
