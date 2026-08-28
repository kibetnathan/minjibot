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
		return h.ping(s, m.ChannelID)
	case "help":
		return h.help(s, m.ChannelID)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func (h *CommandHandler) HandleSlash(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) error {
	switch i.ApplicationCommandData().Name {
	case "ping":
		return h.ping(s, i.ChannelID)
	case "help":
		return h.help(s, i.ChannelID)
	default:
		return fmt.Errorf("unknown command: %s", i.ApplicationCommandData().Name)
	}
}

func (h *CommandHandler) ping(s *discordgo.Session, channelID string) error {
	_, err := s.ChannelMessageSend(channelID, "Pong! 🏓")
	return err
}

func (h *CommandHandler) help(s *discordgo.Session, channelID string) error {
	_, err := s.ChannelMessageSend(channelID, "Available commands: `ping`, `help`")
	return err
}

var SlashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Check bot latency",
	},
	{
		Name:        "help",
		Description: "Show available commands",
	},
}