package model

type State struct {
	AgentID             string `json:"agent_id"`
	ETag                string `json:"etag"`
	ConfigURL           string `json:"config_url"`
	PollURL             string `json:"poll_url"`
	PollIntervalSeconds int    `json:"poll_interval_seconds"`
	LastConfigVersion   int    `json:"last_config_version"`
}
