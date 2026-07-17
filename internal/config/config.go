package config

import (
	"errors"
	"log"
	"os"
	"time"
)

const defaultJWTSecret = "my_super_secret_key_change_me"

type Config struct {
	HTTPAddr string
	DBDSN    string

	AppEnv    string
	JWTSecret string
	JWTIssuer string
	JWTTTL    time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DBDSN:    getEnv("DB_DSN", "root:123456@tcp(127.0.0.1:3306)/login_system?charset=utf8mb4&parseTime=True&loc=Local"),

		AppEnv:    getEnv("APP_ENV", "development"),
		JWTSecret: getEnv("JWT_SECRET", defaultJWTSecret),
		JWTIssuer: getEnv("JWT_ISSUER", "login-system"),
		JWTTTL:    getDurationEnv("JWT_TTL", 24*time.Hour),
	}
}

func (c Config) Validate() error {
	if c.AppEnv == "production" && c.JWTSecret == defaultJWTSecret {
		return errors.New("生产环境必须通过 JWT_SECRET 配置安全密钥")
	}
	if c.DBDSN == "" {
		return errors.New("DB_DSN 不能为空")
	}
	if c.HTTPAddr == "" {
		return errors.New("HTTP_ADDR 不能为空")
	}
	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("配置 %s=%q 不是合法 duration，使用默认值 %s", key, value, fallback)
		return fallback
	}
	return duration
}
