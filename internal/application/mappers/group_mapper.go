package mappers

import (
	"user-management/internal/application/dto"
	"user-management/internal/domain/entities"
)

func ToListGroupResponseDTO(groups []*entities.Group, total int64, page int64, perPage int64) *dto.ListGroupResponseDTO {

	var groupDTOs []*dto.GroupResponseDTO
	for _, group := range groups {
		groupDTOs = append(groupDTOs, ToGroupResponseDTO(group))
	}

	totalPages := calculateTotalPages(total, perPage)

	return &dto.ListGroupResponseDTO{
		Data: groupDTOs,
		Meta: dto.Meta{
			Total:      total,
			Page:       page + 1, // Converte de volta para página baseada em 1 para o usuário
			PerPage:    perPage,
			TotalPages: totalPages,
		},
	}
}

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
		Members: dto.Members,
	}
}
