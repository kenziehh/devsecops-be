package errors

import (
    "fmt"
    "net/http"

    "devsecops-be/pkg/response"
    "github.com/gofiber/fiber/v2"
)

// Custom error types
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Type    string `json:"type"`
    HTTPStatus int `json:"-"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Pre-defined errors
var (
    ErrUserNotFound = &AppError{
        Code:       "USER_NOT_FOUND",
        Message:    "User not found",
        Type:       "NOT_FOUND",
        HTTPStatus: http.StatusNotFound,
    }

    ErrUserAlreadyExists = &AppError{
        Code:       "USER_ALREADY_EXISTS",
        Message:    "User with this email already exists",
        Type:       "CONFLICT",
        HTTPStatus: http.StatusConflict,
    }

    ErrInvalidCredentials = &AppError{
        Code:       "INVALID_CREDENTIALS",
        Message:    "Invalid email or password",
        Type:       "UNAUTHORIZED",
        HTTPStatus: http.StatusUnauthorized,
    }

    ErrInvalidToken = &AppError{
        Code:       "INVALID_TOKEN",
        Message:    "Invalid or expired token",
        Type:       "UNAUTHORIZED",
        HTTPStatus: http.StatusUnauthorized,
    }

    ErrTokenRequired = &AppError{
        Code:       "TOKEN_REQUIRED",
        Message:    "Authorization token is required",
        Type:       "UNAUTHORIZED",
        HTTPStatus: http.StatusUnauthorized,
    }

    ErrInternalServer = &AppError{
        Code:       "INTERNAL_SERVER_ERROR",
        Message:    "Internal server error",
        Type:       "INTERNAL_ERROR",
        HTTPStatus: http.StatusInternalServerError,
    }
)

// Error wrapping functions
func WrapDatabaseError(err error, message string) *AppError {
    return &AppError{
        Code:       "DATABASE_ERROR",
        Message:    fmt.Sprintf("%s: %v", message, err),
        Type:       "INTERNAL_ERROR",
        HTTPStatus: http.StatusInternalServerError,
    }
}

func WrapValidationError(err error, message string) *AppError {
    return &AppError{
        Code:       "VALIDATION_ERROR",
        Message:    fmt.Sprintf("%s: %v", message, err),
        Type:       "BAD_REQUEST",
        HTTPStatus: http.StatusBadRequest,
    }
}

func WrapInternalError(err error, message string) *AppError {
    return &AppError{
        Code:       "INTERNAL_ERROR",
        Message:    fmt.Sprintf("%s: %v", message, err),
        Type:       "INTERNAL_ERROR",
        HTTPStatus: http.StatusInternalServerError,
    }
}

// HTTP Error Handler
func HandleHTTPError(c *fiber.Ctx, err error) error {
    if appErr, ok := err.(*AppError); ok {
        switch appErr.HTTPStatus {
        case http.StatusBadRequest:
            return response.BadRequest(c, appErr.Message, appErr)
        case http.StatusUnauthorized:
            return response.Unauthorized(c, appErr.Message, appErr)
        case http.StatusNotFound:
            return response.NotFound(c, appErr.Message, appErr)
        case http.StatusConflict:
            return response.Conflict(c, appErr.Message, appErr)
        default:
            return response.InternalServerError(c, appErr.Message, appErr)
        }
    }
    
    // Fallback for unknown errors
    return response.InternalServerError(c, "An unexpected error occurred", nil)
}