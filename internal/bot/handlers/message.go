package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/domain/commands"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	"log/slog"
)

const DefaultPrefix = "-"

type MessageHandlerDeps struct {
	Logger       *slog.Logger
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo     repository.UserPermissionRepository
	AuditRepo    repository.AuditLogRepository
}

func RegisterMessageHandler(s *discordgo.Session, deps MessageHandlerDeps, cmdHandler *commands.CommandHandler) {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		onMessageCreate(s, m, deps, cmdHandler)
	})
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate, deps MessageHandlerDeps, cmdHandler *commands.CommandHandler) {
	if m.Author.Bot {
		return
	}

	ctx := context.Background()

	// Ensure guild exists in DB
	guild, err := deps.GuildRepo.GetByID(ctx, m.GuildID)
	if err != nil {
		guild, err = deps.GuildRepo.Create(ctx, dto.CreateGuildParams{
			ID:          m.GuildID,
			Name:        "",
			PremiumTier: 0,
		})
		if err != nil {
			deps.Logger.Error("Failed to create guild", "error", err, "guild_id", m.GuildID)
			return
		}
	}

	// Log the message as audit log
	_, err = deps.AuditRepo.Create(ctx, dto.CreateAuditLogParams{
		GuildID:  m.GuildID,
		Action:   "MESSAGE_CREATE",
		ActorID:  m.Author.ID,
		TargetID: m.ChannelID,
		Metadata: []byte(fmt.Sprintf(`{"message_id":%q,"content":%q,"channel_id":%q}`, m.ID, m.Content, m.ChannelID)),
	})
	if err != nil {
		deps.Logger.Error("Failed to create audit log", "error", err)
	}

	// Get guild settings for prefix
	settings, err := deps.SettingsRepo.Get(ctx, m.GuildID)
	prefix := DefaultPrefix
	if err == nil && settings.Prefix != "" {
		prefix = settings.Prefix
	}

	_ = guild
	_ = settings

	// Check for command
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(args) == 0 {
		return
	}

	if err := cmdHandler.Handle(ctx, s, m, args[0], args[1:]); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
	}
}
