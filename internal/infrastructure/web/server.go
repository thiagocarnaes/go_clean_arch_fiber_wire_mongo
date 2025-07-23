package web

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management/internal/config"
	"user-management/internal/infrastructure/database"
	"user-management/internal/infrastructure/web/controllers"
	"user-management/internal/infrastructure/web/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Server struct {
	app     *fiber.App
	cfg     *config.Config
	log     *logrus.Logger
	mongoDB *database.MongoDB
}

func NewServer(cfg *config.Config,
	UserController *controllers.UserController,
	GroupController *controllers.GroupController,
	log *logrus.Logger,
	mongoDB *database.MongoDB) *Server {

	app := fiber.New()
	routes.SetupRoutes(app, UserController, GroupController)
	return &Server{app: app, cfg: cfg, log: log, mongoDB: mongoDB}
}

func (s *Server) Start() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Iniciar o servidor em uma goroutine
	go func() {
		s.log.WithFields(logrus.Fields{
			"ddsource": "go",
			"service":  "user-management",
			"ddtags":   "env:dev,app:fiber",
			"port":     s.cfg.Port,
		}).Info("Starting server")
		if err := s.app.Listen(s.cfg.Port); err != nil {
			s.log.WithFields(logrus.Fields{
				"ddsource": "go",
				"service":  "user-management",
				"ddtags":   "env:dev,app:fiber",
				"error":    err.Error(),
			}).Error("Failed to start server")
		}
	}()

	// Aguardar sinal de shutdown
	<-stop
	s.log.WithFields(logrus.Fields{
		"ddsource": "go",
		"service":  "user-management",
		"ddtags":   "env:dev,app:fiber",
	}).Info("Shutdown signal received, initiating graceful shutdown")

	// Criar contexto com timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Fechar conexões do Fiber
	if err := s.app.Shutdown(); err != nil {
		s.log.WithFields(logrus.Fields{
			"ddsource": "go",
			"service":  "user-management",
			"ddtags":   "env:dev,app:fiber",
			"error":    err.Error(),
		}).Error("Failed to shutdown Fiber server")
		return err
	}

	// Fechar conexão com MongoDB
	if s.mongoDB == nil {
		s.log.WithFields(logrus.Fields{
			"ddsource": "go",
			"service":  "user-management",
			"ddtags":   "env:dev,app:fiber",
		}).Warn("MongoDB client is nil, skipping disconnect")
	} else {
		if err := s.mongoDB.Client.Disconnect(ctx); err != nil {
			s.log.WithFields(logrus.Fields{
				"ddsource": "go",
				"service":  "user-management",
				"ddtags":   "env:dev,app:fiber",
				"error":    err.Error(),
			}).Error("Failed to disconnect MongoDB client")
			return err
		}
	}

	s.log.WithFields(logrus.Fields{
		"ddsource": "go",
		"service":  "user-management",
		"ddtags":   "env:dev,app:fiber",
	}).Info("Server gracefully shutdown")
	return nil
}

func (s *Server) ServerData() *Server {
	return s
}
