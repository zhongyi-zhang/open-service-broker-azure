package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (
	d *dbmsFeManager,
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
	d *dbmsFeManager,
) getUpdatingParametersSchema() service.InputParametersSchema {
	return getDBMSCommonProvisionParamSchema()
}

type dbmsFeInstanceDetails struct {
	ARMDeploymentName        string `json:"armDeployment"`
	FullyQualifiedDomainName string `json:"fullyQualifiedDomainName"`
	ServerName               string `json:"server"`
	AdministratorLogin       string `json:"administratorLogin"`
}

type secureDBMSFeInstanceDetails struct {
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}
