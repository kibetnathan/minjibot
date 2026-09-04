package dto

import "time"

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
	MessageLoggingEnabled bool
}

type UpdateGuildSettingsParams struct {
	Prefix                string
	Language              string
	AutoModerationEnabled bool
	LoggingChannelID      string
	MessageLoggingEnabled bool
}

type UpsertUserPermissionParams struct {
	UserID          string
	GuildID         string
	Role            string
	PermissionsJSON []byte
}

type CreateAuditLogParams struct {
	GuildID    string
	Action     string
	ActorID    string
	ActorName  string
	TargetID   string
	TargetName string
	Metadata   []byte
}

type CreateUserParams struct {
	UserID       string
	Email        *string
	Passwordhash *string
}

type UpdateUserParams struct {
	Email        *string
	Passwordhash *string
	IsActive     bool
}

type SetBirthdayParams struct {
	GuildID  string
	UserID   string
	Birthday time.Time
}

type SetGuildBirthdayChannelParams struct {
	GuildID   string
	ChannelID string
}

type SetGuildBirthdayRoleParams struct {
	GuildID string
	RoleID  string
}

type CreateDiaryEntryParams struct {
	UserID  string
	Content string
}

type CreateDeletedMessageParams struct {
	GuildID       string
	ChannelID     string
	MessageID     string
	AuthorID      string
	AuthorName    string
	Content       string
	Attachments   []byte
	DeletedBy     string
	DeletedByName string
}

type ListDeletedMessagesParams struct {
	GuildID   string
	ChannelID string
	Limit     int32
	Offset    int32
}
