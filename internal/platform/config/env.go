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
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load("../../.env"); err != nil {
		return nil, err
	}

	slog.Info("env loaded successfully")

	return &Config{
		AppEnv:     os.Getenv("APP_ENV"),
		ServerPort: os.Getenv("SERVER_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}, nil
}
