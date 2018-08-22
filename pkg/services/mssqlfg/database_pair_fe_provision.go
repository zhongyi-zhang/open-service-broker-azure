package mssqlfg

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (d *databasePairFeManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", d.preProvision),
		service.NewProvisioningStep("validatePriDatabase", d.validatePriDatabase),
		service.NewProvisioningStep("validateSecDatabase", d.validateSecDatabase),
		service.NewProvisioningStep(
			"validateFailoverGroup",
			d.validateFailoverGroup,
		),
		service.NewProvisioningStep(
			"deployPriFeARMTemplate",
			d.deployPriFeARMTemplate,
		),
		service.NewProvisioningStep(
			"deployFgFeARMTemplate",
			d.deployFgFeARMTemplate,
		),
		service.NewProvisioningStep(
			"deploySecFeARMTemplate",
			d.deploySecFeARMTemplate,
		),
	)
}

func (d *databasePairFeManager) preProvision(
	_ context.Context,
	_ service.Instance,
) (service.InstanceDetails, error) {
	return &databasePairInstanceDetails{
		PriARMDeploymentName: uuid.NewV4().String(),
		SecARMDeploymentName: uuid.NewV4().String(),
		FgARMDeploymentName:  uuid.NewV4().String(),
	}, nil
}
