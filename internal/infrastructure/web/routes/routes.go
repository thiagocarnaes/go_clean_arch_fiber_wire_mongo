package routes

import (
	"github.com/gofiber/fiber/v2"
	"user-management/internal/infrastructure/web/handlers"
)

func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, groupHandler *handlers.GroupHandler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User routes
	users := v1.Group("/users")
	users.Post("/", userHandler.Create)
	users.Get("/:id", userHandler.Get)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)
	users.Get("/", userHandler.List)

	// Group routes
	groups := v1.Group("/groups")
	groups.Post("/", groupHandler.Create)
	groups.Get("/:id", groupHandler.Get)
	groups.Put("/:id", groupHandler.Update)
	groups.Delete("/:id", groupHandler.Delete)
	groups.Get("/", groupHandler.List)
	groups.Post("/:groupId/members/:userId", groupHandler.AddUser)
	groups.Delete("/:groupId/members/:userId", groupHandler.RemoveUser)
}
