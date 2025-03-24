package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	DBHost           string
	DBPort           string
	PostgresDB       string

	RefreshTokenExpiresIn time.Duration
	RefreshTokenMaxAge    int
	RefreshTokenSecret    string

	AccessTokenExpiresIn time.Duration
	AccessTokenSecret    string

	RedisHost string
	RedisPort string
	RedisDB   int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	requiredEnvVars := []string{
		"POSTGRES_USER", "POSTGRES_PASSWORD", "DB_HOST", "DB_PORT", "POSTGRES_DB",
		"REFRESH_TOKEN_EXPIRATION", "REFRESH_TOKEN_MAXAGE", "REFRESH_TOKEN_SECRET",
		"ACCESS_TOKEN_EXPIRATION", "ACCESS_TOKEN_SECRET",
		"REDIS_HOST", "REDIS_PORT", "REDIS_DB",
	}
	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", env)
		}
	}
	// Parse token expiration
	refreshtokenExpiration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRATION"))
	if err != nil {
		return &Config{}, fmt.Errorf("invalid format for REFRESH TOKEN_EXPIRATION: %v", err)
	}
	accesstokenExpiration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRATION"))
	if err != nil {
		return &Config{}, fmt.Errorf("invalid format for ACCESS TOKEN_EXPIRATION: %v", err)
	}
	// Parse token maxage
	refreshTokenMaxAge, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))
	if err != nil {
		return &Config{}, fmt.Errorf("invalid value for REFRESH_TOKEN_MAXAGE: %v", err)
	}

	// Parse REDIS_DB
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return &Config{}, fmt.Errorf("invalid value for REDIS_DB: %v", err)
	}
	return &Config{
		PostgresUser:          os.Getenv("POSTGRES_USER"),
		PostgresPassword:      os.Getenv("POSTGRES_PASSWORD"),
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		PostgresDB:            os.Getenv("POSTGRES_DB"),
		RefreshTokenExpiresIn: refreshtokenExpiration,
		RefreshTokenMaxAge:    refreshTokenMaxAge,
		RefreshTokenSecret:    os.Getenv("REFRESH_TOKEN_SECRET"),
		AccessTokenExpiresIn:  accesstokenExpiration,
		AccessTokenSecret:     os.Getenv("ACCESS_TOKEN_SECRET"),
		RedisHost:             os.Getenv("REDIS_HOST"),
		RedisPort:             os.Getenv("REDIS_PORT"),
		RedisDB:               redisDB,
	}, nil
}
