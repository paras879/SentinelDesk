package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName    string
	AppPort    string
	AppEnv     string
	AppVersion string

	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	DatabaseURL  string

	// JWT
	JWTSecret string

	// Agent API Key
	AgentAPIKey string

	// Plan & Licensing
	MaxDevices int
	PlanName   string

	// CORS
	AllowedOrigins string
}

func LoadConfig() *Config {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	maxDev, err := strconv.Atoi(getEnv("MAX_DEVICES", "3"))
	if err != nil {
		maxDev = 3
	}

	return &Config{
		// Application
		AppName:    getEnv("APP_NAME", "SentinelDesk"),
		AppPort:    getEnv("PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		AppVersion: getEnv("APP_VERSION", "1.0.0"),

		// Plan & Licensing
		MaxDevices: maxDev,
		PlanName:   getEnv("PLAN_NAME", "Starter"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "admin"),
		DBPassword:  getEnv("DB_PASSWORD", "admin123"),
		DBName:      getEnv("DB_NAME", "sentineldesk"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "SentinelDeskSuperSecretKey123456789"),

		// Agent API Key
		AgentAPIKey: getEnv("AGENT_API_KEY", "SentinelDeskAgentKey123456"),

		// CORS
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
