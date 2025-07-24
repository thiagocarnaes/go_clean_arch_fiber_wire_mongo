package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI     string
	MongoDB      string
	Port         string
	DDSource     string
	DDService    string
	DDTags       string
	DatabaseType string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		MongoURI:     os.Getenv("MONGO_URI"),
		MongoDB:      os.Getenv("MONGO_DB"),
		Port:         os.Getenv("PORT"),
		DDSource:     os.Getenv("DD_SOURCE"),
		DDService:    os.Getenv("DD_SERVICE"),
		DDTags:       os.Getenv("DD_TAGS"),
		DatabaseType: getEnvOrDefault("DATABASE_TYPE", "mongodb"),
	}, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
