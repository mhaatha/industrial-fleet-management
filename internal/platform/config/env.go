package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	ServerPort string

	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	slog.Info("env loaded successfully")

	return &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBName:     getEnv("DB_NAME", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
