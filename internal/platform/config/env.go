package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	ServerPort string

	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	slog.Info("env loaded successfully")

	return &Config{
		AppEnv:        getEnvString("APP_ENV", "development"),
		ServerPort:    getEnvString("SERVER_PORT", "8080"),
		DBName:        getEnvString("DB_NAME", ""),
		DBUser:        getEnvString("DB_USER", ""),
		DBPassword:    getEnvString("DB_PASSWORD", ""),
		DBHost:        getEnvString("DB_HOST", "localhost"),
		DBPort:        getEnvString("DB_PORT", "5432"),
		RedisHost:     getEnvString("REDIS_HOST", "localhost"),
		RedisPort:     getEnvString("REDIS_PORT", "6379"),
		RedisPassword: getEnvString("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
	}, nil
}

func getEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		res, _ := strconv.Atoi(value)

		return res
	}
	return fallback
}
