package api

import (
	"github.com/go-resty/resty/v2"

	"sentineldesk/agent/internal/config"
)

func NewClient() *resty.Client {
	cfg := config.Get()

	client := resty.New()
	client.SetBaseURL(cfg.ServerURL)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("X-Agent-Key", cfg.AgentKey)

	return client
}
