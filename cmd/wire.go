//go:build wireinject
// +build wireinject

package cmd

import (
	"user-management/internal/application/usecases/group"
	"user-management/internal/application/usecases/user"
	"user-management/internal/config"
	"user-management/internal/infrastructure/database"
	"user-management/internal/infrastructure/logger"
	"user-management/internal/infrastructure/repositories"
	"user-management/internal/infrastructure/web"
	"user-management/internal/infrastructure/web/controllers"

	"github.com/google/wire"
)

func InitializeServer() (*web.Server, error) {
	wire.Build(
		logger.NewLogger,
		config.NewConfig,
		database.NewDatabaseManager,
		repositories.NewUserRepository,
		repositories.NewGroupRepository,
		user.NewCreateUserUseCase,
		user.NewGetUserUseCase,
		user.NewUpdateUserUseCase,
		user.NewDeleteUserUseCase,
		user.NewListUsersUseCase,
		group.NewCreateGroupUseCase,
		group.NewGetGroupUseCase,
		group.NewUpdateGroupUseCase,
		group.NewDeleteGroupUseCase,
		group.NewListGroupsUseCase,
		group.NewAddUserToGroupUseCase,
		group.NewRemoveUserFromGroupUseCase,
		controllers.NewUserController,
		controllers.NewGroupController,
		web.NewServer,
	)
	return &web.Server{}, nil
}
