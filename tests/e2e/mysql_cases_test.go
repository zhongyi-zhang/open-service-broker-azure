// +build !unit

package e2e

import uuid "github.com/satori/go.uuid"

func getMySQLTestCases() []e2eTestCase {
	alias := uuid.NewV4().String()
	return []e2eTestCase{
		{
			group:     "mysql",
			name:      "all-in-one",
			serviceID: "3c715189-9843-4d8b-bb21-6ae653ad95c5",
			planID:    "643038f4-0343-4d94-8daf-738334ede7b6",
			provisioningParameters: map[string]interface{}{
				"location":       "southcentralus",
				"resourceGroup":  "placeholder",
				"sslEnforcement": "disabled",
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
		{
			group:     "mysql",
			name:      "dbms-only",
			serviceID: "ef21a7aa-fb6b-457c-b43d-bb0081334332",
			planID:    "de271154-2f6c-4004-94f8-81e37a26178b",
			provisioningParameters: map[string]interface{}{
				"alias":         alias,
				"location":      "eastus",
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
				{
					group:     "mysql",
					name:      "database-only",
					serviceID: "5f91e726-abb2-43db-a96d-4abf2e06ae28",
					planID:    "98e18e2e-6b03-4935-9146-0f71106610a0",
					provisioningParameters: map[string]interface{}{
						"parentAlias": alias,
					},
					bind: true,
				},
			},
		},
	}
}
