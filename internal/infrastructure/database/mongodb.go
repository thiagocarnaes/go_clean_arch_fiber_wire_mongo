package database

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
	"user-management/internal/config"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDB(cfg *config.Config, log *logrus.Logger) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoUri := cfg.MongoURI + "/" + cfg.MongoDB

	client, err := mongo.Connect(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.WithFields(logrus.Fields{
			"ddsource": "go",
			"service":  "user-management",
			"ddtags":   "env:dev,app:fiber",
			"error":    err.Error(),
		}).Error("Failed to connect to MongoDB")
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"ddsource": "go",
			"service":  "user-management",
			"ddtags":   "env:dev,app:fiber",
			"error":    err.Error(),
		}).Error("MongoDB is not available")
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"ddsource": "go",
		"service":  "user-management",
		"ddtags":   "env:dev,app:fiber",
		"uri":      cfg.MongoURI,
	}).Info("Successfully connected to MongoDB")

	db := client.Database(cfg.MongoDB)
	return &MongoDB{Client: client, DB: db}, nil
}
