package middleware

import (
    "devsecops-be/pkg/errors"
    "devsecops-be/pkg/jwt"
    "devsecops-be/pkg/logger"
    "strings"

    "github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtUtil jwt.JWTUtil, log logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            log.Warn(c.Context(), "Missing authorization header", logger.Fields{
                "path": c.Path(),
                "ip":   c.IP(),
            })
            return errors.HandleHTTPError(c, errors.ErrTokenRequired)
        }

        // Check Bearer format
        if !strings.HasPrefix(authHeader, "Bearer ") {
            log.Warn(c.Context(), "Invalid token format", logger.Fields{
                "path":        c.Path(),
                "ip":          c.IP(),
                "auth_header": authHeader[:min(len(authHeader), 20)] + "...",
            })
            return errors.HandleHTTPError(c, errors.ErrInvalidToken)
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")

        claims, err := jwtUtil.ValidateToken(token)
        if err != nil {
            log.Warn(c.Context(), "Token validation failed", logger.Fields{
                "path":  c.Path(),
                "ip":    c.IP(),
                "error": err.Error(),
            })
            return errors.HandleHTTPError(c, err)
        }

        // Check token type
        if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
            log.Warn(c.Context(), "Invalid token type", logger.Fields{
                "path":       c.Path(),
                "ip":         c.IP(),
                "token_type": tokenType,
            })
            return errors.HandleHTTPError(c, errors.ErrInvalidToken)
        }

        // Set user info in context
        userID := int(claims["user_id"].(float64))
        c.Locals("user_id", userID)
        c.Locals("token", token)

        log.Debug(c.Context(), "Token validation successful", logger.Fields{
            "user_id": userID,
            "path":    c.Path(),
        })

        return c.Next()
    }
}

