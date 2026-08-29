package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
)

type MessageDeleteHandlerDeps struct {
	Logger    *slog.Logger
	GuildRepo repository.GuildRepository
	AuditRepo repository.AuditLogRepository
}

func RegisterMessageDeleteHandler(s *discordgo.Session, deps MessageDeleteHandlerDeps) {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDelete) {
		onMessageDelete(s, m, deps)
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageDeleteBulk) {
		onMessageDeleteBulk(s, m, deps)
	})
}

func onMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete, deps MessageDeleteHandlerDeps) {
	ctx := context.Background()

	// The cache may hold the deleted message even though Discord only sends
	// an ID in the event; BeforeDelete is populated by discordgo's state.
	if err := ensureGuild(ctx, deps, m.GuildID); err != nil {
		return
	}

	content := ""
	authorID := ""
	if m.BeforeDelete != nil {
		content = m.BeforeDelete.Content
		if m.BeforeDelete.Author != nil {
			authorID = m.BeforeDelete.Author.ID
		}
	}
	if authorID == "" {
		authorID = "unknown"
	}

	meta, _ := json.Marshal(map[string]any{
		"message_id": m.ID,
		"content":    content,
		"author_id":  authorID,
		"channel_id": m.ChannelID,
		"deleted_by": "unknown",
	})

	_, err := deps.AuditRepo.Create(ctx, dto.CreateAuditLogParams{
		GuildID:  m.GuildID,
		Action:   "MESSAGE_DELETE",
		ActorID:  authorID,
		TargetID: m.ChannelID,
		Metadata: meta,
	})
	if err != nil {
		deps.Logger.Error("Failed to log message delete", "error", err, "guild_id", m.GuildID, "message_id", m.ID)
	}
}

func onMessageDeleteBulk(s *discordgo.Session, m *discordgo.MessageDeleteBulk, deps MessageDeleteHandlerDeps) {
	ctx := context.Background()

	if err := ensureGuild(ctx, deps, m.GuildID); err != nil {
		return
	}

	meta, _ := json.Marshal(map[string]any{
		"message_ids": m.Messages,
		"channel_id":  m.ChannelID,
		"count":       len(m.Messages),
	})

	_, err := deps.AuditRepo.Create(ctx, dto.CreateAuditLogParams{
		GuildID:  m.GuildID,
		Action:   "MESSAGE_DELETE_BULK",
		ActorID:  "unknown",
		TargetID: m.ChannelID,
		Metadata: meta,
	})
	if err != nil {
		deps.Logger.Error("Failed to log bulk message delete", "error", err, "guild_id", m.GuildID)
	}
}

func ensureGuild(ctx context.Context, deps MessageDeleteHandlerDeps, guildID string) error {
	_, err := deps.GuildRepo.GetByID(ctx, guildID)
	if err != nil {
		_, err = deps.GuildRepo.Create(ctx, dto.CreateGuildParams{
			ID:          guildID,
			Name:        "",
			PremiumTier: 0,
		})
		if err != nil {
			deps.Logger.Error("Failed to create guild during delete logging", "error", err, "guild_id", guildID)
			return fmt.Errorf("create guild: %w", err)
		}
	}
	return nil
}
