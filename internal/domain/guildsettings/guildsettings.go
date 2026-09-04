// Package guildsettings defines the domain model for per-guild bot
// configuration (logging channel, prefix, etc.).
package guildsettings

import "time"

type GuildSettings struct {
	GuildID               string
	Prefix                string
	Language              string
	AutoModerationEnabled bool
	LoggingChannelID      string
	UpdatedAt             time.Time
}
