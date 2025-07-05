package utils

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger" // Importa el logger
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// jsonResponse es una estructura auxiliar para enviar respuestas JSON.
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ReadJSON lee el cuerpo de una solicitud HTTP y lo decodifica en la estructura de datos proporcionada.
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// WriteJSON escribe una respuesta JSON al cliente.
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		logger.ErrorLog.Printf("Error al serializar JSON: %v", err)
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		logger.ErrorLog.Printf("Error al escribir respuesta JSON: %v", err)
		return err
	}

	return nil
}

// ErrorJSON envía un error como respuesta JSON.
func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var customErr error

	switch {
	case strings.Contains(err.Error(), "SQLSTATE 23505"):
		customErr = errors.New("valor duplicado viola la restricción única")
		statusCode = http.StatusConflict // 409 Conflict es más apropiado para duplicados
	case strings.Contains(err.Error(), "SQLSTATE 22001"):
		customErr = errors.New("el valor que intenta insertar es demasiado grande")
		statusCode = http.StatusRequestEntityTooLarge // 413 Payload Too Large
	case strings.Contains(err.Error(), "SQLSTATE 23403"):
		customErr = errors.New("violación de clave foránea")
		statusCode = http.StatusConflict
	case errors.Is(err, io.EOF): // Manejo específico para EOF cuando el cuerpo está vacío
		customErr = errors.New("cuerpo de la solicitud vacío")
		statusCode = http.StatusBadRequest
	case strings.Contains(err.Error(), "invalid character"): // Manejo de JSON inválido
		customErr = errors.New("formato JSON inválido")
		statusCode = http.StatusBadRequest
	default:
		customErr = err
	}

	payload := JSONResponse{
		Error:   true,
		Message: customErr.Error(),
	}

	WriteJSON(w, statusCode, payload) // Usar la función WriteJSON
}
