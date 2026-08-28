package handlers

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	"log/slog"
)

type HandlerDeps struct {
	Logger       *slog.Logger
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	AuditRepo    repository.AuditLogRepository
}

func RegisterMessageHandler(s *discordgo.Session, deps HandlerDeps) {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		onMessageCreate(s, m, deps)
	})
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate, deps HandlerDeps) {
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
		Metadata: []byte(fmt.Sprintf(`{"content":%q,"channel_id":%q}`, m.Content, m.ChannelID)),
	})
	if err != nil {
		deps.Logger.Error("Failed to create audit log", "error", err)
	}

	// Get guild settings for prefix
	settings, err := deps.SettingsRepo.Get(ctx, m.GuildID)
	if err != nil {
		deps.Logger.Debug("No guild settings found", "guild_id", m.GuildID)
	}

	_ = guild
	_ = settings

	// TODO: Add command handling logic here using settings.Prefix
}