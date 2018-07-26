package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsRegisteredManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner()
}
