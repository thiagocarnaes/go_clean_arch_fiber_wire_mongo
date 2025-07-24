package user

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type UpdateUserUseCase struct {
	repo repositories.IUserRepository
}

func NewUpdateUserUseCase(repo repositories.IUserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, userID string, userDTO *dto.CreateUserRequestDTO) (*dto.UserResponseDTO, error) {
	existingUser, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := mappers.ToUserEntityFromRequest(userDTO)
	user.ID = existingUser.ID
	errUpdate := uc.repo.Update(ctx, user)
	if errUpdate != nil {
		return nil, errUpdate
	}
	return mappers.ToUserResponseDTO(user), nil
}
