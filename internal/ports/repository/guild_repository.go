package repository

import (
	"context"

	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/entities"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// Guild Repository --
type GuildRepository interface {
	GetByID(ctx context.Context, id string) (entities.Guild, error)
	List(ctx context.Context) ([]entities.Guild, error)
	Create(ctx context.Context, arg dto.CreateGuildParams) (entities.Guild, error)
	Update(ctx context.Context, id string, arg dto.UpdateGuildParams) (entities.Guild, error)
	Delete(ctx context.Context, id string) error
}

type sqlGuildRepository struct {
	store *SQLStore
}

func NewGuildRepository(store *SQLStore) GuildRepository {
	return &sqlGuildRepository{store: store}
}

func (r *sqlGuildRepository) GetByID(ctx context.Context, id string) (entities.Guild, error) {
	guild, err := r.store.queries.GetGuild(ctx, id)
	if err != nil {
		return entities.Guild{}, err
	}

	return toentitiesGuild(guild), nil
}

func (r *sqlGuildRepository) List(ctx context.Context) ([]entities.Guild, error) {
	guilds, err := r.store.queries.ListGuilds(ctx)
	if err != nil {
		return nil, err
	}

	entitiesGuilds := make([]entities.Guild, len(guilds))
	for i, g := range guilds {
		entitiesGuilds[i] = toentitiesGuild(g)
	}

	return entitiesGuilds, nil
}

func (r *sqlGuildRepository) Create(ctx context.Context, arg dto.CreateGuildParams) (entities.Guild, error) {
	guild, err := r.store.queries.CreateGuild(ctx, postgres.CreateGuildParams{
		ID:          arg.ID,
		Name:        arg.Name,
		PremiumTier: arg.PremiumTier,
	})
	if err != nil {
		return entities.Guild{}, err
	}

	return toentitiesGuild(guild), nil
}

func (r *sqlGuildRepository) Update(ctx context.Context, id string, arg dto.UpdateGuildParams) (entities.Guild, error) {
	guild, err := r.store.queries.UpdateGuild(ctx, postgres.UpdateGuildParams{
		ID:          id,
		Name:        arg.Name,
		PremiumTier: arg.PremiumTier,
	})
	if err != nil {
		return entities.Guild{}, err
	}

	return toentitiesGuild(guild), nil
}

func (r *sqlGuildRepository) Delete(ctx context.Context, id string) error {
	return r.store.queries.DeleteGuild(ctx, id)
}

func toentitiesGuild(g postgres.Guild) entities.Guild {
	return entities.Guild{
		ID:          g.ID,
		Name:        g.Name,
		PremiumTier: g.PremiumTier,
		CreatedAt:   g.CreatedAt.Time,
	}
}
