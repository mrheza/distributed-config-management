package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AdminAPIKey string
	AgentAPIKey string
	PollURL     string
	GinMode     string
	DatabaseURL string
	Port        string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AdminAPIKey: os.Getenv("ADMIN_API_KEY"),
		AgentAPIKey: os.Getenv("AGENT_API_KEY"),
		PollURL:     os.Getenv("POLL_URL"),
		GinMode:     os.Getenv("GIN_MODE"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}
}
