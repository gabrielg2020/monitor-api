package handlers

import (
	"strings"

	"github.com/gabrielg2020/monitor-api/internal/entities"
	"github.com/gabrielg2020/monitor-api/internal/models"
)

// setMetricQueryDefaults validates and sets defaults for metric query params
func setMetricQueryDefaults(params *entities.MetricQueryParams) *models.ErrorResponse {
	// Set defaults
	if params.Limit == 0 {
		params.Limit = 100
	}
	if params.Limit > 1000 {
		params.Limit = 1000
	}
	if params.Order == "" {
		params.Order = "DESC"
	} else {
		params.Order = strings.ToUpper(params.Order)
		if params.Order != "ASC" && params.Order != "DESC" {
			return &models.ErrorResponse{
				Error:   "Invalid order parameter",
				Details: "Must be 'ASC' or 'DESC'",
			}
		}
	}

	return nil
}
