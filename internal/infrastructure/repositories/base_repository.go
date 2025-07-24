package repositories

import (
	"fmt"
	"user-management/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	dbManager *database.DatabaseManager
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(dbManager *database.DatabaseManager) *BaseRepository {
	return &BaseRepository{
		dbManager: dbManager,
	}
}

// GetMongoCollection returns a MongoDB collection (type-safe helper)
func (r *BaseRepository) GetMongoCollection(name string) (*mongo.Collection, error) {
	if !r.dbManager.IsConnected() {
		return nil, fmt.Errorf("database is not connected")
	}

	// Get the collection connection
	collectionInterface := r.dbManager.GetCollectionConnection(name)
	if collectionInterface == nil {
		return nil, fmt.Errorf("failed to get collection: %s", name)
	}

	// Type assert to MongoDB collection
	collection, ok := collectionInterface.(*mongo.Collection)
	if !ok {
		return nil, fmt.Errorf("expected MongoDB collection, got %T", collectionInterface)
	}

	return collection, nil
}

// GetDatabaseManager returns the database manager
func (r *BaseRepository) GetDatabaseManager() *database.DatabaseManager {
	return r.dbManager
}
