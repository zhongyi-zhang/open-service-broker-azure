// +build !unit

package e2e

import uuid "github.com/satori/go.uuid"

func getMSSQLTestCases() []e2eTestCase {
	alias := uuid.NewV4().String()
	return []e2eTestCase{
		{ // all-in-one scenario (dtu-based)
			group:     "mssql",
			name:      "all-in-one",
			serviceID: "aa62bb24-1d49-4f2d-905a-d387ae339f3a",
			planID:    "63d62185-d277-4735-96d6-b7cf6a6d128a",
			provisioningParameters: map[string]interface{}{
				"location":      "southcentralus",
				"resourceGroup": "placeholder",
				"firewallRules": []map[string]string{
					{
						"name":           "AllowSome",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "35.0.0.0",
					},
					{
						"name":           "AllowMore",
						"startIPAddress": "35.0.0.1",
						"endIPAddress":   "255.255.255.255",
					},
				},
			},
			bind: true,
		},
		{ // dbms only scenario
			group:     "mssql",
			name:      "dbms-only",
			serviceID: "3d07f78a-e15c-4f26-ae82-62a963a7162d",
			planID:    "d98d557a-983e-4c96-a928-926288583975",
			provisioningParameters: map[string]interface{}{
				"alias":         alias,
				"location":      "southcentralus",
				"resourceGroup": "placeholder",
				"firewallRules": []map[string]string{
					{
						"name":           "AllowAll",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "255.255.255.255",
					},
				},
			},
			childTestCases: []*e2eTestCase{
				{ // db only scenario (dtu-based)
					group:     "mssql",
					name:      "database-only (DTU)",
					serviceID: "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
					planID:    "756ccc03-e701-4336-a5cd-ea0cf22e597c",
					provisioningParameters: map[string]interface{}{
						"parentAlias": alias,
					},
					bind: true,
				},
				{ // db only scenario (vcore-based)
					group:     "mssql",
					name:      "database-only (vCore)",
					serviceID: "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
					planID:    "8bcd1643-b02c-4d71-8860-c31adae10a6b",
					provisioningParameters: map[string]interface{}{
						"parentAlias": alias,
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
				"location":      "southcentralus",
				"resourceGroup": "placeholder",
				"cores":         4,
				"storage":       25,
				"firewallRules": []interface{}{
					map[string]interface{}{
						"name":           "AllowAll",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "255.255.255.255",
					},
				},
			},
		},
	}
}
