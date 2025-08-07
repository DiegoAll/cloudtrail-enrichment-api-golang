package utils

import (
	"bytes"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestMain se ejecuta antes de cualquier función de prueba en este paquete.
// Es ideal para configurar y limpiar recursos compartidos, como la inicialización de loggers.
func TestMain(m *testing.M) {
	// Inicializamos los loggers para que no sean nil durante las pruebas.
	// Esto es crucial, ya que las funciones de utils dependen de ellos.
	logger.InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.DebugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Ejecuta todos los tests en este paquete
	exitCode := m.Run()

	// Termina el proceso con el código de salida de las pruebas
	os.Exit(exitCode)
}

// TestReadJSON_Success prueba la lectura exitosa de una solicitud JSON.
func TestReadJSON_Success(t *testing.T) {
	// Estructura de datos de prueba
	type testData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Crea un body de solicitud JSON
	payload := `{"name": "John Doe", "age": 30}`
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	var td testData
	err := ReadJSON(rr, req, &td)

	if err != nil {
		t.Fatalf("ReadJSON retornó un error inesperado: %v", err)
	}

	if td.Name != "John Doe" {
		t.Errorf("Nombre incorrecto. Se esperaba 'John Doe', se obtuvo '%s'", td.Name)
	}
	if td.Age != 30 {
		t.Errorf("Edad incorrecta. Se esperaba 30, se obtuvo %d", td.Age)
	}
}

// TestReadJSON_BodyTooLarge prueba un error cuando el body es demasiado grande.
func TestReadJSON_BodyTooLarge(t *testing.T) {
	// Genera un body de 2MB para simular un cuerpo demasiado grande
	payload := bytes.Repeat([]byte("a"), 2*1024*1024)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	var data interface{}
	err := ReadJSON(rr, req, &data)

	if err == nil {
		t.Fatal("ReadJSON debería haber retornado un error para un cuerpo demasiado grande, pero retornó nil")
	}
	// Se espera que el error contenga "too large" o un error de JSON inválido
	if !strings.Contains(err.Error(), "http: request body too large") && !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("Error incorrecto para body grande. Se esperaba 'too large' o 'invalid character', se obtuvo: %v", err)
	}
}

// TestReadJSON_MultipleJSONValues prueba un error cuando hay múltiples valores JSON.
func TestReadJSON_MultipleJSONValues(t *testing.T) {
	payload := `{"name": "John Doe"}{"name": "Jane Doe"}`
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	var data map[string]interface{}
	err := ReadJSON(rr, req, &data)

	if err == nil {
		t.Fatal("ReadJSON debería haber retornado un error para múltiples valores JSON, pero retornó nil")
	}
	if err.Error() != "body must have only a single JSON value" {
		t.Errorf("Error incorrecto. Se esperaba 'body must have only a single JSON value', se obtuvo: %v", err)
	}
}

// TestWriteJSON_Success prueba la escritura exitosa de una respuesta JSON.
func TestWriteJSON_Success(t *testing.T) {
	data := JSONResponse{
		Error:   false,
		Message: "success",
		Data:    "some data",
	}

	rr := httptest.NewRecorder()

	err := WriteJSON(rr, http.StatusOK, data)
	if err != nil {
		t.Fatalf("WriteJSON retornó un error inesperado: %v", err)
	}

	// Verifica el status code
	if rr.Code != http.StatusOK {
		t.Errorf("Status code incorrecto. Se esperaba %d, se obtuvo %d", http.StatusOK, rr.Code)
	}

	// Verifica el Content-Type
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type incorrecto. Se esperaba 'application/json', se obtuvo '%s'", rr.Header().Get("Content-Type"))
	}

	// Verifica el cuerpo de la respuesta, sin el carácter de nueva línea
	expectedBody := `{"error":false,"message":"success","data":"some data"}`
	if strings.TrimSpace(rr.Body.String()) != expectedBody {
		t.Errorf("Cuerpo de respuesta incorrecto. Se esperaba %q, se obtuvo %q", expectedBody, rr.Body.String())
	}
}

// TestErrorJSON_GenericError prueba el manejo de un error genérico.
func TestErrorJSON_GenericError(t *testing.T) {
	rr := httptest.NewRecorder()
	err := errors.New("error genérico de prueba")

	ErrorJSON(rr, err)

	// Verifica el status code por defecto
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status code incorrecto. Se esperaba %d, se obtuvo %d", http.StatusBadRequest, rr.Code)
	}

	// Verifica el cuerpo de la respuesta JSON
	var payload JSONResponse
	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("Error al decodificar la respuesta JSON: %v", err)
	}

	if !payload.Error {
		t.Error("Se esperaba que el campo 'error' fuera true")
	}
	if payload.Message != "error genérico de prueba" {
		t.Errorf("Mensaje de error incorrecto. Se esperaba 'error genérico de prueba', se obtuvo '%s'", payload.Message)
	}
}

// TestErrorJSON_DuplicateValueError prueba el manejo de un error de valor duplicado.
func TestErrorJSON_DuplicateValueError(t *testing.T) {
	rr := httptest.NewRecorder()
	err := errors.New("duplicate key value violates unique constraint \"unique_constraint_name\" (SQLSTATE 23505)")

	ErrorJSON(rr, err)

	if rr.Code != http.StatusConflict {
		t.Errorf("Status code incorrecto. Se esperaba %d, se obtuvo %d", http.StatusConflict, rr.Code)
	}

	var payload JSONResponse
	_ = json.NewDecoder(rr.Body).Decode(&payload)

	if payload.Message != "valor duplicado viola la restricción única" {
		t.Errorf("Mensaje de error incorrecto. Se esperaba 'valor duplicado viola la restricción única', se obtuvo '%s'", payload.Message)
	}
}

// TestErrorJSON_EmptyBodyError prueba el manejo de un error de cuerpo de solicitud vacío.
func TestErrorJSON_EmptyBodyError(t *testing.T) {
	rr := httptest.NewRecorder()
	err := io.EOF

	ErrorJSON(rr, err)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status code incorrecto. Se esperaba %d, se obtuvo %d", http.StatusBadRequest, rr.Code)
	}

	var payload JSONResponse
	_ = json.NewDecoder(rr.Body).Decode(&payload)

	if payload.Message != "cuerpo de la solicitud vacío" {
		t.Errorf("Mensaje de error incorrecto. Se esperaba 'cuerpo de la solicitud vacío', se obtuvo '%s'", payload.Message)
	}
}

// TestErrorJSON_InvalidJSONError prueba el manejo de un error de formato JSON inválido.
func TestErrorJSON_InvalidJSONError(t *testing.T) {
	rr := httptest.NewRecorder()
	err := errors.New("invalid character '}' after array element")

	ErrorJSON(rr, err)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Status code incorrecto. Se esperaba %d, se obtuvo %d", http.StatusBadRequest, rr.Code)
	}

	var payload JSONResponse
	_ = json.NewDecoder(rr.Body).Decode(&payload)

	if payload.Message != "formato JSON inválido" {
		t.Errorf("Mensaje de error incorrecto. Se esperaba 'formato JSON inválido', se obtuvo '%s'", payload.Message)
	}
}
