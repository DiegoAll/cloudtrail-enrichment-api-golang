package repository_test

import (
	"cloudtrail-enrichment-api-golang/database/mongo"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// Este archivo contiene las pruebas unitarias para EnrichmentMongoRepository.

func TestMain(m *testing.M) {
	logger.Init()
	m.Run()
}

// TestEnrichmentMongoRepository_InsertLog prueba el método InsertLog.
func TestEnrichmentMongoRepository_InsertLog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Run("insert_log_test", func(mt *mtest.T) {
		repo := mongo.NewEnrichMongoRepository(mt.Client, "testdb", "enriched_events")
		ctx := context.Background()

		eventToInsert := &models.EnrichedEventRecord{
			EventVersion: "1.08",
			EventSource:  "ec2.amazonaws.com",
			EventName:    "StartInstances",
			EventTime:    time.Now().UTC(),
			Enrichment:   models.EnrichmentData{Country: "USA"},
		}

		// --- Caso de éxito ---
		mt.Run("success", func(mt *mtest.T) {
			// Configura el mock solo para esta subprueba
			mt.AddMockResponses(mtest.CreateSuccessResponse())
			err := repo.InsertLog(ctx, eventToInsert)
			if err != nil {
				mt.Fatalf("InsertLog debería tener éxito, pero falló con: %v", err)
			}
		})

		// --- Caso de fallo en la inserción ---
		mt.Run("failure_on_insert", func(mt *mtest.T) {
			// Configura el mock solo para esta subprueba
			mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
				Index:   0,
				Code:    11000,
				Message: "error de clave duplicada simulado",
			}))

			err := repo.InsertLog(ctx, eventToInsert)
			if err == nil {
				mt.Fatal("Se esperaba un error de inserción, pero se obtuvo nil")
			}
			if err.Error() != "error al insertar evento enriquecido: error de clave duplicada simulado" {
				mt.Fatalf("El error obtenido no coincide con el esperado.\nObtenido: %v\nEsperado: %v", err, errors.New("error al insertar evento enriquecido: error de clave duplicada simulado"))
			}
		})
	})
}

// TestEnrichmentMongoRepository_GetLatestLogs prueba el método GetLatestLogs.
func TestEnrichmentMongoRepository_GetLatestLogs(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Run("get_latest_logs_test", func(mt *mtest.T) {
		repo := mongo.NewEnrichMongoRepository(mt.Client, "testdb", "enriched_events")
		ctx := context.Background()

		events := []*models.EnrichedEventRecord{
			{ID: primitive.NewObjectID(), EventName: "EventC", EventTime: time.Now().Add(time.Second * 3)},
			{ID: primitive.NewObjectID(), EventName: "EventB", EventTime: time.Now().Add(time.Second * 2)},
			{ID: primitive.NewObjectID(), EventName: "EventA", EventTime: time.Now().Add(time.Second * 1)},
		}

		// --- Caso de éxito ---
		mt.Run("success", func(mt *mtest.T) {
			// Configura el mock solo para esta subprueba
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "testdb.enriched_events", mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: events[0].ID},
					{Key: "eventName", Value: events[0].EventName},
					{Key: "eventTime", Value: events[0].EventTime},
				},
				bson.D{
					{Key: "_id", Value: events[1].ID},
					{Key: "eventName", Value: events[1].EventName},
					{Key: "eventTime", Value: events[1].EventTime},
				},
				bson.D{
					{Key: "_id", Value: events[2].ID},
					{Key: "eventName", Value: events[2].EventName},
					{Key: "eventTime", Value: events[2].EventTime},
				},
			))
			mt.AddMockResponses(mtest.CreateCursorResponse(0, "testdb.enriched_events", mtest.NextBatch))

			retrievedEvents, err := repo.GetLatestLogs(ctx)
			if err != nil {
				mt.Fatalf("GetLatestLogs debería tener éxito, pero falló con: %v", err)
			}

			if len(retrievedEvents) != 3 {
				mt.Fatalf("Se esperaba 3 eventos, se obtuvieron %d", len(retrievedEvents))
			}
			if retrievedEvents[0].EventName != "EventC" {
				mt.Fatalf("El primer evento no está ordenado correctamente")
			}
		})

		// --- Caso de fallo en la búsqueda ---
		mt.Run("failure_on_find", func(mt *mtest.T) {
			// Configura el mock solo para esta subprueba
			mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "error simulado de consulta"}})

			retrievedEvents, err := repo.GetLatestLogs(ctx)
			if err == nil {
				mt.Fatal("Se esperaba un error, pero se obtuvo nil")
			}
			if retrievedEvents != nil {
				mt.Fatal("Se esperaba una lista de eventos nula en caso de error")
			}
			if err.Error() != "error al obtener los últimos eventos enriquecidos: error simulado de consulta" {
				mt.Fatalf("El error obtenido no coincide con el esperado: %v", err)
			}
		})

		// --- Caso de error en el cursor ---
		mt.Run("failure_on_cursor", func(mt *mtest.T) {
			// Configura el mock solo para esta subprueba
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "testdb.enriched_events", mtest.FirstBatch,
				bson.D{{Key: "eventName", Value: "EventA"}},
			))
			mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}, {Key: "errmsg", Value: "error de cursor simulado"}})

			retrievedEvents, err := repo.GetLatestLogs(ctx)
			if err == nil {
				mt.Fatal("Se esperaba un error, pero se obtuvo nil")
			}
			if retrievedEvents != nil {
				mt.Fatal("Se esperaba una lista de eventos nula en caso de error")
			}
		})
	})
}
