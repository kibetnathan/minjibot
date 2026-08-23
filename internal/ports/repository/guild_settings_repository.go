package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/entities"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// GuildSettings Repository --
type GuildSettingsRepository interface {
	Get(ctx context.Context, guildID string) (entities.GuildSettings, error)
	Upsert(ctx context.Context, arg dto.UpsertGuildSettingsParams) (entities.GuildSettings, error)
	Update(ctx context.Context, guildID string, arg dto.UpdateGuildSettingsParams) (entities.GuildSettings, error)
	Delete(ctx context.Context, guildID string) error
}

type sqlGuildSettingsRepository struct {
	store *SQLStore
}

func NewGuildSettingsRepository(store *SQLStore) GuildSettingsRepository {
	return &sqlGuildSettingsRepository{store: store}
}

func (r *sqlGuildSettingsRepository) Get(ctx context.Context, guildID string) (entities.GuildSettings, error) {
	settings, err := r.store.queries.GetGuildSettings(ctx, guildID)
	if err != nil {
		return entities.GuildSettings{}, err
	}

	return toentitiesGuildSettings(settings), nil
}

func (r *sqlGuildSettingsRepository) Upsert(ctx context.Context, arg dto.UpsertGuildSettingsParams) (entities.GuildSettings, error) {
	settings, err := r.store.queries.UpsertGuildSettings(ctx, postgres.UpsertGuildSettingsParams{
		GuildID:               arg.GuildID,
		Prefix:                arg.Prefix,
		Language:              arg.Language,
		AutoModerationEnabled: arg.AutoModerationEnabled,
		LoggingChannelID: pgtype.Text{
			String: arg.LoggingChannelID,
			Valid:  arg.LoggingChannelID != "",
		},
	})
	if err != nil {
		return entities.GuildSettings{}, err
	}

	return toentitiesGuildSettings(settings), nil
}

func (r *sqlGuildSettingsRepository) Update(ctx context.Context, guildID string, arg dto.UpdateGuildSettingsParams) (entities.GuildSettings, error) {
	settings, err := r.store.queries.UpdateGuildSettings(ctx, postgres.UpdateGuildSettingsParams{
		GuildID:               guildID,
		Prefix:                arg.Prefix,
		Language:              arg.Language,
		AutoModerationEnabled: arg.AutoModerationEnabled,
		LoggingChannelID: pgtype.Text{
			String: arg.LoggingChannelID,
			Valid:  arg.LoggingChannelID != "",
		},
	})
	if err != nil {
		return entities.GuildSettings{}, err
	}

	return toentitiesGuildSettings(settings), nil
}

func (r *sqlGuildSettingsRepository) Delete(ctx context.Context, guildID string) error {
	return r.store.queries.DeleteGuildSettings(ctx, guildID)
}

func toentitiesGuildSettings(s postgres.GuildSetting) entities.GuildSettings {
	return entities.GuildSettings{
		GuildID:               s.GuildID,
		Prefix:                s.Prefix,
		Language:              s.Language,
		AutoModerationEnabled: s.AutoModerationEnabled,
		LoggingChannelID:      s.LoggingChannelID.String,
		UpdatedAt:             s.UpdatedAt.Time,
	}
}
