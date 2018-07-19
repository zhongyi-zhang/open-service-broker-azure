package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (
	d *dbmsRegisteredManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{
			"resourceGroup",
			"location",
			"server",
			"administratorLogin",
			"administratorLoginPassword",
		},
		SecureProperties: []string{
			"administratorLoginPassword",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"resourceGroup": &service.StringPropertySchema{
				Description: "Specifies the resource group of the existing server",
			},
			"location": &service.StringPropertySchema{
				Description: "Specifies the location of the existing server",
			},
			"server": &service.StringPropertySchema{
				Description: "Specifies the name of the existing server",
			},
			"administratorLogin": &service.StringPropertySchema{
				Description: "Specifies the administrator login name" +
					" of the existing server",
			},
			"administratorLoginPassword": &service.StringPropertySchema{
				Description: "Specifies the administrator login password" +
					" of the existing server",
			},
			"tags": &service.ObjectPropertySchema{
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}

func (
	d *dbmsRegisteredManager,
) getUpdatingParametersSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		SecureProperties: []string{
			"administratorLoginPassword",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"administratorLogin": &service.StringPropertySchema{
				Description: "Specifies the administrator login name" +
					" of the existing server",
			},
			"administratorLoginPassword": &service.StringPropertySchema{
				Description: "Specifies the administrator login password" +
					" of the existing server",
			},
		},
	}
}

func (d *dbmsRegisteredManager) GetEmptyInstanceDetails() service.InstanceDetails { // nolint: lll
	return &dbmsInstanceDetails{}
}

func (d *dbmsRegisteredManager) GetEmptyBindingDetails() service.BindingDetails { // nolint: lll
	return nil
}
