package token

import (
	"cloudtrail-enrichment-api-golang/internal/config"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"errors"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// mockTokenDBRepository es una implementación mock de la interfaz TokenDBRepository
// para usar en las pruebas.
type mockTokenDBRepository struct {
	tokens          map[string]*models.Token
	users           map[int]*models.User
	insertCalled    bool
	deleteCalled    bool
	getByHashCalled bool
}

// NewMockTokenDBRepository crea una nueva instancia del mock
func NewMockTokenDBRepository() *mockTokenDBRepository {
	return &mockTokenDBRepository{
		tokens: make(map[string]*models.Token),
		users:  make(map[int]*models.User),
	}
}

// InsertToken simula la inserción de un token en la base de datos.
func (m *mockTokenDBRepository) InsertToken(ctx context.Context, token *models.Token) error {
	m.insertCalled = true
	// Usamos el token en texto plano como clave para simplificar
	m.tokens[token.Token] = token
	return nil
}

// GetTokenByTokenHash simula la búsqueda de un token por su hash.
func (m *mockTokenDBRepository) GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error) {
	m.getByHashCalled = true
	for _, token := range m.tokens {
		if token.TokenHash == tokenHash {
			return token, nil
		}
	}
	return nil, errors.New("token no encontrado")
}

// DeleteTokensByUserID simula la eliminación de tokens por ID de usuario.
func (m *mockTokenDBRepository) DeleteTokensByUserID(ctx context.Context, userID int) error {
	m.deleteCalled = true
	for key, token := range m.tokens {
		if token.UserID == userID {
			delete(m.tokens, key)
		}
	}
	return nil
}

// GetTokenByToken simula la búsqueda de un token por su valor en texto plano.
func (m *mockTokenDBRepository) GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error) {
	token, ok := m.tokens[tokenString]
	if !ok {
		return nil, errors.New("token no encontrado")
	}
	return token, nil
}

