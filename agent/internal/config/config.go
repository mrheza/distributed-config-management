package config

import (
	"os"
	"strconv"

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