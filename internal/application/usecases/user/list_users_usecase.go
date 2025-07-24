package user

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type ListUsersUseCase struct {
	repo repositories.IUserRepository
}

func NewListUsersUseCase(repo repositories.IUserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{repo: repo}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context) ([]*dto.UserResponseDTO, error) {
	users, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var userDTOs []*dto.UserResponseDTO
	for _, user := range users {
		userDTOs = append(userDTOs, mappers.ToUserResponseDTO(user))
	}
	return userDTOs, nil
}
