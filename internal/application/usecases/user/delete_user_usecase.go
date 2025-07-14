package user

import (
	"context"
	"user-management/internal/domain/interfaces/repositories"
)

type DeleteUserUseCase struct {
	repo repositories.UserRepository
}

func NewDeleteUserUseCase(repo repositories.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{repo: repo}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
