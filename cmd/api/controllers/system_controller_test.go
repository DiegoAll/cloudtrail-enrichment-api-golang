package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheck Tests that the health check returns an OK status.
func TestHealthCheck(t *testing.T) {
	controller := NewSystemController()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	controller.HealthCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusOK, rr.Code)
	}
}
