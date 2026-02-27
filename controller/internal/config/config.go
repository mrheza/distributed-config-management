package config

import (
	"fmt"
	"os"
	"strings"

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

func (c *Config) Validate() error {
	missing := make([]string, 0, 6)
	if strings.TrimSpace(c.AdminAPIKey) == "" {
		missing = append(missing, "ADMIN_API_KEY")
	}
	if strings.TrimSpace(c.AgentAPIKey) == "" {
		missing = append(missing, "AGENT_API_KEY")
	}
	if strings.TrimSpace(c.PollURL) == "" {
		missing = append(missing, "POLL_URL")
	}
	if strings.TrimSpace(c.GinMode) == "" {
		missing = append(missing, "GIN_MODE")
	}
	if strings.TrimSpace(c.DatabaseURL) == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if strings.TrimSpace(c.Port) == "" {
		missing = append(missing, "PORT")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env: %s", strings.Join(missing, ", "))
	}

	return nil
}
