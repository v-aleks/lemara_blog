package service

import (
	"context"
	"errors"
	"lemara_blog/internal/domain"
	"lemara_blog/internal/repository"
	"time"

	"github.com/google/uuid"
)

type PostService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

// Метод для создания новой статьи
func (s *PostService) CreatePost(ctx context.Context, req *domain.PostCreateRequest) (*domain.Post, error) {
	// Проверка на пустые поля
	if req.Title == "" || req.Content == "" || req.Author == "" {
		return nil, errors.New("Title, Content and Author fields are required")
	}
	// Проверка на существование статьи с таким ID
	post := domain.Post{
		ID:        uuid.New(),
		Title:     req.Title,
		Content:   req.Content,
		Author:    req.Author,
		CreatedAt: time.Now(),
	}
	err := s.repo.Create(ctx, post)
	if err != nil {
        return nil, err
    }
	return &post, err
}

// Метод для получения статьи по ID
func (s *PostService) GetPostByID(ctx context.Context, id uuid.UUID) (*domain.PostSearchResponse, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
        return nil, err
    }
	return &post, err
}
