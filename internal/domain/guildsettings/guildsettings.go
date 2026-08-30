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
