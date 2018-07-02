package postgresql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func createBasicPlan(
	planID string,
	includeDBParams bool,
) service.PlanProperties {
	td := tierDetails{
		tierName:                "Basic",
		tierShortName:           "B",
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
				ProvisioningParametersSchema: generateProvisioningParamsSchema(
					td,
					includeDBParams,
				),
				UpdatingParametersSchema: generateUpdatingParamsSchema(td),
			},
		},
	}
}

func createGPPlan(
	planID string,
	includeDBParams bool,
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
			"require balanced compute and memory with scalable I/O throughput.",
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
				ProvisioningParametersSchema: generateProvisioningParamsSchema(
					td,
					includeDBParams,
				),
				UpdatingParametersSchema: generateUpdatingParamsSchema(td),
			},
		},
	}
}

func createMemoryOptimizedPlan(
	planID string,
	includeDBParams bool,
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
				ProvisioningParametersSchema: generateProvisioningParamsSchema(
					td,
					includeDBParams,
				),
				UpdatingParametersSchema: generateUpdatingParamsSchema(td),
			},
		},
	}
}

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		// all-in-one
		service.NewService(
			service.ServiceProperties{
				ID:          "4d4e2afa-4eb6-4cbd-a321-35f115281ab2",
				Name:        "azure-postgresql-9-6",
				Description: "Azure Database for PostgreSQL 9.6-- DBMS and single database",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS and single database",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "PostgreSQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "9.6",
				},
			},
			m.allInOneManager,
			service.NewPlan(createBasicPlan("74cb4795-7c06-4ddb-8064-d2cdd3818256", true)),
			service.NewPlan(createGPPlan("e2286a43-0de5-4415-9782-9ed2070cb116", true)),
			service.NewPlan(createMemoryOptimizedPlan("c79ad81b-3000-4abf-a27f-c8a397d34b41", true)),
		),
		// dbms only
		service.NewService(
			service.ServiceProperties{
				ID:             "278c0ee4-7aa6-4f79-953e-3d60034f93b5",
				Name:           "azure-postgresql-9-6-dbms",
				Description:    "Azure Database for PostgreSQL 9.6-- DBMS only",
				ChildServiceID: "20defa86-7dfc-4c3a-aafc-9f106ac56fcb",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6-- DBMS Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "PostgreSQL", "DBMS", "Server", "Database"},
				Extended: map[string]interface{}{
					"version": "9.6",
				},
			},
			m.dbmsManager,
			service.NewPlan(createBasicPlan("1d6067ba-ec51-4078-bdfe-969c622178de", false)),
			service.NewPlan(createGPPlan("d75039e2-f333-472a-b5ed-b43dfbef1771", false)),
			service.NewPlan(createMemoryOptimizedPlan("14986696-b6ac-47ff-8000-203ae3e4ae3b", false)),
		),
		// database only
		service.NewService(
			service.ServiceProperties{
				ID:              "20defa86-7dfc-4c3a-aafc-9f106ac56fcb",
				Name:            "azure-postgresql-9-6-database",
				Description:     "Azure Database for PostgreSQL 9.6-- database only",
				ParentServiceID: "278c0ee4-7aa6-4f79-953e-3d60034f93b5",
				Metadata: service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL 9.6-- Database Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- database only",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "PostgreSQL", "Database"},
				Extended: map[string]interface{}{
					"version": "9.6",
				},
			},
			m.databaseManager,
			service.NewPlan(service.PlanProperties{
				ID:          "ee762481-19e8-49e6-91dc-38f17336789a",
				Name:        "database",
				Description: "A new database added to an existing DBMS",
				Free:        false,
				Stability:   service.StabilityStable,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Azure Database for PostgreSQL-- Database Only",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.databaseManager.getProvisionParametersSchema(),
					},
				},
			}),
		),
	}), nil
}
