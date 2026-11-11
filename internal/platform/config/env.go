package config

import (
	"errors"
	"fmt"
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

	appEnv, err := getEnvString("APP_ENV", "development")
	if err != nil {
		return nil, fmt.Errorf("failed to get APP_ENV: %w", err)
	}

	serverPort, err := getEnvString("SERVER_PORT", "8080")
	if err != nil {
		return nil, fmt.Errorf("failed to get SERVER_PORT: %w", err)
	}

	dbName, err := getEnvString("POSTGRES_DB", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get POSTGRES_DB: %w", err)
	}

	dbUser, err := getEnvString("POSTGRES_USER", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get POSTGRES_USER: %w", err)
	}

	dbPassword, err := getEnvString("POSTGRES_PASSWORD", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get POSTGRES_PASSWORD: %w", err)
	}

	dbHost, err := getEnvString("POSTGRES_HOST", "localhost")
	if err != nil {
		return nil, fmt.Errorf("failed to get POSTGRES_HOST: %w", err)
	}

	dbPort, err := getEnvInt("POSTGRES_PORT", 5432)
	if err != nil {
		return nil, fmt.Errorf("failed to parse POSTGRES_PORT: %w", err)
	}

	redisHost, err := getEnvString("REDIS_HOST", "localhost")
	if err != nil {
		return nil, fmt.Errorf("failed to get REDIS_HOST: %w", err)
	}

	redisPort, err := getEnvString("REDIS_PORT", "6379")
	if err != nil {
		return nil, fmt.Errorf("failed to get REDIS_PORT: %w", err)
	}

	redisPassword, err := getEnvString("REDIS_PASSWORD", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get REDIS_PASSWORD: %w", err)
	}

	redisDB, err := getEnvInt("REDIS_DB", 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse REDIS_DB: %w", err)
	}

	return &Config{
		AppEnv:        appEnv,
		ServerPort:    serverPort,
		DBName:        dbName,
		DBUser:        dbUser,
		DBPassword:    dbPassword,
		DBHost:        dbHost,
		DBPort:        dbPort,
		RedisHost:     redisHost,
		RedisPort:     redisPort,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
	}, nil
}

func getEnvString(key, fallback string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		if key == "POSTGRES_DB" || key == "POSTGRES_USER" {
			if value == "" {
				return "", fmt.Errorf("invalid configuration value: %v can not empty", key)
			}
		}

		return value, nil
	}
	return fallback, nil
}

func getEnvInt(key string, fallback int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		res, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid configuration value: %v is not a valid %v", value, key)
		}

		// Port validation (recommendation: 1024 - 65535)
		if strings.Contains(key, "PORT") && (res < 1024 || res > 65535) {
			return 0, fmt.Errorf("invalid port value: port %v is invalid, please use port 1024 - 65535", res)
		}

		// Redis DB value validation
		if key == "REDIS_DB" && (res < 0 || res > 15) {
			return 0, fmt.Errorf("invalid redis db value: %v is not a valid redis db (use 0-15)", res)
		}

		return res, nil
	}
	return fallback, nil
}
