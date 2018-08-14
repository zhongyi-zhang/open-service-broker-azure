package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (
	d *dbmsRegisteredManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return getDBMSRegisteredProvisionParamSchema()
}

func (
	d *dbmsRegisteredManager,
) getUpdatingParametersSchema() service.InputParametersSchema {
	return getDBMSRegisteredUpdateParamSchema()
}

func (d *dbmsRegisteredManager) GetEmptyInstanceDetails() service.InstanceDetails { // nolint: lll
	return &dbmsInstanceDetails{}
}

func (d *dbmsRegisteredManager) GetEmptyBindingDetails() service.BindingDetails { // nolint: lll
	return nil
}
