package user

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type GetUserUseCase struct {
	repo repositories.IUserRepository
}

func NewGetUserUseCase(repo repositories.IUserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*dto.UserResponseDTO, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mappers.ToUserResponseDTO(user), nil
}
