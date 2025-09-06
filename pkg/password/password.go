package password

import (
    "golang.org/x/crypto/bcrypt"
)

type PasswordUtil interface {
    HashPassword(password string) (string, error)
    CheckPassword(password, hash string) bool
}

type passwordUtil struct {
    cost int
}

func NewPasswordUtil() PasswordUtil {
    return &passwordUtil{
        cost: bcrypt.DefaultCost,
    }
}

func (p *passwordUtil) HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func (p *passwordUtil) CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}