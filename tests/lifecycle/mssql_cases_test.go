// +build !unit

package lifecycle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

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
		deliverProvisioningParametersToChild: deliverServerNameAndLogin,
		childTestCases: []*serviceLifecycleTestCase{
			{
				group:     "mssql",
				name:      "dbms-only-registered",
				serviceID: "97c5a775-333f-42a1-bfca-16819ddf7e2e",
				planID:    "840399dd-5593-493e-80c1-3b21f687997d",
				provisioningParameters: map[string]interface{}{
					"location": "southcentralus",
					"alias":    mssqlDBMSFeAlias,
				},
				childTestCases: []*serviceLifecycleTestCase{
					{ // dtu db only scenario
						group:           "mssql",
						name:            "database-only (DTU)",
						serviceID:       "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
						planID:          "756ccc03-e701-4336-a5cd-ea0cf22e597c",
						testCredentials: testMsSQLCreds,
						provisioningParameters: map[string]interface{}{
							"parentAlias": mssqlDBMSFeAlias,
						},
					},
					{ // vcore db only scenario
						group:           "mssql",
						name:            "database-only (vCore)",
						serviceID:       "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
						planID:          "8bcd1643-b02c-4d71-8860-c31adae10a6b",
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
				deliverProvisioningParametersToChild: deliverDatabaseName,
				childTestCases: []*serviceLifecycleTestCase{
					{
						// db only from existing scenario (dtu-based)
						group:           "mssql",
						name:            "database-only-fe (DTU)",
						serviceID:       "0938a2d1-3490-41fc-a095-d235debff907",
						planID:          "fc2f3117-2539-414d-b5ab-f047fc4c93d4",
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
}

func deliverServerNameAndLogin(
	childPp *map[string]interface{},
	dt service.InstanceDetails,
	svc service.Service,
) {
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
