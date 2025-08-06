package repository_test

import (
	"cloudtrail-enrichment-api-golang/internal/repository" // Importamos el paquete principal
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

// MockAuthRepository es un mock simple para la interfaz AuthRepository
type MockAuthRepository struct {
	InsertUserFunc           func(ctx context.Context, user *models.User) error
	GetUserByEmailFunc       func(ctx context.Context, email string) (*models.User, error)
	GetUserByUUIDFunc        func(ctx context.Context, uuid string) (*models.User, error)
	InsertTokenFunc          func(ctx context.Context, token *models.Token) error
	GetTokenByTokenHashFunc  func(ctx context.Context, tokenHash string) (*models.Token, error)
	DeleteTokensByUserIDFunc func(ctx context.Context, userID int) error
	GetTokenByTokenFunc      func(ctx context.Context, tokenString string) (*models.Token, error)
	GetUserForTokenFunc      func(ctx context.Context, userID int) (*models.User, error)
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

func (m *MockAuthRepository) InsertToken(ctx context.Context, token *models.Token) error {
	return m.InsertTokenFunc(ctx, token)
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

// TestGlobalFunctionsOfAuthRepository prueba las funciones de fachada del paquete repository
func TestGlobalFunctionsOfAuthRepository(t *testing.T) {
	ctx := context.Background()

	// --- Pruebas para InsertUser ---
	t.Run("InsertUser success", func(t *testing.T) {
		mockRepo := &MockAuthRepository{
			InsertUserFunc: func(ctx context.Context, user *models.User) error { return nil },
		}
		repository.SetAuthRepository(mockRepo)
		err := repository.InsertUser(ctx, &models.User{})
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
	})

	t.Run("InsertUser failure", func(t *testing.T) {
		expectedErr := errors.New("error simulado")
		mockRepo := &MockAuthRepository{
			InsertUserFunc: func(ctx context.Context, user *models.User) error { return expectedErr },
		}
		repository.SetAuthRepository(mockRepo)
		err := repository.InsertUser(ctx, &models.User{})
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("se esperaba el error '%v', se obtuvo: %v", expectedErr, err)
		}
	})

	// --- Pruebas para GetUserByEmail ---
	t.Run("GetUserByEmail success", func(t *testing.T) {
		expectedUser := &models.User{Email: "test@example.com"}
		mockRepo := &MockAuthRepository{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) { return expectedUser, nil },
		}
		repository.SetAuthRepository(mockRepo)
		user, err := repository.GetUserByEmail(ctx, "test@example.com")
		if err != nil || user != expectedUser {
			t.Errorf("se esperaba el usuario '%v', se obtuvo: '%v', error: %v", expectedUser, user, err)
		}
	})

	t.Run("GetUserByEmail failure", func(t *testing.T) {
		expectedErr := errors.New("error simulado")
		mockRepo := &MockAuthRepository{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) { return nil, expectedErr },
		}
		repository.SetAuthRepository(mockRepo)
		_, err := repository.GetUserByEmail(ctx, "test@example.com")
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("se esperaba el error '%v', se obtuvo: %v", expectedErr, err)
		}
	})

	// --- Pruebas para GetUserByUUID ---
	t.Run("GetUserByUUID success", func(t *testing.T) {
		testUUID := uuid.New().String()
		expectedUser := &models.User{UUID: testUUID}
		mockRepo := &MockAuthRepository{
			GetUserByUUIDFunc: func(ctx context.Context, uuid string) (*models.User, error) { return expectedUser, nil },
		}
		repository.SetAuthRepository(mockRepo)
		user, err := repository.GetUserByUUID(ctx, testUUID)
		if err != nil || user != expectedUser {
			t.Errorf("se esperaba el usuario '%v', se obtuvo: '%v', error: %v", expectedUser, user, err)
		}
	})

	t.Run("GetUserByUUID failure", func(t *testing.T) {
		expectedErr := errors.New("error simulado")
		mockRepo := &MockAuthRepository{
			GetUserByUUIDFunc: func(ctx context.Context, uuid string) (*models.User, error) { return nil, expectedErr },
		}
		repository.SetAuthRepository(mockRepo)
		_, err := repository.GetUserByUUID(ctx, uuid.New().String())
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("se esperaba el error '%v', se obtuvo: %v", expectedErr, err)
		}
	})

	// --- Pruebas para InsertToken ---
	t.Run("InsertToken success", func(t *testing.T) {
		mockRepo := &MockAuthRepository{
			InsertTokenFunc: func(ctx context.Context, token *models.Token) error { return nil },
		}
		repository.SetAuthRepository(mockRepo)
		err := repository.InsertToken(ctx, &models.Token{})
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
	})

	t.Run("InsertToken failure", func(t *testing.T) {
		expectedErr := errors.New("error simulado")
		mockRepo := &MockAuthRepository{
			InsertTokenFunc: func(ctx context.Context, token *models.Token) error { return expectedErr },
		}
		repository.SetAuthRepository(mockRepo)
		err := repository.InsertToken(ctx, &models.Token{})
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("se esperaba el error '%v', se obtuvo: %v", expectedErr, err)
		}
	})

	// --- Pruebas para GetTokenByTokenHash ---
	t.Run("GetTokenByTokenHash success", func(t *testing.T) {
		expectedToken := &models.Token{TokenHash: "somehash"}
		mockRepo := &MockAuthRepository{
			GetTokenByTokenHashFunc: func(ctx context.Context, tokenHash string) (*models.Token, error) { return expectedToken, nil },
		}
		repository.SetAuthRepository(mockRepo)
		token, err := repository.GetTokenByTokenHash(ctx, "somehash")
		if err != nil || token != expectedToken {
			t.Errorf("se esperaba el token '%v', se obtuvo: '%v', error: %v", expectedToken, token, err)
		}
	})

	t.Run("GetTokenByTokenHash failure", func(t *testing.T) {
		expectedErr := errors.New("error simulado")
		mockRepo := &MockAuthRepository{
			GetTokenByTokenHashFunc: func(ctx context.Context, tokenHash string) (*models.Token, error) { return nil, expectedErr },
		}
		repository.SetAuthRepository(mockRepo)
		_, err := repository.GetTokenByTokenHash(ctx, "somehash")
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("se esperaba el error '%v', se obtuvo: %v", expectedErr, err)
		}
	})

	// --- Pruebas para DeleteTokensByUserID ---
	t.Run("DeleteTokensByUserID success", func(t *testing.T) {
		mockRepo := &MockAuthRepository{
			DeleteTokensByUserIDFunc: func(ctx context.Context, userID int) error { return nil },
		}
		repository.SetAuthRepository(mockRepo)
		err := repository.DeleteTokensByUserID(ctx, 1)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
	})

	t.Run("DeleteTokensByUserID failure", func(t *testing.T) {
		expectedErr := errors.New("error simulado")
		mockRepo := &MockAuthRepository{
			DeleteTokensByUserIDFunc: func(ctx context.Context, userID int) error { return expectedErr },
		}
		repository.SetAuthRepository(mockRepo)
		err := repository.DeleteTokensByUserID(ctx, 1)
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("se esperaba el error '%v', se obtuvo: %v", expectedErr, err)
		}
	})

	// --- Pruebas para GetTokenByToken ---
	t.Run("GetTokenByToken success", func(t *testing.T) {
		expectedToken := &models.Token{Token: "sometoken"}
		mockRepo := &MockAuthRepository{
			GetTokenByTokenFunc: func(ctx context.Context, tokenString string) (*models.Token, error) { return expectedToken, nil },
		}
		repository.SetAuthRepository(mockRepo)
		token, err := repository.GetTokenByToken(ctx, "sometoken")
		if err != nil || token != expectedToken {
			t.Errorf("se esperaba el token '%v', se obtuvo: '%v', error: %v", expectedToken, token, err)
		}
	})

	t.Run("GetTokenByToken failure", func(t *testing.T) {
		expectedErr := errors.New("error simulado")
		mockRepo := &MockAuthRepository{
			GetTokenByTokenFunc: func(ctx context.Context, tokenString string) (*models.Token, error) { return nil, expectedErr },
		}
		repository.SetAuthRepository(mockRepo)
		_, err := repository.GetTokenByToken(ctx, "sometoken")
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("se esperaba el error '%v', se obtuvo: %v", expectedErr, err)
		}
	})

	// --- Pruebas para GetUserForToken ---
	t.Run("GetUserForToken success", func(t *testing.T) {
		expectedUser := &models.User{ID: 1}
		mockRepo := &MockAuthRepository{
			GetUserForTokenFunc: func(ctx context.Context, userID int) (*models.User, error) { return expectedUser, nil },
		}
		repository.SetAuthRepository(mockRepo)
		user, err := repository.GetUserForToken(ctx, 1)
		if err != nil || user != expectedUser {
			t.Errorf("se esperaba el usuario '%v', se obtuvo: '%v', error: %v", expectedUser, user, err)
		}
	})

	t.Run("GetUserForToken failure", func(t *testing.T) {
		expectedErr := errors.New("error simulado")
		mockRepo := &MockAuthRepository{
			GetUserForTokenFunc: func(ctx context.Context, userID int) (*models.User, error) { return nil, expectedErr },
		}
		repository.SetAuthRepository(mockRepo)
		_, err := repository.GetUserForToken(ctx, 1)
		if err == nil || !errors.Is(err, expectedErr) {
			t.Errorf("se esperaba el error '%v', se obtuvo: %v", expectedErr, err)
		}
	})
}
