package user

import (
	"context"
	"user-management/internal/domain/interfaces/repositories"
)

type DeleteUserUseCase struct {
	repo repositories.IUserRepository
}

func NewDeleteUserUseCase(repo repositories.IUserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{repo: repo}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
