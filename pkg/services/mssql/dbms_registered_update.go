package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsRegisteredManager) ValidateUpdatingParameters(
	service.Instance,
) error {
	return nil
}

func (d *dbmsRegisteredManager) GetUpdater(
	service.Plan,
) (service.Updater, error) {
	return service.NewUpdater()
}
