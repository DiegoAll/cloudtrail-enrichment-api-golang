package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"cloudtrail-enrichment-api-golang/internal/config"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token"
	"cloudtrail-enrichment-api-golang/models"
)

// --- Mocks para las dependencias del middleware y servicios ---

// MockAuthService implementa la interfaz services.AuthService para las pruebas.
type MockAuthService struct {
	ValidateTokenForMiddlewareFunc func(ctx context.Context, tokenString string) (*token.User, error)
}

func (m *MockAuthService) RegisterUser(ctx context.Context, user *models.RegisterPayload) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (m *MockAuthService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error) {
	return nil, nil, errors.New("not implemented")
}

func (m *MockAuthService) ValidateTokenForMiddleware(ctx context.Context, tokenString string) (*token.User, error) {
	return m.ValidateTokenForMiddlewareFunc(ctx, tokenString)
}

// MockTokenDBRepository implementa la interfaz token.TokenDBRepository para las pruebas.
type MockTokenDBRepository struct {
	GetTokenByTokenHashFunc  func(ctx context.Context, tokenHash string) (*models.Token, error)
	GetTokenByTokenFunc      func(ctx context.Context, tokenString string) (*models.Token, error)
	DeleteTokensByUserIDFunc func(ctx context.Context, userID int) error
	InsertTokenFunc          func(ctx context.Context, token *models.Token) error
	GetUserForTokenFunc      func(ctx context.Context, userID int) (*models.User, error)
}

func (m *MockTokenDBRepository) GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error) {
	return m.GetTokenByTokenHashFunc(ctx, tokenHash)
}

func (m *MockTokenDBRepository) GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error) {
	return m.GetTokenByTokenFunc(ctx, tokenString)
}

func (m *MockTokenDBRepository) DeleteTokensByUserID(ctx context.Context, userID int) error {
	return m.DeleteTokensByUserIDFunc(ctx, userID)
}

func (m *MockTokenDBRepository) InsertToken(ctx context.Context, token *models.Token) error {
	return m.InsertTokenFunc(ctx, token)
}

func (m *MockTokenDBRepository) GetUserForToken(ctx context.Context, userID int) (*models.User, error) {
	return m.GetUserForTokenFunc(ctx, userID)
}

// TestMain se ejecuta una vez antes de todas las pruebas del paquete.
func TestMain(m *testing.M) {
	// Inicializar el logger para evitar panics.
	log.Println("Initializing logger for tests...")
	logger.Init()

	// Ejecutar todas las pruebas.
	os.Exit(m.Run())
}

// TestAuthTokenMiddleware es la suite de pruebas para el middleware.
func TestAuthTokenMiddleware(t *testing.T) {
	// 1. Definir los casos de prueba
	tests := []struct {
		name               string
		authHeader         string
		mockValidateToken  func(ctx context.Context, tokenString string) (*token.User, error)
		expectedStatusCode int
		shouldCallNext     bool
	}{
		{
			name:               "Solicitud sin token - debe fallar con 401",
			authHeader:         "",
			mockValidateToken:  nil,
			expectedStatusCode: http.StatusUnauthorized,
			shouldCallNext:     false,
		},
		{
			name:               "Solicitud con formato de token inválido - debe fallar con 401",
			authHeader:         "Invalid_Token",
			mockValidateToken:  nil,
			expectedStatusCode: http.StatusUnauthorized,
			shouldCallNext:     false,
		},
		{
			name:       "Solicitud con token inválido - debe fallar con 401",
			authHeader: "Bearer invalid_token",
			mockValidateToken: func(ctx context.Context, tokenString string) (*token.User, error) {
				return nil, errors.New("token inválido")
			},
			expectedStatusCode: http.StatusUnauthorized,
			shouldCallNext:     false,
		},
		{
			name:       "Solicitud con token válido - debe pasar con 200",
			authHeader: "Bearer valid_token",
			mockValidateToken: func(ctx context.Context, tokenString string) (*token.User, error) {
				return &token.User{ID: 1, Email: "test@example.com", Role: "user"}, nil
			},
			expectedStatusCode: http.StatusOK,
			shouldCallNext:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Crear los mocks para el servicio de autenticación
			mockAuthService := &MockAuthService{
				ValidateTokenForMiddlewareFunc: tt.mockValidateToken,
			}

			// Crear una configuración mock, necesaria para instanciar el JWTService
			mockConfig := &config.Config{
				AuthConfig: config.AuthConfig{
					JWTSecret:     "test-secret-key",
					TokenDuration: time.Hour,
				},
				DatabaseConfig: config.DatabaseConfig{
					DBTimeout: time.Second,
				},
			}

			// Creamos una instancia real de JWTService.
			// Esto soluciona el problema de tipos en NewMiddleware.
			realJWTService := token.NewJWTService(mockConfig, &MockTokenDBRepository{})

			// Creamos la instancia del middleware con la instancia real del JWTService y el mock de AuthService.
			mw := NewMiddleware(realJWTService, mockAuthService)

			// Handler "next" mockeado
			nextCalled := false
			var userClaimsFromContext *token.User
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				claims, ok := r.Context().Value("userClaims").(*token.User)
				if ok {
					userClaimsFromContext = claims
				}
				w.WriteHeader(http.StatusOK)
			})

			// Crear una solicitud HTTP simulada
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			// Ejecutar el middleware
			mw.AuthTokenMiddleware(nextHandler).ServeHTTP(rr, req)

			// Verificar el código de estado
			if rr.Code != tt.expectedStatusCode {
				t.Errorf("Código de estado incorrecto. Se esperaba %d, se obtuvo %d", tt.expectedStatusCode, rr.Code)
			}

			// Verificar si se llamó al siguiente handler
			if nextCalled != tt.shouldCallNext {
				t.Errorf("next handler llamado incorrectamente. Se esperaba %v, se obtuvo %v", tt.shouldCallNext, nextCalled)
			}

			// Para el caso de éxito, verificar que los claims del usuario están en el contexto
			if tt.shouldCallNext && userClaimsFromContext == nil {
				t.Error("Se esperaba que los claims del usuario estuvieran en el contexto, pero no se encontraron")
			} else if tt.shouldCallNext && userClaimsFromContext != nil && userClaimsFromContext.Email != "test@example.com" {
				t.Errorf("Claims incorrectos. Se esperaba 'test@example.com', se obtuvo '%s'", userClaimsFromContext.Email)
			}
		})
	}
}
