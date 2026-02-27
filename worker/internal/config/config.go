package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	RequestTimeoutSeconds int
	AgentAPIKey           string
	GinMode               string
	Port                  string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		RequestTimeoutSeconds: getEnvInt("REQUEST_TIMEOUT_SECONDS"),
		AgentAPIKey:           os.Getenv("AGENT_API_KEY"),
		GinMode:               os.Getenv("GIN_MODE"),
		Port:                  os.Getenv("PORT"),
	}
}

func (c *Config) Validate() error {
	missing := make([]string, 0, 3)
	if strings.TrimSpace(c.AgentAPIKey) == "" {
		missing = append(missing, "AGENT_API_KEY")
	}
	if strings.TrimSpace(c.GinMode) == "" {
		missing = append(missing, "GIN_MODE")
	}
	if strings.TrimSpace(c.Port) == "" {
		missing = append(missing, "PORT")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env: %s", strings.Join(missing, ", "))
	}

	if c.RequestTimeoutSeconds <= 0 {
		return fmt.Errorf("invalid REQUEST_TIMEOUT_SECONDS: must be > 0")
	}

	return nil
}

func getEnvInt(k string) int {
	raw := os.Getenv(k)
	if raw == "" {
		return 0
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return v
}
