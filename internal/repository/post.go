package repository

import (
	"context"
	"time"

	"lemara_blog/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Интерфейс для работы с постами
type PostRepository interface {
	Create(ctx context.Context, post domain.Post) error
	//GetByID(ctx context.Context, id int64) (domain.Post, error)
	//GetByTitle(ctx context.Context, title string) (domain.Post, error)
	//Update(ctx context.Context, post domain.Post) error
	//Delete(ctx context.Context, id int64) error
	//List(ctx context.Context, limit, offset int) ([]domain.Post, error)
}

// Тип для работы с постами
type postRepository struct {
	pool *pgxpool.Pool
}

// Конструктор для создания нового экземпляра PostRepository
func NewPostRepository(pool *pgxpool.Pool) PostRepository {
	return &postRepository{pool: pool}
}

func (r *postRepository) Create(ctx context.Context, post domain.Post) error {
	query := `INSERT INTO posts (id, title, content, author, created_at) VALUES ($1, $2, $3, $4, $5)`

	post.CreatedAt = time.Now()

	_, err := r.pool.Exec(
		ctx,
		query,
		post.ID,
		post.Title,
		post.Content,
		post.Author,
		post.CreatedAt)


	return err
}
