package repository

import (
	"context"

	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/entities"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// UserPermission Repository --
type UserPermissionRepository interface {
	Get(ctx context.Context, userID, guildID, role string) (entities.UserPermission, error)
	ListForUser(ctx context.Context, userID, guildID string) ([]entities.UserPermission, error)
	ListForGuild(ctx context.Context, guildID string) ([]entities.UserPermission, error)
	Upsert(ctx context.Context, arg dto.UpsertUserPermissionParams) (entities.UserPermission, error)
	Delete(ctx context.Context, userID, guildID, role string) error
}

type sqlUserPermissionRepository struct {
	store *SQLStore
}

func NewUserPermissionRepository(store *SQLStore) UserPermissionRepository {
	return &sqlUserPermissionRepository{store: store}
}

func (r *sqlUserPermissionRepository) Get(ctx context.Context, userID, guildID, role string) (entities.UserPermission, error) {
	permission, err := r.store.queries.GetUserPermission(ctx, postgres.GetUserPermissionParams{
		UserID:  userID,
		GuildID: guildID,
		Role:    role,
	})
	if err != nil {
		return entities.UserPermission{}, err
	}

	return toentitiesUserPermission(permission), nil
}

func (r *sqlUserPermissionRepository) ListForUser(ctx context.Context, userID, guildID string) ([]entities.UserPermission, error) {
	permissions, err := r.store.queries.ListUserPermissionsForUser(ctx, postgres.ListUserPermissionsForUserParams{
		UserID:  userID,
		GuildID: guildID,
	})
	if err != nil {
		return nil, err
	}

	return toentitiesUserPermissions(permissions), nil
}

func (r *sqlUserPermissionRepository) ListForGuild(ctx context.Context, guildID string) ([]entities.UserPermission, error) {
	permissions, err := r.store.queries.ListUserPermissionsForGuild(ctx, guildID)
	if err != nil {
		return nil, err
	}

	return toentitiesUserPermissions(permissions), nil
}

func (r *sqlUserPermissionRepository) Upsert(ctx context.Context, arg dto.UpsertUserPermissionParams) (entities.UserPermission, error) {
	permission, err := r.store.queries.UpsertUserPermission(ctx, postgres.UpsertUserPermissionParams{
		UserID:          arg.UserID,
		GuildID:         arg.GuildID,
		Role:            arg.Role,
		PermissionsJson: arg.PermissionsJSON,
	})
	if err != nil {
		return entities.UserPermission{}, err
	}

	return toentitiesUserPermission(permission), nil
}

func (r *sqlUserPermissionRepository) Delete(ctx context.Context, userID, guildID, role string) error {
	return r.store.queries.DeleteUserPermission(ctx, postgres.DeleteUserPermissionParams{
		UserID:  userID,
		GuildID: guildID,
		Role:    role,
	})
}

func toentitiesUserPermissions(perms []postgres.UserPermission) []entities.UserPermission {
	entitiesPerms := make([]entities.UserPermission, len(perms))
	for i, p := range perms {
		entitiesPerms[i] = toentitiesUserPermission(p)
	}

	return entitiesPerms
}

func toentitiesUserPermission(p postgres.UserPermission) entities.UserPermission {
	return entities.UserPermission{
		ID:              p.ID,
		UserID:          p.UserID,
		GuildID:         p.GuildID,
		Role:            p.Role,
		PermissionsJSON: p.PermissionsJson,
		CreatedAt:       p.CreatedAt.Time,
	}
}
