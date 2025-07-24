package user

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type CreateUserUseCase struct {
	repo repositories.IUserRepository
}

func NewCreateUserUseCase(repo repositories.IUserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, userDTO *dto.CreateUserRequestDTO) (*dto.UserResponseDTO, error) {
	user := mappers.ToUserEntityFromRequest(userDTO)
	err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return mappers.ToUserResponseDTO(user), nil
}
