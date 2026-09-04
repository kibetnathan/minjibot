package repository

import (
	"context"

	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/deletedmessage"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// DeletedMessageRepository --
type DeletedMessageRepository interface {
	Create(ctx context.Context, arg dto.CreateDeletedMessageParams) (deletedmessage.DeletedMessage, error)
	ListForGuild(ctx context.Context, guildID string, limit, offset int32) ([]deletedmessage.DeletedMessage, error)
	ListForChannel(ctx context.Context, guildID, channelID string, limit, offset int32) ([]deletedmessage.DeletedMessage, error)
	CountForGuild(ctx context.Context, guildID string) (int64, error)
	CountForAllGuilds(ctx context.Context) (map[string]int64, error)
}

type sqlDeletedMessageRepository struct {
	store *SQLStore
}

func NewDeletedMessageRepository(store *SQLStore) DeletedMessageRepository {
	return &sqlDeletedMessageRepository{store: store}
}

func (r *sqlDeletedMessageRepository) Create(ctx context.Context, arg dto.CreateDeletedMessageParams) (deletedmessage.DeletedMessage, error) {
	dm, err := r.store.queries.InsertDeletedMessage(ctx, postgres.InsertDeletedMessageParams{
		GuildID:       arg.GuildID,
		ChannelID:     arg.ChannelID,
		MessageID:     arg.MessageID,
		AuthorID:      arg.AuthorID,
		AuthorName:    arg.AuthorName,
		Content:       arg.Content,
		Attachments:   arg.Attachments,
		DeletedBy:     arg.DeletedBy,
		DeletedByName: arg.DeletedByName,
	})
	if err != nil {
		return deletedmessage.DeletedMessage{}, err
	}
	return toentitiesDeletedMessage(dm), nil
}

func (r *sqlDeletedMessageRepository) ListForGuild(ctx context.Context, guildID string, limit, offset int32) ([]deletedmessage.DeletedMessage, error) {
	rows, err := r.store.queries.ListDeletedMessagesForGuild(ctx, postgres.ListDeletedMessagesForGuildParams{
		GuildID: guildID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	return toentitiesDeletedMessages(rows), nil
}

func (r *sqlDeletedMessageRepository) ListForChannel(ctx context.Context, guildID, channelID string, limit, offset int32) ([]deletedmessage.DeletedMessage, error) {
	rows, err := r.store.queries.ListDeletedMessagesForChannel(ctx, postgres.ListDeletedMessagesForChannelParams{
		GuildID:   guildID,
		ChannelID: channelID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}
	return toentitiesDeletedMessages(rows), nil
}

func (r *sqlDeletedMessageRepository) CountForGuild(ctx context.Context, guildID string) (int64, error) {
	return r.store.queries.CountDeletedMessagesForGuild(ctx, guildID)
}

func (r *sqlDeletedMessageRepository) CountForAllGuilds(ctx context.Context) (map[string]int64, error) {
	rows, err := r.store.queries.CountDeletedMessagesForAllGuilds(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, row := range rows {
		out[row.GuildID] = row.Count
	}
	return out, nil
}

func toentitiesDeletedMessages(rows []postgres.DeletedMessage) []deletedmessage.DeletedMessage {
	out := make([]deletedmessage.DeletedMessage, len(rows))
	for i, r := range rows {
		out[i] = toentitiesDeletedMessage(r)
	}
	return out
}

func toentitiesDeletedMessage(d postgres.DeletedMessage) deletedmessage.DeletedMessage {
	return deletedmessage.DeletedMessage{
		ID:            d.ID,
		GuildID:       d.GuildID,
		ChannelID:     d.ChannelID,
		MessageID:     d.MessageID,
		AuthorID:      d.AuthorID,
		AuthorName:    d.AuthorName,
		Content:       d.Content,
		Attachments:   d.Attachments,
		DeletedBy:     d.DeletedBy,
		DeletedByName: d.DeletedByName,
		CreatedAt:     d.CreatedAt.Time,
	}
}
