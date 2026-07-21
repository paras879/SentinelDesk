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
	configPath := findConfigFile("config.yaml")
	log.Printf("Config file path: %s", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", configPath, err)
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("Failed to parse %s: %v", configPath, err)
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
		cfg.HeartbeatInterval = 15
	}

	log.Printf("Loaded Server URL: %s", cfg.ServerURL)
}

func findConfigFile(name string) string {
	if _, err := os.Stat(name); err == nil {
		return name
	}
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		p := filepath.Join(exeDir, name)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return name
}

func Get() *Config {
	return cfg
}

func GetHeartbeatInterval() int {
	if cfg == nil {
		return 15
	}
	return cfg.HeartbeatInterval
}
