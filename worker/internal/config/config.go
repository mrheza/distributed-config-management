package config

import (
	"os"
	"strconv"

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