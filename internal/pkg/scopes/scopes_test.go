package scopes

import (
	"os"
	"testing"
)

// TestGetTypeScope_Local comprueba que GetTypeScope devuelve "local" cuando la variable
// de entorno SCOPE contiene la palabra "local".
func TestGetTypeScope_Local(t *testing.T) {
	// Configura la variable de entorno SCOPE para la prueba.
	os.Setenv("SCOPE", "local")
	// Asegura que la variable de entorno se limpie al finalizar la prueba.
	defer os.Unsetenv("SCOPE")

	// Llama a la función que se está probando.
	result := GetTypeScope()

	// Comprueba si el resultado es el esperado.
	expected := "local"
	if result != expected {
		t.Errorf("GetTypeScope() devolvió %s, se esperaba %s", result, expected)
	}
}

// TestGetTypeScope_Test comprueba que GetTypeScope devuelve "test" cuando la variable
// de entorno SCOPE contiene la palabra "test".
func TestGetTypeScope_Test(t *testing.T) {
	// Configura la variable de entorno SCOPE.
	os.Setenv("SCOPE", "testing_environment")
	// Limpia la variable de entorno al finalizar.
	defer os.Unsetenv("SCOPE")

	// Llama a la función y verifica el resultado.
	result := GetTypeScope()

	expected := "test"
	if result != expected {
		t.Errorf("GetTypeScope() devolvió %s, se esperaba %s", result, expected)
	}
}

// TestGetTypeScope_Prod comprueba que GetTypeScope devuelve "prod" cuando la variable
// de entorno SCOPE contiene la palabra "prod".
func TestGetTypeScope_Prod(t *testing.T) {
	// Configura la variable de entorno SCOPE.
	os.Setenv("SCOPE", "production")
	// Limpia la variable de entorno al finalizar.
	defer os.Unsetenv("SCOPE")

	// Llama a la función y verifica el resultado.
	result := GetTypeScope()

	expected := "prod"
	if result != expected {
		t.Errorf("GetTypeScope() devolvió %s, se esperaba %s", result, expected)
	}
}

// TestGetTypeScope_Unknown comprueba que GetTypeScope devuelve "unknown" cuando la variable
// de entorno SCOPE no contiene ninguna de las palabras clave.
func TestGetTypeScope_Unknown(t *testing.T) {
	// Configura la variable de entorno SCOPE con un valor no reconocido.
	os.Setenv("SCOPE", "staging")
	// Limpia la variable de entorno al finalizar.
	defer os.Unsetenv("SCOPE")

	// Llama a la función y verifica el resultado.
	result := GetTypeScope()

	expected := "unknown"
	if result != expected {
		t.Errorf("GetTypeScope() devolvió %s, se esperaba %s", result, expected)
	}
}

// TestGetTypeScope_Empty comprueba que GetTypeScope devuelve "unknown" cuando la variable
// de entorno SCOPE no está definida.
func TestGetTypeScope_Empty(t *testing.T) {
	// Asegúrate de que la variable de entorno no esté configurada.
	os.Unsetenv("SCOPE")

	// Llama a la función y verifica el resultado.
	result := GetTypeScope()

	expected := "unknown"
	if result != expected {
		t.Errorf("GetTypeScope() devolvió %s, se esperaba %s", result, expected)
	}
}
