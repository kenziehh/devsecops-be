package auth

import (
	"database/sql"
	"devsecops-be/internal/domain/auth/handler/http"
	"devsecops-be/internal/domain/auth/repository"
	"devsecops-be/internal/domain/auth/service"
	"devsecops-be/pkg/jwt"
	"devsecops-be/pkg/logger"
	"devsecops-be/pkg/password"
	"devsecops-be/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthModule struct {
	Handler *http.AuthHandler
	Service service.AuthService
}

func NewAuthModule(db *sql.DB, jwtUtil jwt.JWTUtil, logger logger.Logger) *AuthModule {
	// Initialize dependencies
	authRepo := repository.NewAuthRepository(db)
	passUtil := password.NewPasswordUtil()
	validator := validator.NewValidator()

	// Initialize service
	authService := service.NewAuthService(authRepo, jwtUtil, passUtil, logger)

	// Initialize handler
	authHandler := http.NewAuthHandler(authService, validator, logger)

	return &AuthModule{
		Handler: authHandler,
		Service: authService,
	}
}

func (m *AuthModule) RegisterRoutes(app *fiber.App) {
	auth := app.Group("/api/v1/auth")

	auth.Post("/login", m.Handler.Login)
	auth.Post("/register", m.Handler.Register)
}
