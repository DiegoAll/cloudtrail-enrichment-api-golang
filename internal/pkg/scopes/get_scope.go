package scopes

import (
	"os"
	"strings"
)

func GetTypeScope() string {

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
