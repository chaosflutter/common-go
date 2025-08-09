package tts

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds TTS configuration
type Config struct {
	AppID     string
	AccessKey string
	Cluster   string
	AudioDir  string
	BaseURL   string
}

// LoadConfig loads TTS configuration from environment variables
func LoadConfig() *Config {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		AppID:     getEnv("VOLCENGINE_TTS_APP_ID", ""),
		AccessKey: getEnv("VOLCENGINE_TTS_ACCESS_KEY", ""),
		Cluster:   getEnv("VOLCENGINE_TTS_CLUSTER", "volcano_tts"),
		AudioDir:  getEnv("TTS_AUDIO_DIR", "./audio"),
		BaseURL:   getEnv("TTS_BASE_URL", "http://localhost:8080"),
	}

	return config
}

// IsValid checks if the configuration is valid
func (c *Config) IsValid() bool {
	return c.AppID != "" && c.AccessKey != ""
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// DefaultConfig returns a default configuration for testing
func DefaultConfig() *Config {
	return &Config{
		AppID:     "your_app_id",
		AccessKey: "your_access_key",
		Cluster:   "volcano_tts",
		AudioDir:  "./audio",
		BaseURL:   "http://localhost:8080",
	}
}
