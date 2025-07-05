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
	middleware *middleware.Middleware
	// productsController *controllers.ProductsController
	authController   *controllers.AuthController
	systemController *controllers.SystemController
	//enrichmentController *controllers.EnrichmentController
}

func main() {
	logger.Init() // Inicializa logger con os.Stdout, os.Stderr
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error al cargar la configuración:", err)
		logger.ErrorLog.Fatalf("Error al cargar la configuración: %v", err)
	}

	logger.InfoLog.Println("Servidor iniciado")
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
		log.Fatal("Error al abrir la conexión con la base de datos:", err)
		logger.ErrorLog.Fatalf("Error al abrir la conexión con la base de datos: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
		logger.ErrorLog.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	logger.InfoLog.Println("Conexión a la base de datos PostgreSQL establecida exitosamente.")

	defer func() {
		if err := db.Close(); err != nil {
			logger.ErrorLog.Printf("Error al cerrar la conexión principal a la base de datos: %v", err)
		}
	}()

	// --- Conexión a MongoDB (para enriquecimiento) ---
	mongoClient, err := mongo.NewMongoClient(config)
	if err != nil {
		log.Fatal("Error al conectar a MongoDB:", err)
		logger.ErrorLog.Fatalf("Error al conectar a MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Client.Disconnect(context.Background()); err != nil {
			logger.ErrorLog.Printf("Error al cerrar la conexión a MongoDB: %v", err)
		}
	}()
	logger.InfoLog.Println("Conexión a MongoDB establecida exitosamente.")

	// Inicialización del repositorio de autenticación
	authRepo := postgresql.NewAuthPostgresRepository(db)
	enrichRepo := mongo.NewEnrichMongoRepository(mongoClient)

	repository.SetAuthRepository(authRepo)         // Setear la implementación global del AuthRepo
	repository.SetEnrichmentRepository(enrichRepo) // Setear la implementación global del AuthRepo

	// AHORA: Creamos una instancia de JWTService, no de JWTToken
	jwtService := token.NewJWTService(config, authRepo) // CAMBIO IMPORTANTE AQUÍ

	// Inicialización de servicios
	// PASAMOS jwtService al servicio de autenticación
	authService := services.NewAuthService(repository.AuthRepo, jwtService) // CAMBIO IMPORTANTE AQUÍ

	// Inicialización de controladores
	authController := controllers.NewAuthController(authService)
	systemController := controllers.NewSystemController()
	// enrichmentController := controllers.NewEnrichmentController(enrichService)

	// PASAMOS jwtService al middleware
	mw := middleware.NewMiddleware(jwtService, authService) // CAMBIO IMPORTANTE AQUÍ

	app := &application{
		config:     config,
		infoLog:    logger.InfoLog,
		errorLog:   logger.ErrorLog,
		middleware: mw,
		// productsController: productsController,
		authController:   authController,
		systemController: systemController,
		// enrichmentController: enrichmentController,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
		logger.ErrorLog.Fatal(err)
	}
}

func (app *application) serve() error {
	app.infoLog.Printf("API escuchando en el puerto %d", app.config.ServerConfig.Port)

	certFile := app.config.ServerConfig.TLS.CertFile
	keyFile := app.config.ServerConfig.TLS.KeyFile

	// Solo cargar certificados TLS si los archivos existen
	if _, err := os.Stat(certFile); os.IsNotExist(err) || (certFile == "" && keyFile == "") {
		app.infoLog.Println("Archivos de certificado TLS no encontrados o no especificados. Iniciando servidor HTTP plano.")
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", app.config.ServerConfig.Port),
			Handler: app.routes(),
		}
		return srv.ListenAndServe()
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("error cargando certificados TLS: %v", err)
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

	app.infoLog.Printf("Iniciando servidor HTTPS con certificados en %s y %s", certFile, keyFile)
	return srv.ListenAndServeTLS("", "")
}
