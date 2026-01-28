package config

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

type Config struct {
	Environment Environment
	Host        string
	Port        string
	LogLevel    slog.Level
}

var (
	Global *Config
	once   sync.Once
)

func init() {
	once.Do(func() {
		Global = Load()
	})
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func loadBase() *Config {
	godotenv.Load()

	return &Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "8080"),
		LogLevel: func() slog.Level {
			switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
			case "debug":
				return slog.LevelDebug
			case "info":
				return slog.LevelInfo
			case "warn":
				return slog.LevelWarn
			case "error":
				return slog.LevelError
			default:
				return slog.LevelInfo
			}
		}(),
	}
}
