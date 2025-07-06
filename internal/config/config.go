package config

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/scopes"
	"encoding/json"
	"fmt"
	"os"
)

// var (
// 	productsAPIConfig ProductsConfig
// )

var (
	appConfig Config
)

func LoadConfig() (*Config, error) {

	var (
		config Config
		raw    []byte
		err    error
	)

	switch scopes.GetTypeScope() {
	case "local":
		//config.ServerConfig.Host = "localhost" // ?
		// raw, err = os.ReadFile("./internal/config/scaffold_config")
		raw, err = os.ReadFile("/home/diegoall/Projects/cloudtrail-enrichment-api-golang/internal/config/scaffold_config")
		fmt.Println("RAW", raw)
		// raw, err = os.ReadFile("./local/scaffold_config")
	case "test":
		raw, err = os.ReadFile("scaffold_config")
	case "prod":
		raw, err = os.ReadFile("scaffold_config")
	default:
		raw, err = os.ReadFile("scaffold_config")
	}

	fmt.Println("Scope:", scopes.GetTypeScope())
	if err != nil {
		logger.ErrorLog.Println("[config] error reading config file", err)
		// log.Error(context.Background(), "[config] error reading falcox configs", err)
		return nil, err
	}

	var tmpConfigs Config
	// if err := json.Unmarshal(raw, &tmpConfigs); err != nil {
	if err := json.Unmarshal(raw, &config); err != nil {
		logger.ErrorLog.Println("[config] error unmarshaling configs", err)
		return nil, err
	}

	// Secrets
	if scopes.GetTypeScope() == "local" {
		appConfig = tmpConfigs
		// os.Setenv("TRIN_USER", tmpConfigs.UebaConfig.Ueba_user)
		// os.Setenv("TRIN_PASSWORD", tmpConfigs.UebaConfig.Ueba_password)
	} else {
		fmt.Println("Otra config HTTP params")
	}

	appConfig = config

	logger.InfoLog.Println("[config] Configs loaded successfully")

	return &appConfig, nil
}

// Funciones de acceso a la configuración

func GetServerConfig() ServerConfig {
	return appConfig.ServerConfig
}

func GetDatabaseConfig() DatabaseConfig {
	return appConfig.DatabaseConfig
}

func GetAuthConfig() AuthConfig {
	return appConfig.AuthConfig
}
