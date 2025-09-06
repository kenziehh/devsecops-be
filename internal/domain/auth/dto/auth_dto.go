package dto

import (
    "github.com/google/uuid"
	"time"
)

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email" example:"user@example.com"`
    Password string `json:"password" validate:"required,min=6,max=100" example:"password123"`
}

type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
    Email    string `json:"email" validate:"required,email,max=255" example:"user@example.com"`
    Password string `json:"password" validate:"required,min=6,max=100" example:"password123"`
}

type AuthResponse struct {
    Token     string    `json:"token"`
    User      UserData  `json:"user"`
    ExpiresAt time.Time `json:"expires_at"`
}


type UserData struct {
    ID        uuid.UUID `db:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}