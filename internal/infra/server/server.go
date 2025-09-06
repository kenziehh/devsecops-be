package server

import (
	"context"
	"devsecops-be/config/env"
	"devsecops-be/pkg/logger"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app    *fiber.App
	logger logger.Logger
}

func NewServer(app *fiber.App, logger logger.Logger) *Server {
	return &Server{app: app, logger: logger}
}

func (s *Server) Start() {
	port := env.GetEnv("PORT", "8000")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		s.logger.Info(context.Background(), "Server starting", logger.Fields{
			"port": port,
			"env":  env.GetEnv("ENV", "development"),
		})

		if err := s.app.Listen(":" + port); err != nil {
			s.logger.Fatal(context.Background(), "Failed to start server", err)
		}
	}()

	<-c
	s.logger.Info(context.Background(), "Shutting down server...")

	if err := s.app.ShutdownWithTimeout(10 * time.Second); err != nil {
		s.logger.Error(context.Background(), "Server forced to shutdown", err)
	}

	s.logger.Info(context.Background(), "Server gracefully stopped")
}
