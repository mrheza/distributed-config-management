package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ControllerBaseURL     string
	ControllerAPIKey      string
	WorkerBaseURL         string
	WorkerAPIKey          string
	PollURL               string
	PollIntervalSeconds   int
	StatePath             string
	MaxBackoffSeconds     int
	BackoffJitterPercent  int
	RequestTimeoutSeconds int
	GinMode               string
	Port                  string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		ControllerBaseURL:     os.Getenv("CONTROLLER_BASE_URL"),
		ControllerAPIKey:      os.Getenv("CONTROLLER_API_KEY"),
		WorkerBaseURL:         os.Getenv("WORKER_BASE_URL"),
		WorkerAPIKey:          os.Getenv("WORKER_API_KEY"),
		PollURL:               os.Getenv("POLL_URL"),
		PollIntervalSeconds:   getEnvInt("POLL_INTERVAL_SECONDS"),
		StatePath:             os.Getenv("STATE_PATH"),
		MaxBackoffSeconds:     getEnvInt("MAX_BACKOFF_SECONDS"),
		BackoffJitterPercent:  getEnvInt("BACKOFF_JITTER_PERCENT"),
		RequestTimeoutSeconds: getEnvInt("REQUEST_TIMEOUT_SECONDS"),
		GinMode:               os.Getenv("GIN_MODE"),
		Port:                  os.Getenv("PORT"),
	}
}

func (c *Config) Validate() error {
	missing := make([]string, 0, 8)
	if strings.TrimSpace(c.ControllerBaseURL) == "" {
		missing = append(missing, "CONTROLLER_BASE_URL")
	}
	if strings.TrimSpace(c.ControllerAPIKey) == "" {
		missing = append(missing, "CONTROLLER_API_KEY")
	}
	if strings.TrimSpace(c.WorkerBaseURL) == "" {
		missing = append(missing, "WORKER_BASE_URL")
	}
	if strings.TrimSpace(c.WorkerAPIKey) == "" {
		missing = append(missing, "WORKER_API_KEY")
	}
	if strings.TrimSpace(c.PollURL) == "" {
		missing = append(missing, "POLL_URL")
	}
	if strings.TrimSpace(c.StatePath) == "" {
		missing = append(missing, "STATE_PATH")
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

	if c.PollIntervalSeconds <= 0 {
		return fmt.Errorf("invalid POLL_INTERVAL_SECONDS: must be > 0")
	}
	if c.MaxBackoffSeconds <= 0 {
		return fmt.Errorf("invalid MAX_BACKOFF_SECONDS: must be > 0")
	}
	if c.BackoffJitterPercent < 0 || c.BackoffJitterPercent > 90 {
		return fmt.Errorf("invalid BACKOFF_JITTER_PERCENT: must be between 0 and 90")
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
