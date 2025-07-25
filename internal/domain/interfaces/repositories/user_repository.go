package repositories

import (
	"context"
	"user-management/internal/domain/entities"
)

type IUserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	List(ctx context.Context, offset int64, limit int64) ([]*entities.User, error)
	Search(ctx context.Context, searchTerm string, offset int64, limit int64) ([]*entities.User, error)
	Count(ctx context.Context) (int64, error)
	CountSearch(ctx context.Context, searchTerm string) (int64, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
}
