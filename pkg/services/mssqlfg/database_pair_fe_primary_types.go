package mssqlfg

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (d *databasePairFePrimaryManager) GetEmptyInstanceDetails() service.InstanceDetails { // nolint: lll
	return &databasePairInstanceDetails{}
}

func (d *databasePairFePrimaryManager) GetEmptyBindingDetails() service.BindingDetails { // nolint: lll
	return &bindingDetails{}
}
