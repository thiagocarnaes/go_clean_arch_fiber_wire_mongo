package user

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type UpdateUserUseCase struct {
	repo repositories.UserRepository
}

func NewUpdateUserUseCase(repo repositories.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, userDTO *dto.UserDTO) error {
	user := mappers.ToUserEntity(userDTO)
	return uc.repo.Update(ctx, user)
}
