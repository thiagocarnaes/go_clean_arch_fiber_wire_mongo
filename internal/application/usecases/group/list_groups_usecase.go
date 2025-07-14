package group

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type ListGroupsUseCase struct {
	repo repositories.GroupRepository
}

func NewListGroupsUseCase(repo repositories.GroupRepository) *ListGroupsUseCase {
	return &ListGroupsUseCase{repo: repo}
}

func (gc *ListGroupsUseCase) Execute(ctx context.Context) ([]*dto.GroupDTO, error) {
	groups, err := gc.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var groupDTOs []*dto.GroupDTO
	for _, group := range groups {
		groupDTOs = append(groupDTOs, mappers.ToGroupDTO(group))
	}
	return groupDTOs, nil
}
