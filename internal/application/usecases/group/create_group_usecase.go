package group

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type CreateGroupUseCase struct {
	repo repositories.IGroupRepository
}

func NewCreateGroupUseCase(repo repositories.IGroupRepository) *CreateGroupUseCase {
	return &CreateGroupUseCase{repo: repo}
}

func (uc *CreateGroupUseCase) Execute(ctx context.Context, groupDTO *dto.CreateGroupRequestDTO) (*dto.GroupResponseDTO, error) {
	group := mappers.ToGroupEntityFromRequest(groupDTO)
	err := uc.repo.Create(ctx, group)
	if err != nil {
		return nil, err
	}
	return mappers.ToGroupResponseDTO(group), nil
}
