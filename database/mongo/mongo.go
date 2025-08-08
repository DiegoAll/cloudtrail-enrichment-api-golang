package mongo

import (
	"cloudtrail-enrichment-api-golang/internal/config"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoInstance encapsulates the MongoDB connection and the specific collection.
// *EnrichMongoRepository  (Alternative)
type MongoInstance struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func BuildMongoURI(cfg *config.MongoDBConfig) string {
	var mongoURI string
	if cfg.Username != "" && cfg.Password != "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	} else {
		// If no credentials, only host, port, and database.
		mongoURI = fmt.Sprintf("mongodb://%s:%d/%s", cfg.Host, cfg.Port, cfg.Database)
	}
	logger.InfoLog.Printf("Building MONGO_URI from configurations: %s", mongoURI)
	return mongoURI
}

func NewMongoClient(mongoURI string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.InfoLog.Printf("[DEBUG] Final MongoDB URI used for connection: %s", mongoURI)

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.ErrorLog.Printf("Error connecting to MongoDB: %v", err)
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.ErrorLog.Printf("Error pinging MongoDB: %v", err)
		// Disconnect the client if the ping fails to free resources.
		if disconnectErr := client.Disconnect(context.Background()); disconnectErr != nil {
			logger.ErrorLog.Printf("Error disconnecting MongoDB client after ping failure: %v", disconnectErr)
		}
		return nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	logger.InfoLog.Println("MongoDB connection established successfully.")
	return client, nil
}

// EnrichmentMongoRepository implements the EnrichmentRepository interface for MongoDB.
type EnrichmentMongoRepository struct {
	mongoInstance *MongoInstance
}

// NewEnrichMongoRepository creates a new instance of EnrichmentMongoRepository.
func NewEnrichMongoRepository(client *mongo.Client, dbName, collectionName string) *EnrichmentMongoRepository {
	logger.InfoLog.Printf("[DEBUG] Connecting to MongoDB. Database: '%s', Collection: '%s'", dbName, collectionName)
	collection := client.Database(dbName).Collection(collectionName)

	return &EnrichmentMongoRepository{
		mongoInstance: &MongoInstance{
			Client:     client,
			Collection: collection,
		},
	}
}

// InsertLog inserts a new enriched event record into MongoDB.
func (m *EnrichmentMongoRepository) InsertLog(ctx context.Context, event *models.EnrichedEventRecord) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // Use context timeout, or define one if nil
	defer cancel()

	_, err := m.mongoInstance.Collection.InsertOne(ctx, event)
	if err != nil {
		logger.ErrorLog.Printf("Error inserting enriched event into MongoDB: %v", err)
		return fmt.Errorf("error inserting enriched event: %w", err)
	}
	logger.InfoLog.Printf("Enriched event inserted into MongoDB. EventSource: %s", event.EventSource)
	return nil
}

// GetLatestLogs retrieves the last 10 enriched event records from MongoDB.
func (m *EnrichmentMongoRepository) GetLatestLogs(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // Use context timeout, or define one if nil
	defer cancel()

	var events []*models.EnrichedEventRecord
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "eventTime", Value: -1}}) // Sort by eventTime descending (-1 most recent)
	findOptions.SetLimit(10)                                   // Limit to last 10

	cursor, err := m.mongoInstance.Collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		logger.ErrorLog.Printf("Error getting latest enriched events from MongoDB: %v", err)
		return nil, fmt.Errorf("error getting latest enriched events: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event models.EnrichedEventRecord
		if err := cursor.Decode(&event); err != nil {
			logger.ErrorLog.Printf("Error decoding MongoDB event: %v", err)
			return nil, fmt.Errorf("error decoding event: %w", err)
		}
		events = append(events, &event)
	}

	if err := cursor.Err(); err != nil {
		logger.ErrorLog.Printf("Error in MongoDB cursor: %v", err)
		return nil, fmt.Errorf("error in MongoDB cursor: %w", err)
	}

	logger.InfoLog.Println("Last 10 enriched events retrieved from MongoDB.")
	return events, nil
}
