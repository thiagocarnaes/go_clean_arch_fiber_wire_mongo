package acceptance

import (
	"context"
	"net/http"
	"testing"
	"time"

	"user-management/internal/application/usecases/group"
	"user-management/internal/application/usecases/user"
	"user-management/internal/config"
	"user-management/internal/infrastructure/database"
	"user-management/internal/infrastructure/logger"
	"user-management/internal/infrastructure/repositories"
	"user-management/internal/infrastructure/web/controllers"
	"user-management/internal/infrastructure/web/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

type TestApp struct {
	App       *fiber.App
	DB        *database.MongoDB
	Container testcontainers.Container
}

func SetupTestApp(t *testing.T) *TestApp {
	ctx := context.Background()

	// Start MongoDB container
	mongoContainer, err := mongodb.Run(ctx, "mongo:7.0")
	require.NoError(t, err)

	// Get connection string
	connectionString, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	// Setup config with test values
	cfg := &config.Config{
		MongoURI:     connectionString,
		MongoDB:      "testdb",
		Port:         "8080",
		DatabaseType: "mongodb",
	}

	// Initialize logger
	log := logger.NewLogger()

	// Initialize database
	db, err := database.NewMongoDB(cfg, log)
	require.NoError(t, err)

	// Initialize repositories
	userRepo, err := repositories.NewUserRepository(db)
	require.NoError(t, err)
	groupRepo, err := repositories.NewGroupRepository(db)
	require.NoError(t, err)

	// Initialize use cases
	createUserUseCase := user.NewCreateUserUseCase(userRepo)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepo)
	deleteUserUseCase := user.NewDeleteUserUseCase(userRepo)
	listUsersUseCase := user.NewListUsersUseCase(userRepo)

	createGroupUseCase := group.NewCreateGroupUseCase(groupRepo)
	getGroupUseCase := group.NewGetGroupUseCase(groupRepo)
	updateGroupUseCase := group.NewUpdateGroupUseCase(groupRepo)
	deleteGroupUseCase := group.NewDeleteGroupUseCase(groupRepo)
	listGroupsUseCase := group.NewListGroupsUseCase(groupRepo)
	addUserToGroupUseCase := group.NewAddUserToGroupUseCase(groupRepo, userRepo)
	removeUserFromGroupUseCase := group.NewRemoveUserFromGroupUseCase(groupRepo)

	// Initialize controllers
	userController := controllers.NewUserController(
		createUserUseCase,
		getUserUseCase,
		updateUserUseCase,
		deleteUserUseCase,
		listUsersUseCase,
	)

	groupController := controllers.NewGroupController(
		createGroupUseCase,
		getGroupUseCase,
		updateGroupUseCase,
		deleteGroupUseCase,
		listGroupsUseCase,
		addUserToGroupUseCase,
		removeUserFromGroupUseCase,
	)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Setup routes
	routes.SetupRoutes(app, userController, groupController)

	return &TestApp{
		App:       app,
		DB:        db,
		Container: mongoContainer,
	}
}

func (ta *TestApp) Cleanup(t *testing.T) {
	ctx := context.Background()

	// Clean up database
	if ta.DB != nil && ta.DB.DB != nil {
		err := ta.DB.DB.Drop(ctx)
		require.NoError(t, err)
	}

	// Stop container
	if ta.Container != nil {
		err := ta.Container.Terminate(ctx)
		require.NoError(t, err)
	}
}

func (ta *TestApp) Request(req *http.Request) (*http.Response, error) {
	return ta.App.Test(req)
}

func (ta *TestApp) ClearDatabase(t *testing.T) {
	ctx := context.Background()

	// Drop all collections
	collections, err := ta.DB.DB.ListCollectionNames(ctx, map[string]interface{}{})
	require.NoError(t, err)

	for _, collection := range collections {
		err := ta.DB.DB.Collection(collection).Drop(ctx)
		require.NoError(t, err)
	}
}

// Wait for container to be ready
func WaitForContainer(container testcontainers.Container, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return container.Start(ctx)
}