// GetUserForToken simula la obtención de un usuario por ID.
func (m *mockTokenDBRepository) GetUserForToken(ctx context.Context, userID int) (*models.User, error) {
	user, ok := m.users[userID]
	if !ok {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

// TestMain se ejecuta antes de cualquier función de prueba en este paquete.
func TestMain(m *testing.M) {
	logger.Init() // Se asume que esta función configura los loggers correctamente
	logger.InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.DebugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	os.Exit(m.Run())
}

// TestNewJWTService valida que el servicio se inicializa correctamente.
func TestNewJWTService(t *testing.T) {
	cfg := &config.Config{}
	repo := &mockTokenDBRepository{}

	service := NewJWTService(cfg, repo)

	if service.Config != cfg {
		t.Error("El servicio no se inicializó con la configuración correcta.")
	}
	if service.TokenRepository != repo {
		t.Error("El servicio no se inicializó con el repositorio correcto.")
	}
}

// TestGenerateJWTToken_Success prueba la generación exitosa de un token.
func TestGenerateJWTToken_Success(t *testing.T) {
	// Configuración de prueba
	cfg := &config.Config{
		AuthConfig: config.AuthConfig{
			JWTSecret:     "secret-key-para-tests",
			TokenDuration: 24 * time.Hour,
		},
		DatabaseConfig: config.DatabaseConfig{
			DBTimeout: 1 * time.Second,
		},
	}
	repo := NewMockTokenDBRepository()
	service := NewJWTService(cfg, repo)

	userID := 1
	email := "test@example.com"
	role := "user"
	user := models.User{ID: userID, Email: email, Role: role}
	repo.users[userID] = &user

	// Se inserta un token de prueba para que sea revocado.
	repo.tokens["oldtoken"] = &models.Token{UserID: userID, Token: "oldtoken"}

	// Generar el token
	signedToken, expiry, err := service.GenerateJWTToken(context.Background(), userID, email, role)

	// Validaciones
	if err != nil {
		t.Fatalf("GenerateJWTToken retornó un error inesperado: %v", err)
	}

	if signedToken == "" {
		t.Error("El token firmado no debería estar vacío.")
	}

	if expiry.IsZero() {
		t.Error("La fecha de expiración no debería ser cero.")
	}

	if !repo.deleteCalled {
		t.Error("No se llamó a la función de eliminación de tokens anteriores.")
	}

	if !repo.insertCalled {
		t.Error("No se llamó a la función de inserción del nuevo token.")
	}
}

// TestGenerateJWTToken_NoSecret prueba el error cuando no hay clave secreta.
func TestGenerateJWTToken_NoSecret(t *testing.T) {
	cfg := &config.Config{
		AuthConfig: config.AuthConfig{
			JWTSecret:     "",
			TokenDuration: 24 * time.Hour,
		},
		DatabaseConfig: config.DatabaseConfig{
			DBTimeout: 1 * time.Second,
		},
	}
	repo := NewMockTokenDBRepository()
	service := NewJWTService(cfg, repo)

	signedToken, _, err := service.GenerateJWTToken(context.Background(), 1, "test@example.com", "user")

	if err == nil {
		t.Error("Se esperaba un error por clave secreta no configurada, pero no se recibió.")
	}
	if signedToken != "" {
		t.Error("El token firmado debería estar vacío cuando hay un error.")
	}
}

// TestValidJWTToken_Success prueba la validación exitosa de un token válido.
func TestValidJWTToken_Success(t *testing.T) {
	// Configuración
	cfg := &config.Config{
		AuthConfig: config.AuthConfig{
			JWTSecret:     "secret-key-para-tests",
			TokenDuration: 24 * time.Hour,
		},
		DatabaseConfig: config.DatabaseConfig{
			DBTimeout: 1 * time.Second,
		},
	}
	repo := NewMockTokenDBRepository()
	service := NewJWTService(cfg, repo)

	userID := 1
	email := "test@example.com"
	role := "user"
	user := models.User{ID: userID, Email: email, Role: role}
	repo.users[userID] = &user

	// Generar y persistir el token
	signedToken, _, _ := service.GenerateJWTToken(context.Background(), userID, email, role)

	// Validar el token
	validUser, err := service.ValidJWTToken(context.Background(), signedToken)

	// Validaciones
	if err != nil {
		t.Fatalf("ValidJWTToken retornó un error inesperado para un token válido: %v", err)
	}

	if validUser.ID != userID || validUser.Email != email || validUser.Role != role {
		t.Errorf("La información del usuario no coincide. Se esperaba ID %d, Email %s, Role %s, se obtuvo ID %d, Email %s, Role %s",
			userID, email, role, validUser.ID, validUser.Email, validUser.Role)
	}

	if !repo.getByHashCalled {
		t.Error("No se llamó a la función de búsqueda de token por hash en la base de datos.")
	}
}

// TestValidJWTToken_Expired prueba que la validación falla con un token expirado.
func TestValidJWTToken_Expired(t *testing.T) {
	// Configuración con una duración muy corta
	cfg := &config.Config{
		AuthConfig: config.AuthConfig{
			JWTSecret:     "secret-key-para-tests",
			TokenDuration: 1 * time.Millisecond,
		},
		DatabaseConfig: config.DatabaseConfig{
			DBTimeout: 1 * time.Second,
		},
	}
	repo := NewMockTokenDBRepository()
	service := NewJWTService(cfg, repo)

	// Generar token
	signedToken, _, _ := service.GenerateJWTToken(context.Background(), 1, "test@example.com", "user")

	// Esperar a que expire
	time.Sleep(2 * time.Millisecond)

	// Validar el token
	user, err := service.ValidJWTToken(context.Background(), signedToken)

	// Validaciones
	if err == nil {
		t.Fatal("Se esperaba un error por token expirado, pero no se recibió.")
	}
	if !strings.Contains(err.Error(), "expirado") {
		t.Errorf("El error no indica que el token está expirado. Se obtuvo: %v", err)
	}
	if user != nil {
		t.Error("El usuario retornado no debería ser nil en caso de error.")
	}
}

// TestValidJWTToken_NotFoundInDB prueba que la validación falla si el token no está en la DB.
func TestValidJWTToken_NotFoundInDB(t *testing.T) {
	// Configuración
	cfg := &config.Config{
		AuthConfig: config.AuthConfig{
			JWTSecret:     "secret-key-para-tests",
			TokenDuration: 24 * time.Hour,
		},
		DatabaseConfig: config.DatabaseConfig{
			DBTimeout: 1 * time.Second,
		},
	}
	// Usamos un mock nuevo sin tokens persistidos.
	repo := NewMockTokenDBRepository()
	service := NewJWTService(cfg, repo)

	// Generar token (no se persistirá en el mock)
	claims := &JWTToken{
		UserID: 1,
		Email:  "test@example.com",
		Role:   "user",
		Expiry: time.Now().Add(24 * time.Hour),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "test@example.com",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(cfg.AuthConfig.JWTSecret))

	// Validar el token
	user, err := service.ValidJWTToken(context.Background(), signedToken)

	// Validaciones
	if err == nil {
		t.Fatal("Se esperaba un error por token no encontrado en DB, pero no se recibió.")
	}
	if !strings.Contains(err.Error(), "token no encontrado o revocado") {
		t.Errorf("El error no indica que el token no se encontró. Se obtuvo: %v", err)
	}
	if user != nil {
		t.Error("El usuario retornado no debería ser nil en caso de error.")
	}
}

// TestExtractJWTToken_Success prueba la extracción exitosa del token.
func TestExtractJWTToken_Success(t *testing.T) {
	service := &JWTService{}
	tokenString := "xyz123abc"
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	extractedToken, err := service.ExtractJWTToken(req)

	if err != nil {
		t.Fatalf("ExtractJWTToken retornó un error inesperado: %v", err)
	}
	if extractedToken != tokenString {
		t.Errorf("Token extraído incorrecto. Se esperaba %s, se obtuvo %s", tokenString, extractedToken)
	}
}

// TestExtractJWTToken_InvalidFormat prueba la extracción con un formato incorrecto.
func TestExtractJWTToken_InvalidFormat(t *testing.T) {
	service := &JWTService{}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidFormat token")

	_, err := service.ExtractJWTToken(req)

	if err == nil {
		t.Fatal("Se esperaba un error por formato inválido, pero no se recibió.")
	}
	if err.Error() != "formato de token inválido" {
		t.Errorf("Mensaje de error incorrecto. Se esperaba 'formato de token inválido', se obtuvo '%s'", err.Error())
	}
}

// TestDeleteByJWTToken_Success prueba la eliminación exitosa de un token.
func TestDeleteByJWTToken_Success(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			DBTimeout: 1 * time.Second,
		},
	}
	repo := NewMockTokenDBRepository()
	service := NewJWTService(cfg, repo)

	// Preparamos un token en el mock para que pueda ser eliminado.
	plainTextToken := "test-token"
	repo.tokens[plainTextToken] = &models.Token{
		UserID: 1,
		Token:  plainTextToken,
	}

	err := service.DeleteByJWTToken(context.Background(), plainTextToken)

	if err != nil {
		t.Fatalf("DeleteByJWTToken retornó un error inesperado: %v", err)
	}
	if !repo.deleteCalled {
		t.Error("La función de eliminación no fue llamada.")
	}
}

