package group

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/interfaces/repositories"
)

type CreateGroupUseCase struct {
	repo repositories.GroupRepository
}

func NewCreateGroupUseCase(repo repositories.GroupRepository) *CreateGroupUseCase {
	return &CreateGroupUseCase{repo: repo}
}

func (uc *CreateGroupUseCase) Execute(ctx context.Context, groupDTO *dto.GroupDTO) error {
	group := mappers.ToGroupEntity(groupDTO)
	return uc.repo.Create(ctx, group)
}
