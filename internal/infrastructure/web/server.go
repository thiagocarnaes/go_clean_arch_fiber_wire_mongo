package web

import (
	"github.com/gofiber/fiber/v2"
	"user-management/internal/config"
	"user-management/internal/infrastructure/web/handlers"
	"user-management/internal/infrastructure/web/routes"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func NewServer(cfg *config.Config, userHandler *handlers.UserHandler, groupHandler *handlers.GroupHandler) *Server {

	app := fiber.New()
	routes.SetupRoutes(app, userHandler, groupHandler)
	return &Server{app: app, cfg: cfg}
}

func (s *Server) Start() error {
	return s.app.Listen(s.cfg.Port)
}
