package v1

import (
	"github.com/milyrock/Surf/internal/models"
)

type Service interface {
	CreateRecord(rec *models.Record) (*models.Record, error)
}


type APIConfig struct {
	Service Service
}

type API struct {
	service Service
}

func NewAPI(config APIConfig) *API {
	return &API{service: config.Service}
}
