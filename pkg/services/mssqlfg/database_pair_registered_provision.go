package mssqlfg

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *databasePairRegisteredManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("getPriDatabase", d.getPriDatabase),
		service.NewProvisioningStep("getSecDatabase", d.getSecDatabase),
		service.NewProvisioningStep("getFailoverGroup", d.getFailoverGroup),
	)
}

func (d *databasePairRegisteredManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	return &databasePairInstanceDetails{
		DatabaseName:      pp.GetString("database"),
		FailoverGroupName: pp.GetString("failoverGroup"),
	}, nil
}

func (d *databasePairRegisteredManager) getPriDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if err := getDatabase(
		ctx,
		&d.databasesClient,
		ppp.GetString("primaryResourceGroup"),
		pdt.PriServerName,
		pp.GetString("database"),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (d *databasePairRegisteredManager) getSecDatabase(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if err := getDatabase(
		ctx,
		&d.databasesClient,
		ppp.GetString("secondaryResourceGroup"),
		pdt.SecServerName,
		pp.GetString("database"),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (d *databasePairRegisteredManager) getFailoverGroup(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pp := instance.ProvisioningParameters
	ppp := instance.Parent.ProvisioningParameters
	pdt := instance.Parent.Details.(*dbmsPairInstanceDetails)
	if err := getFailoverGroup(
		ctx,
		&d.failoverGroupsClient,
		ppp.GetString("primaryResourceGroup"),
		pdt.PriServerName,
		pp.GetString("failoverGroup"),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}
