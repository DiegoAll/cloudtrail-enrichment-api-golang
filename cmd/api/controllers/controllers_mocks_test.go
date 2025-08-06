package controllers

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/token"
	"cloudtrail-enrichment-api-golang/models"
	"context"
)

// MockAuthService para el controlador de autenticación.
type MockAuthService struct {
	RegisterUserFunc     func(ctx context.Context, payload *models.RegisterPayload) (*models.User, error)
	AuthenticateUserFunc func(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error)
	// CORRECCIÓN: Se renombra el campo para evitar el conflicto con el método.
	ValidateTokenForMiddlewareFunc func(ctx context.Context, tokenString string) (*token.User, error)
}

func (m *MockAuthService) RegisterUser(ctx context.Context, payload *models.RegisterPayload) (*models.User, error) {
	if m.RegisterUserFunc != nil {
		return m.RegisterUserFunc(ctx, payload)
	}
	return nil, nil
}

func (m *MockAuthService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error) {
	if m.AuthenticateUserFunc != nil {
		return m.AuthenticateUserFunc(ctx, email, password)
	}
	return nil, nil, nil
}

func (m *MockAuthService) ValidateTokenForMiddleware(ctx context.Context, tokenString string) (*token.User, error) {
	if m.ValidateTokenForMiddlewareFunc != nil {
		return m.ValidateTokenForMiddlewareFunc(ctx, tokenString)
	}
	return nil, nil
}

// MockEnrichmentService es una implementación mock de la interfaz services.EnrichmentService.
type MockEnrichmentService struct {
	EnrichEventFunc      func(ctx context.Context, event *models.Event) ([]*models.EnrichedEventRecord, error)
	Top10QueryEventsFunc func(ctx context.Context) ([]*models.EnrichedEventRecord, error)
}

func (m *MockEnrichmentService) EnrichEvent(ctx context.Context, event *models.Event) ([]*models.EnrichedEventRecord, error) {
	if m.EnrichEventFunc != nil {
		return m.EnrichEventFunc(ctx, event)
	}
	return nil, nil
}

func (m *MockEnrichmentService) Top10QueryEvents(ctx context.Context) ([]*models.EnrichedEventRecord, error) {
	if m.Top10QueryEventsFunc != nil {
		return m.Top10QueryEventsFunc(ctx)
	}
	return nil, nil
}
