package group

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type GetGroupUseCase struct {
	repo repositories.IGroupRepository
}

func NewGetGroupUseCase(repo repositories.IGroupRepository) *GetGroupUseCase {
	return &GetGroupUseCase{repo: repo}
}

func (uc *GetGroupUseCase) Execute(ctx context.Context, id string) (*dto.GroupResponseDTO, error) {
	group, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mappers.ToGroupResponseDTO(group), nil
}
