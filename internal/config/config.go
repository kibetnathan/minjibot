// Package config parses environment variables into a Config struct using
// the caarlos0/env library. All configuration comes from .env.
package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBURL               string `env:"DB_URL"`
	Port                string `env:"PORT"`
	DiscordToken        string `env:"DISCORD_TOKEN"`
	DiscordClientID     string `env:"DISCORD_CLIENT_ID"`
	DiscordClientSecret string `env:"DISCORD_CLIENT_SECRET"`
	AppURL              string `env:"APP_URL"`
	FrontendURL         string `env:"FRONTEND_URL"`
	SessionSecret       string `env:"SESSION_SECRET"`
	// DashboardAdminIDs is the allowlist of Discord user IDs permitted to view
	// the dashboard's guild data (deleted messages, moderation logs). Empty
	// means nobody is authorized — the endpoints fail closed.
	DashboardAdminIDs  []string `env:"DASHBOARD_ADMIN_IDS" envSeparator:","`
	GoogleFactCheckKey string   `env:"GOOGLE_FACTCHECK_API_KEY"`
	GeminiAPIKey       string   `env:"GEMINI_API_KEY"`
	GeminiModel        string   `env:"GEMINI_MODEL"`
	RedditClientID     string   `env:"REDDIT_CLIENT_ID"`
	RedditClientSecret string   `env:"REDDIT_CLIENT_SECRET"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	return cfg, err
}
