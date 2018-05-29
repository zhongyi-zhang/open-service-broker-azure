package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (d *dbmsFeManager) ValidateUpdatingParameters(service.Instance) error {
	return nil
}

func (d *dbmsFeManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater()
}
