package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloudtrail-enrichment-api-golang/internal/config"
	"cloudtrail-enrichment-api-golang/internal/middleware"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token"
	"cloudtrail-enrichment-api-golang/internal/pkg/utils"
	"cloudtrail-enrichment-api-golang/models"
	"cloudtrail-enrichment-api-golang/services"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// --- Mocks para los controladores y servicios ---

// MockSystemController es un mock simple del SystemController.
type MockSystemController struct {
	HealthCheckCalled bool
}

func (m *MockSystemController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	m.HealthCheckCalled = true
	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Error:   false,
		Message: "health_check_ok",
	})
}

// MockAuthController es un mock simple del AuthController.
type MockAuthController struct {
	RegisterUserCalled     bool
	AuthenticateUserCalled bool
}

func (m *MockAuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	m.RegisterUserCalled = true
	utils.WriteJSON(w, http.StatusCreated, utils.JSONResponse{
		Error:   false,
		Message: "register_user_ok",
	})
}

func (m *MockAuthController) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	m.AuthenticateUserCalled = true
	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Error:   false,
		Message: "authenticate_user_ok",
	})
}

// MockEnrichmentController es un mock simple del EnrichmentController.
type MockEnrichmentController struct {
	IngestDataCalled  bool
	QueryEventsCalled bool
}

func (m *MockEnrichmentController) IngestData(w http.ResponseWriter, r *http.Request) {
	m.IngestDataCalled = true
	utils.WriteJSON(w, http.StatusAccepted, utils.JSONResponse{
		Error:   false,
		Message: "ingest_data_ok",
	})
}

func (m *MockEnrichmentController) QueryEvents(w http.ResponseWriter, r *http.Request) {
	m.QueryEventsCalled = true
	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Error:   false,
		Message: "query_events_ok",
	})
}

// --- Mocks para el middleware y servicios ---

// MockAuthService implementa la interfaz services.AuthService para las pruebas.
type MockAuthService struct {
	ValidateTokenForMiddlewareFunc func(ctx context.Context, tokenString string) (*token.User, error)
	RegisterUserCalled             bool
	AuthenticateUserCalled         bool
}

func (m *MockAuthService) RegisterUser(ctx context.Context, user *models.RegisterPayload) (*models.User, error) {
	m.RegisterUserCalled = true
	return &models.User{}, nil
}

func (m *MockAuthService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error) {
	m.AuthenticateUserCalled = true
	return &models.User{}, nil, nil
}

func (m *MockAuthService) ValidateTokenForMiddleware(ctx context.Context, tokenString string) (*token.User, error) {
	return m.ValidateTokenForMiddlewareFunc(ctx, tokenString)
}

var originalInfoLog = logger.InfoLog
var originalErrorLog = logger.ErrorLog

// setupMockLogger inicializa el logger global con un mock.
func setupMockLogger() {
	// Crea un nuevo log.Logger que descarta toda la salida
	silentLogger := log.New(io.Discard, "", 0)
	logger.InfoLog = silentLogger
	logger.ErrorLog = silentLogger
}

// teardownMockLogger restaura el logger original.
func teardownMockLogger() {
	logger.InfoLog = originalInfoLog
	logger.ErrorLog = originalErrorLog
}

// MockApplication es una versión de la aplicación para tests.
type MockApplication struct {
	config               *config.Config
	middleware           *middleware.Middleware
	authController       *MockAuthController
	systemController     *MockSystemController
	enrichmentController *MockEnrichmentController
}

func (app *MockApplication) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(chiMiddleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Grupo de rutas para la versión 1 de la API
	mux.Route("/v1", func(r chi.Router) {
		// Rutas públicas de la V1
		r.Get("/health", app.systemController.HealthCheck)
		r.Post("/signup", app.authController.RegisterUser)
		r.Post("/login", app.authController.AuthenticateUser)

		// Rutas protegidas por el middleware de autenticación de la V1
		r.Route("/enrichment", func(r chi.Router) {
			r.Use(app.middleware.AuthTokenMiddleware)
			r.Post("/", app.enrichmentController.IngestData)
			r.Get("/", app.enrichmentController.QueryEvents)
		})
	})

	return mux
}

// NewMockMiddleware crea un mock de middleware compatible con tu código fuente.
// Se ha eliminado la dependencia de MockLogger.
func NewMockMiddleware(authService services.AuthService) *middleware.Middleware {
	mockJWTService := &token.JWTService{
		Config: &config.Config{
			AuthConfig: config.AuthConfig{
				JWTSecret: "test-secret-key",
			},
		},
	}
	// Se crea el middleware con los dos argumentos que tu función NewMiddleware
	// espera actualmente.
	return middleware.NewMiddleware(mockJWTService, authService)
}

// --- Suite de pruebas para las rutas ---

