package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"user-management/internal/config"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TestConfig holds configuration for integration tests
type TestConfig struct {
	MongoURI         string
	MongoDB          string
	Port             string
	MongoContainer   *mongodb.MongoDBContainer
	UseTestContainer bool
}

// NewTestConfig creates a new test configuration
func NewTestConfig() *TestConfig {
	useTestContainer := os.Getenv("USE_TEST_CONTAINER")
	if useTestContainer == "" {
		useTestContainer = "true" // Por padrão, usar testcontainer
	}

	mongoURI := os.Getenv("TEST_MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDB := os.Getenv("TEST_MONGO_DB")
	if mongoDB == "" {
		mongoDB = "user_management_test"
	}

	// Get a free port for the server
	freePort, err := getFreePort()
	if err != nil {
		panic(fmt.Sprintf("failed to get free port: %v", err))
	}
	port := fmt.Sprintf(":%d", freePort)

	return &TestConfig{
		MongoURI:         mongoURI,
		MongoDB:          mongoDB,
		Port:             port,
		UseTestContainer: useTestContainer == "true",
	}
}

// ToAppConfig converts TestConfig to application Config
func (tc *TestConfig) ToAppConfig() *config.Config {
	return &config.Config{
		MongoURI:     tc.MongoURI,
		MongoDB:      tc.MongoDB,
		Port:         tc.Port,
		DDSource:     "go",
		DDService:    "user-management",
		DDTags:       "env:test,app:fiber",
		DatabaseType: "mongodb",
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

// StartMongoContainer starts a MongoDB test container
func (tc *TestConfig) StartMongoContainer(ctx context.Context) error {
	if !tc.UseTestContainer {
		return nil // Use external MongoDB
	}

	container, err := mongodb.Run(ctx, "mongo:7.0")
	if err != nil {
		return fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	tc.MongoContainer = container

	// Get the connection string
	connectionString, err := container.ConnectionString(ctx)
	if err != nil {
		return fmt.Errorf("failed to get connection string: %w", err)
	}

	tc.MongoURI = connectionString
	return nil
}

// StopMongoContainer stops the MongoDB test container
func (tc *TestConfig) StopMongoContainer(ctx context.Context) error {
	if tc.MongoContainer != nil {
		return testcontainers.TerminateContainer(tc.MongoContainer)
	}
	return nil
}

// SetupTestContainer initializes MongoDB test container if needed
func SetupTestContainer(ctx context.Context) (*TestConfig, error) {
	cfg := NewTestConfig()

	if err := cfg.StartMongoContainer(ctx); err != nil {
		return nil, err
	}

	// Wait for MongoDB to be ready
	if err := WaitForMongo(cfg.MongoURI, 30*time.Second); err != nil {
		cfg.StopMongoContainer(ctx)
		return nil, fmt.Errorf("MongoDB container not ready: %w", err)
	}

	return cfg, nil
}
