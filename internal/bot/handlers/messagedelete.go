package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
)

type MessageDeleteHandlerDeps struct {
	Logger             *slog.Logger
	GuildRepo          repository.GuildRepository
	AuditRepo          repository.AuditLogRepository
	SettingsRepo       repository.GuildSettingsRepository
	DeletedMessageRepo repository.DeletedMessageRepository
}

const logColor = 0xED4245

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
	authorName := ""
	var attachments []byte
	if m.BeforeDelete != nil {
		content = m.BeforeDelete.Content
		if m.BeforeDelete.Author != nil {
			authorID = m.BeforeDelete.Author.ID
			authorName = m.BeforeDelete.Author.Username
		}
		if len(m.BeforeDelete.Attachments) > 0 {
			if b, err := json.Marshal(m.BeforeDelete.Attachments); err == nil {
				attachments = b
			}
		}
	}
	if authorID == "" {
		authorID = "unknown"
	}

	// Persist the deleted message for the dashboard.
	if deps.DeletedMessageRepo != nil {
		if _, err := deps.DeletedMessageRepo.Create(ctx, dto.CreateDeletedMessageParams{
			GuildID:       m.GuildID,
			ChannelID:     m.ChannelID,
			MessageID:     m.ID,
			AuthorID:      authorID,
			AuthorName:    authorName,
			Content:       content,
			Attachments:   attachments,
			DeletedBy:     "unknown",
			DeletedByName: "unknown",
		}); err != nil {
			deps.Logger.Error("Failed to store deleted message", "error", err, "guild_id", m.GuildID, "message_id", m.ID)
		}
	}

	meta, _ := json.Marshal(map[string]any{
		"message_id": m.ID,
		"content":    content,
		"author_id":  authorID,
		"channel_id": m.ChannelID,
		"deleted_by": "unknown",
	})

	if _, err := deps.AuditRepo.Create(ctx, dto.CreateAuditLogParams{
		GuildID:  m.GuildID,
		Action:   "MESSAGE_DELETE",
		ActorID:  authorID,
		TargetID: m.ChannelID,
		Metadata: meta,
	}); err != nil {
		deps.Logger.Error("Failed to log message delete", "error", err, "guild_id", m.GuildID, "message_id", m.ID)
	}

	// Post a concise entry to the configured logging channel, if any.
	logChannel := deps.logChannel(ctx, m.GuildID)
	if logChannel == "" {
		return
	}
	desc := fmt.Sprintf("**Author:** <@%s> (%s)\n**Channel:** <#%s>",
		authorID, trimLog(authorName), m.ChannelID)
	if content != "" {
		desc += fmt.Sprintf("\n\n%s", trimLog(content))
	}
	_, err := s.ChannelMessageSendEmbed(logChannel, &discordgo.MessageEmbed{
		Color:       logColor,
		Title:       "Message deleted",
		Description: desc,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Message ID: %s", m.ID),
		},
	})
	if err != nil {
		deps.Logger.Error("Failed to post message delete to log channel", "error", err, "guild_id", m.GuildID)
	}
}

func onMessageDeleteBulk(s *discordgo.Session, m *discordgo.MessageDeleteBulk, deps MessageDeleteHandlerDeps) {
	ctx := context.Background()

	if err := ensureGuild(ctx, deps, m.GuildID); err != nil {
		return
	}

	// Persist each deleted message (content may be unknowable for bulk deletes).
	if deps.DeletedMessageRepo != nil {
		for _, id := range m.Messages {
			if _, err := deps.DeletedMessageRepo.Create(ctx, dto.CreateDeletedMessageParams{
				GuildID:       m.GuildID,
				ChannelID:     m.ChannelID,
				MessageID:     id,
				AuthorID:      "unknown",
				AuthorName:    "unknown",
				Content:       "",
				Attachments:   nil,
				DeletedBy:     "unknown",
				DeletedByName: "unknown",
			}); err != nil {
				deps.Logger.Error("Failed to store bulk-deleted message", "error", err, "guild_id", m.GuildID, "message_id", id)
			}
		}
	}

	meta, _ := json.Marshal(map[string]any{
		"message_ids": m.Messages,
		"channel_id":  m.ChannelID,
		"count":       len(m.Messages),
	})

	if _, err := deps.AuditRepo.Create(ctx, dto.CreateAuditLogParams{
		GuildID:  m.GuildID,
		Action:   "MESSAGE_DELETE_BULK",
		ActorID:  "unknown",
		TargetID: m.ChannelID,
		Metadata: meta,
	}); err != nil {
		deps.Logger.Error("Failed to log bulk message delete", "error", err, "guild_id", m.GuildID)
	}

	logChannel := deps.logChannel(ctx, m.GuildID)
	if logChannel == "" {
		return
	}
	_, err := s.ChannelMessageSendEmbed(logChannel, &discordgo.MessageEmbed{
		Color:       logColor,
		Title:       fmt.Sprintf("Bulk message deletion (%d messages)", len(m.Messages)),
		Description: fmt.Sprintf("**Channel:** <#%s>", m.ChannelID),
	})
	if err != nil {
		deps.Logger.Error("Failed to post bulk delete to log channel", "error", err, "guild_id", m.GuildID)
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

// logChannel resolves the guild's configured logging channel from settings, if
// any. It returns "" when logging isn't set up or settings can't be read.
func (deps MessageDeleteHandlerDeps) logChannel(ctx context.Context, guildID string) string {
	if deps.SettingsRepo == nil {
		return ""
	}
	settings, err := deps.SettingsRepo.Get(ctx, guildID)
	if err != nil {
		deps.Logger.Error("Failed to read settings for log channel", "error", err, "guild_id", guildID)
		return ""
	}
	return settings.LoggingChannelID
}

// trimLog limits a log string to a sensible size so an embed doesn't blow up.
func trimLog(s string) string {
	const max = 1024
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return strings.TrimSpace(s[:max]) + "…"
}
