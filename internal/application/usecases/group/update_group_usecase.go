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

func (uc *UpdateGroupUseCase) Execute(ctx context.Context, groupDTO *dto.GroupDTO) error {
	group := mappers.ToGroupEntity(groupDTO)
	return uc.repo.Update(ctx, group)
}
