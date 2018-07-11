package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func buildBasicPlan(
	id string,
	includesDBMS bool,
	fe bool,
) service.PlanProperties {

	planDetails := dtuPlanDetails{
		storageInGB: 2,
		defaultDTUs: 5,
		tierName:    "Basic",
		skuMap: map[int64]string{
			5: "Basic",
		},
		includeDBMS: includesDBMS,
	}

	planProperties := service.PlanProperties{
		ID:          id,
		Name:        "basic",
		Description: "Basic Tier, 5 DTUs, 2GB Storage, 7 days point-in-time restore",
		Free:        false,
		Stability:   service.StabilityStable,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Basic Tier",
			Bullets: []string{
				"5 DTUs",
				"Includes 2GB Storage",
				"7 days point-in-time restore",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}
	if fe {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getFeProvisionSchema()
	}

	return planProperties
}

func buildStandardPlan(
	id string,
	includesDBMS bool,
	fe bool,
) service.PlanProperties {
	planDetails := dtuPlanDetails{
		storageInGB: 250,
		allowedDTUs: []int64{
			10, 20, 50, 100, 200, 400, 800, 1600, 3000,
		},
		defaultDTUs: 10,
		tierName:    "Standard",
		skuMap: map[int64]string{
			10:   "S0",
			20:   "S1",
			50:   "S2",
			100:  "S3",
			200:  "S4",
			400:  "S6",
			800:  "S7",
			1600: "S9",
			3000: "S12",
		},
		includeDBMS: includesDBMS,
	}

	planProperties := service.PlanProperties{
		ID:   id,
		Name: "standard",
		Description: "Standard Tier, Up to 3000 DTUs, 250GB Storage, " +
			"35 days point-in-time restore",
		Free:      false,
		Stability: service.StabilityStable,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Standard Tier",
			Bullets: []string{
				"10-3000 DTUs",
				"250GB",
				"35 days point-in-time restore",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}
	if fe {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getFeProvisionSchema()
	}

	return planProperties
}

func buildPremiumPlan(
	id string,
	includesDBMS bool,
	fe bool,
) service.PlanProperties {
	planDetails := dtuPlanDetails{
		storageInGB: 500,
		allowedDTUs: []int64{
			125, 250, 500, 1000, 1750, 4000,
		},
		defaultDTUs: 125,
		tierName:    "Premium",
		skuMap: map[int64]string{
			125:  "P1",
			250:  "P2",
			500:  "P4",
			1000: "P6",
			1750: "P11",
			4000: "P15",
		},
		includeDBMS: includesDBMS,
	}

	planProperties := service.PlanProperties{
		ID:   id,
		Name: "premium",
		Description: "Premium Tier, Up to 4000 DTUs, 500GB Storage, " +
			"35 days point-in-time restore",
		Free:      false,
		Stability: service.StabilityStable,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Premium Tier",
			Bullets: []string{
				"Up to 4000 DTUs",
				"Includes 500GB Storage",
				"35 days point-in-time restore",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}
	if fe {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getFeProvisionSchema()
	}

	return planProperties
}

func buildGeneralPurposePlan(
	id string,
	includesDBMS bool,
	fe bool,
) service.PlanProperties {
	planDetails := vCorePlanDetails{
		tierName:      "GeneralPurpose",
		tierShortName: "GP",
		includeDBMS:   includesDBMS,
	}
	planProperties := service.PlanProperties{
		ID:          id,
		Name:        "general-purpose",
		Description: "Up to 80 vCores, 440 GB memory and 1 TB of storage (preview)",
		Free:        false,
		Stability:   service.StabilityPreview,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "General Purpose (preview)",
			Bullets: []string{
				"Scalable compute and storage options for budget-oriented applications",
				"Up to 80 vCores",
				"Up to 440 GB memory",
				"$187.62 / vCore",
				"7 days point-in-time restore",
				"Currently In Preview",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}
	if fe {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getFeProvisionSchema()
	}

	return planProperties
}

func buildBusinessCriticalPlan(
	id string,
	includesDBMS bool,
	fe bool,
) service.PlanProperties {
	planDetails := vCorePlanDetails{
		tierName:      "BusinessCritical",
		tierShortName: "BC",
		includeDBMS:   includesDBMS,
	}
	planProperties := service.PlanProperties{
		ID:   id,
		Name: "business-critical",
		Description: "Up to 80 vCores, 440 GB memory and 1 TB of storage. " +
			"Local SSD, highest resilience to failures. (preview)",
		Free:      false,
		Stability: service.StabilityPreview,
		Extended: map[string]interface{}{
			"tierDetails": planDetails,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Business Critical (preview)",
			Bullets: []string{
				"Up to 80 vCores",
				"Up to 440 GB memory",
				"$505.50 / vCore",
				"7 days point-in-time restore",
				"Currently In Preview",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: planDetails.getProvisionSchema(),
				UpdatingParametersSchema:     planDetails.getUpdateSchema(),
			},
		},
	}
	if fe {
		planProperties.Schemas.ServiceInstances.ProvisioningParametersSchema =
			planDetails.getFeProvisionSchema()
	}

	return planProperties
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {

	return service.NewCatalog([]service.Service{
		// all-in-one (dbms and database) service
		service.NewService(
			service.ServiceProperties{
				ID:          "aa62bb24-1d49-4f2d-905a-d387ae339f3a",
				Name:        "azure-sql-12-0",
				Description: "Azure SQL Database 12.0-- DBMS and single database",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure SQL Database 12.0",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
					LongDescription:  "Azure SQL Database 12.0-- DBMS and single database",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "SQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.allInOneServiceManager,
			service.NewPlan(
				buildBasicPlan(
					"63d62185-d277-4735-96d6-b7cf6a6d128a",
					true,
					false,
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"e5c5d63d-e32f-47ff-9e57-be72872405be",
					true,
					false,
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"ebc10094-7d57-4e59-86f6-e1204632f0e5",
					true,
					false,
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"fcdce498-a183-4031-96e6-229815a4d75c",
					true,
					false,
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"81300e34-43d8-456c-bd25-7b760592f138",
					true,
					false,
				),
			),
		),
		// dbms only service
		service.NewService(
			service.ServiceProperties{
				ID:             "3d07f78a-e15c-4f26-ae82-62a963a7162d",
				Name:           "azure-sql-12-0-dbms",
				Description:    "Azure SQL 12.0-- DBMS only",
				ChildServiceID: "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure SQL 12.0-- DBMS Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
					LongDescription:  "Azure SQL 12.0-- DBMS only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "SQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.dbmsManager,
			service.NewPlan(service.PlanProperties{
				ID:          "d98d557a-983e-4c96-a928-926288583975",
				Name:        "dbms",
				Description: "Azure SQL Server-- DBMS only",
				Free:        false,
				Stability:   service.StabilityPreview,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure SQL Server-- DBMS Only",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.dbmsManager.getProvisionParametersSchema(),
						UpdatingParametersSchema:     m.dbmsManager.getUpdatingParametersSchema(),
					},
				},
			}),
		),
		// database only service
		service.NewService(
			service.ServiceProperties{
				ID:              "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
				Name:            "azure-sql-12-0-database",
				Description:     "Azure SQL 12.0-- database only",
				Bindable:        true,
				ParentServiceID: "3d07f78a-e15c-4f26-ae82-62a963a7162d", // more parents in fact
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure SQL 12.0-- Database Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
					LongDescription:  "Azure SQL 12.0-- database only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{"Azure", "SQL", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databaseManager,
			service.NewPlan(
				buildBasicPlan(
					"756ccc03-e701-4336-a5cd-ea0cf22e597c",
					false,
					false,
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"f9613acc-6ffd-4c9e-acdf-7631d971e7dc",
					false,
					false,
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"df706b83-cf8e-4e88-bd67-ce7feecef7c8",
					false,
					false,
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"8bcd1643-b02c-4d71-8860-c31adae10a6b",
					false,
					false,
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"9f506da2-4f31-4e1b-85b8-9a5dbf380a0f",
					false,
					false,
				),
			),
		),
		// dbms only registered service
		service.NewService(
			service.ServiceProperties{
				ID:             "97c5a775-333f-42a1-bfca-16819ddf7e2e",
				Name:           "azure-sql-12-0-dbms-registered",
				Description:    "Azure SQL 12.0-- DBMS only registered",
				ChildServiceID: "94e4429c-1dd9-4e50-855f-6af2a0f8756e",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure SQL 12.0-- DBMS Only registered",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
					LongDescription:  "Azure SQL 12.0-- DBMS only registered",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "SQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.dbmsRegisteredManager,
			service.NewPlan(service.PlanProperties{
				ID:          "840399dd-5593-493e-80c1-3b21f687997d",
				Name:        "dbms",
				Description: "Azure SQL Server-- DBMS only",
				Free:        false,
				Stability:   service.StabilityPreview,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure SQL Server-- DBMS Only",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.dbmsRegisteredManager.getProvisionParametersSchema(),
						UpdatingParametersSchema:     m.dbmsRegisteredManager.getUpdatingParametersSchema(),
					},
				},
			}),
		),
		// database only from existing service
		service.NewService(
			service.ServiceProperties{
				ID:              "0938a2d1-3490-41fc-a095-d235debff907",
				Name:            "azure-sql-12-0-database-from-existing",
				Description:     "Azure SQL 12.0-- database only from existing",
				Bindable:        true,
				ParentServiceID: "3d07f78a-e15c-4f26-ae82-62a963a7162d", // more parents in fact
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure SQL 12.0-- Database Only from existing",
					ImageURL:         "https://azure.microsoft.com/svghandler/sql-database/?width=200",
					LongDescription:  "Azure SQL 12.0-- database only from existing",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/sql-database/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Tags: []string{"Azure", "SQL", "Database"},
				Extended: map[string]interface{}{
					"version": "12.0",
				},
			},
			m.databaseFeManager,
			service.NewPlan(
				buildBasicPlan(
					"fc2f3117-2539-414d-b5ab-f047fc4c93d4",
					false,
					true,
				),
			),
			service.NewPlan(
				buildStandardPlan(
					"fb475332-23ee-4aca-953e-55fc97577d01",
					false,
					true,
				),
			),
			service.NewPlan(
				buildPremiumPlan(
					"6323a513-98ca-42ca-9ad5-6e78eff8a8fe",
					false,
					true,
				),
			),
			service.NewPlan(
				buildGeneralPurposePlan(
					"f64950ae-9ed3-4639-afa4-c85b1a2dc759",
					false,
					true,
				),
			),
			service.NewPlan(
				buildBusinessCriticalPlan(
					"666d2a9e-a566-4710-a07f-cf712c43701c",
					false,
					true,
				),
			),
		),
	}), nil
}
