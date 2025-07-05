package scopes

import (
	"os"
	"strings"
)

func GetTypeScope() string {

	// Solo para pruebas en local (descomentar)
	// os.Setenv("SCOPE", "local")
	os.Setenv("CONFIG_DIR", "./local")
	// os.Setenv("CONFIG_DIR", "./internal/config")
	scope := os.Getenv("SCOPE")
	switch {
	case strings.Contains(scope, "local"):
		return "local"
	case strings.Contains(scope, "test"):
		return "test"
	case strings.Contains(scope, "prod"):
		return "prod"
	default:
		return "unknown"
	}
}
