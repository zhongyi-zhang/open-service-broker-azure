package mssql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (d *databaseFeManager) GetEmptyInstanceDetails() service.InstanceDetails {
	return &databaseInstanceDetails{}
}

func (d *databaseFeManager) GetEmptyBindingDetails() service.BindingDetails {
	return &bindingDetails{}
}
