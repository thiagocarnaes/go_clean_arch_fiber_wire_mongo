package database

import (
	"context"
	"user-management/internal/config"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client    *mongo.Client
	DB        *mongo.Database
	config    *config.Config
	logger    *logrus.Logger
	connected bool
}

func NewMongoDB(cfg *config.Config, log *logrus.Logger) (*MongoDB, error) {
	return &MongoDB{
		config:    cfg,
		logger:    log,
		connected: false,
	}, nil
}

// Connect establishes a connection to MongoDB
func (m *MongoDB) Connect(ctx context.Context) error {
	mongoUri := m.config.MongoURI + "/" + m.config.MongoDB

	client, err := mongo.Connect(options.Client().ApplyURI(mongoUri))
	if err != nil {
		m.logger.WithFields(logrus.Fields{
			"ddsource": m.config.DDSource,
			"service":  m.config.DDService,
			"ddtags":   m.config.DDTags,
			"error":    err.Error(),
		}).Error("Failed to connect to MongoDB")
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		m.logger.WithFields(logrus.Fields{
			"ddsource": m.config.DDSource,
			"service":  m.config.DDService,
			"ddtags":   m.config.DDTags,
			"error":    err.Error(),
		}).Error("MongoDB is not available")
		return err
	}

	m.Client = client
	m.DB = client.Database(m.config.MongoDB)
	m.connected = true

	m.logger.WithFields(logrus.Fields{
		"ddsource": m.config.DDSource,
		"service":  m.config.DDService,
		"ddtags":   m.config.DDTags,
		"uri":      m.config.MongoURI,
	}).Info("Successfully connected to MongoDB")

	return nil
}

// Disconnect closes the MongoDB connection
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.Client == nil {
		return nil
	}

	err := m.Client.Disconnect(ctx)
	if err != nil {
		m.logger.WithFields(logrus.Fields{
			"ddsource": m.config.DDSource,
			"service":  m.config.DDService,
			"ddtags":   m.config.DDTags,
			"error":    err.Error(),
		}).Error("Failed to disconnect from MongoDB")
		return err
	}

	m.connected = false
	m.logger.WithFields(logrus.Fields{
		"ddsource": m.config.DDSource,
		"service":  m.config.DDService,
		"ddtags":   m.config.DDTags,
	}).Info("Disconnected from MongoDB")

	return nil
}

// Ping tests the MongoDB connection
func (m *MongoDB) Ping(ctx context.Context) error {
	if m.Client == nil {
		return mongo.ErrClientDisconnected
	}
	return m.Client.Ping(ctx, nil)
}

// GetConnection returns the MongoDB client
func (m *MongoDB) GetConnection() interface{} {
	return m.Client
}

// GetCollectionConnection returns a MongoDB collection
func (m *MongoDB) GetCollectionConnection(name string) interface{} {
	if m.DB == nil {
		return nil
	}
	return m.DB.Collection(name)
}

// IsConnected checks if MongoDB is connected
func (m *MongoDB) IsConnected() bool {
	return m.connected && m.Client != nil
}

// GetDB returns the MongoDB database (for backward compatibility)
func (m *MongoDB) GetDB() *mongo.Database {
	return m.DB
}

// GetClient returns the MongoDB client (for backward compatibility)
func (m *MongoDB) GetClient() *mongo.Client {
	return m.Client
}
