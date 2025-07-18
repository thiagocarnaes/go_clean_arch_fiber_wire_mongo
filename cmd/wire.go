// Code generated by Wire. DO NOT EDIT.
//go:build wireinject
// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"user-management/internal/application/usecases/group"
	"user-management/internal/application/usecases/user"
	"user-management/internal/config"
	"user-management/internal/infrastructure/database"
	"user-management/internal/infrastructure/logger"
	irepos "user-management/internal/infrastructure/repositories"
	"user-management/internal/infrastructure/web"
	"user-management/internal/infrastructure/web/handlers"
)

func InitializeServer() (*web.Server, error) {
	wire.Build(
		logger.NewLogger,
		config.NewConfig,
		database.NewMongoDB,
		irepos.NewUserRepository,
		irepos.NewGroupRepository,
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
		handlers.NewUserHandler,
		handlers.NewGroupHandler,
		web.NewServer,
	)
	return &web.Server{}, nil
}
