package controllers

import (
	"bytes"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// init es una función especial de Go que se ejecuta antes de los tests.
// La usamos para asegurarnos de que el logger esté inicializado.
func init() {
	logger.Init()
}

// TestRegisterUser_Success prueba el registro exitoso de un usuario.
func TestRegisterUser_Success(t *testing.T) {
	mockService := &MockAuthService{
		RegisterUserFunc: func(ctx context.Context, payload *models.RegisterPayload) (*models.User, error) {
			return &models.User{
				UUID:  "test-uuid",
				Email: payload.Email,
				Role:  "user",
			}, nil
		},
	}
	controller := NewAuthController(mockService)

	payload := models.RegisterPayload{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.RegisterUser(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusCreated, rr.Code)
	}
}

// TestRegisterUser_InvalidPayload prueba un payload de registro inválido.
func TestRegisterUser_InvalidPayload(t *testing.T) {
	mockService := &MockAuthService{}
	controller := NewAuthController(mockService)

	req := httptest.NewRequest("POST", "/signup", bytes.NewReader([]byte("not a json")))
	rr := httptest.NewRecorder()

	controller.RegisterUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusBadRequest, rr.Code)
	}
}

// TestAuthenticateUser_Success prueba la autenticación exitosa.
func TestAuthenticateUser_Success(t *testing.T) {
	mockService := &MockAuthService{
		AuthenticateUserFunc: func(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error) {
			user := &models.User{UUID: "test-uuid", Email: email, Role: "user"}
			jwtToken := &token.JWTToken{Token: "test.token"}
			return user, jwtToken, nil
		},
	}
	controller := NewAuthController(mockService)

	payload := models.LoginPayload{
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.AuthenticateUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusOK, rr.Code)
	}
}

// TestAuthenticateUser_Failure prueba una autenticación fallida.
func TestAuthenticateUser_Failure(t *testing.T) {
	mockService := &MockAuthService{
		AuthenticateUserFunc: func(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error) {
			return nil, nil, errors.New("credenciales inválidas")
		},
	}
	controller := NewAuthController(mockService)

	payload := models.LoginPayload{
		Email:    "wrong@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	controller.AuthenticateUser(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusUnauthorized, rr.Code)
	}
}
