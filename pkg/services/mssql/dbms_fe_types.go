package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

// "Fe" means "From existing"
type dbmsFeProvisioningParams struct {
	// Original PP, for updating
	FirewallRules []firewallRule `json:"firewallRules"`
	// Migration PP, for provisioning
	AdministratorLogin string `json:"administratorLogin"`
	ServerName         string `json:"serverName"`
}

type secureDBMSFeProvisioningParams struct {
	// Migration SPP, for provisioning
	AdministratorLoginPassword string `json:"administratorLoginPassword"`
}

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
		},
	}
}

func (
	d *dbmsFeManager,
) getUpdatingParametersSchema() *service.InputParametersSchema {
	return &service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"firewallRules": &service.ArrayPropertySchema{
				Description: "Firewall rules to apply to instance. " +
					"If left unspecified, defaults to only Azure IPs",
				ItemsSchema: &service.ObjectPropertySchema{
					Description: "Individual Firewall Rule",
					RequiredProperties: []string{
						"name",
						"startIPAddress",
						"endIPAddress",
					},
					PropertySchemas: map[string]service.PropertySchema{
						"name": &service.StringPropertySchema{
							Description: "Name of firewall rule",
						},
						"startIPAddress": &service.StringPropertySchema{
							Description:             "Start of firewall rule range",
							CustomPropertyValidator: ipValidator,
						},
						"endIPAddress": &service.StringPropertySchema{
							Description:             "End of firewall rule range",
							CustomPropertyValidator: ipValidator,
						},
					},
					CustomPropertyValidator: firewallRuleValidator,
				},
				DefaultValue: []interface{}{
					map[string]interface{}{
						"name":           "AllowAzure",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "0.0.0.0",
					},
				},
			},
		},
	}
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

func (d *dbmsFeManager) SplitProvisioningParameters(
	cpp map[string]interface{},
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	pp := dbmsFeProvisioningParams{}
	if err := service.GetStructFromMap(cpp, &pp); err != nil {
		return nil, nil, err
	}
	spp := secureDBMSFeProvisioningParams{}
	if err := service.GetStructFromMap(cpp, &spp); err != nil {
		return nil, nil, err
	}
	ppMap, err := service.GetMapFromStruct(pp)
	sppMap, err := service.GetMapFromStruct(spp)
	return ppMap, sppMap, err
}

func (d *dbmsFeManager) SplitBindingParameters(
	params map[string]interface{},
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
