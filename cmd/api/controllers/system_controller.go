package controllers

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/utils"
	"net/http"
)

type SystemController struct {
}

func NewSystemController() *SystemController {
	return &SystemController{}
}

func (sc *SystemController) HealthCheck(w http.ResponseWriter, r *http.Request) {

	response := utils.JSONResponse{
		Error:   false,
		Message: "Service is healthy",
		Data: map[string]string{
			"status": "OK, server is running",
		},
	}

	utils.WriteJSON(w, http.StatusOK, response)
	logger.InfoLog.Println("Health check performed: API is up and running.")
}
