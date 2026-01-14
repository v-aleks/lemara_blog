package repository

import (
	"context"
	"errors"
	"time"

	"lemara_blog/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Интерфейс репозитория пользователей
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id string) error
}

type userRepository struct {
    pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
    return &userRepository{pool: pool}
}

// Создание нового пользователя
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, first_name, last_name, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

    user.CreatedAt = time.Now()
    user.UpdatedAt = user.CreatedAt

    _, err := r.pool.Exec(ctx, query,
        user.ID,
        user.Email,
        user.FirstName,
        user.LastName,
        user.PasswordHash,
        user.CreatedAt,
        user.UpdatedAt,
    )

    return err
}

// Поиск пользователя по ID
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    query := `
        SELECT id, email, first_name, last_name, password_hash, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

    var user domain.User
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.FirstName,
        &user.LastName,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, nil
    }

    return &user, err
}

// Поиск пользователя по email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, email, first_name, last_name, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1 AND deleted_at IS NULL
    `

    var user domain.User
    err := r.pool.QueryRow(ctx, query, email).Scan(
        &user.ID,
        &user.Email,
        &user.FirstName,
        &user.LastName,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, nil
    }

    return &user, err
}

// Обновление пользователя
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users
        SET email = $1, first_name = $2, last_name = $3, password_hash = $4, updated_at = $5
        WHERE id = $6 AND deleted_at IS NULL
    `

    user.UpdatedAt = time.Now()
    _, err := r.pool.Exec(ctx, query,
        user.Email,
        user.FirstName,
        user.LastName,
        user.PasswordHash,
        user.UpdatedAt,
        user.ID,
    )

    return err
}


// Удаление пользователя
func (r *userRepository) Delete(ctx context.Context, id string) error {
    query := `
        UPDATE users
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `

    _, err := r.pool.Exec(ctx, query, time.Now(), id)
    return err
}
