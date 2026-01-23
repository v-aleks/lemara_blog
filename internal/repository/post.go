package repository

import (
	"context"
	"lemara_blog/internal/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Интерфейс для работы с постами
type PostRepository interface {
	Create(ctx context.Context, post domain.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.PostSearchResponse, error)
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

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.PostSearchResponse, error) {
	query := `
				SELECT
					posts.id,
					posts.title,
					posts.content,
					users.id AS author_id, users.email as author_email, users.first_name as author_first_name, users.last_name as author_last_name,
					posts.created_at, posts.updated_at
				FROM posts
				JOIN users ON posts.author = users.id
				WHERE posts.id = $1
			`

	var post domain.PostSearchResponse
	// Вот тут добавить выбор нужных полей в запрос. Добавить в Author.Id поля из UserResponse
	err := r.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Author.ID,
		&post.Author.Email,
		&post.Author.FirstName,
		&post.Author.LastName,
		&post.CreatedAt,
		&post.UpdatedAt,

	)

	if err != nil {
		return domain.PostSearchResponse{}, err
	}

	return post, nil
}
