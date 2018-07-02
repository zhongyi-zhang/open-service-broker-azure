package mysql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func createBasicPlan(
	planID string,
) service.PlanProperties {
	td := tierDetails{
		tierShortName:           "B",
		tierName:                "Basic",
		allowedCores:            []int64{1, 2},
		defaultCores:            1,
		maxStorage:              1024,
		allowedBackupRedundancy: []string{"local"},
	}

	return service.PlanProperties{
		ID:   planID,
		Name: "basic",
		Description: "Basic Tier-- For workloads that require light compute and " +
			"I/O performance.",
		Free:      false,
		Stability: service.StabilityStable,
		Extended: map[string]interface{}{
			"tierDetails": td,
		},
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Basic Tier",
			Bullets:     []string{"Up to 2 vCores", "Variable I/O performance"},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateProvisioningParamsSchema(td),
				UpdatingParametersSchema:     generateUpdatingParamsSchema(td),
			},
		},
	}
}

func createGPPlan(
	planID string,
) service.PlanProperties {

	td := tierDetails{
		tierName:                "GeneralPurpose",
		tierShortName:           "GP",
		allowedCores:            []int64{2, 4, 8, 16, 32},
		defaultCores:            2,
		maxStorage:              2048,
		allowedBackupRedundancy: []string{"local", "geo"},
	}
	extendedPlanData := map[string]interface{}{
		"tierDetails": td,
	}

	return service.PlanProperties{
		ID:   planID,
		Name: "general-purpose",
		Description: "General Purpose Tier-- For most business workloads that " +
			"require balanced compute and memory with scalable I/O throughput. ",
		Free:      false,
		Stability: service.StabilityStable,
		Extended:  extendedPlanData,
		Metadata: service.ServicePlanMetadata{
			DisplayName: "General Purpose Tier",
			Bullets: []string{
				"Up to 32 vCores",
				"Predictable I/O Performance",
				"Local or Geo-Redundant Backups",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateProvisioningParamsSchema(td),
				UpdatingParametersSchema:     generateUpdatingParamsSchema(td),
			},
		},
	}
}

func createMemoryOptimizedPlan(
	planID string,
) service.PlanProperties {

	td := tierDetails{
		tierName:                "MemoryOptimized",
		tierShortName:           "MO",
		allowedCores:            []int64{2, 4, 8, 16},
		defaultCores:            2,
		maxStorage:              2048,
		allowedBackupRedundancy: []string{"local", "geo"},
	}
	extendedPlanData := map[string]interface{}{
		"tierDetails": td,
	}

	return service.PlanProperties{
		ID:   planID,
		Name: "memory-optimized",
		Description: "Memory Optimized Tier-- For high-performance database " +
			"workloads that require in-memory performance for faster transaction " +
			"processing and higher concurrency.",
		Free:      false,
		Stability: service.StabilityStable,
		Extended:  extendedPlanData,
		Metadata: service.ServicePlanMetadata{
			DisplayName: "Memory Optimized Tier",
			Bullets: []string{
				"Up to 16 memory optimized vCores",
				"Predictable I/O Performance",
				"Local or Geo-Redundant Backups",
			},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateProvisioningParamsSchema(td),
				UpdatingParametersSchema:     generateUpdatingParamsSchema(td),
			},
		},
	}
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "3c715189-9843-4d8b-bb21-6ae653ad95c5",
				Name:        "azure-mysql-5-7",
				Description: "Azure Database for MySQL 5.7-- DBMS and single database",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for MySQL 5.7",
					ImageURL:         "https://azure.microsoft.com/svghandler/mysql/?width=200",
					LongDescription:  "Azure Database for MySQL 5.7-- DBMS and single database",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/mysql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "MySQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.allInOneServiceManager,
			service.NewPlan(createBasicPlan("284806af-1689-4d02-8ffb-19509483202f")),
			service.NewPlan(createGPPlan("643038f4-0343-4d94-8daf-738334ede7b6")),
			service.NewPlan(createMemoryOptimizedPlan("18ff0626-7122-4803-a66a-b59b6ccbb795")),
		),
		// dbms only service
		service.NewService(
			service.ServiceProperties{
				ID:             "ef21a7aa-fb6b-457c-b43d-bb0081334332",
				Name:           "azure-mysql-5-7-dbms",
				Description:    "Azure Database for MySQL 5.7-- DBMS only",
				ChildServiceID: "5f91e726-abb2-43db-a96d-4abf2e06ae28",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for MySQL 5.7-- DBMS Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/mysql/?width=200",
					LongDescription:  "Azure Database for MySQL 5.7-- DBMS only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/mysql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "MySQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.dbmsManager,
			service.NewPlan(createBasicPlan("db42bac9-8be2-4354-8c9d-c210dc0f4e3b")),
			service.NewPlan(createGPPlan("a9413ad4-1925-4a65-9352-563128ddef36")),
			service.NewPlan(createMemoryOptimizedPlan("de271154-2f6c-4004-94f8-81e37a26178b")),
		),
		// database only service
		service.NewService(
			service.ServiceProperties{
				ID:              "5f91e726-abb2-43db-a96d-4abf2e06ae28",
				Name:            "azure-mysql-5-7-database",
				Description:     "Azure Database for MySQL 5.7-- database only",
				ParentServiceID: "ef21a7aa-fb6b-457c-b43d-bb0081334332",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for MySQL 5.7-- Database Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/mysql/?width=200",
					LongDescription:  "Azure Database for MySQL 5.7-- database only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/mysql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "MySQL", "Database"},
				Extended: map[string]interface{}{
					"version": "5.7",
				},
			},
			m.databaseManager,
			service.NewPlan(service.PlanProperties{
				ID:          "98e18e2e-6b03-4935-9146-0f71106610a0",
				Name:        "database",
				Description: "A new database added to an existing DBMS",
				Free:        false,
				Stability:   service.StabilityStable,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure Database for MySQL-- Database Only",
				},
			}),
		),
	}), nil
}
