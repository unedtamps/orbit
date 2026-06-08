package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIURL       string
	APIKey       string
	TMDBAPIKey   string
	HostURL      string
	Port         string
	CORSMaxAge   int
	Timeout      time.Duration
	StaticDir    string
	TemplateGlob string
}

func Load() (*Config, error) {
	cfg := &Config{
		APIURL:       getEnv("API_URL", ""),
		APIKey:       getEnv("API_KEY", ""),
		TMDBAPIKey:   getEnv("TMDB_API_KEY", ""),
		HostURL:      getEnv("HOST_URL", "http://localhost:9999"),
		Port:         getEnv("PORT", "9999"),
		CORSMaxAge:   getEnvInt("CORS_MAX_AGE", 300),
		Timeout:      getEnvDuration("PROXY_TIMEOUT", 30*time.Second),
		StaticDir:    getEnv("STATIC_DIR", "./static"),
		TemplateGlob: getEnv("TEMPLATE_GLOB", "templates/*.html"),
	}

	if cfg.APIURL == "" {
		return nil, fmt.Errorf("API_URL environment variable is required")
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable is required")
	}
	if cfg.TMDBAPIKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEY environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
