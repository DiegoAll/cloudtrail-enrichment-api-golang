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

	// var enrichedRecordsResponse []models.EnrichedEventRecord
	enrichedRecordsResponse, err := ec.service.EnrichEvent(r.Context(), &eventInput) // Pasa el contexto del request y el EventInput completo
	if err != nil {
		logger.ErrorLog.Printf("Error al enriquecer eventos: %v", err)
		// Decide el código de estado adecuado. Podría ser 500 si es un error interno del servicio.
		utils.ErrorJSON(w, fmt.Errorf("error al procesar y enriquecer eventos: %w", err), http.StatusInternalServerError)
		return
	}

	payload := utils.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("%d registros procesados y enriquecidos exitosamente", len(enrichedRecordsResponse)),
		Data:    []interface{}{}, // Devolvemos los records enriquecidos que se guardaron
		// Data:    enrichedRecordsResponse,
	}

	if err := utils.WriteJSON(w, http.StatusCreated, payload); err != nil {
		logger.ErrorLog.Println("Error al escribir la respuesta JSON:", err)
	}

}

func (ec *EnrichmentController) QueryEvents(w http.ResponseWriter, r *http.Request) {

	// Validar que el método sea GET
	if r.Method != http.MethodGet {
		logger.ErrorLog.Printf("Método no permitido: %s", r.Method)
		utils.ErrorJSON(w, errors.New("método no permitido"), http.StatusMethodNotAllowed)
		return
	}

	// Llamar al servicio para obtener los últimos 10 logs
	records, err := ec.service.Top10QueryEvents(r.Context())
	if err != nil {
		logger.ErrorLog.Printf("Error en el controlador al consultar eventos: %v", err)
		utils.ErrorJSON(w, fmt.Errorf("error al obtener los últimos 10 eventos: %w", err), http.StatusInternalServerError)
		return
	}

	// Preparar la respuesta
	payload := utils.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Últimos %d eventos enriquecidos obtenidos exitosamente", len(records)),
		Data:    records, // Devolver los registros obtenidos
	}

	// Escribir la respuesta JSON
	if err := utils.WriteJSON(w, http.StatusOK, payload); err != nil {
		logger.ErrorLog.Println("Error al escribir la respuesta JSON:", err)
	}

}
