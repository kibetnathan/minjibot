package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/auditlog"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// AuditLog Repository --
type AuditLogRepository interface {
	Create(ctx context.Context, arg dto.CreateAuditLogParams) (auditlog.AuditLog, error)
	ListForGuild(ctx context.Context, guildID string, limit, offset int32) ([]auditlog.AuditLog, error)
	ListByActor(ctx context.Context, guildID, actorID string, limit, offset int32) ([]auditlog.AuditLog, error)
	CountForGuild(ctx context.Context, guildID string) (int64, error)
	CountForAllGuilds(ctx context.Context) (map[string]int64, error)
	DeleteBefore(ctx context.Context, cutoff time.Time) error
}

type sqlAuditLogRepository struct {
	store *SQLStore
}

func NewAuditLogRepository(store *SQLStore) AuditLogRepository {
	return &sqlAuditLogRepository{store: store}
}

func (r *sqlAuditLogRepository) Create(ctx context.Context, arg dto.CreateAuditLogParams) (auditlog.AuditLog, error) {
	log, err := r.store.queries.CreateAuditLog(ctx, postgres.CreateAuditLogParams{
		GuildID:   arg.GuildID,
		Action:    arg.Action,
		ActorID:   arg.ActorID,
		ActorName: arg.ActorName,
		TargetID: pgtype.Text{
			String: arg.TargetID,
			Valid:  arg.TargetID != "",
		},
		TargetName: arg.TargetName,
		Metadata:   arg.Metadata,
	})
	if err != nil {
		return auditlog.AuditLog{}, err
	}

	return toentitiesAuditLog(log), nil
}

func (r *sqlAuditLogRepository) ListForGuild(ctx context.Context, guildID string, limit, offset int32) ([]auditlog.AuditLog, error) {
	logs, err := r.store.queries.ListAuditLogsForGuild(ctx, postgres.ListAuditLogsForGuildParams{
		GuildID:    guildID,
		PageOffset: offset,
		PageSize:   limit,
	})
	if err != nil {
		return nil, err
	}

	return toentitiesAuditLogs(logs), nil
}

func (r *sqlAuditLogRepository) ListByActor(ctx context.Context, guildID, actorID string, limit, offset int32) ([]auditlog.AuditLog, error) {
	logs, err := r.store.queries.ListAuditLogsByActor(ctx, postgres.ListAuditLogsByActorParams{
		GuildID:    guildID,
		ActorID:    actorID,
		PageOffset: offset,
		PageSize:   limit,
	})
	if err != nil {
		return nil, err
	}

	return toentitiesAuditLogs(logs), nil
}

func (r *sqlAuditLogRepository) CountForGuild(ctx context.Context, guildID string) (int64, error) {
	return r.store.queries.CountAuditLogsForGuild(ctx, guildID)
}

func (r *sqlAuditLogRepository) CountForAllGuilds(ctx context.Context) (map[string]int64, error) {
	rows, err := r.store.queries.CountAuditLogsForAllGuilds(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, row := range rows {
		out[row.GuildID] = row.Count
	}
	return out, nil
}

func (r *sqlAuditLogRepository) DeleteBefore(ctx context.Context, cutoff time.Time) error {
	return r.store.queries.DeleteAuditLogsBefore(ctx, pgtype.Timestamptz{
		Time:  cutoff,
		Valid: true,
	})
}

func toentitiesAuditLogs(logs []postgres.AuditLog) []auditlog.AuditLog {
	entitiesLogs := make([]auditlog.AuditLog, len(logs))
	for i, l := range logs {
		entitiesLogs[i] = toentitiesAuditLog(l)
	}

	return entitiesLogs
}

func toentitiesAuditLog(l postgres.AuditLog) auditlog.AuditLog {
	return auditlog.AuditLog{
		ID:         l.ID,
		GuildID:    l.GuildID,
		Action:     l.Action,
		ActorID:    l.ActorID,
		ActorName:  l.ActorName,
		TargetID:   l.TargetID.String,
		TargetName: l.TargetName,
		Metadata:   l.Metadata,
		CreatedAt:  l.CreatedAt.Time,
	}
}
