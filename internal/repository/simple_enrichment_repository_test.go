package repository_test

import (
	"cloudtrail-enrichment-api-golang/database/mongo"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/repository"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// MockEnrichmentRepository implementa EnrichmentRepository para pruebas
type MockEnrichmentRepository struct {
	ShouldError bool
	Events      []*models.EnrichedEventRecord
}

func (m *MockEnrichmentRepository) InsertLog(ctx context.Context, event *models.EnrichedEventRecord) error {
	if m.ShouldError {
		return errors.New("mock error")
	}
	m.Events = append(m.Events, event)
	return nil
}

func (m *MockEnrichmentRepository) GetLatestLogs(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
	if m.ShouldError {
		return nil, errors.New("mock error")
	}
	return m.Events, nil
}

// TestMain inicializa el logger antes de ejecutar las pruebas
func TestMain(m *testing.M) {
	logger.Init()
	m.Run()
}

// TestSimple_SetEnrichmentRepository prueba la función SetEnrichmentRepository
func TestSimple_SetEnrichmentRepository(t *testing.T) {
	// Crear un mock repository
	mockRepo := &MockEnrichmentRepository{}

	// Llamar a SetEnrichmentRepository
	repository.SetEnrichmentRepository(mockRepo)

	// Verificar que se estableció correctamente probando las funciones auxiliares
	ctx := context.Background()
	event := &models.EnrichedEventRecord{
		EventVersion: "1.08",
		EventSource:  "test.amazonaws.com",
		EventName:    "TestEvent",
		EventTime:    time.Now().UTC(),
		Enrichment:   models.EnrichmentData{Country: "Test"},
	}

	// Probar InsertLog
	err := repository.InsertLog(ctx, event)
	if err != nil {
		t.Fatalf("InsertLog falló: %v", err)
	}

	if len(mockRepo.Events) != 1 {
		t.Fatalf("Esperaba 1 evento, obtuvo %d", len(mockRepo.Events))
	}

	// Probar GetLatestLogs
	events, err := repository.GetLatestLogs(ctx)
	if err != nil {
		t.Fatalf("GetLatestLogs falló: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Esperaba 1 evento, obtuvo %d", len(events))
	}
}

// TestSimple_InsertLogFunction prueba la función auxiliar InsertLog
func TestSimple_InsertLogFunction(t *testing.T) {
	tests := []struct {
		name        string
		shouldError bool
		wantError   bool
	}{
		{
			name:        "éxito",
			shouldError: false,
			wantError:   false,
		},
		{
			name:        "error",
			shouldError: true,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockEnrichmentRepository{ShouldError: tt.shouldError}
			repository.SetEnrichmentRepository(mockRepo)

			ctx := context.Background()
			event := &models.EnrichedEventRecord{
				EventVersion: "1.08",
				EventSource:  "test.amazonaws.com",
			}

			err := repository.InsertLog(ctx, event)
			if (err != nil) != tt.wantError {
				t.Fatalf("InsertLog() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestSimple_GetLatestLogsFunction prueba la función auxiliar GetLatestLogs
func TestSimple_GetLatestLogsFunction(t *testing.T) {
	tests := []struct {
		name        string
		shouldError bool
		wantError   bool
		eventCount  int
	}{
		{
			name:        "éxito_con_eventos",
			shouldError: false,
			wantError:   false,
			eventCount:  2,
		},
		{
			name:        "éxito_sin_eventos",
			shouldError: false,
			wantError:   false,
			eventCount:  0,
		},
		{
			name:        "error",
			shouldError: true,
			wantError:   true,
			eventCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockEnrichmentRepository{ShouldError: tt.shouldError}

			// Preparar eventos de prueba si es necesario
			if tt.eventCount > 0 {
				for i := 0; i < tt.eventCount; i++ {
					mockRepo.Events = append(mockRepo.Events, &models.EnrichedEventRecord{
						EventVersion: "1.08",
						EventSource:  "test.amazonaws.com",
						EventName:    "TestEvent",
					})
				}
			}

			repository.SetEnrichmentRepository(mockRepo)

			ctx := context.Background()
			events, err := repository.GetLatestLogs(ctx)

			if (err != nil) != tt.wantError {
				t.Fatalf("GetLatestLogs() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && len(events) != tt.eventCount {
				t.Fatalf("GetLatestLogs() devolvió %d eventos, esperaba %d", len(events), tt.eventCount)
			}
		})
	}
}

// TestSimple_MongoRepository - Pruebas simples para el repositorio MongoDB
func TestSimple_MongoRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("insert_simple", func(mt *mtest.T) {
		repo := mongo.NewEnrichMongoRepository(mt.Client, "testdb", "enriched_events")
		ctx := context.Background()

		event := &models.EnrichedEventRecord{
			EventVersion: "1.08",
			EventSource:  "test.amazonaws.com",
			EventName:    "TestEvent",
			EventTime:    time.Now().UTC(),
			Enrichment:   models.EnrichmentData{Country: "Test"},
		}

		// Mock simple para InsertOne - éxito
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.InsertLog(ctx, event)
		if err != nil {
			t.Fatalf("InsertLog falló: %v", err)
		}
	})

	mt.Run("insert_error", func(mt *mtest.T) {
		repo := mongo.NewEnrichMongoRepository(mt.Client, "testdb", "enriched_events")
		ctx := context.Background()

		event := &models.EnrichedEventRecord{
			EventVersion: "1.08",
			EventSource:  "test.amazonaws.com",
		}

		// Mock que simula error de inserción
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   0,
			Code:    11000,
			Message: "error de inserción",
		}))

		err := repo.InsertLog(ctx, event)
		if err == nil {
			t.Fatal("Esperaba error de inserción")
		}
	})

	mt.Run("get_logs_empty", func(mt *mtest.T) {
		repo := mongo.NewEnrichMongoRepository(mt.Client, "testdb", "enriched_events")
		ctx := context.Background()

		// Mock para cursor vacío
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "testdb.enriched_events", mtest.FirstBatch))

		events, err := repo.GetLatestLogs(ctx)
		if err != nil {
			t.Fatalf("GetLatestLogs falló: %v", err)
		}

		if len(events) != 0 {
			t.Fatalf("Esperaba 0 eventos, obtuvo %d", len(events))
		}
	})

	mt.Run("get_logs_with_data", func(mt *mtest.T) {
		repo := mongo.NewEnrichMongoRepository(mt.Client, "testdb", "enriched_events")
		ctx := context.Background()

		now := time.Now().UTC()
		eventID := primitive.NewObjectID()

		// Mock para cursor con un documento
		first := mtest.CreateCursorResponse(1, "testdb.enriched_events", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: eventID},
				{Key: "eventVersion", Value: "1.08"},
				{Key: "eventSource", Value: "ec2.amazonaws.com"},
				{Key: "eventName", Value: "StartInstances"},
				{Key: "eventTime", Value: now},
				{Key: "enrichment", Value: bson.D{
					{Key: "country", Value: "USA"},
				}},
			},
		)

		// Mock para finalizar cursor
		getMore := mtest.CreateCursorResponse(0, "testdb.enriched_events", mtest.NextBatch)

		mt.AddMockResponses(first, getMore)

		events, err := repo.GetLatestLogs(ctx)
		if err != nil {
			t.Fatalf("GetLatestLogs falló: %v", err)
		}

		if len(events) != 1 {
			t.Fatalf("Esperaba 1 evento, obtuvo %d", len(events))
		}

		if events[0].EventName != "StartInstances" {
			t.Fatalf("Nombre de evento incorrecto: %s", events[0].EventName)
		}
	})

	mt.Run("get_logs_find_error", func(mt *mtest.T) {
		repo := mongo.NewEnrichMongoRepository(mt.Client, "testdb", "enriched_events")
		ctx := context.Background()

		// Mock que simula error en Find
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "error de consulta",
			Name:    "InternalError",
		}))

		events, err := repo.GetLatestLogs(ctx)
		if err == nil {
			t.Fatal("Esperaba error de consulta")
		}

		if events != nil {
			t.Fatal("Los eventos deberían ser nil en caso de error")
		}
	})
}

// TestSimple_NewEnrichMongoRepository prueba la creación del repositorio
func TestSimple_NewEnrichMongoRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("creation", func(mt *mtest.T) {
		repo := mongo.NewEnrichMongoRepository(mt.Client, "testdb", "testcollection")
		if repo == nil {
			t.Fatal("El repositorio no debería ser nil")
		}

		// Verificar que podemos usar el repositorio
		ctx := context.Background()
		event := &models.EnrichedEventRecord{
			EventVersion: "1.08",
			EventSource:  "test.amazonaws.com",
		}

		// Mock simple
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := repo.InsertLog(ctx, event)
		if err != nil {
			t.Fatalf("El repositorio debería funcionar después de la creación: %v", err)
		}
	})
}
