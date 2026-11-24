package config

import (
	"errors"
	"strconv"
	"testing"

	"github.com/mhaatha/industrial-fleet-management/internal/platform/apperror"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name             string
		appEnv           string
		serverPort       interface{}
		postgresDB       string
		postgresUser     string
		postgresPassword string
		postgresHost     string
		postgresPort     interface{}
		redisHost        string
		redisPort        interface{}
		redisPassword    string
		redisDB          interface{}
		expectErr        bool
		expectedMatch    error
	}{
		{
			name:             "all config is set",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        false,
			expectedMatch:    nil,
		},
		{
			name:             "POSTGRES_DB is not set",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrEnvRequired,
		},
		{
			name:             "POSTGRES_USER is not set",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrEnvRequired,
		},
		{
			name:             "POSTGRES_PASSWORD is not set",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrEnvRequired,
		},
		{
			name:             "SERVER_PORT is not a valid int",
			appEnv:           "production",
			serverPort:       "INVALID_PORT",
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrInvalidDataType,
		},
		{
			name:             "POSTGRES_PORT is not a valid int",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     "INVALID_PORT",
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrInvalidDataType,
		},
		{
			name:             "REDIS_PORT is not a valid int",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        "INVALID_PORT",
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrInvalidDataType,
		},
		{
			name:             "REDIS_DB is not a valid int",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          "INVALID_DB",
			expectErr:        true,
			expectedMatch:    apperror.ErrInvalidDataType,
		},
		{
			name:             "SERVER_PORT is out of range",
			appEnv:           "production",
			serverPort:       999999,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrOutOfRangePort,
		},
		{
			name:             "POSTGRES_PORT is out of range",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     999999,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrOutOfRangePort,
		},
		{
			name:             "REDIS_PORT is out of range",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        999999,
			redisPassword:    "test_password",
			redisDB:          0,
			expectErr:        true,
			expectedMatch:    apperror.ErrOutOfRangePort,
		},
		{
			name:             "REDIS_DB is out of range",
			appEnv:           "production",
			serverPort:       8080,
			postgresDB:       "test_db",
			postgresUser:     "test_user",
			postgresPassword: "test_password",
			postgresHost:     "localhost",
			postgresPort:     5432,
			redisHost:        "localhost",
			redisPort:        6379,
			redisPassword:    "test_password",
			redisDB:          999999,
			expectErr:        true,
			expectedMatch:    apperror.ErrOutOfRangeRedisDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.appEnv)
			t.Setenv("POSTGRES_DB", tt.postgresDB)
			t.Setenv("POSTGRES_USER", tt.postgresUser)
			t.Setenv("POSTGRES_PASSWORD", tt.postgresPassword)
			t.Setenv("POSTGRES_HOST", tt.postgresHost)
			t.Setenv("REDIS_HOST", tt.redisHost)
			t.Setenv("REDIS_PASSWORD", tt.redisPassword)

			// SERVER_PORT setup
			intServerPort, ok := tt.serverPort.(int)
			if !ok {
				strServerPort, ok := tt.serverPort.(string)
				if !ok {
					t.Logf("it should be an int or a string (for test purpose) but it's not one of them: %v", tt.serverPort)
				}

				t.Setenv("SERVER_PORT", strServerPort)
			} else {
				t.Setenv("SERVER_PORT", strconv.Itoa(intServerPort))
			}

			// POSTGRES_PORT setup
			intPostgresPort, ok := tt.postgresPort.(int)
			if !ok {
				strPostgresPort, ok := tt.postgresPort.(string)
				if !ok {
					t.Logf("it should be an int or a string (for test purpose) but it's not one of them: %v", tt.postgresPort)
				}

				t.Setenv("POSTGRES_PORT", strPostgresPort)
			} else {
				t.Setenv("POSTGRES_PORT", strconv.Itoa(intPostgresPort))
			}

			// REDIS_PORT setup
			intRedisPort, ok := tt.redisPort.(int)
			if !ok {
				strRedisPort, ok := tt.redisPort.(string)
				if !ok {
					t.Logf("it should be an int or a string (for test purpose) but it's not one of them: %v", tt.redisPort)
				}

				t.Setenv("REDIS_PORT", strRedisPort)
			} else {
				t.Setenv("REDIS_PORT", strconv.Itoa(intRedisPort))
			}

			// REDIS_DB setup
			intRedisDB, ok := tt.redisDB.(int)
			if !ok {
				strRedisDB, ok := tt.redisDB.(string)
				if !ok {
					t.Logf("it should be an int or a string (for test purpose) but it's not one of them: %v", tt.redisDB)
				}

				t.Setenv("REDIS_DB", strRedisDB)
			} else {
				t.Setenv("REDIS_DB", strconv.Itoa(intRedisDB))
			}

			cfg, err := LoadConfig()
			if err != nil {
				if (err != nil) != tt.expectErr {
					t.Errorf("it shouldn't return error but got an error: %v", err)
				}

				if !errors.Is(err, tt.expectedMatch) {
					t.Errorf("expected error '%v' got '%v'", tt.expectedMatch, err)
				}
			} else {
				if (err != nil) != tt.expectErr {
					t.Errorf("it should return an error but got no error")
				}

				if tt.appEnv != cfg.AppEnv {
					t.Errorf("got APP_ENV value of '%v' want '%v'", cfg.AppEnv, tt.appEnv)
				}

				if tt.serverPort != cfg.ServerPort {
					t.Errorf("got SERVER_PORT value of '%v' want '%v'", cfg.ServerPort, tt.serverPort)
				}

				if tt.postgresDB != cfg.PostgresDB {
					t.Errorf("got POSTGRES_DB value of '%v' want '%v'", cfg.PostgresDB, tt.postgresDB)
				}

				if tt.postgresUser != cfg.PostgresUser {
					t.Errorf("got POSTGRES_USER value of '%v' want '%v'", cfg.PostgresUser, tt.postgresUser)
				}

				if tt.postgresPassword != cfg.PostgresPassword {
					t.Errorf("got POSTGRES_PASSWORD value of '%v' want '%v'", cfg.PostgresPassword, tt.postgresPassword)
				}

				if tt.postgresHost != cfg.PostgresHost {
					t.Errorf("got POSTGRES_HOST value of '%v' want '%v'", cfg.PostgresHost, tt.postgresHost)
				}

				if tt.postgresPort != cfg.PostgresPort {
					t.Errorf("got POSTGRES_PORT value of '%v' want '%v'", cfg.PostgresPort, tt.postgresPort)
				}

				if tt.redisHost != cfg.RedisHost {
					t.Errorf("got REDIS_HOST value of '%v' want '%v'", cfg.RedisHost, tt.redisHost)
				}

				if tt.redisPort != cfg.RedisPort {
					t.Errorf("got REDIS_PORT value of '%v' want '%v'", cfg.RedisPort, tt.redisPort)
				}

				if tt.redisPassword != cfg.RedisPassword {
					t.Errorf("got REDIS_PASSWORD value of '%v' want '%v'", cfg.RedisPassword, tt.redisPassword)
				}

				if tt.redisDB != cfg.RedisDB {
					t.Errorf("got REDIS_DB value of '%v' want '%v'", cfg.RedisDB, tt.redisDB)
				}
			}
		})
	}
}
