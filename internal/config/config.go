package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBURL               string `env:"DB_URL"`
	DiscordToken        string `env:"DISCORD_TOKEN"`
	DiscordClientID     string `env:"DISCORD_CLIENT_ID"`
	DiscordClientSecret string `env:"DISCORD_CLIENT_SECRET"`
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
