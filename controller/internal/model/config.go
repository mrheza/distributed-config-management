package model

type Config struct {
	Version             int    `json:"version"`
	URL                 string `json:"url"`
	PollIntervalSeconds int    `json:"poll_interval_seconds"`
}
