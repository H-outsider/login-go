package config

import (
	"testing"
	"time"
)

func TestValidateRequiresJWTSecretInProduction(t *testing.T) {
	cfg := Config{
		AppEnv:    "production",
		HTTPAddr:  ":8080",
		DBDSN:     "user:pass@tcp(localhost:3306)/app",
		JWTSecret: defaultJWTSecret,
		JWTTTL:    time.Hour,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want production secret error")
	}
}

func TestValidateAllowsDevelopmentDefaultSecret(t *testing.T) {
	cfg := Config{
		AppEnv:    "development",
		HTTPAddr:  ":8080",
		DBDSN:     "user:pass@tcp(localhost:3306)/app",
		JWTSecret: defaultJWTSecret,
		JWTTTL:    time.Hour,
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
