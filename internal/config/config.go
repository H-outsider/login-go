package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr string
	DBDSN    string

	JWTSecret string
	JWTIssuer string
	JWTTTL    time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DBDSN:    getEnv("DB_DSN", "root:123456@tcp(127.0.0.1:3306)/login_system?charset=utf8mb4&parseTime=True&loc=Local"),

		JWTSecret: getEnv("JWT_SECRET", "my_super_secret_key_change_me"),
		JWTIssuer: getEnv("JWT_ISSUER", "login-system"),
		JWTTTL:    24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
