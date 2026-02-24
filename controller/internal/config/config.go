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
	DBPath      string
	Port        string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AdminAPIKey: getEnv("ADMIN_API_KEY", "admin-secret"),
		AgentAPIKey: getEnv("AGENT_API_KEY", "agent-secret"),
		PollURL:     getEnv("POLL_URL", "/config"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		DBPath:      getEnv("DB_PATH", "controller.db"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
