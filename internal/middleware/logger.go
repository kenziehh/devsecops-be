package middleware

import (
    "devsecops-be/pkg/logger"
    "time"

    "github.com/gofiber/fiber/v2"
)

func FiberLogger(log logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()

        // Process request
        err := c.Next()

        // Log request details
        duration := time.Since(start)

        fields := logger.Fields{
            "method":     c.Method(),
            "path":       c.Path(),
            "status":     c.Response().StatusCode(),
            "duration":   duration.Milliseconds(),
            "user_agent": c.Get("User-Agent"),
            "ip":         c.IP(),
        }

        // Add query params if present
        if c.Request().URI().QueryString() != nil {
            fields["query"] = string(c.Request().URI().QueryString())
        }

        // Add user_id if present in context
        if userID := c.Locals("user_id"); userID != nil {
            fields["user_id"] = userID
        }

        message := "HTTP Request"

        switch {
        case c.Response().StatusCode() >= 500:
            log.Error(c.Context(), message, err, fields)
        case c.Response().StatusCode() >= 400:
            log.Warn(c.Context(), message, fields)
        default:
            log.Info(c.Context(), message, fields)
        }

        return err
    }
}
