package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiToken      string
	RecordNames   []string
	Providers     []string
	CheckInterval time.Duration
	EnableIPv6    bool
	DryRun        bool
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  No .env file found, using system environment variables")
	} else {
		log.Println("✅ Configuration loaded from .env file")
	}

	recordsRaw := getEnv("CF_RECORD_NAMES", "")

	return &Config{
		ApiToken:      os.Getenv("CF_API_TOKEN"),
		RecordNames:   strings.Split(recordsRaw, ","),
		Providers:     strings.Split(getEnv("IP_PROVIDERS", "amazon,cloudflare,google"), ","),
		CheckInterval: getDurationEnv("CHECK_INTERVAL", 10*time.Minute),
		EnableIPv6:    os.Getenv("ENABLE_IPV6") == "true",
		DryRun:        os.Getenv("DRY_RUN") == "true",
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		d, err := time.ParseDuration(value)
		if err == nil {
			return d
		}
	}

	return fallback
}
