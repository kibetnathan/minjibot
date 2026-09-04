// Package config parses environment variables into a Config struct using
// the caarlos0/env library. All configuration comes from .env.
package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

// minSessionSecretLen is the shortest SESSION_SECRET we accept. A short secret
// weakens the HMAC that protects session cookies.
const minSessionSecretLen = 16

type Config struct {
	DBURL               string `env:"DB_URL"`
	Port                string `env:"PORT"`
	DiscordToken        string `env:"DISCORD_TOKEN"`
	DiscordClientID     string `env:"DISCORD_CLIENT_ID"`
	DiscordClientSecret string `env:"DISCORD_CLIENT_SECRET"`
	AppURL              string `env:"APP_URL"`
	FrontendURL         string `env:"FRONTEND_URL"`
	SessionSecret       string `env:"SESSION_SECRET"`
	GoogleFactCheckKey  string `env:"GOOGLE_FACTCHECK_API_KEY"`
	GeminiAPIKey        string `env:"GEMINI_API_KEY"`
	GeminiModel         string `env:"GEMINI_MODEL"`
	RedditClientID      string `env:"REDDIT_CLIENT_ID"`
	RedditClientSecret  string `env:"REDDIT_CLIENT_SECRET"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	return cfg, err
}

// ValidateForAPI checks that configuration required by the HTTP API/dashboard
// is present, so the server fails fast at startup instead of silently running
// in a broken or insecure state.
//
// It is intentionally not called for the bot-only entrypoint, which does not
// serve sessions and does not need SESSION_SECRET.
func (c *Config) ValidateForAPI() error {
	secret := strings.TrimSpace(c.SessionSecret)
	if secret == "" {
		// A random per-process fallback would make sessions unverifiable across
		// replicas and wipe them on every restart, so refuse to start instead.
		return fmt.Errorf("SESSION_SECRET is required for the API but is not set")
	}
	if len(secret) < minSessionSecretLen {
		return fmt.Errorf("SESSION_SECRET must be at least %d characters", minSessionSecretLen)
	}
	return nil
}
