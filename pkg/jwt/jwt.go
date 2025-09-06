package jwt

import (
	"devsecops-be/pkg/errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTUtil interface {
    GenerateToken(userID uuid.UUID) (string, time.Time, error)
    ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type jwtUtil struct {
    secretKey      []byte
    accessTokenExp time.Duration
}

func NewJWTUtil() JWTUtil {
    secretKey := os.Getenv("JWT_SECRET_KEY")
    

    accessExp := 6 * time.Hour // Default 6 hours

    if exp := os.Getenv("JWT_ACCESS_EXP_HOURS"); exp != "" {
        if hours, err := strconv.Atoi(exp); err == nil {
            accessExp = time.Duration(hours) * time.Hour
        }
    }

    return &jwtUtil{
        secretKey:      []byte(secretKey),
        accessTokenExp: accessExp,
    }
}

func (j *jwtUtil) GenerateToken(userID uuid.UUID) (string, time.Time, error) {
    expiresAt := time.Now().Add(j.accessTokenExp)
    
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     expiresAt.Unix(),
        "iat":     time.Now().Unix(),
        "type":    "access",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(j.secretKey)
    if err != nil {
        return "", time.Time{}, err
    }
    
    return tokenString, expiresAt, nil
}

func (j *jwtUtil) ValidateToken(tokenString string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.ErrInvalidToken
        }
        return j.secretKey, nil
    })

    if err != nil {
        return nil, errors.ErrInvalidToken
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.ErrInvalidToken
}