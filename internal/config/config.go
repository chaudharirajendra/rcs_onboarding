package config

import (
	"os"
)

type Config struct {
	DSN    string
	JWTKey []byte
}

func LoadConfig() *Config {
	return &Config{
		DSN:    getEnv("DB_DSN", "new_user:password@tcp(localhost:3306)/rcs_onboarding?parseTime=true"),
		JWTKey: []byte(getEnv("JWT_SECRET", "secret_key")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
