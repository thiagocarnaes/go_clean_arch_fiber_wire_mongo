package group

import (
	"context"
	"user-management/internal/domain/interfaces/repositories"
)

type AddUserToGroupUseCase struct {
	groupRepo repositories.IGroupRepository
	userRepo  repositories.IUserRepository
}

func NewAddUserToGroupUseCase(groupRepo repositories.IGroupRepository, userRepo repositories.IUserRepository) *AddUserToGroupUseCase {
	return &AddUserToGroupUseCase{groupRepo: groupRepo, userRepo: userRepo}
}

func (uc *AddUserToGroupUseCase) Execute(ctx context.Context, groupID, userID string) error {
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	return uc.groupRepo.AddUserToGroup(ctx, groupID, userID)
}