func TestRoutes(t *testing.T) {
	// Inicializar el logger con un mock antes de las pruebas.
	setupMockLogger()
	defer teardownMockLogger()

	// 1. Definir los casos de prueba
	tests := []struct {
		name                 string
		method               string
		path                 string
		authRequired         bool
		authPasses           bool
		expectedStatusCode   int
		controllerFlagToTest func(*MockSystemController, *MockAuthController, *MockEnrichmentController) bool
	}{
		// Rutas públicas
		{
			name:               "HealthCheck returns 200 OK",
			method:             "GET",
			path:               "/v1/health",
			authRequired:       false,
			expectedStatusCode: http.StatusOK,
			controllerFlagToTest: func(sc *MockSystemController, ac *MockAuthController, ec *MockEnrichmentController) bool {
				return sc.HealthCheckCalled
			},
		},
		{
			name:               "RegisterUser returns 201 Created",
			method:             "POST",
			path:               "/v1/signup",
			authRequired:       false,
			expectedStatusCode: http.StatusCreated,
			controllerFlagToTest: func(sc *MockSystemController, ac *MockAuthController, ec *MockEnrichmentController) bool {
				return ac.RegisterUserCalled
			},
		},
		{
			name:               "AuthenticateUser returns 200 OK",
			method:             "POST",
			path:               "/v1/login",
			authRequired:       false,
			expectedStatusCode: http.StatusOK,
			controllerFlagToTest: func(sc *MockSystemController, ac *MockAuthController, ec *MockEnrichmentController) bool {
				return ac.AuthenticateUserCalled
			},
		},
		// Rutas protegidas (con autenticación exitosa)
		{
			name:               "IngestData with valid token returns 202 Accepted",
			method:             "POST",
			path:               "/v1/enrichment",
			authRequired:       true,
			authPasses:         true,
			expectedStatusCode: http.StatusAccepted,
			controllerFlagToTest: func(sc *MockSystemController, ac *MockAuthController, ec *MockEnrichmentController) bool {
				return ec.IngestDataCalled
			},
		},
		{
			name:               "QueryEvents with valid token returns 200 OK",
			method:             "GET",
			path:               "/v1/enrichment",
			authRequired:       true,
			authPasses:         true,
			expectedStatusCode: http.StatusOK,
			controllerFlagToTest: func(sc *MockSystemController, ac *MockAuthController, ec *MockEnrichmentController) bool {
				return ec.QueryEventsCalled
			},
		},
		// Rutas protegidas (con autenticación fallida)
		{
			name:               "IngestData with invalid token returns 401 Unauthorized",
			method:             "POST",
			path:               "/v1/enrichment",
			authRequired:       true,
			authPasses:         false,
			expectedStatusCode: http.StatusUnauthorized,
			controllerFlagToTest: func(sc *MockSystemController, ac *MockAuthController, ec *MockEnrichmentController) bool {
				return false
			},
		},
		{
			name:               "QueryEvents with invalid token returns 401 Unauthorized",
			method:             "GET",
			path:               "/v1/enrichment",
			authRequired:       true,
			authPasses:         false,
			expectedStatusCode: http.StatusUnauthorized,
			controllerFlagToTest: func(sc *MockSystemController, ac *MockAuthController, ec *MockEnrichmentController) bool {
				return false
			},
		},
	}

	// 2. Ejecutar las pruebas
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Inicializar mocks y la aplicación de prueba dentro de cada t.Run
			mockAuthService := &MockAuthService{}
			mockMiddleware := NewMockMiddleware(mockAuthService)
			mockSystemController := &MockSystemController{}
			mockAuthController := &MockAuthController{}
			mockEnrichmentController := &MockEnrichmentController{}

			mockApp := &MockApplication{
				config:               &config.Config{},
				middleware:           mockMiddleware,
				authController:       mockAuthController,
				systemController:     mockSystemController,
				enrichmentController: mockEnrichmentController,
			}

			resetMockFlags := func() {
				mockSystemController.HealthCheckCalled = false
				mockAuthController.RegisterUserCalled = false
				mockAuthController.AuthenticateUserCalled = false
				mockEnrichmentController.IngestDataCalled = false
				mockEnrichmentController.QueryEventsCalled = false
			}
			resetMockFlags()

			// Configurar el mock del servicio de autenticación para que pase o falle
			if tt.authRequired {
				if tt.authPasses {
					mockAuthService.ValidateTokenForMiddlewareFunc = func(ctx context.Context, tokenString string) (*token.User, error) {
						return &token.User{ID: 1, Email: "test@example.com", Role: "user"}, nil
					}
				} else {
					mockAuthService.ValidateTokenForMiddlewareFunc = func(ctx context.Context, tokenString string) (*token.User, error) {
						return nil, errors.New("invalid token")
					}
				}
			}

			// Crear una petición HTTP simulada
			req, _ := http.NewRequest(tt.method, tt.path, bytes.NewBufferString(`{}`))
			if tt.authRequired && tt.authPasses {
				req.Header.Set("Authorization", "Bearer mock_token")
			}

			// Crear un ResponseRecorder para capturar la respuesta
			rr := httptest.NewRecorder()

			// Ejecutar la ruta
			mockApp.routes().ServeHTTP(rr, req)

			// Verificar el código de estado
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("para la ruta '%s', se esperaba el código de estado %d, se obtuvo %d", tt.path, tt.expectedStatusCode, rr.Code)
			}

			// Verificar si se llamó al controlador
			if tt.controllerFlagToTest(mockSystemController, mockAuthController, mockEnrichmentController) == false && tt.authPasses {
				t.Errorf("se esperaba que el controlador de la ruta '%s' fuera llamado, pero no lo fue", tt.path)
			}
		})
	}
}
