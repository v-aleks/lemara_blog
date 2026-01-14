package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"lemara_blog/internal/config"

	"lemara_blog/internal/handler"
	"lemara_blog/internal/repository"
	"lemara_blog/internal/service"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Load configuration
    cfg := config.Load()

    // Setup database connection
    dbPool, err := setupDatabase(cfg)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    defer dbPool.Close()

    // Initialize repositories
    userRepo := repository.NewUserRepository(dbPool)

    // Initialize services
    authService := service.NewAuthService(userRepo, &config.Config{
        JWTSecret:     cfg.JWTSecret,
        JWTExpiration: cfg.JWTExpiration,
        BcryptCost:    cfg.BcryptCost,
    })

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)
    userHandler := handler.NewUserHandler(userRepo)
    healthHandler := handler.NewHealthHandler(dbPool)

    // Setup router
    mux := http.NewServeMux()

    // Public routes
    mux.HandleFunc("POST /auth/register", authHandler.Register)
    mux.HandleFunc("POST /auth/login", authHandler.Login)
    mux.HandleFunc("GET /health", healthHandler.Check)

    // Protected routes (with auth middleware)
    protected := http.NewServeMux()
    protected.HandleFunc("GET /api/users/me", userHandler.GetProfile)
    protected.HandleFunc("PUT /api/users/me", userHandler.UpdateProfile) // Остановился тут. Нужно реализовать изменение профиля
    protected.HandleFunc("DELETE /api/users/me", userHandler.DeleteProfile)

    // Вот тут важно подключить защищенные роуты к mux
    mux.Handle("/api/", handler.AuthMiddleware(cfg.JWTSecret)(protected))

    // Setup server
    server := &http.Server{
        Addr:         ":" + cfg.ServerPort,
        Handler:      mux,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Graceful shutdown
    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Printf("Server starting on port %s", cfg.ServerPort)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    <-done
    log.Println("Server shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown failed: %v", err)
    }

    log.Println("Server stopped")
}

func setupDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
    connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        cfg.DBUser,
        cfg.DBPassword,
        cfg.DBHost,
        cfg.DBPort,
        cfg.DBName,
        cfg.DBSSLMode,
    )

    poolConfig, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        return nil, err
    }

    // Connection pool settings
    poolConfig.MaxConns = 25
    poolConfig.MinConns = 5
    poolConfig.MaxConnLifetime = time.Hour
    poolConfig.MaxConnIdleTime = 30 * time.Minute

    return pgxpool.NewWithConfig(context.Background(), poolConfig)
}
