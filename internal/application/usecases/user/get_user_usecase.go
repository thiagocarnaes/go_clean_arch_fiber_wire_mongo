package user

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type GetUserUseCase struct {
	repo repositories.UserRepository
}

func NewGetUserUseCase(repo repositories.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*dto.UserDTO, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mappers.ToUserDTO(user), nil
}
