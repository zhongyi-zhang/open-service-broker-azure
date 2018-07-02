// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres SQL driver
	uuid "github.com/satori/go.uuid"
)

var postgresqlDBMSAlias = uuid.NewV4().String()

var postgresqlTestCases = []serviceLifecycleTestCase{
	{
		group:           "postgresql",
		name:            "all-in-one",
		serviceID:       "4d4e2afa-4eb6-4cbd-a321-35f115281ab2",
		planID:          "c79ad81b-3000-4abf-a27f-c8a397d34b41",
		testCredentials: testPostgreSQLCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
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
			"sslEnforcement": "disabled",
			"extensions": []interface{}{
				"uuid-ossp",
				"postgis",
			},
		},
	},
	{
		group:     "postgresql",
		name:      "dbms-only",
		serviceID: "278c0ee4-7aa6-4f79-953e-3d60034f93b5",
		planID:    "1d6067ba-ec51-4078-bdfe-969c622178de",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    postgresqlDBMSAlias,
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
				group:           "postgresql",
				name:            "database-only",
				serviceID:       "20defa86-7dfc-4c3a-aafc-9f106ac56fcb",
				planID:          "ee762481-19e8-49e6-91dc-38f17336789a",
				testCredentials: testPostgreSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": postgresqlDBMSAlias,
					"extensions": []interface{}{
						"uuid-ossp",
						"postgis",
					},
				},
			},
		},
	},
}

func testPostgreSQLCreds(credentials map[string]interface{}) error {
	var connectionStrTemplate string
	if credentials["sslRequired"].(bool) {
		connectionStrTemplate =
			"postgres://%s:%s@%s/%s?sslmode=require"
	} else {
		connectionStrTemplate = "postgres://%s:%s@%s/%s"
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
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
	rows, err := db.Query(`
			SELECT * from pg_catalog.pg_tables
			WHERE
			schemaname != 'pg_catalog'
			AND schemaname != 'information_schema'
			`)
	if err != nil {
		return fmt.Errorf("error validating the database arguments: %s", err)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return fmt.Errorf(
			`error could not select from pg_catalog'`,
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			`error iterating rows`,
		)
	}
	return nil
}
