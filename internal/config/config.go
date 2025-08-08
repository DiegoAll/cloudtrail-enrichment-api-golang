package config

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/scopes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
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

	fmt.Println("Scope:", scopes.GetTypeScope())

	switch scopes.GetTypeScope() {
	case "local":
		// For local execution, read from the relative path of the scaffold_config file
		raw, err = os.ReadFile("./internal/config/scaffold_config")
		if err != nil {
			logger.ErrorLog.Println("[config] error reading local config file", err)
			return nil, fmt.Errorf("[config] error reading local config file: %w", err)
		}
		if err := json.Unmarshal(raw, &config); err != nil {
			logger.ErrorLog.Println("[config] error unmarshaling local configs", err)
			return nil, fmt.Errorf("[config] error unmarshaling local configs: %w", err)
		}
	case "test", "prod":
		// For containerized environments (test/prod), load configuration from environment variables
		config.ServerConfig.Host = os.Getenv("SERVER_HOST")
		config.ServerConfig.Port, _ = strconv.Atoi(os.Getenv("PORT"))
		config.ServerConfig.TLS.Enabled, _ = strconv.ParseBool(os.Getenv("TLS_ENABLED"))
		config.ServerConfig.TLS.CertFile = os.Getenv("TLS_CERT_FILE")
		config.ServerConfig.TLS.KeyFile = os.Getenv("TLS_KEY_FILE")
		config.ServerConfig.CORS.Enabled, _ = strconv.ParseBool(os.Getenv("CORS_ENABLED"))
		config.ServerConfig.CORS.AllowOrigin = os.Getenv("CORS_ALLOW_ORIGIN")
		config.ServerConfig.CORS.AllowMethods = os.Getenv("CORS_ALLOW_METHODS")
		config.ServerConfig.CORS.AllowHeaders = os.Getenv("CORS_ALLOW_HEADERS")

		config.DatabaseConfig.Host = os.Getenv("DATABASE_HOST")
		config.DatabaseConfig.Port, _ = strconv.Atoi(os.Getenv("DATABASE_PORT"))
		config.DatabaseConfig.Username = os.Getenv("DATABASE_USERNAME")
		config.DatabaseConfig.Password = os.Getenv("DATABASE_PASSWORD")
		config.DatabaseConfig.Database = os.Getenv("DATABASE_NAME")
		config.DatabaseConfig.SSLMode = os.Getenv("DATABASE_SSLMODE")
		dbTimeout, _ := strconv.ParseInt(os.Getenv("DB_TIMEOUT"), 10, 64)
		config.DatabaseConfig.DBTimeout = time.Duration(dbTimeout)
		config.DatabaseConfig.MaxOpenConns, _ = strconv.Atoi(os.Getenv("MAX_OPEN_CONNS"))

		config.MongoDBConfig.Host = os.Getenv("MONGO_HOST")
		config.MongoDBConfig.Port, _ = strconv.Atoi(os.Getenv("MONGO_PORT"))
		config.MongoDBConfig.Username = os.Getenv("MONGO_USERNAME")
		config.MongoDBConfig.Password = os.Getenv("MONGO_PASSWORD")
		config.MongoDBConfig.Database = os.Getenv("MONGO_DATABASE")
		config.MongoDBConfig.Collection = os.Getenv("MONGO_COLLECTION")
		mongoDBTimeout, _ := strconv.ParseInt(os.Getenv("MONGO_DB_TIMEOUT"), 10, 64)
		config.MongoDBConfig.DBTimeout = time.Duration(mongoDBTimeout)

		config.AuthConfig.JWTSecret = os.Getenv("JWT_SECRET")
		config.AuthConfig.JWTPrivateKey = os.Getenv("JWT_PRIVATE_KEY")
		config.AuthConfig.JWTPublicKey = os.Getenv("JWT_PUBLIC_KEY")
		hashCost, _ := strconv.Atoi(os.Getenv("HASH_COST"))
		config.AuthConfig.HashCost = hashCost
		tokenDuration, _ := strconv.ParseInt(os.Getenv("TOKEN_DURATION"), 10, 64)
		config.AuthConfig.TokenDuration = time.Duration(tokenDuration)

		logger.InfoLog.Println("Loading configurations from environment variables.")

		//default:
		// the same or choose params.
	}

	// Secrets
	// if scopes.GetTypeScope() == "local" {
	// 	appConfig = tmpConfigs
	// 	// os.Setenv("TRIN_USER", tmpConfigs.UebaConfig.Ueba_user)
	// 	// os.Setenv("TRIN_PASSWORD", tmpConfigs.UebaConfig.Ueba_password)
	// } else {
	// 	fmt.Println("Other HTTP params config")
	// }

	appConfig = config // Assign the loaded configuration to the global variable
	logger.InfoLog.Println("[config] Configs loaded successfully")

	return &appConfig, nil
}

// Configuration getters functions
func GetServerConfig() ServerConfig {
	return appConfig.ServerConfig
}

func GetDatabaseConfig() DatabaseConfig {
	return appConfig.DatabaseConfig
}

func GetMongoDBConfig() MongoDBConfig {
	return appConfig.MongoDBConfig
}

func GetAuthConfig() AuthConfig {
	return appConfig.AuthConfig
}
