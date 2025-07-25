package units

import (
	"os"
	"path/filepath"
	"testing"

	"user-management/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_Success(t *testing.T) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Set up test environment variables
	testEnvVars := map[string]string{
		"MONGO_URI":     "mongodb://localhost:27017",
		"MONGO_DB":      "test_db",
		"PORT":          "8080",
		"DD_SOURCE":     "test_source",
		"DD_SERVICE":    "test_service",
		"DD_TAGS":       "env:test,version:1.0",
		"DATABASE_TYPE": "mongodb",
	}

	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Create a temporary .env file for testing
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	// Change to temp directory to test .env loading
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create .env file
	envContent := `MONGO_URI=mongodb://localhost:27017
MONGO_DB=test_db
PORT=8080
DD_SOURCE=test_source
DD_SERVICE=test_service
DD_TAGS=env:test,version:1.0
DATABASE_TYPE=mongodb`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	// Test NewConfig
	cfg, err := config.NewConfig()

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoURI)
	assert.Equal(t, "test_db", cfg.MongoDB)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "test_source", cfg.DDSource)
	assert.Equal(t, "test_service", cfg.DDService)
	assert.Equal(t, "env:test,version:1.0", cfg.DDTags)
	assert.Equal(t, "mongodb", cfg.DatabaseType)
}

func TestNewConfig_WithoutEnvFile(t *testing.T) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Set up test environment variables directly (no .env file)
	testEnvVars := map[string]string{
		"MONGO_URI":     "mongodb://prod:27017",
		"MONGO_DB":      "prod_db",
		"PORT":          "3000",
		"DD_SOURCE":     "prod_source",
		"DD_SERVICE":    "prod_service",
		"DD_TAGS":       "env:prod",
		"DATABASE_TYPE": "mongodb",
	}

	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Create temp directory without .env file
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Test NewConfig (should fail to load .env but still work with env vars)
	cfg, err := config.NewConfig()

	// Should return error because .env file doesn't exist
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestNewConfig_DefaultDatabaseType(t *testing.T) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Set up test environment variables without DATABASE_TYPE
	testEnvVars := map[string]string{
		"MONGO_URI":  "mongodb://localhost:27017",
		"MONGO_DB":   "test_db",
		"PORT":       "8080",
		"DD_SOURCE":  "test_source",
		"DD_SERVICE": "test_service",
		"DD_TAGS":    "env:test",
	}

	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Create a temporary .env file for testing
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create .env file without DATABASE_TYPE
	envContent := `MONGO_URI=mongodb://localhost:27017
MONGO_DB=test_db
PORT=8080
DD_SOURCE=test_source
DD_SERVICE=test_service
DD_TAGS=env:test`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	// Test NewConfig
	cfg, err := config.NewConfig()

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "mongodb", cfg.DatabaseType) // Should use default value
}

func TestNewConfig_EmptyEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Clear all relevant environment variables
	envVars := []string{
		"MONGO_URI", "MONGO_DB", "PORT", "DD_SOURCE",
		"DD_SERVICE", "DD_TAGS", "DATABASE_TYPE",
	}

	for _, key := range envVars {
		os.Unsetenv(key)
	}

	// Create a temporary empty .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create empty .env file
	err := os.WriteFile(envFile, []byte(""), 0644)
	require.NoError(t, err)

	// Test NewConfig
	cfg, err := config.NewConfig()

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.MongoURI)
	assert.Equal(t, "", cfg.MongoDB)
	assert.Equal(t, "", cfg.Port)
	assert.Equal(t, "", cfg.DDSource)
	assert.Equal(t, "", cfg.DDService)
	assert.Equal(t, "", cfg.DDTags)
	assert.Equal(t, "mongodb", cfg.DatabaseType) // Should use default value
}

