package main

import (
	"context"
	"devsecops-be/config/fiber"
	"devsecops-be/internal/infra/routes"
	"devsecops-be/internal/infra/server"

	"devsecops-be/pkg/database"
	"devsecops-be/pkg/jwt"
	"devsecops-be/pkg/logger"
)

func main() {
	// Logger
	appLogger := logger.NewLogger()
	appLogger.Info(context.Background(), "Starting Fiber Auth Application")

	// Database
	db, err := database.NewPostgresConnection(appLogger)
	if err != nil {
		appLogger.Fatal(context.Background(), "Failed to connect to database", err)
	}
	defer db.Close()

	// JWT Utility
	jwtUtil := jwt.NewJWTUtil()

	// Fiber App
	fiberApp := fiber.NewFiberApp(appLogger, db, jwtUtil)

	// Register routes
	routes.RegisterRoutes(fiberApp, jwtUtil, appLogger)

	// Server
	server := server.NewServer(fiberApp, appLogger)
	server.Start()
}
