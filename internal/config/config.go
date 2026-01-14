package config

import (
	"os"
	"strconv"
	"time"
)

// Конфигурация приложения
type Config struct {
    DBHost         string
    DBPort         string
    DBUser         string
    DBPassword     string
    DBName         string
    DBSSLMode      string
    ServerPort     string
    JWTSecret      string
    JWTExpiration  time.Duration
    BcryptCost     int
}

func Load() *Config {
	jwtExpiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
    bcryptCost, _ := strconv.Atoi(getEnv("BCRYPT_COST", "10"))

    return &Config{
            DBHost:        getEnv("DB_HOST", "localhost"),
            DBPort:        getEnv("DB_PORT", "5432"),
            DBUser:        getEnv("DB_USER", "postgres"),
            DBPassword:    getEnv("DB_PASSWORD", "postgres"),
            DBName:        getEnv("DB_NAME", "myapp"),
            DBSSLMode:     getEnv("DB_SSL_MODE", "disable"),
            ServerPort:    getEnv("SERVER_PORT", "8080"),
            JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
            JWTExpiration: time.Duration(jwtExpiration) * time.Hour,
            BcryptCost:    bcryptCost,
        }
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
