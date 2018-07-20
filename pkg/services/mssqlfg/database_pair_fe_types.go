package mssqlfg

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (d *databasePairFeManager) GetEmptyInstanceDetails() service.InstanceDetails { // nolint: lll
	return &databasePairInstanceDetails{}
}

func (d *databasePairFeManager) GetEmptyBindingDetails() service.BindingDetails { // nolint: lll
	return &bindingDetails{}
}
