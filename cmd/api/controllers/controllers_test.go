// cmd/api/controllers/controllers_test.go

package controllers

import (
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"net/http"
)

// --- Mocks para los Servicios ---

// MockAuthService implementa services.AuthService para pruebas.
type MockAuthService struct {
	RegisterUserFunc               func(ctx context.Context, payload *models.RegisterPayload) (*models.User, error)
	AuthenticateUserFunc           func(ctx context.Context, email, password string) (*models.User, *models.JWTToken, error)
	ValidateTokenForMiddlewareFunc func(ctx context.Context, tokenString string) (*models.Claims, error)
}

func (m *MockAuthService) RegisterUser(ctx context.Context, payload *models.RegisterPayload) (*models.User, error) {
	return m.RegisterUserFunc(ctx, payload)
}

func (m *MockAuthService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, *models.JWTToken, error) {
	return m.AuthenticateUserFunc(ctx, email, password)
}

func (m *MockAuthService) ValidateTokenForMiddleware(ctx context.Context, tokenString string) (*models.Claims, error) {
	return m.ValidateTokenForMiddlewareFunc(ctx, tokenString)
}

// MockEnrichmentService implementa services.EnrichmentService para pruebas.
type MockEnrichmentService struct {
	EnrichEventFunc      func(ctx context.Context, event *models.Event) ([]models.EnrichedEventRecord, error)
	Top10QueryEventsFunc func(ctx context.Context) ([]models.EnrichedEventRecord, error)
}

func (m *MockEnrichmentService) EnrichEvent(ctx context.Context, event *models.Event) ([]models.EnrichedEventRecord, error) {
	return m.EnrichEventFunc(ctx, event)
}

func (m *MockEnrichmentService) Top10QueryEvents(ctx context.Context) ([]models.EnrichedEventRecord, error) {
	return m.Top10QueryEventsFunc(ctx)
}

// MockJWTService implementa token.JWTService para pruebas.
// Aunque los controladores no lo usan directamente, lo agregamos por si es necesario en otras pruebas.
type MockJWTService struct {
	ExtractJWTTokenFunc  func(r *http.Request) (string, error)
	ValidJWTTokenFunc    func(ctx context.Context, tokenString string) (*models.Claims, error)
	GenerateJWTTokenFunc func(claims *models.Claims) (*models.JWTToken, error)
}

func (m *MockJWTService) ExtractJWTToken(r *http.Request) (string, error) {
	return m.ExtractJWTTokenFunc(r)
}

func (m *MockJWTService) ValidJWTToken(ctx context.Context, tokenString string) (*models.Claims, error) {
	return m.ValidJWTTokenFunc(ctx, tokenString)
}

func (m *MockJWTService) GenerateJWTToken(claims *models.Claims) (*models.JWTToken, error) {
	return m.GenerateJWTTokenFunc(claims)
}
