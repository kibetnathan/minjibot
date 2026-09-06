package handlers

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/commands"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
	"github.com/kibetnathan/minjibot/internal/safe"
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
	// A panic in any command handler must not take down the whole bot.
	defer safe.Recover(deps.Logger, "onInteractionCreate")

	// Component interactions (e.g. help pagination buttons) are routed to the
	// relevant handler directly.
	if i.Type == discordgo.InteractionMessageComponent {
		if err := commands.HelpButtonHandler(s, i); err != nil {
			deps.Logger.Error("Failed to handle component interaction", "custom_id", i.MessageComponentData().CustomID, "error", err)
		}
		return
	}

	if i.Type == discordgo.InteractionModalSubmit {
		if err := commands.HandleBugModalSubmit(s, i); err != nil {
			deps.Logger.Error("Failed to handle modal submit", "custom_id", i.ModalSubmitData().CustomID, "error", err)
		}
		return
	}

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	ctx := context.Background()

	// Command handlers are responsible for their own interaction response.
	// If a handler returns an error before responding, try a followup; if the
	// interaction was never acknowledged, fall back to a direct response.
	if err := cmdHandler.HandleSlash(ctx, s, i); err != nil {
		deps.Logger.Error("Failed to handle slash command", "command", i.ApplicationCommandData().Name, "error", err)

		_, ferr := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Error: %v", err),
		})
		if ferr != nil {
			deps.Logger.Error("Failed to send error followup", "error", ferr)
			if rerr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Error: %v", err),
				},
			}); rerr != nil {
				deps.Logger.Error("Failed to send error response", "error", rerr)
			}
		}
	}
}
