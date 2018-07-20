package mssqlfg

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type planDetails interface {
	getProvisionSchema() service.InputParametersSchema
	getTierProvisionParameters(
		pp service.ProvisioningParameters,
	) (map[string]interface{}, error)
	getUpdateSchema() service.InputParametersSchema
	validateUpdateParameters(service.Instance) error
}

type dtuPlanDetails struct {
	tierName    string
	skuMap      map[int64]string
	allowedDTUs []int64
	defaultDTUs int64
	storageInGB int64
}

func (d dtuPlanDetails) validateUpdateParameters(service.Instance) error {
	return nil // no op
}

func (d dtuPlanDetails) getUpdateSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	// Basic tier is constrained to just 5 DTUs, so don't present this as an
	// option
	if len(d.allowedDTUs) > 0 {
		ips.PropertySchemas["dtus"] = &service.IntPropertySchema{
			AllowedValues: d.allowedDTUs,
			DefaultValue:  ptr.ToInt64(d.defaultDTUs),
			Description: "DTUs are a bundled measure of compute, " +
				"storage, and IO resources.",
		}
	}
	return ips
}

func (d dtuPlanDetails) getProvisionSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		RequiredProperties: []string{
			"failoverGroup",
			"database",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"failoverGroup": &service.StringPropertySchema{
				Description: "The name of the failover group",
			},
			"database": &service.StringPropertySchema{
				Description: "The name of the database",
			},
		},
	}
	// Basic tier is constrained to just 5 DTUs, so don't present this as an
	// option
	if len(d.allowedDTUs) > 0 {
		ips.PropertySchemas["dtus"] = &service.IntPropertySchema{
			AllowedValues: d.allowedDTUs,
			DefaultValue:  ptr.ToInt64(d.defaultDTUs),
			Description: "DTUs are a bundled measure of compute, " +
				"storage, and IO resources.",
		}
	}
	return ips
}

func (d dtuPlanDetails) getTierProvisionParameters(
	pp service.ProvisioningParameters,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = d.getSKU(pp)
	p["tier"] = d.tierName
	// ARM template needs bytes
	p["maxSizeBytes"] = pp.GetInt64("storage") * 1024 * 1024 * 1024
	return p, nil
}

func (d dtuPlanDetails) getSKU(pp service.ProvisioningParameters) string {
	// Basic tier is constrained to just 5 DTUs, if this is the basic tier, there
	// is no dtus param. We can infer this is the case if the tier details don't
	// tell us there's a choice.
	if len(d.allowedDTUs) == 0 {
		return d.skuMap[d.defaultDTUs]
	}
	return d.skuMap[pp.GetInt64("dtus")]
}

type vCorePlanDetails struct {
	tierName      string
	tierShortName string
}

func (v vCorePlanDetails) validateUpdateParameters(
	instance service.Instance,
) error {
	return validateStorageUpdate(
		*instance.ProvisioningParameters,
		*instance.UpdatingParameters,
	)
}

func (v vCorePlanDetails) getUpdateSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{},
	}
	ips.PropertySchemas["cores"] = &service.IntPropertySchema{
		AllowedValues: []int64{2, 4, 8, 16, 24, 32, 48, 80},
		DefaultValue:  ptr.ToInt64(2),
		Description:   "A virtual core represents the logical CPU",
	}
	ips.PropertySchemas["storage"] = &service.IntPropertySchema{
		MinValue:     ptr.ToInt64(5),
		MaxValue:     ptr.ToInt64(1024),
		DefaultValue: ptr.ToInt64(10),
		Description:  "The maximum data storage capacity (in GB)",
	}
	return ips
}

func (v vCorePlanDetails) getProvisionSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		RequiredProperties: []string{
			"failoverGroup",
			"database",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"failoverGroup": &service.StringPropertySchema{
				Description: "The name of the failover group",
			},
			"database": &service.StringPropertySchema{
				Description: "The name of the database",
			},
			"cores": &service.IntPropertySchema{
				AllowedValues: []int64{2, 4, 8, 16, 24, 32, 48, 80},
				DefaultValue:  ptr.ToInt64(2),
				Description:   "A virtual core represents the logical CPU",
			},
			"storage": &service.IntPropertySchema{
				MinValue:     ptr.ToInt64(5),
				MaxValue:     ptr.ToInt64(1024),
				DefaultValue: ptr.ToInt64(10),
				Description:  "The maximum data storage capacity (in GB)",
			},
		},
	}
	return ips
}

func (v vCorePlanDetails) getTierProvisionParameters(
	pp service.ProvisioningParameters,
) (map[string]interface{}, error) {
	p := map[string]interface{}{}
	p["sku"] = v.getSKU(pp)
	p["tier"] = v.tierName
	// ARM template needs bytes
	p["maxSizeBytes"] = pp.GetInt64("storage") * 1024 * 1024 * 1024
	return p, nil
}

func (v vCorePlanDetails) getSKU(pp service.ProvisioningParameters) string {
	return fmt.Sprintf(
		"%s_Gen5_%d",
		v.tierShortName,
		pp.GetInt64("cores"),
	)
}

func validateStorageUpdate(
	pp service.ProvisioningParameters,
	up service.ProvisioningParameters,
) error {
	existingStorage := pp.GetInt64("storage")
	newStorge := up.GetInt64("storage")
	if newStorge < existingStorage {
		return service.NewValidationError(
			"storage",
			fmt.Sprintf(
				`invalid value: cannot reduce storage from %d to %d`,
				existingStorage,
				newStorge,
			),
		)
	}
	return nil
}
