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
		ControllerBaseURL:     getEnv("CONTROLLER_BASE_URL", "http://localhost:8080"),
		ControllerAPIKey:      getEnv("CONTROLLER_API_KEY", "agent-secret"),
		WorkerBaseURL:         getEnv("WORKER_BASE_URL", "http://localhost:8082"),
		WorkerAPIKey:          getEnv("WORKER_API_KEY", "worker-secret"),
		PollURL:               getEnv("POLL_URL", "/config"),
		PollIntervalSeconds:   getEnvInt("POLL_INTERVAL_SECONDS", 30),
		StatePath:             getEnv("STATE_PATH", "data/agent_state.json"),
		MaxBackoffSeconds:     getEnvInt("MAX_BACKOFF_SECONDS", 60),
		BackoffJitterPercent:  getEnvInt("BACKOFF_JITTER_PERCENT", 20),
		RequestTimeoutSeconds: getEnvInt("REQUEST_TIMEOUT_SECONDS", 10),
		GinMode:               getEnv("GIN_MODE", "release"),
		Port:                  getEnv("PORT", "8081"),
	}
}

func getEnv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func getEnvInt(k string, d int) int {
	raw := os.Getenv(k)
	if raw == "" {
		return d
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return d
	}
	return v
}
