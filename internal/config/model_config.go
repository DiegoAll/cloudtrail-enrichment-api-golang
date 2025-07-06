package config

import "time"

type Config struct {
	ServerConfig   ServerConfig   `json:"server_config"`
	DatabaseConfig DatabaseConfig `json:"database_config"`
	MongoDBConfig  MongoDBConfig  `json:"mongodb_config"`
	AuthConfig     AuthConfig     `json:"auth_config"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	TLS  TLS    `json:"tls"`
	CORS CORS   `json:"cors"`
}

type TLS struct {
	Enabled  bool   `json:"enabled"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type CORS struct {
	Enabled      bool   `json:"enabled"`
	AllowOrigin  string `json:"allow_origin"`
	AllowMethods string `json:"allow_methods"`
	AllowHeaders string `json:"allow_headers"`
}

type DatabaseConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	Database     string        `json:"database"`
	SSLMode      string        `json:"ssl_mode"`
	DBTimeout    time.Duration `json:"db_timeout"`
	MaxOpenConns int           `json:"max_open_conns"`
}

type MongoDBConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	Database     string        `json:"database"`
	SSLMode      string        `json:"ssl_mode"`
	DBTimeout    time.Duration `json:"db_timeout"`
	MaxOpenConns int           `json:"max_open_conns"`
}

type AuthConfig struct {
	// Enabled	   bool          `json:"enabled"`
	JWTSecret     string `json:"jwt_secret"`
	JWTPrivateKey string `json:"jwt_private_key"`
	JWTPublicKey  string `json:"jwt_public_key"`
	HashCost      int    `json:"hash_cost"` //store in DB
	// TokenDuration string `json:"token_duration"` // string type
	TokenDuration time.Duration `json:"token_duration"`
}

// rovert
type ConfigLegacy struct {
	Port          int
	CertPath      string
	KeyPath       string
	JWTPrivateKey string
	JWTSecret     string
	HashCost      int
	TokenDuration time.Duration
	DatabaseDSN   string
}
