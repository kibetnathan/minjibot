package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/kibetnathan/minjibot/internal/ports/repository"
)

type CommandHandler struct {
	GuildRepo    repository.GuildRepository
	SettingsRepo repository.GuildSettingsRepository
	PermRepo     repository.UserPermissionRepository
}

func NewCommandHandler(guildRepo repository.GuildRepository, settingsRepo repository.GuildSettingsRepository, permRepo repository.UserPermissionRepository) *CommandHandler {
	return &CommandHandler{
		GuildRepo:    guildRepo,
		SettingsRepo: settingsRepo,
		PermRepo:     permRepo,
	}
}

func (h *CommandHandler) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, cmd string, args []string) error {
	switch cmd {
	case "ping":
		return h.ping(s, m)
	case "help":
		return h.help(s, m)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (h *CommandHandler) ping(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong! 🏓")
	return err
}

func (h *CommandHandler) help(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "Available commands: `ping`, `help`")
	return err
}