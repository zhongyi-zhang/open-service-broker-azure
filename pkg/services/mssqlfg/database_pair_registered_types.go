package mssqlfg

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (d *databasePairRegisteredManager) GetEmptyInstanceDetails() service.InstanceDetails { // nolint: lll
	return &databasePairInstanceDetails{}
}

func (d *databasePairRegisteredManager) GetEmptyBindingDetails() service.BindingDetails { // nolint: lll
	return &bindingDetails{}
}
