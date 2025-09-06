package response

import (
    "github.com/gofiber/fiber/v2"
)

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusOK).JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusCreated).JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func BadRequest(c *fiber.Ctx, message string, error interface{}) error {
    return c.Status(fiber.StatusBadRequest).JSON(Response{
        Success: false,
        Message: message,
        Error:   error,
    })
}

func Unauthorized(c *fiber.Ctx, message string, error interface{}) error {
    return c.Status(fiber.StatusUnauthorized).JSON(Response{
        Success: false,
        Message: message,
        Error:   error,
    })
}

func NotFound(c *fiber.Ctx, message string, error interface{}) error {
    return c.Status(fiber.StatusNotFound).JSON(Response{
        Success: false,
        Message: message,
        Error:   error,
    })
}

func Conflict(c *fiber.Ctx, message string, error interface{}) error {
    return c.Status(fiber.StatusConflict).JSON(Response{
        Success: false,
        Message: message,
        Error:   error,
    })
}

func InternalServerError(c *fiber.Ctx, message string, error interface{}) error {
    return c.Status(fiber.StatusInternalServerError).JSON(Response{
        Success: false,
        Message: message,
        Error:   error,
    })
}
