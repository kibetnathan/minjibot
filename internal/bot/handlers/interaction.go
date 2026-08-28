package handlers

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/domain/commands"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	"log/slog"
)

type InteractionHandlerDeps struct {
	Logger       *slog.Logger
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo     repository.UserPermissionRepository
	AuditRepo    repository.AuditLogRepository
}

func RegisterInteractionHandler(s *discordgo.Session, deps InteractionHandlerDeps, cmdHandler *commands.CommandHandler) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		onInteractionCreate(s, i, deps, cmdHandler)
	})
}

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate, deps InteractionHandlerDeps, cmdHandler *commands.CommandHandler) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	ctx := context.Background()

	// Respond to interaction
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		deps.Logger.Error("Failed to respond to interaction", "error", err)
		return
	}

	if err := cmdHandler.HandleSlash(ctx, s, i); err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Error: %v", err),
		})
	}
}