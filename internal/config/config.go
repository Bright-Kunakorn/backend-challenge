package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds application configuration values.
type Config struct {
	Port           string
	GRPCPort       string
	MongoURI       string
	MongoDatabase  string
	JWTSecret      string
	JWTIssuer      string
	JWTExpiry      time.Duration
	BackgroundTick time.Duration
	Environment    string
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	cfg := Config{
		Port:           getEnv("PORT", "8080"),
		GRPCPort:       getEnv("GRPC_PORT", "50051"),
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:  getEnv("MONGO_DB", "user_service"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTIssuer:      getEnv("JWT_ISSUER", "backend-challenge"),
		JWTExpiry:      parseDuration(getEnv("JWT_EXPIRY", "24h"), 24*time.Hour),
		BackgroundTick: parseDuration(getEnv("USER_COUNT_TICK", "10s"), 10*time.Second),
		Environment:    getEnv("ENVIRONMENT", "development"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET must be provided")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseDuration(value string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return d
}

// MustParseInt reads an integer environment variable.
func MustParseInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return v
}
