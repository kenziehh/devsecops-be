package service

import (
    "context"
    "devsecops-be/internal/domain/auth/dto"
    "devsecops-be/internal/domain/auth/repository"
    "devsecops-be/pkg/errors"
    "devsecops-be/pkg/jwt"
    "devsecops-be/pkg/logger"
    "devsecops-be/pkg/password"
)

type AuthService interface {
    Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
    Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
}

type authService struct {
    authRepo repository.AuthRepository
    jwtUtil  jwt.JWTUtil
    passUtil password.PasswordUtil
    logger   logger.Logger
}

func NewAuthService(
    authRepo repository.AuthRepository, 
    jwtUtil jwt.JWTUtil, 
    passUtil password.PasswordUtil,
    logger logger.Logger,
) AuthService {
    return &authService{
        authRepo: authRepo,
        jwtUtil:  jwtUtil,
        passUtil: passUtil,
        logger:   logger,
    }
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
    s.logger.Info(ctx, "Attempting user login", logger.Fields{
        "email": req.Email,
    })

    user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
    if err != nil {
        if err == errors.ErrUserNotFound {
            s.logger.Warn(ctx, "Login attempt with non-existent email", logger.Fields{
                "email": req.Email,
            })
            return nil, errors.ErrInvalidCredentials
        }
        s.logger.Error(ctx, "Failed to get user during login", err, logger.Fields{
            "email": req.Email,
        })
        return nil, err
    }

    if !s.passUtil.CheckPassword(req.Password, user.Password) {
        s.logger.Warn(ctx, "Login attempt with invalid password", logger.Fields{
            "user_id": user.ID,
            "email":   req.Email,
        })
        return nil, errors.ErrInvalidCredentials
    }

    userData := dto.UserData{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    }

    token, expiresAt, err := s.jwtUtil.GenerateToken(user.ID)
    if err != nil {
        s.logger.Error(ctx, "Failed to generate JWT token", err, logger.Fields{
            "user_id": user.ID,
        })
        return nil, errors.WrapInternalError(err, "failed to generate token")
    }

    s.logger.Info(ctx, "User login successful", logger.Fields{
        "user_id": user.ID,
        "email":   user.Email,
    })

    return &dto.AuthResponse{
        Token:     token,
        User:      userData,
        ExpiresAt: expiresAt,
    }, nil
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
    s.logger.Info(ctx, "Attempting user registration", logger.Fields{
        "name":  req.Name,
        "email": req.Email,
    })

    hashedPassword, err := s.passUtil.HashPassword(req.Password)
    if err != nil {
        s.logger.Error(ctx, "Failed to hash password during registration", err, logger.Fields{
            "email": req.Email,
        })
        return nil, errors.WrapInternalError(err, "failed to process password")
    }

    userData, err := s.authRepo.CreateUser(ctx, req, hashedPassword)
    if err != nil {
        if err == errors.ErrUserAlreadyExists {
            s.logger.Warn(ctx, "Registration attempt with existing email", logger.Fields{
                "email": req.Email,
            })
        } else {
            s.logger.Error(ctx, "Failed to create user during registration", err, logger.Fields{
                "name":  req.Name,
                "email": req.Email,
            })
        }
        return nil, err
    }

    token, expiresAt, err := s.jwtUtil.GenerateToken(userData.ID)
    if err != nil {
        s.logger.Error(ctx, "Failed to generate JWT token during registration", err, logger.Fields{
            "user_id": userData.ID,
        })
        return nil, errors.WrapInternalError(err, "failed to generate token")
    }

    s.logger.Info(ctx, "User registration successful", logger.Fields{
        "user_id": userData.ID,
        "email":   userData.Email,
    })

    return &dto.AuthResponse{
        Token:     token,
        User:      *userData,
        ExpiresAt: expiresAt,
    }, nil
}