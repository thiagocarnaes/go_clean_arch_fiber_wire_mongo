package mappers

import (
	"user-management/internal/application/dto"
	"user-management/internal/domain/entities"
)

func ToUserListResponseDTO(users []*entities.User, total int64, page int64, perPage int64) *dto.UserListResponseDTO {
	var userDTOs []dto.UserResponseDTO
	for _, user := range users {
		userDTOs = append(userDTOs, *ToUserResponseDTO(user))
	}

	totalPages := calculateTotalPages(total, perPage)

	return &dto.UserListResponseDTO{
		Data: userDTOs,
		Meta: dto.Meta{
			Total:      total,
			Page:       page + 1, // Converte de volta para página baseada em 1 para o usuário
			PerPage:    perPage,
			TotalPages: totalPages,
		},
	}
}

func ToUserResponseDTO(user *entities.User) *dto.UserResponseDTO {
	return &dto.UserResponseDTO{
		ID:       user.ID.Hex(),
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
}

func ToUserEntityFromRequest(dto *dto.CreateUserRequestDTO) *entities.User {
	return &entities.User{
		Name:     dto.Name,
		Email:    dto.Email,
		IsActive: dto.IsActive,
	}
}
