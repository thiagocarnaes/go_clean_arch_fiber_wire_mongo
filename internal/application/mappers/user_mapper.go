package mappers

import (
	"user-management/internal/application/dto"
	"user-management/internal/domain/entities"
)

func ToUserDTO(user *entities.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func ToUserEntity(dto *dto.UserDTO) *entities.User {
	return &entities.User{
		ID:    dto.ID,
		Name:  dto.Name,
		Email: dto.Email,
	}
}
