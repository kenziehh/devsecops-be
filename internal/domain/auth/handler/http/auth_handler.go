package http

import (
	"devsecops-be/internal/domain/auth/dto"
	"devsecops-be/internal/domain/auth/service"
	"devsecops-be/pkg/errors"
	"devsecops-be/pkg/logger"
	"devsecops-be/pkg/response"
	"devsecops-be/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
	validator   validator.Validator
	logger      logger.Logger
}

func NewAuthHandler(
	authService service.AuthService,
	validator validator.Validator,
	logger logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn(ctx, "Invalid request body in login", logger.Fields{
			"error": err.Error(),
		})
		return response.BadRequest(c, "Invalid request body", nil)
	}

	if err := h.validator.Validate(req); err != nil {
		h.logger.Warn(ctx, "Validation failed in login", logger.Fields{
			"validation_errors": err,
		})
		return response.BadRequest(c, "Validation failed", err)
	}

	result, err := h.authService.Login(ctx, req)
	if err != nil {
		return errors.HandleHTTPError(c, err)
	}

	return response.Success(c, "Login successful", result)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn(ctx, "Invalid request body in registration", logger.Fields{
			"error": err.Error(),
		})
		return response.BadRequest(c, "Invalid request body", nil)
	}

	if err := h.validator.Validate(req); err != nil {
		h.logger.Warn(ctx, "Validation failed in registration", logger.Fields{
			"validation_errors": err,
		})
		return response.BadRequest(c, "Validation failed", err)
	}

	result, err := h.authService.Register(ctx, req)
	if err != nil {
		return errors.HandleHTTPError(c, err)
	}

	return response.Created(c, "Registration successful", result)
}
