package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	ServerPort string

	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     int

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	return &Config{
		AppEnv:        getEnvString("APP_ENV", "development"),
		ServerPort:    getEnvString("SERVER_PORT", "8080"),
		DBName:        getEnvString("POSTGRES_DB", ""),
		DBUser:        getEnvString("POSTGRES_USER", ""),
		DBPassword:    getEnvString("POSTGRES_PASSWORD", ""),
		DBHost:        getEnvString("POSTGRES_HOST", "localhost"),
		DBPort:        getEnvInt("POSTGRES_PORT", 5432),
		RedisHost:     getEnvString("REDIS_HOST", "localhost"),
		RedisPort:     getEnvString("REDIS_PORT", "6379"),
		RedisPassword: getEnvString("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
	}, nil
}

func getEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		if key == "POSTGRES_DB" || key == "POSTGRES_USER" {
			if value == "" {
				slog.Error("invalid configuration value", "err", fmt.Sprintf("%v can not empty", key))
				os.Exit(1)
			}
		}

		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		res, err := strconv.Atoi(value)
		if err != nil {
			slog.Error("invalid configuration value", "err", fmt.Sprintf("%v is not a valid %v", value, key))
			os.Exit(1)
		}

		// Port validation (recommendation: 1024 - 65535)
		if strings.Contains(key, "PORT") && (res < 1024 || res > 65535) {
			slog.Error("invalid port value", "err", fmt.Sprintf("port %v is not safe, recommended port are 1024 - 65535", res))
			os.Exit(1)
		}

		// Redis DB value validation
		if key == "REDIS_DB" && (res < 0 || res > 15) {
			slog.Error("invalid redis db value", "err", fmt.Sprintf("%v is not a valid redis db (use 0-15)", res))
			os.Exit(1)
		}

		return res
	}
	return fallback
}
