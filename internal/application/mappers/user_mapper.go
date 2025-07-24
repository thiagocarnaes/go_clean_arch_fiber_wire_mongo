package mappers

import (
	"user-management/internal/application/dto"
	"user-management/internal/domain/entities"
)

func ToUserResponseDTO(user *entities.User) *dto.UserResponseDTO {
	return &dto.UserResponseDTO{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}
}

func ToUserEntityFromRequest(dto *dto.CreateUserRequestDTO) *entities.User {
	return &entities.User{
		Name:  dto.Name,
		Email: dto.Email,
	}
}
