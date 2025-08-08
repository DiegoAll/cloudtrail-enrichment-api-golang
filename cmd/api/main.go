package main

import (
	"cloudtrail-enrichment-api-golang/cmd/api/controllers"
	"cloudtrail-enrichment-api-golang/database/mongo"
	"cloudtrail-enrichment-api-golang/database/postgresql"
	"cloudtrail-enrichment-api-golang/internal/config"
	"cloudtrail-enrichment-api-golang/internal/middleware"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token"
	"cloudtrail-enrichment-api-golang/internal/repository"
	"cloudtrail-enrichment-api-golang/services"
	"context"

	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

type application struct {
	config   *config.Config
	infoLog  *log.Logger
	errorLog *log.Logger
	// models     models.Models
	middleware           *middleware.Middleware
	authController       *controllers.AuthController
	systemController     *controllers.SystemController
	enrichmentController *controllers.EnrichmentController
}

func main() {
	logger.Init() // Initializes logger with os.Stdout, os.Stderr
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
		logger.ErrorLog.Fatalf("Error loading configuration: %v", err)
	}

	logger.InfoLog.Println("Server started")
	fmt.Println(config.DatabaseConfig.Database)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DatabaseConfig.Host,
		config.DatabaseConfig.Port,
		config.DatabaseConfig.Username,
		config.DatabaseConfig.Password,
		config.DatabaseConfig.Database,
		config.DatabaseConfig.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
		logger.ErrorLog.Fatalf("Error opening database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to the database:", err)
		logger.ErrorLog.Fatalf("Error connecting to the database: %v", err)
	}
	logger.InfoLog.Println("PostgreSQL database connection established successfully.")

	defer func() {
		if err := db.Close(); err != nil {
			logger.ErrorLog.Printf("Error closing main database connection: %v", err)
		}
	}()

	mongoURI := mongo.BuildMongoURI(&config.MongoDBConfig)

	mongoClient, err := mongo.NewMongoClient(mongoURI, config.MongoDBConfig.DBTimeout)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		logger.ErrorLog.Fatalf("Error connecting to MongoDB: %v", err)
	}
	// logger.InfoLog.Println("Connection to MongoDB successfully established.")

	defer func() {
		// Disconnect the MongoDB client at the end
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			logger.ErrorLog.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Initialization of authentication repository
	authRepo := postgresql.NewAuthPostgresRepository(db)
	// Initialization of enrichment repository
	enrichRepo := mongo.NewEnrichMongoRepository(mongoClient, config.MongoDBConfig.Database, config.MongoDBConfig.Collection)

	repository.SetAuthRepository(authRepo)         // Set the global authRepo implementation
	repository.SetEnrichmentRepository(enrichRepo) // Set the global enrichRepo implementation

	// Create an instance of JWTService.
	jwtService := token.NewJWTService(config, authRepo) // IMPORTANT CHANGE HERE

	// Services initialization
	authService := services.NewAuthService(repository.AuthRepo, jwtService)
	enrichService := services.NewDefaultEnrichmentService(repository.EnrichmentRepo)

	// Controllers initialization
	authController := controllers.NewAuthController(authService)
	systemController := controllers.NewSystemController()
	enrichmentController := controllers.NewEnrichmentController(enrichService)

	// Middleware initialization
	mw := middleware.NewMiddleware(jwtService, authService)

	app := &application{
		config:               config,
		infoLog:              logger.InfoLog,
		errorLog:             logger.ErrorLog,
		middleware:           mw,
		authController:       authController,
		systemController:     systemController,
		enrichmentController: enrichmentController,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
		logger.ErrorLog.Fatal(err)
	}
}

func (app *application) serve() error {
	app.infoLog.Printf("API listening on port %d", app.config.ServerConfig.Port)

	certFile := app.config.ServerConfig.TLS.CertFile
	keyFile := app.config.ServerConfig.TLS.KeyFile

	// Only load TLS certificates if the files exist
	if _, err := os.Stat(certFile); os.IsNotExist(err) || (certFile == "" && keyFile == "") {
		app.infoLog.Println("TLS certificate files not found or not specified. Starting plain HTTP server.")
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", app.config.ServerConfig.Port),
			Handler: app.routes(),
		}
		return srv.ListenAndServe()
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("error loading TLS certificates: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
		ClientAuth:   tls.RequireAnyClientCert,
	}

	srv := &http.Server{
		Addr:      fmt.Sprintf(":%d", app.config.ServerConfig.Port),
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
	}

	app.infoLog.Printf("Starting HTTPS server with certificates at %s and %s", certFile, keyFile)
	return srv.ListenAndServeTLS("", "")
}
