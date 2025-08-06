// cmd/api/controllers/enrichment_controller_test.go

package controllers

import (
	"bytes"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestIngestData_Success prueba la ingesta de datos exitosa.
func TestIngestData_Success(t *testing.T) {
	mockService := &MockEnrichmentService{
		EnrichEventFunc: func(ctx context.Context, event *models.Event) ([]models.EnrichedEventRecord, error) {
			// Simula una respuesta exitosa
			return []models.EnrichedEventRecord{{}}, nil
		},
	}
	controller := NewEnrichmentController(mockService)

	payload := models.Event{
		Records: []models.EventRecord{{}}, // Payload válido
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/enrichment", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.IngestData(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusCreated, rr.Code)
	}
}

// TestIngestData_InvalidPayload prueba un payload de ingesta de datos inválido.
func TestIngestData_InvalidPayload(t *testing.T) {
	mockService := &MockEnrichmentService{}
	controller := NewEnrichmentController(mockService)

	req := httptest.NewRequest("POST", "/enrichment", bytes.NewReader([]byte("not a json")))
	rr := httptest.NewRecorder()

	controller.IngestData(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusBadRequest, rr.Code)
	}
}

// TestQueryEvents_Success prueba la consulta de eventos exitosa.
func TestQueryEvents_Success(t *testing.T) {
	mockService := &MockEnrichmentService{
		Top10QueryEventsFunc: func(ctx context.Context) ([]models.EnrichedEventRecord, error) {
			// Simula una respuesta exitosa
			return []models.EnrichedEventRecord{{ID: "1"}, {ID: "2"}}, nil
		},
	}
	controller := NewEnrichmentController(mockService)

	req := httptest.NewRequest("GET", "/enrichment", nil)
	rr := httptest.NewRecorder()

	controller.QueryEvents(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusOK, rr.Code)
	}
}

// TestQueryEvents_ServiceError prueba un error en el servicio al consultar eventos.
func TestQueryEvents_ServiceError(t *testing.T) {
	mockService := &MockEnrichmentService{
		Top10QueryEventsFunc: func(ctx context.Context) ([]models.EnrichedEventRecord, error) {
			return nil, errors.New("error de base de datos")
		},
	}
	controller := NewEnrichmentController(mockService)

	req := httptest.NewRequest("GET", "/enrichment", nil)
	rr := httptest.NewRecorder()

	controller.QueryEvents(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusInternalServerError, rr.Code)
	}
}
