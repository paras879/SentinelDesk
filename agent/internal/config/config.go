package config

import (
	"log"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerURL         string `yaml:"server_url"`
	AgentKey          string `yaml:"agent_key"`
	HeartbeatInterval int    `yaml:"heartbeat_interval"`
}

var cfg *Config

func Load() {
	configPath := "config.yaml"
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		absPath = configPath
	}

	log.Printf("Config file path: %s", absPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", absPath, err)
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("Failed to parse %s: %v", absPath, err)
	}

	if cfg.ServerURL == "" {
		log.Fatal("server_url is required in config.yaml")
	}
	if _, err := url.ParseRequestURI(cfg.ServerURL); err != nil {
		log.Fatalf("server_url in config.yaml is invalid (%q): %v", cfg.ServerURL, err)
	}
	if cfg.AgentKey == "" {
		log.Fatal("agent_key is required in config.yaml")
	}
	if cfg.HeartbeatInterval <= 0 {
		cfg.HeartbeatInterval = 30
	}

	log.Printf("Loaded Server URL: %s", cfg.ServerURL)
}

func Get() *Config {
	return cfg
}

func GetHeartbeatInterval() int {
	if cfg == nil {
		return 30
	}
	return cfg.HeartbeatInterval
}
