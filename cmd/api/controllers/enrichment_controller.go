package controllers

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/utils"
	"cloudtrail-enrichment-api-golang/models"
	"cloudtrail-enrichment-api-golang/services"
	"errors"
	"fmt"
	"net/http"
)

type EnrichmentController struct {
	service services.EnrichmentService
}

func NewEnrichmentController(service services.EnrichmentService) *EnrichmentController {
	return &EnrichmentController{
		service: service,
	}
}

func (ec *EnrichmentController) IngestData(w http.ResponseWriter, r *http.Request) {
	var eventInput models.Event // Usamos la estructura Event original para la entrada

	err := utils.ReadJSON(w, r, &eventInput)
	if err != nil {
		logger.ErrorLog.Println("Error al leer JSON de entrada:", err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if len(eventInput.Records) == 0 {
		logger.ErrorLog.Println("El JSON no contiene ningún registro en 'Records'")
		utils.ErrorJSON(w, errors.New("el JSON no contiene registros en 'Records'"), http.StatusBadRequest)
		return
	}

	var enrichedRecordsResponse []models.EnrichedEventRecord

	//TODO
	// CALL SERVICE EnrichEvent

	payload := utils.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("%d registros procesados y enriquecidos exitosamente", len(enrichedRecordsResponse)),
		Data:    []interface{}{}, // Devolvemos los records enriquecidos que se guardaron
	}

	if err := utils.WriteJSON(w, http.StatusCreated, payload); err != nil {
		logger.ErrorLog.Println("Error al escribir la respuesta JSON:", err)
	}

}

func (ec *EnrichmentController) QueryEvents(w http.ResponseWriter, r *http.Request) {

}
