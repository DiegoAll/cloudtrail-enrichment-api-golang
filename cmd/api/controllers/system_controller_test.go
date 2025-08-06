// cmd/api/controllers/system_controller_test.go

package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheck prueba que el health check retorna un estado OK.
func TestHealthCheck(t *testing.T) {
	controller := NewSystemController()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	controller.HealthCheck(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Código de estado esperado %d, pero se obtuvo %d", http.StatusOK, rr.Code)
	}

	// Opcional: Podrías verificar que el cuerpo de la respuesta contiene "API está operativa"
}
