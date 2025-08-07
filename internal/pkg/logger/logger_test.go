package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestInit verifica que los loggers se inicialicen correctamente
// y que los prefijos sean los esperados.
func TestInit(t *testing.T) {
	// Llamar a la función Init() para inicializar los loggers
	Init()

	// Redirigir la salida estándar a un buffer para capturar los logs
	var buf bytes.Buffer
	InfoLog.SetOutput(&buf)
	ErrorLog.SetOutput(&buf)
	DebugLog.SetOutput(&buf)

	// Probar que los loggers escriben en el buffer con los prefijos correctos
	buf.Reset()
	InfoLog.Println("Test de InfoLog")
	if !strings.Contains(buf.String(), "INFO\t") {
		t.Error("InfoLog no escribió con el prefijo correcto")
	}

	buf.Reset()
	ErrorLog.Println("Test de ErrorLog")
	if !strings.Contains(buf.String(), "ERROR\t") {
		t.Error("ErrorLog no escribió con el prefijo correcto")
	}

	buf.Reset()
	DebugLog.Println("Test de DebugLog")
	if !strings.Contains(buf.String(), "DEBUG\t") {
		t.Error("DebugLog no escribió con el prefijo correcto")
	}

	// Restaurar los loggers a la salida estándar original para evitar interferencia en otras pruebas
	InfoLog.SetOutput(os.Stderr)
	ErrorLog.SetOutput(os.Stderr)
	DebugLog.SetOutput(os.Stderr)
}
