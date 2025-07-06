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

// MongoInstance encapsula la conexión a MongoDB y la colección específica.
// *EnrichMongoRepository  (Alternativa)
type MongoInstance struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func BuildMongoURI(cfg *config.MongoDBConfig) string {
	var mongoURI string
	if cfg.Username != "" && cfg.Password != "" {
		// Incluye el nombre de la base de datos en la URI si se usan credenciales,
		// y opcionalmente authSource si la base de datos de autenticación es diferente.
		// Asumimos authSource=admin por defecto si las credenciales son proporcionadas.
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	} else {
		// Si no hay credenciales, solo host, puerto y base de datos.
		mongoURI = fmt.Sprintf("mongodb://%s:%d/%s", cfg.Host, cfg.Port, cfg.Database)
	}
	logger.InfoLog.Printf("Construyendo MONGO_URI a partir de configuraciones: %s", mongoURI)
	return mongoURI
}

// NewMongoClient establece una nueva conexión a MongoDB y devuelve una instancia de MongoInstance.
// func NewMongoClient(cfg *config.Config) (*MongoInstance, error) {
// 	// Usar la configuración específica de MongoDB
// 	mongoCfg := cfg.MongoDBConfig

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoCfg.DBTimeout)
// 	defer cancel()

// 	// Intentar obtener la MONGO_URI completa de las variables de entorno primero
// 	mongoURI := os.Getenv("MONGO_URI")
// 	if mongoURI == "" {
// 		// Si MONGO_URI no está definida, construirla desde los parámetros individuales
// 		mongoURI = fmt.Sprintf("mongodb://%s:%d", mongoCfg.Host, mongoCfg.Port)
// 		if mongoCfg.Username != "" && mongoCfg.Password != "" {
// 			mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d", mongoCfg.Username, mongoCfg.Password, mongoCfg.Host, mongoCfg.Port)
// 			// Se recomienda especificar authSource si usas credenciales de usuario/contraseña
// 			mongoURI = fmt.Sprintf("%s/%s?authSource=admin", mongoURI, mongoCfg.Database)
// 		} else {
// 			// Si no hay credenciales, solo añadir la base de datos a la URI
// 			mongoURI = fmt.Sprintf("%s/%s", mongoURI, mongoCfg.Database)
// 		}
// 		logger.InfoLog.Println("Construyendo MONGO_URI a partir de configuraciones individuales.")
// 	} else {
// 		logger.InfoLog.Println("Usando MONGO_URI desde variables de entorno para la conexión a MongoDB.")
// 	}

// 	clientOptions := options.Client().ApplyURI(mongoURI)
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		logger.ErrorLog.Printf("Error al conectar a MongoDB: %v", err)
// 		return nil, fmt.Errorf("error al conectar a MongoDB: %w", err)
// 	}

// 	// Haz ping a la base de datos para verificar la conexión
// 	err = client.Ping(ctx, readpref.Primary())
// 	if err != nil {
// 		logger.ErrorLog.Printf("Error al hacer ping a MongoDB: %v", err)
// 		return nil, fmt.Errorf("error al hacer ping a MongoDB: %w", err)
// 	}

// 	logger.InfoLog.Println("Conexión a MongoDB establecida exitosamente.")

// 	// Aquí asumimos una base de datos y colección específicas para los eventos de enriquecimiento
// 	collection := client.Database(cfg.MongoDBConfig.Database).Collection("enriched_events")

// 	return &MongoInstance{Client: client, Collection: collection}, nil
// }

func NewMongoClient(mongoURI string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.ErrorLog.Printf("Error al conectar a MongoDB: %v", err)
		return nil, fmt.Errorf("error al conectar a MongoDB: %w", err)
	}

	// Haz ping a la base de datos para verificar la conexión
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.ErrorLog.Printf("Error al hacer ping a MongoDB: %v", err)
		// Desconectar el cliente si el ping falla para liberar recursos.
		if disconnectErr := client.Disconnect(context.Background()); disconnectErr != nil {
			logger.ErrorLog.Printf("Error al desconectar el cliente de MongoDB después de un fallo en el ping: %v", disconnectErr)
		}
		return nil, fmt.Errorf("error al hacer ping a MongoDB: %w", err)
	}

	logger.InfoLog.Println("Conexión a MongoDB establecida exitosamente.")
	return client, nil
}

// EnrichmentMongoRepository implementa la interfaz EnrichmentRepository para MongoDB.
type EnrichmentMongoRepository struct {
	mongoInstance *MongoInstance
}

// NewEnrichMongoRepository crea una nueva instancia de EnrichmentMongoRepository.
// func NewEnrichMongoRepository(mi *MongoInstance) *EnrichmentMongoRepository {
// 	return &EnrichmentMongoRepository{mongoInstance: mi}
// }

func NewEnrichMongoRepository(client *mongo.Client, dbName, collectionName string) *EnrichmentMongoRepository {
	collection := client.Database(dbName).Collection(collectionName)

	return &EnrichmentMongoRepository{
		mongoInstance: &MongoInstance{
			Client:     client,
			Collection: collection,
		},
	}
}

// InsertLog inserta un nuevo registro de evento enriquecido en MongoDB.
func (m *EnrichmentMongoRepository) InsertLog(ctx context.Context, event *models.EnrichedEventRecord) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // Usar el timeout del contexto, o definir uno si es nil
	defer cancel()

	_, err := m.mongoInstance.Collection.InsertOne(ctx, event)
	if err != nil {
		logger.ErrorLog.Printf("Error al insertar evento enriquecido en MongoDB: %v", err)
		return fmt.Errorf("error al insertar evento enriquecido: %w", err)
	}
	logger.InfoLog.Printf("Evento enriquecido insertado en MongoDB. EventSource: %s", event.EventSource)
	return nil
}

// GetLatestLogs recupera los últimos 10 registros de eventos enriquecidos de MongoDB.
func (m *EnrichmentMongoRepository) GetLatestLogs(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second) // Usar el timeout del contexto, o definir uno si es nil
	defer cancel()

	var events []*models.EnrichedEventRecord
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "eventTime", Value: -1}}) // Ordenar por fecha de evento descendente (-1 most recent)
	findOptions.SetLimit(10)                                   // Limitar a los últimos 10

	cursor, err := m.mongoInstance.Collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		logger.ErrorLog.Printf("Error al obtener los últimos eventos enriquecidos de MongoDB: %v", err)
		return nil, fmt.Errorf("error al obtener los últimos eventos enriquecidos: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event models.EnrichedEventRecord
		if err := cursor.Decode(&event); err != nil {
			logger.ErrorLog.Printf("Error al decodificar evento de MongoDB: %v", err)
			return nil, fmt.Errorf("error al decodificar evento: %w", err)
		}
		events = append(events, &event)
	}

	if err := cursor.Err(); err != nil {
		logger.ErrorLog.Printf("Error en el cursor de MongoDB: %v", err)
		return nil, fmt.Errorf("error en el cursor de MongoDB: %w", err)
	}

	logger.InfoLog.Println("Últimos 10 eventos enriquecidos obtenidos de MongoDB.")
	return events, nil
}
