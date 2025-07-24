package interfaces

import (
	"context"
)

// Database represents a generic database connection interface
type Database interface {
	// Connect establishes a connection to the database
	Connect(ctx context.Context) error

	// Disconnect closes the database connection
	Disconnect(ctx context.Context) error

	// Ping tests the database connection
	Ping(ctx context.Context) error

	// GetConnection returns the underlying database connection
	GetConnection() any

	// GetCollectionConnection returns a collection/table connection
	GetCollectionConnection(name string) any

	// IsConnected checks if the database is connected
	IsConnected() bool
}
