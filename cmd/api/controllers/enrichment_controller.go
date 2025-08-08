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
	var eventInput models.Event

	err := utils.ReadJSON(w, r, &eventInput)
	if err != nil {
		logger.ErrorLog.Println("Error reading input JSON:", err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if len(eventInput.Records) == 0 {
		logger.ErrorLog.Println("The JSON does not contain any records in 'Records'")
		utils.ErrorJSON(w, errors.New("the JSON does not contain records in 'Records'"), http.StatusBadRequest)
		return
	}

	enrichedRecordsResponse, err := ec.service.EnrichEvent(r.Context(), &eventInput)
	if err != nil {
		logger.ErrorLog.Printf("Error enriching events: %v", err)
		// Decide on the appropriate status code. It could be 500 if it's an internal service error.
		utils.ErrorJSON(w, fmt.Errorf("error processing and enriching events: %w", err), http.StatusInternalServerError)
		return
	}

	payload := utils.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("%d records successfully processed and enriched", len(enrichedRecordsResponse)),
		Data:    []interface{}{}, // We return the enriched records that were saved
		// Data:    enrichedRecordsResponse,
	}

	if err := utils.WriteJSON(w, http.StatusCreated, payload); err != nil {
		logger.ErrorLog.Println("Error writing JSON response:", err)
	}

}

func (ec *EnrichmentController) QueryEvents(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		logger.ErrorLog.Printf("Method not allowed: %s", r.Method)
		utils.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	// Call the service to get the last 10 logs
	records, err := ec.service.Top10QueryEvents(r.Context())
	if err != nil {
		logger.ErrorLog.Printf("Error in controller when querying events: %v", err)
		utils.ErrorJSON(w, fmt.Errorf("error getting the last 10 events: %w", err), http.StatusInternalServerError)
		return
	}

	payload := utils.JSONResponse{
		Error:   false,
		Message: fmt.Sprintf("Last %d enriched events successfully retrieved", len(records)),
		Data:    records, // Return the retrieved records
	}

	if err := utils.WriteJSON(w, http.StatusOK, payload); err != nil {
		logger.ErrorLog.Println("Error writing JSON response:", err)
	}

}
