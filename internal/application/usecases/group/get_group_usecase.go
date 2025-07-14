package group

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type GetGroupUseCase struct {
	repo repositories.GroupRepository
}

func NewGetGroupUseCase(repo repositories.GroupRepository) *GetGroupUseCase {
	return &GetGroupUseCase{repo: repo}
}

func (uc *GetGroupUseCase) Execute(ctx context.Context, id string) (*dto.GroupDTO, error) {
	group, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mappers.ToGroupDTO(group), nil
}
