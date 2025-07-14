package mappers

import (
	"user-management/internal/application/dto"
	"user-management/internal/domain/entities"
)

func ToGroupDTO(group *entities.Group) *dto.GroupDTO {
	return &dto.GroupDTO{
		ID:      group.ID,
		Name:    group.Name,
		Members: group.Members,
	}
}

func ToGroupEntity(dto *dto.GroupDTO) *entities.Group {
	return &entities.Group{
		ID:      dto.ID,
		Name:    dto.Name,
		Members: dto.Members,
	}
}
