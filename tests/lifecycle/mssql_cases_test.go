// +build !unit

package lifecycle

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
	uuid "github.com/satori/go.uuid"
)

var mssqlDBMSAlias = uuid.NewV4().String()
var mssqlDBMSRegisteredAlias = uuid.NewV4().String()

var mssqlTestCases = []serviceLifecycleTestCase{
	{ // all-in-one scenario (dtu-based)
		group:     "mssql",
		name:      "all-in-one (DTU)",
		serviceID: "aa62bb24-1d49-4f2d-905a-d387ae339f3a",
		planID:    "e5c5d63d-e32f-47ff-9e57-be72872405be",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"dtus":     200,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowSome",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "35.0.0.0",
				},
				map[string]interface{}{
					"name":           "AllowMore",
					"startIPAddress": "35.0.0.1",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		testCredentials: testMsSQLCreds,
	},
	{ // dbms only scenario
		group:     "mssql",
		name:      "dbms-only",
		serviceID: "3d07f78a-e15c-4f26-ae82-62a963a7162d",
		planID:    "d98d557a-983e-4c96-a928-926288583975",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    mssqlDBMSAlias,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // db only scenario (dtu-based)
				group:           "mssql",
				name:            "database-only (DTU)",
				serviceID:       "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
				planID:          "756ccc03-e701-4336-a5cd-ea0cf22e597c",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
				},
			},
			{ // db only scenario (vcore-based)
				group:           "mssql",
				name:            "database-only (vCore)",
				serviceID:       "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
				planID:          "8bcd1643-b02c-4d71-8860-c31adae10a6b",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
					"cores":       2,
					"storage":     10,
				},
			},
		},
	},
	{ // all-in-one scenario (vcore-based)
		group:     "mssql",
		name:      "all-in-one (vCore)",
		serviceID: "aa62bb24-1d49-4f2d-905a-d387ae339f3a",
		planID:    "fcdce498-a183-4031-96e6-229815a4d75c",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"cores":    4,
			"storage":  25,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
	},
	{ // dbms only registered scenario
		group:        "mssql",
		name:         "dbms-only-registered",
		serviceID:    "97c5a775-333f-42a1-bfca-16819ddf7e2e",
		planID:       "840399dd-5593-493e-80c1-3b21f687997d",
		preProvision: createSQLServer,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    mssqlDBMSRegisteredAlias,
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // dtu db only scenario
				group:           "mssql",
				name:            "database-only (DTU)",
				serviceID:       "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
				planID:          "756ccc03-e701-4336-a5cd-ea0cf22e597c",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSRegisteredAlias,
				},
			},
			{ // vcore db only scenario
				group:           "mssql",
				name:            "database-only (vCore)",
				serviceID:       "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
				planID:          "8bcd1643-b02c-4d71-8860-c31adae10a6b",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSRegisteredAlias,
					"cores":       int64(2),
					"storage":     int64(10),
				},
			},
		},
	},
	{ // database only from existing scenario
		group:     "mssql",
		name:      "dbms-only",
		serviceID: "3d07f78a-e15c-4f26-ae82-62a963a7162d",
		planID:    "d98d557a-983e-4c96-a928-926288583975",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    mssqlDBMSAlias,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		childTestCases: []*serviceLifecycleTestCase{
			{
				// db only from existing scenario (dtu-based)
				group:           "mssql",
				name:            "database-only-fe (DTU)",
				serviceID:       "0938a2d1-3490-41fc-a095-d235debff907",
				planID:          "fc2f3117-2539-414d-b5ab-f047fc4c93d4",
				testCredentials: testMsSQLCreds,
				preProvision:    createSQLDatabase,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
				},
			},
		},
	},
}

func createSQLServer(
	ctx context.Context,
	resourceGroup string,
	_ *service.Instance,
	pp *map[string]interface{},
) error {
	azureConfig, authorizer, err := getAzureConfigAndAuthorizer()
	if err != nil {
		return err
	}
	serversClient := sqlSDK.NewServersClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	serversClient.Authorizer = authorizer
	firewallRulesClient := sqlSDK.NewFirewallRulesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	firewallRulesClient.Authorizer = authorizer

	serverName := uuid.NewV4().String()
	administratorLogin := generate.NewIdentifier()
	administratorLoginPassword := generate.NewPassword()
	version := "12.0"
	location := (*pp)["location"].(string)
	(*pp)["server"] = serverName
	(*pp)["administratorLogin"] = administratorLogin
	(*pp)["administratorLoginPassword"] = administratorLoginPassword

	server := sqlSDK.Server{
		Location: &location,
		ServerProperties: &sqlSDK.ServerProperties{
			AdministratorLogin:         &administratorLogin,
			AdministratorLoginPassword: &administratorLoginPassword,
			Version:                    &version,
		},
	}
	result, err := serversClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		serverName,
		server,
	)
	if err != nil {
		return fmt.Errorf("error creating sql server: %s", err)
	}
	if err := result.WaitForCompletion(ctx, serversClient.Client); err != nil {
		return fmt.Errorf("error creating sql server: %s", err)
	}

	startIPAddress := "0.0.0.0"
	endIPAddress := "255.255.255.255"
	firewallRule := sqlSDK.FirewallRule{
		FirewallRuleProperties: &sqlSDK.FirewallRuleProperties{
			StartIPAddress: &startIPAddress,
			EndIPAddress:   &endIPAddress,
		},
	}
	if _, err := firewallRulesClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		serverName,
		"all",
		firewallRule,
	); err != nil {
		return fmt.Errorf("error creating firewall rule: %s", err)
	}
	return nil
}

func createSQLDatabase(
	ctx context.Context,
	resourceGroup string,
	parent *service.Instance,
	pp *map[string]interface{},
) error {
	azureConfig, authorizer, err := getAzureConfigAndAuthorizer()
	if err != nil {
		return err
	}
	databasesClient := sqlSDK.NewDatabasesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	databasesClient.Authorizer = authorizer

	dtMap, err := service.GetMapFromStruct(parent.Details)
	if err != nil {
		return err
	}
	serverName := dtMap["server"].(string)
	databaseName := generate.NewIdentifier()
	location := parent.ProvisioningParameters.GetString("location")
	database := sqlSDK.Database{
		Location: &location,
	}
	(*pp)["database"] = databaseName

	result, err := databasesClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		serverName,
		databaseName,
		database,
	)
	if err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	if err := result.WaitForCompletion(ctx, databasesClient.Client); err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	return nil
}

func testMsSQLCreds(credentials map[string]interface{}) error {
	query := url.Values{}
	query.Add("database", credentials["database"].(string))
	query.Add("encrypt", "true")
	query.Add("TrustServerCertificate", "true")

	u := &url.URL{
		Scheme: "sqlserver",
		User: url.UserPassword(
			credentials["username"].(string),
			credentials["password"].(string),
		),
		Host: fmt.Sprintf(
			"%s:%d",
			credentials["host"].(string),
			int(credentials["port"].(float64)),
		),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("mssql", u.String())
	if err != nil {
		return fmt.Errorf("error validating the database arguments: %s", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close() // nolint: errcheck

	rows, err := db.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='CONTROL'") // nolint: lll
	if err != nil {
		return fmt.Errorf(
			`error querying SELECT from table fn_my_permissions: %s`,
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return fmt.Errorf(
			`error user doesn't have permission 'CONTROL'`,
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			`error iterating rows`,
		)
	}

	return nil
}