// TestGetByToken_Success prueba la obtención de un token por su valor en texto plano.
func TestGetByToken_Success(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			DBTimeout: 1 * time.Second,
		},
	}
	repo := NewMockTokenDBRepository()
	service := NewJWTService(cfg, repo)

	plainTextToken := "test-token-to-get"
	expectedToken := &models.Token{UserID: 1, Token: plainTextToken}
	repo.tokens[plainTextToken] = expectedToken

	token, err := service.GetByToken(context.Background(), plainTextToken)

	if err != nil {
		t.Fatalf("GetByToken retornó un error inesperado: %v", err)
	}
	if token.Token != expectedToken.Token {
		t.Errorf("Token incorrecto. Se esperaba %s, se obtuvo %s", expectedToken.Token, token.Token)
	}
}

// TestGetUserForToken_Success prueba la obtención de un usuario para un token.
func TestGetUserForToken_Success(t *testing.T) {
	cfg := &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			DBTimeout: 1 * time.Second,
		},
	}
	repo := NewMockTokenDBRepository()
	service := NewJWTService(cfg, repo)

	userID := 1
	expectedUser := &models.User{ID: userID, Email: "user@example.com"}
	repo.users[userID] = expectedUser

	token := models.Token{UserID: userID}

	user, err := service.GetUserForToken(context.Background(), token)

	if err != nil {
		t.Fatalf("GetUserForToken retornó un error inesperado: %v", err)
	}
	if user.ID != expectedUser.ID {
		t.Errorf("Usuario incorrecto. Se esperaba ID %d, se obtuvo ID %d", expectedUser.ID, user.ID)
	}
}
