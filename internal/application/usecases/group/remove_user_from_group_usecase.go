package group

import (
	"context"
	"user-management/internal/domain/interfaces/repositories"
)

type RemoveUserFromGroupUseCase struct {
	groupRepo repositories.IGroupRepository
}

func NewRemoveUserFromGroupUseCase(groupRepo repositories.IGroupRepository) *RemoveUserFromGroupUseCase {
	return &RemoveUserFromGroupUseCase{groupRepo: groupRepo}
}

func (uc *RemoveUserFromGroupUseCase) Execute(ctx context.Context, groupID, userID string) error {
	return uc.groupRepo.RemoveUserFromGroup(ctx, groupID, userID)
}
