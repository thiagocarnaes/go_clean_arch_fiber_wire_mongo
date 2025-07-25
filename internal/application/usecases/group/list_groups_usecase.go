package group

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type ListGroupsUseCase struct {
	repo repositories.IGroupRepository
}

func NewListGroupsUseCase(repo repositories.IGroupRepository) *ListGroupsUseCase {
	return &ListGroupsUseCase{repo: repo}
}

func (gc *ListGroupsUseCase) Execute(ctx context.Context, input *dto.ListGroupQueryParam) (*dto.ListGroupResponseDTO, error) {

	groups, err := gc.repo.List(ctx, input.Page, input.PerPage)
	if err != nil {
		return nil, err
	}

	total, err := gc.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	groupDTOs := mappers.ToListGroupResponseDTO(groups, total, input.Page, input.PerPage)
	return groupDTOs, nil
}
