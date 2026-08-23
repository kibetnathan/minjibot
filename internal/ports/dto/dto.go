package dto

type CreateGuildParams struct {
	ID          string
	Name        string
	PremiumTier int32
}

type UpdateGuildParams struct {
	Name        string
	PremiumTier int32
}

type UpsertGuildSettingsParams struct {
	GuildID               string
	Prefix                string
	Language              string
	AutoModerationEnabled bool
	LoggingChannelID      string
}

type UpdateGuildSettingsParams struct {
	Prefix                string
	Language              string
	AutoModerationEnabled bool
	LoggingChannelID      string
}

type UpsertUserPermissionParams struct {
	UserID          string
	GuildID         string
	Role            string
	PermissionsJSON []byte
}

type CreateAuditLogParams struct {
	GuildID  string
	Action   string
	ActorID  string
	TargetID string
	Metadata []byte
}
