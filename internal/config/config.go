package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	MongoURI string
	MongoDB  string
	Port     string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		MongoURI: os.Getenv("MONGO_URI"),
		MongoDB:  os.Getenv("MONGO_DB"),
		Port:     os.Getenv("PORT"),
	}, nil
}
