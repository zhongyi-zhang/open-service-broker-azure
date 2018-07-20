package mssqlfg

import "github.com/Azure/open-service-broker-azure/pkg/service"

type databasePairInstanceDetails struct {
	PriARMDeploymentName string `json:"primaryArmDeployment"`
	SecARMDeploymentName string `json:"secondaryArmDeployment"`
	FgARMDeploymentName  string `json:"failoverGroupArmDeployment"`
	FailoverGroupName    string `json:"failoverGroup"`
	DatabaseName         string `json:"database"`
}

func (d *databasePairManager) GetEmptyInstanceDetails() service.InstanceDetails { // nolint: lll
	return &databasePairInstanceDetails{}
}

func (d *databasePairManager) GetEmptyBindingDetails() service.BindingDetails { // nolint: lll
	return &bindingDetails{}
}
