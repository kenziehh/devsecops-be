package fiber

import (
	"devsecops-be/internal/domain/auth"
	"devsecops-be/internal/middleware"
	"devsecops-be/pkg/jwt"
	"devsecops-be/pkg/logger"
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewFiberApp(appLogger logger.Logger, db *sql.DB, jwtUtil jwt.JWTUtil) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			appLogger.Error(c.Context(), "Unhandled error in Fiber", err, logger.Fields{
				"path":   c.Path(),
				"method": c.Method(),
				"ip":     c.IP(),
			})

			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": message,
				"error": fiber.Map{
					"code": "FIBER_ERROR",
					"type": "INTERNAL_ERROR",
				},
			})
		},
		DisableStartupMessage: true,
	})

	// Middleware global
	app.Use(helmet.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: os.Getenv("ENV") == "development"}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     getEnv("CORS_ORIGINS", "*"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))
	app.Use(middleware.FiberLogger(appLogger))

	// Auth module
	authModule := auth.NewAuthModule(db, jwtUtil, appLogger)
	authModule.RegisterRoutes(app)


	return app
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
