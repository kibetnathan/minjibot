package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/birthday"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// GuildBirthdaySettings Repository --
type GuildBirthdaySettingsRepository interface {
	Get(ctx context.Context, guildID string) (birthday.GuildSetting, error)
	SetChannel(ctx context.Context, arg dto.SetGuildBirthdayChannelParams) (birthday.GuildSetting, error)
	SetRole(ctx context.Context, arg dto.SetGuildBirthdayRoleParams) (birthday.GuildSetting, error)
}

type sqlGuildBirthdaySettingsRepository struct {
	store *SQLStore
}

func NewGuildBirthdaySettingsRepository(store *SQLStore) GuildBirthdaySettingsRepository {
	return &sqlGuildBirthdaySettingsRepository{store: store}
}

func (r *sqlGuildBirthdaySettingsRepository) Get(ctx context.Context, guildID string) (birthday.GuildSetting, error) {
	s, err := r.store.queries.GetGuildBirthdaySettings(ctx, guildID)
	if err != nil {
		return birthday.GuildSetting{}, err
	}
	return toentitiesGuildBirthdaySetting(s), nil
}

func (r *sqlGuildBirthdaySettingsRepository) SetChannel(ctx context.Context, arg dto.SetGuildBirthdayChannelParams) (birthday.GuildSetting, error) {
	s, err := r.store.queries.UpsertGuildBirthdayChannel(ctx, postgres.UpsertGuildBirthdayChannelParams{
		GuildID: arg.GuildID,
		ChannelID: pgtype.Text{
			String: arg.ChannelID,
			Valid:  arg.ChannelID != "",
		},
	})
	if err != nil {
		return birthday.GuildSetting{}, err
	}
	return toentitiesGuildBirthdaySetting(s), nil
}

func (r *sqlGuildBirthdaySettingsRepository) SetRole(ctx context.Context, arg dto.SetGuildBirthdayRoleParams) (birthday.GuildSetting, error) {
	s, err := r.store.queries.UpsertGuildBirthdayRole(ctx, postgres.UpsertGuildBirthdayRoleParams{
		GuildID: arg.GuildID,
		RoleID: pgtype.Text{
			String: arg.RoleID,
			Valid:  arg.RoleID != "",
		},
	})
	if err != nil {
		return birthday.GuildSetting{}, err
	}
	return toentitiesGuildBirthdaySetting(s), nil
}

func toentitiesGuildBirthdaySetting(s postgres.GuildBirthdaySetting) birthday.GuildSetting {
	return birthday.GuildSetting{
		GuildID:   s.GuildID,
		ChannelID: s.ChannelID.String,
		RoleID:    s.RoleID.String,
		UpdatedAt: s.UpdatedAt.Time,
	}
}
