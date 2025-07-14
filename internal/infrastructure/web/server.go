package web

import (
	"github.com/gofiber/fiber/v2"
	"user-management/internal/config"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	//, userHandler *handlers.UserHandler, groupHandler *handlers.GroupHandler

	app := fiber.New()
	//routes.SetupRoutes(app, userHandler, groupHandler)
	return &Server{app: app, cfg: cfg}
}

func (s *Server) Start() error {
	return s.app.Listen(s.cfg.Port)
}
