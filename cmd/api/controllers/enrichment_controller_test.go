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
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestIngestData_Success prueba la ingesta de datos exitosa.
func TestIngestData_Success(t *testing.T) {
	mockService := &MockEnrichmentService{
		EnrichEventFunc: func(ctx context.Context, event *models.Event) ([]*models.EnrichedEventRecord, error) {
			return []*models.EnrichedEventRecord{{}}, nil
		},
	}
	controller := NewEnrichmentController(mockService)

	// CORRECCIÓN: Se añaden las etiquetas JSON a la definición del struct anónimo
	// para que coincida con el tipo definido en models/enrichment.go.
	payload := models.Event{
		Records: []struct {
			EventVersion      string                   `json:"eventVersion"`
			UserIdentity      models.UserIdentity      `json:"userIdentity"`
			EventTime         time.Time                `json:"eventTime"`
			EventSource       string                   `json:"eventSource"`
			EventName         string                   `json:"eventName"`
			AwsRegion         string                   `json:"awsRegion"`
			SourceIPAddress   string                   `json:"sourceIPAddress"`
			UserAgent         string                   `json:"userAgent"`
			RequestParameters models.RequestParameters `json:"requestParameters"`
			ResponseElements  models.ResponseElements  `json:"responseElements"`
			Enrichment        models.EnrichmentData    `json:"enrichment"`
		}{{}},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/enrichment", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.IngestData(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("código de estado esperado %d, pero se obtuvo %d", http.StatusCreated, rr.Code)
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
		t.Errorf("código de estado esperado %d, pero se obtuvo %d", http.StatusBadRequest, rr.Code)
	}
}

// TestQueryEvents_Success prueba la consulta de eventos exitosa.
func TestQueryEvents_Success(t *testing.T) {
	mockService := &MockEnrichmentService{
		Top10QueryEventsFunc: func(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
			return []*models.EnrichedEventRecord{{ID: primitive.NewObjectID()}, {ID: primitive.NewObjectID()}}, nil
		},
	}
	controller := NewEnrichmentController(mockService)

	req := httptest.NewRequest("GET", "/enrichment", nil)
	rr := httptest.NewRecorder()

	controller.QueryEvents(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("código de estado esperado %d, pero se obtuvo %d", http.StatusOK, rr.Code)
	}
}

// TestQueryEvents_ServiceError prueba un error en el servicio al consultar eventos.
func TestQueryEvents_ServiceError(t *testing.T) {
	mockService := &MockEnrichmentService{
		Top10QueryEventsFunc: func(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
			return nil, errors.New("error de base de datos")
		},
	}
	controller := NewEnrichmentController(mockService)

	req := httptest.NewRequest("GET", "/enrichment", nil)
	rr := httptest.NewRecorder()

	controller.QueryEvents(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("código de estado esperado %d, pero se obtuvo %d", http.StatusInternalServerError, rr.Code)
	}
}
