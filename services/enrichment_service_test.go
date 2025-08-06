package services

import (
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEnrichmentRepository es un mock para la interfaz EnrichmentRepository.
type MockEnrichmentRepository struct {
	mock.Mock
}

func (m *MockEnrichmentRepository) InsertLog(ctx context.Context, event *models.EnrichedEventRecord) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEnrichmentRepository) GetLatestLogs(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.EnrichedEventRecord), args.Error(1)
}

// TestEnrichEvent_Success prueba el enriquecimiento exitoso de un evento.
func TestEnrichEvent_Success(t *testing.T) {
	mockRepo := new(MockEnrichmentRepository)
	mockRepo.On("InsertLog", mock.Anything, mock.Anything).Return(nil)

	service := NewDefaultEnrichmentService(mockRepo)

	// Creamos la variable de tipo models.Event y le asignamos directamente los datos.
	event := &models.Event{
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
		}{
			{
				SourceIPAddress: "8.8.8.8",
				EventTime:       time.Now(),
			},
		},
	}

	_, err := service.EnrichEvent(context.Background(), event)

	assert.Nil(t, err, "Se esperaba que no hubiera error")
	mockRepo.AssertCalled(t, "InsertLog", mock.Anything, mock.AnythingOfType("*models.EnrichedEventRecord"))
}

// TestEnrichEvent_InsertLogError prueba un error al insertar en la base de datos.
func TestEnrichEvent_InsertLogError(t *testing.T) {
	mockRepo := new(MockEnrichmentRepository)
	mockRepo.On("InsertLog", mock.Anything, mock.Anything).Return(errors.New("error de base de datos"))

	service := NewDefaultEnrichmentService(mockRepo)

	event := &models.Event{
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
		}{
			{
				SourceIPAddress: "8.8.8.8",
				EventTime:       time.Now(),
			},
		},
	}

	_, err := service.EnrichEvent(context.Background(), event)

	assert.NotNil(t, err, "Se esperaba un error al insertar en la base de datos")
	assert.Contains(t, err.Error(), "error al insertar evento enriquecido", "El mensaje de error debe ser específico")

	mockRepo.AssertCalled(t, "InsertLog", mock.Anything, mock.Anything)
}

// --- Pruebas para el método Top10QueryEvents ---

func TestTop10QueryEvents_Success(t *testing.T) {
	mockRepo := new(MockEnrichmentRepository)
	expectedRecords := []*models.EnrichedEventRecord{
		{SourceIPAddress: "1.1.1.1"},
		{SourceIPAddress: "2.2.2.2"},
	}
	mockRepo.On("GetLatestLogs", mock.Anything).Return(expectedRecords, nil)

	service := NewDefaultEnrichmentService(mockRepo)
	records, err := service.Top10QueryEvents(context.Background())

	assert.Nil(t, err, "Se esperaba que no hubiera error")
	assert.Equal(t, 2, len(records), "Se esperaban dos registros")
	assert.Equal(t, expectedRecords, records, "Los registros obtenidos deben ser los mismos que los esperados")
	mockRepo.AssertCalled(t, "GetLatestLogs", mock.Anything)
}

func TestTop10QueryEvents_RepoError(t *testing.T) {
	mockRepo := new(MockEnrichmentRepository)
	mockRepo.On("GetLatestLogs", mock.Anything).Return(([]*models.EnrichedEventRecord)(nil), errors.New("error de DB"))

	service := NewDefaultEnrichmentService(mockRepo)
	_, err := service.Top10QueryEvents(context.Background())

	assert.NotNil(t, err, "Se esperaba un error al consultar")
	assert.Contains(t, err.Error(), "error al obtener los últimos 10 eventos del repositorio", "El mensaje de error debe ser específico")
	mockRepo.AssertCalled(t, "GetLatestLogs", mock.Anything)
}
