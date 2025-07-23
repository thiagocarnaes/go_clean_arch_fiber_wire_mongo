package tests

import (
	"context"
	"os"
	"testing"
	"time"
	"user-management/internal/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TestConfig holds configuration for integration tests
type TestConfig struct {
	MongoURI string
	MongoDB  string
	Port     string
}

// NewTestConfig creates a new test configuration
func NewTestConfig() *TestConfig {
	mongoURI := os.Getenv("TEST_MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDB := os.Getenv("TEST_MONGO_DB")
	if mongoDB == "" {
		mongoDB = "user_management_test"
	}

	port := os.Getenv("TEST_PORT")
	if port == "" {
		port = ":3001"
	}

	return &TestConfig{
		MongoURI: mongoURI,
		MongoDB:  mongoDB,
		Port:     port,
	}
}

// ToAppConfig converts TestConfig to application Config
func (tc *TestConfig) ToAppConfig() *config.Config {
	return &config.Config{
		MongoURI: tc.MongoURI,
		MongoDB:  tc.MongoDB,
		Port:     tc.Port,
	}
}

// WaitForMongo waits for MongoDB to be available
func WaitForMongo(uri string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	// Ping the database
	return client.Ping(ctx, nil)
}

// CleanupDatabase removes all data from the test database
func CleanupDatabase(uri, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	return client.Database(dbName).Drop(ctx)
}

// SkipIfNoMongo skips the test if MongoDB is not available
func SkipIfNoMongo(t *testing.T, uri string) {
	if err := WaitForMongo(uri, 5*time.Second); err != nil {
		t.Skipf("MongoDB not available at %s: %v", uri, err)
	}
}
