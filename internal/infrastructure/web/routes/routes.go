package routes

import (
	"log"
	"user-management/internal/infrastructure/web/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App, UserController *controllers.UserController, GroupController *controllers.GroupController) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Use(logger.New(logger.Config{
		Format:     `{"timestamp":"${time}","status":${status},"method":"${method}","path":"${path}","latency":"${latency}"}` + "\n",
		TimeFormat: "2006-01-02T15:04:05.999Z", // Formato ISO 8601
		TimeZone:   "UTC",
		Output:     log.Writer(), // Enviar logs para o Logrus
	}))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User routes
	users := v1.Group("/users")
	users.Post("/", UserController.Create)
	users.Get("/:id", UserController.Get)
	users.Put("/:id", UserController.Update)
	users.Delete("/:id", UserController.Delete)
	users.Get("/", UserController.List)

	// Group routes
	groups := v1.Group("/groups")
	groups.Post("/", GroupController.Create)
	groups.Get("/:id", GroupController.Get)
	groups.Put("/:id", GroupController.Update)
	groups.Delete("/:id", GroupController.Delete)
	groups.Get("/", GroupController.List)
	groups.Post("/:groupId/members/:userId", GroupController.AddUser)
	groups.Delete("/:groupId/members/:userId", GroupController.RemoveUser)
}