func TestNewConfig_MixedEnvironmentSources(t *testing.T) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Clear all relevant environment variables first
	envVars := []string{
		"MONGO_URI", "MONGO_DB", "PORT", "DD_SOURCE",
		"DD_SERVICE", "DD_TAGS", "DATABASE_TYPE",
	}

	for _, key := range envVars {
		os.Unsetenv(key)
	}

	// Set some environment variables directly after clearing
	os.Setenv("PORT", "9000")

	// Create a temporary .env file with values
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create .env file
	envContent := `MONGO_URI=mongodb://dotenv:27017
MONGO_DB=dotenv_db
DD_SOURCE=dotenv_source
DATABASE_TYPE=postgresql`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	// Test NewConfig
	cfg, err := config.NewConfig()

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	// Values from .env file
	assert.Equal(t, "mongodb://dotenv:27017", cfg.MongoURI)
	assert.Equal(t, "dotenv_db", cfg.MongoDB)
	assert.Equal(t, "dotenv_source", cfg.DDSource)
	assert.Equal(t, "postgresql", cfg.DatabaseType)
	// Value from direct environment variable (set after .env loading)
	assert.Equal(t, "9000", cfg.Port)
	// Not set anywhere
	assert.Equal(t, "", cfg.DDService)
	assert.Equal(t, "", cfg.DDTags)
}

func TestNewConfig_DatabaseTypeWithEmptyValue(t *testing.T) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Set DATABASE_TYPE to empty string
	os.Setenv("DATABASE_TYPE", "")

	// Create a temporary .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create .env file with empty DATABASE_TYPE
	envContent := `DATABASE_TYPE=`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	// Test NewConfig
	cfg, err := config.NewConfig()

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "mongodb", cfg.DatabaseType) // Should use default value when empty
}

func TestNewConfig_DatabaseTypeWithCustomValue(t *testing.T) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Create a temporary .env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create .env file with custom DATABASE_TYPE
	envContent := `DATABASE_TYPE=postgresql`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(t, err)

	// Test NewConfig
	cfg, err := config.NewConfig()

	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgresql", cfg.DatabaseType)
}

func TestNewConfig_InvalidEnvFile(t *testing.T) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Create a temporary directory with a directory named .env (not a file)
	tempDir := t.TempDir()
	envDir := filepath.Join(tempDir, ".env")

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	// Create .env as a directory instead of a file
	err := os.Mkdir(envDir, 0755)
	require.NoError(t, err)

	// Test NewConfig (should fail to load .env)
	cfg, err := config.NewConfig()

	// Should return error because .env is a directory, not a file
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

// Benchmark test for NewConfig
func BenchmarkNewConfig(b *testing.B) {
	// Save original environment
	originalEnv := saveEnvironment()
	defer restoreEnvironment(originalEnv)

	// Set up test environment variables
	testEnvVars := map[string]string{
		"MONGO_URI":     "mongodb://localhost:27017",
		"MONGO_DB":      "test_db",
		"PORT":          "8080",
		"DD_SOURCE":     "test_source",
		"DD_SERVICE":    "test_service",
		"DD_TAGS":       "env:test,version:1.0",
		"DATABASE_TYPE": "mongodb",
	}

	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}

	// Create a temporary .env file
	tempDir := b.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tempDir)

	envContent := `MONGO_URI=mongodb://localhost:27017
MONGO_DB=test_db
PORT=8080
DD_SOURCE=test_source
DD_SERVICE=test_service
DD_TAGS=env:test,version:1.0
DATABASE_TYPE=mongodb`

	err := os.WriteFile(envFile, []byte(envContent), 0644)
	require.NoError(b, err)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg, err := config.NewConfig()
		if err != nil {
			b.Fatal(err)
		}
		if cfg == nil {
			b.Fatal("config is nil")
		}
	}
}

// Helper functions for managing environment state

func saveEnvironment() map[string]string {
	envVars := []string{
		"MONGO_URI", "MONGO_DB", "PORT", "DD_SOURCE",
		"DD_SERVICE", "DD_TAGS", "DATABASE_TYPE",
	}

	saved := make(map[string]string)
	for _, key := range envVars {
		if value, exists := os.LookupEnv(key); exists {
			saved[key] = value
		}
	}
	return saved
}

func restoreEnvironment(saved map[string]string) {
	// Clear all test environment variables
	envVars := []string{
		"MONGO_URI", "MONGO_DB", "PORT", "DD_SOURCE",
		"DD_SERVICE", "DD_TAGS", "DATABASE_TYPE",
	}

	for _, key := range envVars {
		os.Unsetenv(key)
	}

	// Restore original values
	for key, value := range saved {
		os.Setenv(key, value)
	}
}
