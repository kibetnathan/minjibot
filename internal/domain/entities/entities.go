package entities

import "time"

type Guild struct {
	ID          string
	Name        string
	PremiumTier int32
	CreatedAt   time.Time
}

type GuildSettings struct {
	GuildID               string
	Prefix                string
	Language              string
	AutoModerationEnabled bool
	LoggingChannelID      string
	UpdatedAt             time.Time
}

type UserPermission struct {
	ID              int64
	UserID          string
	GuildID         string
	Role            string
	PermissionsJSON []byte
	CreatedAt       time.Time
}

type AuditLog struct {
	ID        int64
	GuildID   string
	Action    string
	ActorID   string
	TargetID  string
	Metadata  []byte
	CreatedAt time.Time
}
