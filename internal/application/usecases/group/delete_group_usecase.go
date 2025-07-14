package group

import (
	"context"
	"user-management/internal/domain/interfaces/repositories"
)

type DeleteGroupUseCase struct {
	repo repositories.IGroupRepository
}

func NewDeleteGroupUseCase(repo repositories.IGroupRepository) *DeleteGroupUseCase {
	return &DeleteGroupUseCase{repo: repo}
}

func (uc *DeleteGroupUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
