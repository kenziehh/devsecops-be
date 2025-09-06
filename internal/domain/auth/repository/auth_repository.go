package repository

import (
	"context"
	"database/sql"
	"devsecops-be/internal/domain/auth/dto"
	"devsecops-be/pkg/errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AuthRepository interface {
    CreateUser(ctx context.Context, user dto.RegisterRequest, hashedPassword string) (*dto.UserData, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserData, error)
}

type User struct {
    ID        uuid.UUID `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    Password  string    `db:"password"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type authRepository struct {
    db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
    return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, user dto.RegisterRequest, hashedPassword string) (*dto.UserData, error) {
    query := `
        INSERT INTO users (name, email, password) 
        VALUES ($1, $2, $3) 
        RETURNING id, name, email, created_at
    `
    
    var userData dto.UserData
    err := r.db.QueryRowContext(ctx, query, user.Name, user.Email, hashedPassword).
        Scan(&userData.ID, &userData.Name, &userData.Email, &userData.CreatedAt)
    
    if err != nil {
        // Handle unique constraint violation (duplicate email)
        if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
            return nil, errors.ErrUserAlreadyExists
        }
        return nil, errors.WrapDatabaseError(err, "failed to create user")
    }
    
    return &userData, nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    query := `
        SELECT id, name, email, password, created_at, updated_at 
        FROM users 
        WHERE email = $1
    `
    
    var user User
    err := r.db.QueryRowContext(ctx, query, email).
        Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.ErrUserNotFound
        }
        return nil, errors.WrapDatabaseError(err, "failed to get user by email")
    }
    
    return &user, nil
}

func (r *authRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserData, error) {
    query := `
        SELECT id, name, email, created_at 
        FROM users 
        WHERE id = $1
    `
    
    var user dto.UserData
    err := r.db.QueryRowContext(ctx, query, id).
        Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.ErrUserNotFound
        }
        return nil, errors.WrapDatabaseError(err, "failed to get user by ID")
    }
    
    return &user, nil
}
