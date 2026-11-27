package env

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func InitEnv() error {
	return godotenv.Load()
}

func Env(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func EnvInt(key string, fallback int) int {
	valStr := Env(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

func EnvDuration(key string, fallback time.Duration) time.Duration {
	valStr := Env(key, "")
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return fallback
}