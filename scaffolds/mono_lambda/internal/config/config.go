package config

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Configuration holds the application configuration settings. The configuration is loaded from
// environment variables.
type Configuration struct {
	Env             string     `env:"ENV,required"`
	LogLevel        slog.Level `env:"LOG_LEVEL,required"`
	DBName          string     `env:"DATABASE_NAME"`
	DBUser          string     `env:"DATABASE_USER"`
	DBPassword      string     `env:"DATABASE_PASSWORD"`
	DBHost          string     `env:"DATABASE_HOST"`
	DBPort          string     `env:"DATABASE_PORT"`
	DBRetryDuration int        `env:"DATABASE_RETRY_DURATION_SECONDS"`
}

// New loads the configuration settings from environment variables and .env file, and returns a
// Configuration struct.
func New() (Configuration, error) {
	_ = godotenv.Load()

	cfg, err := env.ParseAs[Configuration]()
	if err != nil {
		return Configuration{}, fmt.Errorf("[in config.New] failed to parse config: %w", err)
	}

	return cfg, nil
}
