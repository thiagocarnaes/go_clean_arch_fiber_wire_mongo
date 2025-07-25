package repositories

import (
	"context"
	"user-management/internal/domain/entities"
)

type IGroupRepository interface {
	Create(ctx context.Context, group *entities.Group) error
	GetByID(ctx context.Context, id string) (*entities.Group, error)
	List(ctx context.Context, offset int64, limit int64) ([]*entities.Group, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, group *entities.Group) error
	Delete(ctx context.Context, id string) error
	AddUserToGroup(ctx context.Context, groupID, userID string) error
	RemoveUserFromGroup(ctx context.Context, groupID, userID string) error
}
