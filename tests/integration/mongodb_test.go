package integration

import (
	"context"
	"testing"
	"time"
	"user-management/internal/config"
	"user-management/internal/infrastructure/database"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	testInvalidPassword = "wrongpasswd" // Test password for authentication failure scenarios
	testMongoURI        = "mongodb://localhost:27017"
	testService         = "test-service"
	testTags            = "env:test"
)

func TestNewMongoDBSuccess(t *testing.T) {
	// Setup test config with valid MongoDB connection
	cfg := &config.Config{
		MongoURI:  testMongoURI,
		MongoDB:   "test_db",
		DDSource:  "test",
		DDService: testService,
		DDTags:    testTags,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests

	// Test successful connection (requires MongoDB running)
	// Skip if MongoDB is not available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(nil)
	if err != nil {
		t.Skip("MongoDB not available for integration test")
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		t.Skip("MongoDB not responding for integration test")
		return
	}
	client.Disconnect(ctx)

	// Now test our NewMongoDB function
	mongodb, err := database.NewMongoDB(cfg, logger)

	require.NoError(t, err)
	assert.NotNil(t, mongodb)
	assert.NotNil(t, mongodb.Client)
	assert.NotNil(t, mongodb.DB)
	assert.Equal(t, "test_db", mongodb.DB.Name())

	// Cleanup
	mongodb.Client.Disconnect(context.Background())
}

func TestNewMongoDBInvalidURI(t *testing.T) {
	// Test with invalid MongoDB URI
	cfg := &config.Config{
		MongoURI:  "mongodb://invalid-host:99999",
		MongoDB:   "test_db",
		DDSource:  "test",
		DDService: testService,
		DDTags:    testTags,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress error logs for this test

	mongodb, err := database.NewMongoDB(cfg, logger)

	assert.Error(t, err)
	assert.Nil(t, mongodb)
	// Check for any connection-related error message
	errorMsg := err.Error()
	assert.True(t,
		len(errorMsg) > 0,
		"Expected an error message but got empty string")
}

func TestNewMongoDBConnectionRefused(t *testing.T) {
	// Test with URI pointing to non-existent MongoDB instance
	cfg := &config.Config{
		MongoURI:  "mongodb://localhost:27099", // Non-existent port
		MongoDB:   "test_db",
		DDSource:  "test",
		DDService: testService,
		DDTags:    testTags,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress error logs for this test

	mongodb, err := database.NewMongoDB(cfg, logger)

	assert.Error(t, err)
	assert.Nil(t, mongodb)
}

func TestNewMongoDBEmptyDatabaseName(t *testing.T) {
	// Test with empty database name
	cfg := &config.Config{
		MongoURI:  testMongoURI,
		MongoDB:   "", // Empty database name
		DDSource:  "test",
		DDService: testService,
		DDTags:    testTags,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	// Skip if MongoDB is not available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(nil)
	if err != nil {
		t.Skip("MongoDB not available for integration test")
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		t.Skip("MongoDB not responding for integration test")
		return
	}
	client.Disconnect(ctx)

	mongodb, err := database.NewMongoDB(cfg, logger)

	// Should succeed but with empty database name
	require.NoError(t, err)
	assert.NotNil(t, mongodb)
	assert.Equal(t, "", mongodb.DB.Name())

	// Cleanup
	mongodb.Client.Disconnect(context.Background())
}

func TestNewMongoDBMalformedURI(t *testing.T) {
	// Test with malformed MongoDB URI
	cfg := &config.Config{
		MongoURI:  "not-a-valid-uri",
		MongoDB:   "test_db",
		DDSource:  "test",
		DDService: testService,
		DDTags:    testTags,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	mongodb, err := database.NewMongoDB(cfg, logger)

	assert.Error(t, err)
	assert.Nil(t, mongodb)
	assert.Contains(t, err.Error(), "error parsing uri")
}

func TestNewMongoDBWithAuthInvalidCredentials(t *testing.T) {
	// Test with invalid authentication credentials
	cfg := &config.Config{
		MongoURI:  "mongodb://testuser:" + testInvalidPassword + "@localhost:27017", // Test credentials that will fail
		MongoDB:   "test_db",
		DDSource:  "test",
		DDService: testService,
		DDTags:    testTags,
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	mongodb, err := database.NewMongoDB(cfg, logger)

	assert.Error(t, err)
	assert.Nil(t, mongodb)
}

func TestNewMongoDBLoggingFields(t *testing.T) {
	// Test that logging fields are properly used
	cfg := &config.Config{
		MongoURI:  "mongodb://localhost:27099", // Will fail to connect
		MongoDB:   "test_db",
		DDSource:  "test-source",
		DDService: testService,
		DDTags:    "env:test,version:1.0",
	}

	// Create a custom logger to capture log entries
	logger := logrus.New()

	// Use a hook to capture log entries
	var logEntries []*logrus.Entry
	logger.AddHook(&testLogHook{entries: &logEntries})

	mongodb, err := database.NewMongoDB(cfg, logger)

	assert.Error(t, err)
	assert.Nil(t, mongodb)

	// Verify that error was logged with correct fields
	require.NotEmpty(t, logEntries)
	errorEntry := logEntries[len(logEntries)-1] // Get last log entry

	assert.Equal(t, logrus.ErrorLevel, errorEntry.Level)
	// The error can be either "Failed to connect to MongoDB" or "MongoDB is not available"
	// depending on whether the error occurs during connection or ping
	errorMessage := errorEntry.Message
	assert.True(t,
		errorMessage == "Failed to connect to MongoDB" || errorMessage == "MongoDB is not available",
		"Expected 'Failed to connect to MongoDB' or 'MongoDB is not available', got: %s", errorMessage)
	assert.Equal(t, "test-source", errorEntry.Data["ddsource"])
	assert.Equal(t, testService, errorEntry.Data["service"])
	assert.Equal(t, "env:test,version:1.0", errorEntry.Data["ddtags"])
	assert.NotEmpty(t, errorEntry.Data["error"], "Error field should not be empty")
}

// Custom hook to capture log entries for testing
type testLogHook struct {
	entries *[]*logrus.Entry
}

func (hook *testLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *testLogHook) Fire(entry *logrus.Entry) error {
	*hook.entries = append(*hook.entries, entry)
	return nil
}

// Benchmark test for MongoDB connection
func BenchmarkNewMongoDB(b *testing.B) {
	cfg := &config.Config{
		MongoURI:  testMongoURI,
		MongoDB:   "benchmark_db",
		DDSource:  "benchmark",
		DDService: "benchmark-service",
		DDTags:    "env:benchmark",
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress logs during benchmark

	// Skip if MongoDB is not available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(nil)
	if err != nil {
		b.Skip("MongoDB not available for benchmark test")
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		b.Skip("MongoDB not responding for benchmark test")
		return
	}
	client.Disconnect(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mongodb, err := database.NewMongoDB(cfg, logger)
		if err != nil {
			b.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		mongodb.Client.Disconnect(context.Background())
	}
}
