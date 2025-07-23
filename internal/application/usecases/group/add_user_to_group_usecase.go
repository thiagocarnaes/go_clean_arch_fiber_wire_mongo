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
	_, errGroup := uc.groupRepo.GetByID(ctx, groupID)
	if errGroup != nil {
		return errGroup
	}
	_, errUser := uc.userRepo.GetByID(ctx, userID)
	if errUser != nil {
		return errUser
	}
	return uc.groupRepo.AddUserToGroup(ctx, groupID, userID)
}
