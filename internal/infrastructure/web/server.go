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
	app       *fiber.App
	cfg       *config.Config
	log       *logrus.Logger
	dbManager *database.DatabaseManager
}

func NewServer(cfg *config.Config,
	UserController *controllers.UserController,
	GroupController *controllers.GroupController,
	log *logrus.Logger,
	dbManager *database.DatabaseManager) *Server {

	app := fiber.New()
	routes.SetupRoutes(app, UserController, GroupController)
	return &Server{app: app, cfg: cfg, log: log, dbManager: dbManager}
}

func (s *Server) Start() error {
	// Initialize database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.dbManager.Initialize(ctx); err != nil {
		s.log.WithFields(logrus.Fields{
			"ddsource": s.cfg.DDSource,
			"service":  s.cfg.DDService,
			"ddtags":   s.cfg.DDTags,
			"error":    err.Error(),
		}).Error("Failed to initialize database")
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Iniciar o servidor em uma goroutine
	go func() {
		s.log.WithFields(logrus.Fields{
			"ddsource": s.cfg.DDSource,
			"service":  s.cfg.DDService,
			"ddtags":   s.cfg.DDTags,
			"port":     s.cfg.Port,
		}).Info("Starting server")
		if err := s.app.Listen(s.cfg.Port); err != nil {
			s.log.WithFields(logrus.Fields{
				"ddsource": s.cfg.DDSource,
				"service":  s.cfg.DDService,
				"ddtags":   s.cfg.DDTags,
				"error":    err.Error(),
			}).Error("Failed to start server")
		}
	}()

	// Aguardar sinal de shutdown
	<-stop
	s.log.WithFields(logrus.Fields{
		"ddsource": s.cfg.DDSource,
		"service":  s.cfg.DDService,
		"ddtags":   s.cfg.DDTags,
	}).Info("Shutdown signal received, initiating graceful shutdown")

	// Criar contexto com timeout para shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Fechar conexões do Fiber
	if err := s.app.Shutdown(); err != nil {
		s.log.WithFields(logrus.Fields{
			"ddsource": s.cfg.DDSource,
			"service":  s.cfg.DDService,
			"ddtags":   s.cfg.DDTags,
			"error":    err.Error(),
		}).Error("Failed to shutdown Fiber server")
		return err
	}

	// Fechar conexão com o banco de dados
	if s.dbManager != nil {
		if err := s.dbManager.Close(shutdownCtx); err != nil {
			s.log.WithFields(logrus.Fields{
				"ddsource": s.cfg.DDSource,
				"service":  s.cfg.DDService,
				"ddtags":   s.cfg.DDTags,
				"error":    err.Error(),
			}).Error("Failed to close database connection")
			return err
		}
	}

	s.log.WithFields(logrus.Fields{
		"ddsource": s.cfg.DDSource,
		"service":  s.cfg.DDService,
		"ddtags":   s.cfg.DDTags,
	}).Info("Server gracefully shutdown")
	return nil
}

func (s *Server) ServerData() *Server {
	return s
}
