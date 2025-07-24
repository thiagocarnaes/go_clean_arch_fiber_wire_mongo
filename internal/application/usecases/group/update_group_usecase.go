package group

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type UpdateGroupUseCase struct {
	repo repositories.IGroupRepository
}

func NewUpdateGroupUseCase(repo repositories.IGroupRepository) *UpdateGroupUseCase {
	return &UpdateGroupUseCase{repo: repo}
}

func (uc *UpdateGroupUseCase) Execute(ctx context.Context, groupID string, groupDTO *dto.CreateGroupRequestDTO) (*dto.GroupResponseDTO, error) {
	// First get the existing group to preserve members
	existingGroup, err := uc.repo.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Create the updated group with existing members
	group := mappers.ToGroupEntityFromRequest(groupDTO)
	group.ID = existingGroup.ID
	group.Members = groupDTO.Members

	err = uc.repo.Update(ctx, group)
	if err != nil {
		return nil, err
	}

	return mappers.ToGroupResponseDTO(group), nil
}
