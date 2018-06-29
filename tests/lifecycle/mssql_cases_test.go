// +build !unit

package lifecycle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/noop"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
	uuid "github.com/satori/go.uuid"
)

var mssqlDBMSAlias = uuid.NewV4().String()
var mssqlDBMSFeAlias = mssqlDBMSAlias + "-fe"

var mssqlTestCases = []serviceLifecycleTestCase{
	{ // all-in-one scenario (dtu-based)
		group:     "mssql",
		name:      "all-in-one (DTU)",
		serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
		planID:    "2497b7f3-341b-4ac6-82fb-d4a48c005e19",
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
		serviceID: "a7454e0e-be2c-46ac-b55f-8c4278117525",
		planID:    "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
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
				serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				planID:          "8fa8d759-c142-45dd-ae38-b93482ddc04a",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
				},
			},
			{ // db only scenario (vcore-based)
				group:           "mssql",
				name:            "database-only (vCore)",
				serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				planID:          "da591616-77a1-4df8-a493-6c119649bc6b",
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
		serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
		planID:    "c77e86af-f050-4457-a2ff-2b48451888f3",
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
	{ // dbms only from existing scenario
		group:     "mssql",
		name:      "dbms-only",
		serviceID: "a7454e0e-be2c-46ac-b55f-8c4278117525",
		planID:    "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
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
		deliverProvisioningParametersToChild: deliverServerNameAndLogin,
		childTestCases: []*serviceLifecycleTestCase{
			{
				group:     "mssql",
				name:      "dbms-fe-only",
				serviceID: "c9bd94e7-5b7d-4b20-96e6-c5678f99d997",
				planID:    "4e95e962-0142-4117-b212-bcc7aec7f6c2",
				provisioningParameters: map[string]interface{}{
					"location": "southcentralus",
					"alias":    mssqlDBMSFeAlias,
				},
				childTestCases: []*serviceLifecycleTestCase{
					{ // dtu db only scenario
						group:           "mssql",
						name:            "database-only (DTU)",
						serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
						planID:          "8fa8d759-c142-45dd-ae38-b93482ddc04a",
						testCredentials: testMsSQLCreds,
						provisioningParameters: map[string]interface{}{
							"parentAlias": mssqlDBMSFeAlias,
						},
					},
					{ // vcore db only scenario
						group:           "mssql",
						name:            "database-only (vCore)",
						serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
						planID:          "da591616-77a1-4df8-a493-6c119649bc6b",
						testCredentials: testMsSQLCreds,
						provisioningParameters: map[string]interface{}{
							"parentAlias": mssqlDBMSFeAlias,
							"cores":       int64(2),
							"storage":     int64(10),
						},
					},
				},
			},
		},
	},
	{ // database only from existing scenario
		group:     "mssql",
		name:      "dbms-only",
		serviceID: "a7454e0e-be2c-46ac-b55f-8c4278117525",
		planID:    "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
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
				serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				planID:          "8fa8d759-c142-45dd-ae38-b93482ddc04a",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
				},
				deliverProvisioningParametersToChild: deliverDatabaseName,
				childTestCases: []*serviceLifecycleTestCase{
					{
						// db only from existing scenario (dtu-based)
						group:           "mssql",
						name:            "database-only-fe (DTU)",
						serviceID:       "b0b2a2f7-9b5e-4692-8b94-24fe2f6a9a8e",
						planID:          "e5804586-625a-4f67-996f-ca19a14711cc",
						testCredentials: testMsSQLCreds,
						provisioningParameters: map[string]interface{}{
							"parentAlias": mssqlDBMSAlias,
						},
					},
				},
			},
			{ // db only scenario (vcore-based)
				group:           "mssql",
				name:            "database-only (vCore)",
				serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				planID:          "da591616-77a1-4df8-a493-6c119649bc6b",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
					"cores":       2,
					"storage":     10,
				},
			},
		},
	},
}

func deliverServerNameAndLogin(
	childPp *map[string]interface{},
	dt service.InstanceDetails,
	svc service.Service,
) {
	if err := crypto.InitializeGlobalCodec(noop.NewCodec()); err != nil {
		panic(err)
	}
	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		panic(err)
	}
	(*childPp)["server"] = dtMap["server"]
	(*childPp)["administratorLogin"] = dtMap["administratorLogin"]
	// https://play.golang.org/p/fWTCTXCw81P
	x, err := json.Marshal(dtMap["administratorLoginPassword"])
	if err != nil {
		panic(err)
	}
	var y service.SecureString
	err = json.Unmarshal(x, &y)
	if err != nil {
		panic(err)
	}
	(*childPp)["administratorLoginPassword"] = string(y)
}

func deliverDatabaseName(
	childPp *map[string]interface{},
	dt service.InstanceDetails,
	svc service.Service,
) {
	if err := crypto.InitializeGlobalCodec(noop.NewCodec()); err != nil {
		panic(err)
	}
	dtMap, err := service.GetMapFromStruct(dt)
	if err != nil {
		panic(err)
	}
	(*childPp)["database"] = dtMap["database"]
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
