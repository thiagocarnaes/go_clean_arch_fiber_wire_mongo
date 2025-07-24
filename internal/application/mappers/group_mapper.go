package mappers

import (
	"user-management/internal/application/dto"
	"user-management/internal/domain/entities"
)

func ToGroupResponseDTO(group *entities.Group) *dto.GroupResponseDTO {
	return &dto.GroupResponseDTO{
		ID:      group.ID.Hex(),
		Name:    group.Name,
		Members: group.Members,
	}
}

func ToGroupEntityFromRequest(dto *dto.CreateGroupRequestDTO) *entities.Group {
	return &entities.Group{
		Name:    dto.Name,
		Members: []string{},
	}
}
