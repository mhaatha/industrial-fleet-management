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

	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     int

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
}

func LoadConfig() (*Config, error) {
	// isError is a variable
	isError := false

	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	appEnv, err := getEnvString("APP_ENV", "development")
	var appEnvErr error
	if err != nil {
		isError = true
		appEnvErr = fmt.Errorf("failed to get APP_ENV: %w", err)
	}

	serverPort, err := getEnvString("SERVER_PORT", "8080")
	var serverPortErr error
	if err != nil {
		isError = true
		serverPortErr = fmt.Errorf("failed to get SERVER_PORT: %w", err)
	}

	postgresDB, err := getEnvString("POSTGRES_DB", "")
	var postgresDBErr error
	if err != nil {
		isError = true
		postgresDBErr = fmt.Errorf("failed to get POSTGRES_DB: %w", err)
	}

	postgresUser, err := getEnvString("POSTGRES_USER", "")
	var postgresUserErr error
	if err != nil {
		isError = true
		postgresUserErr = fmt.Errorf("failed to get POSTGRES_USER: %w", err)
	}

	postgresPassword, err := getEnvString("POSTGRES_PASSWORD", "")
	var postgresPasswordErr error
	if err != nil {
		isError = true
		postgresPasswordErr = fmt.Errorf("failed to get POSTGRES_PASSWORD: %w", err)
	}

	postgresHost, err := getEnvString("POSTGRES_HOST", "localhost")
	var postgresHostErr error
	if err != nil {
		isError = true
		postgresHostErr = fmt.Errorf("failed to get POSTGRES_HOST: %w", err)
	}

	postgresPort, err := getEnvInt("POSTGRES_PORT", 5432)
	var postgresPortErr error
	if err != nil {
		isError = true
		postgresPortErr = fmt.Errorf("failed to parse POSTGRES_PORT: %w", err)
	}

	redisHost, err := getEnvString("REDIS_HOST", "localhost")
	var redisHostErr error
	if err != nil {
		isError = true
		redisHostErr = fmt.Errorf("failed to get REDIS_HOST: %w", err)
	}

	redisPort, err := getEnvInt("REDIS_PORT", 6379)
	var redisPortErr error
	if err != nil {
		isError = true
		redisPortErr = fmt.Errorf("failed to get REDIS_PORT: %w", err)
	}

	redisPassword, err := getEnvString("REDIS_PASSWORD", "")
	var redisPasswordErr error
	if err != nil {
		isError = true
		redisPasswordErr = fmt.Errorf("failed to get REDIS_PASSWORD: %w", err)
	}

	redisDB, err := getEnvInt("REDIS_DB", 0)
	var redisDBErr error
	if err != nil {
		isError = true
		redisDBErr = fmt.Errorf("failed to parse REDIS_DB: %w", err)
	}

	if isError {
		return nil, errors.Join(
			appEnvErr,
			serverPortErr,
			postgresDBErr,
			postgresUserErr,
			postgresPasswordErr,
			postgresHostErr,
			postgresPortErr,
			redisHostErr,
			redisPortErr,
			redisPasswordErr,
			redisDBErr,
		)
	}

	return &Config{
		AppEnv:           appEnv,
		ServerPort:       serverPort,
		PostgresDB:       postgresDB,
		PostgresUser:     postgresUser,
		PostgresPassword: postgresPassword,
		PostgresHost:     postgresHost,
		PostgresPort:     postgresPort,
		RedisHost:        redisHost,
		RedisPort:        redisPort,
		RedisPassword:    redisPassword,
		RedisDB:          redisDB,
	}, nil
}

func getEnvString(key, fallback string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		if key == "POSTGRES_DB" || key == "POSTGRES_USER" || key == "POSTGRES_PASSWORD" {
			if value == "" {
				return "", fmt.Errorf("invalid configuration value: %v can not empty", key)
			}
		}

		return value, nil
	} else {
		if key == "POSTGRES_DB" || key == "POSTGRES_USER" || key == "POSTGRES_PASSWORD" {
			if value == "" {
				return "", fmt.Errorf("invalid configuration value: %v can not empty", key)
			}
		}
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
