package birthday

import "time"

// Birthday is a member's saved birthday within a guild (birthday add/list/
// celebrate). One row per (guild, user) pair.
type Birthday struct {
	ID        int64
	GuildID   string
	UserID    string
	Birthday  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GuildSetting is the per-guild birthday configuration: the channel where
// automated notices are posted and the temporary role awarded on a birthday.
// Empty strings mean the value hasn't been configured yet.
type GuildSetting struct {
	GuildID   string
	ChannelID string
	RoleID    string
	UpdatedAt time.Time
}
