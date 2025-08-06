package services

import (
	"cloudtrail-enrichment-api-golang/internal/config"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"database/sql"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Inicializa el logger para que no falle en las pruebas.
func init() {
	logger.Init()
}

// MockAuthRepository es un mock para la interfaz repository.AuthRepository.
type MockAuthRepository struct {
	InsertUserFunc           func(ctx context.Context, user *models.User) error
	GetUserByEmailFunc       func(ctx context.Context, email string) (*models.User, error)
	GetUserByUUIDFunc        func(ctx context.Context, uuid string) (*models.User, error)
	InsertTokenFunc          func(ctx context.Context, token *models.Token) error
	DeleteTokensByUserIDFunc func(ctx context.Context, userID int) error
	GetUserForTokenFunc      func(ctx context.Context, userID int) (*models.User, error)
	GetTokenByTokenFunc      func(ctx context.Context, tokenString string) (*models.Token, error)
	GetTokenByTokenHashFunc  func(ctx context.Context, tokenHash string) (*models.Token, error)
}

func (m *MockAuthRepository) InsertUser(ctx context.Context, user *models.User) error {
	return m.InsertUserFunc(ctx, user)
}
func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.GetUserByEmailFunc(ctx, email)
}
func (m *MockAuthRepository) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	return m.GetUserByUUIDFunc(ctx, uuid)
}
func (m *MockAuthRepository) InsertToken(ctx context.Context, t *models.Token) error {
	return m.InsertTokenFunc(ctx, t)
}
func (m *MockAuthRepository) GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error) {
	return m.GetTokenByTokenHashFunc(ctx, tokenHash)
}
func (m *MockAuthRepository) DeleteTokensByUserID(ctx context.Context, userID int) error {
	return m.DeleteTokensByUserIDFunc(ctx, userID)
}
func (m *MockAuthRepository) GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error) {
	return m.GetTokenByTokenFunc(ctx, tokenString)
}
func (m *MockAuthRepository) GetUserForToken(ctx context.Context, userID int) (*models.User, error) {
	return m.GetUserForTokenFunc(ctx, userID)
}

// MockTokenDBRepository es un mock para la interfaz TokenDBRepository del paquete token.
type MockTokenDBRepository struct {
	InsertTokenFunc          func(ctx context.Context, token *models.Token) error
	GetTokenByTokenHashFunc  func(ctx context.Context, tokenHash string) (*models.Token, error)
	DeleteTokensByUserIDFunc func(ctx context.Context, userID int) error
	GetTokenByTokenFunc      func(ctx context.Context, tokenString string) (*models.Token, error)
	GetUserForTokenFunc      func(ctx context.Context, userID int) (*models.User, error)
}

func (m *MockTokenDBRepository) InsertToken(ctx context.Context, t *models.Token) error {
	return m.InsertTokenFunc(ctx, t)
}
func (m *MockTokenDBRepository) GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error) {
	return m.GetTokenByTokenHashFunc(ctx, tokenHash)
}
func (m *MockTokenDBRepository) DeleteTokensByUserID(ctx context.Context, userID int) error {
	return m.DeleteTokensByUserIDFunc(ctx, userID)
}
func (m *MockTokenDBRepository) GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error) {
	return m.GetTokenByTokenFunc(ctx, tokenString)
}
func (m *MockTokenDBRepository) GetUserForToken(ctx context.Context, userID int) (*models.User, error) {
	return m.GetUserForTokenFunc(ctx, userID)
}

// --- Pruebas para el método RegisterUser ---

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := &MockAuthRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return nil, sql.ErrNoRows
		},
		InsertUserFunc: func(ctx context.Context, user *models.User) error {
			return nil
		},
	}

	// Como RegisterUser no usa el JWTService, podemos pasar nil
	service := NewAuthService(mockRepo, nil)

	payload := &models.RegisterPayload{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}
	user, err := service.RegisterUser(context.Background(), payload)

	if err != nil {
		t.Errorf("Se esperaba un registro exitoso, pero se obtuvo un error: %v", err)
	}
	if user == nil {
		t.Error("Se esperaba que se devolviera un usuario, pero se obtuvo nil")
	}
	if user.Email != payload.Email {
		t.Errorf("El email del usuario no coincide. Esperado: %s, Obtenido: '%s'", payload.Email, user.Email)
	}
}

// --- Pruebas para el método AuthenticateUser ---

func TestAuthenticateUser_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &MockAuthRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{
				Email:        email,
				PasswordHash: string(hashedPassword),
				Role:         "user",
				ID:           1,
			}, nil
		},
	}

	// Creamos un mock de las dependencias de JWTService
	mockTokenRepo := &MockTokenDBRepository{
		DeleteTokensByUserIDFunc: func(ctx context.Context, userID int) error { return nil },
		InsertTokenFunc:          func(ctx context.Context, t *models.Token) error { return nil },
	}

	mockConfig := &config.Config{
		AuthConfig: config.AuthConfig{
			JWTSecret:     "secretkey-test",
			TokenDuration: 1 * time.Hour,
		},
	}

	// Instanciamos el JWTService real con las dependencias mockeadas
	jwtService := token.NewJWTService(mockConfig, mockTokenRepo)

	// Usamos el constructor para instanciar el servicio principal
	service := NewAuthService(mockRepo, jwtService)

	user, jwtToken, err := service.AuthenticateUser(context.Background(), "test@example.com", "password123")

	if err != nil {
		t.Errorf("Se esperaba una autenticación exitosa, pero se obtuvo un error: %v", err)
	}
	if user == nil || jwtToken == nil {
		t.Error("Se esperaba que se devolvieran un usuario y un token, pero se obtuvo nil")
	}
	if jwtToken.Token == "" {
		t.Errorf("Se esperaba que el token no estuviera vacío, pero se obtuvo: '%s'", jwtToken.Token)
	}
}

func TestAuthenticateUser_InvalidCredentials(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &MockAuthRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{
				Email:        email,
				PasswordHash: string(hashedPassword),
				Role:         "user",
				ID:           1,
			}, nil
		},
	}

	// Creamos un mock de las dependencias de JWTService
	mockTokenRepo := &MockTokenDBRepository{
		DeleteTokensByUserIDFunc: func(ctx context.Context, userID int) error { return nil },
		InsertTokenFunc:          func(ctx context.Context, t *models.Token) error { return nil },
	}

	mockConfig := &config.Config{
		AuthConfig: config.AuthConfig{
			JWTSecret:     "secretkey-test",
			TokenDuration: 1 * time.Hour,
		},
	}

	// Instanciamos el JWTService real con las dependencias mockeadas
	jwtService := token.NewJWTService(mockConfig, mockTokenRepo)

	// Usamos el constructor para instanciar el servicio principal
	service := NewAuthService(mockRepo, jwtService)

	_, _, err := service.AuthenticateUser(context.Background(), "test@example.com", "wrongpassword")

	if err == nil {
		t.Error("Se esperaba un error por credenciales inválidas, pero se obtuvo nil")
	}
	if err.Error() != "credenciales inválidas" {
		t.Errorf("Mensaje de error inesperado. Esperado: 'credenciales inválidas', Obtenido: '%s'", err.Error())
	}
}
