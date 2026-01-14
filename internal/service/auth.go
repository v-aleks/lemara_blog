package service

import (
	"context"
	"errors"
	"time"

	"lemara_blog/internal/config"
	"lemara_blog/internal/domain"
	"lemara_blog/internal/repository"
	"lemara_blog/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.AuthResponse, error)
    Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
    HashPassword(password string) (string, error)
    ComparePassword(hashedPassword, password string) error
}

type authService struct {
    userRepo repository.UserRepository
    config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
    return &authService{
        userRepo: userRepo,
        config:   config,
    }
}

func (s *authService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.AuthResponse, error) {
    // Check if user already exists
    existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("user already exists")
    }

    // Hash password
    hashedPassword, err := s.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    // Create user
    user := &domain.User{
        ID:           generateID(),
        Email:        req.Email,
        PasswordHash: hashedPassword,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTExpiration)
    if err != nil {
        return nil, err
    }

    return &domain.AuthResponse{
        Token: token,
        User: domain.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        },
    }, nil
}

func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
    // Find user by email
    user, err := s.userRepo.FindByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid credentials")
    }

    // Compare password
    if err := s.ComparePassword(user.PasswordHash, req.Password); err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTExpiration)
    if err != nil {
        return nil, err
    }

    return &domain.AuthResponse{
        Token: token,
        User: domain.UserResponse{
            ID:        user.ID,
            Email:     user.Email,
            CreatedAt: user.CreatedAt,
        },
    }, nil
}

func (s *authService) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.BcryptCost)
    return string(bytes), err
}

func (s *authService) ComparePassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateID() string {
    return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    result := make([]byte, n)
    for i := range result {
        result[i] = letters[time.Now().UnixNano()%int64(len(letters))]
    }
    return string(result)
}
