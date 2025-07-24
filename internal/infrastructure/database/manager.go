package database

import (
	"context"
	"fmt"
	"user-management/internal/config"
	"user-management/internal/domain/interfaces"

	"github.com/sirupsen/logrus"
)

// DatabaseManager manages database connections and provides abstraction
type DatabaseManager struct {
	config   *config.Config
	logger   *logrus.Logger
	database interfaces.Database
	dbType   string
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(cfg *config.Config, log *logrus.Logger) *DatabaseManager {
	return &DatabaseManager{
		config: cfg,
		logger: log,
		dbType: cfg.DatabaseType, // We'll add this to config
	}
}

// Initialize initializes the database connection based on configuration
func (dm *DatabaseManager) Initialize(ctx context.Context) error {
	var db interfaces.Database
	var err error

	switch dm.dbType {
	case "mongodb", "mongo", "":
		// Default to MongoDB for backward compatibility
		db, err = NewMongoDB(dm.config, dm.logger)
		if err != nil {
			return fmt.Errorf("failed to initialize MongoDB: %w", err)
		}
	case "postgresql", "postgres":
		// Future implementation
		return fmt.Errorf("PostgreSQL implementation not yet available")
	case "mysql":
		// Future implementation
		return fmt.Errorf("MySQL implementation not yet available")
	default:
		return fmt.Errorf("unsupported database type: %s", dm.dbType)
	}

	dm.database = db

	// Connect to the database
	if err := dm.database.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	dm.logger.WithFields(logrus.Fields{
		"ddsource": dm.config.DDSource,
		"service":  dm.config.DDService,
		"ddtags":   dm.config.DDTags,
		"type":     dm.dbType,
	}).Info("Database initialized successfully")

	return nil
}

// GetDatabase returns the database interface
func (dm *DatabaseManager) GetDatabase() interfaces.Database {
	return dm.database
}

// GetConnection returns the underlying database connection
func (dm *DatabaseManager) GetConnection() interface{} {
	if dm.database == nil {
		return nil
	}
	return dm.database.GetConnection()
}

// GetCollectionConnection returns a collection/table connection
func (dm *DatabaseManager) GetCollectionConnection(name string) interface{} {
	if dm.database == nil {
		return nil
	}
	return dm.database.GetCollectionConnection(name)
}

// Close closes the database connection
func (dm *DatabaseManager) Close(ctx context.Context) error {
	if dm.database == nil {
		return nil
	}

	if err := dm.database.Disconnect(ctx); err != nil {
		dm.logger.WithFields(logrus.Fields{
			"ddsource": dm.config.DDSource,
			"service":  dm.config.DDService,
			"ddtags":   dm.config.DDTags,
			"error":    err.Error(),
		}).Error("Failed to close database connection")
		return err
	}

	dm.logger.WithFields(logrus.Fields{
		"ddsource": dm.config.DDSource,
		"service":  dm.config.DDService,
		"ddtags":   dm.config.DDTags,
	}).Info("Database connection closed successfully")

	return nil
}

// IsConnected checks if the database is connected
func (dm *DatabaseManager) IsConnected() bool {
	if dm.database == nil {
		return false
	}
	return dm.database.IsConnected()
}
