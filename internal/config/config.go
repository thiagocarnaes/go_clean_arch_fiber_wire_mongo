package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	MongoURI  string
	MongoDB   string
	Port      string
	DDSource  string
	DDService string
	DDTags    string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		MongoURI:  os.Getenv("MONGO_URI"),
		MongoDB:   os.Getenv("MONGO_DB"),
		Port:      os.Getenv("PORT"),
		DDSource:  os.Getenv("DD_SOURCE"),
		DDService: os.Getenv("DD_SERVICE"),
		DDTags:    os.Getenv("DD_TAGS"),
	}, nil
}
