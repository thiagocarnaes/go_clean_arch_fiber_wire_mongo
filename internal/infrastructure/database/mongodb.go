package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"user-management/internal/config"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDB(cfg *config.Config) (*MongoDB, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.MongoDB)
	return &MongoDB{Client: client, DB: db}, nil
}
