// +build !unit

package e2e

import uuid "github.com/satori/go.uuid"

func getPostgreSQLTestCases() []e2eTestCase {
	alias := uuid.NewV4().String()
	return []e2eTestCase{
		{
			group:     "postgresql",
			name:      "all-in-one",
			serviceID: "4d4e2afa-4eb6-4cbd-a321-35f115281ab2",
			planID:    "c79ad81b-3000-4abf-a27f-c8a397d34b41",
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
				"sslEnforcement": "disabled",
				"extensions": []string{
					"uuid-ossp",
					"postgis",
				},
			},
			bind: true,
		},
		{
			group:     "postgresql",
			name:      "dbms-only",
			serviceID: "278c0ee4-7aa6-4f79-953e-3d60034f93b5",
			planID:    "1d6067ba-ec51-4078-bdfe-969c622178de",
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
				{ // database only scenario
					group:     "postgresql",
					name:      "database-only",
					serviceID: "20defa86-7dfc-4c3a-aafc-9f106ac56fcb",
					planID:    "ee762481-19e8-49e6-91dc-38f17336789a",
					provisioningParameters: map[string]interface{}{
						"parentAlias": alias,
						"extensions": []string{
							"uuid-ossp",
							"postgis",
						},
					},
					bind: true,
				},
			},
		},
	}
}
