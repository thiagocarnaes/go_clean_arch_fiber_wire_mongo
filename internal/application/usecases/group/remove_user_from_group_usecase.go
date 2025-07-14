package group

import (
	"context"
	"user-management/internal/domain/interfaces/repositories"
)

type RemoveUserFromGroupUseCase struct {
	groupRepo repositories.GroupRepository
}

func NewRemoveUserFromGroupUseCase(groupRepo repositories.GroupRepository) *RemoveUserFromGroupUseCase {
	return &RemoveUserFromGroupUseCase{groupRepo: groupRepo}
}

func (uc *RemoveUserFromGroupUseCase) Execute(ctx context.Context, groupID, userID string) error {
	return uc.groupRepo.RemoveUserFromGroup(ctx, groupID, userID)
}
