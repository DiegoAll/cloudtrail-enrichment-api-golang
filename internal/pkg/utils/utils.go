package utils

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// JSONResponse is a helper struct to send JSON responses.
type JSONResponse struct {
	Error   bool        `json:"error"`          // Indicates if there was an error
	Message string      `json:"message"`        // Message associated with the response
	Data    interface{} `json:"data,omitempty"` // Optional data
}

// ReadJSON reads the body of an HTTP request and decodes it into the provided data structure.
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain only a single JSON value")
	}

	return nil
}

// WriteJSON writes a JSON response to the client.
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		logger.ErrorLog.Printf("Error serializing JSON: %v", err)
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
		logger.ErrorLog.Printf("Error writing JSON response: %v", err)
		return err
	}

	return nil
}

// ErrorJSON sends an error as a JSON response.
func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var customErr error

	switch {
	case strings.Contains(err.Error(), "SQLSTATE 23505"):
		customErr = errors.New("duplicate value violates unique constraint")
		statusCode = http.StatusConflict // 409 Conflict is more appropriate for duplicates
	case strings.Contains(err.Error(), "SQLSTATE 22001"):
		customErr = errors.New("the value you are trying to insert is too large")
		statusCode = http.StatusRequestEntityTooLarge // 413 Payload Too Large
	case strings.Contains(err.Error(), "SQLSTATE 23403"):
		customErr = errors.New("foreign key violation")
		statusCode = http.StatusConflict
	case errors.Is(err, io.EOF): // Specific handling for EOF when the body is empty
		customErr = errors.New("request body is empty")
		statusCode = http.StatusBadRequest
	case strings.Contains(err.Error(), "invalid character"): // Handling invalid JSON
		customErr = errors.New("invalid JSON format")
		statusCode = http.StatusBadRequest
	default:
		customErr = err
	}

	payload := JSONResponse{
		Error:   true,
		Message: customErr.Error(),
	}

	WriteJSON(w, statusCode, payload)
}
