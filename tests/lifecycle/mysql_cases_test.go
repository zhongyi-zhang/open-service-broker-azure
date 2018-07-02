// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL SQL driver
	uuid "github.com/satori/go.uuid"
)

var mysqlDBMSAlias = uuid.NewV4().String()

var mysqlTestCases = []serviceLifecycleTestCase{
	{
		group:     "mysql",
		name:      "all-in-one",
		serviceID: "3c715189-9843-4d8b-bb21-6ae653ad95c5",
		planID:    "643038f4-0343-4d94-8daf-738334ede7b6",
		provisioningParameters: map[string]interface{}{
			"location":       "southcentralus",
			"sslEnforcement": "disabled",
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
		testCredentials: testMySQLCreds,
	},
	{
		group:     "mysql",
		name:      "dbms-only",
		serviceID: "ef21a7aa-fb6b-457c-b43d-bb0081334332",
		planID:    "de271154-2f6c-4004-94f8-81e37a26178b",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    mysqlDBMSAlias,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // database only scenario
				group:           "mysql",
				name:            "database-only",
				serviceID:       "5f91e726-abb2-43db-a96d-4abf2e06ae28",
				planID:          "98e18e2e-6b03-4935-9146-0f71106610a0",
				testCredentials: testMySQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mysqlDBMSAlias,
				},
			},
		},
	},
}

func testMySQLCreds(credentials map[string]interface{}) error {

	var connectionStrTemplate string
	if credentials["sslRequired"].(bool) {
		connectionStrTemplate =
			"%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=true"
	} else {
		connectionStrTemplate =
			"%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true"
	}

	db, err := sql.Open("mysql", fmt.Sprintf(
		connectionStrTemplate,
		credentials["username"].(string),
		credentials["password"].(string),
		credentials["host"].(string),
		credentials["database"].(string),
	))
	if err != nil {
		return fmt.Errorf("error validating the database arguments: %s", err)
	}
	defer db.Close() // nolint: errcheck
	rows, err := db.Query("SELECT * from INFORMATION_SCHEMA.TABLES")
	if err != nil {
		return fmt.Errorf("error validating the database arguments: %s", err)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return fmt.Errorf(
			`error could not select from INFORMATION_SCHEMA.TABLES'`,
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			`error iterating rows`,
		)
	}
	return nil
}
