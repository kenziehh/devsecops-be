package routes

import (
	"devsecops-be/internal/middleware"
	"devsecops-be/pkg/jwt"
	"devsecops-be/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, jwtUtil jwt.JWTUtil, appLogger logger.Logger) {
	// Public health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Server is healthy",
			"data": fiber.Map{
				"timestamp": time.Now().Unix(),
				"service":   "devsecops-be",
				"version":   "1.0.0",
			},
		})
	})

	// Protected health check
	authMiddleware := middleware.AuthMiddleware(jwtUtil, appLogger)
	app.Get("/health/protected", authMiddleware, func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(int)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Protected endpoint accessible",
			"data": fiber.Map{
				"user_id":   userID,
				"timestamp": time.Now().Unix(),
			},
		})
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		appLogger.Warn(c.Context(), "Route not found", logger.Fields{
			"path":   c.Path(),
			"method": c.Method(),
			"ip":     c.IP(),
		})

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Route not found",
			"error": fiber.Map{
				"code": "ROUTE_NOT_FOUND",
				"type": "NOT_FOUND",
			},
		})
	})
}
